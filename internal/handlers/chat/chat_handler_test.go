package chat_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/souravsspace/texly.chat/internal/handlers/chat"
	"github.com/souravsspace/texly.chat/internal/models"
	botRepo "github.com/souravsspace/texly.chat/internal/repo/bot"
	messageRepo "github.com/souravsspace/texly.chat/internal/repo/message"
	vectorRepo "github.com/souravsspace/texly.chat/internal/repo/vector"
	usage "github.com/souravsspace/texly.chat/internal/services/billing/usage"
	"github.com/souravsspace/texly.chat/internal/services/cache"
	chatSvc "github.com/souravsspace/texly.chat/internal/services/chat"
	"github.com/souravsspace/texly.chat/internal/services/embedding"
	"github.com/souravsspace/texly.chat/internal/services/vector"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestChatHandler_StreamChat(t *testing.T) {
	db := shared.SetupTestDB()

	// Mock Redis
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cacheSvc := cache.NewCacheService(rdb)

	// Mock OpenAI
	openAIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/chat/completions") {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`data: {"id":"chatcmpl-123","choices":[{"index":0,"delta":{"content":"Hello"}}],"created":123}` + "\n\n"))
			w.Write([]byte("data: [DONE]\n\n"))
			return
		}
		if strings.Contains(r.URL.Path, "/embeddings") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"embedding": make([]float32, 1536),
						"index":     0,
					},
				},
				"usage": map[string]interface{}{
					"total_tokens": 10,
				},
			})
			return
		}
	}))
	defer openAIServer.Close()

	// Repos
	repoBot := botRepo.NewBotRepo(db, cacheSvc)
	repoMsg := messageRepo.New(db)
	repoVec := vectorRepo.NewVectorRepository(db)

	// Services
	usageSvc := usage.NewUsageService(db)
	
	embSvc := embedding.NewEmbeddingService("test-key", "text-embedding-3-small", 1536)
	embSvc.SetBaseURL(openAIServer.URL)
	
	// Create Vector Service using real repo and mocked embedding service
	searchSvc := vector.NewSearchService(db, repoVec, embSvc)
	
	// Create Chat Service with mocked BaseURL
	chatService := chatSvc.NewChatService(
		embSvc,
		searchSvc,
		repoMsg,
		"gpt-3.5-turbo",
		0.7,
		3,
		"test-key",
	)
	chatService.SetBaseURL(openAIServer.URL)

	// Handler
	chatHandler := chat.NewChatHandler(repoBot, chatService, usageSvc)

	// Setup Data
	userID := "user_chat_full"
	botID := "bot_chat_full"
	
	shared.TruncateTable(db, "users")
	shared.TruncateTable(db, "bots")
	shared.TruncateTable(db, "usage_records")
	
	// Setup Data
	user := &models.User{
		ID: userID,
		Tier: "pro",
		CreditsBalance: 20.0,
		CreditsAllocated: 20.0,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	bot := &models.Bot{
		ID: botID,
		UserID: userID,
		Name: "Test Bot",
	}
	if err := db.Create(bot).Error; err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Test Route
	r := gin.New()
	r.POST("/api/bots/:id/chat", func(c *gin.Context) {
		// Mock Auth
		c.Set("user_id", userID)
		chatHandler.StreamChat(c)
	})

	// Perform Request
	w := httptest.NewRecorder()
	body, _ := json.Marshal(models.ChatRequest{Message: "Hello"})
	req, _ := http.NewRequest("POST", "/api/bots/"+botID+"/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Hello")

	// Usage should be tracked by Handler -> UsageService
	// Wait, ChatHandler calls usageSvc.TrackChatMessage?
	// Let's verify usage record created
	time.Sleep(100 * time.Millisecond) // async tracking?
	var count int64
	db.Model(&models.UsageRecord{}).Where("user_id = ? AND type = ?", userID, "chat_message").Count(&count)
	assert.Equal(t, int64(1), count)
}
