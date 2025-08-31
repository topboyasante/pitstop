package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/modules/auth/handler"
)

// RegisterRoutes registers all auth-related routes
func RegisterRoutes(router fiber.Router, authHandler *handler.AuthHandler) {
	auth := router.Group("/auth")

	// OAuth routes
	auth.Get("/google", authHandler.GoogleAuth)
	auth.Get("/google/callback", authHandler.GoogleCallback)
}
