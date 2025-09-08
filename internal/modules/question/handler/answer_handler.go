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

// AnswerHandler handles HTTP requests for answers
type AnswerHandler struct {
	answerService *service.AnswerService
}

// NewAnswerHandler creates a new answer handler instance
func NewAnswerHandler(answerService *service.AnswerService) *AnswerHandler {
	return &AnswerHandler{
		answerService: answerService,
	}
}

// CreateAnswer creates a new answer for a question
// @Summary Create an answer
// @Description Create an answer for a specific question
// @Tags answers
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Param request body dto.CreateAnswerRequest true "Answer details"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /questions/{question_id}/answers [post]
func (h *AnswerHandler) CreateAnswer(c *fiber.Ctx) error {
	questionID := c.Params("question_id")
	if strings.TrimSpace(questionID) == "" {
		return response.ValidationErrorJSON(c, "Invalid question ID", "Question ID cannot be empty")
	}

	var req dto.CreateAnswerRequest
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

	answer, err := h.answerService.CreateAnswer(questionID, req)
	if err != nil {
		logger.Error("Failed to create answer", "question_id", questionID, "error", err)
		if strings.Contains(err.Error(), "not found") {
			return response.NotFoundJSON(c, "Question")
		}
		return response.ValidationErrorJSON(c, "Failed to create answer", err.Error())
	}

	logger.Info("Answer created successfully", "answer_id", answer.ID, "question_id", questionID)
	return response.CreatedJSON(c, answer, "Answer created successfully")
}

// GetAnswersByQuestionID retrieves all answers for a question
// @Summary Get answers for a question
// @Description Retrieve all answers for a specific question
// @Tags answers
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Answers per page" default(20)
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /questions/{question_id}/answers [get]
func (h *AnswerHandler) GetAnswersByQuestionID(c *fiber.Ctx) error {
	questionID := c.Params("question_id")
	if strings.TrimSpace(questionID) == "" {
		return response.ValidationErrorJSON(c, "Invalid question ID", "Question ID cannot be empty")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	answers, err := h.answerService.GetAnswersByQuestionID(questionID, page, limit)
	if err != nil {
		logger.Error("Failed to retrieve answers", "question_id", questionID, "error", err)
		if strings.Contains(err.Error(), "not found") {
			return response.NotFoundJSON(c, "Question")
		}
		return response.InternalErrorJSON(c, "Failed to retrieve answers")
	}

	// Create pagination metadata
	meta := response.NewPaginationMeta(answers.Page, answers.Limit, answers.TotalCount, answers.HasNext)

	return response.SuccessJSONWithMeta(c, answers.Answers, "Answers retrieved successfully", meta)
}

// GetAnswer retrieves a specific answer by ID
// @Summary Get an answer by ID
// @Description Retrieve a specific answer
// @Tags answers
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Param answer_id path string true "Answer ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /questions/{question_id}/answers/{answer_id} [get]
func (h *AnswerHandler) GetAnswer(c *fiber.Ctx) error {
	answerID := c.Params("answer_id")
	if strings.TrimSpace(answerID) == "" {
		return response.ValidationErrorJSON(c, "Invalid answer ID", "Answer ID cannot be empty")
	}

	answer, err := h.answerService.GetAnswerByID(answerID)
	if err != nil {
		return response.NotFoundJSON(c, "Answer")
	}

	return response.SuccessJSON(c, answer, "Answer retrieved successfully")
}

// UpdateAnswer updates an answer
// @Summary Update an answer
// @Description Update an answer (only by author)
// @Tags answers
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Param answer_id path string true "Answer ID"
// @Param request body dto.UpdateAnswerRequest true "Updated answer details"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /questions/{question_id}/answers/{answer_id} [put]
func (h *AnswerHandler) UpdateAnswer(c *fiber.Ctx) error {
	answerID := c.Params("answer_id")
	if strings.TrimSpace(answerID) == "" {
		return response.ValidationErrorJSON(c, "Invalid answer ID", "Answer ID cannot be empty")
	}

	var req dto.UpdateAnswerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ValidationErrorJSON(c, "Invalid request body", err.Error())
	}

	// Extract user ID from JWT claims
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		logger.Error("Failed to extract user ID from context", "error", err)
		return response.UnauthorizedJSON(c)
	}

	answer, err := h.answerService.UpdateAnswer(answerID, req, userID)
	if err != nil {
		logger.Error("Failed to update answer", "answer_id", answerID, "error", err)
		if strings.Contains(err.Error(), "not found") {
			return response.NotFoundJSON(c, "Answer")
		}
		if strings.Contains(err.Error(), "unauthorized") {
			return response.UnauthorizedJSON(c)
		}
		return response.ValidationErrorJSON(c, "Failed to update answer", err.Error())
	}

	return response.SuccessJSON(c, answer, "Answer updated successfully")
}

