package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/middleware"
	"github.com/topboyasante/pitstop/internal/modules/user/handler"
)

// RegisterRoutes registers all user-related routes
func RegisterRoutes(router fiber.Router, userHandler *handler.UserHandler) {
	users := router.Group("/users")
	
	// Public routes
	users.Get("/", userHandler.GetAllUsers)
	users.Get("/:id", userHandler.GetUser)
	
	// Protected routes
	protected := users.Group("", middleware.JWTMiddleware(config.Get()))
	protected.Post("/", userHandler.CreateUser)
}
