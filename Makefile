# Code-Factory Makefile

.PHONY: help build run test clean install lint fmt

# Variables
BINARY_NAME=factory
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Default target
help:
	@echo "Code-Factory Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  make build        - Build binary for current platform"
	@echo "  make build-all    - Build for all platforms"
	@echo "  make run          - Build and run"
	@echo "  make test         - Run tests"
	@echo "  make test-verbose - Run tests with verbose output"
	@echo "  make lint         - Run linters"
	@echo "  make fmt          - Format code"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make install      - Install to /usr/local/bin"
	@echo ""

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/factory
	@echo "✓ Build complete: $(BINARY_NAME)"

# Build for all platforms
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p dist
	
	@echo "Building for macOS (Intel)..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/factory
	
	@echo "Building for macOS (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/factory
	
	@echo "Building for Linux (x86_64)..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/factory
	
	@echo "Building for Linux (ARM64)..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/factory
	
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/factory
	
	@echo "✓ All builds complete"
	@ls -lh dist/

# Run the application
run: build
	./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

test-verbose:
	@echo "Running tests (verbose)..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Lint code
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Install to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed successfully"
	@echo "Run '$(BINARY_NAME) --help' to get started"

# Generate checksums for release
checksums:
	@echo "Generating checksums..."
	@cd dist && shasum -a 256 * > SHA256SUMS
	@echo "✓ Checksums generated: dist/SHA256SUMS"

# Development: watch and rebuild on changes
watch:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air
