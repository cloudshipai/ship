// MCP Server implementation for Ship CLI

package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	shipMcp "github.com/cloudshipai/ship/internal/cli/mcp"
	"github.com/cloudshipai/ship/internal/telemetry"
	"github.com/cloudshipai/ship/pkg/ship"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

// ExecutionContext holds options for enhanced MCP execution
type ExecutionContext struct {
	OutputFile    string
	ExecutionLog  string
}

// Global execution context
var globalExecutionContext *ExecutionContext

var mcpCmd = &cobra.Command{
	Use:   "mcp [tool]",
	Short: "Start MCP server for a specific tool or all tools",
	Long: shipMcp.GenerateMCPHelpText(),
	Args: cobra.MaximumNArgs(1),
	RunE: runMCPServer,
}

func init() {
	rootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().Int("port", 0, "Port to listen on (0 for stdio)")
	mcpCmd.Flags().String("host", "localhost", "Host to bind to")
	mcpCmd.Flags().Bool("stdio", true, "Use stdio transport (default)")
	mcpCmd.Flags().StringToString("var", nil, "Environment variables for MCP servers and containers (e.g., --var API_KEY=value --var DEBUG=true)")
	mcpCmd.Flags().StringToString("image-tag", nil, "Override container image tags for tools (e.g., --image-tag trivy=aquasec/trivy:0.50.0 --image-tag checkov=bridgecrew/checkov:3.2.0)")
	mcpCmd.Flags().Bool("version", false, "Show version information for tools")
	mcpCmd.Flags().String("output-file", "", "Write tool output to file (in addition to MCP response)")
	mcpCmd.Flags().String("execution-log", "", "Write execution logs and timing to file")
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	host, _ := cmd.Flags().GetString("host")
	useStdio, _ := cmd.Flags().GetBool("stdio")
	envVars, _ := cmd.Flags().GetStringToString("var")
	imageTags, _ := cmd.Flags().GetStringToString("image-tag")
	showVersion, _ := cmd.Flags().GetBool("version")
	outputFile, _ := cmd.Flags().GetString("output-file")
	executionLog, _ := cmd.Flags().GetString("execution-log")

	// Set global execution context for output options
	globalExecutionContext = &ExecutionContext{
		OutputFile:   outputFile,
		ExecutionLog: executionLog,
	}

	// Set environment variables for Dagger tools to access output options
	if outputFile != "" {
		os.Setenv("SHIP_OUTPUT_FILE", outputFile)
	}
	if executionLog != "" {
		os.Setenv("SHIP_EXECUTION_LOG", executionLog)
	}

	// Handle version requests
	if showVersion {
		return handleVersionRequest(args)
	}

	// Determine which tool to serve
	toolName := "all"
	if len(args) > 0 {
		toolName = args[0]
	}

	// Track MCP command usage
	telemetry.TrackMCPCommand(toolName)

	// Create MCP server with enhanced configuration
	serverName := fmt.Sprintf("ship-%s", toolName)
	s := server.NewMCPServer(serverName, "1.0.0")

	// Set environment variables for containerized tools
	if len(envVars) > 0 {
		setContainerEnvironmentVars(envVars)
	}

	// Set image tag overrides as environment variables
	if len(imageTags) > 0 {
		setImageTagOverrides(imageTags)
	}

	// Add specific tools based on argument using the modular registry
	switch toolName {
	case "all":
		// Register all tools from all categories with enhanced execution wrapper
		shipMcp.RegisterAllTools(s, executeShipCommandWithStabilityEnhancements)
	case "terraform", "security", "aws", "kubernetes", "cloud", "supply-chain":
		// Register tools by category with enhanced execution wrapper
		shipMcp.RegisterToolsByCategory(toolName, s, executeShipCommandWithStabilityEnhancements)
	default:
		// Check if this is a specific tool name with enhanced execution wrapper
		shipMcp.RegisterToolByName(toolName, s, executeShipCommandWithStabilityEnhancements)
		
		// If no tools were registered, check if it's an external MCP server  
		// Note: we can't check s.Tools() as it's not exposed, so we'll assume tool was registered
		// If the tool doesn't exist, the registry function will handle it gracefully
		// External MCP servers are handled by checking the modular external servers
		if shipMcp.IsExternalMCPServer(toolName) {
			return runMCPProxy(cmd, toolName)
		}
	}

	// Add resources for documentation and help
	addResources(s)

	// Add prompts only for 'all' mode
	if toolName == "all" {
		addPrompts(s)
	}

	// Start server with enhanced stability
	if useStdio || port == 0 {
		fmt.Fprintf(os.Stderr, "Starting %s MCP server on stdio with stability enhancements...\n", serverName)
		return serveStdioWithStability(s)
	} else {
		fmt.Fprintf(os.Stderr, "Starting %s MCP server on %s:%d...\n", serverName, host, port)
		return fmt.Errorf("HTTP server not implemented in this version, use --stdio")
	}
}

