package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddPowerpipeTools adds Powerpipe (infrastructure benchmarking) MCP tool implementations
func AddPowerpipeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Powerpipe benchmark tool
	benchmarkTool := mcp.NewTool("powerpipe_benchmark",
		mcp.WithDescription("Run security and compliance benchmarks using Powerpipe"),
		mcp.WithString("benchmark",
			mcp.Description("Benchmark to run (e.g., aws_compliance, kubernetes_compliance)"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table", "csv"),
		),
	)
	s.AddTool(benchmarkTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		benchmark := request.GetString("benchmark", "")
		args := []string{"security", "powerpipe", "benchmark", benchmark}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Powerpipe query tool
	queryTool := mcp.NewTool("powerpipe_query",
		mcp.WithDescription("Execute SQL queries against cloud infrastructure"),
		mcp.WithString("query",
			mcp.Description("SQL query to execute"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table", "csv"),
		),
	)
	s.AddTool(queryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.GetString("query", "")
		args := []string{"security", "powerpipe", "query", query}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Powerpipe server tool
	serverTool := mcp.NewTool("powerpipe_server",
		mcp.WithDescription("Start Powerpipe server for interactive analysis"),
		mcp.WithNumber("port",
			mcp.Description("Port for Powerpipe server"),
		),
		mcp.WithString("host",
			mcp.Description("Host for Powerpipe server"),
		),
	)
	s.AddTool(serverTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "powerpipe", "server"}
		if port := request.GetInt("port", 0); port > 0 {
			args = append(args, "--port", string(rune(port)))
		}
		if host := request.GetString("host", ""); host != "" {
			args = append(args, "--host", host)
		}
		return executeShipCommand(args)
	})

	// Powerpipe get version tool
	getVersionTool := mcp.NewTool("powerpipe_get_version",
		mcp.WithDescription("Get Powerpipe version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "powerpipe", "--version"}
		return executeShipCommand(args)
	})
}