# Phase 1.4: Vector Embedding & Search

## Goal
Generate embeddings for scraped content and store them in SQLite using `sqlite-vec`.

## Backend Tasks

### Step 1: Embedding Service
**File**: `internal/services/embedding/service.go`
- Call OpenAI API (`text-embedding-3-small`).
- Chunk content (1000 chars, 200 overlap).

### Step 2: Vector Storage (SQLite)
**File**: `internal/repo/vector/vector_repo.go`
- **Insert**: Explicit SQL for `vec_items` virtual table.
```go
// Example Insert
db.Exec("INSERT INTO vec_items(rowid, embedding) VALUES (?, ?)", chunkID, embeddingVector)
```
- Note: `sqlite-vec` usually maps a rowid to a vector. You might need to link this `rowid` to your `DocumentChunk.ID`.

### Step 3: Vector Search
**File**: `internal/services/vector/search.go`
- **Query**: Use `vec_distance_cosine` (or L2).
```sql
SELECT
  rowid,
  distance
FROM vec_items
WHERE embedding MATCH ?
ORDER BY distance
LIMIT k;
```
- Retrieve actual content from `document_chunks` table using the returned `rowid`.

## Frontend Tasks
- None (Background process).
