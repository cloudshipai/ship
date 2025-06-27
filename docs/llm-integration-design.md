# LLM Integration Design for Ship CLI

## Overview

This document outlines how to integrate Dagger's native LLM capabilities with Ship CLI's existing tools, particularly focusing on Steampipe investigations and Terraform analysis tools.

## Architecture

### 1. Core LLM Module

Create a central LLM module that leverages Dagger's LLM primitive:

```go
// internal/dagger/modules/llm.go
type LLMModule struct {
    client *dagger.Client
    model  string // e.g., "gpt-4", "claude-3", "llama2"
}
```

### 2. Integration Points

#### A. Steampipe Intelligence Layer

**Purpose**: Natural language infrastructure queries and automated investigation workflows

**Features**:
1. **Natural Language to SQL Translation**
   - User: "Show me all EC2 instances that are running but not tagged with environment"
   - LLM: Generates Steampipe SQL query
   - System: Executes query and returns results

2. **Automated Investigation Chains**
   - User: "Investigate security issues in my AWS account"
   - LLM: Creates investigation workflow:
     ```
     1. Query security groups with open ports
     2. Check for unencrypted S3 buckets
     3. Find IAM users without MFA
     4. Generate remediation recommendations
     ```

3. **Contextual Analysis**
   - Feed query results back to LLM for analysis
   - Generate human-readable summaries and insights
   - Prioritize findings by severity

**Implementation Example**:
```go
func (m *LLMModule) InvestigateWithNaturalLanguage(ctx context.Context, prompt string) (*Investigation, error) {
    // 1. Use LLM to interpret the investigation request
    // 2. Generate Steampipe queries
    // 3. Execute queries via Steampipe module
    // 4. Analyze results with LLM
    // 5. Generate report with recommendations
}
```

#### B. Terraform Tools Intelligence

**Purpose**: Enhanced analysis of Terraform code with contextual understanding

**Features**:

1. **Cost Optimization Assistant**
   - Input: OpenInfraQuote cost analysis results
   - LLM Analysis:
     - Identifies cost optimization opportunities
     - Suggests instance rightsizing
     - Recommends reserved instances or savings plans
     - Provides terraform code modifications

2. **Security Remediation Guide**
   - Input: Checkov security scan results
   - LLM Analysis:
     - Explains security issues in plain language
     - Provides step-by-step remediation
     - Generates compliant Terraform code
     - Prioritizes fixes by risk level

3. **Documentation Enhancement**
   - Input: terraform-docs output
   - LLM Enhancement:
     - Adds contextual explanations
     - Creates usage examples
     - Generates architecture diagrams descriptions
     - Adds best practices section

**Implementation Example**:
```go
func (m *LLMModule) AnalyzeCosts(ctx context.Context, costData string) (*CostOptimization, error) {
    prompt := fmt.Sprintf(`
        Analyze this AWS infrastructure cost data and provide:
        1. Top 3 cost optimization opportunities
        2. Specific Terraform code changes
        3. Estimated monthly savings
        
        Cost Data: %s
    `, costData)
    
    return m.GenerateAnalysis(ctx, prompt)
}
```

### 3. Interactive Agent Mode

**Purpose**: Step-by-step infrastructure management with AI guidance

**Features**:
- Interactive Q&A about infrastructure
- Guided troubleshooting
- Automated fix generation and application

**Example Session**:
```
$ ship investigate --ai-assisted

AI: I'll help you investigate your infrastructure. What would you like to explore?

User: Check if any of my databases are publicly accessible

AI: I'll check for publicly accessible databases across your cloud providers.
    Executing queries...
    
    Found 2 potential issues:
    1. RDS instance "prod-db" has security group allowing 0.0.0.0/0 on port 5432
    2. DynamoDB table "user-data" has no encryption at rest
    
    Would you like me to generate Terraform code to fix these issues?
```

### 4. MCP Server Enhancement

Enhance the existing MCP server with LLM capabilities:

1. **Tool Discovery**: LLM can discover and use Ship CLI tools via MCP
2. **Context Sharing**: Share investigation results across LLM sessions
3. **Workflow Automation**: Create complex multi-step workflows

## Implementation Phases

### Phase 1: Core LLM Module (Week 1)
- [ ] Create LLM module with Dagger's LLM primitive
- [ ] Support for OpenAI, Anthropic, and local models
- [ ] Basic prompt templates for each tool

### Phase 2: Steampipe Integration (Week 2)
- [ ] Natural language to Steampipe SQL
- [ ] Investigation workflow generation
- [ ] Results analysis and summarization

### Phase 3: Terraform Tools Integration (Week 3)
- [ ] Cost optimization analysis
- [ ] Security remediation guides
- [ ] Documentation enhancement

### Phase 4: Interactive Agent (Week 4)
- [ ] Interactive CLI mode
- [ ] Context management
- [ ] Fix generation and application

## Configuration

```yaml
# ~/.ship/config.yaml
llm:
  provider: openai  # or anthropic, ollama
  model: gpt-4
  api_key: ${OPENAI_API_KEY}
  
  # Local model config (optional)
  local:
    provider: ollama
    model: llama2
    endpoint: http://localhost:11434
```

## Example Commands

```bash
# Natural language investigation
ship investigate --prompt "Find all unencrypted storage in my AWS account"

# Cost optimization with AI
ship terraform-tools cost-analysis . --ai-optimize

# Security fix generation
ship terraform-tools checkov-scan . --ai-remediate

# Interactive mode
ship ai-assistant
```

## Benefits

1. **Accessibility**: Non-technical users can perform complex investigations
2. **Automation**: AI can chain multiple tools for comprehensive analysis
3. **Intelligence**: Context-aware recommendations based on best practices
4. **Efficiency**: Faster identification and resolution of issues
5. **Learning**: AI can explain issues and fixes, educating users

## Security Considerations

1. **Data Privacy**: Ensure sensitive data is not sent to cloud LLM providers
2. **Code Review**: All AI-generated code should be reviewed before applying
3. **Access Control**: LLM should respect existing credential boundaries
4. **Audit Trail**: Log all AI-assisted actions for compliance

## Next Steps

1. Implement core LLM module
2. Create proof-of-concept for Steampipe natural language queries
3. Test with various LLM providers (OpenAI, Anthropic, local)
4. Gather feedback and iterate on prompts
5. Build interactive agent mode