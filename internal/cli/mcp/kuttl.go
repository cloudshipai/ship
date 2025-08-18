package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddKuttlTools adds KUTTL (Kubernetes Test Tool) MCP tool implementations using direct Dagger calls
func AddKuttlTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addKuttlToolsDirect(s)
}

// addKuttlToolsDirect adds KUTTL tools using direct Dagger module calls
func addKuttlToolsDirect(s *server.MCPServer) {
	// KUTTL test tool
	testTool := mcp.NewTool("kuttl_test",
		mcp.WithDescription("Run KUTTL tests"),
		mcp.WithString("path",
			mcp.Description("Path to test directory"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithNumber("parallel",
			mcp.Description("Number of parallel tests to run"),
		),
		mcp.WithBoolean("skip_delete",
			mcp.Description("Skip deleting test resources after test"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Validate tests without running them"),
		),
		mcp.WithString("config",
			mcp.Description("Path to KUTTL config file"),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace to run tests in"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Test timeout in seconds"),
		),
	)
	s.AddTool(testTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKuttlModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		kubeconfig := request.GetString("kubeconfig", "")
		parallel := request.GetInt("parallel", 0)
		skipDelete := request.GetBool("skip_delete", false)
		dryRun := request.GetBool("dry_run", false)

		// Choose appropriate method
		var output string
		if dryRun {
			// Validate tests
			output, err = module.ValidateTest(ctx, path)
		} else {
			// Run tests
			output, err = module.RunTest(ctx, path, kubeconfig, parallel, skipDelete)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kuttl test failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// KUTTL test with kind
	testKindTool := mcp.NewTool("kuttl_test_kind",
		mcp.WithDescription("Run KUTTL tests with kind cluster"),
		mcp.WithString("path",
			mcp.Description("Path to test directory"),
			mcp.Required(),
		),
		mcp.WithString("kind_config",
			mcp.Description("Path to kind configuration file"),
		),
		mcp.WithNumber("parallel",
			mcp.Description("Number of parallel tests to run"),
		),
		mcp.WithString("kind_context",
			mcp.Description("Kind context name"),
		),
		mcp.WithBoolean("skip_cluster_delete",
			mcp.Description("Skip deleting kind cluster after test"),
		),
	)
	s.AddTool(testKindTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKuttlModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		kindConfig := request.GetString("kind_config", "")
		parallel := request.GetInt("parallel", 0)

		// Run tests with kind
		output, err := module.RunTestWithKind(ctx, path, kindConfig, parallel)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kuttl test with kind failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// KUTTL version
	versionTool := mcp.NewTool("kuttl_version",
		mcp.WithDescription("Get KUTTL version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKuttlModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get kuttl version: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// KUTTL help
	helpTool := mcp.NewTool("kuttl_help",
		mcp.WithDescription("Get KUTTL help"),
		mcp.WithString("command",
			mcp.Description("Command to get help for (optional)"),
		),
	)
	s.AddTool(helpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKuttlModule(client)

		// Get command to get help for
		command := request.GetString("command", "")

		// Get help
		output, err := module.GetHelp(ctx, command)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get kuttl help: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}