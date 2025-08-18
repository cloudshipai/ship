# GitHub Packages

GitHub package management and security guidance through CLI.

## Description

GitHub Packages is a package hosting service that lets you host your packages privately or publicly. While GitHub doesn't provide a dedicated CLI for package security operations, this tool provides guidance on using GitHub's built-in security features and integrating with existing tools like cosign for package verification.

## MCP Tools

### Package Management
- **`github_packages_list_packages`** - List GitHub Packages in organization using gh API
- **`github_packages_audit_dependencies`** - Guidance on package dependency auditing
- **`github_packages_check_signatures`** - Guidance on package signature verification
- **`github_packages_enforce_policies`** - Guidance on package policy enforcement
- **`github_packages_generate_sbom`** - Guidance on SBOM generation for packages
- **`github_packages_get_version`** - Get GitHub CLI version information

## Real CLI Commands Used

### GitHub CLI API Commands
- `gh api orgs/{org}/packages` - List organization packages
- `gh api orgs/{org}/packages/{type}/{name}` - Get package details
- `gh api orgs/{org}/packages/{type}/{name}/versions` - List package versions
- `gh --version` - Get GitHub CLI version

### Package Management Integration
- Package security through GitHub's security tab
- Dependabot alerts for vulnerability detection
- Organization settings for package policies
- Cosign for package signature verification
- External SBOM tools (syft, cyclonedx-cli)

## Package Types Supported

### Supported Package Types
- **npm** - Node.js packages
- **docker** - Container images
- **maven** - Java packages
- **nuget** - .NET packages
- **rubygems** - Ruby packages

## Security Features

### Vulnerability Management
- Dependabot security alerts
- Vulnerability detection in dependencies
- Security advisories integration
- Repository security tab insights

### Access Control
- Organization-level package policies
- Package visibility controls (public/private/internal)
- Fine-grained permissions
- Team and user access management

### Signing and Verification
- Integration with sigstore/cosign
- Package attestation support
- Signature verification workflows
- Supply chain security best practices

## Use Cases

### Package Distribution
- Host private packages for organizations
- Publish open source packages
- Version management and release workflows
- Multi-language package support

### Security Management
- Monitor package vulnerabilities
- Enforce access policies
- Verify package signatures
- Audit package dependencies

### DevOps Integration
- CI/CD pipeline integration
- Automated package publishing
- Security scanning in workflows
- Policy enforcement automation

## Configuration Examples

### Organization Package Policy
```yaml
# Configure in Organization Settings > Packages
visibility: private
deletion_protection: enabled
package_creation: members_only
vulnerability_alerts: enabled
```

### Package Signature Verification
```bash
# Using cosign for container packages
cosign verify ghcr.io/org/package:tag \
  --certificate-identity-regexp="https://github.com/org/repo" \
  --certificate-oidc-issuer="https://token.actions.githubusercontent.com"
```

### SBOM Generation
```bash
# Using syft for package SBOM
syft packages ghcr.io/org/package:tag -o spdx-json

# Using CycloneDX
cyclonedx-cli analyze ghcr.io/org/package:tag
```

## Integration Points

### GitHub Actions
```yaml
- name: Publish Package
  run: |
    gh auth login --with-token < ${{ secrets.GITHUB_TOKEN }}
    gh api /user/packages
```

### Security Scanning
```yaml
- name: Security Scan
  run: |
    # Dependency scanning via Dependabot
    # Container scanning via trivy/grype
    # License scanning via license detectors
```

### Policy Enforcement
```yaml
- name: Policy Check
  run: |
    # Check organization package policies
    # Verify package compliance
    # Enforce security requirements
```

## Best Practices

### Security Guidelines
- Enable Dependabot alerts for all packages
- Use package signing with cosign
- Implement least-privilege access policies
- Regular security audits and updates

### Management Practices
- Clear package naming conventions
- Version tagging and release strategies
- Documentation and README files
- Deprecation and cleanup policies

### Integration Patterns
- Automated CI/CD publishing
- Security scanning in pipelines
- Policy-as-code implementation
- Monitoring and alerting setup

GitHub Packages integrates with the broader GitHub ecosystem for comprehensive package management, security, and governance capabilities.