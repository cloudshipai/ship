# Ship CLI API Reference

## Cloudship API Endpoints

### Base URL
```
https://staging.cloudshipai.com/api/v1
```

### Authentication

All API requests require authentication using an API key as a bearer token:

```bash
Authorization: Bearer your-api-key
```

Get your API key from: https://app.cloudshipai.com/settings/api-keys

### Endpoints

#### 1. Upload Artifact
```http
POST /artifacts/upload
Authorization: Bearer {api-key}
Content-Type: application/json

Request Body:
{
  "fleet_id": "your-fleet-id",
  "file_name": "analysis.json",
  "file_type": "application/json",
  "content": "base64-encoded-content",
  "metadata": {
    "scan_type": "security_scan",
    "scan_timestamp": "2025-01-30T10:00:00Z",
    "source": "ship-cli/v0.3.5",
    "tags": ["production", "aws"],
    "custom_field": "custom_value"
  }
}

Response 200:
{
  "artifact_id": "art_123abc",
  "download_url": "https://...",
  "version": 1,
  "created_at": "2025-01-30T10:00:00Z"
}
```

#### 2. List Artifacts
```http
GET /artifacts?fleet_id={fleet-id}&limit=50&offset=0&type=security_scan
Authorization: Bearer {api-key}

Response 200:
{
  "artifacts": [
    {
      "id": "art_123abc",
      "file_name": "analysis.json",
      "file_type": "application/json",
      "file_size": 1024,
      "version": 1,
      "created_at": "2025-01-30T10:00:00Z",
      "metadata": {...}
    }
  ],
  "total": 100,
  "limit": 50,
  "offset": 0
}
```

#### 3. Download Artifact
```http
GET /artifacts/{artifact-id}/download
Authorization: Bearer {api-key}

Response 200:
(Binary file content)
```

### Error Responses

#### 401 Unauthorized
```json
{
  "error": "Unauthorized", 
  "message": "Invalid API key"
}
```

#### 413 Payload Too Large
```json
{
  "error": "file_too_large",
  "message": "File size exceeds 100MB limit",
  "max_size": 104857600
}
```

### File Size Limits

- Maximum file size: 100MB (104,857,600 bytes)
- Content is base64 encoded, so actual file size limit is ~75MB before encoding

### Supported Artifact Types

The `scan_type` in metadata can be any of:
- `terraform_plan` - Terraform plan files
- `security_scan` - Security scan results (Trivy, Checkov, etc.)
- `cost_analysis` - Cost analysis results
- `cost_estimate` - Cost estimation results (Infracost)
- `lint_results` - Linting results (TFLint)
- `terraform_docs` - Generated documentation
- `infrastructure_diagram` - Infrastructure diagrams
- `checkov_scan` - Checkov security scan results
- `infracost_estimate` - Infracost cost estimates

## MCP Protocol

Ship CLI includes a built-in MCP (Model Context Protocol) server that exposes all CLI functionality to AI assistants.

### Starting the MCP Server

```bash
ship mcp
```

### Available Tools

The MCP server exposes the following tools to AI assistants:

1. **terraform_lint** - Lint Terraform code
2. **terraform_security_scan** - Security scan using multiple tools
3. **terraform_cost_estimate** - Cost estimation and analysis
4. **terraform_generate_docs** - Generate documentation
5. **terraform_generate_diagram** - Generate infrastructure diagrams
6. **ai_investigate** - Natural language infrastructure queries
7. **cloudship_push** - Push artifacts to CloudShip

### MCP Configuration

#### Claude Desktop
Add to `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "ship-cli": {
      "command": "ship",
      "args": ["mcp"],
      "env": {
        "AWS_PROFILE": "your-profile",
        "CLOUDSHIP_API_KEY": "your-api-key",
        "CLOUDSHIP_FLEET_ID": "your-fleet-id"
      }
    }
  }
}
```

#### Environment Variables

The MCP server respects the following environment variables:
- `AWS_PROFILE` - AWS profile for cloud access
- `CLOUDSHIP_API_KEY` - CloudShip API key for authentication
- `CLOUDSHIP_FLEET_ID` - Default fleet ID for artifact uploads
- `INFRACOST_API_KEY` - Infracost API key for detailed cost estimates