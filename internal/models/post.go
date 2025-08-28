package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Title     string         `gorm:"not null;size:255" json:"title" validate:"required,min=1,max=255"`
	Content   string         `gorm:"type:text" json:"content" validate:"required,min=1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}