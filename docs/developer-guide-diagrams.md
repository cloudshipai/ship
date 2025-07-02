# Developer Guide Diagrams - Ship CLI

This guide provides visual diagrams to help developers understand how to extend and contribute to Ship CLI.

## Adding a New Command

```mermaid
flowchart TD
    START[Want to Add New Command] --> DECIDE{Type of Command?}
    
    DECIDE -->|Tool Integration| TOOL_PATH
    DECIDE -->|AI Feature| AI_PATH
    DECIDE -->|Utility| UTIL_PATH
    
    subgraph "Tool Integration Path"
        TOOL_PATH[Create Tool Module]
        TOOL_PATH --> TOOL_MODULE[internal/dagger/modules/]
        TOOL_MODULE --> TOOL_CMD[internal/cli/]
        TOOL_CMD --> TOOL_TEST[Write Tests]
    end
    
    subgraph "AI Feature Path"
        AI_PATH[Design AI Flow]
        AI_PATH --> AI_PROMPT[Create Prompts]
        AI_PROMPT --> AI_MODULE[Add to LLM Module]
        AI_MODULE --> AI_CMD[Create Command]
        AI_CMD --> AI_TEST[Test with LLMs]
    end
    
    subgraph "Utility Path"
        UTIL_PATH[Design Utility]
        UTIL_PATH --> UTIL_IMPL[Implement Logic]
        UTIL_IMPL --> UTIL_CMD[Add Command]
        UTIL_CMD --> UTIL_TEST[Unit Tests]
    end
    
    TOOL_TEST --> INTEGRATE
    AI_TEST --> INTEGRATE
    UTIL_TEST --> INTEGRATE
    
    INTEGRATE[Add to Root Command]
    INTEGRATE --> DOC[Update Documentation]
    DOC --> PR[Create Pull Request]
```

## Creating a New Dagger Module

```mermaid
classDiagram
    class YourModule {
        -client *dagger.Client
        +NewYourModule(client) *YourModule
        +RunAnalysis(ctx, path, options) (string, error)
        +GetVersion(ctx) (string, error)
        -setupContainer() *dagger.Container
        -executeCommand(container, cmd) (string, error)
    }
    
    class ModuleInterface {
        <<interface>>
        +RunAnalysis(ctx, path, options) (string, error)
        +GetVersion(ctx) (string, error)
    }
    
    class ExampleImplementation {
        +setupContainer() *dagger.Container
        +executeCommand() (string, error)
        +parseOutput(output) (Result, error)
        +handleErrors(err) error
    }
    
    ModuleInterface <|.. YourModule
    YourModule --> ExampleImplementation
```

### Module Implementation Example

```mermaid
sequenceDiagram
    participant CLI
    participant Module
    participant Dagger
    participant Container
    participant Tool
    
    CLI->>Module: NewYourModule(daggerClient)
    Module-->>CLI: module instance
    
    CLI->>Module: RunAnalysis(ctx, path, options)
    Module->>Module: setupContainer()
    Module->>Dagger: Container().From("base:image")
    Dagger-->>Module: container
    
    Module->>Dagger: WithDirectory("/workspace", hostDir)
    Module->>Dagger: WithEnvVariable("KEY", "value")
    Module->>Dagger: WithExec(["tool", "command"])
    
    Module->>Container: Stdout(ctx)
    Container->>Tool: Execute command
    Tool-->>Container: Output
    Container-->>Module: stdout
    
    Module->>Module: parseOutput(stdout)
    Module-->>CLI: results
```

## Adding CloudShip Integration

```mermaid
flowchart LR
    subgraph "Your Tool"
        EXEC[Execute Tool]
        RESULT[Get Results]
        FORMAT[Format Output]
    end
    
    subgraph "Integration Steps"
        ADD_FLAG[Add --push Flag]
        CHECK_FLAG[Check Flag in RunE]
        CALL_PUSH[Call pushToCloudShip]
    end
    
    subgraph "Push Function"
        LOAD_CFG[Load Config]
        CHECK_AUTH[Check API Key]
        ENCODE[Base64 Encode]
        BUILD_REQ[Build Request]
        SEND[Send to API]
        SHOW_URL[Display URL]
    end
    
    EXEC --> RESULT
    RESULT --> FORMAT
    FORMAT --> CHECK_FLAG
    
    CHECK_FLAG -->|--push set| CALL_PUSH
    CHECK_FLAG -->|no --push| END[Display Results]
    
    CALL_PUSH --> LOAD_CFG
    LOAD_CFG --> CHECK_AUTH
    CHECK_AUTH --> ENCODE
    ENCODE --> BUILD_REQ
    BUILD_REQ --> SEND
    SEND --> SHOW_URL
    SHOW_URL --> END
```

