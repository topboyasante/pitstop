package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/api/v1/services"
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
// @Success 303 {object} string "Redirects to Google OAuth"
// @Router /auth [get]
func (ac *AuthController) Authenticate(c *fiber.Ctx) error {
	url := ac.authService.Authenticate()
	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return nil
}

// @Summary OAuth callback handler
// @Description Handles OAuth callback and exchanges authorization code for access token
// @Tags Auth
// @Param code query string true "Authorization code from OAuth provider"
// @Param state query string true "State parameter for CSRF protection"
// @Success 200 {object} map[string]interface{} "Authentication successful with access token"
// @Failure 400 {object} map[string]interface{} "Bad request - missing code or invalid state"
// @Router /auth/callback [get]
func (ac *AuthController) Callback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "authorization code not found",
		})
	}

	token, err := ac.authService.ExchangeCode(code, state)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// For now, just return success with token
	// In production, create JWT and redirect to frontend
	return c.JSON(fiber.Map{
		"message": "Authentication successful",
		"token":   token,
	})
}
