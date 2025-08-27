package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddGitleaksTools adds Gitleaks (fast secret scanning) MCP tool implementations using direct Dagger calls
func AddGitleaksTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGitleaksToolsDirect(s)
}

// addGitleaksToolsDirect adds Gitleaks tools using direct Dagger module calls
func addGitleaksToolsDirect(s *server.MCPServer) {
	// Gitleaks detect tool
	detectTool := mcp.NewTool("gitleaks_detect",
		mcp.WithDescription("Detect secrets in git repository using Gitleaks"),
		mcp.WithString("source_path",
			mcp.Description("Path to git repository or directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("config_path",
			mcp.Description("Path to Gitleaks configuration file"),
		),
		mcp.WithString("report_format",
			mcp.Description("Report format"),
			mcp.Enum("json", "csv", "sarif"),
		),
		mcp.WithString("report_path",
			mcp.Description("Path to write the report"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Show verbose output"),
		),
		mcp.WithBoolean("no_git",
			mcp.Description("Treat git repo as a regular directory and scan those files"),
		),
	)
	s.AddTool(detectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		sourcePath := request.GetString("source_path", "")
		configPath := request.GetString("config_path", "")
		reportFormat := request.GetString("report_format", "")
		reportPath := request.GetString("report_path", "")
		verbose := request.GetBool("verbose", false)
		noGit := request.GetBool("no_git", false)

		if sourcePath == "" {
			return mcp.NewToolResultError("source_path is required"), nil
		}

		// Create Gitleaks module
		gitleaksModule := modules.NewGitleaksModule(client)

		// Set up options
		opts := modules.GitleaksDetectOptions{
			ConfigPath:   configPath,
			ReportFormat: reportFormat,
			ReportPath:   reportPath,
			Verbose:      verbose,
			NoGit:        noGit,
		}

		// Run detection
		result, err := gitleaksModule.Detect(ctx, sourcePath, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Gitleaks detection failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Gitleaks protect tool
	protectTool := mcp.NewTool("gitleaks_protect",
		mcp.WithDescription("Protect git repository with Gitleaks pre-commit scanning"),
		mcp.WithString("source_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithString("config_path",
			mcp.Description("Path to Gitleaks configuration file"),
		),
		mcp.WithBoolean("staged",
			mcp.Description("Scan only staged files"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Show verbose output"),
		),
	)
	s.AddTool(protectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		sourcePath := request.GetString("source_path", "")
		configPath := request.GetString("config_path", "")
		staged := request.GetBool("staged", false)
		verbose := request.GetBool("verbose", false)

		if sourcePath == "" {
			return mcp.NewToolResultError("source_path is required"), nil
		}

		// Create Gitleaks module
		gitleaksModule := modules.NewGitleaksModule(client)

		// Set up options
		opts := modules.GitleaksProtectOptions{
			ConfigPath: configPath,
			Staged:     staged,
			Verbose:    verbose,
		}

		// Run protection
		result, err := gitleaksModule.Protect(ctx, sourcePath, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Gitleaks protection failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}