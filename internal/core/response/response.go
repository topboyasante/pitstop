package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// APIResponse represents a standardized API response structure
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Meta      *MetaInfo   `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// MetaInfo contains pagination and other metadata
type MetaInfo struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
	HasNext    bool  `json:"has_next,omitempty"`
	HasPrev    bool  `json:"has_prev,omitempty"`
}

// NewPaginationMeta creates pagination metadata
func NewPaginationMeta(page, limit int, total int64, hasNext bool) *MetaInfo {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &MetaInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    page > 1,
	}
}

// Success creates a successful response
func Success(data interface{}, message string) *APIResponse {
	return &APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// SuccessWithMeta creates a successful response with metadata (e.g., pagination)
func SuccessWithMeta(data interface{}, message string, meta *MetaInfo) *APIResponse {
	return &APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	}
}

// Error creates an error response
func Error(code, message, details string) *APIResponse {
	return &APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
}

// ValidationError creates a validation error response
func ValidationError(message string, details string) *APIResponse {
	return Error("VALIDATION_ERROR", message, details)
}

// NotFoundError creates a not found error response
func NotFoundError(resource string) *APIResponse {
	return Error("NOT_FOUND", resource+" not found", "")
}

// InternalError creates an internal server error response
func InternalError(message string) *APIResponse {
	return Error("INTERNAL_ERROR", "Internal server error", message)
}

// UnauthorizedError creates an unauthorized error response
func UnauthorizedError() *APIResponse {
	return Error("UNAUTHORIZED", "Authentication required", "")
}

// ForbiddenError creates a forbidden error response
func ForbiddenError() *APIResponse {
	return Error("FORBIDDEN", "Access denied", "")
}

// JSON sends a JSON response using Fiber
func JSON(c *fiber.Ctx, statusCode int, resp *APIResponse) error {
	return c.Status(statusCode).JSON(resp)
}

// SuccessJSON sends a successful JSON response
func SuccessJSON(c *fiber.Ctx, data interface{}, message string) error {
	return JSON(c, fiber.StatusOK, Success(data, message))
}

// SuccessJSONWithMeta sends a successful JSON response with metadata
func SuccessJSONWithMeta(c *fiber.Ctx, data interface{}, message string, meta *MetaInfo) error {
	return JSON(c, fiber.StatusOK, SuccessWithMeta(data, message, meta))
}

// CreatedJSON sends a created JSON response
func CreatedJSON(c *fiber.Ctx, data interface{}, message string) error {
	return JSON(c, fiber.StatusCreated, Success(data, message))
}

// ErrorJSON sends an error JSON response
func ErrorJSON(c *fiber.Ctx, statusCode int, code, message, details string) error {
	return JSON(c, statusCode, Error(code, message, details))
}

// ValidationErrorJSON sends a validation error JSON response
func ValidationErrorJSON(c *fiber.Ctx, message, details string) error {
	return JSON(c, fiber.StatusBadRequest, ValidationError(message, details))
}

// NotFoundJSON sends a not found error JSON response
func NotFoundJSON(c *fiber.Ctx, resource string) error {
	return JSON(c, fiber.StatusNotFound, NotFoundError(resource))
}

// InternalErrorJSON sends an internal server error JSON response
func InternalErrorJSON(c *fiber.Ctx, message string) error {
	return JSON(c, fiber.StatusInternalServerError, InternalError(message))
}

// UnauthorizedJSON sends an unauthorized error JSON response
func UnauthorizedJSON(c *fiber.Ctx) error {
	return JSON(c, fiber.StatusUnauthorized, UnauthorizedError())
}

// ForbiddenJSON sends a forbidden error JSON response
func ForbiddenJSON(c *fiber.Ctx) error {
	return JSON(c, fiber.StatusForbidden, ForbiddenError())
}