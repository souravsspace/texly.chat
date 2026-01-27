# Task: Widget Dashboard Configuration

## Goal
Allow bot owners to customize how their public widget looks and behaves, and provide them with the code snippet to embed it.

## UI Implementation
**File**: `ui/src/routes/bots/$botId/widget.tsx`

### Step 1: Configuration Form
- **Fields**:
    - `Theme Color` (Color picker).
    - `Initial Message` (Text input).
    - `Allowed Origins` (Text area, one domain per line).
    - `Position` (Dropdown: `bottom-right`, `bottom-left`).
- **Save**: PATCH `/api/bots/:id` (or a dedicated config endpoint if we separated it).

### Step 2: Live Preview
- Render a mock version of the widget on the right side of the settings page.
- Update immediately as the user changes form values.

### Step 3: Embed Code Generator
- **Display**: Use a code block component.
- **Template**:
    ```html
    <script src="https://api.texly.chat/widget.js" data-bot-id="CURRENT_BOT_ID"></script>
    ```
- **Copy Button**: One-click copy to clipboard.
