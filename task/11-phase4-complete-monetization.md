# Phase 4: Complete Monetization with Polar Integration

## Goal
Implement three-tier pricing system with Polar.sh integration, credit-based billing, and pay-as-you-go model with 70% profit margin.

---

## Tier Structure

### Free Tier
**Price**: $0/month
**Target**: Individuals, testing, hobby projects
**Limits**:
- 1 bot maximum
- 100 messages/month
- 5 sources per bot
- 10 MB file storage
- 1 allowed origin (widget domain)
- Community support only

### Pro Tier
**Price**: $20/month base subscription + pay-as-you-go
**Target**: Small businesses, freelancers, startups
**Included Credits**: $20 in usage credits/month (refreshed monthly)
**Pay-as-you-go Pricing** (70% profit margin):
- Chat messages: $0.001 per message (cost: $0.0003)
- Vector embeddings: $0.0002 per 1k tokens (cost: $0.00006)
- File storage: $0.10 per GB/month (cost: $0.03)
- Extra bot: $5/month per bot (cost: $1.50)

**Limits**:
- Unlimited messages (pay-as-you-go after credits exhausted)
- 5 bots included, $5/month per additional bot
- 50 sources per bot
- 1 GB file storage included
- 10 allowed origins per bot
- Email support (24h response)

### Enterprise Tier
**Price**: Custom (contact sales)
**Target**: Large organizations, high-volume users
**Features**:
- Unlimited bots (`-1` in database)
- Unlimited messages (`-1`)
- Unlimited sources (`-1`)
- Unlimited file storage (`-1`)
- Unlimited allowed origins (`-1`)
- Custom domain (white-label widget)
- Dedicated support (4h response SLA)
- Custom AI model options
- On-premise deployment option
- SSO/SAML integration

**Flags in Database**: All limits set to `-1` to indicate unlimited

---

## Backend Tasks

### Step 1: Update User & Subscription Models

#### 1.1 Update User Model
- [ ] Update `internal/models/user_model.go`:
  ```go
  type User struct {
      ID    string `json:"id" gorm:"primaryKey"`
      Email string `json:"email" gorm:"uniqueIndex;not null"`
      
      // Subscription fields
      Tier             string  `json:"tier" gorm:"default:'free'"` // "free", "pro", "enterprise"
      PolarCustomerID  string  `json:"polar_customer_id" gorm:"uniqueIndex"`
      SubscriptionID   string  `json:"subscription_id"`
      SubscriptionStatus string `json:"subscription_status"` // "active", "cancelled", "past_due"
      
      // Credits (Pro tier only)
      CreditsBalance   float64 `json:"credits_balance" gorm:"default:0"` // Current balance in USD
      CreditsAllocated float64 `json:"credits_allocated" gorm:"default:0"` // Monthly allocation
      BillingCycleStart time.Time `json:"billing_cycle_start"`
      BillingCycleEnd   time.Time `json:"billing_cycle_end"`
      
      // Usage tracking
      CurrentPeriodUsage float64 `json:"current_period_usage" gorm:"default:0"` // Total usage in USD
      
      // Existing fields...
      CreatedAt time.Time      `json:"created_at"`
      UpdatedAt time.Time      `json:"updated_at"`
  }
  ```

#### 1.2 Create Tier Definition Model
- [ ] Create `internal/models/tier_model.go`:
  ```go
  type TierLimits struct {
      Tier              string
      MaxBots           int     // -1 = unlimited
      MaxMessages       int     // -1 = unlimited
      MaxSources        int     // -1 = unlimited
      MaxStorageGB      int     // -1 = unlimited
      MaxAllowedOrigins int     // -1 = unlimited
      IncludedCredits   float64 // Monthly credits (Pro only)
  }
  
  var Tiers = map[string]TierLimits{
      "free": {
          Tier: "free",
          MaxBots: 1,
          MaxMessages: 100,
          MaxSources: 5,
          MaxStorageGB: 0.01, // 10 MB
          MaxAllowedOrigins: 1,
          IncludedCredits: 0,
      },
      "pro": {
          Tier: "pro",
          MaxBots: 5,
          MaxMessages: -1, // Unlimited (pay-as-you-go)
          MaxSources: 50,
          MaxStorageGB: 1,
          MaxAllowedOrigins: 10,
          IncludedCredits: 20.00,
      },
      "enterprise": {
          Tier: "enterprise",
          MaxBots: -1, // Unlimited
          MaxMessages: -1,
          MaxSources: -1,
          MaxStorageGB: -1,
          MaxAllowedOrigins: -1,
          IncludedCredits: 0, // Custom billing
      },
  }
  ```

