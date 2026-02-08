# Texly.Chat - Architecture & Development Guide

## Project Overview

**Texly.Chat** is a powerful SaaS chatbot platform with built-in Retrieval-Augmented Generation (RAG) capabilities. Users can create custom AI chatbots trained on their own data sources (URLs, PDF files, Excel spreadsheets, and plain text).

### Project Type
- Full-stack SaaS application
- Monolithic backend (Go)
- Separate frontend and widget applications (React)
- Real-time chat with streaming responses
- Embeddable widget for third-party websites

### Key Features
- Custom chatbot creation and management
- RAG (Retrieval-Augmented Generation) for context-aware answers
- Web scraping and URL indexing
- File support (PDF, Excel, text)
- Vector search using `sqlite-vec` and OpenAI embeddings
- Real-time chat via Server-Sent Events (SSE)
- Embeddable chat widget with Shadow DOM isolation
- Session persistence
- Docker containerization

---

## Tech Stack

### Backend
- **Language**: Go 1.25+
- **Framework**: Gin Web Framework (HTTP routing)
- **Database**: SQLite with WAL mode
- **ORM**: GORM (database abstraction)
- **Vector DB**: `sqlite-vec` extension (in-process vector search)
- **AI/ML**: OpenAI API (embeddings & chat completions)
- **File Storage**: MinIO (S3-compatible object storage)
- **Async Processing**: In-memory job queue with goroutine workers

### Frontend
- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite
- **Routing**: TanStack Router (file-based)
- **State Management**: TanStack Query (server state) + Zustand (client state)
- **Styling**: TailwindCSS v4
- **UI Components**: shadcn/ui + Base UI
- **Form Handling**: TanStack React Form
- **Code Quality**: Ultracite (Biome-based linter/formatter)

### Widget
- **Framework**: React 19 with TypeScript
- **Isolation**: Shadow DOM for CSS encapsulation
- **Styling**: TailwindCSS v4
- **Communication**: SSE streaming for real-time chat

### DevOps
- **Containerization**: Docker + Docker Compose
- **Package Manager**: Node.js/Bun
- **Development Server**: Hot reload via `air` (Go) and Vite (JS)

---

## Directory Structure & Key Directories

```
texly.chat/
├── cmd/                          # Application entry points
│   ├── app/                      # Main backend server (cmd/app/main.go)
│   └── ui-types/                 # TypeScript type generation utility
├── internal/                      # Core application logic
│   ├── db/                        # Database connection & migrations
│   ├── models/                    # Data models (User, Bot, Source, Chat, etc.)
│   ├── handlers/                  # HTTP request handlers
│   │   ├── auth/                  # Authentication endpoints
│   │   ├── bot/                   # Bot CRUD operations
│   │   ├── chat/                  # Chat streaming endpoint
│   │   ├── source/                # Source management (URL, file, text)
│   │   ├── user/                  # User profile endpoints
│   │   └── public/                # Public widget API endpoints
│   ├── middleware/                # HTTP middleware (CORS, auth, widget-CORS)
│   ├── repo/                      # Repository layer (data access)
│   │   ├── bot/                   # Bot repository
│   │   ├── source/                # Source repository
│   │   ├── user/                  # User repository
│   │   └── vector/                # Vector search repository
│   ├── services/                  # Business logic & external integrations
│   │   ├── chat/                  # RAG-powered chat with streaming
│   │   ├── embedding/             # OpenAI embeddings service
│   │   ├── extraction/            # PDF, Excel, text extraction
│   │   ├── chunker/               # Text chunking for embeddings
│   │   ├── scraper/               # Web scraping service
│   │   ├── storage/               # MinIO file storage service
│   │   ├── session/               # Session management for widget
│   │   └── vector/                # Vector search service
│   ├── queue/                     # Job queue implementation (in-memory)
│   ├── worker/                    # Background job processor
│   ├── server/                    # HTTP server setup & route registration
│   └── shared/                    # Shared utilities & test helpers
├── ui/                            # React dashboard frontend
│   ├── src/
│   │   ├── routes/                # File-based routes (TanStack Router)
│   │   │   ├── dashboard/         # Main app layout & pages
│   │   │   ├── _auth/             # Authentication pages
│   │   │   └── __root.tsx         # Root layout
│   │   ├── components/            # React components
│   │   ├── api/                   # API client & types
│   │   ├── hooks/                 # Custom React hooks
│   │   ├── providers/             # Context providers & app setup
│   │   ├── stores/                # Zustand stores (client state)
│   │   ├── lib/                   # Utilities & helpers
│   │   └── styles.css             # Global styles
│   ├── public/                    # Static assets
│   └── package.json
├── widget/                        # Embeddable chat widget
│   ├── src/
│   │   ├── components/            # Widget UI components
│   │   ├── api/                   # Public API client
│   │   ├── types.ts               # TypeScript types
│   │   └── index.tsx              # Entry point with Shadow DOM setup
│   └── package.json
├── configs/                       # Configuration management
│   └── config.go                  # Load env vars into typed Config struct
├── scripts/                       # Utility scripts
├── task/                          # Project planning & task documentation
│   ├── 00-overview-architecture.md
│   ├── 00-overview-tech_stack.md
│   └── [phase-specific tasks]
├── data/                          # SQLite database (dev.db)
├── dist/                          # Build output (binary)
├── Makefile                       # Build & development commands
├── Dockerfile                     # Container image definition
├── docker-compose.yml             # Multi-container setup (app + MinIO)
├── go.mod & go.sum                # Go dependencies
├── package.json & package-lock.json # Root-level npm config
└── README.md                      # Project overview

```

