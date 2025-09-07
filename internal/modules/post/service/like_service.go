package service

import (
	"errors"
	"fmt"

	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/modules/post/domain"
	"github.com/topboyasante/pitstop/internal/modules/post/dto"
	"github.com/topboyasante/pitstop/internal/modules/post/repository"
	"github.com/topboyasante/pitstop/internal/shared/events"
	"gorm.io/gorm"
)

// LikeService handles like business logic
type LikeService struct {
	likeRepo *repository.LikeRepository
	postRepo *repository.PostRepository
	eventBus *events.EventBus
}

// NewLikeService creates a new like service instance
func NewLikeService(likeRepo *repository.LikeRepository, postRepo *repository.PostRepository, eventBus *events.EventBus) *LikeService {
	return &LikeService{
		likeRepo: likeRepo,
		postRepo: postRepo,
		eventBus: eventBus,
	}
}

// TogglePostLike toggles a like for a post (like/unlike)
func (s *LikeService) TogglePostLike(postID, userID string) (*dto.LikeToggleResponse, error) {
	// Check if post exists
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to check post: %w", err)
	}

	// Toggle the like
	liked, err := s.likeRepo.ToggleLike(postID, domain.LikableTypePost, userID)
	if err != nil {
		logger.Error("Failed to toggle post like", "post_id", postID, "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to toggle like: %w", err)
	}

	// Get updated like count
	likeCount, err := s.likeRepo.CountLikesByLikable(postID, domain.LikableTypePost)
	if err != nil {
		logger.Error("Failed to get post like count", "post_id", postID, "error", err)
		return nil, fmt.Errorf("failed to get like count: %w", err)
	}

	// Log the action
	action := "liked"
	if !liked {
		action = "unliked"
	}
	logger.Info(fmt.Sprintf("Post %s", action), "post_id", postID, "user_id", userID)

	return &dto.LikeToggleResponse{
		Liked:     liked,
		LikeCount: likeCount,
	}, nil
}

// ToggleCommentLike toggles a like for a comment (like/unlike)
func (s *LikeService) ToggleCommentLike(postID, commentID, userID string) (*dto.LikeToggleResponse, error) {
	// Check if post exists
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to check post: %w", err)
	}

	// TODO: Add comment repository to validate comment exists and belongs to post
	// For now we'll trust the comment ID is valid

	// Toggle the like
	liked, err := s.likeRepo.ToggleLike(commentID, domain.LikableTypeComment, userID)
	if err != nil {
		logger.Error("Failed to toggle comment like", "comment_id", commentID, "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to toggle like: %w", err)
	}

	// Get updated like count
	likeCount, err := s.likeRepo.CountLikesByLikable(commentID, domain.LikableTypeComment)
	if err != nil {
		logger.Error("Failed to get comment like count", "comment_id", commentID, "error", err)
		return nil, fmt.Errorf("failed to get like count: %w", err)
	}

	// Log the action
	action := "liked"
	if !liked {
		action = "unliked"
	}
	logger.Info(fmt.Sprintf("Comment %s", action), "comment_id", commentID, "user_id", userID)

	return &dto.LikeToggleResponse{
		Liked:     liked,
		LikeCount: likeCount,
	}, nil
}

// GetLikesByPost retrieves all likes for a post
func (s *LikeService) GetLikesByPost(postID string) (*dto.LikesResponse, error) {
	// Check if post exists
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to check post: %w", err)
	}

	// Get likes
	likes, err := s.likeRepo.GetLikesByLikable(postID, domain.LikableTypePost)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve likes: %w", err)
	}

	// Convert to response DTOs
	likeResponses := make([]dto.LikeResponse, len(likes))
	for i, like := range likes {
		likeResponses[i] = dto.LikeResponse{
			ID:        like.ID,
			PostID:    like.LikableID, // For backward compatibility in DTO
			UserID:    like.UserID,
			CreatedAt: like.CreatedAt,
		}

		if like.User != nil {
			likeResponses[i].User = &dto.LikeUserResponse{
				Username:    like.User.Username,
				DisplayName: like.User.DisplayName,
				AvatarURL:   like.User.AvatarURL,
			}
		}
	}

	return &dto.LikesResponse{
		Likes:      likeResponses,
		TotalCount: int64(len(likeResponses)),
	}, nil
}

// GetLikesByComment retrieves all likes for a comment
func (s *LikeService) GetLikesByComment(postID, commentID string) (*dto.LikesResponse, error) {
	// Check if post exists
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to check post: %w", err)
	}

	// TODO: Add comment validation

	// Get likes
	likes, err := s.likeRepo.GetLikesByLikable(commentID, domain.LikableTypeComment)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve comment likes: %w", err)
	}

	// Convert to response DTOs
	likeResponses := make([]dto.LikeResponse, len(likes))
	for i, like := range likes {
		likeResponses[i] = dto.LikeResponse{
			ID:        like.ID,
			PostID:    postID, // Include post ID for context
			UserID:    like.UserID,
			CreatedAt: like.CreatedAt,
		}

		if like.User != nil {
			likeResponses[i].User = &dto.LikeUserResponse{
				Username:    like.User.Username,
				DisplayName: like.User.DisplayName,
				AvatarURL:   like.User.AvatarURL,
			}
		}
	}

	return &dto.LikesResponse{
		Likes:      likeResponses,
		TotalCount: int64(len(likeResponses)),
	}, nil
}

// CheckUserLikedPost checks if a user has liked a specific post
func (s *LikeService) CheckUserLikedPost(postID, userID string) (bool, error) {
	exists, err := s.likeRepo.Exists(postID, domain.LikableTypePost, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check user like status: %w", err)
	}
	return exists, nil
}

// CheckUserLikedComment checks if a user has liked a specific comment
func (s *LikeService) CheckUserLikedComment(commentID, userID string) (bool, error) {
	exists, err := s.likeRepo.Exists(commentID, domain.LikableTypeComment, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check user like status: %w", err)
	}
	return exists, nil
}