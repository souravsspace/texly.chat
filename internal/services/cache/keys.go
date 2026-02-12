package cache

import "time"

// Cache key patterns.
// Use fmt.Sprintf to insert the dynamic parts.
const (
	// BotCacheKey caches a single bot by ID.
	// Format: bot:{bot_id}
	BotCacheKey = "bot:%s"

	// BotListCacheKey caches the bot list for a user.
	// Format: bots:user:{user_id}
	BotListCacheKey = "bots:user:%s"

	// UserCacheKey caches a user by ID.
	// Format: user:{user_id}
	UserCacheKey = "user:%s"

	// UserEmailCacheKey caches a user lookup by email.
	// Format: user:email:{email}
	UserEmailCacheKey = "user:email:%s"

	// SourceCacheKey caches a source by ID.
	// Format: source:{source_id}
	SourceCacheKey = "source:%s"

	// SourceListCacheKey caches the source list for a bot.
	// Format: sources:bot:{bot_id}
	SourceListCacheKey = "sources:bot:%s"

	// VectorSearchCacheKey caches vector search results.
	// Format: vector:{bot_id}:{query_hash}
	VectorSearchCacheKey = "vector:%s:%s"

	// SessionCacheKey caches an active chat session.
	// Format: session:{session_id}
	SessionCacheKey = "session:%s"
)

// TTL strategies per cache type.
const (
	BotCacheTTL          = 1 * time.Hour
	BotListCacheTTL      = 15 * time.Minute
	UserCacheTTL         = 1 * time.Hour
	SourceCacheTTL       = 30 * time.Minute
	SourceListCacheTTL   = 15 * time.Minute
	VectorSearchCacheTTL = 5 * time.Minute
	SessionCacheTTL      = 24 * time.Hour
)
