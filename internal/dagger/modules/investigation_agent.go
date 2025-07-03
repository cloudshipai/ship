package modules

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// LLMInterface defines the interface for LLM modules
type LLMInterface interface {
	CreateInvestigationPlan(ctx context.Context, objective string, providers []string) ([]InvestigationStep, error)
	AnalyzeSteampipeResults(ctx context.Context, queryResults string, queryContext string) (string, error)
}

// InvestigationAgent performs investigations with retry and self-correction
type InvestigationAgent struct {
	llm             LLMInterface
	steampipe       *SteampipeModule
	maxRetries      int
	executedQueries []ExecutedQuery
}

// ExecutedQuery tracks query attempts and results
type ExecutedQuery struct {
	Query   string
	Error   error
	Result  string
	Attempt int
}

// NewInvestigationAgent creates a new investigation agent
func NewInvestigationAgent(llm LLMInterface, steampipe *SteampipeModule) *InvestigationAgent {
	return &InvestigationAgent{
		llm:        llm,
		steampipe:  steampipe,
		maxRetries: 3,
	}
}

// InvestigateWithRetry executes an investigation with error correction
func (a *InvestigationAgent) InvestigateWithRetry(ctx context.Context, objective string, provider string, credentials map[string]string) (*InvestigationAgentReport, error) {
	// First, try to get an investigation plan
	steps, err := a.llm.CreateInvestigationPlan(ctx, objective, []string{provider})
	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	var results []StepResult
	
	// Execute each step with retry logic
	for _, step := range steps {
		slog.Info("Executing step", "description", step.Description)
		
		result, err := a.executeStepWithRetry(ctx, step, provider, credentials)
		results = append(results, result)
		
		if err != nil {
			slog.Warn("Step failed after retries", "error", err)
		}
	}

	// Generate final analysis
	analysis, err := a.generateFinalAnalysis(ctx, objective, results)
	if err != nil {
		analysis = "Failed to generate analysis"
	}

	return &InvestigationAgentReport{
		Objective: objective,
		Steps:     steps,
		Results:   results,
		Analysis:  analysis,
	}, nil
}

// executeStepWithRetry executes a single step with retry logic
func (a *InvestigationAgent) executeStepWithRetry(ctx context.Context, step InvestigationStep, provider string, credentials map[string]string) (StepResult, error) {
	var lastError error
	query := step.Query
	
	for attempt := 1; attempt <= a.maxRetries; attempt++ {
		slog.Debug("Attempting query", "attempt", attempt, "query", query)
		
		// Try to execute the query
		result, err := a.steampipe.RunQuery(ctx, provider, query, credentials)
		
		if err == nil {
			// Success!
			return StepResult{
				StepNumber: step.StepNumber,
				Query:      query,
				Result:     result,
				Success:    true,
			}, nil
		}
		
		// Query failed - analyze the error
		lastError = err
		errorMsg := err.Error()
		
		// Log the attempt
		a.executedQueries = append(a.executedQueries, ExecutedQuery{
			Query:   query,
			Error:   err,
			Attempt: attempt,
		})
		
		// Don't retry on the last attempt
		if attempt >= a.maxRetries {
			break
		}
		
		// Try to fix the query based on the error
		fixedQuery, canFix := a.analyzeAndFixQuery(ctx, query, errorMsg, provider)
		if !canFix {
			slog.Warn("Cannot automatically fix query", "error", errorMsg)
			break
		}
		
		slog.Info("Retrying with fixed query", "original", query, "fixed", fixedQuery)
		query = fixedQuery
	}
	
	return StepResult{
		StepNumber: step.StepNumber,
		Query:      query,
		Error:      lastError,
		Success:    false,
	}, lastError
}

// analyzeAndFixQuery uses the LLM to fix a failed query
func (a *InvestigationAgent) analyzeAndFixQuery(ctx context.Context, query string, errorMsg string, provider string) (string, bool) {
	// First try simple pattern-based fixes
	if strings.Contains(errorMsg, "column") && strings.Contains(errorMsg, "does not exist") {
		// Extract the problematic column name
		if strings.Contains(errorMsg, `column "running"`) {
			return strings.ReplaceAll(query, "WHERE running", "WHERE instance_state = 'running'"), true
		}
		if strings.Contains(errorMsg, `column "state_name"`) {
			return strings.ReplaceAll(query, "state_name", "instance_state"), true
		}
		if strings.Contains(errorMsg, `column "state"`) && strings.Contains(query, "aws_ec2_instance") {
			return strings.ReplaceAll(query, " state", " instance_state"), true
		}
		if strings.Contains(errorMsg, "sg.group_id") {
			fixed := strings.ReplaceAll(query, "sg.group_id", "sg->>'GroupId'")
			fixed = strings.ReplaceAll(fixed, "sg.group_name", "sg->>'GroupName'")
			return fixed, true
		}
	}
	
	// If simple fixes don't work, ask the LLM to fix it
	fixPrompt := fmt.Sprintf(`
The following Steampipe query failed with an error. Please provide a corrected version.

Original Query:
%s

Error:
%s

Important:
- For AWS EC2, use 'instance_state' not 'state' or 'state_name'
- For running instances use: WHERE instance_state = 'running'
- For security groups in EC2, use jsonb_array_elements like: jsonb_array_elements(security_groups) as sg
- Then access fields like: sg->>'GroupId'

Provide ONLY the corrected SQL query, no explanation.
`, query, errorMsg)

	response, err := a.llm.AnalyzeSteampipeResults(ctx, fixPrompt, "fix_query")
	if err != nil {
		return "", false
	}
	
	// Clean up the response
	fixedQuery := strings.TrimSpace(response)
	fixedQuery = strings.TrimPrefix(fixedQuery, "```sql")
	fixedQuery = strings.TrimPrefix(fixedQuery, "```")
	fixedQuery = strings.TrimSuffix(fixedQuery, "```")
	fixedQuery = strings.TrimSpace(fixedQuery)
	
	return fixedQuery, true
}

// generateFinalAnalysis creates a summary of the investigation
func (a *InvestigationAgent) generateFinalAnalysis(ctx context.Context, objective string, results []StepResult) (string, error) {
	// Build a summary of successful results
	var successfulResults []string
	var failedSteps []string
	
	for _, result := range results {
		if result.Success {
			successfulResults = append(successfulResults, fmt.Sprintf("Step %d: %s", result.StepNumber, result.Result))
		} else {
			failedSteps = append(failedSteps, fmt.Sprintf("Step %d failed: %v", result.StepNumber, result.Error))
		}
	}
	
	analysisPrompt := fmt.Sprintf(`
Based on the investigation for "%s", here are the results:

Successful queries:
%s

Failed queries:
%s

Please provide a concise analysis of what we learned.
`, objective, strings.Join(successfulResults, "\n"), strings.Join(failedSteps, "\n"))

	return a.llm.AnalyzeSteampipeResults(ctx, strings.Join(successfulResults, "\n\n"), analysisPrompt)
}

// StepResult holds the result of executing a step
type StepResult struct {
	StepNumber int
	Query      string
	Result     string
	Error      error
	Success    bool
}

// InvestigationAgentReport with results
type InvestigationAgentReport struct {
	Objective string
	Steps     []InvestigationStep
	Results   []StepResult
	Analysis  string
}