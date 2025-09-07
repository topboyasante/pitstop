package dto

import (
	"time"
)

// FollowUserResponse represents user data in follow responses (limited fields)
type FollowUserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio,omitempty"`
}

// FollowResponse represents a follow relationship in API responses
type FollowResponse struct {
	ID        string              `json:"id"`
	Follower  *FollowUserResponse `json:"follower,omitempty"`
	Following *FollowUserResponse `json:"following,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
}

// FollowersResponse represents a list of followers
type FollowersResponse struct {
	Followers  []FollowUserResponse `json:"followers"`
	TotalCount int64                `json:"total_count"`
}

// FollowingResponse represents a list of users being followed
type FollowingResponse struct {
	Following  []FollowUserResponse `json:"following"`
	TotalCount int64                `json:"total_count"`
}

// FollowToggleResponse represents the response after toggling a follow
type FollowToggleResponse struct {
	IsFollowing    bool  `json:"is_following"`
	FollowerCount  int64 `json:"follower_count"`
	FollowingCount int64 `json:"following_count"`
}