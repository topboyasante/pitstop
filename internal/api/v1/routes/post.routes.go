package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/provider"
)

func RegisterV1PostRoutes(a fiber.Router, p *provider.Provider) {
	postRoutes := a.Group("/posts")

	{
		postRoutes.Get("/", p.PostController.GetAllPosts)
	}
}