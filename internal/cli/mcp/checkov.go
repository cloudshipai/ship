package mcp

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCheckovTools adds Checkov (Infrastructure as Code static analysis) MCP tool implementations
func AddCheckovTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addCheckovToolsDirect(s)
}

// logDaggerOutput logs tool output and execution details to files if configured
func logDaggerOutput(toolName string, args []string, result string, duration time.Duration) {
	// Get output file paths from environment variables (set by CLI flags)
	outputFile := os.Getenv("SHIP_OUTPUT_FILE")
	executionLog := os.Getenv("SHIP_EXECUTION_LOG")
	
	// Write execution log if requested
	if executionLog != "" {
		logEntry := fmt.Sprintf("[%s] Tool: %s | Args: %s | Duration: %v | Success: true\n",
			time.Now().Format("2006-01-02 15:04:05"),
			toolName,
			strings.Join(args, " "),
			duration)
		
		// Append to execution log file
		logFile, logErr := os.OpenFile(executionLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if logErr == nil {
			logFile.WriteString(logEntry)
			logFile.Close()
		}
	}
	
	// Write output to file if requested
	if outputFile != "" && result != "" {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		separator := strings.Repeat("=", 80)
		outputContent := fmt.Sprintf("\n%s\n=== Ship CLI Dagger Output - %s ===\nTool: %s\nArgs: %s\nDuration: %v\n%s\n\n%s\n",
			separator,
			timestamp,
			toolName,
			strings.Join(args, " "),
			duration,
			separator,
			result)
		
		// Append to output file
		outputFileHandle, fileErr := os.OpenFile(outputFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if fileErr == nil {
			outputFileHandle.WriteString(outputContent)
			outputFileHandle.Close()
		}
	}
}

// addCheckovToolsDirect implements direct Dagger calls for Checkov tools
func addCheckovToolsDirect(s *server.MCPServer) {
	// Checkov scan directory tool
	scanDirectoryTool := mcp.NewTool("checkov_scan_directory",
		mcp.WithDescription("Scan directory for security issues using Checkov"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("framework",
			mcp.Description("Framework to scan (terraform, cloudformation, kubernetes, etc.)"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("cli", "json", "junitxml", "github_failed_only", "sarif", "csv"),
		),
		mcp.WithBoolean("compact",
			mcp.Description("Compact output format"),
		),
		mcp.WithBoolean("quiet",
			mcp.Description("Reduce verbosity"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		startTime := time.Now()
		
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		directory := request.GetString("directory", "")
		framework := request.GetString("framework", "")
		output := request.GetString("output", "")
		compact := request.GetBool("compact", false)
		quiet := request.GetBool("quiet", false)

		// Create Checkov module and scan directory
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.ScanDirectoryWithOptions(ctx, directory, framework, output, compact, quiet)
		
		elapsed := time.Since(startTime)
		
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov directory scan failed: %v", err)), nil
		}

		// Log output and execution details
		args := []string{
			fmt.Sprintf("directory=%s", directory),
			fmt.Sprintf("framework=%s", framework),
			fmt.Sprintf("output=%s", output),
			fmt.Sprintf("compact=%t", compact),
			fmt.Sprintf("quiet=%t", quiet),
		}
		logDaggerOutput("checkov_scan_directory", args, result, elapsed)

		return mcp.NewToolResultText(result), nil
	})

	// Checkov scan file tool
	scanFileTool := mcp.NewTool("checkov_scan_file",
		mcp.WithDescription("Scan specific file(s) for security issues using Checkov"),
		mcp.WithString("file",
			mcp.Description("Path to file to scan"),
			mcp.Required(),
		),
		mcp.WithString("framework",
			mcp.Description("Framework type for the file"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("cli", "json", "junitxml", "sarif"),
		),
	)
	s.AddTool(scanFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		file := request.GetString("file", "")
		framework := request.GetString("framework", "")
		output := request.GetString("output", "")

		// Create Checkov module and scan file
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.ScanFileWithOptions(ctx, file, framework, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov file scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Checkov scan with specific checks tool
	scanWithChecksTool := mcp.NewTool("checkov_scan_with_checks",
		mcp.WithDescription("Scan with specific checks enabled or disabled"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("check",
			mcp.Description("Comma-separated list of check IDs to run"),
		),
		mcp.WithString("skip_check",
			mcp.Description("Comma-separated list of check IDs to skip"),
		),
	)
	s.AddTool(scanWithChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		directory := request.GetString("directory", "")
		check := request.GetString("check", "")
		skipCheck := request.GetString("skip_check", "")

		// Create Checkov module and scan with specific checks
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.ScanWithSpecificChecks(ctx, directory, check, skipCheck)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov scan with checks failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Checkov scan Docker image tool
	scanDockerImageTool := mcp.NewTool("checkov_scan_docker_image",
		mcp.WithDescription("Scan Docker container image for vulnerabilities"),
		mcp.WithString("docker_image",
			mcp.Description("Docker image to scan (name:tag)"),
			mcp.Required(),
		),
		mcp.WithString("dockerfile_path",
			mcp.Description("Path to Dockerfile (optional)"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("cli", "json", "sarif"),
		),
	)
	s.AddTool(scanDockerImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		dockerImage := request.GetString("docker_image", "")
		dockerfilePath := request.GetString("dockerfile_path", "")
		output := request.GetString("output", "")

		// Create Checkov module and scan Docker image
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.ScanDockerImage(ctx, dockerImage, dockerfilePath, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov Docker image scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Checkov scan packages tool
	scanPackagesTool := mcp.NewTool("checkov_scan_packages",
		mcp.WithDescription("Scan package dependencies for vulnerabilities"),
		mcp.WithString("directory",
			mcp.Description("Directory containing package files"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("cli", "json", "sarif"),
		),
	)
	s.AddTool(scanPackagesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		directory := request.GetString("directory", "")
		output := request.GetString("output", "")

		// Create Checkov module and scan packages
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.ScanPackages(ctx, directory, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov package scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Checkov scan secrets tool
	scanSecretsTool := mcp.NewTool("checkov_scan_secrets",
		mcp.WithDescription("Scan for hardcoded secrets in code"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan for secrets"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("cli", "json", "sarif"),
		),
	)
	s.AddTool(scanSecretsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		directory := request.GetString("directory", "")
		output := request.GetString("output", "")

		// Create Checkov module and scan secrets
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.ScanSecrets(ctx, directory, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov secrets scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Checkov scan with config file tool
	scanWithConfigTool := mcp.NewTool("checkov_scan_with_config",
		mcp.WithDescription("Scan using configuration file"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to configuration file"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		directory := request.GetString("directory", "")
		configFile := request.GetString("config_file", "")

		// Create Checkov module and scan with config
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.ScanWithConfig(ctx, directory, configFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov scan with config failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Checkov create config tool
	createConfigTool := mcp.NewTool("checkov_create_config",
		mcp.WithDescription("Generate configuration file from current settings"),
		mcp.WithString("config_path",
			mcp.Description("Path where config file should be created"),
			mcp.Required(),
		),
	)
	s.AddTool(createConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		configPath := request.GetString("config_path", "")

		// Create Checkov module and create config
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.CreateConfig(ctx, configPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov create config failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Checkov download external modules tool
	downloadModulesTool := mcp.NewTool("checkov_download_external_modules",
		mcp.WithDescription("Scan with external module downloading enabled"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithBoolean("download_external_modules",
			mcp.Description("Download external terraform modules"),
		),
	)
	s.AddTool(downloadModulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		directory := request.GetString("directory", "")
		downloadExternalModules := request.GetBool("download_external_modules", false)

		// Create Checkov module and scan with external modules
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.ScanWithExternalModules(ctx, directory, downloadExternalModules)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov scan with external modules failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Checkov get version tool
	getVersionTool := mcp.NewTool("checkov_get_version",
		mcp.WithDescription("Get Checkov version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create Checkov module and get version
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("checkov get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}