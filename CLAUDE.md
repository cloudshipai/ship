# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CloudshipAI CLI ("ship") - A command-line tool that enables both non-technical users and power users to:
- Push artifacts (terraform plans, SBOMs, etc.) to Cloudship for analysis
- Run automated cloud infrastructure investigations using Steampipe
- Host an MCP server for LLM integrations

Key capabilities:
- One-liner artifact pushes with automatic cost/security analysis
- One-command investigations powered by Steampipe + Dagger
- Optional MCP host for offline LLM tool access

## Development Setup

Since this is a Go project in its initial state, you'll need to initialize it:

```bash
# Initialize Go module (if not already done)
go mod init github.com/[organization]/ship

# Download dependencies (after adding imports)
go mod download

# Tidy up dependencies
go mod tidy
```

## Common Commands

### Building
```bash
# Build the CLI binary
go build -o ship ./cmd/ship

# Build with version information
go build -ldflags "-X main.version=v1.0.0" -o ship ./cmd/ship
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -run TestName ./package
```

### Linting and Formatting
```bash
# Format code
go fmt ./...

# Run go vet
go vet ./...

# If golangci-lint is installed
golangci-lint run
```

## Expected Project Structure

For a Go CLI application, the typical structure would be:

```
ship/
├── cmd/
│   └── ship/          # Main application entry point
│       └── main.go
├── internal/          # Private application code
│   ├── cli/          # CLI command implementations
│   ├── config/       # Configuration handling
│   └── ...
├── pkg/              # Public libraries (if any)
├── go.mod           # Go module file
├── go.sum           # Go module checksums
├── Makefile         # Build automation (optional)
└── .golangci.yml    # Linter configuration (optional)
```

## Architecture Notes

The project implements four main commands:

1. **`ship auth`** - Manages authentication tokens and config files
2. **`ship push`** - Uploads artifacts (tfplan, SBOM, etc.) for analysis
3. **`ship investigate`** - Runs Steampipe queries via Dagger containers
4. **`ship mcp`** - Starts MCP server for LLM tool integrations
5. **`ship terraform-tools`** - Runs Terraform analysis tools via Dagger:
   - `cost-analysis` - Estimates costs using OpenInfraQuote
   - `security-scan` - Scans for security issues using Trivy
   - `generate-docs` - Generates documentation using terraform-docs

Key architectural decisions:
- Uses Cobra for CLI framework
- Embeds Dagger CLI for container orchestration
- Configuration stored in `~/.ship/config.yaml`
- Max file upload size: 100MB
- Supports AWS, Cloudflare, and Heroku providers
- Dagger modules run containerized tools without local installation

## Project Documentation

Detailed documentation is available in the `docs/` folder:
- `PRD.md` - Product Requirements Document
- `implementation-plan.md` - Phased development plan
- `technical-spec.md` - Architecture and component specifications
- `development-tasks.md` - Sprint-by-sprint task breakdown
- `api-reference.md` - Cloudship API and MCP protocol specs

## Development Workflow Memories

- Okay, let's begin building it and check off as you go, commit as you go, and write Notion documentation about this project as you go

## License

This project is licensed under the Apache License 2.0.