package billing

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/internal/repo/user"
	"github.com/souravsspace/texly.chat/internal/services/billing/polar"
	"github.com/souravsspace/texly.chat/internal/services/billing/usage"
)

type CheckoutHandler struct {
	polarService *polar.PolarService
	usageService *usage.UsageService
	userRepo     *user.UserRepo
}

func NewCheckoutHandler(polarService *polar.PolarService, usageService *usage.UsageService, userRepo *user.UserRepo) *CheckoutHandler {
	return &CheckoutHandler{
		polarService: polarService,
		usageService: usageService,
		userRepo:     userRepo,
	}
}

// CreateCheckoutSession handles POST /api/billing/checkout
func (h *CheckoutHandler) CreateCheckoutSession(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// In a real app, you might validate the request body for specific plan details,
	// but for now we assume upgrading to "pro".
	url, err := h.polarService.CreateCheckoutSession(user.ID, user.Email)
	if err != nil {
		fmt.Printf("❌ CheckoutHandler: Failed to create checkout session: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// CreatePortalSession handles POST /api/billing/portal
func (h *CheckoutHandler) CreatePortalSession(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

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
	userID := c.GetString("user_id")
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

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
	userID := c.GetString("user_id")
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

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
