package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitSecretsTools adds Git-secrets (git repository secret scanner) MCP tool implementations
func AddGitSecretsTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGitSecretsToolsDirect(s)
}

// addGitSecretsToolsDirect implements direct Dagger calls for Git-secrets tools
func addGitSecretsToolsDirect(s *server.MCPServer) {
	// Git-secrets scan repository tool
	scanRepositoryTool := mcp.NewTool("git_secrets_scan_repository",
		mcp.WithDescription("Scan git repository for secrets using git-secrets"),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		repoPath := request.GetString("repository_path", "")

		// Create Git-secrets module and scan repository
		gitSecretsModule := modules.NewGitSecretsModule(client)
		result, err := gitSecretsModule.ScanRepository(ctx, repoPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("git-secrets scan repository failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		repoPath := request.GetString("repository_path", "")

		// Create Git-secrets module and scan history
		gitSecretsModule := modules.NewGitSecretsModule(client)
		result, err := gitSecretsModule.ScanHistory(ctx, repoPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("git-secrets scan history failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		repoPath := request.GetString("repository_path", "")

		// Create Git-secrets module and install hooks
		gitSecretsModule := modules.NewGitSecretsModule(client)
		result, err := gitSecretsModule.InstallHooks(ctx, repoPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("git-secrets install hooks failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		pattern := request.GetString("pattern", "")
		repoPath := request.GetString("repository_path", "")

		// Create Git-secrets module and add pattern
		gitSecretsModule := modules.NewGitSecretsModule(client)
		result, err := gitSecretsModule.AddPattern(ctx, repoPath, pattern, false)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("git-secrets add pattern failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		repoPath := request.GetString("repository_path", "")

		// Create Git-secrets module and register AWS patterns
		gitSecretsModule := modules.NewGitSecretsModule(client)
		result, err := gitSecretsModule.RegisterAWS(ctx, repoPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("git-secrets register aws failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		repoPath := request.GetString("repository_path", "")

		// Create Git-secrets module and list config
		gitSecretsModule := modules.NewGitSecretsModule(client)
		result, err := gitSecretsModule.ListConfig(ctx, repoPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("git-secrets list config failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		pattern := request.GetString("pattern", "")
		repoPath := request.GetString("repository_path", "")

		// Create Git-secrets module and add allowed pattern
		gitSecretsModule := modules.NewGitSecretsModule(client)
		result, err := gitSecretsModule.AddAllowedPattern(ctx, repoPath, pattern)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("git-secrets add allowed failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}