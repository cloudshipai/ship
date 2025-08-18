package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCfnNagTools adds CFN Nag (CloudFormation template security scanning) MCP tool implementations
func AddCfnNagTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// CFN Nag scan tool
	scanTool := mcp.NewTool("cfn_nag_scan",
		mcp.WithDescription("Scan CloudFormation templates for security issues using cfn_nag_scan"),
		mcp.WithString("input_path",
			mcp.Description("Path to CloudFormation template file or directory"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text"),
		),
		mcp.WithBoolean("debug",
			mcp.Description("Dump information about rule loading"),
		),
	)
	s.AddTool(scanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputPath := request.GetString("input_path", "")
		args := []string{"cfn_nag_scan", "--input-path", inputPath}
		if outputFormat := request.GetString("output_format", ""); outputFormat == "json" {
			args = append(args, "--output-format", "json")
		}
		if request.GetBool("debug", false) {
			args = append(args, "--debug")
		}
		return executeShipCommand(args)
	})

	// CFN Nag scan with profile tool
	scanWithProfileTool := mcp.NewTool("cfn_nag_scan_with_profile",
		mcp.WithDescription("Scan CloudFormation template with specific rule profile"),
		mcp.WithString("input_path",
			mcp.Description("Path to CloudFormation template file or directory"),
			mcp.Required(),
		),
		mcp.WithString("profile_path",
			mcp.Description("Path to profile file specifying rules to apply"),
		),
		mcp.WithString("deny_list_path",
			mcp.Description("Path to deny list file specifying rules to ignore"),
		),
	)
	s.AddTool(scanWithProfileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputPath := request.GetString("input_path", "")
		args := []string{"cfn_nag_scan", "--input-path", inputPath}
		if profilePath := request.GetString("profile_path", ""); profilePath != "" {
			args = append(args, "--profile-path", profilePath)
		}
		if denyListPath := request.GetString("deny_list_path", ""); denyListPath != "" {
			args = append(args, "--deny-list-path", denyListPath)
		}
		return executeShipCommand(args)
	})

	// CFN Nag scan with parameters tool
	scanWithParametersTool := mcp.NewTool("cfn_nag_scan_with_parameters",
		mcp.WithDescription("Scan CloudFormation template with parameter values"),
		mcp.WithString("input_path",
			mcp.Description("Path to CloudFormation template file or directory"),
			mcp.Required(),
		),
		mcp.WithString("parameter_values_path",
			mcp.Description("Path to JSON file with parameter values"),
		),
		mcp.WithString("condition_values_path",
			mcp.Description("Path to JSON file with condition values"),
		),
		mcp.WithString("rule_arguments",
			mcp.Description("Custom rule thresholds (e.g., spcm_threshold:100)"),
		),
	)
	s.AddTool(scanWithParametersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputPath := request.GetString("input_path", "")
		args := []string{"cfn_nag_scan", "--input-path", inputPath}
		if parameterValuesPath := request.GetString("parameter_values_path", ""); parameterValuesPath != "" {
			args = append(args, "--parameter-values-path", parameterValuesPath)
		}
		if conditionValuesPath := request.GetString("condition_values_path", ""); conditionValuesPath != "" {
			args = append(args, "--condition-values-path", conditionValuesPath)
		}
		if ruleArguments := request.GetString("rule_arguments", ""); ruleArguments != "" {
			args = append(args, "--rule-arguments", ruleArguments)
		}
		return executeShipCommand(args)
	})

	// CFN Nag list rules tool
	listRulesTool := mcp.NewTool("cfn_nag_list_rules",
		mcp.WithDescription("List all available CFN Nag rules"),
	)
	s.AddTool(listRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cfn_nag_rules"}
		return executeShipCommand(args)
	})

	// SPCM scan tool (Stelligent Policy Complexity Metrics)
	spcmScanTool := mcp.NewTool("cfn_nag_spcm_scan",
		mcp.WithDescription("Generate Stelligent Policy Complexity Metrics report"),
		mcp.WithString("input_path",
			mcp.Description("Path to CloudFormation template file or directory"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for report"),
			mcp.Enum("json", "html"),
		),
	)
	s.AddTool(spcmScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputPath := request.GetString("input_path", "")
		args := []string{"spcm_scan", "--input-path", inputPath}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output-format", outputFormat)
		}
		return executeShipCommand(args)
	})

	// CFN Nag get version tool
	getVersionTool := mcp.NewTool("cfn_nag_get_version",
		mcp.WithDescription("Get cfn_nag version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cfn_nag_scan", "--version"}
		return executeShipCommand(args)
	})
}