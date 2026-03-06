package database

import (
	"fmt"
	"practise/go_fiber/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config, log *zap.SugaredLogger) (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Errorw("Failed to connect to database",
			"host", cfg.DBHost,
			"port", cfg.DBPort,
			"dbname", cfg.DBName,
			"error", err,
		)
		return nil, err
	}

	log.Infow("Database connected successfully",
		"host", cfg.DBHost,
		"port", cfg.DBPort,
		"dbname", cfg.DBName,
	)
	return db, nil
}
