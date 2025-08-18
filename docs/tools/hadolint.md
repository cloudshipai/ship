# Hadolint

Dockerfile linter for best practices and security analysis.

## Description

Hadolint is a Dockerfile linter that helps you build best practice Docker images. It parses the Dockerfile into an AST and performs rules on top of the AST. It stands on the shoulders of ShellCheck to lint the Bash code inside RUN instructions. The tool supports various output formats and configuration options for CI/CD integration.

## MCP Tools

### Dockerfile Analysis
- **`hadolint_scan_dockerfile`** - Scan Dockerfile for best practices and security issues
- **`hadolint_scan_directory`** - Scan directory for Dockerfiles
- **`hadolint_scan_with_config`** - Scan Dockerfile with custom configuration
- **`hadolint_scan_ignore_rules`** - Scan Dockerfile while ignoring specific rules

### Utility
- **`hadolint_get_version`** - Get Hadolint version information

## Real CLI Commands Used

### Core Commands
- `hadolint <Dockerfile>` - Basic Dockerfile linting
- `hadolint --format <format> <Dockerfile>` - Specify output format
- `hadolint --config <config> <Dockerfile>` - Use custom configuration
- `hadolint --ignore <rule> <Dockerfile>` - Ignore specific rules
- `hadolint --version` - Show version information

### Directory Scanning
- `find <directory> -name 'Dockerfile*' -exec hadolint {} +` - Scan all Dockerfiles in directory
- `find <directory> -name 'Dockerfile*' -exec hadolint --format <format> {} +` - Directory scan with format

### Configuration Options
- `--trusted-registry <registry>` - Allow specific registries in FROM instructions
- `--require-label <label:format>` - Enforce specific label requirements
- `--strict-labels` - Prohibit labels not in defined schema
- `--no-fail` - Don't exit with failure status on rule violations
- `--no-color` - Disable colorized output

## Supported Output Formats

### Available Formats
- **tty** - Default terminal output with colors
- **json** - Structured JSON format for automation
- **checkstyle** - Checkstyle XML format for CI/CD
- **codeclimate** - Code Climate JSON format
- **gitlab_codeclimate** - GitLab Code Quality format
- **gnu** - GNU-style error format
- **codacy** - Codacy platform format
- **sonarqube** - SonarQube generic issue format
- **sarif** - Static Analysis Results Interchange Format

### Format Examples
```bash
# JSON for automation
hadolint --format json Dockerfile

# SARIF for security tools
hadolint --format sarif Dockerfile

# Code Climate for GitLab integration
hadolint --format gitlab_codeclimate Dockerfile
```

## Rule Categories

### Best Practices
- **DL3000-DL3999** - Docker best practices
- **DL4000-DL4999** - Maintainer guidelines
- **DL5000-DL5999** - Security recommendations

### Common Rules
- **DL3003** - Use WORKDIR to switch to a directory
- **DL3006** - Always tag the version of an image explicitly
- **DL3008** - Pin versions in apt get install
- **DL3009** - Delete the apt-get lists after installing something
- **DL3015** - Avoid additional packages by specifying `--no-install-recommends`

### Security Rules
- **DL3001** - For some bash commands it makes no sense running them in a Docker container
- **DL3002** - Last user should not be root
- **DL3004** - Do not use sudo as it leads to unpredictable behavior
- **DL3007** - Using latest is prone to errors if the image will ever update

## Use Cases

### CI/CD Integration
- **Quality Gates**: Block builds with Dockerfile issues
- **Automated Reviews**: Generate feedback on pull requests
- **Compliance Checking**: Ensure Docker best practices
- **Security Scanning**: Identify security anti-patterns

### Development Workflow
- **Local Development**: Lint Dockerfiles before committing
- **IDE Integration**: Real-time feedback during development
- **Team Standards**: Enforce consistent Docker practices
- **Learning Tool**: Understand Docker best practices

### Security Operations
- **Security Reviews**: Identify potential security issues
- **Compliance Audits**: Verify adherence to security policies
- **Base Image Validation**: Ensure secure base image usage
- **Secret Detection**: Prevent hardcoded secrets in Dockerfiles

### Container Optimization
- **Build Optimization**: Improve build times and image sizes
- **Layer Optimization**: Reduce unnecessary layers
- **Cache Efficiency**: Optimize Docker build cache usage
- **Multi-stage Builds**: Validate multi-stage build practices

