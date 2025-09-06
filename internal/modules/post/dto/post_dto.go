package dto

import (
	"time"
)

// PostUserResponse represents user data included with posts (limited fields)
type PostUserResponse struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// CreatePostRequest represents a request to create a new post
type CreatePostRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// UpdatePostRequest represents a request to update a post
type UpdatePostRequest struct {
	Content string `json:"content,omitempty" validate:"omitempty,min=1,max=255"`
}

// PostResponse represents a post in API responses
type PostResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Content   string            `json:"content"`
	User      *PostUserResponse `json:"user"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// PostsResponse represents a paginated list of posts
type PostsResponse struct {
	Posts      []PostResponse `json:"posts"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	HasNext    bool           `json:"has_next"`
}
