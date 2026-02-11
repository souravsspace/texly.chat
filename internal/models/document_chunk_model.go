package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

/*
* DocumentChunk represents a chunk of text with its vector embedding
 */
type DocumentChunk struct {
	ID         string           `json:"id" gorm:"primaryKey"`
	SourceID   string           `json:"source_id" gorm:"not null;index"`
	Content    string           `json:"content" gorm:"not null"`
	ChunkIndex int              `json:"chunk_index"`
	Embedding  *pgvector.Vector `json:"-" gorm:"type:vector(1536)"`
	CreatedAt  time.Time        `json:"created_at"`

	// Relation to Source (for GORM Preload)
	Source Source `json:"source,omitempty" gorm:"foreignKey:SourceID"`
}

/*
* BeforeCreate generates a new UUID for the chunk
 */
func (d *DocumentChunk) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return
}
