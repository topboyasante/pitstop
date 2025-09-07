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
func (r *LikeRepository) Delete(postID, userID string) error {
	result := r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&domain.Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}
	return nil
}

// GetByPostAndUser retrieves a like by post and user
func (r *LikeRepository) GetByPostAndUser(postID, userID string) (*domain.Like, error) {
	var like domain.Like
	err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// Exists checks if a like exists for a post and user
func (r *LikeRepository) Exists(postID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Like{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetLikesByPost retrieves all likes for a post with user information
func (r *LikeRepository) GetLikesByPost(postID string) ([]domain.Like, error) {
	var likes []domain.Like
	err := r.db.Preload("User").Where("post_id = ?", postID).Order("created_at DESC").Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

// CountLikesByPost counts likes for a specific post
func (r *LikeRepository) CountLikesByPost(postID string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Like{}).Where("post_id = ?", postID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ToggleLike toggles a like for a post (like/unlike)
func (r *LikeRepository) ToggleLike(postID, userID string) (bool, error) {
	// Check if like already exists
	exists, err := r.Exists(postID, userID)
	if err != nil {
		return false, err
	}

	if exists {
		// Unlike: delete the like
		err = r.Delete(postID, userID)
		if err != nil {
			return false, err
		}
		return false, nil // false means unliked
	} else {
		// Like: create new like
		like := &domain.Like{
			PostID: postID,
			UserID: userID,
		}
		err = r.Create(like)
		if err != nil {
			return false, err
		}
		return true, nil // true means liked
	}
}