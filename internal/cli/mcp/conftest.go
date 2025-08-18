package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddConftestTools adds Conftest (OPA policy testing) MCP tool implementations
func AddConftestTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addConftestToolsDirect(s)
}

// addConftestToolsDirect implements direct Dagger calls for Conftest tools
func addConftestToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		inputFile := request.GetString("input_file", "")
		policy := request.GetString("policy", "")
		namespace := request.GetString("namespace", "")
		allNamespaces := request.GetBool("all_namespaces", false)
		output := request.GetString("output", "")
		parser := request.GetString("parser", "")

		// Create Conftest module and test with options
		conftestModule := modules.NewConftestModule(client)
		result, err := conftestModule.TestWithOptions(ctx, inputFile, policy, namespace, allNamespaces, output, parser)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("conftest test failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		policy := request.GetString("policy", "")
		showBuiltinErrors := request.GetBool("show_builtin_errors", false)

		// Create Conftest module and verify with options
		conftestModule := modules.NewConftestModule(client)
		result, err := conftestModule.VerifyWithOptions(ctx, policy, showBuiltinErrors)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("conftest verify failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		inputFile := request.GetString("input_file", "")
		parser := request.GetString("parser", "")

		// Create Conftest module and parse file
		conftestModule := modules.NewConftestModule(client)
		result, err := conftestModule.ParseFile(ctx, inputFile, parser)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("conftest parse failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		registryURL := request.GetString("registry_url", "")
		policyDir := request.GetString("policy_dir", "")

		// Create Conftest module and push policies
		conftestModule := modules.NewConftestModule(client)
		result, err := conftestModule.PushPolicies(ctx, registryURL, policyDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("conftest push failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Conftest get version tool
	getVersionTool := mcp.NewTool("conftest_get_version",
		mcp.WithDescription("Get Conftest version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create Conftest module and get version
		conftestModule := modules.NewConftestModule(client)
		result, err := conftestModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("conftest get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}