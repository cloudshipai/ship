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

	"dagger.io/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// EinoInvestigationAgent implements the InvestigationAgent interface using Eino framework
type EinoInvestigationAgent struct {
	client     *dagger.Client
	memory     *AgentMemory
	agent      *react.Agent
	memoryPath string
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

	// Load existing memory if it exists
	if memoryPath != "" {
		if err := os.MkdirAll(filepath.Dir(memoryPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create memory directory: %w", err)
		}

		if data, err := os.ReadFile(memoryPath); err == nil {
			if err := json.Unmarshal(data, memory); err != nil {
				slog.Warn("Failed to load agent memory", "error", err)
			} else {
				slog.Info("Loaded agent memory", "schemas", len(memory.Schemas), "patterns", len(memory.Patterns))
			}
		}
	}

	// Create OpenAI ChatModel
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:  "gpt-4",
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI chat model: %w", err)
	}

	// Create Steampipe tool
	steampipeTool := createSteampipeTool(client)

	// Create React Agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{steampipeTool},
		},
		MaxStep: 10,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create React agent: %w", err)
	}

	return &EinoInvestigationAgent{
		client:     client,
		memory:     memory,
		agent:      agent,
		memoryPath: memoryPath,
	}, nil
}

// SteampipeRequest represents a Steampipe query request
type SteampipeRequest struct {
	Query    string `json:"query"`
	Provider string `json:"provider"`
}

// SteampipeResponse represents a Steampipe query response
type SteampipeResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// createSteampipeTool creates a Steampipe query execution tool for the agent
func createSteampipeTool(client *dagger.Client) tool.BaseTool {
	toolInfo := &schema.ToolInfo{
		Name: "steampipe_query",
		Desc: "Execute Steampipe SQL queries to investigate cloud infrastructure. Use this to query AWS, Azure, or GCP resources using SQL syntax.",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Type:     schema.String,
				Desc:     "The SQL query to execute against Steampipe tables",
				Required: true,
			},
			"provider": {
				Type: schema.String,
				Desc: "Cloud provider (aws, azure, gcp)",
			},
		}),
	}

	return utils.NewTool(toolInfo, func(ctx context.Context, req *SteampipeRequest) (*SteampipeResponse, error) {
		// Default to AWS if no provider specified
		if req.Provider == "" {
			req.Provider = "aws"
		}

		// Execute Steampipe query
		steampipeModule := modules.NewSteampipeModule(client)

		// Get provider credentials from environment
		credentials := getProviderCredentialsFromEnv(req.Provider)

		result, err := steampipeModule.RunQuery(ctx, req.Provider, req.Query, credentials, "json")
		if err != nil {
			return &SteampipeResponse{
				Error: fmt.Sprintf("steampipe query failed: %v", err),
			}, nil // Return error in response, not as error
		}

		return &SteampipeResponse{
			Result: result,
		}, nil
	})
}

// getProviderCredentialsFromEnv gets cloud provider credentials from environment variables
func getProviderCredentialsFromEnv(provider string) map[string]string {
	credentials := make(map[string]string)

	switch provider {
	case "aws":
		if val := os.Getenv("AWS_ACCESS_KEY_ID"); val != "" {
			credentials["AWS_ACCESS_KEY_ID"] = val
		}
		if val := os.Getenv("AWS_SECRET_ACCESS_KEY"); val != "" {
			credentials["AWS_SECRET_ACCESS_KEY"] = val
		}
		if val := os.Getenv("AWS_SESSION_TOKEN"); val != "" {
			credentials["AWS_SESSION_TOKEN"] = val
		}
		if val := os.Getenv("AWS_REGION"); val != "" {
			credentials["AWS_REGION"] = val
		} else {
			credentials["AWS_REGION"] = "us-east-1" // default
		}
		if val := os.Getenv("AWS_PROFILE"); val != "" {
			credentials["AWS_PROFILE"] = val
		}

	case "azure":
		if val := os.Getenv("AZURE_CLIENT_ID"); val != "" {
			credentials["AZURE_CLIENT_ID"] = val
		}
		if val := os.Getenv("AZURE_CLIENT_SECRET"); val != "" {
			credentials["AZURE_CLIENT_SECRET"] = val
		}
		if val := os.Getenv("AZURE_TENANT_ID"); val != "" {
			credentials["AZURE_TENANT_ID"] = val
		}
		if val := os.Getenv("AZURE_SUBSCRIPTION_ID"); val != "" {
			credentials["AZURE_SUBSCRIPTION_ID"] = val
		}

	case "gcp":
		if val := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); val != "" {
			credentials["GOOGLE_APPLICATION_CREDENTIALS"] = val
		}
		if val := os.Getenv("GOOGLE_CLOUD_PROJECT"); val != "" {
			credentials["GOOGLE_CLOUD_PROJECT"] = val
		}
		if val := os.Getenv("GCLOUD_PROJECT"); val != "" {
			credentials["GCLOUD_PROJECT"] = val
		}
	}

	return credentials
}

