# Ship CLI Technical Specification

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                     Ship CLI                            │
├─────────────────────────────────────────────────────────┤
│  Commands Layer                                         │
│  ┌─────────┐ ┌──────────┐ ┌──────────────┐ ┌────────┐ │
│  │  Auth   │ │   Push   │ │ Investigate  │ │  MCP   │ │
│  └─────────┘ └──────────┘ └──────────────┘ └────────┘ │
├─────────────────────────────────────────────────────────┤
│  Core Services                                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐ │
│  │  Config  │ │   HTTP   │ │  Dagger  │ │    MCP    │ │
│  │  Manager │ │  Client  │ │  Engine  │ │   Server  │ │
│  └──────────┘ └──────────┘ └──────────┘ └───────────┘ │
├─────────────────────────────────────────────────────────┤
│  External Dependencies                                  │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐               │
│  │  Docker  │ │Cloudship │ │Steampipe │               │
│  │  Runtime │ │    API   │ │   MCP    │               │
│  └──────────┘ └──────────┘ └──────────┘               │
└─────────────────────────────────────────────────────────┘
```

## Module Structure

```
ship/
├── cmd/
│   └── ship/
│       └── main.go              # Entry point
├── internal/
│   ├── auth/                    # Authentication logic
│   │   ├── auth.go
│   │   ├── token.go
│   │   └── auth_test.go
│   ├── cli/                     # Command implementations
│   │   ├── auth.go
│   │   ├── push.go
│   │   ├── investigate.go
│   │   ├── mcp.go
│   │   └── root.go
│   ├── config/                  # Configuration management
│   │   ├── config.go
│   │   ├── loader.go
│   │   └── config_test.go
│   ├── client/                  # HTTP client
│   │   ├── client.go
│   │   ├── auth.go
│   │   └── retry.go
│   ├── push/                    # Push functionality
│   │   ├── upload.go
│   │   ├── checksum.go
│   │   └── detector.go
│   ├── investigate/             # Investigation logic
│   │   ├── orchestrator.go
│   │   ├── goals.go
│   │   ├── providers.go
│   │   └── steampipe.go
│   ├── dagger/                  # Dagger integration
│   │   ├── engine.go
│   │   ├── container.go
│   │   └── embed.go
│   ├── mcp/                     # MCP server
│   │   ├── server.go
│   │   ├── tools.go
│   │   ├── rpc.go
│   │   └── schema.go
│   └── telemetry/              # Analytics
│       ├── metrics.go
│       └── consent.go
├── pkg/                         # Public packages
│   ├── types/                   # Shared types
│   │   ├── artifact.go
│   │   ├── goal.go
│   │   └── provider.go
│   └── utils/                   # Utilities
│       ├── progress.go
│       ├── terminal.go
│       └── fs.go
├── assets/                      # Embedded assets
│   ├── dagger/                  # Dagger binary
│   └── templates/               # SQL templates
├── scripts/                     # Build scripts
│   ├── build.sh
│   └── embed-dagger.sh
└── test/                        # E2E tests
    └── e2e/
```

## Core Components

### 1. Configuration System

```go
type Config struct {
    Token      string            `yaml:"token"`
    OrgID      string            `yaml:"org_id"`
    DefaultEnv string            `yaml:"default_env"`
    BaseURL    string            `yaml:"base_url"`
    Telemetry  TelemetryConfig   `yaml:"telemetry"`
}

type TelemetryConfig struct {
    Enabled   bool   `yaml:"enabled"`
    SessionID string `yaml:"session_id"`
}
```

**Location**: `~/.ship/config.yaml`

### 2. Authentication Flow

```go
type AuthClient interface {
    Authenticate(token string) (*AuthResponse, error)
    ValidateToken() error
    Logout() error
}

type AuthResponse struct {
    OrgID       string `json:"org_id"`
    UserID      string `json:"user_id"`
    Permissions []string `json:"permissions"`
}
```

### 3. Artifact Push System

```go
type Artifact struct {
    Path     string
    Kind     string
    Env      string
    Tags     map[string]string
    SHA256   string
    Size     int64
}

type PushClient interface {
    Push(ctx context.Context, artifact *Artifact) (*PushResponse, error)
    ValidateSize(size int64) error
    DetectKind(content []byte) (string, error)
}
```

### 4. Investigation Orchestrator

```go
type Investigator interface {
    FetchGoals(env string) ([]Goal, error)
    MapGoalsToQueries(goals []Goal, useAI bool) ([]Query, error)
    ExecuteQueries(queries []Query) ([]Result, error)
    BundleResults(results []Result) (*Bundle, error)
}

type Goal struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Panels      []string  `json:"panels"`
    Provider    string   `json:"provider"`
}
```

### 5. MCP Server Implementation

```go
type MCPServer struct {
    port      int
    tools     map[string]Tool
    engine    *DaggerEngine
}

type Tool interface {
    Name() string
    Schema() json.RawMessage
    Execute(params json.RawMessage) (interface{}, error)
}
```

## API Endpoints

### Authentication
- `POST /v1/auth/validate` - Validate token
- `POST /v1/auth/logout` - Revoke token

### Artifacts
- `POST /v1/artifacts` - Upload artifact
- `GET /v1/artifacts/{sha}` - Get artifact status

### Goals
- `GET /v1/goals?env={env}&panels=enabled` - List goals

## Security Considerations

1. **Token Storage**
   - Tokens stored in config file with 0600 permissions
   - Environment variable takes precedence
   - No tokens in command history

2. **TLS/HTTPS**
   - All API calls use HTTPS
   - Certificate validation enabled
   - Optional proxy support

3. **Artifact Integrity**
   - Client-side SHA-256 calculation
   - Server validates checksum
   - Duplicate detection

## Performance Targets

| Operation | Target | Max |
|-----------|--------|-----|
| CLI startup | <100ms | 200ms |
| Auth validation | <500ms | 1s |
| File upload (10MB) | <5s | 10s |
| First investigate | <150s | 180s |
| Subsequent investigate | <30s | 45s |
| MCP tool call | <100ms | 500ms |

## Error Handling

### Error Types
```go
type ShipError struct {
    Code    string
    Message string
    Details map[string]interface{}
    Err     error
}

const (
    ErrAuthFailed      = "AUTH_FAILED"
    ErrFileTooLarge    = "FILE_TOO_LARGE"
    ErrNetworkTimeout  = "NETWORK_TIMEOUT"
    ErrDockerNotFound  = "DOCKER_NOT_FOUND"
    ErrProviderMissing = "PROVIDER_MISSING"
)
```

### User-Friendly Messages
- Clear error descriptions
- Suggested fixes
- Links to documentation
- Support contact for critical errors

## Testing Strategy

### Unit Test Coverage
- Target: 80% coverage
- Critical paths: 95% coverage
- Mock external dependencies

### Integration Tests
- Real API calls (staging environment)
- Docker container tests
- File system operations

### E2E Scenarios
1. Complete auth flow
2. Push various artifact types
3. Multi-provider investigation
4. MCP tool execution
5. Error recovery scenarios