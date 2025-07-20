# Variables
PROTO_DIR := proto
PB_DIR := pb
BUF := $(shell go env GOPATH)/bin/buf

# Default target
.PHONY: all
all: generate

# Generate Go code from all proto files using buf
.PHONY: generate
generate: $(PB_DIR)
	@echo "Generating Go code from proto files using buf..."
	$(BUF) generate
	@echo "Generation complete!"

# Create pb directory if it doesn't exist
$(PB_DIR):
	mkdir -p $(PB_DIR)

# Clean generated files
.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	rm -f $(PB_DIR)/*.pb.go
	rm -rf $(PB_DIR)/*.pb.gw.go
	rm -f $(PB_DIR)/*.swagger.json
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

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  generate      - Generate Go code from all proto files using buf"
	@echo "  clean         - Remove generated .pb.go files"
	@echo "  install-tools - Install buf"
	@echo "  buf-deps      - Update buf dependencies"
	@echo "  lint          - Lint proto files using buf"
	@echo "  format        - Format proto files using buf"
	@echo "  help          - Show this help message"