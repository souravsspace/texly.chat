package oauth_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/services/oauth"
)

func SetupTestRedis() (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

func SetupMemoryDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{})
	return db
}

func TestStateService(t *testing.T) {
	client, mr := SetupTestRedis()
	defer mr.Close()

	service := oauth.NewStateService(client)

	t.Run("GenerateState", func(t *testing.T) {
		state, err := service.GenerateState()
		assert.NoError(t, err)
		assert.NotEmpty(t, state)

		// Verify stored in Redis
		val, err := client.Get(context.Background(), "oauth_state:"+state).Result()
		assert.NoError(t, err)
		assert.Equal(t, "valid", val)
	})

	t.Run("ValidateState", func(t *testing.T) {
		state, _ := service.GenerateState()
		
		// Valid state
		assert.True(t, service.ValidateState(state))
		
		// Replay should fail
		assert.False(t, service.ValidateState(state))
		
		// Invalid state
		assert.False(t, service.ValidateState("invalid-state"))
	})
}

func TestOAuthService_GetGoogleAuthURL(t *testing.T) {
	db := SetupMemoryDB()
	cfg := configs.Config{
		GoogleClientID:     "test-client-id",
		GoogleClientSecret: "test-client-secret",
		GoogleRedirectURL:  "http://localhost:8080/callback",
	}
	
	service := oauth.NewOAuthService(cfg, db)
	state := "test-state"
	
	url := service.GetGoogleAuthURL(state)
	assert.Contains(t, url, "https://accounts.google.com/o/oauth2/auth")
	assert.Contains(t, url, "client_id=test-client-id")
	assert.Contains(t, url, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback")
	assert.Contains(t, url, "state=test-state")
}

// NOTE: HandleGoogleCallback is hard to test without mocking the Google API response
// and the OAuth2 exchange. The `oauth2` library makes it a bit tricky to mock without
// overriding the Endpoint to a local test server.
// For this test, verifying the URL generation covers the configuration part.
// The StateService test covers the security part.
// Integration testing would be better for the full flow.
