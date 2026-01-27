# Technology Stack & Key Decisions

## Backend (Go)

We chose **Go** for its performance, simplicity, and strong concurrency primitives, which are essential for handling real-time chat and background scraping jobs.

-   **Language**: Go 1.25.2+
-   **Web Framework**: `github.com/gin-gonic/gin`
    -   *Why*: High performance, simple middleware ecosystem.
-   **Database ORM**: `gorm.io/gorm` + `gorm.io/driver/sqlite`
    -   *Why*: Developer productivity for standard CRUD.
-   **Vector Search**: `sqlite-vec` (sqlite extension)
    -   *Why*: High-performance vector search within SQLite. No need for Pinecone/Weaviate.
-   **Background Jobs**: Go Channels (MVP) / SQLite-backed Queue
    -   *Why*: Keeps infrastructure simple (Single Binary). No Redis dependency required for initial scale.
-   **Scraping**: `github.com/gocolly/colly/v2`
    -   *Why*: Mature, fast scraping framework.
-   **Payments**: `github.com/polarsource/polar-go`
    -   *Why*: Handles Merchant of Record complexity.

## Frontend (TypeScript / React)

-   **Framework**: **TanStack Start** (React 19)
-   **Styling**: **Tailwind CSS v4**
-   **Component Library**: **Shadcn UI**
-   **State Management**: **TanStack Query** (Server State) + **Zustand** (Client State)

## AI & Data Pipeline

-   **LLM Model**: `gpt-4o-mini` (General chat) / `gpt-4o` (Complex reasoning)
-   **Embeddings**: `text-embedding-3-small`
-   **Text Splitting**: Recursive Character Splitter (~500 tokens)

## Infrastructure & DevOps

-   **Database**: SQLite (Single file) with **WAL (Write-Ahead-Logging)**.
    -   *Strategy*: Single-node deployment.
-   **Queue**: In-Process (Channels).
-   **Deployment**: Docker Compose.

## Key Architectural Decisions

### 1. SQLite over PostgreSQL
For a self-hostable or easy-to-manage SaaS, SQLite removes the complexity of managing a database server. With modern NVMe SSDs, WAL mode, and `sqlite-vec`, it is a powerful stack for thousands of users.

### 2. "Widget" Architecture
The chat widget is a separate entry point (`src/entry-widget.tsx`) rendered inside a **Shadow DOM** to isolate styles from the host website.

### 3. Asynchronous Ingestion
We use internal Go channels or a simple database table to handle scraping jobs asynchronously, ensuring the API request `POST /api/sources` responds immediately while processing happens in the background.
