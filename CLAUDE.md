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

The project implements five main commands:

1. **`ship auth`** - Manages CloudShip API authentication
2. **`ship push`** - Uploads artifacts to CloudShip for analysis
3. **`ship investigate`** - Runs Steampipe queries via Dagger containers
4. **`ship mcp`** - Starts MCP server for LLM tool integrations
5. **`ship terraform-tools`** - Runs Terraform analysis tools via Dagger:
   - `cost-analysis` - Estimates costs using OpenInfraQuote
   - `security-scan` - Scans for security issues using InfraScan
   - `generate-docs` - Generates documentation using terraform-docs
   - `lint` - Lints Terraform code using TFLint
   - `checkov-scan` - Scans for security issues using Checkov
   - `cost-estimate` - Estimates costs using Infracost
   - `generate-diagram` - Generates infrastructure diagrams using InfraMap

Key architectural decisions:
- Uses Cobra for CLI framework
- Embeds Dagger CLI for container orchestration
- Configuration stored in `~/.ship/config.yaml`
- Max file upload size: 100MB
- Supports AWS, Cloudflare, and Heroku providers
- Dagger modules run containerized tools without local installation
- All terraform-tools commands support --push flag for automatic CloudShip upload
- CloudShip API authentication via API keys (from https://app.cloudshipai.com/settings/api-keys)

## Project Documentation

Detailed documentation is available in the `docs/` folder:
- `PRD.md` - Product Requirements Document
- `implementation-plan.md` - Phased development plan
- `technical-spec.md` - Architecture and component specifications
- `development-tasks.md` - Sprint-by-sprint task breakdown
- `api-reference.md` - Cloudship API and MCP protocol specs

## Development Workflow Memories

- Okay, let's begin building it and check off as you go, commit as you go, and write Notion documentation about this project as you go

## Recent Fixes

### AI Hallucination Fix (July 2, 2025)
- Fixed AI generating non-existent Steampipe tables like `aws_resource_inventory`
- Added `steampipe_tables.go` with real table names for AWS, Azure, GCP
- Enhanced AI prompts to include actual tables and example queries
- Fixed hardcoded fallback in `llm_dagger.go`
- Improved JSON parsing to handle markdown code blocks

### AWS Credentials Fix (July 2, 2025)
- Fixed automatic credential detection from `~/.aws/credentials`
- Changed from profile-based to environment variable authentication
- Credentials now parsed on host and passed as env vars to containers
- Removed problematic credential file mounting

## License

This project is licensed under the Apache License 2.0.