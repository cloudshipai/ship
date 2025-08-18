package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddSigstorePolicyControllerTools adds Sigstore Policy Controller MCP tool implementations using direct Dagger calls
func AddSigstorePolicyControllerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addSigstorePolicyControllerToolsDirect(s)
}

// addSigstorePolicyControllerToolsDirect adds Sigstore Policy Controller tools using direct Dagger module calls
func addSigstorePolicyControllerToolsDirect(s *server.MCPServer) {
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
			mcp.Description("Path to Kubernetes resource file - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("trustroot",
			mcp.Description("Path to Kubernetes TrustRoot resource - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("log_level",
			mcp.Description("Log level for output - NOTE: not supported in current Dagger module"),
			mcp.Enum("debug", "info", "warn", "error"),
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
		module := modules.NewSigstorePolicyControllerModule(client)

		// Get parameters
		policy := request.GetString("policy", "")
		image := request.GetString("image", "")

		// Check for unsupported parameters
		if request.GetString("resource", "") != "" || request.GetString("trustroot", "") != "" ||
			request.GetString("log_level", "") != "" {
			return mcp.NewToolResultError("resource, trustroot, and log_level options not supported in current Dagger module"), nil
		}

		// Test policy
		output, err := module.TestPolicy(ctx, policy, image)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("sigstore test policy failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Get tester version
	versionTool := mcp.NewTool("sigstore_tester_version",
		mcp.WithDescription("Get policy-controller-tester version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSigstorePolicyControllerModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("sigstore get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSigstorePolicyControllerModule(client)

		// Get parameters
		policyFile := request.GetString("policy_file", "")
		namespace := request.GetString("namespace", "")

		// Create policy
		output, err := module.CreatePolicy(ctx, policyFile, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("sigstore create policy failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSigstorePolicyControllerModule(client)

		// Get parameters
		output := request.GetString("output", "")

		// List policies
		result, err := module.ListClusterImagePolicies(ctx, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("sigstore list policies failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSigstorePolicyControllerModule(client)

		// Get parameters
		policyName := request.GetString("policy_name", "")

		// Delete policy
		output, err := module.DeletePolicy(ctx, policyName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("sigstore delete policy failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSigstorePolicyControllerModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")

		// Enable namespace
		output, err := module.EnableNamespace(ctx, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("sigstore enable namespace failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSigstorePolicyControllerModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")

		// Disable namespace
		output, err := module.DisableNamespace(ctx, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("sigstore disable namespace failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSigstorePolicyControllerModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "")

		// Get namespace status
		output, err := module.GetNamespaceStatus(ctx, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("sigstore get namespace status failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSigstorePolicyControllerModule(client)

		// Get parameters
		policyName := request.GetString("policy_name", "")

		// Describe policy
		output, err := module.DescribePolicy(ctx, policyName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("sigstore describe policy failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}