package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddGoldilocksTools adds Goldilocks (Kubernetes resource recommendations) MCP tool implementations using direct Dagger calls
func AddGoldilocksTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGoldilocksToolsDirect(s)
}

// addGoldilocksToolsDirect adds Goldilocks tools using direct Dagger module calls
func addGoldilocksToolsDirect(s *server.MCPServer) {
	// Goldilocks install via Helm
	installHelmTool := mcp.NewTool("goldilocks_install_helm",
		mcp.WithDescription("Install Goldilocks using Helm chart"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for installation"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("release_name",
			mcp.Description("Helm release name"),
		),
	)
	s.AddTool(installHelmTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGoldilocksModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "goldilocks")
		kubeconfig := request.GetString("kubeconfig", "")

		// Install via Helm
		output, err := module.InstallHelm(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to install Goldilocks: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Goldilocks enable namespace tool
	enableNamespaceTool := mcp.NewTool("goldilocks_enable_namespace",
		mcp.WithDescription("Enable Goldilocks monitoring for a namespace"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to enable"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(enableNamespaceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGoldilocksModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")
		kubeconfig := request.GetString("kubeconfig", "")

		// Enable namespace
		output, err := module.EnableNamespace(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to enable namespace: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Goldilocks dashboard access tool (note: port-forwarding is handled differently in containerized environment)
	dashboardTool := mcp.NewTool("goldilocks_dashboard",
		mcp.WithDescription("Get dashboard access information for Goldilocks"),
		mcp.WithString("namespace",
			mcp.Description("Namespace where Goldilocks is installed"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("local_port",
			mcp.Description("Local port for dashboard access"),
		),
	)
	s.AddTool(dashboardTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGoldilocksModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "goldilocks")
		_ = request.GetString("kubeconfig", "") // Not used for version check

		// Get version as a simple check
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get version: %v", err)), nil
		}

		// Return dashboard access instructions
		localPort := request.GetString("local_port", "8080")
		info := fmt.Sprintf("Goldilocks version: %s\n\n", output)
		info += fmt.Sprintf("To access the dashboard, run:\n")
		info += fmt.Sprintf("kubectl -n %s port-forward svc/goldilocks-dashboard %s:80\n", namespace, localPort)
		info += fmt.Sprintf("Then access: http://localhost:%s", localPort)

		return mcp.NewToolResultText(info), nil
	})

	// Goldilocks get VPA recommendations tool
	getRecommendationsTool := mcp.NewTool("goldilocks_get_recommendations",
		mcp.WithDescription("Get VPA recommendations from Goldilocks-enabled namespace"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to get recommendations for"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(getRecommendationsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGoldilocksModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")
		kubeconfig := request.GetString("kubeconfig", "")

		// Get recommendations
		output, err := module.GetRecommendations(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get recommendations: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Goldilocks uninstall tool
	uninstallTool := mcp.NewTool("goldilocks_uninstall",
		mcp.WithDescription("Uninstall Goldilocks using Helm"),
		mcp.WithString("namespace",
			mcp.Description("Namespace where Goldilocks is installed"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("release_name",
			mcp.Description("Helm release name"),
		),
	)
	s.AddTool(uninstallTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGoldilocksModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "goldilocks")
		kubeconfig := request.GetString("kubeconfig", "")

		// Uninstall
		output, err := module.Uninstall(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to uninstall Goldilocks: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}