// Investigate performs an AI-powered infrastructure investigation using the Eino React agent
func (a *EinoInvestigationAgent) Investigate(ctx context.Context, request InvestigationRequest) (*InvestigationResult, error) {
	startTime := time.Now()

	slog.Info("Starting investigation", "prompt", request.Prompt, "provider", request.Provider)

	// Create enhanced prompt with context
	enhancedPrompt := a.enhancePrompt(request.Prompt, request.Provider, request.Region)

	// Create user message
	messages := []*schema.Message{
		schema.SystemMessage(fmt.Sprintf("You are an expert cloud infrastructure analyst specializing in %s. Use the steampipe_query tool to investigate cloud resources and provide detailed insights. Always explain your findings and provide actionable recommendations.", request.Provider)),
		schema.UserMessage(enhancedPrompt),
	}

	// Execute investigation using React agent
	response, err := a.agent.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("agent execution failed: %w", err)
	}

	// Build result
	result := &InvestigationResult{
		Success:    true,
		Summary:    response.Content,
		Duration:   time.Since(startTime).String(),
		Confidence: 0.95, // High confidence for Eino-based queries
		QueryCount: 1,    // Will be updated based on tool calls
		Steps:      []InvestigationStep{},
		Insights:   []Insight{},
		Timestamp:  startTime,
	}

	// Extract insights from the response
	result.Insights = a.extractInsights(response.Content, request.Provider)

	// Update memory with successful investigation
	a.updateMemory(request.Prompt, result)

	// Save memory if path is provided
	if a.memoryPath != "" {
		if err := a.SaveMemory(); err != nil {
			slog.Warn("Failed to save agent memory", "error", err)
		}
	}

	return result, nil
}

// enhancePrompt creates a better prompt with context and examples
func (a *EinoInvestigationAgent) enhancePrompt(prompt, provider, region string) string {
	enhanced := fmt.Sprintf("Cloud Provider: %s", provider)
	if region != "" {
		enhanced += fmt.Sprintf("\nRegion: %s", region)
	}
	enhanced += fmt.Sprintf("\n\nTask: %s", prompt)

	enhanced += `

Please execute Steampipe queries to investigate this request. Provide:
1. Clear findings from the data
2. Security or cost implications
3. Specific actionable recommendations
4. Any compliance concerns

Use appropriate Steampipe tables for ` + provider + ` such as:
- aws_ec2_instance (for EC2 instances)
- aws_s3_bucket (for S3 buckets)  
- aws_iam_user (for IAM users)
- aws_vpc_security_group (for security groups)
- aws_rds_db_instance (for RDS instances)

Format your response with clear sections and specific findings.`

	return enhanced
}

