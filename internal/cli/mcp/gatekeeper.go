package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGatekeeperTools adds Gatekeeper (OPA Kubernetes policy engine) MCP tool implementations using kubectl
func AddGatekeeperTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		version := request.GetString("version", "v3.20.0")
		useHelm := request.GetBool("use_helm", false)
		
		if useHelm {
			args := []string{"helm", "install", "gatekeeper", "gatekeeper/gatekeeper", 
				"--namespace", "gatekeeper-system", "--create-namespace"}
			return executeShipCommand(args)
		} else {
			url := "https://raw.githubusercontent.com/open-policy-agent/gatekeeper/" + version + "/deploy/gatekeeper.yaml"
			args := []string{"kubectl", "apply", "-f", url}
			return executeShipCommand(args)
		}
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
		version := request.GetString("version", "v3.20.0")
		useHelm := request.GetBool("use_helm", false)
		
		if useHelm {
			args := []string{"helm", "delete", "gatekeeper", "--namespace", "gatekeeper-system"}
			return executeShipCommand(args)
		} else {
			url := "https://raw.githubusercontent.com/open-policy-agent/gatekeeper/" + version + "/deploy/gatekeeper.yaml"
			args := []string{"kubectl", "delete", "-f", url}
			return executeShipCommand(args)
		}
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
		templateFile := request.GetString("template_file", "")
		args := []string{"kubectl", "apply", "-f", templateFile}
		return executeShipCommand(args)
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
		constraintFile := request.GetString("constraint_file", "")
		args := []string{"kubectl", "apply", "-f", constraintFile}
		return executeShipCommand(args)
	})

	// Get constraint templates
	getConstraintTemplatesTool := mcp.NewTool("gatekeeper_get_constraint_templates",
		mcp.WithDescription("List Gatekeeper constraint templates"),
	)
	s.AddTool(getConstraintTemplatesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubectl", "get", "constrainttemplates"}
		return executeShipCommand(args)
	})

	// Get constraints
	getConstraintsTool := mcp.NewTool("gatekeeper_get_constraints",
		mcp.WithDescription("List Gatekeeper constraints"),
		mcp.WithString("constraint_type",
			mcp.Description("Specific constraint type to list"),
		),
	)
	s.AddTool(getConstraintsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		constraintType := request.GetString("constraint_type", "")
		if constraintType != "" {
			args := []string{"kubectl", "get", constraintType}
			return executeShipCommand(args)
		} else {
			// List all constraint types by getting constraint templates first
			args := []string{"kubectl", "get", "constrainttemplates", "-o", "name"}
			return executeShipCommand(args)
		}
	})

	// Get Gatekeeper status
	getStatusTool := mcp.NewTool("gatekeeper_get_status",
		mcp.WithDescription("Get Gatekeeper system status"),
	)
	s.AddTool(getStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubectl", "get", "pods", "-n", "gatekeeper-system"}
		return executeShipCommand(args)
	})
}