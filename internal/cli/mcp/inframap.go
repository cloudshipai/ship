package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddInfraMapTools adds InfraMap (infrastructure diagram generator) MCP tool implementations
func AddInfraMapTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// InfraMap generate from state tool
	generateFromStateTool := mcp.NewTool("inframap_generate_from_state",
		mcp.WithDescription("Generate infrastructure diagram from Terraform state file"),
		mcp.WithString("state_file",
			mcp.Description("Path to Terraform state file"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("png", "svg", "pdf", "dot"),
		),
	)
	s.AddTool(generateFromStateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stateFile := request.GetString("state_file", "")
		args := []string{"terraform-tools", "inframap", "--state", stateFile}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// InfraMap generate from HCL tool
	generateFromHCLTool := mcp.NewTool("inframap_generate_from_hcl",
		mcp.WithDescription("Generate infrastructure diagram from Terraform HCL files"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform HCL files"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("png", "svg", "pdf", "dot"),
		),
	)
	s.AddTool(generateFromHCLTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"terraform-tools", "inframap", "--hcl", directory}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// InfraMap generate with options tool
	generateWithOptionsTool := mcp.NewTool("inframap_generate_with_options",
		mcp.WithDescription("Generate infrastructure diagram with custom options"),
		mcp.WithString("input",
			mcp.Description("Input file or directory"),
			mcp.Required(),
		),
		mcp.WithString("provider",
			mcp.Description("Filter by specific provider"),
			mcp.Enum("aws", "google", "azurerm", "digitalocean"),
		),
		mcp.WithBoolean("raw",
			mcp.Description("Show all resources without InfraMap logic"),
		),
		mcp.WithBoolean("clean",
			mcp.Description("Remove unconnected nodes"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("png", "svg", "pdf", "dot"),
		),
	)
	s.AddTool(generateWithOptionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		input := request.GetString("input", "")
		args := []string{"terraform-tools", "inframap", input}
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, "--provider", provider)
		}
		if request.GetBool("raw", false) {
			args = append(args, "--raw")
		}
		if !request.GetBool("clean", true) {
			args = append(args, "--clean=false")
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// InfraMap prune state tool
	pruneStateTool := mcp.NewTool("inframap_prune_state",
		mcp.WithDescription("Remove unnecessary information from Terraform state"),
		mcp.WithString("state_file",
			mcp.Description("Path to Terraform state file to prune"),
			mcp.Required(),
		),
	)
	s.AddTool(pruneStateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stateFile := request.GetString("state_file", "")
		args := []string{"terraform-tools", "inframap", "--prune", stateFile}
		return executeShipCommand(args)
	})
}