// @title Pitstop API
// @version 1.0
// @description A RESTful API built with Go and Fiber
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	docs "github.com/topboyasante/pitstop/docs/v1"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/database"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/middleware"
	"github.com/topboyasante/pitstop/internal/core/redis"
	"github.com/topboyasante/pitstop/internal/modules/auth"
	"github.com/topboyasante/pitstop/internal/modules/post"
	"github.com/topboyasante/pitstop/internal/modules/user"
	"github.com/topboyasante/pitstop/internal/provider"
)

func main() {
	logger.InitGlobal()

	err := config.InitGlobal()
	if err != nil {
		logger.Error("failed to start server - configuration error: %v", err)
	}
	cfg := config.Get()

	db, err := database.Init(cfg)
	if err != nil {
		logger.Fatal("failed to connect to database: %v", err)
		log.Panicf("error: %s", err)
	}

	redisClient, err := redis.Connect(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", "error", err)
		log.Panicf("error: %s", err)
	}

	// Ensure Redis connection is closed on application exit
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("error closing redis connection: %v", err)
		}
	}()

	// Initialize validator
	validator := validator.New()

	// Initialize provider with dependency injection
	provider := provider.NewProvider(db, redisClient, cfg, validator)

	// Update Swagger host dynamically
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
	})

	// Add rate limiting and request logging middleware
	app.Use(middleware.RateLimiter((provider.Redis)))
	app.Use(middleware.RequestLogger())

	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/v1/swagger.json",
		Path:     "docs",
		Title:    "Pitstop API Documentation",
	}))

	v1 := app.Group("/api/v1")

	// Register modular routes
	auth.RegisterRoutes(v1, provider.AuthHandler)
	user.RegisterRoutes(v1, provider.UserHandler)
	post.RegisterRoutes(v1, provider.PostHandler, provider.CommentHandler, provider.LikeHandler)

	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		logger.Fatal("failed to start server: %v", err)
		log.Panicf("error: %s", err)
	}
}
