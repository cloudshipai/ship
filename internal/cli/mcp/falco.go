package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddFalcoTools adds Falco (runtime security monitoring) MCP tool implementations
func AddFalcoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Falco start monitoring tool
	startMonitoringTool := mcp.NewTool("falco_start_monitoring",
		mcp.WithDescription("Start Falco runtime security monitoring"),
		mcp.WithString("config_path",
			mcp.Description("Path to Falco configuration file"),
		),
		mcp.WithString("rules_path",
			mcp.Description("Path to custom Falco rules"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text"),
		),
	)
	s.AddTool(startMonitoringTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "falco", "--start"}
		if configPath := request.GetString("config_path", ""); configPath != "" {
			args = append(args, "--config", configPath)
		}
		if rulesPath := request.GetString("rules_path", ""); rulesPath != "" {
			args = append(args, "--rules", rulesPath)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Falco validate rules tool
	validateRulesTool := mcp.NewTool("falco_validate_rules",
		mcp.WithDescription("Validate Falco rules syntax"),
		mcp.WithString("rules_path",
			mcp.Description("Path to Falco rules file to validate"),
			mcp.Required(),
		),
	)
	s.AddTool(validateRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		rulesPath := request.GetString("rules_path", "")
		args := []string{"security", "falco", "--validate", rulesPath}
		return executeShipCommand(args)
	})

	// Falco test rules tool
	testRulesTool := mcp.NewTool("falco_test_rules",
		mcp.WithDescription("Test Falco rules against sample events"),
		mcp.WithString("rules_path",
			mcp.Description("Path to Falco rules file"),
			mcp.Required(),
		),
		mcp.WithString("events_path",
			mcp.Description("Path to test events file"),
			mcp.Required(),
		),
	)
	s.AddTool(testRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		rulesPath := request.GetString("rules_path", "")
		eventsPath := request.GetString("events_path", "")
		args := []string{"security", "falco", "--test", rulesPath, "--events", eventsPath}
		return executeShipCommand(args)
	})

	// Falco list supported fields tool
	listFieldsTool := mcp.NewTool("falco_list_fields",
		mcp.WithDescription("List supported fields for Falco rules"),
		mcp.WithString("source",
			mcp.Description("Event source to list fields for"),
			mcp.Enum("syscall", "k8s_audit", "docker"),
		),
	)
	s.AddTool(listFieldsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "falco", "--list-fields"}
		if source := request.GetString("source", ""); source != "" {
			args = append(args, "--source", source)
		}
		return executeShipCommand(args)
	})

	// Falco get version tool
	getVersionTool := mcp.NewTool("falco_get_version",
		mcp.WithDescription("Get Falco version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "falco", "--version"}
		return executeShipCommand(args)
	})
}