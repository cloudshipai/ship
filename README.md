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
- **🤖 AI Development**: OpenCode AI coding assistant with chat, generation, and analysis
- **📈 Monitoring**: Grafana integration for dashboards, alerts, and metrics analysis
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

### 1. Ship Framework

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

#### Advanced Integration Patterns

**Bring Your Own MCP Server:**
Perfect for existing applications that already use mcp-go - just add Ship's containerized tools:

```go
import (
    "context"
    "log"
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/internal/tools" 
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    ctx := context.Background()
    
    // Your existing mcp-go server
    mcpServer := server.NewMCPServer("my-app", "1.0.0")
    
    // Add your existing tools
    mcpServer.AddTool(myCustomTool, myHandler)
    
    // Add Ship's containerized infrastructure tools
    shipAdapter := ship.NewMCPAdapter().
        AddTool(tools.NewTFLintTool())
    
    // Attach Ship tools to your existing server
    if err := shipAdapter.AttachToServer(ctx, mcpServer); err != nil {
        log.Fatalf("Failed to attach Ship tools: %v", err)
    }
    defer shipAdapter.Close()
    
    // Now you have both your tools AND Ship's containerized tools
    server.ServeStdio(mcpServer)
}
```

**Selective Integration:**
Only use the Ship capabilities you need:

```go
// Just use Ship's container framework with custom tools
customTool := ship.NewContainerTool("my-scanner", ship.ContainerToolConfig{
    Description: "Custom security scanner",
    Image: "my-org/scanner:latest",
    Parameters: []ship.Parameter{
        {Name: "directory", Type: "string", Description: "Directory to scan", Required: true},
    },
    Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
        // Your custom tool logic here
        return &ship.ToolResult{Content: "Scan completed"}, nil
    },
})

adapter := ship.NewMCPAdapter().AddTool(customTool)
adapter.AttachToServer(ctx, yourExistingMCPServer)
```

📚 **Complete Integration Guide**: See [examples/integration-patterns.md](examples/integration-patterns.md) for detailed integration patterns with code examples for all four usage patterns.

#### Integration with AI Assistants

Configure Ship MCP servers in Claude Code:

```json
{
  "mcpServers": {
    "ship-terraform": {
      "command": "ship",
      "args": ["mcp", "all"]
    },
    "ship-filesystem": {
      "command": "ship",
      "args": ["mcp", "filesystem"],
      "env": {
        "FILESYSTEM_ROOT": "/workspace"
      }
    },
    "ship-search": {
      "command": "ship", 
      "args": ["mcp", "brave-search", "--var", "BRAVE_API_KEY=your_key"]
    }
  }
}
```

#### Ship Framework Mode with External MCP Servers

Use external MCP servers in your Ship framework applications:

```go
package main

import (
    "context"
    "log"
    "github.com/cloudshipai/ship/pkg/ship"
)

func main() {
    ctx := context.Background()
    
    // Create Ship server with external MCP server integration
    server := ship.NewServer("my-app", "1.0.0")
    
    // Add built-in Ship tools
    server.AddTool(tools.NewTFLintTool())
    
    // Add external MCP servers as proxy tools
    filesystemConfig := ship.MCPServerConfig{
        Name:      "filesystem",
        Command:   "npx",
        Args:      []string{"-y", "@modelcontextprotocol/server-filesystem", "/workspace"},
        Transport: "stdio",
        Env: map[string]string{
            "FILESYSTEM_ROOT": "/workspace",
        },
    }
    
    // Create proxy for external MCP server
    proxy := ship.NewMCPProxy(filesystemConfig)
    if err := proxy.Connect(ctx); err != nil {
        log.Fatalf("Failed to connect to filesystem MCP server: %v", err)
    }
    defer proxy.Close()
    
    // Discover and add external tools
    externalTools, err := proxy.DiscoverTools(ctx)
    if err != nil {
        log.Fatalf("Failed to discover external tools: %v", err)
    }
    
    for _, tool := range externalTools {
        server.AddTool(tool)
    }
    
    // Build and start server
    mcpServer := server.Build()
    if err := mcpServer.ServeStdio(); err != nil {
        log.Fatalf("Server failed: %v", err)
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

# Start MCP servers for AI assistant integration
ship mcp all                # All Ship tools
ship mcp lint               # Just TFLint
ship mcp filesystem         # External filesystem MCP server
ship mcp brave-search --var BRAVE_API_KEY=your_api_key  # External search with API key
ship mcp grafana --var GRAFANA_URL=http://localhost:3000 --var GRAFANA_API_KEY=glsa_xyz  # Grafana monitoring

# AI Development Tools
ship opencode chat "explain this terraform module" --model "openai/gpt-4o-mini"
ship opencode generate "create a kubernetes deployment" --output k8s/deployment.yaml
```

