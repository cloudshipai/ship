package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCustodianTools adds Cloud Custodian (cloud governance engine) MCP tool implementations
func AddCustodianTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		args := []string{"custodian", "run", "-s", "out", policyFile}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		return executeShipCommand(args)
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
		outputDir := request.GetString("output_dir", "out")
		args := []string{"custodian", "run", "--dryrun", "-s", outputDir, policyFile}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		return executeShipCommand(args)
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
		args := []string{"custodian", "validate", policyFile}
		return executeShipCommand(args)
	})

	// Custodian schema tool
	schemaTool := mcp.NewTool("custodian_schema",
		mcp.WithDescription("Get Cloud Custodian policy schema for resource types"),
		mcp.WithString("resource_type",
			mcp.Description("AWS resource type (e.g., ec2, s3, iam)"),
		),
	)
	s.AddTool(schemaTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"custodian", "schema"}
		if resourceType := request.GetString("resource_type", ""); resourceType != "" {
			args = append(args, resourceType)
		}
		return executeShipCommand(args)
	})

	// Custodian get version tool
	getVersionTool := mcp.NewTool("custodian_get_version",
		mcp.WithDescription("Get Cloud Custodian version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"custodian", "version"}
		return executeShipCommand(args)
	})
}