## Configuration Examples

### Basic Scanning
```bash
# Scan single Dockerfile
hadolint Dockerfile

# Scan with JSON output
hadolint --format json Dockerfile

# Scan all Dockerfiles in project
find . -name 'Dockerfile*' -exec hadolint {} +
```

### Rule Management
```bash
# Ignore specific rules
hadolint --ignore DL3003 --ignore DL3006 Dockerfile

# Use configuration file
hadolint --config .hadolint.yaml Dockerfile

# Trusted registry example
hadolint --trusted-registry my-company.com:5000 Dockerfile
```

### Configuration File (.hadolint.yaml)
```yaml
ignored:
  - DL3003
  - DL3006
  - DL3009

trusted-registries:
  - my-company.com:5000
  - docker.io

failure-threshold: error

format: json

require-labels:
  - maintainer:email
  - version:semver

strict-labels: true
```

## Integration Patterns

### GitHub Actions
```yaml
- name: Lint Dockerfile
  run: |
    hadolint Dockerfile --format sarif > hadolint-results.sarif
    hadolint Dockerfile --format json > hadolint-results.json
```

### GitLab CI
```yaml
hadolint:
  image: hadolint/hadolint:latest
  script:
    - hadolint --format gitlab_codeclimate Dockerfile > hadolint-report.json
  artifacts:
    reports:
      codequality: hadolint-report.json
```

### Jenkins Pipeline
```groovy
stage('Dockerfile Lint') {
    steps {
        sh 'hadolint --format checkstyle Dockerfile > hadolint-results.xml'
        recordIssues enabledForFailure: true, tools: [checkStyle(pattern: 'hadolint-results.xml')]
    }
}
```

### Docker Integration
```bash
# Run hadolint in container
docker run --rm -i hadolint/hadolint < Dockerfile

# With configuration
docker run --rm -i -v "$PWD"/.hadolint.yaml:/.config/hadolint.yaml hadolint/hadolint < Dockerfile
```

## Advanced Usage

### Custom Rules Configuration
```yaml
# .hadolint.yaml
ignored:
  - DL3008  # Pin versions in apt get install
  - DL3009  # Delete apt-get lists

override:
  error:
    - DL3001  # Critical security issues
    - DL3002
  warning:
    - DL3003  # Use WORKDIR
    - DL3006  # Tag versions
  info:
    - DL3015  # Avoid additional packages

failure-threshold: warning
```

### Multi-stage Build Validation
```bash
# Lint multi-stage Dockerfile
hadolint --format json Dockerfile.multistage

# With specific stage focus
hadolint --require-label stage:name Dockerfile.multistage
```

### Security-focused Configuration
```yaml
# Security-first hadolint configuration
failure-threshold: warning

ignored: []  # Don't ignore any rules

require-labels:
  - maintainer:email
  - security-contact:email
  - version:semver

strict-labels: true

trusted-registries:
  - registry.access.redhat.com
  - gcr.io/distroless
```

## Best Practices

### Rule Management
- **Start Permissive**: Begin with warnings, gradually enforce errors
- **Team Consensus**: Agree on ignored rules as a team
- **Documentation**: Document why specific rules are ignored
- **Regular Review**: Periodically review ignored rules

### CI/CD Integration
- **Fail Fast**: Run hadolint early in build pipeline
- **Multiple Formats**: Generate both human and machine-readable outputs
- **Caching**: Cache hadolint results for unchanged Dockerfiles
- **Parallel Execution**: Lint multiple Dockerfiles concurrently

### Development Workflow
- **Pre-commit Hooks**: Run hadolint before committing
- **IDE Integration**: Use hadolint plugins for real-time feedback
- **Local Testing**: Test Dockerfile changes locally first
- **Incremental Improvement**: Fix issues gradually in existing projects

### Security Considerations
- **Security Rules Priority**: Treat security rules as errors
- **Base Image Validation**: Enforce trusted registry usage
- **Secret Prevention**: Never ignore rules about hardcoded secrets
- **User Permissions**: Always validate final user in containers

Hadolint provides comprehensive Dockerfile analysis to improve security, maintainability, and adherence to Docker best practices.