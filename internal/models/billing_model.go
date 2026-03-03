package models

import "time"

// Subscription represents a localized view of a subscription
type Subscription struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

/*
* BillingUsageResponse represents the billing usage information returned to the frontend
 */
type BillingUsageResponse struct {
	CreditsBalance     float64   `json:"credits_balance"`
	CreditsAllocated   float64   `json:"credits_allocated"`
	CurrentPeriodUsage float64   `json:"current_period_usage"`
	Tier               string    `json:"tier"`
	BillingCycleEnd    time.Time `json:"billing_cycle_end"`
}
