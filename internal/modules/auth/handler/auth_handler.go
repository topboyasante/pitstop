package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/response"
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
// @Success 200 {object} response.APIResponse
// @Router /auth/google [get]
func (h *AuthHandler) GoogleAuth(c *fiber.Ctx) error {
	authURL := h.authService.Authenticate()

	return response.SuccessJSON(c, dto.AuthURLResponse{
		AuthURL: authURL,
	}, "Google OAuth URL generated successfully")
}

// GoogleCallback handles Google OAuth callback
// @Summary Handle Google OAuth callback
// @Description Redirects to frontend with authorization code
// @Tags auth
// @Param code query string true "Authorization code"
// @Param state query string true "CSRF state token"
// @Success 302 "Redirect to frontend"
// @Failure 400 {object} map[string]string
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		// Redirect to frontend with error
		return c.Redirect(config.Get().Server.FrontendURL + "/auth/error?error=missing_parameters")
	}

	// Redirect to frontend with code and state. The frontend will then hit a /exchange endpoint to retrieve the auth tokens
	redirectURL := fmt.Sprintf("%s/auth/callback?code=%s&state=%s", config.Get().Server.FrontendURL, code, state)
	return c.Redirect(redirectURL)
}

// RefreshToken handles token refresh requests
// @Summary Refresh JWT tokens
// @Description Generate new JWT tokens using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ValidationErrorJSON(c, "Invalid request body", err.Error())
	}

	// Simple validation - refresh token is required
	if req.RefreshToken == "" {
		return response.ValidationErrorJSON(c, "Refresh token is required", "refresh_token field cannot be empty")
	}

	tokens, err := h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		logger.Error("Token refresh failed", "error", err)
		return response.ErrorJSON(c, fiber.StatusUnauthorized, "TOKEN_REFRESH_FAILED", "Failed to refresh tokens", err.Error())
	}

	return response.SuccessJSON(c, tokens, "Tokens refreshed successfully")
}

// ExchangeCode handles authorization code to token exchange
// @Summary Exchange authorization code for JWT tokens
// @Description Exchange OAuth authorization code for JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ExchangeCodeRequest true "Code exchange request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /auth/exchange [post]
func (h *AuthHandler) ExchangeCode(c *fiber.Ctx) error {
	var req dto.ExchangeCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ValidationErrorJSON(c, "Invalid request body", err.Error())
	}

	if req.Code == "" || req.State == "" {
		return response.ValidationErrorJSON(c, "Code and state are required", "Both 'code' and 'state' fields must be provided")
	}

	tokens, err := h.authService.ExchangeCode(req.Code, req.State)
	if err != nil {
		logger.Error("Code exchange failed", "error", err)
		return response.ValidationErrorJSON(c, "Failed to exchange authorization code", err.Error())
	}

	return response.SuccessJSON(c, tokens, "Authorization code exchanged successfully")
}

// Me returns the current authenticated user info from JWT token
// @Summary Get current user
// @Description Get current authenticated user information from JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// Get user from database via user service
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		return response.NotFoundJSON(c, "User")
	}

	return response.SuccessJSON(c, user, "User information retrieved successfully")
}
