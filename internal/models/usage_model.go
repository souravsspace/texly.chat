package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
 * UsageRecord tracks a single billable event
 */
type UsageRecord struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	BotID     string    `json:"bot_id" gorm:"index"`
	Type      string    `json:"type"`     // "chat_message", "embedding", "storage", "extra_bot"
	Quantity  float64   `json:"quantity"` // Amount used (e.g., 1 message, 1000 tokens, 0.5 GB)
	Cost      float64   `json:"cost"`     // Cost in USD
	BilledAt  time.Time `json:"billed_at"`
	CreatedAt time.Time `json:"created_at"`
}

/*
 * BeforeCreate generates a new UUID for the usage record
 */
func (u *UsageRecord) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}

// Usage Types
const (
	UsageTypeChatMessage = "chat_message"
	UsageTypeEmbedding   = "embedding"
	UsageTypeStorage     = "storage"
	UsageTypeExtraBot    = "extra_bot"
)
