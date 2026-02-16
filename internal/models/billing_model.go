package models

// Subscription represents a localized view of a subscription
type Subscription struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
