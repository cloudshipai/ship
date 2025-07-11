package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudshipai/ship/internal/agent/tools"
	"dagger.io/dagger"
)

// EinoInvestigationAgent implements the InvestigationAgent interface using Eino framework
type EinoInvestigationAgent struct {
	client       *dagger.Client
	memory       *AgentMemory
	learner      SchemaLearner
	steampipeTool *tools.SteampipeTool
	llmModel     model.ChatModel
	agent        *react.Agent
	memoryPath   string
}

// NewEinoInvestigationAgent creates a new Eino-based investigation agent
func NewEinoInvestigationAgent(ctx context.Context, client *dagger.Client, apiKey, memoryPath string) (*EinoInvestigationAgent, error) {
	// Initialize agent memory
	memory := &AgentMemory{
		Schemas:   make(map[string]TableSchema),
		Patterns:  make(map[string]QueryPattern),
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}

	// Create schema learner
	learner := NewMemorySchemaLearner(client, memory)

	// Create Steampipe tool
	toolMemory := &tools.AgentMemory{
		Successes: make([]tools.QuerySuccess, 0),
		Failures:  make([]tools.QueryFailure, 0),
	}
	steampipeTool := tools.NewSteampipeTool(client, toolMemory)

	// Configure OpenAI model
	config := &openai.ChatModelConfig{
		APIKey: apiKey,
		Model:  "gpt-4",
	}

	llmModel, err := openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI model: %w", err)
	}

	// Create agent tools
	agentTools := []tool.BaseTool{steampipeTool}

	// Create ReAct agent
	reactAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		Model: llmModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: agentTools,
		},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			systemMessage := schema.SystemMessage(`You are an expert cloud infrastructure investigator. Your goal is to help users understand their cloud infrastructure through natural language queries.

CAPABILITIES:
- Execute SQL queries against AWS, Azure, and GCP using Steampipe
- Analyze security configurations and compliance
- Identify cost optimization opportunities  
- Investigate performance and operational issues
- Provide actionable insights and recommendations

APPROACH:
1. Parse the user's natural language request to understand their intent
2. Use the steampipe_query tool to gather relevant infrastructure data
3. Execute multiple targeted queries if needed to get comprehensive information
4. Analyze the results and provide clear, actionable insights
5. Include specific recommendations when security or cost issues are found

QUERY GUIDELINES:
- Always use exact column names from table schemas
- For AWS EC2 instances, use 'instance_state' not 'state' 
- Be specific with WHERE clauses to avoid large result sets
- When unsure about schema, query information_schema first
- Prefer targeted queries over broad SELECT * statements

RESPONSE FORMAT:
- Start with a brief summary of what you found
- List specific findings with numbers/counts
- Highlight any security or cost concerns
- Provide actionable recommendations
- Be concise but thorough in your analysis`)
			
			// Prepend system message to input
			result := make([]*schema.Message, 0, len(input)+1)
			result = append(result, systemMessage)
			result = append(result, input...)
			return result
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ReAct agent: %w", err)
	}

	agent := &EinoInvestigationAgent{
		client:        client,
		memory:        memory,
		learner:       learner,
		steampipeTool: steampipeTool,
		llmModel:      llmModel,
		agent:         reactAgent,
		memoryPath:    memoryPath,
	}

	// Load existing memory if available
	if err := agent.LoadMemory(); err != nil {
		slog.Debug("No existing memory found, starting fresh", "error", err)
	}

	return agent, nil
}

