# Ship CLI Terraform Examples

These examples demonstrate how to use Ship CLI to analyze Terraform configurations of varying complexity.

## Available Examples

### 1. Easy: S3 Bucket (`easy-s3-bucket/`)
- **What it creates**: A single S3 bucket with security best practices
- **Concepts**: Basic resource creation, encryption, lifecycle rules
- **Estimated cost**: ~$1-5/month depending on storage
- **Good for**: Learning Ship CLI basics, simple security scans

### 2. Medium: Web Application (`medium-web-app/`)
- **What it creates**: Scalable web app with VPC, ALB, and Auto Scaling
- **Concepts**: Networking, load balancing, auto scaling
- **Estimated cost**: ~$125/month
- **Good for**: Testing cost analysis, security scanning, diagram generation

### 3. Complex: Multi-Region Production (`complex-multi-region/`)
- **What it creates**: Full production architecture across 2 regions
- **Concepts**: Modules, multi-region, disaster recovery, monitoring
- **Estimated cost**: ~$326/month
- **Good for**: Comprehensive testing of all Ship CLI features

## Quick Start

```bash
# Clone the repository
git clone https://github.com/cloudshipai/ship.git
cd ship/examples/terraform

# Choose an example
cd easy-s3-bucket  # or medium-web-app, or complex-multi-region

# Run Ship CLI commands
ship terraform-tools lint
ship terraform-tools security-scan
ship terraform-tools cost-estimate
ship terraform-tools generate-diagram . --hcl -o diagram.png
```

## Testing All Ship CLI Features

### 1. Linting and Code Quality
```bash
ship terraform-tools lint
```

### 2. Security Scanning
```bash
# Multiple security scanners available
ship terraform-tools security-scan    # InfraScan
ship terraform-tools checkov-scan     # Checkov
```

### 3. Cost Analysis
```bash
# Different cost estimation tools
ship terraform-tools cost-estimate    # Infracost
ship terraform-tools cost-analysis    # OpenInfraQuote
```

### 4. Documentation
```bash
ship terraform-tools generate-docs
```

### 5. Infrastructure Visualization
```bash
ship terraform-tools generate-diagram . --hcl -o infrastructure.png
```

### 6. Push to CloudShip
```bash
# First authenticate
ship auth --api-key YOUR_API_KEY

# Run with automatic push
ship terraform-tools security-scan --push
ship terraform-tools cost-estimate --push --push-tags "example,test"
```

## Example Outputs

Each example includes expected outputs:
- **Lint results**: Code quality issues and best practices
- **Security findings**: Potential vulnerabilities and misconfigurations
- **Cost estimates**: Monthly cost breakdowns
- **Documentation**: Auto-generated README files
- **Diagrams**: Visual infrastructure representations

## For LLMs and AI Assistants

If you're an AI assistant helping a user test Ship CLI:

1. Start with the easy example to verify installation
2. Use the medium example to test all basic features
3. Use the complex example to demonstrate advanced capabilities
4. Always run `ship terraform-tools lint` first to catch syntax errors
5. Generate diagrams to help users visualize their infrastructure
6. Use the `--push` flag to demonstrate CloudShip integration

## Notes

- These examples use AWS provider configurations
- No actual resources are created unless you run `terraform apply`
- All examples follow Terraform best practices
- Each example has its own detailed README

## Contributing

To add new examples:
1. Create a new directory with descriptive name
2. Include main.tf, variables.tf, outputs.tf
3. Add a comprehensive README.md
4. Test all Ship CLI commands
5. Document expected costs and outputs