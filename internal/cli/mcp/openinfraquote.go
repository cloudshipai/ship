package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
		mcp.WithDescription("Match Terraform plan/state JSON to pricing rows using oiq match"),
		mcp.WithString("plan_or_state_json_path",
			mcp.Description("Path to Terraform plan or state JSON file (from terraform show -json)"),
			mcp.Required(),
		),
		mcp.WithString("pricesheet_path",
			mcp.Description("Path to pricing CSV file (download with download_pricesheet)"),
			mcp.Required(),
		),
		mcp.WithString("output_path",
			mcp.Description("Optional: Write matched results to file instead of stdout"),
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
		planOrStateJson := request.GetString("plan_or_state_json_path", "")
		if planOrStateJson == "" {
			return mcp.NewToolResultError("plan_or_state_json_path is required"), nil
		}
		pricesheetPath := request.GetString("pricesheet_path", "")
		if pricesheetPath == "" {
			return mcp.NewToolResultError("pricesheet_path is required"), nil
		}
		// TODO: Add support for output_path in the module
		_ = request.GetString("output_path", "")

		// Match resources
		output, err := module.Match(ctx, pricesheetPath, planOrStateJson)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oiq match failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote price tool
	priceTool := mcp.NewTool("openinfraquote_price",
		mcp.WithDescription("Calculate prices from matched resources using oiq price"),
		mcp.WithString("matched_input_path",
			mcp.Description("Path to matched JSON file from oiq match (or pipe from match)"),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing (e.g., us-east-1) - required for accurate pricing"),
		),
		mcp.WithString("usage_path",
			mcp.Description("Optional: Path to usage.json file for custom usage assumptions"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: json, summary (default), text, markdown, atlantis-comment"),
		),
		mcp.WithString("mq",
			mcp.Description("Optional: Match query for advanced resource filtering"),
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
		matchedInputPath := request.GetString("matched_input_path", "")
		region := request.GetString("region", "us-east-1")
		// TODO: Add support for these parameters in the module
		_ = request.GetString("usage_path", "")
		_ = request.GetString("format", "summary")
		_ = request.GetString("mq", "")

		// Calculate prices
		output, err := module.Price(ctx, region, matchedInputPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oiq price failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote download pricesheet tool
	downloadPricesheetTool := mcp.NewTool("openinfraquote_download_pricesheet",
		mcp.WithDescription("Download AWS pricing data CSV from oiq.terrateam.io"),
		mcp.WithString("dest_path",
			mcp.Description("Output file for pricing CSV (default: ./prices.csv)"),
		),
	)
	s.AddTool(downloadPricesheetTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get parameters
		destPath := request.GetString("dest_path", "./prices.csv")

		// Download prices
		_, err = module.DownloadPrices(ctx, destPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("download pricesheet failed: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Downloaded pricesheet to %s", destPath)), nil
	})

	// OpenInfraQuote print-default-usage tool
	printDefaultUsageTool := mcp.NewTool("openinfraquote_print_default_usage",
		mcp.WithDescription("Print default usage assumptions for cost estimation"),
	)
	s.AddTool(printDefaultUsageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get default usage
		output, err := module.PrintDefaultUsage(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("print-default-usage failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote compare regions tool
	compareRegionsTool := mcp.NewTool("openinfraquote_compare_regions",
		mcp.WithDescription("Compare costs across multiple AWS regions"),
		mcp.WithString("plan_file",
			mcp.Description("Path to Terraform plan JSON file"),
			mcp.Required(),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of AWS regions (e.g., us-east-1,us-west-2,eu-west-1)"),
			mcp.Required(),
		),
	)
	s.AddTool(compareRegionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get parameters
		planFile := request.GetString("plan_file", "")
		if planFile == "" {
			return mcp.NewToolResultError("plan_file is required"), nil
		}
		regionsStr := request.GetString("regions", "")
		if regionsStr == "" {
			return mcp.NewToolResultError("regions is required"), nil
		}

		// Split regions
		regions := []string{}
		for _, r := range strings.Split(regionsStr, ",") {
			trimmed := strings.TrimSpace(r)
			if trimmed != "" {
				regions = append(regions, trimmed)
			}
		}

		if len(regions) == 0 {
			return mcp.NewToolResultError("at least one region is required"), nil
		}

		// Compare regions
		output, err := module.CompareRegions(ctx, planFile, regions)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("compare regions failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote cost gate tool - policy enforcement
	costGateTool := mcp.NewTool("openinfraquote_cost_gate",
		mcp.WithDescription("Enforce cost policies on OIQ JSON output"),
		mcp.WithString("report_json",
			mcp.Description("JSON output from oiq price --format=json"),
			mcp.Required(),
		),
		mcp.WithNumber("max_total_usd",
			mcp.Description("Maximum allowed total monthly cost in USD"),
		),
		mcp.WithNumber("max_monthly_delta_usd",
			mcp.Description("Maximum allowed monthly cost increase in USD"),
		),
	)
	s.AddTool(costGateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		reportJson := request.GetString("report_json", "")
		if reportJson == "" {
			return mcp.NewToolResultError("report_json is required"), nil
		}
		// GetNumber doesn't exist, use GetFloat
		maxTotalUSD := request.GetFloat("max_total_usd", 0)
		maxDeltaUSD := request.GetFloat("max_monthly_delta_usd", 0)

		// Parse JSON report
		var report map[string]interface{}
		if err := json.Unmarshal([]byte(reportJson), &report); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse JSON report: %v", err)), nil
		}

		// Extract pricing information
		var reasons []string
		pass := true

		// Check price difference
		if priceDiff, ok := report["price_diff"].(map[string]interface{}); ok {
			if maxVal, ok := priceDiff["max"].(float64); ok && maxDeltaUSD > 0 {
				if maxVal > maxDeltaUSD {
					pass = false
					reasons = append(reasons, fmt.Sprintf("Cost increase $%.2f exceeds limit $%.2f", maxVal, maxDeltaUSD))
				}
			}
		}

		// Check total price
		if price, ok := report["price"].(map[string]interface{}); ok {
			if maxVal, ok := price["max"].(float64); ok && maxTotalUSD > 0 {
				if maxVal > maxTotalUSD {
					pass = false
					reasons = append(reasons, fmt.Sprintf("Total cost $%.2f exceeds limit $%.2f", maxVal, maxTotalUSD))
				}
			}
		}

		// Build response
		result := map[string]interface{}{
			"pass":    pass,
			"reasons": reasons,
		}

		resultJson, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJson)), nil
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

	// OpenInfraQuote analyze directory tool - handles .tf files properly
	analyzeDirectoryTool := mcp.NewTool("openinfraquote_analyze_directory",
		mcp.WithDescription("Analyze all Terraform files in a directory - automatically generates plan, downloads pricesheet, and estimates costs"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform .tf files"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing (e.g., us-east-1)"),
		),
	)
	s.AddTool(analyzeDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get parameters
		directory := request.GetString("directory", ".")
		region := request.GetString("region", "us-east-1")

		// Analyze directory (handles terraform init, plan, and cost estimation)
		output, err := module.AnalyzeDirectory(ctx, directory, region)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oiq analyze directory failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote analyze plan tool - for existing tfplan.json files
	analyzePlanTool := mcp.NewTool("openinfraquote_analyze_plan",
		mcp.WithDescription("Analyze an existing Terraform plan JSON file - automatically downloads pricesheet and estimates costs"),
		mcp.WithString("plan_file",
			mcp.Description("Path to Terraform plan JSON file (tfplan.json)"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing (e.g., us-east-1)"),
		),
	)
	s.AddTool(analyzePlanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenInfraQuoteModule(client)

		// Get parameters
		planFile := request.GetString("plan_file", "")
		if planFile == "" {
			return mcp.NewToolResultError("plan_file is required"), nil
		}
		region := request.GetString("region", "us-east-1")

		// Analyze plan (automatically downloads pricesheet)
		output, err := module.AnalyzePlan(ctx, planFile, region)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oiq analyze plan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenInfraQuote full pipeline tool
	fullPipelineTool := mcp.NewTool("openinfraquote_full_pipeline",
		mcp.WithDescription("Run full cost estimation pipeline using oiq (requires existing tfplan.json and pricesheet)"),
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