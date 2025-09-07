package repository

import (
	"github.com/google/uuid"
	"github.com/topboyasante/pitstop/internal/modules/post/domain"
	"gorm.io/gorm"
)

// CommentRepository handles comment data operations
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new comment repository instance
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create creates a new comment
func (r *CommentRepository) Create(comment *domain.Comment) error {
	if comment.ID == "" {
		comment.ID = uuid.NewString()
	}
	return r.db.Create(comment).Error
}

// GetByID retrieves a comment by ID
func (r *CommentRepository) GetByID(id string) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.Preload("User").Preload("Parent").Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	
	// Calculate like count
	var likeCount int64
	r.db.Model(&domain.Like{}).Where("likable_id = ? AND likable_type = ?", id, domain.LikableTypeComment).Count(&likeCount)
	comment.LikeCount = likeCount
	
	return &comment, nil
}

// GetByPostID retrieves all comments for a specific post with nested replies
func (r *CommentRepository) GetByPostID(postID string) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at DESC").
		Find(&comments).Error
		
	if err != nil {
		return nil, err
	}
	
	// Calculate like counts for all comments and replies
	for i := range comments {
		var likeCount int64
		r.db.Model(&domain.Like{}).Where("likable_id = ? AND likable_type = ?", comments[i].ID, domain.LikableTypeComment).Count(&likeCount)
		comments[i].LikeCount = likeCount
		
		// Calculate like counts for replies
		for j := range comments[i].Replies {
			var replyLikeCount int64
			r.db.Model(&domain.Like{}).Where("likable_id = ? AND likable_type = ?", comments[i].Replies[j].ID, domain.LikableTypeComment).Count(&replyLikeCount)
			comments[i].Replies[j].LikeCount = replyLikeCount
		}
	}
	
	return comments, nil
}

// GetRepliesByParentID retrieves all replies for a specific comment
func (r *CommentRepository) GetRepliesByParentID(parentID string) ([]domain.Comment, error) {
	var replies []domain.Comment
	err := r.db.Preload("User").
		Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Find(&replies).Error
		
	if err != nil {
		return nil, err
	}
	
	// Calculate like counts for all replies
	for i := range replies {
		var likeCount int64
		r.db.Model(&domain.Like{}).Where("likable_id = ? AND likable_type = ?", replies[i].ID, domain.LikableTypeComment).Count(&likeCount)
		replies[i].LikeCount = likeCount
	}
	
	return replies, nil
}

// Update updates a comment
func (r *CommentRepository) Update(comment *domain.Comment) error {
	return r.db.Save(comment).Error
}

// Delete soft deletes a comment
func (r *CommentRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Comment{}).Error
}

// GetCommentsByUserID retrieves all comments by a specific user
func (r *CommentRepository) GetCommentsByUserID(userID string, page, limit int) ([]domain.Comment, int64, error) {
	var comments []domain.Comment
	var totalCount int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&domain.Comment{}).Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get comments
	if err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, totalCount, nil
}