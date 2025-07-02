# Data Flow and State Management - Ship CLI

This document illustrates how data flows through Ship CLI and how state is managed across different components.

## Overall Data Flow

```mermaid
flowchart TB
    subgraph "Input Layer"
        USER_INPUT[User Command]
        ENV_VARS[Environment Variables]
        CONFIG_FILE[Config File]
        CREDS_FILE[Credential Files]
    end
    
    subgraph "Processing Layer"
        CLI_PARSER[CLI Parser]
        CONFIG_MERGER[Config Merger]
        VALIDATOR[Input Validator]
        TRANSFORMER[Data Transformer]
    end
    
    subgraph "Execution Layer"
        DAGGER_ENGINE[Dagger Engine]
        CONTAINER_EXEC[Container Execution]
        RESULT_PARSER[Result Parser]
    end
    
    subgraph "Output Layer"
        FORMATTER[Output Formatter]
        FILE_WRITER[File Writer]
        API_SENDER[API Sender]
        DISPLAY[Terminal Display]
    end
    
    USER_INPUT --> CLI_PARSER
    ENV_VARS --> CONFIG_MERGER
    CONFIG_FILE --> CONFIG_MERGER
    CREDS_FILE --> CONFIG_MERGER
    
    CLI_PARSER --> VALIDATOR
    CONFIG_MERGER --> VALIDATOR
    VALIDATOR --> TRANSFORMER
    
    TRANSFORMER --> DAGGER_ENGINE
    DAGGER_ENGINE --> CONTAINER_EXEC
    CONTAINER_EXEC --> RESULT_PARSER
    
    RESULT_PARSER --> FORMATTER
    FORMATTER --> FILE_WRITER
    FORMATTER --> API_SENDER
    FORMATTER --> DISPLAY
```

## Command Execution State Machine

```mermaid
stateDiagram-v2
    [*] --> ParseCommand
    
    ParseCommand --> ValidateInput
    ValidateInput --> LoadConfig: Valid
    ValidateInput --> ShowError: Invalid
    
    ShowError --> [*]
    
    LoadConfig --> CheckDocker
    CheckDocker --> InitDagger: Running
    CheckDocker --> DockerError: Not Running
    
    DockerError --> [*]
    
    InitDagger --> CreateModule
    CreateModule --> SetupContainer
    SetupContainer --> ExecuteTool
    
    ExecuteTool --> ProcessOutput
    ProcessOutput --> CheckPushFlag
    
    CheckPushFlag --> PushToCloudShip: --push
    CheckPushFlag --> FormatOutput: No Push
    
    PushToCloudShip --> FormatOutput
    FormatOutput --> DisplayResult
    DisplayResult --> [*]
```

## Configuration Precedence Flow

```mermaid
flowchart LR
    subgraph "Priority Order"
        direction TB
        P1[1. CLI Flags]
        P2[2. Environment Vars]
        P3[3. Config File]
        P4[4. Defaults]
        P1 --> P2
        P2 --> P3
        P3 --> P4
    end
    
    subgraph "Merge Process"
        CHECK_FLAG{Has Flag?}
        CHECK_ENV{Has Env?}
        CHECK_FILE{Has File?}
        USE_DEFAULT[Use Default]
    end
    
    subgraph "Final Config"
        MERGED[Merged Config]
        VALIDATE[Validate]
        READY[Ready to Use]
    end
    
    P1 --> CHECK_FLAG
    CHECK_FLAG -->|No| CHECK_ENV
    CHECK_FLAG -->|Yes| MERGED
    
    CHECK_ENV -->|No| CHECK_FILE
    CHECK_ENV -->|Yes| MERGED
    
    CHECK_FILE -->|No| USE_DEFAULT
    CHECK_FILE -->|Yes| MERGED
    
    USE_DEFAULT --> MERGED
    MERGED --> VALIDATE
    VALIDATE --> READY
```

## Credential Resolution Flow

