package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCfnNagTools adds CFN Nag (CloudFormation template security scanning) MCP tool implementations
func AddCfnNagTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// CFN Nag scan template tool
	scanTemplateTool := mcp.NewTool("cfn_nag_scan_template",
		mcp.WithDescription("Scan CloudFormation template for security issues using CFN Nag"),
		mcp.WithString("template_path",
			mcp.Description("Path to CloudFormation template file"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (json, txt, csv)"),
		),
	)
	s.AddTool(scanTemplateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templatePath := request.GetString("template_path", "")
		args := []string{"security", "cfn-nag", "scan", templatePath}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output-format", outputFormat)
		}
		return executeShipCommand(args)
	})

	// CFN Nag scan directory tool
	scanDirectoryTool := mcp.NewTool("cfn_nag_scan_directory",
		mcp.WithDescription("Scan directory of CloudFormation templates using CFN Nag"),
		mcp.WithString("directory",
			mcp.Description("Directory containing CloudFormation templates"),
			mcp.Required(),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively scan subdirectories"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (json, txt, csv)"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "cfn-nag", "scan-dir", directory}
		if request.GetBool("recursive", false) {
			args = append(args, "--recursive")
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output-format", outputFormat)
		}
		return executeShipCommand(args)
	})

	// CFN Nag scan with rules tool
	scanWithRulesTool := mcp.NewTool("cfn_nag_scan_with_rules",
		mcp.WithDescription("Scan CloudFormation template with custom rules"),
		mcp.WithString("template_path",
			mcp.Description("Path to CloudFormation template file"),
			mcp.Required(),
		),
		mcp.WithString("rules_directory",
			mcp.Description("Directory containing custom CFN Nag rules"),
			mcp.Required(),
		),
		mcp.WithString("profile_path",
			mcp.Description("Path to CFN Nag profile configuration"),
		),
	)
	s.AddTool(scanWithRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templatePath := request.GetString("template_path", "")
		rulesDirectory := request.GetString("rules_directory", "")
		args := []string{"security", "cfn-nag", "scan", templatePath, "--rules-directory", rulesDirectory}
		if profilePath := request.GetString("profile_path", ""); profilePath != "" {
			args = append(args, "--profile-path", profilePath)
		}
		return executeShipCommand(args)
	})

	// CFN Nag generate report tool
	generateReportTool := mcp.NewTool("cfn_nag_generate_report",
		mcp.WithDescription("Generate comprehensive security report for CloudFormation templates"),
		mcp.WithString("template_path",
			mcp.Description("Path to CloudFormation template file or directory"),
			mcp.Required(),
		),
		mcp.WithString("report_format",
			mcp.Description("Report format (html, json, sarif)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for the report"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templatePath := request.GetString("template_path", "")
		args := []string{"security", "cfn-nag", "report", templatePath}
		if reportFormat := request.GetString("report_format", ""); reportFormat != "" {
			args = append(args, "--format", reportFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// CFN Nag validate rules tool
	validateRulesTool := mcp.NewTool("cfn_nag_validate_rules",
		mcp.WithDescription("Validate custom CFN Nag rules"),
		mcp.WithString("rules_directory",
			mcp.Description("Directory containing custom CFN Nag rules"),
			mcp.Required(),
		),
	)
	s.AddTool(validateRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		rulesDirectory := request.GetString("rules_directory", "")
		args := []string{"security", "cfn-nag", "validate-rules", rulesDirectory}
		return executeShipCommand(args)
	})

	// CFN Nag list rules tool
	listRulesTool := mcp.NewTool("cfn_nag_list_rules",
		mcp.WithDescription("List all available CFN Nag rules"),
		mcp.WithBoolean("show_descriptions",
			mcp.Description("Show detailed rule descriptions"),
		),
	)
	s.AddTool(listRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "cfn-nag", "list-rules"}
		if request.GetBool("show_descriptions", false) {
			args = append(args, "--verbose")
		}
		return executeShipCommand(args)
	})

	// CFN Nag get version tool
	getVersionTool := mcp.NewTool("cfn_nag_get_version",
		mcp.WithDescription("Get CFN Nag version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "cfn-nag", "--version"}
		return executeShipCommand(args)
	})
}