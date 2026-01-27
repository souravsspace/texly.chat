# Phase 1.2: Bot Management CRUD

## Goal
Allow users to create, list, update, and delete their chatbots.

## Backend Tasks

### Step 1: Bot Handler
**File**: `internal/handlers/bot/handler.go`
- Implement standard CRUD methods:
    - `CreateBot`: POST /api/bots
    - `ListBots`: GET /api/bots
    - `GetBot`: GET /api/bots/:id
    - `UpdateBot`: PUT /api/bots/:id
    - `DeleteBot`: DELETE /api/bots/:id (Soft delete)

### Step 2: Validation
- Use `go-playground/validator` tags on DTO structs.
- Validate `Model` selection (gpt-4, gpt-3.5, etc.).

## Frontend Tasks

### Step 1: API Client
**File**: `ui/src/api/bots.ts`
- Implement fetch wrapper for Bot endpoints.

### Step 2: Dashboard UI
**File**: `ui/src/routes/dashboard.tsx`
- Grid layout of `BotCard` components.
- `CreateBotDialog` with form (Name, Description).
- Use `TanStack Query` for data fetching (`useQuery(['bots'])`).
