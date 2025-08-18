package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddConftestTools adds Conftest (OPA policy testing) MCP tool implementations
func AddConftestTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Conftest test policies tool
	testTool := mcp.NewTool("conftest_test",
		mcp.WithDescription("Test configuration files against OPA policies"),
		mcp.WithString("input_file",
			mcp.Description("Path to configuration file or directory to test"),
			mcp.Required(),
		),
		mcp.WithString("policy",
			mcp.Description("Path to policy directory (default: policy)"),
		),
		mcp.WithString("namespace",
			mcp.Description("Override default namespace"),
		),
		mcp.WithBoolean("all_namespaces",
			mcp.Description("Look in all namespaces"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "tap", "junit", "github"),
		),
		mcp.WithString("parser",
			mcp.Description("Parser to use for input files"),
			mcp.Enum("yaml", "json", "toml", "hcl1", "hcl2", "dockerfile"),
		),
	)
	s.AddTool(testTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputFile := request.GetString("input_file", "")
		args := []string{"conftest", "test", inputFile}
		
		if policy := request.GetString("policy", ""); policy != "" {
			args = append(args, "--policy", policy)
		}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if request.GetBool("all_namespaces", false) {
			args = append(args, "--all-namespaces")
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if parser := request.GetString("parser", ""); parser != "" {
			args = append(args, "--parser", parser)
		}
		
		return executeShipCommand(args)
	})

	// Conftest verify policies tool
	verifyTool := mcp.NewTool("conftest_verify",
		mcp.WithDescription("Run policy unit tests"),
		mcp.WithString("policy",
			mcp.Description("Path to policy directory"),
		),
		mcp.WithBoolean("show_builtin_errors",
			mcp.Description("Show parsing errors (recommended)"),
		),
	)
	s.AddTool(verifyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"conftest", "verify"}
		
		if policy := request.GetString("policy", ""); policy != "" {
			args = append(args, "--policy", policy)
		}
		if request.GetBool("show_builtin_errors", false) {
			args = append(args, "--show-builtin-errors")
		}
		
		return executeShipCommand(args)
	})

	// Conftest parse configuration files tool
	parseTool := mcp.NewTool("conftest_parse",
		mcp.WithDescription("Parse and print structured data from input files"),
		mcp.WithString("input_file",
			mcp.Description("Path to configuration file to parse"),
			mcp.Required(),
		),
		mcp.WithString("parser",
			mcp.Description("Parser to use for input files"),
			mcp.Enum("yaml", "json", "toml", "hcl1", "hcl2", "dockerfile"),
		),
	)
	s.AddTool(parseTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputFile := request.GetString("input_file", "")
		args := []string{"conftest", "parse", inputFile}
		
		if parser := request.GetString("parser", ""); parser != "" {
			args = append(args, "--parser", parser)
		}
		
		return executeShipCommand(args)
	})

	// Conftest push policies to OCI registry
	pushTool := mcp.NewTool("conftest_push",
		mcp.WithDescription("Push OPA policy bundles to OCI registry"),
		mcp.WithString("registry_url",
			mcp.Description("OCI registry URL to push policies to"),
			mcp.Required(),
		),
		mcp.WithString("policy_dir",
			mcp.Description("Policy directory to push (optional)"),
		),
		mcp.WithString("tag",
			mcp.Description("Tag for the policy bundle (optional)"),
		),
	)
	s.AddTool(pushTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		registryURL := request.GetString("registry_url", "")
		args := []string{"conftest", "push", registryURL}
		
		if policyDir := request.GetString("policy_dir", ""); policyDir != "" {
			args = append(args, policyDir)
		}
		
		return executeShipCommand(args)
	})

	// Conftest get version tool
	getVersionTool := mcp.NewTool("conftest_get_version",
		mcp.WithDescription("Get Conftest version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"conftest", "--version"}
		return executeShipCommand(args)
	})
}