// Investigation tools removed to focus on Terraform analysis workflows

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
		content := shipMcp.GenerateToolsResourceContent()

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

// executeShipCommandWithLogging executes ship commands with enhanced logging and file output
func executeShipCommandWithLogging(args []string) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	
	// Execute the command
	result, err := executeShipCommand(args)
	
	elapsed := time.Since(startTime)
	
	// Write execution log if requested
	if globalExecutionContext != nil && globalExecutionContext.ExecutionLog != "" {
		logEntry := fmt.Sprintf("[%s] Command: ship %s | Duration: %v | Success: %t\n",
			time.Now().Format("2006-01-02 15:04:05"),
			strings.Join(args, " "),
			elapsed,
			err == nil && (result == nil || !result.IsError))
		
		// Append to execution log file
		logFile, logErr := os.OpenFile(globalExecutionContext.ExecutionLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if logErr == nil {
			logFile.WriteString(logEntry)
			logFile.Close()
		}
	}
	
	// Write output to file if requested and we have content
	if globalExecutionContext != nil && globalExecutionContext.OutputFile != "" && result != nil && !result.IsError {
		var outputText string
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(mcp.TextContent); ok {
				outputText = textContent.Text
			}
		}
		
		if outputText != "" {
			// Create timestamped output for multiple invocations
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			separator := strings.Repeat("=", 80)
			outputContent := fmt.Sprintf("\n%s\n=== Ship CLI Output - %s ===\nCommand: ship %s\nDuration: %v\n%s\n\n%s\n",
				separator,
				timestamp,
				strings.Join(args, " "),
				elapsed,
				separator,
				outputText)
			
			// Append to output file
			outputFile, fileErr := os.OpenFile(globalExecutionContext.OutputFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if fileErr == nil {
				outputFile.WriteString(outputContent)
				outputFile.Close()
			}
		}
	}
	
	return result, err
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

	outputStr := string(output)

	// Check if output needs to be chunked
	if needsChunking(outputStr) {
		return createChunkedResponse(outputStr), nil
	}

	return mcp.NewToolResultText(outputStr), nil
}

func handleVersionRequest(args []string) error {
	toolName := "all"
	if len(args) > 0 {
		toolName = args[0]
	}

	fmt.Printf("Ship CLI Version Information for: %s\n", toolName)
	fmt.Println("=====================================")

	switch toolName {
	case "all":
		fmt.Println("Displaying version information for all available tools...")
		// Get versions for key containerized tools
		tools := []string{"checkov", "trivy", "tflint", "terragrunt", "terraform", "ansible"}
		for _, tool := range tools {
			if output, err := getToolVersion(tool); err == nil {
				fmt.Printf("\n%s:\n%s\n", strings.ToUpper(tool), strings.TrimSpace(output))
			} else {
				fmt.Printf("\n%s: Version information unavailable (%v)\n", strings.ToUpper(tool), err)
			}
		}
	default:
		// Get version for specific tool
		if output, err := getToolVersion(toolName); err == nil {
			fmt.Printf("\n%s:\n%s\n", strings.ToUpper(toolName), strings.TrimSpace(output))
		} else {
			return fmt.Errorf("failed to get version for %s: %v", toolName, err)
		}
	}

	return nil
}

func getToolVersion(toolName string) (string, error) {
	// Use the existing executeShipCommand to get version information
	args := []string{toolName, "--version"}
	result, err := executeShipCommand(args)
	if err != nil {
		return "", err
	}

	if result.IsError {
		return "", fmt.Errorf("error getting version for %s", toolName)
	}

	if result != nil && len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			return textContent.Text, nil
		}
	}

	return "", fmt.Errorf("no version information available")
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