```mermaid
flowchart TD
    subgraph "AWS Credential Sources"
        ENV_KEYS[Environment Variables]
        PROFILE_FLAG[--aws-profile Flag]
        CRED_FILE[~/.aws/credentials]
        INSTANCE_ROLE[Instance Role]
    end
    
    subgraph "Resolution Process"
        CHECK_ENV{ENV Set?}
        PARSE_FILE[Parse Credentials]
        EXTRACT_PROFILE[Extract Profile]
        BUILD_MAP[Build Credential Map]
    end
    
    subgraph "Container Injection"
        SET_ENV[Set Container Env]
        MOUNT_FILES[Mount Files]
        CONFIGURE[Configure SDK]
    end
    
    ENV_KEYS --> CHECK_ENV
    CHECK_ENV -->|Yes| BUILD_MAP
    CHECK_ENV -->|No| PARSE_FILE
    
    PROFILE_FLAG --> EXTRACT_PROFILE
    CRED_FILE --> PARSE_FILE
    PARSE_FILE --> EXTRACT_PROFILE
    EXTRACT_PROFILE --> BUILD_MAP
    
    BUILD_MAP --> SET_ENV
    SET_ENV --> CONFIGURE
    
    INSTANCE_ROLE --> CONFIGURE
```

## AI Investigation Data Flow

```mermaid
flowchart LR
    subgraph "Input Processing"
        PROMPT[User Prompt]
        CONTEXT[Context Builder]
        TABLES[Table Registry]
    end
    
    subgraph "LLM Processing"
        ENHANCE[Enhance Prompt]
        SEND_LLM[Send to LLM]
        PARSE_RESP[Parse Response]
    end
    
    subgraph "Plan Execution"
        PLAN[Investigation Plan]
        STEPS[Individual Steps]
        QUERIES[SQL Queries]
    end
    
    subgraph "Result Assembly"
        RAW_RESULTS[Raw Results]
        INSIGHTS[Extract Insights]
        SUMMARY[Build Summary]
    end
    
    PROMPT --> CONTEXT
    TABLES --> CONTEXT
    CONTEXT --> ENHANCE
    
    ENHANCE --> SEND_LLM
    SEND_LLM --> PARSE_RESP
    PARSE_RESP --> PLAN
    
    PLAN --> STEPS
    STEPS --> QUERIES
    QUERIES --> RAW_RESULTS
    
    RAW_RESULTS --> INSIGHTS
    INSIGHTS --> SUMMARY
```

## Container Lifecycle Management

```mermaid
stateDiagram-v2
    [*] --> ContainerRequested
    
    ContainerRequested --> CheckCache
    CheckCache --> CacheHit: Found
    CheckCache --> CacheMiss: Not Found
    
    CacheHit --> ContainerReady
    CacheMiss --> PullImage
    
    PullImage --> CreateContainer
    CreateContainer --> ConfigureContainer
    
    ConfigureContainer --> MountVolumes
    MountVolumes --> SetEnvironment
    SetEnvironment --> ContainerReady
    
    ContainerReady --> ExecuteCommand
    ExecuteCommand --> CaptureOutput
    CaptureOutput --> ProcessResults
    
    ProcessResults --> CleanupDecision
    CleanupDecision --> KeepContainer: Reusable
    CleanupDecision --> RemoveContainer: One-time
    
    KeepContainer --> [*]
    RemoveContainer --> [*]
```

## CloudShip Push Data Flow

```mermaid
sequenceDiagram
    participant Tool
    participant CLI
    participant Config
    participant Encoder
    participant API
    participant User
    
    Tool->>CLI: Tool execution complete
    CLI->>CLI: Check --push flag
    
    alt Push enabled
        CLI->>Config: Load API credentials
        Config-->>CLI: APIKey, FleetID
        
        CLI->>Encoder: Encode artifact
        Encoder-->>CLI: Base64 data
        
        CLI->>CLI: Build metadata
        CLI->>API: POST /v1/artifacts
        
        alt Success
            API-->>CLI: {url: "analysis_url"}
            CLI-->>User: Display URL
        else Failure
            API-->>CLI: Error
            CLI-->>User: Show error
        end
    else No push
        CLI-->>User: Display results only
    end
```

## Error Propagation Flow

```mermaid
flowchart TD
    subgraph "Error Sources"
        DOCKER_ERR[Docker Error]
        PARSE_ERR[Parse Error]
        EXEC_ERR[Execution Error]
        API_ERR[API Error]
        AUTH_ERR[Auth Error]
    end
    
    subgraph "Error Handling"
        CATCH[Error Handler]
        WRAP[Wrap Context]
        LOG[Log Error]
        FORMAT[Format Message]
    end
    
    subgraph "User Feedback"
        ERROR_MSG[Error Message]
        SUGGESTION[Suggested Fix]
        HELP_CMD[Help Command]
    end
    
    DOCKER_ERR --> CATCH
    PARSE_ERR --> CATCH
    EXEC_ERR --> CATCH
    API_ERR --> CATCH
    AUTH_ERR --> CATCH
    
    CATCH --> WRAP
    WRAP --> LOG
    LOG --> FORMAT
    
    FORMAT --> ERROR_MSG
    FORMAT --> SUGGESTION
    FORMAT --> HELP_CMD
```

