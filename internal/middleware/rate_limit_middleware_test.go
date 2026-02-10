package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis_rate/v10"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestRateLimiter(t *testing.T) (*RateLimiter, *gorm.DB, func()) {
	// Setup miniredis for testing
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	// Setup test database
	db := shared.SetupTestDB()

	// Create config with miniredis address
	cfg := configs.Config{
		RedisAddr:     mr.Addr(),
		RedisPassword: "",
		RedisDB:       0,
	}

	// Create rate limiter
	rateLimiter, err := NewRateLimiter(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create rate limiter: %v", err)
	}

	cleanup := func() {
		mr.Close()
	}

	return rateLimiter, db, cleanup
}

func TestNewRateLimiter(t *testing.T) {
	rateLimiter, _, cleanup := setupTestRateLimiter(t)
	defer cleanup()

	assert.NotNil(t, rateLimiter)
	assert.NotNil(t, rateLimiter.limiter)
	assert.NotNil(t, rateLimiter.db)
}

func TestNewRateLimiter_InvalidRedis(t *testing.T) {
	db := shared.SetupTestDB()

	cfg := configs.Config{
		RedisAddr:     "invalid:9999",
		RedisPassword: "",
		RedisDB:       0,
	}

	rateLimiter, err := NewRateLimiter(cfg, db)
	assert.Error(t, err)
	assert.Nil(t, rateLimiter)
	assert.Contains(t, err.Error(), "failed to connect to Redis")
}

func TestDefaultRateLimitConfig(t *testing.T) {
	config := DefaultRateLimitConfig()

	assert.Equal(t, 100, config.FreeLimit)
	assert.Equal(t, 1000, config.ProLimit)
	assert.Equal(t, time.Hour, config.Period)
}

func TestRateLimiter_Allow(t *testing.T) {
	rateLimiter, _, cleanup := setupTestRateLimiter(t)
	defer cleanup()

	ctx := context.Background()
	key := "test:user:123"

	// Test allowing requests within limit
	config := RateLimitConfig{
		FreeLimit: 5,
		ProLimit:  10,
		Period:    time.Minute,
	}

	// Should allow first 5 requests
	for i := 0; i < 5; i++ {
		result, err := rateLimiter.limiter.Allow(ctx, key, redis_rate.Limit{
			Rate:   config.FreeLimit,
			Burst:  config.FreeLimit,
			Period: config.Period,
		})

		assert.NoError(t, err)
		assert.Greater(t, result.Allowed, 0)
	}

	// 6th request should be rate limited
	result, err := rateLimiter.limiter.Allow(ctx, key, redis_rate.Limit{
		Rate:   config.FreeLimit,
		Burst:  config.FreeLimit,
		Period: config.Period,
	})

	assert.NoError(t, err)
	assert.Equal(t, 0, result.Allowed)
	assert.Greater(t, result.RetryAfter, time.Duration(0))
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	rateLimiter, _, cleanup := setupTestRateLimiter(t)
	defer cleanup()

	ctx := context.Background()
	config := RateLimitConfig{
		FreeLimit: 5,
		ProLimit:  10,
		Period:    time.Minute,
	}

	// Use up limit for first key
	key1 := "test:user:123"
	for i := 0; i < 5; i++ {
		_, err := rateLimiter.limiter.Allow(ctx, key1, redis_rate.Limit{
			Rate:   config.FreeLimit,
			Burst:  config.FreeLimit,
			Period: config.Period,
		})
		assert.NoError(t, err)
	}

	// Second key should still have its own limit
	key2 := "test:user:456"
	result, err := rateLimiter.limiter.Allow(ctx, key2, redis_rate.Limit{
		Rate:   config.FreeLimit,
		Burst:  config.FreeLimit,
		Period: config.Period,
	})

	assert.NoError(t, err)
	assert.Greater(t, result.Allowed, 0)
	assert.Equal(t, config.FreeLimit-1, result.Remaining)
}
