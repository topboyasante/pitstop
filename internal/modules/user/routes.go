package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/middleware"
	"github.com/topboyasante/pitstop/internal/modules/user/handler"
)

// RegisterRoutes registers all user-related routes
func RegisterRoutes(router fiber.Router, userHandler *handler.UserHandler, followHandler *handler.FollowHandler) {
	users := router.Group("/users")
	
	// Public routes
	users.Get("/", userHandler.GetAllUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Get("/:user_id/followers", followHandler.GetFollowers)
	users.Get("/:user_id/following", followHandler.GetFollowing)
	
	// Protected routes
	protected := users.Group("", middleware.JWTMiddleware(config.Get()))
	protected.Post("/", userHandler.CreateUser)
	protected.Post("/:user_id/follow", followHandler.ToggleFollow)
	protected.Get("/:user_id/follow/status", followHandler.CheckFollowStatus)
}
