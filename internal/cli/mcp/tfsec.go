package mcp

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTfsecTools adds tfsec Terraform security scanning tools to the MCP server
func AddTfsecTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTfsecToolsDirect(s)
}

// addTfsecToolsDirect adds tfsec tools using direct Dagger module calls
func addTfsecToolsDirect(s *server.MCPServer) {
	// Scan directory
	scanDirTool := mcp.NewTool("tfsec_scan_directory",
		mcp.WithDescription("Scan Terraform directory for security issues"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
		),
		mcp.WithString("format",
			mcp.Description("Output format (json, sarif, junit, html, csv)"),
		),
	)
	s.AddTool(scanDirTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		format := request.GetString("format", "json")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewTfsecModule(client)
		result, err := module.ScanDirectory(ctx, dir, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Scan with severity
	scanSeverityTool := mcp.NewTool("tfsec_scan_with_severity",
		mcp.WithDescription("Scan Terraform with minimum severity threshold"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level"),
			mcp.Required(),
		),
	)
	s.AddTool(scanSeverityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		severity := request.GetString("severity", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewTfsecModule(client)
		result, err := module.ScanWithSeverity(ctx, dir, severity)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Scan with excludes
	scanExcludeTool := mcp.NewTool("tfsec_scan_with_excludes",
		mcp.WithDescription("Scan Terraform with excluded checks"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
		),
		mcp.WithString("excludes",
			mcp.Description("Comma-separated list of check IDs to exclude"),
		),
	)
	s.AddTool(scanExcludeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		excludesStr := request.GetString("excludes", "")
		
		var excludes []string
		if excludesStr != "" {
			excludes = strings.Split(excludesStr, ",")
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewTfsecModule(client)
		result, err := module.ScanWithExcludes(ctx, dir, excludes)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Scan with config
	scanConfigTool := mcp.NewTool("tfsec_scan_with_config",
		mcp.WithDescription("Scan Terraform using a config file"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
		),
		mcp.WithString("config_path",
			mcp.Description("Path to tfsec config file"),
		),
	)
	s.AddTool(scanConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		configPath := request.GetString("config_path", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewTfsecModule(client)
		result, err := module.ScanWithConfig(ctx, dir, configPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Validate tfvars
	validateTfvarsTool := mcp.NewTool("tfsec_validate_tfvars",
		mcp.WithDescription("Validate Terraform variable files"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
		),
		mcp.WithString("tfvars_file",
			mcp.Description("Path to .tfvars file"),
		),
	)
	s.AddTool(validateTfvarsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		tfvarsFile := request.GetString("tfvars_file", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewTfsecModule(client)
		result, err := module.ValidateTfvars(ctx, dir, tfvarsFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("validation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Generate report
	generateReportTool := mcp.NewTool("tfsec_generate_report",
		mcp.WithDescription("Generate security scan report in various formats"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
		),
		mcp.WithString("report_type",
			mcp.Description("Report format"),
			mcp.Required(),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		reportType := request.GetString("report_type", "json")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewTfsecModule(client)
		result, err := module.GenerateReport(ctx, dir, reportType)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("report generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Get version
	versionTool := mcp.NewTool("tfsec_get_version",
		mcp.WithDescription("Get tfsec version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewTfsecModule(client)
		result, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get version: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}