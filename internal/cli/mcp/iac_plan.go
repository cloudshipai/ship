package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddIacPlanTools adds Infrastructure as Code Plan MCP tool implementations
func AddIacPlanTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// IaC generate plan tool
	generatePlanTool := mcp.NewTool("iac_plan_generate",
		mcp.WithDescription("Generate Infrastructure as Code plan"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing IaC files"),
			mcp.Required(),
		),
		mcp.WithString("tool",
			mcp.Description("IaC tool to use"),
			mcp.Required(),
			mcp.Enum("terraform", "terragrunt", "pulumi", "cloudformation"),
		),
		mcp.WithString("var_files",
			mcp.Description("Comma-separated list of variable files"),
		),
		mcp.WithBoolean("destroy",
			mcp.Description("Generate destroy plan"),
		),
	)
	s.AddTool(generatePlanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		tool := request.GetString("tool", "")
		args := []string{"security", "iac-plan", "--generate", "--workdir", workdir, "--tool", tool}
		if varFiles := request.GetString("var_files", ""); varFiles != "" {
			args = append(args, "--var-files", varFiles)
		}
		if destroy := request.GetBool("destroy", false); destroy {
			args = append(args, "--destroy")
		}
		return executeShipCommand(args)
	})

	// IaC validate configuration tool
	validateConfigTool := mcp.NewTool("iac_plan_validate",
		mcp.WithDescription("Validate Infrastructure as Code configuration"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing IaC files"),
			mcp.Required(),
		),
		mcp.WithString("tool",
			mcp.Description("IaC tool to use"),
			mcp.Required(),
			mcp.Enum("terraform", "terragrunt", "pulumi", "cloudformation"),
		),
	)
	s.AddTool(validateConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		tool := request.GetString("tool", "")
		args := []string{"security", "iac-plan", "--validate", "--workdir", workdir, "--tool", tool}
		return executeShipCommand(args)
	})

	// IaC format configuration tool
	formatConfigTool := mcp.NewTool("iac_plan_format",
		mcp.WithDescription("Format Infrastructure as Code configuration"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing IaC files"),
			mcp.Required(),
		),
		mcp.WithString("tool",
			mcp.Description("IaC tool to use"),
			mcp.Required(),
			mcp.Enum("terraform", "terragrunt", "pulumi", "cloudformation"),
		),
		mcp.WithBoolean("check",
			mcp.Description("Check if files are formatted correctly without modifying"),
		),
	)
	s.AddTool(formatConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		tool := request.GetString("tool", "")
		args := []string{"security", "iac-plan", "--format", "--workdir", workdir, "--tool", tool}
		if check := request.GetBool("check", false); check {
			args = append(args, "--check")
		}
		return executeShipCommand(args)
	})

	// IaC analyze plan tool
	analyzePlanTool := mcp.NewTool("iac_plan_analyze",
		mcp.WithDescription("Analyze Infrastructure as Code plan"),
		mcp.WithString("plan_json",
			mcp.Description("Plan JSON content"),
			mcp.Required(),
		),
		mcp.WithString("analysis_types",
			mcp.Description("Comma-separated analysis types (cost, security, compliance)"),
			mcp.Required(),
		),
	)
	s.AddTool(analyzePlanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		planJson := request.GetString("plan_json", "")
		analysisTypes := request.GetString("analysis_types", "")
		args := []string{"security", "iac-plan", "--analyze", "--plan-json", planJson, "--analysis", analysisTypes}
		return executeShipCommand(args)
	})

	// IaC compare plans tool
	comparePlansTool := mcp.NewTool("iac_plan_compare",
		mcp.WithDescription("Compare Infrastructure as Code plans"),
		mcp.WithString("baseline_plan",
			mcp.Description("Baseline plan file path"),
			mcp.Required(),
		),
		mcp.WithString("current_plan",
			mcp.Description("Current plan file path"),
			mcp.Required(),
		),
	)
	s.AddTool(comparePlansTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		baselinePlan := request.GetString("baseline_plan", "")
		currentPlan := request.GetString("current_plan", "")
		args := []string{"security", "iac-plan", "--compare", "--baseline", baselinePlan, "--current", currentPlan}
		return executeShipCommand(args)
	})

	// IaC manage workspace tool
	manageWorkspaceTool := mcp.NewTool("iac_plan_workspace",
		mcp.WithDescription("Manage Infrastructure as Code workspace"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing IaC files"),
			mcp.Required(),
		),
		mcp.WithString("tool",
			mcp.Description("IaC tool to use"),
			mcp.Required(),
			mcp.Enum("terraform", "terragrunt"),
		),
		mcp.WithString("operation",
			mcp.Description("Workspace operation"),
			mcp.Required(),
			mcp.Enum("list", "show", "new", "select", "delete"),
		),
		mcp.WithString("workspace_name",
			mcp.Description("Workspace name (for new/select/delete operations)"),
		),
	)
	s.AddTool(manageWorkspaceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		tool := request.GetString("tool", "")
		operation := request.GetString("operation", "")
		args := []string{"security", "iac-plan", "--workspace", operation, "--workdir", workdir, "--tool", tool}
		if workspaceName := request.GetString("workspace_name", ""); workspaceName != "" {
			args = append(args, "--name", workspaceName)
		}
		return executeShipCommand(args)
	})

	// IaC generate graph tool
	generateGraphTool := mcp.NewTool("iac_plan_graph",
		mcp.WithDescription("Generate Infrastructure as Code dependency graph"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing IaC files"),
			mcp.Required(),
		),
		mcp.WithString("tool",
			mcp.Description("IaC tool to use"),
			mcp.Required(),
			mcp.Enum("terraform", "terragrunt"),
		),
		mcp.WithString("graph_type",
			mcp.Description("Type of graph to generate"),
			mcp.Enum("plan", "apply", "plan-destroy"),
		),
	)
	s.AddTool(generateGraphTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		tool := request.GetString("tool", "")
		args := []string{"security", "iac-plan", "--graph", "--workdir", workdir, "--tool", tool}
		if graphType := request.GetString("graph_type", ""); graphType != "" {
			args = append(args, "--type", graphType)
		}
		return executeShipCommand(args)
	})
}