## Adding AI Table Support

```mermaid
graph TD
    subgraph "When AI Needs New Tables"
        NEW_PROVIDER[New Cloud Provider]
        NEW_SERVICE[New Service Tables]
        UPDATE[Update Existing]
    end
    
    subgraph "Update Process"
        EDIT_TABLES[Edit steampipe_tables.go]
        ADD_TABLES[Add Table Names]
        ADD_EXAMPLES[Add Query Examples]
        TEST_AI[Test AI Prompts]
    end
    
    subgraph "Table Definition"
        TABLE_NAME[Table Name]
        TABLE_SCHEMA[Common Columns]
        EXAMPLE_QUERY[Example Query]
    end
    
    NEW_PROVIDER --> EDIT_TABLES
    NEW_SERVICE --> EDIT_TABLES
    UPDATE --> EDIT_TABLES
    
    EDIT_TABLES --> ADD_TABLES
    EDIT_TABLES --> ADD_EXAMPLES
    
    ADD_TABLES --> TABLE_NAME
    ADD_TABLES --> TABLE_SCHEMA
    ADD_EXAMPLES --> EXAMPLE_QUERY
    
    TABLE_NAME --> TEST_AI
    TABLE_SCHEMA --> TEST_AI
    EXAMPLE_QUERY --> TEST_AI
```

## Testing Strategy

```mermaid
graph TB
    subgraph "Test Types"
        UNIT[Unit Tests]
        INTEGRATION[Integration Tests]
        E2E[End-to-End Tests]
        AI_TEST[AI Tests]
    end
    
    subgraph "Unit Test Targets"
        PARSE[Parsing Logic]
        FORMAT[Formatting]
        VALIDATION[Validation]
    end
    
    subgraph "Integration Test Targets"
        DAGGER_INT[Dagger Integration]
        API_INT[API Integration]
        TOOL_INT[Tool Integration]
    end
    
    subgraph "E2E Test Scenarios"
        FULL_FLOW[Complete Workflows]
        ERROR_FLOW[Error Scenarios]
        PUSH_FLOW[CloudShip Push]
    end
    
    subgraph "AI Test Scenarios"
        PROMPT_TEST[Prompt Generation]
        TABLE_TEST[Table Recognition]
        FALLBACK_TEST[Fallback Behavior]
    end
    
    UNIT --> PARSE
    UNIT --> FORMAT
    UNIT --> VALIDATION
    
    INTEGRATION --> DAGGER_INT
    INTEGRATION --> API_INT
    INTEGRATION --> TOOL_INT
    
    E2E --> FULL_FLOW
    E2E --> ERROR_FLOW
    E2E --> PUSH_FLOW
    
    AI_TEST --> PROMPT_TEST
    AI_TEST --> TABLE_TEST
    AI_TEST --> FALLBACK_TEST
```

## Configuration Management

```mermaid
classDiagram
    class Config {
        +APIKey string
        +FleetID string
        +DefaultRegion string
        +LLMProvider string
        +Load() (*Config, error)
        +Save() error
        +Validate() error
    }
    
    class ConfigSources {
        <<enumeration>>
        FILE
        ENV_VAR
        CLI_FLAG
        DEFAULT
    }
    
    class ConfigLoader {
        +LoadFromFile(path) Config
        +LoadFromEnv() Config
        +MergeConfigs(configs) Config
        +GetConfigPath() string
    }
    
    class Precedence {
        1. CLI Flags
        2. Environment Variables
        3. Config File
        4. Defaults
    }
    
    Config --> ConfigSources
    ConfigLoader --> Config
    ConfigLoader --> Precedence
```

## Error Handling Patterns

```mermaid
flowchart TD
    subgraph "Error Types"
        AUTH_ERR[Authentication Error]
        DOCKER_ERR[Docker Not Running]
        TIMEOUT_ERR[Operation Timeout]
        PARSE_ERR[Parse Error]
        API_ERR[API Error]
    end
    
    subgraph "Error Handling"
        CATCH[Catch Error]
        CLASSIFY[Classify Type]
        ENHANCE[Add Context]
        SUGGEST[Suggest Fix]
        RETURN[Return to User]
    end
    
    subgraph "User Messages"
        AUTH_MSG[Run 'ship auth']
        DOCKER_MSG[Start Docker]
        TIMEOUT_MSG[Try with --timeout]
        PARSE_MSG[Check input format]
        API_MSG[Check credentials]
    end
    
    AUTH_ERR --> CATCH
    DOCKER_ERR --> CATCH
    TIMEOUT_ERR --> CATCH
    PARSE_ERR --> CATCH
    API_ERR --> CATCH
    
    CATCH --> CLASSIFY
    CLASSIFY --> ENHANCE
    ENHANCE --> SUGGEST
    
    SUGGEST --> AUTH_MSG
    SUGGEST --> DOCKER_MSG
    SUGGEST --> TIMEOUT_MSG
    SUGGEST --> PARSE_MSG
    SUGGEST --> API_MSG
    
    AUTH_MSG --> RETURN
    DOCKER_MSG --> RETURN
    TIMEOUT_MSG --> RETURN
    PARSE_MSG --> RETURN
    API_MSG --> RETURN
```

