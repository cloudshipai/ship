package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCloudQueryTools adds CloudQuery (cloud asset inventory) MCP tool implementations
func AddCloudQueryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// CloudQuery sync assets tool
	syncAssetsTool := mcp.NewTool("cloudquery_sync_assets",
		mcp.WithDescription("Sync cloud assets to database using CloudQuery"),
		mcp.WithString("config_path",
			mcp.Description("Path to CloudQuery configuration file"),
			mcp.Required(),
		),
		mcp.WithString("destination",
			mcp.Description("Destination for synced data"),
			mcp.Enum("postgres", "sqlite", "mysql", "bigquery"),
		),
	)
	s.AddTool(syncAssetsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configPath := request.GetString("config_path", "")
		args := []string{"cloud", "cloudquery", "sync", configPath}
		if destination := request.GetString("destination", ""); destination != "" {
			args = append(args, "--destination", destination)
		}
		return executeShipCommand(args)
	})

	// CloudQuery migrate schema tool
	migrateSchemaTool := mcp.NewTool("cloudquery_migrate_schema",
		mcp.WithDescription("Migrate CloudQuery database schema"),
		mcp.WithString("config_path",
			mcp.Description("Path to CloudQuery configuration file"),
			mcp.Required(),
		),
		mcp.WithString("direction",
			mcp.Description("Migration direction"),
			mcp.Enum("up", "down"),
		),
	)
	s.AddTool(migrateSchemaTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configPath := request.GetString("config_path", "")
		args := []string{"cloud", "cloudquery", "migrate", configPath}
		if direction := request.GetString("direction", ""); direction != "" {
			args = append(args, "--direction", direction)
		}
		return executeShipCommand(args)
	})

	// CloudQuery validate config tool
	validateConfigTool := mcp.NewTool("cloudquery_validate_config",
		mcp.WithDescription("Validate CloudQuery configuration file"),
		mcp.WithString("config_path",
			mcp.Description("Path to CloudQuery configuration file"),
			mcp.Required(),
		),
	)
	s.AddTool(validateConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configPath := request.GetString("config_path", "")
		args := []string{"cloud", "cloudquery", "validate", configPath}
		return executeShipCommand(args)
	})

	// CloudQuery list plugins tool
	listPluginsTool := mcp.NewTool("cloudquery_list_plugins",
		mcp.WithDescription("List available CloudQuery plugins"),
		mcp.WithString("plugin_type",
			mcp.Description("Type of plugins to list"),
			mcp.Enum("source", "destination", "all"),
		),
	)
	s.AddTool(listPluginsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloud", "cloudquery", "plugins", "list"}
		if pluginType := request.GetString("plugin_type", ""); pluginType != "" {
			args = append(args, "--type", pluginType)
		}
		return executeShipCommand(args)
	})

	// CloudQuery get version tool
	getVersionTool := mcp.NewTool("cloudquery_get_version",
		mcp.WithDescription("Get CloudQuery version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloud", "cloudquery", "--version"}
		return executeShipCommand(args)
	})
}