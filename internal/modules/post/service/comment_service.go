package service

import (
	"fmt"

	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/modules/post/domain"
	"github.com/topboyasante/pitstop/internal/modules/post/dto"
	"github.com/topboyasante/pitstop/internal/modules/post/repository"
)

// CommentService handles comment business logic
type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
}

// NewCommentService creates a new comment service instance
func NewCommentService(commentRepo *repository.CommentRepository, postRepo *repository.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

// CreateComment creates a new comment on a post
func (s *CommentService) CreateComment(postID, userID string, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	logger.Info("Creating comment", "postID", postID, "userID", userID)

	// Verify post exists
	if _, err := s.postRepo.GetByID(postID); err != nil {
		logger.Error("Post not found", "postID", postID, "error", err)
		return nil, fmt.Errorf("post not found")
	}

	comment := &domain.Comment{
		PostID:   postID,
		UserID:   userID,
		Content:  req.Content,
		ParentID: nil,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		logger.Error("Failed to create comment", "error", err)
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Get the created comment with user info
	createdComment, err := s.commentRepo.GetByID(comment.ID)
	if err != nil {
		logger.Error("Failed to retrieve created comment", "error", err)
		return nil, fmt.Errorf("failed to retrieve created comment: %w", err)
	}

	return s.mapCommentToResponse(createdComment), nil
}

// CreateReply creates a reply to an existing comment
func (s *CommentService) CreateReply(postID, parentCommentID, userID string, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	logger.Info("Creating reply", "postID", postID, "parentCommentID", parentCommentID, "userID", userID)

	// Verify post exists
	if _, err := s.postRepo.GetByID(postID); err != nil {
		logger.Error("Post not found", "postID", postID, "error", err)
		return nil, fmt.Errorf("post not found")
	}

	// Verify parent comment exists and belongs to the post
	parentComment, err := s.commentRepo.GetByID(parentCommentID)
	if err != nil {
		logger.Error("Parent comment not found", "parentCommentID", parentCommentID, "error", err)
		return nil, fmt.Errorf("parent comment not found")
	}

	if parentComment.PostID != postID {
		logger.Error("Parent comment doesn't belong to post", "parentCommentID", parentCommentID, "postID", postID)
		return nil, fmt.Errorf("parent comment doesn't belong to this post")
	}

	comment := &domain.Comment{
		PostID:   postID,
		UserID:   userID,
		Content:  req.Content,
		ParentID: &parentCommentID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		logger.Error("Failed to create reply", "error", err)
		return nil, fmt.Errorf("failed to create reply: %w", err)
	}

	// Get the created comment with user info
	createdComment, err := s.commentRepo.GetByID(comment.ID)
	if err != nil {
		logger.Error("Failed to retrieve created reply", "error", err)
		return nil, fmt.Errorf("failed to retrieve created reply: %w", err)
	}

	return s.mapCommentToResponse(createdComment), nil
}

// GetCommentsByPostID retrieves all comments for a post
func (s *CommentService) GetCommentsByPostID(postID string) ([]dto.CommentResponse, error) {
	logger.Info("Getting comments for post", "postID", postID)

	comments, err := s.commentRepo.GetByPostID(postID)
	if err != nil {
		logger.Error("Failed to retrieve comments", "error", err)
		return nil, fmt.Errorf("failed to retrieve comments: %w", err)
	}

	responses := make([]dto.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		responses = append(responses, *s.mapCommentToResponse(&comment))
	}

	return responses, nil
}

// mapCommentToResponse converts domain comment to DTO response
func (s *CommentService) mapCommentToResponse(comment *domain.Comment) *dto.CommentResponse {
	response := &dto.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		LikeCount: comment.LikeCount,
		CreatedAt: comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if comment.ParentID != nil {
		response.ParentID = comment.ParentID
	}

	if comment.User != nil {
		response.User = &dto.UserResponse{
			Username:    comment.User.Username,
			DisplayName: comment.User.DisplayName,
			AvatarURL:   comment.User.AvatarURL,
		}
	}

	if comment.Parent != nil {
		response.Parent = s.mapCommentToResponse(comment.Parent)
	}

	if len(comment.Replies) > 0 {
		for _, reply := range comment.Replies {
			response.Replies = append(response.Replies, *s.mapCommentToResponse(&reply))
		}
	}

	return response
}