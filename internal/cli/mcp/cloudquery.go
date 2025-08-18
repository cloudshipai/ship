package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCloudQueryTools adds CloudQuery (cloud asset inventory) MCP tool implementations
func AddCloudQueryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// CloudQuery sync resources from sources to destinations
	syncTool := mcp.NewTool("cloudquery_sync",
		mcp.WithDescription("Sync resources from source plugins to destinations"),
		mcp.WithString("config_path",
			mcp.Description("Path to configuration file or directory"),
			mcp.Required(),
		),
		mcp.WithString("log_level",
			mcp.Description("Set logging level"),
			mcp.Enum("trace", "debug", "info", "warn", "error"),
		),
		mcp.WithBoolean("no_migrate",
			mcp.Description("Disable auto-migration before sync"),
		),
	)
	s.AddTool(syncTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configPath := request.GetString("config_path", "")
		args := []string{"cloudquery", "sync", configPath}
		
		if logLevel := request.GetString("log_level", ""); logLevel != "" {
			args = append(args, "--log-level", logLevel)
		}
		if request.GetBool("no_migrate", false) {
			args = append(args, "--no-migrate")
		}
		
		return executeShipCommand(args)
	})

	// CloudQuery migrate destination schema
	migrateTool := mcp.NewTool("cloudquery_migrate",
		mcp.WithDescription("Update destination schema based on source configuration"),
		mcp.WithString("config_path",
			mcp.Description("Path to configuration file or directory"),
			mcp.Required(),
		),
		mcp.WithString("log_level",
			mcp.Description("Set logging level"),
			mcp.Enum("trace", "debug", "info", "warn", "error"),
		),
	)
	s.AddTool(migrateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configPath := request.GetString("config_path", "")
		args := []string{"cloudquery", "migrate", configPath}
		
		if logLevel := request.GetString("log_level", ""); logLevel != "" {
			args = append(args, "--log-level", logLevel)
		}
		
		return executeShipCommand(args)
	})

	// CloudQuery init - generate initial configuration
	initTool := mcp.NewTool("cloudquery_init",
		mcp.WithDescription("Generate initial configuration file"),
		mcp.WithString("source",
			mcp.Description("Source plugin name (e.g., aws, gcp, azure)"),
		),
		mcp.WithString("destination",
			mcp.Description("Destination plugin name (e.g., postgresql, sqlite)"),
		),
	)
	s.AddTool(initTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloudquery", "init"}
		
		if source := request.GetString("source", ""); source != "" {
			args = append(args, "--source", source)
		}
		if destination := request.GetString("destination", ""); destination != "" {
			args = append(args, "--destination", destination)
		}
		
		return executeShipCommand(args)
	})

	// CloudQuery validate configuration
	validateConfigTool := mcp.NewTool("cloudquery_validate_config",
		mcp.WithDescription("Validate CloudQuery configuration files"),
		mcp.WithString("config_path",
			mcp.Description("Path to configuration file or directory"),
			mcp.Required(),
		),
	)
	s.AddTool(validateConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configPath := request.GetString("config_path", "")
		args := []string{"cloudquery", "validate-config", configPath}
		return executeShipCommand(args)
	})

	// CloudQuery test connection
	testConnectionTool := mcp.NewTool("cloudquery_test_connection",
		mcp.WithDescription("Test plugin connections without running full sync"),
		mcp.WithString("config_path",
			mcp.Description("Path to configuration file or directory"),
			mcp.Required(),
		),
	)
	s.AddTool(testConnectionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configPath := request.GetString("config_path", "")
		args := []string{"cloudquery", "test-connection", configPath}
		return executeShipCommand(args)
	})

	// CloudQuery tables - generate table documentation
	tablesTool := mcp.NewTool("cloudquery_tables",
		mcp.WithDescription("Generate documentation for supported tables"),
		mcp.WithString("source",
			mcp.Description("Source plugin name to generate tables for"),
		),
		mcp.WithString("output_dir",
			mcp.Description("Output directory for table documentation"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "markdown"),
		),
	)
	s.AddTool(tablesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloudquery", "tables"}
		
		if source := request.GetString("source", ""); source != "" {
			args = append(args, source)
		}
		if outputDir := request.GetString("output_dir", ""); outputDir != "" {
			args = append(args, "--output-dir", outputDir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		
		return executeShipCommand(args)
	})

	// CloudQuery login
	loginTool := mcp.NewTool("cloudquery_login",
		mcp.WithDescription("Login to CloudQuery Hub"),
	)
	s.AddTool(loginTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloudquery", "login"}
		return executeShipCommand(args)
	})

	// CloudQuery logout
	logoutTool := mcp.NewTool("cloudquery_logout",
		mcp.WithDescription("Logout from CloudQuery Hub"),
	)
	s.AddTool(logoutTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloudquery", "logout"}
		return executeShipCommand(args)
	})

	// CloudQuery plugin install
	pluginInstallTool := mcp.NewTool("cloudquery_plugin_install",
		mcp.WithDescription("Install CloudQuery plugin"),
		mcp.WithString("plugin_name",
			mcp.Description("Plugin name to install"),
			mcp.Required(),
		),
	)
	s.AddTool(pluginInstallTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pluginName := request.GetString("plugin_name", "")
		args := []string{"cloudquery", "plugin", "install", pluginName}
		return executeShipCommand(args)
	})

	// CloudQuery switch
	switchTool := mcp.NewTool("cloudquery_switch",
		mcp.WithDescription("Switch between CloudQuery contexts or configurations"),
	)
	s.AddTool(switchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloudquery", "switch"}
		return executeShipCommand(args)
	})
}