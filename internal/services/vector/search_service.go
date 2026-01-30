package vector

import (
	"context"
	"fmt"

	"github.com/souravsspace/texly.chat/internal/models"
	vectorRepo "github.com/souravsspace/texly.chat/internal/repo/vector"
	"github.com/souravsspace/texly.chat/internal/services/embedding"
	"gorm.io/gorm"
)

/*
* SearchService provides semantic search capabilities
 */
type SearchService struct {
	db               *gorm.DB
	vectorRepo       *vectorRepo.VectorRepository
	embeddingService *embedding.EmbeddingService
}

/*
* NewSearchService creates a new search service instance
 */
func NewSearchService(
	db *gorm.DB,
	vectorRepo *vectorRepo.VectorRepository,
	embeddingService *embedding.EmbeddingService,
) *SearchService {
	return &SearchService{
		db:               db,
		vectorRepo:       vectorRepo,
		embeddingService: embeddingService,
	}
}

/*
* SearchResult represents a single search result with metadata
 */
type SearchResult struct {
	ChunkID      string                 `json:"chunk_id"`
	SourceID     string                 `json:"source_id"`
	Content      string                 `json:"content"`
	Distance     float32                `json:"distance"`
	ChunkIndex   int                    `json:"chunk_index"`
	URL          string                 `json:"url"`
	SourceStatus models.SourceStatus    `json:"source_status"`
	Metadata     map[string]interface{} `json:"metadata"`
}

/*
* SearchSimilar performs semantic search for a text query
 */
func (s *SearchService) SearchSimilar(ctx context.Context, query string, botID string, limit int) ([]SearchResult, error) {
	// Generate embedding for the query
	queryEmbedding, err := s.embeddingService.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	return s.SearchSimilarByEmbedding(ctx, queryEmbedding, botID, limit)
}

/*
* SearchSimilarByEmbedding performs semantic search using a pre-computed embedding
 */
func (s *SearchService) SearchSimilarByEmbedding(ctx context.Context, embedding []float32, botID string, limit int) ([]SearchResult, error) {
	// Perform vector similarity search
	matches, err := s.vectorRepo.SearchSimilar(ctx, embedding, limit*2) // Get more to filter by botID
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	if len(matches) == 0 {
		return []SearchResult{}, nil
	}

	// Extract chunk IDs
	chunkIDs := make([]string, len(matches))
	distanceMap := make(map[string]float32)
	for i, match := range matches {
		chunkIDs[i] = match.ChunkID
		distanceMap[match.ChunkID] = match.Distance
	}

	// Fetch chunk metadata with source preloaded
	var chunks []models.DocumentChunk

	err = s.db.WithContext(ctx).
		Preload("Source").
		Where("id IN ?", chunkIDs).
		Find(&chunks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch chunk metadata: %w", err)
	}

	// Filter by bot ID and deleted_at
	var filteredChunks []models.DocumentChunk
	for _, chunk := range chunks {
		if chunk.Source.BotID == botID && chunk.Source.DeletedAt.Time.IsZero() {
			filteredChunks = append(filteredChunks, chunk)
		}
	}

	// Build results maintaining original order by distance
	results := make([]SearchResult, 0, len(filteredChunks))
	for _, match := range matches {
		for _, chunk := range filteredChunks {
			if chunk.ID == match.ChunkID {
				results = append(results, SearchResult{
					ChunkID:      chunk.ID,
					SourceID:     chunk.SourceID,
					Content:      chunk.Content,
					Distance:     distanceMap[chunk.ID],
					ChunkIndex:   chunk.ChunkIndex,
					URL:          chunk.Source.URL,
					SourceStatus: chunk.Source.Status,
					Metadata: map[string]interface{}{
						"created_at": chunk.CreatedAt,
						"source_url": chunk.Source.URL,
					},
				})
				break
			}
		}

		// Stop if we have enough results
		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

/*
* SearchMultipleBots searches across multiple bots (e.g., for admin features)
 */
func (s *SearchService) SearchMultipleBots(ctx context.Context, query string, botIDs []string, limit int) ([]SearchResult, error) {
	// Generate embedding for the query
	queryEmbedding, err := s.embeddingService.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Perform vector similarity search
	matches, err := s.vectorRepo.SearchSimilar(ctx, queryEmbedding, limit*len(botIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	if len(matches) == 0 {
		return []SearchResult{}, nil
	}

	// Extract chunk IDs
	chunkIDs := make([]string, len(matches))
	distanceMap := make(map[string]float32)
	for i, match := range matches {
		chunkIDs[i] = match.ChunkID
		distanceMap[match.ChunkID] = match.Distance
	}

	// Fetch chunk metadata and join with sources
	var chunks []struct {
		models.DocumentChunk
		Source models.Source `gorm:"foreignKey:SourceID"`
	}

	dbQuery := s.db.WithContext(ctx).
		Table("document_chunks").
		Select("document_chunks.*, sources.*").
		Joins("JOIN sources ON sources.id = document_chunks.source_id").
		Where("document_chunks.id IN ?", chunkIDs).
		Where("sources.deleted_at IS NULL")

	if len(botIDs) > 0 {
		dbQuery = dbQuery.Where("sources.bot_id IN ?", botIDs)
	}

	err = dbQuery.Find(&chunks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chunk metadata: %w", err)
	}

	// Build results
	results := make([]SearchResult, 0, len(chunks))
	for _, match := range matches {
		for _, chunk := range chunks {
			if chunk.ID == match.ChunkID {
				results = append(results, SearchResult{
					ChunkID:      chunk.ID,
					SourceID:     chunk.SourceID,
					Content:      chunk.Content,
					Distance:     distanceMap[chunk.ID],
					ChunkIndex:   chunk.ChunkIndex,
					URL:          chunk.Source.URL,
					SourceStatus: chunk.Source.Status,
					Metadata: map[string]interface{}{
						"created_at": chunk.CreatedAt,
						"source_url": chunk.Source.URL,
						"bot_id":     chunk.Source.BotID,
					},
				})
				break
			}
		}

		if len(results) >= limit {
			break
		}
	}

	return results, nil
}
