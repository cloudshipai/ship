package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddPowerpipeTools adds Powerpipe MCP tool implementations using direct Dagger calls
func AddPowerpipeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addPowerpipeToolsDirect(s)
}

// addPowerpipeToolsDirect adds Powerpipe tools using direct Dagger module calls
func addPowerpipeToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPowerpipeModule(client)

		// Get parameters
		benchmark := request.GetString("benchmark", "")
		if benchmark == "" {
			return mcp.NewToolResultError("benchmark is required"), nil
		}
		modLocation := request.GetString("mod_location", "")

		// Note: Dagger module doesn't support output format and workspace options
		if request.GetString("output", "") != "" || request.GetString("workspace", "") != "" {
			return mcp.NewToolResultError("output format and workspace options not supported in Dagger module"), nil
		}

		// Run benchmark
		output, err := module.RunBenchmark(ctx, benchmark, modLocation)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("powerpipe benchmark run failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPowerpipeModule(client)

		// Get parameters
		modLocation := request.GetString("mod_location", "")

		// Note: Dagger module doesn't support output format option
		if request.GetString("output", "") != "" {
			return mcp.NewToolResultError("output format option not supported in Dagger module"), nil
		}

		// List benchmarks
		output, err := module.ListBenchmarks(ctx, modLocation)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("powerpipe benchmark list failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPowerpipeModule(client)

		// Get parameters
		query := request.GetString("query", "")
		if query == "" {
			return mcp.NewToolResultError("query is required"), nil
		}
		modLocation := request.GetString("mod_location", "")

		// Note: Dagger module doesn't support output format option
		if request.GetString("output", "") != "" {
			return mcp.NewToolResultError("output format option not supported in Dagger module"), nil
		}

		// Run query
		output, err := module.RunQuery(ctx, query, modLocation)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("powerpipe query run failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPowerpipeModule(client)

		// Get parameters
		modLocation := request.GetString("mod_location", "")

		// Note: Dagger module doesn't support output format option
		if request.GetString("output", "") != "" {
			return mcp.NewToolResultError("output format option not supported in Dagger module"), nil
		}

		// List queries
		output, err := module.ListQueries(ctx, modLocation)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("powerpipe query list failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPowerpipeModule(client)

		// Get parameters
		port := request.GetInt("port", 0)
		modLocation := request.GetString("mod_location", "")

		// Note: Dagger module doesn't support listen and workspace options
		if request.GetString("listen", "") != "" || request.GetString("workspace", "") != "" {
			return mcp.NewToolResultError("listen and workspace options not supported in Dagger module"), nil
		}

		// Start server
		output, err := module.StartServer(ctx, modLocation, port)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("powerpipe server failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPowerpipeModule(client)

		// Get parameters
		modLocation := request.GetString("mod_location", "")

		// Note: Dagger module doesn't support output format option
		if request.GetString("output", "") != "" {
			return mcp.NewToolResultError("output format option not supported in Dagger module"), nil
		}

		// List dashboards
		output, err := module.ListDashboards(ctx, modLocation)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("powerpipe dashboard list failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Powerpipe version tool
	versionTool := mcp.NewTool("powerpipe_version",
		mcp.WithDescription("Get Powerpipe version information using real powerpipe CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPowerpipeModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("powerpipe version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}