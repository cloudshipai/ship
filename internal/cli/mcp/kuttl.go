package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKuttlTools adds KUTTL (Kubernetes testing framework) MCP tool implementations using real CLI commands
func AddKuttlTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// KUTTL test tool
	testTool := mcp.NewTool("kuttl_test",
		mcp.WithDescription("Run Kubernetes tests using kubectl kuttl test"),
		mcp.WithString("test_suite",
			mcp.Description("Path to KUTTL test suite directory"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Path to test settings file"),
		),
		mcp.WithString("artifacts_dir",
			mcp.Description("Directory to output test artifacts and logs"),
		),
		mcp.WithString("crd_dir",
			mcp.Description("Directory containing CustomResourceDefinitions to apply before tests"),
		),
		mcp.WithString("manifest_dir",
			mcp.Description("Directory containing manifests to apply before tests"),
		),
		mcp.WithBoolean("start_kind",
			mcp.Description("Start a KIND cluster for testing"),
		),
		mcp.WithBoolean("start_control_plane",
			mcp.Description("Start a local Kubernetes control plane"),
		),
		mcp.WithString("parallel",
			mcp.Description("Maximum number of concurrent tests (default 8)"),
		),
		mcp.WithString("verbosity",
			mcp.Description("Logging verbosity level"),
			mcp.Enum("1", "2"),
		),
	)
	s.AddTool(testTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		testSuite := request.GetString("test_suite", "")
		args := []string{"kubectl", "kuttl", "test", testSuite}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if artifactsDir := request.GetString("artifacts_dir", ""); artifactsDir != "" {
			args = append(args, "--artifacts-dir", artifactsDir)
		}
		if crdDir := request.GetString("crd_dir", ""); crdDir != "" {
			args = append(args, "--crd-dir", crdDir)
		}
		if manifestDir := request.GetString("manifest_dir", ""); manifestDir != "" {
			args = append(args, "--manifest-dir", manifestDir)
		}
		if request.GetBool("start_kind", false) {
			args = append(args, "--start-kind")
		}
		if request.GetBool("start_control_plane", false) {
			args = append(args, "--start-control-plane")
		}
		if parallel := request.GetString("parallel", ""); parallel != "" {
			args = append(args, "--parallel", parallel)
		}
		if verbosity := request.GetString("verbosity", ""); verbosity != "" {
			if verbosity == "1" {
				args = append(args, "-v")
			} else if verbosity == "2" {
				args = append(args, "-vv")
			}
		}
		
		return executeShipCommand(args)
	})

	// KUTTL test with KIND cluster
	testKindTool := mcp.NewTool("kuttl_test_kind",
		mcp.WithDescription("Run KUTTL tests with automated KIND cluster setup"),
		mcp.WithString("test_suite",
			mcp.Description("Path to KUTTL test suite directory"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Path to test settings file"),
		),
		mcp.WithString("artifacts_dir",
			mcp.Description("Directory to output test artifacts and logs"),
		),
		mcp.WithString("parallel",
			mcp.Description("Maximum number of concurrent tests"),
		),
	)
	s.AddTool(testKindTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		testSuite := request.GetString("test_suite", "")
		args := []string{"kubectl", "kuttl", "test", "--start-kind", testSuite}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if artifactsDir := request.GetString("artifacts_dir", ""); artifactsDir != "" {
			args = append(args, "--artifacts-dir", artifactsDir)
		}
		if parallel := request.GetString("parallel", ""); parallel != "" {
			args = append(args, "--parallel", parallel)
		}
		
		return executeShipCommand(args)
	})

	// KUTTL version tool
	versionTool := mcp.NewTool("kuttl_version",
		mcp.WithDescription("Get KUTTL version information using kubectl kuttl version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubectl", "kuttl", "version"}
		return executeShipCommand(args)
	})

	// KUTTL help tool
	helpTool := mcp.NewTool("kuttl_help",
		mcp.WithDescription("Get KUTTL help information"),
		mcp.WithString("command",
			mcp.Description("Get help for specific command (test)"),
		),
	)
	s.AddTool(helpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubectl", "kuttl", "help"}
		
		if command := request.GetString("command", ""); command != "" {
			args = append(args, command)
		}
		
		return executeShipCommand(args)
	})
}