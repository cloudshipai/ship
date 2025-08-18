package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddAWSPricingTools adds AWS Pricing (official AWS CLI pricing commands) MCP tool implementations
func AddAWSPricingTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// AWS pricing describe services tool
	describeServicesTool := mcp.NewTool("aws_pricing_describe_services",
		mcp.WithDescription("Get metadata for AWS services and their pricing attributes"),
		mcp.WithString("service_code",
			mcp.Description("AWS service code (e.g., AmazonEC2, AmazonS3) - leave empty to list all services"),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("max_items",
			mcp.Description("Maximum number of items to return"),
		),
	)
	s.AddTool(describeServicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"aws", "pricing", "describe-services"}
		if serviceCode := request.GetString("service_code", ""); serviceCode != "" {
			args = append(args, "--service-code", serviceCode)
		}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if maxItemsStr := request.GetString("max_items", ""); maxItemsStr != "" {
			args = append(args, "--max-items", maxItemsStr)
		}
		return executeShipCommand(args)
	})

	// AWS pricing get attribute values tool
	getAttributeValuesTool := mcp.NewTool("aws_pricing_get_attribute_values",
		mcp.WithDescription("Get available attribute values for AWS service pricing filters"),
		mcp.WithString("service_code",
			mcp.Description("AWS service code (e.g., AmazonEC2, AmazonS3)"),
			mcp.Required(),
		),
		mcp.WithString("attribute_name",
			mcp.Description("Attribute name (e.g., instanceType, volumeType, location)"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("max_items",
			mcp.Description("Maximum number of items to return"),
		),
	)
	s.AddTool(getAttributeValuesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serviceCode := request.GetString("service_code", "")
		attributeName := request.GetString("attribute_name", "")
		args := []string{"aws", "pricing", "get-attribute-values", "--service-code", serviceCode, "--attribute-name", attributeName}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if maxItemsStr := request.GetString("max_items", ""); maxItemsStr != "" {
			args = append(args, "--max-items", maxItemsStr)
		}
		return executeShipCommand(args)
	})

	// AWS pricing get products tool
	getProductsTool := mcp.NewTool("aws_pricing_get_products",
		mcp.WithDescription("Get AWS pricing information for products that match filter criteria"),
		mcp.WithString("service_code",
			mcp.Description("AWS service code (e.g., AmazonEC2, AmazonS3)"),
			mcp.Required(),
		),
		mcp.WithString("filters",
			mcp.Description("JSON string of filter criteria (e.g., '[{\"Type\":\"TERM_MATCH\",\"Field\":\"location\",\"Value\":\"US East (N. Virginia)\"}]')"),
		),
		mcp.WithString("format_version",
			mcp.Description("Format version for response"),
			mcp.Enum("aws_v1"),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("max_items",
			mcp.Description("Maximum number of items to return"),
		),
	)
	s.AddTool(getProductsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serviceCode := request.GetString("service_code", "")
		args := []string{"aws", "pricing", "get-products", "--service-code", serviceCode}
		if filters := request.GetString("filters", ""); filters != "" {
			args = append(args, "--filters", filters)
		}
		if formatVersion := request.GetString("format_version", ""); formatVersion != "" {
			args = append(args, "--format-version", formatVersion)
		}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if maxItemsStr := request.GetString("max_items", ""); maxItemsStr != "" {
			args = append(args, "--max-items", maxItemsStr)
		}
		return executeShipCommand(args)
	})

	// AWS CLI version tool
	getVersionTool := mcp.NewTool("aws_pricing_get_version",
		mcp.WithDescription("Get AWS CLI version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"aws", "--version"}
		return executeShipCommand(args)
	})
}