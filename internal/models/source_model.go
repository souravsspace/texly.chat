package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
* SourceStatus represents the processing status of a source
 */
type SourceStatus string

const (
	SourceStatusPending    SourceStatus = "pending"
	SourceStatusProcessing SourceStatus = "processing"
	SourceStatusCompleted  SourceStatus = "completed"
	SourceStatusFailed     SourceStatus = "failed"
)

/*
* Source represents a data source (URL) for a bot
 */
type Source struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	BotID        string         `json:"bot_id" gorm:"not null;index"`
	URL          string         `json:"url" gorm:"not null"`
	Status       SourceStatus   `json:"status" gorm:"not null;default:'pending'"`
	ErrorMessage string         `json:"error_message"`
	ProcessedAt  *time.Time     `json:"processed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

/*
* BeforeCreate generates a new UUID for the source
 */
func (s *Source) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	if s.Status == "" {
		s.Status = SourceStatusPending
	}
	return
}

/*
* CreateSourceRequest holds data for creating a new source
 */
type CreateSourceRequest struct {
	URL string `json:"url" binding:"required,url"`
}
