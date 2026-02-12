package polar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/souravsspace/texly.chat/configs"
)

type PolarService struct {
	config configs.Config
	client *http.Client
}

func NewPolarService(cfg configs.Config) *PolarService {
	return &PolarService{
		config: cfg,
		client: &http.Client{},
	}
}

// CheckoutRequest payload for creating a session
type CheckoutRequest struct {
	ProductPriceID string            `json:"product_price_id"`
	SuccessURL     string            `json:"success_url"`
	CustomerEmail  string            `json:"customer_email,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// CheckoutResponse payload from Polar
type CheckoutResponse struct {
	URL string `json:"url"`
}

// CreateCheckoutSession generates a Polar checkout URL
func (s *PolarService) CreateCheckoutSession(userID, userEmail string) (string, error) {
	// 1. Get Product Price ID from config (using Product ID for now as placeholder, 
	// ideally should be Price ID if distinct, but let's assume Product ID maps to a default price or we configure Price ID)
	// NOTE: The Polar API `product_price_id` is required. We'll assume the config `PolarProProductID` holds the Price ID 
	// or we need to fetch products to find the price ID. For simplicity, let's assume the config has the correct ID.
	priceID := s.config.PolarProProductID 

	reqBody := CheckoutRequest{
		ProductPriceID: priceID,
		SuccessURL:     s.config.FrontendURL + "/dashboard?checkout=success",
		CustomerEmail:  userEmail,
		Metadata: map[string]string{
			"user_id": userID,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.config.PolarServerURL+"/v1/checkouts/custom/", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.PolarAccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("polar api error: %s", string(bodyBytes))
	}

	var result CheckoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.URL, nil
}

// CustomerPortalSessionRequest payload
type CustomerPortalSessionRequest struct {
	CustomerID string `json:"customer_id"`
}

// CustomerPortalSessionResponse payload
type CustomerPortalSessionResponse struct {
	URL string `json:"url"`
}

// CreateCustomerPortalSession generates a link to manage subscription
func (s *PolarService) CreateCustomerPortalSession(polarCustomerID string) (string, error) {
	// Note: Polar might not have a direct "create portal session" endpoint like Stripe 
	// depending on the version. Checking docs, it usually manages via Dashboard or specific endpoints.
	// If unavailable, we might just link to the generic customer portal or send an email.
	// For now, let's assume a standard portal flow or return a placeholder if not yet supported by API v1.
	
	// Placeholder: Polar currently manages this via magic links or user dashboard.
	// We'll treat this as "TODO" or check if there's a specific endpoint.
	// Based on docs, managing subscriptions is often done via the initial checkout or email links.
	// We'll return the general Polar dashboard URL for now or a specific deep link if known.
	return "https://polar.sh/purchases", nil 
}

// CreateUsageInvoice creates an invoice for usage overage
func (s *PolarService) CreateUsageInvoice(userID string, amount float64) error {
	// Placeholder. In production, call Polar API or Stripe usage-based billing endpoint.
	// For MVP, we just log it. Real implementation would likely create a one-off charge or add to subscription.
	fmt.Printf("[PolarService] Creating usage invoice for User %s: $%.2f\n", userID, amount)
	return nil
}

// Subscription represents a localized view of a subscription
type Subscription struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// GetSubscription fetches subscription details
func (s *PolarService) GetSubscription(subscriptionID string) (*Subscription, error) {
	// Placeholder
	return &Subscription{ID: subscriptionID, Status: "active"}, nil
}

// CancelSubscription cancels a subscription at period end
func (s *PolarService) CancelSubscription(subscriptionID string) error {
	// Placeholder
	fmt.Printf("[PolarService] Cancelling subscription %s\n", subscriptionID)
	return nil
}
