package message

import (
	"context"
	"fmt"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/gorm"
)

/*
 * MessageRepository handles database operations for messages
 */
type MessageRepository struct {
	db *gorm.DB
}

/*
 * New creates a new MessageRepository instance
 */
func New(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

/*
 * Create saves a new message to the database
 */
func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {
	if err := r.db.WithContext(ctx).Create(message).Error; err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

/*
 * GetBySessionID retrieves all messages for a session
 */
func (r *MessageRepository) GetBySessionID(ctx context.Context, sessionID string) ([]models.Message, error) {
	var messages []models.Message
	if err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages by session: %w", err)
	}
	return messages, nil
}

/*
 * GetByBotID retrieves all messages for a bot
 */
func (r *MessageRepository) GetByBotID(ctx context.Context, botID string) ([]models.Message, error) {
	var messages []models.Message
	if err := r.db.WithContext(ctx).
		Where("bot_id = ?", botID).
		Order("created_at DESC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages by bot: %w", err)
	}
	return messages, nil
}

/*
 * GetDailyStats retrieves daily message statistics for a bot within a date range
 */
func (r *MessageRepository) GetDailyStats(ctx context.Context, botID string, startDate, endDate time.Time) ([]models.MessageStats, error) {
	var stats []models.MessageStats

	query := `
		SELECT 
			bot_id,
			TO_CHAR(DATE(created_at), 'YYYY-MM-DD') as date,
			COUNT(*) as message_count,
			SUM(token_count) as total_tokens,
			SUM(CASE WHEN role = 'user' THEN 1 ELSE 0 END) as user_messages,
			SUM(CASE WHEN role = 'assistant' THEN 1 ELSE 0 END) as bot_messages,
			COUNT(DISTINCT session_id) as unique_session
		FROM messages
		WHERE bot_id = $1
		AND created_at BETWEEN $2 AND $3
		AND deleted_at IS NULL
		GROUP BY bot_id, DATE(created_at)
		ORDER BY date DESC
	`

	type rawStats struct {
		BotID         string
		Date          string
		MessageCount  int
		TotalTokens   int
		UserMessages  int
		BotMessages   int
		UniqueSession int
	}

	var rawResults []rawStats
	if err := r.db.WithContext(ctx).Raw(query, botID, startDate, endDate).Scan(&rawResults).Error; err != nil {
		return nil, fmt.Errorf("failed to get daily stats: %w", err)
	}

	// Convert raw results to proper types
	for _, raw := range rawResults {
		parsedDate, err := time.Parse("2006-01-02 15:04:05", raw.Date)
		if err != nil {
			// Try alternative format
			parsedDate, err = time.Parse("2006-01-02", raw.Date)
			if err != nil {
				continue
			}
		}

		stats = append(stats, models.MessageStats{
			BotID:         raw.BotID,
			Date:          parsedDate,
			MessageCount:  raw.MessageCount,
			TotalTokens:   raw.TotalTokens,
			UserMessages:  raw.UserMessages,
			BotMessages:   raw.BotMessages,
			UniqueSession: raw.UniqueSession,
		})
	}

	return stats, nil
}

/*
 * GetBotAnalytics retrieves overall analytics for a bot
 */
