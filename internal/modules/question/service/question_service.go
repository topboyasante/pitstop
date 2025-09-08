package service

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/modules/question/domain"
	"github.com/topboyasante/pitstop/internal/modules/question/dto"
	"github.com/topboyasante/pitstop/internal/modules/question/repository"
	"github.com/topboyasante/pitstop/internal/shared/events"
)

// QuestionService handles question business logic
type QuestionService struct {
	questionRepo *repository.QuestionRepository
	validator    *validator.Validate
	eventBus     *events.EventBus
}

// NewQuestionService creates a new question service instance
func NewQuestionService(questionRepo *repository.QuestionRepository, validator *validator.Validate, eventBus *events.EventBus) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
		validator:    validator,
		eventBus:     eventBus,
	}
}

// CreateQuestion creates a new question
func (s *QuestionService) CreateQuestion(req dto.CreateQuestionRequest) (*dto.QuestionResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	question := &domain.Question{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
	}

	if err := s.questionRepo.Create(question); err != nil {
		logger.Error("Failed to create question", "error", err)
		return nil, fmt.Errorf("failed to create question: %w", err)
	}

	logger.Info("Question created successfully", "question_id", question.ID)

	// Parse tags for response
	tags := parseTags(question.Tags)

	return &dto.QuestionResponse{
		ID:           question.ID,
		UserID:       question.UserID,
		Title:        question.Title,
		Content:      question.Content,
		Tags:         tags,
		IsAnswered:   question.IsAnswered,
		CommentCount: question.CommentCount,
		LikeCount:    question.LikeCount,
		AnswerCount:  question.AnswerCount,
		CreatedAt:    question.CreatedAt,
		UpdatedAt:    question.UpdatedAt,
	}, nil
}

// GetQuestionByID retrieves a question by ID
func (s *QuestionService) GetQuestionByID(id string) (*dto.QuestionResponse, error) {
	question, err := s.questionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	// Parse tags for response
	tags := parseTags(question.Tags)

	response := &dto.QuestionResponse{
		ID:           question.ID,
		UserID:       question.UserID,
		Title:        question.Title,
		Content:      question.Content,
		Tags:         tags,
		IsAnswered:   question.IsAnswered,
		CommentCount: question.CommentCount,
		LikeCount:    question.LikeCount,
		AnswerCount:  question.AnswerCount,
		CreatedAt:    question.CreatedAt,
		UpdatedAt:    question.UpdatedAt,
	}

	if question.User != nil {
		response.User = &dto.QuestionUserResponse{
			Username:    question.User.Username,
			DisplayName: question.User.DisplayName,
			AvatarURL:   question.User.AvatarURL,
		}
	}

	return response, nil
}

// GetAllQuestions retrieves all questions with pagination
func (s *QuestionService) GetAllQuestions(page, limit int) (*dto.QuestionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	questions, totalCount, err := s.questionRepo.GetAll(page, limit)
	if err != nil {
		logger.Error("Failed to retrieve questions", "error", err)
		return nil, fmt.Errorf("failed to retrieve questions: %w", err)
	}

	// Convert to response format
	questionResponses := make([]dto.QuestionResponse, len(questions))
	for i, question := range questions {
		tags := parseTags(question.Tags)

		questionResponses[i] = dto.QuestionResponse{
			ID:           question.ID,
			UserID:       question.UserID,
			Title:        question.Title,
			Content:      question.Content,
			Tags:         tags,
			IsAnswered:   question.IsAnswered,
			CommentCount: question.CommentCount,
			LikeCount:    question.LikeCount,
			AnswerCount:  question.AnswerCount,
			CreatedAt:    question.CreatedAt,
			UpdatedAt:    question.UpdatedAt,
		}

		if question.User != nil {
			questionResponses[i].User = &dto.QuestionUserResponse{
				Username:    question.User.Username,
				DisplayName: question.User.DisplayName,
				AvatarURL:   question.User.AvatarURL,
			}
		}
	}

	hasNext := int64(page*limit) < totalCount

	return &dto.QuestionsResponse{
		Questions:  questionResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasNext:    hasNext,
	}, nil
}

// GetQuestionsByTag retrieves questions filtered by tag with pagination
func (s *QuestionService) GetQuestionsByTag(tag string, page, limit int) (*dto.QuestionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	questions, totalCount, err := s.questionRepo.GetByTag(tag, page, limit)
	if err != nil {
		logger.Error("Failed to retrieve questions by tag", "tag", tag, "error", err)
		return nil, fmt.Errorf("failed to retrieve questions by tag: %w", err)
	}

	// Convert to response format
	questionResponses := make([]dto.QuestionResponse, len(questions))
	for i, question := range questions {
		tags := parseTags(question.Tags)

		questionResponses[i] = dto.QuestionResponse{
			ID:           question.ID,
			UserID:       question.UserID,
			Title:        question.Title,
			Content:      question.Content,
			Tags:         tags,
			IsAnswered:   question.IsAnswered,
			CommentCount: question.CommentCount,
			LikeCount:    question.LikeCount,
			AnswerCount:  question.AnswerCount,
			CreatedAt:    question.CreatedAt,
			UpdatedAt:    question.UpdatedAt,
		}

		if question.User != nil {
			questionResponses[i].User = &dto.QuestionUserResponse{
				Username:    question.User.Username,
				DisplayName: question.User.DisplayName,
				AvatarURL:   question.User.AvatarURL,
			}
		}
	}

	hasNext := int64(page*limit) < totalCount

	return &dto.QuestionsResponse{
		Questions:  questionResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasNext:    hasNext,
	}, nil
}

// UpdateQuestion updates a question
func (s *QuestionService) UpdateQuestion(id string, req dto.UpdateQuestionRequest) (*dto.QuestionResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	question, err := s.questionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	// Update fields if provided
	if req.Title != "" {
		question.Title = req.Title
	}
	if req.Content != "" {
		question.Content = req.Content
	}
	if req.Tags != "" {
		question.Tags = req.Tags
	}

	if err := s.questionRepo.Update(question); err != nil {
		logger.Error("Failed to update question", "question_id", id, "error", err)
		return nil, fmt.Errorf("failed to update question: %w", err)
	}

	logger.Info("Question updated successfully", "question_id", id)

	// Parse tags for response
	tags := parseTags(question.Tags)

	return &dto.QuestionResponse{
		ID:           question.ID,
		UserID:       question.UserID,
		Title:        question.Title,
		Content:      question.Content,
		Tags:         tags,
		IsAnswered:   question.IsAnswered,
		CommentCount: question.CommentCount,
		LikeCount:    question.LikeCount,
		AnswerCount:  question.AnswerCount,
		CreatedAt:    question.CreatedAt,
		UpdatedAt:    question.UpdatedAt,
	}, nil
}

// DeleteQuestion deletes a question
func (s *QuestionService) DeleteQuestion(id string) error {
	if err := s.questionRepo.Delete(id); err != nil {
		logger.Error("Failed to delete question", "question_id", id, "error", err)
		return fmt.Errorf("failed to delete question: %w", err)
	}

	logger.Info("Question deleted successfully", "question_id", id)
	return nil
}

// parseTags splits comma-separated tags into a slice
func parseTags(tags string) []string {
	if tags == "" {
		return []string{}
	}
	
	tagList := strings.Split(tags, ",")
	result := make([]string, 0, len(tagList))
	
	for _, tag := range tagList {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}