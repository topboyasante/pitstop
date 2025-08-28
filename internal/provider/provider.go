package provider

import (
	"github.com/go-playground/validator/v10"
	"github.com/topboyasante/pitstop/internal/api/v1/controllers"
	"github.com/topboyasante/pitstop/internal/api/v1/repositories"
	"github.com/topboyasante/pitstop/internal/api/v1/services"
	"github.com/topboyasante/pitstop/internal/config"
	"gorm.io/gorm"
)

// The provider is the central place where all dependencies are initialized and injected.
type Provider struct {
	PostController *controllers.PostController
	AuthController *controllers.AuthController
	DB             *gorm.DB
	config         *config.Config
}

func NewProvider(db *gorm.DB, validator *validator.Validate, config *config.Config) *Provider {
	// Initialize repositories
	postRepo := repositories.NewPostRepository(db)

	// Initialize services
	postService := services.NewPostService(postRepo, validator)
	authService := services.NewAuthService(config)

	// Initialize controllers
	postController := controllers.NewPostController(postService)
	authController := controllers.NewAuthController(authService)

	return &Provider{
		PostController: postController,
		AuthController: authController,
		DB:             db,
		config:         config,
	}
}
