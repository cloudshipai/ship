#!/bin/bash

# Local Install Script for Ship Development
# This script sets up the local development environment

set -e

echo "🚢 Setting up Ship local development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [[ "$(echo "$GO_VERSION 1.21" | tr " " "\n" | sort -V | head -n 1)" != "1.21" ]]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go 1.21+ first."
    exit 1
fi

echo "✅ Go version $GO_VERSION detected"

# Check if Docker is installed (for GitHub MCP server)
if ! command -v docker &> /dev/null; then
    echo "⚠️  Docker is not installed. GitHub MCP server will not work."
    echo "   Install Docker to use GitHub MCP server functionality."
fi

# Check if Node.js/npm is installed (for external MCP servers)
if ! command -v npm &> /dev/null; then
    echo "⚠️  Node.js/npm is not installed. Some external MCP servers may not work."
    echo "   Install Node.js to use npm-based MCP servers."
fi

# Check if uv is installed (for AWS Labs MCP servers)
if ! command -v uv &> /dev/null; then
    echo "⚠️  uv is not installed. AWS Labs MCP servers will not work."
    echo "   Install uv to use AWS Labs MCP servers:"
    echo "   curl -LsSf https://astral.sh/uv/install.sh | sh"
fi

# Build Ship CLI
echo "🔨 Building Ship CLI..."
go build -o bin/ship ./cmd/ship

if [ $? -eq 0 ]; then
    echo "✅ Ship CLI built successfully"
    echo "   Binary location: bin/ship"
else
    echo "❌ Failed to build Ship CLI"
    exit 1
fi

# Build CLI package
echo "🔨 Building CLI package..."
go build ./internal/cli

if [ $? -eq 0 ]; then
    echo "✅ CLI package built successfully"
else
    echo "❌ Failed to build CLI package"
    exit 1
fi

# Create bin directory if it doesn't exist
mkdir -p bin

# Make the script executable
chmod +x bin/ship

echo ""
echo "🎉 Ship local development environment setup complete!"
echo ""
echo "Usage:"
echo "  ./bin/ship --help                    # Show help"
echo "  ./bin/ship modules list              # List all modules"
echo "  ./bin/ship modules info slack        # Show Slack MCP server info"
echo "  ./bin/ship modules info github       # Show GitHub MCP server info"
echo ""
echo "Testing:"
echo "  ./bin/ship mcp slack --var SLACK_MCP_XOXC_TOKEN=your_token"
echo "  ./bin/ship mcp github --var GITHUB_PERSONAL_ACCESS_TOKEN=your_token"
echo ""
echo "Note: Make sure to set appropriate environment variables for testing."
