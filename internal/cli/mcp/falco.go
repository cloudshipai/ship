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
		args := []string{"falco"}
		if configPath := request.GetString("config_path", ""); configPath != "" {
			args = append(args, "-c", configPath)
		}
		if rulesPath := request.GetString("rules_path", ""); rulesPath != "" {
			args = append(args, "-r", rulesPath)
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
		args := []string{"falco", "-V", rulesPath}
		return executeShipCommand(args)
	})

	// Falco dry run tool
	dryRunTool := mcp.NewTool("falco_dry_run",
		mcp.WithDescription("Run Falco in dry-run mode without processing events"),
		mcp.WithString("config_path",
			mcp.Description("Path to Falco configuration file"),
		),
		mcp.WithString("rules_path",
			mcp.Description("Path to custom Falco rules"),
		),
	)
	s.AddTool(dryRunTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"falco", "--dry-run"}
		if configPath := request.GetString("config_path", ""); configPath != "" {
			args = append(args, "-c", configPath)
		}
		if rulesPath := request.GetString("rules_path", ""); rulesPath != "" {
			args = append(args, "-r", rulesPath)
		}
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
		args := []string{"falco", "--list"}
		if source := request.GetString("source", ""); source != "" {
			args = append(args, source)
		}
		return executeShipCommand(args)
	})

	// Falco get version tool
	getVersionTool := mcp.NewTool("falco_get_version",
		mcp.WithDescription("Get Falco version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"falco", "--version"}
		return executeShipCommand(args)
	})

	// Falco list rules tool
	listRulesTool := mcp.NewTool("falco_list_rules",
		mcp.WithDescription("List all loaded Falco rules"),
		mcp.WithString("rules_path",
			mcp.Description("Path to custom Falco rules"),
		),
	)
	s.AddTool(listRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"falco", "-L"}
		if rulesPath := request.GetString("rules_path", ""); rulesPath != "" {
			args = append(args, "-r", rulesPath)
		}
		return executeShipCommand(args)
	})

	// Falco describe rule tool
	describeRuleTool := mcp.NewTool("falco_describe_rule",
		mcp.WithDescription("Show description of a specific Falco rule"),
		mcp.WithString("rule_name",
			mcp.Description("Name of the rule to describe"),
			mcp.Required(),
		),
		mcp.WithString("rules_path",
			mcp.Description("Path to custom Falco rules"),
		),
	)
	s.AddTool(describeRuleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleName := request.GetString("rule_name", "")
		args := []string{"falco", "-l", ruleName}
		if rulesPath := request.GetString("rules_path", ""); rulesPath != "" {
			args = append(args, "-r", rulesPath)
		}
		return executeShipCommand(args)
	})
}