package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitHubAdminTools adds GitHub administration MCP tool implementations using gh CLI
func AddGitHubAdminTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGitHubAdminToolsDirect(s)
}

// addGitHubAdminToolsDirect implements direct Dagger calls for GitHub Admin tools
func addGitHubAdminToolsDirect(s *server.MCPServer) {
	// GitHub list organization repositories
	listOrgReposTool := mcp.NewTool("github_admin_list_org_repos",
		mcp.WithDescription("List GitHub organization repositories using gh CLI"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("visibility",
			mcp.Description("Repository visibility filter"),
			mcp.Enum("public", "private", "internal"),
		),
	)
	s.AddTool(listOrgReposTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		organization := request.GetString("organization", "")
		visibility := request.GetString("visibility", "")

		// Create GitHub Admin module and list org repos
		githubAdminModule := modules.NewGitHubAdminModule(client)
		result, err := githubAdminModule.ListOrgReposSimple(ctx, organization, visibility)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github admin list org repos failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// GitHub create repository in organization
	createOrgRepoTool := mcp.NewTool("github_admin_create_org_repo",
		mcp.WithDescription("Create repository in GitHub organization"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name"),
			mcp.Required(),
		),
		mcp.WithString("visibility",
			mcp.Description("Repository visibility"),
			mcp.Enum("public", "private", "internal"),
		),
		mcp.WithString("description",
			mcp.Description("Repository description"),
		),
	)
	s.AddTool(createOrgRepoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		organization := request.GetString("organization", "")
		repoName := request.GetString("repo_name", "")
		visibility := request.GetString("visibility", "")
		description := request.GetString("description", "")

		// Create GitHub Admin module and create org repo
		githubAdminModule := modules.NewGitHubAdminModule(client)
		result, err := githubAdminModule.CreateOrgRepoSimple(ctx, organization, repoName, visibility, description)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github admin create org repo failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// GitHub get repository information
	getRepoInfoTool := mcp.NewTool("github_admin_get_repo_info",
		mcp.WithDescription("Get GitHub repository information"),
		mcp.WithString("repository",
			mcp.Description("Repository in format owner/repo"),
			mcp.Required(),
		),
	)
	s.AddTool(getRepoInfoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		repository := request.GetString("repository", "")

		// Create GitHub Admin module and get repo info
		githubAdminModule := modules.NewGitHubAdminModule(client)
		result, err := githubAdminModule.GetRepoInfoSimple(ctx, repository)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github admin get repo info failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// GitHub list issues in organization
	listOrgIssuesTool := mcp.NewTool("github_admin_list_org_issues",
		mcp.WithDescription("List issues across organization repositories"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("state",
			mcp.Description("Issue state filter"),
			mcp.Enum("open", "closed", "all"),
		),
	)
	s.AddTool(listOrgIssuesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		organization := request.GetString("organization", "")
		state := request.GetString("state", "")

		// Create GitHub Admin module and list org issues
		githubAdminModule := modules.NewGitHubAdminModule(client)
		result, err := githubAdminModule.ListOrgIssuesSimple(ctx, organization, state)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github admin list org issues failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// GitHub list pull requests in organization
	listOrgPRsTool := mcp.NewTool("github_admin_list_org_prs",
		mcp.WithDescription("List pull requests across organization repositories"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("state",
			mcp.Description("PR state filter"),
			mcp.Enum("open", "closed", "merged", "all"),
		),
	)
	s.AddTool(listOrgPRsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		organization := request.GetString("organization", "")
		state := request.GetString("state", "")

		// Create GitHub Admin module and list org PRs
		githubAdminModule := modules.NewGitHubAdminModule(client)
		result, err := githubAdminModule.ListOrgPRsSimple(ctx, organization, state)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github admin list org prs failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// GitHub CLI version
	getVersionTool := mcp.NewTool("github_admin_get_version",
		mcp.WithDescription("Get GitHub CLI version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create GitHub Admin module and get version
		githubAdminModule := modules.NewGitHubAdminModule(client)
		result, err := githubAdminModule.GetVersionSimple(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github admin get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}