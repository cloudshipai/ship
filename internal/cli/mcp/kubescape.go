package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKubescapeTools adds Kubescape MCP tool implementations
func AddKubescapeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Kubescape scan cluster tool
	scanClusterTool := mcp.NewTool("kubescape_scan_cluster",
		mcp.WithDescription("Scan Kubernetes cluster using Kubescape"),
		mcp.WithString("framework",
			mcp.Description("Security framework to use"),
			mcp.Enum("nsa", "mitre", "cis", "all"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(scanClusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kubescape", "cluster"}
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "--framework", framework)
		}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		return executeShipCommand(args)
	})

	// Kubescape scan manifests tool
	scanManifestsTool := mcp.NewTool("kubescape_scan_manifests",
		mcp.WithDescription("Scan Kubernetes manifests using Kubescape"),
		mcp.WithString("manifests_dir",
			mcp.Description("Directory containing Kubernetes manifests"),
		),
		mcp.WithString("framework",
			mcp.Description("Security framework to use"),
			mcp.Enum("nsa", "mitre", "cis", "all"),
		),
	)
	s.AddTool(scanManifestsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kubescape", "manifests"}
		if manifestsDir := request.GetString("manifests_dir", ""); manifestsDir != "" {
			args = append(args, manifestsDir)
		}
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "--framework", framework)
		}
		return executeShipCommand(args)
	})

	// Kubescape scan Helm chart tool
	scanHelmTool := mcp.NewTool("kubescape_scan_helm",
		mcp.WithDescription("Scan Helm chart using Kubescape"),
		mcp.WithString("chart_path",
			mcp.Description("Path to Helm chart directory"),
		),
		mcp.WithString("framework",
			mcp.Description("Security framework to use"),
			mcp.Enum("nsa", "mitre", "cis", "all"),
		),
	)
	s.AddTool(scanHelmTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kubescape", "helm"}
		if chartPath := request.GetString("chart_path", ""); chartPath != "" {
			args = append(args, chartPath)
		}
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "--framework", framework)
		}
		return executeShipCommand(args)
	})

	// Kubescape scan repository tool
	scanRepositoryTool := mcp.NewTool("kubescape_scan_repository",
		mcp.WithDescription("Scan repository for Kubernetes security issues using Kubescape"),
		mcp.WithString("repo_path",
			mcp.Description("Path to repository directory"),
		),
		mcp.WithString("framework",
			mcp.Description("Security framework to use"),
			mcp.Enum("nsa", "mitre", "cis", "all"),
		),
	)
	s.AddTool(scanRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kubescape", "repo"}
		if repoPath := request.GetString("repo_path", ""); repoPath != "" {
			args = append(args, repoPath)
		}
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "--framework", framework)
		}
		return executeShipCommand(args)
	})

	// Kubescape generate report tool
	generateReportTool := mcp.NewTool("kubescape_generate_report",
		mcp.WithDescription("Generate security report using Kubescape"),
		mcp.WithString("input_path",
			mcp.Description("Path to scan results or manifests"),
		),
		mcp.WithString("format",
			mcp.Description("Report format"),
			mcp.Enum("json", "junit", "pdf", "html"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kubescape", "report"}
		if inputPath := request.GetString("input_path", ""); inputPath != "" {
			args = append(args, inputPath)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})
}