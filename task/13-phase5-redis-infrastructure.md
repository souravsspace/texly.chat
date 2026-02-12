# Phase 5: Infrastructure & Performance (Redis Layer)

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
- [ ] Add `github.com/redis/go-redis/v9` to `go.mod`
- [ ] Add Redis configuration to `configs/config.go`:
  ```go
  RedisURL      string // e.g., "redis://localhost:6379"
  RedisPassword string
  RedisDB       int
  RedisTTL      int // Default cache TTL in seconds
  ```
- [ ] Add environment variables to `.env.local`:
  ```
  REDIS_URL=redis://localhost:6379
  REDIS_PASSWORD=
  REDIS_DB=0
  REDIS_TTL=3600
  ```

#### 1.2 Create Redis Client
- [ ] Create `internal/db/redis.go`:
  - [ ] Initialize Redis client with connection pooling
  - [ ] Implement health check function
  - [ ] Add graceful shutdown
  - [ ] Export singleton client instance
- [ ] Update `cmd/app/main.go` to initialize Redis on startup
- [ ] Add Redis health check to server startup

---

### Step 2: Caching Layer Implementation

#### 2.1 Create Cache Service
- [ ] Create `internal/services/cache/cache_service.go`:
  - [ ] `Get(ctx, key) (string, error)` - Get cached value
  - [ ] `Set(ctx, key, value, ttl) error` - Set cache with TTL
  - [ ] `Delete(ctx, key) error` - Invalidate cache
  - [ ] `DeletePattern(ctx, pattern) error` - Bulk invalidation
  - [ ] `Exists(ctx, key) (bool, error)` - Check if key exists
  - [ ] `GetJSON(ctx, key, dest) error` - Get and unmarshal JSON
  - [ ] `SetJSON(ctx, key, value, ttl) error` - Marshal and set JSON

#### 2.2 Cache Keys Strategy
- [ ] Define cache key patterns in `internal/services/cache/keys.go`:
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
- [ ] Update Bot Repository (`internal/repo/bot/bot_repo.go`):
  - [ ] Check cache before database read
  - [ ] Store result in cache after database read
  - [ ] Invalidate cache on update/delete
- [ ] Update User Repository:
  - [ ] Cache user profiles
  - [ ] Cache authentication lookups
- [ ] Update Source Repository:
  - [ ] Cache source metadata
  - [ ] Invalidate on status changes
- [ ] Update Vector Search:
  - [ ] Cache frequent query results (hash query as key)
  - [ ] Set shorter TTL (5-10 minutes)

---

### Step 3: Write Buffering for High Load

#### 3.1 Create Write Buffer Service
- [ ] Create `internal/services/writebuffer/buffer_service.go`:
  - [ ] Use Redis Lists as write queues
  - [ ] `QueueWrite(ctx, operation, data) error` - Add to buffer
  - [ ] `ProcessBuffer(ctx) error` - Batch process writes
  - [ ] Implement retry logic for failed writes
  - [ ] Add dead letter queue for permanently failed writes

#### 3.2 Batch Write Worker
- [ ] Create `internal/worker/write_worker.go`:
  - [ ] Poll Redis queue every 100ms
  - [ ] Batch up to 100 writes per transaction
  - [ ] Use PostgreSQL transaction for atomic batch commit
  - [ ] Log failures and retry up to 3 times
  - [ ] Move to dead letter queue after max retries

#### 3.3 Update Chat Handler
- [ ] Modify `internal/handlers/chat/chat_handler.go`:
  - [ ] Queue chat message writes instead of direct DB write
  - [ ] Return immediately to user (async write)
  - [ ] Ensure message ID returned for tracking

---

### Step 4: Rate Limiting

#### 4.1 Create Rate Limiter Middleware
- [ ] Create `internal/middleware/rate_limit.go`:
  - [ ] Use Redis for distributed rate limiting
  - [ ] Implement sliding window algorithm
  - [ ] Support multiple limits:
    - [ ] Global: 1000 req/min per server
    - [ ] Per User: 100 req/min
    - [ ] Per IP: 200 req/min
    - [ ] Per Bot (widget): 500 req/min
  - [ ] Return `429 Too Many Requests` with `Retry-After` header
  - [ ] Whitelist enterprise users (bypass limits)

