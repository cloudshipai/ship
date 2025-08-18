package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKyvernoMultitenantTools adds Kyverno multi-tenant policy MCP tool implementations using real CLI commands
func AddKyvernoMultitenantTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Create tenant namespace with labels
	createTenantNamespaceTool := mcp.NewTool("kyverno_multitenant_create_namespace",
		mcp.WithDescription("Create namespace for tenant with appropriate labels using kubectl"),
		mcp.WithString("tenant_name",
			mcp.Description("Name of the tenant"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace name for the tenant"),
			mcp.Required(),
		),
	)
	s.AddTool(createTenantNamespaceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tenantName := request.GetString("tenant_name", "")
		namespace := request.GetString("namespace", "")
		
		// Create namespace with tenant label
		args := []string{"kubectl", "create", "namespace", namespace}
		result, err := executeShipCommand(args)
		if err != nil {
			return result, err
		}
		
		// Label the namespace with tenant information
		labelArgs := []string{"kubectl", "label", "namespace", namespace, 
			fmt.Sprintf("tenant=%s", tenantName),
			fmt.Sprintf("kyverno.io/tenant=%s", tenantName)}
		return executeShipCommand(labelArgs)
	})

	// Apply namespace isolation policy
	applyNamespaceIsolationTool := mcp.NewTool("kyverno_multitenant_namespace_isolation",
		mcp.WithDescription("Apply Kyverno policy for namespace isolation using kubectl"),
		mcp.WithString("policy_file",
			mcp.Description("Path to namespace isolation policy YAML file"),
			mcp.Required(),
		),
	)
	s.AddTool(applyNamespaceIsolationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		args := []string{"kubectl", "apply", "-f", policyFile}
		return executeShipCommand(args)
	})

	// Create ResourceQuota for tenant
	createResourceQuotaTool := mcp.NewTool("kyverno_multitenant_create_quota",
		mcp.WithDescription("Create ResourceQuota for tenant namespace using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Tenant namespace"),
			mcp.Required(),
		),
		mcp.WithString("cpu_limit",
			mcp.Description("CPU limit (e.g., '4')"),
		),
		mcp.WithString("memory_limit",
			mcp.Description("Memory limit (e.g., '8Gi')"),
		),
		mcp.WithString("pods_limit",
			mcp.Description("Maximum number of pods"),
		),
	)
	s.AddTool(createResourceQuotaTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		
		// Build ResourceQuota YAML
		quotaYAML := fmt.Sprintf(`
apiVersion: v1
kind: ResourceQuota
metadata:
  name: tenant-quota
  namespace: %s
spec:
  hard:`, namespace)
		
		if cpu := request.GetString("cpu_limit", ""); cpu != "" {
			quotaYAML += fmt.Sprintf("\n    requests.cpu: '%s'", cpu)
			quotaYAML += fmt.Sprintf("\n    limits.cpu: '%s'", cpu)
		}
		if memory := request.GetString("memory_limit", ""); memory != "" {
			quotaYAML += fmt.Sprintf("\n    requests.memory: %s", memory)
			quotaYAML += fmt.Sprintf("\n    limits.memory: %s", memory)
		}
		if pods := request.GetString("pods_limit", ""); pods != "" {
			quotaYAML += fmt.Sprintf("\n    pods: '%s'", pods)
		}
		
		// Apply using kubectl with inline YAML
		args := []string{"sh", "-c", fmt.Sprintf("echo '%s' | kubectl apply -f -", quotaYAML)}
		return executeShipCommand(args)
	})

	// Apply Kyverno generate policy for automatic resource creation
	applyGeneratePolicyTool := mcp.NewTool("kyverno_multitenant_generate_policy",
		mcp.WithDescription("Apply Kyverno generate policy for automatic resource creation in tenant namespaces"),
		mcp.WithString("policy_file",
			mcp.Description("Path to Kyverno generate policy file"),
			mcp.Required(),
		),
	)
	s.AddTool(applyGeneratePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		
		// First validate the policy
		validateArgs := []string{"kyverno", "apply", policyFile, "--dry-run"}
		_, err := executeShipCommand(validateArgs)
		if err != nil {
			return nil, fmt.Errorf("policy validation failed: %v", err)
		}
		
		// Apply the policy to cluster
		args := []string{"kubectl", "apply", "-f", policyFile}
		return executeShipCommand(args)
	})

	// List tenant namespaces
	listTenantNamespacesTool := mcp.NewTool("kyverno_multitenant_list_namespaces",
		mcp.WithDescription("List namespaces for a specific tenant using kubectl"),
		mcp.WithString("tenant_name",
			mcp.Description("Name of the tenant"),
			mcp.Required(),
		),
	)
	s.AddTool(listTenantNamespacesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tenantName := request.GetString("tenant_name", "")
		args := []string{"kubectl", "get", "namespaces", "-l", fmt.Sprintf("tenant=%s", tenantName)}
		return executeShipCommand(args)
	})

	// Get tenant policies
	getTenantPoliciesTool := mcp.NewTool("kyverno_multitenant_get_policies",
		mcp.WithDescription("Get Kyverno policies affecting a tenant's namespaces"),
		mcp.WithString("namespace",
			mcp.Description("Tenant namespace"),
			mcp.Required(),
		),
	)
	s.AddTool(getTenantPoliciesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		
		// Get policies in the namespace
		args := []string{"kubectl", "get", "policy", "-n", namespace}
		return executeShipCommand(args)
	})
}