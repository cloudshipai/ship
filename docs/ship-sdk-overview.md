# Ship SDK Overview

The Ship SDK is a comprehensive framework for building Model Context Protocol (MCP) servers that run tools securely in containers using Dagger. It's designed to provide maximum flexibility while maintaining security and reliability.

## What is the Ship SDK?

The Ship SDK enables developers to:

- **Build Custom MCP Servers**: Create MCP servers with custom tools for AI assistants like Claude, Cursor, and others
- **Run Tools Securely**: All tools execute in isolated containers via Dagger, ensuring security and consistency
- **Leverage Ship Tools**: Use the curated collection of CloudShip AI team tools or build your own
- **Flexible Architecture**: Choose from three usage patterns based on your needs

## Core Components

### Tools (`pkg/ship/tool.go`)
The foundation of the SDK - tools that can be executed by AI assistants:
- **Container Tools**: Run commands in Docker containers via Dagger
- **Custom Tools**: Implement your own tool logic
- **Ship Tools**: Pre-built tools for Terraform, security scanning, documentation generation

### Registry (`pkg/ship/registry.go`)
Centralized management of tools, prompts, and resources:
- **Thread-safe**: Concurrent access support
- **Flexible Registration**: Add tools individually or import entire registries
- **Discovery**: List and retrieve tools dynamically

### Server (`pkg/ship/server.go`)
MCP server builder with fluent API:
- **Builder Pattern**: Chainable methods for server configuration
- **MCP Integration**: Automatic integration with Model Context Protocol
- **Engine Management**: Handles Dagger engine lifecycle

## Three Usage Patterns

### 1. Pure Ship SDK (Custom Tools Only)
Build completely custom MCP servers with no Ship tools:

```go
server := ship.NewServer("my-server", "1.0.0").
    AddContainerTool("my-tool", ship.ContainerToolConfig{
        Description: "My custom tool",
        Image: "alpine:latest",
        Parameters: []ship.Parameter{
            {Name: "input", Type: "string", Required: true},
        },
        Execute: myToolExecutor,
    }).
    Build()
```

### 2. Cherry-Pick Ship Tools
Select specific Ship tools and combine with custom tools:

```go
server := ship.NewServer("hybrid-server", "1.0.0").
    AddTool(tools.NewTFLintTool()).
    AddTool(tools.NewCheckovTool()).
    AddContainerTool("my-custom-tool", config).
    Build()
```

### 3. Everything Plus
Use all Ship tools and extend with custom tools:

```go
server := all.AddAllTools(
    ship.NewServer("full-server", "1.0.0").
    AddContainerTool("my-extension", config),
).Build()
```

## Why Choose Ship SDK?

### Security First
- **Container Isolation**: All tools run in isolated Docker containers
- **No Local Dependencies**: Tools don't require local installation
- **Dagger Integration**: Leverages Dagger's security model

### Developer Experience
- **Fluent API**: Chainable, intuitive builder pattern
- **Type Safety**: Full Go type safety with comprehensive interfaces
- **Rich Examples**: Complete examples for all usage patterns

### Production Ready
- **Thread Safety**: Concurrent execution support
- **Error Handling**: Comprehensive error handling and reporting
- **Testing**: Full test coverage with integration tests

## Getting Started

See the [Quick Start Guide](ship-sdk-quickstart.md) for step-by-step instructions to build your first MCP server with the Ship SDK.