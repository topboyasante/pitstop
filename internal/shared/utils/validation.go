package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ParseValidationErrors converts Go validator errors to simple field: message map
func ParseValidationErrors(err error) map[string]string {
	details := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			message := getValidationMessage(fieldError)
			details[field] = message
		}
	}

	return details
}

// getValidationMessage returns user-friendly validation messages
func getValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", fieldError.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters long", fieldError.Param())
	case "len":
		return fmt.Sprintf("Must be exactly %s characters long", fieldError.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fieldError.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fieldError.Param())
	case "gt":
		return fmt.Sprintf("Must be greater than %s", fieldError.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", fieldError.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fieldError.Param())
	case "uuid":
		return "Must be a valid UUID"
	case "url":
		return "Must be a valid URL"
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "numeric":
		return "Must be numeric"
	default:
		return "Invalid value"
	}
}

// SendValidationErrors sends a structured validation error response
func SendValidationErrors(c *fiber.Ctx, err error, requestID string) error {
	validationDetails := ParseValidationErrors(err)
	errorCount := len(validationDetails)
	
	var message string
	if errorCount == 1 {
		message = "The request contains 1 validation error"
	} else {
		message = fmt.Sprintf("The request contains %d validation errors", errorCount)
	}

	// Convert map[string]string to map[string]any
	details := make(map[string]any)
	for k, v := range validationDetails {
		details[k] = v
	}

	return SendValidationError(c, message, requestID, details)
}

// ValidateStruct validates a struct and returns formatted error response if invalid
func ValidateStruct(c *fiber.Ctx, validator *validator.Validate, data any, requestID string) error {
	if err := validator.Struct(data); err != nil {
		return SendValidationErrors(c, err, requestID)
	}
	return nil
}