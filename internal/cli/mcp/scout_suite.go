package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddScoutSuiteTools adds Scout Suite (multi-cloud security auditing) MCP tool implementations
func AddScoutSuiteTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Scout Suite scan AWS tool
	scanAWSTool := mcp.NewTool("scout_suite_scan_aws",
		mcp.WithDescription("Scan AWS environment for security issues using Scout Suite"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use for scanning"),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of AWS regions to scan"),
		),
		mcp.WithString("services",
			mcp.Description("Comma-separated list of AWS services to scan"),
		),
	)
	s.AddTool(scanAWSTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "scout-suite", "scan", "--provider", "aws"}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if regions := request.GetString("regions", ""); regions != "" {
			args = append(args, "--regions", regions)
		}
		if services := request.GetString("services", ""); services != "" {
			args = append(args, "--services", services)
		}
		return executeShipCommand(args)
	})

	// Scout Suite scan Azure tool
	scanAzureTool := mcp.NewTool("scout_suite_scan_azure",
		mcp.WithDescription("Scan Azure environment for security issues using Scout Suite"),
		mcp.WithString("subscription_id",
			mcp.Description("Azure subscription ID to scan"),
		),
		mcp.WithString("tenant_id",
			mcp.Description("Azure tenant ID"),
		),
	)
	s.AddTool(scanAzureTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "scout-suite", "scan", "--provider", "azure"}
		if subscriptionID := request.GetString("subscription_id", ""); subscriptionID != "" {
			args = append(args, "--subscription", subscriptionID)
		}
		if tenantID := request.GetString("tenant_id", ""); tenantID != "" {
			args = append(args, "--tenant", tenantID)
		}
		return executeShipCommand(args)
	})

	// Scout Suite scan GCP tool
	scanGCPTool := mcp.NewTool("scout_suite_scan_gcp",
		mcp.WithDescription("Scan Google Cloud Platform for security issues using Scout Suite"),
		mcp.WithString("project_id",
			mcp.Description("GCP project ID to scan"),
			mcp.Required(),
		),
		mcp.WithString("credentials_file",
			mcp.Description("Path to GCP service account credentials file"),
		),
	)
	s.AddTool(scanGCPTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID := request.GetString("project_id", "")
		args := []string{"security", "scout-suite", "scan", "--provider", "gcp", "--project", projectID}
		if credentialsFile := request.GetString("credentials_file", ""); credentialsFile != "" {
			args = append(args, "--credentials", credentialsFile)
		}
		return executeShipCommand(args)
	})

	// Scout Suite generate report tool
	generateReportTool := mcp.NewTool("scout_suite_generate_report",
		mcp.WithDescription("Generate comprehensive security report from Scout Suite scan"),
		mcp.WithString("scan_results",
			mcp.Description("Path to Scout Suite scan results"),
			mcp.Required(),
		),
		mcp.WithString("report_format",
			mcp.Description("Report format (html, json, csv)"),
		),
		mcp.WithString("output_dir",
			mcp.Description("Output directory for the report"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scanResults := request.GetString("scan_results", "")
		args := []string{"security", "scout-suite", "report", scanResults}
		if reportFormat := request.GetString("report_format", ""); reportFormat != "" {
			args = append(args, "--format", reportFormat)
		}
		if outputDir := request.GetString("output_dir", ""); outputDir != "" {
			args = append(args, "--output", outputDir)
		}
		return executeShipCommand(args)
	})

	// Scout Suite list rules tool
	listRulesTool := mcp.NewTool("scout_suite_list_rules",
		mcp.WithDescription("List available Scout Suite security rules"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, azure, gcp)"),
			mcp.Required(),
		),
		mcp.WithString("service",
			mcp.Description("Specific service to list rules for"),
		),
	)
	s.AddTool(listRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provider := request.GetString("provider", "")
		args := []string{"security", "scout-suite", "list-rules", "--provider", provider}
		if service := request.GetString("service", ""); service != "" {
			args = append(args, "--service", service)
		}
		return executeShipCommand(args)
	})

	// Scout Suite validate rules tool
	validateRulesTool := mcp.NewTool("scout_suite_validate_rules",
		mcp.WithDescription("Validate custom Scout Suite rules"),
		mcp.WithString("rules_path",
			mcp.Description("Path to custom rules directory"),
			mcp.Required(),
		),
		mcp.WithString("provider",
			mcp.Description("Cloud provider for rule validation (aws, azure, gcp)"),
			mcp.Required(),
		),
	)
	s.AddTool(validateRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		rulesPath := request.GetString("rules_path", "")
		provider := request.GetString("provider", "")
		args := []string{"security", "scout-suite", "validate-rules", rulesPath, "--provider", provider}
		return executeShipCommand(args)
	})

	// Scout Suite compare scans tool
	compareScansTool := mcp.NewTool("scout_suite_compare_scans",
		mcp.WithDescription("Compare two Scout Suite scan results"),
		mcp.WithString("baseline_scan",
			mcp.Description("Path to baseline scan results"),
			mcp.Required(),
		),
		mcp.WithString("current_scan",
			mcp.Description("Path to current scan results"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for comparison report"),
		),
	)
	s.AddTool(compareScansTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		baselineScan := request.GetString("baseline_scan", "")
		currentScan := request.GetString("current_scan", "")
		args := []string{"security", "scout-suite", "compare", baselineScan, currentScan}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// Scout Suite get version tool
	getVersionTool := mcp.NewTool("scout_suite_get_version",
		mcp.WithDescription("Get Scout Suite version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "scout-suite", "--version"}
		return executeShipCommand(args)
	})
}