// extractInsights extracts actionable insights from the investigation response
func (a *EinoInvestigationAgent) extractInsights(content, provider string) []Insight {
	insights := []Insight{}

	contentLower := strings.ToLower(content)

	// Security insights for public access (first priority)
	if strings.Contains(contentLower, "0.0.0.0/0") || strings.Contains(contentLower, "public") {
		insights = append(insights, Insight{
			Type:           "security",
			Title:          "Public Access Detected",
			Description:    "Found resources with public access that may pose security risks",
			Severity:       "high",
			Recommendation: "Review and restrict public access to essential services only",
		})
	}

	// Security insights for unencrypted resources (second priority)
	if strings.Contains(contentLower, "unencrypted") || strings.Contains(contentLower, "no encryption") {
		insights = append(insights, Insight{
			Type:           "security",
			Title:          "Encryption Issue",
			Description:    "Found resources without proper encryption",
			Severity:       "high",
			Recommendation: "Enable encryption for sensitive data and storage",
		})
	}

	// Cost insights (third priority)
	if strings.Contains(contentLower, "unused") || strings.Contains(contentLower, "idle") || strings.Contains(contentLower, "stopped") || strings.Contains(contentLower, "cost") {
		insights = append(insights, Insight{
			Type:           "cost",
			Title:          "Cost Optimization Opportunity",
			Description:    "Found unused or idle resources that may be costing money",
			Severity:       "medium",
			Recommendation: "Consider terminating or rightsizing unused resources",
		})
	}

	// Compliance insights (fourth priority)
	if strings.Contains(contentLower, "compliance") || strings.Contains(contentLower, "regulation") {
		insights = append(insights, Insight{
			Type:           "compliance",
			Title:          "Compliance Issue",
			Description:    "Found potential compliance concerns",
			Severity:       "high",
			Recommendation: "Review compliance requirements and implement necessary controls",
		})
	}

	return insights
}

// updateMemory updates the agent memory with successful investigation patterns
func (a *EinoInvestigationAgent) updateMemory(prompt string, result *InvestigationResult) {
	// Record successful query pattern
	pattern := QueryPattern{
		Prompt:     prompt,
		Provider:   "aws", // Default for now
		Success:    result.Success,
		Confidence: result.Confidence,
		Timestamp:  time.Now(),
	}

	// Store pattern with a key based on prompt keywords
	key := fmt.Sprintf("pattern_%d", len(a.memory.Patterns))
	a.memory.Patterns[key] = pattern

	// Record success
	success := QuerySuccess{
		Prompt:     prompt,
		QueryCount: result.QueryCount,
		Duration:   result.Duration,
		Confidence: result.Confidence,
		Timestamp:  time.Now(),
	}
	a.memory.Successes = append(a.memory.Successes, success)

	// Update last update time
	a.memory.LastUpdate = time.Now().Format(time.RFC3339)
}

// SaveMemory saves the agent memory to disk
func (a *EinoInvestigationAgent) SaveMemory() error {
	if a.memoryPath == "" {
		return nil
	}

	data, err := json.MarshalIndent(a.memory, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal memory: %w", err)
	}

	return os.WriteFile(a.memoryPath, data, 0644)
}

// LoadMemory loads the agent memory from disk
func (a *EinoInvestigationAgent) LoadMemory() error {
	if a.memoryPath == "" {
		return nil
	}

	data, err := os.ReadFile(a.memoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, initialize with empty memory
			a.memory = &AgentMemory{
				Schemas:   make(map[string]TableSchema),
				Patterns:  make(map[string]QueryPattern),
				Successes: make([]QuerySuccess, 0),
				Failures:  make([]QueryFailure, 0),
			}
			return nil
		}
		return fmt.Errorf("failed to read memory file: %w", err)
	}

	err = json.Unmarshal(data, a.memory)
	if err != nil {
		return fmt.Errorf("failed to unmarshal memory: %w", err)
	}

	return nil
}

// GetMemory returns the current agent memory
func (a *EinoInvestigationAgent) GetMemory() *AgentMemory {
	return a.memory
}

