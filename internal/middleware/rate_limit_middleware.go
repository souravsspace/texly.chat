package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
	"gorm.io/gorm"
)

type RateLimiter struct {
	limiter *redis_rate.Limiter
	db      *gorm.DB
}

// RateLimitConfig defines rate limiting configuration for different user tiers
type RateLimitConfig struct {
	FreeLimit int           // Requests per period for free users
	ProLimit  int           // Requests per period for pro users
	Period    time.Duration // Time period for rate limiting
}

// DefaultRateLimitConfig returns sensible defaults for rate limiting
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		FreeLimit: 100,       // 100 requests per hour for free users
		ProLimit:  1000,      // 1000 requests per hour for pro users
		Period:    time.Hour, // 1 hour window
	}
}

// NewRateLimiter creates a new rate limiter with Redis backend
func NewRateLimiter(cfg configs.Config, db *gorm.DB) (*RateLimiter, error) {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Create rate limiter
	limiter := redis_rate.NewLimiter(rdb)

	return &RateLimiter{
		limiter: limiter,
		db:      db,
	}, nil
}

// RateLimitMiddleware applies rate limiting based on user tier
func (rl *RateLimiter) RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			// For public endpoints without auth, use IP-based rate limiting
			rl.rateLimitByIP(c, config)
			return
		}

		// Get user from database to check tier
		var user models.User
		if err := rl.db.First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
			c.Abort()
			return
		}

		// Determine rate limit based on user tier
		// TODO: Add Tier field to User model and check user.Tier
		// For now, use free tier limits for all users
		limit := config.FreeLimit

		// Create rate limit key
		key := fmt.Sprintf("rate_limit:user:%s", user.ID)

		// Check rate limit
		ctx := c.Request.Context()
		result, err := rl.limiter.Allow(ctx, key, redis_rate.Limit{
			Rate:   limit,
			Burst:  limit,
			Period: config.Period,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(result.RetryAfter).Unix()))

		// Check if rate limit exceeded
		if result.Allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": result.RetryAfter.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// rateLimitByIP applies IP-based rate limiting for unauthenticated requests
func (rl *RateLimiter) rateLimitByIP(c *gin.Context, config RateLimitConfig) {
	// Get client IP
	clientIP := c.ClientIP()
	key := fmt.Sprintf("rate_limit:ip:%s", clientIP)

	// Use free tier limits for IP-based rate limiting
	limit := config.FreeLimit

	// Check rate limit
	ctx := c.Request.Context()
	result, err := rl.limiter.Allow(ctx, key, redis_rate.Limit{
		Rate:   limit,
		Burst:  limit,
		Period: config.Period,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
		c.Abort()
		return
	}

	// Set rate limit headers
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
	c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(result.RetryAfter).Unix()))

	// Check if rate limit exceeded
	if result.Allowed == 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       "Rate limit exceeded",
			"retry_after": result.RetryAfter.Seconds(),
		})
		c.Abort()
		return
	}

	c.Next()
}

// PublicRateLimitMiddleware provides rate limiting for public widget endpoints
func (rl *RateLimiter) PublicRateLimitMiddleware() gin.HandlerFunc {
	config := RateLimitConfig{
		FreeLimit: 60, // 60 requests per hour for public endpoints
		ProLimit:  60, // Same limit for consistency
		Period:    time.Hour,
	}

	return func(c *gin.Context) {
		// Use IP-based rate limiting for public endpoints
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:public:%s", clientIP)

		// Check rate limit
		ctx := c.Request.Context()
		result, err := rl.limiter.Allow(ctx, key, redis_rate.Limit{
			Rate:   config.FreeLimit,
			Burst:  config.FreeLimit,
			Period: config.Period,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.FreeLimit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(result.RetryAfter).Unix()))

		// Check if rate limit exceeded
		if result.Allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": result.RetryAfter.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
