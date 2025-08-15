# Ship SDK FAQ

## General Questions

### What is the Ship SDK?

The Ship SDK is a framework for building Model Context Protocol (MCP) servers that run tools securely in containers using Dagger. It's developed by the CloudShip AI team and provides both pre-built tools and a framework for custom tool development.

### How is this different from other MCP servers?

Key differences:
- **Container Security**: All tools run in isolated Docker containers
- **Dagger Integration**: Leverages Dagger for consistent, reproducible builds
- **Curated Tools**: Includes professionally maintained infrastructure tools
- **Three Usage Patterns**: Pure framework, cherry-pick tools, or everything plus extensions

### Do I need Docker installed?

Yes, Docker is required because Ship SDK uses Dagger, which orchestrates Docker containers. However, tools don't require local installation - they run in containers.

### What's the relationship to CloudShip AI?

The Ship SDK is developed by CloudShip AI as an open-source framework. It includes tools that work with CloudShip's platform but can be used independently.

## Usage Questions

### Which usage pattern should I choose?

- **Pure Ship SDK**: Choose this if you need completely custom tools with no CloudShip dependencies
- **Cherry-Pick**: Choose this if you want some Ship tools plus your own custom tools
- **Everything Plus**: Choose this if you want all Ship tools and plan to extend with custom tools

### Can I mix usage patterns?

Yes! You can create multiple servers with different patterns or combine patterns within a single server using registries.

### How do I know which tools are available?

Check the [Ship Tools Reference](ship-tools-reference.md) for current tools. You can also list available tools programmatically:

```go
registry := all.AllRegistry()
tools := registry.ListTools()
fmt.Printf("Available tools: %v\n", tools)
```

## Technical Questions

### How do I handle tool failures?

Ship SDK provides two levels of error handling:

```go
result, err := tool.Execute(ctx, params, engine)
if err != nil {
    // Container execution failed - serious error
    log.Printf("Execution failed: %v", err)
    return err
}

if result.Error != nil {
    // Tool reported issues (e.g., lint violations) - may be expected
    log.Printf("Tool found issues: %v", result.Error)
}
```

### Can I run tools in parallel?

Yes! Tools are designed for concurrent execution:

```go
var wg sync.WaitGroup
results := make(chan *ship.ToolResult, len(tools))

for _, tool := range tools {
    wg.Add(1)
    go func(t ship.Tool) {
        defer wg.Done()
        result, err := t.Execute(ctx, params, engine)
        if err == nil {
            results <- result
        }
    }(tool)
}

wg.Wait()
close(results)
```

### How do I debug tool execution?

1. **Check container logs**: Dagger provides detailed container execution logs
2. **Use metadata**: Tools provide rich metadata about execution
3. **Test locally**: Run container images directly with Docker
4. **Enable debug mode**: Some tools support verbose output

### Can I use my own container images?

Absolutely! Create custom tools with any container image:

```go
customTool := ship.NewContainerTool("my-tool", ship.ContainerToolConfig{
    Description: "My custom tool",
    Image:       "my-registry/my-tool:v1.0.0",
    Parameters: []ship.Parameter{
        {Name: "input", Type: "string", Required: true},
    },
    Execute: myExecuteFunction,
})
```

### How do I handle secrets and credentials?

Use Dagger's secret management:

```go
Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
    apiKey := params["api_key"].(string)
    
    result, err := engine.Container().
        From("my-tool:latest").
        WithSecretVariable("API_KEY", engine.Secret(apiKey)).
        WithExec([]string{"./my-tool"}).
        CombinedOutput(ctx)
        
    return &ship.ToolResult{Content: result}, err
}
```

## Development Questions

### How do I contribute a new Ship tool?

1. Create the tool implementation in `pkg/tools/`
2. Add comprehensive tests
3. Update appropriate collection in `pkg/tools/all/`
4. Submit a pull request with documentation

See the [Ship Tools Reference](ship-tools-reference.md) for detailed guidelines.

### Can I modify existing Ship tools?

Yes, wrap them with custom logic:

