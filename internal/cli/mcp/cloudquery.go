package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCloudQueryTools adds CloudQuery (cloud asset inventory) MCP tool implementations
func AddCloudQueryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addCloudQueryToolsDirect(s)
}

// addCloudQueryToolsDirect implements direct Dagger calls for CloudQuery tools
func addCloudQueryToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		configPath := request.GetString("config_path", "")
		logLevel := request.GetString("log_level", "")
		noMigrate := request.GetBool("no_migrate", false)

		// Create CloudQuery module and sync
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.SyncWithOptions(ctx, configPath, logLevel, noMigrate)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery sync failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		configPath := request.GetString("config_path", "")
		logLevel := request.GetString("log_level", "")

		// Create CloudQuery module and migrate
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.MigrateWithOptions(ctx, configPath, logLevel)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery migrate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		source := request.GetString("source", "")
		destination := request.GetString("destination", "")

		// Create CloudQuery module and init
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.InitConfig(ctx, source, destination)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery init failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		configPath := request.GetString("config_path", "")

		// Create CloudQuery module and validate config
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.ValidateConfig(ctx, configPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery validate config failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		configPath := request.GetString("config_path", "")

		// Create CloudQuery module and test connection
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.TestConnection(ctx, configPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery test connection failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		source := request.GetString("source", "")
		outputDir := request.GetString("output_dir", "")
		format := request.GetString("format", "")

		// Create CloudQuery module and get tables
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.GetTables(ctx, source, outputDir, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery get tables failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// CloudQuery login
	loginTool := mcp.NewTool("cloudquery_login",
		mcp.WithDescription("Login to CloudQuery Hub"),
	)
	s.AddTool(loginTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create CloudQuery module and login
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.Login(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery login failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// CloudQuery logout
	logoutTool := mcp.NewTool("cloudquery_logout",
		mcp.WithDescription("Logout from CloudQuery Hub"),
	)
	s.AddTool(logoutTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create CloudQuery module and logout
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.Logout(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery logout failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		pluginName := request.GetString("plugin_name", "")

		// Create CloudQuery module and install plugin
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.InstallPlugin(ctx, pluginName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery install plugin failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// CloudQuery switch
	switchTool := mcp.NewTool("cloudquery_switch",
		mcp.WithDescription("Switch between CloudQuery contexts or configurations"),
	)
	s.AddTool(switchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create CloudQuery module and switch
		cloudQueryModule := modules.NewCloudQueryModule(client)
		result, err := cloudQueryModule.Switch(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudquery switch failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}