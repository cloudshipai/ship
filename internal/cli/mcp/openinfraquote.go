package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddOpenInfraQuoteTools adds OpenInfraQuote (infrastructure cost estimation) MCP tool implementations using real oiq CLI
func AddOpenInfraQuoteTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// OpenInfraQuote match tool
	matchTool := mcp.NewTool("openinfraquote_match",
		mcp.WithDescription("Match Terraform resources to pricing using real oiq CLI"),
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
		pricesheet := request.GetString("pricesheet", "")
		tfplanJson := request.GetString("tfplan_json", "")
		args := []string{"oiq", "match", "--pricesheet", pricesheet, tfplanJson}
		return executeShipCommand(args)
	})

	// OpenInfraQuote price tool
	priceTool := mcp.NewTool("openinfraquote_price",
		mcp.WithDescription("Calculate prices from matched resources using real oiq CLI"),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing (e.g., us-east-1)"),
		),
		mcp.WithString("input_file",
			mcp.Description("Input file with matched resources (or use stdin)"),
		),
	)
	s.AddTool(priceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"oiq", "price"}
		
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		if inputFile := request.GetString("input_file", ""); inputFile != "" {
			args = append(args, inputFile)
		}
		
		return executeShipCommand(args)
	})

	// OpenInfraQuote download prices tool
	downloadPricesTool := mcp.NewTool("openinfraquote_download_prices",
		mcp.WithDescription("Download AWS pricing data using curl and gunzip"),
		mcp.WithString("output_file",
			mcp.Description("Output file for pricing CSV (default: prices.csv)"),
		),
	)
	s.AddTool(downloadPricesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		outputFile := request.GetString("output_file", "prices.csv")
		// Using curl and gunzip to download pricing data
		args := []string{"sh", "-c", "curl -s https://oiq.terrateam.io/prices.csv.gz | gunzip > " + outputFile}
		return executeShipCommand(args)
	})

	// OpenInfraQuote help tool
	helpTool := mcp.NewTool("openinfraquote_help",
		mcp.WithDescription("Get OpenInfraQuote help information using real oiq CLI"),
	)
	s.AddTool(helpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"oiq", "--help"}
		return executeShipCommand(args)
	})
	
	// OpenInfraQuote full pipeline tool
	fullPipelineTool := mcp.NewTool("openinfraquote_full_pipeline",
		mcp.WithDescription("Run full cost estimation pipeline using real oiq CLI"),
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
		tfplanJson := request.GetString("tfplan_json", "")
		pricesheet := request.GetString("pricesheet", "")
		region := request.GetString("region", "")
		
		// Running the full pipeline: match | price
		args := []string{"sh", "-c", 
			"oiq match --pricesheet " + pricesheet + " " + tfplanJson + " | oiq price --region " + region}
		return executeShipCommand(args)
	})
}