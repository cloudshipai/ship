package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddOpenInfraQuoteTools adds OpenInfraQuote (infrastructure cost estimation) MCP tool implementations
func AddOpenInfraQuoteTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// OpenInfraQuote estimate tool
	estimateTool := mcp.NewTool("openinfraquote_estimate",
		mcp.WithDescription("Estimate infrastructure costs using OpenInfraQuote"),
		mcp.WithString("terraform_dir",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "csv"),
		),
	)
	s.AddTool(estimateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		terraformDir := request.GetString("terraform_dir", "")
		args := []string{"terraform", "openinfraquote", "estimate", terraformDir}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// OpenInfraQuote breakdown tool
	breakdownTool := mcp.NewTool("openinfraquote_breakdown",
		mcp.WithDescription("Get detailed cost breakdown by resource"),
		mcp.WithString("terraform_dir",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("group_by",
			mcp.Description("Group costs by resource type or service"),
			mcp.Enum("resource", "service"),
		),
	)
	s.AddTool(breakdownTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		terraformDir := request.GetString("terraform_dir", "")
		args := []string{"terraform", "openinfraquote", "breakdown", terraformDir}
		if groupBy := request.GetString("group_by", ""); groupBy != "" {
			args = append(args, "--group-by", groupBy)
		}
		return executeShipCommand(args)
	})

	// OpenInfraQuote compare tool
	compareTool := mcp.NewTool("openinfraquote_compare",
		mcp.WithDescription("Compare costs between different Terraform configurations"),
		mcp.WithString("baseline_dir",
			mcp.Description("Baseline Terraform directory"),
			mcp.Required(),
		),
		mcp.WithString("comparison_dir",
			mcp.Description("Comparison Terraform directory"),
			mcp.Required(),
		),
	)
	s.AddTool(compareTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		baselineDir := request.GetString("baseline_dir", "")
		comparisonDir := request.GetString("comparison_dir", "")
		args := []string{"terraform", "openinfraquote", "compare", baselineDir, comparisonDir}
		return executeShipCommand(args)
	})

	// OpenInfraQuote get version tool
	getVersionTool := mcp.NewTool("openinfraquote_get_version",
		mcp.WithDescription("Get OpenInfraQuote version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform", "openinfraquote", "--version"}
		return executeShipCommand(args)
	})
}