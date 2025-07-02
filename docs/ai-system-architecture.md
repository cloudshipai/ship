# AI System Architecture - Ship CLI

This document provides detailed diagrams of the AI components and their interactions within Ship CLI.

## AI Components Overview

```mermaid
graph TB
    subgraph "AI Commands"
        AI_INVESTIGATE[ai-investigate]
        AI_AGENT[ai-agent]
        AI_SERVICES[ai-services]
    end
    
    subgraph "LLM Integration Layer"
        LLM_MODULE[LLMModule]
        LLM_DAGGER[DaggerLLMModule]
        LLM_WITH_TOOLS[LLMWithToolsModule]
        LLM_SERVICE[LLMServiceOrchestrator]
    end
    
    subgraph "Context Providers"
        STEAMPIPE_TABLES[SteampipeTables]
        TOOL_REGISTRY[ToolRegistry]
        PROMPT_TEMPLATES[PromptTemplates]
    end
    
    subgraph "LLM Providers"
        OPENAI[OpenAI GPT-4]
        ANTHROPIC[Claude]
        OLLAMA[Local Ollama]
    end
    
    subgraph "Execution"
        PLAN_EXECUTOR[PlanExecutor]
        TOOL_EXECUTOR[ToolExecutor]
        QUERY_ENGINE[QueryEngine]
    end
    
    AI_INVESTIGATE --> LLM_DAGGER
    AI_AGENT --> LLM_WITH_TOOLS
    AI_SERVICES --> LLM_SERVICE
    
    LLM_DAGGER --> STEAMPIPE_TABLES
    LLM_WITH_TOOLS --> TOOL_REGISTRY
    LLM_SERVICE --> PROMPT_TEMPLATES
    
    LLM_MODULE --> OPENAI
    LLM_MODULE --> ANTHROPIC
    LLM_MODULE --> OLLAMA
    
    LLM_DAGGER --> PLAN_EXECUTOR
    LLM_WITH_TOOLS --> TOOL_EXECUTOR
    PLAN_EXECUTOR --> QUERY_ENGINE
```

## AI Investigation Detailed Flow

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant LLM
    participant Tables
    participant Steampipe
    participant Results
    
    User->>CLI: "Find all S3 buckets with versioning"
    CLI->>Tables: GetCommonSteampipeTables("aws")
    Tables-->>CLI: [aws_s3_bucket, aws_ec2_instance, ...]
    
    CLI->>Tables: GetSteampipeTableExamples("aws")
    Tables-->>CLI: {"Check S3": "SELECT name FROM aws_s3_bucket", ...}
    
    CLI->>LLM: CreateInvestigationPlan(prompt, tables, examples)
    
    Note over LLM: LLM sees:<br/>- User prompt<br/>- Available tables<br/>- Example queries
    
    LLM-->>CLI: JSON Plan with Steps
    
    loop For each step
        CLI->>Steampipe: ExecuteQuery(step.query)
        Steampipe-->>CLI: Query Results
        CLI->>Results: Append(step, results)
    end
    
    CLI->>Results: FormatInsights()
    Results-->>User: Display findings
```

## AI Hallucination Prevention

```mermaid
flowchart TD
    subgraph "Before Fix"
        OLD_PROMPT[Generic Prompt]
        OLD_LLM[LLM Guesses]
        OLD_FALLBACK[Hardcoded Fake Table]
        OLD_ERROR[Query Fails]
        
        OLD_PROMPT --> OLD_LLM
        OLD_LLM --> OLD_FALLBACK
        OLD_FALLBACK --> OLD_ERROR
    end
    
    subgraph "After Fix"
        NEW_PROMPT[Prompt with Tables]
        TABLE_LIST[Real Table Names]
        EXAMPLES[Query Examples]
        NEW_LLM[LLM Uses Real Tables]
        VALID_QUERY[Valid Query]
        SUCCESS[Query Succeeds]
        
        TABLE_LIST --> NEW_PROMPT
        EXAMPLES --> NEW_PROMPT
        NEW_PROMPT --> NEW_LLM
        NEW_LLM --> VALID_QUERY
        VALID_QUERY --> SUCCESS
    end
    
    OLD_ERROR -.->|Fixed| SUCCESS
```

## AI Agent Tool Selection

```mermaid
stateDiagram-v2
    [*] --> ReceiveTask
    
    ReceiveTask --> AnalyzeTask: Parse Intent
    
    AnalyzeTask --> SecurityTask: Security Keywords
    AnalyzeTask --> CostTask: Cost Keywords
    AnalyzeTask --> DocTask: Documentation Keywords
    AnalyzeTask --> GeneralTask: Other
    
    SecurityTask --> SelectSecurityTools
    CostTask --> SelectCostTools
    DocTask --> SelectDocTools
    GeneralTask --> SelectGeneralTools
    
    SelectSecurityTools --> PlanExecution: [Checkov, Trivy]
    SelectCostTools --> PlanExecution: [Infracost, OpenInfraQuote]
    SelectDocTools --> PlanExecution: [terraform-docs]
    SelectGeneralTools --> PlanExecution: [Steampipe, TFLint]
    
    PlanExecution --> ExecuteTools
    ExecuteTools --> CollectResults
    CollectResults --> GenerateReport
    GenerateReport --> [*]
