package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKyvernoTools adds Kyverno policy management MCP tool implementations
func AddKyvernoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Kyverno install tool
	installTool := mcp.NewTool("kyverno_install",
		mcp.WithDescription("Install Kyverno policy engine in Kubernetes cluster"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for Kyverno installation"),
		),
		mcp.WithString("version",
			mcp.Description("Kyverno version to install"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "kyverno", "install"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if version := request.GetString("version", ""); version != "" {
			args = append(args, "--version", version)
		}
		return executeShipCommand(args)
	})

	// Kyverno create policy tool
	createPolicyTool := mcp.NewTool("kyverno_create_policy",
		mcp.WithDescription("Create Kyverno policy from template or definition"),
		mcp.WithString("policy_name",
			mcp.Description("Name of the Kyverno policy"),
			mcp.Required(),
		),
		mcp.WithString("policy_type",
			mcp.Description("Type of policy (validate, mutate, generate)"),
			mcp.Required(),
		),
		mcp.WithString("resource_kind",
			mcp.Description("Kubernetes resource kind to target"),
			mcp.Required(),
		),
	)
	s.AddTool(createPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyName := request.GetString("policy_name", "")
		policyType := request.GetString("policy_type", "")
		resourceKind := request.GetString("resource_kind", "")
		args := []string{"kubernetes", "kyverno", "create-policy", policyName, "--type", policyType, "--kind", resourceKind}
		return executeShipCommand(args)
	})

	// Kyverno apply policy tool
	applyPolicyTool := mcp.NewTool("kyverno_apply_policy",
		mcp.WithDescription("Apply Kyverno policy to cluster"),
		mcp.WithString("policy_file",
			mcp.Description("Path to Kyverno policy YAML file"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(applyPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		args := []string{"kubernetes", "kyverno", "apply", policyFile}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		return executeShipCommand(args)
	})

	// Kyverno validate policy tool
	validatePolicyTool := mcp.NewTool("kyverno_validate_policy",
		mcp.WithDescription("Validate Kyverno policy syntax and logic"),
		mcp.WithString("policy_file",
			mcp.Description("Path to Kyverno policy YAML file"),
			mcp.Required(),
		),
	)
	s.AddTool(validatePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		args := []string{"kubernetes", "kyverno", "validate", policyFile}
		return executeShipCommand(args)
	})

	// Kyverno test policy tool
	testPolicyTool := mcp.NewTool("kyverno_test_policy",
		mcp.WithDescription("Test Kyverno policy against resource manifests"),
		mcp.WithString("policy_file",
			mcp.Description("Path to Kyverno policy YAML file"),
			mcp.Required(),
		),
		mcp.WithString("resource_file",
			mcp.Description("Path to resource YAML file to test against"),
			mcp.Required(),
		),
	)
	s.AddTool(testPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		resourceFile := request.GetString("resource_file", "")
		args := []string{"kubernetes", "kyverno", "test", policyFile, resourceFile}
		return executeShipCommand(args)
	})

	// Kyverno get version tool
	getVersionTool := mcp.NewTool("kyverno_get_version",
		mcp.WithDescription("Get Kyverno version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "kyverno", "--version"}
		return executeShipCommand(args)
	})
}