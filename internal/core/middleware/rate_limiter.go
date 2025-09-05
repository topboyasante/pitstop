package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/shared/utils"
)

func RateLimiter(redisClient *redis.Client) fiber.Handler {
	// Define public routes that don't require JWT authentication
	publicRoutes := []string{
		"/api/v1/auth/google",
		"/api/v1/auth/google/callback",
		"/api/v1/auth/exchange",
		"/api/v1/auth/refresh",
		"/api/v1/docs",
		"/health",
		"/docs",
	}

	return func(c *fiber.Ctx) error {
		// Generate or get request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = utils.GenerateRequestID()
		}

		path := c.Path()

		// Check if this is a public route
		isPublic := false
		for _, publicRoute := range publicRoutes {
			if strings.HasPrefix(path, publicRoute) {
				isPublic = true
				break
			}
		}

		// Get user ID from JWT middleware if available (for user-specific rate limiting)
		userID := ""
		if userIDLocal := c.Locals("userID"); userIDLocal != nil {
			if uid, ok := userIDLocal.(string); ok {
				userID = uid
			}
		}

		// Get IP address
		ip := c.IP()

		// Current hour for Redis key
		hour := time.Now().Format("2006-01-02-15")
		ctx := context.Background()

		var limits [][4]any

		if isPublic {
			// For public routes, apply IP-only rate limiting
			limits = [][4]any{
				{true, fmt.Sprintf("rate_limit:public:ip:%s:%s", ip, hour), 100, "public route rate limit by IP"},
			}
		} else {
			// For authenticated routes, apply user-based and IP-based rate limiting
			limits = [][4]any{
				{userID != "", fmt.Sprintf("rate_limit:auth:user:%s:%s", userID, hour), 1000, "authenticated rate limit by user"},
				{true, fmt.Sprintf("rate_limit:auth:ip:%s:%s", ip, hour), 500, "authenticated rate limit by IP"},
			}
		}

		// we don't need the index, so we comment it out
		for _, limit := range limits {
			// Check if the condition is true. If it's true, we need to enforce this limit. If false, skip to the next limit.
			shouldCheckThisLimit := limit[0].(bool)
			if !shouldCheckThisLimit {
				continue
			}

			key := limit[1].(string)
			maxRequests := limit[2].(int)
			resource := limit[3].(string)

			// Increase the count for the key provided. if the key does not exist, it will be created with count value 1
			count, err := redisClient.Incr(ctx, key).Result()
			if err != nil {
				logger.Error("Request failed with error",
					"event", "request.error",
					"request_id", requestID,
					"path", c.Path(),
					"method", c.Method(),
					"ip", c.IP(),
					"user_agent", c.Get("User-Agent"),
					"error", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error",
				})
			}

			// Set expiration for the key to 1 hour if it's newly created
			// After one hour, the key will be automatically deleted, hence resetting the count
			// We check if count == 1 to avoid resetting the expiration on every request
			// which would defeat the purpose of rate limiting
			if count == 1 {
				if err := redisClient.Expire(ctx, key, time.Hour).Err(); err != nil {
					logger.Error("Failed to set expiration on rate limit key",
						"event", "redis.error",
						"request_id", requestID,
						"key", key,
						"error", err.Error())

					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to set rate limit expiration",
					})
				}
			}

			// Calculate remaining requests and reset time
			remaining := max(maxRequests-int(count), 0)
			resetTime := time.Now().Add(time.Hour).Truncate(time.Hour).Unix()

			// Determine if the request is allowed
			allowed := count <= int64(maxRequests)

			// Set rate limit headers
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
			c.Set("X-RateLimit-Resource", resource)

			// If limit exceeded, return 429 Too Many Requests
			if !allowed {
				logger.Warn("Request rate limit exceeded",
					"event", "request.rate_limited",
					"request_id", requestID,
					"path", c.Path(),
					"method", c.Method(),
					"ip", c.IP(),
					"user_agent", c.Get("User-Agent"))

				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Rate limit exceeded. Try again later.",
				})
			}
		}

		return c.Next()
	}
}
