package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
* SourceType represents the type of data source
 */
type SourceType string

const (
	SourceTypeURL  SourceType = "url"
	SourceTypeText SourceType = "text"
	SourceTypeFile SourceType = "file"
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
* Source represents a data source for a bot
 */
type Source struct {
	ID                 string         `json:"id" gorm:"primaryKey"`
	BotID              string         `json:"bot_id" gorm:"not null;index"`
	SourceType         SourceType     `json:"source_type" gorm:"not null;default:'url'"`
	URL                string         `json:"url"`
	FilePath           string         `json:"file_path"`           // MinIO object path
	OriginalFilename   string         `json:"original_filename"`   // Original uploaded filename
	ContentType        string         `json:"content_type"`        // MIME type
	Status             SourceStatus   `json:"status" gorm:"not null;default:'pending'"`
	ProcessingProgress int            `json:"processing_progress"` // 0-100
	ErrorMessage       string         `json:"error_message"`
	ProcessedAt        *time.Time     `json:"processed_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

/*
* BeforeCreate generates a new UUID for the source and sets defaults
 */
func (s *Source) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	if s.Status == "" {
		s.Status = SourceStatusPending
	}
	if s.SourceType == "" {
		s.SourceType = SourceTypeURL
	}
	if s.ProcessingProgress == 0 {
		s.ProcessingProgress = 0
	}
	return
}

/*
* CreateSourceRequest holds data for creating a new URL source
 */
type CreateSourceRequest struct {
	URL string `json:"url" binding:"required,url"`
}

/*
* CreateTextSourceRequest holds data for creating a text source
 */
type CreateTextSourceRequest struct {
	Text string `json:"text" binding:"required"`
	Name string `json:"name"` // Optional name for the text source
}
