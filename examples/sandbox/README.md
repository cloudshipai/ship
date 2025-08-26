# Ship Sandbox Environments

This directory contains sandbox environments for testing various Ship tools and their capabilities.

## Available Sandboxes

### 1. Vulnerable App (`vulnerable-app/`)
**Purpose**: Test security scanning tools (Trivy, Checkov)
**Tools**: Trivy, Checkov
**Focus**: Security vulnerability detection, infrastructure security scanning

Contains intentionally vulnerable Terraform configuration including:
- Public S3 buckets
- Unencrypted storage
- Overly permissive security groups
- Hardcoded passwords
- Public RDS instances

### 2. Cost Demo (`cost-demo/`)
**Purpose**: Test cost analysis tools
**Tools**: OpenInfraQuote
**Focus**: Infrastructure cost estimation and analysis

Contains multi-tier infrastructure with:
- Various EC2 instance types (micro, medium, large)
- RDS database instances
- Load balancers
- Different storage configurations

### 3. Container App (`container-app/`)
**Purpose**: Test SBOM generation and container scanning
**Tools**: Syft, Trivy
**Focus**: Software Bill of Materials, supply chain analysis, container security

Contains a Node.js application with:
- Vulnerable dependencies (lodash 4.17.20)
- Standard web application stack
- Dockerfile for containerization
- Package.json with various dependencies

### 4. Terraform Quality (`terraform-quality/`)
**Purpose**: Test Terraform linting and quality analysis
**Tools**: TFLint
**Focus**: Code quality, best practices, Terraform standards

Contains Terraform configuration with common issues:
- Hardcoded values
- Missing variable descriptions
- Deprecated arguments
- Poor naming conventions
- Missing security configurations

## Running Tests

Each sandbox can be tested using the Ship CLI:

```bash
# Security scanning
ship terraform-tools checkov-scan --directory examples/sandbox/vulnerable-app
ship terraform-tools security-scan --directory examples/sandbox/vulnerable-app

# Cost analysis
ship terraform-tools cost-analysis --directory examples/sandbox/cost-demo

# SBOM generation
ship terraform-tools syft-generate-sbom --target examples/sandbox/container-app

# Quality analysis
ship terraform-tools lint --directory examples/sandbox/terraform-quality
```

## Test Expectations

### Security Tools
- **Checkov**: Should detect 8+ security issues in vulnerable-app
- **Trivy**: Should find vulnerabilities in container dependencies

### Cost Tools
- **OpenInfraQuote**: Should provide cost estimates for multi-tier infrastructure

### Quality Tools
- **TFLint**: Should identify 10+ code quality issues in terraform-quality

### SBOM Tools
- **Syft**: Should generate comprehensive SBOM with package inventory

## Testing Each Sandbox

### Manual Testing

#### 1. Vulnerable App Testing
```bash
# Test Checkov security scanning
cd examples/sandbox/vulnerable-app
ship terraform-tools checkov-scan --directory .

# Expected results:
# - 8+ HIGH/CRITICAL security findings
# - Unencrypted storage violations
# - Public access violations
# - Hardcoded secrets detection

# Test with different output formats
ship terraform-tools checkov-scan --directory . --output json
ship terraform-tools checkov-scan --directory . --output sarif

# Test Trivy infrastructure scanning
ship terraform-tools security-scan --directory .
```

#### 2. Cost Demo Testing
```bash
# Test cost analysis with OpenInfraQuote
cd examples/sandbox/cost-demo

# Generate Terraform plan first
terraform init
terraform plan -out=tfplan.binary
terraform show -json tfplan.binary > tfplan.json

# Run cost analysis
ship terraform-tools cost-analysis --plan-file tfplan.json --region us-east-1

# Expected results:
# - Cost breakdown for multiple instance types
# - Database cost estimates
# - Load balancer pricing
# - Storage cost analysis
```

#### 3. Container App Testing
```bash
# Test SBOM generation
cd examples/sandbox/container-app

# Generate SBOM for directory
ship terraform-tools syft-generate-sbom --target .

# Build container image for scanning
docker build -t sandbox-app:latest .

# Generate SBOM for container
ship terraform-tools syft-generate-sbom --target docker:sandbox-app:latest

# Test vulnerability scanning
ship terraform-tools trivy-scan-image --image sandbox-app:latest

# Expected results:
# - Package inventory with Node.js dependencies
# - Vulnerable lodash package detection
# - Container layer analysis
# - License information
```

#### 4. Terraform Quality Testing
```bash
# Test TFLint quality analysis
cd examples/sandbox/terraform-quality

# Run comprehensive linting
ship terraform-tools lint --directory .

# Test with different configurations
ship terraform-tools lint --directory . --enable-rule terraform_deprecated_interpolation
ship terraform-tools lint --directory . --format json

# Expected results:
# - 10+ code quality violations
# - Hardcoded value warnings
# - Missing variable descriptions
# - Deprecated argument usage
# - Naming convention violations
```

### Automated Testing Strategy

#### Test Validation Criteria

Each sandbox should validate:

1. **Tool Execution Success**
   - Command exits with expected code (0 or 1 for findings)
   - Output is generated and parseable
   - No runtime errors or crashes

2. **Expected Findings**
   - Security tools detect known vulnerabilities
   - Cost tools provide numeric estimates
   - Quality tools identify code issues
   - SBOM tools generate package inventories

3. **Output Format Validation**
   - JSON output is valid JSON
   - SARIF output follows SARIF schema
   - Human-readable output contains expected sections

#### Sample Test Commands

```bash
# Quick validation of all sandboxes
./scripts/test-sandboxes.sh

# Individual sandbox validation
go test -v ./internal/dagger/modules -run TestVulnerableAppSandbox
go test -v ./internal/dagger/modules -run TestCostDemoSandbox
go test -v ./internal/dagger/modules -run TestContainerAppSandbox
go test -v ./internal/dagger/modules -run TestTerraformQualitySandbox
```

### Continuous Integration Testing

The sandboxes support CI/CD validation:

1. **Pre-commit Testing**: Validate sandbox integrity
2. **Integration Testing**: Test tool execution in containers
3. **Regression Testing**: Ensure consistent tool behavior
4. **Performance Testing**: Monitor execution times

### Troubleshooting

#### Common Issues

1. **Docker Not Available**
   - Ensure Docker daemon is running
   - Check Dagger client connectivity

2. **Missing Dependencies**
   - Some tools require specific setup (AWS credentials, etc.)
   - Check tool-specific requirements

3. **Output Parsing Failures**
   - Validate JSON with `jq` or similar tools
   - Check for unexpected error messages in output

#### Debug Commands

```bash
# Enable verbose output
ship terraform-tools checkov-scan --directory examples/sandbox/vulnerable-app --verbose

# Check Dagger status
dagger version

# Validate Terraform syntax
terraform validate examples/sandbox/vulnerable-app
```

## Integration Testing

These sandboxes are used by the Go test suite in `internal/dagger/modules/sandbox_test.go` to validate tool functionality and ensure consistent behavior across updates.