// Investigate performs an infrastructure investigation using the Eino agent
func (a *EinoInvestigationAgent) Investigate(ctx context.Context, request InvestigationRequest) (*InvestigationResult, error) {
	start := time.Now()
	
	slog.Info("Starting investigation", "prompt", request.Prompt, "provider", request.Provider)

	// Learn schemas for the provider if not already cached
	if err := a.ensureSchemaLearned(ctx, request.Provider, request.Credentials); err != nil {
		slog.Warn("Failed to learn schemas", "error", err)
		// Continue anyway, agent might still work with basic knowledge
	}

	// Create investigation context
	ctxWithCreds := context.WithValue(ctx, "credentials", request.Credentials)
	ctxWithProvider := context.WithValue(ctxWithCreds, "provider", request.Provider)
	ctxWithStart := context.WithValue(ctxWithProvider, "start_time", start)

	// Enhance the user prompt with schema information
	enhancedPrompt := a.enhancePromptWithContext(request)

	// Execute the agent
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: enhancedPrompt,
		},
	}

	// Generate the agent response
	response, err := a.agent.Generate(ctxWithStart, messages)
	if err != nil {
		return nil, fmt.Errorf("agent execution failed: %w", err)
	}

	// Parse the agent's response
	result := a.parseAgentResponse(response, request, start)

	// Save updated memory
	if err := a.SaveMemory(); err != nil {
		slog.Warn("Failed to save agent memory", "error", err)
	}

	slog.Info("Investigation completed", "duration", time.Since(start), "steps", len(result.Steps))
	return result, nil
}

// ensureSchemaLearned ensures that schemas are learned for the provider
func (a *EinoInvestigationAgent) ensureSchemaLearned(ctx context.Context, provider string, credentials map[string]string) error {
	// For now, always try to learn schemas
	// TODO: Implement proper schema caching

	// Learn schemas
	slog.Info("Learning schemas for provider", "provider", provider)
	return a.learner.LearnSchema(ctx, provider, credentials)
}

// enhancePromptWithContext adds relevant context to the user prompt
func (a *EinoInvestigationAgent) enhancePromptWithContext(request InvestigationRequest) string {
	var enhanced strings.Builder
	
	enhanced.WriteString(fmt.Sprintf("INVESTIGATION REQUEST:\n%s\n\n", request.Prompt))
	enhanced.WriteString(fmt.Sprintf("TARGET PROVIDER: %s\n", request.Provider))
	
	if request.Region != "" {
		enhanced.WriteString(fmt.Sprintf("REGION: %s\n", request.Region))
	}

	// Add schema context for relevant tables
	relevantTables := a.identifyRelevantTables(request.Prompt, request.Provider)
	if len(relevantTables) > 0 {
		enhanced.WriteString(fmt.Sprintf("\nRelevant tables: %s\n", strings.Join(relevantTables, ", ")))
	}

	// Add lessons learned from previous failures
	if len(a.memory.Failures) > 0 {
		enhanced.WriteString("\nKNOWN ISSUES TO AVOID:\n")
		recentFailures := a.getRecentFailures(5) // Last 5 failures
		for _, failure := range recentFailures {
			if failure.LessonLearned != "" {
				enhanced.WriteString(fmt.Sprintf("- %s\n", failure.LessonLearned))
			}
		}
	}

	enhanced.WriteString("\nPlease investigate this request step by step using the steampipe_query tool.")
	
	return enhanced.String()
}

// identifyRelevantTables identifies which tables are likely relevant to the query
func (a *EinoInvestigationAgent) identifyRelevantTables(prompt, provider string) []string {
	prompt = strings.ToLower(prompt)
	var relevant []string

	// Common table mappings based on keywords
	tableKeywords := map[string][]string{
		"aws_ec2_instance":        {"ec2", "instance", "server", "compute", "vm"},
		"aws_s3_bucket":          {"s3", "bucket", "storage", "object"},
		"aws_rds_db_instance":    {"rds", "database", "db", "mysql", "postgres"},
		"aws_vpc_security_group": {"security group", "sg", "firewall", "rules", "ports"},
		"aws_iam_user":           {"iam", "user", "identity", "access"},
		"aws_iam_role":           {"role", "permission", "policy"},
		"aws_vpc":                {"vpc", "network", "subnet"},
		"aws_lambda_function":    {"lambda", "function", "serverless"},
	}

	for table, keywords := range tableKeywords {
		for _, keyword := range keywords {
			if strings.Contains(prompt, keyword) {
				relevant = append(relevant, table)
				break
			}
		}
	}

	// Always include account table for basic info
	if provider == "aws" {
		relevant = append(relevant, "aws_account")
	}

	return relevant
}

