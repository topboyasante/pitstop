package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Data      any       `json:"data"`
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

// SendSuccessResponse sends a standardized success response
func SendSuccessResponse(c *fiber.Ctx, data any, requestID string) error {
	response := SuccessResponse{
		Data:      data,
		RequestID: requestID,
		Timestamp: time.Now(),
	}

	c.Set("X-Request-ID", requestID)
	return c.JSON(response)
}

// SendCreatedResponse sends a 201 Created response
func SendCreatedResponse(c *fiber.Ctx, data any, requestID string) error {
	response := SuccessResponse{
		Data:      data,
		RequestID: requestID,
		Timestamp: time.Now(),
	}

	c.Set("X-Request-ID", requestID)
	return c.Status(fiber.StatusCreated).JSON(response)
}
