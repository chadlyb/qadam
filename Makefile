# QADAM Tools Makefile

# Variables
VERSION ?= dev
BINARY_DIR = bin
EXTRACT_BINARY = extract
BUILD_BINARY = build

# Go build flags
LDFLAGS = -ldflags="-s -w -X main.version=$(VERSION)"

# Default target
.PHONY: all
all: clean build

# Build all binaries for current platform
.PHONY: build
build: $(BINARY_DIR)
	go build $(LDFLAGS) -o $(BINARY_DIR)/$(EXTRACT_BINARY) ./cmd/extract
	go build $(LDFLAGS) -o $(BINARY_DIR)/$(BUILD_BINARY) ./cmd/build

# Build for all platforms
.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(EXTRACT_BINARY)-linux-amd64 ./cmd/extract
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(BUILD_BINARY)-linux-amd64 ./cmd/build
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(EXTRACT_BINARY)-windows-amd64.exe ./cmd/extract
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(BUILD_BINARY)-windows-amd64.exe ./cmd/build
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(EXTRACT_BINARY)-darwin-amd64 ./cmd/extract
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(BUILD_BINARY)-darwin-amd64 ./cmd/build

# Create binary directory
$(BINARY_DIR):
	mkdir -p $(BINARY_DIR)

# Run all tests
.PHONY: test
test:
	@echo "Running all tests..."
	go test -v ./...

# Run tests with verbose output
.PHONY: test-verbose
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -count=1 ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run language model tests specifically
.PHONY: test-language-model
test-language-model:
	@echo "Testing language model..."
	go test -v ./shared -run "TestLanguageModel"

# Run string extraction tests specifically
.PHONY: test-extraction
test-extraction:
	@echo "Testing string extraction..."
	go test -v ./cmd/extract -run "TestQGetStrings"

# Run benchmarks
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. ./shared

# Run tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	go test -race -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Create release packages
.PHONY: release
release: build-all
	@echo "Creating release packages..."
	cd $(BINARY_DIR) && tar -czf qadam-$(VERSION)-linux-amd64.tar.gz $(EXTRACT_BINARY)-linux-amd64 $(BUILD_BINARY)-linux-amd64 README.md
	cd $(BINARY_DIR) && tar -czf qadam-$(VERSION)-darwin-amd64.tar.gz $(EXTRACT_BINARY)-darwin-amd64 $(BUILD_BINARY)-darwin-amd64 README.md
	cd $(BINARY_DIR) && zip qadam-$(VERSION)-windows-amd64.zip $(EXTRACT_BINARY)-windows-amd64.exe $(BUILD_BINARY)-windows-amd64.exe README.md

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build              - Build binaries for current platform"
	@echo "  build-all          - Build binaries for all platforms (Linux, Windows, macOS)"
	@echo "  test               - Run all tests"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  test-language-model- Test language model specifically"
	@echo "  test-extraction    - Test string extraction specifically"
	@echo "  benchmark          - Run benchmarks"
	@echo "  test-race          - Run tests with race detection"
	@echo "  clean              - Clean build artifacts"
	@echo "  deps               - Download and tidy dependencies"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linter"
	@echo "  release            - Create release packages"
	@echo "  help               - Show this help" 