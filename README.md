# Ship CLI

CloudshipAI CLI - A powerful command-line tool that brings enterprise-grade infrastructure analysis tools to your fingertips, all running in containers without local installations.

## 🚀 Features

- **🔍 Terraform Linting**: Catch errors and enforce best practices with TFLint
- **🛡️ Security Scanning**: Multi-cloud security analysis with Checkov and Trivy
- **💰 Cost Estimation**: Estimate infrastructure costs with Infracost and OpenInfraQuote
- **📝 Documentation Generation**: Auto-generate beautiful Terraform module documentation
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
```

### 3. CI/CD Integration

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
| **TFLint** | `ship terraform-tools lint` | Terraform linter for syntax and best practices | `ghcr.io/terraform-linters/tflint` |
| **Checkov** | `ship terraform-tools checkov-scan` | Comprehensive security and compliance scanner | `bridgecrew/checkov` |
| **Infracost** | `ship terraform-tools cost-estimate` | Cloud cost estimation with breakdown | `infracost/infracost` |
| **Trivy** | `ship terraform-tools security-scan` | Vulnerability scanner for IaC | `aquasec/trivy` |
| **terraform-docs** | `ship terraform-tools generate-docs` | Auto-generate module documentation | `quay.io/terraform-docs/terraform-docs` |
| **OpenInfraQuote** | `ship terraform-tools cost-analysis` | Alternative cost analysis tool | `gruebel/openinfraquote` |

## 📋 Command Reference

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

## 📈 Roadmap

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