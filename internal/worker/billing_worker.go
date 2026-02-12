package worker

import (
	"context"
	"fmt"
	"time"

	billing "github.com/souravsspace/texly.chat/internal/services/billing/core"
)

/*
* StartDailyBillingJob runs the billing cycle check daily
* It should be run in a goroutine
 */
func StartDailyBillingJob(ctx context.Context, billingSvc *billing.BillingCycleService) {
	// Calculate time until next midnight
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	initialDelay := nextMidnight.Sub(now)

	fmt.Printf("[BillingWorker] Daily billing job scheduled to start in %v\n", initialDelay)

	// Timer for the first run
	timer := time.NewTimer(initialDelay)
	defer timer.Stop()

	// Ticker for subsequent runs (24 hours)
	// Note: Ticker might drift, for strict cron use a cron library.
	// For MVP, this is sufficient.


	// Wait for the first run at midnight
	select {
	case <-ctx.Done():
		fmt.Println("[BillingWorker] Stopping billing job early")
		return
	case <-timer.C:
		runBillingCheck(billingSvc)
	}

	// Then run every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[BillingWorker] Stopping billing job")
			return
		case <-ticker.C:
			runBillingCheck(billingSvc)
		}
	}
}

func runBillingCheck(billingSvc *billing.BillingCycleService) {
	fmt.Println("[BillingWorker] Running daily cycle check...")
	if err := billingSvc.ProcessDueCycles(); err != nil {
		fmt.Printf("[BillingWorker] Error processing cycles: %v\n", err)
	}
}
