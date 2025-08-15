# Ship SDK Quick Start Guide

This guide will help you build your first MCP server using the Ship SDK in just a few minutes.

## Prerequisites

- Go 1.21 or later
- Docker (for Dagger containers)
- Basic understanding of Go and MCP concepts

## Installation

Add the Ship SDK to your Go project:

```bash
go mod init my-mcp-server
go get github.com/cloudshipai/ship/pkg/ship
go get github.com/cloudshipai/ship/pkg/tools  # Optional: for Ship tools
```

## Quick Start: Echo Server

Let's build a simple MCP server with an echo tool:

### 1. Create main.go

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
            {
                Name:        "message",
                Type:        "string",
                Description: "Message to echo",
                Required:    true,
            },
        },
        Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
            message := params["message"].(string)
            
            result, err := engine.Container().
                From("alpine:latest").
                WithExec([]string{"echo", message}).
                Stdout(ctx)
                
            if err != nil {
                return &ship.ToolResult{Error: err}, err
            }
            
            return &ship.ToolResult{
                Content: result,
                Metadata: map[string]interface{}{
                    "tool": "echo",
                    "message": message,
                },
            }, nil
        },
    })

    // Build and start the server
    server := ship.NewServer("echo-server", "1.0.0").
        AddTool(echoTool).
        Build()

    // Start serving over stdio (MCP protocol)
    if err := server.ServeStdio(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

### 2. Build and Run

```bash
go build -o echo-server .
./echo-server
```

Your MCP server is now running and ready to be used by AI assistants!

## Using Ship Tools

To use pre-built Ship tools, import the tools package:

```go
package main

import (
    "log"
    
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/pkg/tools"
    "github.com/cloudshipai/ship/pkg/tools/all"
)

func main() {
    // Option 1: Add specific Ship tools
    server := ship.NewServer("terraform-server", "1.0.0").
        AddTool(tools.NewTFLintTool()).
        Build()

    // Option 2: Use convenience functions
    server = all.AddTerraformTools(
        ship.NewServer("terraform-server", "1.0.0"),
    ).Build()

    // Option 3: Add all Ship tools
    server = all.AddAllTools(
        ship.NewServer("full-server", "1.0.0"),
    ).Build()

    if err := server.ServeStdio(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

## Integration with AI Assistants

### Claude Code

Add your MCP server to Claude Code's configuration:

```json
{
  "mcpServers": {
    "my-server": {
      "command": "/path/to/my-mcp-server"
    }
  }
}
```

### Cursor

Configure in Cursor's MCP settings:

```json
{
  "mcpServers": {
    "my-server": {
      "command": "/path/to/my-mcp-server",
      "args": []
    }
  }
}
```

## Next Steps

- Read the [API Reference](ship-sdk-api-reference.md) for detailed documentation
- Check out [Usage Patterns](ship-sdk-usage-patterns.md) for advanced patterns
- Browse [Examples](../examples/) for complete working examples
- See [Ship Tools](ship-tools-reference.md) for available pre-built tools

## Need Help?

- Check the [FAQ](ship-sdk-faq.md)
- Review [examples/](../examples/) for working code
- Open an issue on GitHub for bugs or feature requests