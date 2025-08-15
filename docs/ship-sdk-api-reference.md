# Ship SDK API Reference

This document provides detailed API documentation for the Ship SDK components.

## Core Interfaces

### Tool Interface

The fundamental interface that all tools must implement:

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() []Parameter
    Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error)
}
```

**Methods:**
- `Name()`: Returns the unique tool identifier
- `Description()`: Returns human-readable tool description
- `Parameters()`: Returns parameter definitions for the tool
- `Execute()`: Executes the tool with given parameters

### Parameter

Defines tool input parameters:

```go
type Parameter struct {
    Name        string      `json:"name"`
    Type        string      `json:"type"`
    Description string      `json:"description"`
    Required    bool        `json:"required"`
    Enum        []string    `json:"enum,omitempty"`
    Default     interface{} `json:"default,omitempty"`
}
```

**Fields:**
- `Name`: Parameter identifier
- `Type`: Data type (`string`, `boolean`, `number`, `object`, `array`)
- `Description`: Human-readable description
- `Required`: Whether parameter is mandatory
- `Enum`: Allowed values (optional)
- `Default`: Default value (optional)

### ToolResult

Represents tool execution results:

```go
type ToolResult struct {
    Content  string                 `json:"content"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
    Error    error                  `json:"error,omitempty"`
}
```

**Fields:**
- `Content`: Primary tool output
- `Metadata`: Additional structured data
- `Error`: Execution error (if any)

## Container Tools

### ContainerTool

Pre-built implementation for container-based tools:

```go
func NewContainerTool(name string, config ContainerToolConfig) Tool
```

### ContainerToolConfig

Configuration for container tools:

```go
type ContainerToolConfig struct {
    Description string
    Image       string
    Parameters  []Parameter
    Execute     func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error)
}
```

**Fields:**
- `Description`: Tool description
- `Image`: Docker image to use
- `Parameters`: Parameter definitions
- `Execute`: Execution function

## Registry

### Registry Interface

Manages tool collections:

```go
type Registry interface {
    RegisterTool(tool Tool) error
    GetTool(name string) (Tool, error)
    ListTools() []string
    ToolCount() int
    ImportRegistry(other Registry) error
}
```

**Methods:**
- `RegisterTool()`: Add a tool to the registry
- `GetTool()`: Retrieve tool by name
- `ListTools()`: List all tool names
- `ToolCount()`: Get total tool count
- `ImportRegistry()`: Import tools from another registry

### Creating Registries

```go
// Create new registry
registry := ship.NewRegistry()

// Global registry for Ship tools
defaultRegistry := ship.DefaultRegistry
```

## Server Builder

### ServerBuilder

Fluent API for building MCP servers:

```go
type ServerBuilder struct {
    // private fields
}
```

### Creating Servers

```go
// Create new server builder
builder := ship.NewServer(name, version string) *ServerBuilder
```

### Builder Methods

**Tool Management:**
```go
AddTool(tool Tool) *ServerBuilder
AddContainerTool(name string, config ContainerToolConfig) *ServerBuilder
ImportRegistry(registry Registry) *ServerBuilder
```

**Configuration:**
```go
WithEngine(engine *dagger.Engine) *ServerBuilder
WithExecutor(executor ToolExecutor) *ServerBuilder
```

**Building:**
```go
Build() *Server
```

## Server

### Server Interface

The built MCP server:

```go
type Server interface {
    Start(ctx context.Context) error
    Close() error
    ServeStdio() error
    GetRegistry() Registry
    GetEngine() *dagger.Engine
}
```

**Methods:**
- `Start()`: Initialize server resources
- `Close()`: Clean up server resources
- `ServeStdio()`: Start MCP server over stdin/stdout
- `GetRegistry()`: Access tool registry
- `GetEngine()`: Access Dagger engine

## Convenience Functions

### Tool Collections

The `pkg/tools/all` package provides convenience functions:

```go
// Registry functions
func TerraformRegistry() Registry
func SecurityRegistry() Registry
func DocsRegistry() Registry
func AllRegistry() Registry

// Builder functions
func AddTerraformTools(builder *ServerBuilder) *ServerBuilder
func AddSecurityTools(builder *ServerBuilder) *ServerBuilder
func AddDocsTools(builder *ServerBuilder) *ServerBuilder
func AddAllTools(builder *ServerBuilder) *ServerBuilder
```

## Error Handling

### Common Error Types

- `ErrToolNotFound`: Tool not found in registry
- `ErrToolAlreadyExists`: Tool already registered
- `ErrInvalidParameter`: Invalid parameter value
- `ErrExecutionFailed`: Tool execution failed

### Error Handling Patterns

```go
// Tool execution with error handling
result, err := tool.Execute(ctx, params, engine)
if err != nil {
    // Handle execution error
    log.Printf("Tool execution failed: %v", err)
    return err
}

// Check result for tool-specific errors
if result.Error != nil {
    // Handle tool-specific error
    log.Printf("Tool reported error: %v", result.Error)
}
```

## Best Practices

### Tool Implementation

1. **Parameter Validation**: Always validate input parameters
2. **Error Handling**: Provide meaningful error messages
3. **Container Security**: Use minimal, specific container images
4. **Resource Cleanup**: Ensure proper cleanup of resources

### Server Configuration

1. **Tool Naming**: Use clear, descriptive tool names
2. **Registry Organization**: Group related tools in registries
3. **Engine Management**: Reuse Dagger engines when possible
4. **Lifecycle Management**: Properly start and close servers

### Performance

1. **Container Reuse**: Leverage Dagger's container caching
2. **Parallel Execution**: Design tools for concurrent execution
3. **Resource Limits**: Set appropriate container resource limits
4. **Caching**: Use Dagger's caching for repeated operations

## Examples

See the [examples/](../examples/) directory for complete working examples of all API features.