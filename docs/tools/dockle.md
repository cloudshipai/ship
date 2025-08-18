# Dockle

Container image linter for security and best practices.

## Description

Dockle is a container image linter that helps you to build secure container images by analyzing images against security best practices. It checks for vulnerabilities, misconfigurations, and adherence to best practices in container images.

## MCP Tools

### Image Scanning
- **`dockle_scan_image`** - Scan container images for security and best practices
- **`dockle_scan_tarball`** - Scan container image tarball files

### JSON Output
- **`dockle_scan_json`** - Scan container images with JSON output format
- **`dockle_scan_tarball_json`** - Scan tarball files with JSON output format

## Real CLI Commands Used

- `dockle [IMAGE_NAME]` - Scan container image
- `dockle --input [TARBALL_PATH]` - Scan from image tarball
- `dockle -f json [IMAGE_NAME]` - Scan with JSON output
- `dockle -f json -o results.json [IMAGE_NAME]` - Scan with JSON output to file

## Security Checks

Dockle performs various security checks including:

### Dockerfile Best Practices
- Use COPY instead of ADD
- Avoid running as root user
- Use specific version tags, not 'latest'
- Remove unnecessary packages

### Security Configuration
- Check for exposed sensitive files
- Verify proper user permissions
- Scan for secret leaks in image layers
- Validate security configurations

### Image Composition
- Check for known vulnerabilities
- Verify image metadata
- Analyze layer composition
- Check for unnecessary files

## Severity Levels

- **FATAL**: Critical security issues that must be fixed
- **WARN**: Important issues that should be addressed
- **INFO**: Minor issues with potential performance impact
- **PASS**: No problems detected

## Use Cases

### CI/CD Integration
- Automated container security scanning
- Quality gates for container deployment
- Security compliance verification
- Build-time security checks

### Development Workflow
- Local container security testing
- Pre-commit security validation
- Security best practices enforcement
- Developer security education

### Production Security
- Container registry scanning
- Runtime security verification
- Compliance reporting
- Security audit trails

## Output Formats

### Standard Output
Human-readable format with color-coded severity levels and detailed explanations.

### JSON Output
Machine-readable format for integration with other tools and automated processing.

## Integration Examples

### GitHub Actions
```yaml
- name: Scan container image
  run: dockle myapp:latest
```

### GitLab CI
```yaml
security_scan:
  script:
    - dockle -f json -o dockle-results.json $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
```

### Docker Tarball Scanning
```bash
# Save image to tarball
docker save myapp:latest -o myapp.tar

# Scan the tarball
dockle --input myapp.tar
```

Works with any OCI-compliant container image and supports both local and remote image scanning.