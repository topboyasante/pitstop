package models

import "time"

// SuccessResponse represents the standard success response format
type SuccessResponse struct {
	Data      any       `json:"data" example:"{}"`
	RequestID string    `json:"request_id" example:"req_8n3mN9pQ2x"`
	Timestamp time.Time `json:"timestamp" example:"2023-12-01T10:30:00Z"`
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the error information
type ErrorDetail struct {
	Code      string         `json:"code" example:"VALIDATION_FAILED"`
	Message   string         `json:"message" example:"The request contains 2 validation errors"`
	RequestID string         `json:"request_id" example:"req_8n3mN9pQ2x"`
	Timestamp time.Time      `json:"timestamp" example:"2023-12-01T10:30:00Z"`
	Details   map[string]any `json:"details,omitempty"`
}

// AuthSuccessResponse represents a successful authentication response
type AuthSuccessResponse struct {
	Data      AuthData  `json:"data"`
	RequestID string    `json:"request_id" example:"req_8n3mN9pQ2x"`
	Timestamp time.Time `json:"timestamp" example:"2023-12-01T10:30:00Z"`
}

// AuthData contains authentication response data
type AuthData struct {
	Message string `json:"message" example:"Authentication successful"`
	Token   string `json:"token" example:"ya29.a0AfH6SMC..."`
}

// ValidationErrorResponse represents a validation error response
type ValidationErrorResponse struct {
	Error ValidationErrorDetail `json:"error"`
}

// ValidationErrorDetail contains validation error information
type ValidationErrorDetail struct {
	Code      string            `json:"code" example:"VALIDATION_FAILED"`
	Message   string            `json:"message" example:"The request contains 2 validation errors"`
	RequestID string            `json:"request_id" example:"req_8n3mN9pQ2x"`
	Timestamp time.Time         `json:"timestamp" example:"2023-12-01T10:30:00Z"`
	Details   map[string]string `json:"details"`
}

// BadRequestErrorResponse represents a bad request error response
type BadRequestErrorResponse struct {
	Error BadRequestErrorDetail `json:"error"`
}

// BadRequestErrorDetail contains bad request error information
type BadRequestErrorDetail struct {
	Code      string    `json:"code" example:"INVALID_REQUEST"`
	Message   string    `json:"message" example:"Authorization code is required"`
	RequestID string    `json:"request_id" example:"req_8n3mN9pQ2x"`
	Timestamp time.Time `json:"timestamp" example:"2023-12-01T10:30:00Z"`
}
