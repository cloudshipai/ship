package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddGitSecretsTools adds git-secrets (AWS secret scanning) MCP tool implementations using direct Dagger calls
func AddGitSecretsTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGitSecretsToolsDirect(s)
}

// addGitSecretsToolsDirect adds git-secrets tools using direct Dagger module calls
func addGitSecretsToolsDirect(s *server.MCPServer) {
	// Git-secrets scan tool
	scanTool := mcp.NewTool("git_secrets_scan",
		mcp.WithDescription("Scan git repository for AWS secrets and credentials"),
		mcp.WithString("source_path",
			mcp.Description("Path to git repository to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("text", "json"),
		),
		mcp.WithBoolean("scan_history",
			mcp.Description("Scan entire git history"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively scan subdirectories"),
		),
	)
	s.AddTool(scanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		sourcePath := request.GetString("source_path", "")
		outputFormat := request.GetString("output_format", "")
		scanHistory := request.GetBool("scan_history", false)
		recursive := request.GetBool("recursive", false)

		if sourcePath == "" {
			return mcp.NewToolResultError("source_path is required"), nil
		}

		// Create git-secrets module
		gitSecretsModule := modules.NewGitSecretsModule(client)

		// Set up options
		opts := modules.GitSecretsScanOptions{
			OutputFormat: outputFormat,
			ScanHistory:  scanHistory,
			Recursive:    recursive,
		}

		// Run scan
		result, err := gitSecretsModule.Scan(ctx, sourcePath, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("git-secrets scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Git-secrets install hooks tool
	installHooksTool := mcp.NewTool("git_secrets_install_hooks",
		mcp.WithDescription("Install git-secrets hooks in a repository"),
		mcp.WithString("repo_path",
			mcp.Description("Path to git repository"),
			mcp.Required(),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force installation, overwriting existing hooks"),
		),
	)
	s.AddTool(installHooksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		repoPath := request.GetString("repo_path", "")
		force := request.GetBool("force", false)

		if repoPath == "" {
			return mcp.NewToolResultError("repo_path is required"), nil
		}

		// Create git-secrets module
		gitSecretsModule := modules.NewGitSecretsModule(client)

		// Install hooks
		result, err := gitSecretsModule.InstallHooks(ctx, repoPath, force)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("git-secrets hook installation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}