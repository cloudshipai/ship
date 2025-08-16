package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCloudsplainingTools adds Cloudsplaining (AWS IAM policy scanner) MCP tool implementations
func AddCloudsplainingTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Cloudsplaining scan account authorization tool
	scanAccountTool := mcp.NewTool("cloudsplaining_scan_account",
		mcp.WithDescription("Scan AWS account IAM authorization using Cloudsplaining"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use for scanning"),
			mcp.Required(),
		),
	)
	s.AddTool(scanAccountTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profile := request.GetString("profile", "")
		args := []string{"security", "cloudsplaining", "--scan-account", "--profile", profile}
		return executeShipCommand(args)
	})

	// Cloudsplaining scan policy file tool
	scanPolicyFileTool := mcp.NewTool("cloudsplaining_scan_policy_file",
		mcp.WithDescription("Scan IAM policy file using Cloudsplaining"),
		mcp.WithString("policy_path",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
	)
	s.AddTool(scanPolicyFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "cloudsplaining", "--scan-policy", policyPath}
		return executeShipCommand(args)
	})

	// Cloudsplaining create report tool
	createReportTool := mcp.NewTool("cloudsplaining_create_report",
		mcp.WithDescription("Create HTML report from Cloudsplaining scan results"),
		mcp.WithString("results_path",
			mcp.Description("Path to scan results JSON file"),
			mcp.Required(),
		),
		mcp.WithString("output_path",
			mcp.Description("Output path for HTML report"),
		),
	)
	s.AddTool(createReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resultsPath := request.GetString("results_path", "")
		args := []string{"security", "cloudsplaining", "--create-report", resultsPath}
		if outputPath := request.GetString("output_path", ""); outputPath != "" {
			args = append(args, "--output", outputPath)
		}
		return executeShipCommand(args)
	})

	// Cloudsplaining scan with minimization tool
	scanWithMinimizationTool := mcp.NewTool("cloudsplaining_scan_minimized",
		mcp.WithDescription("Scan AWS account with policy minimization suggestions"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use for scanning"),
			mcp.Required(),
		),
		mcp.WithString("minimize_statement_id",
			mcp.Description("Statement ID to minimize"),
		),
	)
	s.AddTool(scanWithMinimizationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profile := request.GetString("profile", "")
		args := []string{"security", "cloudsplaining", "--scan-account", "--minimize", "--profile", profile}
		if statementId := request.GetString("minimize_statement_id", ""); statementId != "" {
			args = append(args, "--statement-id", statementId)
		}
		return executeShipCommand(args)
	})
}