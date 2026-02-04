# Texly Chat Widget

A lightweight, embeddable React chat widget for Texly bots.

## Features

- 🚀 **Lightweight**: Minimal bundle size (~200KB)
- 🎨 **Customizable**: Theme colors, position, and initial message
- 🔒 **Isolated**: Uses Shadow DOM for CSS isolation
- 💬 **Real-time**: SSE streaming for instant responses
- 📱 **Responsive**: Works on desktop and mobile
- 🔄 **Session Persistence**: Maintains conversation across page reloads

## Build

```bash
bun install
bun run build
```

The compiled widget will be output to `dist/texly-widget.js`.

## Usage

Add the following script tag to any website:

```html
<script 
  src="https://your-domain.com/widget/texly-widget.js" 
  data-bot-id="your-bot-id">
</script>
```

## Configuration

The widget automatically fetches its configuration from the bot settings, including:

- Theme color
- Initial greeting message
- Position (bottom-right, bottom-left, etc.)
- Bot avatar

## Development

```bash
bun run dev
```

## Testing

Open `test.html` in a browser (after starting the main server) to test the widget locally.

## Architecture

- **index.tsx**: Entry point with Shadow DOM setup
- **App.tsx**: Main app logic and state management
- **components/**: React components (Launcher, ChatWindow)
- **api/**: API client for public endpoints
- **utils/**: Helper utilities (session management)
- **types.ts**: TypeScript interfaces

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)
- Mobile browsers
