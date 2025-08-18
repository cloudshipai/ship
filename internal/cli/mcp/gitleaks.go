package mcp

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddGitleaksTools adds Gitleaks (secret detection in code and git history) MCP tool implementations using direct Dagger calls
func AddGitleaksTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGitleaksToolsDirect(s)
}

// addGitleaksToolsDirect adds Gitleaks tools using direct Dagger module calls
func addGitleaksToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGitleaksModule(client)

		// Get path
		path := request.GetString("path", ".")
		if path == "" {
			path = "."
		}

		// If config is specified, use ScanWithConfig
		if config := request.GetString("config", ""); config != "" {
			output, err := module.ScanWithConfig(ctx, path, config)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to scan with config: %v", err)), nil
			}
			return mcp.NewToolResultText(output), nil
		}

		// Use ScanGitRepo for git repositories
		output, err := module.ScanGitRepo(ctx, path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to scan git repository: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGitleaksModule(client)

		// Get path
		path := request.GetString("path", ".")
		if path == "" {
			path = "."
		}

		// If config is specified, use ScanWithConfig
		if config := request.GetString("config", ""); config != "" {
			output, err := module.ScanWithConfig(ctx, path, config)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to scan with config: %v", err)), nil
			}
			return mcp.NewToolResultText(output), nil
		}

		// Check if it's a file or directory
		absPath, _ := filepath.Abs(path)
		
		// Try as directory first, then as file
		output, err := module.ScanDirectory(ctx, path)
		if err != nil {
			// Try as file
			output, err = module.ScanFile(ctx, absPath)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to scan path: %v", err)), nil
			}
		}

		return mcp.NewToolResultText(output), nil
	})

	// Gitleaks detect secrets from stdin
	stdinScanTool := mcp.NewTool("gitleaks_stdin",
		mcp.WithDescription("Detect secrets from standard input"),
		mcp.WithString("input",
			mcp.Description("Input text to scan for secrets"),
			mcp.Required(),
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
	s.AddTool(stdinScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGitleaksModule(client)

		// Get input text
		input := request.GetString("input", "")
		if input == "" {
			return mcp.NewToolResultError("input text is required"), nil
		}

		// Scan stdin
		output, err := module.ScanStdin(ctx, input)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to scan stdin: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Gitleaks version information
	versionTool := mcp.NewTool("gitleaks_version",
		mcp.WithDescription("Display Gitleaks version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGitleaksModule(client)

		// Get version
		version, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get version: %v", err)), nil
		}

		return mcp.NewToolResultText(version), nil
	})
}