package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// GenerateRequestID creates a readable request ID with format: req_<random_string>.
// This must be used in each controller to generate a unique request ID for logging and tracing
func GenerateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	randomString := base64.URLEncoding.EncodeToString(b)[:10]
	return fmt.Sprintf("req_%s", randomString)
}

// ExtractUserIDFromContext extracts user ID from Fiber context locals
func ExtractUserIDFromContext(c *fiber.Ctx) (string, error) {
	userIDLocal := c.Locals("userID")
	if userIDLocal == nil {
		return "", errors.New("user ID not found in context")
	}
	
	userID, ok := userIDLocal.(string)
	if !ok || userID == "" {
		return "", errors.New("invalid user ID in context")
	}
	
	return userID, nil
}