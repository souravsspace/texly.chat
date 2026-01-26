# Texly Development Roadmap

## Executive Summary

Texly is an AI-powered customer support chatbot platform (texly.chat) designed as a SaaS alternative to SiteGPT. The current codebase provides foundational authentication and basic CRUD using **Go (Backend)** and **React 19/TanStack Start (Frontend in SPA mode)**. This roadmap outlines the strategic path from MVP to enterprise-grade platform, with all business logic implemented in Go and the frontend consuming REST APIs.

---

## Current State Analysis

### Backend Stack
- **Language**: Go 1.25.2
- **Web Framework**: Gin (HTTP routing)
- **ORM**: GORM
- **Database**: SQLite (needs migration to PostgreSQL)
- **Authentication**: JWT-based auth implemented

### Frontend Stack
- **Language**: TypeScript
- **Framework**: React 19
- **Router**: TanStack Start (SPA mode)
- **Styling**: Tailwind CSS v4
- **Component Library**: Shadcn UI
- **Data Fetching**: Will consume Go REST APIs

### Implemented Features
✅ User registration and authentication  
✅ JWT token management  
✅ Basic user profile endpoints  
✅ Placeholder "Post" CRUD (to be replaced)  
✅ Docker configuration  

### Critical Gaps (95% of Product)
❌ Bot/Chatbot management  
❌ Training data ingestion (scraping, files)  
❌ Vector embeddings and storage  
❌ AI/LLM integration (openai gpt 5.2)
❌ RAG (Retrieval-Augmented Generation) logic  
❌ Chat interface and streaming  
❌ Embeddable widget  
❌ Analytics and monitoring  
❌ Payment integration (polar.sh) 
❌ Advanced features (lead capture, escalation, integrations)  

---

## Target Architecture

### Backend Architecture (Go)
```
┌─────────────────────────────────────────────────────────┐
│                     Gin HTTP Server                      │
├─────────────────────────────────────────────────────────┤
│  Handlers (API Controllers)                             │
│  - Auth  - Bots  - Sources  - Chat  - Analytics        │
├─────────────────────────────────────────────────────────┤
│  Services (Business Logic)                              │
│  - Scraper  - Parser  - Embeddings  - LLM  - Vector    │
├─────────────────────────────────────────────────────────┤
│  Models (GORM)                                          │
│  - User  - Bot  - Source  - Conversation  - Message    │
├─────────────────────────────────────────────────────────┤
│  Infrastructure                                          │
│  - PostgreSQL + pgvector  - Redis/Asynq  - OpenAI API  │
└─────────────────────────────────────────────────────────┘
```

### Frontend Architecture (TanStack Start SPA)
```
┌─────────────────────────────────────────────────────────┐
│              TanStack Start (SPA Mode)                   │
├─────────────────────────────────────────────────────────┤
│  Routes                                                  │
│  - /login  - /dashboard  - /bots  - /chat              │
├─────────────────────────────────────────────────────────┤
│  Components                                              │
│  - Auth Forms  - Bot List  - Chat Interface            │
├─────────────────────────────────────────────────────────┤
│  API Layer (fetch/axios)                                │
│  - Calls Go REST endpoints                              │
│  - JWT token management                                 │
└─────────────────────────────────────────────────────────┘
```

### Key Architectural Decisions
1. **Database**: Migrate from SQLite to **PostgreSQL with pgvector extension** for unified vector storage
2. **Job Queue**: **Redis + Asynq** for background tasks (scraping, embedding generation)
3. **AI Provider**: **OpenAI API** (GPT-4, embeddings) with abstraction layer for future providers
4. **Frontend-Backend**: **Pure REST API** communication, no GraphQL or tRPC
5. **Widget Deployment**: Separate lightweight vanilla JS build for embed widget

---

## Development Phases

## Phase 1: MVP Foundation (Weeks 1-6)

**Goal**: Users can create a chatbot, train it on a single URL, and chat with it.

**Success Criteria**:
- ✅ User creates a bot with custom name
- ✅ User adds a URL to scrape
- ✅ Bot answers questions based on scraped content
- ✅ Basic chat interface works

---

### 1.1 Database Migration & Core Models

**Priority**: 🔴 Critical  
**Complexity**: Medium  
**Duration**: 3 days  

#### Backend Tasks

**Step 1: Migrate to PostgreSQL**
```bash
# Files to modify:
- internal/config/database.go
- docker-compose.yml
- .env.example
```

**Actions**:
1. Update GORM driver from `sqlite` to `postgres`
2. Add PostgreSQL connection string to config
3. Update Docker Compose with PostgreSQL service
4. Add pgvector extension initialization

**Step 2: Create Core Models**
```go
// internal/models/bot.go
type Bot struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      uint      `gorm:"not null;index"`
    Name        string    `gorm:"not null"`
    Description string
    SystemPrompt string   `gorm:"type:text"`
    Model       string    `gorm:"default:'gpt-4o-mini'"`
    Temperature float32   `gorm:"default:0.7"`
    IsActive    bool      `gorm:"default:true"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    User        User      `gorm:"foreignKey:UserID"`
}

// internal/models/source.go
type Source struct {
    ID          uint      `gorm:"primaryKey"`
    BotID       uint      `gorm:"not null;index"`
    Type        string    `gorm:"not null"` // "url", "file", "text"
    SourceURL   string    
    FileName    string
    Status      string    `gorm:"default:'pending'"` // pending, processing, completed, failed
    TotalChunks int       `gorm:"default:0"`
    ErrorMsg    string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Bot         Bot       `gorm:"foreignKey:BotID"`
}

// internal/models/document_chunk.go
type DocumentChunk struct {
    ID         uint      `gorm:"primaryKey"`
    SourceID   uint      `gorm:"not null;index"`
    BotID      uint      `gorm:"not null;index"`
    Content    string    `gorm:"type:text;not null"`
    Embedding  pgvector.Vector `gorm:"type:vector(1536)"` // OpenAI embedding dimension
    Metadata   datatypes.JSON
    TokenCount int
    CreatedAt  time.Time
    Source     Source    `gorm:"foreignKey:SourceID"`
}
```

**Step 3: Auto-Migration**
```go
// internal/config/database.go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Bot{},
        &models.Source{},
        &models.DocumentChunk{},
    )
}
```

#### Frontend Tasks

**No frontend changes needed** - models are backend only.

---

### 1.2 Bot Management CRUD

**Priority**: 🔴 Critical  
**Complexity**: Medium  
**Duration**: 4 days  

#### Backend Tasks

**Step 1: Create Bot Handler**
```go
// internal/handlers/bot/handler.go
type BotHandler struct {
    db *gorm.DB
}

