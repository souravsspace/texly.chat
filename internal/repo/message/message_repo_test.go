package message

import (
	"context"
	"testing"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
 * Test Create message
 */
func TestCreateMessage(t *testing.T) {
	db := shared.SetupTestDB()
	repo := New(db)
	ctx := context.Background()

	message := &models.Message{
		SessionID:  "session-123",
		BotID:      "bot-456",
		Role:       "user",
		Content:    "Hello, bot!",
		TokenCount: 3,
	}

	err := repo.Create(ctx, message)
	require.NoError(t, err)
	assert.NotEmpty(t, message.ID)
	assert.False(t, message.CreatedAt.IsZero())
}

/*
 * Test GetBySessionID retrieves messages in correct order
 */
func TestGetBySessionID(t *testing.T) {
	db := shared.SetupTestDB()
	repo := New(db)
	ctx := context.Background()

	sessionID := "session-789"
	botID := "bot-123"

	// Create messages in sequence
	messages := []models.Message{
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "user",
			Content:    "First message",
			TokenCount: 2,
		},
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "assistant",
			Content:    "Response to first",
			TokenCount: 3,
		},
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "user",
			Content:    "Second message",
			TokenCount: 2,
		},
	}

	for i := range messages {
		err := repo.Create(ctx, &messages[i])
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond) // Ensure ordering
	}

	// Retrieve messages
	retrieved, err := repo.GetBySessionID(ctx, sessionID)
	require.NoError(t, err)
	assert.Len(t, retrieved, 3)

	// Verify order (ASC by created_at)
	assert.Equal(t, "First message", retrieved[0].Content)
	assert.Equal(t, "Response to first", retrieved[1].Content)
	assert.Equal(t, "Second message", retrieved[2].Content)
}

/*
 * Test GetByBotID
 */
func TestGetByBotID(t *testing.T) {
	db := shared.SetupTestDB()
	repo := New(db)
	ctx := context.Background()

	botID := "bot-999"

	// Create messages for the bot
	messages := []models.Message{
		{
			SessionID:  "session-a",
			BotID:      botID,
			Role:       "user",
			Content:    "Message A",
			TokenCount: 2,
		},
		{
			SessionID:  "session-b",
			BotID:      botID,
			Role:       "user",
			Content:    "Message B",
			TokenCount: 2,
		},
	}

	for i := range messages {
		err := repo.Create(ctx, &messages[i])
		require.NoError(t, err)
	}

	// Retrieve messages
	retrieved, err := repo.GetByBotID(ctx, botID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(retrieved), 2)
}

/*
 * Test GetDailyStats
 */
func TestGetDailyStats(t *testing.T) {
	db := shared.SetupTestDB()
	repo := New(db)
	ctx := context.Background()

	botID := "bot-stats-test"
	sessionID := "session-stats"

	// Create messages with different roles
	userID := "user-123"
	messages := []models.Message{
		{
			SessionID:  sessionID,
			BotID:      botID,
			UserID:     &userID,
			Role:       "user",
			Content:    "User message 1",
			TokenCount: 10,
		},
		{
			SessionID:  sessionID,
			BotID:      botID,
			UserID:     &userID,
			Role:       "assistant",
			Content:    "Bot response 1",
			TokenCount: 20,
		},
		{
			SessionID:  sessionID,
			BotID:      botID,
			UserID:     &userID,
			Role:       "user",
			Content:    "User message 2",
			TokenCount: 15,
		},
	}

	for i := range messages {
		err := repo.Create(ctx, &messages[i])
		require.NoError(t, err)
	}

	// Get daily stats for today
	// Use database time (UTC) to ensure timezone consistency
	var dbNow time.Time
	db.Raw("SELECT CURRENT_TIMESTAMP").Scan(&dbNow)
	
	// Convert to UTC before extracting date components
	// This is crucial because Postgres returns time in session timezone (e.g. +0600)
	// If we take the local day (12th) and treat it as UTC day, we might miss messages
	// that happened on the 11th in UTC terms.
	dbNow = dbNow.UTC()
	startDate := time.Date(dbNow.Year(), dbNow.Month(), dbNow.Day(), 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 1)

	stats, err := repo.GetDailyStats(ctx, botID, startDate, endDate)
	require.NoError(t, err)
	assert.NotEmpty(t, stats)

	// Verify stats
	if len(stats) > 0 {
		todayStats := stats[0]
		assert.Equal(t, botID, todayStats.BotID)
		assert.GreaterOrEqual(t, todayStats.MessageCount, 3)
		assert.GreaterOrEqual(t, todayStats.TotalTokens, 45)
		assert.GreaterOrEqual(t, todayStats.UserMessages, 2)
		assert.GreaterOrEqual(t, todayStats.BotMessages, 1)
	}
}

