# Phase 2.3: Sitemap Crawler

## Goal
Allow users to train a bot on an entire website by effectively crawling its `sitemap.xml`.

## Backend Tasks

### Step 1: Sitemap Parser
**File**: `internal/services/scraper/sitemap.go`
- **Functionality**:
    - Input: Base URL or Sitemap URL.
    - Logic:
        1. Try `robots.txt` to find sitemap.
        2. Detect `sitemap.xml` location.
        3. Parse XML to extract list of URLs.
        4. Filter out non-content URLs (images, PDFs unless supported).

### Step 2: Recursive Job Enqueueing
- **Logic**:
    - For each discovered URL, create a `Source` (or sub-source if modeling hierarchy).
    - Push independent `ScrapeJob` to the Queue.
    - *Note*: Implement rate-limiting per domain to avoid getting IP banned.

### Step 3: API Endpoint
- `POST /api/bots/:id/sources/sitemap`
- Payload: `{ url: "https://example.com" }`

## Frontend Tasks

### Step 1: Sitemap Dialog
- Add tab in "Add Source" dialog for "Crawl Site".
- Warning/Info text about time to process.
