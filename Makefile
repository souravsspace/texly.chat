.PHONY: dev build clean install dev-ui dev-api build-ui build-api

# Development (run both servers)
dev:
	@echo "Starting development servers..."
	@make -j2 dev-ui dev-api

dev-ui:
	@cd ui && bun dev

dev-api:
	@go run ./cmd/app

# Build production binary with embedded UI
build: build-ui build-api
	@echo "✓ Build complete! Binary at: dist/texly.chat"

build-ui:
	@echo "Building frontend..."
	@cd ui && bun install && bun run build
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
	@echo "✓ Dependencies installed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf dist/ ui/dist/ ui/.output/
	@echo "✓ Cleaned"

# Run tests
test:
	@go test -v ./...

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
	@go run cmd/ui-types/main.go > ui/src/api.types.d.ts
	@echo "✓ Types generated at ui/src/api.types.d.ts"
