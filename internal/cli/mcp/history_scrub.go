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

// AddHistoryScrubTools adds Git history cleanup MCP tool implementations using direct Dagger calls
func AddHistoryScrubTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addHistoryScrubToolsDirect(s)
}

// addHistoryScrubToolsDirect adds History Scrub tools using direct Dagger module calls
func addHistoryScrubToolsDirect(s *server.MCPServer) {
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
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform dry run without making changes"),
		),
	)
	s.AddTool(bfgRemoveLargeFilesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHistoryScrubModule(client)

		// Get parameters
		repoPath := request.GetString("repository_path", "")
		if repoPath == "" {
			return mcp.NewToolResultError("repository_path is required"), nil
		}
		
		sizeThreshold := request.GetString("size_threshold", "")
		if sizeThreshold == "" {
			return mcp.NewToolResultError("size_threshold is required"), nil
		}

		dryRun := request.GetBool("dry_run", false)

		// Use BFG to remove large files (creating secrets file with size pattern)
		secretsFile := fmt.Sprintf("--strip-blobs-bigger-than %s", sizeThreshold)
		output, err := module.RemoveSecretsWithBFG(ctx, repoPath, secretsFile, dryRun)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to remove large files: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform dry run without making changes"),
		),
	)
	s.AddTool(bfgReplaceTextTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHistoryScrubModule(client)

		// Get parameters
		repoPath := request.GetString("repository_path", "")
		if repoPath == "" {
			return mcp.NewToolResultError("repository_path is required"), nil
		}
		
		replacementsFile := request.GetString("replacements_file", "")
		if replacementsFile == "" {
			return mcp.NewToolResultError("replacements_file is required"), nil
		}

		dryRun := request.GetBool("dry_run", false)

		// Remove secrets with BFG
		output, err := module.RemoveSecretsWithBFG(ctx, repoPath, replacementsFile, dryRun)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to replace text: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform dry run without making changes"),
		),
	)
	s.AddTool(filterRepoRemovePathTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHistoryScrubModule(client)

		// Get parameters
		repoPath := request.GetString("repository_path", "")
		if repoPath == "" {
			return mcp.NewToolResultError("repository_path is required"), nil
		}
		
		pathToRemove := request.GetString("path_to_remove", "")
		if pathToRemove == "" {
			return mcp.NewToolResultError("path_to_remove is required"), nil
		}

		dryRun := request.GetBool("dry_run", false)

		// Prepare patterns file
		patternsFile := pathToRemove
		if request.GetBool("invert_paths", false) {
			patternsFile = "--invert-paths " + pathToRemove
		}

		// Remove secrets with git-filter
		output, err := module.RemoveSecretsWithGitFilter(ctx, repoPath, patternsFile, dryRun)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to remove path: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHistoryScrubModule(client)

		// Get parameters
		repoPath := request.GetString("repository_path", "")
		if repoPath == "" {
			return mcp.NewToolResultError("repository_path is required"), nil
		}
		
		searchPattern := request.GetString("search_pattern", "")
		if searchPattern == "" {
			return mcp.NewToolResultError("search_pattern is required"), nil
		}

		// Verify history is clean using the search pattern
		scanTool := "git-log"
		if request.GetBool("search_all_branches", false) {
			scanTool = fmt.Sprintf("git-log-all:%s", searchPattern)
		} else {
			scanTool = fmt.Sprintf("git-log:%s", searchPattern)
		}

		output, err := module.VerifyHistoryClean(ctx, repoPath, scanTool)
		if err != nil && !strings.Contains(err.Error(), "no matches found") {
			return mcp.NewToolResultError(fmt.Sprintf("failed to search history: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHistoryScrubModule(client)

		// Get parameters
		repoPath := request.GetString("repository_path", "")
		if repoPath == "" {
			return mcp.NewToolResultError("repository_path is required"), nil
		}
		
		backupPath := request.GetString("backup_path", "")
		if backupPath == "" {
			return mcp.NewToolResultError("backup_path is required"), nil
		}

		// Create bare clone for backup
		output, err := module.CreateBareClone(ctx, repoPath, backupPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create backup: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewHistoryScrubModule(client)

		// Get parameters
		repoPath := request.GetString("repository_path", "")
		if repoPath == "" {
			return mcp.NewToolResultError("repository_path is required"), nil
		}

		// Analyze repo size after cleanup
		output, err := module.AnalyzeRepoSize(ctx, repoPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to analyze repo: %v", err)), nil
		}

		// Add cleanup info
		if request.GetBool("aggressive", false) {
			output += "\n\nTo run aggressive garbage collection:\ngit reflog expire --expire=now --all && git gc --prune=now --aggressive"
		} else {
			output += "\n\nTo run garbage collection:\ngit reflog expire --expire=now --all && git gc --prune=now"
		}

		return mcp.NewToolResultText(output), nil
	})
}