package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddFleetTools adds Fleet GitOps MCP tool implementations using kubectl
func AddFleetTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		gitrepoFile := request.GetString("gitrepo_file", "")
		args := []string{"kubectl", "apply", "-f", gitrepoFile}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		}
		return executeShipCommand(args)
	})

	// Fleet get GitRepos status tool
	getGitReposTool := mcp.NewTool("fleet_get_gitrepos",
		mcp.WithDescription("Get Fleet GitRepos status using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace (default: fleet-local)"),
		),
	)
	s.AddTool(getGitReposTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "fleet-local")
		args := []string{"kubectl", "get", "gitrepo", "-n", namespace}
		return executeShipCommand(args)
	})

	// Fleet get bundles tool
	getBundlesTool := mcp.NewTool("fleet_get_bundles",
		mcp.WithDescription("Get Fleet Bundles status using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(getBundlesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubectl", "get", "bundles"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		} else {
			args = append(args, "--all-namespaces")
		}
		return executeShipCommand(args)
	})

	// Fleet get bundle deployments tool
	getBundleDeploymentsTool := mcp.NewTool("fleet_get_bundledeployments",
		mcp.WithDescription("Get Fleet BundleDeployments status using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(getBundleDeploymentsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubectl", "get", "bundledeployments"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		} else {
			args = append(args, "--all-namespaces")
		}
		return executeShipCommand(args)
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
		gitrepoName := request.GetString("gitrepo_name", "")
		namespace := request.GetString("namespace", "fleet-local")
		args := []string{"kubectl", "describe", "gitrepo", gitrepoName, "-n", namespace}
		return executeShipCommand(args)
	})

	// Fleet install tool (using Helm)
	installTool := mcp.NewTool("fleet_install",
		mcp.WithDescription("Install Fleet using Helm"),
		mcp.WithString("version",
			mcp.Description("Fleet version to install (e.g., v0.13.0)"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		version := request.GetString("version", "v0.13.0")
		
		// Install Fleet CRD first
		args := []string{"helm", "-n", "cattle-fleet-system", "install", "--create-namespace", "--wait",
			"fleet-crd", "https://github.com/rancher/fleet/releases/download/" + version + "/fleet-crd-" + version[1:] + ".tgz"}
		return executeShipCommand(args)
	})
}