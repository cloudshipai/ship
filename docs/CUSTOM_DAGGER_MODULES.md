# Creating Custom Dagger Modules for Ship CLI 🚢

This guide shows you how to create your own Dagger modules and integrate them into Ship CLI.

## Table of Contents
- [Overview](#overview)
- [Module Structure](#module-structure)
- [Step-by-Step Guide](#step-by-step-guide)
- [Example: Creating a Python Linter Module](#example-creating-a-python-linter-module)
- [Integration with Ship CLI](#integration-with-ship-cli)
- [Best Practices](#best-practices)

## Overview

Ship CLI uses Dagger modules to run containerized tools without requiring local installation. You can extend Ship CLI by creating custom Dagger modules for your own tools and workflows.

## Module Structure

A Dagger module in Ship CLI follows this structure:

```
internal/dagger/modules/
├── your_module.go         # Your module implementation
├── existing_modules.go    # Existing modules for reference
└── ...
```

## Step-by-Step Guide

### 1. Create Your Module File

Create a new file in `internal/dagger/modules/` directory:

```go
// internal/dagger/modules/python_tools.go
package modules

import (
    "context"
    "fmt"
    "dagger.io/dagger"
)

// PythonTools provides Python development tools via Dagger
type PythonTools struct {
    client *dagger.Client
}

// NewPythonTools creates a new PythonTools instance
func NewPythonTools(client *dagger.Client) *PythonTools {
    return &PythonTools{client: client}
}
```

### 2. Define Module Methods

Add methods for each tool you want to expose:

```go
// RunPylint runs pylint on Python code
func (pt *PythonTools) RunPylint(ctx context.Context, dir string) (string, error) {
    // Create container with Python and pylint
    container := pt.client.Container().
        From("python:3.11-slim").
        WithExec([]string{"pip", "install", "pylint"}).
        WithMountedDirectory("/workspace", pt.client.Host().Directory(dir)).
        WithWorkdir("/workspace")

    // Run pylint
    output, err := container.
        WithExec([]string{"pylint", ".", "--output-format=json"}).
        Stdout(ctx)

    if err != nil {
        return "", fmt.Errorf("pylint failed: %w", err)
    }

    return output, nil
}

// RunBlack formats Python code
func (pt *PythonTools) RunBlack(ctx context.Context, dir string, fix bool) (string, error) {
    args := []string{"black", "."}
    if !fix {
        args = append(args, "--check", "--diff")
    }

    container := pt.client.Container().
        From("python:3.11-slim").
        WithExec([]string{"pip", "install", "black"}).
        WithMountedDirectory("/workspace", pt.client.Host().Directory(dir)).
        WithWorkdir("/workspace")

    output, err := container.
        WithExec(args).
        Stdout(ctx)

    if err != nil {
        return "", fmt.Errorf("black failed: %w", err)
    }

    return output, nil
}
```

### 3. Create CLI Command

Create a new command file in `internal/cli/`:

```go
// internal/cli/python_tools_cmd.go
package cli

import (
    "context"
    "fmt"
    "log/slog"
    
    "github.com/spf13/cobra"
    "dagger.io/dagger"
    "github.com/yourusername/ship/internal/dagger/modules"
)

func pythonToolsCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "python-tools",
        Short: "Run Python development tools",
    }

    cmd.AddCommand(pythonLintCmd())
    cmd.AddCommand(pythonFormatCmd())

    return cmd
}

func pythonLintCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "lint [directory]",
        Short: "Lint Python code with pylint",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            dir := args[0]

            // Initialize Dagger client
            client, err := dagger.Connect(ctx, dagger.WithLogOutput(slog.Default().Handler()))
            if err != nil {
                return fmt.Errorf("failed to connect to dagger: %w", err)
            }
            defer client.Close()

            // Run pylint
            pt := modules.NewPythonTools(client)
            output, err := pt.RunPylint(ctx, dir)
            if err != nil {
                return err
            }

            fmt.Println(output)
            return nil
        },
    }
}

func pythonFormatCmd() *cobra.Command {
    var fix bool
    
    cmd := &cobra.Command{
        Use:   "format [directory]",
        Short: "Format Python code with Black",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            dir := args[0]

            client, err := dagger.Connect(ctx, dagger.WithLogOutput(slog.Default().Handler()))
            if err != nil {
                return fmt.Errorf("failed to connect to dagger: %w", err)
            }
            defer client.Close()

            pt := modules.NewPythonTools(client)
            output, err := pt.RunBlack(ctx, dir, fix)
            if err != nil {
                return err
            }

            fmt.Println(output)
            return nil
        },
    }

    cmd.Flags().BoolVar(&fix, "fix", false, "Apply formatting fixes")
    return cmd
}
```

### 4. Register Command

Add your command to the root command in `internal/cli/root.go`:

```go
func newRootCmd() *cobra.Command {
    // ... existing code ...

    rootCmd.AddCommand(
        authCmd(),
        pushCmd(),
        queryCmd(),
        investigateCmd(),
        mcpCmd(),
        terraformToolsCmd(),
        pythonToolsCmd(), // Add your new command here
    )

    return rootCmd
}
```

## Example: Creating a Python Linter Module

Here's a complete example of a Python linter module that integrates with Ship CLI:

### Module Implementation

```go
// internal/dagger/modules/python_linter.go
package modules

import (
    "context"
    "encoding/json"
    "fmt"
    "path/filepath"
    
    "dagger.io/dagger"
)

type PythonLinter struct {
    client *dagger.Client
}

func NewPythonLinter(client *dagger.Client) *PythonLinter {
    return &PythonLinter{client: client}
}

type LintResult struct {
    File    string `json:"file"`
    Line    int    `json:"line"`
    Column  int    `json:"column"`
    Message string `json:"message"`
    Type    string `json:"type"`
}

func (pl *PythonLinter) LintWithMultipleTools(ctx context.Context, dir string) (map[string][]LintResult, error) {
    results := make(map[string][]LintResult)

    // Base container with multiple Python tools
    baseContainer := pl.client.Container().
        From("python:3.11-slim").
        WithExec([]string{"pip", "install", "pylint", "flake8", "mypy", "bandit"}).
        WithMountedDirectory("/workspace", pl.client.Host().Directory(dir)).
        WithWorkdir("/workspace")

    // Run pylint
    pylintOutput, _ := baseContainer.
        WithExec([]string{"pylint", ".", "--output-format=json", "--exit-zero"}).
        Stdout(ctx)
    
    var pylintResults []LintResult
    json.Unmarshal([]byte(pylintOutput), &pylintResults)
    results["pylint"] = pylintResults

    // Run flake8
    flake8Output, _ := baseContainer.
        WithExec([]string{"flake8", ".", "--format=json", "--exit-zero"}).
        Stdout(ctx)
    
    var flake8Results []LintResult
    json.Unmarshal([]byte(flake8Output), &flake8Results)
    results["flake8"] = flake8Results

    // Run security scan with bandit
    banditOutput, _ := baseContainer.
        WithExec([]string{"bandit", "-r", ".", "-f", "json", "--exit-zero"}).
        Stdout(ctx)
    
    var banditResults []LintResult
    json.Unmarshal([]byte(banditOutput), &banditResults)
    results["bandit"] = banditResults

    return results, nil
}

// Add CloudShip upload capability
func (pl *PythonLinter) LintAndUpload(ctx context.Context, dir string, apiKey string) error {
    results, err := pl.LintWithMultipleTools(ctx, dir)
    if err != nil {
        return err
    }

    // Convert results to CloudShip format
    artifact := map[string]interface{}{
        "type":    "python-lint-results",
        "results": results,
        "path":    dir,
    }

    // Use existing CloudShip upload logic
    return uploadToCloudShip(artifact, apiKey)
}
```

### CLI Command Implementation

```go
// internal/cli/python_lint_cmd.go
package cli

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "dagger.io/dagger"
    "github.com/yourusername/ship/internal/dagger/modules"
)

func pythonLintCmd() *cobra.Command {
    var push bool
    var outputFormat string

    cmd := &cobra.Command{
        Use:   "python-lint [directory]",
        Short: "Comprehensive Python linting with multiple tools",
        Long: `Run multiple Python linting tools including:
- pylint: Code quality and style checker
- flake8: Style guide enforcement
- mypy: Static type checker
- bandit: Security issue scanner`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            dir := args[0]

            slog.Info("Initializing Dagger engine...")
            client, err := dagger.Connect(ctx, dagger.WithLogOutput(slog.Default().Handler()))
            if err != nil {
                return fmt.Errorf("failed to connect to dagger: %w", err)
            }
            defer client.Close()

            linter := modules.NewPythonLinter(client)

            if push {
                apiKey := viper.GetString("cloudship.api_key")
                if apiKey == "" {
                    return fmt.Errorf("CloudShip API key not configured. Run 'ship auth' first")
                }
                
                slog.Info("Running linters and uploading to CloudShip...")
                return linter.LintAndUpload(ctx, dir, apiKey)
            }

            slog.Info("Running Python linters...")
            results, err := linter.LintWithMultipleTools(ctx, dir)
            if err != nil {
                return fmt.Errorf("linting failed: %w", err)
            }

            // Output results based on format
            switch outputFormat {
            case "json":
                output, _ := json.MarshalIndent(results, "", "  ")
                fmt.Println(string(output))
            default:
                // Pretty print results
                for tool, issues := range results {
                    fmt.Printf("\n=== %s Results ===\n", tool)
                    if len(issues) == 0 {
                        fmt.Printf("✅ No issues found!\n")
                    } else {
                        for _, issue := range issues {
                            fmt.Printf("⚠️  %s:%d:%d - %s (%s)\n", 
                                issue.File, issue.Line, issue.Column, 
                                issue.Message, issue.Type)
                        }
                    }
                }
            }

            return nil
        },
    }

    cmd.Flags().BoolVar(&push, "push", false, "Push results to CloudShip")
    cmd.Flags().StringVar(&outputFormat, "output", "text", "Output format (text|json)")

    return cmd
}
```

## Integration with Ship CLI

### 1. Add to Root Command

```go
// internal/cli/root.go
func newRootCmd() *cobra.Command {
    rootCmd := &cobra.Command{
        Use:   "ship",
        Short: "Ship CLI - Cloud infrastructure analysis and automation",
    }

    rootCmd.AddCommand(
        // ... existing commands ...
        pythonLintCmd(), // Add your custom command
    )

    return rootCmd
}
```

### 2. Test Your Module

```bash
# Build Ship CLI
go build -o ship ./cmd/ship

# Test your new command
./ship python-lint ./my-python-project

# Test with CloudShip upload
./ship python-lint ./my-python-project --push

# Test with JSON output
./ship python-lint ./my-python-project --output json
```

## Best Practices

### 1. Container Optimization

```go
// Cache dependencies in a separate layer
func (m *MyModule) BuildContainer(ctx context.Context) *dagger.Container {
    return m.client.Container().
        From("base-image:latest").
        // Copy dependency files first
        WithFile("/tmp/requirements.txt", m.client.Host().File("requirements.txt")).
        // Install dependencies (this layer gets cached)
        WithExec([]string{"pip", "install", "-r", "/tmp/requirements.txt"}).
        // Then mount the actual code
        WithMountedDirectory("/workspace", m.client.Host().Directory("."))
}
```

### 2. Error Handling

```go
func (m *MyModule) RunTool(ctx context.Context, dir string) (string, error) {
    container := m.BuildContainer(ctx)
    
    // Capture both stdout and stderr
    stdout, err := container.Stdout(ctx)
    if err != nil {
        stderr, _ := container.Stderr(ctx)
        return "", fmt.Errorf("tool failed: %w\nstderr: %s", err, stderr)
    }
    
    return stdout, nil
}
```

### 3. Progress Reporting

```go
func (m *MyModule) RunWithProgress(ctx context.Context, dir string) error {
    slog.Info("Starting analysis...", "directory", dir)
    
    // Step 1
    slog.Info("Installing dependencies...")
    // ... code ...
    
    // Step 2
    slog.Info("Running analysis...")
    // ... code ...
    
    slog.Info("✓ Analysis complete!")
    return nil
}
```

### 4. Configuration Support

```go
type ModuleConfig struct {
    ToolVersion string
    ConfigFile  string
    Parallel    bool
}

func (m *MyModule) RunWithConfig(ctx context.Context, dir string, config ModuleConfig) error {
    container := m.client.Container().From(fmt.Sprintf("tool:%s", config.ToolVersion))
    
    if config.ConfigFile != "" {
        container = container.WithFile("/config", m.client.Host().File(config.ConfigFile))
    }
    
    // ... rest of implementation
}
```

### 5. Multi-Stage Operations

```go
func (m *MyModule) ComplexWorkflow(ctx context.Context, dir string) error {
    // Stage 1: Build
    buildContainer := m.client.Container().
        From("golang:1.21").
        WithMountedDirectory("/src", m.client.Host().Directory(dir)).
        WithWorkdir("/src").
        WithExec([]string{"go", "build", "-o", "app"})

    // Stage 2: Test
    testContainer := buildContainer.
        WithExec([]string{"go", "test", "./..."})

    // Stage 3: Security Scan
    scanContainer := m.client.Container().
        From("aquasec/trivy:latest").
        WithMountedDirectory("/workspace", m.client.Host().Directory(dir)).
        WithExec([]string{"trivy", "fs", "/workspace"})

    // Execute all stages
    _, err := scanContainer.Sync(ctx)
    return err
}
```

## Advanced Features

### 1. Cross-Module Integration

```go
// Use outputs from one module as inputs to another
func (m *MyModule) IntegrateWithSteampipe(ctx context.Context) error {
    // Get Steampipe results
    steampipe := modules.NewSteampipe(m.client)
    queryResults, err := steampipe.Query(ctx, "SELECT * FROM aws_s3_bucket")
    
    // Use results in your module
    return m.ProcessResults(ctx, queryResults)
}
```

### 2. Caching and Performance

```go
func (m *MyModule) WithCache(ctx context.Context, dir string) *dagger.Container {
    cacheVolume := m.client.CacheVolume("my-module-cache")
    
    return m.client.Container().
        From("base-image").
        WithMountedCache("/cache", cacheVolume).
        WithEnvVariable("CACHE_DIR", "/cache")
}
```

### 3. Secret Management

```go
func (m *MyModule) WithSecrets(ctx context.Context, apiKey string) *dagger.Container {
    secret := m.client.SetSecret("api-key", apiKey)
    
    return m.client.Container().
        From("base-image").
        WithSecretVariable("API_KEY", secret)
}
```

## Troubleshooting

### Common Issues

1. **Container build failures**: Check Docker daemon is running
2. **Permission errors**: Ensure mounted directories are readable
3. **Network issues**: Container may need internet access for package installation
4. **Output parsing**: Always handle both stdout and stderr

### Debugging Tips

```go
// Enable verbose logging
func (m *MyModule) Debug(ctx context.Context) error {
    container := m.client.Container().
        From("base-image").
        WithEnvVariable("DEBUG", "true").
        WithExec([]string{"sh", "-c", "echo 'Debug info:' && env"})
    
    output, _ := container.Stdout(ctx)
    slog.Debug("Container environment", "output", output)
    return nil
}
```

## Contributing

When contributing custom modules:

1. Follow Go coding standards
2. Add comprehensive tests
3. Document all public methods
4. Include usage examples
5. Update the main README

## Resources

- [Dagger Documentation](https://docs.dagger.io/)
- [Ship CLI Architecture](./technical-spec.md)
- [Contributing Guidelines](../CONTRIBUTING.md)
- [Example Modules](../internal/dagger/modules/)

---

Happy module building! 🚀