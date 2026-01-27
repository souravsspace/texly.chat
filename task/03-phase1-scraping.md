# Phase 1.3: URL Scraper Service & Queue

## Goal
Ingest training data from websites using a simple, robust queue system (initially in-memory or SQLite-backed) to handle scraping jobs without external infrastructure like Redis.

## Backend Tasks

### Step 1: Simple Job Queue
**File**: `internal/queue/queue.go`
- **Design**:
    - Build a simple `JobQueue` interface.
    - Implementation 1: **In-Memory** using Go Channels (good for MVP/Single Instance).
    - (Optional for Phase 1, Required for Scale): **SQLite-backed** table `jobs` if persistence is needed immediately.
    - **Recommendation**: Start with a buffered channel `make(chan Job, 100)` for simplicity in MVP.

```go
type Job struct {
    SourceID string
    URL      string
}
// Worker pool to process jobs
```

### Step 2: Scraper Service
**File**: `internal/services/scraper/service.go`
- Libs: `gocolly/colly/v2` or `net/http` + `goquery`
- **Logic**:
    1. Fetch URL.
    2. Extract semantic content (Title, Body).
    3. Clean HTML (remove scripts, styles, navs).
    4. Return text content.

### Step 3: Source Handler & Worker Integration
**File**: `internal/worker/worker.go`
- **Worker**: Listens to the `JobQueue`.
- **Process**:
    1. Update `Source` status -> `processing`.
    2. Run Scraper.
    3. Chunk content (e.g., 500-1000 tokens).
    4. Save `DocumentChunk`s (text only initially, or call Embedding service immediately).
    5. Update `Source` status -> `completed` / `failed`.

**File**: `internal/handlers/source/handler.go`
- `POST /api/bots/:id/sources`
- Create `Source` record in DB.
- Push job to `JobQueue`.

## Frontend Tasks

### Step 1: Source Management
**File**: `ui/src/routes/bots/$botId/sources.tsx` (or similar)
- **Add Source**: Simple form taking a URL.
- **List Sources**: Table showing URL, Status (Pending, Processing, Completed, Failed), and Date.
- **Polling**: Auto-refresh list every few seconds to show status changes.
