package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKubeHunterTools adds Kube-hunter (Kubernetes penetration testing) MCP tool implementations
func AddKubeHunterTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Kube-hunter scan cluster tool
	scanClusterTool := mcp.NewTool("kube_hunter_scan_cluster",
		mcp.WithDescription("Scan Kubernetes cluster for security vulnerabilities"),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(scanClusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kube-hunter", "--remote"}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Kube-hunter scan network tool
	scanNetworkTool := mcp.NewTool("kube_hunter_scan_network",
		mcp.WithDescription("Scan network for Kubernetes services"),
		mcp.WithString("cidr",
			mcp.Description("CIDR range to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(scanNetworkTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cidr := request.GetString("cidr", "")
		args := []string{"security", "kube-hunter", "--cidr", cidr}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Kube-hunter scan pod tool
	scanPodTool := mcp.NewTool("kube_hunter_scan_pod",
		mcp.WithDescription("Scan from within a Kubernetes pod"),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(scanPodTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kube-hunter", "--pod"}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Kube-hunter scan with custom hunters tool
	scanWithHuntersTool := mcp.NewTool("kube_hunter_scan_with_hunters",
		mcp.WithDescription("Scan using specific hunters"),
		mcp.WithString("hunters",
			mcp.Description("Comma-separated list of hunters to run"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(scanWithHuntersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		hunters := request.GetString("hunters", "")
		args := []string{"security", "kube-hunter", "--include-hunter", hunters}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		return executeShipCommand(args)
	})

	// Kube-hunter scan with severity filter tool
	scanWithSeverityTool := mcp.NewTool("kube_hunter_scan_with_severity",
		mcp.WithDescription("Scan and filter results by severity"),
		mcp.WithString("min_severity",
			mcp.Description("Minimum severity level to report"),
			mcp.Required(),
			mcp.Enum("low", "medium", "high"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(scanWithSeverityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		minSeverity := request.GetString("min_severity", "")
		args := []string{"security", "kube-hunter", "--severity", minSeverity}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		return executeShipCommand(args)
	})

	// Kube-hunter get version tool
	getVersionTool := mcp.NewTool("kube_hunter_get_version",
		mcp.WithDescription("Get kube-hunter version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kube-hunter", "--version"}
		return executeShipCommand(args)
	})
}