package post

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/middleware"
	"github.com/topboyasante/pitstop/internal/modules/post/handler"
)

// RegisterRoutes registers all post-related routes
func RegisterRoutes(router fiber.Router, postHandler *handler.PostHandler, commentHandler *handler.CommentHandler, likeHandler *handler.LikeHandler) {
	posts := router.Group("/posts")
	
	// Public routes
	posts.Get("/", postHandler.GetAllPosts)
	posts.Get("/:id", postHandler.GetPost)
	posts.Get("/:post_id/comments", commentHandler.GetComments)
	posts.Get("/:post_id/likes", likeHandler.GetLikesByPost)
	
	// Protected routes
	protected := posts.Group("", middleware.JWTMiddleware(config.Get()))
	protected.Post("/", postHandler.CreatePost)
	protected.Post("/:post_id/comments", commentHandler.CreateComment)
	protected.Post("/:post_id/comments/:parent_comment_id/reply", commentHandler.CreateReply)
	protected.Post("/:post_id/like", likeHandler.ToggleLike)
	protected.Get("/:post_id/like/status", likeHandler.CheckUserLiked)
}
