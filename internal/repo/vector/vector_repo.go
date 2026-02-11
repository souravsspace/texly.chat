package vector

import (
	"context"
	"fmt"

	"github.com/pgvector/pgvector-go"
	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/gorm"
)

/*
 * VectorRepository handles vector storage and search operations using pgvector
 */
type VectorRepository struct {
	db *gorm.DB
}

/*
 * NewVectorRepository creates a new vector repository instance
 */
func NewVectorRepository(db *gorm.DB) *VectorRepository {
	return &VectorRepository{db: db}
}

/*
 * VectorData represents a chunk ID with its embedding vector
 */
type VectorData struct {
	ChunkID   string
	Embedding []float32
}

/*
 * VectorMatch represents a search result with similarity score
 */
type VectorMatch struct {
	ChunkID  string
	Distance float32
}

/*
 * Initialize creates vector indexes for better performance
 * This should be called after initial data loading
 */
func (r *VectorRepository) Initialize(ctx context.Context, dimension int) error {
	// Create IVFFlat index for faster similarity search (recommended for production)
	// Note: This requires some data to be present for clustering
	// For now, we'll just ensure the vector extension is enabled (done in db.Connect)
	
	// Optionally create index if you have enough data (>1000 rows recommended)
	// indexQuery := `
	// CREATE INDEX IF NOT EXISTS document_chunks_embedding_idx 
	// ON document_chunks 
	// USING ivfflat (embedding vector_cosine_ops)
	// WITH (lists = 100);
	// `
	// return r.db.WithContext(ctx).Exec(indexQuery).Error

	fmt.Println("✅ Vector repository initialized (indexes can be added later)")
	return nil
}

/*
 * InsertEmbedding inserts or updates an embedding for a chunk
 */
func (r *VectorRepository) InsertEmbedding(ctx context.Context, chunkID string, embedding []float32) error {
	vec := pgvector.NewVector(embedding)
	
	result := r.db.WithContext(ctx).
		Model(&models.DocumentChunk{}).
		Where("id = ?", chunkID).
		Update("embedding", &vec)
	
	if result.Error != nil {
		return fmt.Errorf("failed to insert embedding for chunk %s: %w", chunkID, result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("chunk %s not found", chunkID)
	}
	
	return nil
}

/*
 * BulkInsertEmbeddings inserts multiple embeddings efficiently
 */
func (r *VectorRepository) BulkInsertEmbeddings(ctx context.Context, data []VectorData) error {
	if len(data) == 0 {
		return nil
	}

	// Use a transaction for atomicity
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range data {
			vec := pgvector.NewVector(item.Embedding)
			
			result := tx.Model(&models.DocumentChunk{}).
				Where("id = ?", item.ChunkID).
				Update("embedding", &vec)
			
			if result.Error != nil {
				return fmt.Errorf("failed to insert embedding for chunk %s: %w", item.ChunkID, result.Error)
			}
		}
		return nil
	})
}

/*
 * SearchSimilar performs cosine similarity search using pgvector
 * Returns the most similar chunks ordered by distance (ascending)
 */
func (r *VectorRepository) SearchSimilar(ctx context.Context, embedding []float32, limit int) ([]VectorMatch, error) {
	vec := pgvector.NewVector(embedding)
	
	var results []struct {
		ID       string
		Distance float32
	}
	
	// Use cosine distance operator <=> for similarity search
	// Lower distance = more similar
	err := r.db.WithContext(ctx).
		Table("document_chunks").
		Select("id, embedding <=> ? as distance", vec).
		Where("embedding IS NOT NULL").
		Order("distance").
		Limit(limit).
		Find(&results).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to execute similarity search: %w", err)
	}
	
	// Convert to VectorMatch
	matches := make([]VectorMatch, len(results))
	for i, r := range results {
		matches[i] = VectorMatch{
			ChunkID:  r.ID,
			Distance: r.Distance,
		}
	}
	
	return matches, nil
}

/*
 * DeleteByChunkID deletes an embedding by setting it to NULL
 * The chunk record itself remains for potential re-embedding
 */
func (r *VectorRepository) DeleteByChunkID(ctx context.Context, chunkID string) error {
	return r.DeleteByChunkIDs(ctx, []string{chunkID})
}

/*
 * DeleteByChunkIDs deletes multiple embeddings
 */
func (r *VectorRepository) DeleteByChunkIDs(ctx context.Context, chunkIDs []string) error {
	if len(chunkIDs) == 0 {
		return nil
	}
	
	result := r.db.WithContext(ctx).
		Model(&models.DocumentChunk{}).
		Where("id IN ?", chunkIDs).
		Update("embedding", nil)
	
	if result.Error != nil {
		return fmt.Errorf("failed to delete embeddings: %w", result.Error)
	}
	
	return nil
}

/*
 * Exists checks if an embedding exists for a chunk ID
 */
func (r *VectorRepository) Exists(ctx context.Context, chunkID string) (bool, error) {
	var count int64
	
	err := r.db.WithContext(ctx).
		Model(&models.DocumentChunk{}).
		Where("id = ? AND embedding IS NOT NULL", chunkID).
		Count(&count).Error
	
	if err != nil {
		return false, fmt.Errorf("failed to check embedding existence: %w", err)
	}
	
	return count > 0, nil
}
