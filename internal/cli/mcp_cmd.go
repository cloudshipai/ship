package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP (Model Context Protocol) server",
	Long: `Start an MCP server that exposes Ship CLI tools for use with AI assistants.
	
This allows AI assistants like Claude, Cursor, or other MCP-compatible clients to use Ship CLI
functionality directly. The server exposes tools for Terraform analysis, infrastructure investigation,
and cloud resource management.`,
	RunE: runMCPServer,
}

func init() {
	rootCmd.AddCommand(mcpCmd)
	
	mcpCmd.Flags().Int("port", 0, "Port to listen on (0 for stdio)")
	mcpCmd.Flags().String("host", "localhost", "Host to bind to")
	mcpCmd.Flags().Bool("stdio", true, "Use stdio transport (default)")
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	host, _ := cmd.Flags().GetString("host")
	useStdio, _ := cmd.Flags().GetBool("stdio")
	
	// Create MCP server
	s := server.NewMCPServer(
		"ship-cli", 
		"1.0.0",
	)
	
	// Add tools for terraform analysis
	addTerraformTools(s)
	
	// Add tools for infrastructure investigation
	addInvestigationTools(s)
	
	// Add tools for general cloud operations
	addCloudTools(s)
	
	// Add resources for documentation and help
	addResources(s)
	
	// Add prompts for common use cases
	addPrompts(s)
	
	// Start server
	if useStdio || port == 0 {
		fmt.Fprintf(os.Stderr, "Starting Ship CLI MCP server on stdio...\n")
		return server.ServeStdio(s)
	} else {
		fmt.Fprintf(os.Stderr, "Starting Ship CLI MCP server on %s:%d...\n", host, port)
		// Note: HTTP serving may require additional setup in this library version
		return fmt.Errorf("HTTP server not implemented in this version, use --stdio")
	}
}

func addTerraformTools(s *server.MCPServer) {
	// Terraform Lint Tool
	lintTool := mcp.NewTool("terraform_lint",
		mcp.WithDescription("Run TFLint on Terraform code to check for syntax errors and best practices"),
		mcp.WithString("directory", 
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("config", 
			mcp.Description("Path to TFLint configuration file"),
		),
	)
	
	s.AddTool(lintTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "lint"}
		
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		
		return executeShipCommand(args)
	})
	
	// Terraform Security Scan Tool
	securityTool := mcp.NewTool("terraform_security_scan",
		mcp.WithDescription("Run Checkov security scan on Terraform code"),
		mcp.WithString("directory", 
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("framework", 
			mcp.Description("Security framework to use (terraform, arm, etc.)"),
		),
	)
	
	s.AddTool(securityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "checkov-scan"}
		
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "--framework", framework)
		}
		
		return executeShipCommand(args)
	})
	
	// Terraform Cost Estimation Tool
	costTool := mcp.NewTool("terraform_cost_estimate",
		mcp.WithDescription("Estimate infrastructure costs for Terraform code"),
		mcp.WithString("directory", 
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("cloud", 
			mcp.Description("Cloud provider (aws, azure, gcp)"),
		),
	)
	
	s.AddTool(costTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "cost-estimate"}
		
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if cloud := request.GetString("cloud", ""); cloud != "" {
			args = append(args, "--cloud", cloud)
		}
		
		return executeShipCommand(args)
	})
	
	// Terraform Documentation Tool
	docsTool := mcp.NewTool("terraform_generate_docs",
		mcp.WithDescription("Generate documentation for Terraform modules"),
		mcp.WithString("directory", 
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("format", 
			mcp.Description("Output format (markdown, json)"),
			mcp.Enum("markdown", "json"),
		),
		mcp.WithBoolean("show_examples", 
			mcp.Description("Include examples in documentation"),
		),
	)
	
	s.AddTool(docsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "generate-docs"}
		
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if showExamples := request.GetBool("show_examples", false); showExamples {
			args = append(args, "--show-examples")
		}
		
		return executeShipCommand(args)
	})
}

