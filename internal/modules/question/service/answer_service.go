package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/modules/question/domain"
	"github.com/topboyasante/pitstop/internal/modules/question/dto"
	"github.com/topboyasante/pitstop/internal/modules/question/repository"
	"github.com/topboyasante/pitstop/internal/shared/events"
)

// AnswerService handles answer business logic
type AnswerService struct {
	answerRepo   *repository.AnswerRepository
	questionRepo *repository.QuestionRepository
	validator    *validator.Validate
	eventBus     *events.EventBus
}

// NewAnswerService creates a new answer service instance
func NewAnswerService(answerRepo *repository.AnswerRepository, questionRepo *repository.QuestionRepository, validator *validator.Validate, eventBus *events.EventBus) *AnswerService {
	return &AnswerService{
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
		validator:    validator,
		eventBus:     eventBus,
	}
}

// CreateAnswer creates a new answer
func (s *AnswerService) CreateAnswer(questionID string, req dto.CreateAnswerRequest) (*dto.AnswerResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify question exists
	_, err := s.questionRepo.GetByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	answer := &domain.Answer{
		QuestionID: questionID,
		UserID:     req.UserID,
		Content:    req.Content,
	}

	if err := s.answerRepo.Create(answer); err != nil {
		logger.Error("Failed to create answer", "error", err)
		return nil, fmt.Errorf("failed to create answer: %w", err)
	}

	logger.Info("Answer created successfully", "answer_id", answer.ID, "question_id", questionID)

	return &dto.AnswerResponse{
		ID:         answer.ID,
		QuestionID: answer.QuestionID,
		UserID:     answer.UserID,
		Content:    answer.Content,
		IsAccepted: answer.IsAccepted,
		LikeCount:  answer.LikeCount,
		CreatedAt:  answer.CreatedAt,
		UpdatedAt:  answer.UpdatedAt,
	}, nil
}

// GetAnswerByID retrieves an answer by ID
func (s *AnswerService) GetAnswerByID(id string) (*dto.AnswerResponse, error) {
	answer, err := s.answerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("answer not found: %w", err)
	}

	response := &dto.AnswerResponse{
		ID:         answer.ID,
		QuestionID: answer.QuestionID,
		UserID:     answer.UserID,
		Content:    answer.Content,
		IsAccepted: answer.IsAccepted,
		LikeCount:  answer.LikeCount,
		CreatedAt:  answer.CreatedAt,
		UpdatedAt:  answer.UpdatedAt,
	}

	if answer.User != nil {
		response.User = &dto.QuestionUserResponse{
			Username:    answer.User.Username,
			DisplayName: answer.User.DisplayName,
			AvatarURL:   answer.User.AvatarURL,
		}
	}

	return response, nil
}

// GetAnswersByQuestionID retrieves all answers for a question with pagination
func (s *AnswerService) GetAnswersByQuestionID(questionID string, page, limit int) (*dto.AnswersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Verify question exists
	_, err := s.questionRepo.GetByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	answers, totalCount, err := s.answerRepo.GetByQuestionID(questionID, page, limit)
	if err != nil {
		logger.Error("Failed to retrieve answers", "question_id", questionID, "error", err)
		return nil, fmt.Errorf("failed to retrieve answers: %w", err)
	}

	// Convert to response format
	answerResponses := make([]dto.AnswerResponse, len(answers))
	for i, answer := range answers {
		answerResponses[i] = dto.AnswerResponse{
			ID:         answer.ID,
			QuestionID: answer.QuestionID,
			UserID:     answer.UserID,
			Content:    answer.Content,
			IsAccepted: answer.IsAccepted,
			LikeCount:  answer.LikeCount,
			CreatedAt:  answer.CreatedAt,
			UpdatedAt:  answer.UpdatedAt,
		}

		if answer.User != nil {
			answerResponses[i].User = &dto.QuestionUserResponse{
				Username:    answer.User.Username,
				DisplayName: answer.User.DisplayName,
				AvatarURL:   answer.User.AvatarURL,
			}
		}
	}

	hasNext := int64(page*limit) < totalCount

	return &dto.AnswersResponse{
		Answers:    answerResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasNext:    hasNext,
	}, nil
}

