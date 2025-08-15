# Ship SDK Usage Patterns

This guide covers advanced usage patterns and best practices for the Ship SDK.

## The Three Core Patterns

### Pattern 1: Pure Ship SDK (Custom Tools Only)

**When to use:** Building specialized MCP servers with completely custom functionality.

**Benefits:**
- Complete control over tool behavior
- Minimal dependencies
- Custom business logic integration

**Example:**
```go
package main

import (
    "context"
    "fmt"
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/pkg/dagger"
)

func main() {
    // Custom database migration tool
    migrationTool := ship.NewContainerTool("db-migrate", ship.ContainerToolConfig{
        Description: "Run database migrations",
        Image:       "migrate/migrate:latest",
        Parameters: []ship.Parameter{
            {Name: "direction", Type: "string", Required: true, Enum: []string{"up", "down"}},
            {Name: "steps", Type: "number", Required: false, Default: 1},
        },
        Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
            direction := params["direction"].(string)
            steps := 1
            if s, ok := params["steps"].(float64); ok {
                steps = int(s)
            }
            
            args := []string{"migrate", "-path", "/migrations", "-database", "postgres://...", direction}
            if direction != "up" {
                args = append(args, fmt.Sprintf("%d", steps))
            }
            
            output, err := engine.Container().
                From("migrate/migrate:latest").
                WithMountedDirectory("/migrations", "./migrations").
                WithExec(args).
                CombinedOutput(ctx)
                
            return &ship.ToolResult{
                Content: output,
                Metadata: map[string]interface{}{
                    "direction": direction,
                    "steps": steps,
                },
            }, err
        },
    })

    server := ship.NewServer("migration-server", "1.0.0").
        AddTool(migrationTool).
        Build()
        
    server.ServeStdio()
}
```

### Pattern 2: Cherry-Pick Ship Tools

**When to use:** Need some Ship tools but want to add custom extensions.

**Benefits:**
- Leverage proven Ship tools
- Add custom business logic
- Selective tool inclusion

**Example:**
```go
package main

import (
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/pkg/tools"
)

func main() {
    // Custom compliance checker
    complianceChecker := ship.NewContainerTool("compliance-check", ship.ContainerToolConfig{
        Description: "Check infrastructure compliance",
        Image:       "custom/compliance-scanner:latest",
        Parameters: []ship.Parameter{
            {Name: "framework", Type: "string", Required: true, Enum: []string{"SOC2", "HIPAA", "PCI"}},
            {Name: "strict", Type: "boolean", Required: false, Default: false},
        },
        Execute: complianceExecutor,
    })

    server := ship.NewServer("compliance-server", "1.0.0").
        AddTool(tools.NewTFLintTool()).           // Ship tool
        AddTool(tools.NewCheckovTool()).          // Ship tool  
        AddTool(complianceChecker).               // Custom tool
        Build()
        
    server.ServeStdio()
}
```

### Pattern 3: Everything Plus

**When to use:** Want all Ship tools plus custom extensions for comprehensive workflows.

**Benefits:**
- Complete Ship toolchain
- Custom extensions
- One-stop solution

**Example:**
```go
package main

import (
    "github.com/cloudshipai/ship/pkg/ship"
    "github.com/cloudshipai/ship/pkg/tools/all"
)

func main() {
    // Custom deployment orchestrator
    deploymentTool := ship.NewContainerTool("deploy", ship.ContainerToolConfig{
        Description: "Deploy infrastructure with rollback capability",
        Image:       "custom/deployer:latest",
        Parameters: []ship.Parameter{
            {Name: "environment", Type: "string", Required: true},
            {Name: "rollback", Type: "boolean", Required: false},
        },
        Execute: deploymentExecutor,
    })

    server := all.AddAllTools(
        ship.NewServer("comprehensive-server", "1.0.0").
        AddTool(deploymentTool),
    ).Build()
    
    server.ServeStdio()
}
```

## Advanced Patterns

### Multi-Registry Architecture

