package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitleaksTools adds Gitleaks (secret detection in code and git history) MCP tool implementations using real gitleaks CLI commands
func AddGitleaksTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Gitleaks scan git repository for secrets
	gitScanTool := mcp.NewTool("gitleaks_git",
		mcp.WithDescription("Scan git repositories for secrets"),
		mcp.WithString("path",
			mcp.Description("Path to git repository (default: current directory)"),
		),
		mcp.WithString("config",
			mcp.Description("Config file path"),
		),
		mcp.WithString("baseline_path",
			mcp.Description("Path to baseline with issues that can be ignored"),
		),
		mcp.WithString("report_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "csv", "junit", "sarif", "template"),
		),
		mcp.WithString("report_path",
			mcp.Description("Report file destination"),
		),
		mcp.WithString("log_opts",
			mcp.Description("Git log options for commit range (e.g., --all commitA..commitB)"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Show verbose output"),
		),
		mcp.WithString("exit_code",
			mcp.Description("Exit code when leaks are found (default 1)"),
		),
	)
	s.AddTool(gitScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"gitleaks", "git"}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "-c", config)
		}
		if baselinePath := request.GetString("baseline_path", ""); baselinePath != "" {
			args = append(args, "-b", baselinePath)
		}
		if reportFormat := request.GetString("report_format", ""); reportFormat != "" {
			args = append(args, "-f", reportFormat)
		}
		if reportPath := request.GetString("report_path", ""); reportPath != "" {
			args = append(args, "-r", reportPath)
		}
		if logOpts := request.GetString("log_opts", ""); logOpts != "" {
			args = append(args, "--log-opts", logOpts)
		}
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		if exitCode := request.GetString("exit_code", ""); exitCode != "" {
			args = append(args, "--exit-code", exitCode)
		}
		if path := request.GetString("path", ""); path != "" {
			args = append(args, path)
		}
		
		return executeShipCommand(args)
	})

	// Gitleaks scan directories or files for secrets
	dirScanTool := mcp.NewTool("gitleaks_dir",
		mcp.WithDescription("Scan directories or files for secrets"),
		mcp.WithString("path",
			mcp.Description("Path to directory or file to scan (default: current directory)"),
		),
		mcp.WithString("config",
			mcp.Description("Config file path"),
		),
		mcp.WithString("baseline_path",
			mcp.Description("Path to baseline with issues that can be ignored"),
		),
		mcp.WithString("report_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "csv", "junit", "sarif", "template"),
		),
		mcp.WithString("report_path",
			mcp.Description("Report file destination"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Show verbose output"),
		),
		mcp.WithString("exit_code",
			mcp.Description("Exit code when leaks are found (default 1)"),
		),
	)
	s.AddTool(dirScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"gitleaks", "dir"}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "-c", config)
		}
		if baselinePath := request.GetString("baseline_path", ""); baselinePath != "" {
			args = append(args, "-b", baselinePath)
		}
		if reportFormat := request.GetString("report_format", ""); reportFormat != "" {
			args = append(args, "-f", reportFormat)
		}
		if reportPath := request.GetString("report_path", ""); reportPath != "" {
			args = append(args, "-r", reportPath)
		}
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		if exitCode := request.GetString("exit_code", ""); exitCode != "" {
			args = append(args, "--exit-code", exitCode)
		}
		if path := request.GetString("path", ""); path != "" {
			args = append(args, path)
		}
		
		return executeShipCommand(args)
	})

	// Gitleaks detect secrets from stdin
	stdinScanTool := mcp.NewTool("gitleaks_stdin",
		mcp.WithDescription("Detect secrets from standard input"),
		mcp.WithString("config",
			mcp.Description("Config file path"),
		),
		mcp.WithString("baseline_path",
			mcp.Description("Path to baseline with issues that can be ignored"),
		),
		mcp.WithString("report_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "csv", "junit", "sarif", "template"),
		),
		mcp.WithString("report_path",
			mcp.Description("Report file destination"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Show verbose output"),
		),
		mcp.WithString("exit_code",
			mcp.Description("Exit code when leaks are found (default 1)"),
		),
	)
	s.AddTool(stdinScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"gitleaks", "stdin"}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "-c", config)
		}
		if baselinePath := request.GetString("baseline_path", ""); baselinePath != "" {
			args = append(args, "-b", baselinePath)
		}
		if reportFormat := request.GetString("report_format", ""); reportFormat != "" {
			args = append(args, "-f", reportFormat)
		}
		if reportPath := request.GetString("report_path", ""); reportPath != "" {
			args = append(args, "-r", reportPath)
		}
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		if exitCode := request.GetString("exit_code", ""); exitCode != "" {
			args = append(args, "--exit-code", exitCode)
		}
		
		return executeShipCommand(args)
	})

	// Gitleaks version information
	versionTool := mcp.NewTool("gitleaks_version",
		mcp.WithDescription("Display Gitleaks version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"gitleaks", "version"}
		return executeShipCommand(args)
	})
}