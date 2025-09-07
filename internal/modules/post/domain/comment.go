package domain

import (
	"time"

	userDomain "github.com/topboyasante/pitstop/internal/modules/user/domain"
)

// Comment represents a comment entity
type Comment struct {
	ID        string                   `gorm:"primarykey" json:"id"`
	PostID    string                   `gorm:"not null" json:"post_id" validate:"required"`
	UserID    string                   `gorm:"not null" json:"user_id" validate:"required"`
	User      *userDomain.User         `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	ParentID  *string                  `gorm:"default:null" json:"parent_id,omitempty"`
	Parent    *Comment                 `gorm:"foreignKey:ParentID;references:ID" json:"parent,omitempty"`
	Replies   []Comment                `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Content   string                   `gorm:"type:text" json:"content" validate:"required"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

// TableName specifies the table name for the Comment model
func (Comment) TableName() string {
	return "comments"
}