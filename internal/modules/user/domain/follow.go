package domain

import (
	"time"
)

// Follow represents a following relationship between users.
// This implements a self-referential many-to-many relationship where:
// - FollowerID is the user who is following someone
// - FollowingID is the user being followed
//
// Example: If Alice follows Bob:
// - FollowerID = Alice's ID
// - FollowingID = Bob's ID
//
// This allows us to query:
// - Alice's following list (where FollowerID = Alice's ID)
// - Bob's followers list (where FollowingID = Bob's ID)
type Follow struct {
	ID          string    `gorm:"primarykey" json:"id"`
	FollowerID  string    `gorm:"not null;index:idx_follower_following,unique" json:"follower_id" validate:"required"`
	FollowingID string    `gorm:"not null;index:idx_follower_following,unique" json:"following_id" validate:"required"`
	Follower    *User     `gorm:"foreignKey:FollowerID;references:ID" json:"follower,omitempty"`
	Following   *User     `gorm:"foreignKey:FollowingID;references:ID" json:"following,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName specifies the table name for the Follow model
func (Follow) TableName() string {
	return "follows"
}