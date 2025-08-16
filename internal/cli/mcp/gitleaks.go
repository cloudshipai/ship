package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitleaksTools adds Gitleaks MCP tool implementations
func AddGitleaksTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Gitleaks scan directory tool
	scanDirTool := mcp.NewTool("gitleaks_scan_directory",
		mcp.WithDescription("Scan a directory for secrets using Gitleaks"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
	)
	s.AddTool(scanDirTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gitleaks"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// Gitleaks scan file tool
	scanFileTool := mcp.NewTool("gitleaks_scan_file",
		mcp.WithDescription("Scan a specific file for secrets using Gitleaks"),
		mcp.WithString("file_path",
			mcp.Description("Path to the file to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"security", "gitleaks", "--file", filePath}
		return executeShipCommand(args)
	})

	// Gitleaks scan git repository tool
	scanGitTool := mcp.NewTool("gitleaks_scan_git_repo",
		mcp.WithDescription("Scan a git repository for secrets using Gitleaks"),
		mcp.WithString("repository",
			mcp.Description("Path to git repository (default: current directory)"),
		),
	)
	s.AddTool(scanGitTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gitleaks"}
		if repo := request.GetString("repository", ""); repo != "" {
			args = append(args, repo)
		}
		args = append(args, "--git")
		return executeShipCommand(args)
	})

	// Gitleaks scan with config tool
	scanConfigTool := mcp.NewTool("gitleaks_scan_with_config",
		mcp.WithDescription("Scan using custom Gitleaks configuration"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("config",
			mcp.Description("Path to Gitleaks configuration file"),
			mcp.Required(),
		),
	)
	s.AddTool(scanConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gitleaks"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		config := request.GetString("config", "")
		args = append(args, "--config", config)
		return executeShipCommand(args)
	})

	// Gitleaks scan from stdin tool
	scanStdinTool := mcp.NewTool("gitleaks_scan_stdin",
		mcp.WithDescription("Scan input from standard input using Gitleaks"),
		mcp.WithString("input_data",
			mcp.Description("Data to scan for secrets"),
			mcp.Required(),
		),
		mcp.WithString("report_format",
			mcp.Description("Output format for the report"),
			mcp.Enum("json", "csv", "junit", "sarif"),
		),
	)
	s.AddTool(scanStdinTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gitleaks", "stdin"}
		if format := request.GetString("report_format", ""); format != "" {
			args = append(args, "--report-format", format)
		}
		return executeShipCommand(args)
	})

	// Gitleaks scan with baseline tool
	scanWithBaselineTool := mcp.NewTool("gitleaks_scan_with_baseline",
		mcp.WithDescription("Scan with baseline file to ignore known issues"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("baseline_path",
			mcp.Description("Path to baseline file with ignorable issues"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithBaselineTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gitleaks"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		baseline := request.GetString("baseline_path", "")
		args = append(args, "--baseline-path", baseline)
		return executeShipCommand(args)
	})

	// Gitleaks scan with rules tool
	scanWithRulesTool := mcp.NewTool("gitleaks_scan_with_rules",
		mcp.WithDescription("Scan with specific rules enabled"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("enable_rule",
			mcp.Description("Rule ID to enable specifically"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gitleaks"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		rule := request.GetString("enable_rule", "")
		args = append(args, "--enable-rule", rule)
		return executeShipCommand(args)
	})
}