// DeleteAnswer deletes an answer
// @Summary Delete an answer
// @Description Delete an answer (only by author)
// @Tags answers
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Param answer_id path string true "Answer ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /questions/{question_id}/answers/{answer_id} [delete]
func (h *AnswerHandler) DeleteAnswer(c *fiber.Ctx) error {
	answerID := c.Params("answer_id")
	if strings.TrimSpace(answerID) == "" {
		return response.ValidationErrorJSON(c, "Invalid answer ID", "Answer ID cannot be empty")
	}

	// Extract user ID from JWT claims
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		logger.Error("Failed to extract user ID from context", "error", err)
		return response.UnauthorizedJSON(c)
	}

	if err := h.answerService.DeleteAnswer(answerID, userID); err != nil {
		logger.Error("Failed to delete answer", "answer_id", answerID, "error", err)
		if strings.Contains(err.Error(), "not found") {
			return response.NotFoundJSON(c, "Answer")
		}
		if strings.Contains(err.Error(), "unauthorized") {
			return response.UnauthorizedJSON(c)
		}
		return response.InternalErrorJSON(c, "Failed to delete answer")
	}

	return response.SuccessJSON(c, nil, "Answer deleted successfully")
}

// AcceptAnswer marks an answer as accepted
// @Summary Accept an answer
// @Description Mark an answer as accepted (only by question author)
// @Tags answers
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Param answer_id path string true "Answer ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /questions/{question_id}/answers/{answer_id}/accept [post]
func (h *AnswerHandler) AcceptAnswer(c *fiber.Ctx) error {
	questionID := c.Params("question_id")
	answerID := c.Params("answer_id")

	if strings.TrimSpace(questionID) == "" {
		return response.ValidationErrorJSON(c, "Invalid question ID", "Question ID cannot be empty")
	}

	if strings.TrimSpace(answerID) == "" {
		return response.ValidationErrorJSON(c, "Invalid answer ID", "Answer ID cannot be empty")
	}

	// Extract user ID from JWT claims
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		logger.Error("Failed to extract user ID from context", "error", err)
		return response.UnauthorizedJSON(c)
	}

	if err := h.answerService.AcceptAnswer(answerID, questionID, userID); err != nil {
		logger.Error("Failed to accept answer", "answer_id", answerID, "question_id", questionID, "error", err)
		if strings.Contains(err.Error(), "not found") {
			return response.NotFoundJSON(c, "Question or Answer")
		}
		if strings.Contains(err.Error(), "unauthorized") {
			return response.UnauthorizedJSON(c)
		}
		return response.ValidationErrorJSON(c, "Failed to accept answer", err.Error())
	}

	return response.SuccessJSON(c, nil, "Answer accepted successfully")
}

// UnacceptAnswer marks an answer as not accepted
// @Summary Unaccept an answer
// @Description Mark an answer as not accepted (only by question author)
// @Tags answers
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Param answer_id path string true "Answer ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /questions/{question_id}/answers/{answer_id}/unaccept [post]
func (h *AnswerHandler) UnacceptAnswer(c *fiber.Ctx) error {
	questionID := c.Params("question_id")
	answerID := c.Params("answer_id")

	if strings.TrimSpace(questionID) == "" {
		return response.ValidationErrorJSON(c, "Invalid question ID", "Question ID cannot be empty")
	}

	if strings.TrimSpace(answerID) == "" {
		return response.ValidationErrorJSON(c, "Invalid answer ID", "Answer ID cannot be empty")
	}

	// Extract user ID from JWT claims
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		logger.Error("Failed to extract user ID from context", "error", err)
		return response.UnauthorizedJSON(c)
	}

	if err := h.answerService.UnacceptAnswer(answerID, questionID, userID); err != nil {
		logger.Error("Failed to unaccept answer", "answer_id", answerID, "question_id", questionID, "error", err)
		if strings.Contains(err.Error(), "not found") {
			return response.NotFoundJSON(c, "Question or Answer")
		}
		if strings.Contains(err.Error(), "unauthorized") {
			return response.UnauthorizedJSON(c)
		}
		return response.ValidationErrorJSON(c, "Failed to unaccept answer", err.Error())
	}

	return response.SuccessJSON(c, nil, "Answer unaccepted successfully")
}