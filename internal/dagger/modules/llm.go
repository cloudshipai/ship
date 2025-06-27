package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// LLMModule provides AI-powered analysis and generation capabilities
type LLMModule struct {
	client   *dagger.Client
	provider string
	model    string
}

// NewLLMModule creates a new LLM module
func NewLLMModule(client *dagger.Client, provider, model string) *LLMModule {
	return &LLMModule{
		client:   client,
		provider: provider,
		model:    model,
	}
}

// AnalyzeSteampipeResults analyzes Steampipe query results and provides insights
func (m *LLMModule) AnalyzeSteampipeResults(ctx context.Context, queryResults string, queryContext string) (string, error) {
	prompt := fmt.Sprintf(`
You are an infrastructure security and compliance expert. Analyze these Steampipe query results and provide:

1. A summary of key findings
2. Security or compliance issues identified
3. Prioritized recommendations for remediation
4. Any cost optimization opportunities

Query Context: %s

Query Results:
%s

Provide a clear, actionable analysis.
`, queryContext, queryResults)

	return m.runPrompt(ctx, prompt)
}

// GenerateSteampipeQuery converts natural language to Steampipe SQL
func (m *LLMModule) GenerateSteampipeQuery(ctx context.Context, naturalLanguageQuery string, provider string) (string, error) {
	prompt := fmt.Sprintf(`
You are a Steampipe SQL expert. Convert this natural language query into a valid Steampipe SQL query for the %s provider.

Natural Language Query: %s

Requirements:
1. Use appropriate Steampipe tables for %s
2. Return only the SQL query, no explanation
3. Ensure the query is syntactically correct
4. Include appropriate columns for meaningful results

SQL Query:
`, provider, naturalLanguageQuery, provider)

	return m.runPrompt(ctx, prompt)
}

// AnalyzeTerraformCosts provides cost optimization recommendations
func (m *LLMModule) AnalyzeTerraformCosts(ctx context.Context, costAnalysis string) (string, error) {
	prompt := fmt.Sprintf(`
You are a cloud cost optimization expert. Analyze this Terraform cost analysis and provide:

1. Top 3 cost optimization opportunities with estimated savings
2. Specific Terraform code changes to implement optimizations
3. Best practices for ongoing cost management
4. Any architectural changes that could reduce costs

Cost Analysis Data:
%s

Provide specific, actionable recommendations with code examples.
`, costAnalysis)

	return m.runPrompt(ctx, prompt)
}

// GenerateSecurityFixes creates Terraform code to fix security issues
func (m *LLMModule) GenerateSecurityFixes(ctx context.Context, securityScanResults string) (string, error) {
	prompt := fmt.Sprintf(`
You are a cloud security expert. Based on these security scan results, generate:

1. A prioritized list of security issues by severity
2. Terraform code fixes for each issue
3. Explanation of why each fix is important
4. Any additional security hardening recommendations

Security Scan Results:
%s

For each fix, provide the exact Terraform code that should be added or modified.
`, securityScanResults)

	return m.runPrompt(ctx, prompt)
}

// EnhanceDocumentation improves Terraform module documentation
func (m *LLMModule) EnhanceDocumentation(ctx context.Context, basicDocs string, moduleCode string) (string, error) {
	prompt := fmt.Sprintf(`
You are a technical documentation expert. Enhance this Terraform module documentation:

Current Documentation:
%s

Module Code Summary:
%s

Please add:
1. Clear usage examples with common scenarios
2. Best practices for using this module
3. Common pitfalls and how to avoid them
4. Integration examples with other modules
5. Troubleshooting guide

Keep the existing content and enhance it with the above sections.
`, basicDocs, moduleCode)

	return m.runPrompt(ctx, prompt)
}

// CreateInvestigationPlan generates a comprehensive investigation plan
func (m *LLMModule) CreateInvestigationPlan(ctx context.Context, objective string, availableProviders []string) ([]InvestigationStep, error) {
	providersStr := strings.Join(availableProviders, ", ")
	prompt := fmt.Sprintf(`
You are a cloud infrastructure investigator. Create a step-by-step investigation plan for:

Objective: %s
Available Providers: %s

Return a JSON array of investigation steps, each with:
- "step_number": sequential number
- "description": what this step investigates
- "provider": which cloud provider to query
- "query": the Steampipe SQL query to run
- "expected_insights": what we hope to learn

Focus on security, compliance, cost, and performance aspects.
`, objective, providersStr)

	result, err := m.runPrompt(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var steps []InvestigationStep
	if err := json.Unmarshal([]byte(result), &steps); err != nil {
		return nil, fmt.Errorf("failed to parse investigation plan: %w", err)
	}

	return steps, nil
}

// runPrompt executes a prompt using the configured LLM
func (m *LLMModule) runPrompt(ctx context.Context, prompt string) (string, error) {
	// This is a simplified implementation
	// In production, this would use Dagger's LLM primitive

	container := m.client.Container()

	switch m.provider {
	case "openai":
		// Use OpenAI API
		container = container.
			From("curlimages/curl:latest").
			WithEnvVariable("OPENAI_API_KEY", "${OPENAI_API_KEY}").
			WithExec([]string{
				"curl", "-s", "-X", "POST",
				"https://api.openai.com/v1/chat/completions",
				"-H", "Content-Type: application/json",
				"-H", "Authorization: Bearer ${OPENAI_API_KEY}",
				"-d", fmt.Sprintf(`{
					"model": "%s",
					"messages": [{"role": "user", "content": %q}],
					"temperature": 0.7
				}`, m.model, prompt),
			})

	case "anthropic":
		// Use Anthropic API
		container = container.
			From("curlimages/curl:latest").
			WithEnvVariable("ANTHROPIC_API_KEY", "${ANTHROPIC_API_KEY}").
			WithExec([]string{
				"curl", "-s", "-X", "POST",
				"https://api.anthropic.com/v1/messages",
				"-H", "Content-Type: application/json",
				"-H", "x-api-key: ${ANTHROPIC_API_KEY}",
				"-H", "anthropic-version: 2023-06-01",
				"-d", fmt.Sprintf(`{
					"model": "%s",
					"messages": [{"role": "user", "content": %q}],
					"max_tokens": 4096
				}`, m.model, prompt),
			})

	case "ollama":
		// Use local Ollama
		container = container.
			From("curlimages/curl:latest").
			WithExec([]string{
				"curl", "-s", "-X", "POST",
				"http://host.docker.internal:11434/api/generate",
				"-d", fmt.Sprintf(`{
					"model": "%s",
					"prompt": %q,
					"stream": false
				}`, m.model, prompt),
			})

	default:
		return "", fmt.Errorf("unsupported LLM provider: %s", m.provider)
	}

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run LLM prompt: %w", err)
	}

	// Parse response based on provider
	// This is simplified - real implementation would properly parse JSON responses
	return output, nil
}

// InvestigationStep represents a step in an investigation plan
type InvestigationStep struct {
	StepNumber       int    `json:"step_number"`
	Description      string `json:"description"`
	Provider         string `json:"provider"`
	Query            string `json:"query"`
	ExpectedInsights string `json:"expected_insights"`
}
