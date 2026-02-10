package analytics

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/models"
	messageRepo "github.com/souravsspace/texly.chat/internal/repo/message"
	"github.com/souravsspace/texly.chat/internal/services/analytics"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * Helper function to setup test router with auth middleware mock
 */
func setupTestRouter(handler *AnalyticsHandler, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})

	return router
}

/*
 * Test GetBotAnalytics handler
 */
func TestGetBotAnalytics_Handler(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := analytics.NewAnalyticsService(msgRepo)
	handler := NewAnalyticsHandler(service)

	userID := "user-123"
	router := setupTestRouter(handler, userID)
	router.GET("/api/analytics/bots/:id", handler.GetBotAnalytics)

	// Create a bot and messages
	bot := &models.Bot{
		UserID:       userID,
		Name:         "Test Bot",
		SystemPrompt: "You are helpful",
	}
	err := db.Create(bot).Error
	require.NoError(t, err)

	// Create messages
	ctx := context.Background()
	messages := []models.Message{
		{
			SessionID:  "session-1",
			BotID:      bot.ID,
			Role:       "user",
			Content:    "Hello",
			TokenCount: 5,
		},
		{
			SessionID:  "session-1",
			BotID:      bot.ID,
			Role:       "assistant",
			Content:    "Hi!",
			TokenCount: 10,
		},
	}

	for i := range messages {
		err := msgRepo.Create(ctx, &messages[i])
		require.NoError(t, err)
	}

	// Make request
	req := httptest.NewRequest("GET", "/api/analytics/bots/"+bot.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.BotAnalytics
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, bot.ID, response.BotID)
	assert.GreaterOrEqual(t, response.TotalMessages, 2)
}

/*
 * Test GetBotDailyStats handler
 */
func TestGetBotDailyStats_Handler(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := analytics.NewAnalyticsService(msgRepo)
	handler := NewAnalyticsHandler(service)

	userID := "user-456"
	router := setupTestRouter(handler, userID)
	router.GET("/api/analytics/bots/:id/daily", handler.GetBotDailyStats)

	// Create a bot
	bot := &models.Bot{
		UserID:       userID,
		Name:         "Daily Stats Bot",
		SystemPrompt: "You are helpful",
	}
	err := db.Create(bot).Error
	require.NoError(t, err)

	// Create messages
	ctx := context.Background()
	message := &models.Message{
		SessionID:  "session-daily",
		BotID:      bot.ID,
		UserID:     &userID,
		Role:       "user",
		Content:    "Test",
		TokenCount: 5,
	}
	err = msgRepo.Create(ctx, message)
	require.NoError(t, err)

	// Make request with default days (30)
	req := httptest.NewRequest("GET", "/api/analytics/bots/"+bot.ID+"/daily", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.MessageStats
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response)
}

/*
 * Test GetBotDailyStats with custom days parameter
 */