### Key Directories Explained

**`internal/`** - Contains all backend business logic:
- **models**: GORM-based database models with automatic ID generation (UUID)
- **handlers**: HTTP controllers structured by domain (auth, bot, chat, source, user, public)
- **services**: Orchestrate complex business logic and external integrations
- **repo**: Database access layer using GORM
- **middleware**: HTTP middleware for auth, CORS, logging
- **queue & worker**: Async job processing for scraping, embedding, extraction
- **db**: Database initialization, migrations, and vector table setup

**`ui/`** - React frontend dashboard:
- File-based routing with TanStack Router
- `dashboard/` contains main app pages after authentication
- `_auth/` contains login/signup pages
- Component-driven architecture with shadcn/ui

**`widget/`** - Standalone embeddable chat widget:
- Builds as a self-contained `texly-widget.js` file
- Shadow DOM for CSS isolation
- Session management for conversation persistence
- Communicates via public API endpoints

**`configs/`** - Centralized configuration:
- Single `config.go` file loads environment variables with sensible defaults
- Supports both development and production configurations via `.env` files

---

## Build & Development Commands

All commands are defined in `Makefile` and can be run with `make <command>`:

### Development

```bash
make dev           # Start both frontend and backend servers with hot reload
make dev-ui        # Start frontend only (Vite) on port 3000
make dev-api       # Start backend only (Go with air) on port 8080
make install       # Install all dependencies (Go modules + npm)
```

### Building

```bash
make build         # Build complete production binary with embedded UI
make build-ui      # Build React dashboard frontend
make build-widget  # Build embeddable chat widget
make build-api     # Build Go backend binary

# Individual steps during build
make clean         # Remove all build artifacts
```

### Testing & Code Quality

```bash
make test          # Run all Go tests with verbose output
make test-coverage # Run tests with coverage report
make fmt           # Format code (Go fmt + Biome/Ultracite for JS/TS)
make ui-types      # Generate TypeScript types from Go endpoints
```

### Docker Operations

```bash
make docker-up     # Start full stack (App + MinIO) with Docker Compose
make docker-down   # Stop all Docker containers
make docker-logs   # View Docker container logs
make docker-build  # Build Docker image only
make docker-clean  # Remove containers, volumes, and images
```

### Available npm/bun Scripts

**In `ui/` directory:**
```bash
bun dev            # Start Vite dev server (port 3000)
bun build          # Build for production
bun fix            # Format with Ultracite/Biome
bun check          # Check code style
```

**In `widget/` directory:**
```bash
bun dev            # Start widget dev server
bun build          # Build widget to dist/texly-widget.js
```

---

## High-Level Architecture & Data Flow

### System Architecture Diagram
```
┌─────────────────────────┐         ┌──────────────────┐
│  React Dashboard        │         │  Embeddable      │
│  (TanStack Router)      │         │  Widget (Shadow) │
└────────────┬────────────┘         └────────┬─────────┘
             │                               │
             └───────────────┬───────────────┘
                             │
                    ┌────────▼────────┐
                    │  Go/Gin API     │
                    │  (Port 8080)    │
                    └────────┬────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
    ┌────▼────┐      ┌──────▼──────┐      ┌────▼─────┐
    │ SQLite  │      │  MinIO S3   │      │ OpenAI   │
    │ + vec   │      │  Storage    │      │ API      │
    │ (WAL)   │      │             │      │          │
    └────┬────┘      └─────────────┘      └──────────┘
         │
    ┌────▼────────────────────┐
    │ Background Worker Pool  │
    │ (Goroutine + Queue)     │
    │ - Scrape URLs           │
    │ - Extract files         │
    │ - Generate embeddings   │
    │ - Chunk content         │
    └─────────────────────────┘
```

