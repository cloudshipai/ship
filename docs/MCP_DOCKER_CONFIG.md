# Ship MCP Server - Docker Configuration Guide

This guide shows you how to configure Ship's MCP servers using Docker with Claude Desktop and other AI assistants.

## Prerequisites

- Docker installed and running
- Claude Desktop or other MCP-compatible AI assistant
- Ship Docker image: `ghcr.io/cloudshipai/ship:latest`

## Quick Start

### Build or Pull Ship Docker Image

```bash
# Build from source
docker build -t ghcr.io/cloudshipai/ship:latest .

# Or pull from Docker Hub (when available)
docker pull ghcr.io/cloudshipai/ship:latest

# Verify it works
docker run --rm ghcr.io/cloudshipai/ship:latest version
```

## Claude Desktop Configuration

### Configuration File Location

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

### Docker MCP Configuration Format

Edit your `claude_desktop_config.json` file and add Ship MCP servers:

```json
{
  "mcpServers": {
    "ship-security": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "--group-add=999",
        "-v",
        "/var/run/docker.sock:/var/run/docker.sock",
        "-v",
        "${workspaceFolder}:/workspace",
        "ghcr.io/cloudshipai/ship:latest",
        "mcp",
        "security"
      ]
    },
    "ship-terraform": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "--group-add=999",
        "-v",
        "/var/run/docker.sock:/var/run/docker.sock",
        "-v",
        "${workspaceFolder}:/workspace",
        "ghcr.io/cloudshipai/ship:latest",
        "mcp",
        "terraform"
      ]
    },
    "ship-all": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "--group-add=999",
        "-v",
        "/var/run/docker.sock:/var/run/docker.sock",
        "-v",
        "${workspaceFolder}:/workspace",
        "ghcr.io/cloudshipai/ship:latest",
        "mcp",
        "all"
      ]
    }
  }
}
```

## Available Tool Categories

### Individual Tools

Configure individual tools for lightweight, focused analysis:

#### Security Tools
```json
{
  "mcpServers": {
    "ship-semgrep": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "semgrep"]
    },
    "ship-trivy": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "trivy"]
    },
    "ship-gitleaks": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "gitleaks"]
    },
    "ship-checkov": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "checkov"]
    }
  }
}
```

#### Terraform Tools
```json
{
  "mcpServers": {
    "ship-tflint": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "tflint"]
    },
    "ship-terraform-docs": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "terraform-docs"]
    },
    "ship-tfsec": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "tfsec"]
    }
  }
}
```

#### Kubernetes Tools
```json
{
  "mcpServers": {
    "ship-kubescape": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "kubescape"]
    },
    "ship-kube-bench": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "kube-bench"]
    }
  }
}
```

### Tool Collections

Configure tool collections for comprehensive analysis:

```json
{
  "mcpServers": {
    "ship-security": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "security"],
      "description": "All security tools (31 tools)"
    },
    "ship-terraform": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "terraform"],
      "description": "All Terraform tools (7 tools)"
    },
    "ship-kubernetes": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "kubernetes"],
      "description": "All Kubernetes tools (9 tools)"
    },
    "ship-cloud": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "cloud"],
      "description": "All cloud tools"
    },
    "ship-finops": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "finops"],
      "description": "All FinOps tools"
    },
    "ship-supply-chain": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "supply-chain"],
      "description": "All supply chain security tools"
    },
    "ship-all": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--group-add=999", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "${workspaceFolder}:/workspace", "ghcr.io/cloudshipai/ship:latest", "mcp", "all"],
      "description": "All 56+ tools"
    }
  }
}
```

## Configuration with Environment Variables

For tools requiring credentials (AWS, cloud services, etc.):

```json
{
  "mcpServers": {
    "ship-aws-tools": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "--group-add=999",
        "-v",
        "/var/run/docker.sock:/var/run/docker.sock",
        "-v",
        "${workspaceFolder}:/workspace",
        "ghcr.io/cloudshipai/ship:latest",
        "mcp",
        "prowler",
        "--var",
        "AWS_PROFILE=default",
        "--var",
        "AWS_REGION=us-east-1"
      ],
      "env": {
        "AWS_PROFILE": "default",
        "AWS_REGION": "us-east-1"
      }
    }
  }
}
```

## Complete Tool List

