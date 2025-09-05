package redis

import (
	"context"
	"crypto/tls"

	"github.com/redis/go-redis/v9"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/logger"
)

var ctx = context.Background()

func Connect(config *config.Config) (*redis.Client, error) {
	logger.Info("Attempting to connect to Redis", "url", config.Redis.URL)

	opts, err := redis.ParseURL(config.Redis.URL)
	if err != nil {
		logger.Error("Failed to parse Redis URL", "error", err.Error())
		return nil, err
	}

	// Configure TLS for production Redis connections
	if opts.TLSConfig != nil {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         opts.Addr,
		}
	}

	client := redis.NewClient(opts)

	// Test connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed to connect to Redis",
			"url", config.Redis.URL,
			"error", err.Error())
		return nil, err
	}

	logger.Info("Connected to Redis successfully", "response", pong)
	return client, nil
}
