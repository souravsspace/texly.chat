package polar

import (
	"context"
	"fmt"

	polargo "github.com/polarsource/polar-go"
	"github.com/polarsource/polar-go/models/components"
	"github.com/polarsource/polar-go/models/operations"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/models"
)

type PolarService struct {
	client *polargo.Polar
	config configs.Config
}

func NewPolarService(cfg configs.Config, opts ...polargo.SDKOption) *PolarService {
	// Initialize Polar client with access token
	options := []polargo.SDKOption{
		polargo.WithSecurity(cfg.PolarAccessToken),
		polargo.WithServerURL(cfg.PolarServerURL),
	}
	options = append(options, opts...)

	client := polargo.New(options...)

	return &PolarService{
		client: client,
		config: cfg,
	}
}

// CreateCheckoutSession generates a Polar checkout URL for a Pro subscription
func (s *PolarService) CreateCheckoutSession(userID, userEmail string) (string, error) {
	ctx := context.Background()

	// Create checkout session
	// We use the PolarProProductID from config.
	// We pass user_id in metadata to link the subscription back to the user in webhooks.
	req := components.CheckoutCreate{
		Products:      []string{s.config.PolarProProductID},
		SuccessURL:    polargo.String(s.config.FrontendURL + "/dashboard/billing?success=true"),
		CustomerEmail: polargo.String(userEmail),
		Metadata: map[string]components.CheckoutCreateMetadata{
			"user_id": components.CreateCheckoutCreateMetadataStr(userID),
		},
	}

	resp, err := s.client.Checkouts.Create(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	if resp.Checkout == nil {
		return "", fmt.Errorf("checkout response is empty")
	}

	return resp.Checkout.URL, nil
}

// CreateCustomerPortalSession generates a link to manage subscription
func (s *PolarService) CreateCustomerPortalSession(polarCustomerID string) (string, error) {
	ctx := context.Background()

	// Use CustomerSessions to create a portal session
	// The SDK uses a union type for the request body, we must use the specific constructor.
	req := operations.CreateCustomerSessionsCreateCustomerSessionCreateCustomerSessionCustomerIDCreate(
		components.CustomerSessionCustomerIDCreate{
			CustomerID: polarCustomerID,
		},
	)

	resp, err := s.client.CustomerSessions.Create(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create customer portal session: %w", err)
	}

	if resp.CustomerSession == nil {
		return "", fmt.Errorf("customer session response is empty")
	}

	return "https://polar.sh/portal/" + resp.CustomerSession.Token, nil
}

// CreateUsageInvoice creates a usage invoice (if applicable in this flow)
// For Polar, we trigger a checkout session for the specific amount/product to "top up" or pay.
func (s *PolarService) CreateUsageInvoice(userID string, amount float64) (string, error) {
	ctx := context.Background()

	// Assuming we have a "Top Up" product ID in config or we use the Pro one as fallback for now.
	productID := s.config.PolarCreditsProductID
	if productID == "" {
		productID = s.config.PolarProProductID
	}

	req := components.CheckoutCreate{
		Products:   []string{productID}, // specific product for credits
		SuccessURL: polargo.String(s.config.FrontendURL + "/dashboard/billing?success=true&type=usage"),
		Metadata: map[string]components.CheckoutCreateMetadata{
			"user_id": components.CreateCheckoutCreateMetadataStr(userID),
			"type":    components.CreateCheckoutCreateMetadataStr("usage_charge"),
			"amount":  components.CreateCheckoutCreateMetadataStr(fmt.Sprintf("%.2f", amount)),
		},
	}

	resp, err := s.client.Checkouts.Create(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create usage checkout: %w", err)
	}

	if resp.Checkout == nil {
		return "", fmt.Errorf("checkout response is empty")
	}

	return resp.Checkout.URL, nil
}

// GetSubscription fetches subscription details
func (s *PolarService) GetSubscription(subscriptionID string) (*models.Subscription, error) {
	ctx := context.Background()

	resp, err := s.client.Subscriptions.Get(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if resp.Subscription == nil {
		return nil, fmt.Errorf("subscription not found or empty response")
	}

	return &models.Subscription{
		ID:     resp.Subscription.ID,
		Status: string(resp.Subscription.Status),
	}, nil
}

// CancelSubscription cancels a subscription
func (s *PolarService) CancelSubscription(subscriptionID string) error {
	ctx := context.Background()

	// Use Revoke for immediate cancellation
	_, err := s.client.Subscriptions.Revoke(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}
	return nil
}