Organize tools into logical registries:

```go
// Infrastructure registry
infraRegistry := ship.NewRegistry()
infraRegistry.RegisterTool(tools.NewTFLintTool())
infraRegistry.RegisterTool(tools.NewCheckovTool())

// Security registry  
securityRegistry := ship.NewRegistry()
securityRegistry.RegisterTool(tools.NewTrivyTool())
securityRegistry.RegisterTool(customVulnScanner)

// Combined server
server := ship.NewServer("multi-registry-server", "1.0.0").
    ImportRegistry(infraRegistry).
    ImportRegistry(securityRegistry).
    Build()
```

### Environment-Specific Tools

Build different servers for different environments:

```go
func buildDevServer() *ship.Server {
    return ship.NewServer("dev-server", "1.0.0").
        AddTool(tools.NewTFLintTool()).
        AddTool(mockDeploymentTool).
        Build()
}

func buildProdServer() *ship.Server {
    return all.AddAllTools(
        ship.NewServer("prod-server", "1.0.0").
        AddTool(prodDeploymentTool).
        AddTool(alertingTool),
    ).Build()
}
```

### Workflow Orchestration

Chain tools for complex workflows:

```go
workflowTool := ship.NewContainerTool("tf-workflow", ship.ContainerToolConfig{
    Description: "Complete Terraform workflow: lint -> security -> plan -> deploy",
    Image:       "alpine:latest",
    Parameters: []ship.Parameter{
        {Name: "directory", Type: "string", Required: true},
        {Name: "auto_deploy", Type: "boolean", Required: false},
    },
    Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
        dir := params["directory"].(string)
        autoDeploy := params["auto_deploy"].(bool)
        
        var results []string
        
        // 1. Lint
        lintResult, err := tflintTool.Execute(ctx, map[string]interface{}{"directory": dir}, engine)
        if err != nil {
            return nil, fmt.Errorf("lint failed: %w", err)
        }
        results = append(results, "✓ Lint passed")
        
        // 2. Security scan
        secResult, err := checkovTool.Execute(ctx, map[string]interface{}{"directory": dir}, engine)
        if err != nil {
            return nil, fmt.Errorf("security scan failed: %w", err)
        }
        results = append(results, "✓ Security scan passed")
        
        // 3. Plan
        planResult, err := terraformPlanTool.Execute(ctx, map[string]interface{}{"directory": dir}, engine)
        if err != nil {
            return nil, fmt.Errorf("plan failed: %w", err)
        }
        results = append(results, "✓ Plan generated")
        
        // 4. Deploy (if auto_deploy)
        if autoDeploy {
            deployResult, err := terraformApplyTool.Execute(ctx, map[string]interface{}{"directory": dir}, engine)
            if err != nil {
                return nil, fmt.Errorf("deploy failed: %w", err)
            }
            results = append(results, "✓ Deployed successfully")
        }
        
        return &ship.ToolResult{
            Content: strings.Join(results, "\n"),
            Metadata: map[string]interface{}{
                "steps_completed": len(results),
                "auto_deployed": autoDeploy,
            },
        }, nil
    },
})
```

## Integration Patterns

### With Existing Infrastructure

```go
// Tool that integrates with existing monitoring
monitoringTool := ship.NewContainerTool("monitor", ship.ContainerToolConfig{
    Description: "Send metrics to existing monitoring system",
    Image:       "custom/monitoring-client:latest",
    Parameters: []ship.Parameter{
        {Name: "metric_name", Type: "string", Required: true},
        {Name: "value", Type: "number", Required: true},
        {Name: "tags", Type: "object", Required: false},
    },
    Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
        // Send metrics to Datadog, Prometheus, etc.
        return sendMetrics(ctx, params, engine)
    },
})
```

### With CI/CD Systems

