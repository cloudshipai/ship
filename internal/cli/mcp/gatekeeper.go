package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGatekeeperTools adds Gatekeeper (OPA Kubernetes policy engine) MCP tool implementations using kubectl
func AddGatekeeperTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGatekeeperToolsDirect(s)
}

// addGatekeeperToolsDirect implements direct Dagger calls for Gatekeeper tools
func addGatekeeperToolsDirect(s *server.MCPServer) {
	// Gatekeeper install tool using kubectl
	installTool := mcp.NewTool("gatekeeper_install",
		mcp.WithDescription("Install Gatekeeper using kubectl"),
		mcp.WithString("version",
			mcp.Description("Gatekeeper version to install (default: v3.20.0)"),
		),
		mcp.WithBoolean("use_helm",
			mcp.Description("Use Helm for installation instead of kubectl"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		version := request.GetString("version", "v3.20.0")
		useHelm := request.GetBool("use_helm", false)

		// Create Gatekeeper module and install
		gatekeeperModule := modules.NewGatekeeperModule(client)
		result, err := gatekeeperModule.InstallGatekeeper(ctx, version, useHelm)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("gatekeeper install failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Gatekeeper uninstall tool
	uninstallTool := mcp.NewTool("gatekeeper_uninstall",
		mcp.WithDescription("Uninstall Gatekeeper from cluster"),
		mcp.WithString("version",
			mcp.Description("Gatekeeper version to uninstall (default: v3.20.0)"),
		),
		mcp.WithBoolean("use_helm",
			mcp.Description("Use Helm for uninstallation instead of kubectl"),
		),
	)
	s.AddTool(uninstallTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		version := request.GetString("version", "v3.20.0")
		useHelm := request.GetBool("use_helm", false)

		// Create Gatekeeper module and uninstall
		gatekeeperModule := modules.NewGatekeeperModule(client)
		result, err := gatekeeperModule.UninstallGatekeeper(ctx, version, useHelm)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("gatekeeper uninstall failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Apply constraint template
	applyConstraintTemplateTool := mcp.NewTool("gatekeeper_apply_constraint_template",
		mcp.WithDescription("Apply Gatekeeper constraint template using kubectl"),
		mcp.WithString("template_file",
			mcp.Description("Path to constraint template YAML file"),
			mcp.Required(),
		),
	)
	s.AddTool(applyConstraintTemplateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		templateFile := request.GetString("template_file", "")

		// Create Gatekeeper module and apply constraint template
		gatekeeperModule := modules.NewGatekeeperModule(client)
		result, err := gatekeeperModule.ApplyConstraintTemplate(ctx, templateFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("gatekeeper apply constraint template failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Apply constraint
	applyConstraintTool := mcp.NewTool("gatekeeper_apply_constraint",
		mcp.WithDescription("Apply Gatekeeper constraint using kubectl"),
		mcp.WithString("constraint_file",
			mcp.Description("Path to constraint YAML file"),
			mcp.Required(),
		),
	)
	s.AddTool(applyConstraintTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		constraintFile := request.GetString("constraint_file", "")

		// Create Gatekeeper module and apply constraint
		gatekeeperModule := modules.NewGatekeeperModule(client)
		result, err := gatekeeperModule.ApplyConstraint(ctx, constraintFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("gatekeeper apply constraint failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Get constraint templates
	getConstraintTemplatesTool := mcp.NewTool("gatekeeper_get_constraint_templates",
		mcp.WithDescription("List Gatekeeper constraint templates"),
	)
	s.AddTool(getConstraintTemplatesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create Gatekeeper module and get constraint templates
		gatekeeperModule := modules.NewGatekeeperModule(client)
		result, err := gatekeeperModule.GetConstraintTemplates(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("gatekeeper get constraint templates failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Get constraints
	getConstraintsTool := mcp.NewTool("gatekeeper_get_constraints",
		mcp.WithDescription("List Gatekeeper constraints"),
		mcp.WithString("constraint_type",
			mcp.Description("Specific constraint type to list"),
		),
	)
	s.AddTool(getConstraintsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		constraintType := request.GetString("constraint_type", "")

		// Create Gatekeeper module and get constraints
		gatekeeperModule := modules.NewGatekeeperModule(client)
		result, err := gatekeeperModule.GetConstraints(ctx, constraintType)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("gatekeeper get constraints failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Get Gatekeeper status
	getStatusTool := mcp.NewTool("gatekeeper_get_status",
		mcp.WithDescription("Get Gatekeeper system status"),
	)
	s.AddTool(getStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create Gatekeeper module and get status
		gatekeeperModule := modules.NewGatekeeperModule(client)
		result, err := gatekeeperModule.GetStatus(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("gatekeeper get status failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}