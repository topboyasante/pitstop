package dto

import "time"

// CreateUserRequest represents OAuth user data
type CreateUserRequest struct {
	ProviderID  string `json:"provider_id" validate:"required"`
	Provider    string `json:"provider" validate:"required,oneof=google facebook github"`
	FirstName   string `json:"first_name" validate:"omitempty,max=255"`
	LastName    string `json:"last_name" validate:"omitempty,max=255"`
	Email       string `json:"email" validate:"required,email,max=255"`
	AvatarURL   string `json:"avatar_url" validate:"omitempty,url,max=500"`
	Locale      string `json:"locale" validate:"omitempty,max=10"`
}

// UpdateUserRequest represents a request to update user profile
type UpdateUserRequest struct {
	FirstName   string `json:"first_name" validate:"omitempty,max=255"`
	LastName    string `json:"last_name" validate:"omitempty,max=255"`
	Username    string `json:"username" validate:"omitempty,min=3,max=100,alphanum"`
	DisplayName string `json:"display_name" validate:"omitempty,max=150"`
	Bio         string `json:"bio" validate:"omitempty,max=500"`
	AvatarURL   string `json:"avatar_url" validate:"omitempty,url,max=500"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID          string    `json:"id"`
	Provider    string    `json:"provider"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username,omitempty"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	FullName    string    `json:"full_name"`
	IsOAuth     bool      `json:"is_oauth"`
	CreatedAt   time.Time `json:"created_at"`
}

// UsersResponse represents a paginated list of users
type UsersResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	HasNext    bool           `json:"has_next"`
}

// OAuthCallbackResponse represents the response after OAuth callback
type OAuthCallbackResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
	IsNewUser   bool         `json:"is_new_user"`
}

// GoogleOAuthProfile represents Google OAuth profile data
type GoogleOAuthProfile struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	VerifiedEmail bool   `json:"verified_email"`
}
