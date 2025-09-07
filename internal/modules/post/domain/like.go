package domain

import (
	"time"

	userDomain "github.com/topboyasante/pitstop/internal/modules/user/domain"
)

// Like represents a polymorphic like entity that can be applied to different types of content.
// This uses a polymorphic association pattern where:
// - LikableID stores the ID of the entity being liked (post ID, comment ID, etc.)
// - LikableType identifies what type of entity it is ("post", "comment", etc.)
//
// This approach allows us to:
// 1. Use a single likes table for all content types instead of separate tables
// 2. Reuse the same like functionality across posts, comments, and future content
// 3. Maintain consistent like behavior and counting across all "likable" entities
//
// Examples:
// - Like a post: LikableID="post-123", LikableType="post"
// - Like a comment: LikableID="comment-456", LikableType="comment"
type Like struct {
	ID          string                   `gorm:"primarykey" json:"id"`
	LikableID   string                   `gorm:"not null;index:idx_likable_user_like,unique" json:"likable_id" validate:"required"`
	LikableType string                   `gorm:"not null;index:idx_likable_user_like,unique" json:"likable_type" validate:"required,oneof=post comment"`
	UserID      string                   `gorm:"not null;index:idx_likable_user_like,unique" json:"user_id" validate:"required"`
	User        *userDomain.User         `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
}

// TableName specifies the table name for the Like model
func (Like) TableName() string {
	return "likes"
}

// Like type constants
const (
	LikableTypePost    = "post"
	LikableTypeComment = "comment"
)