// POST /api/bots - Create bot
func (h *BotHandler) CreateBot(c *gin.Context)

// GET /api/bots - List user's bots
func (h *BotHandler) ListBots(c *gin.Context)

// GET /api/bots/:id - Get single bot
func (h *BotHandler) GetBot(c *gin.Context)

// PUT /api/bots/:id - Update bot
func (h *BotHandler) UpdateBot(c *gin.Context)

// DELETE /api/bots/:id - Delete bot (soft delete)
func (h *BotHandler) DeleteBot(c *gin.Context)
```

**Step 2: Add Routes**
```go
// internal/routes/router.go
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
    api := r.Group("/api")
    
    // Authenticated routes
    auth := api.Group("")
    auth.Use(middleware.AuthMiddleware())
    {
        botHandler := bot.NewBotHandler(db)
        auth.POST("/bots", botHandler.CreateBot)
        auth.GET("/bots", botHandler.ListBots)
        auth.GET("/bots/:id", botHandler.GetBot)
        auth.PUT("/bots/:id", botHandler.UpdateBot)
        auth.DELETE("/bots/:id", botHandler.DeleteBot)
    }
}
```

**Step 3: Validation & DTOs**
```go
// internal/dto/bot.go
type CreateBotRequest struct {
    Name         string  `json:"name" binding:"required,min=3,max=100"`
    Description  string  `json:"description" binding:"max=500"`
    SystemPrompt string  `json:"systemPrompt"`
    Model        string  `json:"model" binding:"omitempty,oneof=gpt-4 gpt-4-turbo gpt-3.5-turbo"`
    Temperature  float32 `json:"temperature" binding:"omitempty,min=0,max=2"`
}
```

#### Frontend Tasks

**Step 1: Create Bot API Client**
```typescript
// ui/src/lib/api/bots.ts
export interface Bot {
  id: number;
  name: string;
  description: string;
  systemPrompt: string;
  model: string;
  temperature: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export const botsApi = {
  create: async (data: CreateBotRequest): Promise<Bot> => {
    const res = await fetch('/api/bots', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      body: JSON.stringify(data)
    });
    return res.json();
  },
  
  list: async (): Promise<Bot[]> => {
    const res = await fetch('/api/bots', {
      headers: { 'Authorization': `Bearer ${getToken()}` }
    });
    return res.json();
  },
  
  // ... other CRUD methods
};
```

**Step 2: Create Dashboard Route**
```typescript
// ui/src/routes/dashboard.tsx
import { createFileRoute } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import { botsApi } from '@/lib/api/bots';

export const Route = createFileRoute('/dashboard')({
  component: Dashboard
});

function Dashboard() {
  const { data: bots, isLoading } = useQuery({
    queryKey: ['bots'],
    queryFn: botsApi.list
  });

  return (
    <div className="container mx-auto p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Your Chatbots</h1>
        <CreateBotDialog />
      </div>
      
      {isLoading ? (
        <BotsSkeleton />
      ) : (
        <div className="grid grid-cols-3 gap-4">
          {bots?.map(bot => (
            <BotCard key={bot.id} bot={bot} />
          ))}
        </div>
      )}
    </div>
  );
}
```

**Step 3: Create Bot Form Component**
```typescript
// ui/src/components/create-bot-dialog.tsx
import { useState } from 'react';
import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

export function CreateBotDialog() {
  const [open, setOpen] = useState(false);
  
  const mutation = useMutation({
    mutationFn: botsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries(['bots']);
      setOpen(false);
    }
  });

  // Form implementation...
}
```

---

### 1.3 URL Scraper Service

**Priority**: 🔴 Critical  
**Complexity**: Complex  
**Duration**: 5 days  

#### Backend Tasks

**Step 1: Install Dependencies**
```bash
go get github.com/gocolly/colly/v2
go get github.com/PuerkitoBio/goquery
```

**Step 2: Create Scraper Service**
```go
// internal/services/scraper/scraper.go
package scraper

import (
    "github.com/gocolly/colly/v2"
    "strings"
)

type Scraper struct {
    maxDepth int
    timeout  time.Duration
}

type ScrapedContent struct {
    URL      string
    Title    string
    Content  string
    Metadata map[string]string
}

func NewScraper() *Scraper {
    return &Scraper{
        maxDepth: 3,
        timeout:  30 * time.Second,
    }
}

func (s *Scraper) ScrapeURL(url string) (*ScrapedContent, error) {
    c := colly.NewCollector(
        colly.MaxDepth(1),
        colly.Async(false),
    )
    
    var content ScrapedContent
    content.URL = url
    content.Metadata = make(map[string]string)
    
    // Extract title
    c.OnHTML("title", func(e *colly.HTMLElement) {
        content.Title = e.Text
    })
    
    // Extract main content
    c.OnHTML("body", func(e *colly.HTMLElement) {
        // Remove script and style tags
        e.DOM.Find("script, style, nav, footer, iframe").Remove()
        content.Content = cleanText(e.Text)
    })
    
    // Extract metadata
    c.OnHTML("meta[name='description']", func(e *colly.HTMLElement) {
        content.Metadata["description"] = e.Attr("content")
    })
    
    err := c.Visit(url)
    if err != nil {
        return nil, fmt.Errorf("failed to scrape: %w", err)
    }
    
    return &content, nil
}

func cleanText(text string) string {
    // Remove extra whitespace
    text = strings.TrimSpace(text)
    lines := strings.Split(text, "\n")
    
    var cleaned []string
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" {
            cleaned = append(cleaned, line)
        }
    }
    
    return strings.Join(cleaned, "\n")
}
```

**Step 3: Create Source Handler**
```go
// internal/handlers/source/handler.go
type SourceHandler struct {
    db      *gorm.DB
    scraper *scraper.Scraper
    queue   *asynq.Client
}

