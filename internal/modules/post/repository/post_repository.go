package repository

import (
	"github.com/google/uuid"
	"github.com/topboyasante/pitstop/internal/modules/post/domain"
	"gorm.io/gorm"
)

// PostRepository handles post data operations
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository instance
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create creates a new post
func (r *PostRepository) Create(post *domain.Post) error {
	if post.ID == "" {
		post.ID = uuid.NewString()
	}
	return r.db.Create(post).Error
}

// GetByID retrieves a post by ID
func (r *PostRepository) GetByID(id string) (*domain.Post, error) {
	var post domain.Post
	err := r.db.Preload("User").Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetAll retrieves all posts with pagination
func (r *PostRepository) GetAll(page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var totalCount int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&domain.Post{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get posts
	if err := r.db.Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, totalCount, nil
}

// Update updates a post
func (r *PostRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

// Delete soft deletes a post
func (r *PostRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Post{}).Error
}