## Module Communication Protocol

```mermaid
sequenceDiagram
    participant CLI as CLI Command
    participant Engine as Dagger Engine
    participant Module as Tool Module
    participant Container as Container
    participant FS as File System
    
    CLI->>Engine: Initialize
    Engine->>Module: Create(client)
    
    CLI->>Module: Execute(params)
    Module->>Container: From(image)
    Module->>Container: WithDirectory(host, container)
    
    alt Need credentials
        Module->>Container: WithEnvVariable(key, value)
    end
    
    Module->>Container: WithExec(command)
    Container->>Container: Run command
    
    Container->>FS: Write output
    Container->>Module: Return stdout/stderr
    
    Module->>Module: Parse output
    Module->>CLI: Return results
    
    CLI->>CLI: Format output
    CLI->>User: Display
```

## State Persistence

```mermaid
graph TB
    subgraph "Persistent State"
        CONFIG_STATE[~/.ship/config.yaml]
        CACHE_STATE[~/.ship/cache/]
        LOGS_STATE[~/.ship/logs/]
    end
    
    subgraph "Session State"
        CMD_FLAGS[Command Flags]
        ENV_STATE[Environment]
        RUNTIME_STATE[Runtime Config]
    end
    
    subgraph "Container State"
        CONTAINER_FS[Container FS]
        CONTAINER_ENV[Container Env]
        CONTAINER_NET[Container Network]
    end
    
    CONFIG_STATE --> RUNTIME_STATE
    CMD_FLAGS --> RUNTIME_STATE
    ENV_STATE --> RUNTIME_STATE
    
    RUNTIME_STATE --> CONTAINER_ENV
    CACHE_STATE --> CONTAINER_FS
    
    CONTAINER_FS --> LOGS_STATE
```

## Output Formatting Pipeline

```mermaid
flowchart LR
    subgraph "Raw Output"
        JSON_OUT[JSON Output]
        TEXT_OUT[Text Output]
        ERROR_OUT[Error Output]
    end
    
    subgraph "Parsers"
        JSON_PARSE[JSON Parser]
        REGEX_PARSE[Regex Parser]
        CUSTOM_PARSE[Custom Parser]
    end
    
    subgraph "Formatters"
        TABLE_FMT[Table Format]
        JSON_FMT[JSON Format]
        YAML_FMT[YAML Format]
        MD_FMT[Markdown Format]
    end
    
    subgraph "Writers"
        STDOUT[Stdout]
        FILE[File]
        PIPE[Pipe]
    end
    
    JSON_OUT --> JSON_PARSE
    TEXT_OUT --> REGEX_PARSE
    ERROR_OUT --> CUSTOM_PARSE
    
    JSON_PARSE --> TABLE_FMT
    JSON_PARSE --> JSON_FMT
    REGEX_PARSE --> TABLE_FMT
    CUSTOM_PARSE --> MD_FMT
    
    TABLE_FMT --> STDOUT
    JSON_FMT --> FILE
    YAML_FMT --> FILE
    MD_FMT --> STDOUT
    
    STDOUT --> PIPE
```

## Resource Cleanup Flow

```mermaid
stateDiagram-v2
    [*] --> CommandComplete
    
    CommandComplete --> CheckResources
    
    CheckResources --> HasContainers: Yes
    CheckResources --> NoCleanup: No
    
    HasContainers --> CheckKeepFlag
    CheckKeepFlag --> KeepContainers: --keep-containers
    CheckKeepFlag --> RemoveContainers: Default
    
    RemoveContainers --> StopContainers
    StopContainers --> RemoveImages
    RemoveImages --> CleanCache
    
    KeepContainers --> LogKept
    
    CleanCache --> CleanupComplete
    LogKept --> CleanupComplete
    NoCleanup --> CleanupComplete
    
    CleanupComplete --> [*]
```

---

These data flow diagrams illustrate:
- How data moves through the system
- State management across components
- Configuration resolution and precedence
- Credential handling and injection
- Container lifecycle management
- Error propagation patterns
- Output formatting pipeline
- Resource cleanup processes

This completes the comprehensive visual documentation of Ship CLI's architecture, making it easier for developers to understand the system's behavior and extend it effectively.