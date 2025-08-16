package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddActionlintTools adds Actionlint (GitHub Actions linter) MCP tool implementations
func AddActionlintTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Actionlint scan directory tool
	scanDirectoryTool := mcp.NewTool("actionlint_scan_directory",
		mcp.WithDescription("Scan directory for GitHub Actions workflow issues"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan for workflow files"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text", "sarif"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "actionlint", directory}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Actionlint scan specific file tool
	scanFileTool := mcp.NewTool("actionlint_scan_file",
		mcp.WithDescription("Scan specific GitHub Actions workflow file"),
		mcp.WithString("file_path",
			mcp.Description("Path to workflow file to scan"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text", "sarif"),
		),
	)
	s.AddTool(scanFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"security", "actionlint", "--file", filePath}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Actionlint scan with shellcheck tool
	scanWithShellcheckTool := mcp.NewTool("actionlint_scan_with_shellcheck",
		mcp.WithDescription("Scan workflows with shellcheck integration"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithBoolean("shellcheck",
			mcp.Description("Enable shellcheck integration"),
		),
	)
	s.AddTool(scanWithShellcheckTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "actionlint", directory}
		if request.GetBool("shellcheck", false) {
			args = append(args, "--shellcheck")
		}
		return executeShipCommand(args)
	})

	// Actionlint get version tool
	getVersionTool := mcp.NewTool("actionlint_get_version",
		mcp.WithDescription("Get Actionlint version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "actionlint", "--version"}
		return executeShipCommand(args)
	})
}