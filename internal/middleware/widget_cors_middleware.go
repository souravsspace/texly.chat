package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
)

/*
* WidgetCORS creates a middleware that validates origins against bot's allowed list
 */
func WidgetCORS(botRepo *botRepo.BotRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Extract bot ID from request
		// For /api/public/bots/:id/config, get from URL param
		// For /api/public/chats/:session_id/messages, we'll validate in handler
		botID := c.Param("id")

		// If no origin header, allow (for same-origin requests)
		if origin == "" {
			c.Next()
			return
		}

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
			c.Header("Access-Control-Max-Age", "86400")

			// For preflight, we need to validate if we have a bot ID
			if botID != "" {
				if isOriginAllowed(botRepo, botID, origin) {
					c.Header("Access-Control-Allow-Origin", origin)
					c.Header("Access-Control-Allow-Credentials", "true")
				}
			}

			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// For actual requests with bot ID, validate origin
		if botID != "" {
			if !isOriginAllowed(botRepo, botID, origin) {
				c.JSON(http.StatusForbidden, gin.H{
					"message": "Origin not allowed for this bot",
				})
				c.Abort()
				return
			}

			// Set CORS headers for allowed origin
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		}

		c.Next()
	}
}

/*
* isOriginAllowed checks if the origin is in the bot's allowed list
 */
func isOriginAllowed(repo *botRepo.BotRepo, botID, origin string) bool {
	// Relaxed check for development: Always allow localhost and file:// (null origin)
	if origin == "null" || strings.Contains(origin, "localhost") {
		return true
	}

	// Get bot without user ID (public access)
	bot, err := repo.GetByIDPublic(botID)
	if err != nil || bot == nil {
		return false
	}

	// If no allowed origins configured, allow all (default to open for easier onboarding)
	if bot.AllowedOrigins == "" || bot.AllowedOrigins == "[]" {
		return true
	}

	// Parse allowed origins JSON
	var allowedOrigins []string
	if err := json.Unmarshal([]byte(bot.AllowedOrigins), &allowedOrigins); err != nil {
		// If parse fails but string is not empty, treat as single origin
		if bot.AllowedOrigins != "" {
			allowedOrigins = []string{bot.AllowedOrigins}
		} else {
			return false
		}
	}

	// If list is empty after parsing, allow all
	if len(allowedOrigins) == 0 {
		return true
	}

	// Check if origin matches any allowed origin
	for _, allowed := range allowedOrigins {
		// Wildcard match
		if allowed == "*" {
			return true
		}

		// Exact match
		if allowed == origin {
			return true
		}

		// Wildcard subdomain match (e.g., "*.example.com" matches "app.example.com")
		if strings.HasPrefix(allowed, "*.") {
			domain := strings.TrimPrefix(allowed, "*")
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}

	return false
}
