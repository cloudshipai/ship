package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddFleetTools adds Fleet GitOps MCP tool implementations using kubectl
func AddFleetTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addFleetToolsDirect(s)
}

// addFleetToolsDirect implements direct Dagger calls for Fleet tools
func addFleetToolsDirect(s *server.MCPServer) {
	// Fleet apply GitRepo tool
	applyGitRepoTool := mcp.NewTool("fleet_apply_gitrepo",
		mcp.WithDescription("Apply Fleet GitRepo configuration using kubectl"),
		mcp.WithString("gitrepo_file",
			mcp.Description("Path to GitRepo YAML file"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Target Kubernetes namespace (default: fleet-local)"),
		),
	)
	s.AddTool(applyGitRepoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		gitrepoFile := request.GetString("gitrepo_file", "")
		namespace := request.GetString("namespace", "")

		// Create Fleet module and apply GitRepo
		fleetModule := modules.NewFleetModule(client)
		result, err := fleetModule.ApplyGitRepo(ctx, gitrepoFile, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fleet apply gitrepo failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Fleet get GitRepos status tool
	getGitReposTool := mcp.NewTool("fleet_get_gitrepos",
		mcp.WithDescription("Get Fleet GitRepos status using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace (default: fleet-local)"),
		),
	)
	s.AddTool(getGitReposTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		namespace := request.GetString("namespace", "fleet-local")

		// Create Fleet module and get GitRepos
		fleetModule := modules.NewFleetModule(client)
		result, err := fleetModule.GetGitReposSimple(ctx, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fleet get gitrepos failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Fleet get bundles tool
	getBundlesTool := mcp.NewTool("fleet_get_bundles",
		mcp.WithDescription("Get Fleet Bundles status using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(getBundlesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		namespace := request.GetString("namespace", "")

		// Create Fleet module and get bundles
		fleetModule := modules.NewFleetModule(client)
		result, err := fleetModule.GetBundlesSimple(ctx, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fleet get bundles failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Fleet get bundle deployments tool
	getBundleDeploymentsTool := mcp.NewTool("fleet_get_bundledeployments",
		mcp.WithDescription("Get Fleet BundleDeployments status using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(getBundleDeploymentsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		namespace := request.GetString("namespace", "")

		// Create Fleet module and get bundle deployments
		fleetModule := modules.NewFleetModule(client)
		result, err := fleetModule.GetBundleDeploymentsSimple(ctx, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fleet get bundledeployments failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Fleet describe GitRepo tool
	describeGitRepoTool := mcp.NewTool("fleet_describe_gitrepo",
		mcp.WithDescription("Describe Fleet GitRepo using kubectl"),
		mcp.WithString("gitrepo_name",
			mcp.Description("GitRepo name to describe"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace (default: fleet-local)"),
		),
	)
	s.AddTool(describeGitRepoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		gitrepoName := request.GetString("gitrepo_name", "")
		namespace := request.GetString("namespace", "fleet-local")

		// Create Fleet module and describe GitRepo
		fleetModule := modules.NewFleetModule(client)
		result, err := fleetModule.DescribeGitRepoSimple(ctx, gitrepoName, namespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fleet describe gitrepo failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Fleet install tool (using Helm)
	installTool := mcp.NewTool("fleet_install",
		mcp.WithDescription("Install Fleet using Helm"),
		mcp.WithString("version",
			mcp.Description("Fleet version to install (e.g., v0.13.0)"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		version := request.GetString("version", "v0.13.0")

		// Create Fleet module and install
		fleetModule := modules.NewFleetModule(client)
		result, err := fleetModule.InstallWithVersion(ctx, version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fleet install failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}