package dto

import (
	"time"
)

// QuestionUserResponse represents user data included with questions (limited fields)
type QuestionUserResponse struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// CreateQuestionRequest represents a request to create a new question
type CreateQuestionRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	Title   string `json:"title" validate:"required,min=5,max=255"`
	Content string `json:"content" validate:"required,min=10"`
	Tags    string `json:"tags,omitempty"` // Comma-separated tags
}

// UpdateQuestionRequest represents a request to update a question
type UpdateQuestionRequest struct {
	Title   string `json:"title,omitempty" validate:"omitempty,min=5,max=255"`
	Content string `json:"content,omitempty" validate:"omitempty,min=10"`
	Tags    string `json:"tags,omitempty"`
}

// QuestionResponse represents a question in API responses
type QuestionResponse struct {
	ID           string                `json:"id"`
	UserID       string                `json:"user_id"`
	Title        string                `json:"title"`
	Content      string                `json:"content"`
	Tags         []string              `json:"tags"` // Will be split from comma-separated string
	IsAnswered   bool                  `json:"is_answered"`
	User         *QuestionUserResponse `json:"user"`
	CommentCount int64                 `json:"comment_count"`
	LikeCount    int64                 `json:"like_count"`
	AnswerCount  int64                 `json:"answer_count"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// QuestionsResponse represents a paginated list of questions
type QuestionsResponse struct {
	Questions  []QuestionResponse `json:"questions"`
	TotalCount int64              `json:"total_count"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	HasNext    bool               `json:"has_next"`
}

// CreateAnswerRequest represents a request to create a new answer
type CreateAnswerRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,min=10"`
}

// UpdateAnswerRequest represents a request to update an answer
type UpdateAnswerRequest struct {
	Content string `json:"content,omitempty" validate:"omitempty,min=10"`
}

// AcceptAnswerRequest represents a request to accept an answer
type AcceptAnswerRequest struct {
	IsAccepted bool `json:"is_accepted"`
}

// AnswerResponse represents an answer in API responses
type AnswerResponse struct {
	ID         string                `json:"id"`
	QuestionID string                `json:"question_id"`
	UserID     string                `json:"user_id"`
	Content    string                `json:"content"`
	IsAccepted bool                  `json:"is_accepted"`
	User       *QuestionUserResponse `json:"user"`
	LikeCount  int64                 `json:"like_count"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

// AnswersResponse represents a paginated list of answers
type AnswersResponse struct {
	Answers    []AnswerResponse `json:"answers"`
	TotalCount int64            `json:"total_count"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	HasNext    bool             `json:"has_next"`
}