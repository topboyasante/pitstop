package provider

import (
	"github.com/go-playground/validator/v10"
	"github.com/topboyasante/pitstop/internal/api/v1/controllers"
	"github.com/topboyasante/pitstop/internal/api/v1/repositories"
	"github.com/topboyasante/pitstop/internal/api/v1/services"
	"gorm.io/gorm"
)

type Provider struct {
	PostController *controllers.PostController
	DB             *gorm.DB
}

func NewProvider(db *gorm.DB, validator *validator.Validate) *Provider {
	// Initialize repositories
	postRepo := repositories.NewPostRepository(db)

	// Initialize services
	postService := services.NewPostService(postRepo, validator)

	// Initialize controllers
	postController := controllers.NewPostController(postService)

	return &Provider{
		PostController: postController,
		DB:             db,
	}
}