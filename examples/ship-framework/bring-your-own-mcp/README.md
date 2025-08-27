# Bring Your Own MCP Server Example

This example demonstrates how to integrate Ship Framework tools into your **existing mcp-go server**. This is perfect when you already have an MCP server and want to add Ship's containerized tools to it.

## Architecture

```
Your Existing MCP Server
├── Your Native MCP Tools (pure mcp-go)
└── Ship Tools (via MCPAdapter)
    ├── Container-based tools
    ├── Dagger integration
    └── Ship tool interfaces
```

## Key Components

### 1. Your Existing MCP Server
```go
// Create your own mcp-go server as usual
mcpServer := server.NewMCPServer("my-server", "1.0.0")

// Add your existing native tools
addCustomMCPTool(mcpServer)
```

### 2. Ship MCPAdapter
```go
// Create adapter and add Ship tools
adapter := ship.NewMCPAdapter().
    AddTool(&CustomTool{})

// Attach Ship tools to your existing server
adapter.AttachToServer(ctx, mcpServer)
```

### 3. Mixed Tool Types
- **Native MCP Tools**: Direct mcp-go implementations
- **Ship Framework Tools**: Containerized tools using Ship's interfaces

## Running the Example

```bash
# Install dependencies  
go mod tidy

# Start the mixed MCP server
go run main.go
```

## Testing with Claude Desktop

```json
{
  "mcpServers": {
    "mixed-tools": {
      "command": "go",
      "args": ["run", "main.go"], 
      "cwd": "/path/to/ship/examples/ship-framework/bring-your-own-mcp"
    }
  }
}
```

Then ask Claude to:
- "Use the native-tool to process 'Hello World'"
- "Use custom-hello to greet 'Alice'"

## Use Cases

### 1. Migration Strategy
```go
// Gradually migrate existing tools to Ship framework
adapter := ship.NewMCPAdapter()
    .AddTool(&LegacyToolWrapper{})  // Wrap existing logic
    .AddTool(&NewContainerTool{})   // New containerized tools
```

### 2. Specialized Integration
```go
// Add specific Ship tools to domain-specific servers
adapter := ship.NewMCPAdapter()
    .AddTool(&TerraformValidateTool{})
    .AddTool(&DockerLintTool{})
```

### 3. Tool Router Pattern
```go
// Route different tool types to different servers
router := ship.NewToolRouter().
    AddRoute("terraform", terraformAdapter).
    AddRoute("security", securityAdapter)

router.RouteToServer(ctx, map[string]*server.MCPServer{
    "terraform": terraformServer,
    "security":  securityServer,
})
```

## Benefits of MCPAdapter

1. **Incremental Adoption**: Add Ship tools to existing servers
2. **Tool Mixing**: Combine native and containerized tools
3. **Container Benefits**: Get Dagger's isolation and caching
4. **Ship Ecosystem**: Access to Ship's tool interfaces and utilities

## Advanced Patterns

### Multiple Adapters
```go
// Different adapters for different tool categories
securityAdapter := ship.NewMCPAdapter().AddTool(&TrivyTool{})
infraAdapter := ship.NewMCPAdapter().AddTool(&TerraformTool{})

securityAdapter.AttachToServer(ctx, securityServer)
infraAdapter.AttachToServer(ctx, infraServer)
```

### Conditional Tool Loading
```go
adapter := ship.NewMCPAdapter()

if os.Getenv("ENABLE_DOCKER") == "true" {
    adapter.AddTool(&DockerLintTool{})
}

if os.Getenv("ENABLE_TERRAFORM") == "true" {
    adapter.AddTool(&TerraformValidateTool{})
}
```

This pattern allows you to **keep your existing MCP server architecture** while adding Ship's containerized tool capabilities where needed.