# Contributing to Ship CLI - Creating Custom Dagger Modules

Welcome to the Ship CLI community! This guide will help you create custom Dagger modules to extend Ship CLI's functionality.

## 🚀 Quick Start - Creating Your First Module

### Step 1: Module Structure

Create a new directory for your module with this structure:

```
my-awesome-tool/
├── dagger.json         # Dagger module configuration
├── main.go            # Module implementation
├── module.yaml        # Ship CLI integration metadata
└── README.md          # Documentation
```

### Step 2: Write Your Dagger Module

**main.go** - Example security scanner module:
```go
package main

import (
    "context"
    "fmt"
    "dagger.io/dagger"
)

type MyAwesomeTool struct{}

// Scan performs a custom security scan
func (m *MyAwesomeTool) Scan(
    ctx context.Context,
    // Directory to scan
    source *dagger.Directory,
    // Scan sensitivity level
    // +default="medium"
    sensitivity string,
) (string, error) {
    return dag.Container().
        From("alpine:latest").
        WithExec([]string{"apk", "add", "--no-cache", "your-tool"}).
        WithMountedDirectory("/src", source).
        WithWorkdir("/src").
        WithExec([]string{"your-tool", "scan", "--sensitivity", sensitivity}).
        Stdout(ctx)
}

// GenerateReport creates a detailed report
func (m *MyAwesomeTool) GenerateReport(
    ctx context.Context,
    source *dagger.Directory,
    // Output format
    // +default="json"
    format string,
) (*dagger.File, error) {
    return dag.Container().
        From("alpine:latest").
        WithExec([]string{"apk", "add", "--no-cache", "your-tool"}).
        WithMountedDirectory("/src", source).
        WithWorkdir("/src").
        WithExec([]string{"your-tool", "report", "--format", format, "-o", "report." + format}).
        File("report." + format)
}
```

### Step 3: Create Module Metadata

**module.yaml** - Ship CLI integration config:
```yaml
apiVersion: ship.cloudship.ai/v1
kind: Module
metadata:
  name: my-awesome-tool
  version: "1.0.0"
  description: "Custom security scanner for Ship CLI"
  author: "Your Name <you@example.com>"
  repository: "https://github.com/yourusername/ship-awesome-tool"
  
spec:
  type: dagger
  category: security  # Options: security, cost, documentation, analysis, utility
  
  # Dagger module configuration
  dagger:
    module: "."
    functions:
      - name: "scan"
        description: "Run security scan"
        cli_command: "awesome-scan"
      - name: "generateReport"
        description: "Generate detailed report"
        cli_command: "awesome-report"
    
  # CLI integration
  commands:
    - name: "awesome-scan"
      description: "Run awesome security scan"
      flags:
        - name: "sensitivity"
          type: "string"
          default: "medium"
          description: "Scan sensitivity level"
          enum: ["low", "medium", "high", "critical"]
          
    - name: "awesome-report"
      description: "Generate security report"
      flags:
        - name: "format"
          type: "string"
          default: "json"
          description: "Output format"
          enum: ["json", "yaml", "html", "markdown"]
          
  # AI/LLM integration
  ai_tools:
    - name: "awesome_security_scan"
      description: "Perform custom security analysis"
      usage: |
        {"tool": "awesome-tool", "action": "scan", "params": {"path": ".", "sensitivity": "high"}}
    
  # Requirements
  requirements:
    - "docker"
    - "dagger >= 0.18.0"
```

### Step 4: Create dagger.json

**dagger.json**:
```json
{
  "name": "my-awesome-tool",
  "sdk": "go",
  "source": ".",
  "engineVersion": "v0.18.10"
}
```

### Step 5: Document Your Module

**README.md**:
```markdown
# My Awesome Tool for Ship CLI

Custom security scanner that integrates with Ship CLI's ecosystem.

## Installation

```bash
# Coming soon: Install via Ship CLI
ship modules install github.com/yourusername/ship-awesome-tool

# Current method: Clone to modules directory
git clone https://github.com/yourusername/ship-awesome-tool ~/.ship/modules/my-awesome-tool
```

## Usage

```bash
# Run security scan
ship awesome-scan --sensitivity high

# Generate report
ship awesome-report --format markdown
```

## AI Integration

This module is available to Ship CLI's AI agents:

```bash
ship ai-agent --task "Use awesome-tool to scan for security issues"
```
```

## 📦 Module Categories

### Security Modules
- Vulnerability scanners
- Compliance checkers
- Secret detection tools
- Policy validators

### Cost Analysis Modules
- Cloud cost calculators
- Resource optimization tools
- Budget tracking
- Usage analytics

### Documentation Modules
- Diagram generators
- API documentation tools
- Architecture visualizers
- README generators

### Infrastructure Modules
- Resource provisioning
- Configuration management
- Deployment tools
- Migration assistants

## 🔧 Advanced Module Features

### 1. Using External Tools

