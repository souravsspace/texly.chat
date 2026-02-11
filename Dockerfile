# Multi-stage Dockerfile for Texly.Chat - Optimized for minimal size
# Stage 1: Build the UI using Bun
FROM oven/bun:1.2-alpine AS ui-builder

# Set Node.js memory limit to prevent OOM errors during build
# Install Node.js (v22+) for build stability and Vite 7 support
RUN apk add --no-cache --repository http://dl-cdn.alpinelinux.org/alpine/v3.21/main nodejs

# Set Node.js memory limit (increased to 4GB for heavy builds)
ENV NODE_OPTIONS="--max-old-space-size=4096"

WORKDIR /app/ui

# Copy UI package files
COPY ui/package.json ui/bun.lock* ./

# Install dependencies
# Install dependencies with cache
RUN --mount=type=cache,target=/root/.bun/install/cache bun install --frozen-lockfile

# Copy UI source code
COPY ui/ ./

# Build the UI using Node.js directly for better memory management
RUN node node_modules/vite/bin/vite.js build && \
    rm -rf dist && \
    mv .output/public dist

# Stage 2: Build the Go application with embedded UI
FROM golang:1.25-bookworm AS go-builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    git \
    && rm -rf /var/lib/apt/lists/*


# Enable CGO (no special flags needed for Debian/Glibc)
ENV CGO_CFLAGS=""

WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./

# Download Go dependencies
# Download Go dependencies with cache
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy the entire Go source code
COPY . .

# Copy built UI from previous stage
COPY --from=ui-builder /app/ui/dist ./ui/dist

# Build the Go binary with CGO enabled (dynamic link for Debian)
# Strip symbols and debug info for smaller binary
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags 'osusergo netgo' \
    -o texly.chat \
    ./cmd/app && \
    strip texly.chat

# Stage 3: Create minimal runtime image
FROM debian:bookworm-slim

# Install only essential runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN useradd -m -U appuser && \
    chown -R appuser:appuser /app

# Copy only the binary from builder
COPY --from=go-builder --chown=appuser:appuser /app/texly.chat /app/texly.chat

# Switch to non-root user for security
USER appuser

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/texly.chat"]