#### 4.2 Apply Rate Limiting
- [ ] Add to chat endpoint (highest load)
- [ ] Add to public widget endpoints
- [ ] Add to file upload endpoints
- [ ] Add to source creation endpoints
- [ ] Skip rate limiting for health checks

---

### Step 5: Session Management (Redis-backed)

#### 5.1 Migrate Sessions to Redis
- [ ] Update `internal/services/session/session_service.go`:
  - [ ] Store sessions in Redis instead of PostgreSQL
  - [ ] Set TTL for auto-expiration (24 hours)
  - [ ] Use session ID as key
  - [ ] Store full session data as JSON
- [ ] Keep chat history in PostgreSQL (permanent storage)
- [ ] Use Redis for active session tracking only

---

### Step 6: Connection Pooling & Monitoring

#### 6.1 Optimize Database Connections
- [ ] Update `internal/db/db.go`:
  - [ ] Set max open connections: 25
  - [ ] Set max idle connections: 10
  - [ ] Set connection max lifetime: 1 hour
  - [ ] Add connection pool metrics

#### 6.2 Redis Connection Pooling
- [ ] Configure Redis client pool:
  - [ ] Max connections: 100
  - [ ] Min idle connections: 10
  - [ ] Connection timeout: 5s
  - [ ] Read/Write timeout: 3s

#### 6.3 Add Monitoring Endpoint
- [ ] Create `internal/handlers/health/health_handler.go`:
  - [ ] `/health/db` - PostgreSQL status & pool stats
  - [ ] `/health/redis` - Redis status & pool stats
  - [ ] `/health/cache` - Cache hit/miss rates
  - [ ] Return metrics in JSON format

---

### Step 7: Cache Invalidation Strategy

#### 7.1 Implement Cache Invalidation
- [ ] On bot update:
  - [ ] Delete `bot:{bot_id}`
  - [ ] Delete `bots:user:{user_id}`
  - [ ] Delete all `vector:{bot_id}:*` (pattern delete)
- [ ] On source update:
  - [ ] Delete `source:{source_id}`
  - [ ] Delete `vector:{bot_id}:*`
- [ ] On user update:
  - [ ] Delete `user:{user_id}`
- [ ] On session expire:
  - [ ] Redis auto-expires (TTL-based)

#### 7.2 Cache Warming
- [ ] Warm cache on server startup:
  - [ ] Load top 100 most active bots
  - [ ] Cache frequently accessed user data
- [ ] Background job to refresh popular queries

---

## DevOps Tasks

### Step 1: Update Docker Compose
- [ ] Add Redis service to `docker-compose.yml`:
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
- [ ] Add `redis-data` volume
- [ ] Update app service to depend on Redis

### Step 2: Update Makefile
- [ ] Add `redis-cli` helper command
- [ ] Add cache flush command for development
- [ ] Update `make dev` to start Redis

---

## Testing Tasks

### Unit Tests
- [ ] Test cache service operations
- [ ] Test write buffer queuing
- [ ] Test rate limiter logic
- [ ] Test cache invalidation patterns

### Integration Tests
- [ ] Test cache-aside pattern with real Redis
- [ ] Test write buffer under load
- [ ] Test rate limiting with concurrent requests
- [ ] Test cache hit/miss ratios

### Load Tests
- [ ] Simulate 10k concurrent writes
- [ ] Measure database load with/without Redis
- [ ] Verify cache reduces DB queries by >70%
- [ ] Ensure no data loss in write buffer

---

## Success Metrics

- [ ] Handle 10k+ concurrent write requests without errors
- [ ] Cache hit rate >70% for read operations
- [ ] API response time <200ms (p95)
- [ ] Database write load reduced by >80%
- [ ] Rate limiting prevents abuse (no 500 errors from overload)
- [ ] Zero data loss in buffered writes

---

## Rollback Plan

If Redis causes issues:
1. Feature flag to disable Redis caching
2. Fall back to direct PostgreSQL reads/writes
3. Keep rate limiting in-memory (per-server basis)
4. Monitor for performance degradation

---

**Created**: February 11, 2025
