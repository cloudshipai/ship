package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddScoutSuiteTools adds Scout Suite MCP tool implementations using real scout CLI commands
func AddScoutSuiteTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Scout Suite scan AWS tool
	scanAWSTool := mcp.NewTool("scout_suite_scan_aws",
		mcp.WithDescription("Scan AWS environment for security issues using real scout CLI"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use for scanning"),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of AWS regions to scan"),
		),
		mcp.WithString("services",
			mcp.Description("Comma-separated list of AWS services to scan"),
		),
		mcp.WithString("report_dir",
			mcp.Description("Directory to save the report"),
		),
		mcp.WithString("exceptions",
			mcp.Description("Path to exceptions file"),
		),
	)
	s.AddTool(scanAWSTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"scout", "aws"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if regions := request.GetString("regions", ""); regions != "" {
			args = append(args, "--regions", regions)
		}
		if services := request.GetString("services", ""); services != "" {
			args = append(args, "--services", services)
		}
		if reportDir := request.GetString("report_dir", ""); reportDir != "" {
			args = append(args, "--report-dir", reportDir)
		}
		if exceptions := request.GetString("exceptions", ""); exceptions != "" {
			args = append(args, "--exceptions", exceptions)
		}
		
		return executeShipCommand(args)
	})

	// Scout Suite scan Azure tool
	scanAzureTool := mcp.NewTool("scout_suite_scan_azure",
		mcp.WithDescription("Scan Azure environment for security issues using real scout CLI"),
		mcp.WithString("subscriptions",
			mcp.Description("Azure subscription IDs to scan (space-separated)"),
		),
		mcp.WithString("tenant_id",
			mcp.Description("Azure tenant ID"),
		),
		mcp.WithString("username",
			mcp.Description("Azure username"),
		),
		mcp.WithString("password",
			mcp.Description("Azure password"),
		),
		mcp.WithBoolean("cli",
			mcp.Description("Use Azure CLI for authentication"),
		),
		mcp.WithBoolean("service_principal",
			mcp.Description("Use service principal authentication"),
		),
		mcp.WithString("report_dir",
			mcp.Description("Directory to save the report"),
		),
	)
	s.AddTool(scanAzureTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"scout", "azure"}
		
		if subscriptions := request.GetString("subscriptions", ""); subscriptions != "" {
			args = append(args, "--subscriptions", subscriptions)
		}
		if tenantID := request.GetString("tenant_id", ""); tenantID != "" {
			args = append(args, "--tenant-id", tenantID)
		}
		if username := request.GetString("username", ""); username != "" {
			args = append(args, "--username", username)
		}
		if password := request.GetString("password", ""); password != "" {
			args = append(args, "--password", password)
		}
		if request.GetBool("cli", false) {
			args = append(args, "--cli")
		}
		if request.GetBool("service_principal", false) {
			args = append(args, "--service-principal")
		}
		if reportDir := request.GetString("report_dir", ""); reportDir != "" {
			args = append(args, "--report-dir", reportDir)
		}
		
		return executeShipCommand(args)
	})

	// Scout Suite scan GCP tool
	scanGCPTool := mcp.NewTool("scout_suite_scan_gcp",
		mcp.WithDescription("Scan Google Cloud Platform for security issues using real scout CLI"),
		mcp.WithString("project_id",
			mcp.Description("GCP project ID to scan"),
		),
		mcp.WithString("folder_id",
			mcp.Description("GCP folder ID to scan"),
		),
		mcp.WithString("organization_id",
			mcp.Description("GCP organization ID to scan"),
		),
		mcp.WithString("service_account",
			mcp.Description("Path to GCP service account key file"),
		),
		mcp.WithBoolean("user_account",
			mcp.Description("Use user account for authentication"),
		),
		mcp.WithBoolean("all_projects",
			mcp.Description("Scan all accessible projects"),
		),
		mcp.WithString("report_dir",
			mcp.Description("Directory to save the report"),
		),
	)
	s.AddTool(scanGCPTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"scout", "gcp"}
		
		if projectID := request.GetString("project_id", ""); projectID != "" {
			args = append(args, "--project-id", projectID)
		}
		if folderID := request.GetString("folder_id", ""); folderID != "" {
			args = append(args, "--folder-id", folderID)
		}
		if orgID := request.GetString("organization_id", ""); orgID != "" {
			args = append(args, "--organization-id", orgID)
		}
		if serviceAccount := request.GetString("service_account", ""); serviceAccount != "" {
			args = append(args, "--service-account", serviceAccount)
		}
		if request.GetBool("user_account", false) {
			args = append(args, "--user-account")
		}
		if request.GetBool("all_projects", false) {
			args = append(args, "--all-projects")
		}
		if reportDir := request.GetString("report_dir", ""); reportDir != "" {
			args = append(args, "--report-dir", reportDir)
		}
		
		return executeShipCommand(args)
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
		provider := request.GetString("provider", "")
		args := []string{"scout", provider, "--serve"}
		
		if reportName := request.GetString("report_name", ""); reportName != "" {
			args = append(args, reportName)
		}
		if host := request.GetString("host", ""); host != "" {
			args = append(args, "--host", host)
		}
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "--port", port)
		}
		
		return executeShipCommand(args)
	})

	// Scout Suite help tool
	helpTool := mcp.NewTool("scout_suite_help",
		mcp.WithDescription("Get Scout Suite help information using real scout CLI"),
		mcp.WithString("provider",
			mcp.Description("Get help for specific provider (aws, azure, gcp)"),
		),
	)
	s.AddTool(helpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"scout", "--help"}
		
		if provider := request.GetString("provider", ""); provider != "" {
			args = []string{"scout", provider, "--help"}
		}
		
		return executeShipCommand(args)
	})
}