### Key Data Flow Patterns

#### 1. **Ingestion Pipeline (URL/File → Vector DB)**
```
User uploads URL/File
    ↓
Handler validates & creates Source (status: pending)
    ↓
Job enqueued to internal queue
    ↓
Worker picks up job
    ↓
Extract content (scrape URL or parse file)
    ↓
Chunk text (~800 tokens per chunk)
    ↓
Generate embeddings via OpenAI API
    ↓
Store vectors in sqlite-vec & metadata in document_chunks table
    ↓
Update Source status to active
```

#### 2. **Chat & RAG Flow**
```
User sends message to widget/dashboard
    ↓
Chat handler receives message + bot_id
    ↓
ChatService performs vector search:
  - Convert user message to embedding
  - Query sqlite-vec with cosine similarity
  - Retrieve top 5 relevant chunks
    ↓
Build system prompt with context chunks
    ↓
Stream chat completion from OpenAI API
    ↓
SSE stream tokens back to client
    ↓
Client displays streaming response
```

#### 3. **Widget Session Flow**
```
Widget script loaded on third-party site
    ↓
Fetch bot config (public API)
    ↓
Create session (returns session_id)
    ↓
User sends message
    ↓
POST to /api/public/chats/:session_id/messages
    ↓
Backend validates CORS via allowed_origins
    ↓
Process chat with RAG + streaming
    ↓
SSE streams response
```

### Architectural Patterns

#### Repository Pattern
- Data access encapsulated in `repo/` packages
- Each entity (User, Bot, Source, Vector) has its own repository
- Repositories use GORM for type-safe queries

#### Service Layer Pattern
- Business logic abstracted into services
- Services coordinate multiple repositories and external APIs
- Examples: `ChatService`, `EmbeddingService`, `StorageService`

#### Handler/Controller Pattern
- HTTP handlers in `handlers/` receive requests
- Handlers delegate to services
- Responses marshaled to JSON

#### Middleware Chain
- Authentication middleware validates JWT tokens
- CORS middleware handles cross-origin requests
- Widget-specific CORS checks allowed_origins

#### Job Queue Pattern
- In-memory queue with buffering (capacity: 100)
- Worker pool with configurable workers (default: 3)
- Jobs contain source ID and are processed sequentially

#### Streaming Pattern
- Chat responses use Server-Sent Events (SSE)
- Go channels used for streaming tokens
- Frontend reconstructs response from token stream

---

## Important Patterns & Conventions

### Go Code Patterns

**Dependency Injection**: Services and handlers receive dependencies via constructor functions:
```go
func New(db *gorm.DB, cfg configs.Config) *Server {
    // ...
}
```

**Error Handling**: Errors wrapped with context using `fmt.Errorf`:
```go
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

**GORM Models**: Each model implements `BeforeCreate` hook for UUID generation:
```go
func (b *Bot) BeforeCreate(tx *gorm.DB) (err error) {
    if b.ID == "" {
        b.ID = uuid.New().String()
    }
    return
}
```

**JSON Marshaling**: Struct tags control API serialization:
```go
type Bot struct {
    ID    string `json:"id" gorm:"primaryKey"`
    Name  string `json:"name" gorm:"not null"`
}
```

### TypeScript/React Patterns

**Type Generation**: TypeScript types generated from Go handlers via `make ui-types`
- Ensures API contract consistency
- Prevents runtime mismatches

**TanStack Router**: File-based routing with type-safe navigation
- Routes auto-generated from file structure
- Supports nested layouts and protected routes

**TanStack Query**: Server state management
- Automatic caching and synchronization
- Optimistic updates for better UX

**Zustand Stores**: Lightweight client state
- Global state without Redux complexity
- Used for UI state (modals, sidebar, etc.)

**Component Structure**: Atomic design with shadcn/ui components
- Composed UI from reusable base components
- Consistent theming with Tailwind

### Database Patterns

**Soft Deletes**: Models include `DeletedAt` field:
```go
DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
```

**Indexing**: Foreign keys automatically indexed for performance
```go
UserID string `gorm:"not null;index"`
```

**JSON Columns**: Complex data stored as JSON:
```go
WidgetConfig string `gorm:"type:text"` // JSON-encoded WidgetConfig
```

**Vector Tables**: Virtual `sqlite-vec` tables for embeddings
- Separate from traditional relational schema
- Joined with metadata via rowid

---

## Environment Configuration

Configuration is managed via `.env.local` (development) or `.env.prd` (production).

### Key Configuration Variables

```bash
# Server
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL=./data/dev.db

