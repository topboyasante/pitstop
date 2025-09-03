package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/topboyasante/pitstop/internal/core/logger"
	"golang.org/x/oauth2"
)

// The API configuration structure
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	OAuth    oauth2.Config
	Redis    RedisConfig
}

// Server configuration structure
type ServerConfig struct {
	Port        string
	JWTSecret   string
	JWTIssuer   string
	FrontendURL string
}

// Database configuration structure
type DatabaseConfig struct {
	Host           string
	DatabaseName   string
	User           string
	Password       string
	SslMode        string
	ChannelBinding string
}

type RedisConfig struct {
	URL string
}

// getEnvWithDefault retrieves an environment variable or returns a default value if not set.
// It logs whether the actual environment variable was used or if it fell back to the default.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Info("Using default value for environment variable", "key", key, "default", defaultValue)
		return defaultValue
	}
	logger.Info("Loaded environment variable", "key", key)
	return value
}

// New creates and initializes a new Config instance by loading environment variables.
// It attempts to load from a .env file first, then reads required and optional environment variables.
// Returns a fully configured Config struct or an error if required variables are missing.
func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Unable to load .env file", "error", err)
	} else {
		logger.Info("Successfully loaded .env file")
	}

	port := getEnv("PORT", "8080")
	dbHost := getEnv("PGHOST", "")
	dbName := getEnv("PGDATABASE", "")
	dbUser := getEnv("PGUSER", "")
	dbPassword := getEnv("PGPASSWORD", "")
	dbSslMode := getEnv("PGSSLMODE", "disable")
	dbChannelBinding := getEnv("PGCHANNELBINDING", "disable")
	oauthClientID := getEnv("OAUTH_CLIENT_ID", "")
	oauthClientSecret := getEnv("OAUTH_CLIENT_SECRET", "")
	oauthRedirectURI := getEnv("OAUTH_REDIRECT_URI", "")
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")
	jwtSecret := getEnv("JWT_SECRET", "dummy")
	jwtIssuer := getEnv("JWT_ISSUER", "pitstop")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")

	logger.Info("Configuration loaded successfully",
		"server_port", port,
		"database_configured", true)

	return &Config{
		Server: ServerConfig{
			Port:        port,
			JWTSecret:   jwtSecret,
			JWTIssuer:   jwtIssuer,
			FrontendURL: frontendURL,
		},
		Database: DatabaseConfig{
			Host:           dbHost,
			DatabaseName:   dbName,
			User:           dbUser,
			Password:       dbPassword,
			SslMode:        dbSslMode,
			ChannelBinding: dbChannelBinding,
		},
		OAuth: oauth2.Config{
			ClientID:     oauthClientID,
			ClientSecret: oauthClientSecret,
			RedirectURL:  oauthRedirectURI,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
		},
		Redis: RedisConfig{
			URL: redisURL,
		},
	}, nil
}

// Global config instance
var GlobalConfig *Config

// InitGlobal initializes the global config instance
func InitGlobal() error {
	var err error
	GlobalConfig, err = New()
	return err
}

// Get returns the global config instance
func Get() *Config {
	return GlobalConfig
}
