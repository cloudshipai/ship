# Ship CLI

CloudshipAI CLI - A powerful command-line tool that brings enterprise-grade infrastructure analysis tools to your fingertips, all running in containers without local installations.

> **🤖 For LLMs and AI Assistants**: Complete installation and usage instructions specifically designed for AI consumption are available at [llms.txt](https://raw.githubusercontent.com/cloudshipai/ship/main/llms.txt). This includes MCP server setup, integration examples, and best practices for AI-driven infrastructure analysis.

## 🚀 Features

- **🔍 Terraform Linting**: Catch errors and enforce best practices with TFLint
- **🛡️ Security Scanning**: Multi-cloud security analysis with Checkov and Trivy
- **💰 Cost Estimation**: Estimate infrastructure costs with Infracost and OpenInfraQuote
- **📝 Documentation Generation**: Auto-generate beautiful Terraform module documentation
- **📊 Infrastructure Diagrams**: Visualize your infrastructure with InfraMap integration
- **🧠 AI-Powered Infrastructure Investigation**: Query your cloud infrastructure using natural language
- **🔎 Real-time Cloud Analysis**: Investigate live AWS, Azure, and GCP resources with Steampipe
- **🤖 AI Assistant Integration**: Built-in MCP server for Claude Desktop, Cursor, and other AI tools
- **🔌 Extensible Module System**: Add custom tools and Dagger functions without modifying core CLI
- **🐳 Containerized Tools**: All tools run in containers via Dagger - no local installations needed
- **☁️ Cloud Integration**: Seamlessly works with AWS, Azure, GCP, and other cloud providers
- **🔧 CI/CD Ready**: Perfect for integration into your existing pipelines

## 📦 Installation

### Quick Install (Recommended)

```bash
# Linux/macOS - Install latest release with wget
wget -qO- https://github.com/cloudshipai/ship/releases/latest/download/ship_$(uname -s)_$(uname -m).tar.gz | tar xz && sudo mv ship /usr/local/bin/

# Or with curl
curl -sSL https://github.com/cloudshipai/ship/releases/latest/download/ship_$(uname -s)_$(uname -m).tar.gz | tar xz && sudo mv ship /usr/local/bin/

# Verify installation
ship version
```

### Platform-Specific Downloads

```bash
# Linux x86_64
wget https://github.com/cloudshipai/ship/releases/latest/download/ship_Linux_x86_64.tar.gz

# macOS Intel
wget https://github.com/cloudshipai/ship/releases/latest/download/ship_Darwin_x86_64.tar.gz

# macOS Apple Silicon
wget https://github.com/cloudshipai/ship/releases/latest/download/ship_Darwin_arm64.tar.gz

# Linux ARM64
wget https://github.com/cloudshipai/ship/releases/latest/download/ship_Linux_arm64.tar.gz

# Windows
wget https://github.com/cloudshipai/ship/releases/latest/download/ship_Windows_x86_64.zip
```

### From Source
```bash
# Clone the repository
git clone https://github.com/cloudshipai/ship.git
cd ship

# Build and install
go build -o ship ./cmd/ship
sudo mv ship /usr/local/bin/

# Verify installation
ship version
```

### Using Go Install
```bash
go install github.com/cloudshipai/ship/cmd/ship@latest
```

## 🏃 Quick Start

### 1. Basic Usage

```bash
# Navigate to your Terraform project
cd your-terraform-project

# Run a comprehensive analysis
ship terraform-tools lint                # Check for errors and best practices
ship terraform-tools checkov-scan        # Security scanning
ship terraform-tools cost-estimate       # Estimate AWS/Azure/GCP costs
ship terraform-tools generate-docs       # Generate documentation
```

### 2. Real-World Example

```bash
# Clone a sample Terraform project
git clone https://github.com/terraform-aws-modules/terraform-aws-vpc.git
cd terraform-aws-vpc/examples/simple

# Run all tools
ship terraform-tools lint
ship terraform-tools checkov-scan
ship terraform-tools security-scan
ship terraform-tools cost-estimate
ship terraform-tools generate-docs > README.md
ship terraform-tools generate-diagram . --hcl -o infrastructure.png
```

### 3. Generate Infrastructure Diagrams

Visualize your infrastructure with InfraMap integration:

```bash
# Generate diagram from Terraform files (no state file needed!)
ship terraform-tools generate-diagram . --hcl --format png -o infrastructure.png

# Generate from existing state file
ship terraform-tools generate-diagram terraform.tfstate -o current-state.png

# Generate SVG for web documentation
ship terraform-tools generate-diagram . --hcl --format svg -o architecture.svg

# Filter by provider (AWS only)
ship terraform-tools generate-diagram terraform.tfstate --provider aws -o aws-resources.png

# Show all resources without filtering (raw mode)
ship terraform-tools generate-diagram . --hcl --raw -o complete-diagram.png

# Real-world example
cd /path/to/your/terraform/project
ship terraform-tools generate-diagram . --hcl -o docs/infrastructure-diagram.png
```

### 4. AI-Powered Infrastructure Investigation

Ship CLI offers multiple AI-powered approaches to analyze and investigate your infrastructure:

#### Basic AI Investigation

Query your live cloud infrastructure using natural language with Steampipe-powered analysis:

```bash
# Configure AWS credentials (Ship CLI will use your existing AWS config)
export AWS_PROFILE=your-profile  # or use default

# Ask questions about your infrastructure in natural language
ship ai-investigate --prompt "Show me all my S3 buckets with their creation dates and regions" --execute

ship ai-investigate --prompt "Check for security issues in my AWS account" --execute

ship ai-investigate --prompt "List all running EC2 instances with their IP addresses" --execute

ship ai-investigate --prompt "Show me any unused or idle resources that might be costing money" --execute

ship ai-investigate --prompt "Find all publicly accessible RDS instances" --execute
```

#### Autonomous AI Agent

Let the AI agent autonomously investigate your infrastructure using multiple tools:

```bash
# Run a comprehensive security audit with autonomous decision-making
ship ai-agent --task "Perform complete security audit of AWS infrastructure"

# Cost optimization with detailed analysis
ship ai-agent --task "Optimize costs for our production environment" --max-steps 15

# Documentation and compliance check
ship ai-agent --task "Document all Terraform modules and check for compliance issues"

# Interactive mode - approve each tool use
ship ai-agent --task "Analyze security posture and fix critical issues" --approve-each
```

#### Microservices-Based AI Investigation

Run AI investigation with each tool as a separate scalable service:

```bash
# Launch AI with microservices architecture
ship ai-services --task "Audit security across all AWS resources"

# Show service endpoints for debugging
ship ai-services --task "Generate cost report with optimization recommendations" --show-endpoints

# Keep services running for other tools to use
ship ai-services --task "Document infrastructure and analyze patterns" --keep-services

# Export service endpoints for integration
ship ai-services --task "Full infrastructure analysis" --export-endpoints services.json
```

#### How it works:
1. **Natural Language Processing**: Ship analyzes your prompt to understand what you're looking for
2. **Dynamic Query Generation**: Automatically generates appropriate Steampipe SQL queries
3. **Multi-Step Investigation**: Creates comprehensive investigation plans with multiple related queries
4. **Real-Time Analysis**: Executes queries against your live cloud infrastructure
5. **Intelligent Insights**: Provides security findings, cost optimization tips, and actionable recommendations
6. **Tool Orchestration**: AI agents can autonomously use Steampipe, OpenInfraQuote, Terraform-docs, InfraMap, and security scanners
7. **Service Architecture**: Optional microservices mode for enterprise-scale deployments

#### Supported AI Capabilities:
- **Security Analysis**: "Check for security vulnerabilities", "Find open security groups", "Show unencrypted resources"
- **Cost Optimization**: "Find unused resources", "Show expensive instances", "Identify idle resources"
- **Resource Inventory**: "List all S3 buckets", "Show running instances", "Find RDS databases"
- **Infrastructure Visualization**: "Generate infrastructure diagram", "Create visual documentation", "Diagram dependencies"
- **Compliance Checks**: "Check encryption status", "Verify MFA settings", "Audit logging configuration"
- **Autonomous Investigation**: AI agent can chain multiple tools to solve complex problems
- **Service-Based Architecture**: Run tools as HTTP services for better scalability

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
- **Infrastructure Investigation**: "Check my AWS account for security issues"
- **Terraform Analysis**: "Analyze this Terraform code for costs and security"
- **Cost Optimization**: "Find unused resources in my cloud account"
- **Documentation**: "Generate docs for this Terraform module"
- **Compliance Audits**: "Run a compliance check on my infrastructure"

**Available MCP Tools:**
- `ai_investigate` - Natural language infrastructure investigation
- `terraform_lint` - Code linting and best practices
- `terraform_security_scan` - Security analysis
- `terraform_cost_estimate` - Cost estimation
- `terraform_generate_docs` - Documentation generation
- `cloudship_push` - Upload artifacts for AI analysis

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
          wget https://github.com/cloudship/ship/releases/latest/download/ship-linux-amd64
          chmod +x ship-linux-amd64
          sudo mv ship-linux-amd64 /usr/local/bin/ship
      
      - name: Run Security Scan
        run: ship terraform-tools checkov-scan
      
      - name: Estimate Costs
        run: ship terraform-tools cost-estimate
        env:
          INFRACOST_API_KEY: ${{ secrets.INFRACOST_API_KEY }}
```

## 🛠️ Available Tools

| Tool | Command | Description | Docker Image |
|------|---------|-------------|--------------|
| **Steampipe + AI** | `ship ai-investigate` | AI-powered cloud infrastructure investigation | `turbot/steampipe:latest` |
| **TFLint** | `ship terraform-tools lint` | Terraform linter for syntax and best practices | `ghcr.io/terraform-linters/tflint` |
| **Checkov** | `ship terraform-tools checkov-scan` | Comprehensive security and compliance scanner | `bridgecrew/checkov` |
| **Infracost** | `ship terraform-tools cost-estimate` | Cloud cost estimation with breakdown | `infracost/infracost` |
| **Trivy** | `ship terraform-tools security-scan` | Vulnerability scanner for IaC | `aquasec/trivy` |
| **terraform-docs** | `ship terraform-tools generate-docs` | Auto-generate module documentation | `quay.io/terraform-docs/terraform-docs` |
| **OpenInfraQuote** | `ship terraform-tools cost-analysis` | Alternative cost analysis tool | `gruebel/openinfraquote` |

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

### AI-Powered Investigation
```bash
# Natural language infrastructure investigation
ship ai-investigate --prompt "Show me all S3 buckets" --execute
ship ai-investigate --prompt "Check for security issues" --execute  
ship ai-investigate --prompt "Find unused resources costing money" --execute

# Use specific AWS profile/region
ship ai-investigate --prompt "List running instances" --aws-profile prod --aws-region us-west-2 --execute

# Preview queries without execution
ship ai-investigate --prompt "Security audit" --provider aws
```

### Linting
```bash
# Basic linting
ship terraform-tools lint

# Lint specific directory
ship terraform-tools lint ./modules/vpc

# Lint with custom config
ship terraform-tools lint --config .tflint.hcl
```

### Security Scanning
```bash
# Checkov scan (recommended)
ship terraform-tools checkov-scan

# Trivy scan (alternative)
ship terraform-tools security-scan

# Scan specific frameworks
ship terraform-tools checkov-scan --framework terraform,arm
```

### Cost Estimation
```bash
# Estimate costs for current directory
ship terraform-tools cost-estimate

# Estimate with specific cloud provider
ship terraform-tools cost-estimate --cloud aws

# Compare costs between branches
ship terraform-tools cost-estimate --compare-to main
```

### Documentation
```bash
# Generate markdown documentation
ship terraform-tools generate-docs

# Generate JSON output
ship terraform-tools generate-docs --format json

# Include examples in docs
ship terraform-tools generate-docs --show-examples
```

### AI Infrastructure Investigation
```bash
# Basic investigation with natural language
ship ai-investigate --prompt "Show me my S3 buckets"

# Execute the generated queries (add --execute to run)
ship ai-investigate --prompt "Check for security issues" --execute

# Use specific cloud provider
ship ai-investigate --prompt "List running instances" --provider aws --execute

# Use specific AWS profile and region
ship ai-investigate --prompt "Find unused EBS volumes" --aws-profile production --aws-region us-west-2 --execute

# Cost analysis investigation
ship ai-investigate --prompt "Show me expensive resources that might be optimized" --execute

# Security-focused investigation
ship ai-investigate --prompt "Find all publicly accessible resources" --execute

# Compliance investigation
ship ai-investigate --prompt "Check encryption status across all resources" --execute
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

## 🧩 Creating Custom Modules (Community Extensions)

Ship CLI is designed to be extensible! You can create custom Dagger modules to add new tools and capabilities.

### Quick Example: Create Your Own Tool

```go
// my-tool/main.go
package main

import (
    "context"
    "dagger.io/dagger"
)

type MyTool struct{}

func (m *MyTool) Analyze(ctx context.Context, source *dagger.Directory) (string, error) {
    return dag.Container().
        From("alpine:latest").
        WithMountedDirectory("/src", source).
        WithExec([]string{"analyze-command"}).
        Stdout(ctx)
}
```

### 📚 Full Documentation

See [CONTRIBUTING.md](CONTRIBUTING.md) for a complete guide on:
- Creating custom Dagger modules
- Module structure and metadata
- AI/LLM integration for your tools
- Testing and publishing modules
- Joining the community registry

### 🎯 Module Ideas We'd Love to See

- **Cloud Security Scanner**: Deep security analysis for AWS/Azure/GCP
- **Kubernetes Analyzer**: K8s manifest validation and cluster analysis
- **Database Tools**: Schema validation, migration checks, documentation
- **Performance Profiler**: Infrastructure performance analysis
- **Compliance Checkers**: SOC2, HIPAA, PCI-DSS validators
- **Custom Cost Analyzers**: Organization-specific cost allocation

### 🤝 Community

- **Submit Your Module**: Add to our [community registry](https://github.com/cloudshipai/community-modules)
- **Get Help**: Join our [Discord](https://discord.gg/cloudship)
- **Share Ideas**: Open an [issue](https://github.com/cloudshipai/ship/issues)

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