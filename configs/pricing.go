package configs

/*
 * Pricing Configuration
 *
 * Central source of truth for all billing, pricing, and tier limits.
 * Change values here to adjust pricing across the entire application.
 *
 * Model: $20/month subscription with $20 in usage credits.
 * - Credits are always charged even if unused (subscription fee).
 * - Unused credits do NOT roll over to the next billing cycle.
 * - Once credits are exhausted, pay-as-you-go pricing applies.
 */

// ---------------------------------------------------------------------------
// Subscription Pricing
// ---------------------------------------------------------------------------

const (
	// ProMonthlyPriceCents is the base subscription fee in cents ($20.00)
	ProMonthlyPriceCents = 2000

	// ProIncludedCredits is the monthly credit allocation in USD
	ProIncludedCredits = 20.00
)

// ---------------------------------------------------------------------------
// Pay-As-You-Go Unit Pricing (70% profit margin)
// ---------------------------------------------------------------------------

const (
	// PricePerMessage is the charge per chat message in USD
	// Cost: $0.0003 → Price: $0.001
	PricePerMessage = 0.001

	// CostPerMessage is the actual cost per chat message in USD
	CostPerMessage = 0.0003

	// PricePerEmbedding1KTokens is the charge per 1K embedding tokens in USD
	// Cost: $0.00006 → Price: $0.0002
	PricePerEmbedding1KTokens = 0.0002

	// CostPerEmbedding1KTokens is the actual cost per 1K embedding tokens in USD
	CostPerEmbedding1KTokens = 0.00006

	// PricePerGBStorageMonthly is the charge per GB of storage per month in USD
	// Cost: $0.03 → Price: $0.10
	PricePerGBStorageMonthly = 0.10

	// CostPerGBStorageMonthly is the actual cost per GB of storage per month in USD
	CostPerGBStorageMonthly = 0.03

	// PricePerExtraBotMonthly is the charge per additional bot per month in USD
	// Cost: $1.50 → Price: $5.00
	PricePerExtraBotMonthly = 5.00

	// CostPerExtraBotMonthly is the actual cost per additional bot per month in USD
	CostPerExtraBotMonthly = 1.50
)

// ---------------------------------------------------------------------------
// Billing Rules
// ---------------------------------------------------------------------------

const (
	// GracePeriodDays is the number of days after payment failure before downgrade
	GracePeriodDays = 7

	// MinChargeThresholdUSD is the minimum overage amount before triggering a charge
	MinChargeThresholdUSD = 0.50

	// CreditsRollOver controls whether unused credits carry to the next cycle
	CreditsRollOver = false
)

// ---------------------------------------------------------------------------
// Tier Limits
// ---------------------------------------------------------------------------

// TierLimits defines the resource limits for a subscription tier.
// A value of -1 means unlimited.
type TierLimits struct {
	Tier             string
	MaxBots          int
	MaxMessagesPerMo int     // -1 = unlimited (pay-as-you-go)
	MaxSourcesPerBot int     // -1 = unlimited
	MaxStorageGB     float64 // -1 = unlimited
	MaxOriginsPerBot int     // -1 = unlimited
	IncludedCredits  float64 // monthly credit allocation in USD
	IncludedBots     int     // bots included before extra-bot charges apply
}

const (
	TierFree       = "free"
	TierPro        = "pro"
	TierEnterprise = "enterprise"
)

// Tiers maps tier names to their resource limits.
var Tiers = map[string]TierLimits{
	TierFree: {
		Tier:             TierFree,
		MaxBots:          1,
		MaxMessagesPerMo: 100,
		MaxSourcesPerBot: 5,
		MaxStorageGB:     0.01, // 10 MB
		MaxOriginsPerBot: 1,
		IncludedCredits:  0,
		IncludedBots:     1,
	},
	TierPro: {
		Tier:             TierPro,
		MaxBots:          -1, // unlimited (extra bots charged)
		MaxMessagesPerMo: -1, // unlimited (pay-as-you-go)
		MaxSourcesPerBot: 50,
		MaxStorageGB:     1,
		MaxOriginsPerBot: 10,
		IncludedCredits:  ProIncludedCredits,
		IncludedBots:     5, // 5 included, $5/mo per additional
	},
	TierEnterprise: {
		Tier:             TierEnterprise,
		MaxBots:          -1,
		MaxMessagesPerMo: -1,
		MaxSourcesPerBot: -1,
		MaxStorageGB:     -1,
		MaxOriginsPerBot: -1,
		IncludedCredits:  0, // custom billing
		IncludedBots:     -1,
	},
}

// ---------------------------------------------------------------------------
// Cost Calculation Helpers
// ---------------------------------------------------------------------------

// CalculateMessageCost returns the total price for a given message count.
func CalculateMessageCost(count int) float64 {
	return float64(count) * PricePerMessage
}

// CalculateEmbeddingCost returns the total price for a given token count.
func CalculateEmbeddingCost(tokens int) float64 {
	return float64(tokens) / 1000.0 * PricePerEmbedding1KTokens
}

// CalculateStorageCost returns the monthly price for a given storage amount in GB.
func CalculateStorageCost(gb float64) float64 {
	return gb * PricePerGBStorageMonthly
}

// CalculateExtraBotCost returns the monthly price for extra bots beyond the included count.
func CalculateExtraBotCost(totalBots int, tier string) float64 {
	t, ok := Tiers[tier]
	if !ok {
		return 0
	}
	if t.IncludedBots == -1 {
		return 0 // unlimited
	}
	extra := totalBots - t.IncludedBots
	if extra <= 0 {
		return 0
	}
	return float64(extra) * PricePerExtraBotMonthly
}

// GetTierLimits returns the limits for a given tier name, defaulting to Free.
func GetTierLimits(tier string) TierLimits {
	if t, ok := Tiers[tier]; ok {
		return t
	}
	return Tiers[TierFree]
}
