package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user entity
type User struct {
	ID             string         `gorm:"primarykey" json:"id"`
	ProviderID     string         `gorm:"not null;size:100;index" json:"provider_id"`
	Provider       string         `gorm:"not null;size:50;index" json:"provider" validate:"required,oneof=google facebook github"`
	FirstName      string         `gorm:"size:255" json:"first_name" validate:"omitempty,max=255"`
	LastName       string         `gorm:"size:255" json:"last_name" validate:"omitempty,max=255"`
	Username       string         `gorm:"uniqueIndex;size:100" json:"username" validate:"omitempty,min=3,max=100,alphanum"`
	Email          string         `gorm:"uniqueIndex;not null;size:255" json:"email" validate:"required,email,max=255"`
	DisplayName    string         `gorm:"size:150" json:"display_name" validate:"omitempty,max=150"`
	Bio            string         `gorm:"size:500" json:"bio" validate:"omitempty,max=500"`
	AvatarURL      string         `gorm:"size:500" json:"avatar_url" validate:"omitempty,url,max=500"`
	Locale         string         `gorm:"size:10" json:"locale" validate:"omitempty,max=10"`
	FollowerCount  int64          `gorm:"-" json:"follower_count"`
	FollowingCount int64          `gorm:"-" json:"following_count"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}


