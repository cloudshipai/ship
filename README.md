# Ship MCP Framework

A collection of CloudShip AI team curated MCP servers that run on top of Dagger engine with the ability to use it as a framework to build MCP servers that run securely in containers.

Ship is primarily an MCP (Model Context Protocol) framework for building AI assistant integrations, with optional CLI capabilities for direct usage.

> **🤖 For LLMs and AI Assistants**: Complete installation and usage instructions specifically designed for AI consumption are available in [llms.txt](./llms.txt). This includes MCP server setup, integration examples, and best practices for AI-driven infrastructure analysis with all 7 MCP tools.

## 🚀 Features

### Ship MCP Framework
- **🏗️ MCP Server Builder**: Fluent API for building custom MCP servers for AI assistants
- **🔧 Container Tool Framework**: Run any tool securely in Docker containers via Dagger
- **📦 Pre-built Ship Tools**: Curated collection of infrastructure tools ready to use
- **🎯 Three Usage Patterns**: Pure framework, cherry-pick tools, or everything plus custom extensions
- **🔒 Security First**: All tools run in isolated containers with no local dependencies
- **⚡ Performance Optimized**: Leverages Dagger's caching and parallel execution
- **🤖 AI Assistant Ready**: Built for Claude, Cursor, and other MCP-compatible AI tools
- **🧪 Test Coverage**: Comprehensive test suite with integration tests
- **📚 Rich Documentation**: Complete API reference and usage examples

