package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/modules/auth/dto"
	"github.com/topboyasante/pitstop/internal/modules/auth/service"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// GoogleAuth initiates Google OAuth authentication
// @Summary Initiate Google OAuth
// @Description Get Google OAuth authorization URL
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.AuthURLResponse
// @Router /auth/google [get]
func (h *AuthHandler) GoogleAuth(c *fiber.Ctx) error {
	authURL := h.authService.Authenticate()

	return c.JSON(dto.AuthURLResponse{
		AuthURL: authURL,
	})
}

// GoogleCallback handles Google OAuth callback
// @Summary Handle Google OAuth callback
// @Description Exchange authorization code for access token
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "CSRF state token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing code or state parameter",
		})
	}

	token, err := h.authService.ExchangeCode(code, state)
	if err != nil {
		logger.Error("OAuth callback failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to exchange authorization code",
		})
	}

	return c.JSON(fiber.Map{
		"access_token": token,
	})
}
