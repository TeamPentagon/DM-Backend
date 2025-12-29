# DM-Backend Makefile
# Common development tasks automation

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary name
BINARY_NAME=dm-backend
BINARY_UNIX=$(BINARY_NAME)_unix

# Directories
BUILD_DIR=./build
COVERAGE_DIR=./coverage

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^##/ {c=$$0} /^[a-zA-Z_-]+:/ {gsub(/:.*/,""); print "\t" $$0 "\t" substr(c, 4)}' $(MAKEFILE_LIST)

## all: Build and test the application
.PHONY: all
all: test build

## build: Build the application
.PHONY: build
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v .

## run: Run the application
.PHONY: run
run:
	@echo "Running..."
	$(GOCMD) run main.go

## test: Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-cover: Run tests with coverage report
.PHONY: test-cover
test-cover:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"

## test-short: Run tests in short mode (skip slow tests)
.PHONY: test-short
test-short:
	@echo "Running short tests..."
	$(GOTEST) -v -short ./...

## bench: Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	rm -rf Database/

## deps: Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## vet: Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

## lint: Run linter (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## proto: Generate protobuf files
.PHONY: proto
proto:
	@echo "Generating protobuf files..."
	@if command -v protoc >/dev/null 2>&1; then \
		cd internal/model && \
		protoc --go_out=. --go_opt=paths=source_relative user.proto && \
		protoc --go_out=. --go_opt=paths=source_relative chat.proto; \
		echo "Protobuf files generated successfully"; \
	else \
		echo "protoc not installed. Please install Protocol Buffers compiler."; \
		exit 1; \
	fi

## check: Run all checks (fmt, vet, lint, test)
.PHONY: check
check: fmt vet lint test
	@echo "All checks passed!"

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t dm-backend:latest .

## docker-run: Run Docker container
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8085:8085 dm-backend:latest

## install: Install the binary
.PHONY: install
install:
	@echo "Installing..."
	$(GOCMD) install .

## cross-compile: Build for Linux (from any platform)
.PHONY: cross-compile
cross-compile:
	@echo "Cross-compiling for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v .

## release: Create a release build
.PHONY: release
release: clean test
	@echo "Creating release build..."
	@mkdir -p $(BUILD_DIR)/release
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 .
	# macOS
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 .
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Release builds created in $(BUILD_DIR)/release/"

## dev: Start development mode with auto-reload (requires air)
.PHONY: dev
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not installed. Install with:"; \
		echo "  go install github.com/cosmtrek/air@latest"; \
		echo ""; \
		echo "Falling back to normal run..."; \
		$(GOCMD) run main.go; \
	fi
