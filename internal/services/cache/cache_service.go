package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrCacheMiss is returned when a key is not found in the cache.
var ErrCacheMiss = errors.New("cache: key not found")

/*
* CacheService provides reusable caching operations backed by Redis.
* It wraps the Redis client with JSON serialization and standardized error handling.
 */
type CacheService struct {
	client  *redis.Client
	enabled bool
}

/*
* NewCacheService creates a new CacheService.
* If client is nil, the service operates in pass-through mode (all reads miss, writes are no-ops).
 */
func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{
		client:  client,
		enabled: client != nil,
	}
}

/*
* Get retrieves a string value from cache.
 */
func (c *CacheService) Get(ctx context.Context, key string) (string, error) {
	if c == nil || !c.enabled {
		return "", ErrCacheMiss
	}
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrCacheMiss
	}
	return val, err
}

/*
* Set stores a string value in cache with a TTL.
 */
func (c *CacheService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if c == nil || !c.enabled {
		return nil
	}
	return c.client.Set(ctx, key, value, ttl).Err()
}

/*
* GetJSON retrieves a value from cache and unmarshals it into dest.
* Returns ErrCacheMiss if the key does not exist.
 */
func (c *CacheService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	if c == nil || !c.enabled {
		return ErrCacheMiss
	}
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	}
	if err != nil {
		return fmt.Errorf("cache get error: %w", err)
	}
	return json.Unmarshal([]byte(val), dest)
}

/*
* SetJSON marshals value to JSON and stores it in cache with a TTL.
 */
func (c *CacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if c == nil || !c.enabled {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

/*
* Delete removes a single key from cache.
 */
func (c *CacheService) Delete(ctx context.Context, key string) error {
	if c == nil || !c.enabled {
		return nil
	}
	return c.client.Del(ctx, key).Err()
}

/*
* DeletePattern removes all keys matching a glob-style pattern.
* Uses SCAN to avoid blocking the server on large key sets.
 */
func (c *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	if c == nil || !c.enabled {
		return nil
	}

	var cursor uint64
	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("cache scan error: %w", err)
		}
		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("cache delete error: %w", err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

/*
* Exists checks whether a key exists in cache.
 */
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	if c == nil || !c.enabled {
		return false, nil
	}
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

/*
* IsEnabled returns whether the cache service is operational.
 */
func (c *CacheService) IsEnabled() bool {
	return c != nil && c.enabled
}