```

## AI Services Microservices Architecture

```mermaid
graph LR
    subgraph "AI Orchestrator"
        ORCHESTRATOR[Service Orchestrator]
        TASK_QUEUE[Task Queue]
        HTTP_CLIENT[HTTP Client]
    end
    
    subgraph "Service Containers"
        subgraph "Steampipe Service"
            SP_HTTP[:8001]
            SP_ENGINE[Query Engine]
        end
        
        subgraph "Cost Service"
            COST_HTTP[:8002]
            COST_ENGINE[Cost Calculator]
        end
        
        subgraph "Docs Service"
            DOCS_HTTP[:8003]
            DOCS_ENGINE[Doc Generator]
        end
        
        subgraph "Security Service"
            SEC_HTTP[:8004]
            SEC_ENGINE[Security Scanner]
        end
        
        subgraph "Diagram Service"
            DIAG_HTTP[:8005]
            DIAG_ENGINE[Diagram Builder]
        end
    end
    
    subgraph "Service Discovery"
        REGISTRY[Service Registry]
        HEALTH[Health Checks]
    end
    
    ORCHESTRATOR --> TASK_QUEUE
    TASK_QUEUE --> HTTP_CLIENT
    
    HTTP_CLIENT --> SP_HTTP
    HTTP_CLIENT --> COST_HTTP
    HTTP_CLIENT --> DOCS_HTTP
    HTTP_CLIENT --> SEC_HTTP
    HTTP_CLIENT --> DIAG_HTTP
    
    SP_HTTP --> SP_ENGINE
    COST_HTTP --> COST_ENGINE
    DOCS_HTTP --> DOCS_ENGINE
    SEC_HTTP --> SEC_ENGINE
    DIAG_HTTP --> DIAG_ENGINE
    
    REGISTRY --> ORCHESTRATOR
    HEALTH --> REGISTRY
```

## LLM Provider Selection

```mermaid
flowchart TD
    subgraph "Provider Selection"
        CHECK_ENV{Check Environment}
        CHECK_FLAG{Check CLI Flag}
        DEFAULT[Default Provider]
    end
    
    subgraph "OpenAI"
        OPENAI_KEY{OPENAI_API_KEY?}
        OPENAI_MODEL[gpt-4/gpt-3.5]
        OPENAI_CLIENT[OpenAI Client]
    end
    
    subgraph "Anthropic"
        ANTHROPIC_KEY{ANTHROPIC_API_KEY?}
        ANTHROPIC_MODEL[claude-3]
        ANTHROPIC_CLIENT[Anthropic Client]
    end
    
    subgraph "Ollama"
        OLLAMA_CHECK{Ollama Running?}
        OLLAMA_MODEL[llama2/mistral]
        OLLAMA_CLIENT[Ollama Client]
    end
    
    CHECK_FLAG --> CHECK_ENV
    CHECK_ENV --> OPENAI_KEY
    CHECK_ENV --> ANTHROPIC_KEY
    CHECK_ENV --> OLLAMA_CHECK
    
    OPENAI_KEY -->|Yes| OPENAI_MODEL
    OPENAI_KEY -->|No| DEFAULT
    
    ANTHROPIC_KEY -->|Yes| ANTHROPIC_MODEL
    ANTHROPIC_KEY -->|No| DEFAULT
    
    OLLAMA_CHECK -->|Yes| OLLAMA_MODEL
    OLLAMA_CHECK -->|No| DEFAULT
    
    OPENAI_MODEL --> OPENAI_CLIENT
    ANTHROPIC_MODEL --> ANTHROPIC_CLIENT
    OLLAMA_MODEL --> OLLAMA_CLIENT
    
    DEFAULT --> OPENAI_CLIENT
```

## Prompt Engineering Architecture

```mermaid
classDiagram
    class PromptBuilder {
        +BuildInvestigationPrompt()
        +BuildAgentPrompt()
        +BuildToolPrompt()
        -addSystemContext()
        -addTableContext()
        -addExamples()
        -formatJSON()
    }
    
    class TableContext {
        +GetTables(provider)
        +GetExamples(provider)
        +GetSchemas(table)
    }
    
    class SystemPrompts {
        +InvestigatorRole
        +AgentRole
        +ToolUserRole
        +JSONFormatter
    }
    
    class PromptTemplate {
        +Template string
        +Variables map
        +Render() string
    }
    
    class ValidationRules {
        +ValidateTables()
        +ValidateQueries()
        +ValidateJSON()
    }
    
    PromptBuilder --> TableContext
    PromptBuilder --> SystemPrompts
    PromptBuilder --> PromptTemplate
    PromptBuilder --> ValidationRules
    
    TableContext --> SteampipeTables
