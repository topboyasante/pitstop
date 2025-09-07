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

// GetByID retrieves a post by ID with comments and comment count
func (r *PostRepository) GetByID(id string) (*domain.Post, error) {
	var post domain.Post
	err := r.db.Preload("User").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").
				Preload("Replies", func(db *gorm.DB) *gorm.DB {
					return db.Preload("User").Order("created_at ASC")
				}).
				Where("parent_id IS NULL").
				Order("created_at DESC")
		}).
		Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}

	// Calculate total comment count (including replies)
	var commentCount int64
	r.db.Model(&domain.Comment{}).Where("post_id = ?", id).Count(&commentCount)
	post.CommentCount = commentCount

	// Calculate like count
	var likeCount int64
	r.db.Model(&domain.Like{}).Where("likable_id = ? AND likable_type = ?", id, domain.LikableTypePost).Count(&likeCount)
	post.LikeCount = likeCount

	return &post, nil
}

// GetAll retrieves all posts with pagination and comment counts
func (r *PostRepository) GetAll(page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var totalCount int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&domain.Post{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get posts with comment counts
	if err := r.db.Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// Calculate comment and like counts for each post
	for i := range posts {
		var commentCount int64
		r.db.Model(&domain.Comment{}).Where("post_id = ?", posts[i].ID).Count(&commentCount)
		posts[i].CommentCount = commentCount
		
		var likeCount int64
		r.db.Model(&domain.Like{}).Where("likable_id = ? AND likable_type = ?", posts[i].ID, domain.LikableTypePost).Count(&likeCount)
		posts[i].LikeCount = likeCount
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
