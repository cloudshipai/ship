package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddPMapperTools adds PMapper (AWS IAM privilege escalation analysis) MCP tool implementations
func AddPMapperTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// PMapper create graph tool
	createGraphTool := mcp.NewTool("pmapper_create_graph",
		mcp.WithDescription("Create IAM privilege graph using PMapper"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
			mcp.Required(),
		),
	)
	s.AddTool(createGraphTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profile := request.GetString("profile", "")
		args := []string{"security", "pmapper", "--create-graph", "--profile", profile}
		return executeShipCommand(args)
	})

	// PMapper query access tool
	queryAccessTool := mcp.NewTool("pmapper_query_access",
		mcp.WithDescription("Query IAM access permissions using PMapper"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
			mcp.Required(),
		),
		mcp.WithString("principal",
			mcp.Description("IAM principal (user/role) to query"),
			mcp.Required(),
		),
		mcp.WithString("action",
			mcp.Description("AWS action to check (e.g., s3:GetObject)"),
			mcp.Required(),
		),
		mcp.WithString("resource",
			mcp.Description("AWS resource ARN to check"),
			mcp.Required(),
		),
	)
	s.AddTool(queryAccessTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profile := request.GetString("profile", "")
		principal := request.GetString("principal", "")
		action := request.GetString("action", "")
		resource := request.GetString("resource", "")
		args := []string{"security", "pmapper", "--query", "--profile", profile, "--principal", principal, "--action", action, "--resource", resource}
		return executeShipCommand(args)
	})

	// PMapper find privilege escalation tool
	findPrivEscTool := mcp.NewTool("pmapper_find_privilege_escalation",
		mcp.WithDescription("Find privilege escalation paths using PMapper"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
			mcp.Required(),
		),
		mcp.WithString("principal",
			mcp.Description("IAM principal to analyze for privilege escalation"),
			mcp.Required(),
		),
	)
	s.AddTool(findPrivEscTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profile := request.GetString("profile", "")
		principal := request.GetString("principal", "")
		args := []string{"security", "pmapper", "--privesc", "--profile", profile, "--principal", principal}
		return executeShipCommand(args)
	})

	// PMapper visualize graph tool
	visualizeGraphTool := mcp.NewTool("pmapper_visualize_graph",
		mcp.WithDescription("Visualize IAM privilege graph using PMapper"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (svg, png, dot)"),
			mcp.Enum("svg", "png", "dot"),
		),
	)
	s.AddTool(visualizeGraphTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profile := request.GetString("profile", "")
		args := []string{"security", "pmapper", "--visualize", "--profile", profile}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		return executeShipCommand(args)
	})

	// PMapper list principals tool
	listPrincipalsTool := mcp.NewTool("pmapper_list_principals",
		mcp.WithDescription("List IAM principals using PMapper"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
			mcp.Required(),
		),
	)
	s.AddTool(listPrincipalsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profile := request.GetString("profile", "")
		args := []string{"security", "pmapper", "--list-principals", "--profile", profile}
		return executeShipCommand(args)
	})

	// PMapper check admin access tool
	checkAdminAccessTool := mcp.NewTool("pmapper_check_admin_access",
		mcp.WithDescription("Check if principal has admin access using PMapper"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
			mcp.Required(),
		),
		mcp.WithString("principal",
			mcp.Description("IAM principal to check for admin access"),
			mcp.Required(),
		),
	)
	s.AddTool(checkAdminAccessTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		profile := request.GetString("profile", "")
		principal := request.GetString("principal", "")
		args := []string{"security", "pmapper", "--check-admin", "--profile", profile, "--principal", principal}
		return executeShipCommand(args)
	})
}