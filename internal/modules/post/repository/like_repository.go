package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/topboyasante/pitstop/internal/modules/post/domain"
	"gorm.io/gorm"
)

// LikeRepository handles like data operations
type LikeRepository struct {
	db *gorm.DB
}

// NewLikeRepository creates a new like repository instance
func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

// Create creates a new like (toggle on)
func (r *LikeRepository) Create(like *domain.Like) error {
	if like.ID == "" {
		like.ID = uuid.NewString()
	}
	return r.db.Create(like).Error
}

// Delete removes a like (toggle off)
func (r *LikeRepository) Delete(likableID, likableType, userID string) error {
	result := r.db.Where("likable_id = ? AND likable_type = ? AND user_id = ?", likableID, likableType, userID).Delete(&domain.Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}
	return nil
}

// GetByLikableAndUser retrieves a like by likable entity and user
func (r *LikeRepository) GetByLikableAndUser(likableID, likableType, userID string) (*domain.Like, error) {
	var like domain.Like
	err := r.db.Where("likable_id = ? AND likable_type = ? AND user_id = ?", likableID, likableType, userID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// Exists checks if a like exists for a likable entity and user
func (r *LikeRepository) Exists(likableID, likableType, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Like{}).Where("likable_id = ? AND likable_type = ? AND user_id = ?", likableID, likableType, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetLikesByLikable retrieves all likes for a likable entity with user information
func (r *LikeRepository) GetLikesByLikable(likableID, likableType string) ([]domain.Like, error) {
	var likes []domain.Like
	err := r.db.Preload("User").Where("likable_id = ? AND likable_type = ?", likableID, likableType).Order("created_at DESC").Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

// CountLikesByLikable counts likes for a specific likable entity
func (r *LikeRepository) CountLikesByLikable(likableID, likableType string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Like{}).Where("likable_id = ? AND likable_type = ?", likableID, likableType).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ToggleLike toggles a like for a likable entity (like/unlike)
func (r *LikeRepository) ToggleLike(likableID, likableType, userID string) (bool, error) {
	// Check if like already exists
	exists, err := r.Exists(likableID, likableType, userID)
	if err != nil {
		return false, err
	}

	if exists {
		// Unlike: delete the like
		err = r.Delete(likableID, likableType, userID)
		if err != nil {
			return false, err
		}
		return false, nil // false means unliked
	} else {
		// Like: create new like
		like := &domain.Like{
			LikableID:   likableID,
			LikableType: likableType,
			UserID:      userID,
		}
		err = r.Create(like)
		if err != nil {
			return false, err
		}
		return true, nil // true means liked
	}
}

// Legacy methods for backward compatibility - these delegate to the polymorphic methods

// DeletePostLike removes a like from a post (legacy method)
func (r *LikeRepository) DeletePostLike(postID, userID string) error {
	return r.Delete(postID, domain.LikableTypePost, userID)
}

// ExistsForPost checks if a like exists for a post and user (legacy method)
func (r *LikeRepository) ExistsForPost(postID, userID string) (bool, error) {
	return r.Exists(postID, domain.LikableTypePost, userID)
}

// GetLikesByPost retrieves all likes for a post with user information (legacy method)
func (r *LikeRepository) GetLikesByPost(postID string) ([]domain.Like, error) {
	return r.GetLikesByLikable(postID, domain.LikableTypePost)
}

// CountLikesByPost counts likes for a specific post (legacy method)
func (r *LikeRepository) CountLikesByPost(postID string) (int64, error) {
	return r.CountLikesByLikable(postID, domain.LikableTypePost)
}

// TogglePostLike toggles a like for a post (legacy method)
func (r *LikeRepository) TogglePostLike(postID, userID string) (bool, error) {
	return r.ToggleLike(postID, domain.LikableTypePost, userID)
}