package db

import (
	"fmt"

	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
 * Connect establishes a connection to the PostgreSQL database with pgvector extension
 */
func Connect(dsn string) (*gorm.DB, error) {
	// Open connection with PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable pgvector extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		return nil, fmt.Errorf("failed to enable pgvector extension: %w", err)
	}
	fmt.Println("✅ pgvector extension enabled successfully")

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	// Set connection pool settings (will be overridden by config)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	// Verify database connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("✅ PostgreSQL connection established successfully")

	return db, nil
}

/*
 * Migrate applies database migrations
 */
func Migrate(db *gorm.DB) error {
	// Run GORM AutoMigrate for all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Bot{},
		&models.Source{},
		&models.DocumentChunk{},
		&models.Message{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	fmt.Println("✅ Database migrations completed successfully")
	return nil
}

