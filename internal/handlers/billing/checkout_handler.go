package billing

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/models"
	"github.com/souravsspace/texly.chat/internal/services/billing/polar"
	"github.com/souravsspace/texly.chat/internal/services/billing/usage"
)

type CheckoutHandler struct {
	polarService *polar.PolarService
	usageService *usage.UsageService
}

func NewCheckoutHandler(polarService *polar.PolarService, usageService *usage.UsageService) *CheckoutHandler {
	return &CheckoutHandler{
		polarService: polarService,
		usageService: usageService,
	}
}

// CreateCheckoutSession handles POST /api/billing/checkout
func (h *CheckoutHandler) CreateCheckoutSession(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	// In a real app, you might validate the request body for specific plan details,
	// but for now we assume upgrading to "pro".
	url, err := h.polarService.CreateCheckoutSession(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// CreatePortalSession handles POST /api/billing/portal
func (h *CheckoutHandler) CreatePortalSession(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	if user.PolarCustomerID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not have a billing customer ID"})
		return
	}

	url, err := h.polarService.CreateCustomerPortalSession(*user.PolarCustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create portal session: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GetUsage handles GET /api/billing/usage
func (h *CheckoutHandler) GetUsage(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	currentUsage, err := h.usageService.GetCurrentUsage(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get usage stats: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"credits_balance":      user.CreditsBalance,
		"credits_allocated":    user.CreditsAllocated,
		"current_period_usage": currentUsage,
		"tier":                 user.Tier,
		"billing_cycle_end":    user.BillingCycleEnd,
	})
}

// PayUsage handles POST /api/billing/pay-usage
// Allows user to manually pay off usage or top up credits
func (h *CheckoutHandler) PayUsage(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	url, err := h.polarService.CreateUsageInvoice(user.ID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment session: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
