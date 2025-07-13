// MCP Server implementation for Ship CLI

package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP (Model Context Protocol) server",
	Long: `Start an MCP server that exposes Ship CLI tools for use with AI assistants.
	
This allows AI assistants like Claude, Cursor, or other MCP-compatible clients to use Ship CLI
functionality directly. The server exposes tools for Terraform analysis, cost estimation,
security scanning, documentation generation, and cloud resource management.`,
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

	// Investigation tools removed - focusing on Terraform analysis

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
		mcp.WithString("format",
			mcp.Description("Output format: default, json, compact"),
			mcp.Enum("default", "json", "compact"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save lint results"),
		),
	)

	s.AddTool(lintTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "lint"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})

	// Terraform Security Scan Tool (Checkov)
	checkovTool := mcp.NewTool("terraform_checkov_scan",
		mcp.WithDescription("Run Checkov security scan on Terraform code for policy compliance"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: cli, json, junit, sarif"),
			mcp.Enum("cli", "json", "junit", "sarif"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save scan results"),
		),
	)

	s.AddTool(checkovTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "checkov-scan"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})

	// Terraform Security Scan Tool (Alternative)
	securityTool := mcp.NewTool("terraform_security_scan",
		mcp.WithDescription("Run alternative security scan on Terraform code using Trivy"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		// mcp.WithBoolean("push",
		//	mcp.Description("Push results to Cloudship for analysis"),
		// ),
	)

	s.AddTool(securityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "security-scan"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		// Push functionality disabled during staging
		// if push := request.GetBool("push", false); push {
		//	args = append(args, "--push")
		// }

		return executeShipCommand(args)
	})

	// Terraform Cost Analysis Tool
	costAnalysisTool := mcp.NewTool("terraform_cost_analysis",
		mcp.WithDescription("Analyze infrastructure costs using OpenInfraQuote"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing (e.g., us-east-1, us-west-2)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: json, table"),
			mcp.Enum("json", "table"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save cost analysis"),
		),
	)

	s.AddTool(costAnalysisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "cost-analysis"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})

	// Terraform Documentation Tool
	docsTool := mcp.NewTool("terraform_generate_docs",
		mcp.WithDescription("Generate documentation for Terraform modules using terraform-docs"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("filename",
			mcp.Description("Filename to save documentation as (default README.md)"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save documentation"),
		),
	)

	s.AddTool(docsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "generate-docs"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if filename := request.GetString("filename", ""); filename != "" {
			args = append(args, "--filename", filename)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})

	// Terraform Diagram Generation Tool
	diagramTool := mcp.NewTool("terraform_generate_diagram",
		mcp.WithDescription("Generate infrastructure diagrams from Terraform state"),
		mcp.WithString("input",
			mcp.Description("Input directory or file containing Terraform files (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: png, svg, pdf, dot"),
			mcp.Enum("png", "svg", "pdf", "dot"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save diagram"),
		),
		mcp.WithBoolean("hcl",
			mcp.Description("Generate from HCL files instead of state file"),
		),
		mcp.WithString("provider",
			mcp.Description("Filter by specific provider (aws, google, azurerm)"),
			mcp.Enum("aws", "google", "azurerm"),
		),
	)

	s.AddTool(diagramTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "generate-diagram"}

		if input := request.GetString("input", ""); input != "" {
			args = append(args, input)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if hcl := request.GetBool("hcl", false); hcl {
			args = append(args, "--hcl")
		}
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, "--provider", provider)
		}

		return executeShipCommand(args)
	})
}

// Investigation tools removed to focus on Terraform analysis workflows

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
		if result != nil && len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(mcp.TextContent); ok {
				helpText = textContent.Text
			}
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
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
- **cost-analysis**: Cost analysis with OpenInfraQuote
- **generate-docs**: Generate documentation with terraform-docs
- **generate-diagram**: Generate infrastructure diagrams with InfraMap

## Cloud Operations
- **push**: Upload artifacts to Cloudship for AI analysis
- **auth**: Manage authentication and configuration

