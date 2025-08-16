package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddOpenSCAPTools adds OpenSCAP (security compliance scanning) MCP tool implementations
func AddOpenSCAPTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// OpenSCAP scan system tool
	scanSystemTool := mcp.NewTool("openscap_scan_system",
		mcp.WithDescription("Scan system for security compliance using OpenSCAP"),
		mcp.WithString("profile",
			mcp.Description("Security profile to use (e.g., cis, stig, pci-dss)"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (xml, html, json)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for scan results"),
		),
	)
	s.AddTool(scanSystemTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profile := request.GetString("profile", "")
		args := []string{"security", "openscap", "scan", "--profile", profile}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// OpenSCAP evaluate OVAL tool
	evaluateOVALTool := mcp.NewTool("openscap_evaluate_oval",
		mcp.WithDescription("Evaluate OVAL definitions for system compliance"),
		mcp.WithString("oval_file",
			mcp.Description("Path to OVAL definitions file"),
			mcp.Required(),
		),
		mcp.WithString("variables_file",
			mcp.Description("Path to OVAL variables file"),
		),
	)
	s.AddTool(evaluateOVALTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ovalFile := request.GetString("oval_file", "")
		args := []string{"security", "openscap", "eval", "--oval", ovalFile}
		if variablesFile := request.GetString("variables_file", ""); variablesFile != "" {
			args = append(args, "--variables", variablesFile)
		}
		return executeShipCommand(args)
	})

	// OpenSCAP generate report tool
	generateReportTool := mcp.NewTool("openscap_generate_report",
		mcp.WithDescription("Generate compliance report from scan results"),
		mcp.WithString("results_file",
			mcp.Description("Path to scan results file"),
			mcp.Required(),
		),
		mcp.WithString("report_format",
			mcp.Description("Report format (html, pdf, json)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for the report"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resultsFile := request.GetString("results_file", "")
		args := []string{"security", "openscap", "report", resultsFile}
		if reportFormat := request.GetString("report_format", ""); reportFormat != "" {
			args = append(args, "--format", reportFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// OpenSCAP validate content tool
	validateContentTool := mcp.NewTool("openscap_validate_content",
		mcp.WithDescription("Validate SCAP content for correctness"),
		mcp.WithString("content_file",
			mcp.Description("Path to SCAP content file"),
			mcp.Required(),
		),
		mcp.WithString("content_type",
			mcp.Description("Content type (xccdf, oval, ds)"),
		),
	)
	s.AddTool(validateContentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		contentFile := request.GetString("content_file", "")
		args := []string{"security", "openscap", "validate", contentFile}
		if contentType := request.GetString("content_type", ""); contentType != "" {
			args = append(args, "--type", contentType)
		}
		return executeShipCommand(args)
	})

	// OpenSCAP list profiles tool
	listProfilesTool := mcp.NewTool("openscap_list_profiles",
		mcp.WithDescription("List available security profiles"),
		mcp.WithString("content_file",
			mcp.Description("Path to SCAP content file"),
		),
	)
	s.AddTool(listProfilesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "openscap", "list-profiles"}
		if contentFile := request.GetString("content_file", ""); contentFile != "" {
			args = append(args, contentFile)
		}
		return executeShipCommand(args)
	})

	// OpenSCAP remediate system tool
	remediateSystemTool := mcp.NewTool("openscap_remediate_system",
		mcp.WithDescription("Apply remediation based on scan results"),
		mcp.WithString("results_file",
			mcp.Description("Path to scan results file"),
			mcp.Required(),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform a dry run without making changes"),
		),
	)
	s.AddTool(remediateSystemTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resultsFile := request.GetString("results_file", "")
		args := []string{"security", "openscap", "remediate", resultsFile}
		if request.GetBool("dry_run", false) {
			args = append(args, "--dry-run")
		}
		return executeShipCommand(args)
	})

	// OpenSCAP get version tool
	getVersionTool := mcp.NewTool("openscap_get_version",
		mcp.WithDescription("Get OpenSCAP version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "openscap", "--version"}
		return executeShipCommand(args)
	})
}