## 🛠️ Available Infrastructure Tools

### Built-in Ship Tools

| Tool | Ship Framework | Description | Container Image |
|------|----------------|-------------|-----------------|
| **TFLint** | `tools.NewTFLintTool()` | Terraform linter for syntax and best practices | `ghcr.io/terraform-linters/tflint` |
| **Checkov** | `tools.NewCheckovTool()` | Security scanning for Terraform | `bridgecrew/checkov` |
| **Trivy** | `tools.NewTrivyTool()` | Security scanning for Terraform | `aquasec/trivy` |
| **OpenInfraQuote** | `tools.NewCostTool()` | Cost analysis for infrastructure | `cloudshipai/openinfraquote` |
| **terraform-docs** | `tools.NewDocsTool()` | Documentation generation | `quay.io/terraform-docs/terraform-docs` |
| **InfraMap** | `tools.NewDiagramTool()` | Infrastructure diagrams | `cycloidio/inframap` |
| **OpenCode** | `ship opencode chat` | AI coding assistant with chat, generation, analysis | `node:18` + `opencode-ai` |

### External MCP Servers

Ship can proxy external MCP servers, discovering their tools dynamically:

| Server | Description | Variables | Example Usage |
|--------|-------------|-----------|---------------|
| **filesystem** | File and directory operations | `FILESYSTEM_ROOT` (optional) | `ship mcp filesystem --var FILESYSTEM_ROOT=/custom/path` |
| **memory** | Persistent knowledge storage | `MEMORY_STORAGE_PATH`, `MEMORY_MAX_SIZE` (optional) | `ship mcp memory --var MEMORY_STORAGE_PATH=/data` |
| **brave-search** | Web search capabilities | `BRAVE_API_KEY` (required), `BRAVE_SEARCH_COUNT` (optional) | `ship mcp brave-search --var BRAVE_API_KEY=your_key` |
| **grafana** | Grafana monitoring and visualization | `GRAFANA_URL` (required), `GRAFANA_API_KEY` (optional), `GRAFANA_USERNAME`/`GRAFANA_PASSWORD` (optional) | `ship mcp grafana --var GRAFANA_URL=http://localhost:3000 --var GRAFANA_API_KEY=glsa_xyz` |

> **Note**: External MCP servers are automatically installed via npm when needed. Tools are discovered dynamically at runtime.

## ⚙️ Environment Variables and Configuration

### --var Flag System

Ship supports passing environment variables to both containerized tools and external MCP servers using the `--var` flag:

```bash
# Single variable
ship mcp brave-search --var BRAVE_API_KEY=your_api_key

# Multiple variables
ship mcp memory --var MEMORY_STORAGE_PATH=/data --var MEMORY_MAX_SIZE=100MB

# Variables for containerized tools
ship mcp cost --var AWS_REGION=us-east-1 --var DEBUG=true
```

### Variable Types

**Framework-Defined Variables**: Each tool and external MCP server defines its own variables with:
- **Required vs Optional**: Some variables are mandatory, others have defaults
- **Default Values**: Optional variables often have sensible defaults
- **Secret Handling**: API keys and sensitive data are marked as secrets
- **Validation**: Variables are validated before starting tools

**Examples by Tool**:

```bash
# Filesystem operations (all optional)
ship mcp filesystem --var FILESYSTEM_ROOT=/custom/path

# Memory storage (all optional) 
ship mcp memory --var MEMORY_STORAGE_PATH=/data --var MEMORY_MAX_SIZE=100MB

# Brave search (API key required)
ship mcp brave-search --var BRAVE_API_KEY=your_key --var BRAVE_SEARCH_COUNT=20

# Grafana monitoring (URL required, auth optional)
ship mcp grafana --var GRAFANA_URL=http://localhost:3000 --var GRAFANA_API_KEY=glsa_xyz

# OpenCode AI assistant (AI provider API keys via environment)
ship opencode chat "analyze this code" --model "openai/gpt-4o-mini"
export OPENAI_API_KEY=sk-...  # Set before running OpenCode

# Containerized tools (any environment variable)
ship mcp all --var AWS_PROFILE=production --var AWS_REGION=us-west-2
```

### Variable Discovery

Use `ship modules info <tool-name>` to see available variables:

```bash
# See variables for external MCP servers
ship modules info filesystem
ship modules info memory
ship modules info brave-search
ship modules info grafana

# See information about built-in tools
ship modules info lint
ship modules info cost
```

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