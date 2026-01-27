# Phase 3: Analytics & Scaling

## Goal
Provide insights into chatbot usage and ensure system stability.

## Backend Tasks

### Step 1: Analytics Service
- **Queries**: Aggregation over `messages` and `conversations`.
    - "Messages per Day": `SELECT date(created_at), count(*) ... GROUP BY 1`.
    - "Token Usage": Sum of `token_count`.

### Step 2: Rate Limiting
- Use Redis-based rate limiter (e.g., `go-redis/redis_rate`).
- Middleware: `RateLimitMiddleware`.
    - Free: 100 req/hour.
    - Pro: 1000 req/hour.

## Frontend Tasks

### Step 1: Analytics Dashboard
- Charts: Line chart (Activity), Pie chart (Sources).
- Library: `recharts` or `visx`.
