# Variables
PROTO_DIR := proto
API_DIR := api
BUF := $(shell go env GOPATH)/bin/buf
CMD_DIR := cmd/server

# Default target
.PHONY: all
all: generate

# Generate Go code from all proto files using buf
.PHONY: generate
generate: $(API_DIR)
	@echo "Generating Go code from proto files using buf..."
	$(BUF) generate
	@echo "Generation complete!"

# Create api directory if it doesn't exist
$(API_DIR):
	mkdir -p $(API_DIR)

# Clean generated files
.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	rm -f $(API_DIR)/*.pb.go
	rm -rf $(API_DIR)/*.pb.gw.go
	rm -f $(API_DIR)/*.swagger.json
	@echo "Clean complete!"

# Install buf if needed
.PHONY: install-tools
install-tools:
	@echo "Installing buf..."
	go install github.com/bufbuild/buf/cmd/buf@latest
	@echo "buf installed successfully!"

# Initialize buf dependencies
.PHONY: buf-deps
buf-deps:
	@echo "Updating buf dependencies..."
	$(BUF) dep update
	@echo "Dependencies updated!"

# Lint proto files
.PHONY: lint
lint:
	@echo "Linting proto files..."
	$(BUF) lint
	@echo "Linting complete!"

# Format proto files
.PHONY: format
format:
	@echo "Formatting proto files..."
	$(BUF) format -w
	@echo "Formatting complete!"

# Build and run targets
.PHONY: build
build:
	@echo "Building server..."
	go build -o bin/server ./$(CMD_DIR)
	@echo "Build complete!"

.PHONY: run
run:
	@echo "Running server..."
	go run ./$(CMD_DIR)

.PHONY: dev
dev: generate run

# Database targets
.PHONY: db-up
db-up:
	@echo "Starting database services..."
	docker-compose up -d dgraph-zero dgraph-alpha redis-master redis-replica redis-sentinel
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Database services started!"

.PHONY: db-down
db-down:
	@echo "Stopping database services..."
	docker-compose down
	@echo "Database services stopped!"

.PHONY: db-logs
db-logs:
	@echo "Showing database logs..."
	docker-compose logs -f dgraph-zero dgraph-alpha redis-master

.PHONY: db-status
db-status:
	@echo "Database service status:"
	docker-compose ps

.PHONY: db-clean
db-clean:
	@echo "Cleaning database volumes..."
	docker-compose down -v
	docker volume prune -f
	@echo "Database volumes cleaned!"

.PHONY: db-reset
db-reset: db-down db-clean db-up
	@echo "Database reset complete!"

# Development environment targets
.PHONY: dev-up
dev-up: db-up
	@echo "Starting full development environment..."
	docker-compose up -d
	@echo "Development environment ready!"
	@echo "Services available at:"
	@echo "  - Dgraph Alpha (HTTP): http://localhost:8080"
	@echo "  - Dgraph Alpha (gRPC): localhost:9080"
	@echo "  - Dgraph Ratel (UI): http://localhost:8000"
	@echo "  - Redis: localhost:6379"
	@echo "  - Redis Commander (UI): http://localhost:8081 (admin/admin)"

.PHONY: dev-down
dev-down:
	@echo "Stopping development environment..."
	docker-compose down
	@echo "Development environment stopped!"

# Testing targets
.PHONY: test
test:
	@echo "Running unit tests..."
	go test -v ./...

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./internal/database/

.PHONY: test-all
test-all: test test-integration

# Dependency management
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated!"

.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies updated!"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Development:"
	@echo "  dev-up        - Start full development environment (databases + UIs)"
	@echo "  dev-down      - Stop development environment"
	@echo "  generate      - Generate Go code from proto files"
	@echo "  build         - Build the server binary"
	@echo "  run           - Run the server directly"
	@echo "  dev           - Generate and run (development workflow)"
	@echo ""
	@echo "Database:"
	@echo "  db-up         - Start database services only"
	@echo "  db-down       - Stop database services"
	@echo "  db-logs       - Show database logs"
	@echo "  db-status     - Show database service status"
	@echo "  db-clean      - Clean database volumes"
	@echo "  db-reset      - Reset database (down + clean + up)"
	@echo ""
	@echo "Testing:"
	@echo "  test          - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-all      - Run all tests"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  deps-update   - Update all dependencies"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean         - Remove generated files"
	@echo "  lint          - Lint proto files"
	@echo "  format        - Format proto files"
	@echo "  install-tools - Install buf"
	@echo "  buf-deps      - Update buf dependencies"
	@echo "  help          - Show this help message"
