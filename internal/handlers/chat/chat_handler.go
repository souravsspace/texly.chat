package chat

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/models"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	"github.com/souravsspace/texly.chat/internal/services/chat"
)

/*
 * ChatHandler handles HTTP requests for chat endpoints
 */
type ChatHandler struct {
	botRepo     *botRepo.BotRepo
	chatService *chat.ChatService
}

/*
 * NewChatHandler creates a new chat handler instance
 */
func NewChatHandler(botRepo *botRepo.BotRepo, chatService *chat.ChatService) *ChatHandler {
	return &ChatHandler{
		botRepo:     botRepo,
		chatService: chatService,
	}
}

/*
 * StreamChat handles POST /api/bots/:id/chat
 * Streams bot responses using Server-Sent Events (SSE)
 */
func (h *ChatHandler) StreamChat(c *gin.Context) {
	// Get user ID from auth middleware
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Get bot ID from URL params
	botID := c.Param("id")
	if botID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bot ID is required"})
		return
	}

	// Parse request body
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	// Validate bot ownership
	bot, err := h.botRepo.GetByID(botID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch bot"})
		return
	}
	if bot == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bot not found"})
		return
	}

	// Validate chat service is available
	if h.chatService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Chat service not available"})
		return
	}

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
	// Generate a session ID for this chat (in dashboard, session = one conversation)
	sessionID := botID + "-" + userID
	tokenChan, errChan := h.chatService.StreamChat(
		c.Request.Context(),
		botID,
		bot.SystemPrompt,
		req.Message,
		sessionID,
		&userID,
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
