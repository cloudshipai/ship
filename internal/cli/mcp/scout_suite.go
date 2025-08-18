package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddScoutSuiteTools adds Scout Suite MCP tool implementations using direct Dagger calls
func AddScoutSuiteTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addScoutSuiteToolsDirect(s)
}

// addScoutSuiteToolsDirect adds Scout Suite tools using direct Dagger module calls
func addScoutSuiteToolsDirect(s *server.MCPServer) {
	// Scout Suite scan AWS tool
	scanAWSTool := mcp.NewTool("scout_suite_scan_aws",
		mcp.WithDescription("Scan AWS environment for security issues using real scout CLI"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use for scanning"),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of AWS regions to scan - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("services",
			mcp.Description("Comma-separated list of AWS services to scan - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("report_dir",
			mcp.Description("Directory to save the report - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("exceptions",
			mcp.Description("Path to exceptions file - NOTE: not supported in current Dagger module"),
		),
	)
	s.AddTool(scanAWSTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewScoutSuiteModule(client)

		// Get parameters
		profile := request.GetString("profile", "default")

		// Check for unsupported parameters
		if request.GetString("regions", "") != "" || request.GetString("services", "") != "" ||
			request.GetString("report_dir", "") != "" || request.GetString("exceptions", "") != "" {
			return mcp.NewToolResultError("regions, services, report_dir, and exceptions options not supported in current Dagger module"), nil
		}

		// Scan AWS
		output, err := module.ScanAWS(ctx, profile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scout suite AWS scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Scout Suite scan Azure tool
	scanAzureTool := mcp.NewTool("scout_suite_scan_azure",
		mcp.WithDescription("Scan Azure environment for security issues using real scout CLI"),
		mcp.WithString("subscriptions",
			mcp.Description("Azure subscription IDs to scan (space-separated) - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("tenant_id",
			mcp.Description("Azure tenant ID - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("username",
			mcp.Description("Azure username - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("password",
			mcp.Description("Azure password - NOTE: not supported in current Dagger module"),
		),
		mcp.WithBoolean("cli",
			mcp.Description("Use Azure CLI for authentication - NOTE: not supported in current Dagger module"),
		),
		mcp.WithBoolean("service_principal",
			mcp.Description("Use service principal authentication - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("report_dir",
			mcp.Description("Directory to save the report - NOTE: not supported in current Dagger module"),
		),
	)
	s.AddTool(scanAzureTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewScoutSuiteModule(client)

		// Check for unsupported parameters
		if request.GetString("subscriptions", "") != "" || request.GetString("tenant_id", "") != "" ||
			request.GetString("username", "") != "" || request.GetString("password", "") != "" ||
			request.GetBool("cli", false) || request.GetBool("service_principal", false) ||
			request.GetString("report_dir", "") != "" {
			return mcp.NewToolResultError("advanced Azure options not supported in current Dagger module - uses environment variables"), nil
		}

		// Scan Azure
		output, err := module.ScanAzure(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scout suite Azure scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Scout Suite scan GCP tool
	scanGCPTool := mcp.NewTool("scout_suite_scan_gcp",
		mcp.WithDescription("Scan Google Cloud Platform for security issues using real scout CLI"),
		mcp.WithString("project_id",
			mcp.Description("GCP project ID to scan"),
		),
		mcp.WithString("folder_id",
			mcp.Description("GCP folder ID to scan - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("organization_id",
			mcp.Description("GCP organization ID to scan - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("service_account",
			mcp.Description("Path to GCP service account key file - NOTE: not supported in current Dagger module"),
		),
		mcp.WithBoolean("user_account",
			mcp.Description("Use user account for authentication - NOTE: not supported in current Dagger module"),
		),
		mcp.WithBoolean("all_projects",
			mcp.Description("Scan all accessible projects - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("report_dir",
			mcp.Description("Directory to save the report - NOTE: not supported in current Dagger module"),
		),
	)
	s.AddTool(scanGCPTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewScoutSuiteModule(client)

		// Get parameters
		projectID := request.GetString("project_id", "")
		if projectID == "" {
			return mcp.NewToolResultError("project_id is required"), nil
		}

		// Check for unsupported parameters
		if request.GetString("folder_id", "") != "" || request.GetString("organization_id", "") != "" ||
			request.GetString("service_account", "") != "" || request.GetBool("user_account", false) ||
			request.GetBool("all_projects", false) || request.GetString("report_dir", "") != "" {
			return mcp.NewToolResultError("advanced GCP options not supported in current Dagger module"), nil
		}

		// Scan GCP
		output, err := module.ScanGCP(ctx, projectID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scout suite GCP scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Scout Suite serve report tool
	serveReportTool := mcp.NewTool("scout_suite_serve_report",
		mcp.WithDescription("Serve Scout Suite report using real scout CLI"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, azure, gcp)"),
			mcp.Required(),
		),
		mcp.WithString("report_name",
			mcp.Description("Name of report to serve"),
		),
		mcp.WithString("host",
			mcp.Description("Host to bind server to (default 127.0.0.1)"),
		),
		mcp.WithString("port",
			mcp.Description("Port to bind server to (default 8080)"),
		),
	)
	s.AddTool(serveReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewScoutSuiteModule(client)

		// Get parameters
		provider := request.GetString("provider", "")
		reportName := request.GetString("report_name", "")
		host := request.GetString("host", "")
		port := request.GetString("port", "")

		// Serve report
		output, err := module.ServeReport(ctx, provider, reportName, host, port)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scout suite serve report failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Scout Suite help tool
	helpTool := mcp.NewTool("scout_suite_help",
		mcp.WithDescription("Get Scout Suite help information using real scout CLI"),
		mcp.WithString("provider",
			mcp.Description("Get help for specific provider (aws, azure, gcp)"),
		),
	)
	s.AddTool(helpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewScoutSuiteModule(client)

		// Get parameters
		provider := request.GetString("provider", "")

		// Get help
		output, err := module.Help(ctx, provider)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scout suite help failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}