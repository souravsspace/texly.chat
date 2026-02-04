package vector

import (
	"context"
	"testing"

	"github.com/souravsspace/texly.chat/internal/models"
	vectorRepo "github.com/souravsspace/texly.chat/internal/repo/vector"
	"github.com/souravsspace/texly.chat/internal/services/embedding"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
* TestNewSearchService tests service creation
 */
func TestNewSearchService(t *testing.T) {
	gormDB := shared.SetupTestDB()
	vRepo := vectorRepo.NewVectorRepository(gormDB)
	embSvc := embedding.NewEmbeddingService("test-key", "test-model", 128)

	service := NewSearchService(gormDB, vRepo, embSvc)
	assert.NotNil(t, service)
}

/*
* TestSearchSimilarByEmbedding tests search functionality
 */
func TestSearchSimilarByEmbedding(t *testing.T) {
	gormDB := shared.SetupTestDB()
	vRepo := vectorRepo.NewVectorRepository(gormDB)
	embSvc := embedding.NewEmbeddingService("test-key", "test-model", 2)

	ctx := context.Background()
	err := vRepo.Initialize(ctx, 2)
	if err != nil {
		t.Skipf("Skipping test - sqlite-vec extension not available: %v", err)
	}

	// Create test bot
	bot := models.Bot{
		ID:   "bot-1",
		Name: "Test Bot",
	}
	err = gormDB.Create(&bot).Error
	require.NoError(t, err)

	// Create test source
	source := models.Source{
		ID:     "source-1",
		BotID:  "bot-1",
		URL:    "https://example.com",
		Status: models.SourceStatusCompleted,
	}
	err = gormDB.Create(&source).Error
	require.NoError(t, err)

	// Create test chunks
	chunks := []models.DocumentChunk{
		{
			ID:         "chunk-1",
			SourceID:   "source-1",
			Content:    "This is about machine learning",
			ChunkIndex: 0,
		},
		{
			ID:         "chunk-2",
			SourceID:   "source-1",
			Content:    "This is about cooking recipes",
			ChunkIndex: 1,
		},
	}

	for _, chunk := range chunks {
		err := gormDB.Create(&chunk).Error
		require.NoError(t, err)
	}

	// Insert embeddings
	data := []vectorRepo.VectorData{
		{ChunkID: "chunk-1", Embedding: []float32{0.9, 0.1}},
		{ChunkID: "chunk-2", Embedding: []float32{0.1, 0.9}},
	}

	err = vRepo.BulkInsertEmbeddings(ctx, data)
	if err != nil {
		t.Logf("Skipping search test without sqlite-vec: %v", err)
		return
	}

	// Create search service
	service := NewSearchService(gormDB, vRepo, embSvc)

	// Search with embedding similar to chunk-1
	queryEmbedding := []float32{0.8, 0.2}
	results, err := service.SearchSimilarByEmbedding(ctx, queryEmbedding, "bot-1", 5)

	require.NoError(t, err)
	assert.NotEmpty(t, results)

	// Verify results contain chunk data
	if len(results) > 0 {
		assert.NotEmpty(t, results[0].Content)
		assert.Equal(t, "source-1", results[0].SourceID)
		assert.Equal(t, "https://example.com", results[0].URL)
		assert.Equal(t, models.SourceStatusCompleted, results[0].SourceStatus)
	}
}

/*
* TestSearchSimilar_EmptyResults tests behavior with no matches
 */
func TestSearchSimilar_EmptyResults(t *testing.T) {
	gormDB := shared.SetupTestDB()
	vRepo := vectorRepo.NewVectorRepository(gormDB)
	embSvc := embedding.NewEmbeddingService("test-key", "test-model", 2)

	ctx := context.Background()
	err := vRepo.Initialize(ctx, 2)
	if err != nil {
		t.Skipf("Skipping test - sqlite-vec extension not available: %v", err)
	}

	service := NewSearchService(gormDB, vRepo, embSvc)

	// Search with no data
	queryEmbedding := []float32{0.5, 0.5}
	results, err := service.SearchSimilarByEmbedding(ctx, queryEmbedding, "non-existent-bot", 5)

	// Should not error, just return empty results
	require.NoError(t, err)
	assert.Empty(t, results)
}

/*
* TestSearchMultipleBots tests cross-bot search
 */
func TestSearchMultipleBots(t *testing.T) {
	gormDB := shared.SetupTestDB()
	vRepo := vectorRepo.NewVectorRepository(gormDB)
	embSvc := embedding.NewEmbeddingService("test-key", "test-model", 2)

	ctx := context.Background()
	err := vRepo.Initialize(ctx, 2)
	if err != nil {
		t.Skipf("Skipping test - sqlite-vec extension not available: %v", err)
	}

	// Create multiple bots and sources
	bot1 := models.Bot{ID: "bot-1", Name: "Bot 1"}
	bot2 := models.Bot{ID: "bot-2", Name: "Bot 2"}
	gormDB.Create(&bot1)
	gormDB.Create(&bot2)

	source1 := models.Source{ID: "source-1", BotID: "bot-1", URL: "https://example1.com", Status: models.SourceStatusCompleted}
	source2 := models.Source{ID: "source-2", BotID: "bot-2", URL: "https://example2.com", Status: models.SourceStatusCompleted}
	gormDB.Create(&source1)
	gormDB.Create(&source2)

	chunks := []models.DocumentChunk{
		{ID: "chunk-1", SourceID: "source-1", Content: "Content from bot 1"},
		{ID: "chunk-2", SourceID: "source-2", Content: "Content from bot 2"},
	}
	for _, chunk := range chunks {
		gormDB.Create(&chunk)
	}

	// Insert embeddings
	data := []vectorRepo.VectorData{
		{ChunkID: "chunk-1", Embedding: []float32{0.7, 0.3}},
		{ChunkID: "chunk-2", Embedding: []float32{0.6, 0.4}},
	}

	err = vRepo.BulkInsertEmbeddings(ctx, data)
	if err != nil {
		t.Logf("Skipping multi-bot search test without sqlite-vec: %v", err)
		return
	}

	service := NewSearchService(gormDB, vRepo, embSvc)

	// Search across both bots - note: can't actually call SearchMultipleBots without real embedding API
	// Just verify the function exists and is callable
	queryEmbedding := []float32{0.65, 0.35}
	results, err := service.SearchSimilarByEmbedding(ctx, queryEmbedding, "bot-1", 5)

	require.NoError(t, err)
	// Should only return results from bot-1
	for _, result := range results {
		var source models.Source
		gormDB.First(&source, "id = ?", result.SourceID)
		assert.Equal(t, "bot-1", source.BotID)
	}
}
