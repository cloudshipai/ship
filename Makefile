# Makefile for Ship CLI

# Variables
BINARY_NAME := ship
GO_MODULE := github.com/cloudshipai/ship
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS := -ldflags "-s -w -X '$(GO_MODULE)/cmd/ship.version=$(VERSION)' -X '$(GO_MODULE)/cmd/ship.commit=$(shell git rev-parse HEAD)' -X '$(GO_MODULE)/cmd/ship.date=$(BUILD_TIME)'"

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
.PHONY: help
help:
	@echo "Ship CLI Makefile"
	@echo "================="
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

## build: Build the binary for current platform
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/ship

## build-all: Build binaries for all platforms
.PHONY: build-all
build-all:
	@echo "Building for all platforms..."
	goreleaser build --snapshot --clean

## test: Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v -short ./...

## test-integration: Run all tests including integration
.PHONY: test-integration
test-integration:
	@echo "Running all tests..."
	go test -v ./...

## lint: Run linters
.PHONY: lint
lint:
	@echo "Running linters..."
	golangci-lint run ./...

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist/

## install: Install binary to /usr/local/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/

## local-install: Install binary to ~/.local/bin (no sudo required)
.PHONY: local-install
local-install: build
	@echo "Installing $(BINARY_NAME) to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	cp $(BINARY_NAME) ~/.local/bin/
	@echo "✅ $(BINARY_NAME) installed to ~/.local/bin/"
	@echo ""
	@echo "Make sure ~/.local/bin is in your PATH:"
	@echo "  export PATH=\"\$$HOME/.local/bin:\$$PATH\""

## uninstall: Remove binary from /usr/local/bin
.PHONY: uninstall
uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

## local-uninstall: Remove binary from ~/.local/bin
.PHONY: local-uninstall
local-uninstall:
	@echo "Removing $(BINARY_NAME) from ~/.local/bin..."
	rm -f ~/.local/bin/$(BINARY_NAME)

## deps: Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

## deps-update: Update dependencies
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

## release-check: Check if ready for release
.PHONY: release-check
release-check:
	@echo "Checking release readiness..."
	@echo "Version: $(VERSION)"
	@echo "Running tests..."
	@make test
	@echo "Checking GoReleaser config..."
	goreleaser check
	@echo "✅ Ready for release!"

## release-snapshot: Create a snapshot release (no upload)
.PHONY: release-snapshot
release-snapshot: clean
	@echo "Creating snapshot release..."
	goreleaser release --snapshot --clean

## release-local: Create a local release (no upload)
.PHONY: release-local
release-local: release-check
	@echo "Creating local release..."
	@if [ -z "$(GITHUB_TOKEN)" ]; then \
		echo "❌ GITHUB_TOKEN not set. Running in local mode..."; \
		goreleaser release --snapshot --clean --skip=publish; \
	else \
		echo "✅ GITHUB_TOKEN found"; \
		goreleaser release --clean --skip=validate; \
	fi

## release: Create and publish a release
.PHONY: release
release: release-check
	@echo "Creating release $(VERSION)..."
	@if [ -z "$(GITHUB_TOKEN)" ]; then \
		echo "❌ Error: GITHUB_TOKEN environment variable not set"; \
		echo ""; \
		echo "To create a release, run:"; \
		echo "  export GITHUB_TOKEN=your-github-personal-access-token"; \
		echo "  make release"; \
		exit 1; \
	fi
	@echo "✅ GITHUB_TOKEN found"
	@echo "🚀 Creating release..."
	goreleaser release --clean

## release-major: Create a major release (v1.0.0 -> v2.0.0)
.PHONY: release-major
release-major:
	@scripts/release.sh major

## release-minor: Create a minor release (v1.0.0 -> v1.1.0)
.PHONY: release-minor
release-minor:
	@scripts/release.sh minor

## release-patch: Create a patch release (v1.0.0 -> v1.0.1)
.PHONY: release-patch
release-patch:
	@scripts/release.sh patch

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t cloudshipai/ship:latest -t cloudshipai/ship:$(VERSION) .

## docker-push: Push Docker image to registry
.PHONY: docker-push
docker-push: docker-build
	@echo "Pushing Docker images..."
	docker push cloudshipai/ship:latest
	docker push cloudshipai/ship:$(VERSION)

## docker-run: Run Docker container
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it cloudshipai/ship:latest

## docker-test: Test Docker container
.PHONY: docker-test
docker-test:
	@echo "Testing Docker container..."
	docker run --rm cloudshipai/ship:latest version
	docker run --rm --group-add=999 -v /var/run/docker.sock:/var/run/docker.sock cloudshipai/ship:latest dagger-test

.PHONY: all
all: clean deps lint test build