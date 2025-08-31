package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Username       string         `gorm:"uniqueIndex;not null;size:50" json:"username" validate:"required,min=3,max=50"`
	Bio            string         `gorm:"size:500" json:"bio"`
	Location       string         `gorm:"size:100" json:"location"`
	Reputation     int            `gorm:"default:0" json:"reputation"`
	FollowerCount  int            `gorm:"default:0" json:"follower_count"`
	FollowingCount int            `gorm:"default:0" json:"following_count"`
	IsVerified     bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}