# Grype

Vulnerability scanner for container images and filesystems.

## Description

Grype is a vulnerability scanner for container images and filesystems developed by Anchore. It scans for vulnerabilities in OS packages, language-specific packages, and application dependencies. Grype uses a comprehensive vulnerability database and supports various input sources including container images, directories, archives, and SBOMs.

## MCP Tools

### Vulnerability Scanning
- **`grype_scan`** - Scan target for vulnerabilities (image, directory, archive, SBOM)
- **`grype_explain`** - Get detailed information about specific CVE from scan results

### Database Management
- **`grype_db_status`** - Report current status of vulnerability database
- **`grype_db_check`** - Check if database updates are available
- **`grype_db_update`** - Update vulnerability database to latest version
- **`grype_db_list`** - Show databases available for download
- **`grype_db_import`** - Import vulnerability database archive

### Utility
- **`grype_version`** - Display Grype version information

## Real CLI Commands Used

### Core Scanning Commands
- `grype <target>` - Scan container image, directory, or SBOM
- `grype <target> -o <format>` - Specify output format
- `grype <target> --fail-on <severity>` - Exit with error on specified severity
- `grype <target> --only-fixed` - Show only vulnerabilities with confirmed fixes
- `grype <target> --only-notfixed` - Show vulnerabilities without confirmed fixes

### Database Commands
- `grype db status` - Show database status
- `grype db check` - Check for database updates
- `grype db update` - Update vulnerability database
- `grype db list` - List available databases
- `grype db import <archive>` - Import database archive

### Advanced Options
- `grype <target> --scope all-layers` - Scan all container image layers
- `grype <target> --exclude <path>` - Exclude specific file paths
- `grype <target> --add-cpes-if-none` - Generate CPE information if missing
- `grype <target> --distro <distro:version>` - Specify distribution for SBOM scanning
- `grype --version` - Show version information

## Supported Input Sources

### Container Images
```bash
# Docker images
grype ubuntu:latest
grype alpine:3.15

# Registry images
grype registry:docker.io/library/ubuntu:latest

# Local tar archives
grype docker-archive:ubuntu-latest.tar
grype oci-archive:ubuntu-latest.tar
```

### Filesystems and Directories
```bash
# Directory scanning
grype dir:/path/to/directory

# Current directory
grype dir:.

# Specific application directories
grype dir:/usr/local/bin
```

### SBOMs (Software Bill of Materials)
```bash
# SPDX SBOM
grype sbom:./sbom.spdx.json

# CycloneDX SBOM
grype sbom:./sbom.cyclonedx.json

# Syft-generated SBOM
grype sbom:./syft-output.json
```

## Output Formats

### Supported Formats
- **table** - Human-readable table format (default)
- **json** - Machine-readable JSON output
- **cyclonedx** - CycloneDX SBOM format
- **sarif** - Static Analysis Results Interchange Format
- **template** - Custom Go template format

### Format Examples
```bash
# JSON output for automation
grype ubuntu:latest -o json

# SARIF for IDE integration
grype . -o sarif

# CycloneDX for supply chain
grype myapp:latest -o cyclonedx
```

## Severity Levels

### Vulnerability Severities
- **negligible** - Minimal risk
- **low** - Low risk
- **medium** - Moderate risk
- **high** - Significant risk
- **critical** - Severe risk

### Severity Filtering
```bash
# Fail on medium or higher
grype ubuntu:latest --fail-on medium

# Show only critical vulnerabilities
grype myapp:latest --fail-on critical
```

## Use Cases

### CI/CD Integration
- **Security Gates**: Block deployments with critical vulnerabilities
- **Vulnerability Reporting**: Generate security reports in pipelines
- **Compliance Verification**: Ensure images meet security standards
- **Automated Scanning**: Continuous vulnerability assessment

### Development Workflow
- **Local Development**: Scan images before pushing to registry
- **Dependency Analysis**: Identify vulnerable dependencies
- **Base Image Selection**: Choose secure base images
- **Security Testing**: Validate application security posture

### Production Operations
- **Runtime Scanning**: Monitor deployed workloads for vulnerabilities
- **Incident Response**: Quickly assess vulnerability impact
- **Compliance Reporting**: Generate audit reports
- **Risk Assessment**: Prioritize vulnerability remediation

### Supply Chain Security
- **SBOM Analysis**: Scan software bills of materials
- **Third-party Components**: Assess vendor software security
- **Artifact Verification**: Validate software supply chain integrity
- **Vulnerability Tracking**: Monitor dependency vulnerabilities

## Configuration Examples

### Basic Scanning
```bash
# Scan container image
grype alpine:latest

# Scan local directory
grype dir:/path/to/project

# Scan with JSON output
grype myapp:v1.2.3 -o json > vulnerabilities.json
```

### Advanced Configuration
```bash
# Comprehensive image scan
grype ubuntu:latest \
  --scope all-layers \
  --fail-on medium \
  --only-fixed \
  -o json

# Directory scan with exclusions
grype dir:. \
  --exclude './node_modules/**' \
  --exclude './vendor/**' \
  --fail-on high
```

### SBOM Analysis
```bash
# Scan SBOM with distribution context
grype sbom:./alpine-sbom.json \
  --distro alpine:3.15 \
  --add-cpes-if-none

# Generate vulnerability report from SBOM
grype sbom:./app-dependencies.json \
  -o sarif > security-report.sarif
```

## Integration Patterns

### GitHub Actions
```yaml
- name: Vulnerability Scan
  run: |
    grype . -o sarif > grype-results.sarif
    grype . --fail-on critical
```

### Jenkins Pipeline
```groovy
stage('Security Scan') {
    steps {
        sh 'grype $IMAGE_NAME --fail-on high -o json > vulnerabilities.json'
        archiveArtifacts artifacts: 'vulnerabilities.json'
    }
}
```

### Docker Integration
```dockerfile
# Multi-stage build with security scanning
FROM alpine:latest AS scanner
RUN apk add --no-cache grype
COPY --from=builder /app /scan-target
RUN grype dir:/scan-target --fail-on critical
```

## Database Management

### Database Operations
```bash
# Check database status
grype db status

# Update to latest database
grype db update

# List available databases
grype db list

# Import custom database
grype db import ./custom-db.tar.gz
```

### Database Configuration
- **Automatic Updates**: Database updates automatically by default
- **Offline Mode**: Support for air-gapped environments
- **Custom Databases**: Import organization-specific vulnerability data
- **Version Management**: Track database versions and updates

## Best Practices

### Scanning Strategy
- **Layer-by-layer Analysis**: Use `--scope all-layers` for comprehensive scans
- **Severity Thresholds**: Set appropriate `--fail-on` levels for different environments
- **Fixed vs. Unfixed**: Use `--only-fixed` to focus on actionable vulnerabilities
- **Regular Updates**: Keep vulnerability database current

### Performance Optimization
- **Exclude Irrelevant Paths**: Use `--exclude` to skip test files and dependencies
- **SBOM-based Scanning**: Use pre-generated SBOMs for faster scans
- **Caching**: Leverage image layer caching for repeated scans
- **Parallel Processing**: Run scans in parallel for multiple targets

### Security Workflow
- **Pre-commit Scanning**: Scan during development
- **CI/CD Gates**: Implement security gates in pipelines
- **Production Monitoring**: Regular scanning of deployed images
- **Vulnerability Tracking**: Monitor and remediate identified issues

Grype provides comprehensive vulnerability scanning capabilities for modern software development and deployment workflows, supporting various input formats and integration patterns.