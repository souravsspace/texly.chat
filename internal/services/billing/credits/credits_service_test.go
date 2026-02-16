package credits_test

import (
	"testing"

	billing "github.com/souravsspace/texly.chat/internal/services/billing/credits"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestCreditsService(t *testing.T) {
	// Setup DB
	db := shared.SetupTestDB()
	// defer shared.CleanupTestDB(db) // Keep for inspection if needed, or rely on distinct test runs

	// Clear relevant tables
	shared.TruncateTable(db, "users")

	svc := billing.NewCreditsService(db)

	userID := "user_test_credits"

	t.Run("AddCredits", func(t *testing.T) {
		// Mock User
		db.Exec("INSERT INTO users (id, email, credits_balance, tier) VALUES (?, ?, ?, ?)", userID, "test@example.com", 0, "pro")

		err := svc.AddCredits(userID, 50.0)
		assert.NoError(t, err)

		bal, err := svc.GetCreditsBalance(userID)
		assert.NoError(t, err)
		assert.Equal(t, 50.0, bal)
	})

	t.Run("RefreshMonthlyCredits", func(t *testing.T) {
		// Set allocation to 20 (Pro default) or rely on existing logic
		// Update user allocated credits
		db.Exec("UPDATE users SET credits_allocated = 20.0 WHERE id = ?", userID)

		err := svc.RefreshMonthlyCredits(userID)
		assert.NoError(t, err)

		bal, err := svc.GetCreditsBalance(userID)
		assert.NoError(t, err)
		assert.Equal(t, 20.0, bal) // Should reset to 20
	})

	// Cleanup
	shared.TruncateTable(db, "users")
}
