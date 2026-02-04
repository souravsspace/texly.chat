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

func TestCreateBot_WithWidgetConfig(t *testing.T) {
	db := setupTestDB()
	r := setupRouter(db)

	widgetConfig := &models.WidgetConfig{
		ThemeColor:     "#ff5733",
		InitialMessage: "Welcome to our bot!",
		Position:       "bottom-left",
		BotAvatar:      "https://example.com/avatar.png",
	}

	reqBody := models.CreateBotRequest{
		Name:           "Widget Bot",
		SystemPrompt:   "You are a helpful assistant",
		AllowedOrigins: []string{"https://example.com", "https://app.example.com"},
		WidgetConfig:   widgetConfig,
	}
	jsonValue, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bots", bytes.NewBuffer(jsonValue))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdBot models.Bot
	json.Unmarshal(w.Body.Bytes(), &createdBot)
	assert.Equal(t, "Widget Bot", createdBot.Name)
	assert.NotEmpty(t, createdBot.AllowedOrigins)
	assert.NotEmpty(t, createdBot.WidgetConfig)

	// Verify JSON parsing
	var allowedOrigins []string
	json.Unmarshal([]byte(createdBot.AllowedOrigins), &allowedOrigins)
	assert.Len(t, allowedOrigins, 2)
	assert.Contains(t, allowedOrigins, "https://example.com")

	var parsedConfig models.WidgetConfig
	json.Unmarshal([]byte(createdBot.WidgetConfig), &parsedConfig)
	assert.Equal(t, "#ff5733", parsedConfig.ThemeColor)
	assert.Equal(t, "Welcome to our bot!", parsedConfig.InitialMessage)
}

func TestUpdateBot_WithWidgetConfig(t *testing.T) {
	db := setupTestDB()
	r := setupRouter(db)

	// Create initial bot
	botInstance := models.Bot{UserID: "test-user-id", Name: "Test Bot"}
	db.Create(&botInstance)

	// Update with widget config
	widgetConfig := &models.WidgetConfig{
		ThemeColor:     "#6366f1",
		InitialMessage: "Updated message",
		Position:       "bottom-right",
		BotAvatar:      "",
	}

	updateReq := models.UpdateBotRequest{
		Name:           "Updated Bot",
		AllowedOrigins: []string{"https://newdomain.com"},
		WidgetConfig:   widgetConfig,
	}
	jsonValue, _ := json.Marshal(updateReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/bots/"+botInstance.ID, bytes.NewBuffer(jsonValue))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedBot models.Bot
	db.First(&updatedBot, "id = ?", botInstance.ID)
	assert.Equal(t, "Updated Bot", updatedBot.Name)
	assert.NotEmpty(t, updatedBot.AllowedOrigins)
	assert.NotEmpty(t, updatedBot.WidgetConfig)
}
