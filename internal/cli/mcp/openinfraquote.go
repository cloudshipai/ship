package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddOpenInfraQuoteTools adds OpenInfraQuote (infrastructure cost estimation) MCP tool implementations using direct Dagger calls
func AddOpenInfraQuoteTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addOpenInfraQuoteToolsDirect(s)
}

// addOpenInfraQuoteToolsDirect adds OpenInfraQuote tools using direct Dagger module calls
func addOpenInfraQuoteToolsDirect(s *server.MCPServer) {
	// OpenInfraQuote match tool
	matchTool := mcp.NewTool("openinfraquote_match",
		mcp.WithDescription("Match Terraform resources to pricing using oiq"),
		mcp.WithString("pricesheet",
			mcp.Description("Path to pricing CSV file"),
			mcp.Required(),
		),
		mcp.WithString("tfplan_json",
			mcp.Description("Path to Terraform plan JSON file"),
			mcp.Required(),
		),
	)
	s.AddTool(matchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get parameters
		pricesheet := request.GetString("pricesheet", "")
		if pricesheet == "" {
			return mcp.NewToolResultError("pricesheet is required"), nil
		}
		tfplanJson := request.GetString("tfplan_json", "")
		if tfplanJson == "" {
			return mcp.NewToolResultError("tfplan_json is required"), nil
		}

		// Match resources
		output, err := module.Match(ctx, pricesheet, tfplanJson)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oiq match failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote price tool
	priceTool := mcp.NewTool("openinfraquote_price",
		mcp.WithDescription("Calculate prices from matched resources using oiq"),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing (e.g., us-east-1)"),
		),
		mcp.WithString("input_file",
			mcp.Description("Input file with matched resources (or use stdin)"),
		),
	)
	s.AddTool(priceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get parameters
		region := request.GetString("region", "")
		inputFile := request.GetString("input_file", "")

		// Calculate prices
		output, err := module.Price(ctx, region, inputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oiq price failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote download prices tool
	downloadPricesTool := mcp.NewTool("openinfraquote_download_prices",
		mcp.WithDescription("Download AWS pricing data"),
		mcp.WithString("output_file",
			mcp.Description("Output file for pricing CSV (default: prices.csv)"),
		),
	)
	s.AddTool(downloadPricesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get parameters
		outputFile := request.GetString("output_file", "prices.csv")

		// Download prices
		output, err := module.DownloadPrices(ctx, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("download prices failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote help tool
	helpTool := mcp.NewTool("openinfraquote_help",
		mcp.WithDescription("Get OpenInfraQuote help information using oiq"),
	)
	s.AddTool(helpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get help
		output, err := module.GetHelp(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oiq help failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote full pipeline tool
	fullPipelineTool := mcp.NewTool("openinfraquote_full_pipeline",
		mcp.WithDescription("Run full cost estimation pipeline using oiq"),
		mcp.WithString("tfplan_json",
			mcp.Description("Path to Terraform plan JSON file"),
			mcp.Required(),
		),
		mcp.WithString("pricesheet",
			mcp.Description("Path to pricing CSV file"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing (e.g., us-east-1)"),
			mcp.Required(),
		),
	)
	s.AddTool(fullPipelineTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get parameters
		tfplanJson := request.GetString("tfplan_json", "")
		if tfplanJson == "" {
			return mcp.NewToolResultError("tfplan_json is required"), nil
		}
		pricesheet := request.GetString("pricesheet", "")
		if pricesheet == "" {
			return mcp.NewToolResultError("pricesheet is required"), nil
		}
		region := request.GetString("region", "")
		if region == "" {
			return mcp.NewToolResultError("region is required"), nil
		}

		// Run full pipeline
		output, err := module.FullPipeline(ctx, tfplanJson, pricesheet, region)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oiq full pipeline failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}