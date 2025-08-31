package dto

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Bio      string `json:"bio,omitempty" validate:"max=500"`
	Location string `json:"location,omitempty" validate:"max=100"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Bio      string `json:"bio,omitempty" validate:"max=500"`
	Location string `json:"location,omitempty" validate:"max=100"`
}