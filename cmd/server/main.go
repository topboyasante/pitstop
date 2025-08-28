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
	"github.com/topboyasante/pitstop/internal/api/v1/routes"
	"github.com/topboyasante/pitstop/internal/config"
	"github.com/topboyasante/pitstop/internal/database"
	"github.com/topboyasante/pitstop/internal/logger"
	"github.com/topboyasante/pitstop/internal/provider"
	_ "github.com/topboyasante/pitstop/docs/v1"
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

	// Initialize validator
	validator := validator.New()

	// Initialize provider with dependency injection
	provider := provider.NewProvider(db, validator, config)

	app := fiber.New()

	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/v1/swagger.json",
		Path:     "docs",
		Title:    "Pitstop API Documentation",
	}))

	v1 := app.Group("/api/v1")
	routes.RegisterV1PostRoutes(v1, provider)
	routes.RegisterV1AuthRoutes(v1, provider)

	if err := app.Listen(":" + config.Server.Port); err != nil {
		logger.Fatal("failed to start server: %v", err)
		log.Panicf("error: %s", err)
	}
}