func addInvestigationTools(s *server.MCPServer) {
	// AI Infrastructure Investigation Tool
	investigateTool := mcp.NewTool("ai_investigate",
		mcp.WithDescription("Investigate cloud infrastructure using natural language queries powered by Steampipe"),
		mcp.WithString("prompt", 
			mcp.Required(),
			mcp.Description("Natural language description of what to investigate (e.g., 'Show me all S3 buckets', 'Check for security issues')"),
		),
		mcp.WithString("provider", 
			mcp.Description("Cloud provider to investigate (aws, azure, gcp)"),
			mcp.Enum("aws", "azure", "gcp"),
		),
		mcp.WithString("aws_profile", 
			mcp.Description("AWS profile to use for authentication"),
		),
		mcp.WithString("aws_region", 
			mcp.Description("AWS region to focus on"),
		),
		mcp.WithBoolean("execute", 
			mcp.Description("Whether to execute the generated queries (default: false, will only show plan)"),
		),
	)
	
	s.AddTool(investigateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"ai-investigate"}
		
		if prompt, err := request.RequireString("prompt"); err == nil {
			args = append(args, "--prompt", prompt)
		} else {
			return mcp.NewToolResultError("prompt parameter is required"), nil
		}
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, "--provider", provider)
		}
		if awsProfile := request.GetString("aws_profile", ""); awsProfile != "" {
			args = append(args, "--aws-profile", awsProfile)
		}
		if awsRegion := request.GetString("aws_region", ""); awsRegion != "" {
			args = append(args, "--aws-region", awsRegion)
		}
		if execute := request.GetBool("execute", false); execute {
			args = append(args, "--execute")
		}
		
		return executeShipCommand(args)
	})
}

func addCloudTools(s *server.MCPServer) {
	// Push artifacts to Cloudship for analysis
	pushTool := mcp.NewTool("cloudship_push",
		mcp.WithDescription("Upload and analyze infrastructure artifacts with Cloudship AI"),
		mcp.WithString("file", 
			mcp.Required(),
			mcp.Description("Path to the file to upload (Terraform plan, SBOM, etc.)"),
		),
		mcp.WithString("type", 
			mcp.Description("Type of artifact being uploaded"),
			mcp.Enum("terraform-plan", "sbom", "dockerfile", "kubernetes-manifest"),
		),
	)
	
	s.AddTool(pushTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"push"}
		
		if file, err := request.RequireString("file"); err == nil {
			args = append(args, file)
		} else {
			return mcp.NewToolResultError("file parameter is required"), nil
		}
		if artifactType := request.GetString("type", ""); artifactType != "" {
			args = append(args, "--type", artifactType)
		}
		
		return executeShipCommand(args)
	})
}

