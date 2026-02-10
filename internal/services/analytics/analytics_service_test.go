package analytics

import (
	"context"
	"testing"

	"github.com/souravsspace/texly.chat/internal/models"
	messageRepo "github.com/souravsspace/texly.chat/internal/repo/message"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * Test NewAnalyticsService initialization
 */
func TestNewAnalyticsService(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := NewAnalyticsService(msgRepo)

	assert.NotNil(t, service)
	assert.NotNil(t, service.messageRepo)
}

/*
 * Test GetBotAnalytics
 */
func TestGetBotAnalytics(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := NewAnalyticsService(msgRepo)
	ctx := context.Background()

	// Create a bot and messages
	bot := &models.Bot{
		UserID:       "user-123",
		Name:         "Test Bot",
		SystemPrompt: "You are helpful",
	}
	err := db.Create(bot).Error
	require.NoError(t, err)

	// Create some messages
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
			Content:    "Hi there!",
			TokenCount: 10,
		},
	}

	for i := range messages {
		err := msgRepo.Create(ctx, &messages[i])
		require.NoError(t, err)
	}

	// Get analytics
	analytics, err := service.GetBotAnalytics(ctx, bot.ID)
	require.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.Equal(t, bot.ID, analytics.BotID)
	assert.GreaterOrEqual(t, analytics.TotalMessages, 2)
	assert.GreaterOrEqual(t, analytics.TotalTokens, 15)
}

/*
 * Test GetBotDailyStats
 */
func TestGetBotDailyStats(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := NewAnalyticsService(msgRepo)
	ctx := context.Background()

	// Create a bot and messages
	bot := &models.Bot{
		UserID:       "user-456",
		Name:         "Stats Bot",
		SystemPrompt: "You are helpful",
	}
	err := db.Create(bot).Error
	require.NoError(t, err)

	// Create messages
	userID := "user-456"
	messages := []models.Message{
		{
			SessionID:  "session-stats",
			BotID:      bot.ID,
			UserID:     &userID,
			Role:       "user",
			Content:    "Test message",
			TokenCount: 5,
		},
		{
			SessionID:  "session-stats",
			BotID:      bot.ID,
			UserID:     &userID,
			Role:       "assistant",
			Content:    "Test response",
			TokenCount: 10,
		},
	}

	for i := range messages {
		err := msgRepo.Create(ctx, &messages[i])
		require.NoError(t, err)
	}

	// Get daily stats for last 7 days
	stats, err := service.GetBotDailyStats(ctx, bot.ID, 7)
	require.NoError(t, err)
	assert.NotNil(t, stats)
	// Stats should have at least today's data if messages were created
	if len(stats) > 0 {
		assert.GreaterOrEqual(t, stats[0].MessageCount, 1)
	}
}

/*
 * Test GetUserAnalytics
 */
func TestGetUserAnalytics(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := NewAnalyticsService(msgRepo)
	ctx := context.Background()

	userID := "user-analytics-789"

	// Create multiple bots for the user
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

	for i := range bots {
		err := db.Create(&bots[i]).Error
		require.NoError(t, err)

		// Create a message for each bot
		message := &models.Message{
			SessionID:  "session-" + bots[i].ID,
			BotID:      bots[i].ID,
			UserID:     &userID,
			Role:       "user",
			Content:    "Hello bot",
			TokenCount: 5,
		}
		err = msgRepo.Create(ctx, message)
		require.NoError(t, err)
	}

	// Get user analytics
	analytics, err := service.GetUserAnalytics(ctx, userID)
	require.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.GreaterOrEqual(t, len(analytics), 2)

	// Verify analytics for each bot
	for _, bot := range bots {
		botAnalytics, exists := analytics[bot.ID]
		assert.True(t, exists)
		assert.GreaterOrEqual(t, botAnalytics.TotalMessages, 1)
	}
}

/*
 * Test GetSessionMessages
 */
func TestGetSessionMessages(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := NewAnalyticsService(msgRepo)
	ctx := context.Background()

	sessionID := "session-test-messages"
	botID := "bot-test"

	// Create messages
	messages := []models.Message{
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "user",
			Content:    "First message",
			TokenCount: 3,
		},
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "assistant",
			Content:    "First response",
			TokenCount: 4,
		},
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "user",
			Content:    "Second message",
			TokenCount: 3,
		},
	}

	for i := range messages {
		err := msgRepo.Create(ctx, &messages[i])
		require.NoError(t, err)
	}

	// Get session messages
	retrieved, err := service.GetSessionMessages(ctx, sessionID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(retrieved), 3)

	// Verify order
	assert.Equal(t, "First message", retrieved[0].Content)
	assert.Equal(t, "First response", retrieved[1].Content)
	assert.Equal(t, "Second message", retrieved[2].Content)
}

/*
 * Test GetBotAnalytics with no data
 */
func TestGetBotAnalytics_NoData(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := NewAnalyticsService(msgRepo)
	ctx := context.Background()

	analytics, err := service.GetBotAnalytics(ctx, "non-existent-bot")
	require.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.Equal(t, 0, analytics.TotalMessages)
}

/*
 * Test GetUserAnalytics with no bots
 */
func TestGetUserAnalytics_NoBots(t *testing.T) {
	db := shared.SetupTestDB()
	msgRepo := messageRepo.New(db)
	service := NewAnalyticsService(msgRepo)
	ctx := context.Background()

	analytics, err := service.GetUserAnalytics(ctx, "user-no-bots")
	require.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.Equal(t, 0, len(analytics))
}
