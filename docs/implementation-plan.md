# Ship CLI Implementation Plan

## Phase 1: Core Foundation (Week 1-2)

### 1.1 Project Setup
- [ ] Initialize Go module structure
- [ ] Set up CLI framework (Cobra)
- [ ] Configure build pipeline
- [ ] Set up testing framework
- [ ] Configure linting (golangci-lint)

### 1.2 Configuration Management
- [ ] Implement config file handling (`~/.ship/config.yaml`)
- [ ] Environment variable support (`SHIP_TOKEN`)
- [ ] Config command: `ship config`
- [ ] Tests for config management

### 1.3 Authentication Module
- [ ] Implement `ship auth --token` command
- [ ] Token validation with backend
- [ ] Implement `ship auth --logout`
- [ ] Secure token storage
- [ ] Tests for auth flows

### 1.4 HTTP Client Foundation
- [ ] Create HTTP client with retry logic
- [ ] Authentication middleware
- [ ] Error handling and logging
- [ ] Response parsing utilities
- [ ] Tests for HTTP client

## Phase 2: Push Command (Week 2-3)

### 2.1 File Handling
- [ ] File size validation (100MB limit)
- [ ] SHA-256 checksum calculation
- [ ] File type detection logic
- [ ] Stdin support (`-`)

### 2.2 Push Command Implementation
- [ ] Implement `ship push` command
- [ ] Support `--kind` flag
- [ ] Support `--env` flag
- [ ] Support `--tag` flag
- [ ] Multipart upload for large files

### 2.3 Backend Integration
- [ ] API endpoint integration
- [ ] Progress bar for uploads
- [ ] Error handling and retries
- [ ] Success response handling (SHA & URL)
- [ ] Integration tests

## Phase 3: Dagger Integration (Week 3-4)

### 3.1 Dagger Embedding
- [ ] Bundle Dagger CLI binary
- [ ] Dagger engine initialization
- [ ] Container runtime detection (Docker/Podman)
- [ ] Fallback messaging for missing runtime

### 3.2 Container Management
- [ ] Composite image building
- [ ] Volume mounting for credentials
- [ ] Container lifecycle management
- [ ] Output collection from containers

## Phase 4: Investigate Command (Week 4-5)

### 4.1 Goals API Integration
- [ ] Implement goals fetching (`GET /v1/goals`)
- [ ] Goal parsing and validation
- [ ] Environment filtering

### 4.2 Provider Support
- [ ] AWS provider implementation
- [ ] Cloudflare provider implementation
- [ ] Heroku provider implementation
- [ ] Provider plugin architecture

### 4.3 Steampipe Integration
- [ ] Steampipe container setup
- [ ] Plugin auto-installation
- [ ] SQL query execution
- [ ] Result collection and bundling

### 4.4 Investigation Orchestration
- [ ] Progress tracking and display
- [ ] Error handling and reporting
- [ ] Result bundling and compression
- [ ] Automatic artifact push

## Phase 5: MCP Host (Week 6-7)

### 5.1 MCP Server Implementation
- [ ] JSON-RPC server (TCP & stdio)
- [ ] Tool registration system
- [ ] Request routing

### 5.2 MCP Tools
- [ ] Implement `goal.list` tool
- [ ] Implement `goal.run` tool
- [ ] Implement `steampipe.query` tool
- [ ] Tool schema generation

### 5.3 Integration Features
- [ ] Hot-reload container support
- [ ] Session management
- [ ] Tool documentation generator

## Phase 6: AI Integration (Week 6-7)

### 6.1 LLM Support
- [ ] OpenAI integration (optional flag)
- [ ] Local LLM support investigation
- [ ] Goal-to-SQL mapping logic
- [ ] Prompt template system

### 6.2 Built-in Templates
- [ ] Default SQL templates per provider
- [ ] Template selection logic
- [ ] Template customization support

## Phase 7: Distribution & Documentation (Week 8)

### 7.1 Build & Distribution
- [ ] Cross-platform builds (Linux, macOS, Windows)
- [ ] Binary size optimization
- [ ] Homebrew formula
- [ ] Installation scripts

### 7.2 Documentation
- [ ] CLI reference documentation
- [ ] Getting started guide
- [ ] CI/CD integration examples
- [ ] MCP integration guide
- [ ] Troubleshooting guide

### 7.3 Testing & Quality
- [ ] End-to-end testing suite
- [ ] Performance benchmarks
- [ ] Security audit
- [ ] Beta testing program

## Phase 8: Telemetry & Polish (Week 8)

### 8.1 Telemetry
- [ ] Anonymous metrics collection
- [ ] First-run opt-in prompt
- [ ] Telemetry configuration

### 8.2 User Experience
- [ ] Enhanced error messages
- [ ] Interactive prompts
- [ ] Command suggestions
- [ ] Update notifications

## Testing Strategy

### Unit Tests
- Config management
- Authentication
- File handling
- HTTP client
- Goal parsing

### Integration Tests
- Auth flow with backend
- File uploads
- Goal API integration
- Steampipe queries
- MCP tool invocations

### End-to-End Tests
- Complete auth → push → investigate flow
- MCP server with tool execution
- Multi-provider investigations
- CI/CD scenarios

## Risk Mitigation

1. **Dagger Binary Size**: Monitor total binary size, consider dynamic download if >50MB
2. **Windows Support**: Early testing on WSL, clear documentation for Docker Desktop
3. **Provider Licensing**: Audit all dependencies, separate AGPL components
4. **Performance**: Implement caching for repeated operations, connection pooling