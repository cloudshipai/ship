package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCfnNagTools adds CFN Nag (CloudFormation template security scanning) MCP tool implementations
func AddCfnNagTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addCfnNagToolsDirect(s)
}

// addCfnNagToolsDirect implements direct Dagger calls for CFN Nag tools
func addCfnNagToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		inputPath := request.GetString("input_path", "")
		outputFormat := request.GetString("output_format", "")
		debug := request.GetBool("debug", false)

		// Create CFN Nag module and scan
		cfnNagModule := modules.NewCfnNagModule(client)
		result, err := cfnNagModule.Scan(ctx, inputPath, outputFormat, debug)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cfn-nag scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		inputPath := request.GetString("input_path", "")
		profilePath := request.GetString("profile_path", "")
		denyListPath := request.GetString("deny_list_path", "")

		// Create CFN Nag module and scan with profile
		cfnNagModule := modules.NewCfnNagModule(client)
		result, err := cfnNagModule.ScanWithProfile(ctx, inputPath, profilePath, denyListPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cfn-nag scan with profile failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		inputPath := request.GetString("input_path", "")
		parameterValuesPath := request.GetString("parameter_values_path", "")
		conditionValuesPath := request.GetString("condition_values_path", "")
		ruleArguments := request.GetString("rule_arguments", "")

		// Create CFN Nag module and scan with parameters
		cfnNagModule := modules.NewCfnNagModule(client)
		result, err := cfnNagModule.ScanWithParameters(ctx, inputPath, parameterValuesPath, conditionValuesPath, ruleArguments)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cfn-nag scan with parameters failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// CFN Nag list rules tool
	listRulesTool := mcp.NewTool("cfn_nag_list_rules",
		mcp.WithDescription("List all available CFN Nag rules"),
	)
	s.AddTool(listRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create CFN Nag module and list rules
		cfnNagModule := modules.NewCfnNagModule(client)
		result, err := cfnNagModule.ListRules(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cfn-nag list rules failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		inputPath := request.GetString("input_path", "")
		outputFormat := request.GetString("output_format", "")

		// Create CFN Nag module and run SPCM scan
		cfnNagModule := modules.NewCfnNagModule(client)
		result, err := cfnNagModule.SPCMScan(ctx, inputPath, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cfn-nag SPCM scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// CFN Nag get version tool
	getVersionTool := mcp.NewTool("cfn_nag_get_version",
		mcp.WithDescription("Get cfn_nag version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create CFN Nag module and get version
		cfnNagModule := modules.NewCfnNagModule(client)
		result, err := cfnNagModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cfn-nag get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}