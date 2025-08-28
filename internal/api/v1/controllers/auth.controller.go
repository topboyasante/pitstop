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

func (ac *AuthController) Authenticate(c *fiber.Ctx) error {
	url := ac.authService.Authenticate()
	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return c.JSON(url)
}
