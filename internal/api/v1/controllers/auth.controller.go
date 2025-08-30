package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/api/v1/services"
	"github.com/topboyasante/pitstop/internal/logger"
	"github.com/topboyasante/pitstop/internal/utils"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// @Summary Start OAuth authentication
// @Description Redirects user to Google OAuth for authentication
// @Tags Auth
// @Success 303 "Redirects to Google OAuth"
// @Header 303 {string} X-Request-ID "Request ID for tracing"
// @Router /auth [get]
func (ac *AuthController) Authenticate(c *fiber.Ctx) error {
	requestID := utils.GenerateRequestID()
	logger.Info("Authentication endpoint called",
		"event", "auth.endpoint_accessed",
		"request_id", requestID,
		"path", "/auth",
		"method", "GET",
		"ip", c.IP(),
		"user_agent", c.Get("User-Agent"))

	url := ac.authService.Authenticate()

	logger.Info("Authentication redirect sent",
		"event", "auth.redirect_sent",
		"request_id", requestID,
		"redirect_url", url,
		"status_code", fiber.StatusSeeOther)

	c.Set("X-Request-ID", requestID)
	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)

	return nil
}

// @Summary OAuth callback handler
// @Description Handles OAuth callback and exchanges authorization code for access token
// @Tags Auth
// @Param code query string true "Authorization code from OAuth provider"
// @Param state query string true "State parameter for CSRF protection"
// @Success 200 {object} models.AuthSuccessResponse "Authentication successful with access token"
// @Failure 400 {object} models.BadRequestErrorResponse "Bad request - missing code or invalid state"
// @Failure 503 {object} models.ErrorResponse "External service error - OAuth provider failure"
// @Header 200,400,503 {string} X-Request-ID "Request ID for tracing"
// @Router /auth/callback [get]
func (ac *AuthController) Callback(c *fiber.Ctx) error {
	requestID := utils.GenerateRequestID()
	code := c.Query("code")
	state := c.Query("state")

	logger.Info("OAuth callback received",
		"event", "auth.callback_received",
		"request_id", requestID,
		"path", "/auth/callback",
		"method", "GET",
		"ip", c.IP(),
		"has_code", code != "")

	if code == "" {
		logger.Error("OAuth callback failed",
			"event", "auth.callback_failed",
			"request_id", requestID,
			"reason", "missing_authorization_code",
			"ip", c.IP())
		return utils.SendBadRequestError(c, "Authorization code is required", requestID)
	}

	token, err := ac.authService.ExchangeCode(code, state)
	if err != nil {
		logger.Error("OAuth callback failed",
			"event", "auth.callback_failed",
			"request_id", requestID,
			"reason", "token_exchange_error",
			"error", err.Error(),
			"ip", c.IP())
		return utils.SendExternalServiceError(c, "Failed to exchange authorization code", requestID, map[string]any{
			"service": "google_oauth",
		})
	}

	logger.Info("OAuth callback successful",
		"event", "auth.callback_successful",
		"request_id", requestID,
		"ip", c.IP())

	// For now, just return success with token
	// In production, create JWT and redirect to frontend
	return utils.SendSuccessResponse(c, map[string]any{
		"message": "Authentication successful",
		"token":   token,
	}, requestID)
}
