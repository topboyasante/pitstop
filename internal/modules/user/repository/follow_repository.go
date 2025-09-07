package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/topboyasante/pitstop/internal/modules/user/domain"
	"gorm.io/gorm"
)

// FollowRepository handles follow data operations
type FollowRepository struct {
	db *gorm.DB
}

// NewFollowRepository creates a new follow repository instance
func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

// Create creates a new follow relationship
func (r *FollowRepository) Create(follow *domain.Follow) error {
	if follow.ID == "" {
		follow.ID = uuid.NewString()
	}
	return r.db.Create(follow).Error
}

// Delete removes a follow relationship
func (r *FollowRepository) Delete(followerID, followingID string) error {
	result := r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&domain.Follow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("follow relationship not found")
	}
	return nil
}

// Exists checks if a follow relationship exists
func (r *FollowRepository) Exists(followerID, followingID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetFollowers retrieves all followers for a user with user information
func (r *FollowRepository) GetFollowers(userID string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Table("users").
		Select("users.*").
		Joins("INNER JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Order("follows.created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Calculate follower/following counts for each user
	for i := range users {
		r.calculateFollowCounts(&users[i])
	}

	return users, nil
}

// GetFollowing retrieves all users that a user is following with user information
func (r *FollowRepository) GetFollowing(userID string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Table("users").
		Select("users.*").
		Joins("INNER JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ?", userID).
		Order("follows.created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Calculate follower/following counts for each user
	for i := range users {
		r.calculateFollowCounts(&users[i])
	}

	return users, nil
}

// CountFollowers counts how many followers a user has
func (r *FollowRepository) CountFollowers(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Follow{}).Where("following_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountFollowing counts how many users a user is following
func (r *FollowRepository) CountFollowing(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ToggleFollow toggles a follow relationship (follow/unfollow)
func (r *FollowRepository) ToggleFollow(followerID, followingID string) (bool, error) {
	// Prevent self-following
	if followerID == followingID {
		return false, errors.New("users cannot follow themselves")
	}

	// Check if follow already exists
	exists, err := r.Exists(followerID, followingID)
	if err != nil {
		return false, err
	}

	if exists {
		// Unfollow: delete the follow
		err = r.Delete(followerID, followingID)
		if err != nil {
			return false, err
		}
		return false, nil // false means unfollowed
	} else {
		// Follow: create new follow
		follow := &domain.Follow{
			FollowerID:  followerID,
			FollowingID: followingID,
		}
		err = r.Create(follow)
		if err != nil {
			return false, err
		}
		return true, nil // true means followed
	}
}

// calculateFollowCounts calculates and sets follower/following counts for a user
func (r *FollowRepository) calculateFollowCounts(user *domain.User) {
	// Get follower count
	followerCount, _ := r.CountFollowers(user.ID)
	user.FollowerCount = followerCount

	// Get following count
	followingCount, _ := r.CountFollowing(user.ID)
	user.FollowingCount = followingCount
}