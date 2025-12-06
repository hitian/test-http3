# HTTP/3 Checker Makefile

# Application name
APP_NAME := http3check

# Version (can be overridden with make VERSION=x.x.x)
VERSION ?= 1.0.0

# Build directories
BUILD_DIR := dist
BIN_DIR := bin

# Go build flags
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"
BUILD_FLAGS := $(LDFLAGS) -trimpath

# Supported platforms
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	@echo "Building $(APP_NAME) for current platform..."
	@mkdir -p $(BIN_DIR)
	go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME) .
	@echo "Build complete: $(BIN_DIR)/$(APP_NAME)"

# Build for all platforms
.PHONY: build-all
build-all: clean
	@echo "Building $(APP_NAME) v$(VERSION) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_NAME=$(APP_NAME)-$$OS-$$ARCH; \
		if [ $$OS = "windows" ]; then OUTPUT_NAME=$$OUTPUT_NAME.exe; fi; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$$OUTPUT_NAME .; \
		if [ $$? -eq 0 ]; then \
			echo "✅ Built $(BUILD_DIR)/$$OUTPUT_NAME"; \
		else \
			echo "❌ Failed to build for $$OS/$$ARCH"; \
			exit 1; \
		fi; \
	done
	@echo "All builds complete!"

# Build for Linux platforms only
.PHONY: build-linux
build-linux: clean
	@echo "Building $(APP_NAME) for Linux platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .
	@echo "Linux builds complete!"

# Build for macOS platforms only
.PHONY: build-macos
build-macos: clean
	@echo "Building $(APP_NAME) for macOS platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .
	@echo "macOS builds complete!"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(BIN_DIR)
	@rm -f $(APP_NAME)
	@rm -f $(APP_NAME)-*
	@echo "Clean complete!"

# Show build information
.PHONY: info
info:
	@echo "Application: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Dir: $(BUILD_DIR)"
	@echo "Binary Dir: $(BIN_DIR)"
	@echo "Platforms: $(PLATFORMS)"
	@echo ""
	@echo "Available targets:"
	@echo "  build       - Build for current platform"
	@echo "  build-all   - Build for all platforms"
	@echo "  build-linux - Build for Linux platforms only"
	@echo "  build-macos - Build for macOS platforms only"
	@echo "  deps        - Install dependencies"
	@echo "  test        - Run tests"
	@echo "  test-race   - Run tests with race detection"
	@echo "  clean       - Clean build artifacts"
	@echo "  info        - Show this information"

# Install to local bin (requires GOPATH/bin or GOBIN in PATH)
.PHONY: install
install:
	@echo "Installing $(APP_NAME)..."
	go install $(BUILD_FLAGS) .
	@echo "Installed $(APP_NAME) to $$(go env GOPATH)/bin/"

# Create release archives
.PHONY: release
release: build-all
	@echo "Creating release archives..."
	@cd $(BUILD_DIR) && \
	for binary in $(APP_NAME)-*; do \
		if [ -f "$$binary" ]; then \
			echo "Creating archive for $$binary..."; \
			tar -czf "$$binary.tar.gz" "$$binary"; \
			echo "✅ Created $$binary.tar.gz"; \
		fi; \
	done
	@echo "Release archives complete!"

# Help target
.PHONY: help
help: info