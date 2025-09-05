package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/response"
	"github.com/topboyasante/pitstop/internal/shared/utils"
)

// ErrorHandler is a middleware that handles panics and provides consistent error responses
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Generate or get request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = utils.GenerateRequestID()
		}

		// Log the error with context
		logger.Error("Request failed with error",
			"event", "request.error",
			"request_id", requestID,
			"path", c.Path(),
			"method", c.Method(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
			"error", err.Error())

		// Check if it's a Fiber error
		if e, ok := err.(*fiber.Error); ok {
			return response.ErrorJSON(c, e.Code, "FIBER_ERROR", e.Message, requestID)
		}

		// Default internal server error
		return response.InternalErrorJSON(c, "Request ID: "+requestID)
	}
}