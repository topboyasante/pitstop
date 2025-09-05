package provider

import (
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/logger"
	authHandler "github.com/topboyasante/pitstop/internal/modules/auth/handler"
	authService "github.com/topboyasante/pitstop/internal/modules/auth/service"
	postHandler "github.com/topboyasante/pitstop/internal/modules/post/handler"
	postRepository "github.com/topboyasante/pitstop/internal/modules/post/repository"
	postService "github.com/topboyasante/pitstop/internal/modules/post/service"
	userHandler "github.com/topboyasante/pitstop/internal/modules/user/handler"
	userRepository "github.com/topboyasante/pitstop/internal/modules/user/repository"
	userService "github.com/topboyasante/pitstop/internal/modules/user/service"
	"github.com/topboyasante/pitstop/internal/shared/events"
	"gorm.io/gorm"
)

// Provider is the central dependency injection container for the modular monolith
type Provider struct {
	// Database
	DB    *gorm.DB
	Redis *redis.Client

	// Shared services
	Config    *config.Config
	Validator *validator.Validate
	EventBus  *events.EventBus

	// Handlers
	AuthHandler *authHandler.AuthHandler
	UserHandler *userHandler.UserHandler
	PostHandler *postHandler.PostHandler

	// Module dependencies (can be accessed by other modules if needed)
	AuthService *authService.AuthService
	UserService *userService.UserService
	PostService *postService.PostService
}

// NewProvider creates and initializes the dependency injection container
func NewProvider(db *gorm.DB, redis *redis.Client, cfg *config.Config, validator *validator.Validate) *Provider {
	// Initialize event bus
	eventBus := events.NewEventBus()

	// Initialize User module
	userRepo := userRepository.NewUserRepository(db)
	userService := userService.NewUserService(userRepo, validator, eventBus)
	userHandler := userHandler.NewUserHandler(userService)

	// Initialize Post module
	postRepo := postRepository.NewPostRepository(db)
	postService := postService.NewPostService(postRepo, validator, eventBus)
	postHandler := postHandler.NewPostHandler(postService)

	// Initialize Auth module (depends on user service)
	authService := authService.NewAuthService(cfg, redis, eventBus, validator, userService)
	authHandler := authHandler.NewAuthHandler(authService)

	// Set up event subscribers
	setupEventSubscribers(eventBus, authService)

	return &Provider{
		DB:        db,
		Redis:     redis,
		Config:    cfg,
		Validator: validator,
		EventBus:  eventBus,

		AuthHandler: authHandler,
		UserHandler: userHandler,
		PostHandler: postHandler,

		AuthService: authService,
		UserService: userService,
		PostService: postService,
	}
}

// setupEventSubscribers configures cross-module event handlers
func setupEventSubscribers(eventBus *events.EventBus, authService *authService.AuthService) {
	eventBus.Subscribe("AuthenticationSuccessful", func(event events.Event) {
		userEvent := event.(*events.AuthenticationSuccessful)
		_ = userEvent
		logger.Info("Recieved a pblished event",
			"event", event)

		// create a user record in the db if we have to
		// Generate the neccessary tokens they need
		//return the tokens to the client or sth
	})

	// Example: When a user registers, other modules can react
	// eventBus.Subscribe("UserRegistered", func(event events.Event) {
	// 	userEvent := event.(*events.UserRegistered)
	// 	_ = userEvent
	// 	// Could trigger welcome email, create default garage, etc.
	// })
}
