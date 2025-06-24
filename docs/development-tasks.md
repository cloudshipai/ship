# Ship CLI Development Tasks

## Immediate Next Steps (Sprint 1)

### 1. Project Initialization
```bash
# Commands to run:
go mod init github.com/cloudship/ship
go get github.com/spf13/cobra
go get github.com/spf13/viper
go get github.com/fatih/color
go get github.com/schollz/progressbar/v3
```

### 2. Core Structure Setup

#### Task 1: Create Base CLI Structure
- [ ] Create `cmd/ship/main.go`
- [ ] Set up Cobra root command
- [ ] Add version information
- [ ] Implement help system

#### Task 2: Config Module
- [ ] Create `internal/config/config.go`
- [ ] Implement YAML parsing
- [ ] Add environment variable support
- [ ] Create config file writer
- [ ] Add config validation

#### Task 3: Auth Command
- [ ] Create `internal/cli/auth.go`
- [ ] Implement `--token` flag
- [ ] Add token validation
- [ ] Implement `--logout` flag
- [ ] Create secure file permissions

### 3. First Testable Milestone

**Goal**: Working `ship auth` command that:
1. Accepts a token via flag
2. Validates it with a mock endpoint
3. Saves to `~/.ship/config.yaml`
4. Can logout and clear credentials

**Test Commands**:
```bash
# Build
go build -o ship ./cmd/ship

# Test auth
./ship auth --token sk_test_123
./ship auth --logout

# Verify config
cat ~/.ship/config.yaml
```

## Sprint 2: Push Command

### Core Components Needed

1. **HTTP Client** (`internal/client/`)
   - Retryable client
   - Auth middleware
   - Progress tracking

2. **File Handler** (`internal/push/`)
   - Size validation
   - SHA-256 calculation
   - Type detection

3. **Push Command** (`internal/cli/push.go`)
   - File/stdin input
   - Flag parsing
   - Upload orchestration

### Implementation Order

1. **Day 1-2**: HTTP client with auth
2. **Day 3-4**: File handling utilities
3. **Day 5-6**: Push command integration
4. **Day 7**: Testing and refinement

## Sprint 3: Dagger Integration

### Prerequisites
- Research Dagger Go SDK
- Design embedding strategy
- Plan container architecture

### Tasks
1. **Embed Dagger Binary**
   - Build script for embedding
   - Version management
   - Platform-specific builds

2. **Engine Wrapper**
   - Initialize Dagger
   - Container lifecycle
   - Output streaming

3. **Provider Containers**
   - AWS Steampipe image
   - Cloudflare Steampipe image
   - Heroku Steampipe image

## Testing Checklist

### Unit Tests (per module)
```go
// Example test structure
func TestAuthCommand(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        want    error
        wantCfg *Config
    }{
        // Test cases
    }
}
```

### Integration Test Scenarios
1. **Auth Flow**
   - Valid token → save → use in request
   - Invalid token → error
   - Logout → clear config

2. **Push Flow**
   - Small file → success
   - Large file → rejection
   - Network failure → retry

3. **Investigate Flow**
   - Goals fetch → query mapping → execution
   - Provider plugin installation
   - Result bundling

## Code Quality Gates

### Pre-commit Checks
```bash
# Format
go fmt ./...

# Lint
golangci-lint run

# Test
go test ./...

# Build
go build -o ship ./cmd/ship
```

### CI Pipeline Requirements
- All tests pass
- 80% code coverage
- No linting errors
- Binary builds for all platforms
- Size under 50MB

## Documentation Requirements

### For Each Feature
1. **User Documentation**
   - Command examples
   - Flag descriptions
   - Common workflows

2. **API Documentation**
   - Endpoint specs
   - Request/response examples
   - Error codes

3. **Code Documentation**
   - Package comments
   - Public API godocs
   - Complex algorithm explanations

## Dependency Management

### Approved Dependencies
- CLI: `cobra`, `viper`
- HTTP: `standard library`
- Progress: `progressbar`
- Colors: `fatih/color`
- Testing: `testify`

### Evaluation Criteria for New Deps
1. License compatibility (Apache 2.0 preferred)
2. Maintenance status
3. Binary size impact
4. Security track record

## Release Checklist

### Version 0.1.0 (MVP)
- [ ] Auth command working
- [ ] Push command working
- [ ] Basic investigate (hardcoded queries)
- [ ] Linux/macOS binaries
- [ ] Installation script
- [ ] Basic documentation

### Version 0.2.0
- [ ] Full investigate with goals API
- [ ] All 3 providers (AWS, Cloudflare, Heroku)
- [ ] Windows support
- [ ] Homebrew formula

### Version 1.0.0
- [ ] MCP server
- [ ] AI integration
- [ ] Production-ready error handling
- [ ] Comprehensive documentation
- [ ] Performance optimizations