// POST /api/bots/:id/sources
func (h *SourceHandler) AddSource(c *gin.Context) {
    var req dto.AddSourceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    botID := c.Param("id")
    
    // Create source record
    source := &models.Source{
        BotID:     parseBotID(botID),
        Type:      req.Type,
        SourceURL: req.URL,
        Status:    "pending",
    }
    
    if err := h.db.Create(source).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to create source"})
        return
    }
    
    // Enqueue processing job
    task := tasks.NewProcessSourceTask(source.ID)
    _, err := h.queue.Enqueue(task)
    if err != nil {
        // Update status to failed
        h.db.Model(source).Update("status", "failed")
        c.JSON(500, gin.H{"error": "Failed to queue processing"})
        return
    }
    
    c.JSON(200, source)
}
```

**Step 4: Create Background Worker**
```go
// internal/worker/tasks/process_source.go
func ProcessSourceTask(ctx context.Context, task *asynq.Task) error {
    var payload SourcePayload
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return err
    }
    
    // Get source from DB
    var source models.Source
    db.First(&source, payload.SourceID)
    
    // Update status
    db.Model(&source).Update("status", "processing")
    
    // Scrape content
    scraper := scraper.NewScraper()
    content, err := scraper.ScrapeURL(source.SourceURL)
    if err != nil {
        db.Model(&source).Updates(map[string]interface{}{
            "status": "failed",
            "error_msg": err.Error(),
        })
        return err
    }
    
    // Pass to embedding service
    embeddingService := embedding.NewService()
    chunks, err := embeddingService.ProcessContent(content.Content, source.ID, source.BotID)
    if err != nil {
        return err
    }
    
    // Update source
    db.Model(&source).Updates(map[string]interface{}{
        "status": "completed",
        "total_chunks": len(chunks),
    })
    
    return nil
}
```

#### Frontend Tasks

**Step 1: Create Add Source Dialog**
```typescript
// ui/src/components/add-source-dialog.tsx
export function AddSourceDialog({ botId }: { botId: number }) {
  const [url, setUrl] = useState('');
  
  const mutation = useMutation({
    mutationFn: (url: string) => sourcesApi.add(botId, { type: 'url', url }),
    onSuccess: () => {
      toast.success('Source added and processing started');
    }
  });
  
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Add Training Data</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Training Source</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <Input
            placeholder="https://example.com/docs"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
          />
          <Button onClick={() => mutation.mutate(url)}>
            Add URL
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
```

**Step 2: Create Sources List Component**
```typescript
// ui/src/components/sources-list.tsx
export function SourcesList({ botId }: { botId: number }) {
  const { data: sources } = useQuery({
    queryKey: ['sources', botId],
    queryFn: () => sourcesApi.list(botId),
    refetchInterval: 5000 // Poll every 5s for status updates
  });
  
  return (
    <div className="space-y-2">
      {sources?.map(source => (
        <div key={source.id} className="flex items-center justify-between p-4 border rounded">
          <div>
            <p className="font-medium">{source.sourceUrl || source.fileName}</p>
            <p className="text-sm text-gray-500">{source.totalChunks} chunks</p>
          </div>
          <Badge variant={getStatusVariant(source.status)}>
            {source.status}
          </Badge>
        </div>
      ))}
    </div>
  );
}
```

---

### 1.4 Vector Embedding Pipeline

**Priority**: 🔴 Critical  
**Complexity**: Complex  
**Duration**: 6 days  

#### Backend Tasks

**Step 1: Install OpenAI SDK**
```bash
go get github.com/sashabaranov/go-openai
go get github.com/pgvector/pgvector-go
```

**Step 2: Create Embedding Service**
```go
// internal/services/embedding/service.go
package embedding

import (
    "context"
    "github.com/sashabaranov/go-openai"
)

type Service struct {
    client *openai.Client
    db     *gorm.DB
}

const (
    ChunkSize    = 1000  // ~500 tokens
    ChunkOverlap = 200
)

func NewService(apiKey string, db *gorm.DB) *Service {
    return &Service{
        client: openai.NewClient(apiKey),
        db:     db,
    }
}

func (s *Service) ProcessContent(content string, sourceID, botID uint) ([]*models.DocumentChunk, error) {
    // Chunk the content
    chunks := s.chunkText(content)
    
    var documentChunks []*models.DocumentChunk
    
    for i, chunk := range chunks {
        // Generate embedding
        embedding, err := s.generateEmbedding(chunk)
        if err != nil {
            return nil, fmt.Errorf("failed to generate embedding: %w", err)
        }
        
        // Create document chunk
        docChunk := &models.DocumentChunk{
            SourceID:   sourceID,
            BotID:      botID,
            Content:    chunk,
            Embedding:  pgvector.NewVector(embedding),
            TokenCount: estimateTokens(chunk),
            Metadata:   datatypes.JSON([]byte(fmt.Sprintf(`{"chunk_index": %d}`, i))),
        }
        
        if err := s.db.Create(docChunk).Error; err != nil {
            return nil, err
        }
        
        documentChunks = append(documentChunks, docChunk)
    }
    
    return documentChunks, nil
}

func (s *Service) generateEmbedding(text string) ([]float32, error) {
    resp, err := s.client.CreateEmbeddings(
        context.Background(),
        openai.EmbeddingRequestStrings{
            Input: []string{text},
            Model: openai.AdaEmbeddingV2,
        },
    )
    if err != nil {
        return nil, err
    }
    
    return resp.Data[0].Embedding, nil
}

func (s *Service) chunkText(text string) []string {
    // Simple chunking by character count
    var chunks []string
    runes := []rune(text)
    
    for i := 0; i < len(runes); i += (ChunkSize - ChunkOverlap) {
        end := i + ChunkSize
        if end > len(runes) {
            end = len(runes)
        }
        chunks = append(chunks, string(runes[i:end]))
        
        if end == len(runes) {
            break
        }
    }
    
    return chunks
}

func estimateTokens(text string) int {
    // Rough estimation: 1 token ≈ 4 characters
    return len(text) / 4
}
```

**Step 3: Create Vector Search Service**
```go
// internal/services/vector/search.go
package vector

import (
    "fmt"
    "github.com/pgvector/pgvector-go"
)

type SearchService struct {
    db              *gorm.DB
    embeddingService *embedding.Service
}

func NewSearchService(db *gorm.DB, embSvc *embedding.Service) *SearchService {
    return &SearchService{
        db:              db,
        embeddingService: embSvc,
    }
}

func (s *SearchService) SearchSimilar(botID uint, query string, limit int) ([]*models.DocumentChunk, error) {
    // Generate embedding for query
    queryEmbedding, err := s.embeddingService.generateEmbedding(query)
    if err != nil {
        return nil, err
    }
    
    var chunks []*models.DocumentChunk
    
    // Perform vector similarity search using pgvector
    err = s.db.
        Where("bot_id = ?", botID).
        Order(fmt.Sprintf("embedding <=> '%s'", pgvector.NewVector(queryEmbedding))).
        Limit(limit).
        Find(&chunks).
        Error
    
    if err != nil {
        return nil, err
    }
    
    return chunks, nil
}
```

**Step 4: Update Worker to Use Embedding Service**
```go
// internal/worker/tasks/process_source.go (update)
func ProcessSourceTask(ctx context.Context, task *asynq.Task) error {
    // ... existing scraping code ...
    
    // NEW: Process embeddings
    embeddingService := embedding.NewService(os.Getenv("OPENAI_API_KEY"), db)
    chunks, err := embeddingService.ProcessContent(content.Content, source.ID, source.BotID)
    if err != nil {
        db.Model(&source).Updates(map[string]interface{}{
            "status": "failed",
            "error_msg": err.Error(),
        })
        return err
    }
    
    // Update source
    db.Model(&source).Updates(map[string]interface{}{
        "status": "completed",
        "total_chunks": len(chunks),
    })
    
    return nil
}
```

#### Frontend Tasks

**No immediate frontend changes** - embedding happens in background. Frontend only sees source status updates (already implemented in 1.3).

---

### 1.5 Chat Interface & RAG Logic

**Priority**: 🔴 Critical  
**Complexity**: Complex  
**Duration**: 7 days  

#### Backend Tasks

**Step 1: Create Conversation Models**
```go
// internal/models/conversation.go
type Conversation struct {
    ID        uint      `gorm:"primaryKey"`
    BotID     uint      `gorm:"not null;index"`
    SessionID string    `gorm:"uniqueIndex;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Bot       Bot       `gorm:"foreignKey:BotID"`
    Messages  []Message `gorm:"foreignKey:ConversationID"`
}

type Message struct {
    ID             uint      `gorm:"primaryKey"`
    ConversationID uint      `gorm:"not null;index"`
    Role           string    `gorm:"not null"` // "user" or "assistant"
    Content        string    `gorm:"type:text;not null"`
    TokenCount     int
    SourceChunks   datatypes.JSON // IDs of chunks used
    CreatedAt      time.Time
    Conversation   Conversation `gorm:"foreignKey:ConversationID"`
}
```

**Step 2: Create LLM Service**
```go
// internal/services/llm/service.go
package llm

import (
    "context"
    "github.com/sashabaranov/go-openai"
)

type Service struct {
    client *openai.Client
}

type ChatRequest struct {
    Model       string
    Messages    []openai.ChatCompletionMessage
    Temperature float32
    Stream      bool
}

func NewService(apiKey string) *Service {
    return &Service{
        client: openai.NewClient(apiKey),
    }
}

func (s *Service) Chat(ctx context.Context, req *ChatRequest) (string, error) {
    resp, err := s.client.CreateChatCompletion(
        ctx,
        openai.ChatCompletionRequest{
            Model:       req.Model,
            Messages:    req.Messages,
            Temperature: req.Temperature,
            MaxTokens:   1000,
        },
    )
    
    if err != nil {
        return "", err
    }
    
    return resp.Choices[0].Message.Content, nil
}

func (s *Service) ChatStream(ctx context.Context, req *ChatRequest) (*openai.ChatCompletionStream, error) {
    stream, err := s.client.CreateChatCompletionStream(
        ctx,
        openai.ChatCompletionRequest{
            Model:       req.Model,
            Messages:    req.Messages,
            Temperature: req.Temperature,
            MaxTokens:   1000,
            Stream:      true,
        },
    )
    
    return stream, err
}
```

**Step 3: Create Chat Handler with RAG**
```go
// internal/handlers/chat/handler.go
type ChatHandler struct {
    db           *gorm.DB
    llmService   *llm.Service
    vectorSearch *vector.SearchService
}

type ChatRequest struct {
    Message   string `json:"message" binding:"required"`
    SessionID string `json:"sessionId"`
}

// POST /api/bots/:id/chat
func (h *ChatHandler) Chat(c *gin.Context) {
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    botID, _ := strconv.Atoi(c.Param("id"))
    
    // Get bot
    var bot models.Bot
    if err := h.db.First(&bot, botID).Error; err != nil {
        c.JSON(404, gin.H{"error": "Bot not found"})
        return
    }
    
    // Get or create conversation
    var conversation models.Conversation
    if req.SessionID != "" {
        h.db.Where("session_id = ?", req.SessionID).FirstOrCreate(&conversation, models.Conversation{
            BotID:     uint(botID),
            SessionID: req.SessionID,
        })
    } else {
        conversation = models.Conversation{
            BotID:     uint(botID),
            SessionID: generateSessionID(),
        }
        h.db.Create(&conversation)
    }
    
    // Save user message
    userMsg := models.Message{
        ConversationID: conversation.ID,
        Role:           "user",
        Content:        req.Message,
    }
    h.db.Create(&userMsg)
    
    // RAG: Search for relevant context
    relevantChunks, err := h.vectorSearch.SearchSimilar(uint(botID), req.Message, 5)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to search context"})
        return
    }
    
    // Build context
    var contextParts []string
    var chunkIDs []uint
    for _, chunk := range relevantChunks {
        contextParts = append(contextParts, chunk.Content)
        chunkIDs = append(chunkIDs, chunk.ID)
    }
    context := strings.Join(contextParts, "\n\n---\n\n")
    
    // Build messages for LLM
    messages := []openai.ChatCompletionMessage{
        {
            Role: openai.ChatMessageRoleSystem,
            Content: fmt.Sprintf(`%s

Use the following context to answer the user's question. If the context doesn't contain relevant information, say so politely.

Context:
%s`, bot.SystemPrompt, context),
        },
        {
            Role:    openai.ChatMessageRoleUser,
            Content: req.Message,
        },
    }
    
    // Get LLM response
    response, err := h.llmService.Chat(c.Request.Context(), &llm.ChatRequest{
        Model:       bot.Model,
        Messages:    messages,
        Temperature: bot.Temperature,
    })
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate response"})
        return
    }
    
    // Save assistant message
    chunkIDsJSON, _ := json.Marshal(chunkIDs)
    assistantMsg := models.Message{
        ConversationID: conversation.ID,
        Role:           "assistant",
        Content:        response,
        SourceChunks:   datatypes.JSON(chunkIDsJSON),
    }
    h.db.Create(&assistantMsg)
    
    c.JSON(200, gin.H{
        "message":   response,
        "sessionId": conversation.SessionID,
    })
}

// POST /api/bots/:id/chat/stream - Streaming version
func (h *ChatHandler) ChatStream(c *gin.Context) {
    // Similar to Chat() but uses Server-Sent Events
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    
    // ... RAG logic same as above ...
    
    stream, err := h.llmService.ChatStream(c.Request.Context(), &llm.ChatRequest{
        Model:       bot.Model,
        Messages:    messages,
        Temperature: bot.Temperature,
        Stream:      true,
    })
    
    if err != nil {
        c.SSEvent("error", err.Error())
        return
    }
    defer stream.Close()
    
    var fullResponse strings.Builder
    
    for {
        response, err := stream.Recv()
        if errors.Is(err, io.EOF) {
            break
        }
        if err != nil {
            c.SSEvent("error", err.Error())
            return
        }
        
        content := response.Choices[0].Delta.Content
        fullResponse.WriteString(content)
        
        c.SSEvent("message", content)
        c.Writer.Flush()
    }
    
    // Save complete message
    assistantMsg := models.Message{
        ConversationID: conversation.ID,
        Role:           "assistant",
        Content:        fullResponse.String(),
    }
    h.db.Create(&assistantMsg)
    
    c.SSEvent("done", conversation.SessionID)
}

func generateSessionID() string {
    return fmt.Sprintf("session_%d_%s", time.Now().Unix(), randomString(8))
}
```

**Step 4: Add Chat Routes**
```go
// internal/routes/router.go (add to existing routes)
chatHandler := chat.NewChatHandler(db, llmService, vectorSearch)
auth.POST("/bots/:id/chat", chatHandler.Chat)
auth.POST("/bots/:id/chat/stream", chatHandler.ChatStream)
auth.GET("/bots/:id/conversations", chatHandler.ListConversations)
auth.GET("/conversations/:id/messages", chatHandler.GetMessages)
```

#### Frontend Tasks

**Step 1: Create Chat API Client**
```typescript
// ui/src/lib/api/chat.ts
export interface ChatMessage {
  id: number;
  role: 'user' | 'assistant';
  content: string;
  createdAt: string;
}

export const chatApi = {
  send: async (botId: number, message: string, sessionId?: string) => {
    const res = await fetch(`/api/bots/${botId}/chat`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      body: JSON.stringify({ message, sessionId })
    });
    return res.json();
  },
  
  sendStream: async (
    botId: number, 
    message: string, 
    onChunk: (text: string) => void,
    sessionId?: string
  ) => {
    const res = await fetch(`/api/bots/${botId}/chat/stream`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`
      },
      body: JSON.stringify({ message, sessionId })
    });
    
    const reader = res.body?.getReader();
    const decoder = new TextDecoder();
    
    while (true) {
      const { done, value } = await reader!.read();
      if (done) break;
      
      const chunk = decoder.decode(value);
      const lines = chunk.split('\n');
      
      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6);
          if (data === '[DONE]') continue;
          onChunk(data);
        }
      }
    }
  }
};
```

**Step 2: Create Chat Interface Route**
```typescript
// ui/src/routes/bots/$botId/chat.tsx
import { createFileRoute } from '@tanstack/react-router';
import { useState, useRef, useEffect } from 'react';
import { chatApi } from '@/lib/api/chat';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { ScrollArea } from '@/components/ui/scroll-area';