### Available Infrastructure Tools
- **🔍 Terraform Linting**: TFLint for catching errors and enforcing best practices
- **🛡️ Security Scanning**: Checkov and Trivy for multi-cloud security analysis
- **💰 Cost Estimation**: OpenInfraQuote and Infracost for infrastructure cost analysis
- **📝 Documentation Generation**: terraform-docs for beautiful module documentation
- **📊 Infrastructure Diagrams**: InfraMap for visualizing infrastructure
- **🐳 Containerized Execution**: All tools run via Dagger - no local installations needed
- **☁️ Multi-Cloud Support**: Works with AWS, Azure, GCP, and other cloud providers

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
- [Available Tools](#-available-tools)
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

### 1. Ship SDK Framework

Build your own MCP servers using the Ship SDK:

#### Quick Example: Echo Server

```bash
# Create a new Go project
mkdir my-mcp-server && cd my-mcp-server
go mod init my-mcp-server
go get github.com/cloudshipai/ship
```

Create `main.go`:

```go
package main

import (
    "context"
    "log"
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/pkg/dagger"
)

func main() {
    // Create a simple echo tool
    echoTool := ship.NewContainerTool("echo", ship.ContainerToolConfig{
        Description: "Echo a message",
        Image:       "alpine:latest",
        Parameters: []ship.Parameter{
            {Name: "message", Type: "string", Description: "Message to echo", Required: true},
        },
        Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
            message := params["message"].(string)
            result, err := engine.Container().From("alpine:latest").WithExec([]string{"echo", message}).Stdout(ctx)
            if err != nil {
                return &ship.ToolResult{Error: err}, err
            }
            return &ship.ToolResult{Content: result}, nil
        },
    })

    // Build and start the MCP server
    server := ship.NewServer("echo-server", "1.0.0").AddTool(echoTool).Build()
    if err := server.ServeStdio(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

```bash
# Build and run your MCP server
go build -o echo-server .
./echo-server
```

#### Three Usage Patterns

**1. Single Tool MCP Server:**
```bash
# Start MCP server with just TFLint
ship mcp lint
```

**2. Multiple Tools MCP Server:**
```bash
# Start MCP server with Checkov and TFLint
ship mcp checkov lint
```

**3. All Tools MCP Server:**
```bash
# Start MCP server with all Ship tools
ship mcp all
```

#### Integration with AI Assistants

Configure your custom MCP server in Claude Code:

```json
{
  "mcpServers": {
    "ship-terraform": {
      "command": "ship",
      "args": ["mcp", "all"]
    }
  }
}
```

📚 **Learn More**: Check out the [Ship Documentation](docs/) for complete guides and tool reference.

### 2. AI Assistant Integration (MCP)

Configure your custom MCP servers built with Ship in Claude Desktop or other AI assistants:

```json
{
  "mcpServers": {
    "my-infrastructure-server": {
      "command": "/path/to/my-mcp-server",
      "env": {
        "AWS_PROFILE": "your-profile"
      }
    }
  }
}
```

**Ship Framework wraps mcp-go** to provide:
- **Container-based Tools**: All tools run securely in Docker containers
- **Fluent Builder API**: Easy-to-use APIs for building MCP servers
- **Pre-built Ship Tools**: Infrastructure tools ready to use in your servers
- **Three Usage Patterns**: Pure framework, cherry-pick tools, or everything plus extensions

### 3. Optional CLI Usage

Ship also includes a CLI for direct usage of infrastructure tools:

```bash
# Navigate to your Terraform project
cd your-terraform-project

# Run analysis tools
ship tf lint                # TFLint for syntax and best practices
ship tf checkov             # Security scanning
ship tf cost                # Cost estimation
ship tf docs                # Generate documentation
ship tf diagram . --hcl -o infrastructure.png  # Generate diagrams
```

## 🛠️ Available Infrastructure Tools

| Tool | Ship SDK | Description | Container Image |
|------|----------|-------------|-----------------|
| **TFLint** | `tools.NewTFLintTool()` | Terraform linter for syntax and best practices | `ghcr.io/terraform-linters/tflint` |
| **Checkov** | `tools.NewCheckovTool()` | Security and compliance scanner | `bridgecrew/checkov` |
| **Trivy** | `tools.NewTrivyTool()` | Vulnerability scanner for IaC | `aquasec/trivy` |
| **OpenInfraQuote** | `tools.NewCostTool()` | Infrastructure cost analysis | `gruebel/openinfraquote` |
| **terraform-docs** | `tools.NewDocsTool()` | Auto-generate module documentation | `quay.io/terraform-docs/terraform-docs` |
| **InfraMap** | `tools.NewDiagramTool()` | Infrastructure diagram generation | `cycloid/inframap:latest` |

## 🔧 Ship Framework Integration

### mcp-go Integration

Ship Framework wraps the official [mcp-go](https://github.com/modelcontextprotocol/mcp-go) library to provide:

- **Enhanced Container Tools**: All tools run in isolated Docker containers via Dagger
- **Fluent Builder API**: Easy-to-use APIs for building MCP servers
- **Security by Default**: No local tool installations required
- **Pre-built Infrastructure Tools**: Ready-to-use Terraform analysis tools

### Example: Wrapping mcp-go

```go
package main

import (
    "context"
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/pkg/tools"
    "github.com/modelcontextprotocol/mcp-go/mcp"
)

func main() {
    // Ship builds on mcp-go's foundation
    server := ship.NewServer("infrastructure-server", "1.0.0").
        AddTool(tools.NewTFLintTool()).         // Pre-built Ship tool
        AddTool(tools.NewCheckovTool()).        // Security scanning
        AddContainerTool("custom", config).     // Your custom tools
        Build()
    
    // Leverages mcp-go's stdio transport
    server.ServeStdio()
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
server.AddTool(tools.NewTFLintTool().WithEnv(map[string]string{
    "AWS_PROFILE": "production",
    "AWS_REGION": "us-west-2",
}))

// Azure credentials
server.AddTool(tools.NewCheckovTool().WithEnv(map[string]string{
    "ARM_CLIENT_ID": "your-client-id",
    "ARM_CLIENT_SECRET": "your-client-secret",
}))
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
- [Dynamic Module Discovery](docs/dynamic-module-discovery.md) - Extensible module system
- [Dagger Modules](docs/dagger-modules.md) - How to add new tools
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

## 🌐 Extensibility

Ship Framework is designed to be extensible:

### Custom Container Tools

Easily add any containerized tool to your MCP server:

```go
// Add a custom security scanner
customTool := ship.NewContainerTool("security-scan", ship.ContainerToolConfig{
    Description: "Custom security analysis",
    Image:       "my-org/security-scanner:latest",
    Parameters: []ship.Parameter{
        {Name: "directory", Type: "string", Description: "Directory to scan", Required: true},
        {Name: "severity", Type: "string", Description: "Minimum severity level"},
    },
    Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
        dir := params["directory"].(string)
        severity := params["severity"].(string)
        
        result, err := engine.Container().
            From("my-org/security-scanner:latest").
            WithWorkdir("/workspace").
            WithExec([]string{"scan", "--dir", dir, "--severity", severity}).
            Stdout(ctx)
            
        return &ship.ToolResult{Content: result}, err
    },
})

server := ship.NewServer("security-server", "1.0.0").
    AddTool(customTool).
    Build()
```

### Community Ideas

- **Cloud Security Scanners**: Deep analysis for AWS/Azure/GCP
- **Kubernetes Tools**: K8s manifest validation and cluster analysis  
- **Database Tools**: Schema validation, migration checks
- **Compliance Checkers**: SOC2, HIPAA, PCI-DSS validators
- **Custom Cost Analyzers**: Organization-specific cost allocation

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