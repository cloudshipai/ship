// MCP Server implementation for Ship CLI

package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"github.com/cloudshipai/ship/internal/telemetry"
	"github.com/cloudshipai/ship/pkg/ship"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

// hardcodedMCPServers contains the built-in external MCP server configurations
var hardcodedMCPServers = map[string]ship.MCPServerConfig{
	"filesystem": {
		Name:      "filesystem",
		Command:   "npx",
		Args:      []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "FILESYSTEM_ROOT",
				Description: "Root directory for filesystem operations (overrides /tmp default)",
				Required:    false,
				Default:     "/tmp",
			},
		},
	},
	"memory": {
		Name:      "memory",
		Command:   "npx",
		Args:      []string{"-y", "@modelcontextprotocol/server-memory"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "MEMORY_STORAGE_PATH",
				Description: "Path for persistent memory storage",
				Required:    false,
				Default:     "/tmp/mcp-memory",
			},
			{
				Name:        "MEMORY_MAX_SIZE",
				Description: "Maximum memory storage size (e.g., 100MB)",
				Required:    false,
				Default:     "50MB",
			},
		},
	},
	"brave-search": {
		Name:      "brave-search",
		Command:   "npx",
		Args:      []string{"-y", "@modelcontextprotocol/server-brave-search"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "BRAVE_API_KEY",
				Description: "Brave Search API key for search functionality",
				Required:    true,
				Secret:      true,
			},
			{
				Name:        "BRAVE_SEARCH_COUNT",
				Description: "Number of search results to return (default: 10)",
				Required:    false,
				Default:     "10",
			},
		},
	},
}

var mcpCmd = &cobra.Command{
	Use:   "mcp [tool]",
	Short: "Start MCP server for a specific tool or all tools",
	Long: `Start an MCP server that exposes specific Ship CLI tools for AI assistants.

Available tools:
  lint       - TFLint for syntax and best practices
  checkov    - Checkov security scanning  
  trivy      - Trivy security scanning
  cost       - OpenInfraQuote cost analysis
  docs       - terraform-docs documentation
  diagram    - InfraMap diagram generation
  all        - All tools (default if no tool specified)

External MCP Servers:
  filesystem     - Filesystem operations MCP server
  memory         - Memory/knowledge storage MCP server
  brave-search   - Brave search MCP server

Examples:
  ship mcp lint        # MCP server for just TFLint
  ship mcp checkov     # MCP server for just Checkov
  ship mcp all         # MCP server for all tools
  ship mcp filesystem     # Proxy filesystem operations MCP server
  ship mcp memory         # Proxy memory/knowledge storage MCP server
  ship mcp brave-search --var BRAVE_API_KEY=your_api_key   # Proxy Brave search with API key
  ship mcp cost --var AWS_REGION=us-east-1 --var DEBUG=true  # Pass multiple environment variables`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMCPServer,
}

func init() {
	rootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().Int("port", 0, "Port to listen on (0 for stdio)")
	mcpCmd.Flags().String("host", "localhost", "Host to bind to")
	mcpCmd.Flags().Bool("stdio", true, "Use stdio transport (default)")
	mcpCmd.Flags().StringToString("var", nil, "Environment variables for MCP servers and containers (e.g., --var API_KEY=value --var DEBUG=true)")
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	host, _ := cmd.Flags().GetString("host")
	useStdio, _ := cmd.Flags().GetBool("stdio")
	envVars, _ := cmd.Flags().GetStringToString("var")

	// Determine which tool to serve
	toolName := "all"
	if len(args) > 0 {
		toolName = args[0]
	}

	// Track MCP command usage
	telemetry.TrackMCPCommand(toolName)

	// Create MCP server
	serverName := fmt.Sprintf("ship-%s", toolName)
	s := server.NewMCPServer(serverName, "1.0.0")

	// Set environment variables for containerized tools
	if len(envVars) > 0 {
		setContainerEnvironmentVars(envVars)
	}
	
	// Add specific tools based on argument
	switch toolName {
	case "lint":
		addLintTool(s)
	case "checkov":
		addCheckovTool(s)
	case "trivy", "security":
		addTrivyTool(s)
	case "cost":
		addCostTool(s)
	case "docs":
		addDocsTool(s)
	case "diagram":
		addDiagramTool(s)
	case "all":
		addTerraformTools(s)
	default:
		// Check if this is an external MCP server
		if isExternalMCPServer(toolName) {
			return runMCPProxy(cmd, toolName)
		}
		return fmt.Errorf("unknown tool: %s. Available: lint, checkov, trivy, cost, docs, diagram, all, filesystem, memory, brave-search", toolName)
	}

	// Add resources for documentation and help
	addResources(s)

	// Add prompts only for 'all' mode
	if toolName == "all" {
		addPrompts(s)
	}

	// Start server
	if useStdio || port == 0 {
		fmt.Fprintf(os.Stderr, "Starting %s MCP server on stdio...\n", serverName)
		return server.ServeStdio(s)
	} else {
		fmt.Fprintf(os.Stderr, "Starting %s MCP server on %s:%d...\n", serverName, host, port)
		return fmt.Errorf("HTTP server not implemented in this version, use --stdio")
	}
}

