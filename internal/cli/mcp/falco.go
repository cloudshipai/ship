package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddFalcoTools adds Falco (runtime security monitoring) MCP tool implementations
func AddFalcoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addFalcoToolsDirect(s)
}

// addFalcoToolsDirect implements direct Dagger calls for Falco tools
func addFalcoToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		configPath := request.GetString("config_path", "")
		rulesPath := request.GetString("rules_path", "")
		outputFormat := request.GetString("output_format", "")

		// Create Falco module and start monitoring
		falcoModule := modules.NewFalcoModule(client)
		result, err := falcoModule.StartMonitoring(ctx, configPath, rulesPath, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("falco start monitoring failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		rulesPath := request.GetString("rules_path", "")

		// Create Falco module and validate rules
		falcoModule := modules.NewFalcoModule(client)
		result, err := falcoModule.ValidateRulesSimple(ctx, rulesPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("falco validate rules failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		configPath := request.GetString("config_path", "")
		rulesPath := request.GetString("rules_path", "")

		// Create Falco module and run dry run
		falcoModule := modules.NewFalcoModule(client)
		result, err := falcoModule.DryRunSimple(ctx, configPath, rulesPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("falco dry run failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		source := request.GetString("source", "")

		// Create Falco module and list fields
		falcoModule := modules.NewFalcoModule(client)
		result, err := falcoModule.ListFieldsWithSource(ctx, source)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("falco list fields failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Falco get version tool
	getVersionTool := mcp.NewTool("falco_get_version",
		mcp.WithDescription("Get Falco version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create Falco module and get version
		falcoModule := modules.NewFalcoModule(client)
		result, err := falcoModule.GetVersionSimple(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("falco get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Falco list rules tool
	listRulesTool := mcp.NewTool("falco_list_rules",
		mcp.WithDescription("List all loaded Falco rules"),
		mcp.WithString("rules_path",
			mcp.Description("Path to custom Falco rules"),
		),
	)
	s.AddTool(listRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		rulesPath := request.GetString("rules_path", "")

		// Create Falco module and list rules
		falcoModule := modules.NewFalcoModule(client)
		result, err := falcoModule.ListRulesSimple(ctx, rulesPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("falco list rules failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		ruleName := request.GetString("rule_name", "")
		rulesPath := request.GetString("rules_path", "")

		// Create Falco module and describe rule
		falcoModule := modules.NewFalcoModule(client)
		result, err := falcoModule.DescribeRuleSimple(ctx, ruleName, rulesPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("falco describe rule failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}