# Ship MCP Framework

A comprehensive DevOps toolkit that provides **63 essential infrastructure tools** through AI-enabled MCP (Model Context Protocol) servers. Ship runs all tools securely in containers via Dagger, eliminating local dependencies.

Ship is primarily an **MCP framework** for building AI assistant integrations, with optional CLI capabilities for direct usage.

> **🤖 For LLMs and AI Assistants**: Complete installation and usage instructions specifically designed for AI consumption are available in [llms.txt](./llms.txt). This includes MCP server setup, integration examples, and best practices for AI-driven infrastructure analysis with all **63+ tools**.

## 🚀 Features

### Ship MCP Framework
- **🏗️ MCP Server Builder**: Fluent API for building custom MCP servers for AI assistants
- **🔧 Container Tool Framework**: Run any tool securely in Docker containers via Dagger
- **📦 Pre-built Ship Tools**: **63 essential DevOps tools** ready to use
- **🎯 Multiple Usage Patterns**: Pure framework, cherry-pick tools, or everything plus custom extensions
- **🔒 Security First**: All tools run in isolated containers with no local dependencies
- **⚡ Performance Optimized**: Leverages Dagger's caching and parallel execution
- **🤖 AI Assistant Ready**: Built for Claude, Cursor, and other MCP-compatible AI tools
- **🧪 Test Coverage**: Comprehensive test suite with integration tests
- **📚 Rich Documentation**: Complete API reference and usage examples

### DevOps Toolkit Coverage
- **🔍 Terraform Development**: `tflint`, `terraform-docs`, `checkov`, `tfsec`, `openinfraquote`
- **🛡️ Container Security**: `trivy`, `grype`, `syft`, `dockle`, `cosign`, `hadolint`
- **🔐 Secret Management**: `trufflehog`, `gitleaks`, `git-secrets`, `sops`
- **☸️ Kubernetes Operations**: `kubescape`, `kube-bench`, `velero`, `falco`, `goldilocks`
- **☁️ Cloud Security**: `prowler`, `scout-suite`, `steampipe`, `cloudquery`, `custodian`
- **🌐 Web Application Testing**: `nuclei`, `zap`, `nikto`, `nmap`
- **👨‍💻 Development & CI/CD**: `semgrep`, `actionlint`, `opencode`
- **📦 Supply Chain Security**: `cosign`, `syft`, `dependency-track`
- **🔗 External Integrations**: 16 MCP servers (`postgresql`, `playwright`, `bitbucket`, etc.)

## 🎯 Why Ship MCP Framework?

Ship is designed for the **AI-first infrastructure era** where AI assistants need secure, reliable access to infrastructure tools:

- **🤖 AI Assistant Native**: Built specifically for Claude, Cursor, and other AI assistants
- **🔒 Security by Design**: All tools run in isolated containers - no local tool installations
- **📦 Curated & Tested**: CloudShip AI team maintains and tests all included tools
- **🏗️ Framework First**: Extensible architecture for custom tool development
- **⚡ Developer Experience**: Simple APIs with comprehensive documentation

## 📚 Table of Contents

