package bot_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/handlers/bot"
	"github.com/souravsspace/texly.chat/internal/models"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database and migrates the schema.
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Bot{})
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	repo := botRepo.NewBotRepo(db)
	handler := bot.NewBotHandler(repo)

	// Mock Auth middleware by setting user_id in context
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})

	r.POST("/api/bots", handler.CreateBot)
	r.GET("/api/bots", handler.ListBots)
	r.GET("/api/bots/:id", handler.GetBot)
	r.PUT("/api/bots/:id", handler.UpdateBot)
	r.DELETE("/api/bots/:id", handler.DeleteBot)

	return r
}

func TestCreateBot(t *testing.T) {
	db := setupTestDB()
	r := setupRouter(db)

	reqBody := models.CreateBotRequest{
		Name:         "Test Bot",
		SystemPrompt: "You are a helpful assistant",
	}
	jsonValue, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bots", bytes.NewBuffer(jsonValue))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdBot models.Bot
	json.Unmarshal(w.Body.Bytes(), &createdBot)
	assert.Equal(t, "Test Bot", createdBot.Name)
	assert.Equal(t, "test-user-id", createdBot.UserID)
}

func TestGetBots(t *testing.T) {
	db := setupTestDB()
	r := setupRouter(db)

	// Seed db
	db.Create(&models.Bot{UserID: "test-user-id", Name: "Bot 1"})
	db.Create(&models.Bot{UserID: "other-user", Name: "Bot 2"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/bots", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var bots []models.Bot
	json.Unmarshal(w.Body.Bytes(), &bots)
	assert.Len(t, bots, 1)
	assert.Equal(t, "Bot 1", bots[0].Name)
}

func TestUpdateBot(t *testing.T) {
	db := setupTestDB()
	r := setupRouter(db)

	// Seed db
	botInstance := models.Bot{UserID: "test-user-id", Name: "Old Name"}
	db.Create(&botInstance)

	updateReq := models.UpdateBotRequest{Name: "New Name"}
	jsonValue, _ := json.Marshal(updateReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/bots/"+botInstance.ID, bytes.NewBuffer(jsonValue))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var updatedBot models.Bot
	db.First(&updatedBot, "id = ?", botInstance.ID)
	assert.Equal(t, "New Name", updatedBot.Name)
}

func TestDeleteBot(t *testing.T) {
	db := setupTestDB()
	r := setupRouter(db)

	botInstance := models.Bot{UserID: "test-user-id", Name: "To Delete"}
	db.Create(&botInstance)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/bots/"+botInstance.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var count int64
	// Check for soft delete logic. Gorm Delete marks DeletedAt if model has gorm.DeletedAt.
	// Find will ignore deleted records by default.
	db.Model(&models.Bot{}).Where("id = ?", botInstance.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
