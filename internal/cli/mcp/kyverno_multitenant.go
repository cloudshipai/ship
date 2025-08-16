package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKyvernoMultitenantTools adds Kyverno multi-tenant policy MCP tool implementations
func AddKyvernoMultitenantTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Kyverno multitenant setup tool
	setupTool := mcp.NewTool("kyverno_multitenant_setup",
		mcp.WithDescription("Setup Kyverno for multi-tenant environment"),
		mcp.WithString("tenant_name",
			mcp.Description("Name of the tenant"),
			mcp.Required(),
		),
		mcp.WithString("namespace_prefix",
			mcp.Description("Namespace prefix for tenant"),
			mcp.Required(),
		),
	)
	s.AddTool(setupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tenantName := request.GetString("tenant_name", "")
		namespacePrefix := request.GetString("namespace_prefix", "")
		args := []string{"kubernetes", "kyverno-multitenant", "setup", tenantName, "--namespace-prefix", namespacePrefix}
		return executeShipCommand(args)
	})

	// Kyverno multitenant create policies tool
	createPolicesTool := mcp.NewTool("kyverno_multitenant_create_policies",
		mcp.WithDescription("Create tenant-specific Kyverno policies"),
		mcp.WithString("tenant_name",
			mcp.Description("Name of the tenant"),
			mcp.Required(),
		),
		mcp.WithString("policy_set",
			mcp.Description("Policy set to apply (basic, strict, custom)"),
			mcp.Required(),
		),
	)
	s.AddTool(createPolicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tenantName := request.GetString("tenant_name", "")
		policySet := request.GetString("policy_set", "")
		args := []string{"kubernetes", "kyverno-multitenant", "create-policies", tenantName, "--policy-set", policySet}
		return executeShipCommand(args)
	})

	// Kyverno multitenant isolate resources tool
	isolateResourcesTool := mcp.NewTool("kyverno_multitenant_isolate_resources",
		mcp.WithDescription("Create resource isolation policies for tenant"),
		mcp.WithString("tenant_name",
			mcp.Description("Name of the tenant"),
			mcp.Required(),
		),
		mcp.WithString("resource_types",
			mcp.Description("Comma-separated list of resource types to isolate"),
		),
	)
	s.AddTool(isolateResourcesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tenantName := request.GetString("tenant_name", "")
		args := []string{"kubernetes", "kyverno-multitenant", "isolate", tenantName}
		if resourceTypes := request.GetString("resource_types", ""); resourceTypes != "" {
			args = append(args, "--resources", resourceTypes)
		}
		return executeShipCommand(args)
	})

	// Kyverno multitenant enforce quotas tool
	enforceQuotasTool := mcp.NewTool("kyverno_multitenant_enforce_quotas",
		mcp.WithDescription("Enforce resource quotas for tenant"),
		mcp.WithString("tenant_name",
			mcp.Description("Name of the tenant"),
			mcp.Required(),
		),
		mcp.WithString("cpu_limit",
			mcp.Description("CPU limit for tenant (e.g., 4, 8)"),
		),
		mcp.WithString("memory_limit",
			mcp.Description("Memory limit for tenant (e.g., 8Gi, 16Gi)"),
		),
	)
	s.AddTool(enforceQuotasTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tenantName := request.GetString("tenant_name", "")
		args := []string{"kubernetes", "kyverno-multitenant", "enforce-quotas", tenantName}
		if cpuLimit := request.GetString("cpu_limit", ""); cpuLimit != "" {
			args = append(args, "--cpu", cpuLimit)
		}
		if memoryLimit := request.GetString("memory_limit", ""); memoryLimit != "" {
			args = append(args, "--memory", memoryLimit)
		}
		return executeShipCommand(args)
	})

	// Kyverno multitenant validate tenant tool
	validateTenantTool := mcp.NewTool("kyverno_multitenant_validate_tenant",
		mcp.WithDescription("Validate tenant configuration and policies"),
		mcp.WithString("tenant_name",
			mcp.Description("Name of the tenant to validate"),
			mcp.Required(),
		),
	)
	s.AddTool(validateTenantTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tenantName := request.GetString("tenant_name", "")
		args := []string{"kubernetes", "kyverno-multitenant", "validate", tenantName}
		return executeShipCommand(args)
	})

	// Kyverno multitenant get version tool
	getVersionTool := mcp.NewTool("kyverno_multitenant_get_version",
		mcp.WithDescription("Get Kyverno multitenant tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "kyverno-multitenant", "--version"}
		return executeShipCommand(args)
	})
}