- [Installation](#-installation)
- [Quick Start](#-quick-start)
  - [Ship MCP Framework](#1-ship-mcp-framework)
  - [AI Assistant Integration](#2-ai-assistant-integration)
  - [Optional CLI Usage](#3-optional-cli-usage)
- [Available Tools Reference](#-available-tools-reference)
- [Documentation](#-documentation)
- [Contributing](#-contributing)
- [License](#-license)

## 📦 Installation

### Quick Install (CLI)

```bash
# Install Ship CLI with one command
curl -fsSL https://raw.githubusercontent.com/cloudshipai/ship/main/install.sh | bash
```

### Install with Go

```bash
# Install directly with Go  
go install github.com/cloudshipai/ship/cmd/ship@latest
```

## 🏃 Quick Start

### 1. Ship MCP Framework

Build your own MCP servers using the Ship framework:

#### Quick Example: Infrastructure MCP Server

```bash
# Create a new Go project
mkdir my-infrastructure-server && cd my-infrastructure-server
go mod init my-infrastructure-server
go get github.com/cloudshipai/ship
```

Create `main.go`:

```go
package main

import (
    "log"
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/internal/tools"
)

func main() {
    // Build MCP server with Ship's pre-built infrastructure tools
    server := ship.NewServer("infrastructure-server", "1.0.0").
        AddTool(tools.NewTFLintTool()).     // Terraform linting
        Build()

    // Start the MCP server
    if err := server.ServeStdio(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

```bash
# Build and run your MCP server
go build -o infrastructure-server .
./infrastructure-server
```

#### Four Integration Patterns

**1. Ship-First (Recommended):**
```go
// Ship manages the entire MCP server
server := ship.NewServer("my-server", "1.0.0").
    AddTool(tools.NewTFLintTool()).
    Build()
server.ServeStdio()
```

**2. Bring Your Own MCP Server:**
```go
// Add Ship tools to existing mcp-go server  
shipAdapter := ship.NewMCPAdapter().
    AddTool(tools.NewTFLintTool())
shipAdapter.AttachToServer(ctx, existingMCPServer)
```

**3. Tool Router (Advanced):**
```go
// Route different tools to different servers
router := ship.NewToolRouter().
    AddRoute("terraform", terraformAdapter).
    AddRoute("security", securityAdapter)
```

**4. CLI Usage (Optional):**
```bash
# Direct CLI access to tools
ship mcp all  # Start MCP server with all tools
```

📚 **Complete Integration Guide**: See [examples/integration-patterns.md](examples/integration-patterns.md) for detailed integration patterns with code examples.

### 2. AI Assistant Integration (MCP)

Configure your custom MCP servers built with Ship in Claude Desktop or other AI assistants:

```json
{
  "mcpServers": {
    "ship-all-tools": {
      "command": "ship",
      "args": ["mcp", "all"]
    },
    "ship-terraform": {
      "command": "ship",
      "args": ["mcp", "terraform"],
      "env": {
        "AWS_PROFILE": "your-profile"
      }
    },
    "ship-security": {
      "command": "ship",
      "args": ["mcp", "security"]
    },
    "ship-external": {
      "command": "ship",
      "args": ["mcp", "filesystem", "--var", "FILESYSTEM_ROOT=/workspace"]
    }
  }
}
```

**Ship Framework wraps mcp-go** to provide:
- **Container-based Tools**: All tools run securely in Docker containers
- **Fluent Builder API**: Easy-to-use APIs for building MCP servers
- **Pre-built Ship Tools**: 63 infrastructure tools ready to use in your servers
- **Multiple Usage Patterns**: Pure framework, cherry-pick tools, or everything plus extensions

### 3. Optional CLI Usage

Ship also includes a CLI for direct usage of infrastructure tools:

```bash
# Navigate to your project directory
cd your-infrastructure-project

# Start MCP servers for AI assistant integration (MAIN USE CASE)
ship mcp all                # All 63 tools across all categories
ship mcp terraform          # All Terraform tools (tflint, terraform-docs, checkov, etc.)
ship mcp security           # All security tools (trivy, grype, nuclei, etc.)
ship mcp kubernetes         # All Kubernetes tools (kubescape, velero, etc.)

# External MCP server proxying
ship mcp filesystem --var FILESYSTEM_ROOT=/workspace
ship mcp brave-search --var BRAVE_API_KEY=your_api_key
ship mcp postgresql --var POSTGRES_CONNECTION_STRING=postgresql://...

# AI Development Tools
ship opencode chat "explain this terraform module" --model "openai/gpt-4o-mini"
ship opencode generate "create a kubernetes deployment" --output k8s/deployment.yaml
ship opencode analyze app.py --model "anthropic/claude-3-5-sonnet-20241022"

# Tool information and discovery
ship modules list           # List all available tools
ship modules info terraform # See details about terraform tools
```

## 🛠️ Available Tools Reference

Ship provides **63 essential DevOps tools** across security, infrastructure, cloud, and development workflows. All tools run in isolated containers via Dagger.

### Quick Reference by Workflow

| Workflow | Primary Tools | Supporting Tools |
|----------|---------------|------------------|
| **Terraform Development** | `tflint`, `terraform-docs`, `checkov`, `tfsec` | `inframap`, `terrascan`, `openinfraquote` |
| **Container Security** | `trivy`, `grype`, `syft` | `dockle`, `cosign`, `hadolint` |
| **Secret Management** | `trufflehog`, `gitleaks`, `git-secrets` | `sops` |
| **Kubernetes Operations** | `kubescape`, `kube-bench`, `velero` | `falco`, `kyverno`, `goldilocks` |
| **Cloud Security** | `prowler`, `scout-suite`, `steampipe` | `cloudquery`, `custodian` |
| **Web Application Testing** | `nuclei`, `zap`, `nikto` | `nmap` |
| **Development & CI/CD** | `semgrep`, `actionlint`, `gitleaks` | `hadolint`, `conftest` |

### Tool Categories Summary

| Category | Count | Key Tools | Purpose |
|----------|-------|-----------|---------|
| **Terraform** | 11 | `tflint`, `terraform-docs`, `checkov` | IaC development & security |
| **Security** | 31 | `trivy`, `trufflehog`, `kubescape` | Vulnerability & secret scanning |
| **Cloud** | 17 | `prowler`, `terraformer`, `packer` | Cloud infrastructure & governance |
| **Supply Chain** | 3 | `cosign`, `syft`, `dependency-track` | Supply chain security |
| **AWS** | 4 | `cloudsplaining`, `prowler`, `parliament` | AWS security & IAM |
| **Development** | 1 | `opencode` | AI-powered development |
| **External MCPs** | 16 | `postgresql`, `playwright`, `bitbucket` | Third-party integrations |

### Usage Examples

```bash
# Terraform workflow
ship mcp terraform  # All Terraform tools (tflint, terraform-docs, checkov, etc.)

# Container security pipeline  
ship mcp security   # All security tools (trivy, grype, syft, etc.)

# Kubernetes operations
ship mcp kubernetes # All K8s tools (kubescape, velero, goldilocks, etc.)

# Cloud security assessment
ship mcp cloud      # All cloud tools (prowler, scout-suite, etc.)

# Full DevOps toolkit
ship mcp all        # All 63 tools across all categories
```

📋 **Complete Tools Reference**: See [docs/tools-reference-table.md](docs/tools-reference-table.md) for detailed tool descriptions and usage guidance.

## ⚙️ Environment Variables and Configuration

### --var Flag System

Ship supports passing environment variables to both containerized tools and external MCP servers using the `--var` flag:

```bash
# Single variable
ship mcp brave-search --var BRAVE_API_KEY=your_api_key

# Multiple variables
ship mcp postgresql --var POSTGRES_CONNECTION_STRING=postgresql://user:pass@host:5432/db

# Variables for containerized tools
ship mcp all --var AWS_REGION=us-east-1 --var DEBUG=true
```

### Variable Discovery

Use `ship modules info <tool-name>` to see available variables:

```bash
# See variables for external MCP servers
ship modules info filesystem
ship modules info postgresql
ship modules info brave-search

# See information about built-in tools
ship modules info terraform
ship modules info security
```

## 🔧 Ship Framework Integration

### mcp-go Integration

Ship Framework wraps the official [mcp-go](https://github.com/modelcontextprotocol/mcp-go) library to provide:

- **Enhanced Container Tools**: All tools run in isolated Docker containers via Dagger
- **Fluent Builder API**: Easy-to-use APIs for building MCP servers
- **Security by Default**: No local tool installations required
- **Pre-built Infrastructure Tools**: Ready-to-use Terraform analysis tools

### Example: Framework Usage

```go
package main

import (
    "log"
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/internal/tools"
)

func main() {
    // Ship builds on mcp-go's foundation
    server := ship.NewServer("infrastructure-server", "1.0.0").
        AddTool(tools.NewTFLintTool()).         // Pre-built Ship tool
        Build()
    
    // Leverages mcp-go's stdio transport
    if err := server.ServeStdio(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

## 🔐 Container Security

Ship Framework provides security through containerization:

- **No Local Dependencies**: All tools run in containers
- **Isolated Execution**: Each tool runs in its own container
- **Dagger Engine**: Secure container orchestration
- **Credential Passthrough**: Environment variables passed securely to containers

### Supported Cloud Providers

```go
// AWS credentials passed automatically from environment
server := ship.NewServer("aws-server", "1.0.0").
    AddTool(tools.NewTFLintTool()).
    Build()

// Environment variables are automatically passed to containers
// Set AWS_PROFILE, AWS_REGION etc. in your environment
```

## 🏗️ Framework Architecture

Ship MCP Framework is built on:
- **mcp-go**: Official Model Context Protocol implementation
- **Dagger Engine**: Container orchestration and caching
- **Builder Pattern**: Fluent APIs for server construction
- **Registry Pattern**: Extensible tool registration system

### Benefits
- **Consistency**: Same tool versions across all environments
- **Isolation**: No conflicts with local installations  
- **Security**: Tools run in sandboxed containers
- **Simplicity**: No need to install or manage tool versions

## 🤝 Contributing

We welcome contributions! See our [Contributing Guide](CONTRIBUTING.md) for details.

### Adding New Tools

1. Create a new module in `internal/dagger/modules/`
2. Add CLI command in `internal/cli/`
3. Update documentation
4. Submit a pull request

## 📚 Documentation

### CLI Documentation
- [CLI Reference](docs/cli-reference.md) - Complete command reference
- [MCP Integration Guide](docs/mcp-integration.md) - AI assistant integration setup
- [External MCP Servers](docs/mcp-external-servers.md) - Proxying external MCP servers with --var flags
- [Tools Use Cases Guide](docs/tools-use-cases.md) - When to use each tool (75 pages)
- [Tools Reference Table](docs/tools-reference-table.md) - Complete tools reference with examples
- [Dynamic Module Discovery](docs/dynamic-module-discovery.md) - Extensible module system
- [Development Guide](docs/development-tasks.md) - For contributors
- [Technical Spec](docs/technical-spec.md) - Architecture and design

### Ship SDK Documentation
- [📖 Ship SDK Overview](docs/ship-sdk-overview.md) - Framework introduction and concepts
- [🚀 Quick Start Guide](docs/ship-sdk-quickstart.md) - Step-by-step getting started
- [📋 API Reference](docs/ship-sdk-api-reference.md) - Complete API documentation
- [🎯 Usage Patterns](docs/ship-sdk-usage-patterns.md) - Advanced patterns and best practices
- [🛠️ Ship Tools Reference](docs/ship-tools-reference.md) - Pre-built tools documentation
- [❓ FAQ](docs/ship-sdk-faq.md) - Frequently asked questions

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run integration tests
go test -v ./internal/dagger/modules/

# Test specific module
go test -v -run TestTFLintModule ./internal/dagger/modules/
```

## 📈 Roadmap

- [ ] Enhanced mcp-go integration with streaming support
- [ ] Policy as Code with Open Policy Agent containers
- [ ] Web UI for MCP server management and testing
- [ ] More pre-built infrastructure tools
- [ ] Kubernetes and cloud-native tooling support
- [ ] Integration with more cloud providers

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Ship Framework wouldn't be possible without these amazing open source projects:
- [mcp-go](https://github.com/modelcontextprotocol/mcp-go) - Official Model Context Protocol implementation
- [Dagger](https://dagger.io) - For containerized execution
- [Cobra](https://github.com/spf13/cobra) - For CLI framework
- All the individual tool maintainers

---

**Built with ❤️ by the CloudshipAI team**