package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddPMapperTools adds PMapper (AWS IAM privilege escalation analysis) MCP tool implementations using real CLI commands
func AddPMapperTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		args := []string{"pmapper"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if account := request.GetString("account", ""); account != "" {
			args = append(args, "--account", account)
		}
		
		args = append(args, "graph", "create")
		return executeShipCommand(args)
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
		queryString := request.GetString("query_string", "")
		args := []string{"pmapper"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if account := request.GetString("account", ""); account != "" {
			args = append(args, "--account", account)
		}
		
		args = append(args, "query", queryString)
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		args := []string{"pmapper"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if account := request.GetString("account", ""); account != "" {
			args = append(args, "--account", account)
		}
		
		queryString := "preset privesc " + target
		args = append(args, "query", queryString)
		return executeShipCommand(args)
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
		args := []string{"pmapper"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if account := request.GetString("account", ""); account != "" {
			args = append(args, "--account", account)
		}
		
		args = append(args, "visualize")
		
		if filetype := request.GetString("filetype", ""); filetype != "" {
			args = append(args, "--filetype", filetype)
		}
		
		return executeShipCommand(args)
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
		action := request.GetString("action", "")
		args := []string{"pmapper"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if account := request.GetString("account", ""); account != "" {
			args = append(args, "--account", account)
		}
		
		queryString := "who can do " + action
		args = append(args, "query", queryString)
		return executeShipCommand(args)
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
		action := request.GetString("action", "")
		args := []string{"pmapper"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if account := request.GetString("account", ""); account != "" {
			args = append(args, "--account", account)
		}
		
		args = append(args, "argquery")
		
		if request.GetBool("skip_admin", false) {
			args = append(args, "-s")
		}
		
		args = append(args, "--action", action)
		
		if condition := request.GetString("condition", ""); condition != "" {
			args = append(args, "--condition", condition)
		}
		
		return executeShipCommand(args)
	})




}