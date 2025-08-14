# Ship CLI

A powerful command-line tool that brings enterprise-grade infrastructure analysis tools to your fingertips, all running in containers without local installations.

> **🤖 For LLMs and AI Assistants**: Complete installation and usage instructions specifically designed for AI consumption are available in [llms.txt](./llms.txt). This includes MCP server setup, integration examples, and best practices for AI-driven infrastructure analysis with all 7 MCP tools.

## 🚀 Features

- **🔍 Terraform Linting**: Catch errors and enforce best practices with TFLint
- **🛡️ Security Scanning**: Multi-cloud security analysis with Checkov and Trivy
- **💰 Cost Estimation**: Estimate infrastructure costs with Infracost and OpenInfraQuote
- **📝 Documentation Generation**: Auto-generate beautiful Terraform module documentation
- **📊 Infrastructure Diagrams**: Visualize your infrastructure with InfraMap integration
- **🔍 Infrastructure Diagram Generation**: Create visual diagrams from Terraform configurations
- **🤖 AI Assistant Integration**: Built-in MCP server for Claude Desktop, Cursor, and other AI tools
- **🔌 Extensible Module System**: Add custom tools and Dagger functions without modifying core CLI
- **🐳 Containerized Tools**: All tools run in containers via Dagger - no local installations needed
- **☁️ Cloud Integration**: Seamlessly works with AWS, Azure, GCP, and other cloud providers
- **🔧 CI/CD Ready**: Perfect for integration into your existing pipelines

## 📊 Privacy & Telemetry

Ship CLI collects anonymous usage telemetry to help us improve the tool. This data is completely anonymous and helps us understand which features are most valuable.

**What we collect:**
- Command usage patterns (e.g., which MCP tools are used)
- Tool execution frequency
- Anonymous system identifiers (no personal information)

**What we DON'T collect:**
- File contents or code
- Personal information
- Project names or paths
- Error details or sensitive data

**Opt-out anytime:**
```bash
# Disable telemetry via environment variable
export SHIP_TELEMETRY=false

# Or disable in config
ship vars set telemetry.enabled false
```

