package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddIacPlanTools adds Infrastructure as Code planning MCP tool implementations using real IaC tools
func AddIacPlanTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Terraform plan tool
	terraformPlanTool := mcp.NewTool("iac_plan_terraform_plan",
		mcp.WithDescription("Generate Terraform execution plan"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("var_file",
			mcp.Description("Path to variable file"),
		),
		mcp.WithString("out_file",
			mcp.Description("Output plan to specified file"),
		),
		mcp.WithBoolean("destroy",
			mcp.Description("Generate destroy plan"),
		),
		mcp.WithBoolean("detailed_exitcode",
			mcp.Description("Enable detailed exit codes"),
		),
	)
	s.AddTool(terraformPlanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		args := []string{"sh", "-c", "cd " + workdir + " && terraform plan"}
		
		if varFile := request.GetString("var_file", ""); varFile != "" {
			args = []string{"sh", "-c", "cd " + workdir + " && terraform plan -var-file=" + varFile}
		}
		if outFile := request.GetString("out_file", ""); outFile != "" {
			// Append to existing command
			currentCmd := args[2]
			args[2] = currentCmd + " -out=" + outFile
		}
		if request.GetBool("destroy", false) {
			currentCmd := args[2]
			args[2] = currentCmd + " -destroy"
		}
		if request.GetBool("detailed_exitcode", false) {
			currentCmd := args[2]
			args[2] = currentCmd + " -detailed-exitcode"
		}
		
		return executeShipCommand(args)
	})

	// Terraform validate tool
	terraformValidateTool := mcp.NewTool("iac_plan_terraform_validate",
		mcp.WithDescription("Validate Terraform configuration syntax"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithBoolean("json",
			mcp.Description("Output validation results in JSON format"),
		),
	)
	s.AddTool(terraformValidateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		args := []string{"sh", "-c", "cd " + workdir + " && terraform validate"}
		
		if request.GetBool("json", false) {
			args = []string{"sh", "-c", "cd " + workdir + " && terraform validate -json"}
		}
		
		return executeShipCommand(args)
	})

	// Terraform format tool
	terraformFormatTool := mcp.NewTool("iac_plan_terraform_format",
		mcp.WithDescription("Format Terraform configuration files"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithBoolean("check",
			mcp.Description("Check if files are formatted correctly without modifying"),
		),
		mcp.WithBoolean("diff",
			mcp.Description("Show differences between original and formatted files"),
		),
	)
	s.AddTool(terraformFormatTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		args := []string{"sh", "-c", "cd " + workdir + " && terraform fmt"}
		
		if request.GetBool("check", false) {
			args = []string{"sh", "-c", "cd " + workdir + " && terraform fmt -check"}
		}
		if request.GetBool("diff", false) {
			currentCmd := args[2]
			args[2] = currentCmd + " -diff"
		}
		
		return executeShipCommand(args)
	})

	// Terraform show plan tool
	terraformShowTool := mcp.NewTool("iac_plan_terraform_show",
		mcp.WithDescription("Show Terraform plan in human-readable format"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("plan_file",
			mcp.Description("Path to plan file to show"),
		),
		mcp.WithBoolean("json",
			mcp.Description("Output plan in JSON format"),
		),
	)
	s.AddTool(terraformShowTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		args := []string{"sh", "-c", "cd " + workdir + " && terraform show"}
		
		if planFile := request.GetString("plan_file", ""); planFile != "" {
			args = []string{"sh", "-c", "cd " + workdir + " && terraform show " + planFile}
		}
		if request.GetBool("json", false) {
			currentCmd := args[2]
			args[2] = currentCmd + " -json"
		}
		
		return executeShipCommand(args)
	})

	// Terraform workspace management tool
	terraformWorkspaceTool := mcp.NewTool("iac_plan_terraform_workspace",
		mcp.WithDescription("Manage Terraform workspaces"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing Terraform files"),
			mcp.Required(),
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
	s.AddTool(terraformWorkspaceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		operation := request.GetString("operation", "")
		
		args := []string{"sh", "-c", "cd " + workdir + " && terraform workspace " + operation}
		
		if workspaceName := request.GetString("workspace_name", ""); workspaceName != "" && (operation == "new" || operation == "select" || operation == "delete") {
			args = []string{"sh", "-c", "cd " + workdir + " && terraform workspace " + operation + " " + workspaceName}
		}
		
		return executeShipCommand(args)
	})

	// Terraform graph tool
	terraformGraphTool := mcp.NewTool("iac_plan_terraform_graph",
		mcp.WithDescription("Generate Terraform dependency graph"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("graph_type",
			mcp.Description("Type of graph to generate"),
			mcp.Enum("plan", "apply", "plan-destroy"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for graph (e.g., graph.dot)"),
		),
	)
	s.AddTool(terraformGraphTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		
		var graphCmd string
		if graphType := request.GetString("graph_type", ""); graphType != "" {
			graphCmd = "terraform graph -type=" + graphType
		} else {
			graphCmd = "terraform graph"
		}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			graphCmd += " > " + outputFile
		}
		
		args := []string{"sh", "-c", "cd " + workdir + " && " + graphCmd}
		return executeShipCommand(args)
	})

	// Terraform init tool
	terraformInitTool := mcp.NewTool("iac_plan_terraform_init",
		mcp.WithDescription("Initialize Terraform working directory"),
		mcp.WithString("workdir",
			mcp.Description("Working directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithBoolean("upgrade",
			mcp.Description("Upgrade modules and plugins"),
		),
		mcp.WithBoolean("reconfigure",
			mcp.Description("Reconfigure backend ignoring saved configuration"),
		),
	)
	s.AddTool(terraformInitTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workdir := request.GetString("workdir", "")
		args := []string{"sh", "-c", "cd " + workdir + " && terraform init"}
		
		if request.GetBool("upgrade", false) {
			args = []string{"sh", "-c", "cd " + workdir + " && terraform init -upgrade"}
		}
		if request.GetBool("reconfigure", false) {
			currentCmd := args[2]
			args[2] = currentCmd + " -reconfigure"
		}
		
		return executeShipCommand(args)
	})
}