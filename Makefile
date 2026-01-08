.PHONY: all build test lint clean install release help

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

BINARY := factory
MAIN := ./cmd/factory

all: build

## build: Build the binary
build:
	@echo "Building $(BINARY)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(MAIN)

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

## test-integration: Run integration tests
test-integration: build
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/integration/...

## lint: Run linters
lint:
	@echo "Running linters..."
	go vet ./...
	@which staticcheck > /dev/null && staticcheck ./... || echo "staticcheck not installed"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
	rm -f coverage.out

## install: Install binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY)..."
	cp $(BINARY) $(GOPATH)/bin/

## install-local: Install to /usr/local/bin (requires sudo)
install-local: build
	@echo "Installing $(BINARY) to /usr/local/bin..."
	sudo cp $(BINARY) /usr/local/bin/

## release: Create a release using goreleaser
release:
	@echo "Creating release..."
	goreleaser release --clean

## snapshot: Create a snapshot release (no publish)
snapshot:
	@echo "Creating snapshot..."
	goreleaser release --snapshot --clean

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':'

default: help
