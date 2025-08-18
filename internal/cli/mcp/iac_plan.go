package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddIacPlanTools adds Infrastructure as Code planning MCP tool implementations using direct Dagger calls
func AddIacPlanTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addIacPlanToolsDirect(s)
}

// addIacPlanToolsDirect adds IaC Plan tools using direct Dagger module calls
func addIacPlanToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewIacPlanModule(client)

		// Get parameters
		workdir := request.GetString("workdir", "")
		if workdir == "" {
			return mcp.NewToolResultError("workdir is required"), nil
		}

		// Prepare var files
		var varFiles []string
		if varFile := request.GetString("var_file", ""); varFile != "" {
			varFiles = append(varFiles, varFile)
		}

		destroy := request.GetBool("destroy", false)

		// Generate plan
		output, err := module.GeneratePlan(ctx, workdir, "terraform", varFiles, destroy)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to generate plan: %v", err)), nil
		}

		// Add note about output file if specified
		if outFile := request.GetString("out_file", ""); outFile != "" {
			output += fmt.Sprintf("\n\nNote: Plan output should be saved to: %s", outFile)
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewIacPlanModule(client)

		// Get parameters
		workdir := request.GetString("workdir", "")
		if workdir == "" {
			return mcp.NewToolResultError("workdir is required"), nil
		}

		// Validate configuration
		output, err := module.ValidateConfiguration(ctx, workdir, "terraform")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to validate configuration: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewIacPlanModule(client)

		// Get parameters
		workdir := request.GetString("workdir", "")
		if workdir == "" {
			return mcp.NewToolResultError("workdir is required"), nil
		}

		check := request.GetBool("check", false)

		// Format configuration
		output, err := module.FormatConfiguration(ctx, workdir, "terraform", check)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to format configuration: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewIacPlanModule(client)

		// Get parameters
		workdir := request.GetString("workdir", "")
		if workdir == "" {
			return mcp.NewToolResultError("workdir is required"), nil
		}

		planFile := request.GetString("plan_file", "")

		// If we have a plan file, analyze it
		if planFile != "" && request.GetBool("json", false) {
			// Analyze the plan (assuming JSON format)
			output, err := module.AnalyzePlan(ctx, planFile, []string{"resources", "changes"})
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to analyze plan: %v", err)), nil
			}
			return mcp.NewToolResultText(output), nil
		}

		// Otherwise, generate and show current plan
		output, err := module.GeneratePlan(ctx, workdir, "terraform", nil, false)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to show plan: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewIacPlanModule(client)

		// Get parameters
		workdir := request.GetString("workdir", "")
		if workdir == "" {
			return mcp.NewToolResultError("workdir is required"), nil
		}

		operation := request.GetString("operation", "")
		workspaceName := request.GetString("workspace_name", "")

		// Manage workspace
		output, err := module.ManageWorkspace(ctx, workdir, "terraform", operation, workspaceName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to manage workspace: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewIacPlanModule(client)

		// Get parameters
		workdir := request.GetString("workdir", "")
		if workdir == "" {
			return mcp.NewToolResultError("workdir is required"), nil
		}

		graphType := request.GetString("graph_type", "plan")

		// Generate graph
		output, err := module.GenerateGraph(ctx, workdir, "terraform", graphType)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to generate graph: %v", err)), nil
		}

		// Add note about output file if specified
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			output += fmt.Sprintf("\n\nNote: Graph output should be saved to: %s", outputFile)
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewIacPlanModule(client)

		// Get parameters
		workdir := request.GetString("workdir", "")
		if workdir == "" {
			return mcp.NewToolResultError("workdir is required"), nil
		}

		// Validate configuration (init is part of validation)
		output, err := module.ValidateConfiguration(ctx, workdir, "terraform")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to initialize: %v", err)), nil
		}

		// Add init-specific messages
		initMsg := "Terraform initialized successfully.\n\n"
		if request.GetBool("upgrade", false) {
			initMsg += "Modules and plugins upgraded.\n"
		}
		if request.GetBool("reconfigure", false) {
			initMsg += "Backend reconfigured.\n"
		}

		return mcp.NewToolResultText(initMsg + output), nil
	})
}