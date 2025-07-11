package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
)

// AgentMemory represents simplified agent memory to avoid import cycles
type AgentMemory struct {
	Successes []QuerySuccess
	Failures  []QueryFailure
}

// QuerySuccess represents a successful query execution
type QuerySuccess struct {
	OriginalIntent string `json:"original_intent"`
	GeneratedQuery string `json:"generated_query"`
	ResultCount    int    `json:"result_count"`
	ExecutionTime  string `json:"execution_time"`
	Provider       string `json:"provider"`
	Timestamp      string `json:"timestamp"`
	PatternUsed    string `json:"pattern_used,omitempty"`
}

// QueryFailure represents a failed query with learning information
type QueryFailure struct {
	OriginalIntent string `json:"original_intent"`
	GeneratedQuery string `json:"generated_query"`
	ErrorMessage   string `json:"error_message"`
	ErrorType      string `json:"error_type"`
	Provider       string `json:"provider"`
	Timestamp      string `json:"timestamp"`
	LessonLearned  string `json:"lesson_learned"`
}


// SteampipeTool implements the schema.Tool interface for Steampipe query execution
type SteampipeTool struct {
	client   *dagger.Client
	module   *modules.SteampipeModule
	memory   *AgentMemory
	toolInfo *schema.ToolInfo
}

// NewSteampipeTool creates a new Steampipe tool for the Eino agent
func NewSteampipeTool(client *dagger.Client, memory *AgentMemory) *SteampipeTool {
	// Create tool parameters schema
	params := map[string]*schema.ParameterInfo{
		"provider": {
			Type:     "string",
			Desc:     "Cloud provider (aws, azure, gcp)",
			Required: true,
			Enum:     []string{"aws", "azure", "gcp"},
		},
		"query": {
			Type:     "string",
			Desc:     "SQL query to execute against Steampipe tables",
			Required: true,
		},
		"credentials": {
			Type:     "object",
			Desc:     "Cloud provider credentials (handled automatically)",
			Required: false,
		},
	}

	toolInfo := &schema.ToolInfo{
		Name: "steampipe_query",
		Desc: `Execute SQL queries against cloud infrastructure using Steampipe.
This tool can query AWS, Azure, and GCP resources to gather information about:
- EC2 instances, S3 buckets, RDS databases
- Security groups, IAM users and roles
- Cost and usage information
- Compliance and security findings
- Network configurations and VPCs

The tool automatically handles authentication and provides structured results.`,
		ParamsOneOf: schema.NewParamsOneOfByParams(params),
	}

	return &SteampipeTool{
		client:   client,
		module:   modules.NewSteampipeModule(client),
		memory:   memory,
		toolInfo: toolInfo,
	}
}

// Info returns the tool information for Eino (implements tool.BaseTool)
func (t *SteampipeTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return t.toolInfo, nil
}

// Name returns the tool name for compatibility
func (t *SteampipeTool) Name() string {
	return t.toolInfo.Name
}

// Description returns the tool description for compatibility
func (t *SteampipeTool) Description() string {
	return t.toolInfo.Desc
}

// SteampipeRequest represents the input to the Steampipe tool
type SteampipeRequest struct {
	Provider    string            `json:"provider"`
	Query       string            `json:"query"`
	Credentials map[string]string `json:"credentials"`
}

// SteampipeResponse represents the output from the Steampipe tool
type SteampipeResponse struct {
	Success       bool                     `json:"success"`
	Results       []map[string]interface{} `json:"results"`
	RowCount      int                      `json:"row_count"`
	ExecutionTime string                   `json:"execution_time"`
	Query         string                   `json:"query"`
	Error         string                   `json:"error,omitempty"`
	Insights      []string                 `json:"insights,omitempty"`
}

// InvokableRun executes the Steampipe tool with the given input (implements tool.InvokableTool)
func (t *SteampipeTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	slog.Debug("SteampipeTool.InvokableRun", "input", argumentsInJSON)

	// Parse the input JSON
	var req SteampipeRequest
	if err := json.Unmarshal([]byte(argumentsInJSON), &req); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	// Validate the request
	if req.Provider == "" {
		return "", fmt.Errorf("provider is required")
	}
	if req.Query == "" {
		return "", fmt.Errorf("query is required")
	}

	// Start timing (would use actual time tracking)
	_ = ctx.Value("start_time")

	// Prepare credentials
	credentials := req.Credentials
	if credentials == nil {
		credentials = t.getProviderCredentials(req.Provider)
	}

	// Validate and potentially improve the query
	improvedQuery, err := t.improveQuery(req.Query, req.Provider)
	if err != nil {
		slog.Warn("Query validation failed", "error", err, "original_query", req.Query)
		// Continue with original query, let Steampipe handle the error
		improvedQuery = req.Query
	}

	// Execute the query
	result, err := t.module.RunQuery(ctx, req.Provider, improvedQuery, credentials, "json")
	if err != nil {
		// Record the failure for learning
		t.recordFailure(req.Query, req.Provider, err.Error())
		
		response := SteampipeResponse{
			Success:       false,
			Query:         improvedQuery,
			Error:         err.Error(),
			ExecutionTime: "failed",
		}
		
		responseJSON, _ := json.Marshal(response)
		return string(responseJSON), nil // Return structured error instead of Go error
	}

	// Parse the results
	var results []map[string]interface{}
	if result != "" {
		if err := json.Unmarshal([]byte(result), &results); err != nil {
			slog.Debug("Failed to parse results as JSON array, treating as single result", "error", err)
			// Try to parse as single object or treat as raw string
			results = []map[string]interface{}{
				{"result": result},
			}
		}
	}

	// Record the success for learning
	t.recordSuccess(req.Query, improvedQuery, req.Provider, len(results))

	// Generate insights about the results
	insights := t.generateInsights(results, req.Provider, req.Query)

	response := SteampipeResponse{
		Success:       true,
		Results:       results,
		RowCount:      len(results),
		Query:         improvedQuery,
		ExecutionTime: "completed", // Would use actual timing
		Insights:      insights,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(responseJSON), nil
}

// improveQuery attempts to improve the query based on learned patterns
func (t *SteampipeTool) improveQuery(query, provider string) (string, error) {
	// Basic query validation
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("empty query")
	}

	// Ensure query doesn't have multiple statements
	if strings.Count(query, ";") > 1 {
		statements := strings.Split(query, ";")
		query = strings.TrimSpace(statements[0])
		slog.Info("Multiple statements detected, using only the first", "original_count", len(statements))
	}

	// Apply learned improvements based on provider
	improved := t.applyProviderSpecificImprovements(query, provider)
	
	return improved, nil
}

