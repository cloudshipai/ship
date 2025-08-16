package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGatekeeperTools adds Gatekeeper (Kubernetes policy engine) MCP tool implementations
func AddGatekeeperTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Gatekeeper install tool
	installTool := mcp.NewTool("gatekeeper_install",
		mcp.WithDescription("Install Gatekeeper in Kubernetes cluster"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for installation (default: gatekeeper-system)"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform a dry run without making changes"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gatekeeper", "install"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if request.GetBool("dry_run", false) {
			args = append(args, "--dry-run")
		}
		return executeShipCommand(args)
	})

	// Gatekeeper validate policy tool
	validatePolicyTool := mcp.NewTool("gatekeeper_validate_policy",
		mcp.WithDescription("Validate Gatekeeper policy constraints"),
		mcp.WithString("policy_path",
			mcp.Description("Path to policy file or directory"),
			mcp.Required(),
		),
	)
	s.AddTool(validatePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "gatekeeper", "validate", policyPath}
		return executeShipCommand(args)
	})

	// Gatekeeper scan cluster tool
	scanClusterTool := mcp.NewTool("gatekeeper_scan_cluster",
		mcp.WithDescription("Scan Kubernetes cluster for policy violations"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to scan (default: all namespaces)"),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to specific policy to check"),
		),
	)
	s.AddTool(scanClusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gatekeeper", "scan"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if policyPath := request.GetString("policy_path", ""); policyPath != "" {
			args = append(args, "--policy", policyPath)
		}
		return executeShipCommand(args)
	})

	// Gatekeeper apply constraints tool
	applyConstraintsTool := mcp.NewTool("gatekeeper_apply_constraints",
		mcp.WithDescription("Apply Gatekeeper constraints to cluster"),
		mcp.WithString("constraints_path",
			mcp.Description("Path to constraints file or directory"),
			mcp.Required(),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform a dry run without applying changes"),
		),
	)
	s.AddTool(applyConstraintsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		constraintsPath := request.GetString("constraints_path", "")
		args := []string{"security", "gatekeeper", "apply", constraintsPath}
		if request.GetBool("dry_run", false) {
			args = append(args, "--dry-run")
		}
		return executeShipCommand(args)
	})

	// Gatekeeper generate template tool
	generateTemplateTool := mcp.NewTool("gatekeeper_generate_template",
		mcp.WithDescription("Generate Gatekeeper constraint template"),
		mcp.WithString("template_name",
			mcp.Description("Name for the constraint template"),
			mcp.Required(),
		),
		mcp.WithString("output_path",
			mcp.Description("Output path for generated template"),
		),
	)
	s.AddTool(generateTemplateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateName := request.GetString("template_name", "")
		args := []string{"security", "gatekeeper", "generate", "--template", templateName}
		if outputPath := request.GetString("output_path", ""); outputPath != "" {
			args = append(args, "--output", outputPath)
		}
		return executeShipCommand(args)
	})

	// Gatekeeper list violations tool
	listViolationsTool := mcp.NewTool("gatekeeper_list_violations",
		mcp.WithDescription("List current policy violations in cluster"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to check (default: all namespaces)"),
		),
		mcp.WithString("resource_type",
			mcp.Description("Filter by resource type (pod, deployment, service, etc.)"),
		),
	)
	s.AddTool(listViolationsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gatekeeper", "violations"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if resourceType := request.GetString("resource_type", ""); resourceType != "" {
			args = append(args, "--resource", resourceType)
		}
		return executeShipCommand(args)
	})

	// Gatekeeper get version tool
	getVersionTool := mcp.NewTool("gatekeeper_get_version",
		mcp.WithDescription("Get Gatekeeper version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gatekeeper", "--version"}
		return executeShipCommand(args)
	})
}