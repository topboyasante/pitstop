package service

import (
	"errors"
	"fmt"

	"github.com/topboyasante/pitstop/internal/core/logger"
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

// ToggleLike toggles a like for a post (like/unlike)
func (s *LikeService) ToggleLike(postID, userID string) (*dto.LikeToggleResponse, error) {
	// Check if post exists
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to check post: %w", err)
	}

	// Toggle the like
	liked, err := s.likeRepo.ToggleLike(postID, userID)
	if err != nil {
		logger.Error("Failed to toggle like", "post_id", postID, "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to toggle like: %w", err)
	}

	// Get updated like count
	likeCount, err := s.likeRepo.CountLikesByPost(postID)
	if err != nil {
		logger.Error("Failed to get like count", "post_id", postID, "error", err)
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
	likes, err := s.likeRepo.GetLikesByPost(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve likes: %w", err)
	}

	// Convert to response DTOs
	likeResponses := make([]dto.LikeResponse, len(likes))
	for i, like := range likes {
		likeResponses[i] = dto.LikeResponse{
			ID:        like.ID,
			PostID:    like.PostID,
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

// CheckUserLiked checks if a user has liked a specific post
func (s *LikeService) CheckUserLiked(postID, userID string) (bool, error) {
	exists, err := s.likeRepo.Exists(postID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check user like status: %w", err)
	}
	return exists, nil
}