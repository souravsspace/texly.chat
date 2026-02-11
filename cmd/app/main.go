package main

import (
	"log"

	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/db"
	"github.com/souravsspace/texly.chat/internal/server"
)

/*
* main is the entry point of the application
 */
func main() {
	cfg := configs.Load()

	gormDb, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// Configure connection pool from config
	sqlDB, err := gormDb.DB()
	if err != nil {
		log.Fatalf("failed to get underlying DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(cfg.DatabaseMaxConns)
	sqlDB.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)

	if err := db.Migrate(gormDb); err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}

	srv := server.New(gormDb, cfg)
	log.Printf("🚀 Server starting on port %s\n", cfg.Port)
	if err := srv.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
