package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddSteampipeTools adds Steampipe (cloud asset querying) MCP tool implementations using direct Dagger calls
// NOTE: Steampipe is typically configured as an external MCP server via npx @turbot/steampipe-mcp
// These tools provide Dagger-based execution as an alternative
func AddSteampipeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addSteampipeToolsDirect(s)
}

// addSteampipeToolsDirect adds Steampipe tools using direct Dagger module calls
func addSteampipeToolsDirect(s *server.MCPServer) {
	// Steampipe query tool
	queryTool := mcp.NewTool("steampipe_query",
		mcp.WithDescription("Execute SQL query against cloud resources using real Steampipe CLI"),
		mcp.WithString("query",
			mcp.Description("SQL query to execute"),
			mcp.Required(),
		),
		mcp.WithString("plugin",
			mcp.Description("Plugin to install and use for the query"),
			mcp.Required(),
		),
	)
	s.AddTool(queryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// Get parameters
		query := request.GetString("query", "")
		plugin := request.GetString("plugin", "")

		// Execute query
		output, err := module.Query(ctx, query, plugin)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe query failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Steampipe interactive query tool
	interactiveQueryTool := mcp.NewTool("steampipe_query_interactive",
		mcp.WithDescription("Start interactive SQL query session (limited in containers) using real Steampipe CLI"),
		mcp.WithString("plugin",
			mcp.Description("Plugin to install for the interactive session"),
			mcp.Required(),
		),
	)
	s.AddTool(interactiveQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// Get parameters
		plugin := request.GetString("plugin", "")

		// Start interactive query
		output, err := module.QueryInteractive(ctx, plugin)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe interactive query failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Steampipe plugin list tool
	pluginListTool := mcp.NewTool("steampipe_plugin_list",
		mcp.WithDescription("List installed and available plugins using real Steampipe CLI"),
	)
	s.AddTool(pluginListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// List plugins
		output, err := module.ListPlugins(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe plugin list failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Steampipe plugin install tool
	pluginInstallTool := mcp.NewTool("steampipe_plugin_install",
		mcp.WithDescription("Install plugin using real Steampipe CLI"),
		mcp.WithString("plugin",
			mcp.Description("Plugin name to install (e.g., aws, azure, gcp)"),
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

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// Get parameters
		plugin := request.GetString("plugin", "")

		// Install plugin
		output, err := module.InstallPlugin(ctx, plugin)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe plugin install failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Steampipe plugin update tool
	pluginUpdateTool := mcp.NewTool("steampipe_plugin_update",
		mcp.WithDescription("Update plugin using real Steampipe CLI"),
		mcp.WithString("plugin",
			mcp.Description("Plugin name to update (e.g., aws, azure, gcp)"),
			mcp.Required(),
		),
	)
	s.AddTool(pluginUpdateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// Get parameters
		plugin := request.GetString("plugin", "")

		// Update plugin
		output, err := module.UpdatePlugin(ctx, plugin)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe plugin update failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Steampipe plugin uninstall tool
	pluginUninstallTool := mcp.NewTool("steampipe_plugin_uninstall",
		mcp.WithDescription("Uninstall plugin using real Steampipe CLI"),
		mcp.WithString("plugin",
			mcp.Description("Plugin name to uninstall (e.g., aws, azure, gcp)"),
			mcp.Required(),
		),
	)
	s.AddTool(pluginUninstallTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// Get parameters
		plugin := request.GetString("plugin", "")

		// Uninstall plugin
		output, err := module.UninstallPlugin(ctx, plugin)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe plugin uninstall failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Steampipe service start tool
	serviceStartTool := mcp.NewTool("steampipe_service_start",
		mcp.WithDescription("Start Steampipe service using real Steampipe CLI"),
		mcp.WithNumber("port",
			mcp.Description("Port to run the service on (default: 9193)"),
		),
	)
	s.AddTool(serviceStartTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// Get parameters
		port := int(request.GetFloat("port", 9193))

		// Start service
		output, err := module.StartService(ctx, port)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe service start failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Steampipe service status tool
	serviceStatusTool := mcp.NewTool("steampipe_service_status",
		mcp.WithDescription("Check Steampipe service status using real Steampipe CLI"),
	)
	s.AddTool(serviceStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// Get service status
		output, err := module.GetServiceStatus(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe service status failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Steampipe service stop tool
	serviceStopTool := mcp.NewTool("steampipe_service_stop",
		mcp.WithDescription("Stop Steampipe service using real Steampipe CLI"),
	)
	s.AddTool(serviceStopTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// Stop service
		output, err := module.StopService(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe service stop failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Steampipe version tool
	versionTool := mcp.NewTool("steampipe_version",
		mcp.WithDescription("Get Steampipe version information using real Steampipe CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSteampipeModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Steampipe get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}