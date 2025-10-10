.PHONY: build build-all test lint clean install fmt tidy help

# Binary name
BINARY_NAME=forkspacer

# Build variables
VERSION?=dev
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/forkspacer/cli/cmd.version=$(VERSION) -X github.com/forkspacer/cli/cmd.gitCommit=$(GIT_COMMIT) -X github.com/forkspacer/cli/cmd.buildDate=$(BUILD_DATE)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build binary for current platform
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

build-all: ## Build for all platforms
	@echo "Building for multiple platforms..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe .
	@echo "Build complete!"

test: ## Run tests
	$(GOTEST) -v ./...

lint: ## Run linters
	@echo "Running go vet..."
	@$(GOVET) ./...
	@echo "Checking formatting..."
	@if [ -n "$$(gofmt -s -l .)" ]; then \
		echo "Code is not formatted. Run 'make fmt'"; \
		gofmt -s -l .; \
		exit 1; \
	fi
	@echo "Linting passed!"

fmt: ## Format code
	$(GOFMT) ./...

tidy: ## Tidy go.mod
	$(GOMOD) tidy

install: build ## Install binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "Installed successfully!"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)*
	@rm -f coverage.out
	@echo "Clean complete!"

run: build ## Build and run version command
	./$(BINARY_NAME) version

verify: lint test build ## Run all verification steps
	@echo "All verification steps passed!"

.DEFAULT_GOAL := help
