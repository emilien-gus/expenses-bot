package db

import (
	"fmt"
	"log"

	"github.com/ai-bot-expenses/backend/internal/config"
	"github.com/ai-bot-expenses/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init connects to PostgreSQL using the configuration and performs auto-migration.
func Init(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established. Running auto-migrations...")
	if err := db.AutoMigrate(&models.Expense{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return db, nil
}
