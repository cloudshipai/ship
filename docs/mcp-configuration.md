# Ship CLI MCP Server Configuration

The Ship CLI includes a built-in MCP (Model Context Protocol) server that exposes all CLI functionality as tools that AI assistants like Claude Code can use directly.

## Quick Setup

1. **Install Ship CLI** (if not already installed):
   ```bash
   go install github.com/cloudshipai/ship/cmd/ship@latest
   ```

2. **Create MCP Configuration** in your project or home directory:

### Example `.mcp.json` Configuration

```json
{
  "mcpServers": {
    "ship-cli": {
      "command": "ship",
      "args": ["mcp", "--stdio"],
      "env": {
        "AWS_REGION": "us-east-1",
        "OPENAI_API_KEY": "your-openai-api-key-here",
        "AWS_ACCESS_KEY_ID": "your-aws-access-key",
        "AWS_SECRET_ACCESS_KEY": "your-aws-secret-key"
      }
    }
  }
}
```

### Full Configuration Example

For a complete setup with multiple MCP servers:

```json
{
  "mcpServers": {
    "ship-cli": {
      "command": "ship",
      "args": ["mcp", "--stdio"],
      "env": {
        "AWS_REGION": "us-east-1",
        "OPENAI_API_KEY": "sk-your-openai-key-here",
        "AWS_ACCESS_KEY_ID": "AKIA...",
        "AWS_SECRET_ACCESS_KEY": "your-secret-key"
      }
    },
    "notion": {
      "command": "npx",
      "args": ["-y", "mcp-remote", "https://mcp.notion.com/sse"]
    },
    "linear": {
      "command": "npx", 
      "args": ["-y", "mcp-remote", "https://mcp.linear.app/sse"]
    }
  }
}
```

## Environment Variables

The Ship CLI MCP server uses these environment variables:

### Required for AI Investigation
- **`OPENAI_API_KEY`** - Your OpenAI API key for AI-powered infrastructure analysis
- **`AWS_ACCESS_KEY_ID`** - AWS access key for cloud resource access
- **`AWS_SECRET_ACCESS_KEY`** - AWS secret key for authentication
- **`AWS_REGION`** - AWS region to focus on (defaults to us-east-1)

### Optional
- **`AWS_SESSION_TOKEN`** - For temporary credentials/assumed roles
- **`AWS_PROFILE`** - AWS profile name (if using profiles instead of keys)

## Available Tools

Once configured, Claude Code will have access to these Ship CLI tools:

### 🔧 Terraform Analysis Tools
- **`terraform_lint`** - Run TFLint on Terraform code for syntax and best practices
- **`terraform_checkov_scan`** - Security scanning with Checkov for policy compliance
- **`terraform_security_scan`** - Alternative security scanning using Trivy
- **`terraform_cost_analysis`** - Infrastructure cost analysis using OpenInfraQuote
- **`terraform_generate_docs`** - Generate documentation using terraform-docs
- **`terraform_generate_diagram`** - Generate infrastructure diagrams from Terraform state

### 🤖 AI-Powered Investigation
- **`ai_investigate`** - Natural language infrastructure queries
  - Example: "Show me all S3 buckets with public access"
  - Example: "Find EC2 instances without proper backup tags"
  - Example: "Check for security groups allowing 0.0.0.0/0 access"

### ☁️ Cloud Operations
- **`cloudship_push`** - Upload artifacts to Cloudship for analysis

### 📚 Resources & Help
- **`ship://help`** - Complete CLI documentation
- **`ship://tools`** - List of available tools and examples

## Tool Parameters

All terraform-tools support these common parameters:
- **`directory`** - Directory containing Terraform files (default: current directory)
- **`push`** - Currently disabled during staging phase

### Tool-Specific Parameters
- **`terraform_lint`**: `config` (path to TFLint config file)
- **`terraform_checkov_scan`**: `framework` (terraform, cloudformation, etc.)
- **`terraform_cost_analysis`**: `provider` (aws, azure, gcp)
- **`terraform_generate_docs`**: `format` (markdown, json, table)
- **`terraform_generate_diagram`**: `output_path` (path for generated diagram)
- **`ai_investigate`**: `prompt` (required), `provider` (aws/azure/gcp), `aws_profile`, `aws_region`

## Usage Examples

### In Claude Code Chat

Once configured, you can use Ship CLI tools directly in conversation:

```
Can you analyze my Terraform code for security issues?
```

Claude Code will automatically use the `terraform_security_scan` tool.

```
Show me all my AWS S3 buckets and check which ones have public access
```

Claude Code will use the `ai_investigate` tool with Steampipe queries.

### Direct Tool Usage

You can also explicitly request specific tools:

```
Please use the terraform_lint tool to check my Terraform code in the current directory
```

```
Use ai_investigate to find all EC2 instances that have been stopped for more than 30 days
```

## Configuration Locations

Place your `.mcp.json` file in one of these locations:

1. **Project root** - For project-specific configuration
2. **Home directory** (`~/.mcp.json`) - For global configuration
3. **Claude Code config directory** - Following Claude Code's configuration

## Troubleshooting

### MCP Server Not Starting
- Ensure `ship` binary is in your PATH
- Check that all required environment variables are set
- Verify Docker is running (required for containerized tools)

### AWS Credentials Issues
- Set AWS credentials as environment variables in the MCP config
- Test credentials: `aws sts get-caller-identity`
- Ensure AWS region is specified

### AI Investigation Timeouts
- Increase timeout in prompts: "Use ai_investigate with a 5-minute timeout"
- Check OpenAI API key is valid and has sufficient credits
- Verify network connectivity to OpenAI and AWS

### Permission Errors
- Ensure AWS credentials have necessary permissions:
  - EC2: `DescribeInstances`, `DescribeSecurityGroups`
  - S3: `ListAllMyBuckets`, `GetBucketLocation`
  - IAM: `ListUsers`, `GetUser` (for security checks)

## Security Best Practices

1. **Use IAM Roles** when possible instead of access keys
2. **Rotate credentials** regularly
3. **Limit permissions** to only what's needed for your analysis
4. **Never commit** `.mcp.json` with real credentials to version control
5. **Use environment variables** or credential files instead of hardcoding

## Advanced Configuration

### Custom AWS Profile
```json
{
  "mcpServers": {
    "ship-cli": {
      "command": "ship",
      "args": ["mcp", "--stdio"],
      "env": {
        "AWS_PROFILE": "production",
        "AWS_REGION": "us-west-2",
        "OPENAI_API_KEY": "your-key-here"
      }
    }
  }
}
```

### Multiple Environments
```json
{
  "mcpServers": {
    "ship-dev": {
      "command": "ship",
      "args": ["mcp", "--stdio"],
      "env": {
        "AWS_PROFILE": "dev",
        "AWS_REGION": "us-east-1",
        "OPENAI_API_KEY": "your-key-here"
      }
    },
    "ship-prod": {
      "command": "ship", 
      "args": ["mcp", "--stdio"],
      "env": {
        "AWS_PROFILE": "production",
        "AWS_REGION": "us-west-2",
        "OPENAI_API_KEY": "your-key-here"
      }
    }
  }
}
```

## Getting Started

1. Copy the example `.mcp.json` configuration above
2. Replace placeholder values with your actual credentials
3. Place the file in your project root or home directory
4. Restart Claude Code
5. Start using Ship CLI tools in your conversations!

For more information, see the [Ship CLI documentation](../README.md) and [MCP Protocol specification](https://modelcontextprotocol.io/).