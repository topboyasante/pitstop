package post

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/modules/post/handler"
)

// RegisterRoutes registers all post-related routes
func RegisterRoutes(router fiber.Router, postHandler *handler.PostHandler) {
	posts := router.Group("/posts")
	
	posts.Get("/", postHandler.GetAllPosts)
	posts.Post("/", postHandler.CreatePost)        // Requires auth
	posts.Get("/:id", postHandler.GetPost)
}
