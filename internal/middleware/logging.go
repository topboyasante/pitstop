package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/logger"
	"github.com/topboyasante/pitstop/internal/utils"
)

// RequestLogger creates middleware that logs all requests with structured data
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID if not present
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = utils.GenerateRequestID()
			c.Set("X-Request-ID", requestID)
		}

		start := time.Now()

		// Log request start
		logger.Info("Request started",
			"event", "request.started",
			"request_id", requestID,
			"path", c.Path(),
			"method", c.Method(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"))

		// Process request
		err := c.Next()

		duration := time.Since(start)

		// Log request completion
		logger.Info("Request completed",
			"event", "request.completed",
			"request_id", requestID,
			"path", c.Path(),
			"method", c.Method(),
			"status_code", c.Response().StatusCode(),
			"duration_ms", duration.Milliseconds(),
			"ip", c.IP())

		return err
	}
}