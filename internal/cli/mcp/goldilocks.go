package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGoldilocksTools adds Goldilocks (Kubernetes resource recommendations) MCP tool implementations using kubectl and Helm
func AddGoldilocksTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Goldilocks install via Helm
	installHelmTool := mcp.NewTool("goldilocks_install_helm",
		mcp.WithDescription("Install Goldilocks using Helm chart"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for installation"),
			mcp.Required(),
		),
		mcp.WithString("release_name",
			mcp.Description("Helm release name"),
		),
	)
	s.AddTool(installHelmTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "goldilocks")
		releaseName := request.GetString("release_name", "goldilocks")
		
		// First add repo and create namespace, then install
		args := []string{"sh", "-c", 
			"helm repo add fairwinds-stable https://charts.fairwinds.com/stable && " +
			"kubectl create namespace " + namespace + " --dry-run=client -o yaml | kubectl apply -f - && " +
			"helm install " + releaseName + " --namespace " + namespace + " fairwinds-stable/goldilocks"}
		return executeShipCommand(args)
	})

	// Goldilocks enable namespace tool
	enableNamespaceTool := mcp.NewTool("goldilocks_enable_namespace",
		mcp.WithDescription("Enable Goldilocks monitoring for a namespace"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to enable"),
			mcp.Required(),
		),
	)
	s.AddTool(enableNamespaceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		args := []string{"kubectl", "label", "ns", namespace, "goldilocks.fairwinds.com/enabled=true", "--overwrite"}
		return executeShipCommand(args)
	})

	// Goldilocks dashboard access tool
	dashboardTool := mcp.NewTool("goldilocks_dashboard",
		mcp.WithDescription("Port-forward to access Goldilocks dashboard"),
		mcp.WithString("namespace",
			mcp.Description("Namespace where Goldilocks is installed"),
		),
		mcp.WithString("local_port",
			mcp.Description("Local port for dashboard access"),
		),
	)
	s.AddTool(dashboardTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "goldilocks")
		localPort := request.GetString("local_port", "8080")
		args := []string{"kubectl", "-n", namespace, "port-forward", "svc/goldilocks-dashboard", localPort + ":80"}
		return executeShipCommand(args)
	})

	// Goldilocks get VPA recommendations tool
	getRecommendationsTool := mcp.NewTool("goldilocks_get_recommendations",
		mcp.WithDescription("Get VPA recommendations from Goldilocks-enabled namespace"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to get recommendations for"),
			mcp.Required(),
		),
	)
	s.AddTool(getRecommendationsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		args := []string{"kubectl", "get", "vpa", "-n", namespace, "-o", "yaml"}
		return executeShipCommand(args)
	})

	// Goldilocks uninstall tool
	uninstallTool := mcp.NewTool("goldilocks_uninstall",
		mcp.WithDescription("Uninstall Goldilocks using Helm"),
		mcp.WithString("namespace",
			mcp.Description("Namespace where Goldilocks is installed"),
		),
		mcp.WithString("release_name",
			mcp.Description("Helm release name"),
		),
	)
	s.AddTool(uninstallTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "goldilocks")
		releaseName := request.GetString("release_name", "goldilocks")
		args := []string{"helm", "uninstall", releaseName, "--namespace", namespace}
		return executeShipCommand(args)
	})
}