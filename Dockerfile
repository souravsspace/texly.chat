# Multi-stage Dockerfile for Texly.Chat - Optimized for minimal size
# Stage 1: Build the UI using Bun
FROM oven/bun:1.2-alpine AS ui-builder

# Set Node.js memory limit to prevent OOM errors during build
ENV NODE_OPTIONS="--max-old-space-size=4096"

WORKDIR /app/ui

# Copy UI package files
COPY ui/package.json ui/bun.lock* ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy UI source code
COPY ui/ ./

# Build the UI (generates .output/public which is moved to dist)
RUN bun run build && \
    rm -rf dist && \
    mv .output/public dist

# Stage 2: Build the Go application with embedded UI
FROM golang:1.25-alpine AS go-builder

# Install build dependencies (required for CGO and sqlite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the entire Go source code
COPY . .

# Copy built UI from previous stage
COPY --from=ui-builder /app/ui/dist ./ui/dist

# Build the Go binary with CGO enabled (required for sqlite-vec)
# Strip symbols and debug info for smaller binary
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static-pie'" \
    -tags 'osusergo netgo static_build' \
    -o texly.chat \
    ./cmd/app && \
    strip texly.chat

# Stage 3: Create minimal runtime image
FROM alpine:3.21

# Install only essential runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite-libs \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Create data directory for SQLite database
RUN mkdir -p /app/data && \
    adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

# Copy only the binary from builder
COPY --from=go-builder --chown=appuser:appuser /app/texly.chat /app/texly.chat

# Switch to non-root user for security
USER appuser

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/texly.chat"]
