.PHONY: build clean install uninstall test deps run help release build-all checksums

BINARY_NAME=pull-vids
# The tag itself on a tagged commit, otherwise tag-distance-hash or a bare hash.
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO_FILES=$(shell find . -name '*.go' -type f)
# Tags carry a leading v, but the banner and every package manager add or want
# the bare number, so strip it here once.
LDFLAGS=-ldflags "-s -w -X main.version=$(patsubst v%,%,$(VERSION))"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go binary
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete! Binary: ./$(BINARY_NAME)"

build-all: ## Build for all platforms (Linux, macOS, Windows)
	@echo "Building for all platforms..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✓ Build complete! Binaries in ./dist/"
	@ls -lh dist/

release: build-all ## Create release archives for distribution
	@echo "Creating release archives..."
	@mkdir -p dist/releases
	cd dist && tar -czf releases/$(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && echo "✓ Linux AMD64"
	cd dist && tar -czf releases/$(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64 && echo "✓ Linux ARM64"
	cd dist && tar -czf releases/$(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && echo "✓ macOS AMD64"
	cd dist && tar -czf releases/$(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && echo "✓ macOS ARM64 (Apple Silicon)"
	cd dist && zip -q releases/$(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe && echo "✓ Windows AMD64"
	@echo ""
	@echo "Release archives created in dist/releases/:"
	@ls -lh dist/releases/

install: build ## Install the binary to /usr/local/bin (Unix/macOS)
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Installed! Run '$(BINARY_NAME)' from anywhere."

uninstall: ## Uninstall the binary from /usr/local/bin
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Uninstalled!"

clean: ## Remove built binaries and dist directory
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	@echo "✓ Clean complete!"

test: ## Run tests
	go test -v ./...

deps: ## Download Go dependencies
	go mod download
	go mod tidy

run: build ## Build and run with help
	./$(BINARY_NAME) --help

checksums: release ## Generate SHA256 checksums for releases
	@echo "Generating checksums..."
	@cd dist/releases && shasum -a 256 * > SHA256SUMS
	@echo "✓ Checksums saved to dist/releases/SHA256SUMS"
	@cat dist/releases/SHA256SUMS

.DEFAULT_GOAL := help
