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

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  generate      - Generate Go code from all proto files using buf"
	@echo "  build         - Build the server binary"
	@echo "  run           - Run the server"
	@echo "  dev           - Generate and run (development workflow)"
	@echo "  clean         - Remove generated .pb.go files"
	@echo "  install-tools - Install buf"
	@echo "  buf-deps      - Update buf dependencies"
	@echo "  lint          - Lint proto files using buf"
	@echo "  format        - Format proto files using buf"
	@echo "  help          - Show this help message"