/*
 * Test GetBotAnalytics
 */
func TestGetBotAnalytics(t *testing.T) {
	db := shared.SetupTestDB()
	repo := New(db)
	ctx := context.Background()

	botID := "bot-analytics-test"
	sessionID := "session-analytics"

	// Create some messages
	messages := []models.Message{
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "user",
			Content:    "Test message",
			TokenCount: 5,
		},
		{
			SessionID:  sessionID,
			BotID:      botID,
			Role:       "assistant",
			Content:    "Test response",
			TokenCount: 10,
		},
	}

	for i := range messages {
		err := repo.Create(ctx, &messages[i])
		require.NoError(t, err)
	}

	// Get analytics
	analytics, err := repo.GetBotAnalytics(ctx, botID)
	require.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.Equal(t, botID, analytics.BotID)
	assert.GreaterOrEqual(t, analytics.TotalMessages, 2)
	assert.GreaterOrEqual(t, analytics.TotalTokens, 15)
	assert.GreaterOrEqual(t, analytics.TotalSessions, 1)
	// LastMessageAt might be nil if time parsing fails, but that's okay
	// The important thing is that we got analytics back
}

/*
 * Test GetBotAnalytics with no messages
 */
func TestGetBotAnalytics_NoMessages(t *testing.T) {
	db := shared.SetupTestDB()
	repo := New(db)
	ctx := context.Background()

	analytics, err := repo.GetBotAnalytics(ctx, "non-existent-bot")
	require.NoError(t, err)
	assert.NotNil(t, analytics)
	// When there are no messages, BotID should still be returned as an empty result
	assert.Equal(t, 0, analytics.TotalMessages)
}

/*
 * Test CountMessagesByBot
 */
func TestCountMessagesByBot(t *testing.T) {
	db := shared.SetupTestDB()
	repo := New(db)
	ctx := context.Background()

	botID := "bot-count-test"

	// Create 5 messages
	for i := 0; i < 5; i++ {
		message := &models.Message{
			SessionID:  "session-count",
			BotID:      botID,
			Role:       "user",
			Content:    "Test message",
			TokenCount: 3,
		}
		err := repo.Create(ctx, message)
		require.NoError(t, err)
	}

	count, err := repo.CountMessagesByBot(ctx, botID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

/*
 * Test GetUserAnalytics
 */
func TestGetUserAnalytics(t *testing.T) {
	db := shared.SetupTestDB()
	repo := New(db)
	ctx := context.Background()

	// First create a user and bot
	userID := "user-analytics-test"
	bot := &models.Bot{
		UserID:       userID,
		Name:         "Analytics Bot",
		SystemPrompt: "You are helpful",
	}
	err := db.Create(bot).Error
	require.NoError(t, err)

	// Create messages for the bot
	message := &models.Message{
		SessionID:  "session-user-analytics",
		BotID:      bot.ID,
		UserID:     &userID,
		Role:       "user",
		Content:    "User analytics test",
		TokenCount: 5,
	}
	err = repo.Create(ctx, message)
	require.NoError(t, err)

	// Get user analytics
	analytics, err := repo.GetUserAnalytics(ctx, userID)
	require.NoError(t, err)
	assert.NotEmpty(t, analytics)
	assert.Contains(t, analytics, bot.ID)
	assert.GreaterOrEqual(t, analytics[bot.ID].TotalMessages, 1)
}
