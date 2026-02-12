package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

/*
* HealthHandler exposes health and monitoring endpoints
 */
type HealthHandler struct {
	db          *gorm.DB
	redisClient *redis.Client
}

/*
* NewHealthHandler creates a new HealthHandler
 */
func NewHealthHandler(db *gorm.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redisClient: redisClient}
}

/*
* GetBasicHealth returns a basic ok response
 */
func (h *HealthHandler) GetBasicHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

/*
* GetDBHealth returns PostgreSQL connection pool stats
 */
func (h *HealthHandler) GetDBHealth(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}

	stats := sqlDB.Stats()
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"pool": gin.H{
			"maxOpenConnections": stats.MaxOpenConnections,
			"openConnections":    stats.OpenConnections,
			"inUse":              stats.InUse,
			"idle":               stats.Idle,
			"waitCount":          stats.WaitCount,
			"waitDuration":       stats.WaitDuration.String(),
		},
	})
}

/*
* GetRedisHealth returns Redis connection pool stats
 */
func (h *HealthHandler) GetRedisHealth(c *gin.Context) {
	if h.redisClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "redis not configured"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}

	stats := h.redisClient.PoolStats()
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"pool": gin.H{
			"hits":       stats.Hits,
			"misses":     stats.Misses,
			"timeouts":   stats.Timeouts,
			"totalConns": stats.TotalConns,
			"idleConns":  stats.IdleConns,
			"staleConns": stats.StaleConns,
		},
	})
}
