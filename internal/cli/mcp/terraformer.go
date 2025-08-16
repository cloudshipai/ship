package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTerraformerTools adds Terraformer MCP tool implementations
func AddTerraformerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Terraformer import tool
	importTool := mcp.NewTool("terraformer_import",
		mcp.WithDescription("Import existing infrastructure to Terraform using Terraformer"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, gcp, azure, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("services",
			mcp.Description("Comma-separated list of services to import"),
			mcp.Required(),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of regions"),
		),
	)
	s.AddTool(importTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provider := request.GetString("provider", "")
		services := request.GetString("services", "")
		args := []string{"terraform", "terraformer", "import", "--provider", provider, "--services", services}
		if regions := request.GetString("regions", ""); regions != "" {
			args = append(args, "--regions", regions)
		}
		return executeShipCommand(args)
	})

	// Terraformer plan tool
	planTool := mcp.NewTool("terraformer_plan",
		mcp.WithDescription("Generate Terraform plan from imported infrastructure"),
		mcp.WithString("directory",
			mcp.Description("Directory containing imported Terraform files"),
			mcp.Required(),
		),
	)
	s.AddTool(planTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"terraform", "terraformer", "plan", directory}
		return executeShipCommand(args)
	})

	// Terraformer validate tool
	validateTool := mcp.NewTool("terraformer_validate",
		mcp.WithDescription("Validate imported Terraform configuration"),
		mcp.WithString("directory",
			mcp.Description("Directory containing imported Terraform files"),
			mcp.Required(),
		),
	)
	s.AddTool(validateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"terraform", "terraformer", "validate", directory}
		return executeShipCommand(args)
	})

	// Terraformer list providers tool
	listProvidersTool := mcp.NewTool("terraformer_list_providers",
		mcp.WithDescription("List supported providers and services for Terraformer"),
	)
	s.AddTool(listProvidersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform", "terraformer", "list-providers"}
		return executeShipCommand(args)
	})

	// Terraformer refresh tool
	refreshTool := mcp.NewTool("terraformer_refresh",
		mcp.WithDescription("Refresh Terraform state for imported resources"),
		mcp.WithString("directory",
			mcp.Description("Directory containing imported Terraform files"),
			mcp.Required(),
		),
	)
	s.AddTool(refreshTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"terraform", "terraformer", "refresh", directory}
		return executeShipCommand(args)
	})

	// Terraformer get version tool
	getVersionTool := mcp.NewTool("terraformer_get_version",
		mcp.WithDescription("Get Terraformer version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform", "terraformer", "--version"}
		return executeShipCommand(args)
	})
}