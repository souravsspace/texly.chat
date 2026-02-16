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
		// allow /v1/checkouts or /v1/checkouts/custom depending on SDK version
		if r.URL.Path == "/v1/checkouts/" || r.URL.Path == "/v1/checkouts/custom/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			// Polar SDK expects the Checkout object directly in the body (likely).
			response := map[string]interface{}{
				"id":                         "checkout_123",
				"created_at":                 time.Now().Format(time.RFC3339),
				"payment_processor":          "stripe",
				"status":                     "open",
				"client_secret":              "secret_123",
				"url":                        "http://mock-checkout-url",
				"expires_at":                 time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"success_url":                "http://example.com/success",
				"amount":                     500,
				"tax_amount":                 0,
				"currency":                   "usd",
				"total_amount":               500,
				"net_amount":                 500,
				"discount_amount":            0,
				"organization_id":            "org_123",
				"allow_discount_codes":       true,
				"require_billing_address":    false,
				"is_discount_applicable":     false,
				"is_free_product_price":      false,
				"is_payment_required":        true,
				"is_payment_setup_required":  false,
				"is_payment_form_required":   true,
				"customer_ip_address":        nil,
				"customer_name":              "Test User",
				"customer_email":             "test@example.com",
				"customer_billing_address":   nil,
				"customer_id":                "cust_123",
				"is_business_customer":       false,
				"billing_address_fields":     map[string]interface{}{}, // Changed to map
				"customer_metadata":          map[string]interface{}{}, // Changed to map
				"payment_processor_metadata": map[string]interface{}{},
				"metadata":                   map[string]interface{}{},
				"products":                   []interface{}{},
				"product_price_id":           "price_123",
				"product_id":                 "prod_123",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		// Default interaction
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
			ID:                 userID,
			Tier:               "pro",
			CreditsBalance:     0.0,
			CreditsAllocated:   20.0,
			BillingCycleStart:  cycleStart,
			BillingCycleEnd:    yesterday,
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
