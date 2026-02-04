package models

import (
	"time"

	"github.com/google/uuid"
)

/*
* ChatSession represents an anonymous user session for the widget
 */
type ChatSession struct {
	ID             string    `json:"id"`
	BotID          string    `json:"bot_id"`
	CreatedAt      time.Time `json:"created_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

/*
* CreateSessionRequest holds data for creating a new chat session
 */
type CreateSessionRequest struct {
	BotID string `json:"bot_id" binding:"required"`
}

/*
* SessionResponse represents the response when creating a session
 */
type SessionResponse struct {
	SessionID string    `json:"session_id"`
	BotID     string    `json:"bot_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

/*
* NewChatSession creates a new chat session with a 24-hour expiration
 */
func NewChatSession(botID string) *ChatSession {
	now := time.Now()
	return &ChatSession{
		ID:             uuid.New().String(),
		BotID:          botID,
		CreatedAt:      now,
		LastActivityAt: now,
		ExpiresAt:      now.Add(24 * time.Hour),
	}
}

/*
* IsExpired checks if the session has expired
 */
func (s *ChatSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

/*
* UpdateActivity updates the last activity timestamp
 */
func (s *ChatSession) UpdateActivity() {
	s.LastActivityAt = time.Now()
}
