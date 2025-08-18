package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddKyvernoMultitenantTools adds Kyverno multi-tenant policy MCP tool implementations using direct Dagger calls
func AddKyvernoMultitenantTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addKyvernoMultitenantToolsDirect(s)
}

// addKyvernoMultitenantToolsDirect adds Kyverno multitenant tools using direct Dagger module calls
func addKyvernoMultitenantToolsDirect(s *server.MCPServer) {
	// Create tenant namespace with labels
	createTenantNamespaceTool := mcp.NewTool("kyverno_multitenant_create_namespace",
		mcp.WithDescription("Create namespace for tenant with appropriate labels"),
		mcp.WithString("tenant_name",
			mcp.Description("Name of the tenant"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace name for the tenant"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(createTenantNamespaceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoMultitenantModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")
		if namespace == "" {
			return mcp.NewToolResultError("namespace is required"), nil
		}
		kubeconfig := request.GetString("kubeconfig", "")

		// Create tenant namespace
		output, err := module.CreateTenantNamespace(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("create tenant namespace failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Apply namespace isolation policy
	applyNamespaceIsolationTool := mcp.NewTool("kyverno_multitenant_namespace_isolation",
		mcp.WithDescription("Apply Kyverno policy for namespace isolation"),
		mcp.WithString("tenant_name",
			mcp.Description("Name of the tenant for policy generation"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(applyNamespaceIsolationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoMultitenantModule(client)

		// Get parameters
		tenantName := request.GetString("tenant_name", "")
		if tenantName == "" {
			return mcp.NewToolResultError("tenant_name is required"), nil
		}
		kubeconfig := request.GetString("kubeconfig", "")

		// Create tenant policies (includes isolation)
		output, err := module.CreateTenantPolicies(ctx, tenantName, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("create tenant policies failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Create resource quota for tenant
	createResourceQuotaTool := mcp.NewTool("kyverno_multitenant_create_quota",
		mcp.WithDescription("Create resource quota for tenant namespace"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to create quota for"),
			mcp.Required(),
		),
		mcp.WithString("cpu_limit",
			mcp.Description("CPU limit (e.g., '2' or '1000m')"),
		),
		mcp.WithString("memory_limit",
			mcp.Description("Memory limit (e.g., '4Gi' or '2048Mi')"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(createResourceQuotaTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoMultitenantModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")
		if namespace == "" {
			return mcp.NewToolResultError("namespace is required"), nil
		}
		cpuLimit := request.GetString("cpu_limit", "1")
		memoryLimit := request.GetString("memory_limit", "2Gi")
		kubeconfig := request.GetString("kubeconfig", "")

		// Create resource quota
		output, err := module.CreateResourceQuota(ctx, namespace, cpuLimit, memoryLimit, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("create resource quota failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Generate and apply multitenant policies
	applyGeneratePolicyTool := mcp.NewTool("kyverno_multitenant_generate_policy",
		mcp.WithDescription("Generate and apply comprehensive multitenant policies"),
		mcp.WithString("tenants_config",
			mcp.Description("Path to tenants configuration file"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(applyGeneratePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoMultitenantModule(client)

		// Get parameters
		tenantsConfig := request.GetString("tenants_config", "")
		if tenantsConfig == "" {
			return mcp.NewToolResultError("tenants_config is required"), nil
		}
		kubeconfig := request.GetString("kubeconfig", "")

		// Validate multitenant setup
		output, err := module.ValidateMultitenantSetup(ctx, tenantsConfig, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("validate multitenant setup failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// List tenant namespaces
	listTenantNamespacesTool := mcp.NewTool("kyverno_multitenant_list_namespaces",
		mcp.WithDescription("List all tenant namespaces"),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(listTenantNamespacesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoMultitenantModule(client)

		// Get parameters
		kubeconfig := request.GetString("kubeconfig", "")

		// List tenant namespaces
		output, err := module.ListTenantNamespaces(ctx, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("list tenant namespaces failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Get tenant policies
	getTenantPoliciesTool := mcp.NewTool("kyverno_multitenant_get_policies",
		mcp.WithDescription("Get Kyverno policies for a specific tenant namespace"),
		mcp.WithString("namespace",
			mcp.Description("Tenant namespace to get policies for"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(getTenantPoliciesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoMultitenantModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")
		if namespace == "" {
			return mcp.NewToolResultError("namespace is required"), nil
		}
		kubeconfig := request.GetString("kubeconfig", "")

		// Get tenant policies
		output, err := module.GetTenantPolicies(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get tenant policies failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}