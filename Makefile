.PHONY: build clean install test help

BINARY_NAME=pull-vids
GO_FILES=$(shell find . -name '*.go' -type f)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go binary
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go
	@echo "Build complete! Binary: ./$(BINARY_NAME)"

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe main.go
	@echo "Build complete! Binaries in ./dist/"

install: build ## Install the binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "Installed! Run '$(BINARY_NAME)' from anywhere."

clean: ## Remove built binaries
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	@echo "Clean complete!"

test: ## Run tests
	go test -v ./...

deps: ## Download Go dependencies
	go mod download
	go mod tidy

run: build ## Build and run with example (requires URL argument)
	./$(BINARY_NAME) --help

.DEFAULT_GOAL := help