func (r *MessageRepository) GetBotAnalytics(ctx context.Context, botID string) (*models.BotAnalytics, error) {
	type rawAnalytics struct {
		BotID              string
		TotalMessages      int
		TotalTokens        int
		TotalSessions      int
		AvgMessagesPerDay  float64
		AvgTokensPerDay    float64
		AvgMessagesSession float64
		LastMessageAt      *string
	}

	var raw rawAnalytics

	query := `
		SELECT 
			bot_id,
			COUNT(*) as total_messages,
			SUM(token_count) as total_tokens,
			COUNT(DISTINCT session_id) as total_sessions,
			CAST(COUNT(*) AS FLOAT) / NULLIF(EXTRACT(EPOCH FROM (MAX(created_at) - MIN(created_at))) / 86400, 0) as avg_messages_per_day,
			CAST(SUM(token_count) AS FLOAT) / NULLIF(EXTRACT(EPOCH FROM (MAX(created_at) - MIN(created_at))) / 86400, 0) as avg_tokens_per_day,
			CAST(COUNT(*) AS FLOAT) / NULLIF(COUNT(DISTINCT session_id), 0) as avg_messages_session,
			TO_CHAR(MAX(created_at), 'YYYY-MM-DD HH24:MI:SS') as last_message_at
		FROM messages
		WHERE bot_id = $1
		AND deleted_at IS NULL
		GROUP BY bot_id
	`

	if err := r.db.WithContext(ctx).Raw(query, botID).Scan(&raw).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return empty analytics if no messages exist
			return &models.BotAnalytics{BotID: botID}, nil
		}
		return nil, fmt.Errorf("failed to get bot analytics: %w", err)
	}

	// If no messages, return empty analytics
	if raw.TotalMessages == 0 {
		return &models.BotAnalytics{BotID: botID}, nil
	}

	analytics := &models.BotAnalytics{
		BotID:              raw.BotID,
		TotalMessages:      raw.TotalMessages,
		TotalTokens:        raw.TotalTokens,
		TotalSessions:      raw.TotalSessions,
		AvgMessagesPerDay:  raw.AvgMessagesPerDay,
		AvgTokensPerDay:    raw.AvgTokensPerDay,
		AvgMessagesSession: raw.AvgMessagesSession,
	}

	// Parse LastMessageAt if available
	if raw.LastMessageAt != nil && *raw.LastMessageAt != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", *raw.LastMessageAt)
		if err == nil {
			analytics.LastMessageAt = &parsedTime
		}
	}

	return analytics, nil
}

/*
 * GetUserAnalytics retrieves analytics for all bots owned by a user
 */
func (r *MessageRepository) GetUserAnalytics(ctx context.Context, userID string) (map[string]*models.BotAnalytics, error) {
	type rawResult struct {
		BotID              string
		TotalMessages      int
		TotalTokens        int
		TotalSessions      int
		AvgMessagesPerDay  float64
		AvgTokensPerDay    float64
		AvgMessagesSession float64
		LastMessageAt      *string
	}

	var results []rawResult

	query := `
		SELECT 
			m.bot_id,
			COUNT(*) as total_messages,
			SUM(m.token_count) as total_tokens,
			COUNT(DISTINCT m.session_id) as total_sessions,
			CAST(COUNT(*) AS FLOAT) / NULLIF(EXTRACT(EPOCH FROM (MAX(m.created_at) - MIN(m.created_at))) / 86400, 0) as avg_messages_per_day,
			CAST(SUM(m.token_count) AS FLOAT) / NULLIF(EXTRACT(EPOCH FROM (MAX(m.created_at) - MIN(m.created_at))) / 86400, 0) as avg_tokens_per_day,
			CAST(COUNT(*) AS FLOAT) / NULLIF(COUNT(DISTINCT m.session_id), 0) as avg_messages_session,
			TO_CHAR(MAX(m.created_at), 'YYYY-MM-DD HH24:MI:SS') as last_message_at
		FROM messages m
		INNER JOIN bots b ON m.bot_id = b.id
		WHERE b.user_id = $1
		AND m.deleted_at IS NULL
		AND b.deleted_at IS NULL
		GROUP BY m.bot_id
	`

	if err := r.db.WithContext(ctx).Raw(query, userID).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get user analytics: %w", err)
	}

	analytics := make(map[string]*models.BotAnalytics)
	for _, result := range results {
		botAnalytics := &models.BotAnalytics{
			BotID:              result.BotID,
			TotalMessages:      result.TotalMessages,
			TotalTokens:        result.TotalTokens,
			TotalSessions:      result.TotalSessions,
			AvgMessagesPerDay:  result.AvgMessagesPerDay,
			AvgTokensPerDay:    result.AvgTokensPerDay,
			AvgMessagesSession: result.AvgMessagesSession,
		}

		// Parse LastMessageAt if available
		if result.LastMessageAt != nil && *result.LastMessageAt != "" {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", *result.LastMessageAt)
			if err == nil {
				botAnalytics.LastMessageAt = &parsedTime
			}
		}

		analytics[result.BotID] = botAnalytics
	}

	return analytics, nil
}

/*
 * CountMessagesByBot counts total messages for a bot
 */
func (r *MessageRepository) CountMessagesByBot(ctx context.Context, botID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Message{}).
		Where("bot_id = ?", botID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}
	return count, nil
}
