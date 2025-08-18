package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddKyvernoTools adds Kyverno policy management MCP tool implementations using direct Dagger calls
func AddKyvernoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addKyvernoToolsDirect(s)
}

// addKyvernoToolsDirect adds Kyverno tools using direct Dagger module calls
func addKyvernoToolsDirect(s *server.MCPServer) {
	// Kyverno install tool
	installTool := mcp.NewTool("kyverno_install",
		mcp.WithDescription("Install Kyverno in Kubernetes cluster using Helm"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for Kyverno installation (default: kyverno)"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "kyverno")
		kubeconfig := request.GetString("kubeconfig", "")

		// Install Kyverno
		output, err := module.Install(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kyverno install failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kyverno apply policy tool
	applyPolicyTool := mcp.NewTool("kyverno_apply",
		mcp.WithDescription("Apply Kyverno policies to cluster"),
		mcp.WithString("policies_path",
			mcp.Description("Path to directory containing Kyverno policy files"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(applyPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoModule(client)

		// Get parameters
		policiesPath := request.GetString("policies_path", "")
		if policiesPath == "" {
			return mcp.NewToolResultError("policies_path is required"), nil
		}
		kubeconfig := request.GetString("kubeconfig", "")

		// Apply policies
		output, err := module.ApplyPolicies(ctx, policiesPath, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kyverno apply policies failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kyverno test policy tool
	testPolicyTool := mcp.NewTool("kyverno_test",
		mcp.WithDescription("Test Kyverno policies against resources"),
		mcp.WithString("policies_path",
			mcp.Description("Path to directory containing Kyverno policy files"),
			mcp.Required(),
		),
		mcp.WithString("resources_path",
			mcp.Description("Path to directory containing Kubernetes resource files for testing"),
			mcp.Required(),
		),
	)
	s.AddTool(testPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoModule(client)

		// Get parameters
		policiesPath := request.GetString("policies_path", "")
		if policiesPath == "" {
			return mcp.NewToolResultError("policies_path is required"), nil
		}
		resourcesPath := request.GetString("resources_path", "")
		if resourcesPath == "" {
			return mcp.NewToolResultError("resources_path is required"), nil
		}

		// Test policies
		output, err := module.TestPolicies(ctx, policiesPath, resourcesPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kyverno test policies failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kyverno create cluster role tool
	createClusterRoleTool := mcp.NewTool("kyverno_create_cluster_role",
		mcp.WithDescription("Create cluster role for Kyverno"),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(createClusterRoleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoModule(client)

		// Get parameters
		kubeconfig := request.GetString("kubeconfig", "")

		// Create cluster role
		output, err := module.CreateClusterRole(ctx, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kyverno create cluster role failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kyverno version tool
	versionTool := mcp.NewTool("kyverno_version",
		mcp.WithDescription("Get Kyverno CLI version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get kyverno version: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kyverno list policies tool
	listPoliciesTool := mcp.NewTool("kyverno_list_policies",
		mcp.WithDescription("List Kyverno policies in cluster"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to list policies from (empty for all namespaces)"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(listPoliciesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")
		kubeconfig := request.GetString("kubeconfig", "")

		// List policies
		output, err := module.ListPolicies(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kyverno list policies failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kyverno apply policy file tool
	applyPolicyFileTool := mcp.NewTool("kyverno_apply_policy_file",
		mcp.WithDescription("Apply Kyverno policy file to cluster using kubectl"),
		mcp.WithString("file",
			mcp.Description("Path to Kyverno policy YAML file"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace for namespaced policies"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform dry run without applying"),
		),
	)
	s.AddTool(applyPolicyFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoModule(client)

		// Get parameters
		file := request.GetString("file", "")
		if file == "" {
			return mcp.NewToolResultError("file is required"), nil
		}
		namespace := request.GetString("namespace", "")
		kubeconfig := request.GetString("kubeconfig", "")
		dryRun := request.GetBool("dry_run", false)

		// Apply policy file
		output, err := module.ApplyPolicyFile(ctx, file, namespace, kubeconfig, dryRun)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kyverno apply policy file failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kyverno get policy reports tool
	getPolicyReportsTool := mcp.NewTool("kyverno_get_policy_reports",
		mcp.WithDescription("Get Kyverno policy reports"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to get policy reports from (empty for all namespaces)"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(getPolicyReportsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")
		kubeconfig := request.GetString("kubeconfig", "")

		// Get policy reports
		output, err := module.GetPolicyReports(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kyverno get policy reports failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kyverno status tool
	statusTool := mcp.NewTool("kyverno_status",
		mcp.WithDescription("Get Kyverno installation status"),
		mcp.WithString("namespace",
			mcp.Description("Namespace where Kyverno is installed (default: kyverno)"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(statusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKyvernoModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "kyverno")
		kubeconfig := request.GetString("kubeconfig", "")

		// Get status
		output, err := module.GetStatus(ctx, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kyverno status failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}