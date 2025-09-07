package domain

import (
	"time"

	userDomain "github.com/topboyasante/pitstop/internal/modules/user/domain"
)

// Post represents a post entity
type Post struct {
	ID           string                   `gorm:"primarykey" json:"id"`
	UserID       string                   `gorm:"not null" json:"user_id" validate:"required"`
	User         *userDomain.User         `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Content      string                   `gorm:"type:text" json:"content" validate:"required"`
	Comments     []Comment                `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	CommentCount int64                    `gorm:"-" json:"comment_count"`
	LikeCount    int64                    `gorm:"-" json:"like_count"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

// TableName specifies the table name for the Post model
func (Post) TableName() string {
	return "posts"
}
