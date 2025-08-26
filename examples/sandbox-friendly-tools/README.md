# Sandbox-Friendly Security Tools Examples

This directory contains examples demonstrating tools that can run in basic sandbox environments without requiring API keys or external authentication.

## Self-Contained Tools

### 1. Trivy - Vulnerability Scanner
- **Container**: `aquasec/trivy:latest`
- **Requirements**: None (uses built-in vulnerability database)
- **Examples**: Image scanning, filesystem scanning, SBOM generation

### 2. Checkov - Multi-Cloud Security Scanner
- **Container**: `bridgecrew/checkov:latest`
- **Requirements**: None (uses built-in security policies)
- **Examples**: Terraform scanning, Kubernetes manifest scanning

### 3. Nuclei - Vulnerability Scanner
- **Container**: `projectdiscovery/nuclei:latest`
- **Requirements**: None (uses `-disable-update-check` flag)
- **Examples**: URL scanning, template-based vulnerability detection

### 4. Gitleaks - Secret Scanner
- **Container**: `zricethezav/gitleaks:latest`
- **Requirements**: None (uses built-in secret patterns)
- **Examples**: Directory scanning, git repository scanning

### 5. Semgrep - Static Analysis
- **Container**: `semgrep/semgrep:latest`
- **Requirements**: None (uses built-in rule sets)
- **Examples**: Code scanning, OWASP Top 10 detection

### 6. Syft - SBOM Generator
- **Container**: `anchore/syft:latest`
- **Requirements**: None (generates SBOMs from source)
- **Examples**: Directory scanning, container image SBOM generation

### 7. Grype - Vulnerability Scanner
- **Container**: `anchore/grype:latest`
- **Requirements**: None (uses built-in vulnerability database)
- **Examples**: SBOM vulnerability scanning, image vulnerability scanning

## Testing Examples

Each subdirectory contains:
- Sample target files/projects for scanning
- Expected output examples
- Command-line usage demonstrations
- Container execution examples

## API-Dependent Tools (Not Sandbox-Compatible)

These tools require external authentication and are not suitable for basic sandboxes:
- **Infracost**: Requires INFRACOST_API_KEY + cloud credentials
- **Snyk**: Requires SNYK_TOKEN environment variable
- **Prowler**: Requires cloud provider credentials
- **Scout Suite**: Requires cloud authentication
- **Cloudsplaining**: Requires AWS credentials