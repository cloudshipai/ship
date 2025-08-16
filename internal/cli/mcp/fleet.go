package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddFleetTools adds Fleet GitOps MCP tool implementations
func AddFleetTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Fleet deploy tool
	deployTool := mcp.NewTool("fleet_deploy",
		mcp.WithDescription("Deploy applications using Fleet GitOps"),
		mcp.WithString("git_repo",
			mcp.Description("Git repository URL containing Fleet configuration"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Target Kubernetes namespace"),
		),
		mcp.WithString("cluster",
			mcp.Description("Target cluster name"),
		),
	)
	s.AddTool(deployTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		gitRepo := request.GetString("git_repo", "")
		args := []string{"kubernetes", "fleet", "deploy", gitRepo}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if cluster := request.GetString("cluster", ""); cluster != "" {
			args = append(args, "--cluster", cluster)
		}
		return executeShipCommand(args)
	})

	// Fleet status tool
	statusTool := mcp.NewTool("fleet_status",
		mcp.WithDescription("Check Fleet deployment status"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to check"),
		),
	)
	s.AddTool(statusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "fleet", "status"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		return executeShipCommand(args)
	})

	// Fleet sync tool
	syncTool := mcp.NewTool("fleet_sync",
		mcp.WithDescription("Force synchronization of Fleet managed resources"),
		mcp.WithString("app_name",
			mcp.Description("Application name to sync"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(syncTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appName := request.GetString("app_name", "")
		args := []string{"kubernetes", "fleet", "sync", appName}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		return executeShipCommand(args)
	})

	// Fleet rollback tool
	rollbackTool := mcp.NewTool("fleet_rollback",
		mcp.WithDescription("Rollback Fleet deployment to previous version"),
		mcp.WithString("app_name",
			mcp.Description("Application name to rollback"),
			mcp.Required(),
		),
		mcp.WithString("revision",
			mcp.Description("Target revision for rollback"),
		),
	)
	s.AddTool(rollbackTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		appName := request.GetString("app_name", "")
		args := []string{"kubernetes", "fleet", "rollback", appName}
		if revision := request.GetString("revision", ""); revision != "" {
			args = append(args, "--revision", revision)
		}
		return executeShipCommand(args)
	})

	// Fleet validate tool
	validateTool := mcp.NewTool("fleet_validate",
		mcp.WithDescription("Validate Fleet configuration and manifests"),
		mcp.WithString("config_path",
			mcp.Description("Path to Fleet configuration"),
			mcp.Required(),
		),
	)
	s.AddTool(validateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configPath := request.GetString("config_path", "")
		args := []string{"kubernetes", "fleet", "validate", configPath}
		return executeShipCommand(args)
	})

	// Fleet get version tool
	getVersionTool := mcp.NewTool("fleet_get_version",
		mcp.WithDescription("Get Fleet version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "fleet", "--version"}
		return executeShipCommand(args)
	})
}