The telemetry system is powered by [PostHog](https://posthog.com) and uses industry-standard privacy practices.

## 📚 Table of Contents

- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Demo](#-demo)
- [Available Tools](#️-available-tools)
- [Command Reference](#-command-reference)
- [Authentication](#-authentication)
- [Using External Dagger Modules](#-using-external-dagger-modules)
- [Contributing](#-contributing)
- [License](#-license)

## 🎬 Demo

### Terraform Tools in Action

![Ship CLI Terraform Tools Demo](./demos/terraform-tools-clean.gif)

> This demo shows Ship CLI running terraform-docs, tflint, and security scanning on a Terraform module - all without any local tool installations!

### OpenInfraQuote - Advanced Cost Analysis

![OpenInfraQuote Cost Analysis Demo](./demos/openinfraquote-real-demo.gif)

> OpenInfraQuote provides highly accurate AWS cost estimation by analyzing your Terraform plans against real AWS pricing data. It supports 100+ AWS resource types with region-specific pricing!

## 📦 Installation

### Quick Install with Go

```bash
# Install directly with Go
go install github.com/cloudshipai/ship/cmd/ship@latest

# Verify installation
ship version
```

### From Source

```bash
# Clone the repository
git clone https://github.com/cloudshipai/ship.git
cd ship

# Build and install
go build -o ship ./cmd/ship
sudo mv ship /usr/local/bin/

# Or just run directly
go run ./cmd/ship [command]
```

## 🏃 Quick Start

### 1. Basic Usage

```bash
# Navigate to your Terraform project
cd your-terraform-project

# Run a comprehensive analysis
ship tf lint                # Check for errors and best practices
ship tf checkov             # Security scanning
ship tf cost                # Estimate AWS/Azure/GCP costs
ship tf docs                # Generate documentation
```

### 2. Real-World Example

```bash
# Clone a sample Terraform project
git clone https://github.com/terraform-aws-modules/terraform-aws-vpc.git
cd terraform-aws-vpc/examples/simple

# Run all tools
ship tf lint
ship tf checkov
ship tf trivy
ship tf cost
ship tf docs > README.md
ship tf diagram . --hcl -o infrastructure.png
```

### 3. Generate Infrastructure Diagrams

Visualize your infrastructure with InfraMap integration:

```bash
# Generate diagram from Terraform files (no state file needed!)
ship tf diagram . --hcl --format png -o infrastructure.png

# Generate from existing state file
ship tf diagram terraform.tfstate -o current-state.png

# Generate SVG for web documentation
ship tf diagram . --hcl --format svg -o architecture.svg

# Filter by provider (AWS only)
ship tf diagram terraform.tfstate --provider aws -o aws-resources.png

# Show all resources without filtering (raw mode)
ship tf diagram . --hcl --raw -o complete-diagram.png

# Real-world example
cd /path/to/your/terraform/project
ship tf diagram . --hcl -o docs/infrastructure-diagram.png
```

### 4. AI Assistant Integration (MCP)

Ship CLI includes a built-in MCP (Model Context Protocol) server that makes all functionality available to AI assistants like Claude Desktop and Cursor:

```bash
# Start MCP server for AI assistant integration
ship mcp

# Configure in Claude Desktop (claude_desktop_config.json):
{
  "mcpServers": {
    "ship-cli": {
      "command": "ship",
      "args": ["mcp"],
      "env": {
        "AWS_PROFILE": "your-profile"
      }
    }
  }
}
```

**What AI assistants can do with Ship CLI:**
- **Terraform Analysis**: "Analyze this Terraform code for costs and security"
- **Documentation**: "Generate docs for this Terraform module"
- **Cost Analysis**: "Estimate infrastructure costs for this project"
- **Security Scanning**: "Run security scans on this Terraform configuration"
- **Infrastructure Diagrams**: "Generate a visual diagram of this infrastructure"

**Available MCP Tools:**
- `lint` - Code linting and best practices
- `checkov` - Security analysis with Checkov
- `trivy` - Security analysis with Trivy
- `cost` - Cost estimation with OpenInfraQuote
- `docs` - Documentation generation
- `diagram` - Infrastructure diagram generation

**Pre-built Workflows:**
- `security_audit` - Comprehensive security audit process
- `cost_optimization` - Cost optimization analysis workflow

See the [MCP Integration Guide](docs/mcp-integration.md) for complete setup instructions.

### 5. CI/CD Integration

```yaml
# GitHub Actions Example
name: Terraform Analysis
on: [pull_request]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Ship CLI
        run: |
          go install github.com/cloudshipai/ship/cmd/ship@latest
      
      - name: Run Security Scan
        run: ship tf checkov
      
      - name: Estimate Costs
        run: ship tf cost
        env:
          INFRACOST_API_KEY: ${{ secrets.INFRACOST_API_KEY }}
```

## 🛠️ Available Tools

| Tool | Command | Description | Docker Image |
|------|---------|-------------|--------------|
| **InfraMap** | `ship tf diagram` | Infrastructure diagram generation | `cycloid/inframap:latest` |
| **TFLint** | `ship tf lint` | Terraform linter for syntax and best practices | `ghcr.io/terraform-linters/tflint` |
| **Checkov** | `ship tf checkov` | Comprehensive security and compliance scanner | `bridgecrew/checkov` |
| **Infracost** | `ship tf infracost` | Cloud cost estimation with breakdown | `infracost/infracost` |
| **Trivy** | `ship tf trivy` | Vulnerability scanner for IaC | `aquasec/trivy` |
| **terraform-docs** | `ship tf docs` | Auto-generate module documentation | `quay.io/terraform-docs/terraform-docs` |
| **OpenInfraQuote** | `ship tf cost` | Alternative cost analysis tool | `gruebel/openinfraquote` |

## 📋 Command Reference

### Module Management
```bash
# List all available modules (built-in, user, project)
ship modules list

# Show detailed information about a module
ship modules info terraform-tools

# Create a new custom module template
ship modules new my-custom-tool --type docker --description "My custom analysis tool"

# Filter modules by type or source
ship modules list --type docker --source user
ship modules list --trusted  # Show only trusted modules
```

### Infrastructure Diagrams
```bash
# Generate diagram from Terraform files
ship tf diagram . --hcl --format png

# Generate from state file
ship tf diagram terraform.tfstate --format svg

# Filter by provider
ship tf diagram . --hcl --provider aws --format png
```

### Linting
```bash
# Basic linting
ship tf lint

# Lint specific directory
ship tf lint ./modules/vpc

# Lint with custom config
ship tf lint --config .tflint.hcl
```

### Security Scanning
```bash
# Checkov scan (recommended)
ship tf checkov

# Trivy scan (alternative)
ship tf trivy

# Scan specific frameworks
ship tf checkov --framework terraform,arm
```

### Cost Estimation

#### Using Infracost
```bash
# Estimate costs for current directory
ship tf infracost

# Estimate with specific cloud provider
ship tf infracost --cloud aws

# Compare costs between branches
ship tf infracost --compare-to main
```

#### Using OpenInfraQuote (More Accurate)
```bash
# Analyze costs with OpenInfraQuote
ship tf cost

# Analyze specific plan file
ship tf cost terraform.tfplan.json

# Use specific AWS region for pricing
ship tf cost --region us-west-2
```

**OpenInfraQuote Features:**
- 🎯 **Accurate Pricing**: Uses real-time AWS pricing API data
- 📊 **Detailed Breakdown**: Shows costs per resource with hourly/monthly rates
- 🌍 **Region-Specific**: Accounts for regional price variations
- 📈 **100+ Resources**: Supports EC2, RDS, S3, ELB, Lambda, and more
- 🔄 **JSON Output**: Machine-readable format for automation

### Documentation
```bash
# Generate markdown documentation
ship tf docs

# Generate with specific filename
ship tf docs --filename API.md

# Save to specific output file
ship tf docs --output documentation.md
```

### Infrastructure Diagram Generation
```bash
# Generate PNG diagram from Terraform HCL files
ship tf diagram . --hcl --format png -o infrastructure.png

# Generate SVG diagram from state file
ship tf diagram terraform.tfstate --format svg -o current-state.svg

# Generate DOT format for programmatic processing
ship tf diagram . --hcl --format dot -o infrastructure.dot

# Filter by specific cloud provider
ship tf diagram . --hcl --provider aws --format png

# Generate PDF for documentation
ship tf diagram . --hcl --format pdf -o docs/infrastructure.pdf
```

## 🔐 Authentication

### AWS
```bash
# Ship CLI automatically uses your AWS credentials from:
# 1. Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
# 2. AWS credentials file (~/.aws/credentials)
# 3. IAM role (when running on EC2/ECS/Lambda)
```

### Azure
```bash
# Set Azure credentials
export ARM_CLIENT_ID="your-client-id"
export ARM_CLIENT_SECRET="your-client-secret"
export ARM_SUBSCRIPTION_ID="your-subscription-id"
export ARM_TENANT_ID="your-tenant-id"
```

### GCP
```bash
# Set GCP credentials
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
```

### Infracost
```bash
# Get free API key
infracost auth login

# Or set directly
export INFRACOST_API_KEY="your-api-key"
```

## 🏗️ Architecture

Ship CLI uses Dagger to run all tools in containers, providing:
- **Consistency**: Same tool versions across all environments
- **Isolation**: No conflicts with local installations
- **Security**: Tools run in sandboxed containers
- **Simplicity**: No need to install or manage tool versions

## 🤝 Contributing

We welcome contributions! See our [Contributing Guide](CONTRIBUTING.md) for details.

### Adding New Tools

1. Create a new module in `internal/dagger/modules/`
2. Add CLI command in `internal/cli/`
3. Update documentation
4. Submit a pull request

## 📚 Documentation

- [CLI Reference](docs/cli-reference.md) - Complete command reference
- [MCP Integration Guide](docs/mcp-integration.md) - AI assistant integration setup
- [Dynamic Module Discovery](docs/dynamic-module-discovery.md) - Extensible module system
- [Dagger Modules](docs/dagger-modules.md) - How to add new tools
- [Development Guide](docs/development-tasks.md) - For contributors
- [Technical Spec](docs/technical-spec.md) - Architecture and design

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run integration tests
go test -v ./internal/dagger/modules/

# Test specific module
go test -v -run TestTFLintModule ./internal/dagger/modules/
```

## 🧩 Using External Dagger Modules

Ship CLI is designed to be extensible! You can use any Dagger module without modifying Ship CLI itself.

### Using Published Dagger Modules

```bash
# Use any Dagger module directly
ship run dagger call --mod github.com/username/my-module@v1.0.0 analyze --source .

# Example: Using a custom security scanner
ship run dagger call --mod github.com/security/scanner@latest scan \
  --directory . \
  --severity high

# Example: Custom cost analyzer
ship run dagger call --mod github.com/finops/analyzer@v2.1.0 estimate \
  --terraform-dir . \
  --currency USD
```

### Creating Your Own Dagger Module

1. **Initialize a new Dagger module:**
```bash
dagger init --sdk=go my-custom-tool
cd my-custom-tool
```

2. **Define your tool's functionality:**
```go
// main.go
package main

import (
    "context"
    "dagger.io/dagger"
)

type MyCustomTool struct{}

// Analyze runs custom analysis on source code
func (m *MyCustomTool) Analyze(
    ctx context.Context,
    // Directory containing code to analyze
    source *dagger.Directory,
    // +optional
    // Output format (json, text, markdown)
    format string,
) (string, error) {
    return dag.Container().
        From("alpine:latest").
        WithMountedDirectory("/src", source).
        WithWorkdir("/src").
        WithExec([]string{"your-analysis-command", "--format", format}).
        Stdout(ctx)
}
```

3. **Publish your module:**
```bash
# Push to GitHub
git init
git add .
git commit -m "Initial module"
git remote add origin https://github.com/yourusername/my-custom-tool
git push -u origin main
git tag v1.0.0
git push --tags
```

4. **Use your module with Ship CLI:**
```bash
# Now anyone can use your module!
ship run dagger call --mod github.com/yourusername/my-custom-tool@v1.0.0 \
  analyze --source . --format json
```

### Module Ideas We'd Love to See

- **Cloud Security Scanner**: Deep security analysis for AWS/Azure/GCP
- **Kubernetes Analyzer**: K8s manifest validation and cluster analysis
- **Database Tools**: Schema validation, migration checks, documentation
- **Performance Profiler**: Infrastructure performance analysis
- **Compliance Checkers**: SOC2, HIPAA, PCI-DSS validators
- **Custom Cost Analyzers**: Organization-specific cost allocation

### 🤝 Community

- **Share Your Modules**: Tag them with `#ship-cli` on GitHub
- **Get Help**: Open an [issue](https://github.com/cloudshipai/ship/issues)
- **Contribute**: See our [Contributing Guide](CONTRIBUTING.md)

## 📈 Roadmap

- [ ] Dynamic module discovery and installation (`ship modules install`)
- [ ] Support for Atlantis integration
- [ ] Policy as Code with Open Policy Agent
- [ ] Custom tool configurations
- [ ] Web UI for results visualization
- [ ] Integration with more cloud providers

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Ship CLI wouldn't be possible without these amazing open source projects:
- [Dagger](https://dagger.io) - For containerized execution
- [Cobra](https://github.com/spf13/cobra) - For CLI framework
- All the individual tool maintainers

---

**Built with ❤️ by the CloudshipAI team**