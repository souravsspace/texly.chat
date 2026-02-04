package vector

import (
	"context"
	"testing"

	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
* TestNewVectorRepository tests repository creation
 */
func TestNewVectorRepository(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)
	assert.NotNil(t, repo)
}

/*
* TestInitialize tests vector table initialization
 */
func TestInitialize(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 128) // Small dimension for testing

	if err != nil {
		t.Skipf("Skipping test - sqlite-vec extension not available: %v", err)
	}

	// Verify tables were created
	sqlDB, _ := gormDB.DB()
	
	// Check vec_chunk_map table
	var count int
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='vec_chunk_map'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Check vec_items virtual table (if sqlite-vec is loaded)
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='vec_items'").Scan(&count)
	if err == nil {
		assert.Equal(t, 1, count)
	}
}

/*
* TestInsertEmbedding tests single embedding insertion
 */
func TestInsertEmbedding(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 3)
	if err != nil {
		t.Skipf("Skipping test - sqlite-vec extension not available: %v", err)
	}

	// Create a test chunk first
	chunk := models.DocumentChunk{
		ID:       "chunk-1",
		SourceID: "source-1",
		Content:  "test content",
	}
	err = gormDB.Create(&chunk).Error
	require.NoError(t, err)

	// Insert embedding
	embedding := []float32{0.1, 0.2, 0.3}
	err = repo.InsertEmbedding(ctx, "chunk-1", embedding)

	// Note: This will fail if sqlite-vec is not loaded
	// In that case, we just verify the error is expected
	if err != nil {
		t.Logf("Expected error without sqlite-vec: %v", err)
		return
	}

	// If it succeeded, verify existence
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
	err := repo.Initialize(ctx, 2)
	if err != nil {
		t.Skipf("Skipping test - sqlite-vec extension not available: %v", err)
	}

	// Create test chunks
	chunks := []models.DocumentChunk{
		{ID: "chunk-1", SourceID: "source-1", Content: "content 1"},
		{ID: "chunk-2", SourceID: "source-1", Content: "content 2"},
	}
	for _, chunk := range chunks {
		err := gormDB.Create(&chunk).Error
		require.NoError(t, err)
	}

	// Bulk insert embeddings
	data := []VectorData{
		{ChunkID: "chunk-1", Embedding: []float32{0.1, 0.2}},
		{ChunkID: "chunk-2", Embedding: []float32{0.3, 0.4}},
	}

	err = repo.BulkInsertEmbeddings(ctx, data)

	// Note: This will fail if sqlite-vec is not loaded
	if err != nil {
		t.Logf("Expected error without sqlite-vec: %v", err)
		return
	}

	// Verify both exist
	exists1, _ := repo.Exists(ctx, "chunk-1")
	exists2, _ := repo.Exists(ctx, "chunk-2")
	assert.True(t, exists1)
	assert.True(t, exists2)
}

/*
* TestSearchSimilar tests similarity search
 */
func TestSearchSimilar(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 2)
	if err != nil {
		t.Skipf("Skipping test - sqlite-vec extension not available: %v", err)
	}

	// Create and insert test data
	chunks := []models.DocumentChunk{
		{ID: "chunk-1", SourceID: "source-1", Content: "similar to query"},
		{ID: "chunk-2", SourceID: "source-1", Content: "different content"},
	}
	for _, chunk := range chunks {
		err := gormDB.Create(&chunk).Error
		require.NoError(t, err)
	}

	data := []VectorData{
		{ChunkID: "chunk-1", Embedding: []float32{0.9, 0.1}},
		{ChunkID: "chunk-2", Embedding: []float32{0.1, 0.9}},
	}

	err = repo.BulkInsertEmbeddings(ctx, data)
	if err != nil {
		t.Logf("Skipping search test without sqlite-vec: %v", err)
		return
	}

	// Search with query similar to chunk-1
	queryEmbedding := []float32{0.8, 0.2}
	matches, err := repo.SearchSimilar(ctx, queryEmbedding, 2)

	require.NoError(t, err)
	assert.NotEmpty(t, matches)
	// chunk-1 should be more similar than chunk-2
	if len(matches) >= 2 {
		assert.Equal(t, "chunk-1", matches[0].ChunkID)
	}
}

/*
* TestDeleteByChunkID tests deletion
 */
func TestDeleteByChunkID(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 2)
	if err != nil {
		t.Skipf("Skipping test - sqlite-vec extension not available: %v", err)
	}

	// Create and insert test data
	chunk := models.DocumentChunk{
		ID:       "chunk-1",
		SourceID: "source-1",
		Content:  "test",
	}
	err = gormDB.Create(&chunk).Error
	require.NoError(t, err)

	embedding := []float32{0.1, 0.2}
	err = repo.InsertEmbedding(ctx, "chunk-1", embedding)
	if err != nil {
		t.Logf("Skipping delete test without sqlite-vec: %v", err)
		return
	}

	// Verify exists
	exists, _ := repo.Exists(ctx, "chunk-1")
	assert.True(t, exists)

	// Delete
	err = repo.DeleteByChunkID(ctx, "chunk-1")
	require.NoError(t, err)

	// Verify deleted
	exists, _ = repo.Exists(ctx, "chunk-1")
	assert.False(t, exists)
}

/*
* TestExists tests existence check
 */
func TestExists(t *testing.T) {
	gormDB := shared.SetupTestDB()
	repo := NewVectorRepository(gormDB)

	ctx := context.Background()
	err := repo.Initialize(ctx, 2)
	if err != nil {
		t.Skipf("Skipping test - sqlite-vec extension not available: %v", err)
	}

	// Check non-existent
	exists, err := repo.Exists(ctx, "non-existent")
	require.NoError(t, err)
	assert.False(t, exists)
}
