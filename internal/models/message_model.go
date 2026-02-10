package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
 * Message represents a chat message exchanged between user and bot
 */
type Message struct {
	ID         string         `json:"id" gorm:"primaryKey"`
	SessionID  string         `json:"session_id" gorm:"not null;index"`
	BotID      string         `json:"bot_id" gorm:"not null;index"`
	UserID     *string        `json:"user_id" gorm:"index"` // Null for widget users
	Role       string         `json:"role" gorm:"not null"` // "user" or "assistant"
	Content    string         `json:"content" gorm:"type:text;not null"`
	TokenCount int            `json:"token_count" gorm:"default:0"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

/*
 * BeforeCreate generates a UUID for the message if not set
 */
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

/*
 * MessageStats represents aggregated statistics for a bot
 */
type MessageStats struct {
	BotID         string    `json:"bot_id"`
	Date          time.Time `json:"date"`
	MessageCount  int       `json:"message_count"`
	TotalTokens   int       `json:"total_tokens"`
	UserMessages  int       `json:"user_messages"`
	BotMessages   int       `json:"bot_messages"`
	UniqueSession int       `json:"unique_sessions"`
}

/*
 * BotAnalytics represents overall analytics for a bot
 */
type BotAnalytics struct {
	BotID              string     `json:"bot_id"`
	TotalMessages      int        `json:"total_messages"`
	TotalTokens        int        `json:"total_tokens"`
	TotalSessions      int        `json:"total_sessions"`
	AvgMessagesPerDay  float64    `json:"avg_messages_per_day"`
	AvgTokensPerDay    float64    `json:"avg_tokens_per_day"`
	AvgMessagesSession float64    `json:"avg_messages_per_session"`
	LastMessageAt      *time.Time `json:"last_message_at"`
}
