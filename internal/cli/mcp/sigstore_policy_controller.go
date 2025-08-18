package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSigstorePolicyControllerTools adds Sigstore Policy Controller MCP tool implementations using real CLI tools
func AddSigstorePolicyControllerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Test policy against image using policy-controller-tester
	testPolicyTool := mcp.NewTool("sigstore_test_policy",
		mcp.WithDescription("Test Sigstore policy against container image using real policy-controller-tester"),
		mcp.WithString("policy",
			mcp.Description("Path to ClusterImagePolicy file or URL"),
			mcp.Required(),
		),
		mcp.WithString("image",
			mcp.Description("Container image to test against policy"),
			mcp.Required(),
		),
		mcp.WithString("resource",
			mcp.Description("Path to Kubernetes resource file"),
		),
		mcp.WithString("trustroot",
			mcp.Description("Path to Kubernetes TrustRoot resource"),
		),
		mcp.WithString("log_level",
			mcp.Description("Log level for output"),
			mcp.Enum("debug", "info", "warn", "error"),
		),
	)
	s.AddTool(testPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policy := request.GetString("policy", "")
		image := request.GetString("image", "")
		args := []string{"policy-tester", "--policy", policy, "--image", image}
		if resource := request.GetString("resource", ""); resource != "" {
			args = append(args, "-resource", resource)
		}
		if trustroot := request.GetString("trustroot", ""); trustroot != "" {
			args = append(args, "-trustroot", trustroot)
		}
		if logLevel := request.GetString("log_level", ""); logLevel != "" {
			args = append(args, "-log-level", logLevel)
		}
		return executeShipCommand(args)
	})

	// Get tester version
	versionTool := mcp.NewTool("sigstore_tester_version",
		mcp.WithDescription("Get policy-controller-tester version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"policy-tester", "--version"}
		return executeShipCommand(args)
	})

	// Create ClusterImagePolicy using kubectl
	createPolicyTool := mcp.NewTool("sigstore_create_policy",
		mcp.WithDescription("Create ClusterImagePolicy using kubectl"),
		mcp.WithString("policy_file",
			mcp.Description("Path to ClusterImagePolicy YAML file"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace for the policy (optional)"),
		),
	)
	s.AddTool(createPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		args := []string{"kubectl", "apply", "-f", policyFile}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		}
		return executeShipCommand(args)
	})

	// List ClusterImagePolicies using kubectl
	listPolicesTool := mcp.NewTool("sigstore_list_policies",
		mcp.WithDescription("List ClusterImagePolicies using kubectl"),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("yaml", "json", "wide", "name"),
		),
	)
	s.AddTool(listPolicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubectl", "get", "clusterimagepolicy"}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "-o", output)
		}
		return executeShipCommand(args)
	})

	// Delete ClusterImagePolicy using kubectl
	deletePolicyTool := mcp.NewTool("sigstore_delete_policy",
		mcp.WithDescription("Delete ClusterImagePolicy using kubectl"),
		mcp.WithString("policy_name",
			mcp.Description("Name of the ClusterImagePolicy to delete"),
			mcp.Required(),
		),
	)
	s.AddTool(deletePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyName := request.GetString("policy_name", "")
		args := []string{"kubectl", "delete", "clusterimagepolicy", policyName}
		return executeShipCommand(args)
	})

	// Enable policy enforcement for namespace using kubectl
	enableNamespaceTool := mcp.NewTool("sigstore_enable_namespace",
		mcp.WithDescription("Enable Sigstore policy enforcement for namespace using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to enable policy enforcement"),
			mcp.Required(),
		),
	)
	s.AddTool(enableNamespaceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		args := []string{"kubectl", "label", "namespace", namespace, "policy.sigstore.dev/include=true"}
		return executeShipCommand(args)
	})

	// Disable policy enforcement for namespace using kubectl
	disableNamespaceTool := mcp.NewTool("sigstore_disable_namespace",
		mcp.WithDescription("Disable Sigstore policy enforcement for namespace using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to disable policy enforcement"),
			mcp.Required(),
		),
	)
	s.AddTool(disableNamespaceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		args := []string{"kubectl", "label", "namespace", namespace, "policy.sigstore.dev/exclude=true"}
		return executeShipCommand(args)
	})

	// Get policy enforcement status for namespace using kubectl
	getNamespaceStatusTool := mcp.NewTool("sigstore_get_namespace_status",
		mcp.WithDescription("Get Sigstore policy enforcement status for namespace using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to check policy enforcement status"),
			mcp.Required(),
		),
	)
	s.AddTool(getNamespaceStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		args := []string{"kubectl", "get", "namespace", namespace, "--show-labels"}
		return executeShipCommand(args)
	})

	// Describe ClusterImagePolicy using kubectl
	describePolicyTool := mcp.NewTool("sigstore_describe_policy",
		mcp.WithDescription("Describe ClusterImagePolicy using kubectl"),
		mcp.WithString("policy_name",
			mcp.Description("Name of the ClusterImagePolicy to describe"),
			mcp.Required(),
		),
	)
	s.AddTool(describePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyName := request.GetString("policy_name", "")
		args := []string{"kubectl", "describe", "clusterimagepolicy", policyName}
		return executeShipCommand(args)
	})
}