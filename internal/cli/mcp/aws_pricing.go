package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddAWSPricingTools adds AWS Pricing (cost calculator) MCP tool implementations
func AddAWSPricingTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// AWS pricing estimate tool
	estimateTool := mcp.NewTool("aws_pricing_estimate",
		mcp.WithDescription("Estimate AWS pricing for resources"),
		mcp.WithString("service",
			mcp.Description("AWS service name (e.g., ec2, s3, rds)"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing"),
		),
		mcp.WithString("instance_type",
			mcp.Description("Instance type for compute services"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "csv"),
		),
	)
	s.AddTool(estimateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service := request.GetString("service", "")
		args := []string{"aws", "pricing", "estimate", service}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		if instanceType := request.GetString("instance_type", ""); instanceType != "" {
			args = append(args, "--instance-type", instanceType)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// AWS pricing compare tool
	compareTool := mcp.NewTool("aws_pricing_compare",
		mcp.WithDescription("Compare AWS pricing across regions or instance types"),
		mcp.WithString("service",
			mcp.Description("AWS service name"),
			mcp.Required(),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of regions to compare"),
		),
		mcp.WithString("instance_types",
			mcp.Description("Comma-separated list of instance types to compare"),
		),
	)
	s.AddTool(compareTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service := request.GetString("service", "")
		args := []string{"aws", "pricing", "compare", service}
		if regions := request.GetString("regions", ""); regions != "" {
			args = append(args, "--regions", regions)
		}
		if instanceTypes := request.GetString("instance_types", ""); instanceTypes != "" {
			args = append(args, "--instance-types", instanceTypes)
		}
		return executeShipCommand(args)
	})

	// AWS pricing list services tool
	listServicesTool := mcp.NewTool("aws_pricing_list_services",
		mcp.WithDescription("List available AWS services for pricing"),
		mcp.WithString("category",
			mcp.Description("Service category filter"),
		),
	)
	s.AddTool(listServicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"aws", "pricing", "list-services"}
		if category := request.GetString("category", ""); category != "" {
			args = append(args, "--category", category)
		}
		return executeShipCommand(args)
	})

	// AWS pricing get version tool
	getVersionTool := mcp.NewTool("aws_pricing_get_version",
		mcp.WithDescription("Get AWS pricing tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"aws", "pricing", "--version"}
		return executeShipCommand(args)
	})
}