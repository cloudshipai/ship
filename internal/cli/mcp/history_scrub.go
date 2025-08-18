package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddHistoryScrubTools adds Git history cleanup MCP tool implementations using real tools
func AddHistoryScrubTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// BFG Repo Cleaner - remove large files
	bfgRemoveLargeFilesTool := mcp.NewTool("history_scrub_bfg_remove_large_files",
		mcp.WithDescription("Remove large files from git history using BFG Repo-Cleaner"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository (.git directory)"),
			mcp.Required(),
		),
		mcp.WithString("size_threshold",
			mcp.Description("Size threshold (e.g., 1M, 100K, 50M)"),
			mcp.Required(),
		),
	)
	s.AddTool(bfgRemoveLargeFilesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		sizeThreshold := request.GetString("size_threshold", "")
		args := []string{"bfg", "--strip-blobs-bigger-than", sizeThreshold, repoPath}
		return executeShipCommand(args)
	})

	// BFG Repo Cleaner - replace text/secrets
	bfgReplaceTextTool := mcp.NewTool("history_scrub_bfg_replace_text",
		mcp.WithDescription("Replace sensitive text in git history using BFG Repo-Cleaner"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository (.git directory)"),
			mcp.Required(),
		),
		mcp.WithString("replacements_file",
			mcp.Description("Path to file containing text replacements (format: secret==>REMOVED)"),
			mcp.Required(),
		),
	)
	s.AddTool(bfgReplaceTextTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		replacementsFile := request.GetString("replacements_file", "")
		args := []string{"bfg", "--replace-text", replacementsFile, repoPath}
		return executeShipCommand(args)
	})

	// Git filter-repo - remove files/paths
	filterRepoRemovePathTool := mcp.NewTool("history_scrub_filter_repo_remove_path",
		mcp.WithDescription("Remove files/paths from git history using git-filter-repo"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithString("path_to_remove",
			mcp.Description("Path/pattern to remove from history"),
			mcp.Required(),
		),
		mcp.WithBoolean("invert_paths",
			mcp.Description("Keep only the specified paths (remove everything else)"),
		),
	)
	s.AddTool(filterRepoRemovePathTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		pathToRemove := request.GetString("path_to_remove", "")
		
		// Change to repository directory and run git-filter-repo
		args := []string{"sh", "-c", "cd " + repoPath + " && git filter-repo --path " + pathToRemove}
		
		if request.GetBool("invert_paths", false) {
			// Use --path with invert to keep only the specified paths
			args = []string{"sh", "-c", "cd " + repoPath + " && git filter-repo --path " + pathToRemove + " --invert-paths"}
		}
		
		return executeShipCommand(args)
	})

	// Git history search using git log
	searchHistoryTool := mcp.NewTool("history_scrub_search_history",
		mcp.WithDescription("Search git history for sensitive patterns using git log"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithString("search_pattern",
			mcp.Description("Pattern to search for in commit history"),
			mcp.Required(),
		),
		mcp.WithBoolean("search_all_branches",
			mcp.Description("Search all branches (default: current branch only)"),
		),
	)
	s.AddTool(searchHistoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		searchPattern := request.GetString("search_pattern", "")
		
		gitArgs := "git log -S\"" + searchPattern + "\" --oneline"
		if request.GetBool("search_all_branches", false) {
			gitArgs += " --all"
		}
		
		args := []string{"sh", "-c", "cd " + repoPath + " && " + gitArgs}
		return executeShipCommand(args)
	})

	// Git repository backup using git clone
	backupRepositoryTool := mcp.NewTool("history_scrub_backup_repository",
		mcp.WithDescription("Create backup of repository before cleanup using git clone"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithString("backup_path",
			mcp.Description("Path for backup location"),
			mcp.Required(),
		),
		mcp.WithBoolean("bare_clone",
			mcp.Description("Create bare clone backup"),
		),
	)
	s.AddTool(backupRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		backupPath := request.GetString("backup_path", "")
		
		if request.GetBool("bare_clone", false) {
			args := []string{"git", "clone", "--bare", repoPath, backupPath}
			return executeShipCommand(args)
		} else {
			args := []string{"git", "clone", repoPath, backupPath}
			return executeShipCommand(args)
		}
	})

	// Git cleanup post-processing
	gitCleanupTool := mcp.NewTool("history_scrub_git_cleanup",
		mcp.WithDescription("Run git cleanup commands after history rewriting"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithBoolean("aggressive",
			mcp.Description("Run aggressive garbage collection"),
		),
	)
	s.AddTool(gitCleanupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoPath := request.GetString("repository_path", "")
		
		gcArgs := "git reflog expire --expire=now --all && git gc --prune=now"
		if request.GetBool("aggressive", false) {
			gcArgs += " --aggressive"
		}
		
		args := []string{"sh", "-c", "cd " + repoPath + " && " + gcArgs}
		return executeShipCommand(args)
	})
}