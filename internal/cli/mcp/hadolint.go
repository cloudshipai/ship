package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddHadolintTools adds Hadolint (Dockerfile linter) MCP tool implementations using direct Dagger calls
func AddHadolintTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addHadolintToolsDirect(s)
}

// addHadolintToolsDirect adds Hadolint tools using direct Dagger module calls
func addHadolintToolsDirect(s *server.MCPServer) {
	// Hadolint scan Dockerfile tool
	scanDockerfileTool := mcp.NewTool("hadolint_scan_dockerfile",
		mcp.WithDescription("Scan Dockerfile for best practices and security issues"),
		mcp.WithString("dockerfile_path",
			mcp.Description("Path to Dockerfile to scan"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("tty", "json", "checkstyle", "codeclimate", "gitlab_codeclimate", "gnu", "codacy", "sonarqube", "sarif"),
		),
	)
	s.AddTool(scanDockerfileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHadolintModule(client)

		// Get Dockerfile path
		dockerfilePath := request.GetString("dockerfile_path", "")
		if dockerfilePath == "" {
			return mcp.NewToolResultError("dockerfile_path is required"), nil
		}

		// Scan Dockerfile
		output, err := module.ScanDockerfile(ctx, dockerfilePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to scan Dockerfile: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Hadolint scan directory tool
	scanDirectoryTool := mcp.NewTool("hadolint_scan_directory",
		mcp.WithDescription("Scan directory for Dockerfiles"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan for Dockerfiles"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("tty", "json", "checkstyle", "codeclimate", "gitlab_codeclimate", "gnu", "codacy", "sonarqube", "sarif"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHadolintModule(client)

		// Get directory
		directory := request.GetString("directory", "")
		if directory == "" {
			return mcp.NewToolResultError("directory is required"), nil
		}

		// Scan directory
		output, err := module.ScanDirectory(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to scan directory: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Hadolint scan with config tool
	scanWithConfigTool := mcp.NewTool("hadolint_scan_with_config",
		mcp.WithDescription("Scan Dockerfile with custom Hadolint configuration"),
		mcp.WithString("dockerfile_path",
			mcp.Description("Path to Dockerfile"),
			mcp.Required(),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to Hadolint configuration file"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("tty", "json", "checkstyle", "codeclimate", "gitlab_codeclimate", "gnu", "codacy", "sonarqube", "sarif"),
		),
	)
	s.AddTool(scanWithConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHadolintModule(client)

		// Get parameters
		dockerfilePath := request.GetString("dockerfile_path", "")
		if dockerfilePath == "" {
			return mcp.NewToolResultError("dockerfile_path is required"), nil
		}
		
		configFile := request.GetString("config_file", "")
		if configFile == "" {
			return mcp.NewToolResultError("config_file is required"), nil
		}

		// Scan with config
		output, err := module.ScanWithConfig(ctx, dockerfilePath, configFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to scan with config: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Hadolint scan with ignore rules tool
	scanIgnoreRulesTool := mcp.NewTool("hadolint_scan_ignore_rules",
		mcp.WithDescription("Scan Dockerfile while ignoring specific rules"),
		mcp.WithString("dockerfile_path",
			mcp.Description("Path to Dockerfile to scan"),
			mcp.Required(),
		),
		mcp.WithString("ignore_rules",
			mcp.Description("Comma-separated list of rules to ignore (e.g., DL3003,DL3006)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("tty", "json", "checkstyle", "codeclimate", "gitlab_codeclimate", "gnu", "codacy", "sonarqube", "sarif"),
		),
	)
	s.AddTool(scanIgnoreRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHadolintModule(client)

		// Get parameters
		dockerfilePath := request.GetString("dockerfile_path", "")
		if dockerfilePath == "" {
			return mcp.NewToolResultError("dockerfile_path is required"), nil
		}

		// Parse ignore rules
		var ignoreRules []string
		if ignoreRulesStr := request.GetString("ignore_rules", ""); ignoreRulesStr != "" {
			rules := strings.Split(ignoreRulesStr, ",")
			for _, rule := range rules {
				ignoreRules = append(ignoreRules, strings.TrimSpace(rule))
			}
		}

		// Scan with ignore rules
		output, err := module.ScanIgnoreRules(ctx, dockerfilePath, ignoreRules)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to scan with ignore rules: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Hadolint get version tool
	getVersionTool := mcp.NewTool("hadolint_get_version",
		mcp.WithDescription("Get Hadolint version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHadolintModule(client)

		// Get version
		version, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get version: %v", err)), nil
		}

		return mcp.NewToolResultText(version), nil
	})
}