## Examples
- ` + "`ship terraform-tools lint`" + ` - Lint current directory
- ` + "`ship terraform-tools generate-diagram . --hcl --format png`" + ` - Generate infrastructure diagram
- ` + "`ship push terraform.tfplan --type terraform-plan`" + ` - Upload Terraform plan for analysis
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
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
		return &mcp.GetPromptResult{
			Description: "Comprehensive security audit workflow",
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: `Please perform a comprehensive security audit of my Terraform infrastructure. Follow these steps:

1. Run terraform_checkov_scan to identify security issues in infrastructure-as-code
2. Run terraform_security_scan for additional security analysis  
3. Use terraform_lint to check for configuration best practices

4. Summarize all findings with:
   - Critical security issues requiring immediate attention
   - Recommendations for improvement
   - Best practices to implement

Please be thorough and provide actionable recommendations.`,
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
		return &mcp.GetPromptResult{
			Description: "Cost optimization analysis workflow",
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: `Help me optimize costs for my Terraform infrastructure:

1. Use terraform_cost_analysis to analyze current cost projections
2. Review Terraform configurations for cost optimization opportunities
3. Use terraform_lint to identify inefficient resource configurations

4. Provide a prioritized list of cost-saving recommendations:
   - Quick wins (resource rightsizing, unused resources)
   - Medium-term optimizations (reserved instances, storage classes)
   - Long-term architectural improvements

Include estimated cost savings where possible.`,
					},
				},
			},
		}, nil
	})
}

// Maximum tokens allowed in MCP response (conservative estimate)
const maxMCPTokens = 20000

// Rough estimation: 1 token ≈ 4 characters for typical text
const charsPerToken = 4

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

	outputStr := string(output)

	// Check if output needs to be chunked
	if needsChunking(outputStr) {
		return createChunkedResponse(outputStr), nil
	}

	return mcp.NewToolResultText(outputStr), nil
}

func needsChunking(text string) bool {
	return utf8.RuneCountInString(text) > (maxMCPTokens * charsPerToken)
}

func createChunkedResponse(text string) *mcp.CallToolResult {
	maxChunkSize := maxMCPTokens * charsPerToken

	// Split text into chunks, preferring to break at newlines
	chunks := smartChunk(text, maxChunkSize)

	if len(chunks) <= 1 {
		return mcp.NewToolResultText(text)
	}

	// Create a summary response with information about chunking
	summary := fmt.Sprintf(`Output is large (%d characters, ~%d tokens) and has been summarized.

TOTAL CHUNKS: %d

SUMMARY OF FIRST CHUNK:
%s

--- [Content continues in additional chunks] ---

To see the full output, you can:
1. Run the command with smaller scope (specific directory/subset)
2. Use filtering options if available
3. Process chunks individually if needed

FIRST CHUNK PREVIEW (showing first %d characters):
%s`,
		utf8.RuneCountInString(text),
		utf8.RuneCountInString(text)/charsPerToken,
		len(chunks),
		getChunkSummary(chunks[0]),
		maxChunkSize/4, // Show 1/4 of max chunk size as preview
		truncateText(chunks[0], maxChunkSize/4),
	)

	return mcp.NewToolResultText(summary)
}

func smartChunk(text string, maxSize int) []string {
	if utf8.RuneCountInString(text) <= maxSize {
		return []string{text}
	}

	var chunks []string
	lines := strings.Split(text, "\n")

	currentChunk := ""
	currentSize := 0

	for _, line := range lines {
		lineSize := utf8.RuneCountInString(line) + 1 // +1 for newline

		// If adding this line would exceed the chunk size, start a new chunk
		if currentSize+lineSize > maxSize && currentChunk != "" {
			chunks = append(chunks, strings.TrimSuffix(currentChunk, "\n"))
			currentChunk = line + "\n"
			currentSize = lineSize
		} else {
			currentChunk += line + "\n"
			currentSize += lineSize
		}
	}

	// Add the last chunk if it has content
	if currentChunk != "" {
		chunks = append(chunks, strings.TrimSuffix(currentChunk, "\n"))
	}

	return chunks
}

func getChunkSummary(chunk string) string {
	lines := strings.Split(chunk, "\n")
	if len(lines) == 0 {
		return "Empty content"
	}

	// Try to identify the type of content
	firstLine := strings.TrimSpace(lines[0])

	if strings.Contains(chunk, "CRITICAL") || strings.Contains(chunk, "HIGH") {
		return "Security scan results with findings"
	} else if strings.Contains(chunk, "resource \"") {
		return "Terraform configuration analysis"
	} else if strings.Contains(chunk, "$") && strings.Contains(chunk, "cost") {
		return "Cost analysis results"
	} else if strings.Contains(chunk, "Error:") || strings.Contains(chunk, "Warning:") {
		return "Tool output with errors/warnings"
	} else {
		return fmt.Sprintf("Tool output starting with: %s", truncateText(firstLine, 100))
	}
}

func truncateText(text string, maxLen int) string {
	if utf8.RuneCountInString(text) <= maxLen {
		return text
	}

	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}

	return string(runes[:maxLen]) + "..."
}
