package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddPMapperTools adds PMapper (AWS IAM privilege escalation analysis) MCP tool implementations using direct Dagger calls
func AddPMapperTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addPMapperToolsDirect(s)
}

// addPMapperToolsDirect adds PMapper tools using direct Dagger module calls
func addPMapperToolsDirect(s *server.MCPServer) {
	// PMapper create graph tool
	createGraphTool := mcp.NewTool("pmapper_graph_create",
		mcp.WithDescription("Create IAM privilege graph using real pmapper CLI"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("account",
			mcp.Description("AWS account number"),
		),
	)
	s.AddTool(createGraphTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPMapperModule(client)

		// Get parameters
		profile := request.GetString("profile", "")

		// Note: Dagger module doesn't support account parameter
		if request.GetString("account", "") != "" {
			return mcp.NewToolResultError("account parameter not supported in Dagger module"), nil
		}

		// Create graph
		output, err := module.CreateGraph(ctx, profile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("pmapper create graph failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// PMapper query tool
	queryTool := mcp.NewTool("pmapper_query",
		mcp.WithDescription("Query IAM permissions using real pmapper CLI"),
		mcp.WithString("query_string",
			mcp.Description("Query string (e.g., 'who can do iam:CreateUser')"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("account",
			mcp.Description("AWS account number"),
		),
	)
	s.AddTool(queryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		queryString := request.GetString("query_string", "")
		if queryString == "" {
			return mcp.NewToolResultError("query_string is required"), nil
		}

		// Note: Generic query not directly supported in Dagger module
		// Dagger module has specific methods for different query types
		return mcp.NewToolResultError("generic query not supported in Dagger module - use specific query functions instead"), nil
	})

	// PMapper privilege escalation query tool
	privEscTool := mcp.NewTool("pmapper_query_privesc",
		mcp.WithDescription("Find privilege escalation paths using real pmapper CLI preset query"),
		mcp.WithString("target",
			mcp.Description("Target principal or wildcard (*) for all principals"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("account",
			mcp.Description("AWS account number"),
		),
	)
	s.AddTool(privEscTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPMapperModule(client)

		// Get parameters
		target := request.GetString("target", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		profile := request.GetString("profile", "")

		// Note: Dagger module doesn't support account parameter
		if request.GetString("account", "") != "" {
			return mcp.NewToolResultError("account parameter not supported in Dagger module"), nil
		}

		// Find privilege escalation
		output, err := module.FindPrivilegeEscalation(ctx, profile, target)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("pmapper privilege escalation query failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// PMapper visualize tool
	visualizeTool := mcp.NewTool("pmapper_visualize",
		mcp.WithDescription("Visualize IAM privilege graph using real pmapper CLI"),
		mcp.WithString("filetype",
			mcp.Description("Output file type (svg, png, etc.)"),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("account",
			mcp.Description("AWS account number"),
		),
	)
	s.AddTool(visualizeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPMapperModule(client)

		// Get parameters
		profile := request.GetString("profile", "")
		filetype := request.GetString("filetype", "")

		// Note: Dagger module doesn't support account parameter
		if request.GetString("account", "") != "" {
			return mcp.NewToolResultError("account parameter not supported in Dagger module"), nil
		}

		// Visualize graph
		output, err := module.VisualizeGraph(ctx, profile, filetype)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("pmapper visualize failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// PMapper query who can do action tool
	queryWhoCanTool := mcp.NewTool("pmapper_query_who_can",
		mcp.WithDescription("Query who can perform specific action using real pmapper CLI"),
		mcp.WithString("action",
			mcp.Description("AWS action to check (e.g., iam:CreateUser)"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("account",
			mcp.Description("AWS account number"),
		),
	)
	s.AddTool(queryWhoCanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPMapperModule(client)

		// Get parameters
		action := request.GetString("action", "")
		if action == "" {
			return mcp.NewToolResultError("action is required"), nil
		}
		profile := request.GetString("profile", "")

		// Note: "who can do" query not directly supported in Dagger module
		// Use QueryAccess with wildcard principal instead
		output, err := module.QueryAccess(ctx, profile, "*", action, "")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("pmapper who can query failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// PMapper advanced query tool
	argqueryTool := mcp.NewTool("pmapper_argquery",
		mcp.WithDescription("Advanced query with conditions using real pmapper CLI"),
		mcp.WithString("action",
			mcp.Description("AWS action to check (e.g., ec2:RunInstances)"),
			mcp.Required(),
		),
		mcp.WithString("condition",
			mcp.Description("Condition to check (e.g., ec2:InstanceType=c6gd.16xlarge)"),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("account",
			mcp.Description("AWS account number"),
		),
		mcp.WithBoolean("skip_admin",
			mcp.Description("Skip reporting current admin users (-s flag)"),
		),
	)
	s.AddTool(argqueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPMapperModule(client)

		// Get parameters
		action := request.GetString("action", "")
		if action == "" {
			return mcp.NewToolResultError("action is required"), nil
		}
		profile := request.GetString("profile", "")
		condition := request.GetString("condition", "")

		// Note: Dagger module doesn't support advanced argquery with conditions and skip_admin
		if request.GetString("account", "") != "" || request.GetBool("skip_admin", false) || condition != "" {
			return mcp.NewToolResultError("advanced argquery options not supported in Dagger module"), nil
		}

		// Use QueryAccess as a fallback for basic action queries
		output, err := module.QueryAccess(ctx, profile, "*", action, "")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("pmapper argquery failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})




}