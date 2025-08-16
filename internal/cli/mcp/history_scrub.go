package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddHistoryScrubTools adds History Scrub (Git history cleanup) MCP tool implementations
func AddHistoryScrubTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// History Scrub scan repository tool
	scanRepositoryTool := mcp.NewTool("history_scrub_scan_repository",
		mcp.WithDescription("Scan git repository for sensitive data in history"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository (default: current directory)"),
		),
		mcp.WithString("patterns_file",
			mcp.Description("Path to file containing search patterns"),
		),
	)
	s.AddTool(scanRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "history-scrub", "scan"}
		if repoPath := request.GetString("repository_path", ""); repoPath != "" {
			args = append(args, repoPath)
		}
		if patternsFile := request.GetString("patterns_file", ""); patternsFile != "" {
			args = append(args, "--patterns", patternsFile)
		}
		return executeShipCommand(args)
	})

	// History Scrub clean secrets tool
	cleanSecretsTool := mcp.NewTool("history_scrub_clean_secrets",
		mcp.WithDescription("Remove sensitive data from git history"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithString("patterns_file",
			mcp.Description("Path to file containing patterns to remove"),
			mcp.Required(),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform a dry run without making changes"),
		),
	)
	s.AddTool(cleanSecretsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		patternsFile := request.GetString("patterns_file", "")
		args := []string{"security", "history-scrub", "clean", repoPath, "--patterns", patternsFile}
		if request.GetBool("dry_run", false) {
			args = append(args, "--dry-run")
		}
		return executeShipCommand(args)
	})

	// History Scrub remove file tool
	removeFileTool := mcp.NewTool("history_scrub_remove_file",
		mcp.WithDescription("Remove specific file from git history"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithString("file_path",
			mcp.Description("Path to file to remove from history"),
			mcp.Required(),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force removal without confirmation"),
		),
	)
	s.AddTool(removeFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		filePath := request.GetString("file_path", "")
		args := []string{"security", "history-scrub", "remove-file", repoPath, filePath}
		if request.GetBool("force", false) {
			args = append(args, "--force")
		}
		return executeShipCommand(args)
	})

	// History Scrub verify cleanup tool
	verifyCleanupTool := mcp.NewTool("history_scrub_verify_cleanup",
		mcp.WithDescription("Verify that sensitive data has been removed from history"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithString("patterns_file",
			mcp.Description("Path to file containing verification patterns"),
		),
	)
	s.AddTool(verifyCleanupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		args := []string{"security", "history-scrub", "verify", repoPath}
		if patternsFile := request.GetString("patterns_file", ""); patternsFile != "" {
			args = append(args, "--patterns", patternsFile)
		}
		return executeShipCommand(args)
	})

	// History Scrub backup repository tool
	backupRepositoryTool := mcp.NewTool("history_scrub_backup_repository",
		mcp.WithDescription("Create backup of repository before cleanup"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithString("backup_path",
			mcp.Description("Path for backup location"),
			mcp.Required(),
		),
	)
	s.AddTool(backupRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		backupPath := request.GetString("backup_path", "")
		args := []string{"security", "history-scrub", "backup", repoPath, backupPath}
		return executeShipCommand(args)
	})
}