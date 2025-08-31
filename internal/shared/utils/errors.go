package utils

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ErrorCode represents different types of errors
type ErrorCode string

const (
	// Client errors (4xx)
	ErrCodeInvalidRequest   ErrorCode = "INVALID_REQUEST"
	ErrCodeMissingParameter ErrorCode = "MISSING_PARAMETER" 
	ErrCodeInvalidParameter ErrorCode = "INVALID_PARAMETER"
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeResourceExists   ErrorCode = "RESOURCE_EXISTS"
	ErrCodeRateLimit        ErrorCode = "RATE_LIMIT_EXCEEDED"

	// Server errors (5xx)
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeExternalService    ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeDatabaseError      ErrorCode = "DATABASE_ERROR"
)

// APIError represents a structured API error
type APIError struct {
	Code      ErrorCode      `json:"code"`
	Message   string         `json:"message"`
	RequestID string         `json:"request_id"`
	Timestamp time.Time      `json:"timestamp"`
	Details   map[string]any `json:"details,omitempty"`
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// HTTPStatusFromErrorCode maps error codes to HTTP status codes
func HTTPStatusFromErrorCode(code ErrorCode) int {
	switch code {
	case ErrCodeInvalidRequest, ErrCodeMissingParameter, ErrCodeInvalidParameter:
		return fiber.StatusBadRequest
	case ErrCodeValidationFailed:
		return fiber.StatusUnprocessableEntity
	case ErrCodeUnauthorized:
		return fiber.StatusUnauthorized
	case ErrCodeForbidden:
		return fiber.StatusForbidden
	case ErrCodeNotFound:
		return fiber.StatusNotFound
	case ErrCodeResourceExists:
		return fiber.StatusConflict
	case ErrCodeRateLimit:
		return fiber.StatusTooManyRequests
	case ErrCodeInternalError, ErrCodeDatabaseError:
		return fiber.StatusInternalServerError
	case ErrCodeServiceUnavailable, ErrCodeExternalService:
		return fiber.StatusServiceUnavailable
	default:
		return fiber.StatusInternalServerError
	}
}

// NewAPIError creates a new APIError
func NewAPIError(code ErrorCode, message string, requestID string, details ...map[string]any) APIError {
	apiErr := APIError{
		Code:      code,
		Message:   message,
		RequestID: requestID,
		Timestamp: time.Now(),
	}

	if len(details) > 0 {
		apiErr.Details = details[0]
	}

	return apiErr
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(code ErrorCode, message string, requestID string, details ...map[string]any) ErrorResponse {
	return ErrorResponse{
		Error: NewAPIError(code, message, requestID, details...),
	}
}

// SendErrorResponse sends a standardized error response with proper status code and headers
func SendErrorResponse(c *fiber.Ctx, code ErrorCode, message string, requestID string, details ...map[string]any) error {
	statusCode := HTTPStatusFromErrorCode(code)
	errorResponse := NewErrorResponse(code, message, requestID, details...)

	c.Set("X-Request-ID", requestID)
	return c.Status(statusCode).JSON(errorResponse)
}

// Common error response helpers
func SendBadRequestError(c *fiber.Ctx, message string, requestID string, details ...map[string]any) error {
	return SendErrorResponse(c, ErrCodeInvalidRequest, message, requestID, details...)
}

func SendValidationError(c *fiber.Ctx, message string, requestID string, details ...map[string]any) error {
	return SendErrorResponse(c, ErrCodeValidationFailed, message, requestID, details...)
}

func SendNotFoundError(c *fiber.Ctx, resource string, requestID string) error {
	return SendErrorResponse(c, ErrCodeNotFound, fmt.Sprintf("%s not found", resource), requestID)
}

func SendUnauthorizedError(c *fiber.Ctx, requestID string) error {
	return SendErrorResponse(c, ErrCodeUnauthorized, "Authentication required", requestID)
}

func SendForbiddenError(c *fiber.Ctx, requestID string) error {
	return SendErrorResponse(c, ErrCodeForbidden, "Access denied", requestID)
}

func SendConflictError(c *fiber.Ctx, resource string, requestID string) error {
	return SendErrorResponse(c, ErrCodeResourceExists, fmt.Sprintf("%s already exists", resource), requestID)
}

func SendInternalError(c *fiber.Ctx, requestID string, details ...map[string]any) error {
	return SendErrorResponse(c, ErrCodeInternalError, "An internal error occurred", requestID, details...)
}

func SendExternalServiceError(c *fiber.Ctx, message string, requestID string, details ...map[string]any) error {
	return SendErrorResponse(c, ErrCodeExternalService, message, requestID, details...)
}
