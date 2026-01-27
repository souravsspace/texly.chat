package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
* DocumentChunk represents a chunk of text with its vector embedding
 */
type DocumentChunk struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	SourceID  string    `json:"source_id" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"not null"`
	ChunkIndex int       `json:"chunk_index"`
	CreatedAt time.Time `json:"created_at"`
	// Embedding is stored in a virtual table for search, but we might keep it here or handle via raw SQL
	// For sqlite-vec, we typically use a virtual table `vec_items` or similar.
	// This struct represents the "Metadata" side, or we can map it to the virtual table if GORM supports it (tricky).
	// For now, we will store standard metadata here. Vector data usually lives in a separate virtual table managed via raw SQL.
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
