package health

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/modules/health/handler"
)

// RegisterRoutes registers all health-related routes
func RegisterRoutes(router fiber.Router, healthHandler *handler.HealthHandler) {
	// Health check endpoint - public access
	router.Get("/health", healthHandler.HealthCheck)
}