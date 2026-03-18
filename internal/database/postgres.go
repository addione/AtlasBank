package database

import (
	"fmt"

	"github.com/atlasbank/api/internal/config"
	"github.com/atlasbank/api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// RunMigrations runs GORM auto-migrations
// Note: For production, consider using golang-migrate CLI tool or SQL migration files
func RunMigrations(db *gorm.DB) error {
	// Auto-migrate models
	err := db.AutoMigrate(
		&models.User{},
		&models.OTP{},
		&models.Account{},
		&models.Transaction{},
		&models.Notification{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Database migrations completed successfully")
	return nil
}
