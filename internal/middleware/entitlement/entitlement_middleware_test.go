package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	middleware "github.com/souravsspace/texly.chat/internal/middleware/entitlement"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestEntitlementMiddleware(t *testing.T) {
	db := shared.SetupTestDB()

	// Setup user helper
	setupUser := func(id, tier string) {
		shared.TruncateTable(db, "users")
		shared.TruncateTable(db, "bots")
		shared.TruncateTable(db, "usage_records")

		db.Exec("INSERT INTO users (id, email, tier) VALUES (?, ?, ?)", id, "mid@example.com", tier)
	}

	// Create test engine
	r := gin.New()

	// Mock Auth Middleware to inject user
	mockAuth := func(userID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			var user models.User
			if err := db.First(&user, "id = ?", userID).Error; err != nil {
				c.AbortWithStatus(401)
				return
			}
			c.Set("user", &user)
			c.Set("userID", user.ID)
			c.Next()
		}
	}

	entitlementMw := middleware.NewEntitlementMiddleware(db)

	t.Run("BotCreationLimit_Free_Exceeded", func(t *testing.T) {
		userID := "user_free_bot"
		setupUser(userID, "free")

		// Create 1 bot (Limit is 1 for Free)
		db.Exec("INSERT INTO bots (id, user_id, name) VALUES (?, ?, ?)", "bot1", userID, "Bot 1")

		// Add route
		r.POST("/bots_free", mockAuth(userID), entitlementMw.EnforceLimit("bot_creation"), func(c *gin.Context) {
			c.Status(200)
		})

		req, _ := http.NewRequest("POST", "/bots_free", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "Bot limit reached")
	})

	t.Run("BotCreationLimit_Pro_Allowed", func(t *testing.T) {
		userID := "user_pro_bot"
		setupUser(userID, "pro")

		// Create 3 bots (Pro limit 5)
		for i := 0; i < 3; i++ {
			db.Exec("INSERT INTO bots (id, user_id, name) VALUES (?, ?, ?)", string(rune(i)), userID, "Bot")
		}

		r.POST("/bots_pro", mockAuth(userID), entitlementMw.EnforceLimit("bot_creation"), func(c *gin.Context) {
			c.Status(200)
		})

		req, _ := http.NewRequest("POST", "/bots_pro", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("MessageLimit_Free_Exceeded", func(t *testing.T) {
		userID := "user_free_msg"
		setupUser(userID, "free")

		// Insert 101 messages (Limit 100)
		// UsageRecord counts monthly usage.
		// Middleware calls `getMonthlyMessages` which queries usage records.
		// Note: implementation details matter. Does EnforceLimit count messages or check usage?
		// "monthlyMessages := getMonthlyMessages(user.ID)"

		// Insert UsageRecord instead of standard messages if that's what it counts
		// Logic check: EntitlementMiddleware usually iterates `usage_records` where event_type='chat_message'
		currentMonth := time.Now()
		for i := 0; i < 101; i++ {
			db.Create(&models.UsageRecord{
				UserID:    userID,
				Type:      "chat_message",
				CreatedAt: currentMonth,
			})
		}

		r.POST("/msg_free", mockAuth(userID), entitlementMw.EnforceLimit("message_send"), func(c *gin.Context) {
			c.Status(200)
		})

		req, _ := http.NewRequest("POST", "/msg_free", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
	})
}
