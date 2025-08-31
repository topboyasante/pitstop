package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/shared/utils"
)

func RateLimiter(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate or get request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = utils.GenerateRequestID()
		}

		//  Get the API KEY from headers
		apiKey := c.Get("X-API-KEY")
		if apiKey == "" {
			logger.Error("Request failed with error",
				"event", "request.error",
				"request_id", requestID,
				"path", c.Path(),
				"method", c.Method(),
				"ip", c.IP(),
				"user_agent", c.Get("User-Agent"),
				"error", "API key is required")

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "API key is required",
			})
		}

		// Generate a unique key for the API key and current hour
		key := fmt.Sprintf("rate_limit:%s:%s", apiKey, time.Now().Format("2006-01-02-15"))

		ctx := context.Background()

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

		remaining := max(1000-int(count), 0)
		allowed := count <= 1000

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", "1000")
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

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

		return c.Next()
	}
}
