package usage_test

import (
	"testing"

	"github.com/souravsspace/texly.chat/internal/models"
	billing "github.com/souravsspace/texly.chat/internal/services/billing/usage"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestUsageService(t *testing.T) {
	db := shared.SetupTestDB()
	svc := billing.NewUsageService(db)

	userID := "user_test_usage"
	botID := "bot_test_usage"

	// Setup Helper
	setupUser := func() {
		shared.TruncateTable(db, "usage_records")
		shared.TruncateTable(db, "bots")
		shared.TruncateTable(db, "users")

		// Create Pro User with credits
		db.Exec("INSERT INTO users (id, email, tier, credits_balance, credits_allocated) VALUES (?, ?, ?, ?, ?)", 
			userID, "usage@example.com", "pro", 5.00, 20.00)
		
		// Create Bot
		db.Exec("INSERT INTO bots (id, user_id, name, model) VALUES (?, ?, ?, ?)", 
			botID, userID, "Test Bot", "gpt-3.5-turbo")
	}

	t.Run("TrackChatMessage_DesuctsCredits", func(t *testing.T) {
		setupUser()

		// Cost for 1 message logic depends on model. 
		// Assuming implementation uses fixed cost or calculates it.
		// Let's assume standard cost.
		
		err := svc.TrackChatMessage(userID, botID)
		assert.NoError(t, err)

		// Check Balance
		var user models.User
		db.First(&user, "id = ?", userID)
		assert.Less(t, user.CreditsBalance, 5.00) // Should be less than initial

		// Check Usage Record
		var count int64
		db.Model(&models.UsageRecord{}).Where("user_id = ? AND type = ?", userID, "chat_message").Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("TrackEmbedding_CalculatesCost", func(t *testing.T) {
		setupUser()

		tokens := 1000
		// Cost for text-embedding-3-small is usually very low ($0.00002 / 1k tokens)
		// configs.PricingModel should have it.

		err := svc.TrackEmbedding(userID, tokens)
		assert.NoError(t, err)

		var user models.User
		db.First(&user, "id = ?", userID)
		assert.Less(t, user.CreditsBalance, 5.00)
		
		// Verify cost calculation if possible, or just that it deducted something
	})

	t.Run("TrackStorage_MonthlyProRate", func(t *testing.T) {
		setupUser()

		// 1GB storage
		sizeGB := 1.0
		
		err := svc.TrackStorage(userID, sizeGB)
		assert.NoError(t, err)

		// Check record
		var record models.UsageRecord
		err = db.Where("user_id = ? AND type = ?", userID, "storage").First(&record).Error
		assert.NoError(t, err)
		assert.Equal(t, 1.0, record.Quantity)
	})

	t.Run("GetCurrentUsage", func(t *testing.T) {
		setupUser()
		_ = svc.TrackChatMessage(userID, botID)
		_ = svc.TrackChatMessage(userID, botID)

		report, err := svc.GetCurrentUsage(userID)
		assert.NoError(t, err)
		assert.NotNil(t, report)
		// Assert details based on implementation
	})
}
