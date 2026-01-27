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
		&models.DocumentChunk{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	// 2. Create Virtual Table for Vector Search (sqlite-vec)
	// We use vec0 virtual table. `embedding float[1536]` matches OpenAI's small embedding model.
	// We map rowid to DocumentChunk.ID (Wait, rowid in virtual tables is Int64, our IDs are UUID strings).
	// Strategy: We will store the UUID in an auxiliary column or just map a numeric ID.
	// Simple approach: vec0 tables have a rowid. We can add a `vec_id` integer column to DocumentChunk.
	// For now, let's create the table with an explicit string ID if supported, or just standard rowid.
	// sqlite-vec supports "rowid" (int64).
	createVirtualTableQuery := `
	CREATE VIRTUAL TABLE IF NOT EXISTS vec_items USING vec0(
		embedding float[1536]
	);
	`
	if err := db.Exec(createVirtualTableQuery).Error; err != nil {
		// Log warning but don't fail if extension is missing (it will fail syntax error)
		fmt.Printf("Warning: Failed to create virtual table 'vec_items'. Ensure sqlite-vec is loaded. Error: %v\n", err)
	}

	return nil
}
