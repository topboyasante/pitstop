package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/response"
	"github.com/topboyasante/pitstop/internal/modules/user/dto"
	"github.com/topboyasante/pitstop/internal/modules/user/service"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAllUsers retrieves all users
// @Summary Get all users
// @Description Retrieve a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Users per page" default(20)
// @Success 200 {object} response.APIResponse
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	users, err := h.userService.GetAllUsers(page, limit)
	if err != nil {
		logger.Error("Failed to retrieve users", "error", err)
		return response.InternalErrorJSON(c, "Failed to retrieve users")
	}

	// Create pagination metadata
	meta := &response.MetaInfo{
		Page:    users.Page,
		Limit:   users.Limit,
		Total:   users.TotalCount,
		HasNext: users.HasNext,
	}

	return response.SuccessJSONWithMeta(c, users.Users, "Users retrieved successfully", meta)
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User details"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ValidationErrorJSON(c, "Invalid request body", err.Error())
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		logger.Error("Failed to create user", "error", err)
		return response.ValidationErrorJSON(c, "Failed to create user", err.Error())
	}

	logger.Info("User created successfully", "user_id", user.ID)
	return response.CreatedJSON(c, user, "User created successfully")
}

// GetUser retrieves a specific user by ID
// @Summary Get a user by ID
// @Description Retrieve a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return response.NotFoundJSON(c, "User")
	}

	return response.SuccessJSON(c, user, "User retrieved successfully")
}
