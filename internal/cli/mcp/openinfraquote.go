package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddOpenInfraQuoteTools adds OpenInfraQuote (cost estimation) MCP tool implementations using direct Dagger calls
func AddOpenInfraQuoteTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addOpenInfraQuoteToolsDirect(s)
}

// addOpenInfraQuoteToolsDirect adds OpenInfraQuote tools using direct Dagger module calls
func addOpenInfraQuoteToolsDirect(s *server.MCPServer) {
	// OpenInfraQuote estimate tool
	estimateTool := mcp.NewTool("openinfraquote_estimate",
		mcp.WithDescription("Generate cost estimates for Terraform infrastructure"),
		mcp.WithString("terraform_path",
			mcp.Description("Path to Terraform configuration directory"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for cost estimation"),
			mcp.Enum("table", "json", "html", "csv"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path (optional, prints to stdout if not specified)"),
		),
		mcp.WithString("terraform_plan_file",
			mcp.Description("Path to Terraform plan file (JSON format)"),
		),
		mcp.WithString("currency",
			mcp.Description("Currency for cost estimation"),
			mcp.Enum("USD", "EUR", "GBP", "CAD", "AUD", "INR", "JPY"),
		),
		mcp.WithString("region",
			mcp.Description("Cloud region for pricing"),
		),
		mcp.WithBoolean("show_skipped",
			mcp.Description("Show skipped resources in output"),
		),
		mcp.WithBoolean("sync_usage_file",
			mcp.Description("Sync usage file with missing resources"),
		),
		mcp.WithString("usage_file",
			mcp.Description("Path to usage file for accurate cost estimation"),
		),
	)
	s.AddTool(estimateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		terraformPath := request.GetString("terraform_path", "")
		outputFormat := request.GetString("output_format", "")
		outputFile := request.GetString("output_file", "")
		terraformPlanFile := request.GetString("terraform_plan_file", "")
		currency := request.GetString("currency", "")
		region := request.GetString("region", "")
		showSkipped := request.GetBool("show_skipped", false)
		syncUsageFile := request.GetBool("sync_usage_file", false)
		usageFile := request.GetString("usage_file", "")

		if terraformPath == "" {
			return mcp.NewToolResultError("terraform_path is required"), nil
		}

		// Set defaults
		if outputFormat == "" {
			outputFormat = "table"
		}
		if currency == "" {
			currency = "USD"
		}

		// Create OpenInfraQuote module
		openInfraQuoteModule := modules.NewOpenInfraQuoteModule(client)

		// Set up options
		opts := modules.OpenInfraQuoteOptions{
			OutputFormat:      outputFormat,
			OutputFile:        outputFile,
			TerraformPlanFile: terraformPlanFile,
			Currency:          currency,
			Region:            region,
			ShowSkipped:       showSkipped,
			SyncUsageFile:     syncUsageFile,
			UsageFile:         usageFile,
		}

		// Generate cost estimate
		result, err := openInfraQuoteModule.Estimate(ctx, terraformPath, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("OpenInfraQuote estimation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// OpenInfraQuote diff tool
	diffTool := mcp.NewTool("openinfraquote_diff",
		mcp.WithDescription("Show cost difference between two Terraform configurations or states"),
		mcp.WithString("path1",
			mcp.Description("Path to first Terraform configuration"),
			mcp.Required(),
		),
		mcp.WithString("path2",
			mcp.Description("Path to second Terraform configuration"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for cost diff"),
			mcp.Enum("table", "json", "html"),
		),
		mcp.WithString("currency",
			mcp.Description("Currency for cost comparison"),
			mcp.Enum("USD", "EUR", "GBP", "CAD", "AUD", "INR", "JPY"),
		),
		mcp.WithBoolean("show_all_projects",
			mcp.Description("Show all projects in diff output"),
		),
	)
	s.AddTool(diffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		path1 := request.GetString("path1", "")
		path2 := request.GetString("path2", "")
		outputFormat := request.GetString("output_format", "")
		currency := request.GetString("currency", "")
		showAllProjects := request.GetBool("show_all_projects", false)

		if path1 == "" || path2 == "" {
			return mcp.NewToolResultError("both path1 and path2 are required"), nil
		}

		// Set defaults
		if outputFormat == "" {
			outputFormat = "table"
		}
		if currency == "" {
			currency = "USD"
		}

		// Create OpenInfraQuote module
		openInfraQuoteModule := modules.NewOpenInfraQuoteModule(client)

		// Set up options
		opts := modules.OpenInfraQuoteDiffOptions{
			OutputFormat:      outputFormat,
			Currency:          currency,
			ShowAllProjects:   showAllProjects,
		}

		// Generate cost diff
		result, err := openInfraQuoteModule.Diff(ctx, path1, path2, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("OpenInfraQuote diff failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}