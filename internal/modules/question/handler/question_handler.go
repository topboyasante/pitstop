package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/response"
	"github.com/topboyasante/pitstop/internal/modules/question/dto"
	"github.com/topboyasante/pitstop/internal/modules/question/service"
	utils "github.com/topboyasante/pitstop/internal/shared/utils"
)

// QuestionHandler handles HTTP requests for questions
type QuestionHandler struct {
	questionService *service.QuestionService
}

// NewQuestionHandler creates a new question handler instance
func NewQuestionHandler(questionService *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

// GetAllQuestions retrieves all questions
// @Summary Get all questions
// @Description Retrieve a paginated list of questions
// @Tags questions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Questions per page" default(20)
// @Success 200 {object} response.APIResponse
// @Router /questions [get]
func (h *QuestionHandler) GetAllQuestions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	questions, err := h.questionService.GetAllQuestions(page, limit)
	if err != nil {
		logger.Error("Failed to retrieve questions", "error", err)
		return response.InternalErrorJSON(c, "Failed to retrieve questions")
	}

	// Create pagination metadata
	meta := response.NewPaginationMeta(questions.Page, questions.Limit, questions.TotalCount, questions.HasNext)

	return response.SuccessJSONWithMeta(c, questions.Questions, "Questions retrieved successfully", meta)
}

// CreateQuestion creates a new question
// @Summary Create a new question
// @Description Create a new question
// @Tags questions
// @Accept json
// @Produce json
// @Param request body dto.CreateQuestionRequest true "Question details"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Security BearerAuth
// @Router /questions [post]
func (h *QuestionHandler) CreateQuestion(c *fiber.Ctx) error {
	var req dto.CreateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ValidationErrorJSON(c, "Invalid request body", err.Error())
	}

	// Extract user ID from JWT claims
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		logger.Error("Failed to extract user ID from context", "error", err)
		return response.UnauthorizedJSON(c)
	}

	req.UserID = userID

	question, err := h.questionService.CreateQuestion(req)
	if err != nil {
		logger.Error("Failed to create question", "error", err)
		return response.ValidationErrorJSON(c, "Failed to create question", err.Error())
	}

	logger.Info("Question created successfully", "question_id", question.ID)
	return response.CreatedJSON(c, question, "Question created successfully")
}

// GetQuestion retrieves a specific question by ID
// @Summary Get a question by ID
// @Description Retrieve a specific question
// @Tags questions
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /questions/{id} [get]
func (h *QuestionHandler) GetQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		return response.ValidationErrorJSON(c, "Invalid question ID", "ID cannot be empty")
	}

	question, err := h.questionService.GetQuestionByID(id)
	if err != nil {
		return response.NotFoundJSON(c, "Question")
	}

	return response.SuccessJSON(c, question, "Question retrieved successfully")
}

// UpdateQuestion updates a question
// @Summary Update a question
// @Description Update a question (only by author)
// @Tags questions
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Param request body dto.UpdateQuestionRequest true "Updated question details"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /questions/{id} [put]
func (h *QuestionHandler) UpdateQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		return response.ValidationErrorJSON(c, "Invalid question ID", "ID cannot be empty")
	}

	var req dto.UpdateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ValidationErrorJSON(c, "Invalid request body", err.Error())
	}

	question, err := h.questionService.UpdateQuestion(id, req)
	if err != nil {
		logger.Error("Failed to update question", "question_id", id, "error", err)
		if strings.Contains(err.Error(), "not found") {
			return response.NotFoundJSON(c, "Question")
		}
		return response.ValidationErrorJSON(c, "Failed to update question", err.Error())
	}

	return response.SuccessJSON(c, question, "Question updated successfully")
}

// DeleteQuestion deletes a question
// @Summary Delete a question
// @Description Delete a question (only by author)
// @Tags questions
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /questions/{id} [delete]
func (h *QuestionHandler) DeleteQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		return response.ValidationErrorJSON(c, "Invalid question ID", "ID cannot be empty")
	}

	if err := h.questionService.DeleteQuestion(id); err != nil {
		logger.Error("Failed to delete question", "question_id", id, "error", err)
		if strings.Contains(err.Error(), "not found") {
			return response.NotFoundJSON(c, "Question")
		}
		return response.InternalErrorJSON(c, "Failed to delete question")
	}

	return response.SuccessJSON(c, nil, "Question deleted successfully")
}

// GetQuestionsByTag retrieves questions filtered by tag
// @Summary Get questions by tag
// @Description Retrieve questions filtered by a specific tag
// @Tags questions
// @Accept json
// @Produce json
// @Param tag query string true "Tag to filter by"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Questions per page" default(20)
// @Success 200 {object} response.APIResponse
// @Router /questions/tag [get]
func (h *QuestionHandler) GetQuestionsByTag(c *fiber.Ctx) error {
	tag := c.Query("tag")
	if strings.TrimSpace(tag) == "" {
		return response.ValidationErrorJSON(c, "Tag parameter is required", "Tag cannot be empty")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	questions, err := h.questionService.GetQuestionsByTag(tag, page, limit)
	if err != nil {
		logger.Error("Failed to retrieve questions by tag", "tag", tag, "error", err)
		return response.InternalErrorJSON(c, "Failed to retrieve questions by tag")
	}

	// Create pagination metadata
	meta := response.NewPaginationMeta(questions.Page, questions.Limit, questions.TotalCount, questions.HasNext)

	return response.SuccessJSONWithMeta(c, questions.Questions, "Questions retrieved successfully", meta)
}