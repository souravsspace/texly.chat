# texly.chat

A powerful SaaS chatbot platform with built-in Retrieval-Augmented Generation (RAG) capabilities, allowing users to create AI bots trained on their own data (URLs, Files).

## Features

- 🤖 **Custom Chatbots**: Create and manage multiple AI chatbots.
- 📚 **RAG Support**: Automatically indexes content for context-aware answers.
- 🕷️ **Web Scraping**: Crawl and index websites as knowledge sources.
- 📄 **File Support**: Upload and process PDF and Excel documents.
- 🔍 **Vector Search**: Built-in vector search using `sqlite-vec` and OpenAI embeddings.
- ⚡ **Real-time Chat**: Streaming responses via Server-Sent Events (SSE).
- 🐳 **Containerized**: Fully Dockerized for easy deployment.

## Tech Stack

### Backend
- **Language**: Go 1.25+
- **Framework**: Gin Web Framework
- **Database**: SQLite (with WAL mode & `sqlite-vec` extension)
- **ORM**: GORM
- **AI**: OpenAI API (Embeddings & Chat Completions)
- **Storage**: MinIO (S3 Compatible) for file storage

### Frontend
- **Framework**: React 19 (Vite)
- **Routing**: TanStack Router
- **State Management**: TanStack Query & Zustand
- **Styling**: TailwindCSS v4
- **Components**: shadcn/ui

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js / Bun
- Docker & Docker Compose
- OpenAI API Key

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/souravsspace/texly.chat.git
   cd texly.chat
   ```

2. **Environment Setup**
   ```bash
   cp .env.example .env.local
   # Edit .env.local and add your OPENAI_API_KEY
   ```

3. **Run with Docker (Recommended)**
   To start the full stack (App + MinIO):
   ```bash
   make docker-up
   ```
   The app will be available at http://localhost:8080.

### Local Development

To run the backend and frontend separately for development:

1. **Install Dependencies**
   ```bash
   make install
   ```

2. **Run Development Servers**
   ```bash
   make dev
   ```
   - Frontend: http://localhost:3000
   - Backend: http://localhost:8080

## Build

To build the production binary with embedded UI:

```bash
make build
```
The binary will be output to `dist/texly.chat`.

## Architecture

- **`cmd/`**: Application entry points.
- **`internal/`**: Core application logic (Handlers, Services, Repositories).
- **`ui/`**: React frontend application.
- **`data/`**: SQLite database storage.
