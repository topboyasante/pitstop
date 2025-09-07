package dto

// CreateCommentRequest represents a comment creation request
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// CommentResponse represents a comment response
type CommentResponse struct {
	ID        string          `json:"id"`
	PostID    string          `json:"post_id"`
	UserID    string          `json:"user_id"`
	User      *UserResponse   `json:"user,omitempty"`
	ParentID  *string         `json:"parent_id,omitempty"`
	Parent    *CommentResponse `json:"parent,omitempty"`
	Replies   []CommentResponse `json:"replies,omitempty"`
	Content   string          `json:"content"`
	LikeCount int64           `json:"like_count"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

// UserResponse represents a simplified user response for comments
type UserResponse struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}