// getRecentFailures returns the most recent failures for learning
func (a *EinoInvestigationAgent) getRecentFailures(count int) []QueryFailure {
	failures := a.memory.Failures
	if len(failures) <= count {
		return failures
	}
	return failures[len(failures)-count:]
}

// parseAgentResponse converts the agent's response into an InvestigationResult
func (a *EinoInvestigationAgent) parseAgentResponse(response *schema.Message, request InvestigationRequest, start time.Time) *InvestigationResult {
	result := &InvestigationResult{
		Success:    true,
		Steps:      []InvestigationStep{},
		Summary:    response.Content,
		Insights:   []Insight{},
		QueryCount: 0,
		Duration:   time.Since(start).String(),
		Confidence: 0.8, // Default confidence
	}

	// Extract insights from the response
	result.Insights = a.extractInsights(response.Content, request.Provider)

	// Count queries from memory
	result.QueryCount = len(a.memory.Successes) + len(a.memory.Failures)

	return result
}

// extractInsights parses insights from the agent's response
func (a *EinoInvestigationAgent) extractInsights(content, provider string) []Insight {
	var insights []Insight

	// Look for common security patterns in the response
	if strings.Contains(strings.ToLower(content), "security group") && strings.Contains(content, "0.0.0.0/0") {
		insights = append(insights, Insight{
			Type:           "security",
			Severity:       "high",
			Title:          "Open Security Groups Detected",
			Description:    "Found security groups allowing access from 0.0.0.0/0",
			Impact:         "Potential unauthorized access to resources",
			Recommendation: "Review and restrict security group rules to specific IP ranges",
			Confidence:     0.9,
		})
	}

	if strings.Contains(strings.ToLower(content), "unencrypted") {
		insights = append(insights, Insight{
			Type:           "security",
			Severity:       "medium",
			Title:          "Unencrypted Resources Found",
			Description:    "Some resources are not encrypted",
			Impact:         "Data may be vulnerable if compromised",
			Recommendation: "Enable encryption for all sensitive data stores",
			Confidence:     0.8,
		})
	}

	// Look for cost optimization opportunities
	if strings.Contains(strings.ToLower(content), "stopped") && strings.Contains(content, "instance") {
		insights = append(insights, Insight{
			Type:           "cost",
			Severity:       "medium",
			Title:          "Stopped Instances Found",
			Description:    "Found stopped EC2 instances that may be incurring costs",
			Impact:         "Unnecessary costs for unused resources",
			Recommendation: "Consider terminating unused instances or using scheduled start/stop",
			Confidence:     0.7,
		})
	}

	return insights
}

// GetMemory returns the agent's current memory
func (a *EinoInvestigationAgent) GetMemory() *AgentMemory {
	return a.memory
}

// SaveMemory persists the agent's memory to disk
func (a *EinoInvestigationAgent) SaveMemory() error {
	if a.memoryPath == "" {
		return nil // Memory persistence disabled
	}

	// Ensure directory exists
	dir := filepath.Dir(a.memoryPath)
	if err := ensureDir(dir); err != nil {
		return fmt.Errorf("failed to create memory directory: %w", err)
	}

	// Update timestamp
	a.memory.LastUpdate = time.Now().Format(time.RFC3339)

	// Marshal to JSON
	data, err := json.MarshalIndent(a.memory, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal memory: %w", err)
	}

	// Write to file
	if err := writeFile(a.memoryPath, data); err != nil {
		return fmt.Errorf("failed to write memory file: %w", err)
	}

	slog.Debug("Agent memory saved", "path", a.memoryPath)
	return nil
}

// LoadMemory loads the agent's memory from disk
func (a *EinoInvestigationAgent) LoadMemory() error {
	if a.memoryPath == "" {
		return nil // Memory persistence disabled
	}

	data, err := readFile(a.memoryPath)
	if err != nil {
		return fmt.Errorf("failed to read memory file: %w", err)
	}

	if err := json.Unmarshal(data, a.memory); err != nil {
		return fmt.Errorf("failed to unmarshal memory: %w", err)
	}

	slog.Debug("Agent memory loaded", "path", a.memoryPath, "schemas", len(a.memory.Schemas))
	return nil
}

// File I/O functions for agent memory persistence
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}