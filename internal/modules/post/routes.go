package post

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/middleware"
	"github.com/topboyasante/pitstop/internal/modules/post/handler"
)

// RegisterRoutes registers all post-related routes
func RegisterRoutes(router fiber.Router, postHandler *handler.PostHandler) {
	posts := router.Group("/posts")
	
	// Public routes
	posts.Get("/", postHandler.GetAllPosts)
	posts.Get("/:id", postHandler.GetPost)
	
	// Protected routes
	protected := posts.Group("", middleware.JWTMiddleware(config.Get()))
	protected.Post("/", postHandler.CreatePost)
}
