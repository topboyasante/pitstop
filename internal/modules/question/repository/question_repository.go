package repository

import (
	"github.com/google/uuid"
	postDomain "github.com/topboyasante/pitstop/internal/modules/post/domain"
	"github.com/topboyasante/pitstop/internal/modules/question/domain"
	"gorm.io/gorm"
)

// QuestionRepository handles question data operations
type QuestionRepository struct {
	db *gorm.DB
}

// NewQuestionRepository creates a new question repository instance
func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// Create creates a new question
func (r *QuestionRepository) Create(question *domain.Question) error {
	if question.ID == "" {
		question.ID = uuid.NewString()
	}
	return r.db.Create(question).Error
}

// GetByID retrieves a question by ID with related data
func (r *QuestionRepository) GetByID(id string) (*domain.Question, error) {
	var question domain.Question
	err := r.db.Preload("User").
		Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, err
	}

	// Calculate answer count
	var answerCount int64
	r.db.Model(&domain.Answer{}).Where("question_id = ?", id).Count(&answerCount)
	question.AnswerCount = answerCount

	// Calculate comment count
	var commentCount int64
	r.db.Model(&postDomain.Comment{}).Where("commentable_id = ? AND commentable_type = ?", id, "question").Count(&commentCount)
	question.CommentCount = commentCount

	// Calculate like count
	var likeCount int64
	r.db.Model(&postDomain.Like{}).Where("likable_id = ? AND likable_type = ?", id, "question").Count(&likeCount)
	question.LikeCount = likeCount

	// Check if question is answered (has at least one accepted answer)
	var acceptedAnswerCount int64
	r.db.Model(&domain.Answer{}).Where("question_id = ? AND is_accepted = ?", id, true).Count(&acceptedAnswerCount)
	question.IsAnswered = acceptedAnswerCount > 0

	return &question, nil
}

// GetAll retrieves all questions with pagination
func (r *QuestionRepository) GetAll(page, limit int) ([]domain.Question, int64, error) {
	var questions []domain.Question
	var totalCount int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&domain.Question{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get questions
	if err := r.db.Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&questions).Error; err != nil {
		return nil, 0, err
	}

	// Calculate counts for each question
	for i := range questions {
		var answerCount int64
		r.db.Model(&domain.Answer{}).Where("question_id = ?", questions[i].ID).Count(&answerCount)
		questions[i].AnswerCount = answerCount

		var commentCount int64
		r.db.Model(&postDomain.Comment{}).Where("commentable_id = ? AND commentable_type = ?", questions[i].ID, "question").Count(&commentCount)
		questions[i].CommentCount = commentCount

		var likeCount int64
		r.db.Model(&postDomain.Like{}).Where("likable_id = ? AND likable_type = ?", questions[i].ID, "question").Count(&likeCount)
		questions[i].LikeCount = likeCount

		// Check if question is answered
		var acceptedAnswerCount int64
		r.db.Model(&domain.Answer{}).Where("question_id = ? AND is_accepted = ?", questions[i].ID, true).Count(&acceptedAnswerCount)
		questions[i].IsAnswered = acceptedAnswerCount > 0
	}

	return questions, totalCount, nil
}

// GetByTag retrieves questions filtered by tag with pagination
func (r *QuestionRepository) GetByTag(tag string, page, limit int) ([]domain.Question, int64, error) {
	var questions []domain.Question
	var totalCount int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&domain.Question{}).Where("tags LIKE ?", "%"+tag+"%").Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get questions
	if err := r.db.Preload("User").
		Where("tags LIKE ?", "%"+tag+"%").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&questions).Error; err != nil {
		return nil, 0, err
	}

	// Calculate counts for each question
	for i := range questions {
		var answerCount int64
		r.db.Model(&domain.Answer{}).Where("question_id = ?", questions[i].ID).Count(&answerCount)
		questions[i].AnswerCount = answerCount

		var commentCount int64
		r.db.Model(&postDomain.Comment{}).Where("commentable_id = ? AND commentable_type = ?", questions[i].ID, "question").Count(&commentCount)
		questions[i].CommentCount = commentCount

		var likeCount int64
		r.db.Model(&postDomain.Like{}).Where("likable_id = ? AND likable_type = ?", questions[i].ID, "question").Count(&likeCount)
		questions[i].LikeCount = likeCount

		// Check if question is answered
		var acceptedAnswerCount int64
		r.db.Model(&domain.Answer{}).Where("question_id = ? AND is_accepted = ?", questions[i].ID, true).Count(&acceptedAnswerCount)
		questions[i].IsAnswered = acceptedAnswerCount > 0
	}

	return questions, totalCount, nil
}

// Update updates a question
func (r *QuestionRepository) Update(question *domain.Question) error {
	return r.db.Save(question).Error
}

// Delete soft deletes a question
func (r *QuestionRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Question{}).Error
}