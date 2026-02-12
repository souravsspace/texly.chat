package billing_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/souravsspace/texly.chat/internal/models"
	billing "github.com/souravsspace/texly.chat/internal/services/billing/core"
	credits "github.com/souravsspace/texly.chat/internal/services/billing/credits"
	polar "github.com/souravsspace/texly.chat/internal/services/billing/polar"
	"github.com/souravsspace/texly.chat/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestBillingCycleService(t *testing.T) {
	db := shared.SetupTestDB()
	
	// Mock Polar Server
	polarServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/checkouts/custom/" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"url": "http://mock-checkout-url"})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer polarServer.Close()

	cfg := shared.GetTestConfig()
	cfg.PolarServerURL = polarServer.URL
	cfg.PolarAccessToken = "test-token"
	
	creditsSvc := credits.NewCreditsService(db)
	polarSvc := polar.NewPolarService(cfg)
	cycleSvc := billing.NewBillingCycleService(db, creditsSvc, polarSvc)

	userID := "user_test_cycle"

	// Helper to setup user
	setupUser := func() {
		shared.TruncateTable(db, "usage_records")
		shared.TruncateTable(db, "users")

		// Create user with cycle ending YESTERDAY (-25h to be safe)
		yesterday := time.Now().Add(-25 * time.Hour)
		cycleStart := yesterday.Add(-30 * 24 * time.Hour)
		
		user := &models.User{
			ID:                userID,
			Tier:              "pro",
			CreditsBalance:    0.0,
			CreditsAllocated:  20.0,
			BillingCycleStart: cycleStart,
			BillingCycleEnd:   yesterday,
			CurrentPeriodUsage: 25.0,
		}
		if err := db.Create(user).Error; err != nil {
			panic(err)
		} 
			// 25 usage > 20 allocated. Need to invoice 5.0. Balance 0.
	}

	t.Run("ProcessDueCycles_InvoicesAndResets", func(t *testing.T) {
		setupUser()
		
		// Debug check
		var u models.User
		db.First(&u, "id = ?", userID)
		t.Logf("User in DB: ID=%s CycleEnd=%v Now=%v", u.ID, u.BillingCycleEnd, time.Now())

		// Run ProcessDueCycles
		err := cycleSvc.ProcessDueCycles()
		assert.NoError(t, err)

		// Verify User Updated
		var user models.User
		db.First(&user, "id = ?", userID)

		// Credits should be reset to 20
		assert.Equal(t, 20.0, user.CreditsBalance)

		// Usage counter reset
		assert.Equal(t, 0.0, user.CurrentPeriodUsage)
		
		// Next cycle should be +1 month
		assert.True(t, user.BillingCycleEnd.After(time.Now()))
	})

	t.Run("ProcessDueCycles_NoOverage", func(t *testing.T) {
		setupUser()
		
		// Set usage < allocated
		db.Exec("UPDATE users SET current_period_usage = 10.0 WHERE id = ?", userID)
		
		err := cycleSvc.ProcessDueCycles()
		assert.NoError(t, err)

		var user models.User
		db.First(&user, "id = ?", userID)

		// Credits reset to 20
		assert.Equal(t, 20.0, user.CreditsBalance)
	})
}
