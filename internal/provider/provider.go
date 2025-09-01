package provider

import (
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/modules/auth/handler"
	"github.com/topboyasante/pitstop/internal/modules/auth/service"
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

	// Auth module
	AuthHandler *handler.AuthHandler

	// Module dependencies (can be accessed by other modules if needed)
	AuthService *service.AuthService
}

// NewProvider creates and initializes the dependency injection container
func NewProvider(db *gorm.DB, redis *redis.Client, cfg *config.Config, validator *validator.Validate) *Provider {
	// Initialize event bus
	eventBus := events.NewEventBus()

	// Initialize Auth module
	authService := service.NewAuthService(cfg, validator)
	authHandler := handler.NewAuthHandler(authService)

	// Set up event subscribers
	setupEventSubscribers(eventBus, authService)

	return &Provider{
		DB:        db,
		Redis:     redis,
		Config:    cfg,
		Validator: validator,
		EventBus:  eventBus,

		AuthHandler: authHandler,

		AuthService: authService,
	}
}

// setupEventSubscribers configures cross-module event handlers
func setupEventSubscribers(eventBus *events.EventBus, authService *service.AuthService) {
	// Example: When a user registers, other modules can react
	eventBus.Subscribe("UserRegistered", func(event events.Event) {
		userEvent := event.(*events.UserRegistered)
		_ = userEvent
		// Could trigger welcome email, create default garage, etc.
	})
}
