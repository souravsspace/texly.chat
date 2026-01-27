# Phase 4: Monetization (Polar.sh)

## Goal
Implement paid subscriptions using Polar.sh.

## Backend Tasks

### Step 1: Polar Integration
- **Lib**: `github.com/polarsource/polar-go`.
- **Webhook**: Listen for `subscription.created`, `subscription.updated`.
- **Database**: Add `plan_id` / `subscription_status` to `User` model.

### Step 2: Entitlement Check
- Middleware to enforce limits based on Plan.
    - Active Bots limit.
    - Messages/month limit.

## Frontend Tasks

### Step 1: Upgrade UI
- Pricing Table.
- Buttons link to Polar.sh Checkout URL.
- "Manage Subscription" link to Polar Customer Portal.