## Contributing Workflow

```mermaid
gitGraph
    commit id: "main"
    branch feature/your-feature
    checkout feature/your-feature
    commit id: "Add module"
    commit id: "Add command"
    commit id: "Add tests"
    commit id: "Update docs"
    checkout main
    merge feature/your-feature
    commit id: "Release v1.x"
```

## Module Communication

```mermaid
sequenceDiagram
    participant Command
    participant Engine
    participant Module
    participant Container
    participant CloudShip
    
    Command->>Engine: Initialize Dagger
    Engine-->>Command: client
    
    Command->>Module: NewModule(client)
    Module-->>Command: instance
    
    Command->>Module: Execute(params)
    Module->>Container: Setup & Run
    Container-->>Module: results
    
    Module-->>Command: output
    
    opt Push to CloudShip
        Command->>CloudShip: Upload(results)
        CloudShip-->>Command: URL
    end
    
    Command-->>User: Display
```

## Debugging Ship CLI

```mermaid
flowchart TD
    subgraph "Debug Levels"
        INFO[--log-level info]
        DEBUG[--log-level debug]
        TRACE[--log-level trace]
    end
    
    subgraph "Debug Tools"
        DAGGER_LOG[DAGGER_LOG=debug]
        SHIP_DEBUG[SHIP_DEBUG=1]
        VERBOSE[-v flag]
    end
    
    subgraph "Common Issues"
        NO_OUTPUT[No Output]
        TIMEOUT[Timeouts]
        AUTH_FAIL[Auth Failures]
        CONTAINER_FAIL[Container Errors]
    end
    
    subgraph "Debug Steps"
        CHECK_DOCKER[Check Docker]
        CHECK_CREDS[Check Credentials]
        CHECK_NETWORK[Check Network]
        CHECK_LOGS[Check Logs]
    end
    
    NO_OUTPUT --> DEBUG
    TIMEOUT --> TRACE
    AUTH_FAIL --> CHECK_CREDS
    CONTAINER_FAIL --> CHECK_DOCKER
    
    DEBUG --> CHECK_LOGS
    TRACE --> DAGGER_LOG
    CHECK_CREDS --> SHIP_DEBUG
    CHECK_DOCKER --> VERBOSE
```

## Performance Optimization

```mermaid
graph LR
    subgraph "Optimization Areas"
        CACHE[Container Caching]
        PARALLEL[Parallel Execution]
        LAZY[Lazy Loading]
        STREAM[Stream Processing]
    end
    
    subgraph "Caching Strategy"
        IMG_CACHE[Image Cache]
        BUILD_CACHE[Build Cache]
        RESULT_CACHE[Result Cache]
    end
    
    subgraph "Parallel Patterns"
        MULTI_QUERY[Multiple Queries]
        MULTI_TOOL[Multiple Tools]
        BATCH_PROC[Batch Processing]
    end
    
    subgraph "Performance Gains"
        STARTUP[Faster Startup]
        EXEC[Faster Execution]
        MEMORY[Lower Memory]
    end
    
    CACHE --> IMG_CACHE
    CACHE --> BUILD_CACHE
    CACHE --> RESULT_CACHE
    
    PARALLEL --> MULTI_QUERY
    PARALLEL --> MULTI_TOOL
    PARALLEL --> BATCH_PROC
    
    IMG_CACHE --> STARTUP
    BUILD_CACHE --> EXEC
    RESULT_CACHE --> EXEC
    MULTI_QUERY --> EXEC
    LAZY --> MEMORY
    STREAM --> MEMORY
```

---

This developer guide provides:
- Step-by-step visual guides for adding new features
- Code structure and patterns to follow
- Testing strategies
- Configuration management approaches
- Error handling patterns
- Contributing workflow
- Debugging techniques
- Performance optimization strategies

These diagrams serve as a visual reference for developers looking to understand and extend Ship CLI.