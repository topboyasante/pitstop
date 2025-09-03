package dto

import "time"

// UserResponse represents a user in API responses
type UserResponse struct {
	ID             uint      `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Bio            string    `json:"bio"`
	Location       string    `json:"location"`
	Reputation     int       `json:"reputation"`
	FollowerCount  int       `json:"follower_count"`
	FollowingCount int       `json:"following_count"`
	IsVerified     bool      `json:"is_verified"`
	CreatedAt      time.Time `json:"created_at"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// AuthURLResponse represents the OAuth authorization URL response
type AuthURLResponse struct {
	AuthURL string `json:"auth_url"`
}

// JWTTokenResponse represents a JWT token response
type JWTTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}