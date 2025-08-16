package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddInfrascanTools adds Infrascan MCP tool implementations
func AddInfrascanTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Infrascan scan infrastructure tool
	scanInfrastructureTool := mcp.NewTool("infrascan_scan_infrastructure",
		mcp.WithDescription("Scan infrastructure for security vulnerabilities and misconfigurations"),
		mcp.WithString("target",
			mcp.Description("Target infrastructure (directory, cloud account, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("scan_type",
			mcp.Description("Type of scan (iac, cloud, network)"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (json, yaml, table)"),
		),
	)
	s.AddTool(scanInfrastructureTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"security", "infrascan", "scan", target}
		if scanType := request.GetString("scan_type", ""); scanType != "" {
			args = append(args, "--type", scanType)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Infrascan compliance check tool
	complianceCheckTool := mcp.NewTool("infrascan_compliance_check",
		mcp.WithDescription("Check infrastructure compliance against security frameworks"),
		mcp.WithString("target",
			mcp.Description("Target infrastructure to check"),
			mcp.Required(),
		),
		mcp.WithString("framework",
			mcp.Description("Compliance framework (cis, nist, pci-dss, sox)"),
		),
	)
	s.AddTool(complianceCheckTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"security", "infrascan", "compliance", target}
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "--framework", framework)
		}
		return executeShipCommand(args)
	})

	// Infrascan generate report tool
	generateReportTool := mcp.NewTool("infrascan_generate_report",
		mcp.WithDescription("Generate comprehensive security report for infrastructure"),
		mcp.WithString("scan_results",
			mcp.Description("Path to scan results file"),
			mcp.Required(),
		),
		mcp.WithString("report_format",
			mcp.Description("Report format (html, pdf, json)"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scanResults := request.GetString("scan_results", "")
		args := []string{"security", "infrascan", "report", scanResults}
		if format := request.GetString("report_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Infrascan baseline tool
	baselineTool := mcp.NewTool("infrascan_baseline",
		mcp.WithDescription("Create security baseline for infrastructure configuration"),
		mcp.WithString("target",
			mcp.Description("Target infrastructure for baseline"),
			mcp.Required(),
		),
		mcp.WithString("baseline_name",
			mcp.Description("Name for the baseline configuration"),
		),
	)
	s.AddTool(baselineTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"security", "infrascan", "baseline", target}
		if name := request.GetString("baseline_name", ""); name != "" {
			args = append(args, "--name", name)
		}
		return executeShipCommand(args)
	})

	// Infrascan get version tool
	getVersionTool := mcp.NewTool("infrascan_get_version",
		mcp.WithDescription("Get Infrascan version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "infrascan", "--version"}
		return executeShipCommand(args)
	})
}