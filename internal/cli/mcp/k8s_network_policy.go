package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddK8sNetworkPolicyTools adds Kubernetes network policy management MCP tool implementations
func AddK8sNetworkPolicyTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// K8s network policy create tool
	createPolicyTool := mcp.NewTool("k8s_network_policy_create",
		mcp.WithDescription("Create Kubernetes network policy"),
		mcp.WithString("policy_name",
			mcp.Description("Name of the network policy"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
			mcp.Required(),
		),
		mcp.WithString("pod_selector",
			mcp.Description("Pod selector labels (key=value format)"),
		),
	)
	s.AddTool(createPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyName := request.GetString("policy_name", "")
		namespace := request.GetString("namespace", "")
		args := []string{"kubernetes", "network-policy", "create", policyName, "--namespace", namespace}
		if podSelector := request.GetString("pod_selector", ""); podSelector != "" {
			args = append(args, "--pod-selector", podSelector)
		}
		return executeShipCommand(args)
	})

	// K8s network policy validate tool
	validatePolicyTool := mcp.NewTool("k8s_network_policy_validate",
		mcp.WithDescription("Validate Kubernetes network policy configuration"),
		mcp.WithString("policy_file",
			mcp.Description("Path to network policy YAML file"),
			mcp.Required(),
		),
	)
	s.AddTool(validatePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		args := []string{"kubernetes", "network-policy", "validate", policyFile}
		return executeShipCommand(args)
	})

	// K8s network policy test tool
	testPolicyTool := mcp.NewTool("k8s_network_policy_test",
		mcp.WithDescription("Test network connectivity with current policies"),
		mcp.WithString("source_pod",
			mcp.Description("Source pod name"),
			mcp.Required(),
		),
		mcp.WithString("target_pod",
			mcp.Description("Target pod name"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(testPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sourcePod := request.GetString("source_pod", "")
		targetPod := request.GetString("target_pod", "")
		args := []string{"kubernetes", "network-policy", "test", sourcePod, targetPod}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		return executeShipCommand(args)
	})

	// K8s network policy analyze tool
	analyzePolicyTool := mcp.NewTool("k8s_network_policy_analyze",
		mcp.WithDescription("Analyze network policies for gaps and conflicts"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to analyze"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (json, yaml, table)"),
		),
	)
	s.AddTool(analyzePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "network-policy", "analyze"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// K8s network policy generate tool
	generatePolicyTool := mcp.NewTool("k8s_network_policy_generate",
		mcp.WithDescription("Generate network policy from observed traffic"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
			mcp.Required(),
		),
		mcp.WithString("duration",
			mcp.Description("Observation duration (e.g., 5m, 1h)"),
		),
	)
	s.AddTool(generatePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		args := []string{"kubernetes", "network-policy", "generate", "--namespace", namespace}
		if duration := request.GetString("duration", ""); duration != "" {
			args = append(args, "--duration", duration)
		}
		return executeShipCommand(args)
	})

	// K8s network policy get version tool
	getVersionTool := mcp.NewTool("k8s_network_policy_get_version",
		mcp.WithDescription("Get Kubernetes network policy tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "network-policy", "--version"}
		return executeShipCommand(args)
	})
}