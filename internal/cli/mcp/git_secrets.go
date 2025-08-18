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
		// Change to repository directory first, then run git secrets
		args := []string{"sh", "-c", "cd " + repoPath + " && git secrets --scan"}
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
		args := []string{"sh", "-c", "cd " + repoPath + " && git secrets --scan-history"}
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
		args := []string{"sh", "-c", "cd " + repoPath + " && git secrets --install"}
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
		args := []string{"sh", "-c", "cd " + repoPath + " && git secrets --add '" + pattern + "'"}
		return executeShipCommand(args)
	})

	// Git-secrets register AWS patterns
	registerAwsTool := mcp.NewTool("git_secrets_register_aws",
		mcp.WithDescription("Register AWS-specific secret patterns"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
	)
	s.AddTool(registerAwsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		args := []string{"sh", "-c", "cd " + repoPath + " && git secrets --register-aws"}
		return executeShipCommand(args)
	})

	// Git-secrets list configuration
	listConfigTool := mcp.NewTool("git_secrets_list_config",
		mcp.WithDescription("List current git-secrets configuration"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
	)
	s.AddTool(listConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		args := []string{"sh", "-c", "cd " + repoPath + " && git secrets --list"}
		return executeShipCommand(args)
	})

	// Git-secrets add allowed pattern
	addAllowedTool := mcp.NewTool("git_secrets_add_allowed",
		mcp.WithDescription("Add allowed pattern to prevent false positives"),
		mcp.WithString("pattern",
			mcp.Description("Pattern to allow"),
			mcp.Required(),
		),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
	)
	s.AddTool(addAllowedTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pattern := request.GetString("pattern", "")
		repoPath := request.GetString("repository_path", "")
		args := []string{"sh", "-c", "cd " + repoPath + " && git secrets --add -a '" + pattern + "'"}
		return executeShipCommand(args)
	})
}