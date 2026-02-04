package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
* WidgetConfig holds configuration for the embeddable widget
 */
type WidgetConfig struct {
	ThemeColor     string `json:"theme_color"`     // Primary color for widget UI (e.g., "#6366f1")
	InitialMessage string `json:"initial_message"` // Welcome message shown when widget opens
	Position       string `json:"position"`        // Widget position: "bottom-right" or "bottom-left"
	BotAvatar      string `json:"bot_avatar"`      // URL to bot avatar image
}

/*
* Bot represents a user's chatbot
 */
type Bot struct {
	ID             string         `json:"id" gorm:"primaryKey"`
	UserID         string         `json:"user_id" gorm:"not null;index"`
	Name           string         `json:"name" gorm:"not null"`
	SystemPrompt   string         `json:"system_prompt"`
	AllowedOrigins string         `json:"allowed_origins" gorm:"type:text"` // JSON array of whitelisted domains
	WidgetConfig   string          `json:"widget_config" gorm:"type:text"`   // JSON-encoded WidgetConfig
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

/*
* BeforeCreate generates a new UUID for the bot
 */
func (b *Bot) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}

/*
* CreateBotRequest holds data for creating a new bot
 */
type CreateBotRequest struct {
	Name           string        `json:"name" binding:"required"`
	SystemPrompt   string        `json:"system_prompt"`
	AllowedOrigins []string      `json:"allowed_origins"` // Optional: whitelisted domains for widget
	WidgetConfig   *WidgetConfig   `json:"widget_config"`   // Optional: widget configuration
}

/*
* UpdateBotRequest holds data for updating an existing bot
 */
type UpdateBotRequest struct {
	Name           string        `json:"name"`
	SystemPrompt   string        `json:"system_prompt"`
	AllowedOrigins []string      `json:"allowed_origins"` // Optional: whitelisted domains for widget
	WidgetConfig   *WidgetConfig   `json:"widget_config"`   // Optional: widget configuration
}
