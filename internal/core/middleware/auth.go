package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/shared/utils"
)

// JWTMiddleware validates JWT tokens from Authorization header
func JWTMiddleware(config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Debug("JWT middleware validating request")

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Warn("Missing Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Check for Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Warn("Invalid Authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := tokenParts[1]

		// Validate JWT token
		token, err := utils.ValidateJWTToken(config, tokenString)
		if err != nil {
			logger.Error("JWT token validation failed", "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Extract claims
		userID, audience, _, err := utils.ExtractClaims(token)
		if err != nil {
			logger.Error("Failed to extract JWT claims", "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Store user info in context for route handlers
		c.Locals("userID", userID)
		c.Locals("audience", audience)

		logger.Debug("JWT middleware validation successful", "userID", userID, "audience", audience)
		return c.Next()
	}
}

// OptionalJWTMiddleware validates JWT tokens but doesn't block requests if missing
func OptionalJWTMiddleware(config *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No token provided, continue without user context
			return c.Next()
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// Invalid format, continue without user context
			return c.Next()
		}

		tokenString := tokenParts[1]
		token, err := utils.ValidateJWTToken(config, tokenString)
		if err != nil {
			// Invalid token, continue without user context
			return c.Next()
		}

		userID, audience, _, err := utils.ExtractClaims(token)
		if err != nil {
			// Invalid claims, continue without user context
			return c.Next()
		}

		// Store user info in context if token is valid
		c.Locals("userID", userID)
		c.Locals("audience", audience)

		logger.Debug("Optional JWT middleware found valid token", "userID", userID)
		return c.Next()
	}
}