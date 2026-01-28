package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/*
* Connect establishes a connection to the SQLite database
 */
func Connect(path string) (*gorm.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	// Enable WAL mode for concurrency
	if err := db.Exec("PRAGMA journal_mode=WAL;").Error; err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Attempt to load sqlite-vec extension
	// Note: This relies on the extension being available in the system path or pre-loaded by the driver.
	// Since we are using standard gorm driver, we might need to handle this more robustly in production.
	// For now, we log if it fails but don't crash, as user might not have it installed yet.
	if err := db.Exec("SELECT vec_version();").Error; err != nil {
		fmt.Printf("Warning: sqlite-vec extension not loaded: %v\n", err)
	} else {
		fmt.Println("✅ sqlite-vec extension loaded successfully")
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
