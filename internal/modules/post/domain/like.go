package domain

import (
	"time"

	userDomain "github.com/topboyasante/pitstop/internal/modules/user/domain"
)

// Like represents a like entity for posts
type Like struct {
	ID        string                   `gorm:"primarykey" json:"id"`
	PostID    string                   `gorm:"not null;index:idx_post_user_like,unique" json:"post_id" validate:"required"`
	UserID    string                   `gorm:"not null;index:idx_post_user_like,unique" json:"user_id" validate:"required"`
	User      *userDomain.User         `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Post      *Post                    `gorm:"foreignKey:PostID;references:ID" json:"post,omitempty"`
	CreatedAt time.Time                `json:"created_at"`
}

// TableName specifies the table name for the Like model
func (Like) TableName() string {
	return "likes"
}