```

## AI Error Recovery

```mermaid
stateDiagram-v2
    [*] --> ExecuteAIPlan
    
    ExecuteAIPlan --> ParseJSON
    ParseJSON --> JSONValid: Success
    ParseJSON --> JSONInvalid: Failed
    
    JSONInvalid --> ExtractFromMarkdown: Try Markdown
    ExtractFromMarkdown --> JSONValid: Found JSON
    ExtractFromMarkdown --> UseFallback: Still Failed
    
    JSONValid --> ExecuteSteps
    UseFallback --> ExecuteSteps: Default Query
    
    ExecuteSteps --> QueryError: Query Failed
    ExecuteSteps --> QuerySuccess: Query OK
    
    QueryError --> CheckTimeout: Exit Code 41?
    QueryError --> CheckSyntax: SQL Error?
    QueryError --> CheckAuth: Auth Error?
    
    CheckTimeout --> RetryWithTimeout: Increase Timeout
    CheckSyntax --> FixQuery: Validate Tables
    CheckAuth --> ReloadCreds: Refresh Credentials
    
    RetryWithTimeout --> QuerySuccess
    FixQuery --> QuerySuccess
    ReloadCreds --> QuerySuccess
    
    QuerySuccess --> NextStep
    NextStep --> ExecuteSteps: More Steps
    NextStep --> Complete: Done
    
    Complete --> [*]
```

## Tool Registry and Discovery

```mermaid
graph TB
    subgraph "Tool Registry"
        REGISTRY[(Tool Database)]
        
        subgraph "Tool Metadata"
            NAME[Tool Name]
            DESC[Description]
            CONTAINER[Container Image]
            COMMANDS[Commands]
            INPUT[Input Schema]
            OUTPUT[Output Schema]
        end
    end
    
    subgraph "Tool Categories"
        SECURITY[Security Tools]
        COST[Cost Tools]
        DOCS[Documentation]
        INFRA[Infrastructure]
        QUALITY[Code Quality]
    end
    
    subgraph "Discovery"
        LIST[List Available]
        SEARCH[Search by Tag]
        MATCH[Match by Task]
    end
    
    subgraph "Registration"
        BUILTIN[Built-in Tools]
        COMMUNITY[Community Tools]
        CUSTOM[Custom Tools]
    end
    
    REGISTRY --> NAME
    REGISTRY --> DESC
    REGISTRY --> CONTAINER
    REGISTRY --> COMMANDS
    REGISTRY --> INPUT
    REGISTRY --> OUTPUT
    
    SECURITY --> REGISTRY
    COST --> REGISTRY
    DOCS --> REGISTRY
    INFRA --> REGISTRY
    QUALITY --> REGISTRY
    
    LIST --> REGISTRY
    SEARCH --> REGISTRY
    MATCH --> REGISTRY
    
    BUILTIN --> REGISTRY
    COMMUNITY --> REGISTRY
    CUSTOM --> REGISTRY
```

## AI Context Management

```mermaid
flowchart LR
    subgraph "Context Sources"
        USER_PROMPT[User Prompt]
        CLOUD_PROVIDER[Cloud Provider]
        PREV_RESULTS[Previous Results]
        ERROR_CONTEXT[Error Messages]
    end
    
    subgraph "Context Builder"
        COLLECTOR[Context Collector]
        FILTER[Relevance Filter]
        LIMITER[Token Limiter]
        FORMATTER[Context Formatter]
    end
    
    subgraph "Context Types"
        TABLES_CTX[Table Context]
        EXAMPLES_CTX[Example Context]
        ERROR_CTX[Error Context]
        HISTORY_CTX[History Context]
    end
    
    subgraph "LLM Input"
        SYSTEM_MSG[System Message]
        USER_MSG[User Message]
        CONTEXT_MSG[Context Message]
    end
    
    USER_PROMPT --> COLLECTOR
    CLOUD_PROVIDER --> COLLECTOR
    PREV_RESULTS --> COLLECTOR
    ERROR_CONTEXT --> COLLECTOR
    
    COLLECTOR --> FILTER
    FILTER --> LIMITER
    LIMITER --> FORMATTER
    
    FORMATTER --> TABLES_CTX
    FORMATTER --> EXAMPLES_CTX
    FORMATTER --> ERROR_CTX
    FORMATTER --> HISTORY_CTX
    
    TABLES_CTX --> CONTEXT_MSG
    EXAMPLES_CTX --> CONTEXT_MSG
    ERROR_CTX --> CONTEXT_MSG
    HISTORY_CTX --> CONTEXT_MSG
    
    SYSTEM_MSG --> LLM
    USER_MSG --> LLM
    CONTEXT_MSG --> LLM
```

---

These diagrams provide detailed insights into:
- How the AI components are structured
- The flow of AI-powered investigations
- How hallucination prevention works
- Tool selection logic for AI agents
- Microservices architecture for AI services
- LLM provider selection and fallbacks
- Prompt engineering architecture
- Error recovery mechanisms
- Tool registry and discovery
- Context management for AI interactions

Together with the main architecture diagrams, this provides a comprehensive visual guide to understanding and extending Ship CLI's AI capabilities.