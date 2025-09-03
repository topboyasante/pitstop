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
// @Description Returns authorization code for frontend to exchange
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "CSRF state token"
// @Success 200 {object} dto.AuthCodeResponse
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

	// Just return the code and state - let frontend exchange for tokens
	return c.JSON(dto.AuthCodeResponse{
		Code:  code,
		State: state,
	})
}

// RefreshToken handles token refresh requests
// @Summary Refresh JWT tokens
// @Description Generate new JWT tokens using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} dto.JWTTokenResponse
// @Failure 400 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Simple validation - refresh token is required
	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	tokens, err := h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		logger.Error("Token refresh failed", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to refresh tokens",
		})
	}

	return c.JSON(tokens)
}

// ExchangeCode handles authorization code to token exchange
// @Summary Exchange authorization code for JWT tokens
// @Description Exchange OAuth authorization code for JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ExchangeCodeRequest true "Code exchange request"
// @Success 200 {object} dto.JWTTokenResponse
// @Failure 400 {object} map[string]string
// @Router /auth/exchange [post]
func (h *AuthHandler) ExchangeCode(c *fiber.Ctx) error {
	var req dto.ExchangeCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Code == "" || req.State == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code and state are required",
		})
	}

	tokens, err := h.authService.ExchangeCode(req.Code, req.State)
	if err != nil {
		logger.Error("Code exchange failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to exchange authorization code",
		})
	}

	return c.JSON(tokens)
}

// Me returns the current authenticated user info from JWT token
// @Summary Get current user
// @Description Get current authenticated user information from JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	audience := c.Locals("audience").(string)

	return c.JSON(fiber.Map{
		"user_id":  userID,
		"audience": audience,
		"message":  "Authentication successful",
	})
}
