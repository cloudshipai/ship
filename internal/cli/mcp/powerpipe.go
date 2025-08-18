package mcp

import (
	"context"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddPowerpipeTools adds Powerpipe MCP tool implementations using real powerpipe CLI commands
func AddPowerpipeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Powerpipe benchmark run tool
	benchmarkRunTool := mcp.NewTool("powerpipe_benchmark_run",
		mcp.WithDescription("Run security and compliance benchmarks using real powerpipe CLI"),
		mcp.WithString("benchmark",
			mcp.Description("Benchmark to run"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("pretty", "plain", "yaml", "json"),
		),
		mcp.WithString("mod_location",
			mcp.Description("Workspace working directory"),
		),
		mcp.WithString("workspace",
			mcp.Description("Powerpipe workspace profile"),
		),
	)
	s.AddTool(benchmarkRunTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		benchmark := request.GetString("benchmark", "")
		args := []string{"powerpipe", "benchmark", "run", benchmark}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if modLocation := request.GetString("mod_location", ""); modLocation != "" {
			args = append(args, "--mod-location", modLocation)
		}
		if workspace := request.GetString("workspace", ""); workspace != "" {
			args = append(args, "--workspace", workspace)
		}
		
		return executeShipCommand(args)
	})

	// Powerpipe benchmark list tool
	benchmarkListTool := mcp.NewTool("powerpipe_benchmark_list",
		mcp.WithDescription("List available benchmarks using real powerpipe CLI"),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("pretty", "plain", "yaml", "json"),
		),
		mcp.WithString("mod_location",
			mcp.Description("Workspace working directory"),
		),
	)
	s.AddTool(benchmarkListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"powerpipe", "benchmark", "list"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if modLocation := request.GetString("mod_location", ""); modLocation != "" {
			args = append(args, "--mod-location", modLocation)
		}
		
		return executeShipCommand(args)
	})

	// Powerpipe query run tool
	queryRunTool := mcp.NewTool("powerpipe_query_run",
		mcp.WithDescription("Execute SQL queries against cloud infrastructure using real powerpipe CLI"),
		mcp.WithString("query",
			mcp.Description("Query to run"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("pretty", "plain", "yaml", "json"),
		),
		mcp.WithString("mod_location",
			mcp.Description("Workspace working directory"),
		),
	)
	s.AddTool(queryRunTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.GetString("query", "")
		args := []string{"powerpipe", "query", "run", query}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if modLocation := request.GetString("mod_location", ""); modLocation != "" {
			args = append(args, "--mod-location", modLocation)
		}
		
		return executeShipCommand(args)
	})

	// Powerpipe query list tool
	queryListTool := mcp.NewTool("powerpipe_query_list",
		mcp.WithDescription("List available queries using real powerpipe CLI"),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("pretty", "plain", "yaml", "json"),
		),
		mcp.WithString("mod_location",
			mcp.Description("Workspace working directory"),
		),
	)
	s.AddTool(queryListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"powerpipe", "query", "list"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if modLocation := request.GetString("mod_location", ""); modLocation != "" {
			args = append(args, "--mod-location", modLocation)
		}
		
		return executeShipCommand(args)
	})

	// Powerpipe server start tool
	serverTool := mcp.NewTool("powerpipe_server",
		mcp.WithDescription("Start Powerpipe server using real powerpipe CLI"),
		mcp.WithString("listen",
			mcp.Description("Accept connections from specified address (default local)"),
			mcp.Enum("local", "network"),
		),
		mcp.WithNumber("port",
			mcp.Description("Port for Powerpipe server"),
		),
		mcp.WithString("mod_location",
			mcp.Description("Workspace working directory"),
		),
		mcp.WithString("workspace",
			mcp.Description("Powerpipe workspace profile"),
		),
	)
	s.AddTool(serverTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"powerpipe", "server"}
		
		if listen := request.GetString("listen", ""); listen != "" {
			args = append(args, "--listen", listen)
		}
		if port := request.GetInt("port", 0); port > 0 {
			args = append(args, "--port", strconv.Itoa(port))
		}
		if modLocation := request.GetString("mod_location", ""); modLocation != "" {
			args = append(args, "--mod-location", modLocation)
		}
		if workspace := request.GetString("workspace", ""); workspace != "" {
			args = append(args, "--workspace", workspace)
		}
		
		return executeShipCommand(args)
	})

	// Powerpipe dashboard list tool
	dashboardListTool := mcp.NewTool("powerpipe_dashboard_list",
		mcp.WithDescription("List available dashboards using real powerpipe CLI"),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("pretty", "plain", "yaml", "json"),
		),
		mcp.WithString("mod_location",
			mcp.Description("Workspace working directory"),
		),
	)
	s.AddTool(dashboardListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"powerpipe", "dashboard", "list"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if modLocation := request.GetString("mod_location", ""); modLocation != "" {
			args = append(args, "--mod-location", modLocation)
		}
		
		return executeShipCommand(args)
	})

	// Powerpipe version tool
	versionTool := mcp.NewTool("powerpipe_version",
		mcp.WithDescription("Get Powerpipe version information using real powerpipe CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"powerpipe", "--version"}
		return executeShipCommand(args)
	})
}