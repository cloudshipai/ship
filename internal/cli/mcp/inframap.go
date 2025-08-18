package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddInfraMapTools adds InfraMap (infrastructure diagram generator) MCP tool implementations using real CLI commands
func AddInfraMapTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// InfraMap generate tool
	generateTool := mcp.NewTool("inframap_generate",
		mcp.WithDescription("Generate infrastructure graph using inframap generate"),
		mcp.WithString("input",
			mcp.Description("Path to Terraform state file, HCL file, or directory"),
			mcp.Required(),
		),
		mcp.WithBoolean("hcl",
			mcp.Description("Force HCL input type"),
		),
		mcp.WithBoolean("tfstate",
			mcp.Description("Force Terraform state input type"),
		),
		mcp.WithBoolean("connections",
			mcp.Description("Enable connections in graph (default: true)"),
		),
		mcp.WithBoolean("raw",
			mcp.Description("Show configuration without InfraMap processing"),
		),
		mcp.WithBoolean("clean",
			mcp.Description("Remove unconnected nodes (default: true)"),
		),
	)
	s.AddTool(generateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		input := request.GetString("input", "")
		args := []string{"inframap", "generate", input}
		
		if request.GetBool("hcl", false) {
			args = append(args, "--hcl")
		}
		if request.GetBool("tfstate", false) {
			args = append(args, "--tfstate")
		}
		if !request.GetBool("connections", true) {
			args = append(args, "--connections=false")
		}
		if request.GetBool("raw", false) {
			args = append(args, "--raw")
		}
		if !request.GetBool("clean", true) {
			args = append(args, "--clean=false")
		}
		
		return executeShipCommand(args)
	})

	// InfraMap prune tool
	pruneTool := mcp.NewTool("inframap_prune",
		mcp.WithDescription("Remove unnecessary information from Terraform state/HCL using inframap prune"),
		mcp.WithString("input",
			mcp.Description("Path to Terraform state file or HCL file to prune"),
			mcp.Required(),
		),
	)
	s.AddTool(pruneTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		input := request.GetString("input", "")
		args := []string{"inframap", "prune", input}
		return executeShipCommand(args)
	})


}