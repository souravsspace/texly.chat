package vector

import (
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

/*
* VectorRepository handles vector storage and search operations using sqlite-vec
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
* Initialize creates the necessary tables for vector storage
* This should be called during database migration
 */
func (r *VectorRepository) Initialize(ctx context.Context, dimension int) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Create mapping table to link rowid to chunk_id
	createMapTableQuery := `
	CREATE TABLE IF NOT EXISTS vec_chunk_map (
		rowid INTEGER PRIMARY KEY AUTOINCREMENT,
		chunk_id TEXT NOT NULL UNIQUE,
		FOREIGN KEY(chunk_id) REFERENCES document_chunks(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_vec_chunk_map_chunk_id ON vec_chunk_map(chunk_id);
	`

	if _, err := sqlDB.ExecContext(ctx, createMapTableQuery); err != nil {
		return fmt.Errorf("failed to create vec_chunk_map table: %w", err)
	}

	// Create virtual table for vector embeddings
	createVecTableQuery := fmt.Sprintf(`
	CREATE VIRTUAL TABLE IF NOT EXISTS vec_items USING vec0(
		embedding float[%d]
	);
	`, dimension)

	if _, err := sqlDB.ExecContext(ctx, createVecTableQuery); err != nil {
		return fmt.Errorf("failed to create vec_items table: %w", err)
	}

	return nil
}

/*
* InsertEmbedding inserts a single embedding for a chunk
 */
func (r *VectorRepository) InsertEmbedding(ctx context.Context, chunkID string, embedding []float32) error {
	return r.BulkInsertEmbeddings(ctx, []VectorData{{ChunkID: chunkID, Embedding: embedding}})
}

/*
* BulkInsertEmbeddings inserts multiple embeddings efficiently in a transaction
 */
func (r *VectorRepository) BulkInsertEmbeddings(ctx context.Context, data []VectorData) error {
	if len(data) == 0 {
		return nil
	}

	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Start transaction
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare statements
	insertMapStmt, err := tx.PrepareContext(ctx, "INSERT INTO vec_chunk_map (chunk_id) VALUES (?) RETURNING rowid")
	if err != nil {
		return fmt.Errorf("failed to prepare map insert statement: %w", err)
	}
	defer insertMapStmt.Close()

	insertVecStmt, err := tx.PrepareContext(ctx, "INSERT INTO vec_items (rowid, embedding) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare vec insert statement: %w", err)
	}
	defer insertVecStmt.Close()

	// Insert each embedding
	for _, item := range data {
		// Insert into mapping table and get rowid
		var rowid int64
		err := insertMapStmt.QueryRowContext(ctx, item.ChunkID).Scan(&rowid)
		if err != nil {
			return fmt.Errorf("failed to insert chunk mapping for %s: %w", item.ChunkID, err)
		}

		// Convert embedding to bytes
		embeddingBytes := float32SliceToBytes(item.Embedding)

		// Insert into vector table
		if _, err := insertVecStmt.ExecContext(ctx, rowid, embeddingBytes); err != nil {
			return fmt.Errorf("failed to insert embedding for rowid %d: %w", rowid, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

/*
* SearchSimilar performs cosine similarity search
 */
func (r *VectorRepository) SearchSimilar(ctx context.Context, embedding []float32, limit int) ([]VectorMatch, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	embeddingBytes := float32SliceToBytes(embedding)

	// Query using vec_distance_cosine for similarity search
	query := `
	SELECT 
		m.chunk_id,
		distance
	FROM vec_items v
	JOIN vec_chunk_map m ON v.rowid = m.rowid
	WHERE v.embedding MATCH ?
	  AND k = ?
	ORDER BY distance
	`

	rows, err := sqlDB.QueryContext(ctx, query, embeddingBytes, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	var matches []VectorMatch
	for rows.Next() {
		var match VectorMatch
		if err := rows.Scan(&match.ChunkID, &match.Distance); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return matches, nil
}

/*
* DeleteByChunkID deletes an embedding by chunk ID
 */
func (r *VectorRepository) DeleteByChunkID(ctx context.Context, chunkID string) error {
	return r.DeleteByChunkIDs(ctx, []string{chunkID})
}

/*
* DeleteByChunkIDs deletes multiple embeddings by chunk IDs
 */
func (r *VectorRepository) DeleteByChunkIDs(ctx context.Context, chunkIDs []string) error {
	if len(chunkIDs) == 0 {
		return nil
	}

	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Start transaction
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get rowids for the chunk IDs
	placeholders := make([]string, len(chunkIDs))
	args := make([]interface{}, len(chunkIDs))
	for i, id := range chunkIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("SELECT rowid FROM vec_chunk_map WHERE chunk_id IN (%s)", strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to query rowids: %w", err)
	}

	var rowids []int64
	for rows.Next() {
		var rowid int64
		if err := rows.Scan(&rowid); err != nil {
			rows.Close()
			return fmt.Errorf("failed to scan rowid: %w", err)
		}
		rowids = append(rowids, rowid)
	}
	rows.Close()

	if len(rowids) == 0 {
		return nil // Nothing to delete
	}

	// Delete from vec_items
	rowidPlaceholders := make([]string, len(rowids))
	rowidArgs := make([]interface{}, len(rowids))
	for i, id := range rowids {
		rowidPlaceholders[i] = "?"
		rowidArgs[i] = id
	}

	deleteVecQuery := fmt.Sprintf("DELETE FROM vec_items WHERE rowid IN (%s)", strings.Join(rowidPlaceholders, ","))
	if _, err := tx.ExecContext(ctx, deleteVecQuery, rowidArgs...); err != nil {
		return fmt.Errorf("failed to delete from vec_items: %w", err)
	}

	// Delete from vec_chunk_map
	deleteMapQuery := fmt.Sprintf("DELETE FROM vec_chunk_map WHERE chunk_id IN (%s)", strings.Join(placeholders, ","))
	if _, err := tx.ExecContext(ctx, deleteMapQuery, args...); err != nil {
		return fmt.Errorf("failed to delete from vec_chunk_map: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

/*
* Exists checks if an embedding exists for a chunk ID
 */
func (r *VectorRepository) Exists(ctx context.Context, chunkID string) (bool, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return false, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	var count int
	err = sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM vec_chunk_map WHERE chunk_id = ?", chunkID).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return count > 0, nil
}

/*
* Helper function to convert float32 slice to bytes for sqlite-vec
 */
func float32SliceToBytes(floats []float32) []byte {
	bytes := make([]byte, len(floats)*4)
	for i, f := range floats {
		bits := math.Float32bits(f)
		binary.LittleEndian.PutUint32(bytes[i*4:], bits)
	}
	return bytes
}
