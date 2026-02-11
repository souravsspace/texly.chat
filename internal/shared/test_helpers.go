package shared

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
 * SetupTestDB creates a PostgreSQL test database connection
 * Requires DATABASE_URL_TEST environment variable or uses default test database
 */
func SetupTestDB() *gorm.DB {
	gin.SetMode(gin.TestMode)

	// Get test database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL_TEST")
	if dbURL == "" {
		// Default to local PostgreSQL test database
		dbURL = "postgres://texly:texly_dev@localhost:5432/texly_test?sslmode=disable"
	}

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL test database: %v\nMake sure PostgreSQL is running: docker-compose up -d postgres", err)
	}

	// Enable pgvector extension (only if not already enabled)
	var extExists bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'vector')").Scan(&extExists)
	if !extExists {
		if err := db.Exec("CREATE EXTENSION vector").Error; err != nil {
			log.Printf("Warning: Could not enable pgvector extension: %v", err)
		}
	}

	// Clean all tables AND enum types to ensure clean state
	// PostgreSQL creates custom enum types that persist even after DROP TABLE
	// Drop enum types first to prevent "duplicate key value violates unique constraint" errors
	enumTypes := []string{"source_type", "source_status"}
	for _, enumType := range enumTypes {
		db.Exec(fmt.Sprintf("DROP TYPE IF EXISTS %s CASCADE", enumType))
	}

	// Drop tables in reverse dependency order to avoid foreign key issues
	// document_chunks depends on sources, messages/sources depend on bots, bots depends on users
	tables := []string{"document_chunks", "messages", "sources", "bots", "users"}
	for _, table := range tables {
		db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
	}

	// Run migrations in proper dependency order
	// DocumentChunk must come after Source because it has a foreign key reference
	if err := db.AutoMigrate(
		&models.User{},
		&models.Bot{},
		&models.Source{},
		&models.Message{},
		&models.DocumentChunk{},
	); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func GetTestConfig() configs.Config {
	return configs.Config{
		DatabaseURL:          "postgres://texly:texly_dev@localhost:5432/texly_test?sslmode=disable",
		DatabaseMaxConns:     10,
		DatabaseMaxIdleConns: 2,
		Port:                 "8080",
		JWTSecret:            "testsecret",
		EmbeddingDimension:   1536,
	}
}

/*
 * CleanupTestDB cleans up the test database after tests
 */
func CleanupTestDB(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

/*
 * TruncateTable truncates a specific table for cleanup between tests
 */
func TruncateTable(db *gorm.DB, tableName string) error {
	return db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)).Error
}
