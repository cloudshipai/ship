# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CloudshipAI CLI ("ship") - A fully functional command-line tool that enables both non-technical users and power users to:
- Run comprehensive Terraform analysis tools in containerized environments
- Generate infrastructure documentation and diagrams  
- Push artifacts (terraform plans, SBOMs, etc.) to Cloudship for analysis
- Host an MCP server for LLM integrations

The project is production-ready with a complete architecture using:
- **Go 1.23** with proper module structure (`github.com/cloudshipai/ship`)
- **Cobra CLI framework** for command structure
- **Dagger** for containerized tool execution
- **Viper** for configuration management
- **MCP (Model Context Protocol)** server implementation

## Development Commands

### Building
```bash
# Build using Makefile (recommended)
make build

# Build directly with Go
go build -o ship ./cmd/ship

# Build all platforms
make build-all
```

### Testing
```bash
# Run tests (excludes integration tests)
make test

# Run all tests including integration
make test-integration

# Run specific module tests
go test -v ./internal/dagger/modules/

# Run with coverage
go test -cover ./...
```

### Linting and Formatting
```bash
# Run linters using golangci-lint (configured in .golangci.yml)
make lint

# Format code
make fmt

# Check dependencies
make deps
```

### Development Tools
```bash
# Install dependencies
make deps

# Clean build artifacts
make clean

# Check release readiness
make release-check
```

## Project Structure

```
ship/
├── cmd/ship/main.go           # Entry point with version handling
├── internal/
│   ├── cli/                   # Cobra command implementations
│   │   ├── root.go           # Root command and logger setup
│   │   ├── auth.go           # Authentication commands
│   │   ├── push.go           # Artifact upload commands
│   │   ├── mcp_cmd.go        # MCP server command
│   │   └── terraform_tools_cmd.go  # Terraform analysis tools
│   ├── auth/                  # Authentication logic
│   ├── config/                # Configuration management (~/.ship/config.yaml)
│   ├── cloudship/             # CloudShip API client
│   ├── dagger/                # Dagger engine and modules
│   │   └── modules/           # Individual tool implementations
│   │       ├── tflint.go     # TFLint integration
│   │       ├── checkov.go    # Checkov security scanning
│   │       ├── infracost.go  # Cost estimation
│   │       ├── inframap.go   # Diagram generation
│   │       └── ...
│   └── logger/                # Structured logging
├── examples/                  # Sample Terraform projects for testing
├── docs/                      # Comprehensive documentation
├── demos/                     # Generated demo GIFs and scripts
└── Makefile                   # Build automation
```

## Architecture Overview

The CLI uses a **containerized execution model** via Dagger:

1. **Command Layer**: Cobra commands parse user input and flags
2. **Dagger Engine**: Orchestrates containerized tool execution
3. **Tool Modules**: Individual Go modules wrapping Docker-based tools
4. **Configuration**: Centralized config management with Viper
5. **MCP Server**: Model Context Protocol server for AI assistant integration

### Key Architectural Decisions

- **Containerization**: All tools run in Docker containers via Dagger (no local installs)
- **Modularity**: Each tool is implemented as a separate Dagger module
- **Configuration**: Uses `~/.ship/config.yaml` for settings and API keys
- **Logging**: Structured logging with configurable levels and file output
- **Error Handling**: Comprehensive error handling with user-friendly messages
- **Testing**: Extensive test coverage including integration tests

## Available Tools and Commands

The CLI provides comprehensive Terraform analysis through these tools:

### Core Commands
- **`ship auth`** - Manage CloudShip API authentication
- **`ship push`** - Upload artifacts to CloudShip for analysis
- **`ship mcp`** - Start MCP server for AI assistant integration
- **`ship modules`** - Manage and discover external Dagger modules
- **`ship terraform-tools`** - Run containerized Terraform analysis tools

### Terraform Tools (via Dagger containers)
- **`lint`** - TFLint for syntax and best practices
- **`checkov-scan`** - Security and compliance scanning with Checkov
- **`security-scan`** - Alternative security scanning with Trivy
- **`cost-estimate`** - Cost estimation with Infracost
- **`cost-analysis`** - Alternative cost analysis with OpenInfraQuote
- **`generate-docs`** - Documentation generation with terraform-docs
- **`generate-diagram`** - Infrastructure diagrams with InfraMap

All terraform-tools commands support:
- `--push` flag for automatic CloudShip upload
- Output redirection and format options
- Directory specification for multi-module projects

## Development Context

- **Current Status**: Production-ready CLI with full tool suite implemented
- **Go Version**: 1.23.0 with toolchain 1.23.10
- **Key Dependencies**: Dagger v0.18.10, Cobra v1.9.1, Viper v1.20.1
- **Testing**: Comprehensive test suite including integration tests
- **Documentation**: Complete docs in `/docs` folder and demo GIFs in `/demos`
- **Examples**: Multiple Terraform examples in `/examples` for testing

## Important Implementation Details

- **Module Path**: `github.com/cloudshipai/ship` (not a placeholder)
- **Version Handling**: Uses ldflags injection and build info fallback in `cmd/ship/main.go:18`
- **Configuration**: Stored in `~/.ship/config.yaml` via Viper
- **Docker Requirement**: Dagger requires Docker daemon for containerized execution
- **Error Handling**: Uses structured logging with `internal/logger` package
- **MCP Integration**: Full MCP server implementation for AI assistant integration

## License

This project is licensed under the Apache License 2.0.