# JWT Authentication
JWT_SECRET=your-super-secret-key-change-this-in-production

# OpenAI Configuration
OPENAI_API_KEY=your-openai-api-key-here
EMBEDDING_MODEL=text-embedding-3-small        # Default
EMBEDDING_DIMENSION=1536                      # For text-embedding-3-small
OPENAI_CHAT_MODEL=gpt-4o-mini                 # Chat model
CHAT_TEMPERATURE=1.0                          # Randomness (0-2)
MAX_CONTEXT_CHUNKS=5                          # RAG context window

# MinIO Object Storage
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=texly-uploads
MINIO_USE_SSL=false

# Upload Limits
MAX_UPLOAD_SIZE_MB=100
```

**Default Values**: All configuration has sensible defaults. If `.env.local` is missing, the application uses defaults (SQLite at `data/dev.db`, port 8080, etc.).

---

## Existing Documentation

### Root-Level Documentation
- **`README.md`**: Quick start guide, feature list, tech stack overview
- **`Makefile`**: Build & development command definitions
- **`task/` directory**: Project planning documents
  - `00-overview-architecture.md`: Detailed architecture diagrams & flow
  - `00-overview-tech_stack.md`: Technology justification
  - `task.md`: Current sprint tasks

### Frontend Documentation
- **`ui/.claude/CLAUDE.md`**: Frontend-specific code standards (Ultracite/Biome)
- **`ui/README.md`**: TanStack Start template info
- **`ui/AGENTS.md`**: Agent-specific guidance

### Widget Documentation
- **`widget/README.md`**: Embedding guide, configuration, usage examples

---

## Key Files to Understand

### Backend Entry & Setup
- **`cmd/app/main.go`**: Application entry point (config load, DB connect, server start)
- **`internal/server/server.go`**: HTTP server setup, route registration, service initialization
- **`configs/config.go`**: Configuration loading and defaults

### Core Services
- **`internal/services/chat/chat_service.go`**: RAG orchestration & OpenAI streaming
- **`internal/services/embedding/embedding.go`**: OpenAI embedding generation
- **`internal/services/vector/search.go`**: Vector similarity search
- **`internal/services/storage/minio.go`**: File upload/download management
- **`internal/worker/worker.go`**: Background job processor (scrape, extract, embed)

### Database & Models
- **`internal/db/db.go`**: Connection, migrations, vector table initialization
- **`internal/models/bot_model.go`**: Bot definition with widget config
- **`internal/models/source_model.go`**: URL/file/text sources
- **`internal/models/session_model.go`**: Chat session tracking

### API Handlers
- **`internal/handlers/chat/chat_handler.go`**: Chat streaming endpoint
- **`internal/handlers/source/source_handler.go`**: Source creation & management
- **`internal/handlers/public/public_handler.go`**: Widget public API
- **`internal/handlers/auth/auth_handler.go`**: Sign up/login/JWT

### Frontend Entry & Layout
- **`ui/src/router.tsx`**: TanStack Router setup
- **`ui/src/routes/__root.tsx`**: Root layout component
- **`ui/src/routes/dashboard/`**: Main app pages
- **`ui/src/providers/`**: Context providers (React Query, Router, etc.)

### Widget
- **`widget/src/index.tsx`**: Widget entry point with Shadow DOM setup
- **`widget/src/App.tsx`**: Widget app logic & state
- **`widget/src/api/`**: Public API client

---

## Development Workflow

### Starting Development

1. **Install dependencies**
   ```bash
   make install
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env.local
   # Edit .env.local and set OPENAI_API_KEY
   ```

3. **Start development servers**
   ```bash
   make dev
   # Frontend: http://localhost:3000
   # Backend: http://localhost:8080
   ```

4. **Hot reload**
   - Backend: Auto-reloaded via `air` on file changes
   - Frontend: Auto-reloaded via Vite HMR

### Running Full Stack with Docker

```bash
make docker-up
# App available at http://localhost:8080
# MinIO console at http://localhost:9001 (minioadmin/minioadmin)
```

### Database Migrations

Migrations are automatically run at startup via GORM's `AutoMigrate()` in `internal/db/db.go`. To add new models:
1. Create model struct with GORM tags
2. Add to `Migrate()` function
3. Restart application

### Adding New Endpoints

1. Create handler in appropriate `handlers/` subdirectory
2. Add repository methods if needed in `repo/`
3. Register routes in `internal/server/server.go`
4. Regenerate TypeScript types: `make ui-types`
5. Create frontend components to consume new endpoint

### Testing

```bash
make test              # Run all tests
make test-coverage     # Generate coverage report
```

Tests are co-located with implementation files (`*_test.go`).

---

## Debugging & Common Issues

### Backend Debugging
- Logs printed to stdout (useful in Docker: `make docker-logs`)
- Set `ENVIRONMENT=development` for verbose logging
- Use Go debugger with your IDE for breakpoints

### Frontend Debugging
- React DevTools browser extension
- TanStack Router DevTools (visible in dev mode)
- Network tab to inspect API calls and SSE streams

### Database Issues
- SQLite database file at `data/dev.db`
- For debugging, use `.sqlite` extension in VS Code
- WAL mode creates `.db-wal` and `.db-shm` files (normal)

### Vector Search Not Working
- Ensure OpenAI API key is set
- Check `sqlite-vec` extension loaded: `internal/db/db.go`
- Verify embeddings generated: check `document_chunks` table

### Widget Not Loading
- Verify bot exists and is published
- Check CORS: `allowed_origins` must include your domain
- Inspect browser console for CORS errors

---

## Performance Considerations

### Backend Optimization
- **Database**: WAL mode enabled for concurrent read/write
- **Indexing**: Foreign keys indexed by default
- **Embeddings**: Cached in vector table, no re-computation
- **Streaming**: SSE avoids response buffering

### Frontend Optimization
- **Code splitting**: TanStack Router provides route-level code splitting
- **Lazy loading**: Components lazy-loaded per route
- **Query caching**: TanStack Query caches API responses
- **CSS**: TailwindCSS v4 tree-shakes unused styles

### Scaling Considerations
- **Database**: Single SQLite instance suitable for ~1000 concurrent connections
- **Vector Search**: `sqlite-vec` suitable for 10M+ vectors on single machine
- **File Storage**: MinIO can be scaled horizontally
- **Worker Pool**: Configurable worker count (default: 3)
- **Production**: Consider migrating to PostgreSQL + pgvector for multi-node

---

## Security Considerations

### Authentication
- JWT tokens with configurable secret
- Tokens validated on protected endpoints
- Passwords hashed (bcrypt or similar)

### Authorization
- Users can only access their own bots
- Widget has CORS validation via `allowed_origins`

### Data Privacy
- Soft deletes preserve historical data
- Vector embeddings stored locally (no cloud sync by default)
- Files stored in MinIO (can be self-hosted)

### Input Validation
- GORM binding validates request payloads
- File uploads size-limited (configurable)
- URL validation before scraping

---

## Code Quality Standards

### Go Code
- Follow `go fmt` for formatting
- Use `gofmt` and tools like `golangci-lint` for linting
- Write table-driven tests for complex logic
- Wrap errors with context

### TypeScript/React
- Enforce via **Ultracite** (Biome-based strict preset)
- Run `bun fix` before committing
- Type-safe components with explicit prop types
- No `any` type - use `unknown` instead
- Prefer functional components

### Documentation
- Comment exported functions (functions starting with uppercase)
- Document complex algorithms with multi-line comments
- Keep README and task files up-to-date
- Use meaningful variable/function names instead of comments

---

## Resources & Links

- **Go Docs**: https://golang.org/doc/
- **GORM**: https://gorm.io/docs/
- **Gin Framework**: https://gin-gonic.com/
- **TanStack Router**: https://tanstack.com/router/latest
- **React 19**: https://react.dev/
- **Tailwind CSS v4**: https://tailwindcss.com/docs/v4
- **OpenAI API**: https://platform.openai.com/docs/
- **sqlite-vec**: https://github.com/asg017/sqlite-vec
- **MinIO**: https://docs.min.io/

---

## Next Steps for New Contributors

1. Read the root README.md for quick overview
2. Review `task/00-overview-architecture.md` for system design
3. Set up development environment: `make install && make dev`
4. Pick a task from `task/task.md`
5. Explore relevant source files listed in "Key Files" section
6. Test changes locally before committing
7. Run `make fmt` and `bun fix` in ui/ before committing

---

Generated: 2025-02-06
Last Updated: Based on current codebase state
