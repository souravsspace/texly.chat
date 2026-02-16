package billing

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/services/billing/credits"
	"github.com/souravsspace/texly.chat/internal/services/billing/polar"
	"gorm.io/gorm"
)

type WebhookHandler struct {
	db             *gorm.DB
	config         configs.Config
	creditsService *credits.CreditsService
	polarService   *polar.PolarService
}

func NewWebhookHandler(db *gorm.DB, cfg configs.Config, creditsService *credits.CreditsService, polarService *polar.PolarService) *WebhookHandler {
	return &WebhookHandler{
		db:             db,
		config:         cfg,
		creditsService: creditsService,
		polarService:   polarService,
	}
}

// HandleWebhook handles POST /api/billing/webhook
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	// 1. Verify Signature (TODO: Implement real signature check using h.config.PolarWebhookSecret)
	// Polar sends `Polar-Webhook-Signature` header.
	// For now, we trust the payload structure but in prod MUST verify signature.

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	eventType, ok := event["type"].(string)
	if !ok {
		// Polar events have `type` field
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	// 2. Route by event type
	switch eventType {
	case "subscription.created":
		err = h.handleSubscriptionCreated(event)
	case "subscription.updated":
		err = h.handleSubscriptionUpdated(event)
	case "subscription.cancelled": // or "subscription.canceled" check docs
		err = h.handleSubscriptionCancelled(event)
		// Add payment.succeeded etc.
	}

	if err != nil {
		// Log error but return 200 so Polar doesn't keep retrying exponentially if it's a logic error
		// Unless it's transient.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

func (h *WebhookHandler) handleSubscriptionCreated(event map[string]interface{}) error {
	data, _ := event["data"].(map[string]interface{})
	metadata, _ := data["metadata"].(map[string]interface{})
	userID, _ := metadata["user_id"].(string)

	if userID == "" {
		return nil // Can't link to user
	}

	// Update user to Pro
	updates := map[string]interface{}{
		"tier":                configs.TierPro,
		"subscription_status": "active",
		"polar_customer_id":   data["customer_id"], // check exact field name
		"subscription_id":     data["id"],
		// Initialize credits
		"credits_balance":   configs.ProIncludedCredits,
		"credits_allocated": configs.ProIncludedCredits,
	}

	// Handle dates
	if startStr, ok := data["current_period_start"].(string); ok {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			updates["billing_cycle_start"] = t
		}
	}
	if endStr, ok := data["current_period_end"].(string); ok {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			updates["billing_cycle_end"] = t
		}
	}

	return h.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (h *WebhookHandler) handleSubscriptionUpdated(event map[string]interface{}) error {
	// Similar logic to update status or dates
	return nil
}

func (h *WebhookHandler) handleSubscriptionCancelled(event map[string]interface{}) error {
	data, _ := event["data"].(map[string]interface{})
	// Find user by subscription_id if userID not in metadata (though it should be)
	subID, _ := data["id"].(string)
	if subID == "" {
		return nil
	}

	return h.db.Model(&models.User{}).Where("subscription_id = ?", subID).
		Update("subscription_status", "cancelled").Error
}
