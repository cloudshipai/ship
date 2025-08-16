package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddInfracostTools adds Infracost MCP tool implementations
func AddInfracostTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Infracost breakdown directory tool
	breakdownDirTool := mcp.NewTool("infracost_breakdown_directory",
		mcp.WithDescription("Generate cost breakdown for Terraform directory"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
	)
	s.AddTool(breakdownDirTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "cost-estimate"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// Infracost breakdown plan tool
	breakdownPlanTool := mcp.NewTool("infracost_breakdown_plan",
		mcp.WithDescription("Generate cost breakdown from Terraform plan file"),
		mcp.WithString("plan_file",
			mcp.Description("Path to Terraform plan JSON file"),
			mcp.Required(),
		),
	)
	s.AddTool(breakdownPlanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		planFile := request.GetString("plan_file", "")
		args := []string{"terraform-tools", "cost-estimate", "--plan", planFile}
		return executeShipCommand(args)
	})

	// Infracost diff tool
	diffTool := mcp.NewTool("infracost_diff",
		mcp.WithDescription("Show cost difference between current and planned state"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
	)
	s.AddTool(diffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "cost-estimate", "--diff"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// Infracost breakdown with config tool
	breakdownConfigTool := mcp.NewTool("infracost_breakdown_config",
		mcp.WithDescription("Generate cost breakdown using Infracost config file"),
		mcp.WithString("config_file",
			mcp.Description("Path to Infracost config file"),
			mcp.Required(),
		),
	)
	s.AddTool(breakdownConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configFile := request.GetString("config_file", "")
		args := []string{"terraform-tools", "cost-estimate", "--config", configFile}
		return executeShipCommand(args)
	})

	// Infracost generate HTML report tool
	htmlReportTool := mcp.NewTool("infracost_generate_html",
		mcp.WithDescription("Generate HTML cost report"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output",
			mcp.Description("Output file path for HTML report"),
		),
	)
	s.AddTool(htmlReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "cost-estimate", "--format", "html"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})

	// Infracost generate table report tool
	tableReportTool := mcp.NewTool("infracost_generate_table",
		mcp.WithDescription("Generate table format cost report"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output",
			mcp.Description("Output file path for table report"),
		),
	)
	s.AddTool(tableReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "cost-estimate", "--format", "table"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})

	// Infracost get version tool
	getVersionTool := mcp.NewTool("infracost_get_version",
		mcp.WithDescription("Get the version of Infracost"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "cost-estimate", "--version"}
		return executeShipCommand(args)
	})

	// Infracost get pricing tool
	getPricingTool := mcp.NewTool("infracost_get_pricing",
		mcp.WithDescription("Get cloud pricing information"),
		mcp.WithString("service",
			mcp.Description("Cloud service to get pricing for (e.g., aws, azure, gcp)"),
		),
	)
	s.AddTool(getPricingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "cost-estimate", "--pricing"}
		if service := request.GetString("service", ""); service != "" {
			args = append(args, "--service", service)
		}
		return executeShipCommand(args)
	})
}