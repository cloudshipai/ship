package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddK8sNetworkPolicyTools adds Kubernetes network policy management MCP tool implementations using direct Dagger calls
func AddK8sNetworkPolicyTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addK8sNetworkPolicyToolsDirect(s)
}

// addK8sNetworkPolicyToolsDirect adds K8s network policy tools using direct Dagger module calls
func addK8sNetworkPolicyToolsDirect(s *server.MCPServer) {
	// kubectl network policy management
	kubectlNetworkPolicyTool := mcp.NewTool("k8s_network_policy_kubectl",
		mcp.WithDescription("Manage network policies using kubectl"),
		mcp.WithString("action",
			mcp.Description("Action to perform"),
			mcp.Required(),
			mcp.Enum("get", "describe", "delete", "create", "apply"),
		),
		mcp.WithString("resource",
			mcp.Description("Resource name or file path (for create/apply/describe/delete)"),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "wide", "name"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(kubectlNetworkPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewK8sNetworkPolicyModule(client)

		// Get parameters
		action := request.GetString("action", "")
		if action == "" {
			return mcp.NewToolResultError("action is required"), nil
		}

		resource := request.GetString("resource", "")
		namespace := request.GetString("namespace", "")
		outputFormat := request.GetString("output", "")
		kubeconfig := request.GetString("kubeconfig", "")

		// Execute kubectl command
		output, err := module.KubectlNetworkPolicy(ctx, action, resource, namespace, outputFormat, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kubectl command failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Netfetch network policy scanner
	netfetchScanTool := mcp.NewTool("k8s_network_policy_netfetch_scan",
		mcp.WithDescription("Scan Kubernetes network policies using netfetch"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to scan (omit for entire cluster)"),
		),
		mcp.WithBoolean("dryrun",
			mcp.Description("Run scan without applying changes"),
		),
		mcp.WithBoolean("cilium",
			mcp.Description("Scan Cilium network policies"),
		),
		mcp.WithString("target",
			mcp.Description("Scan a specific policy by name"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(netfetchScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewK8sNetworkPolicyModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")
		dryrun := request.GetBool("dryrun", false)
		cilium := request.GetBool("cilium", false)
		target := request.GetString("target", "")
		kubeconfig := request.GetString("kubeconfig", "")

		// Execute netfetch scan
		output, err := module.NetfetchScan(ctx, namespace, dryrun, cilium, target, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("netfetch scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Netfetch dashboard
	netfetchDashTool := mcp.NewTool("k8s_network_policy_netfetch_dash",
		mcp.WithDescription("Launch netfetch dashboard for network policy visualization"),
		mcp.WithString("port",
			mcp.Description("Port number for dashboard (default: 8080)"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(netfetchDashTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewK8sNetworkPolicyModule(client)

		// Get parameters
		port := request.GetString("port", "")
		kubeconfig := request.GetString("kubeconfig", "")

		// Launch dashboard
		output, err := module.NetfetchDashboard(ctx, port, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("netfetch dashboard failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Netpol-analyzer evaluation
	netpolAnalyzerEvalTool := mcp.NewTool("k8s_network_policy_netpol_eval",
		mcp.WithDescription("Evaluate network connectivity using netpol-analyzer"),
		mcp.WithString("dirpath",
			mcp.Description("Directory path containing Kubernetes resources"),
			mcp.Required(),
		),
		mcp.WithString("source",
			mcp.Description("Source pod name"),
			mcp.Required(),
		),
		mcp.WithString("destination",
			mcp.Description("Destination pod name"),
			mcp.Required(),
		),
		mcp.WithString("port",
			mcp.Description("Port number"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose output"),
		),
	)
	s.AddTool(netpolAnalyzerEvalTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewK8sNetworkPolicyModule(client)

		// Get parameters
		dirpath := request.GetString("dirpath", "")
		if dirpath == "" {
			return mcp.NewToolResultError("dirpath is required"), nil
		}

		source := request.GetString("source", "")
		if source == "" {
			return mcp.NewToolResultError("source is required"), nil
		}

		destination := request.GetString("destination", "")
		if destination == "" {
			return mcp.NewToolResultError("destination is required"), nil
		}

		port := request.GetString("port", "")
		verbose := request.GetBool("verbose", false)

		// Evaluate connectivity
		output, err := module.NetpolEval(ctx, dirpath, source, destination, port, verbose)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("netpol eval failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Netpol-analyzer list connections
	netpolAnalyzerListTool := mcp.NewTool("k8s_network_policy_netpol_list",
		mcp.WithDescription("List all allowed connections using netpol-analyzer"),
		mcp.WithString("dirpath",
			mcp.Description("Directory path containing Kubernetes resources"),
			mcp.Required(),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose output"),
		),
		mcp.WithBoolean("quiet",
			mcp.Description("Enable quiet mode"),
		),
	)
	s.AddTool(netpolAnalyzerListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewK8sNetworkPolicyModule(client)

		// Get parameters
		dirpath := request.GetString("dirpath", "")
		if dirpath == "" {
			return mcp.NewToolResultError("dirpath is required"), nil
		}

		verbose := request.GetBool("verbose", false)
		quiet := request.GetBool("quiet", false)

		// List connections
		output, err := module.NetpolList(ctx, dirpath, verbose, quiet)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("netpol list failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Netpol-analyzer diff
	netpolAnalyzerDiffTool := mcp.NewTool("k8s_network_policy_netpol_diff",
		mcp.WithDescription("Compare network policies between two directories using netpol-analyzer"),
		mcp.WithString("dir1",
			mcp.Description("First directory containing Kubernetes resources"),
			mcp.Required(),
		),
		mcp.WithString("dir2",
			mcp.Description("Second directory containing Kubernetes resources"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format: md, csv, text"),
			mcp.Enum("md", "csv", "text"),
		),
	)
	s.AddTool(netpolAnalyzerDiffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewK8sNetworkPolicyModule(client)

		// Get parameters
		dir1 := request.GetString("dir1", "")
		if dir1 == "" {
			return mcp.NewToolResultError("dir1 is required"), nil
		}

		dir2 := request.GetString("dir2", "")
		if dir2 == "" {
			return mcp.NewToolResultError("dir2 is required"), nil
		}

		outputFormat := request.GetString("output_format", "")

		// Compare directories
		output, err := module.NetpolDiff(ctx, dir1, dir2, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("netpol diff failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}