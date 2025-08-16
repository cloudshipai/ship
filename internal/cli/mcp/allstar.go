package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddAllstarTools adds Allstar (Kubernetes security policy enforcement) MCP tool implementations
func AddAllstarTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Allstar scan cluster tool
	scanClusterTool := mcp.NewTool("allstar_scan_cluster",
		mcp.WithDescription("Scan Kubernetes cluster for security policy violations using Allstar"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to scan (default: all namespaces)"),
		),
		mcp.WithString("config_path",
			mcp.Description("Path to Allstar configuration file"),
		),
	)
	s.AddTool(scanClusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "allstar", "scan"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if configPath := request.GetString("config_path", ""); configPath != "" {
			args = append(args, "--config", configPath)
		}
		return executeShipCommand(args)
	})

	// Allstar validate policies tool
	validatePolicesTool := mcp.NewTool("allstar_validate_policies",
		mcp.WithDescription("Validate Allstar security policies"),
		mcp.WithString("policy_path",
			mcp.Description("Path to policy directory or file"),
			mcp.Required(),
		),
	)
	s.AddTool(validatePolicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "allstar", "validate", policyPath}
		return executeShipCommand(args)
	})

	// Allstar generate policy tool
	generatePolicyTool := mcp.NewTool("allstar_generate_policy",
		mcp.WithDescription("Generate Allstar security policy templates"),
		mcp.WithString("policy_type",
			mcp.Description("Type of policy to generate (binary_artifacts, outside_collaborators, branch_protection, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("output_path",
			mcp.Description("Output path for generated policy"),
		),
	)
	s.AddTool(generatePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyType := request.GetString("policy_type", "")
		args := []string{"security", "allstar", "generate", "--type", policyType}
		if outputPath := request.GetString("output_path", ""); outputPath != "" {
			args = append(args, "--output", outputPath)
		}
		return executeShipCommand(args)
	})

	// Allstar check compliance tool
	checkComplianceTool := mcp.NewTool("allstar_check_compliance",
		mcp.WithDescription("Check repository compliance with Allstar policies"),
		mcp.WithString("repository",
			mcp.Description("Repository to check (org/repo format)"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level to report (low, medium, high, critical)"),
		),
	)
	s.AddTool(checkComplianceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repository := request.GetString("repository", "")
		args := []string{"security", "allstar", "check", repository}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		return executeShipCommand(args)
	})

	// Allstar install tool
	installTool := mcp.NewTool("allstar_install",
		mcp.WithDescription("Install Allstar in Kubernetes cluster"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for installation (default: allstar-system)"),
		),
		mcp.WithBoolean("dry_run",
			mcp.Description("Perform a dry run without making changes"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "allstar", "install"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if request.GetBool("dry_run", false) {
			args = append(args, "--dry-run")
		}
		return executeShipCommand(args)
	})

	// Allstar get version tool
	getVersionTool := mcp.NewTool("allstar_get_version",
		mcp.WithDescription("Get Allstar version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "allstar", "--version"}
		return executeShipCommand(args)
	})
}