package dto

import (
	"time"
)

// LikeUserResponse represents user data included with likes (limited fields)
type LikeUserResponse struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// LikeResponse represents a like in API responses
type LikeResponse struct {
	ID        string            `json:"id"`
	PostID    string            `json:"post_id"`
	UserID    string            `json:"user_id"`
	User      *LikeUserResponse `json:"user,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// LikesResponse represents a list of likes for a post
type LikesResponse struct {
	Likes      []LikeResponse `json:"likes"`
	TotalCount int64          `json:"total_count"`
}

// LikeToggleResponse represents the response after toggling a like
type LikeToggleResponse struct {
	Liked     bool  `json:"liked"`
	LikeCount int64 `json:"like_count"`
}