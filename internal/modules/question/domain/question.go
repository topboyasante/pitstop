package domain

import (
	"time"

	postDomain "github.com/topboyasante/pitstop/internal/modules/post/domain"
	userDomain "github.com/topboyasante/pitstop/internal/modules/user/domain"
)

// Question represents a question entity
type Question struct {
	ID           string                     `gorm:"primarykey" json:"id"`
	UserID       string                     `gorm:"not null" json:"user_id" validate:"required"`
	User         *userDomain.User           `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Title        string                     `gorm:"type:varchar(255)" json:"title" validate:"required"`
	Content      string                     `gorm:"type:text" json:"content" validate:"required"`
	Tags         string                     `gorm:"type:varchar(500)" json:"tags"` // Comma-separated tags
	IsAnswered   bool                       `gorm:"default:false" json:"is_answered"`
	Comments     []postDomain.Comment       `gorm:"polymorphic:Commentable;polymorphicValue:question" json:"comments,omitempty"`
	CommentCount int64                      `gorm:"-" json:"comment_count"`
	LikeCount    int64                      `gorm:"-" json:"like_count"`
	AnswerCount  int64                      `gorm:"-" json:"answer_count"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
}

// TableName specifies the table name for the Question model
func (Question) TableName() string {
	return "questions"
}

// Answer represents an answer to a question
type Answer struct {
	ID         string               `gorm:"primarykey" json:"id"`
	QuestionID string               `gorm:"not null" json:"question_id" validate:"required"`
	Question   *Question            `gorm:"foreignKey:QuestionID;references:ID" json:"question,omitempty"`
	UserID     string               `gorm:"not null" json:"user_id" validate:"required"`
	User       *userDomain.User     `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Content    string               `gorm:"type:text" json:"content" validate:"required"`
	IsAccepted bool                 `gorm:"default:false" json:"is_accepted"`
	Comments   []postDomain.Comment `gorm:"polymorphic:Commentable;polymorphicValue:answer" json:"comments,omitempty"`
	LikeCount  int64                `gorm:"-" json:"like_count"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

// TableName specifies the table name for the Answer model
func (Answer) TableName() string {
	return "answers"
}