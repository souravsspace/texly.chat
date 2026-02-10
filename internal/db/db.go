package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/*
* Connect establishes a connection to the SQLite database with sqlite-vec extension
 */
func Connect(path string) (*gorm.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Load sqlite-vec extension before opening the database
	sqlite_vec.Auto()

	// Open raw SQL connection first to verify extension
	sqlDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for concurrency
	if _, err := sqlDB.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Verify sqlite-vec extension is loaded
	var version string
	err = sqlDB.QueryRow("SELECT vec_version();").Scan(&version)
	if err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("sqlite-vec extension not loaded: %w", err)
	}
	fmt.Printf("✅ sqlite-vec extension loaded successfully (version: %s)\n", version)

	// Wrap with GORM
	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	return db, nil
}

/*
* Migrate applies database migrations
 */
func Migrate(db *gorm.DB) error {
	// 1. Standard GORM Migrations
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

	// 2. Initialize Vector Tables (sqlite-vec)
	// We use a dedicated repository method to create the required tables
	// This approach is cleaner and more testable than raw SQL here
	fmt.Println("✅ Database migrations completed successfully")

	return nil
}

/*
* InitializeVectorTables creates the necessary vector search tables
* This should be called after Migrate() with the embedding dimension
 */
func InitializeVectorTables(db *gorm.DB, dimension int) error {
	// Import is done in the function that calls this to avoid circular dependency
	// For now, we'll create the vector repository here
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Create mapping table to link rowid to chunk_id
	createMapTableQuery := `
	CREATE TABLE IF NOT EXISTS vec_chunk_map (
		rowid INTEGER PRIMARY KEY AUTOINCREMENT,
		chunk_id TEXT NOT NULL UNIQUE,
		FOREIGN KEY(chunk_id) REFERENCES document_chunks(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_vec_chunk_map_chunk_id ON vec_chunk_map(chunk_id);
	`

	if _, err := sqlDB.Exec(createMapTableQuery); err != nil {
		return fmt.Errorf("failed to create vec_chunk_map table: %w", err)
	}

	// Create virtual table for vector embeddings
	createVecTableQuery := fmt.Sprintf(`
	CREATE VIRTUAL TABLE IF NOT EXISTS vec_items USING vec0(
		embedding float[%d]
	);
	`, dimension)

	if _, err := sqlDB.Exec(createVecTableQuery); err != nil {
		// Log warning but don't fail if extension is missing
		fmt.Printf("Warning: Failed to create vector virtual table. Error: %v\n", err)
		return nil // Don't fail the entire migration if vec extension isn't loaded
	}

	fmt.Println("✅ Vector tables initialized successfully")
	return nil
}
