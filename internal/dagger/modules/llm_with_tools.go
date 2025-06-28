package modules

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger.io/dagger"
)

// LLMWithToolsModule extends LLM with ability to use other modules as tools
type LLMWithToolsModule struct {
	client          *dagger.Client
	model           string
	steampipeModule *SteampipeModule
	openInfraModule *OpenInfraQuoteModule
	terraformDocs   *TerraformDocsModule
}

// NewLLMWithToolsModule creates an LLM that can use other modules as tools
func NewLLMWithToolsModule(client *dagger.Client, model string) *LLMWithToolsModule {
	return &LLMWithToolsModule{
		client:          client,
		model:           model,
		steampipeModule: NewSteampipeModule(client),
		openInfraModule: NewOpenInfraQuoteModule(client),
		terraformDocs:   NewTerraformDocsModule(client),
	}
}

// InvestigateWithTools performs a complete investigation using all available tools
func (m *LLMWithToolsModule) InvestigateWithTools(ctx context.Context, objective string) (*InvestigationReport, error) {
	// Define available tools for the LLM
	toolsPrompt := `
You have access to the following tools:

1. STEAMPIPE_QUERY: Execute SQL queries against cloud infrastructure
   Usage: {"tool": "steampipe", "action": "query", "provider": "aws", "sql": "SELECT ..."}

2. COST_ANALYSIS: Analyze costs of Terraform plans  
   Usage: {"tool": "openinfraquote", "action": "analyze", "file": "path/to/tfplan.json"}

3. TERRAFORM_DOCS: Generate documentation for Terraform modules
   Usage: {"tool": "terraform-docs", "action": "generate", "path": "path/to/module"}

4. SECURITY_SCAN: Scan for security issues
   Usage: {"tool": "checkov", "action": "scan", "path": "path/to/code"}

To use a tool, respond with a JSON object containing the tool request.
After receiving results, you can use more tools or provide final analysis.
`

	// Initial prompt with objective and available tools
	systemPrompt := "You are an infrastructure investigator with access to real tools. " + toolsPrompt

	// Create LLM with tool-use system prompt
	llm := m.client.LLM(dagger.LLMOpts{
		Model: m.model,
	}).WithSystemPrompt(systemPrompt)

	// Start investigation
	conversation := llm.WithPrompt(fmt.Sprintf("Investigate: %s\n\nFirst, plan what tools you'll use.", objective))

	// Execute up to 5 tool uses
	var toolResults []ToolResult
	for i := 0; i < 5; i++ {
		// Get LLM response
		synced, err := conversation.Sync(ctx)
		if err != nil {
			return nil, fmt.Errorf("LLM sync failed: %w", err)
		}

		response, err := synced.LastReply(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get LLM response: %w", err)
		}

		// Check if response contains tool request
		var toolRequest LLMToolRequest
		if err := json.Unmarshal([]byte(response), &toolRequest); err == nil {
			// Execute the requested tool
			result, err := m.executeTool(ctx, toolRequest)
			if err != nil {
				// Tell LLM about the error
				conversation = conversation.WithPrompt(fmt.Sprintf("Tool error: %v", err))
				continue
			}

			toolResults = append(toolResults, result)

			// Feed results back to LLM
			conversation = conversation.WithPrompt(fmt.Sprintf(
				"Tool '%s' returned:\n%s\n\nWhat would you like to do next?",
				toolRequest.Tool, result.Output,
			))
		} else {
			// LLM provided final analysis
			return &InvestigationReport{
				Objective: objective,
				ToolsUsed: toolResults,
				Analysis:  response,
			}, nil
		}
	}

	// Ask for final summary
	finalConv := conversation.WithPrompt("Please provide a final summary of your investigation findings.")
	synced, err := finalConv.Sync(ctx)
	if err != nil {
		return nil, err
	}

	analysis, err := synced.LastReply(ctx)
	if err != nil {
		return nil, err
	}

	return &InvestigationReport{
		Objective: objective,
		ToolsUsed: toolResults,
		Analysis:  analysis,
	}, nil
}

