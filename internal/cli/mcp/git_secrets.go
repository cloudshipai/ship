package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitSecretsTools adds Git-secrets (git repository secret scanner) MCP tool implementations
func AddGitSecretsTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Git-secrets scan repository tool
	scanRepositoryTool := mcp.NewTool("git_secrets_scan_repository",
		mcp.WithDescription("Scan git repository for secrets using git-secrets"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		args := []string{"security", "git-secrets", "--scan", repoPath}
		return executeShipCommand(args)
	})

	// Git-secrets scan history tool
	scanHistoryTool := mcp.NewTool("git_secrets_scan_history",
		mcp.WithDescription("Scan git repository history for secrets"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
	)
	s.AddTool(scanHistoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		args := []string{"security", "git-secrets", "--scan-history", repoPath}
		return executeShipCommand(args)
	})

	// Git-secrets install hooks tool
	installHooksTool := mcp.NewTool("git_secrets_install_hooks",
		mcp.WithDescription("Install git-secrets hooks in repository"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
	)
	s.AddTool(installHooksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		args := []string{"security", "git-secrets", "--install", repoPath}
		return executeShipCommand(args)
	})

	// Git-secrets add pattern tool
	addPatternTool := mcp.NewTool("git_secrets_add_pattern",
		mcp.WithDescription("Add secret pattern to git-secrets"),
		mcp.WithString("pattern",
			mcp.Description("Regular expression pattern to detect"),
			mcp.Required(),
		),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
	)
	s.AddTool(addPatternTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pattern := request.GetString("pattern", "")
		repoPath := request.GetString("repository_path", "")
		args := []string{"security", "git-secrets", "--add", pattern, repoPath}
		return executeShipCommand(args)
	})
}