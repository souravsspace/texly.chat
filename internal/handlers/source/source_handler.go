package source

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/queue"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	sourceRepo "github.com/souravsspace/texly.chat/internal/repo/source"
)

/*
* SourceHandler handles HTTP requests for sources
 */
type SourceHandler struct {
	sourceRepo *sourceRepo.SourceRepo
	botRepo    *botRepo.BotRepo
	jobQueue   queue.JobQueue
}

/*
* NewSourceHandler creates a new source handler
 */
func NewSourceHandler(sourceRepo *sourceRepo.SourceRepo, botRepo *botRepo.BotRepo, jobQueue queue.JobQueue) *SourceHandler {
	return &SourceHandler{
		sourceRepo: sourceRepo,
		botRepo:    botRepo,
		jobQueue:   jobQueue,
	}
}

/*
* CreateSource handles POST /api/bots/:id/sources
 */
func (h *SourceHandler) CreateSource(c *gin.Context) {
	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get bot ID from URL
	botID := c.Param("id")

	// Verify bot ownership
	bot, err := h.botRepo.GetByID(botID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	if bot == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Parse request
	var req models.CreateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create source
	source := &models.Source{
		BotID:  botID,
		URL:    req.URL,
		Status: models.SourceStatusPending,
	}

	if err := h.sourceRepo.Create(source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source"})
		return
	}

	// Enqueue job for processing
	job := queue.Job{
		SourceID: source.ID,
		BotID:    botID,
		URL:      req.URL,
	}

	if err := h.jobQueue.Enqueue(job); err != nil {
		// Source created but job failed to enqueue - update status to failed
		_ = h.sourceRepo.UpdateStatus(source.ID, models.SourceStatusFailed, "Failed to queue processing job")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue processing job"})
		return
	}

	c.JSON(http.StatusCreated, source)
}

/*
* ListSources handles GET /api/bots/:id/sources
 */
func (h *SourceHandler) ListSources(c *gin.Context) {
	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get bot ID from URL
	botID := c.Param("id")

	// Verify bot ownership
	bot, err := h.botRepo.GetByID(botID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	if bot == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get sources
	sources, err := h.sourceRepo.ListByBotID(botID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sources"})
		return
	}

	c.JSON(http.StatusOK, sources)
}

/*
* GetSource handles GET /api/bots/:id/sources/:sourceId
 */
func (h *SourceHandler) GetSource(c *gin.Context) {
	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get IDs from URL
	botID := c.Param("id")
	sourceID := c.Param("sourceId")

	// Verify bot ownership
	bot, err := h.botRepo.GetByID(botID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	if bot == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get source (also verifies it belongs to the bot)
	source, err := h.sourceRepo.GetByBotIDAndSourceID(botID, sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	c.JSON(http.StatusOK, source)
}

/*
* DeleteSource handles DELETE /api/bots/:id/sources/:sourceId
 */
func (h *SourceHandler) DeleteSource(c *gin.Context) {
	// Get authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get IDs from URL
	botID := c.Param("id")
	sourceID := c.Param("sourceId")

	// Verify bot ownership
	bot, err := h.botRepo.GetByID(botID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	if bot == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Verify source belongs to bot
	_, err = h.sourceRepo.GetByBotIDAndSourceID(botID, sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	// Delete source
	if err := h.sourceRepo.Delete(sourceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Source deleted successfully"})
}
