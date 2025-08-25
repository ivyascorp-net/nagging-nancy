# Nagging Nancy - Makefile
# Common commands for building, testing, and managing the project

# Build variables
BINARY_NAME := nancy
MAIN_PACKAGE := ./cmd/nancy
BUILD_DIR := ./build
DIST_DIR := ./dist

# Version information (can be overridden)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Go build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT) -w -s"
BUILD_FLAGS := -trimpath

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Nagging Nancy - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development commands
.PHONY: build
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: install
install: build ## Install the binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install $(BUILD_FLAGS) $(LDFLAGS) $(MAIN_PACKAGE)
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

.PHONY: run
run: ## Run the application
	go run $(MAIN_PACKAGE)

.PHONY: run-daemon
run-daemon: ## Run the daemon in foreground mode for debugging
	go run $(MAIN_PACKAGE) daemon start --foreground

.PHONY: start-daemon
start-daemon: ## Start the daemon in background mode
	go run $(MAIN_PACKAGE) daemon start

.PHONY: dev
dev: build ## Build and run in development mode
	./$(BUILD_DIR)/$(BINARY_NAME)

# Testing commands
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@mkdir -p $(BUILD_DIR)
	go test -v -race -coverprofile=$(BUILD_DIR)/coverage.out ./...
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report: $(BUILD_DIR)/coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -v -race ./...

.PHONY: bench
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Code quality commands
.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

.PHONY: check
check: fmt vet lint test ## Run all code quality checks

# Build and release commands
.PHONY: build-all
build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	
	# Windows
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PACKAGE)
	
	@echo "Built all platforms in $(DIST_DIR)/"

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)

.PHONY: build-macos
build-macos: ## Build for macOS
	@echo "Building for macOS..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "Building for Windows..."
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

.PHONY: release
release: check build-all ## Create a release (runs checks and builds all platforms)
	@echo "Release $(VERSION) built successfully!"
	@echo "Binaries available in $(DIST_DIR)/"

# Utility commands
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	go clean

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

.PHONY: update-deps
update-deps: ## Update all dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

.PHONY: vendor
vendor: ## Create vendor directory
	@echo "Creating vendor directory..."
	go mod vendor

# Development utilities
.PHONY: debug
debug: ## Build with debug symbols
	@echo "Building with debug symbols..."
	@mkdir -p $(BUILD_DIR)
	go build -gcflags="all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME)-debug $(MAIN_PACKAGE)
	@echo "Debug binary: $(BUILD_DIR)/$(BINARY_NAME)-debug"

.PHONY: profile
profile: ## Build with profiling enabled
	@echo "Building with profiling..."
	@mkdir -p $(BUILD_DIR)
	go build -tags=profile -o $(BUILD_DIR)/$(BINARY_NAME)-profile $(MAIN_PACKAGE)

.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"

# Docker commands (if needed in future)
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t nancy:$(VERSION) .
	docker tag nancy:$(VERSION) nancy:latest

# Git helpers
.PHONY: tag
tag: ## Create and push a git tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make tag VERSION=v1.0.0"; exit 1; fi
	@echo "Creating tag $(VERSION)..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

# Test and development helpers
.PHONY: test-notifications
test-notifications: build ## Test notification system
	@echo "Testing notifications..."
	./$(BUILD_DIR)/$(BINARY_NAME) test notification

.PHONY: install-notifications
install-notifications: ## Install platform-specific notification dependencies
	@echo "Installing notification dependencies..."
	./scripts/install-notifications.sh

.PHONY: add-test-reminder
add-test-reminder: build ## Add a test reminder
	@echo "Adding test reminder..."
	./$(BUILD_DIR)/$(BINARY_NAME) add "Test reminder" --time "$(shell date -d '+1 minute' '+%H:%M')" --priority high

.PHONY: daemon-start
daemon-start: build ## Start daemon in background
	./$(BUILD_DIR)/$(BINARY_NAME) daemon start

.PHONY: daemon-stop
daemon-stop: build ## Stop daemon
	./$(BUILD_DIR)/$(BINARY_NAME) daemon stop

.PHONY: daemon-status
daemon-status: build ## Check daemon status
	./$(BUILD_DIR)/$(BINARY_NAME) daemon status

# Documentation
.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	@which godoc > /dev/null || (echo "godoc not found. Install with: go install golang.org/x/tools/cmd/godoc@latest" && exit 1)
	@echo "Documentation server: http://localhost:6060/pkg/github.com/ivyascorp-net/nagging-nancy/"
	godoc -http=:6060

# Default target when just running 'make'
.DEFAULT_GOAL := help