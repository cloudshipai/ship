# Ship Tools Reference

This document provides a comprehensive reference for all available Ship tools in the SDK.

## Overview

Ship tools are pre-built, containerized tools created by the CloudShip AI team. They provide secure, reliable infrastructure tooling that can be used standalone or integrated into custom MCP servers.

## Terraform Tools

### TFLint

**Name:** `tflint`  
**Description:** Run TFLint on Terraform code to check for syntax errors and best practices  
**Image:** `ghcr.io/terraform-linters/tflint:latest`

**Parameters:**
- `directory` (string, optional): Directory containing Terraform files (default: ".")
- `format` (string, optional): Output format: cli, json, junit, sarif (default: "cli")
- `init` (boolean, optional): Initialize TFLint before running (default: true)
- `output` (string, optional): Output file to save scan results

**Usage:**
```go
import "github.com/cloudshipai/ship/pkg/tools"

server := ship.NewServer("terraform-server", "1.0.0").
    AddTool(tools.NewTFLintTool()).
    Build()
```

**Example Execution:**
```go
result, err := tflintTool.Execute(ctx, map[string]interface{}{
    "directory": "./terraform",
    "format": "json",
    "init": true,
}, engine)
```

### Checkov (Coming Soon)

**Name:** `checkov`  
**Description:** Run Checkov security scan on Terraform code for policy compliance  
**Status:** Planned for next release

### Terraform Docs (Coming Soon)

**Name:** `terraform-docs`  
**Description:** Generate documentation for Terraform modules  
**Status:** Planned for next release

### Cost Analysis (Coming Soon)

**Name:** `cost-analysis`  
**Description:** Analyze infrastructure costs using OpenInfraQuote  
**Status:** Planned for next release

## Security Tools (Coming Soon)

### Trivy

**Name:** `trivy`  
**Description:** Security scanner for vulnerabilities and misconfigurations  
**Status:** Planned for next release

### InfraScan

**Name:** `infrascan`  
**Description:** Infrastructure security scanner  
**Status:** Planned for next release

## Documentation Tools (Coming Soon)

### Diagram Generator

**Name:** `generate-diagram`  
**Description:** Generate infrastructure diagrams from Terraform state  
**Status:** Planned for next release

## Using Tool Collections

### Terraform Registry

Get all Terraform tools:

```go
import "github.com/cloudshipai/ship/pkg/tools/all"

registry := all.TerraformRegistry()
server := ship.NewServer("terraform-server", "1.0.0").
    ImportRegistry(registry).
    Build()
```

### Convenience Functions

Add tool collections to server builders:

```go
// Add Terraform tools
server := all.AddTerraformTools(
    ship.NewServer("terraform-server", "1.0.0"),
).Build()

// Add all available tools
server := all.AddAllTools(
    ship.NewServer("full-server", "1.0.0"),
).Build()
```

## Tool Configuration

### Common Patterns

**Directory-based tools:**
Most Ship tools work with directories containing relevant files:

```go
params := map[string]interface{}{
    "directory": "./infrastructure",  // Path to files
}
```

**Output formatting:**
Many tools support multiple output formats:

```go
params := map[string]interface{}{
    "format": "json",  // json, cli, junit, sarif
    "output": "/tmp/results.json",  // Save to file
}
```

**Initialization:**
Some tools need initialization:

```go
params := map[string]interface{}{
    "init": true,  // Run initialization first
}
```

### Error Handling

Ship tools follow consistent error handling:

```go
result, err := tool.Execute(ctx, params, engine)
if err != nil {
    // Container execution failed
    log.Printf("Tool execution failed: %v", err)
    return err
}

if result.Error != nil {
    // Tool reported an error (e.g., lint violations)
    log.Printf("Tool reported issues: %v", result.Error)
}

// Check metadata for additional information
if violations, ok := result.Metadata["violations"].(int); ok {
    log.Printf("Found %d violations", violations)
}
```

