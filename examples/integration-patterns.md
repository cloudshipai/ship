# Ship Integration Patterns

This document shows different ways to integrate Ship's containerized tools with existing mcp-go applications.

## Pattern 1: Ship-First (Current Default)

The simplest pattern where Ship manages everything:

```go
package main

import (
    "log"
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/internal/tools"
)

func main() {
    // Ship manages the entire MCP server
    server := ship.NewServer("my-infrastructure-server", "1.0.0").
        AddTool(tools.NewTFLintTool()).
        AddTool(tools.NewCheckovTool()).
        Build()

    if err := server.ServeStdio(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

## Pattern 2: Bring Your Own MCP Server

For users who already have an mcp-go server and want to add Ship tools:

```go
package main

import (
    "context"
    "log"
    
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/internal/tools"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    ctx := context.Background()
    
    // Your existing mcp-go server setup
    mcpServer := server.NewMCPServer("my-app", "1.0.0")
    
    // Add your regular MCP tools
    mcpServer.AddTool(mcp.NewTool("my-custom-tool"), myCustomHandler)
    mcpServer.AddTool(mcp.NewTool("another-tool"), anotherHandler)
    
    // Create Ship adapter and add containerized infrastructure tools
    shipAdapter := ship.NewMCPAdapter().
        AddTool(tools.NewTFLintTool()).
        AddTool(tools.NewCheckovTool()).
        AddContainerTool("custom-security", ship.ContainerToolConfig{
            Description: "Custom security scanner",
            Image:       "my-org/security:latest",
            // ... config
        })
    
    // Attach Ship tools to your existing MCP server
    if err := shipAdapter.AttachToServer(ctx, mcpServer); err != nil {
        log.Fatalf("Failed to attach Ship tools: %v", err)
    }
    
    defer shipAdapter.Close()
    
    // Serve using your preferred method
    if err := server.ServeStdio(mcpServer); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}

func myCustomHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Your existing tool logic
    return mcp.NewToolResultText("Hello from custom tool"), nil
}

func anotherHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // More custom logic
    return mcp.NewToolResultText("Another tool result"), nil
}
```

## Pattern 3: Tool Routing (Advanced)

For complex setups where you want different tools routed to different MCP servers:

```go
package main

import (
    "context"
    "log"
    
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/internal/tools"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    ctx := context.Background()
    
    // Create separate MCP servers for different concerns
    securityServer := server.NewMCPServer("security-tools", "1.0.0")
    infrastructureServer := server.NewMCPServer("infra-tools", "1.0.0")
    docsServer := server.NewMCPServer("docs-tools", "1.0.0")
    
    // Create Ship adapters for different tool categories
    securityAdapter := ship.NewMCPAdapter().
        AddTool(tools.NewCheckovTool()).
        AddTool(tools.NewTrivyTool())
        
    infraAdapter := ship.NewMCPAdapter().
        AddTool(tools.NewTFLintTool()).
        AddTool(tools.NewCostTool())
        
    docsAdapter := ship.NewMCPAdapter().
        AddTool(tools.NewDocsTool()).
        AddTool(tools.NewDiagramTool())
    
    // Create tool router
    router := ship.NewToolRouter().
        AddRoute("security", securityAdapter).
        AddRoute("infrastructure", infraAdapter).
        AddRoute("docs", docsAdapter)
    
    // Route tools to their respective servers
    routes := map[string]*server.MCPServer{
        "security":      securityServer,
        "infrastructure": infrastructureServer,
        "docs":          docsServer,
    }
    
    if err := router.RouteToServer(ctx, routes); err != nil {
        log.Fatalf("Failed to route tools: %v", err)
    }
    
    // In a real app, you'd serve these on different ports or using different transports
    // This is just an example
    go func() { server.ServeStdio(securityServer) }()
    go func() { server.ServeStdio(infrastructureServer) }()
    server.ServeStdio(docsServer) // Main thread
}
```

## Pattern 4: Selective Tool Integration

For users who only want specific Ship capabilities:

```go
package main

import (
    "context"
    "log"
    
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/pkg/dagger"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    ctx := context.Background()
    mcpServer := server.NewMCPServer("selective-app", "1.0.0")
    
    // Your existing tools
    mcpServer.AddTool(mcp.NewTool("existing-tool"), existingHandler)
    
    // Only use Ship's Dagger engine and container capabilities
    engine, err := dagger.NewEngine(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()
    
    // Create a custom containerized tool using Ship's framework
    customTool := ship.NewContainerTool("my-linter", ship.ContainerToolConfig{
        Description: "My custom linter",
        Image:       "my-org/linter:latest",
        Parameters: []ship.Parameter{
            {Name: "file", Type: "string", Description: "File to lint", Required: true},
        },
        Execute: func(ctx context.Context, params map[string]interface{}, eng *dagger.Engine) (*ship.ToolResult, error) {
            file := params["file"].(string)
            result, err := eng.Container().
                From("my-org/linter:latest").
                WithWorkdir("/workspace").
                WithExec([]string{"lint", file}).
                Stdout(ctx)
            
            return &ship.ToolResult{Content: result}, err
        },
    })
    
    // Add just this one Ship tool to your MCP server
    adapter := ship.NewMCPAdapter().WithEngine(engine).AddTool(customTool)
    if err := adapter.AttachToServer(ctx, mcpServer); err != nil {
        log.Fatal(err)
    }
    
    if err := server.ServeStdio(mcpServer); err != nil {
        log.Fatal(err)
    }
}

func existingHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    return mcp.NewToolResultText("Existing functionality"), nil
}
```

## Benefits of These Patterns

1. **Ship-First**: Fastest to get started, Ship handles everything
2. **Bring Your Own**: Perfect for existing codebases, gradual migration
3. **Tool Routing**: Microservices-style architecture for complex deployments  
4. **Selective**: Cherry-pick only the Ship capabilities you need

## Key Advantages

- **Zero Vendor Lock-in**: Use as much or as little of Ship as you want
- **Gradual Adoption**: Start with one tool, expand over time
- **Framework Agnostic**: Works with any mcp-go setup
- **Container Benefits**: Get Dagger's security and isolation regardless of pattern
- **Composable**: Mix and match Ship tools with your existing tools

Choose the pattern that best fits your existing architecture and migration path!