// Individual tool functions
func addLintTool(s *server.MCPServer) {
	lintTool := mcp.NewTool("lint",
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
		args := []string{"tf", "lint"}

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
}

func addCheckovTool(s *server.MCPServer) {
	checkovTool := mcp.NewTool("checkov",
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
		args := []string{"tf", "checkov"}

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
}

func addTrivyTool(s *server.MCPServer) {
	trivyTool := mcp.NewTool("trivy",
		mcp.WithDescription("Run Trivy security scan on Terraform code using Trivy"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
	)

	s.AddTool(trivyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "trivy"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}

		return executeShipCommand(args)
	})
}

func addCostTool(s *server.MCPServer) {
	costTool := mcp.NewTool("cost",
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

	s.AddTool(costTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "cost"}

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
}

func addDocsTool(s *server.MCPServer) {
	docsTool := mcp.NewTool("docs",
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
		args := []string{"tf", "docs"}

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
}

func addDiagramTool(s *server.MCPServer) {
	diagramTool := mcp.NewTool("diagram",
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
		args := []string{"tf", "diagram"}

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

func addTerraformTools(s *server.MCPServer) {
	// Terraform Lint Tool
	lintTool := mcp.NewTool("lint",
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
		args := []string{"tf", "lint"}

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
	checkovTool := mcp.NewTool("checkov",
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
		args := []string{"tf", "checkov"}

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
	securityTool := mcp.NewTool("trivy",
		mcp.WithDescription("Run alternative security scan on Terraform code using Trivy"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
	)

	s.AddTool(securityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "trivy"}

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
	costAnalysisTool := mcp.NewTool("cost",
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
		args := []string{"tf", "cost"}

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
	docsTool := mcp.NewTool("docs",
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
		args := []string{"tf", "docs"}

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
	diagramTool := mcp.NewTool("diagram",
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
		args := []string{"tf", "diagram"}

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
	// Cloud tools functionality removed - focusing on Terraform analysis only
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
- **checkov**: Security scanning with Checkov
- **trivy**: Alternative security scanning with Trivy
- **cost**: Cost analysis with OpenInfraQuote
- **docs**: Generate documentation with terraform-docs
- **diagram**: Generate infrastructure diagrams with InfraMap


## Examples
- ` + "`ship tf lint`" + ` - Lint current directory
- ` + "`ship tf diagram . --hcl --format png`" + ` - Generate infrastructure diagram
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

// isExternalMCPServer checks if the tool name matches a hardcoded external MCP server
func isExternalMCPServer(toolName string) bool {
	_, exists := hardcodedMCPServers[toolName]
	return exists
}

// runMCPProxy starts an MCP proxy server for external MCP servers
func runMCPProxy(cmd *cobra.Command, serverName string) error {
	useStdio, _ := cmd.Flags().GetBool("stdio")
	port, _ := cmd.Flags().GetInt("port")
	envVars, _ := cmd.Flags().GetStringToString("var")
	
	// Get hardcoded server configuration
	mcpConfig, exists := hardcodedMCPServers[serverName]
	if !exists {
		return fmt.Errorf("external MCP server '%s' not found in hardcoded configurations", serverName)
	}
	
	// Validate and merge environment variables
	if err := validateAndMergeVariables(&mcpConfig, envVars); err != nil {
		return fmt.Errorf("variable validation failed: %w", err)
	}
	
	ctx := context.Background()
	
	// Create and connect proxy
	proxy := ship.NewMCPProxy(mcpConfig)
	if err := proxy.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to external MCP server: %w", err)
	}
	defer proxy.Close()
	
	// Discover tools from the external server
	tools, err := proxy.DiscoverTools(ctx)
	if err != nil {
		return fmt.Errorf("failed to discover tools from external server: %w", err)
	}
	
	fmt.Fprintf(os.Stderr, "Discovered %d tools from external MCP server\n", len(tools))
	
	// Create a Ship MCP server with the discovered tools
	shipServer := ship.NewServer(fmt.Sprintf("ship-proxy-%s", serverName), "1.0.0")
	for _, tool := range tools {
		shipServer.AddTool(tool)
	}
	mcpServer := shipServer.Build()
	defer mcpServer.Close()
	
	// Start the server
	if err := mcpServer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}
	
	// Get the mcp-go server instance
	serverInstance := mcpServer.GetMCPGoServer()
	if serverInstance == nil {
		return fmt.Errorf("failed to get MCP server instance")
	}
	
	// Start the proxy server
	if useStdio || port == 0 {
		fmt.Fprintf(os.Stderr, "Starting Ship proxy for %s on stdio...\n", serverName)
		fmt.Fprintf(os.Stderr, "Available tools: %v\n", mcpServer.GetRegistry().ListTools())
		return server.ServeStdio(serverInstance)
	} else {
		return fmt.Errorf("HTTP server not implemented in this version, use --stdio")
	}
}

// validateAndMergeVariables validates required variables and merges user-provided vars with config
func validateAndMergeVariables(config *ship.MCPServerConfig, userVars map[string]string) error {
	if config.Variables == nil {
		return nil
	}
	
	// Check for required variables
	for _, variable := range config.Variables {
		if variable.Required {
			// Check if provided by user
			if _, exists := userVars[variable.Name]; !exists {
				// Check if has default value
				if variable.Default == "" {
					return fmt.Errorf("required variable %s is missing (use --var %s=value)", 
						variable.Name, variable.Name)
				}
			}
		}
	}
	
	// Merge variables into config.Env
	if config.Env == nil {
		config.Env = make(map[string]string)
	}
	
	// First, set defaults for variables that aren't provided
	for _, variable := range config.Variables {
		if _, exists := userVars[variable.Name]; !exists && variable.Default != "" {
			config.Env[variable.Name] = variable.Default
		}
	}
	
	// Then, override with user-provided values
	for key, value := range userVars {
		config.Env[key] = value
	}
	
	return nil
}

// setContainerEnvironmentVars sets environment variables for containerized tools
func setContainerEnvironmentVars(envVars map[string]string) {
	for key, value := range envVars {
		os.Setenv(key, value)
	}
}

// showVariableHelp displays information about available variables for a tool
func showVariableHelp(serverName string) {
	config, exists := hardcodedMCPServers[serverName]
	if !exists || len(config.Variables) == 0 {
		return
	}
	
	fmt.Fprintf(os.Stderr, "\nAvailable variables for %s:\n", serverName)
	for _, variable := range config.Variables {
		required := ""
		if variable.Required {
			required = " (required)"
		}
		
		secret := ""
		if variable.Secret {
			secret = " (secret)"
		}
		
		defaultInfo := ""
		if variable.Default != "" {
			defaultInfo = fmt.Sprintf(" [default: %s]", variable.Default)
		}
		
		fmt.Fprintf(os.Stderr, "  --var %s=value%s%s%s\n    %s\n", 
			variable.Name, defaultInfo, required, secret, variable.Description)
	}
	fmt.Fprintf(os.Stderr, "\n")
}
