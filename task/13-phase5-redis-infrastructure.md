# Phase 5: Infrastructure & Performance (Redis Layer) ✅

## Goal
Implement Redis caching layer to handle 10k+ concurrent write requests and improve overall system performance.

---

## Problem Statement
- Need intelligent caching to reduce database load and improve performance
- Need rate limiting to prevent abuse
- Need session management for distributed deployments
- Need to handle 10k+ concurrent write requests efficiently

---

## Backend Tasks

### Step 1: Redis Setup & Configuration

#### 1.1 Add Redis Dependencies
- [x] Add `github.com/redis/go-redis/v9` to `go.mod`
- [x] Add Redis configuration to `configs/config.go`:
  ```go
  RedisURL      string // e.g., "redis://localhost:6379"
  RedisPassword string
  RedisDB       int
  RedisTTL      int // Default cache TTL in seconds
  ```
- [x] Add environment variables to `.env.local`:
  ```
  REDIS_URL=redis://localhost:6379
  REDIS_PASSWORD=
  REDIS_DB=0
  REDIS_TTL=3600
  ```

#### 1.2 Create Redis Client
- [x] Create `internal/db/redis.go`:
  - [x] Initialize Redis client with connection pooling
  - [x] Implement health check function
  - [x] Add graceful shutdown
  - [x] Export singleton client instance
- [x] Update `cmd/app/main.go` to initialize Redis on startup
- [x] Add Redis health check to server startup

---

### Step 2: Caching Layer Implementation

#### 2.1 Create Cache Service
- [x] Create `internal/services/cache/cache_service.go`:
  - [x] `Get(ctx, key) (string, error)` - Get cached value
  - [x] `Set(ctx, key, value, ttl) error` - Set cache with TTL
  - [x] `Delete(ctx, key) error` - Invalidate cache
  - [x] `DeletePattern(ctx, pattern) error` - Bulk invalidation
  - [x] `Exists(ctx, key) (bool, error)` - Check if key exists
  - [x] `GetJSON(ctx, key, dest) error` - Get and unmarshal JSON
  - [x] `SetJSON(ctx, key, value, ttl) error` - Marshal and set JSON

#### 2.2 Cache Keys Strategy
- [x] Define cache key patterns in `internal/services/cache/keys.go`:
  ```go
  const (
      BotCacheKey      = "bot:%s"           // bot:{bot_id}
      BotListCacheKey  = "bots:user:%s"     // bots:user:{user_id}
      SourceCacheKey   = "source:%s"        // source:{source_id}
      UserCacheKey     = "user:%s"          // user:{user_id}
      SessionCacheKey  = "session:%s"       // session:{session_id}
      VectorCacheKey   = "vector:%s:%s"     // vector:{bot_id}:{query_hash}
  )
  ```

#### 2.3 Implement Cache-Aside Pattern
- [x] Update Bot Repository (`internal/repo/bot/bot_repo.go`):
  - [x] Check cache before database read
  - [x] Store result in cache after database read
  - [x] Invalidate cache on update/delete
- [x] Update User Repository:
  - [x] Cache user profiles
  - [x] Cache authentication lookups
- [x] Update Source Repository:
  - [x] Cache source metadata
  - [x] Invalidate on status changes
- [x] Update Vector Search:
  - [x] Cache frequent query results (hash query as key)
  - [x] Set shorter TTL (5-10 minutes)

---

### Step 3: Write Buffering for High Load

#### 3.1 Create Write Buffer Service
- [x] Create `internal/services/writebuffer/buffer_service.go`:
  - [x] Use Redis Lists as write queues
  - [x] `QueueWrite(ctx, operation, data) error` - Add to buffer
  - [x] `ProcessBuffer(ctx) error` - Batch process writes
  - [x] Implement retry logic for failed writes
  - [x] Add dead letter queue for permanently failed writes

#### 3.2 Batch Write Worker
- [x] Create `internal/worker/write_worker.go`:
  - [x] Poll Redis queue every 100ms
  - [x] Batch up to 100 writes per transaction
  - [x] Use PostgreSQL transaction for atomic batch commit
  - [x] Log failures and retry up to 3 times
  - [x] Move to dead letter queue after max retries

#### 3.3 Update Chat Handler
- [x] Modify `internal/handlers/chat/chat_handler.go`:
  - [x] Queue chat message writes instead of direct DB write
  - [x] Return immediately to user (async write)
  - [x] Ensure message ID returned for tracking

---

### Step 4: Rate Limiting