#### 1.3 Create Usage Tracking Model
- [ ] Create `internal/models/usage_model.go`:
  ```go
  type UsageRecord struct {
      ID        string    `json:"id" gorm:"primaryKey"`
      UserID    string    `json:"user_id" gorm:"not null;index"`
      BotID     string    `json:"bot_id" gorm:"index"`
      Type      string    `json:"type"` // "chat_message", "embedding", "storage", "extra_bot"
      Quantity  float64   `json:"quantity"` // Amount used
      Cost      float64   `json:"cost"` // Cost in USD
      BilledAt  time.Time `json:"billed_at"`
      CreatedAt time.Time `json:"created_at"`
  }
  ```

---

### Step 2: Pricing Configuration

#### 2.1 Pricing Constants (DONE — `configs/pricing.go`)
- [x] All pricing constants, tier limits, billing rules, and cost helpers are defined in `configs/pricing.go`
- [x] This is the **single source of truth** — billing services import from here
- [x] Key values:
  - Pro subscription: $20/month (always charged, credits don't roll over)
  - Included credits: $20/month
  - Message: $0.001 | Embedding/1K: $0.0002 | Storage/GB: $0.10 | Extra bot: $5.00
  - 70% profit margin on all pay-as-you-go usage

---

### Step 3: Polar.sh Integration (Enhanced)

#### 3.1 Polar Configuration (DONE — `configs/config.go`)
- [x] Added to `configs/config.go`:
  ```go
  PolarAccessToken    string
  PolarWebhookSecret  string
  PolarOrganizationID string
  PolarProProductID   string
  PolarServerURL      string // defaults to https://api.polar.sh
  ```

#### 3.2 Polar Webhook Handler
- [ ] Update `internal/handlers/billing/webhook_handler.go`:
  - [ ] Handle `subscription.created`:
    - [ ] Update user tier to "pro"
    - [ ] Set `subscription_status` to "active"
    - [ ] Allocate initial credits ($20)
    - [ ] Set billing cycle dates
    - [ ] Send welcome email
  - [ ] Handle `subscription.updated`:
    - [ ] Update subscription status
    - [ ] Handle plan changes (upgrade/downgrade)
  - [ ] Handle `subscription.cancelled`:
    - [ ] Set `subscription_status` to "cancelled"
    - [ ] Gracefully downgrade at cycle end
    - [ ] Keep data but enforce free tier limits
  - [ ] Handle `subscription.renewed`:
    - [ ] Refresh credits to $20
    - [ ] Reset usage counter
    - [ ] Update billing cycle dates
  - [ ] Handle `payment.succeeded`:
    - [ ] Add pay-as-you-go credits to balance
    - [ ] Log transaction
  - [ ] Handle `payment.failed`:
    - [ ] Set status to "past_due"
    - [ ] Send email notification
    - [ ] Grace period: 7 days before downgrade

#### 3.3 Polar API Service
- [ ] Create `internal/services/billing/polar_service.go`:
  - [ ] `CreateCheckoutSession(userID, tier) (string, error)` - Generate checkout URL
  - [ ] `CreateCustomerPortalSession(userID) (string, error)` - Manage subscription link
  - [ ] `GetSubscription(subscriptionID) (*Subscription, error)` - Fetch current subscription
  - [ ] `CancelSubscription(subscriptionID) error` - Cancel at period end
  - [ ] `CreateUsageInvoice(userID, amount) error` - Charge for overage (Pro tier)

---

### Step 4: Usage Tracking & Metering

#### 4.1 Create Usage Service
- [x] Create `internal/services/billing/usage_service.go`:
  - [ ] `TrackChatMessage(userID, botID) error`
    - [ ] Check if user has credits or is free/enterprise
    - [ ] Deduct from credits if Pro
    - [ ] Record usage in `usage_records` table
    - [ ] If credits exhausted, create Polar invoice
  - [ ] `TrackEmbedding(userID, tokens) error`
    - [ ] Calculate cost
    - [ ] Deduct from credits if Pro
    - [ ] Record usage
  - [ ] `TrackStorage(userID, sizeGB) error`
    - [ ] Calculate monthly storage cost
    - [ ] Deduct from credits if Pro
  - [ ] `GetCurrentUsage(userID) (*UsageReport, error)`
    - [ ] Return current period usage breakdown

#### 4.2 Credit Management
- [x] Create `internal/services/billing/credits_service.go`:
  - [ ] `DeductCredits(userID, amount) error`
    - [ ] Check balance
    - [ ] Deduct if sufficient
    - [ ] If insufficient, create overage invoice via Polar
  - [ ] `AddCredits(userID, amount) error` - Add credits (payment received)
  - [ ] `RefreshMonthlyCredits(userID) error` - Reset credits on billing cycle
  - [ ] `GetCreditsBalance(userID) (float64, error)`

#### 4.3 Integrate Usage Tracking
- [x] Update `internal/services/chat/chat_service.go`:
  - [ ] Call `usageService.TrackChatMessage()` after each message
- [ ] Update `internal/services/embedding/embedding.go`:
  - [ ] Call `usageService.TrackEmbedding()` after embedding generation
- [ ] Update `internal/services/storage/minio.go`:
  - [ ] Call `usageService.TrackStorage()` after file upload

---

### Step 5: Entitlement Middleware (Enhanced)

#### 5.1 Update Entitlement Middleware
- [x] Update `internal/middleware/entitlement.go`:
  ```go
  func EnforceLimit(limitType string) gin.HandlerFunc {
      return func(c *gin.Context) {
          user := c.MustGet("user").(*models.User)
          tier := models.Tiers[user.Tier]
          
          switch limitType {
          case "bot_creation":
              if tier.MaxBots == -1 {
                  // Unlimited (Enterprise)
                  c.Next()
                  return
              }
              botCount := getBotCount(user.ID)
              if botCount >= tier.MaxBots {
                  c.JSON(403, gin.H{"error": "Bot limit reached. Upgrade to create more bots."})
                  c.Abort()
                  return
              }
          
          case "message_send":
              if tier.MaxMessages == -1 {
                  // Unlimited (Pro/Enterprise with credits or paid)
                  c.Next()
                  return
              }
              monthlyMessages := getMonthlyMessages(user.ID)
              if monthlyMessages >= tier.MaxMessages {
                  c.JSON(403, gin.H{"error": "Message limit reached. Upgrade to Pro for unlimited messages."})
                  c.Abort()
                  return
              }
          
          // Similar checks for sources, storage, origins...
          }
          
          c.Next()
      }
  }
  ```

#### 5.2 Apply Middleware to Routes
- [x] Bot creation: `EnforceLimit("bot_creation")`
- [ ] Message send: `EnforceLimit("message_send")`
- [ ] Source creation: `EnforceLimit("source_creation")`
- [ ] File upload: `EnforceLimit("storage")`

---

### Step 6: Enterprise Tier Management

#### 6.1 Manual Enterprise Activation
- [ ] Create admin endpoint (protected):
  - [ ] `POST /api/admin/users/:id/upgrade-enterprise`
  - [ ] Requires admin authentication
  - [ ] Sets user tier to "enterprise"
  - [ ] Sets all limits to -1
  - [ ] Sends confirmation email

#### 6.2 Enterprise Contact Flow
- [ ] Create contact form endpoint:
  - [ ] `POST /api/contact/enterprise`
  - [ ] Captures user info and requirements
  - [ ] Sends email to sales team
  - [ ] Auto-reply to user with next steps

---

### Step 7: Billing Cycle Management

#### 7.1 Create Billing Cron Job
- [x] Create `internal/worker/billing_worker.go`:
  - [ ] Run daily at midnight UTC
  - [ ] Find users with billing cycle end date = today
  - [ ] Refresh Pro tier credits ($20)
  - [ ] Reset usage counters
  - [ ] Update billing cycle dates (+1 month)
  - [ ] Send billing summary email

#### 7.2 Overage Billing
- [ ] At end of billing cycle:
  - [ ] Calculate total usage beyond credits
  - [ ] Create Polar invoice for overage
  - [ ] Charge via Polar payment method on file
  - [ ] Reset credits for next cycle

---

## Frontend Tasks

### Step 1: Pricing Page

#### 1.1 Create Pricing Page
- [ ] Create `ui/src/routes/pricing/index.tsx`:
  - [ ] Three-column pricing table (Free, Pro, Enterprise)
  - [ ] Feature comparison matrix
  - [ ] Clear CTA buttons:
    - [ ] Free: "Get Started" → Sign up
    - [ ] Pro: "Subscribe for $20/month" → Polar checkout
    - [ ] Enterprise: "Contact Sales" → Contact form modal
  - [ ] FAQ section
  - [ ] "Pay-as-you-go" pricing details for Pro tier

#### 1.2 Pricing Table Component
- [ ] Create `ui/src/components/pricing/pricing-table.tsx`:
  - [ ] Responsive design (mobile-friendly)
  - [ ] Highlight Pro tier as "Most Popular"
  - [ ] Show monthly/annual toggle (future: annual discount)
  - [ ] Feature icons and checkmarks
  - [ ] Tooltip explanations for complex features

---

### Step 2: Dashboard Usage & Billing UI

#### 2.1 Current Plan Widget
- [ ] Create `ui/src/components/dashboard/current-plan.tsx`:
  - [ ] Display current tier (Free/Pro/Enterprise)
  - [ ] Show usage progress bars:
    - [ ] Bots: 3/5 (Pro)
    - [ ] Messages this month: 450 (unlimited with credits)
    - [ ] Storage: 0.5 GB / 1 GB
    - [ ] Credits remaining: $12.45 / $20.00
  - [ ] "Upgrade" button if on Free
  - [ ] "Manage Subscription" link to Polar portal

#### 2.2 Usage Dashboard
- [ ] Create `ui/src/routes/dashboard/billing/usage.tsx`:
  - [ ] Current billing cycle dates
  - [ ] Credits balance (Pro tier)
  - [ ] Usage breakdown chart:
    - [ ] Messages sent
    - [ ] Embeddings generated
    - [ ] Storage used
    - [ ] Cost breakdown
  - [ ] Usage history table
  - [ ] Export to CSV button

#### 2.3 Subscription Management
- [ ] Create `ui/src/routes/dashboard/billing/subscription.tsx`:
  - [ ] Current plan details
  - [ ] Next billing date
  - [ ] Payment method (via Polar portal)
  - [ ] "Upgrade to Enterprise" CTA
  - [ ] "Cancel Subscription" button (with confirmation)
  - [ ] Link to Polar customer portal for payment updates

---

### Step 3: Upgrade Prompts & Paywalls

#### 3.1 Limit Reached Modals
- [ ] Create `ui/src/components/modals/upgrade-modal.tsx`:
  - [ ] Trigger when user hits limit
  - [ ] Show current usage vs limit
  - [ ] "Upgrade to Pro" CTA
  - [ ] Pricing comparison
  - [ ] "Maybe Later" button

#### 3.2 In-App Upgrade Prompts
- [ ] Bot creation page:
  - [ ] Show bot limit at top
  - [ ] Disable "Create Bot" button when limit reached
  - [ ] Show upgrade prompt
- [ ] Source creation:
  - [ ] Warn when approaching source limit
- [ ] File upload:
  - [ ] Show storage usage bar
  - [ ] Block upload if limit exceeded

---

### Step 4: API Integration

#### 4.1 Update API Client
- [ ] Update `ui/src/api/billing.ts`:
  ```ts
  export const getCurrentUsage = () => apiClient.get('/api/billing/usage');
  export const getSubscription = () => apiClient.get('/api/billing/subscription');
  export const createCheckoutSession = (tier: 'pro') => 
      apiClient.post('/api/billing/checkout', { tier });
  export const getCustomerPortalURL = () => 
      apiClient.get('/api/billing/portal');
  export const cancelSubscription = () => 
      apiClient.post('/api/billing/cancel');
  ```

#### 4.2 React Query Hooks
- [ ] Create `ui/src/hooks/use-billing.ts`:
  ```ts
  export const useCurrentUsage = () => useQuery(['usage'], getCurrentUsage);
  export const useSubscription = () => useQuery(['subscription'], getSubscription);
  export const useCheckout = () => useMutation(createCheckoutSession);
  ```

---

### Step 5: Enterprise Contact Form

#### 5.1 Contact Form Component
- [ ] Create `ui/src/components/forms/enterprise-contact-form.tsx`:
  - [ ] Fields: Name, Email, Company, Use Case, Expected Volume
  - [ ] Validation with TanStack Form
  - [ ] Submit to `/api/contact/enterprise`
  - [ ] Success message: "We'll be in touch within 24 hours"

#### 5.2 Add to Pricing Page
- [ ] Modal or dedicated page for enterprise contact
- [ ] Link from pricing table "Contact Sales" button

---

## Email Notifications

### Email Templates
- [ ] Welcome email (new Pro subscriber)
- [ ] Credits running low (Pro tier, <20% remaining)
- [ ] Credits exhausted (overage billing will apply)
- [ ] Monthly billing summary (usage + charges)
- [ ] Payment failed (grace period notice)
- [ ] Subscription cancelled (data retention policy)
- [ ] Enterprise inquiry confirmation

---

## Testing Tasks

### Unit Tests
- [ ] Test tier limit enforcement
- [ ] Test credit deduction logic
- [ ] Test usage tracking accuracy
- [ ] Test billing cycle refresh

### Integration Tests
- [ ] Test Polar webhook handling
- [ ] Test checkout flow (sandbox)
- [ ] Test overage billing
- [ ] Test subscription cancellation

### Load Tests
- [ ] Simulate high usage (credit exhaustion)
- [ ] Test concurrent usage tracking

---

## Success Metrics

- [ ] Successful Polar integration (webhooks working)
- [ ] Accurate usage metering (<1% error rate)
- [ ] Credit-based billing working correctly
- [ ] 70% profit margin maintained
- [ ] Upgrade conversion rate >3% (Free → Pro)
- [ ] <1% billing disputes
- [ ] Enterprise inquiries processed within 24h

---

## Documentation

- [ ] Update README with pricing tiers
- [ ] Document Polar setup process
- [ ] Add billing FAQ for users
- [ ] Create admin guide for enterprise activation

---

**Created**: February 11, 2025
