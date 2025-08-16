package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddPackerTools adds Packer MCP tool implementations
func AddPackerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Packer build tool
	buildTool := mcp.NewTool("packer_build",
		mcp.WithDescription("Build machine images using Packer configuration"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
		mcp.WithString("variables",
			mcp.Description("Variables to pass to Packer build (key=value format)"),
		),
	)
	s.AddTool(buildTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateFile := request.GetString("template_file", "")
		args := []string{"cloud", "packer", "build", templateFile}
		if vars := request.GetString("variables", ""); vars != "" {
			args = append(args, "--var", vars)
		}
		return executeShipCommand(args)
	})

	// Packer validate tool
	validateTool := mcp.NewTool("packer_validate",
		mcp.WithDescription("Validate Packer configuration template"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
	)
	s.AddTool(validateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateFile := request.GetString("template_file", "")
		args := []string{"cloud", "packer", "validate", templateFile}
		return executeShipCommand(args)
	})

	// Packer inspect tool
	inspectTool := mcp.NewTool("packer_inspect",
		mcp.WithDescription("Inspect and analyze Packer template configuration"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
	)
	s.AddTool(inspectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateFile := request.GetString("template_file", "")
		args := []string{"cloud", "packer", "inspect", templateFile}
		return executeShipCommand(args)
	})

	// Packer fix tool
	fixTool := mcp.NewTool("packer_fix",
		mcp.WithDescription("Fix and upgrade Packer template to current version"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
	)
	s.AddTool(fixTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateFile := request.GetString("template_file", "")
		args := []string{"cloud", "packer", "fix", templateFile}
		return executeShipCommand(args)
	})

	// Packer console tool
	consoleTool := mcp.NewTool("packer_console",
		mcp.WithDescription("Open Packer console for template debugging"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
	)
	s.AddTool(consoleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateFile := request.GetString("template_file", "")
		args := []string{"cloud", "packer", "console", templateFile}
		return executeShipCommand(args)
	})

	// Packer get version tool
	getVersionTool := mcp.NewTool("packer_get_version",
		mcp.WithDescription("Get Packer version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloud", "packer", "--version"}
		return executeShipCommand(args)
	})
}