// runMCPProxy starts an MCP proxy server for external MCP servers
func runMCPProxy(cmd *cobra.Command, serverName string) error {
	useStdio, _ := cmd.Flags().GetBool("stdio")
	port, _ := cmd.Flags().GetInt("port")
	envVars, _ := cmd.Flags().GetStringToString("var")

	// Get external server configuration
	mcpConfig, exists := shipMcp.GetExternalMCPServer(serverName)
	if !exists {
		return fmt.Errorf("external MCP server '%s' not found in configurations", serverName)
	}

	// Validate and merge environment variables
	if err := validateAndMergeVariables(&mcpConfig, envVars); err != nil {
		showVariableHelp(serverName)
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

// setImageTagOverrides sets environment variables for container image tag overrides
func setImageTagOverrides(imageTags map[string]string) {
	for toolName, imageTag := range imageTags {
		// Convert tool name to uppercase and set as environment variable
		// Format: SHIP_IMAGE_TAG_<TOOLNAME>=<image:tag>
		envKey := fmt.Sprintf("SHIP_IMAGE_TAG_%s", strings.ToUpper(toolName))
		os.Setenv(envKey, imageTag)
	}
}

// showVariableHelp displays information about available variables for a tool
func showVariableHelp(serverName string) {
	config, exists := shipMcp.GetExternalMCPServer(serverName)
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

// executeShipCommandWithStabilityEnhancements wraps command execution with connection stability improvements
func executeShipCommandWithStabilityEnhancements(args []string) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	
	// Create a context with extended timeout for long-running operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	
	// Execute the command with context
	result, err := executeShipCommandWithContext(ctx, args)
	
	elapsed := time.Since(startTime)
	
	// Enhanced logging for debugging transport issues
	if globalExecutionContext != nil && globalExecutionContext.ExecutionLog != "" {
		status := "success"
		if err != nil {
			status = fmt.Sprintf("error: %v", err)
		} else if result != nil && result.IsError {
			status = "tool_error"
		}
		
		logEntry := fmt.Sprintf("[%s] Command: ship %s | Duration: %v | Status: %s | PID: %d\n",
			time.Now().Format("2006-01-02 15:04:05"),
			strings.Join(args, " "),
			elapsed,
			status,
			os.Getpid())
		
		// Write to execution log with error handling
		if logFile, logErr := os.OpenFile(globalExecutionContext.ExecutionLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); logErr == nil {
			logFile.WriteString(logEntry)
			logFile.Close()
		}
	}
	
	// Enhanced file output with stability markers
	if globalExecutionContext != nil && globalExecutionContext.OutputFile != "" && result != nil && !result.IsError {
		var outputText string
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(mcp.TextContent); ok {
				outputText = textContent.Text
			}
		}
		
		if outputText != "" {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			separator := strings.Repeat("=", 80)
			outputContent := fmt.Sprintf("\n%s\n=== Ship CLI Output - %s ===\nCommand: ship %s\nDuration: %v\nPID: %d\nMCP Transport: STABLE\n%s\n\n%s\n",
				separator,
				timestamp,
				strings.Join(args, " "),
				elapsed,
				os.Getpid(),
				separator,
				outputText)
			
			// Write to output file with error handling
			if outputFile, fileErr := os.OpenFile(globalExecutionContext.OutputFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); fileErr == nil {
				outputFile.WriteString(outputContent)
				outputFile.Close()
			}
		}
	}
	
	return result, err
}

// executeShipCommandWithContext executes ship commands with context cancellation support
func executeShipCommandWithContext(ctx context.Context, args []string) (*mcp.CallToolResult, error) {
	// Get the current binary path
	executable, err := os.Executable()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get executable path: %v", err)), nil
	}

	// Create command with context
	cmd := exec.CommandContext(ctx, executable, args...)
	
	// Set up process group to handle cleanup properly
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	
	// Execute the command with context timeout
	output, err := cmd.CombinedOutput()

	// Handle context timeout specifically
	if ctx.Err() == context.DeadlineExceeded {
		return mcp.NewToolResultError("Command timed out after 10 minutes. This may indicate a hung container or network issue."), nil
	}

	if err != nil {
		// Enhanced error reporting for transport debugging
		return mcp.NewToolResultError(fmt.Sprintf("Command failed: %s\n\nOutput:\n%s\n\nContext: %v", 
			err.Error(), string(output), ctx.Err())), nil
	}

	outputStr := string(output)

	// Check if output needs to be chunked
	if needsChunking(outputStr) {
		return createChunkedResponse(outputStr), nil
	}

	return mcp.NewToolResultText(outputStr), nil
}

// serveStdioWithStability wraps the stdio server with connection stability improvements
func serveStdioWithStability(s *server.MCPServer) error {
	// Add signal handling to detect broken pipes early
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Monitor for stdio connection health
	go func() {
		// Simple keepalive mechanism - write periodic status to stderr
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				// Write keepalive message to stderr (not interfering with stdio protocol)
				fmt.Fprintf(os.Stderr, "[MCP-KEEPALIVE] Server running, PID: %d\n", os.Getpid())
			case <-ctx.Done():
				return
			}
		}
	}()
	
	// Enhanced error handling for stdio transport
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "[MCP-PANIC] Server panic recovered: %v\n", r)
		}
	}()
	
	// Start the standard stdio server with enhanced monitoring
	err := server.ServeStdio(s)
	
	// Enhanced error reporting
	if err != nil {
		fmt.Fprintf(os.Stderr, "[MCP-ERROR] Server error: %v\n", err)
	}
	
	return err
}