// AcceptAnswer marks an answer as accepted
func (s *AnswerService) AcceptAnswer(answerID, questionID, userID string) error {
	// Verify the question exists and belongs to the user
	question, err := s.questionRepo.GetByID(questionID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	if question.UserID != userID {
		return fmt.Errorf("unauthorized: only the question author can accept answers")
	}

	// Verify the answer exists and belongs to the question
	answer, err := s.answerRepo.GetByID(answerID)
	if err != nil {
		return fmt.Errorf("answer not found: %w", err)
	}

	if answer.QuestionID != questionID {
		return fmt.Errorf("answer does not belong to this question")
	}

	if err := s.answerRepo.AcceptAnswer(answerID, questionID); err != nil {
		logger.Error("Failed to accept answer", "answer_id", answerID, "question_id", questionID, "error", err)
		return fmt.Errorf("failed to accept answer: %w", err)
	}

	logger.Info("Answer accepted successfully", "answer_id", answerID, "question_id", questionID)
	return nil
}

// UnacceptAnswer marks an answer as not accepted
func (s *AnswerService) UnacceptAnswer(answerID, questionID, userID string) error {
	// Verify the question exists and belongs to the user
	question, err := s.questionRepo.GetByID(questionID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	if question.UserID != userID {
		return fmt.Errorf("unauthorized: only the question author can unaccept answers")
	}

	// Verify the answer exists and belongs to the question
	answer, err := s.answerRepo.GetByID(answerID)
	if err != nil {
		return fmt.Errorf("answer not found: %w", err)
	}

	if answer.QuestionID != questionID {
		return fmt.Errorf("answer does not belong to this question")
	}

	if err := s.answerRepo.UnacceptAnswer(answerID, questionID); err != nil {
		logger.Error("Failed to unaccept answer", "answer_id", answerID, "question_id", questionID, "error", err)
		return fmt.Errorf("failed to unaccept answer: %w", err)
	}

	logger.Info("Answer unaccepted successfully", "answer_id", answerID, "question_id", questionID)
	return nil
}

// UpdateAnswer updates an answer
func (s *AnswerService) UpdateAnswer(id string, req dto.UpdateAnswerRequest, userID string) (*dto.AnswerResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	answer, err := s.answerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("answer not found: %w", err)
	}

	// Verify the answer belongs to the user
	if answer.UserID != userID {
		return nil, fmt.Errorf("unauthorized: only the answer author can update this answer")
	}

	// Update fields if provided
	if req.Content != "" {
		answer.Content = req.Content
	}

	if err := s.answerRepo.Update(answer); err != nil {
		logger.Error("Failed to update answer", "answer_id", id, "error", err)
		return nil, fmt.Errorf("failed to update answer: %w", err)
	}

	logger.Info("Answer updated successfully", "answer_id", id)

	return &dto.AnswerResponse{
		ID:         answer.ID,
		QuestionID: answer.QuestionID,
		UserID:     answer.UserID,
		Content:    answer.Content,
		IsAccepted: answer.IsAccepted,
		LikeCount:  answer.LikeCount,
		CreatedAt:  answer.CreatedAt,
		UpdatedAt:  answer.UpdatedAt,
	}, nil
}

// DeleteAnswer deletes an answer
func (s *AnswerService) DeleteAnswer(id string, userID string) error {
	answer, err := s.answerRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("answer not found: %w", err)
	}

	// Verify the answer belongs to the user
	if answer.UserID != userID {
		return fmt.Errorf("unauthorized: only the answer author can delete this answer")
	}

	if err := s.answerRepo.Delete(id); err != nil {
		logger.Error("Failed to delete answer", "answer_id", id, "error", err)
		return fmt.Errorf("failed to delete answer: %w", err)
	}

	logger.Info("Answer deleted successfully", "answer_id", id)
	return nil
}