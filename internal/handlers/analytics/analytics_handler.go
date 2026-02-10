package analytics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/services/analytics"
)

/*
 * AnalyticsHandler handles HTTP requests for analytics endpoints
 */
type AnalyticsHandler struct {
	analyticsService *analytics.AnalyticsService
}

/*
 * NewAnalyticsHandler creates a new analytics handler instance
 */
func NewAnalyticsHandler(analyticsService *analytics.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

/*
 * GetBotAnalytics handles GET /api/analytics/bots/:id
 * Returns overall analytics for a specific bot
 */
func (h *AnalyticsHandler) GetBotAnalytics(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	botID := c.Param("id")
	if botID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bot ID is required"})
		return
	}

	// TODO: Verify bot ownership - for now we trust the user_id from auth

	analytics, err := h.analyticsService.GetBotAnalytics(c.Request.Context(), botID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

/*
 * GetBotDailyStats handles GET /api/analytics/bots/:id/daily
 * Returns daily message statistics for a bot
 * Query params: days (default: 30)
 */
func (h *AnalyticsHandler) GetBotDailyStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	botID := c.Param("id")
	if botID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bot ID is required"})
		return
	}

	// Parse days parameter
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		days = 30
	}
	if days > 365 {
		days = 365 // Cap at 1 year
	}

	stats, err := h.analyticsService.GetBotDailyStats(c.Request.Context(), botID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch daily stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

/*
 * GetUserAnalytics handles GET /api/analytics/user
 * Returns analytics for all bots owned by the current user
 */
func (h *AnalyticsHandler) GetUserAnalytics(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	analytics, err := h.analyticsService.GetUserAnalytics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

/*
 * GetSessionMessages handles GET /api/analytics/sessions/:id/messages
 * Returns all messages for a specific session
 */
func (h *AnalyticsHandler) GetSessionMessages(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Session ID is required"})
		return
	}

	messages, err := h.analyticsService.GetSessionMessages(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
