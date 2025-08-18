# Checkov MCP Tool

Checkov is a static code analysis tool for infrastructure as code (IaC) that scans cloud infrastructure configurations to find misconfigurations before they're deployed.

## Description

Checkov provides comprehensive security and compliance scanning for:
- Terraform, CloudFormation, Kubernetes, Helm charts
- Dockerfile, Serverless framework
- ARM templates, OpenAPI specifications
- Multiple compliance frameworks (CIS, PCI-DSS, HIPAA, SOC2, etc.)

## MCP Functions

### `checkov_scan_directory`
Scan directory for IaC security issues using real checkov CLI.

**Parameters:**
- `directory`: Directory to scan (default: current directory)
- `framework`: Specific framework to scan (terraform, cloudformation, kubernetes, etc.)
- `check`: Specific check ID to run
- `skip_check`: Comma-separated list of checks to skip
- `output_format`: Output format (json, junitxml, github_failed_only, sarif, csv)

**CLI Command:** `checkov -d <directory> [--framework <framework>] [--check <check>] [--skip-check <skip_check>] [-o <output_format>]`

### `checkov_scan_file`
Scan specific file using real checkov CLI.

**Parameters:**
- `file` (required): File path to scan
- `framework`: Specific framework to scan
- `check`: Specific check ID to run
- `skip_check`: Comma-separated list of checks to skip
- `output_format`: Output format (json, junitxml, github_failed_only, sarif, csv)

**CLI Command:** `checkov -f <file> [--framework <framework>] [--check <check>] [--skip-check <skip_check>] [-o <output_format>]`

### `checkov_scan_terraform_plan`
Scan Terraform plan file using real checkov CLI.

**Parameters:**
- `plan_file` (required): Path to Terraform plan JSON file
- `check`: Specific check ID to run
- `skip_check`: Comma-separated list of checks to skip
- `output_format`: Output format

**CLI Command:** `checkov -f <plan_file> --framework terraform_plan [--check <check>] [--skip-check <skip_check>] [-o <output_format>]`

### `checkov_scan_with_baseline`
Scan with baseline for comparison using real checkov CLI.

**Parameters:**
- `directory`: Directory to scan
- `baseline`: Path to baseline file
- `output_format`: Output format

**CLI Command:** `checkov -d <directory> --baseline <baseline> [-o <output_format>]`

### `checkov_list_checks`
List all available checks using real checkov CLI.

**Parameters:**
- `framework`: Filter checks by framework

**CLI Command:** `checkov --list [--framework <framework>]`

### `checkov_create_baseline`
Create baseline file using real checkov CLI.

**Parameters:**
- `directory`: Directory to scan for baseline
- `output_file`: Output baseline file path

**CLI Command:** `checkov -d <directory> --create-baseline --output-baseline <output_file>`

### `checkov_scan_with_external_checks`
Scan with external custom checks using real checkov CLI.

**Parameters:**
- `directory`: Directory to scan
- `external_checks_dir`: Directory containing custom checks
- `output_format`: Output format

**CLI Command:** `checkov -d <directory> --external-checks-dir <external_checks_dir> [-o <output_format>]`

### `checkov_scan_repository`
Scan git repository using real checkov CLI.

**Parameters:**
- `repo_url` (required): Git repository URL
- `branch`: Branch to scan
- `framework`: Specific framework to scan

**CLI Command:** `checkov --repo-url <repo_url> [--branch <branch>] [--framework <framework>]`

### `checkov_scan_docker_image`
Scan Docker image using real checkov CLI.

**Parameters:**
- `image` (required): Docker image name
- `dockerfile_path`: Path to Dockerfile
- `output_format`: Output format

**CLI Command:** `checkov --docker-image <image> [--dockerfile-path <dockerfile_path>] [-o <output_format>]`

### `checkov_version`
Get Checkov version information.

**CLI Command:** `checkov --version`

## Common Use Cases

1. **Pre-deployment Scanning**: Scan IaC before deployment
2. **CI/CD Integration**: Automated security checks in pipelines
3. **Compliance Validation**: Check against compliance frameworks
4. **Custom Policies**: Apply organization-specific security policies
5. **Drift Detection**: Compare against baselines

## Integration with Ship CLI

All Checkov tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Perform automated IaC security scanning
- Validate compliance requirements
- Generate security reports
- Apply custom security policies

The tools use containerized execution via Dagger for consistent, isolated security scanning.