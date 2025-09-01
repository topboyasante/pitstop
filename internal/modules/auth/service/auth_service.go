package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/logger"
)

type AuthService struct {
	config    *config.Config
	redis     *redis.Client
	validator *validator.Validate
}

// NewAuthService creates a new instance of AuthService with the provided configuration
func NewAuthService(config *config.Config, redis *redis.Client, validator *validator.Validate) *AuthService {
	logger.Info("Initializing auth service")
	return &AuthService{
		config:    config,
		redis:     redis,
		validator: validator,
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

// ExchangeCode validates the state token and exchanges the authorization code for an access token
func (as *AuthService) ExchangeCode(code, state string) (string, error) {
	logger.Info("Token exchange initiated",
		"event", "auth.token_exchange_started")

	if !as.ValidateState(state) {
		logger.Error("Token exchange failed",
			"event", "auth.token_exchange_failed",
			"reason", "invalid_state")
		return "", fmt.Errorf("invalid state")
	}

	token, err := as.config.OAuth.Exchange(context.TODO(), code)
	if err != nil {
		logger.Error("Token exchange failed",
			"event", "auth.token_exchange_failed",
			"reason", "oauth_exchange_error",
			"error", err.Error())
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	logger.Info("Token exchange successful",
		"event", "auth.token_exchange_completed",
		"token_type", token.TokenType)

	return token.AccessToken, nil
}
