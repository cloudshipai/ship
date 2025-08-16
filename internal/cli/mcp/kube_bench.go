package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKubeBenchTools adds Kube-bench (Kubernetes CIS benchmark) MCP tool implementations
func AddKubeBenchTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Kube-bench run benchmark tool
	runBenchmarkTool := mcp.NewTool("kube_bench_run_benchmark",
		mcp.WithDescription("Run CIS Kubernetes benchmark using kube-bench"),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("benchmark",
			mcp.Description("Benchmark version to run"),
			mcp.Enum("cis-1.6", "cis-1.20", "cis-1.23", "cis-1.24"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "junit", "text"),
		),
	)
	s.AddTool(runBenchmarkTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kube-bench", "--run"}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		if benchmark := request.GetString("benchmark", ""); benchmark != "" {
			args = append(args, "--benchmark", benchmark)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Kube-bench run specific check tool
	runSpecificCheckTool := mcp.NewTool("kube_bench_run_specific_check",
		mcp.WithDescription("Run specific CIS benchmark check"),
		mcp.WithString("check_id",
			mcp.Description("Specific check ID to run"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "junit", "text"),
		),
	)
	s.AddTool(runSpecificCheckTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		checkID := request.GetString("check_id", "")
		args := []string{"security", "kube-bench", "--check", checkID}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Kube-bench run for specific node type tool
	runForNodeTypeTool := mcp.NewTool("kube_bench_run_for_node_type",
		mcp.WithDescription("Run benchmark for specific Kubernetes node type"),
		mcp.WithString("node_type",
			mcp.Description("Type of Kubernetes node"),
			mcp.Required(),
			mcp.Enum("master", "node", "controlplane", "etcd", "policies"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "junit", "text"),
		),
	)
	s.AddTool(runForNodeTypeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeType := request.GetString("node_type", "")
		args := []string{"security", "kube-bench", "--targets", nodeType}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Kube-bench run with custom config tool
	runWithConfigTool := mcp.NewTool("kube_bench_run_with_config",
		mcp.WithDescription("Run benchmark with custom configuration"),
		mcp.WithString("config_path",
			mcp.Description("Path to custom kube-bench configuration"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(runWithConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configPath := request.GetString("config_path", "")
		args := []string{"security", "kube-bench", "--config", configPath}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		return executeShipCommand(args)
	})

	// Kube-bench get version tool
	getVersionTool := mcp.NewTool("kube_bench_get_version",
		mcp.WithDescription("Get kube-bench version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kube-bench", "--version"}
		return executeShipCommand(args)
	})
}