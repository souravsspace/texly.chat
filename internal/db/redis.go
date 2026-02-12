package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

/*
* InitRedis initializes the Redis client singleton with connection pooling and health check.
* This should be called once at application startup. Subsequent calls are no-ops.
 */
func InitRedis(redisURL string, poolSize, minIdleConns int) error {
	var initErr error
	redisOnce.Do(func() {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			initErr = fmt.Errorf("failed to parse Redis URL: %w", err)
			return
		}

		// Configure connection pool
		opt.PoolSize = poolSize
		opt.MinIdleConns = minIdleConns
		opt.ConnMaxLifetime = time.Hour
		opt.DialTimeout = 5 * time.Second
		opt.ReadTimeout = 3 * time.Second
		opt.WriteTimeout = 3 * time.Second

		redisClient = redis.NewClient(opt)

		// Health check
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			initErr = fmt.Errorf("failed to ping Redis: %w", err)
			return
		}

		fmt.Printf("✅ Redis connected (pool: %d, idle: %d)\n", poolSize, minIdleConns)
	})
	return initErr
}

/*
* GetRedisClient returns the singleton Redis client.
* Must call InitRedis before using this.
 */
func GetRedisClient() *redis.Client {
	return redisClient
}

/*
* CloseRedis gracefully closes the Redis connection.
 */
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

/*
* PingRedis checks if the Redis connection is alive.
 */
func PingRedis(ctx context.Context) error {
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return redisClient.Ping(ctx).Err()
}
