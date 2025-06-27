# MCP (Model Context Protocol) Integration

## Overview

Ship CLI includes a built-in MCP server that exposes all Ship functionality as tools, resources, and prompts for AI assistants. This allows AI assistants like Claude Desktop, Cursor, and other MCP-compatible clients to use Ship CLI directly.

## How It Works

### 1. MCP Server Architecture

```
AI Assistant (Claude Desktop)
        ↓ (MCP Protocol)
Ship CLI MCP Server
        ↓ (Command Execution)
Ship CLI Tools (terraform-tools, ai-investigate, etc.)
        ↓ (Containerized Execution)
Dagger + Docker Containers
```

### 2. Available MCP Components

#### **Tools** (Actions the AI can take)
- `terraform_lint` - Run TFLint on Terraform code
- `terraform_security_scan` - Run Checkov security analysis
- `terraform_cost_estimate` - Estimate infrastructure costs
- `terraform_generate_docs` - Generate module documentation
- `ai_investigate` - Natural language infrastructure investigation
- `cloudship_push` - Upload artifacts for AI analysis

#### **Resources** (Information the AI can access)
- `ship://help` - Complete Ship CLI help and usage
- `ship://tools` - List of all available tools and capabilities

#### **Prompts** (Pre-built workflows)
- `security_audit` - Comprehensive security audit workflow
- `cost_optimization` - Cost optimization analysis workflow

### 3. Protocol Flow

1. **AI Assistant Request**: Claude asks to "investigate my AWS security"
2. **MCP Tool Call**: Assistant calls `ai_investigate` tool with appropriate parameters
3. **Ship Execution**: MCP server executes `ship ai-investigate --prompt "security" --execute`
4. **Steampipe Analysis**: Ship runs Steampipe queries against live AWS infrastructure
5. **Result Parsing**: MCP server formats results for the AI assistant
6. **AI Analysis**: Assistant analyzes results and provides insights

## Getting Started

### Prerequisites

1. **Ship CLI installed** and configured
2. **AWS credentials** configured (for infrastructure investigation)
3. **MCP-compatible AI assistant** (Claude Desktop, Cursor, etc.)

### Step 1: Start the MCP Server

```bash
# Start Ship CLI MCP server on stdio (default)
ship mcp

# The server will output:
# Starting Ship CLI MCP server on stdio...
```

### Step 2: Configure Your AI Assistant

#### For Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "ship-cli": {
      "command": "ship",
      "args": ["mcp"],
      "env": {
        "AWS_PROFILE": "your-aws-profile"
      }
    }
  }
}
```

#### For Cursor

Add to your Cursor MCP settings:

```json
{
  "mcp": {
    "servers": {
      "ship-cli": {
        "command": "ship mcp",
        "env": {
          "AWS_PROFILE": "your-aws-profile"
        }
      }
    }
  }
}
```

### Step 3: Restart Your AI Assistant

Restart Claude Desktop or Cursor to load the new MCP server configuration.

### Step 4: Verify Connection

In your AI assistant, you should now see Ship CLI tools available. Try asking:

> "What Ship CLI tools are available?"

The assistant should be able to access the `ship://tools` resource and list all available functionality.

## Usage Examples

### Infrastructure Investigation

**You:** "Investigate my AWS account for security issues"

**Assistant Actions:**
1. Calls `ai_investigate` tool with `prompt="Check for security issues"` and `execute=true`
2. Ship CLI runs Steampipe queries against your AWS account
3. Returns findings about open security groups, unencrypted resources, etc.
4. Assistant analyzes and provides recommendations

### Terraform Analysis

**You:** "Analyze the Terraform code in my current directory for security issues and cost"

**Assistant Actions:**
1. Calls `terraform_security_scan` tool for security analysis
2. Calls `terraform_cost_estimate` tool for cost estimation  
3. Combines results and provides comprehensive analysis

### Cost Optimization

**You:** "Help me find ways to reduce my AWS costs"

**Assistant Actions:**
1. Uses the `cost_optimization` prompt template
2. Calls `ai_investigate` to find unused resources
3. Calls `terraform_cost_estimate` if Terraform code is present
4. Provides prioritized cost-saving recommendations

### Documentation Generation

**You:** "Generate documentation for this Terraform module"

**Assistant Actions:**
1. Calls `terraform_generate_docs` tool
2. Formats the output appropriately
3. Can suggest improvements or additional documentation

## Advanced Configuration

### Environment Variables

Set these environment variables to customize the MCP server:

```bash
# AWS Configuration
export AWS_PROFILE=production
export AWS_REGION=us-west-2

# Ship Configuration  
export SHIP_CONFIG_DIR=~/.ship

# Debug Mode
export SHIP_DEBUG=true
```

### Multiple AWS Profiles

To work with multiple AWS profiles, configure separate MCP servers:

```json
{
  "mcpServers": {
    "ship-production": {
      "command": "ship",
      "args": ["mcp"],
      "env": {
        "AWS_PROFILE": "production"
      }
    },
    "ship-staging": {
      "command": "ship", 
      "args": ["mcp"],
      "env": {
        "AWS_PROFILE": "staging"
      }
    }
  }
}
```

### Custom Tool Parameters

The MCP tools support all the same parameters as the CLI commands:

```bash
# CLI equivalent:
ship ai-investigate --prompt "List S3 buckets" --provider aws --aws-region us-west-2 --execute

# MCP tool call equivalent:
ai_investigate(
  prompt="List S3 buckets",
  provider="aws", 
  aws_region="us-west-2",
  execute=true
)
```

## Troubleshooting

### Connection Issues

1. **Server not starting**:
   ```bash
   # Test server manually
   ship mcp
   # Should show: "Starting Ship CLI MCP server on stdio..."
   ```

2. **Tools not visible in AI assistant**:
   - Check MCP configuration file syntax
   - Restart the AI assistant
   - Verify Ship CLI is in your PATH

3. **AWS authentication errors**:
   ```bash
   # Test AWS credentials
   aws sts get-caller-identity --profile your-profile
   
   # Test Ship CLI AWS access
   ship ai-investigate --prompt "test connection" --execute
   ```

### Debug Mode

Enable debug logging:

```bash
SHIP_DEBUG=true ship mcp
```

This will show detailed information about MCP requests and responses.

### Common Issues

1. **"Module not found" errors**: Ensure Ship CLI is properly installed and in your PATH
2. **AWS permission errors**: Verify your AWS credentials have necessary permissions for Steampipe queries
3. **Timeout issues**: Large infrastructure investigations may take time; the AI assistant should wait for completion

## Security Considerations

1. **Credential Access**: The MCP server runs with your local credentials - ensure your AI assistant is trusted
2. **Command Execution**: AI assistants can execute Ship CLI commands - review what tools you expose
3. **Network Access**: Infrastructure investigation tools make network calls to cloud APIs
4. **Data Exposure**: Investigation results may contain sensitive infrastructure information

## Extending the MCP Server

The MCP server automatically exposes new Ship CLI functionality:

1. **New CLI commands** are automatically available as MCP tools
2. **Custom modules** added via the module system become MCP tools
3. **Additional prompts** can be added to provide new AI workflows

See the [Dynamic Module Discovery](./dynamic-module-discovery.md) documentation for details on extending Ship CLI functionality.