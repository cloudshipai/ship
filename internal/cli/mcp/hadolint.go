package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddHadolintTools adds Hadolint (Dockerfile linter) MCP tool implementations
func AddHadolintTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		dockerfilePath := request.GetString("dockerfile_path", "")
		args := []string{"hadolint", dockerfilePath}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		// Hadolint doesn't support directory scanning directly - scan individual Dockerfiles
		args := []string{"sh", "-c", "find " + directory + " -name 'Dockerfile*' -exec hadolint {} +"}
		if format := request.GetString("format", ""); format != "" {
			// Add format to find command
			args = []string{"sh", "-c", "find " + directory + " -name 'Dockerfile*' -exec hadolint --format " + format + " {} +"}
		}
		return executeShipCommand(args)
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
		dockerfilePath := request.GetString("dockerfile_path", "")
		configFile := request.GetString("config_file", "")
		args := []string{"hadolint", "--config", configFile, dockerfilePath}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
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
		dockerfilePath := request.GetString("dockerfile_path", "")
		args := []string{"hadolint", dockerfilePath}
		
		if ignoreRules := request.GetString("ignore_rules", ""); ignoreRules != "" {
			// Split comma-separated rules and add each as --ignore flag
			rules := strings.Split(ignoreRules, ",")
			for _, rule := range rules {
				args = append(args, "--ignore", strings.TrimSpace(rule))
			}
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		
		return executeShipCommand(args)
	})

	// Hadolint get version tool
	getVersionTool := mcp.NewTool("hadolint_get_version",
		mcp.WithDescription("Get Hadolint version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"hadolint", "--version"}
		return executeShipCommand(args)
	})
}