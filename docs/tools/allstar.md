# Allstar - GitHub App Security Policy Enforcement

**IMPORTANT**: Allstar is NOT a command-line tool. It is a GitHub App service that provides automated security policy enforcement for GitHub repositories.

## Overview

Allstar is an open source GitHub App installed on organizations or repositories to set and enforce security policies. Unlike other tools in Ship CLI, Allstar runs as a service and cannot be executed as a CLI command.

## What Allstar Does

- **Automated Security Policy Enforcement**: Continuously monitors repositories for security policy violations
- **Branch Protection**: Ensures proper branch protection rules are configured
- **Binary Artifacts**: Detects and prevents binary files in repositories
- **Outside Collaborators**: Monitors external collaborators on repositories
- **Security.md Files**: Ensures security policies are documented
- **Dangerous Workflow Detection**: Identifies risky GitHub Actions workflows
- **OSSF Scorecard Integration**: Provides security scoring for repositories

## Available MCP Function

Since Allstar is not a CLI tool, Ship CLI provides only an informational function:

### `allstar_info`
**Description**: Get information about Allstar GitHub App service

**Parameters**: None

**Returns**: Information about Allstar and installation instructions

**Example Usage**:
```bash
allstar_info()
```

## Installation and Setup

Allstar is installed as a GitHub App, not as a CLI tool:

### 1. Install the GitHub App
```bash
# Visit the GitHub App installation page
# https://github.com/apps/allstar
```

### 2. Configure Policies
Create `.allstar/` directory in your repository or organization:

```yaml
# .allstar/branch_protection.yaml
enforcement: "on"  # or "off", "log"
protectionLevel: "admin_push"  # or "non_admin_push", "admin_bypass"
requireBranchBeDeleted: true
```

```yaml
# .allstar/binary_artifacts.yaml  
enforcement: "on"
allowedBinaryFiles:
  - "*.jpg"
  - "*.png"
  - "docs/*.pdf"
```

```yaml
# .allstar/outside_collaborators.yaml
enforcement: "on"
allowedCollaborators:
  - "username1"
  - "username2"
```

### 3. Global Configuration
```yaml
# .allstar/allstar.yaml
enforcement: "on"
enableCheckboxes: true
notificationsEnabled: true
issueLabel: "security"
issuePriority: "high"
```

## Policy Types

### Branch Protection
- Enforces branch protection rules
- Prevents direct pushes to protected branches
- Requires pull request reviews
- Blocks force pushes

### Binary Artifacts
- Detects binary files in repositories
- Prevents accidental commit of executables
- Configurable allow/deny lists
- Supports file extension patterns

### Outside Collaborators
- Monitors external repository access
- Tracks collaborator permissions
- Alerts on unauthorized access
- Maintains collaborator whitelist

### Security.md
- Ensures security policy documentation
- Validates security contact information
- Checks for vulnerability disclosure process
- Monitors security.md file presence

### Dangerous Workflows
- Analyzes GitHub Actions workflows
- Detects potentially dangerous patterns
- Identifies security risks in CI/CD
- Prevents workflow-based attacks

## Configuration Examples

### Repository-Level Configuration
```yaml
# .allstar/allstar.yaml in specific repository
enforcement: "log"  # Only log violations, don't create issues
repositories:
  - "my-repo"
excludePrivateRepos: false
```

### Organization-Level Configuration
```yaml
# .allstar/allstar.yaml in .allstar repository
enforcement: "on"
repositories:
  - "*"  # Apply to all repositories
excludeRepositories:
  - "legacy-repo"
  - "archived-*"
```

## Integration with Ship CLI

Since Allstar is not a CLI tool, the Ship CLI MCP integration provides:

1. **Information Function**: `allstar_info` returns setup instructions
2. **Documentation**: This reference guide for understanding Allstar
3. **No Direct Execution**: Cannot run Allstar commands through Ship CLI

## Monitoring and Alerts

Allstar creates GitHub issues when policy violations are detected:

- **Issue Creation**: Automated issue filing for violations
- **Issue Updates**: Status updates as violations are resolved
- **Labels and Priorities**: Configurable issue categorization
- **Notifications**: GitHub notifications for security events

## Troubleshooting

### Common Issues

1. **App Not Responding**
   - Check GitHub App installation status
   - Verify repository permissions
   - Review App logs in GitHub settings

2. **Policies Not Enforcing**
   - Validate YAML syntax in `.allstar/` files
   - Check enforcement levels (on/off/log)
   - Verify App has necessary permissions

3. **Too Many Issues Created**
   - Adjust enforcement to "log" mode temporarily
   - Configure exclude patterns
   - Set appropriate issue labels

### Getting Help

- **Official Repository**: https://github.com/ossf/allstar
- **Documentation**: https://github.com/ossf/allstar/tree/main/docs
- **OSSF Community**: https://openssf.org/community/
- **GitHub Discussions**: Available in the Allstar repository

## Alternative CLI Tools

For command-line security scanning, consider these alternatives available in Ship CLI:

- **GitHub Security**: `git_secrets` for secret detection
- **Policy Validation**: `conftest` or `gatekeeper` for policy testing
- **Repository Scanning**: `trivy`, `semgrep`, or `gitleaks` for code analysis
- **Compliance Checking**: `scout_suite` or `prowler` for broader security assessment

## References

- **Official Repository**: https://github.com/ossf/allstar
- **GitHub App Page**: https://github.com/apps/allstar
- **OSSF Security Initiative**: https://openssf.org/
- **Policy Configuration Guide**: https://github.com/ossf/allstar/blob/main/docs/policies.md