package public

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/models"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	"github.com/souravsspace/texly.chat/internal/services/chat"
	"github.com/souravsspace/texly.chat/internal/services/session"
)

/*
 * PublicHandler handles HTTP requests for public widget endpoints
 */
type PublicHandler struct {
	botRepo        *botRepo.BotRepo
	sessionService *session.SessionService
	chatService    *chat.ChatService
}

/*
 * NewPublicHandler creates a new public handler instance
 */
func NewPublicHandler(
	botRepo *botRepo.BotRepo,
	sessionService *session.SessionService,
	chatService *chat.ChatService,
) *PublicHandler {
	return &PublicHandler{
		botRepo:        botRepo,
		sessionService: sessionService,
		chatService:    chatService,
	}
}

/*
 * GetWidgetConfig handles GET /api/public/bots/:id/config
 * Returns widget configuration for embedding
 */
func (h *PublicHandler) GetWidgetConfig(c *gin.Context) {
	botID := c.Param("id")
	if botID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bot ID is required"})
		return
	}

	// Get bot without user authentication
	bot, err := h.botRepo.GetByIDPublic(botID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch bot"})
		return
	}
	if bot == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bot not found"})
		return
	}

	// Parse widget config
	var widgetConfig models.WidgetConfig
	if bot.WidgetConfig != "" {
		if err := json.Unmarshal([]byte(bot.WidgetConfig), &widgetConfig); err != nil {
			// Use default config if parsing fails
			widgetConfig = models.WidgetConfig{
				ThemeColor:     "#6366f1",
				InitialMessage: "Hi! How can I help you today?",
				Position:       "bottom-right",
				BotAvatar:      "",
			}
		}
	} else {
		// Default config
		widgetConfig = models.WidgetConfig{
			ThemeColor:     "#6366f1",
			InitialMessage: "Hi! How can I help you today?",
			Position:       "bottom-right",
			BotAvatar:      "",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            bot.ID,
		"name":          bot.Name,
		"widget_config": widgetConfig,
	})
}

/*
 * CreateSession handles POST /api/public/chats
 * Creates a new anonymous chat session
 */
func (h *PublicHandler) CreateSession(c *gin.Context) {
	var req models.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	// Verify bot exists
	bot, err := h.botRepo.GetByIDPublic(req.BotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch bot"})
		return
	}
	if bot == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bot not found"})
		return
	}

	// Create session
	session := h.sessionService.CreateSession(req.BotID)

	response := models.SessionResponse{
		SessionID: session.ID,
		BotID:     session.BotID,
		ExpiresAt: session.ExpiresAt,
	}

	c.JSON(http.StatusCreated, response)
}

/*
 * StreamChatPublic handles POST /api/public/chats/:session_id/messages
 * Streams bot responses for anonymous sessions using SSE
 */
func (h *PublicHandler) StreamChatPublic(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Session ID is required"})
		return
	}

	// Get and validate session
	chatSession, err := h.sessionService.GetSession(sessionID)
	if err != nil {
		if err == session.ErrSessionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Session not found"})
			return
		}
		if err == session.ErrSessionExpired {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Session has expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to validate session"})
		return
	}

	// Parse request body
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	// Get bot details
	bot, err := h.botRepo.GetByIDPublic(chatSession.BotID)
	if err != nil || bot == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch bot"})
		return
	}

	// Validate chat service is available
	if h.chatService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Chat service not available"})
		return
	}

	// Update session activity
	h.sessionService.UpdateActivity(sessionID)

	// Stream tokens using manual SSE writing
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Streaming not supported"})
		return
	}

	// Start streaming from chat service
	tokenChan, errChan := h.chatService.StreamChat(
		c.Request.Context(),
		bot.ID,
		bot.SystemPrompt,
		req.Message,
	)

	// Stream tokens
	for {
		select {
		case token, ok := <-tokenChan:
			if !ok {
				// Channel closed, stream done
				response := models.ChatTokenResponse{
					Type: "done",
				}
				data, _ := json.Marshal(response)
				c.Writer.Write([]byte("data: " + string(data) + "\n\n"))
				flusher.Flush()
				return
			}

			// Send token
			response := models.ChatTokenResponse{
				Type:    "token",
				Content: token,
			}
			data, _ := json.Marshal(response)
			c.Writer.Write([]byte("data: " + string(data) + "\n\n"))
			flusher.Flush()

		case err := <-errChan:
			if err != nil {
				// Send error
				response := models.ChatTokenResponse{
					Type:  "error",
					Error: err.Error(),
				}
				data, _ := json.Marshal(response)
				c.Writer.Write([]byte("data: " + string(data) + "\n\n"))
				flusher.Flush()
			}
			return
		}
	}
}
