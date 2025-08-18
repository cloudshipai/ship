package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSteampipeTools adds Steampipe (cloud asset querying) MCP tool implementations using real steampipe CLI commands
// NOTE: Steampipe is typically configured as an external MCP server via npx @turbot/steampipe-mcp
// These tools provide Dagger-based execution as an alternative
func AddSteampipeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Steampipe query tool
	queryTool := mcp.NewTool("steampipe_query",
		mcp.WithDescription("Execute SQL query against cloud resources using real steampipe CLI"),
		mcp.WithString("query",
			mcp.Description("SQL query to execute"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("line", "csv", "json", "table", "snapshot"),
		),
		mcp.WithString("export",
			mcp.Description("Export query output to a file"),
		),
		mcp.WithBoolean("header",
			mcp.Description("Include column headers in output"),
		),
		mcp.WithString("timing",
			mcp.Description("Show query execution timing"),
			mcp.Enum("off", "on", "verbose"),
		),
		mcp.WithString("search_path",
			mcp.Description("Set custom search path for connections"),
		),
		mcp.WithNumber("query_timeout",
			mcp.Description("Query timeout in seconds"),
		),
	)
	s.AddTool(queryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.GetString("query", "")
		args := []string{"steampipe", "query", query}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if export := request.GetString("export", ""); export != "" {
			args = append(args, "--export", export)
		}
		if request.GetBool("header", true) {
			args = append(args, "--header=true")
		} else {
			args = append(args, "--header=false")
		}
		if timing := request.GetString("timing", ""); timing != "" {
			args = append(args, "--timing", timing)
		}
		if searchPath := request.GetString("search_path", ""); searchPath != "" {
			args = append(args, "--search-path", searchPath)
		}
		if timeout := request.GetInt("query_timeout", 0); timeout > 0 {
			args = append(args, "--query-timeout", fmt.Sprintf("%d", timeout))
		}
		
		return executeShipCommand(args)
	})

	// Steampipe interactive query tool
	interactiveQueryTool := mcp.NewTool("steampipe_query_interactive",
		mcp.WithDescription("Start interactive SQL query shell using real steampipe CLI"),
		mcp.WithString("workspace",
			mcp.Description("Steampipe workspace profile to use"),
		),
	)
	s.AddTool(interactiveQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"steampipe", "query"}
		
		if workspace := request.GetString("workspace", ""); workspace != "" {
			args = append(args, "--workspace", workspace)
		}
		
		return executeShipCommand(args)
	})

	// Steampipe plugin list tool
	pluginListTool := mcp.NewTool("steampipe_plugin_list",
		mcp.WithDescription("List installed Steampipe plugins using real steampipe CLI"),
	)
	s.AddTool(pluginListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"steampipe", "plugin", "list"}
		return executeShipCommand(args)
	})

	// Steampipe plugin install tool
	pluginInstallTool := mcp.NewTool("steampipe_plugin_install",
		mcp.WithDescription("Install Steampipe plugins using real steampipe CLI"),
		mcp.WithString("plugin_name",
			mcp.Description("Name of the plugin to install (e.g., aws, azure, gcp)"),
			mcp.Required(),
		),
		mcp.WithString("version",
			mcp.Description("Specific version to install (e.g., @0.107.0)"),
		),
		mcp.WithBoolean("skip_config",
			mcp.Description("Skip creating the default config file"),
		),
	)
	s.AddTool(pluginInstallTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pluginName := request.GetString("plugin_name", "")
		args := []string{"steampipe", "plugin", "install"}
		
		if version := request.GetString("version", ""); version != "" {
			args = append(args, pluginName+"@"+version)
		} else {
			args = append(args, pluginName)
		}
		
		if request.GetBool("skip_config", false) {
			args = append(args, "--skip-config")
		}
		
		return executeShipCommand(args)
	})

	// Steampipe plugin update tool
	pluginUpdateTool := mcp.NewTool("steampipe_plugin_update",
		mcp.WithDescription("Update Steampipe plugins using real steampipe CLI"),
		mcp.WithString("plugin_name",
			mcp.Description("Name of specific plugin to update (leave empty to update all)"),
		),
		mcp.WithBoolean("all",
			mcp.Description("Update all installed plugins"),
		),
	)
	s.AddTool(pluginUpdateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"steampipe", "plugin", "update"}
		
		if request.GetBool("all", false) {
			args = append(args, "--all")
		} else if pluginName := request.GetString("plugin_name", ""); pluginName != "" {
			args = append(args, pluginName)
		}
		
		return executeShipCommand(args)
	})

	// Steampipe plugin uninstall tool
	pluginUninstallTool := mcp.NewTool("steampipe_plugin_uninstall",
		mcp.WithDescription("Uninstall Steampipe plugins using real steampipe CLI"),
		mcp.WithString("plugin_name",
			mcp.Description("Name of the plugin to uninstall"),
			mcp.Required(),
		),
	)
	s.AddTool(pluginUninstallTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pluginName := request.GetString("plugin_name", "")
		args := []string{"steampipe", "plugin", "uninstall", pluginName}
		return executeShipCommand(args)
	})

	// Steampipe service start tool
	serviceStartTool := mcp.NewTool("steampipe_service_start",
		mcp.WithDescription("Start Steampipe database service using real steampipe CLI"),
		mcp.WithString("database_listen",
			mcp.Description("Database connection scope"),
			mcp.Enum("local", "network"),
		),
		mcp.WithNumber("database_port",
			mcp.Description("Database service port (default 9193)"),
		),
		mcp.WithString("database_password",
			mcp.Description("Database password for the session"),
		),
		mcp.WithBoolean("foreground",
			mcp.Description("Run service in the foreground"),
		),
		mcp.WithBoolean("show_password",
			mcp.Description("View database connection password"),
		),
	)
	s.AddTool(serviceStartTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"steampipe", "service", "start"}
		
		if listen := request.GetString("database_listen", ""); listen != "" {
			args = append(args, "--database-listen", listen)
		}
		if port := request.GetInt("database_port", 0); port > 0 {
			args = append(args, "--database-port", fmt.Sprintf("%d", port))
		}
		if password := request.GetString("database_password", ""); password != "" {
			args = append(args, "--database-password", password)
		}
		if request.GetBool("foreground", false) {
			args = append(args, "--foreground")
		}
		if request.GetBool("show_password", false) {
			args = append(args, "--show-password")
		}
		
		return executeShipCommand(args)
	})

	// Steampipe service status tool
	serviceStatusTool := mcp.NewTool("steampipe_service_status",
		mcp.WithDescription("Check Steampipe service status using real steampipe CLI"),
		mcp.WithBoolean("all",
			mcp.Description("Show status of all running services"),
		),
	)
	s.AddTool(serviceStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"steampipe", "service", "status"}
		
		if request.GetBool("all", false) {
			args = append(args, "--all")
		}
		
		return executeShipCommand(args)
	})

	// Steampipe service stop tool
	serviceStopTool := mcp.NewTool("steampipe_service_stop",
		mcp.WithDescription("Stop Steampipe service using real steampipe CLI"),
		mcp.WithBoolean("force",
			mcp.Description("Force service shutdown"),
		),
	)
	s.AddTool(serviceStopTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"steampipe", "service", "stop"}
		
		if request.GetBool("force", false) {
			args = append(args, "--force")
		}
		
		return executeShipCommand(args)
	})

	// Steampipe version tool
	versionTool := mcp.NewTool("steampipe_version",
		mcp.WithDescription("Get Steampipe version information using real steampipe CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"steampipe", "--version"}
		return executeShipCommand(args)
	})
}