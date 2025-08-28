package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/provider"
)

func RegisterV1AuthRoutes(a fiber.Router, p *provider.Provider) {
	authRoutes := a.Group("/auth")

	{
		authRoutes.Get("/", p.AuthController.Authenticate)
		authRoutes.Get("/callback", p.AuthController.Callback)
	}
}
