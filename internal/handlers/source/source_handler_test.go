package source_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/handlers/source"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/queue"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	sourceRepo "github.com/souravsspace/texly.chat/internal/repo/source"
	"github.com/souravsspace/texly.chat/internal/services/storage"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Bot{}, &models.Source{})
	return db
}

func setupRouter(db *gorm.DB, jobQueue queue.JobQueue) *gin.Engine {
	r := gin.Default()
	sRepo := sourceRepo.NewSourceRepo(db)
	bRepo := botRepo.NewBotRepo(db)

	// Create a mock MinIO storage service (will fail to connect but that's OK for these tests)
	// For URL source tests, we don't actually use storage
	storageService, _ := storage.NewMinIOStorageService(
		"localhost:9000",
		"minioadmin",
		"minioadmin",
		"test-bucket",
		false,
		100,
	)

	handler := source.NewSourceHandler(sRepo, bRepo, jobQueue, storageService, 100)

	// Mock Auth middleware
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})

	r.POST("/api/bots/:id/sources", handler.CreateSource)
	r.GET("/api/bots/:id/sources", handler.ListSources)
	r.GET("/api/bots/:id/sources/:sourceId", handler.GetSource)
	r.DELETE("/api/bots/:id/sources/:sourceId", handler.DeleteSource)

	return r
}

func TestCreateSource(t *testing.T) {
	db := setupTestDB()
	jobQueue := queue.NewInMemoryQueue(10, 1)
	defer jobQueue.Stop()

	r := setupRouter(db, jobQueue)

	// Create a bot first
	bot := &models.Bot{UserID: "test-user-id", Name: "Test Bot"}
	db.Create(bot)

	reqBody := models.CreateSourceRequest{
		URL: "https://example.com",
	}
	jsonValue, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bots/"+bot.ID+"/sources", bytes.NewBuffer(jsonValue))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdSource models.Source
	json.Unmarshal(w.Body.Bytes(), &createdSource)
	assert.Equal(t, "https://example.com", createdSource.URL)
	assert.Equal(t, bot.ID, createdSource.BotID)
	assert.Equal(t, models.SourceStatusPending, createdSource.Status)
}

func TestCreateSource_InvalidBot(t *testing.T) {
	db := setupTestDB()
	jobQueue := queue.NewInMemoryQueue(10, 1)
	defer jobQueue.Stop()

	r := setupRouter(db, jobQueue)

	reqBody := models.CreateSourceRequest{
		URL: "https://example.com",
	}
	jsonValue, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bots/non-existent-bot/sources", bytes.NewBuffer(jsonValue))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code) // Bot not found or not owned returns 403
}

func TestCreateSource_Unauthorized(t *testing.T) {
	db := setupTestDB()
	jobQueue := queue.NewInMemoryQueue(10, 1)
	defer jobQueue.Stop()

	// Create bot owned by different user
	bot := &models.Bot{UserID: "other-user", Name: "Other Bot"}
	db.Create(bot)

	r := setupRouter(db, jobQueue)

	reqBody := models.CreateSourceRequest{
		URL: "https://example.com",
	}
	jsonValue, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bots/"+bot.ID+"/sources", bytes.NewBuffer(jsonValue))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestListSources(t *testing.T) {
	db := setupTestDB()
	jobQueue := queue.NewInMemoryQueue(10, 1)
	defer jobQueue.Stop()

	r := setupRouter(db, jobQueue)

	// Create bot and sources
	bot := &models.Bot{UserID: "test-user-id", Name: "Test Bot"}
	db.Create(bot)

	source1 := &models.Source{BotID: bot.ID, URL: "https://example1.com", Status: models.SourceStatusPending}
	source2 := &models.Source{BotID: bot.ID, URL: "https://example2.com", Status: models.SourceStatusCompleted}
	db.Create(source1)
	db.Create(source2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/bots/"+bot.ID+"/sources", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var sources []models.Source
	json.Unmarshal(w.Body.Bytes(), &sources)
	assert.Len(t, sources, 2)
}

func TestGetSource(t *testing.T) {
	db := setupTestDB()
	jobQueue := queue.NewInMemoryQueue(10, 1)
	defer jobQueue.Stop()

	r := setupRouter(db, jobQueue)

	// Create bot and source
	bot := &models.Bot{UserID: "test-user-id", Name: "Test Bot"}
	db.Create(bot)

	testSource := &models.Source{BotID: bot.ID, URL: "https://example.com", Status: models.SourceStatusPending}
	db.Create(testSource)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/bots/"+bot.ID+"/sources/"+testSource.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var source models.Source
	json.Unmarshal(w.Body.Bytes(), &source)
	assert.Equal(t, testSource.ID, source.ID)
	assert.Equal(t, "https://example.com", source.URL)
}

func TestDeleteSource(t *testing.T) {
	db := setupTestDB()
	jobQueue := queue.NewInMemoryQueue(10, 1)
	defer jobQueue.Stop()

	r := setupRouter(db, jobQueue)

	// Create bot and source
	bot := &models.Bot{UserID: "test-user-id", Name: "Test Bot"}
	db.Create(bot)

	testSource := &models.Source{BotID: bot.ID, URL: "https://example.com", Status: models.SourceStatusPending}
	db.Create(testSource)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/bots/"+bot.ID+"/sources/"+testSource.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify soft delete
	var count int64
	db.Model(&models.Source{}).Where("id = ?", testSource.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteSource_WrongBot(t *testing.T) {
	db := setupTestDB()
	jobQueue := queue.NewInMemoryQueue(10, 1)
	defer jobQueue.Stop()

	r := setupRouter(db, jobQueue)

	// Create two bots
	bot1 := &models.Bot{UserID: "test-user-id", Name: "Bot 1"}
	bot2 := &models.Bot{UserID: "test-user-id", Name: "Bot 2"}
	db.Create(bot1)
	db.Create(bot2)

	// Create source for bot1
	testSource := &models.Source{BotID: bot1.ID, URL: "https://example.com", Status: models.SourceStatusPending}
	db.Create(testSource)

	// Try to delete using bot2's ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/bots/"+bot2.ID+"/sources/"+testSource.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
