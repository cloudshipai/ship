package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddConftestTools adds Conftest (OPA policy testing) MCP tool implementations
func AddConftestTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Conftest test policies tool
	testPolicesTool := mcp.NewTool("conftest_test_policies",
		mcp.WithDescription("Test configuration files against OPA policies using Conftest"),
		mcp.WithString("input_path",
			mcp.Description("Path to configuration file or directory to test"),
			mcp.Required(),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to OPA policy files"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "tap", "junit"),
		),
	)
	s.AddTool(testPolicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputPath := request.GetString("input_path", "")
		args := []string{"security", "conftest", "test", inputPath}
		if policyPath := request.GetString("policy_path", ""); policyPath != "" {
			args = append(args, "--policy", policyPath)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Conftest verify policies tool
	verifyPolicesTool := mcp.NewTool("conftest_verify_policies",
		mcp.WithDescription("Verify OPA policies against test data"),
		mcp.WithString("policy_path",
			mcp.Description("Path to OPA policy files"),
			mcp.Required(),
		),
		mcp.WithString("data_path",
			mcp.Description("Path to test data files"),
		),
	)
	s.AddTool(verifyPolicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "conftest", "verify", "--policy", policyPath}
		if dataPath := request.GetString("data_path", ""); dataPath != "" {
			args = append(args, "--data", dataPath)
		}
		return executeShipCommand(args)
	})

	// Conftest push policies tool
	pushPolicesTool := mcp.NewTool("conftest_push_policies",
		mcp.WithDescription("Push OPA policies to OCI registry"),
		mcp.WithString("policy_path",
			mcp.Description("Path to policy files to push"),
			mcp.Required(),
		),
		mcp.WithString("registry_url",
			mcp.Description("OCI registry URL to push policies to"),
			mcp.Required(),
		),
	)
	s.AddTool(pushPolicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		registryURL := request.GetString("registry_url", "")
		args := []string{"security", "conftest", "push", registryURL, "--policy", policyPath}
		return executeShipCommand(args)
	})

	// Conftest get version tool
	getVersionTool := mcp.NewTool("conftest_get_version",
		mcp.WithDescription("Get Conftest version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "conftest", "--version"}
		return executeShipCommand(args)
	})
}