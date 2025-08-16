package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddOSSFScorecardTools adds OSSF Scorecard MCP tool implementations
func AddOSSFScorecardTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// OSSF Scorecard score repository tool
	scoreRepositoryTool := mcp.NewTool("ossf_scorecard_score_repository",
		mcp.WithDescription("Score repository security using OSSF Scorecard"),
		mcp.WithString("repo_url",
			mcp.Description("Repository URL to score"),
			mcp.Required(),
		),
		mcp.WithString("github_token",
			mcp.Description("GitHub token for API access"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("default", "json", "sarif"),
		),
	)
	s.AddTool(scoreRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoURL := request.GetString("repo_url", "")
		githubToken := request.GetString("github_token", "")
		args := []string{"security", "ossf-scorecard", "--repo", repoURL, "--token", githubToken}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// OSSF Scorecard score with specific checks tool
	scoreWithChecksTool := mcp.NewTool("ossf_scorecard_score_checks",
		mcp.WithDescription("Score repository with specific security checks"),
		mcp.WithString("repo_url",
			mcp.Description("Repository URL to score"),
			mcp.Required(),
		),
		mcp.WithString("checks",
			mcp.Description("Comma-separated list of specific checks to run"),
			mcp.Required(),
		),
		mcp.WithString("github_token",
			mcp.Description("GitHub token for API access"),
			mcp.Required(),
		),
	)
	s.AddTool(scoreWithChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoURL := request.GetString("repo_url", "")
		checks := request.GetString("checks", "")
		githubToken := request.GetString("github_token", "")
		args := []string{"security", "ossf-scorecard", "--repo", repoURL, "--checks", checks, "--token", githubToken}
		return executeShipCommand(args)
	})

	// OSSF Scorecard list checks tool
	listChecksTool := mcp.NewTool("ossf_scorecard_list_checks",
		mcp.WithDescription("List all available OSSF Scorecard security checks"),
	)
	s.AddTool(listChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "ossf-scorecard", "--list-checks"}
		return executeShipCommand(args)
	})

	// OSSF Scorecard get version tool
	getVersionTool := mcp.NewTool("ossf_scorecard_get_version",
		mcp.WithDescription("Get OSSF Scorecard version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "ossf-scorecard", "--version"}
		return executeShipCommand(args)
	})
}