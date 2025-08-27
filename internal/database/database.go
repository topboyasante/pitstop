package database

import (
	"fmt"

	"github.com/topboyasante/pitstop/internal/config"
	"github.com/topboyasante/pitstop/internal/logger"
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

	return db, nil
}
