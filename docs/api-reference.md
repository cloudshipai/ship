# Ship CLI API Reference

## Cloudship API Endpoints

### Base URL
```
https://api.cloudship.ai/v1
```

### Authentication

All API requests require authentication using a bearer token:

```bash
Authorization: Bearer sk_prod_xxxxxxxxxxxx
```

### Endpoints

#### 1. Validate Token
```http
POST /auth/validate
Authorization: Bearer {token}

Response 200:
{
  "org_id": "org_123abc",
  "user_id": "user_456def",
  "permissions": ["push", "investigate", "mcp"],
  "valid": true
}

Response 401:
{
  "error": "invalid_token",
  "message": "Token is invalid or expired"
}
```

#### 2. Push Artifact
```http
POST /artifacts
Authorization: Bearer {token}
Content-Type: multipart/form-data

Form Data:
- file: (binary)
- kind: "tfplan" | "sbom" | "csv" | "steampipe"
- env: "prod" | "staging" | "dev"
- tags: {"key": "value", ...}
- sha256: "abc123..."

Response 201:
{
  "id": "art_789ghi",
  "sha256": "abc123...",
  "size": 1048576,
  "kind": "tfplan",
  "env": "prod",
  "url": "https://app.cloudship.ai/artifacts/art_789ghi",
  "created_at": "2025-01-24T12:00:00Z"
}

Response 413:
{
  "error": "file_too_large",
  "message": "File size exceeds 100MB limit",
  "max_size": 104857600
}
```

#### 3. Get Goals
```http
GET /goals?env={env}&panels=enabled
Authorization: Bearer {token}

Response 200:
{
  "goals": [
    {
      "id": "goal_cost_optimization",
      "name": "Cost Optimization",
      "description": "Identify cost saving opportunities",
      "panels": ["unused_resources", "rightsizing"],
      "provider": "aws",
      "priority": 1
    },
    {
      "id": "goal_security_hardening",
      "name": "Security Hardening",
      "description": "Improve security posture",
      "panels": ["public_exposure", "encryption_status"],
      "provider": "aws",
      "priority": 2
    }
  ],
  "environment": "prod",
  "total": 15
}
```

#### 4. Get Artifact Status
```http
GET /artifacts/{sha256}
Authorization: Bearer {token}

Response 200:
{
  "id": "art_789ghi",
  "sha256": "abc123...",
  "status": "processed",
  "insights": {
    "cost_impact": "+$1,234/month",
    "security_findings": 3,
    "performance_score": 85
  },
  "created_at": "2025-01-24T12:00:00Z",
  "processed_at": "2025-01-24T12:05:00Z"
}
```

### Error Responses

All errors follow this format:

```json
{
  "error": "error_code",
  "message": "Human readable message",
  "details": {
    "field": "additional context"
  }
}
```

Common error codes:
- `invalid_token` - Authentication failed
- `rate_limited` - Too many requests
- `file_too_large` - Upload exceeds size limit
- `invalid_kind` - Unsupported artifact type
- `org_suspended` - Organization access suspended
- `server_error` - Internal server error

### Rate Limits

- Authentication: 10 requests/minute
- Push: 100 requests/hour
- Goals: 60 requests/minute
- General: 1000 requests/hour

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1706102400
```

## MCP Protocol

### Tool Schemas

#### goal.list
```json
{
  "name": "goal.list",
  "description": "List available goals for an environment",
  "parameters": {
    "type": "object",
    "properties": {
      "env": {
        "type": "string",
        "enum": ["prod", "staging", "dev"],
        "description": "Environment to query"
      }
    },
    "required": ["env"]
  }
}
```

#### goal.run
```json
{
  "name": "goal.run",
  "description": "Execute investigation for a specific goal",
  "parameters": {
    "type": "object",
    "properties": {
      "goal_id": {
        "type": "string",
        "description": "Goal ID to execute"
      },
      "env": {
        "type": "string",
        "enum": ["prod", "staging", "dev"],
        "description": "Environment to investigate"
      },
      "provider": {
        "type": "string",
        "enum": ["aws", "cloudflare", "heroku"],
        "description": "Cloud provider"
      }
    },
    "required": ["goal_id", "env", "provider"]
  }
}
```

#### steampipe.query
```json
{
  "name": "steampipe.query",
  "description": "Execute a raw Steampipe SQL query",
  "parameters": {
    "type": "object",
    "properties": {
      "sql": {
        "type": "string",
        "description": "SQL query to execute"
      },
      "provider": {
        "type": "string",
        "enum": ["aws", "cloudflare", "heroku"],
        "description": "Provider context for the query"
      }
    },
    "required": ["sql", "provider"]
  }
}
```

### MCP Communication

#### JSON-RPC Request
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "goal.list",
  "params": {
    "env": "prod"
  }
}
```

#### JSON-RPC Response
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "goals": [
      {
        "id": "goal_cost_optimization",
        "name": "Cost Optimization",
        "panels": ["unused_resources", "rightsizing"]
      }
    ]
  }
}
```

#### Error Response
```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "error": {
    "code": -32601,
    "message": "Method not found",
    "data": {
      "method": "invalid.method"
    }
  }
}
```