package service

import (
	"errors"
	"fmt"

	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/modules/user/dto"
	"github.com/topboyasante/pitstop/internal/modules/user/repository"
	"github.com/topboyasante/pitstop/internal/shared/events"
	"gorm.io/gorm"
)

// FollowService handles follow business logic
type FollowService struct {
	followRepo *repository.FollowRepository
	userRepo   *repository.UserRepository
	eventBus   *events.EventBus
}

// NewFollowService creates a new follow service instance
func NewFollowService(followRepo *repository.FollowRepository, userRepo *repository.UserRepository, eventBus *events.EventBus) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
		eventBus:   eventBus,
	}
}

// ToggleFollow toggles a follow relationship (follow/unfollow)
func (s *FollowService) ToggleFollow(followerID, followingID string) (*dto.FollowToggleResponse, error) {
	// Check if the user being followed exists
	followingUser, err := s.userRepo.GetByID(followingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	// Toggle the follow
	isFollowing, err := s.followRepo.ToggleFollow(followerID, followingID)
	if err != nil {
		logger.Error("Failed to toggle follow", "follower_id", followerID, "following_id", followingID, "error", err)
		return nil, fmt.Errorf("failed to toggle follow: %w", err)
	}

	// Get updated counts
	followerCount, err := s.followRepo.CountFollowers(followingID)
	if err != nil {
		logger.Error("Failed to get follower count", "user_id", followingID, "error", err)
		return nil, fmt.Errorf("failed to get follower count: %w", err)
	}

	followingCount, err := s.followRepo.CountFollowing(followerID)
	if err != nil {
		logger.Error("Failed to get following count", "user_id", followerID, "error", err)
		return nil, fmt.Errorf("failed to get following count: %w", err)
	}

	// Log the action
	action := "followed"
	if !isFollowing {
		action = "unfollowed"
	}
	logger.Info(fmt.Sprintf("User %s %s", action, followingUser.Username), "follower_id", followerID, "following_id", followingID)

	return &dto.FollowToggleResponse{
		IsFollowing:    isFollowing,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
	}, nil
}

// GetFollowers retrieves all followers for a user
func (s *FollowService) GetFollowers(userID string) (*dto.FollowersResponse, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	// Get followers
	followers, err := s.followRepo.GetFollowers(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve followers: %w", err)
	}

	// Convert to response DTOs
	followerResponses := make([]dto.FollowUserResponse, len(followers))
	for i, follower := range followers {
		followerResponses[i] = dto.FollowUserResponse{
			ID:          follower.ID,
			Username:    follower.Username,
			DisplayName: follower.DisplayName,
			AvatarURL:   follower.AvatarURL,
			Bio:         follower.Bio,
		}
	}

	return &dto.FollowersResponse{
		Followers:  followerResponses,
		TotalCount: int64(len(followerResponses)),
	}, nil
}

// GetFollowing retrieves all users that a user is following
func (s *FollowService) GetFollowing(userID string) (*dto.FollowingResponse, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	// Get following
	following, err := s.followRepo.GetFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve following: %w", err)
	}

	// Convert to response DTOs
	followingResponses := make([]dto.FollowUserResponse, len(following))
	for i, user := range following {
		followingResponses[i] = dto.FollowUserResponse{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			AvatarURL:   user.AvatarURL,
			Bio:         user.Bio,
		}
	}

	return &dto.FollowingResponse{
		Following:  followingResponses,
		TotalCount: int64(len(followingResponses)),
	}, nil
}

// CheckUserFollowing checks if a user is following another user
func (s *FollowService) CheckUserFollowing(followerID, followingID string) (bool, error) {
	exists, err := s.followRepo.Exists(followerID, followingID)
	if err != nil {
		return false, fmt.Errorf("failed to check follow status: %w", err)
	}
	return exists, nil
}