package domain

import (
	"time"
)

// Post represents a post entity
type Post struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    string    `gorm:"not null" json:"user_id" validate:"required"`
	Content   string    `gorm:"type:text" json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for the Post model
func (Post) TableName() string {
	return "posts"
}
