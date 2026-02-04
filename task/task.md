# Texly Development Master Task List

This file tracks the progress of the entire roadmap. Refer to the specific task files in this directory for detailed implementation instructions.

## Phase 1: Foundation & MVP

- [x] **1.1 Database & Vector Foundation** (`01-phase1-database_setup.md`)
    - [x] Configure GORM with SQLite (`internal/config/database.go`)
    - [x] Enable WAL mode for concurrency
    - [x] Integrate `sqlite-vec` extension for vector search
    - [x] Define `User` Model (Refine/Verify)
    - [x] Define `Bot` Model (`internal/models/bot.go`)
    - [x] Define `DocumentChunk` Model (`internal/models/document_chunk.go`) (Vector storage)
    - [x] Run AutoMigrate & Vector Table Creation (Virtual Tables)

- [x] **1.2 Bot Management** (`02-phase1-bot_management.md`)
    - [x] Implement Bot Handlers (Create, List, Get, Update, Delete) (`internal/handlers/bot`)
    - [x] Setup API Routes & Validation
    - [x] Create Frontend API Client (`ui/src/api/index.ts`)
    - [x] Build Dashboard UI (`ui/src/routes/dashboard/index.tsx`)
    - [x] Create CreateBotDialog Component (`ui/src/routes/dashboard/_components/create-bot-dialog.tsx`)

- [x] **1.3 Scraping & Queue System** (`03-phase1-scraping.md`)
    - [x] Design Job Queue Interface (In-memory/SQLite-backed) (`internal/queue`)
    - [x] Implement `Scraper` Service (`internal/services/scraper/scraper.go`)
    - [x] Create `ProcessSource` Worker Task (`internal/worker/tasks`)
    - [x] Implement Source Handler (Add URL) (`internal/handlers/source`)
    - [x] Build Source Management UI (Add Dialog, List)

- [x] **1.4 Embeddings & Vector Search** (`04-phase1-embeddings.md`)
    - [x] Implement Embedding Service using OpenAI (`internal/services/embedding`)
    - [x] Create Vector Repository (`internal/repo/vector/vector_repo.go`)
    - [x] Implement Vector Search Service (`internal/services/vector`)
    - [x] Integrate with scraper worker for automatic embedding generation

- [x] **1.5 Chat Interface** (`05-phase1-chat.md`)
    - [x] Implement Chat Handler with RAG Logic (`internal/handlers/chat`)
    - [x] Implement Streaming Response (SSE)
    - [x] Build Chat UI Component (`ui/src/routes/bots/$botId/chat.tsx`)

## Phase 2: Growth Features

- [/] **2.1 Widget Backend & API** (`06_phase2_widget_backend.md`)
    - [x] Update Bot Model (AllowedOrigins, Config)
    - [x] Implement Public Config Endpoint
    - [x] Implement Public Chat & Streaming Routes
    - [x] Add CORS/Security Middleware

- [x] **2.2 Widget Client (React)** (`07_phase2_widget_client.md`)
    - [x] Setup Widget Project (Vite/ShadowDOM)
    - [x] Implement Chat UI (Launcher, Window)
    - [x] Implement State & API Client
    - [x] Build Embed Script loader

- [ ] **2.3 Widget Dashboard** (`08_phase2_widget_dashboard.md`)
    - [ ] Build Configuration Form (Theme, Origins)
    - [ ] Implement Live Preview
    - [ ] Add Embed Code Generator/Display

- [ ] **2.2 File Uploads** (`09-phase2-files.md`)
    - [ ] Implement File Parser (PDF/DOCX) (`internal/services/parser`)
    - [ ] Create File Upload Handler (`internal/handlers/source`)
    - [ ] Build File Upload UI Component

- [ ] **2.3 Sitemap Crawler** (`12-phase2-sitemap.md`)
    - [ ] Enhance Scraper to parse `sitemap.xml`
    - [ ] Implement Recursive Crawling Job

## Phase 3: Scaling

- [ ] **3.1 Analytics & Security** (`10-phase3-analytics.md`)
    - [ ] Implement Analytics Service (SQL Aggregations)
    - [ ] Implement Rate Limiting Middleware (Redis)
    - [ ] Build Analytics Dashboard Charts

## Phase 4: Business Logic

- [ ] **4.1 Monetization** (`11-phase4-monetization.md`)
    - [ ] Integrate Polar.sh Go SDK
    - [ ] Implement Webhook Handler (`internal/handlers/billing`)
    - [ ] Implement Usage/Entitlement Middleware
    - [ ] Build Subscription/Upgrade UI
