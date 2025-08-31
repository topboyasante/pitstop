// @title Pitstop API
// @version 1.0
// @description A RESTful API built with Go and Fiber
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	_ "github.com/topboyasante/pitstop/docs/v1"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/database"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"github.com/topboyasante/pitstop/internal/core/middleware"
	"github.com/topboyasante/pitstop/internal/core/redis"
	"github.com/topboyasante/pitstop/internal/modules/auth"
	"github.com/topboyasante/pitstop/internal/provider"
)

func main() {
	logger.InitGlobal()

	config, err := config.New()
	if err != nil {
		logger.Error("failed to start server - configuration error: %v", err)
	}

	db, err := database.Init(config)
	if err != nil {
		logger.Fatal("failed to connect to database: %v", err)
		log.Panicf("error: %s", err)
	}

	redisClient, err := redis.Connect(config)
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
	provider := provider.NewProvider(db, redisClient, config, validator)

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

	if err := app.Listen(":" + config.Server.Port); err != nil {
		logger.Fatal("failed to start server: %v", err)
		log.Panicf("error: %s", err)
	}
}
