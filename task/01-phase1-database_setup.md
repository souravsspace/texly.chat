# Phase 1.1: Database Setup & Core Models

## Goal
Initialize the SQLite database with **WAL mode** enabled for high concurrency and the **sqlite-vec** extension for vector search. Define the core data models.

## Backend Tasks

### Step 1: SQLite Configuration
**File**: `internal/db/db.go` (or `internal/config/database.go`)
- [ ] **Connect**: Use `gorm.io/driver/sqlite`.
- [ ] **WAL Mode**: Execute `PRAGMA journal_mode=WAL;` immediately after connection.
- [ ] **Vector Extension**: Load `sqlite-vec` extension (ensure the dynamic library is present or statically linked).
    - If using `mattn/go-sqlite3` directly or via GORM, ensure `sqlite_vec` is loaded.
    - Validate with `SELECT vec_version();`.

### Step 2: Core Models
**File**: `internal/models/bot.go`
```go
type Bot struct {
    ID           string    `gorm:"primaryKey"` // UUID
    UserID       string    `gorm:"not null;index"`
    Name         string    `gorm:"not null"`
    SystemPrompt string    `gorm:"type:text"`
    Config       string    `gorm:"type:json"` // Store extra settings
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

**File**: `internal/models/document_chunk.go`
- This model stores the **metadata** and text content. The vector itself lives in a virtual table for efficiency.
```go
type DocumentChunk struct {
    ID        string    `gorm:"primaryKey"` // UUID
    SourceID  string    `gorm:"not null;index"`
    Content   string    `gorm:"not null"` 
    ChunkIndex int
    CreatedAt time.Time
    // Note: Embedding is NOT stored here to avoid bloating the standard table.
}
```

### Step 3: Vector Virtual Table
- GORM does not support creating Virtual Tables via `AutoMigrate`.
- **Action**: Execute Raw SQL migration.
```sql
CREATE VIRTUAL TABLE IF NOT EXISTS vec_items USING vec0(
    id TEXT PRIMARY KEY,
    embedding float[1536]
);
```
- **Strategy**: When inserting a chunk, insert into both `document_chunks` (GORM) and `vec_items` (Raw SQL) using the same UUID.

### Step 4: Run Migrations
- Update `Migrate` function to run `AutoMigrate` for structs and `Exec` for `vec_items`.

## Frontend Tasks
- None for this step.
