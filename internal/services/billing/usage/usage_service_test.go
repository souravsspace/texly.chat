package usage_test

import (
	"testing"
	"time"

	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/services/billing/usage"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestTrackChatMessage(t *testing.T) {
	db := shared.SetupTestDB()
	service := usage.NewUsageService(db)

	user := models.User{
		ID:             "user_1",
		CreditsBalance: 10.0,
		Tier:           configs.TierPro,
	}
	db.Create(&user)

	// Track 1 message
	err := service.TrackChatMessage(user.ID, "bot_1")
	assert.NoError(t, err)

	// Verify usage record
	var record models.UsageRecord
	err = db.First(&record).Error
	assert.NoError(t, err)
	assert.Equal(t, "chat_message", record.Type)
	assert.Equal(t, 1.0, record.Quantity)
	// configs.CalculateMessageCost(1) = PricePerMessage
	assert.InDelta(t, configs.PricePerMessage, record.Cost, 0.0001)

	// Verify user balance deduction
	var updatedUser models.User
	db.First(&updatedUser, "id = ?", user.ID)
	// Balance should decrease by cost
	expectedBalance := 10.0 - configs.PricePerMessage
	assert.InDelta(t, expectedBalance, updatedUser.CreditsBalance, 0.0001)
}

func TestTrackEmbedding(t *testing.T) {
	db := shared.SetupTestDB()
	service := usage.NewUsageService(db)

	// User created with 5.0 credits
	user := models.User{
		ID:             "user_2",
		CreditsBalance: 5.0,
		Tier:           configs.TierPro,
	}
	db.Create(&user)

	// Track 1000 tokens
	err := service.TrackEmbedding(user.ID, 1000)
	assert.NoError(t, err)

	var record models.UsageRecord
	err = db.First(&record, "user_id = ?", user.ID).Error
	assert.NoError(t, err)

	expectedCost := configs.PricePerEmbedding1KTokens
	assert.InDelta(t, expectedCost, record.Cost, 0.00001)

	var updatedUser models.User
	db.First(&updatedUser, "id = ?", user.ID)
	assert.InDelta(t, 5.0-expectedCost, updatedUser.CreditsBalance, 0.0001)
}

func TestGetCurrentUsage(t *testing.T) {
	db := shared.SetupTestDB()
	service := usage.NewUsageService(db)

	userID := "user_3"
	// Create some usage records manually
	db.Create(&models.UsageRecord{
		UserID:    userID,
		Cost:      1.50,
		BilledAt:  time.Time{}, // Not billed yet
		CreatedAt: time.Now(),
	})
	db.Create(&models.UsageRecord{
		UserID:    userID,
		Cost:      0.50,
		BilledAt:  time.Time{}, // Not billed yet
		CreatedAt: time.Now(),
	})
	// Add a billed record (should be ignored)
	db.Create(&models.UsageRecord{
		UserID:    userID,
		Cost:      10.00,
		BilledAt:  time.Now(),
		CreatedAt: time.Now(),
	})

	total, err := service.GetCurrentUsage(userID)
	assert.NoError(t, err)
	assert.InDelta(t, 2.00, total, 0.0001)
}
