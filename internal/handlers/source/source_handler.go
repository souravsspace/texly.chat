package source

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/queue"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	sourceRepo "github.com/souravsspace/texly.chat/internal/repo/source"
	"github.com/souravsspace/texly.chat/internal/services/storage"
)

/*
* SourceHandler handles HTTP requests for sources
 */
type SourceHandler struct {
	sourceRepo   *sourceRepo.SourceRepo
	botRepo      *botRepo.BotRepo
	jobQueue     queue.JobQueue
	storageSvc   *storage.MinIOStorageService
	maxUploadMB  int
}

/*
* NewSourceHandler creates a new source handler
 */
func NewSourceHandler(sourceRepo *sourceRepo.SourceRepo, botRepo *botRepo.BotRepo, jobQueue queue.JobQueue, storageSvc *storage.MinIOStorageService, maxUploadMB int) *SourceHandler {
	return &SourceHandler{
		sourceRepo:  sourceRepo,
		botRepo:     botRepo,
		jobQueue:    jobQueue,
		storageSvc:  storageSvc,
		maxUploadMB: maxUploadMB,
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
		BotID:      botID,
		SourceType: models.SourceTypeURL,
		URL:        req.URL,
		Status:     models.SourceStatusPending,
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

/*
* UploadFileSource handles POST /api/bots/:id/sources/upload
 */
func (h *SourceHandler) UploadFileSource(c *gin.Context) {
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

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Validate file type
	if err := h.storageSvc.ValidateFileType(header.Filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate file size
	if err := h.storageSvc.ValidateFileSize(header.Size); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create source record first (to get ID for MinIO path)
	source := &models.Source{
		BotID:            botID,
		SourceType:       models.SourceTypeFile,
		OriginalFilename: header.Filename,
		ContentType:      storage.GetContentType(header.Filename),
		Status:           models.SourceStatusPending,
	}

	if err := h.sourceRepo.Create(source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source"})
		return
	}

	// Generate MinIO object name
	objectName := h.storageSvc.GenerateObjectName(source.ID, header.Filename)

	// Upload file to MinIO
	ctx := context.Background()
	if err := h.storageSvc.UploadFile(ctx, objectName, file, header.Size, source.ContentType); err != nil {
		// Failed to upload - delete source record and return error
		_ = h.sourceRepo.Delete(source.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// Update source with file path
	source.FilePath = objectName
	if err := h.sourceRepo.Create(source); err != nil {
		// Cleanup MinIO file if DB update fails
		_ = h.storageSvc.DeleteFile(ctx, objectName)
		_ = h.sourceRepo.Delete(source.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update source"})
		return
	}

	// Enqueue job for processing
	job := queue.Job{
		SourceID: source.ID,
		BotID:    botID,
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
* CreateTextSource handles POST /api/bots/:id/sources/text
 */
func (h *SourceHandler) CreateTextSource(c *gin.Context) {
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
	var req models.CreateTextSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate text size (use same limit as file uploads)
	textSize := int64(len(req.Text))
	maxBytes := int64(h.maxUploadMB) * 1024 * 1024
	if textSize > maxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Text size exceeds maximum allowed size (%d MB)", h.maxUploadMB)})
		return
	}

	if textSize == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text cannot be empty"})
		return
	}

	// Generate a filename for the text source
	filename := "text-source.txt"
	if req.Name != "" {
		filename = req.Name + ".txt"
	}

	// Create source record
	source := &models.Source{
		BotID:            botID,
		SourceType:       models.SourceTypeText,
		OriginalFilename: filename,
		ContentType:      "text/plain",
		Status:           models.SourceStatusPending,
	}

	if err := h.sourceRepo.Create(source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source"})
		return
	}

	// Upload text to MinIO as a file
	ctx := context.Background()
	objectName := h.storageSvc.GenerateObjectName(source.ID, filename)

	// Convert text to reader
	textReader := strings.NewReader(req.Text)

	if err := h.storageSvc.UploadFile(ctx, objectName, textReader, textSize, "text/plain"); err != nil {
		// Failed to upload - delete source record
		_ = h.sourceRepo.Delete(source.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save text: %v", err)})
		return
	}

	// Update source with file path
	source.FilePath = objectName
	if err := h.sourceRepo.Create(source); err != nil {
		// Cleanup MinIO file if DB update fails
		_ = h.storageSvc.DeleteFile(ctx, objectName)
		_ = h.sourceRepo.Delete(source.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update source"})
		return
	}

	// Enqueue job for processing
	job := queue.Job{
		SourceID: source.ID,
		BotID:    botID,
	}

	if err := h.jobQueue.Enqueue(job); err != nil {
		// Source created but job failed to enqueue - update status to failed
		_ = h.sourceRepo.UpdateStatus(source.ID, models.SourceStatusFailed, "Failed to queue processing job")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue processing job"})
		return
	}

	c.JSON(http.StatusCreated, source)
}
