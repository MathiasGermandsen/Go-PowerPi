package database

import (
	"fmt"
	"time"

	"Power-Pi/config"
	"Power-Pi/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to retrieve sql.DB from gorm")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := db.AutoMigrate(&PowerTable{}); err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to auto-migrate")
	}

	DB = db
	logger.Log.Info().Msg("Database connected with pooled connections")
}