func addResources(s *server.MCPServer) {
	// Help resource
	helpResource := mcp.NewResource("ship://help",
		"Ship CLI Help",
		mcp.WithResourceDescription("Complete help and usage information for Ship CLI"),
		mcp.WithMIMEType("text/markdown"),
	)
	
	s.AddResource(helpResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		result, err := executeShipCommand([]string{"--help"})
		if err != nil {
			return nil, err
		}
		
		// Extract text from result - the result should be a simple text response
		var helpText string
		if result != nil {
			helpText = result.Content[0].Text
		}
		
		return []mcp.ResourceContents{
			{
				URI:      "ship://help",
				MIMEType: "text/markdown",
				Text:     helpText,
			},
		}, nil
	})
	
	// Available tools resource
	toolsResource := mcp.NewResource("ship://tools",
		"Available Ship CLI Tools",
		mcp.WithResourceDescription("List of all available Ship CLI tools and their capabilities"),
		mcp.WithMIMEType("text/markdown"),
	)
	
	s.AddResource(toolsResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content := `# Ship CLI Tools

## Terraform Tools
- **lint**: Run TFLint for syntax and best practices
- **checkov-scan**: Security scanning with Checkov
- **security-scan**: Alternative security scanning with Trivy
- **cost-estimate**: Cost estimation with Infracost
- **cost-analysis**: Alternative cost analysis with OpenInfraQuote
- **generate-docs**: Generate documentation with terraform-docs

## Infrastructure Investigation
- **ai-investigate**: Natural language infrastructure investigation using Steampipe

## Cloud Operations
- **push**: Upload artifacts to Cloudship for AI analysis
- **auth**: Manage authentication and configuration

## Examples
- ` + "`ship terraform-tools lint`" + ` - Lint current directory
- ` + "`ship ai-investigate --prompt \"Show me all S3 buckets\" --execute`" + ` - Investigate S3 buckets
- ` + "`ship push terraform.tfplan --type terraform-plan`" + ` - Upload Terraform plan for analysis
`
		
		return []mcp.ResourceContents{
			{
				URI:      "ship://tools",
				MIMEType: "text/markdown",
				Text:     content,
			},
		}, nil
	})
}

func addPrompts(s *server.MCPServer) {
	// Security audit prompt
	securityPrompt := mcp.NewPrompt("security_audit",
		mcp.WithPromptDescription("Comprehensive security audit of cloud infrastructure"),
		mcp.WithArgument("provider", 
			mcp.ArgumentDescription("Cloud provider to audit (aws, azure, gcp)"),
		),
	)
	
	s.AddPrompt(securityPrompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		provider := "aws"
		if p, ok := request.Params.Arguments["provider"]; ok && p != "" {
			provider = p
		}
		
		return &mcp.GetPromptResult{
			Description: "Comprehensive security audit workflow",
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf(`Please perform a comprehensive security audit of my %s infrastructure. Follow these steps:

1. First, use the ai_investigate tool to check for security issues:
   - Find publicly accessible resources
   - Check for unencrypted storage
   - Look for overly permissive security groups
   
2. If there are Terraform files, run terraform_security_scan to check infrastructure-as-code

3. Summarize all findings with:
   - Critical security issues requiring immediate attention
   - Recommendations for improvement
   - Best practices to implement

Please be thorough and provide actionable recommendations.`, provider),
					},
				},
			},
		}, nil
	})
	
	// Cost optimization prompt
	costPrompt := mcp.NewPrompt("cost_optimization",
		mcp.WithPromptDescription("Identify cost optimization opportunities"),
		mcp.WithArgument("provider", 
			mcp.ArgumentDescription("Cloud provider to analyze (aws, azure, gcp)"),
		),
	)
	
	s.AddPrompt(costPrompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		provider := "aws"
		if p, ok := request.Params.Arguments["provider"]; ok && p != "" {
			provider = p
		}
		
		return &mcp.GetPromptResult{
			Description: "Cost optimization analysis workflow",
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf(`Help me optimize costs for my %s infrastructure:

1. Use ai_investigate to find unused or idle resources:
   - Unused EBS volumes
   - Idle EC2 instances
   - Unattached elastic IPs
   - Empty S3 buckets with storage classes

2. If Terraform code is available, use terraform_cost_estimate to get cost projections

3. Provide a prioritized list of cost-saving recommendations:
   - Quick wins (low effort, high impact)
   - Medium-term optimizations
   - Long-term architectural improvements

Include estimated cost savings where possible.`, provider),
					},
				},
			},
		}, nil
	})
}

func executeShipCommand(args []string) (*mcp.CallToolResult, error) {
	// Get the current binary path
	executable, err := os.Executable()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get executable path: %v", err)), nil
	}
	
	// Execute the ship command
	cmd := exec.Command(executable, args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Command failed: %s\n\nOutput:\n%s", err.Error(), string(output))), nil
	}
	
	return mcp.NewToolResultText(string(output)), nil
}