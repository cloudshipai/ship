package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddActionlintTools adds Actionlint (GitHub Actions linter) MCP tool implementations
func AddActionlintTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		args := []string{"actionlint"}
		
		// Add specific workflow files if provided
		if workflowFiles := request.GetString("workflow_files", ""); workflowFiles != "" {
			// Split comma-separated files and add them
			files := strings.Split(workflowFiles, ",")
			for _, file := range files {
				file = strings.TrimSpace(file)
				if file != "" {
					args = append(args, file)
				}
			}
		}
		
		if formatTemplate := request.GetString("format_template", ""); formatTemplate != "" {
			args = append(args, "-format", formatTemplate)
		}
		if ignorePatterns := request.GetString("ignore_patterns", ""); ignorePatterns != "" {
			// Split comma-separated patterns and add each with -ignore flag
			patterns := strings.Split(ignorePatterns, ",")
			for _, pattern := range patterns {
				pattern = strings.TrimSpace(pattern)
				if pattern != "" {
					args = append(args, "-ignore", pattern)
				}
			}
		}
		if request.GetBool("color", false) {
			args = append(args, "-color")
		}
		
		return executeShipCommand(args)
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
		args := []string{"actionlint"}
		
		// Add specific workflow files if provided
		if workflowFiles := request.GetString("workflow_files", ""); workflowFiles != "" {
			files := strings.Split(workflowFiles, ",")
			for _, file := range files {
				file = strings.TrimSpace(file)
				if file != "" {
					args = append(args, file)
				}
			}
		}
		
		if shellcheckPath := request.GetString("shellcheck_path", ""); shellcheckPath != "" {
			args = append(args, "-shellcheck", shellcheckPath)
		}
		if pyflakesPath := request.GetString("pyflakes_path", ""); pyflakesPath != "" {
			args = append(args, "-pyflakes", pyflakesPath)
		}
		if request.GetBool("color", false) {
			args = append(args, "-color")
		}
		
		return executeShipCommand(args)
	})

	// Actionlint get version tool
	getVersionTool := mcp.NewTool("actionlint_get_version",
		mcp.WithDescription("Get Actionlint version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"actionlint", "-version"}
		return executeShipCommand(args)
	})
}