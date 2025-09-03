package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/middleware"
	"github.com/topboyasante/pitstop/internal/modules/auth/handler"
)

// RegisterRoutes registers all auth-related routes
func RegisterRoutes(router fiber.Router, authHandler *handler.AuthHandler, config *config.Config) {
	auth := router.Group("/auth")

	// OAuth routes
	auth.Get("/google", authHandler.GoogleAuth)
	auth.Get("/google/callback", authHandler.GoogleCallback)
	
	// JWT token routes
	auth.Post("/refresh", authHandler.RefreshToken)
	
	// Protected routes (require JWT authentication)
	protected := auth.Group("", middleware.JWTMiddleware(config))
	protected.Get("/me", authHandler.Me)
}