// executeTool runs the requested tool and returns results
func (m *LLMWithToolsModule) executeTool(ctx context.Context, request LLMToolRequest) (ToolResult, error) {
	switch request.Tool {
	case "steampipe":
		// Execute Steampipe query
		provider := request.Params["provider"]
		sql := request.Params["sql"]

		output, err := m.steampipeModule.RunQuery(ctx, provider, sql, nil)
		return ToolResult{
			Tool:   "steampipe",
			Action: request.Action,
			Output: output,
			Error:  err,
		}, err

	case "openinfraquote":
		// Run cost analysis
		file := request.Params["file"]
		region := request.Params["region"]
		if region == "" {
			region = "us-east-1"
		}

		output, err := m.openInfraModule.AnalyzePlan(ctx, file, region)
		return ToolResult{
			Tool:   "openinfraquote",
			Action: request.Action,
			Output: output,
			Error:  err,
		}, err

	case "terraform-docs":
		// Generate documentation
		path := request.Params["path"]

		output, err := m.terraformDocs.GenerateMarkdown(ctx, path)
		return ToolResult{
			Tool:   "terraform-docs",
			Action: request.Action,
			Output: output,
			Error:  err,
		}, err

	default:
		return ToolResult{}, fmt.Errorf("unknown tool: %s", request.Tool)
	}
}

// LLMToolRequest represents a tool use request from the LLM
type LLMToolRequest struct {
	Tool   string            `json:"tool"`
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

// ToolResult represents the output from a tool
type ToolResult struct {
	Tool   string
	Action string
	Output string
	Error  error
}

// InvestigationReport contains the complete investigation results
type InvestigationReport struct {
	Objective string
	ToolsUsed []ToolResult
	Analysis  string
}

// Example: Automated Security Audit
func (m *LLMWithToolsModule) SecurityAudit(ctx context.Context, provider string) (*SecurityAuditReport, error) {
	objective := fmt.Sprintf(`
Perform a comprehensive security audit of %s infrastructure:
1. Use Steampipe to find publicly accessible resources
2. Check for unencrypted storage
3. Identify overly permissive security groups
4. Analyze any Terraform code for security issues
5. Provide prioritized recommendations
`, provider)

	report, err := m.InvestigateWithTools(ctx, objective)
	if err != nil {
		return nil, err
	}

	// Parse findings into structured report
	return &SecurityAuditReport{
		Provider:      provider,
		ToolsExecuted: len(report.ToolsUsed),
		Findings:      report.Analysis,
		Timestamp:     ctx.Value("timestamp").(string),
	}, nil
}

// SecurityAuditReport contains security audit results
type SecurityAuditReport struct {
	Provider      string
	ToolsExecuted int
	Findings      string
	Timestamp     string
}

// Example: Cost Optimization Analysis
func (m *LLMWithToolsModule) CostOptimization(ctx context.Context, tfplanPath string) (*CostReport, error) {
	objective := fmt.Sprintf(`
Analyze infrastructure costs and provide optimization recommendations:
1. Use OpenInfraQuote to analyze the Terraform plan at %s
2. Query current resource utilization with Steampipe
3. Identify unused or underutilized resources
4. Suggest cost-saving changes with specific Terraform code
`, tfplanPath)

	report, err := m.InvestigateWithTools(ctx, objective)
	if err != nil {
		return nil, err
	}

	return &CostReport{
		PlanFile:        tfplanPath,
		Analysis:        report.Analysis,
		ToolsUsed:       len(report.ToolsUsed),
		Recommendations: extractRecommendations(report.Analysis),
	}, nil
}

// CostReport contains cost analysis results
type CostReport struct {
	PlanFile        string
	Analysis        string
	ToolsUsed       int
	Recommendations []string
}

func extractRecommendations(analysis string) []string {
	// In production, this would parse the LLM output
	return []string{"Recommendations would be extracted here"}
}
