package vector

import (
	"context"
	"testing"

	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Helper to generate a valid 1536-dimension embedding for testing
func generateTestEmbedding(seed float32) []float32 {
	embedding := make([]float32, 1536)
	for i := range embedding {
		embedding[i] = seed + float32(i)*0.0001
	}
	return embedding
}

// Helper to create test bot and source
func setupTestBotAndSource(t *testing.T, gormDB *gorm.DB) (string, string) {
	bot := models.Bot{
		ID:   "test-bot-1",
		Name: "Test Bot",
	}
	err := gormDB.Create(&bot).Error
	require.NoError(t, err)

	source := models.Source{
		ID:     "test-source-1",
		BotID:  "test-bot-1",
		URL:    "https://example.com",
		Status: models.SourceStatusCompleted,
	}
	err = gormDB.Create(&source).Error
	require.NoError(t, err)

	return bot.ID, source.ID
}

/*
 * TestNewVectorRepository tests repository creation
 */
func TestNewVectorRepository(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)
	assert.NotNil(t, repo)
}

/*
 * TestInitialize tests vector repository initialization
 */
func TestInitialize(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 1536)

	require.NoError(t, err)
}

/*
 * TestInsertEmbedding tests single embedding insertion
 */
func TestInsertEmbedding(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 1536)
	require.NoError(t, err)

	// Create parent records
	_, sourceID := setupTestBotAndSource(t, gormDB)

	// Create a test chunk
	chunk := models.DocumentChunk{
		ID:       "chunk-1",
		SourceID: sourceID,
		Content:  "test content",
	}
	err = gormDB.Create(&chunk).Error
	require.NoError(t, err)

	// Insert embedding
	embedding := generateTestEmbedding(0.1)
	err = repo.InsertEmbedding(ctx, "chunk-1", embedding)
	require.NoError(t, err)

	// Verify existence
	exists, err := repo.Exists(ctx, "chunk-1")
	require.NoError(t, err)
	assert.True(t, exists)
}

/*
 * TestBulkInsertEmbeddings tests bulk embedding insertion
 */
func TestBulkInsertEmbeddings(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 1536)
	require.NoError(t, err)

	// Create parent records
	_, sourceID := setupTestBotAndSource(t, gormDB)

	// Create test chunks
	chunks := []models.DocumentChunk{
		{ID: "chunk-1", SourceID: sourceID, Content: "content 1"},
		{ID: "chunk-2", SourceID: sourceID, Content: "content 2"},
	}
	for _, chunk := range chunks {
		err := gormDB.Create(&chunk).Error
		require.NoError(t, err)
	}

	// Bulk insert embeddings
	data := []VectorData{
		{ChunkID: "chunk-1", Embedding: generateTestEmbedding(0.1)},
		{ChunkID: "chunk-2", Embedding: generateTestEmbedding(0.3)},
	}

	err = repo.BulkInsertEmbeddings(ctx, data)
	require.NoError(t, err)

	// Verify both exist
	exists1, err := repo.Exists(ctx, "chunk-1")
	require.NoError(t, err)
	assert.True(t, exists1)

	exists2, err := repo.Exists(ctx, "chunk-2")
	require.NoError(t, err)
	assert.True(t, exists2)
}

/*
 * TestSearchSimilar tests similarity search
 */
func TestSearchSimilar(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 1536)
	require.NoError(t, err)

	// Create parent records
	_, sourceID := setupTestBotAndSource(t, gormDB)

	// Create and insert test data
	chunks := []models.DocumentChunk{
		{ID: "chunk-1", SourceID: sourceID, Content: "similar to query"},
		{ID: "chunk-2", SourceID: sourceID, Content: "different content"},
	}
	for _, chunk := range chunks {
		err := gormDB.Create(&chunk).Error
		require.NoError(t, err)
	}

	data := []VectorData{
		{ChunkID: "chunk-1", Embedding: generateTestEmbedding(0.9)},
		{ChunkID: "chunk-2", Embedding: generateTestEmbedding(0.1)},
	}

	err = repo.BulkInsertEmbeddings(ctx, data)
	require.NoError(t, err)

	// Search
	queryEmbedding := generateTestEmbedding(0.85) // Similar to chunk-1
	matches, err := repo.SearchSimilar(ctx, queryEmbedding, 5)

	require.NoError(t, err)
	assert.NotEmpty(t, matches)

	// Verify chunk-1 is more similar (lower distance)
	if len(matches) >= 2 {
		assert.Equal(t, "chunk-1", matches[0].ChunkID)
		assert.Less(t, matches[0].Distance, matches[1].Distance)
	}
}

/*
 * TestDeleteByChunkID tests deletion
 */
func TestDeleteByChunkID(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 1536)
	require.NoError(t, err)

	// Create parent records
	_, sourceID := setupTestBotAndSource(t, gormDB)

	// Create chunk and insert embedding
	chunk := models.DocumentChunk{
		ID:       "chunk-1",
		SourceID: sourceID,
		Content:  "test content",
	}
	err = gormDB.Create(&chunk).Error
	require.NoError(t, err)

	embedding := generateTestEmbedding(0.5)
	err = repo.InsertEmbedding(ctx, "chunk-1", embedding)
	require.NoError(t, err)

	// Verify exists
	exists, err := repo.Exists(ctx, "chunk-1")
	require.NoError(t, err)
	assert.True(t, exists)

	// Delete
	err = repo.DeleteByChunkID(ctx, "chunk-1")
	require.NoError(t, err)

	// Verify deleted
	exists, err = repo.Exists(ctx, "chunk-1")
	require.NoError(t, err)
	assert.False(t, exists)
}

/*
 * TestExists tests existence checking  
 */
func TestExists(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 1536)
	require.NoError(t, err)

	// Check non-existent
	exists, err := repo.Exists(ctx, "non-existent")
	require.NoError(t, err)
	assert.False(t, exists)

	// Create parent records
	_, sourceID := setupTestBotAndSource(t, gormDB)

	// Create chunk with embedding
	chunk := models.DocumentChunk{
		ID:       "chunk-1",
		SourceID: sourceID,
		Content:  "test content",
	}
	err = gormDB.Create(&chunk).Error
	require.NoError(t, err)

	embedding := generateTestEmbedding(0.5)
	err = repo.InsertEmbedding(ctx, "chunk-1", embedding)
	require.NoError(t, err)

	// Check exists
	exists, err = repo.Exists(ctx, "chunk-1")
	require.NoError(t, err)
	assert.True(t, exists)
}
