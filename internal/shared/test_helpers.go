package shared

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB() *gorm.DB {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Bot{}, &models.Source{}, &models.DocumentChunk{}, &models.Message{}); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func GetTestConfig() configs.Config {
	return configs.Config{
		Port:      "8080",
		JwtSecret: "testsecret",
	}
}