// identifyRelevantTables identifies relevant Steampipe tables based on prompt and provider
func (a *EinoInvestigationAgent) identifyRelevantTables(prompt, provider string) []string {
	promptLower := strings.ToLower(prompt)
	tables := []string{}

	// Always include account table for provider context
	switch provider {
	case "aws":
		tables = append(tables, "aws_account")

		// EC2 related tables
		if strings.Contains(promptLower, "instance") || strings.Contains(promptLower, "ec2") || strings.Contains(promptLower, "server") || strings.Contains(promptLower, "compute") {
			tables = append(tables, "aws_ec2_instance")
		}

		// Security group related
		if strings.Contains(promptLower, "security") || strings.Contains(promptLower, "firewall") || strings.Contains(promptLower, "0.0.0.0") || strings.Contains(promptLower, "port") {
			tables = append(tables, "aws_vpc_security_group")
		}

		// S3 related tables
		if strings.Contains(promptLower, "bucket") || strings.Contains(promptLower, "s3") || strings.Contains(promptLower, "storage") {
			tables = append(tables, "aws_s3_bucket")
		}

		// IAM related tables
		if strings.Contains(promptLower, "user") || strings.Contains(promptLower, "role") || strings.Contains(promptLower, "iam") || strings.Contains(promptLower, "permission") || strings.Contains(promptLower, "policy") {
			tables = append(tables, "aws_iam_user")
			tables = append(tables, "aws_iam_role")
		}

		// Lambda related tables
		if strings.Contains(promptLower, "lambda") || strings.Contains(promptLower, "function") || strings.Contains(promptLower, "serverless") {
			tables = append(tables, "aws_lambda_function")
		}

		// RDS related tables
		if strings.Contains(promptLower, "database") || strings.Contains(promptLower, "rds") || strings.Contains(promptLower, "mysql") || strings.Contains(promptLower, "postgres") {
			tables = append(tables, "aws_rds_db_instance")
		}

		// VPC related tables
		if strings.Contains(promptLower, "vpc") || strings.Contains(promptLower, "network") || strings.Contains(promptLower, "subnet") {
			tables = append(tables, "aws_vpc")
		}

		// EBS related tables
		if strings.Contains(promptLower, "volume") || strings.Contains(promptLower, "ebs") || strings.Contains(promptLower, "disk") {
			tables = append(tables, "aws_ebs_volume")
		}

	case "azure":
		tables = append(tables, "azure_subscription")

		// Add Azure-specific table identification logic
		if strings.Contains(promptLower, "vm") || strings.Contains(promptLower, "virtual machine") || strings.Contains(promptLower, "compute") {
			tables = append(tables, "azure_compute_virtual_machine")
		}

		if strings.Contains(promptLower, "storage") || strings.Contains(promptLower, "blob") {
			tables = append(tables, "azure_storage_account")
		}

	case "gcp":
		tables = append(tables, "gcp_project")

		// Add GCP-specific table identification logic
		if strings.Contains(promptLower, "instance") || strings.Contains(promptLower, "compute") || strings.Contains(promptLower, "gce") {
			tables = append(tables, "gcp_compute_instance")
		}

		if strings.Contains(promptLower, "storage") || strings.Contains(promptLower, "bucket") || strings.Contains(promptLower, "gcs") {
			tables = append(tables, "gcp_storage_bucket")
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	result := []string{}
	for _, table := range tables {
		if !seen[table] {
			seen[table] = true
			result = append(result, table)
		}
	}

	return result
}

// enhancePromptWithContext enhances the prompt with context from memory and failures
func (a *EinoInvestigationAgent) enhancePromptWithContext(request InvestigationRequest) string {
	enhanced := fmt.Sprintf("TARGET PROVIDER: %s\n", request.Provider)

	if request.Region != "" {
		enhanced += fmt.Sprintf("REGION: %s\n", request.Region)
	}

	enhanced += fmt.Sprintf("USER REQUEST: %s\n\n", request.Prompt)

	// Add context from memory failures
	if len(a.memory.Failures) > 0 {
		enhanced += "KNOWN ISSUES TO AVOID:\n"
		for _, failure := range a.memory.Failures {
			enhanced += fmt.Sprintf("- %s (Lesson: %s)\n", failure.ErrorType, failure.LessonLearned)
		}
		enhanced += "\n"
	}

	// Add suggested approach
	enhanced += "Please investigate using appropriate Steampipe tables and provide:\n"
	enhanced += "1. Clear findings from the data\n"
	enhanced += "2. Security and cost implications\n"
	enhanced += "3. Specific actionable recommendations\n"
	enhanced += "4. Any compliance concerns\n"

	return enhanced
}
