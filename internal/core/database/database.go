package database

import (
	"fmt"

	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/logger"
	postDomain "github.com/topboyasante/pitstop/internal/modules/post/domain"
	userDomain "github.com/topboyasante/pitstop/internal/modules/user/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg *config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.Database.Host +
		" user=" + cfg.Database.User +
		" password=" + cfg.Database.Password +
		" dbname=" + cfg.Database.DatabaseName +
		" sslmode=" + cfg.Database.SslMode +
		" channel_binding=" + cfg.Database.ChannelBinding

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	logger.Info("successfully connected to database")

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// runMigrations runs all database migrations
func runMigrations(db *gorm.DB) error {
	logger.Info("Running database migrations")

	err := db.AutoMigrate(
		&userDomain.User{},
		&postDomain.Post{},
	)

	if err != nil {
		logger.Error("Failed to run migrations", "error", err)
		return err
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
