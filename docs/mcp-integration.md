# MCP (Model Context Protocol) Integration

## Overview

Ship CLI includes a built-in MCP server that exposes 7 powerful Terraform analysis tools as MCP tools for AI assistants. This allows AI assistants like Claude Code, Cursor, and other MCP-compatible clients to use Ship CLI directly for infrastructure analysis and Terraform workflows.

## How It Works

### 1. MCP Server Architecture

```
AI Assistant (Claude Code)
        ↓ (MCP Protocol)
Ship CLI MCP Server
        ↓ (Command Execution)
Ship CLI Terraform Tools
        ↓ (Containerized Execution)
Dagger + Docker Containers
```

### 2. Available MCP Tools

#### **Terraform Analysis Tools**
- `terraform_lint` - Run TFLint on Terraform code for syntax and best practices
- `terraform_checkov_scan` - Run Checkov security scan for policy compliance
- `terraform_security_scan` - Run Trivy security scan for vulnerabilities
- `terraform_cost_analysis` - Analyze infrastructure costs with OpenInfraQuote
- `terraform_generate_docs` - Generate documentation with terraform-docs
- `terraform_generate_diagram` - Generate infrastructure diagrams with InfraMap
- `cloudship_push` - Upload artifacts for AI analysis

#### **Resources** (Information the AI can access)
- `ship://help` - Complete Ship CLI help and usage
- `ship://tools` - List of all available tools and capabilities

#### **Prompts** (Pre-built workflows)
- `security_audit` - Comprehensive security audit workflow
- `cost_optimization` - Cost optimization analysis workflow

### 3. Protocol Flow

1. **AI Assistant Request**: Claude asks to "analyze this Terraform code for security"
2. **MCP Tool Call**: Assistant calls `terraform_checkov_scan` tool with appropriate parameters
3. **Ship Execution**: MCP server executes `ship terraform-tools checkov-scan`
4. **Containerized Analysis**: Ship runs Checkov in a Dagger container
5. **Result Parsing**: MCP server formats results for the AI assistant
6. **AI Analysis**: Assistant analyzes results and provides insights

## Getting Started

### Prerequisites

1. **Ship CLI installed** and configured
2. **MCP-compatible AI assistant** (Claude Code, Cursor, etc.)

### Step 1: Start the MCP Server

```bash
# Start Ship CLI MCP server on stdio (default)
ship mcp

# The server will output:
# Starting Ship CLI MCP server on stdio...
```

### Step 2: Configure Your AI Assistant

#### For Claude Code

Add to your Claude Code configuration (`~/.config/claude-code/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "ship-cli": {
      "command": "/path/to/ship",
      "args": ["mcp"]
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
        "command": "ship mcp"
      }
    }
  }
}
```

### Step 3: Restart Your AI Assistant

Restart Claude Code or Cursor to load the new MCP server configuration.

### Step 4: Verify Connection

In your AI assistant, you should now see Ship CLI tools available. Try asking:

> "What Ship CLI tools are available?"

The assistant should be able to access the `ship://tools` resource and list all available functionality.

## Usage Examples

### Terraform Security Analysis

**You:** "Analyze the Terraform code in ./examples/aws-web-app for security issues"

**Assistant Actions:**
1. Calls `terraform_checkov_scan` tool with `directory="./examples/aws-web-app"`
2. Calls `terraform_security_scan` tool for additional security analysis
3. Combines results and provides comprehensive security analysis

### Terraform Cost Analysis

**You:** "Estimate costs for this Terraform infrastructure"

**Assistant Actions:**
1. Calls `terraform_cost_analysis` tool for the specified directory
2. Analyzes cost breakdown by resource type
3. Provides cost optimization recommendations

### Infrastructure Documentation

**You:** "Generate documentation for this Terraform module"

**Assistant Actions:**
1. Calls `terraform_generate_docs` tool
2. Formats the output appropriately
3. Can suggest improvements or additional documentation

### Infrastructure Diagrams

**You:** "Create a visual diagram of this Terraform infrastructure"

**Assistant Actions:**
1. Calls `terraform_generate_diagram` tool with `hcl: true`
2. Generates PNG/SVG diagram of infrastructure
3. Explains the infrastructure relationships

### Complete Infrastructure Analysis

**You:** "Perform a complete analysis of this Terraform project"

**Assistant Actions:**
1. Calls `terraform_lint` for code quality
2. Calls `terraform_checkov_scan` for security compliance
3. Calls `terraform_cost_analysis` for cost estimation
4. Calls `terraform_generate_docs` for documentation
5. Calls `terraform_generate_diagram` for visualization
6. Provides comprehensive analysis and recommendations

## Advanced Configuration

### Environment Variables

Set these environment variables to customize the MCP server:

```bash
# AWS Configuration (for cost analysis)
export AWS_REGION=us-west-2

# Ship Configuration  
export SHIP_CONFIG_DIR=~/.ship

# Debug Mode
export SHIP_DEBUG=true
```

### Custom Tool Parameters

The MCP tools support all the same parameters as the CLI commands:

```bash
# CLI equivalent:
ship terraform-tools cost-analysis --region us-west-2 --format json

# MCP tool call equivalent:
terraform_cost_analysis({
  directory: "./infrastructure",
  region: "us-west-2",
  format: "json"
})
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

3. **Terraform analysis errors**:
   ```bash
   # Test Ship CLI functionality
   ship terraform-tools lint
   
   # Check if Terraform files are valid
   terraform validate
   ```

### Debug Mode

Enable debug logging:

```bash
SHIP_DEBUG=true ship mcp
```

This will show detailed information about MCP requests and responses.

### Common Issues

1. **"Module not found" errors**: Ensure Ship CLI is properly installed and in your PATH
2. **Terraform validation errors**: Ensure your Terraform code is syntactically valid
3. **Large output handling**: The MCP server automatically chunks large outputs for better handling

## Security Considerations

1. **Local Execution**: All tools run locally in containers - no cloud credentials needed for most tools
2. **Command Execution**: AI assistants can execute Ship CLI commands - all tools are read-only analysis
3. **Data Privacy**: All analysis happens locally unless using CloudShip push features
4. **Container Isolation**: Tools run in isolated Docker containers via Dagger

## Available Tools Reference

For complete tool documentation with parameters and examples, see [llms.txt](../llms.txt).

### Quick Reference

| Tool | Purpose | Key Parameters |
|------|---------|---------------|
| `terraform_lint` | Code quality analysis | `directory`, `format`, `output` |
| `terraform_checkov_scan` | Security compliance | `directory`, `format`, `output` |
| `terraform_security_scan` | Vulnerability scanning | `directory` |
| `terraform_cost_analysis` | Cost estimation | `directory`, `region`, `format` |
| `terraform_generate_docs` | Documentation | `directory`, `filename`, `output` |
| `terraform_generate_diagram` | Infrastructure diagrams | `input`, `format`, `hcl`, `provider` |
| `cloudship_push` | Artifact upload | `file`, `type` |

## Extending the MCP Server

The MCP server automatically exposes new Ship CLI functionality:

1. **New CLI commands** are automatically available as MCP tools
2. **Custom modules** added via the module system become MCP tools
3. **Additional prompts** can be added to provide new AI workflows

See the [Dynamic Module Discovery](./dynamic-module-discovery.md) documentation for details on extending Ship CLI functionality.