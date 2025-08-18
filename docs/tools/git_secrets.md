# Git Secrets

Prevents committing secrets and credentials to git repositories.

## Description

Git-secrets is a tool developed by AWS Labs that prevents you from committing passwords and other sensitive information to a git repository. It scans commits, commit messages, and merge operations for secrets using configurable regular expression patterns.

## MCP Tools

### Repository Scanning
- **`git_secrets_scan_repository`** - Scan git repository for secrets
- **`git_secrets_scan_history`** - Scan entire git repository history for secrets

### Configuration Management
- **`git_secrets_install_hooks`** - Install git-secrets hooks in repository
- **`git_secrets_add_pattern`** - Add secret detection pattern
- **`git_secrets_add_allowed`** - Add allowed pattern to prevent false positives
- **`git_secrets_list_config`** - List current git-secrets configuration

### AWS Integration
- **`git_secrets_register_aws`** - Register AWS-specific secret patterns

## Real CLI Commands Used

- `git secrets --scan` - Scan files for secrets
- `git secrets --scan-history` - Scan repository history
- `git secrets --install` - Install git hooks
- `git secrets --add '<pattern>'` - Add prohibited pattern
- `git secrets --add -a '<pattern>'` - Add allowed pattern
- `git secrets --register-aws` - Register AWS patterns
- `git secrets --list` - List configuration

## Security Patterns

### AWS Credentials
- AWS Access Key IDs
- AWS Secret Access Keys
- AWS Session Tokens
- RDS passwords
- EC2 private keys

### Generic Secrets
- SSH private keys
- API keys and tokens
- Database passwords
- Certificate private keys
- OAuth secrets

## Use Cases

### Pre-commit Protection
- Prevent accidental credential commits
- Block sensitive data before it reaches remote repositories
- Maintain clean git history
- Protect against credential leaks

### Repository Auditing
- Scan existing repositories for leaked credentials
- Historical analysis of potential security issues
- Compliance verification
- Security assessment of codebases

### Team Security
- Enforce security policies across development teams
- Prevent credential sharing through git
- Standardize secret detection across projects
- Reduce security incidents

## Installation and Setup

### Global Installation
```bash
# Install via Homebrew (macOS)
brew install git-secrets

# Manual installation
git clone https://github.com/awslabs/git-secrets.git
cd git-secrets
make install
```

### Repository Setup
```bash
# Install hooks in current repository
git secrets --install

# Register AWS patterns
git secrets --register-aws

# Add custom patterns
git secrets --add 'password\s*=\s*.+'
git secrets --add 'api[_-]?key\s*=\s*.+'
```

## Integration Examples

### CI/CD Pipelines
```yaml
# GitHub Actions example
- name: Scan for secrets
  run: |
    git secrets --install
    git secrets --register-aws
    git secrets --scan-history
```

### Pre-commit Hooks
```yaml
# .pre-commit-config.yaml
- repo: https://github.com/awslabs/git-secrets
  rev: master
  hooks:
    - id: git-secrets
      entry: git-secrets --scan
      types: [text]
```

Works with any git repository and integrates with development workflows to prevent credential leaks.