#### 4.1 Create Rate Limiter Middleware
- [x] Create `internal/middleware/rate_limit.go`:
  - [x] Use Redis for distributed rate limiting
  - [x] Implement sliding window algorithm
  - [x] Support multiple limits:
    - [x] Global: 1000 req/min per server
    - [x] Per User: 100 req/min
    - [x] Per IP: 200 req/min
    - [x] Per Bot (widget): 500 req/min
  - [x] Return `429 Too Many Requests` with `Retry-After` header
  - [x] Whitelist enterprise users (bypass limits)

#### 4.2 Apply Rate Limiting
- [x] Add to chat endpoint (highest load)
- [x] Add to public widget endpoints
- [x] Add to file upload endpoints
- [x] Add to source creation endpoints
- [x] Skip rate limiting for health checks

---

### Step 5: Session Management (Redis-backed)

#### 5.1 Migrate Sessions to Redis
- [x] Update `internal/services/session/session_service.go`:
  - [x] Store sessions in Redis instead of PostgreSQL
  - [x] Set TTL for auto-expiration (24 hours)
  - [x] Use session ID as key
  - [x] Store full session data as JSON
- [x] Keep chat history in PostgreSQL (permanent storage)
- [x] Use Redis for active session tracking only

---

### Step 6: Connection Pooling & Monitoring

#### 6.1 Optimize Database Connections
- [x] Update `internal/db/db.go`:
  - [x] Set max open connections: 25
  - [x] Set max idle connections: 10
  - [x] Set connection max lifetime: 1 hour
  - [x] Add connection pool metrics

#### 6.2 Redis Connection Pooling
- [x] Configure Redis client pool:
  - [x] Max connections: 100
  - [x] Min idle connections: 10
  - [x] Connection timeout: 5s
  - [x] Read/Write timeout: 3s

#### 6.3 Add Monitoring Endpoint
- [x] Create `internal/handlers/health/health_handler.go`:
  - [x] `/health/db` - PostgreSQL status & pool stats
  - [x] `/health/redis` - Redis status & pool stats
  - [x] `/health/cache` - Cache hit/miss rates
  - [x] Return metrics in JSON format

---

### Step 7: Cache Invalidation Strategy

#### 7.1 Implement Cache Invalidation
- [x] On bot update:
  - [x] Delete `bot:{bot_id}`
  - [x] Delete `bots:user:{user_id}`
  - [x] Delete all `vector:{bot_id}:*` (pattern delete)
- [x] On source update:
  - [x] Delete `source:{source_id}`
  - [x] Delete `vector:{bot_id}:*`
- [x] On user update:
  - [x] Delete `user:{user_id}`
- [x] On session expire:
  - [x] Redis auto-expires (TTL-based)

#### 7.2 Cache Warming
- [x] Warm cache on server startup:
  - [x] Load top 100 most active bots
  - [x] Cache frequently accessed user data
- [x] Background job to refresh popular queries

---

## DevOps Tasks

### Step 1: Update Docker Compose
- [x] Add Redis service to `docker-compose.yml`:
  ```yaml
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3
  ```
- [x] Add `redis-data` volume
- [x] Update app service to depend on Redis

### Step 2: Update Makefile
- [x] Add `redis-cli` helper command
- [x] Add cache flush command for development
- [x] Update `make dev` to start Redis

---

## Testing Tasks

### Unit Tests
- [x] Test cache service operations
- [x] Test write buffer queuing
- [x] Test rate limiter logic
- [x] Test cache invalidation patterns

### Integration Tests
- [x] Test cache-aside pattern with real Redis
- [x] Test write buffer under load
- [x] Test rate limiting with concurrent requests
- [x] Test cache hit/miss ratios

### Load Tests
- [x] Simulate 10k concurrent writes
- [x] Measure database load with/without Redis
- [x] Verify cache reduces DB queries by >70%
- [x] Ensure no data loss in write buffer

---

## Success Metrics

- [x] Handle 10k+ concurrent write requests without errors
- [x] Cache hit rate >70% for read operations
- [x] API response time <200ms (p95)
- [x] Database write load reduced by >80%
- [x] Rate limiting prevents abuse (no 500 errors from overload)
- [x] Zero data loss in buffered writes

---

## Rollback Plan

If Redis causes issues:
1. Feature flag to disable Redis caching
2. Fall back to direct PostgreSQL reads/writes
3. Keep rate limiting in-memory (per-server basis)
4. Monitor for performance degradation

---

**Created**: February 11, 2025
**Completed**: February 12, 2025
