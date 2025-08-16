package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKuttlTools adds KUTTL (Kubernetes testing framework) MCP tool implementations
func AddKuttlTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// KUTTL test tool
	testTool := mcp.NewTool("kuttl_test",
		mcp.WithDescription("Run Kubernetes tests using KUTTL framework"),
		mcp.WithString("test_suite",
			mcp.Description("Path to KUTTL test suite directory"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for tests"),
		),
		mcp.WithString("timeout",
			mcp.Description("Test timeout duration (e.g., 5m, 10m)"),
		),
	)
	s.AddTool(testTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		testSuite := request.GetString("test_suite", "")
		args := []string{"kubernetes", "kuttl", "test", testSuite}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if timeout := request.GetString("timeout", ""); timeout != "" {
			args = append(args, "--timeout", timeout)
		}
		return executeShipCommand(args)
	})

	// KUTTL validate tool
	validateTool := mcp.NewTool("kuttl_validate",
		mcp.WithDescription("Validate KUTTL test configuration and manifests"),
		mcp.WithString("test_suite",
			mcp.Description("Path to KUTTL test suite directory"),
			mcp.Required(),
		),
	)
	s.AddTool(validateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		testSuite := request.GetString("test_suite", "")
		args := []string{"kubernetes", "kuttl", "validate", testSuite}
		return executeShipCommand(args)
	})

	// KUTTL init tool
	initTool := mcp.NewTool("kuttl_init",
		mcp.WithDescription("Initialize new KUTTL test suite"),
		mcp.WithString("project_name",
			mcp.Description("Name of the test project"),
			mcp.Required(),
		),
		mcp.WithString("directory",
			mcp.Description("Directory to create test suite in"),
		),
	)
	s.AddTool(initTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName := request.GetString("project_name", "")
		args := []string{"kubernetes", "kuttl", "init", projectName}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, "--dir", dir)
		}
		return executeShipCommand(args)
	})

	// KUTTL generate tool
	generateTool := mcp.NewTool("kuttl_generate",
		mcp.WithDescription("Generate KUTTL test cases from existing Kubernetes manifests"),
		mcp.WithString("manifest_path",
			mcp.Description("Path to Kubernetes manifest files"),
			mcp.Required(),
		),
		mcp.WithString("output_dir",
			mcp.Description("Output directory for generated tests"),
		),
	)
	s.AddTool(generateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		manifestPath := request.GetString("manifest_path", "")
		args := []string{"kubernetes", "kuttl", "generate", manifestPath}
		if outputDir := request.GetString("output_dir", ""); outputDir != "" {
			args = append(args, "--output", outputDir)
		}
		return executeShipCommand(args)
	})

	// KUTTL report tool
	reportTool := mcp.NewTool("kuttl_report",
		mcp.WithDescription("Generate test report from KUTTL test results"),
		mcp.WithString("results_path",
			mcp.Description("Path to KUTTL test results"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Report format (html, json, junit)"),
		),
	)
	s.AddTool(reportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resultsPath := request.GetString("results_path", "")
		args := []string{"kubernetes", "kuttl", "report", resultsPath}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// KUTTL get version tool
	getVersionTool := mcp.NewTool("kuttl_get_version",
		mcp.WithDescription("Get KUTTL version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "kuttl", "--version"}
		return executeShipCommand(args)
	})
}