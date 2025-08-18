package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddAWSPricingTools adds AWS Pricing (official AWS CLI pricing commands) MCP tool implementations
func AddAWSPricingTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addAWSPricingToolsDirect(s)
}

// addAWSPricingToolsDirect implements direct Dagger calls for AWS pricing tools
func addAWSPricingToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		serviceCode := request.GetString("service_code", "")
		maxItems := request.GetString("max_items", "")
		// Note: profile parameter ignored - AWS CLI in container handles auth via env vars

		// Create AWS pricing module and describe services
		awsModule := modules.NewAWSPricingModule(client)
		result, err := awsModule.DescribeServices(ctx, serviceCode, maxItems)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("describe services failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		serviceCode := request.GetString("service_code", "")
		attributeName := request.GetString("attribute_name", "")
		maxItems := request.GetString("max_items", "")
		// Note: profile parameter ignored - AWS CLI in container handles auth via env vars

		// Create AWS pricing module and get attribute values
		awsModule := modules.NewAWSPricingModule(client)
		result, err := awsModule.GetAttributeValues(ctx, serviceCode, attributeName, maxItems)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get attribute values failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		serviceCode := request.GetString("service_code", "")
		filters := request.GetString("filters", "")
		formatVersion := request.GetString("format_version", "")
		maxItems := request.GetString("max_items", "")
		// Note: profile parameter ignored - AWS CLI in container handles auth via env vars

		// Create AWS pricing module and get products
		awsModule := modules.NewAWSPricingModule(client)
		result, err := awsModule.GetProducts(ctx, serviceCode, filters, formatVersion, maxItems)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get products failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// AWS CLI version tool
	getVersionTool := mcp.NewTool("aws_pricing_get_version",
		mcp.WithDescription("Get AWS CLI version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create AWS pricing module and get version
		awsModule := modules.NewAWSPricingModule(client)
		result, err := awsModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get AWS CLI version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}