```go
func (m *MyModule) AnalyzeWithTool(
    ctx context.Context,
    source *dagger.Directory,
) (string, error) {
    return dag.Container().
        From("my-tool:latest").
        WithMountedDirectory("/workspace", source).
        WithWorkdir("/workspace").
        WithExec([]string{"my-tool", "analyze"}).
        Stdout(ctx)
}
```

### 2. Multi-Container Workflows

```go
func (m *MyModule) ComplexWorkflow(
    ctx context.Context,
    source *dagger.Directory,
) (string, error) {
    // First container: Analysis
    analyzed := dag.Container().
        From("analyzer:latest").
        WithMountedDirectory("/src", source).
        WithExec([]string{"analyze", "--output", "/tmp/results.json"}).
        File("/tmp/results.json")
    
    // Second container: Processing
    return dag.Container().
        From("processor:latest").
        WithFile("/input/results.json", analyzed).
        WithExec([]string{"process", "/input/results.json"}).
        Stdout(ctx)
}
```

### 3. AI Tool Integration

```go
// Add AI-friendly function with structured output
func (m *MyModule) AIAnalysis(
    ctx context.Context,
    source *dagger.Directory,
    prompt string,
) (string, error) {
    result := // ... your analysis logic
    
    // Return structured JSON for AI consumption
    return fmt.Sprintf(`{
        "status": "success",
        "findings": %s,
        "recommendations": %s,
        "severity": "medium"
    }`, findings, recommendations), nil
}
```

## 🧪 Testing Your Module

### Local Testing

```bash
# Test with Dagger CLI directly
cd my-awesome-tool
dagger call scan --source=. --sensitivity=high

# Test with Ship CLI (after installation)
ship awesome-scan --sensitivity high
```

### Integration Testing

```go
// test/integration_test.go
func TestModuleIntegration(t *testing.T) {
    ctx := context.Background()
    
    // Initialize module
    module := &MyAwesomeTool{}
    
    // Test scan function
    result, err := module.Scan(ctx, testDir, "medium")
    assert.NoError(t, err)
    assert.Contains(t, result, "expected output")
}
```

## 📤 Publishing Your Module

### 1. GitHub Repository

Create a public repository with your module code:

```bash
git init
git add .
git commit -m "Initial module implementation"
git remote add origin https://github.com/yourusername/ship-awesome-tool
git push -u origin main
```

### 2. Tag a Release

```bash
git tag v1.0.0
git push origin v1.0.0
```

### 3. Submit to Community Registry

1. Fork [cloudshipai/community-modules](https://github.com/cloudshipai/community-modules)
2. Add your module to `registry.yaml`:

```yaml
modules:
  - name: my-awesome-tool
    repository: https://github.com/yourusername/ship-awesome-tool
    version: v1.0.0
    description: "Custom security scanner"
    author: "Your Name"
    category: security
    verified: false  # Will be set to true after review
```

3. Submit a Pull Request

## 🤝 Module Guidelines

### Do's
- ✅ Use containerized tools for isolation
- ✅ Provide clear documentation
- ✅ Include examples in your README
- ✅ Handle errors gracefully
- ✅ Output structured data (JSON/YAML)
- ✅ Support multiple output formats
- ✅ Make your module AI-friendly

### Don'ts
- ❌ Don't require local tool installation
- ❌ Don't hardcode credentials
- ❌ Don't modify files without user consent
- ❌ Don't use privileged containers without clear documentation
- ❌ Don't include sensitive data in logs

## 🔮 Future: Dynamic Module Loading

We're working on automatic module discovery and installation:

```bash
# Install from GitHub (coming soon)
ship modules install github.com/user/ship-module

# Install from registry (coming soon)
ship modules install awesome-tool

# List available modules (coming soon)
ship modules search security

# Update modules (coming soon)
ship modules update
```

## 💡 Module Ideas

Looking for inspiration? Here are some module ideas the community would love:

- **Cloud Security Scanner**: Scan AWS/Azure/GCP for misconfigurations
- **Kubernetes Analyzer**: Analyze K8s manifests and running clusters
- **Database Schema Validator**: Validate and document database schemas
- **API Test Runner**: Run and validate API tests
- **Performance Profiler**: Profile infrastructure performance
- **Compliance Checker**: Check compliance with various standards
- **Cost Forecaster**: Predict future cloud costs
- **Dependency Analyzer**: Analyze project dependencies
- **Log Aggregator**: Aggregate and analyze logs
- **Metric Collector**: Collect and visualize metrics

## 🆘 Getting Help

- **Discord**: Join our community at [discord.gg/cloudship](https://discord.gg/cloudship)
- **GitHub Issues**: [github.com/cloudshipai/ship/issues](https://github.com/cloudshipai/ship/issues)
- **Documentation**: [docs.cloudship.ai](https://docs.cloudship.ai)

## 📄 License

Modules should be licensed under MIT, Apache 2.0, or a compatible open-source license.

---

Ready to build something awesome? We can't wait to see what you create! 🚀