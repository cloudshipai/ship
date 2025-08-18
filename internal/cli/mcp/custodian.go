package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCustodianTools adds Cloud Custodian (cloud governance engine) MCP tool implementations
func AddCustodianTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addCustodianToolsDirect(s)
}

// addCustodianToolsDirect implements direct Dagger calls for custodian tools
func addCustodianToolsDirect(s *server.MCPServer) {
	// Custodian run policy tool
	runPolicyTool := mcp.NewTool("custodian_run_policy",
		mcp.WithDescription("Run Cloud Custodian policy for cloud governance"),
		mcp.WithString("policy_file",
			mcp.Description("Path to custodian policy YAML file"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("AWS region to run policy in"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(runPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		region := request.GetString("region", "")
		
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create custodian module and run policy
		custodianModule := modules.NewCustodianModule(client)
		result, err := custodianModule.RunPolicy(ctx, policyFile, "out")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("custodian run failed: %v", err)), nil
		}

		// Note: Region parameter not yet supported in Dagger module
		if region != "" {
			result = fmt.Sprintf("Note: --region %s parameter not yet implemented in Dagger module\n\n%s", region, result)
		}

		return mcp.NewToolResultText(result), nil
	})

	// Custodian dry run policy tool
	dryRunTool := mcp.NewTool("custodian_dry_run",
		mcp.WithDescription("Dry run Cloud Custodian policy (preview mode)"),
		mcp.WithString("policy_file",
			mcp.Description("Path to custodian policy YAML file"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("AWS region to run policy in"),
		),
		mcp.WithString("output_dir",
			mcp.Description("Output directory for results (default: out)"),
		),
	)
	s.AddTool(dryRunTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		region := request.GetString("region", "")
		
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create custodian module and run dry run
		custodianModule := modules.NewCustodianModule(client)
		result, err := custodianModule.DryRun(ctx, policyFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("custodian dry run failed: %v", err)), nil
		}

		// Note: Region parameter not yet supported in Dagger module
		if region != "" {
			result = fmt.Sprintf("Note: --region %s parameter not yet implemented in Dagger module\n\n%s", region, result)
		}

		return mcp.NewToolResultText(result), nil
	})

	// Custodian validate policy tool
	validateTool := mcp.NewTool("custodian_validate_policy",
		mcp.WithDescription("Validate Cloud Custodian policy syntax"),
		mcp.WithString("policy_file",
			mcp.Description("Path to custodian policy YAML file"),
			mcp.Required(),
		),
	)
	s.AddTool(validateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create custodian module and validate policy
		custodianModule := modules.NewCustodianModule(client)
		result, err := custodianModule.ValidatePolicy(ctx, policyFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("custodian validate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Custodian schema tool
	schemaTool := mcp.NewTool("custodian_schema",
		mcp.WithDescription("Get Cloud Custodian policy schema for resource types"),
		mcp.WithString("resource_type",
			mcp.Description("AWS resource type (e.g., ec2, s3, iam)"),
		),
	)
	s.AddTool(schemaTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resourceType := request.GetString("resource_type", "")
		
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create custodian module and get schema
		custodianModule := modules.NewCustodianModule(client)
		result, err := custodianModule.Schema(ctx, resourceType)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("custodian schema failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Custodian get version tool
	getVersionTool := mcp.NewTool("custodian_get_version",
		mcp.WithDescription("Get Cloud Custodian version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create custodian module and get version
		custodianModule := modules.NewCustodianModule(client)
		result, err := custodianModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("custodian version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}