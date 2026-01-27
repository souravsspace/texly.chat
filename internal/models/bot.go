package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
* Bot represents a user's chatbot
 */
type Bot struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	UserID       string         `json:"user_id" gorm:"not null;index"`
	Name         string         `json:"name" gorm:"not null"`
	SystemPrompt string         `json:"system_prompt"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
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
