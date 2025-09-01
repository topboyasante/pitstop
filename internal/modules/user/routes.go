package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/modules/user/handler"
)

// RegisterRoutes registers all user-related routes
func RegisterRoutes(router fiber.Router, userHandler *handler.UserHandler) {
	users := router.Group("/users")
	
	users.Get("/", userHandler.GetAllUsers)
	users.Post("/", userHandler.CreateUser)        // Requires auth
	users.Get("/:id", userHandler.GetUser)
}
