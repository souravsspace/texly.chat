package bot

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/models"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	"gorm.io/gorm"
)

type BotHandler struct {
	repo *botRepo.BotRepo
}

func NewBotHandler(repo *botRepo.BotRepo) *BotHandler {
	return &BotHandler{repo: repo}
}

// CreateBot - POST /api/bots
func (h *BotHandler) CreateBot(c *gin.Context) {
	userID := c.GetString("user_id") // Assumes Auth middleware sets this
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var req models.CreateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	bot := models.Bot{
		UserID:       userID,
		Name:         req.Name,
		SystemPrompt: req.SystemPrompt,
	}

	if err := h.repo.Create(&bot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create bot"})
		return
	}

	c.JSON(http.StatusCreated, bot)
}

// ListBots - GET /api/bots
func (h *BotHandler) ListBots(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	bots, err := h.repo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch bots"})
		return
	}

	c.JSON(http.StatusOK, bots)
}

// GetBot - GET /api/bots/:id
func (h *BotHandler) GetBot(c *gin.Context) {
	userID := c.GetString("user_id")
	botID := c.Param("id")

	bot, err := h.repo.GetByID(botID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch bot"})
		return
	}
	if bot == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bot not found"})
		return
	}

	c.JSON(http.StatusOK, bot)
}

// UpdateBot - PUT /api/bots/:id
func (h *BotHandler) UpdateBot(c *gin.Context) {
	userID := c.GetString("user_id")
	botID := c.Param("id")

	var req models.UpdateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	bot, err := h.repo.GetByID(botID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch bot"})
		return
	}
	if bot == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bot not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		bot.Name = req.Name
	}
	// SystemPrompt can be empty, but usually check if present.
	// For simplicity, we update if provided in JSON (even empty string).
	// But struct zero value check is tricky.
	// Let's assume we update if needed. Alternatively map fields.
	bot.SystemPrompt = req.SystemPrompt

	if err := h.repo.Update(bot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update bot"})
		return
	}

	c.JSON(http.StatusOK, bot)
}

// DeleteBot - DELETE /api/bots/:id
func (h *BotHandler) DeleteBot(c *gin.Context) {
	userID := c.GetString("user_id")
	botID := c.Param("id")

	err := h.repo.Delete(botID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Bot not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete bot"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bot deleted successfully"})
}
