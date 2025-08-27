// @title Pitstop API
// @version 1.0
// @description A RESTful API built with Go and Fiber
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"log"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/api/v1/routes"
	"github.com/topboyasante/pitstop/internal/config"
	"github.com/topboyasante/pitstop/internal/database"
	"github.com/topboyasante/pitstop/internal/logger"
)

func main() {
	logger.InitGlobal()

	config, err := config.New()
	if err != nil {
		logger.Error("failed to start server - configuration error: %v", err)
	}

	_, err = database.Init(config)
	if err != nil {
		logger.Fatal("failed to connect to database: %v", err)
		log.Panicf("error: %s", err)
	}

	app := fiber.New()

	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/v1/swagger.json",
		Path:     "docs",
	}))

	v1 := app.Group("/api/v1")
	routes.RegisterPostRoutes(v1)

	if err := app.Listen(":" + config.Server.Port); err != nil {
		logger.Fatal("failed to start server: %v", err)
		log.Panicf("error: %s", err)
	}
}
