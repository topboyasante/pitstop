package repository

import (
	"github.com/google/uuid"
	postDomain "github.com/topboyasante/pitstop/internal/modules/post/domain"
	"github.com/topboyasante/pitstop/internal/modules/question/domain"
	"gorm.io/gorm"
)

// AnswerRepository handles answer data operations
type AnswerRepository struct {
	db *gorm.DB
}

// NewAnswerRepository creates a new answer repository instance
func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

// Create creates a new answer
func (r *AnswerRepository) Create(answer *domain.Answer) error {
	if answer.ID == "" {
		answer.ID = uuid.NewString()
	}
	return r.db.Create(answer).Error
}

// GetByID retrieves an answer by ID with related data
func (r *AnswerRepository) GetByID(id string) (*domain.Answer, error) {
	var answer domain.Answer
	err := r.db.Preload("User").Preload("Question").
		Where("id = ?", id).First(&answer).Error
	if err != nil {
		return nil, err
	}

	// Calculate like count
	var likeCount int64
	r.db.Model(&postDomain.Like{}).Where("likable_id = ? AND likable_type = ?", id, "answer").Count(&likeCount)
	answer.LikeCount = likeCount

	return &answer, nil
}

// GetByQuestionID retrieves all answers for a question with pagination
func (r *AnswerRepository) GetByQuestionID(questionID string, page, limit int) ([]domain.Answer, int64, error) {
	var answers []domain.Answer
	var totalCount int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&domain.Answer{}).Where("question_id = ?", questionID).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get answers - accepted answers first, then ordered by creation date
	if err := r.db.Preload("User").
		Where("question_id = ?", questionID).
		Offset(offset).
		Limit(limit).
		Order("is_accepted DESC, created_at ASC").
		Find(&answers).Error; err != nil {
		return nil, 0, err
	}

	// Calculate like counts
	for i := range answers {
		var likeCount int64
		r.db.Model(&postDomain.Like{}).Where("likable_id = ? AND likable_type = ?", answers[i].ID, "answer").Count(&likeCount)
		answers[i].LikeCount = likeCount
	}

	return answers, totalCount, nil
}

// AcceptAnswer marks an answer as accepted and unmarks other accepted answers for the same question
func (r *AnswerRepository) AcceptAnswer(answerID, questionID string) error {
	// Start transaction
	tx := r.db.Begin()

	// First, unmark all other accepted answers for this question
	if err := tx.Model(&domain.Answer{}).
		Where("question_id = ? AND id != ?", questionID, answerID).
		Update("is_accepted", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Mark the specified answer as accepted
	if err := tx.Model(&domain.Answer{}).
		Where("id = ?", answerID).
		Update("is_accepted", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update question's is_answered status
	if err := tx.Model(&domain.Question{}).
		Where("id = ?", questionID).
		Update("is_answered", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UnacceptAnswer marks an answer as not accepted
func (r *AnswerRepository) UnacceptAnswer(answerID, questionID string) error {
	// Start transaction
	tx := r.db.Begin()

	// Unmark the answer as accepted
	if err := tx.Model(&domain.Answer{}).
		Where("id = ?", answerID).
		Update("is_accepted", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Check if there are any other accepted answers for this question
	var acceptedCount int64
	if err := tx.Model(&domain.Answer{}).
		Where("question_id = ? AND is_accepted = ?", questionID, true).
		Count(&acceptedCount).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update question's is_answered status based on remaining accepted answers
	isAnswered := acceptedCount > 0
	if err := tx.Model(&domain.Question{}).
		Where("id = ?", questionID).
		Update("is_answered", isAnswered).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Update updates an answer
func (r *AnswerRepository) Update(answer *domain.Answer) error {
	return r.db.Save(answer).Error
}

// Delete soft deletes an answer
func (r *AnswerRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Answer{}).Error
}