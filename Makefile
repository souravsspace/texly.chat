.PHONY: dev build clean install dev-ui dev-api build-ui build-widget build-api docker-up docker-down docker-logs docker-build docker-clean

# Development (run both servers)
dev:
	@echo "Starting development servers..."
	@make -j2 dev-ui dev-api

dev-ui:
	@cd ui && bun dev

dev-api:
	@if command -v air > /dev/null; then \
		ENVIRONMENT=development air; \
	else \
		echo "Air not found, using go run... (Run 'go install github.com/air-verse/air@latest' for hot reload)"; \
		ENVIRONMENT=development go run ./cmd/app; \
	fi

# Build production binary with embedded UI
build: build-widget build-ui build-api
	@echo "✓ Build complete! Binary at: dist/texly.chat"

build-widget:
	@echo "Building widget..."
	@cd widget && bun install && bun run build

build-ui:
	@echo "Building frontend..."
	@cd ui && NODE_OPTIONS="--max-old-space-size=4096" bun install && NODE_OPTIONS="--max-old-space-size=4096" bun run build
	@rm -rf ui/dist
	@mv ui/.output/public ui/dist

build-api:
	@echo "Building backend..."
	@CGO_ENABLED=1 go build -ldflags="-s -w" -o dist/texly.chat ./cmd/app

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@cd ui && bun install
	@cd widget && bun install
	@echo "✓ Dependencies installed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf dist/ ui/dist/ ui/.output/ widget/dist/
	@echo "✓ Cleaned"

# Run tests
test:
	@go test -v -p 1 ./...

# Run tests with coverage
test-coverage:
	@mkdir -p coverage
	@go test -v -coverprofile=coverage/coverage-$(shell date +%Y-%m-%d).out ./...
	@go tool cover -func=coverage/coverage-$(shell date +%Y-%m-%d).out


# Format code
fmt:
	@go fmt ./...
	@cd ui && bun fix

# Generate TypeScript types
ui-types:
	@echo "Generating Typescript types..."
	@go run cmd/ui-types/main.go > ui/src/api/index.types.ts
	@echo "✓ Types generated at ui/src/api/index.types.ts"

# Build and start full stack (MinIO + App)
docker-up:
	@echo "Starting full stack with Docker..."
	@docker-compose --env-file .env.prd up -d
	@echo "✓ Full stack running"
	@echo "App: http://localhost:8080"
	@echo "MinIO Console: http://localhost:9001 (minioadmin/minioadmin)"

# Stop all Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose --env-file .env.prd down
	@echo "✓ Docker containers stopped"

# View Docker logs
docker-logs:
	@docker-compose --env-file .env.prd logs -f

# Build Docker image only
docker-build:
	@echo "Building Docker image..."
	@docker-compose --env-file .env.prd build server
	@echo "✓ Docker image built successfully"

# Clean everything (containers, volumes, images)
docker-clean:
	@echo "Cleaning Docker containers, images, and volumes..."
	@docker-compose --env-file .env.prd down -v --rmi local
	@echo "✓ Docker cleaned"

# Redis helpers
redis-cli:
	@docker exec -it texly-redis redis-cli

redis-flush:
	@docker exec -it texly-redis redis-cli FLUSHALL
	@echo "✓ Redis cache flushed"

redis-monitor:
	@docker exec -it texly-redis redis-cli MONITOR

cache-stats:
	@curl -s http://localhost:8080/health/redis | python3 -m json.tool

