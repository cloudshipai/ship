# Ship CLI Reference

## Overview

Ship CLI is a command-line tool that provides Terraform analysis capabilities using containerized tools via Dagger.

## Installation

```bash
# Install from source
go install github.com/cloudship/ship/cmd/ship@latest

# Or download binary from releases
curl -L https://github.com/cloudship/ship/releases/latest/download/ship_$(uname -s)_$(uname -m) -o ship
chmod +x ship
sudo mv ship /usr/local/bin/
```

## Commands

### terraform-tools

The `terraform-tools` command provides various Terraform analysis capabilities.

#### Cost Analysis

Analyze Terraform costs using OpenInfraQuote:

```bash
# Analyze a directory
ship terraform-tools cost-analysis

# Analyze a specific plan file
ship terraform-tools cost-analysis plan.json

# Analyze a different directory
ship terraform-tools cost-analysis ./infrastructure
```

#### Security Scanning

##### Trivy Security Scan

Scan for security issues using Trivy:

```bash
# Scan current directory
ship terraform-tools security-scan

# Scan specific directory
ship terraform-tools security-scan ./modules
```

##### Checkov Security Scan

Comprehensive multi-cloud security scanning:

```bash
# Scan current directory
ship terraform-tools checkov-scan

# Scan specific directory
ship terraform-tools checkov-scan ./terraform
```

#### Linting

Lint Terraform code using TFLint:

```bash
# Lint current directory
ship terraform-tools lint

# Lint specific directory
ship terraform-tools lint ./environments/prod
```

#### Documentation Generation

Generate documentation using terraform-docs:

```bash
# Generate docs for current directory
ship terraform-tools generate-docs

# Generate docs for specific module
ship terraform-tools generate-docs ./modules/vpc
```

#### Cost Estimation

Estimate infrastructure costs using Infracost:

```bash
# Set API key (required for full functionality)
export INFRACOST_API_KEY=your-api-key

# Estimate costs for current directory
ship terraform-tools cost-estimate

# Estimate costs for specific directory
ship terraform-tools cost-estimate ./environments/staging
```

### dagger

Direct Dagger operations:

```bash
# Run a Dagger container
ship dagger run <image> <command>

# Show Dagger version
ship dagger version
```

### dagger-steampipe

Steampipe integration for cloud queries:

```bash
# Run a Steampipe query
ship dagger-steampipe query "select * from aws_s3_bucket"

# List available tables
ship dagger-steampipe tables
```

## Environment Variables

### Required for Some Features

- `INFRACOST_API_KEY`: Required for Infracost cost estimation (get one at https://www.infracost.io)
- `CLOUDSHIP_API_KEY`: Required for pushing artifacts to Cloudship

### Optional Cloud Credentials

The tools will automatically use cloud credentials if available:

- AWS: `~/.aws/credentials` or `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`
- Azure: `~/.azure` or `AZURE_*` environment variables
- GCP: `GOOGLE_APPLICATION_CREDENTIALS` pointing to service account JSON

## Output Formats

Most tools support JSON output for programmatic use:

```bash
# All commands output JSON by default for parsing
ship terraform-tools lint ./modules | jq '.issues[]'

# Cost analysis outputs JSON
ship terraform-tools cost-analysis | jq '.totalMonthlyCost'

# Security scans output JSON
ship terraform-tools checkov-scan | jq '.results.failed_checks'
```

## Examples

### Complete Terraform Module Analysis

```bash
# 1. Generate documentation
ship terraform-tools generate-docs ./modules/eks > README.md

# 2. Lint the code
ship terraform-tools lint ./modules/eks

# 3. Security scan with Trivy
ship terraform-tools security-scan ./modules/eks

# 4. Security scan with Checkov
ship terraform-tools checkov-scan ./modules/eks

# 5. Estimate costs
ship terraform-tools cost-estimate ./modules/eks
```

### CI/CD Pipeline Integration

```yaml
# GitHub Actions example
- name: Terraform Analysis
  run: |
    # Lint
    ship terraform-tools lint
    
    # Security scan
    ship terraform-tools checkov-scan
    
    # Cost estimation
    ship terraform-tools cost-estimate
```

### Pre-commit Hook

```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: terraform-lint
        name: Terraform Lint
        entry: ship terraform-tools lint
        language: system
        files: \.tf$
```

## Troubleshooting

### Common Issues

1. **"Dagger engine initialization failed"**
   - Ensure Docker is running
   - Check Docker permissions: `docker ps`

2. **"INFRACOST_API_KEY not set"**
   - Sign up at https://www.infracost.io
   - Export the API key: `export INFRACOST_API_KEY=your-key`

3. **"No Terraform files found"**
   - Ensure you're in a directory with `.tf` files
   - Specify the correct directory path

### Debug Mode

Enable debug logging:

```bash
export DAGGER_LOG_LEVEL=debug
ship terraform-tools lint --debug
```

## Tool Versions

Ship CLI uses the following containerized tools:

| Tool | Container Image | Purpose |
|------|----------------|---------|
| TFLint | `ghcr.io/terraform-linters/tflint:latest` | Terraform linting |
| Checkov | `bridgecrew/checkov:latest` | Security scanning |
| Infracost | `infracost/infracost:latest` | Cost estimation |
| Trivy | `aquasec/trivy:latest` | Security scanning |
| terraform-docs | `quay.io/terraform-docs/terraform-docs:latest` | Documentation |
| OpenInfraQuote | `ghcr.io/initech-consulting/openinfraquote:latest` | Cost analysis |

## Support

- GitHub Issues: https://github.com/cloudship/ship/issues
- Documentation: https://docs.cloudship.ai
- Community: https://community.cloudship.ai