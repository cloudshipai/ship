package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddK8sNetworkPolicyTools adds Kubernetes network policy management MCP tool implementations using real CLI tools
func AddK8sNetworkPolicyTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
	)
	s.AddTool(kubectlNetworkPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		action := request.GetString("action", "")
		resource := request.GetString("resource", "")
		
		var args []string
		switch action {
		case "get":
			args = []string{"kubectl", "get", "networkpolicy"}
			if resource != "" {
				args = append(args, resource)
			}
		case "describe":
			args = []string{"kubectl", "describe", "networkpolicy"}
			if resource != "" {
				args = append(args, resource)
			}
		case "delete":
			args = []string{"kubectl", "delete", "networkpolicy", resource}
		case "create":
			args = []string{"kubectl", "create", "-f", resource}
		case "apply":
			args = []string{"kubectl", "apply", "-f", resource}
		}
		
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "-o", output)
		}
		
		return executeShipCommand(args)
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
		args := []string{"netfetch", "scan"}
		
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, namespace)
		}
		if request.GetBool("dryrun", false) {
			args = append(args, "--dryrun")
		}
		if request.GetBool("cilium", false) {
			args = append(args, "--cilium")
		}
		if target := request.GetString("target", ""); target != "" {
			args = append(args, "--target", target)
		}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		
		return executeShipCommand(args)
	})

	// Netfetch dashboard
	netfetchDashTool := mcp.NewTool("k8s_network_policy_netfetch_dash",
		mcp.WithDescription("Launch netfetch dashboard for network policy visualization"),
		mcp.WithString("port",
			mcp.Description("Port number for dashboard (default: 8080)"),
		),
	)
	s.AddTool(netfetchDashTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"netfetch", "dash"}
		
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "--port", port)
		}
		
		return executeShipCommand(args)
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
		dirpath := request.GetString("dirpath", "")
		source := request.GetString("source", "")
		destination := request.GetString("destination", "")
		
		args := []string{"netpol-analyzer", "eval", "--dirpath", dirpath, "-s", source, "-d", destination}
		
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "-p", port)
		}
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		
		return executeShipCommand(args)
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
		dirpath := request.GetString("dirpath", "")
		args := []string{"netpol-analyzer", "list", "--dirpath", dirpath}
		
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		if request.GetBool("quiet", false) {
			args = append(args, "-q")
		}
		
		return executeShipCommand(args)
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
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose output"),
		),
	)
	s.AddTool(netpolAnalyzerDiffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir1 := request.GetString("dir1", "")
		dir2 := request.GetString("dir2", "")
		args := []string{"netpol-analyzer", "diff", "--dir1", dir1, "--dir2", dir2}
		
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		
		return executeShipCommand(args)
	})
}