### Terraform Tools (7)
- `tflint` - Terraform linter
- `terraform-docs` - Terraform documentation generator
- `inframap` - Infrastructure visualization
- `iac-plan` - Infrastructure as code planning
- `terraformer` - Infrastructure import and management
- `tfstate-reader` - Terraform state analysis
- `openinfraquote` - Infrastructure cost estimation

### Security Tools (31)
- `trivy` - Comprehensive vulnerability scanner
- `syft` - SBOM generation tool
- `checkov` - Infrastructure as code static analysis
- `terrascan` - IaC security scanner
- `tfsec` - Terraform-specific security scanner
- `semgrep` - Static analysis for security
- `actionlint` - GitHub Actions workflow linter
- `conftest` - OPA policy testing
- `kube-bench` - Kubernetes CIS benchmark
- `kube-hunter` - Kubernetes penetration testing
- `falco` - Runtime security monitoring
- `nuclei` - Fast vulnerability scanner
- `zap` - OWASP ZAP web application scanner
- `nmap` - Network exploration and security auditing
- `git-secrets` - Git repository secret scanner
- `trufflehog` - Advanced secret scanning
- `gitleaks` - Fast secret scanning
- `kubescape` - Kubernetes security scanner
- `dockle` - Container image linter
- `sops` - Secrets management
- `ossf-scorecard` - OSSF security scorecard
- `steampipe` - Cloud asset querying
- `cfn-nag` - CloudFormation security linter
- `gatekeeper` - OPA Gatekeeper policy validation
- `license-detector` - Software license detection
- `openscap` - Security compliance scanning
- `scout-suite` - Multi-cloud security auditing
- `powerpipe` - Infrastructure benchmarking
- `infrascan` - Infrastructure security scanning
- `github-admin` - GitHub administration tools
- `github-packages` - GitHub Packages security

### Kubernetes Tools (9)
- `velero` - Kubernetes backup and restore
- `goldilocks` - Kubernetes resource recommendations
- `fleet` - GitOps for Kubernetes
- `kuttl` - Kubernetes testing framework
- `litmus` - Chaos engineering for Kubernetes
- `cert-manager` - Certificate management
- `k8s-network-policy` - Kubernetes network policy management
- `kyverno` - Kubernetes policy management
- `kyverno-multitenant` - Multi-tenant Kyverno policies

### Cloud & Infrastructure Tools (3)
- `cloudquery` - Cloud asset inventory
- `custodian` - Cloud governance engine
- `packer` - Machine image building

### AWS IAM Tools (6)
- `cloudsplaining` - AWS IAM policy scanner
- `parliament` - AWS IAM policy linter
- `pmapper` - AWS IAM privilege escalation analysis
- `policy-sentry` - AWS IAM policy generator
- `prowler` - Multi-cloud security assessment
- `aws-iam-rotation` - AWS IAM credential rotation

### Supply Chain Tools (2)
- `cosign` - Container signing and verification
- `dependency-track` - OWASP Dependency-Track SBOM analysis

## Troubleshooting

### Docker Socket Permission Denied

If you get permission errors, adjust the `--group-add` value to match your Docker group GID:

```bash
# Find your Docker group GID
stat -c '%g' /var/run/docker.sock

# Use the result in your config
"args": ["run", "-i", "--rm", "--group-add=YOUR_GID", ...]
```

### Workspace Folder Not Found

Ensure the workspace mount path is correct:
- `${workspaceFolder}` works in VS Code
- For Claude Desktop, use absolute paths: `"/path/to/your/project:/workspace"`

### Container Not Starting

Check Docker is running:
```bash
docker ps
docker run --rm ghcr.io/cloudshipai/ship:latest version
```

## After Configuration

1. Save the configuration file
2. Restart Claude Desktop (or your AI assistant)
3. The Ship MCP tools will be available in your assistant

## Testing Your Configuration

```bash
# Test individual tool
docker run --rm -i --group-add=999 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd):/workspace \
  ghcr.io/cloudshipai/ship:latest mcp semgrep

# Test with workspace
cd your-project
docker run --rm -i --group-add=999 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd):/workspace \
  ghcr.io/cloudshipai/ship:latest mcp tflint
```

## Resources

- [Ship GitHub Repository](https://github.com/cloudshipai/ship)
- [Model Context Protocol Documentation](https://modelcontextprotocol.io)
- [Claude Desktop Documentation](https://claude.ai/desktop)
- [Docker Documentation](https://docs.docker.com)
