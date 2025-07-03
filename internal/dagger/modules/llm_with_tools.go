package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

// LLMWithToolsModule extends LLM with ability to use other modules as tools
type LLMWithToolsModule struct {
	client          *dagger.Client
	model           string
	steampipeModule *SteampipeModule
	openInfraModule *OpenInfraQuoteModule
	terraformDocs   *TerraformDocsModule
	infraMapModule  *InfraMapModule
}

// NewLLMWithToolsModule creates an LLM that can use other modules as tools
func NewLLMWithToolsModule(client *dagger.Client, model string) *LLMWithToolsModule {
	return &LLMWithToolsModule{
		client:          client,
		model:           model,
		steampipeModule: NewSteampipeModule(client),
		openInfraModule: NewOpenInfraQuoteModule(client),
		terraformDocs:   NewTerraformDocsModule(client),
		infraMapModule:  NewInfraMapModule(client),
	}
}

// InvestigateWithTools performs a complete investigation using all available tools
func (m *LLMWithToolsModule) InvestigateWithTools(ctx context.Context, objective string) (*InvestigationReport, error) {
	// Define available tools for the LLM
	toolsPrompt := `
You have access to the following tools:

1. STEAMPIPE_QUERY: Execute SQL queries against cloud infrastructure
   Usage: {"tool": "steampipe", "action": "query", "params": {"provider": "aws", "sql": "SELECT ..."}}
   - Query real AWS/Azure/GCP resources
   - Use information_schema.columns to discover table schemas
   - Check tags to identify Terraform-managed resources

2. TERRAFORM_PLAN: Run terraform plan and capture output
   Usage: {"tool": "terraform", "action": "plan", "params": {"path": ".", "format": "json"}}
   
3. TERRAFORM_STATE: Read terraform state file
   Usage: {"tool": "terraform", "action": "show-state", "params": {"path": "terraform.tfstate"}}

4. COST_ANALYSIS: Analyze costs of Terraform plans  
   Usage: {"tool": "openinfraquote", "action": "analyze", "params": {"file": "path/to/tfplan.json", "region": "us-east-1"}}

5. TERRAFORM_DOCS: Generate documentation for Terraform modules
   Usage: {"tool": "terraform-docs", "action": "generate", "params": {"path": "path/to/module"}}

6. SECURITY_SCAN: Scan for security issues with Checkov
   Usage: {"tool": "checkov", "action": "scan", "params": {"path": "path/to/code"}}

7. TFLINT: Lint Terraform code
   Usage: {"tool": "tflint", "action": "lint", "params": {"path": "."}}

8. INFRACOST: Estimate costs for Terraform changes
   Usage: {"tool": "infracost", "action": "breakdown", "params": {"path": "."}}

9. INFRAMAP_DIAGRAM: Generate infrastructure diagrams from Terraform state or HCL
   Usage: {"tool": "inframap", "action": "diagram", "params": {"input": "terraform.tfstate", "format": "png"}}
   Or for HCL: {"tool": "inframap", "action": "diagram-hcl", "params": {"path": ".", "format": "svg"}}

IMPORTANT: When you want to use a tool, you MUST respond with ONLY a valid JSON object. Do not include any text before or after the JSON.
Example correct response:
{"tool": "steampipe", "action": "query", "params": {"provider": "aws", "sql": "SELECT count(*) FROM aws_s3_bucket"}}

After I execute the tool and return results, you can then use another tool or provide your final analysis.
DO NOT describe your plan - just execute it by returning JSON tool requests.
`

	// Initial prompt with objective and available tools
	systemPrompt := "You are an infrastructure investigator with access to real tools. Always respond with JSON when using tools, never with explanatory text. " + toolsPrompt

	// Create LLM with tool-use system prompt
	llm := m.client.LLM(dagger.LLMOpts{
		Model: m.model,
	}).WithSystemPrompt(systemPrompt)

	// Start investigation
	conversation := llm.WithPrompt(fmt.Sprintf("Task: %s\n\nRespond with a JSON tool request to begin investigation. Do not explain, just return JSON.", objective))

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

		slog.Debug("LLM response", "iteration", i+1, "response_length", len(response))
		
		// Try to extract JSON from response (handle cases where LLM adds text)
		jsonStart := strings.Index(response, "{")
		jsonEnd := strings.LastIndex(response, "}")
		
		if jsonStart >= 0 && jsonEnd > jsonStart {
			jsonStr := response[jsonStart:jsonEnd+1]
			
			var toolRequest LLMToolRequest
			if err := json.Unmarshal([]byte(jsonStr), &toolRequest); err == nil && toolRequest.Tool != "" {
				slog.Info("Executing tool", "tool", toolRequest.Tool, "action", toolRequest.Action)
				
				// Execute the requested tool
				result, err := m.executeTool(ctx, toolRequest)
				if err != nil {
					// Tell LLM about the error
					conversation = conversation.WithPrompt(fmt.Sprintf("Tool error: %v\n\nTry another tool or provide your analysis.", err))
					toolResults = append(toolResults, ToolResult{
						Tool:   toolRequest.Tool,
						Action: toolRequest.Action,
						Error:  err,
					})
					continue
				}

				toolResults = append(toolResults, result)

				// Feed results back to LLM
				conversation = conversation.WithPrompt(fmt.Sprintf(
					"Tool '%s' executed successfully. Output:\n%s\n\nRespond with another JSON tool request to continue, or provide your final analysis as plain text.",
					toolRequest.Tool, result.Output,
				))
			} else {
				// Failed to parse as tool request, assume it's the final analysis
				return &InvestigationReport{
					Objective: objective,
					ToolsUsed: toolResults,
					Analysis:  response,
				}, nil
			}
		} else {
			// No JSON found, assume it's the final analysis
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
		
		// Get credentials for the provider
		credentials := getProviderCredentials(provider)

		output, err := m.steampipeModule.RunQuery(ctx, provider, sql, credentials, "json")
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

	case "terraform":
		// Handle terraform operations
		switch request.Action {
		case "plan":
			path := request.Params["path"]
			if path == "" {
				path = "."
			}
			// TODO: Implement terraform plan execution
			return ToolResult{
				Tool:   "terraform",
				Action: request.Action,
				Output: "Terraform plan would be executed here",
				Error:  fmt.Errorf("terraform plan not yet implemented"),
			}, nil
			
		case "show-state":
			statePath := request.Params["path"]
			// TODO: Implement terraform state reading
			return ToolResult{
				Tool:   "terraform",
				Action: request.Action,
				Output: fmt.Sprintf("Terraform state would be read from: %s", statePath),
				Error:  fmt.Errorf("terraform state reading not yet implemented"),
			}, nil
			
		default:
			return ToolResult{}, fmt.Errorf("unknown terraform action: %s", request.Action)
		}
		
	case "checkov":
		// Security scanning
		path := request.Params["path"]
		if path == "" {
			path = "."
		}
		// Use existing terraform-tools checkov functionality
		// TODO: Create CheckovModule and integrate
		return ToolResult{
			Tool:   "checkov",
			Action: request.Action,
			Output: "Checkov scan would run here",
			Error:  fmt.Errorf("checkov integration pending"),
		}, nil
		
	case "tflint":
		// Terraform linting
		path := request.Params["path"]
		if path == "" {
			path = "."
		}
		// TODO: Create TFLintModule
		return ToolResult{
			Tool:   "tflint",
			Action: request.Action,
			Output: "TFLint would run here",
			Error:  fmt.Errorf("tflint integration pending"),
		}, nil
		
	case "infracost":
		// Cost estimation
		path := request.Params["path"]
		if path == "" {
			path = "."
		}
		// TODO: Create InfracostModule
		return ToolResult{
			Tool:   "infracost",
			Action: request.Action,
			Output: "Infracost would run here",
			Error:  fmt.Errorf("infracost integration pending"),
		}, nil

	case "inframap":
		// Generate infrastructure diagram
		switch request.Action {
		case "diagram":
			input := request.Params["input"]
			format := request.Params["format"]
			if format == "" {
				format = "png"
			}

			output, err := m.infraMapModule.GenerateFromState(ctx, input, format)
			return ToolResult{
				Tool:   "inframap",
				Action: request.Action,
				Output: output,
				Error:  err,
			}, err

		case "diagram-hcl":
			path := request.Params["path"]
			format := request.Params["format"]
			if format == "" {
				format = "png"
			}

			output, err := m.infraMapModule.GenerateFromHCL(ctx, path, format)
			return ToolResult{
				Tool:   "inframap",
				Action: request.Action,
				Output: output,
				Error:  err,
			}, err

		default:
			return ToolResult{}, fmt.Errorf("unknown inframap action: %s", request.Action)
		}

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

// getProviderCredentials returns credentials for cloud providers
func getProviderCredentials(provider string) map[string]string {
	creds := make(map[string]string)
	
	switch provider {
	case "aws":
		// AWS credentials from environment
		if v := os.Getenv("AWS_ACCESS_KEY_ID"); v != "" {
			creds["AWS_ACCESS_KEY_ID"] = v
		}
		if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v != "" {
			creds["AWS_SECRET_ACCESS_KEY"] = v
		}
		if v := os.Getenv("AWS_SESSION_TOKEN"); v != "" {
			creds["AWS_SESSION_TOKEN"] = v
		}
		if v := os.Getenv("AWS_REGION"); v != "" {
			creds["AWS_REGION"] = v
		} else {
			creds["AWS_REGION"] = "us-east-1"
		}
		
		// If no environment credentials, try to load from ~/.aws/credentials
		if creds["AWS_ACCESS_KEY_ID"] == "" {
			if homeDir := os.Getenv("HOME"); homeDir != "" {
				credFile := filepath.Join(homeDir, ".aws", "credentials")
				if content, err := os.ReadFile(credFile); err == nil {
					// Simple parsing for default profile
					lines := strings.Split(string(content), "\n")
					inDefaultProfile := false
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if line == "[default]" {
							inDefaultProfile = true
							continue
						}
						if strings.HasPrefix(line, "[") && line != "[default]" {
							inDefaultProfile = false
							continue
						}
						if inDefaultProfile {
							if strings.HasPrefix(line, "aws_access_key_id") {
								parts := strings.SplitN(line, "=", 2)
								if len(parts) == 2 {
									creds["AWS_ACCESS_KEY_ID"] = strings.TrimSpace(parts[1])
								}
							}
							if strings.HasPrefix(line, "aws_secret_access_key") {
								parts := strings.SplitN(line, "=", 2)
								if len(parts) == 2 {
									creds["AWS_SECRET_ACCESS_KEY"] = strings.TrimSpace(parts[1])
								}
							}
						}
					}
				}
			}
		}
	case "azure":
		// Azure credentials
		if v := os.Getenv("AZURE_TENANT_ID"); v != "" {
			creds["AZURE_TENANT_ID"] = v
		}
		if v := os.Getenv("AZURE_CLIENT_ID"); v != "" {
			creds["AZURE_CLIENT_ID"] = v
		}
		if v := os.Getenv("AZURE_CLIENT_SECRET"); v != "" {
			creds["AZURE_CLIENT_SECRET"] = v
		}
	case "gcp":
		// GCP credentials
		if v := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); v != "" {
			creds["GOOGLE_APPLICATION_CREDENTIALS"] = v
		}
	}
	
	return creds
}
