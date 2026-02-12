package chat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/models"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * Setup test environment
 */
func setupTest(t *testing.T) (*gin.Engine, *ChatHandler, *botRepo.BotRepo) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := shared.SetupTestDB()

	// Create bot repo with nil cache (pass-through)
	repo := botRepo.NewBotRepo(db, nil)

	// Create handler with repo injected
	handler := NewChatHandler(repo, nil) // nil chat service for basic tests

	// Create router
	router := gin.New()

	return router, handler, repo
}

/*
 * Test StreamChat requires authentication
 */
func TestStreamChat_RequiresAuth(t *testing.T) {
	router, handler, _ := setupTest(t)

	router.POST("/api/bots/:id/chat", handler.StreamChat)

	reqBody := models.ChatRequest{
		Message: "Hello",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/bots/bot-123/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

/*
 * Test StreamChat requires valid bot ID
 */
func TestStreamChat_RequiresBotID(t *testing.T) {
	router, handler, _ := setupTest(t)

	router.POST("/api/bots/:id/chat", func(c *gin.Context) {
		c.Set("user_id", "user-123") // Mock auth
		handler.StreamChat(c)
	})

	reqBody := models.ChatRequest{
		Message: "Hello",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/bots//chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail due to empty botId
	assert.NotEqual(t, http.StatusOK, w.Code)
}

/*
 * Test StreamChat validates request body
 */
func TestStreamChat_ValidatesRequestBody(t *testing.T) {
	router, handler, _ := setupTest(t)

	router.POST("/api/bots/:id/chat", func(c *gin.Context) {
		c.Set("user_id", "user-123") // Mock auth
		handler.StreamChat(c)
	})

	// Invalid JSON
	req := httptest.NewRequest("POST", "/api/bots/bot-123/chat", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

/*
 * Test StreamChat requires message in request
 */
func TestStreamChat_RequiresMessage(t *testing.T) {
	router, handler, _ := setupTest(t)

	router.POST("/api/bots/:id/chat", func(c *gin.Context) {
		c.Set("user_id", "user-123") // Mock auth
		handler.StreamChat(c)
	})

	// Empty message
	reqBody := models.ChatRequest{
		Message: "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/bots/bot-123/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

/*
 * Test StreamChat validates bot ownership
 */
func TestStreamChat_ValidatesBotOwnership(t *testing.T) {
	router, handler, repo := setupTest(t)

	// Create a bot owned by different user
	bot := &models.Bot{
		UserID:       "other-user",
		Name:         "Test Bot",
		SystemPrompt: "You are helpful",
	}
	err := repo.Create(bot)
	require.NoError(t, err)

	router.POST("/api/bots/:id/chat", func(c *gin.Context) {
		c.Set("user_id", "user-123") // Different user
		handler.StreamChat(c)
	})

	reqBody := models.ChatRequest{
		Message: "Hello",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/bots/"+bot.ID+"/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 404 (bot not found for this user)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

/*
 * Test StreamChat returns error when chat service unavailable
 */
func TestStreamChat_ChatServiceUnavailable(t *testing.T) {
	router, handler, repo := setupTest(t)

	// Create a bot
	bot := &models.Bot{
		UserID:       "user-123",
		Name:         "Test Bot",
		SystemPrompt: "You are helpful",
	}
	err := repo.Create(bot)
	require.NoError(t, err)

	router.POST("/api/bots/:id/chat", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		handler.StreamChat(c)
	})

	reqBody := models.ChatRequest{
		Message: "Hello",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/bots/"+bot.ID+"/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return service unavailable (chat service is nil)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

/*
 * Test StreamChat with non-existent bot
 */
func TestStreamChat_BotNotFound(t *testing.T) {
	router, handler, _ := setupTest(t)

	router.POST("/api/bots/:id/chat", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		handler.StreamChat(c)
	})

	reqBody := models.ChatRequest{
		Message: "Hello",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/bots/non-existent-bot/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

/*
 * Test ChatHandler initialization
 */
func TestNewChatHandler(t *testing.T) {
	db := shared.SetupTestDB()

	repo := botRepo.NewBotRepo(db, nil)
	handler := NewChatHandler(repo, nil)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.botRepo)
	assert.Nil(t, handler.chatService) // Can be nil
}
