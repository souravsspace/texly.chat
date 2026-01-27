# Phase 1.5: Chat Interface & RAG

## Goal
Enable users to chat with their bots using RAG (Retrieval-Augmented Generation).

## Backend Tasks

### Step 1: Chat Handler
**File**: `internal/handlers/chat/handler.go`
- POST `/api/bots/:id/chat`
- Logic:
    1.  Receive user message.
    2.  Generate embedding for message.
    3.  **Search**: Query SQLite `vec_items` for relevant chunks.
    4.  **Context**: Retrieve text content for those chunks.
    5.  **Prompt**: "System Prompt + Context + User Message".
    6.  **LLM**: Call OpenAI Chat Completion (Stream=true).

### Step 2: Streaming (SSE)
- Use Server-Sent Events (SSE) to stream tokens back to the client.
- `c.SSEvent("message", content)` in Gin.

## Frontend Tasks

### Step 1: Chat UI
**File**: `ui/src/routes/bots/$botId/chat.tsx`
- Layout: Message list (scrollable), Input area.
- State: `messages` array.
- Streaming: Use `fetch` with `ReadableStream` to process SSE chunks and update UI in real-time.

### Step 2: Markdown Rendering
- Use `react-markdown` to render bot responses (bold, lists, code blocks).
