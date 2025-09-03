package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/logger"
	authdto "github.com/topboyasante/pitstop/internal/modules/auth/dto"
	"github.com/topboyasante/pitstop/internal/modules/user/dto"
	"github.com/topboyasante/pitstop/internal/modules/user/service"
	"github.com/topboyasante/pitstop/internal/shared/events"
	"github.com/topboyasante/pitstop/internal/shared/utils"
)

type AuthService struct {
	config      *config.Config
	redis       *redis.Client
	validator   *validator.Validate
	eventBus    *events.EventBus
	userService *service.UserService
}

// NewAuthService creates a new instance of AuthService with the provided configuration
func NewAuthService(config *config.Config, redis *redis.Client, eventBus *events.EventBus, validator *validator.Validate, userService *service.UserService) *AuthService {
	logger.Info("Initializing auth service")
	return &AuthService{
		config:      config,
		redis:       redis,
		validator:   validator,
		eventBus:    eventBus,
		userService: userService,
	}
}

// Authenticate generates a CSRF state token and returns the OAuth authorization URL
func (as *AuthService) Authenticate() string {
	logger.Info("OAuth authentication initiated",
		"event", "auth.oauth_started",
		"provider", "google")

	state := as.generateState()

	// Store state in Redis with 10 minute expiration
	key := fmt.Sprintf("oauth:state:%s", state)
	err := as.redis.Set(context.Background(), key, "1", 10*time.Minute).Err()
	if err != nil {
		logger.Error("Failed to store OAuth state in Redis",
			"event", "auth.state_store_failed",
			"provider", "google",
			"operation", "redis_set",
			"ttl_minutes", 10,
			"error", err)
		// Continue anyway - could fall back to in-memory if needed
	} else {
		logger.Info("OAuth state stored successfully",
			"event", "auth.state_stored",
			"provider", "google",
			"operation", "redis_set",
			"ttl_minutes", 10)
	}

	url := as.config.OAuth.AuthCodeURL(state)

	logger.Info("OAuth URL generated",
		"event", "auth.oauth_url_generated",
		"provider", "google",
		"redirect_host", "accounts.google.com",
		"state_stored", err == nil)

	return url
}

// generateState creates a cryptographically secure random state token for CSRF protection
func (as *AuthService) generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// ValidateState verifies the CSRF state token and removes it to prevent reuse
func (as *AuthService) ValidateState(state string) bool {
	key := fmt.Sprintf("oauth:state:%s", state)

	// Check if state exists in Redis
	exists, err := as.redis.Get(context.Background(), key).Result()
	if err != nil || exists == "" {
		logger.Warn("State validation failed",
			"event", "auth.state_validation_failed",
			"provider", "google",
			"operation", "redis_get",
			"reason", "state_not_found_or_expired",
			"redis_error", err != nil,
			"error", err)
		return false
	}

	// Delete the state immediately (one-time use)
	err = as.redis.Del(context.Background(), key).Err()
	if err != nil {
		logger.Error("Failed to delete OAuth state from Redis",
			"event", "auth.state_delete_failed",
			"provider", "google",
			"operation", "redis_del",
			"error", err)
		// Continue anyway - state was valid
	}

	logger.Info("State validation successful",
		"event", "auth.state_validated",
		"provider", "google",
		"operation", "redis_del",
		"state_deleted", err == nil)
	return true
}

// ExchangeCode validates the state token and exchanges the authorization code for JWT tokens
func (as *AuthService) ExchangeCode(code, state string) (*authdto.JWTTokenResponse, error) {
	logger.Info("Token exchange initiated",
		"event", "auth.token_exchange_started")

	if !as.ValidateState(state) {
		logger.Error("Token exchange failed",
			"event", "auth.token_exchange_failed",
			"reason", "invalid_state")
		return nil, fmt.Errorf("invalid state")
	}

	token, err := as.config.OAuth.Exchange(context.TODO(), code)
	if err != nil {
		logger.Error("Token exchange failed",
			"event", "auth.token_exchange_failed",
			"reason", "oauth_exchange_error",
			"error", err.Error())
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	profile, err := as.GetGoogleProfile(token.AccessToken)
	if err != nil {
		return nil, err
	}

	// Create or get existing user synchronously
	userReq := dto.CreateUserRequest{
		Provider:   "google",
		ProviderID: profile.ID,
		FirstName:  profile.FirstName,
		LastName:   profile.LastName,
		Email:      profile.Email,
		AvatarURL:  profile.Picture,
		Locale:     profile.Locale,
	}

	user, err := as.userService.CreateUser(userReq)
	if err != nil {
		logger.Error("Failed to create/get user during OAuth",
			"event", "auth.user_creation_failed",
			"provider", "google",
			"google_id", profile.ID,
			"email", profile.Email,
			"error", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info("User created/retrieved for OAuth",
		"event", "auth.user_ready",
		"provider", "google",
		"internal_user_id", user.ID,
		"google_id", profile.ID)

	// Generate JWT tokens using internal user ID
	accessToken, refreshToken, expiresAt, err := utils.CreateJWTTokens(as.config, user.ID, "web")
	if err != nil {
		logger.Error("Failed to create JWT tokens",
			"event", "auth.jwt_creation_failed",
			"internal_user_id", user.ID,
			"error", err)
		return nil, fmt.Errorf("failed to create JWT tokens: %w", err)
	}

	// Publish event after successful user creation and JWT generation
	event := events.NewAuthenticationSuccessful(
		"google",
		profile.ID,
		profile.Email,
		profile.FirstName,
		profile.LastName,
		profile.Picture,
		profile.Locale,
	)
	as.eventBus.Publish("AuthenticationSuccessful", event)

	logger.Info("Token exchange successful",
		"event", "auth.token_exchange_completed",
		"internal_user_id", user.ID,
		"google_id", profile.ID,
		"expiresAt", expiresAt)

	return &authdto.JWTTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
	}, nil
}

// GetGoogleProfile fetches user profile from Google using the access token
func (as *AuthService) GetGoogleProfile(accessToken string) (*dto.GoogleOAuthProfile, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to fetch Google profile",
			"event", "auth.profile_fetch_failed",
			"error", err)
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Google profile fetch returned non-200 status",
			"event", "auth.profile_fetch_failed",
			"status_code", resp.StatusCode)
		return nil, fmt.Errorf("profile fetch failed with status: %d", resp.StatusCode)
	}

	var profile dto.GoogleOAuthProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		logger.Error("Failed to decode Google profile response",
			"event", "auth.profile_decode_failed",
			"error", err)
		return nil, fmt.Errorf("failed to decode profile: %w", err)
	}

	logger.Info("Google profile fetched successfully",
		"event", "auth.profile_fetched",
		"email", profile.Email)

	return &profile, nil
}

// RefreshTokens validates the refresh token and creates new JWT tokens
func (as *AuthService) RefreshTokens(refreshToken string) (*authdto.JWTTokenResponse, error) {
	logger.Info("Token refresh initiated",
		"event", "auth.token_refresh_started")

	accessToken, newRefreshToken, expiresAt, err := utils.RefreshToken(as.config, refreshToken)
	if err != nil {
		logger.Error("Token refresh failed",
			"event", "auth.token_refresh_failed",
			"error", err)
		return nil, fmt.Errorf("failed to refresh tokens: %w", err)
	}

	logger.Info("Token refresh successful",
		"event", "auth.token_refresh_completed",
		"expiresAt", expiresAt)

	return &authdto.JWTTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
	}, nil
}