func TestGetBotDailyStats_CustomDays(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := analytics.NewAnalyticsService(msgRepo)
	handler := NewAnalyticsHandler(service)

	userID := "user-789"
	router := setupTestRouter(handler, userID)
	router.GET("/api/analytics/bots/:id/daily", handler.GetBotDailyStats)

	// Create a bot
	bot := &models.Bot{
		UserID:       userID,
		Name:         "Custom Days Bot",
		SystemPrompt: "You are helpful",
	}
	err := db.Create(bot).Error
	require.NoError(t, err)

	// Make request with custom days parameter
	req := httptest.NewRequest("GET", "/api/analytics/bots/"+bot.ID+"/daily?days=7", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.MessageStats
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	// Response might be nil or empty slice - both are valid for no data
	if response != nil {
		assert.Equal(t, 0, len(response))
	}
}

/*
 * Test GetUserAnalytics handler
 */
func TestGetUserAnalytics_Handler(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := analytics.NewAnalyticsService(msgRepo)
	handler := NewAnalyticsHandler(service)

	userID := "user-all-analytics"
	router := setupTestRouter(handler, userID)
	router.GET("/api/analytics/user", handler.GetUserAnalytics)

	// Create bots
	bots := []models.Bot{
		{
			UserID:       userID,
			Name:         "Bot 1",
			SystemPrompt: "You are helpful",
		},
		{
			UserID:       userID,
			Name:         "Bot 2",
			SystemPrompt: "You are friendly",
		},
	}

	ctx := context.Background()
	for i := range bots {
		err := db.Create(&bots[i]).Error
		require.NoError(t, err)

		// Create a message for each bot
		message := &models.Message{
			SessionID:  "session-" + bots[i].ID,
			BotID:      bots[i].ID,
			UserID:     &userID,
			Role:       "user",
			Content:    "Hello",
			TokenCount: 5,
		}
		err = msgRepo.Create(ctx, message)
		require.NoError(t, err)
	}

	// Make request
	req := httptest.NewRequest("GET", "/api/analytics/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]*models.BotAnalytics
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(response), 2)
}

/*
 * Test GetSessionMessages handler
 */
func TestGetSessionMessages_Handler(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := analytics.NewAnalyticsService(msgRepo)
	handler := NewAnalyticsHandler(service)

	userID := "user-session-msgs"
	router := setupTestRouter(handler, userID)
	router.GET("/api/analytics/sessions/:id/messages", handler.GetSessionMessages)

	sessionID := "session-messages-test"
	botID := "bot-test"

	// Create messages
	ctx := context.Background()
	messages := []models.Message{
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "user",
			Content:    "First",
			TokenCount: 2,
		},
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "assistant",
			Content:    "Second",
			TokenCount: 2,
		},
	}

	for i := range messages {
		err := msgRepo.Create(ctx, &messages[i])
		require.NoError(t, err)
	}

	// Make request
	req := httptest.NewRequest("GET", "/api/analytics/sessions/"+sessionID+"/messages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(response), 2)
	assert.Equal(t, "First", response[0].Content)
}

/*
 * Test unauthorized access
 */
func TestGetBotAnalytics_Unauthorized(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := analytics.NewAnalyticsService(msgRepo)
	handler := NewAnalyticsHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// No auth middleware - user_id will be empty
	router.GET("/api/analytics/bots/:id", handler.GetBotAnalytics)

	req := httptest.NewRequest("GET", "/api/analytics/bots/some-bot", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

/*
 * Test missing bot ID
 */
func TestGetBotAnalytics_MissingBotID(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := analytics.NewAnalyticsService(msgRepo)
	handler := NewAnalyticsHandler(service)

	userID := "user-test"
	router := setupTestRouter(handler, userID)
	router.GET("/api/analytics/bots/:id", handler.GetBotAnalytics)

	// Request with empty bot ID
	req := httptest.NewRequest("GET", "/api/analytics/bots/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 404 (route not found)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

/*
 * Test invalid days parameter
 */
func TestGetBotDailyStats_InvalidDays(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := analytics.NewAnalyticsService(msgRepo)
	handler := NewAnalyticsHandler(service)

	userID := "user-invalid-days"
	router := setupTestRouter(handler, userID)
	router.GET("/api/analytics/bots/:id/daily", handler.GetBotDailyStats)

	bot := &models.Bot{
		UserID:       userID,
		Name:         "Test Bot",
		SystemPrompt: "You are helpful",
	}
	err := db.Create(bot).Error
	require.NoError(t, err)

	// Request with invalid days parameter (should default to 30)
	req := httptest.NewRequest("GET", "/api/analytics/bots/"+bot.ID+"/daily?days=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should still return 200 with default days
	assert.Equal(t, http.StatusOK, w.Code)
}

/*
 * Test days parameter cap at 365
 */
func TestGetBotDailyStats_DaysCap(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := analytics.NewAnalyticsService(msgRepo)
	handler := NewAnalyticsHandler(service)

	userID := "user-days-cap"
	router := setupTestRouter(handler, userID)
	router.GET("/api/analytics/bots/:id/daily", handler.GetBotDailyStats)

	bot := &models.Bot{
		UserID:       userID,
		Name:         "Test Bot",
		SystemPrompt: "You are helpful",
	}
	err := db.Create(bot).Error
	require.NoError(t, err)

	// Request with days > 365 (should cap at 365)
	req := httptest.NewRequest("GET", "/api/analytics/bots/"+bot.ID+"/daily?days=1000", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 200 with capped days
	assert.Equal(t, http.StatusOK, w.Code)
}
