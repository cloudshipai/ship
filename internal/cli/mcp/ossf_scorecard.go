package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddOSSFScorecardTools adds OSSF Scorecard MCP tool implementations using direct Dagger calls
func AddOSSFScorecardTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addOSSFScorecardToolsDirect(s)
}

// addOSSFScorecardToolsDirect adds OSSF Scorecard tools using direct Dagger module calls
func addOSSFScorecardToolsDirect(s *server.MCPServer) {
	// OSSF Scorecard score repository tool
	scoreRepositoryTool := mcp.NewTool("ossf_scorecard_score_repository",
		mcp.WithDescription("Score repository security using OSSF Scorecard"),
		mcp.WithString("repo_url",
			mcp.Description("Repository URL to score"),
			mcp.Required(),
		),
		mcp.WithString("github_token",
			mcp.Description("GitHub token for API access (optional if using environment variable)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("default", "json", "sarif"),
		),
	)
	s.AddTool(scoreRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSSFScorecardModule(client)

		// Get parameters
		repoURL := request.GetString("repo_url", "")
		if repoURL == "" {
			return mcp.NewToolResultError("repo_url is required"), nil
		}
		githubToken := request.GetString("github_token", "")

		// Score repository
		output, err := module.ScoreRepository(ctx, repoURL, githubToken)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scorecard failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
			mcp.Description("GitHub token for API access (optional if using environment variable)"),
		),
	)
	s.AddTool(scoreWithChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSSFScorecardModule(client)

		// Get parameters
		repoURL := request.GetString("repo_url", "")
		if repoURL == "" {
			return mcp.NewToolResultError("repo_url is required"), nil
		}
		checksStr := request.GetString("checks", "")
		if checksStr == "" {
			return mcp.NewToolResultError("checks is required"), nil
		}
		githubToken := request.GetString("github_token", "")

		// Parse comma-separated checks
		checks := strings.Split(checksStr, ",")
		for i, check := range checks {
			checks[i] = strings.TrimSpace(check)
		}

		// Score with checks
		output, err := module.ScoreWithChecks(ctx, repoURL, checks, githubToken)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scorecard with checks failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OSSF Scorecard list checks tool
	listChecksTool := mcp.NewTool("ossf_scorecard_list_checks",
		mcp.WithDescription("List all available OSSF Scorecard security checks"),
	)
	s.AddTool(listChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSSFScorecardModule(client)

		// List checks
		output, err := module.ListChecks(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("list checks failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OSSF Scorecard get version tool
	getVersionTool := mcp.NewTool("ossf_scorecard_get_version",
		mcp.WithDescription("Get OSSF Scorecard version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSSFScorecardModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}