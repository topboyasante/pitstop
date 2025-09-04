package dto

import "time"

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
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostsResponse represents a paginated list of posts
type PostsResponse struct {
	Posts      []PostResponse `json:"posts"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	HasNext    bool           `json:"has_next"`
}
