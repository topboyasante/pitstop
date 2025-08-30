package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/provider"
	"github.com/topboyasante/pitstop/internal/logger"
)

func RegisterV1AuthRoutes(a fiber.Router, p *provider.Provider) {
	logger.Info("Registering auth routes", 
		"event", "routes.auth_registered",
		"version", "v1")
	
	authRoutes := a.Group("/auth")

	{
		authRoutes.Get("/", p.AuthController.Authenticate)
		authRoutes.Get("/callback", p.AuthController.Callback)
	}
	
	logger.Info("Auth routes registered successfully", 
		"event", "routes.auth_registration_complete",
		"endpoints", []string{"/auth", "/auth/callback"})
}