```go
// Tool that integrates with CI/CD
cicdTool := ship.NewContainerTool("trigger-pipeline", ship.ContainerToolConfig{
    Description: "Trigger CI/CD pipeline",
    Image:       "alpine/curl:latest",
    Parameters: []ship.Parameter{
        {Name: "pipeline_id", Type: "string", Required: true},
        {Name: "branch", Type: "string", Required: false, Default: "main"},
    },
    Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
        // Trigger GitHub Actions, GitLab CI, Jenkins, etc.
        return triggerPipeline(ctx, params, engine)
    },
})
```

## Performance Patterns

### Container Optimization

```go
// Use multi-stage builds for smaller containers
optimizedTool := ship.NewContainerTool("optimized-tool", ship.ContainerToolConfig{
    Description: "Performance optimized tool",
    Image:       "custom/optimized:latest",  // Multi-stage build
    Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
        // Use cached base container
        container := engine.Container().
            From("custom/optimized:latest").
            WithMountedCache("/tmp/cache", engine.CacheVolume("tool-cache"))
            
        // Tool logic here
        return executeOptimized(ctx, container, params)
    },
})
```

### Parallel Execution

```go
// Tool that runs multiple checks in parallel
parallelTool := ship.NewContainerTool("parallel-checks", ship.ContainerToolConfig{
    Description: "Run multiple checks in parallel",
    Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
        // Run checks in parallel using goroutines
        var wg sync.WaitGroup
        results := make(chan string, 3)
        
        checks := []func(){
            func() { runLintCheck(ctx, engine, results) },
            func() { runSecurityCheck(ctx, engine, results) },
            func() { runTestCheck(ctx, engine, results) },
        }
        
        for _, check := range checks {
            wg.Add(1)
            go func(c func()) {
                defer wg.Done()
                c()
            }(check)
        }
        
        wg.Wait()
        close(results)
        
        var allResults []string
        for result := range results {
            allResults = append(allResults, result)
        }
        
        return &ship.ToolResult{
            Content: strings.Join(allResults, "\n"),
        }, nil
    },
})
```

## Testing Patterns

### Mock Tools for Testing

```go
func createMockTool(name string, response string) ship.Tool {
    return ship.NewContainerTool(name, ship.ContainerToolConfig{
        Description: fmt.Sprintf("Mock %s tool", name),
        Image:       "alpine:latest",
        Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
            return &ship.ToolResult{
                Content: response,
                Metadata: map[string]interface{}{"mock": true},
            }, nil
        },
    })
}

func TestWorkflow(t *testing.T) {
    server := ship.NewServer("test-server", "1.0.0").
        AddTool(createMockTool("lint", "✓ Lint passed")).
        AddTool(createMockTool("security", "✓ Security passed")).
        Build()
        
    // Test workflow
    // ...
}
```

## Security Patterns

### Input Validation

```go
secureToolExecutor := func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
    // Validate inputs
    if err := validateInputs(params); err != nil {
        return &ship.ToolResult{Error: err}, err
    }
    
    // Sanitize file paths
    directory := sanitizePath(params["directory"].(string))
    
    // Run in restricted container
    result, err := engine.Container().
        From("alpine:latest").
        WithUser("nobody").  // Non-root user
        WithMountedDirectory("/workspace", directory).
        WithWorkdir("/workspace").
        WithExec([]string{"./safe-command"}).
        CombinedOutput(ctx)
        
    return &ship.ToolResult{Content: result}, err
}
```

### Secret Management

```go
secretAwareTool := ship.NewContainerTool("deploy", ship.ContainerToolConfig{
    Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
        // Never log secrets
        apiKey := params["api_key"].(string)
        
        result, err := engine.Container().
            From("deployment:latest").
            WithSecretVariable("API_KEY", engine.Secret(apiKey)).  // Use Dagger secrets
            WithExec([]string{"deploy"}).
            CombinedOutput(ctx)
            
        // Redact secrets from output
        cleanResult := redactSecrets(result)
        
        return &ship.ToolResult{Content: cleanResult}, err
    },
})
```

These patterns provide a foundation for building robust, maintainable MCP servers with the Ship SDK.