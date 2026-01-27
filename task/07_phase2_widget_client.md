# Task: Widget Client (React App)

## Goal
Build a lightweight, standalone React application that can be embedded on any website. It must interact with the Public API and be isolated from the host site's CSS.

## Project Setup
- [ ] **Initialize Widget App** `ui/widget`
    - Create a new directory or use Vite's multi-page/library mode.
    - **Crucial**: Configure build to output a single JS file (e.g., `widget.js`) if possible, or a minimal set of files.
    - Use `ShadowRoot` (Shadow DOM) for style isolation.

## Components
**File**: `ui/widget/src/App.tsx` (conceptual path)

### Step 1: Launcher
- **Floating Button**: Positioned largely by configuration (e.g., bottom-right).
- **Animation**: Smooth transition to open state.

### Step 2: Chat Interface
- **Chat Window**:
    - Header (Bot Name, Avatar, Close Button).
    - Message List (User & AI messages).
    - Input Area (Textarea, Send button).
- **State Management**:
    - `messages`: Array of `{role, content}`.
    - `isOpen`: Boolean.
    - `sessionId`: Store in `sessionStorage` or `localStorage` to persist across page reloads.

### Step 3: Embed Script
**File**: `ui/public/embed.js` (or similar)
- Create a script that:
    1.  Reads `data-bot-id` from its own script tag.
    2.  Creates a `div` container.
    3.  Attaches Shadow DOM.
    4.  Injects the Widget App styles and script into the Shadow DOM.

## API Client
- Implement `ui/widget/src/api.ts` to call `/api/public` endpoints.