```go
func NewEnhancedTFLint() ship.Tool {
    baseTool := tools.NewTFLintTool()
    
    return ship.NewContainerTool("enhanced-tflint", ship.ContainerToolConfig{
        Description: "TFLint with custom enhancements",
        Parameters:  baseTool.Parameters(),
        Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
            result, err := baseTool.Execute(ctx, params, engine)
            if err != nil {
                return result, err
            }
            
            // Add custom logic here
            enhanced := addCustomLogic(result)
            return enhanced, nil
        },
    })
}
```

### How do I test my custom tools?

Create test environments with mock engines:

```go
func TestMyTool(t *testing.T) {
    tool := NewMyTool()
    
    // Use test engine or mock
    mockEngine := createTestEngine()
    
    result, err := tool.Execute(context.Background(), map[string]interface{}{
        "test_param": "test_value",
    }, mockEngine)
    
    assert.NoError(t, err)
    assert.Contains(t, result.Content, "expected output")
}
```

## Integration Questions

### How do I integrate with Claude Code?

Add your server to Claude Code's MCP configuration:

```json
{
  "mcpServers": {
    "my-ship-server": {
      "command": "/path/to/my-ship-server"
    }
  }
}
```

### How do I integrate with Cursor?

Configure in Cursor's MCP settings:

```json
{
  "mcpServers": {
    "my-ship-server": {
      "command": "/path/to/my-ship-server",
      "args": []
    }
  }
}
```

### Can I use this with other AI assistants?

Yes! The Ship SDK builds standard MCP servers that work with any MCP-compatible AI assistant.

### How do I handle different environments (dev/staging/prod)?

Create environment-specific servers:

```go
func buildServerForEnvironment(env string) *ship.Server {
    builder := ship.NewServer(fmt.Sprintf("%s-server", env), "1.0.0")
    
    switch env {
    case "dev":
        return builder.AddTool(mockDeployTool).Build()
    case "prod":
        return all.AddAllTools(builder.AddTool(prodDeployTool)).Build()
    default:
        return builder.Build()
    }
}
```

## Performance Questions

### How do I optimize container performance?

1. **Use small base images**: Alpine, distroless, or scratch
2. **Leverage caching**: Use Dagger's caching mechanisms
3. **Reuse containers**: Don't recreate containers unnecessarily
4. **Multi-stage builds**: Optimize your container images

### How do I handle large files?

Use mounted directories efficiently:

```go
result, err := engine.Container().
    From("my-tool:latest").
    WithMountedDirectory("/workspace", "./large-directory").
    WithWorkdir("/workspace").
    WithExec([]string{"./process-files"}).
    CombinedOutput(ctx)
```

### Can I limit resource usage?

Dagger inherits Docker's resource limits. Set limits via Docker daemon configuration or container orchestration platforms.

## Troubleshooting

### My tool execution is slow

Common causes:
- Large container images (use smaller bases)
- Cold container starts (leverage caching)
- Inefficient file operations (optimize mounted directories)
- Network latency (use local registries)

### I'm getting permission errors

Check:
- Container user permissions
- File system permissions on mounted directories
- Docker daemon permissions
- SELinux/AppArmor policies

### Tools are not finding files

Ensure:
- Correct directory mounting: `WithMountedDirectory("/workspace", ".")`
- Proper working directory: `WithWorkdir("/workspace")`
- Relative paths are correct within container

### How do I get help?

1. Check this FAQ and other documentation
2. Review examples in the `examples/` directory
3. Open an issue on GitHub
4. Join the CloudShip community discussions

## Roadmap Questions

### What new tools are planned?

See the [Ship Tools Reference](ship-tools-reference.md) for tools marked as "Coming Soon". Priorities include:
- Checkov security scanning
- Terraform documentation generation
- Infrastructure cost analysis
- Vulnerability scanning with Trivy

### Can I request new features?

Yes! Open a GitHub issue with:
- Clear description of the need
- Expected usage patterns
- Sample use cases

### How do I stay updated?

- Watch the GitHub repository for releases
- Follow CloudShip AI for announcements
- Check the changelog for new features