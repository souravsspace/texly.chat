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

	return db, nil
}

/*
* Migrate applies database migrations
*/
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Post{})
}
