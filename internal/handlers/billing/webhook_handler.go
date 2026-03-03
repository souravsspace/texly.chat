package billing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	// Read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("❌ Webhook: Failed to read body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Verify signature if webhook secret is configured
	if h.config.PolarWebhookSecret != "" {
		if err := h.verifySignature(c.Request.Header, body); err != nil {
			log.Printf("❌ Webhook: Signature verification failed: %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid signature"})
			return
		}
	}

	// Verify signature if webhook secret is configured
	if h.config.PolarWebhookSecret != "" {
		if err := h.verifySignature(c.Request.Header, body); err != nil {
			log.Printf("❌ Webhook: Signature verification failed: %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid signature"})
			return
		}
	}

	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("❌ Webhook: Invalid JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	eventType, ok := event["type"].(string)
	if !ok {
		log.Printf("⚠️  Webhook: No type field, ignoring")
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	log.Printf("📨 Webhook: Received event type: %s", eventType)

	// Route by event type
	switch eventType {
	case "subscription.created":
		err = h.handleSubscriptionCreated(event)
	case "subscription.updated":
		err = h.handleSubscriptionUpdated(event)
	case "subscription.active":
		err = h.handleSubscriptionActive(event)
	case "subscription.canceled": // Polar uses US spelling
		err = h.handleSubscriptionCanceled(event)
	case "order.created":
		err = h.handleOrderCreated(event)
	default:
		log.Printf("ℹ️  Webhook: Unhandled event type: %s", eventType)
	}

	if err != nil {
		log.Printf("❌ Webhook: Handler error for %s: %v", eventType, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("✅ Webhook: Successfully processed %s", eventType)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// verifySignature verifies the webhook signature using Standard Webhooks spec
// https://github.com/standard-webhooks/standard-webhooks
// Polar uses standard webhooks format with separate headers:
// - Webhook-Signature: v1,<base64_signature>
// - Webhook-Timestamp: <unix_timestamp>
// - Webhook-Id: <event_id>
func (h *WebhookHandler) verifySignature(headers http.Header, body []byte) error {
	signature := headers.Get("Webhook-Signature")
	timestamp := headers.Get("Webhook-Timestamp")
	webhookID := headers.Get("Webhook-Id")

	if signature == "" {
		return fmt.Errorf("missing Webhook-Signature header")
	}
	if timestamp == "" {
		return fmt.Errorf("missing Webhook-Timestamp header")
	}
	if webhookID == "" {
		return fmt.Errorf("missing Webhook-Id header")
	}

	// Parse signature header: "v1,<base64_signature>" or "v1 <base64_signature>"
	parts := strings.Split(signature, ",")
	if len(parts) < 2 {
		parts = strings.Split(signature, " ")
	}
	if len(parts) < 2 {
		return fmt.Errorf("invalid signature format: expected 'v1,<signature>'")
	}

	version := strings.TrimSpace(parts[0])
	sig := strings.TrimSpace(parts[1])

	if version != "v1" {
		return fmt.Errorf("unsupported webhook signature version: %s", version)
	}

	// Verify timestamp is recent (within 5 minutes)
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}
	if time.Now().Unix()-ts > 300 {
		return fmt.Errorf("timestamp too old")
	}

	// Compute expected signature
	// Standard Webhooks spec: HMAC-SHA256(secret, webhook_id + "." + timestamp + "." + body)
	secret := []byte(h.config.PolarWebhookSecret)

	// Build the signed content: webhook_id.timestamp.body
	signedContent := webhookID + "." + timestamp + "." + string(body)

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signedContent))
	expectedSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func (h *WebhookHandler) handleSubscriptionCreated(event map[string]interface{}) error {
	data, _ := event["data"].(map[string]interface{})
	metadata, _ := data["metadata"].(map[string]interface{})
	userID, _ := metadata["user_id"].(string)

	if userID == "" {
		log.Printf("⚠️  Webhook: subscription.created has no user_id in metadata")
		return nil // Can't link to user
	}

	log.Printf("📝 Webhook: Processing subscription.created for user %s", userID)

	// Update user to Pro
	updates := map[string]interface{}{
		"tier":                configs.TierPro,
		"subscription_status": "active",
		"polar_customer_id":   data["customer_id"],
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

	err := h.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		log.Printf("❌ Webhook: Failed to update user %s: %v", userID, err)
		return err
	}

	log.Printf("✅ Webhook: User %s upgraded to Pro", userID)
	return nil
}

func (h *WebhookHandler) handleSubscriptionActive(event map[string]interface{}) error {
	data, _ := event["data"].(map[string]interface{})
	metadata, _ := data["metadata"].(map[string]interface{})
	userID, _ := metadata["user_id"].(string)

	if userID == "" {
		// Try to find user by subscription_id
		subID, _ := data["id"].(string)
		if subID == "" {
			log.Printf("⚠️  Webhook: subscription.active has no user_id or subscription_id")
			return nil
		}

		var user models.User
		if err := h.db.Where("subscription_id = ?", subID).First(&user).Error; err != nil {
			log.Printf("⚠️  Webhook: Could not find user for subscription %s", subID)
			return nil
		}
		userID = user.ID
	}

	log.Printf("📝 Webhook: Processing subscription.active for user %s", userID)

	// Update subscription status to active
	updates := map[string]interface{}{
		"tier":                configs.TierPro,
		"subscription_status": "active",
	}

	// Update billing cycle dates if provided
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

	// Refresh credits if it's a new billing cycle
	status, _ := data["status"].(string)
	if status == "active" {
		updates["credits_balance"] = configs.ProIncludedCredits
		updates["credits_allocated"] = configs.ProIncludedCredits
		updates["current_period_usage"] = 0.0
	}

	err := h.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		log.Printf("❌ Webhook: Failed to activate subscription for user %s: %v", userID, err)
		return err
	}

	log.Printf("✅ Webhook: Subscription activated for user %s", userID)
	return nil
}

func (h *WebhookHandler) handleSubscriptionUpdated(event map[string]interface{}) error {
	data, _ := event["data"].(map[string]interface{})
	metadata, _ := data["metadata"].(map[string]interface{})
	userID, _ := metadata["user_id"].(string)

	if userID == "" {
		// Try to find user by subscription_id
		subID, _ := data["id"].(string)
		if subID == "" {
			log.Printf("⚠️  Webhook: subscription.updated has no user_id or subscription_id")
			return nil
		}

		var user models.User
		if err := h.db.Where("subscription_id = ?", subID).First(&user).Error; err != nil {
			log.Printf("⚠️  Webhook: Could not find user for subscription %s", subID)
			return nil
		}
		userID = user.ID
	}

	log.Printf("📝 Webhook: Processing subscription.updated for user %s", userID)

	updates := map[string]interface{}{}

	// Update status if changed
	if status, ok := data["status"].(string); ok {
		updates["subscription_status"] = status
	}

	// Update billing cycle dates
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

	if len(updates) > 0 {
		err := h.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
		if err != nil {
			log.Printf("❌ Webhook: Failed to update subscription for user %s: %v", userID, err)
			return err
		}
		log.Printf("✅ Webhook: Subscription updated for user %s", userID)
	}

	return nil
}

func (h *WebhookHandler) handleSubscriptionCanceled(event map[string]interface{}) error {
	data, _ := event["data"].(map[string]interface{})
	metadata, _ := data["metadata"].(map[string]interface{})
	userID, _ := metadata["user_id"].(string)

	if userID == "" {
		// Try to find user by subscription_id
		subID, _ := data["id"].(string)
		if subID == "" {
			log.Printf("⚠️  Webhook: subscription.canceled has no user_id or subscription_id")
			return nil
		}

		var user models.User
		if err := h.db.Where("subscription_id = ?", subID).First(&user).Error; err != nil {
			log.Printf("⚠️  Webhook: Could not find user for subscription %s", subID)
			return nil
		}
		userID = user.ID
	}

	log.Printf("📝 Webhook: Processing subscription.canceled for user %s", userID)

	// Update subscription status to canceled
	// Note: We keep tier as "pro" until billing cycle ends
	// This allows users to continue using Pro features until period end
	updates := map[string]interface{}{
		"subscription_status": "canceled",
	}

	err := h.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		log.Printf("❌ Webhook: Failed to cancel subscription for user %s: %v", userID, err)
		return err
	}

	log.Printf("✅ Webhook: Subscription canceled for user %s (tier remains Pro until period end)", userID)
	return nil
}