## Tool Metadata

Ship tools provide rich metadata in results:

### TFLint Metadata

```go
{
    "tool": "tflint",
    "directory": "./terraform",
    "format": "json", 
    "violations": 5,
    "files_scanned": 12,
    "duration_ms": 1234
}
```

### Common Metadata Fields

- `tool`: Tool name
- `directory`: Scanned directory  
- `format`: Output format used
- `violations`: Number of issues found
- `files_scanned`: Number of files processed
- `duration_ms`: Execution time in milliseconds
- `exit_code`: Container exit code

## Custom Tool Integration

### Extending Ship Tools

Wrap Ship tools with custom logic:

```go
func NewEnhancedTFLintTool() ship.Tool {
    baseTool := tools.NewTFLintTool()
    
    return ship.NewContainerTool("enhanced-tflint", ship.ContainerToolConfig{
        Description: "TFLint with custom post-processing",
        Image:       "alpine:latest",  // Wrapper container
        Parameters:  baseTool.Parameters(),
        Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
            // Run base TFLint
            result, err := baseTool.Execute(ctx, params, engine)
            if err != nil {
                return result, err
            }
            
            // Custom post-processing
            enhanced := enhanceResults(result.Content)
            
            return &ship.ToolResult{
                Content: enhanced,
                Metadata: result.Metadata,
            }, nil
        },
    })
}
```

### Tool Composition

Combine multiple Ship tools:

```go
func NewComplianceChecker() ship.Tool {
    return ship.NewContainerTool("compliance-check", ship.ContainerToolConfig{
        Description: "Run comprehensive compliance checks",
        Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
            var results []string
            
            // Run TFLint
            tflintResult, err := tools.NewTFLintTool().Execute(ctx, params, engine)
            if err == nil {
                results = append(results, "✓ TFLint passed")
            } else {
                results = append(results, "✗ TFLint failed")
            }
            
            // Run Checkov (when available)
            // checkovResult, err := tools.NewCheckovTool().Execute(ctx, params, engine)
            
            // Combine results
            return &ship.ToolResult{
                Content: strings.Join(results, "\n"),
                Metadata: map[string]interface{}{
                    "checks_run": len(results),
                    "compliance_score": calculateScore(results),
                },
            }, nil
        },
    })
}
```

## Tool Development

### Contributing New Tools

To contribute a new Ship tool:

1. **Create Tool Implementation:**
```go
// pkg/tools/mytool.go
func NewMyTool() ship.Tool {
    return ship.NewContainerTool("mytool", ship.ContainerToolConfig{
        Description: "My awesome tool",
        Image:       "custom/mytool:latest",
        Parameters: []ship.Parameter{
            {Name: "input", Type: "string", Required: true},
        },
        Execute: executeMyTool,
    })
}
```

2. **Add Tests:**
```go
// pkg/tools/mytool_test.go
func TestMyTool(t *testing.T) {
    tool := NewMyTool()
    // Test implementation
}
```

3. **Update Collections:**
```go
// pkg/tools/all/terraform.go (or appropriate collection)
func TerraformRegistry() ship.Registry {
    registry := ship.NewRegistry()
    registry.RegisterTool(NewTFLintTool())
    registry.RegisterTool(NewMyTool())  // Add your tool
    return registry
}
```

### Tool Guidelines

**Container Images:**
- Use official, minimal images when possible
- Pin to specific versions for reproducibility
- Document image requirements

**Parameters:**
- Follow consistent naming conventions
- Provide good descriptions and examples
- Use enums for limited value sets
- Set sensible defaults

**Error Handling:**
- Distinguish between execution errors and tool findings
- Provide actionable error messages
- Include relevant context in metadata

**Security:**
- Run containers with minimal privileges
- Validate and sanitize inputs
- Never log sensitive information
- Use Dagger secrets for credentials

This reference will be updated as new Ship tools are added to the collection.