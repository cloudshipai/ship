package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitHubAdminTools adds GitHub administration MCP tool implementations using gh CLI
func AddGitHubAdminTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		organization := request.GetString("organization", "")
		args := []string{"gh", "repo", "list", organization}
		if visibility := request.GetString("visibility", ""); visibility != "" {
			args = append(args, "--visibility", visibility)
		}
		return executeShipCommand(args)
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
		organization := request.GetString("organization", "")
		repoName := request.GetString("repo_name", "")
		args := []string{"gh", "repo", "create", organization + "/" + repoName}
		
		if visibility := request.GetString("visibility", ""); visibility != "" {
			args = append(args, "--" + visibility)
		}
		if description := request.GetString("description", ""); description != "" {
			args = append(args, "--description", description)
		}
		
		return executeShipCommand(args)
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
		repository := request.GetString("repository", "")
		args := []string{"gh", "repo", "view", repository}
		return executeShipCommand(args)
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
		organization := request.GetString("organization", "")
		args := []string{"gh", "issue", "list", "--search", "org:" + organization}
		if state := request.GetString("state", ""); state != "" {
			args = append(args, "--state", state)
		}
		return executeShipCommand(args)
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
		organization := request.GetString("organization", "")
		args := []string{"gh", "pr", "list", "--search", "org:" + organization}
		if state := request.GetString("state", ""); state != "" {
			args = append(args, "--state", state)
		}
		return executeShipCommand(args)
	})

	// GitHub CLI version
	getVersionTool := mcp.NewTool("github_admin_get_version",
		mcp.WithDescription("Get GitHub CLI version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"gh", "--version"}
		return executeShipCommand(args)
	})
}