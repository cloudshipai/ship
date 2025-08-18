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
			mcp.Enum("nsa", "mitre", "cis"),
		),
		mcp.WithString("kube_context",
			mcp.Description("Kubernetes context to use"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "junit", "pdf", "html", "sarif"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose output"),
		),
		mcp.WithBoolean("submit",
			mcp.Description("Submit results to Armo platform"),
		),
		mcp.WithBoolean("enable_host_scan",
			mcp.Description("Enable host scanning"),
		),
	)
	s.AddTool(scanClusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubescape", "scan"}
		
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "framework", framework)
		}
		if kubeContext := request.GetString("kube_context", ""); kubeContext != "" {
			args = append(args, "--kube-context", kubeContext)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		if request.GetBool("verbose", false) {
			args = append(args, "--verbose")
		}
		if request.GetBool("submit", false) {
			args = append(args, "--submit")
		}
		if request.GetBool("enable_host_scan", false) {
			args = append(args, "--enable-host-scan")
		}
		
		return executeShipCommand(args)
	})

	// Kubescape scan manifests tool
	scanManifestsTool := mcp.NewTool("kubescape_scan_manifests",
		mcp.WithDescription("Scan Kubernetes manifests using Kubescape"),
		mcp.WithString("manifests_path",
			mcp.Description("Path to Kubernetes manifests (files or directory)"),
			mcp.Required(),
		),
		mcp.WithString("framework",
			mcp.Description("Security framework to use"),
			mcp.Enum("nsa", "mitre", "cis"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "junit", "pdf", "html", "sarif"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose output"),
		),
	)
	s.AddTool(scanManifestsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		manifestsPath := request.GetString("manifests_path", "")
		args := []string{"kubescape", "scan", manifestsPath}
		
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "framework", framework)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		if request.GetBool("verbose", false) {
			args = append(args, "--verbose")
		}
		
		return executeShipCommand(args)
	})

	// Kubescape get version tool
	getVersionTool := mcp.NewTool("kubescape_get_version",
		mcp.WithDescription("Get Kubescape version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubescape", "version"}
		return executeShipCommand(args)
	})

}