export const Route = createFileRoute('/bots/$botId/chat')({
  component: ChatInterface
});

function ChatInterface() {
  const { botId } = Route.useParams();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [sessionId, setSessionId] = useState<string>();
  const scrollRef = useRef<HTMLDivElement>(null);

  const handleSend = async () => {
    if (!input.trim()) return;
    
    const userMessage = {
      id: Date.now(),
      role: 'user' as const,
      content: input,
      createdAt: new Date().toISOString()
    };
    
    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);
    
    // Create placeholder for assistant response
    const assistantMessage = {
      id: Date.now() + 1,
      role: 'assistant' as const,
      content: '',
      createdAt: new Date().toISOString()
    };
    setMessages(prev => [...prev, assistantMessage]);
    
    try {
      await chatApi.sendStream(
        Number(botId),
        userMessage.content,
        (chunk) => {
          setMessages(prev => {
            const newMessages = [...prev];
            const lastMsg = newMessages[newMessages.length - 1];
            lastMsg.content += chunk;
            return newMessages;
          });
        },
        sessionId
      );
    } catch (error) {
      console.error('Chat error:', error);
    } finally {
      setIsLoading(false);
    }
  };
  
  useEffect(() => {
    scrollRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  return (
    <div className="flex flex-col h-screen">
      <div className="border-b p-4">
        <h1 className="text-xl font-bold">Chat with Bot</h1>
      </div>
      
      <ScrollArea className="flex-1 p-4">
        <div className="space-y-4 max-w-3xl mx-auto">
          {messages.map((msg) => (
            <div
              key={msg.id}
              className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-[80%] rounded-lg p-3 ${
                  msg.role === 'user'
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-100 text-gray-900'
                }`}
              >
                <p className="whitespace-pre-wrap">{msg.content}</p>
              </div>
            </div>
          ))}
          <div ref={scrollRef} />
        </div>
      </ScrollArea>
      
      <div className="border-t p-4">
        <div className="max-w-3xl mx-auto flex gap-2">
          <Input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleSend()}
            placeholder="Type your message..."
            disabled={isLoading}
          />
          <Button onClick={handleSend} disabled={isLoading}>
            Send
          </Button>
        </div>
      </div>
    </div>
  );
}
```

**Step 3: Add Chat Link to Bot Card**
```typescript
// ui/src/components/bot-card.tsx
export function BotCard({ bot }: { bot: Bot }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{bot.name}</CardTitle>
        <CardDescription>{bot.description}</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex gap-2">
          <Button asChild>
            <Link to={`/bots/${bot.id}/chat`}>
              Open Chat
            </Link>
          </Button>
          <Button variant="outline" asChild>
            <Link to={`/bots/${bot.id}/settings`}>
              Settings
            </Link>
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
```

---

## Phase 2: Growth Features (Weeks 7-12)

**Goal**: Enhance the platform with embeddable widgets, file uploads, advanced scraping, and better UX.

**Success Criteria**:
- ✅ Users can embed chatbot on their websites
- ✅ Users can upload PDF/DOCX files
- ✅ Sitemap crawler works for entire websites
- ✅ Quick prompts guide users
- ✅ Daily email summaries sent

---

### 2.1 Embeddable Chat Widget

**Priority**: 🟡 High  
**Complexity**: High  
**Duration**: 6 days  

#### Backend Tasks

**Step 1: Create Public Chat Endpoint**
```go
// internal/handlers/chat/public.go
type PublicChatHandler struct {
    db           *gorm.DB
    llmService   *llm.Service
    vectorSearch *vector.SearchService
}

// POST /api/public/chat/:botId
func (h *PublicChatHandler) Chat(c *gin.Context) {
    // Similar to private chat but:
    // 1. No authentication required
    // 2. Rate limiting by IP
    // 3. Check if bot is public/active
    
    botID, _ := strconv.Atoi(c.Param("botId"))
    
    var bot models.Bot
    if err := h.db.Where("id = ? AND is_active = ?", botID, true).First(&bot).Error; err != nil {
        c.JSON(404, gin.H{"error": "Bot not found"})
        return
    }
    
    // Rest same as private chat endpoint
    // ... RAG + LLM logic ...
}
```

**Step 2: Add CORS Configuration**
```go
// main.go
import "github.com/gin-contrib/cors"

func main() {
    r := gin.Default()
    
    // Allow widget to be embedded anywhere
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"POST", "GET"},
        AllowHeaders:     []string{"Content-Type"},
        AllowCredentials: false,
    }))
    
    // Routes...
}
```

**Step 3: Create Widget Configuration Endpoint**
```go
// GET /api/bots/:id/widget-config
func (h *BotHandler) GetWidgetConfig(c *gin.Context) {
    botID := c.Param("id")
    
    var bot models.Bot
    if err := h.db.First(&bot, botID).Error; err != nil {
        c.JSON(404, gin.H{"error": "Bot not found"})
        return
    }
    
    config := map[string]interface{}{
        "botId":       bot.ID,
        "name":        bot.Name,
        "description": bot.Description,
        "theme": map[string]string{
            "primaryColor":   "#3B82F6",
            "backgroundColor": "#FFFFFF",
        },
    }
    
    c.JSON(200, config)
}
```

#### Frontend Tasks

**Step 1: Create Widget Build**
```typescript
// widget/src/index.ts
// Separate build configuration for widget
class TexlyWidget {
  private botId: string;
  private container: HTMLElement;
  private iframe: HTMLIFrameElement;
  
  constructor(botId: string) {
    this.botId = botId;
    this.init();
  }
  
  private init() {
    // Create floating chat button
    this.container = document.createElement('div');
    this.container.id = 'texly-widget';
    this.container.style.cssText = `
      position: fixed;
      bottom: 20px;
      right: 20px;
      z-index: 999999;
    `;
    
    const button = document.createElement('button');
    button.innerHTML = '💬';
    button.style.cssText = `
      width: 60px;
      height: 60px;
      border-radius: 50%;
      background: #3B82F6;
      border: none;
      cursor: pointer;
      font-size: 24px;
    `;
    
    button.onclick = () => this.toggle();
    this.container.appendChild(button);
    document.body.appendChild(this.container);
  }
  
  private toggle() {
    if (!this.iframe) {
      this.iframe = document.createElement('iframe');
      this.iframe.src = `https://texly.chat/embed/${this.botId}`;
      this.iframe.style.cssText = `
        position: fixed;
        bottom: 90px;
        right: 20px;
        width: 400px;
        height: 600px;
        border: none;
        border-radius: 12px;
        box-shadow: 0 4px 20px rgba(0,0,0,0.15);
        z-index: 999999;
      `;
      document.body.appendChild(this.iframe);
    } else {
      this.iframe.remove();
      this.iframe = null;
    }
  }
}

// Auto-initialize
(function() {
  const script = document.currentScript as HTMLScriptElement;
  const botId = script.getAttribute('data-bot-id');
  if (botId) {
    new TexlyWidget(botId);
  }
})();
```

**Step 2: Create Embed Chat Route**
```typescript
// ui/src/routes/embed/$botId.tsx
export const Route = createFileRoute('/embed/$botId')({
  component: EmbedChat
});

function EmbedChat() {
  const { botId } = Route.useParams();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  
  // Similar to regular chat but:
  // 1. Minimal UI (no navigation)
  // 2. Uses public API endpoint
  // 3. Stores session in localStorage
  
  return (
    <div className="h-screen flex flex-col bg-white">
      <div className="p-3 border-b">
        <h2 className="font-semibold text-sm">Chat Support</h2>
      </div>
      
      <ScrollArea className="flex-1 p-3">
        {/* Messages */}
      </ScrollArea>
      
      <div className="p-3 border-t">
        {/* Input */}
      </div>
    </div>
  );
}
```

**Step 3: Generate Embed Code UI**
```typescript
// ui/src/components/embed-code-dialog.tsx
export function EmbedCodeDialog({ bot }: { bot: Bot }) {
  const embedCode = `<script 
  src="https://texly.chat/widget.js" 
  data-bot-id="${bot.id}"
></script>`;
  
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">Get Embed Code</Button>
      </DialogTrigger>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Embed {bot.name} on Your Website</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <p className="text-sm text-gray-600">
            Copy and paste this code before the closing &lt;/body&gt; tag on your website:
          </p>
          <pre className="p-4 bg-gray-100 rounded-lg overflow-x-auto">
            <code>{embedCode}</code>
          </pre>
          <Button onClick={() => navigator.clipboard.writeText(embedCode)}>
            Copy to Clipboard
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
```

---

### 2.2 File Upload Training

**Priority**: 🟡 High  
**Complexity**: Medium  
**Duration**: 5 days  

#### Backend Tasks

**Step 1: Install File Processing Libraries**
```bash
go get github.com/ledongthuc/pdf
go get github.com/nguyenthenguyen/docx
go get code.sajari.com/docconv
```

**Step 2: Create File Parser Service**
```go
// internal/services/parser/parser.go
package parser

import (
    "github.com/ledongthuc/pdf"
    "code.sajari.com/docconv"
)

type Parser struct{}

func NewParser() *Parser {
    return &Parser{}
}

func (p *Parser) ParseFile(filepath string, fileType string) (string, error) {
    switch fileType {
    case "application/pdf":
        return p.parsePDF(filepath)
    case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
        return p.parseDOCX(filepath)
    case "text/plain":
        return p.parseTXT(filepath)
    default:
        return "", fmt.Errorf("unsupported file type: %s", fileType)
    }
}

func (p *Parser) parsePDF(filepath string) (string, error) {
    f, r, err := pdf.Open(filepath)
    if err != nil {
        return "", err
    }
    defer f.Close()
    
    var buf bytes.Buffer
    b, err := r.GetPlainText()
    if err != nil {
        return "", err
    }
    buf.ReadFrom(b)
    
    return buf.String(), nil
}

func (p *Parser) parseDOCX(filepath string) (string, error) {
    res, err := docconv.ConvertPath(filepath)
    if err != nil {
        return "", err
    }
    return res.Body, nil
}

func (p *Parser) parseTXT(filepath string) (string, error) {
    content, err := os.ReadFile(filepath)
    if err != nil {
        return "", err
    }
    return string(content), nil
}
```

**Step 3: Create File Upload Handler**
```go
// internal/handlers/source/file.go
// POST /api/bots/:id/sources/upload
func (h *SourceHandler) UploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": "No file uploaded"})
        return
    }
    
    // Validate file size (max 10MB)
    if file.Size > 10*1024*1024 {
        c.JSON(400, gin.H{"error": "File too large (max 10MB)"})
        return
    }
    
    // Validate file type
    allowedTypes := map[string]bool{
        "application/pdf":  true,
        "text/plain":       true,
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
    }
    
    if !allowedTypes[file.Header.Get("Content-Type")] {
        c.JSON(400, gin.H{"error": "Unsupported file type"})
        return
    }
    
    botID, _ := strconv.Atoi(c.Param("id"))
    
    // Save file
    filename := fmt.Sprintf("%d_%d_%s", botID, time.Now().Unix(), file.Filename)
    filepath := fmt.Sprintf("./uploads/%s", filename)
    
    if err := c.SaveUploadedFile(file, filepath); err != nil {
        c.JSON(500, gin.H{"error": "Failed to save file"})
        return
    }
    
    // Create source record
    source := &models.Source{
        BotID:    uint(botID),
        Type:     "file",
        FileName: filename,
        Status:   "pending",
    }
    h.db.Create(source)
    
    // Enqueue processing
    task := tasks.NewProcessFileTask(source.ID, filepath, file.Header.Get("Content-Type"))
    h.queue.Enqueue(task)
    
    c.JSON(200, source)
}
```

**Step 4: Create File Processing Worker**
```go
// internal/worker/tasks/process_file.go
func ProcessFileTask(ctx context.Context, task *asynq.Task) error {
    var payload FilePayload
    json.Unmarshal(task.Payload(), &payload)
    
    var source models.Source
    db.First(&source, payload.SourceID)
    db.Model(&source).Update("status", "processing")
    
    // Parse file
    parser := parser.NewParser()
    content, err := parser.ParseFile(payload.FilePath, payload.FileType)
    if err != nil {
        db.Model(&source).Updates(map[string]interface{}{
            "status":    "failed",
            "error_msg": err.Error(),
        })
        return err
    }
    
    // Generate embeddings (reuse existing service)
    embeddingService := embedding.NewService(os.Getenv("OPENAI_API_KEY"), db)
    chunks, err := embeddingService.ProcessContent(content, source.ID, source.BotID)
    if err != nil {
        return err
    }
    
    db.Model(&source).Updates(map[string]interface{}{
        "status":       "completed",
        "total_chunks": len(chunks),
    })
    
    return nil
}
```

#### Frontend Tasks

**Step 1: Create File Upload Component**
```typescript
// ui/src/components/file-upload.tsx
export function FileUpload({ botId }: { botId: number }) {
  const [uploading, setUploading] = useState(false);
  
  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    
    const formData = new FormData();
    formData.append('file', file);
    
    setUploading(true);
    try {
      await fetch(`/api/bots/${botId}/sources/upload`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${getToken()}`
        },
        body: formData
      });
      toast.success('File uploaded and processing started');
    } catch (error) {
      toast.error('Failed to upload file');
    } finally {
      setUploading(false);
    }
  };
  
  return (
    <div className="border-2 border-dashed rounded-lg p-8 text-center">
      <input
        type="file"
        accept=".pdf,.txt,.docx"
        onChange={handleFileUpload}
        disabled={uploading}
        className="hidden"
        id="file-upload"
      />
      <label htmlFor="file-upload" className="cursor-pointer">
        <div className="text-gray-600">
          <Upload className="w-12 h-12 mx-auto mb-4" />
          <p className="font-medium">Click to upload or drag and drop</p>
          <p className="text-sm">PDF, TXT, or DOCX (max 10MB)</p>
        </div>
      </label>
    </div>
  );
}
```

---

### 2.3 Sitemap Crawler

**Priority**: 🟡 High  
**Complexity**: High  
**Duration**: 6 days  

#### Backend Tasks

**Step 1: Create Sitemap Parser**
```go
// internal/services/scraper/sitemap.go
type SitemapURL struct {
    Loc        string
    LastMod    string
    ChangeFreq string
    Priority   float64
}

func (s *Scraper) ParseSitemap(sitemapURL string) ([]SitemapURL, error) {
    resp, err := http.Get(sitemapURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var sitemap struct {
        URLs []SitemapURL `xml:"url"`
    }
    
    if err := xml.Unmarshal(body, &sitemap); err != nil {
        return nil, err
    }
    
    return sitemap.URLs, nil
}
```

**Step 2: Create Bulk Scraping Worker**
```go
// internal/worker/tasks/process_sitemap.go
func ProcessSitemapTask(ctx context.Context, task *asynq.Task) error {
    var payload SitemapPayload
    json.Unmarshal(task.Payload(), &payload)
    
    var source models.Source
    db.First(&source, payload.SourceID)
    
    // Parse sitemap
    scraper := scraper.NewScraper()
    urls, err := scraper.ParseSitemap(source.SourceURL)
    if err != nil {
        db.Model(&source).Update("status", "failed")
        return err
    }
    
    db.Model(&source).Update("status", "processing")
    
    // Create child sources for each URL
    for _, url := range urls {
        childSource := models.Source{
            BotID:     source.BotID,
            Type:      "url",
            SourceURL: url.Loc,
            Status:    "pending",
        }
        db.Create(&childSource)
        
        // Enqueue scraping task
        task := tasks.NewProcessSourceTask(childSource.ID)
        queue.Enqueue(task)
    }
    
    db.Model(&source).Updates(map[string]interface{}{
        "status":       "completed",
        "total_chunks": len(urls),
    })
    
    return nil
}
```

---

### 2.4 Quick Prompts Feature

**Priority**: 🟢 Medium  
**Complexity**: Simple  
**Duration**: 2 days  

#### Backend Tasks

**Step 1: Add Quick Prompts to Bot Model**
```go
// internal/models/bot.go (update)
type Bot struct {
    // ... existing fields ...
    QuickPrompts datatypes.JSON `gorm:"type:jsonb"`
}

// Example JSON structure:
// ["How do I reset my password?", "What are your pricing plans?", "How do I contact support?"]
```

**Step 2: Update Bot Endpoints**
```go
// internal/dto/bot.go (update)
type UpdateBotRequest struct {
    // ... existing fields ...
    QuickPrompts []string `json:"quickPrompts"`
}
```

#### Frontend Tasks

**Step 1: Add Quick Prompts to Chat UI**
```typescript
// ui/src/routes/bots/$botId/chat.tsx (update)
function ChatInterface() {
  const { data: bot } = useQuery({
    queryKey: ['bot', botId],
    queryFn: () => botsApi.get(Number(botId))
  });
  
  const quickPrompts = bot?.quickPrompts || [];
  
  return (
    <div className="flex flex-col h-screen">
      {/* ... existing chat UI ... */}
      
      {messages.length === 0 && quickPrompts.length > 0 && (
        <div className="p-4 space-y-2">
          <p className="text-sm text-gray-600">Quick questions:</p>
          <div className="grid grid-cols-2 gap-2">
            {quickPrompts.map((prompt, i) => (
              <Button
                key={i}
                variant="outline"
                className="text-left h-auto py-3"
                onClick={() => setInput(prompt)}
              >
                {prompt}
              </Button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
```

---

### 2.5 Daily Email Summaries

**Priority**: 🟢 Medium  
**Complexity**: Medium  
**Duration**: 4 days  

#### Backend Tasks

**Step 1: Install Email Library**
```bash
go get gopkg.in/gomail.v2
```

**Step 2: Create Email Service**
```go
// internal/services/email/service.go
package email

import "gopkg.in/gomail.v2"

type Service struct {
    dialer *gomail.Dialer
}

func NewService(host string, port int, username, password string) *Service {
    return &Service{
        dialer: gomail.NewDialer(host, port, username, password),
    }
}

func (s *Service) SendDailySummary(to string, summary DailySummary) error {
    m := gomail.NewMessage()
    m.SetHeader("From", "noreply@texly.chat")
    m.SetHeader("To", to)
    m.SetHeader("Subject", "Your Daily Texly Summary")
    
    htmlBody := fmt.Sprintf(`
        <h2>Daily Summary for %s</h2>
        <p>Total conversations: %d</p>
        <p>Total messages: %d</p>
        <p>Most asked question: %s</p>
    `, summary.BotName, summary.TotalConversations, summary.TotalMessages, summary.TopQuestion)
    
    m.SetBody("text/html", htmlBody)
    
    return s.dialer.DialAndSend(m)
}
```

#### Frontend Tasks

**Step 1: Email Preferences UI**
```typescript
// ui/src/settings/notifications.tsx
export function NotificationSettings() {
  const { data: profile } = useQuery({ queryKey: ['profile'] });
  
  return (
    <div className="space-y-4">
      <h3 className="text-lg font-medium">Email Notifications</h3>
      <div className="flex items-center justify-between">
        <label>Daily Activity Summary</label>
        <Switch 
          checked={profile?.emailPreferences.dailySummary}
          onCheckedChange={(checked) => updatePreferences({ dailySummary: checked })}
        />
      </div>
    </div>
  );
}
```

---

## Phase 3: Scaling & Enterprise Features (Weeks 13-18)

**Goal**: Prepare the platform for scale with advanced analytics, team collaboration, and robust security.

**Success Criteria**:
- ✅ Admins can view detailed usage analytics
- ✅ Users can invite team members with specific roles
- ✅ System is protected against abuse (Rate Limiting)
- ✅ Audit logs track all critical actions

---

### 3.1 Analytics Dashboard

**Priority**: 🟡 High
**Complexity**: High
**Duration**: 5 days

#### Backend Tasks

**Step 1: Analytics Service**
```go
// internal/services/analytics/service.go
type AnalyticsService struct {
    db *gorm.DB
}

func (s *AnalyticsService) GetBotStats(botID uint, period string) (*BotStats, error) {
    // SQL aggregation for messages, tokens, and active users
    // SELECT DATE(created_at), COUNT(*) FROM messages ... GROUP BY DATE(created_at)
    return &BotStats{}, nil
}
```

#### Frontend Tasks

**Step 1: Dashboard Charts**
- Implement Recharts for visualization.
- Show "Messages per Day", "Token Usage", "Top Sources".

---

### 3.2 User Roles & Teams

**Priority**: 🟢 Medium
**Complexity**: Medium
**Duration**: 5 days

#### Backend Tasks

**Step 1: Team Models**
```go
// internal/models/team.go
type Team struct {
    ID    uint
    Name  string
    Users []User `gorm:"many2many:team_users;"`
}

type TeamUser struct {
    TeamID uint
    UserID uint
    Role   string // "admin", "editor", "viewer"
}
```

**Step 2: RBAC Middleware**
- Middleware to check if user has permission for requested resource based on role.

---

### 3.3 Advanced Security

**Priority**: 🟢 Medium
**Complexity**: Medium
**Duration**: 3 days

#### Backend Tasks

**Step 1: Rate Limiting (Redis)**
- Implement `redis_rate` for API endpoints.
- Limit: 100 req/min for free tier, 1000 req/min for pro.

**Step 2: Audit Logs**
- Log critical actions (delete bot, update API key, invite user).
- Store in `audit_logs` table.

---

## Phase 4: Monetization & Integration (Weeks 19-24)

**Goal**: Generate revenue and integrate with external workflows.

**Success Criteria**:
- ✅ Users can subscribe to Paid plans via Polar.sh
- ✅ Subscriptions enforce usage limits
- ✅ Users can connect bots to Slack and Discord

---

### 4.1 Polar.sh Subscription

**Priority**: 🔴 Critical
**Complexity**: High
**Duration**: 5 days

#### Backend Tasks

**Step 1: Install Polar SDK**
```bash
go get github.com/polarsource/polar-go
```

**Step 2: Webhook Handler**
```go
// internal/handlers/billing/polar.go
import "github.com/polarsource/polar-go"

func (h *BillingHandler) HandlePolarWebhook(c *gin.Context) {
    // Verify signature
    // Handle "subscription.created", "subscription.updated", "subscription.canceled"
    // Update user's plan in DB
}
```

**Step 3: Usage Enforcement**
- Middleware to check limits (e.g. "Max 5 bots") against current plan.

#### Frontend Tasks

**Step 1: Pricing Page**
- Display "Free", "Pro", "Enterprise" tiers.
- Link "Upgrade" buttons to Polar.sh checkout URLs.

---

### 4.2 Third-party Integrations

**Priority**: 🟡 High
**Complexity**: High
**Duration**: 6 days

#### Backend Tasks

**Step 1: Slack Integration**
- OAuth2 flow for "Add to Slack".
- Event subscription for mentioning bot.
- Post-back logic to answer in thread.

**Step 2: Discord Integration**
- Bot Token management.
- Discord Gateway (or webhook interaction).

---

## Future Outlook (Phase 5+)

- **Voice Support**: Twilio / ElevenLabs integration for phone support.
- **Mobile App**: React Native wrapper for on-the-go management.
- **Fine-tuning**: UI for fine-tuning open source models (Llama 3) on user data.