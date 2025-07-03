package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

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
	provider := "aws" // Default
	if len(providers) > 0 {
		provider = providers[0]
	}

	// Get available Steampipe tables for context
	availableTables := GetCommonSteampipeTables(provider)
	tableList := ""
	for _, table := range availableTables {
		tableList += "- " + table + "\n"
	}
	
	// Get example queries with REAL column names
	examples := GetSteampipeTableExamples(provider)
	exampleQueries := ""
	for desc, query := range examples {
		exampleQueries += fmt.Sprintf("- %s: %s\n", desc, query)
	}
	
	// Column information will be provided dynamically in the objective
	columnInfo := ""

	promptText := fmt.Sprintf(`
You are a cloud infrastructure investigator. Create a step-by-step investigation plan for:

Objective: %s
Cloud Provider: %s

IMPORTANT: You MUST use these actual Steampipe tables for %s:
%s

Example queries with CORRECT column names:
%s

%s

Return a JSON array of investigation steps, each with:
- "step_number": sequential number  
- "description": what this step investigates
- "provider": which cloud provider to query ("%s")
- "query": the Steampipe SQL query to run (MUST use actual tables and column names from examples above)
- "expected_insights": what we hope to learn

IMPORTANT STEAMPIPE QUERY RULES:
1. Use ONLY the exact column names provided in the schema information above
2. Boolean columns compare with true/false not 'true'/'false' strings  
3. For JSONB columns, use appropriate operators (->>, ->, jsonb_array_elements)
4. Table names always start with provider prefix (aws_, azure_, gcp_)
5. If the objective mentions specific column information, trust it completely
6. ONLY ONE SQL STATEMENT PER QUERY - no semicolons except at the end
7. For complex queries, break them into multiple steps

SPECIAL INSTRUCTIONS FOR TERRAFORM COMPARISONS:
If the user asks about comparing Terraform with cloud resources:
- Focus on querying actual cloud resources via Steampipe
- Query resource tags to identify Terraform-managed resources (look for tags like 'terraform', 'managed-by', etc.)
- Check for common Terraform resource types: EC2 instances, S3 buckets, RDS databases, IAM roles
- Note: Direct Terraform state comparison requires the ai-agent command with cross-module access

Focus on security, compliance, cost, and performance aspects.
Generate exactly 3-5 steps that thoroughly investigate the objective.
NEVER use made-up table names or column names. Only use what's provided above.
`, objective, provider, provider, tableList, exampleQueries, columnInfo, provider)

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

	// Parse the LLM response
	responseText, err := synced.LastReply(ctx)
	if err != nil {
		// Fallback with a real table if LLM fails
		slog.Warn("LLM failed, using fallback", "error", err)
		
		// Use a real table based on provider
		fallbackTable := "aws_account"
		if provider == "azure" {
			fallbackTable = "azure_subscription"
		} else if provider == "gcp" {
			fallbackTable = "gcp_project"
		}
		
		steps := []InvestigationStep{
			{
				StepNumber:       1,
				Description:      "Basic " + provider + " account overview",
				Provider:         provider,
				Query:            "SELECT * FROM " + fallbackTable + " LIMIT 1",
				ExpectedInsights: "Account information and configuration",
			},
		}
		return steps, nil
	}
	
	previewLen := 200
	if len(responseText) < previewLen {
		previewLen = len(responseText)
	}
	slog.Debug("LLM raw response received", "response_length", len(responseText), "response_preview", responseText[:previewLen])
	
	// Clean up the response - remove markdown code blocks if present
	cleanedResponse := responseText
	
	// First, check if there's a JSON array somewhere in the response
	// Look for the first '[' that starts a JSON array
	arrayStart := strings.Index(responseText, "[")
	if arrayStart >= 0 {
		// Find the matching closing bracket
		bracketCount := 0
		arrayEnd := -1
		for i := arrayStart; i < len(responseText); i++ {
			if responseText[i] == '[' {
				bracketCount++
			} else if responseText[i] == ']' {
				bracketCount--
				if bracketCount == 0 {
					arrayEnd = i + 1
					break
				}
			}
		}
		
		if arrayEnd > arrayStart {
			// Extract just the JSON array portion
			cleanedResponse = strings.TrimSpace(responseText[arrayStart:arrayEnd])
			slog.Debug("Extracted JSON array from response", "json_length", len(cleanedResponse))
		}
	} else if strings.Contains(responseText, "```json") {
		// Find the start after ```json
		start := strings.Index(responseText, "```json") + 7
		// Find the end before the last ```
		end := strings.LastIndex(responseText, "```")
		if start > 0 && end > start {
			// Extract the JSON portion and trim whitespace
			cleanedResponse = strings.TrimSpace(responseText[start:end])
		}
	} else if strings.Contains(responseText, "```") {
		// Handle plain ``` blocks
		start := strings.Index(responseText, "```") + 3
		end := strings.LastIndex(responseText, "```")
		if start > 0 && end > start {
			cleanedResponse = strings.TrimSpace(responseText[start:end])
		}
	}
	
	// Try to parse the JSON response
	var steps []InvestigationStep
	if err := json.Unmarshal([]byte(cleanedResponse), &steps); err != nil {
		slog.Debug("Failed to parse LLM JSON response", "error", err, "cleaned_response", cleanedResponse)
		
		// Use template-based fallback based on the objective
		objectiveLower := strings.ToLower(objective)
		if strings.Contains(objectiveLower, "ec2") || strings.Contains(objectiveLower, "instance") {
			// Use real EC2 queries
			if strings.Contains(objectiveLower, "running") {
				steps = []InvestigationStep{
					{
						StepNumber:       1,
						Description:      "Count running EC2 instances",
						Provider:         provider,
						Query:            "SELECT COUNT(*) as count FROM aws_ec2_instance WHERE instance_state = 'running'",
						ExpectedInsights: "Number of running instances",
					},
					{
						StepNumber:       2,
						Description:      "List running EC2 instances with details",
						Provider:         provider,
						Query:            "SELECT instance_id, instance_type, instance_state, region, vpc_id FROM aws_ec2_instance WHERE instance_state = 'running'",
						ExpectedInsights: "Details of running instances",
					},
				}
			} else {
				steps = []InvestigationStep{
					{
						StepNumber:       1,
						Description:      "List all EC2 instances",
						Provider:         provider,
						Query:            "SELECT instance_id, instance_type, instance_state, region, vpc_id FROM aws_ec2_instance",
						ExpectedInsights: "Overview of all EC2 instances",
					},
				}
			}
		} else if strings.Contains(objectiveLower, "s3") || strings.Contains(objectiveLower, "bucket") {
			steps = []InvestigationStep{
				{
					StepNumber:       1,
					Description:      "Count S3 buckets",
					Provider:         provider,
					Query:            "SELECT COUNT(*) as count FROM aws_s3_bucket",
					ExpectedInsights: "Total number of S3 buckets",
				},
				{
					StepNumber:       2,
					Description:      "List S3 buckets with details",
					Provider:         provider,
					Query:            "SELECT name, region, creation_date FROM aws_s3_bucket",
					ExpectedInsights: "S3 bucket inventory",
				},
			}
		} else {
			// Generic fallback
			fallbackTable := "aws_account"
			if provider == "azure" {
				fallbackTable = "azure_subscription"  
			} else if provider == "gcp" {
				fallbackTable = "gcp_project"
			}
			
			steps = []InvestigationStep{
				{
					StepNumber:       1,
					Description:      "Investigation based on: " + objective,
					Provider:         provider,
					Query:            "SELECT * FROM " + fallbackTable + " LIMIT 1",
					ExpectedInsights: "Basic account overview",
				},
			}
		}
		
		slog.Warn("Using fallback investigation plan due to JSON parsing failure")
	} else {
		// Successfully parsed steps
		slog.Info("Successfully parsed investigation steps", "step_count", len(steps))
		for i, step := range steps {
			slog.Debug("Parsed step", "number", i+1, "query", step.Query)
		}
	}

	return steps, nil
}
