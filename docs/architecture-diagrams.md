# Ship CLI Architecture Diagrams

This document provides visual representations of the Ship CLI architecture, workflows, and module system using Mermaid diagrams.

## Table of Contents
1. [High-Level Architecture](#high-level-architecture)
2. [Core Components](#core-components)
3. [User Workflows](#user-workflows)
4. [Module System](#module-system)
5. [AI Investigation Flow](#ai-investigation-flow)
6. [CloudShip Integration](#cloudship-integration)
7. [Dagger Container Orchestration](#dagger-container-orchestration)
8. [Adding New Modules](#adding-new-modules)

## High-Level Architecture

```mermaid
graph TB
    subgraph "User Interface"
        CLI[Ship CLI]
        User[User]
    end
    
    subgraph "Core Components"
        CMD[Command Layer]
        CONFIG[Configuration]
        AUTH[Authentication]
    end
    
    subgraph "Execution Engine"
        DAGGER[Dagger Engine]
        CONTAINERS[Container Runtime]
    end
    
    subgraph "Tool Modules"
        TF[Terraform Tools]
        STEAMPIPE[Steampipe]
        AI[AI Services]
        MCP[MCP Server]
    end
    
    subgraph "External Services"
        CLOUDSHIP[CloudShip API]
        AWS[AWS]
        AZURE[Azure]
        GCP[GCP]
        OPENAI[OpenAI/Anthropic]
    end
    
    User --> CLI
    CLI --> CMD
    CMD --> CONFIG
    CMD --> AUTH
    CMD --> DAGGER
    DAGGER --> CONTAINERS
    CONTAINERS --> TF
    CONTAINERS --> STEAMPIPE
    CONTAINERS --> AI
    CONTAINERS --> MCP
    
    AUTH --> CLOUDSHIP
    STEAMPIPE --> AWS
    STEAMPIPE --> AZURE
    STEAMPIPE --> GCP
    AI --> OPENAI
    TF --> CLOUDSHIP
```

## Core Components

```mermaid
classDiagram
    class RootCommand {
        +Execute()
        +AddCommand()
    }
    
    class TerraformToolsCmd {
        +lint()
        +securityScan()
        +costEstimate()
        +generateDocs()
        +generateDiagram()
        +checkovScan()
        +costAnalysis()
    }
    
    class AIInvestigateCmd {
        +runInvestigation()
        +createPlan()
        +executePlan()
    }
    
    class QueryCmd {
        +runQuery()
        +parseCredentials()
    }
    
    class AuthCmd {
        +authenticate()
        +logout()
    }
    
    class PushCmd {
        +uploadArtifact()
        +validateFile()
    }
    
    class Config {
        +Load()
        +Save()
        +APIKey string
        +FleetID string
    }
    
    class DaggerEngine {
        +NewEngine()
        +NewSteampipeModule()
        +NewTerraformModule()
        +NewLLMModule()
        +Close()
    }
    
    RootCommand <|-- TerraformToolsCmd
    RootCommand <|-- AIInvestigateCmd
    RootCommand <|-- QueryCmd
    RootCommand <|-- AuthCmd
    RootCommand <|-- PushCmd
    
    TerraformToolsCmd --> DaggerEngine
    AIInvestigateCmd --> DaggerEngine
    QueryCmd --> DaggerEngine
    AuthCmd --> Config
    PushCmd --> Config
```

## User Workflows

### Terraform Analysis Workflow

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Dagger
    participant Container
    participant Tool
    participant CloudShip
    
    User->>CLI: ship terraform-tools lint --push
    CLI->>Dagger: Initialize engine
    Dagger->>Container: Create TFLint container
    Container->>Tool: Run linting
    Tool-->>Container: Results
    Container-->>Dagger: Output
    Dagger-->>CLI: Formatted results
    
    opt --push flag
        CLI->>CloudShip: Upload results
        CloudShip-->>CLI: Analysis URL
    end
    
    CLI-->>User: Display results
```

### AI Investigation Workflow

```mermaid
flowchart LR
    subgraph "Input"
        PROMPT[Natural Language Prompt]
        PROVIDER[Cloud Provider]
    end
    
    subgraph "AI Processing"
        LLM[LLM Service]
        PLAN[Investigation Plan]
        TABLES[Steampipe Tables]
    end
    
    subgraph "Execution"
        QUERY[SQL Queries]
        STEAMPIPE[Steampipe Engine]
        RESULTS[Query Results]
    end
    
    subgraph "Output"
        INSIGHTS[Insights]
        NEXT[Next Steps]
    end
    
    PROMPT --> LLM
    PROVIDER --> LLM
    TABLES --> LLM
    LLM --> PLAN
    PLAN --> QUERY
    QUERY --> STEAMPIPE
    STEAMPIPE --> RESULTS
    RESULTS --> INSIGHTS
    INSIGHTS --> NEXT
```

## Module System

```mermaid
graph TB
    subgraph "Dagger Module Interface"
        MODULE[Module Interface]
        MODULE --> INIT[Initialize]
        MODULE --> EXEC[Execute]
        MODULE --> OUTPUT[Get Output]
    end
    
    subgraph "Built-in Modules"
        STEAMPIPE_MOD[SteampipeModule]
        TF_MOD[TerraformModule]
        LLM_MOD[LLMModule]
        INFRA_MOD[InfraMapModule]
        COST_MOD[CostModule]
    end
    
    subgraph "Container Configuration"
        BASE[Base Image]
        ENV[Environment Vars]
        MOUNT[Volume Mounts]
        EXEC_CMD[Exec Commands]
    end
    
    MODULE --> STEAMPIPE_MOD
    MODULE --> TF_MOD
    MODULE --> LLM_MOD
    MODULE --> INFRA_MOD
    MODULE --> COST_MOD
    
    STEAMPIPE_MOD --> BASE
    TF_MOD --> ENV
    LLM_MOD --> MOUNT
    INFRA_MOD --> EXEC_CMD
```

## AI Investigation Flow

```mermaid
stateDiagram-v2
    [*] --> ReceivePrompt
    ReceivePrompt --> LoadTables: Get Steampipe Tables
    LoadTables --> GeneratePlan: Send to LLM
    GeneratePlan --> ParseJSON: Extract Steps
    
    ParseJSON --> ValidPlan: Success
    ParseJSON --> Fallback: Parse Error
    
    ValidPlan --> ExecuteSteps
    Fallback --> ExecuteSteps: Use Default Query
    
    ExecuteSteps --> RunQuery: For Each Step
    RunQuery --> CollectResults
    CollectResults --> ExecuteSteps: More Steps
    CollectResults --> FormatOutput: All Done
    
    FormatOutput --> [*]
```

## CloudShip Integration

```mermaid
flowchart TD
    subgraph "Ship CLI"
        TOOL[Tool Execution]
        RESULT[Tool Results]
        PUSH[Push Handler]
    end
    
    subgraph "Preparation"
        B64[Base64 Encode]
        META[Add Metadata]
        TAGS[Add Tags]
    end
    
    subgraph "CloudShip API"
        AUTH_CHECK{Authenticated?}
        API[/v1/artifacts]
        RESPONSE[Analysis URL]
    end
    
    TOOL --> RESULT
    RESULT --> PUSH
    PUSH --> B64
    B64 --> META
    META --> TAGS
    TAGS --> AUTH_CHECK
    
    AUTH_CHECK -->|Yes| API
    AUTH_CHECK -->|No| ERROR[Auth Error]
    
    API --> RESPONSE
    RESPONSE --> USER[Display to User]
```

## Dagger Container Orchestration

```mermaid
graph LR
    subgraph "Host System"
        CLI[Ship CLI]
        CREDS[AWS Credentials]
        TF_FILES[Terraform Files]
    end
    
    subgraph "Dagger Engine"
        CLIENT[Dagger Client]
        GRAPH[Execution Graph]
    end
    
    subgraph "Containers"
        subgraph "Steampipe Container"
            SP_BASE[turbot/steampipe]
            SP_PLUGIN[AWS Plugin]
            SP_ENV[Environment Vars]
        end
        
        subgraph "Terraform Container"
            TF_BASE[hashicorp/terraform]
            TF_TOOLS[tflint/tfdocs/etc]
        end
        
        subgraph "Cost Analysis"
            COST_BASE[infracost/infracost]
            COST_API[API Keys]
        end
    end
    
    CLI --> CLIENT
    CLIENT --> GRAPH
    GRAPH --> SP_BASE
    GRAPH --> TF_BASE
    GRAPH --> COST_BASE
    
    CREDS -.->|Parse & Pass| SP_ENV
    TF_FILES -.->|Mount| TF_TOOLS
```

## Adding New Modules

### Module Creation Workflow

```mermaid
flowchart TD
    START[Start] --> CREATE[Create Module File]
    CREATE --> IMPLEMENT[Implement Interface]
    
    subgraph "Module Implementation"
        STRUCT[Define Module Struct]
        NEW[NewModule Function]
        METHODS[Implement Methods]
    end
    
    IMPLEMENT --> STRUCT
    STRUCT --> NEW
    NEW --> METHODS
    
    METHODS --> CONTAINER[Configure Container]
    
    subgraph "Container Setup"
        IMAGE[Select Base Image]
        DEPS[Install Dependencies]
        ENV_SETUP[Set Environment]
        EXEC_LOGIC[Define Execution]
    end
    
    CONTAINER --> IMAGE
    IMAGE --> DEPS
    DEPS --> ENV_SETUP
    ENV_SETUP --> EXEC_LOGIC
    
    EXEC_LOGIC --> INTEGRATE[Integrate with CLI]
    
    subgraph "CLI Integration"
        CMD_FILE[Create Command File]
        FLAGS[Define Flags]
        HANDLER[Write Handler]
        REGISTER[Register Command]
    end
    
    INTEGRATE --> CMD_FILE
    CMD_FILE --> FLAGS
    FLAGS --> HANDLER
    HANDLER --> REGISTER
    
    REGISTER --> TEST[Test Module]
    TEST --> DOCUMENT[Document Usage]
    DOCUMENT --> END[End]
```

### Module Interface

```mermaid
classDiagram
    class Module {
        <<interface>>
        +Initialize(ctx, client)
        +Execute(ctx, params)
        +GetOutput() string
    }
    
    class CustomModule {
        -client DaggerClient
        -config ModuleConfig
        +NewCustomModule() CustomModule
        +Initialize(ctx, client)
        +Execute(ctx, params)
        +GetOutput() string
        -setupContainer()
        -runCommand()
        -parseResults()
    }
    
    class ModuleConfig {
        +BaseImage string
        +Commands []string
        +Environment map[string]string
        +Mounts []Mount
        +OutputFormat string
    }
    
    class DaggerClient {
        +Container() Container
        +Host() Host
    }
    
    class Container {
        +From(image) Container
        +WithExec(cmd) Container
        +WithEnvVariable(key, val) Container
        +WithDirectory(path, dir) Container
        +Stdout() string
        +Stderr() string
    }
    
    Module <|.. CustomModule
    CustomModule --> ModuleConfig
    CustomModule --> DaggerClient
    DaggerClient --> Container
```

## Credential Flow

```mermaid
flowchart TD
    subgraph "Credential Sources"
        ENV[Environment Variables]
        FILE[~/.aws/credentials]
        PROFILE[AWS Profile]
    end
    
    subgraph "Ship CLI Processing"
        CHECK{Credentials Set?}
        PARSE[Parse Credentials File]
        EXTRACT[Extract Keys]
        MAP[Create Credential Map]
    end
    
    subgraph "Container Injection"
        ENV_VARS[Set Env Variables]
        MOUNT[Mount Files]
        CONFIG[Configure SDK]
    end
    
    subgraph "Tool Usage"
        STEAMPIPE[Steampipe AWS]
        TERRAFORM[Terraform AWS]
        AWS_CLI[AWS CLI]
    end
    
    ENV --> CHECK
    FILE --> CHECK
    PROFILE --> CHECK
    
    CHECK -->|No Env Vars| PARSE
    CHECK -->|Has Env Vars| MAP
    
    PARSE --> EXTRACT
    EXTRACT --> MAP
    MAP --> ENV_VARS
    
    ENV_VARS --> CONFIG
    CONFIG --> STEAMPIPE
    CONFIG --> TERRAFORM
    CONFIG --> AWS_CLI
```

## Error Handling Flow

```mermaid
stateDiagram-v2
    [*] --> ExecuteCommand
    
    ExecuteCommand --> CheckDocker
    CheckDocker --> DockerError: Not Running
    CheckDocker --> InitDagger: Running
    
    DockerError --> [*]: Exit with Error
    
    InitDagger --> DaggerError: Init Failed
    InitDagger --> RunModule: Success
    
    DaggerError --> [*]: Exit with Error
    
    RunModule --> ParseOutput
    ParseOutput --> Success: Valid Output
    ParseOutput --> ParseError: Invalid Output
    
    Success --> CheckPushFlag
    ParseError --> FallbackOutput
    
    FallbackOutput --> CheckPushFlag
    
    CheckPushFlag --> Push: --push Set
    CheckPushFlag --> Display: No Push
    
    Push --> AuthError: Not Authenticated
    Push --> UploadSuccess: Authenticated
    
    AuthError --> Display: Show Auth Message
    UploadSuccess --> Display: Show URL
    
    Display --> [*]
```

## MCP Server Architecture

```mermaid
graph TB
    subgraph "MCP Clients"
        CLAUDE[Claude Desktop]
        CURSOR[Cursor IDE]
        OTHER[Other MCP Clients]
    end
    
    subgraph "Ship MCP Server"
        SERVER[MCP Server Process]
        ROUTER[Tool Router]
        
        subgraph "Available Tools"
            QUERY_TOOL[steampipe_query]
            LINT_TOOL[terraform_lint]
            DOCS_TOOL[terraform_docs]
            SCAN_TOOL[security_scan]
            COST_TOOL[cost_analysis]
            DIAGRAM_TOOL[inframap_diagram]
        end
    end
    
    subgraph "Execution Layer"
        DAGGER_ENG[Dagger Engine]
        MODULES[Tool Modules]
    end
    
    CLAUDE --> SERVER
    CURSOR --> SERVER
    OTHER --> SERVER
    
    SERVER --> ROUTER
    ROUTER --> QUERY_TOOL
    ROUTER --> LINT_TOOL
    ROUTER --> DOCS_TOOL
    ROUTER --> SCAN_TOOL
    ROUTER --> COST_TOOL
    ROUTER --> DIAGRAM_TOOL
    
    QUERY_TOOL --> DAGGER_ENG
    LINT_TOOL --> DAGGER_ENG
    DOCS_TOOL --> DAGGER_ENG
    SCAN_TOOL --> DAGGER_ENG
    COST_TOOL --> DAGGER_ENG
    DIAGRAM_TOOL --> DAGGER_ENG
    
    DAGGER_ENG --> MODULES
```

## Community Module Integration

```mermaid
flowchart LR
    subgraph "Community"
        DEV[Module Developer]
        REPO[GitHub Repository]
        REGISTRY[Module Registry]
    end
    
    subgraph "Module Structure"
        CODE[module.go]
        CONFIG[module.yaml]
        DOCS[README.md]
        TESTS[tests/]
    end
    
    subgraph "Ship CLI"
        DISCOVER[Module Discovery]
        INSTALL[Module Install]
        LOAD[Module Loader]
        EXEC[Module Executor]
    end
    
    subgraph "User"
        USER_CLI[ship modules install]
        USER_RUN[ship run-module]
    end
    
    DEV --> CODE
    DEV --> CONFIG
    DEV --> DOCS
    DEV --> TESTS
    
    CODE --> REPO
    CONFIG --> REPO
    DOCS --> REPO
    TESTS --> REPO
    
    REPO --> REGISTRY
    REGISTRY --> DISCOVER
    
    USER_CLI --> DISCOVER
    DISCOVER --> INSTALL
    INSTALL --> LOAD
    
    USER_RUN --> LOAD
    LOAD --> EXEC
```

---

This comprehensive set of diagrams covers:
- Overall architecture and component relationships
- User workflows for different features
- Module system and extensibility
- AI investigation process
- CloudShip integration details
- Container orchestration with Dagger
- How to add new modules
- Credential management flow
- Error handling patterns
- MCP server architecture
- Community module integration

Each diagram provides a different perspective on how Ship CLI works, making it easier for developers to understand and extend the system.