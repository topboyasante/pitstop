package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/api/v1/controller"
)

func RegisterPostRoutes(a fiber.Router) {
	postController := controller.NewPostController()
	postRoutes := a.Group("/posts")

	{
		postRoutes.Get("/", postController.GetAllPosts)
	}
}
