package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/models"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	usage "github.com/souravsspace/texly.chat/internal/services/billing/usage"
	"github.com/souravsspace/texly.chat/internal/services/chat"
)

/*
 * ChatHandler handles HTTP requests for chat endpoints
 */
type ChatHandler struct {
	botRepo     *botRepo.BotRepo
	chatService *chat.ChatService
	usageSvc    *usage.UsageService
}

/*
 * NewChatHandler creates a new chat handler instance
 */
func NewChatHandler(botRepo *botRepo.BotRepo, chatService *chat.ChatService, usageSvc *usage.UsageService) *ChatHandler {
	return &ChatHandler{
		botRepo:     botRepo,
		chatService: chatService,
		usageSvc:    usageSvc,
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

	// Track usage (billing owner is Bot.UserID)
	if h.usageSvc != nil {
		if err := h.usageSvc.TrackChatMessage(bot.UserID, bot.ID); err != nil {
			// Failed to track usage/charge - should we block?
			// For now, log error and proceed, or block if strict
			// Ideally we block if payment required
			fmt.Printf("Error tracking chat usage: %v\n", err)
			// If we want to block on insufficient funds (already handled in UsageService which returns error? 
			// UsageService returns error if DB error. It updates balance but doesn't strictly block UNLESS logic added.
			// Current UsageService implementation returns error only on DB failure.
			// Pro tier logic: "Deduct from credits if available".
			// If 0 credits, we assume pay-as-you-go and just record usage.
			// So we don't block here unless DB fails.
		}
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
