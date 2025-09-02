package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/logger"
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
// @Success 200 {object} github.com/topboyasante/pitstop/internal/modules/user/dto.UsersResponse
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	users, err := h.userService.GetAllUsers(page, limit)
	if err != nil {
		logger.Error("Failed to retrieve users", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	return c.JSON(users)
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body github.com/topboyasante/pitstop/internal/modules/user/dto.CreateUserRequest true "User details"
// @Success 201 {object} github.com/topboyasante/pitstop/internal/modules/user/dto.UserResponse
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		logger.Error("Failed to create user", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Info("User created successfully", "user_id", user.ID)
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUser retrieves a specific user by ID
// @Summary Get a user by ID
// @Description Retrieve a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} github.com/topboyasante/pitstop/internal/modules/user/dto.UserResponse
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}
