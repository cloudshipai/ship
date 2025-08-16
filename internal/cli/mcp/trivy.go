package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTrivyTools adds Trivy MCP tool implementations
func AddTrivyTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Trivy scan image tool
	scanImageTool := mcp.NewTool("trivy_scan_image",
		mcp.WithDescription("Scan container image for vulnerabilities using Trivy"),
		mcp.WithString("image_name",
			mcp.Description("Container image name to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"security", "trivy", "image", imageName}
		return executeShipCommand(args)
	})

	// Trivy scan filesystem tool
	scanFilesystemTool := mcp.NewTool("trivy_scan_filesystem",
		mcp.WithDescription("Scan filesystem for vulnerabilities using Trivy"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
	)
	s.AddTool(scanFilesystemTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "trivy", "fs"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// Trivy scan repository tool
	scanRepositoryTool := mcp.NewTool("trivy_scan_repository",
		mcp.WithDescription("Scan git repository for vulnerabilities using Trivy"),
		mcp.WithString("repo_url",
			mcp.Description("Git repository URL to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoURL := request.GetString("repo_url", "")
		args := []string{"security", "trivy", "repo", repoURL}
		return executeShipCommand(args)
	})

	// Trivy scan config tool
	scanConfigTool := mcp.NewTool("trivy_scan_config",
		mcp.WithDescription("Scan configuration files for security issues using Trivy"),
		mcp.WithString("directory",
			mcp.Description("Directory containing configuration files (default: current directory)"),
		),
	)
	s.AddTool(scanConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "trivy", "config"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// Trivy scan SBOM tool
	scanSBOMTool := mcp.NewTool("trivy_scan_sbom",
		mcp.WithDescription("Scan SBOM file for vulnerabilities using Trivy"),
		mcp.WithString("sbom_path",
			mcp.Description("Path to SBOM file"),
			mcp.Required(),
		),
	)
	s.AddTool(scanSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sbomPath := request.GetString("sbom_path", "")
		args := []string{"security", "trivy", "sbom", sbomPath}
		return executeShipCommand(args)
	})

	// Trivy scan Kubernetes tool
	scanKubernetesTool := mcp.NewTool("trivy_scan_kubernetes",
		mcp.WithDescription("Scan Kubernetes cluster for vulnerabilities using Trivy"),
		mcp.WithString("cluster_name",
			mcp.Description("Kubernetes cluster name or context (default: current context)"),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to scan (default: all namespaces)"),
		),
	)
	s.AddTool(scanKubernetesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "trivy", "k8s"}
		if clusterName := request.GetString("cluster_name", ""); clusterName != "" {
			args = append(args, "--context", clusterName)
		}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		} else {
			args = append(args, "cluster")
		}
		return executeShipCommand(args)
	})

	// Trivy scan with scanners tool
	scanWithScannersTool := mcp.NewTool("trivy_scan_with_scanners",
		mcp.WithDescription("Scan with specific scanner types using Trivy"),
		mcp.WithString("target",
			mcp.Description("Target to scan (image name, directory, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("scanners",
			mcp.Description("Comma-separated list of scanners (vuln,secret,misconfig,license)"),
			mcp.Required(),
		),
		mcp.WithString("scan_type",
			mcp.Description("Type of scan to perform"),
			mcp.Enum("image", "fs", "repo"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithScannersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scanType := request.GetString("scan_type", "")
		target := request.GetString("target", "")
		scanners := request.GetString("scanners", "")
		args := []string{"security", "trivy", scanType, "--scanners", scanners, target}
		return executeShipCommand(args)
	})

	// Trivy get version tool
	getVersionTool := mcp.NewTool("trivy_get_version",
		mcp.WithDescription("Get Trivy version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "trivy", "--version"}
		return executeShipCommand(args)
	})
}