func (h *WebhookHandler) handleOrderCreated(event map[string]interface{}) error {
	data, _ := event["data"].(map[string]interface{})
	metadata, _ := data["metadata"].(map[string]interface{})
	userID, _ := metadata["user_id"].(string)

	if userID == "" {
		log.Printf("⚠️  Webhook: order.created has no user_id in metadata")
		return nil
	}

	log.Printf("📝 Webhook: Processing order.created for user %s", userID)

	// Get order details
	orderID, _ := data["id"].(string)
	amount, _ := data["amount"].(float64) // Amount in cents
	productID, _ := data["product_id"].(string)

	log.Printf("💰 Webhook: Order %s created - Amount: $%.2f, Product: %s",
		orderID, amount/100, productID)

	// Check if this is a credit top-up purchase
	if productID == h.config.PolarCreditsProductID {
		// This is a credit purchase - add credits to user balance
		creditsToAdd := amount / 100 // Convert cents to dollars (1:1 with credits)

		err := h.db.Model(&models.User{}).
			Where("id = ?", userID).
			UpdateColumn("credits_balance", gorm.Expr("credits_balance + ?", creditsToAdd)).
			Error

		if err != nil {
			log.Printf("❌ Webhook: Failed to add credits for user %s: %v", userID, err)
			return err
		}

		log.Printf("✅ Webhook: Added $%.2f credits to user %s", creditsToAdd, userID)
	}

	// Log successful payment (could be used for analytics/reporting)
	log.Printf("✅ Webhook: Order processed for user %s", userID)
	return nil
}
