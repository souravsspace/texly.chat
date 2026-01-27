# Task: Widget Backend & Public API

## Goal
Implement the backend infrastructure required to support the public embeddable widget. This includes public-facing secure endpoints and database schema updates.

## Database Changes
- [ ] **Update Bot Model** `internal/models/bot.go`
    - Add `AllowedOrigins` (string/json) to store whitelisted domains.
    - Add `WidgetConfig` (embedded struct or json) containing:
        - `ThemeColor` (string)
        - `InitialMessage` (string)
        - `Position` (string)
    - Run `go run cmd/migrate/main.go`.

## API Endpoints
**File**: `internal/handlers/chat/public.go` (and routes in `internal/server/routes.go`)

### Step 1: Public Configuration
- **GET** `/api/public/bots/:id/config`
- **Logic**:
    1.  Validate `BotID`.
    2.  Check `Origin` header against `Bot.AllowedOrigins`.
    3.  Return config JSON (Name, Avatar, Theme, etc.).

### Step 2: Chat Handling (Public)
- **POST** `/api/public/chats`
    - Start a new anonymous session.
    - Return a session token/ID.
- **POST** `/api/public/chats/:session_id/messages`
    - Logic similar to authenticated chat but uses session config.
    - **Streaming**: Re-use the SSE logic from Phase 1.

### Step 3: Security Middleware
- Implement `PublicCORSMiddleware` to dynamically check `Origin` header against the Bot's allowed list.
