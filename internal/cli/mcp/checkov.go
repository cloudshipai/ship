package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCheckovTools adds Checkov (Infrastructure as Code static analysis) MCP tool implementations
func AddCheckovTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		directory := request.GetString("directory", "")
		args := []string{"checkov", "--directory", directory}
		
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "--framework", framework)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if request.GetBool("compact", false) {
			args = append(args, "--compact")
		}
		if request.GetBool("quiet", false) {
			args = append(args, "--quiet")
		}
		
		return executeShipCommand(args)
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
		file := request.GetString("file", "")
		args := []string{"checkov", "--file", file}
		
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "--framework", framework)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		args := []string{"checkov", "--directory", directory}
		
		if check := request.GetString("check", ""); check != "" {
			args = append(args, "--check", check)
		}
		if skipCheck := request.GetString("skip_check", ""); skipCheck != "" {
			args = append(args, "--skip-check", skipCheck)
		}
		
		return executeShipCommand(args)
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
		dockerImage := request.GetString("docker_image", "")
		args := []string{"checkov", "--docker-image", dockerImage, "--framework", "sca_image"}
		
		if dockerfilePath := request.GetString("dockerfile_path", ""); dockerfilePath != "" {
			args = append(args, "--dockerfile-path", dockerfilePath)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		args := []string{"checkov", "--directory", directory, "--framework", "sca_package"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		args := []string{"checkov", "--directory", directory, "--framework", "secrets"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		configFile := request.GetString("config_file", "")
		args := []string{"checkov", "--directory", directory, "--config-file", configFile}
		
		return executeShipCommand(args)
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
		configPath := request.GetString("config_path", "")
		args := []string{"checkov", "--create-config", configPath}
		
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		args := []string{"checkov", "--directory", directory}
		
		if request.GetBool("download_external_modules", false) {
			args = append(args, "--download-external-modules", "true")
		}
		
		return executeShipCommand(args)
	})

	// Checkov get version tool
	getVersionTool := mcp.NewTool("checkov_get_version",
		mcp.WithDescription("Get Checkov version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"checkov", "--version"}
		return executeShipCommand(args)
	})
}