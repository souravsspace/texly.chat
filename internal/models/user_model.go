package models

import "time"

/*
* User represents a registered user in the system
 */
type User struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-"` // Nullable for OAuth users, removed gorm:not null
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	GoogleID     *string   `json:"google_id" gorm:"unique"`
	AuthProvider string    `json:"auth_provider" gorm:"default:'email'"` // email, google

	// Subscription fields
	Tier               string    `json:"tier" gorm:"default:'free'"` // "free", "pro", "enterprise"
	PolarCustomerID    *string   `json:"polar_customer_id" gorm:"uniqueIndex"`
	SubscriptionID     string    `json:"subscription_id"`
	SubscriptionStatus string    `json:"subscription_status"` // "active", "cancelled", "past_due"
	BillingCycleStart  time.Time `json:"billing_cycle_start"`
	BillingCycleEnd    time.Time `json:"billing_cycle_end"`

	// Credits & Usage (Pro tier)
	CreditsBalance     float64 `json:"credits_balance" gorm:"default:0"`      // Current balance in USD
	CreditsAllocated   float64 `json:"credits_allocated" gorm:"default:0"`    // Monthly allocation
	CurrentPeriodUsage float64 `json:"current_period_usage" gorm:"default:0"` // Total usage in USD this period

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/*
* LoginRequest holds the credentials for user login
 */
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
* SignupRequest holds data for creating a new user
 */
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

/*
* AuthResponse is the response payload for successful authentication
 */
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
