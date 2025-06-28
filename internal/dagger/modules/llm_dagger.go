package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// DaggerLLMModule uses Dagger's native LLM features
type DaggerLLMModule struct {
	client *dagger.Client
	model  string
}

// NewDaggerLLMModule creates a new module using Dagger's LLM primitives
func NewDaggerLLMModule(client *dagger.Client, model string) *DaggerLLMModule {
	if model == "" {
		model = "gpt-4" // Default model
	}
	return &DaggerLLMModule{
		client: client,
		model:  model,
	}
}

// AnalyzeSteampipeResults analyzes query results using Dagger LLM
func (m *DaggerLLMModule) AnalyzeSteampipeResults(ctx context.Context, queryResults string, queryContext string) (string, error) {
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

	// Use native Dagger LLM
	llm := m.client.LLM(dagger.LLMOpts{
		Model: m.model,
	})

	// Set up the prompt and get response
	llmWithPrompt := llm.WithPrompt(prompt)

	// Sync to execute the LLM
	synced, err := llmWithPrompt.Sync(ctx)
	if err != nil {
		return "", fmt.Errorf("LLM sync failed: %w", err)
	}

	// Get the last reply
	response, err := synced.LastReply(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get LLM response: %w", err)
	}

	return response, nil
}

// GenerateSteampipeQuery converts natural language to SQL using Dagger LLM
func (m *DaggerLLMModule) GenerateSteampipeQuery(ctx context.Context, naturalLanguageQuery string, provider string) (QueryPlan, error) {
	_ = ctx // Mark as used
	promptText := fmt.Sprintf(`
You are a Steampipe SQL expert. Convert this natural language query into a valid Steampipe SQL query for the %s provider.

Natural Language Query: %s

Generate a JSON response with:
{
  "query": "the SQL query",
  "description": "what this query does",
  "expected_columns": ["list", "of", "columns"],
  "filters_applied": ["list", "of", "filters"]
}

Requirements:
1. Use appropriate Steampipe tables for %s
2. Ensure the query is syntactically correct
3. Include appropriate columns for meaningful results
`, provider, naturalLanguageQuery, provider)

	// Use native Dagger LLM
	llm := m.client.LLM(dagger.LLMOpts{
		Model: m.model,
	})

	// Get SQL query from LLM
	llmWithPrompt := llm.WithPrompt(promptText)
	synced, err := llmWithPrompt.Sync(ctx)
	if err != nil {
		return QueryPlan{}, fmt.Errorf("LLM sync failed: %w", err)
	}

	response, err := synced.LastReply(ctx)
	if err != nil {
		return QueryPlan{}, fmt.Errorf("failed to get LLM response: %w", err)
	}

	// For now, return a simple plan with the response
	// In production, we'd parse the JSON response
	plan := QueryPlan{
		Query:           response,
		Description:     fmt.Sprintf("Query for: %s", naturalLanguageQuery),
		ExpectedColumns: []string{"*"},
		FiltersApplied:  []string{},
	}

	return plan, nil
}

// QueryPlan represents a generated query with metadata
type QueryPlan struct {
	Query           string   `json:"query"`
	Description     string   `json:"description"`
	ExpectedColumns []string `json:"expected_columns"`
	FiltersApplied  []string `json:"filters_applied"`
}

// CreateInvestigationPlan generates a comprehensive investigation plan using Dagger LLM
func (m *DaggerLLMModule) CreateInvestigationPlan(ctx context.Context, objective string, providers []string) ([]InvestigationStep, error) {
	_ = ctx           // Mark as used
	provider := "aws" // Default
	if len(providers) > 0 {
		provider = providers[0]
	}

	promptText := fmt.Sprintf(`
You are a cloud infrastructure investigator. Create a step-by-step investigation plan for:

Objective: %s
Cloud Provider: %s

Return a JSON array of investigation steps, each with:
- "step_number": sequential number  
- "description": what this step investigates
- "provider": which cloud provider to query ("%s")
- "query": the Steampipe SQL query to run
- "expected_insights": what we hope to learn

Focus on security, compliance, cost, and performance aspects.
Generate exactly 3-5 steps that thoroughly investigate the objective.
`, objective, provider, provider)

	// Use native Dagger LLM with system prompt for better JSON output
	llm := m.client.LLM(dagger.LLMOpts{
		Model: m.model,
	})

	// Add system prompt for JSON formatting
	llmWithPrompt := llm.
		WithSystemPrompt("You are a cloud infrastructure expert. Always respond with valid JSON arrays when asked for investigation plans.").
		WithPrompt(promptText)

	synced, err := llmWithPrompt.Sync(ctx)
	if err != nil {
		return nil, fmt.Errorf("LLM sync failed: %w", err)
	}

	_, err = synced.LastReply(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM response: %w", err)
	}

	// For now, return example steps
	// In production, we'd parse the JSON response from LLM
	// TODO: Parse response JSON into []InvestigationStep
	steps := []InvestigationStep{
		{
			StepNumber:       1,
			Description:      "Analyze resource security based on: " + objective,
			Provider:         provider,
			Query:            "SELECT * FROM " + provider + "_resource_inventory",
			ExpectedInsights: "Overview of all resources",
		},
	}

	return steps, nil
}
