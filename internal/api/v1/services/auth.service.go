package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/topboyasante/pitstop/internal/config"
)

type AuthService struct {
	config *config.Config
	states map[string]bool // In production, use Redis or database
}

// NewAuthService creates a new instance of AuthService with the provided configuration
func NewAuthService(config *config.Config) *AuthService {
	return &AuthService{
		config: config,
		states: make(map[string]bool),
	}
}

// Authenticate generates a CSRF state token and returns the OAuth authorization URL
func (as *AuthService) Authenticate() string {
	state := as.generateState()
	as.states[state] = true
	url := as.config.OAuth.AuthCodeURL(state)
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
		return true
	}
	return false
}

// ExchangeCode validates the state token and exchanges the authorization code for an access token
func (as *AuthService) ExchangeCode(code, state string) (string, error) {
	if !as.ValidateState(state) {
		return "", fmt.Errorf("invalid state")
	}

	token, err := as.config.OAuth.Exchange(context.TODO(), code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	return token.AccessToken, nil
}
