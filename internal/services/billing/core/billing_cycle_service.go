package billing

import (
	"fmt"
	"time"

	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
	credits "github.com/souravsspace/texly.chat/internal/services/billing/credits"
	polar "github.com/souravsspace/texly.chat/internal/services/billing/polar"
	"gorm.io/gorm"
)

/*
* BillingCycleService manages the monthly billing cycles for users
* It handles credit resets, overage calculation, and invoicing
 */
type BillingCycleService struct {
	db             *gorm.DB
	creditsService *credits.CreditsService
	polarService   *polar.PolarService
}

/*
* NewBillingCycleService creates a new instance
 */
func NewBillingCycleService(db *gorm.DB, creditsService *credits.CreditsService, polarService *polar.PolarService) *BillingCycleService {
	return &BillingCycleService{
		db:             db,
		creditsService: creditsService,
		polarService:   polarService,
	}
}

/*
* ProcessDueCycles finds users whose billing cycle has ended and processes them
* expected to be called by a daily cron job
 */
func (s *BillingCycleService) ProcessDueCycles() error {
	var users []models.User
	now := time.Now()

	// Find users with billing_cycle_end in the past
	// And tier is 'pro' (Free tier doesn't need credit reset logic usually, or handled differently)
	// We might want to process 'free' too if we track stats, but for billing purposes 'pro' is key.
	if err := s.db.Where("tier = ? AND billing_cycle_end <= ?", "pro", now).Find(&users).Error; err != nil {
		return fmt.Errorf("failed to fetch due users: %w", err)
	}

	fmt.Printf("[BillingCycle] Processing %d users for billing cycle reset\n", len(users))

	for _, user := range users {
		if err := s.processUserCycle(&user); err != nil {
			fmt.Printf("[BillingCycle] Error processing user %s: %v\n", user.ID, err)
			// Continue with other users
		}
	}

	return nil
}

/*
* processUserCycle handles the logic for a single user
 */
func (s *BillingCycleService) processUserCycle(user *models.User) error {
	// 1. Calculate Overage
	// Overage is usage beyond allocated credits
	// If usage > allocated, the difference is billable
	// NOTE: This assumes CurrentPeriodUsage tracks Total Value Consumed.
	overage := 0.0
	limit := user.CreditsAllocated
	if user.CurrentPeriodUsage > limit {
		overage = user.CurrentPeriodUsage - limit
	}

	// 2. Bill Overage if any
	if overage > 0.01 { // Ignore tiny amounts
		fmt.Printf("[BillingCycle] User %s has overage of $%.2f. Triggering invoice.\n", user.ID, overage)
		invoiceURL, err := s.polarService.CreateUsageInvoice(user.ID, overage)
		if err != nil {
			return fmt.Errorf("failed to create invoice: %w", err)
		}
		// In a real automated system, we'd email this URL or auto-charge if payment method on file.
		// For now, we just log it.
		// TODO: Send email notification with invoiceURL
		fmt.Printf("[BillingCycle] Invoice created for user %s: %s\n", user.ID, invoiceURL)
	}

	// 3. Reset Credits for next month
	// RefreshMonthlyCredits resets balance to user.CreditsAllocated (or default tier amount)
	// We need to ensure user.CreditsAllocated is set correctly (e.g. from Pricing config)
	// If it's 0 (legacy), we fetch from config
	if user.CreditsAllocated == 0 {
		tierLimits := configs.GetTierLimits(user.Tier)
		user.CreditsAllocated = tierLimits.IncludedCredits
	}

	// Reset balance to allocated amount
	user.CreditsBalance = user.CreditsAllocated

	// 4. Reset Usage Counters
	user.CurrentPeriodUsage = 0

	// 5. Advance Cycle Dates
	// If BillingCycleStart is unset, set it to now
	if user.BillingCycleStart.IsZero() {
		user.BillingCycleStart = time.Now()
	}
	// If BillingCycleEnd is unset or far past, catch up
	// We just want next month.
	// Logic: NewStart = OldEnd. NewEnd = NewStart + 1 Month.
	// But if OldEnd is very old (e.g. system down), we might want to skip to Now.
	// For simplicity, let's align to Now if it's too far off, or stick to schedule.
	// Sticky schedule is better for subscriptions.
	newStart := user.BillingCycleEnd
	if newStart.IsZero() {
		newStart = time.Now()
	}
	newEnd := newStart.AddDate(0, 1, 0)

	user.BillingCycleStart = newStart
	user.BillingCycleEnd = newEnd

	// 6. Save User
	if err := s.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to save user updates: %w", err)
	}

	fmt.Printf("[BillingCycle] specific user %s cycle reset. Next cycle ends: %s\n", user.ID, newEnd)
	return nil
}
