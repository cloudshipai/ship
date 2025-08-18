package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddKubescapeTools adds Kubescape (Kubernetes security scanner) MCP tool implementations using direct Dagger calls
func AddKubescapeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addKubescapeToolsDirect(s)
}

// addKubescapeToolsDirect adds Kubescape tools using direct Dagger module calls
func addKubescapeToolsDirect(s *server.MCPServer) {
	// Kubescape scan cluster
	scanClusterTool := mcp.NewTool("kubescape_scan_cluster",
		mcp.WithDescription("Scan Kubernetes cluster using Kubescape"),
		mcp.WithString("framework",
			mcp.Description("Security framework to use"),
			mcp.Enum("nsa", "mitre", "armobest", "devopsbest", "cis-v1.23-t1.0.1", "cis-eks-t1.2.0", "cis-aks-t1.2.0", "cis-gke-t1.2.0", "cis-rke2-t1.2.0", "pci-dss-v3.2.1", "soc2", "iso27001"),
		),
		mcp.WithString("severity_threshold",
			mcp.Description("Severity threshold (low, medium, high, critical)"),
			mcp.Enum("low", "medium", "high", "critical"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("pretty-printer", "json", "junit", "prometheus", "pdf", "html", "sarif"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace to scan (empty for all)"),
		),
		mcp.WithBoolean("exclude_kube_system",
			mcp.Description("Exclude kube-system namespace"),
		),
	)
	s.AddTool(scanClusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubescapeModule(client)

		// Get parameters
		framework := request.GetString("framework", "nsa")
		severityThreshold := request.GetString("severity_threshold", "")
		format := request.GetString("format", "pretty-printer")
		kubeconfig := request.GetString("kubeconfig", "")

		// Scan cluster
		output, err := module.ScanCluster(ctx, kubeconfig, framework, format, severityThreshold)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kubescape cluster scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kubescape scan manifests
	scanManifestsTool := mcp.NewTool("kubescape_scan_manifests",
		mcp.WithDescription("Scan Kubernetes manifest files using Kubescape"),
		mcp.WithString("path",
			mcp.Description("Path to directory containing manifest files"),
			mcp.Required(),
		),
		mcp.WithString("framework",
			mcp.Description("Security framework to use"),
			mcp.Enum("nsa", "mitre", "armobest", "devopsbest", "cis-v1.23-t1.0.1", "cis-eks-t1.2.0", "cis-aks-t1.2.0", "cis-gke-t1.2.0", "cis-rke2-t1.2.0", "pci-dss-v3.2.1", "soc2", "iso27001"),
		),
		mcp.WithString("severity_threshold",
			mcp.Description("Severity threshold (low, medium, high, critical)"),
			mcp.Enum("low", "medium", "high", "critical"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("pretty-printer", "json", "junit", "prometheus", "pdf", "html", "sarif"),
		),
		mcp.WithBoolean("include_helm",
			mcp.Description("Treat as Helm chart"),
		),
	)
	s.AddTool(scanManifestsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubescapeModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		framework := request.GetString("framework", "nsa")
		severityThreshold := request.GetString("severity_threshold", "")
		format := request.GetString("format", "pretty-printer")
		includeHelm := request.GetBool("include_helm", false)

		// Choose appropriate scan method
		var output string
		if includeHelm {
			// Scan as Helm chart
			output, err = module.ScanHelm(ctx, path, framework, format)
		} else {
			// Scan as regular manifests
			output, err = module.ScanManifests(ctx, path, framework, format, severityThreshold)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kubescape manifest scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kubescape get version
	getVersionTool := mcp.NewTool("kubescape_get_version",
		mcp.WithDescription("Get Kubescape version"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubescapeModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get kubescape version: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}