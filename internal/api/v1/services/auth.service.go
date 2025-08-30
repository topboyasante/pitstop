package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/topboyasante/pitstop/internal/config"
	"github.com/topboyasante/pitstop/internal/logger"
)

type AuthService struct {
	config *config.Config
	states map[string]bool // In production, use Redis or database
}

// NewAuthService creates a new instance of AuthService with the provided configuration
func NewAuthService(config *config.Config) *AuthService {
	logger.Info("Initializing auth service")
	return &AuthService{
		config: config,
		states: make(map[string]bool),
	}
}

// Authenticate generates a CSRF state token and returns the OAuth authorization URL
func (as *AuthService) Authenticate() string {
	logger.Info("OAuth authentication initiated", 
		"event", "auth.oauth_started",
		"provider", "google")
	
	state := as.generateState()
	as.states[state] = true
	url := as.config.OAuth.AuthCodeURL(state)
	
	logger.Info("OAuth URL generated", 
		"event", "auth.oauth_url_generated",
		"provider", "google")
	
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
	if as.states[state] {
		delete(as.states, state) // Use once
		logger.Info("State validation successful", 
			"event", "auth.state_validated")
		return true
	}
	logger.Warn("State validation failed", 
		"event", "auth.state_validation_failed",
		"reason", "invalid_state")
	return false
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