// applyProviderSpecificImprovements applies known fixes for common issues
func (t *SteampipeTool) applyProviderSpecificImprovements(query, provider string) string {
	if provider != "aws" {
		return query
	}

	// Apply common AWS-specific fixes based on learned patterns
	improvements := map[string]string{
		" state =":        " instance_state =",
		" state_name =":   " instance_state =", 
		"WHERE running":   "WHERE instance_state = 'running'",
		"WHERE stopped":   "WHERE instance_state = 'stopped'",
		"sg.group_id":     "sg->>'GroupId'",
		"sg.group_name":   "sg->>'GroupName'",
	}

	improved := query
	for old, new := range improvements {
		if strings.Contains(improved, old) {
			improved = strings.ReplaceAll(improved, old, new)
			slog.Debug("Applied query improvement", "old", old, "new", new)
		}
	}

	return improved
}

// recordSuccess records a successful query execution for learning
func (t *SteampipeTool) recordSuccess(originalQuery, executedQuery, provider string, resultCount int) {
	if t.memory == nil {
		return
	}

	success := QuerySuccess{
		OriginalIntent: originalQuery,
		GeneratedQuery: executedQuery,
		ResultCount:    resultCount,
		ExecutionTime:  "completed",
		Provider:       provider,
		Timestamp:      "now", // Would use actual timestamp
	}

	t.memory.Successes = append(t.memory.Successes, success)
	
	// Keep only recent successes (last 100)
	if len(t.memory.Successes) > 100 {
		t.memory.Successes = t.memory.Successes[len(t.memory.Successes)-100:]
	}
}

// recordFailure records a failed query execution for learning
func (t *SteampipeTool) recordFailure(query, provider, errorMsg string) {
	if t.memory == nil {
		return
	}

	// Determine error type
	errorType := "unknown"
	if strings.Contains(errorMsg, "column") && strings.Contains(errorMsg, "does not exist") {
		errorType = "schema"
	} else if strings.Contains(errorMsg, "syntax") {
		errorType = "syntax"
	} else if strings.Contains(errorMsg, "authentication") || strings.Contains(errorMsg, "access") {
		errorType = "auth"
	} else if strings.Contains(errorMsg, "timeout") {
		errorType = "timeout"
	}

	failure := QueryFailure{
		OriginalIntent: query,
		GeneratedQuery: query,
		ErrorMessage:   errorMsg,
		ErrorType:      errorType,
		Provider:       provider,
		Timestamp:      "now",
		LessonLearned:  t.generateLessonFromError(errorMsg),
	}

	t.memory.Failures = append(t.memory.Failures, failure)
	
	// Keep only recent failures (last 50)
	if len(t.memory.Failures) > 50 {
		t.memory.Failures = t.memory.Failures[len(t.memory.Failures)-50:]
	}
}

// generateLessonFromError creates a lesson learned from an error
func (t *SteampipeTool) generateLessonFromError(errorMsg string) string {
	if strings.Contains(errorMsg, `column "state"`) {
		return "Use 'instance_state' instead of 'state' for EC2 instance queries"
	}
	if strings.Contains(errorMsg, `column "running"`) {
		return "Use 'instance_state = \"running\"' instead of 'running' column"
	}
	if strings.Contains(errorMsg, "group_id") {
		return "Use JSONB operators for security group fields: sg->>'GroupId'"
	}
	
	return "Query failed - need to improve schema understanding"
}

// generateInsights creates insights from query results
func (t *SteampipeTool) generateInsights(results []map[string]interface{}, provider, query string) []string {
	var insights []string
	
	if len(results) == 0 {
		insights = append(insights, "No results found - consider broadening the query scope")
		return insights
	}

	// Add basic insights based on result patterns
	if len(results) == 1 {
		insights = append(insights, fmt.Sprintf("Found 1 result"))
	} else {
		insights = append(insights, fmt.Sprintf("Found %d results", len(results)))
	}

	// Provider-specific insights
	if provider == "aws" && strings.Contains(strings.ToLower(query), "ec2") {
		insights = append(insights, "Consider checking instance security groups and tags")
	}
	
	if provider == "aws" && strings.Contains(strings.ToLower(query), "s3") {
		insights = append(insights, "Consider checking bucket encryption and public access settings")
	}

	return insights
}

// getProviderCredentials gets credentials for the specified provider
func (t *SteampipeTool) getProviderCredentials(provider string) map[string]string {
	// This would integrate with the existing credential loading logic
	// For now, return empty map and let Dagger handle credential loading
	return make(map[string]string)
}