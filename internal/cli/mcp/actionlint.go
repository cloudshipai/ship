package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddActionlintTools adds Actionlint (GitHub Actions linter) MCP tool implementations
func AddActionlintTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addActionlintToolsDirect(s)
}

// addActionlintToolsDirect implements direct Dagger calls for actionlint tools
func addActionlintToolsDirect(s *server.MCPServer) {
	// Actionlint scan workflows tool (basic usage)
	scanWorkflowsTool := mcp.NewTool("actionlint_scan_workflows",
		mcp.WithDescription("Scan GitHub Actions workflow files for issues"),
		mcp.WithString("workflow_files",
			mcp.Description("Comma-separated list of workflow file paths (leave empty to scan all)"),
		),
		mcp.WithString("format_template",
			mcp.Description("Go template for formatting output"),
		),
		mcp.WithString("ignore_patterns",
			mcp.Description("Comma-separated regex patterns to ignore errors"),
		),
		mcp.WithBoolean("color",
			mcp.Description("Enable colored output"),
		),
	)
	s.AddTool(scanWorkflowsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workflowFiles := request.GetString("workflow_files", "")
		formatTemplate := request.GetString("format_template", "")
		ignorePatterns := request.GetString("ignore_patterns", "")
		color := request.GetBool("color", false)
		
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create actionlint module
		actionlintModule := modules.NewActionlintModule(client)
		
		var result string
		if workflowFiles != "" {
			// Scan specific workflow files
			files := strings.Split(workflowFiles, ",")
			var cleanFiles []string
			for _, file := range files {
				file = strings.TrimSpace(file)
				if file != "" {
					cleanFiles = append(cleanFiles, file)
				}
			}
			// Use current directory as workspace
			result, err = actionlintModule.ScanSpecificFiles(ctx, ".", cleanFiles, formatTemplate, ignorePatterns, color)
		} else {
			// Scan all workflows in directory
			result, err = actionlintModule.ScanDirectoryWithOptions(ctx, ".", formatTemplate, ignorePatterns, color)
		}
		
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("actionlint scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Actionlint scan with external tools tool
	scanWithExternalToolsTool := mcp.NewTool("actionlint_scan_with_external_tools",
		mcp.WithDescription("Scan workflows with shellcheck and pyflakes integration"),
		mcp.WithString("workflow_files",
			mcp.Description("Comma-separated list of workflow file paths (leave empty to scan all)"),
		),
		mcp.WithString("shellcheck_path",
			mcp.Description("Path to shellcheck executable"),
		),
		mcp.WithString("pyflakes_path",
			mcp.Description("Path to pyflakes executable"),
		),
		mcp.WithBoolean("color",
			mcp.Description("Enable colored output"),
		),
	)
	s.AddTool(scanWithExternalToolsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workflowFiles := request.GetString("workflow_files", "")
		shellcheckPath := request.GetString("shellcheck_path", "")
		pyflakesPath := request.GetString("pyflakes_path", "")
		color := request.GetBool("color", false)
		
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create actionlint module and scan with external tools
		actionlintModule := modules.NewActionlintModule(client)
		
		result, err := actionlintModule.ScanWithExternalTools(ctx, ".", shellcheckPath, pyflakesPath, color)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("actionlint external tools scan failed: %v", err)), nil
		}

		// Add note if specific files were requested but not supported
		if workflowFiles != "" {
			result = fmt.Sprintf("Note: Specific workflow files (%s) not yet supported with external tools, scanned directory instead.\n\n%s", workflowFiles, result)
		}

		return mcp.NewToolResultText(result), nil
	})

	// Actionlint get version tool
	getVersionTool := mcp.NewTool("actionlint_get_version",
		mcp.WithDescription("Get Actionlint version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create actionlint module and get version
		actionlintModule := modules.NewActionlintModule(client)
		result, err := actionlintModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("actionlint version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}