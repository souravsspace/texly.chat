package public

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/souravsspace/texly.chat/internal/models"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	"github.com/souravsspace/texly.chat/internal/services/session"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestHandler() (*PublicHandler, *gin.Engine, *botRepo.BotRepo, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	testDB := shared.SetupTestDB()

	repo := botRepo.NewBotRepo(testDB, nil)
	sessionService := session.NewSessionService()

	handler := NewPublicHandler(repo, sessionService, nil) // chatService is nil for these tests
	router := gin.New()

	return handler, router, repo, testDB
}

func TestPublicHandler_GetWidgetConfig(t *testing.T) {
	handler, router, repo, db := setupTestHandler()
	router.GET("/api/public/bots/:id/config", handler.GetWidgetConfig)

	// Create a test bot with widget config
	botID := uuid.New().String()
	widgetConfig := models.WidgetConfig{
		ThemeColor:     "#ff5733",
		InitialMessage: "Hello from widget!",
		Position:       "bottom-left",
		BotAvatar:      "https://example.com/avatar.png",
	}
	widgetConfigJSON, _ := json.Marshal(widgetConfig)

	bot := &models.Bot{
		ID:           botID,
		UserID:       "user-1",
		Name:         "Test Widget Bot",
		SystemPrompt: "Test prompt",
		WidgetConfig: string(widgetConfigJSON),
	}
	repo.Create(bot)
	defer db.Unscoped().Delete(bot)

	// Test successful config retrieval
	req := httptest.NewRequest("GET", "/api/public/bots/"+botID+"/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, botID, response["id"])
	assert.Equal(t, "Test Widget Bot", response["name"])
	assert.NotNil(t, response["widget_config"])
}

func TestPublicHandler_GetWidgetConfig_NotFound(t *testing.T) {
	handler, router, _, _ := setupTestHandler()
	router.GET("/api/public/bots/:id/config", handler.GetWidgetConfig)

	req := httptest.NewRequest("GET", "/api/public/bots/non-existent-id/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPublicHandler_GetWidgetConfig_DefaultConfig(t *testing.T) {
	handler, router, repo, db := setupTestHandler()
	router.GET("/api/public/bots/:id/config", handler.GetWidgetConfig)

	// Create a bot without widget config
	botID := uuid.New().String()
	bot := &models.Bot{
		ID:           botID,
		UserID:       "user-1",
		Name:         "Bot Without Config",
		SystemPrompt: "Test prompt",
	}
	repo.Create(bot)
	defer db.Unscoped().Delete(bot)

	req := httptest.NewRequest("GET", "/api/public/bots/"+botID+"/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	widgetConfig := response["widget_config"].(map[string]interface{})
	assert.Equal(t, "#6366f1", widgetConfig["theme_color"])
	assert.Equal(t, "Hi! How can I help you today?", widgetConfig["initial_message"])
	assert.Equal(t, "bottom-right", widgetConfig["position"])
}

func TestPublicHandler_CreateSession(t *testing.T) {
	handler, router, repo, db := setupTestHandler()
	router.POST("/api/public/chats", handler.CreateSession)

	// Create a test bot
	botID := uuid.New().String()
	bot := &models.Bot{
		ID:           botID,
		UserID:       "user-1",
		Name:         "Test Bot",
		SystemPrompt: "Test prompt",
	}
	repo.Create(bot)
	defer db.Unscoped().Delete(bot)

	// Test session creation
	reqBody := models.CreateSessionRequest{
		BotID: botID,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/public/chats", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.SessionResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotEmpty(t, response.SessionID)
	assert.Equal(t, botID, response.BotID)
	assert.False(t, response.ExpiresAt.IsZero())
}

func TestPublicHandler_CreateSession_BotNotFound(t *testing.T) {
	handler, router, _, _ := setupTestHandler()
	router.POST("/api/public/chats", handler.CreateSession)

	reqBody := models.CreateSessionRequest{
		BotID: "non-existent-bot",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/public/chats", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPublicHandler_CreateSession_InvalidRequest(t *testing.T) {
	handler, router, _, _ := setupTestHandler()
	router.POST("/api/public/chats", handler.CreateSession)

	// Send invalid JSON
	req := httptest.NewRequest("POST", "/api/public/chats", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPublicHandler_StreamChatPublic_SessionNotFound(t *testing.T) {
	handler, router, _, _ := setupTestHandler()
	router.POST("/api/public/chats/:session_id/messages", handler.StreamChatPublic)

	reqBody := models.ChatRequest{
		Message: "Hello",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/public/chats/non-existent-session/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPublicHandler_StreamChatPublic_InvalidRequest(t *testing.T) {
	handler, router, repo, db := setupTestHandler()
	router.POST("/api/public/chats/:session_id/messages", handler.StreamChatPublic)

	// Create bot and session
	botID := uuid.New().String()
	bot := &models.Bot{
		ID:           botID,
		UserID:       "user-1",
		Name:         "Test Bot",
		SystemPrompt: "Test prompt",
	}
	repo.Create(bot)
	defer db.Unscoped().Delete(bot)

	chatSession := handler.sessionService.CreateSession(botID)

	// Send invalid JSON
	req := httptest.NewRequest("POST", "/api/public/chats/"+chatSession.ID+"/messages", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
