package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddProwlerTools adds Prowler MCP tool implementations
func AddProwlerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Prowler scan AWS tool
	scanAWSTool := mcp.NewTool("prowler_scan_aws",
		mcp.WithDescription("Scan AWS account for security issues using Prowler"),
		mcp.WithString("provider",
			mcp.Description("AWS provider configuration"),
		),
		mcp.WithString("region",
			mcp.Description("AWS region to scan"),
		),
	)
	s.AddTool(scanAWSTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "prowler", "aws"}
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, "--provider", provider)
		}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		return executeShipCommand(args)
	})

	// Prowler scan Azure tool
	scanAzureTool := mcp.NewTool("prowler_scan_azure",
		mcp.WithDescription("Scan Azure subscription for security issues using Prowler"),
	)
	s.AddTool(scanAzureTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "prowler", "azure"}
		return executeShipCommand(args)
	})

	// Prowler scan GCP tool
	scanGCPTool := mcp.NewTool("prowler_scan_gcp",
		mcp.WithDescription("Scan GCP project for security issues using Prowler"),
		mcp.WithString("project_id",
			mcp.Description("GCP project ID to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanGCPTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectId := request.GetString("project_id", "")
		args := []string{"security", "prowler", "gcp", "--project", projectId}
		return executeShipCommand(args)
	})

	// Prowler scan Kubernetes tool
	scanKubernetesTool := mcp.NewTool("prowler_scan_kubernetes",
		mcp.WithDescription("Scan Kubernetes cluster for security issues using Prowler"),
		mcp.WithString("kubeconfig_path",
			mcp.Description("Path to kubeconfig file"),
			mcp.Required(),
		),
	)
	s.AddTool(scanKubernetesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kubeconfigPath := request.GetString("kubeconfig_path", "")
		args := []string{"security", "prowler", "kubernetes", "--kubeconfig", kubeconfigPath}
		return executeShipCommand(args)
	})

	// Prowler scan with compliance tool
	scanWithComplianceTool := mcp.NewTool("prowler_scan_compliance",
		mcp.WithDescription("Scan with specific compliance framework using Prowler"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, azure, gcp)"),
			mcp.Required(),
			mcp.Enum("aws", "azure", "gcp"),
		),
		mcp.WithString("compliance",
			mcp.Description("Compliance framework (cis, pci, gdpr, hipaa, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("Cloud region to scan"),
		),
	)
	s.AddTool(scanWithComplianceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provider := request.GetString("provider", "")
		compliance := request.GetString("compliance", "")
		args := []string{"security", "prowler", provider, "--compliance", compliance}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		return executeShipCommand(args)
	})

	// Prowler scan specific services tool
	scanSpecificServicesTool := mcp.NewTool("prowler_scan_services",
		mcp.WithDescription("Scan specific cloud services using Prowler"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, azure, gcp)"),
			mcp.Required(),
			mcp.Enum("aws", "azure", "gcp"),
		),
		mcp.WithString("services",
			mcp.Description("Comma-separated list of services to scan"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("Cloud region to scan"),
		),
	)
	s.AddTool(scanSpecificServicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provider := request.GetString("provider", "")
		services := request.GetString("services", "")
		args := []string{"security", "prowler", provider, "--services", services}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		return executeShipCommand(args)
	})

	// Prowler comprehensive security scan tool
	comprehensiveSecurityScanTool := mcp.NewTool("prowler_comprehensive_security_scan",
		mcp.WithDescription("Run comprehensive multi-cloud security scan with advanced options"),
		mcp.WithString("accounts",
			mcp.Description("Comma-separated list of cloud account IDs to scan"),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of regions to scan (default: all)"),
		),
		mcp.WithString("severity",
			mcp.Description("Severity levels to include"),
			mcp.Enum("critical", "high", "medium", "low", "critical,high,medium", "critical,high"),
		),
		mcp.WithString("compliance_frameworks",
			mcp.Description("Comma-separated compliance frameworks"),
			mcp.Enum("cis_aws", "aws_foundational_security", "nist_800_53", "soc2", "pci_dss"),
		),
		mcp.WithString("output_formats",
			mcp.Description("Comma-separated output formats"),
			mcp.Enum("json", "html", "csv", "xlsx", "junit", "sarif"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Output directory for results"),
		),
		mcp.WithString("exclude_checks",
			mcp.Description("Comma-separated check IDs to exclude"),
		),
		mcp.WithString("include_checks",
			mcp.Description("Comma-separated specific check IDs to include"),
		),
	)
	s.AddTool(comprehensiveSecurityScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "prowler", "aws", "--comprehensive"}
		
		if accounts := request.GetString("accounts", ""); accounts != "" {
			args = append(args, "--accounts", accounts)
		}
		if regions := request.GetString("regions", ""); regions != "" {
			args = append(args, "--regions", regions)
		}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if compliance := request.GetString("compliance_frameworks", ""); compliance != "" {
			args = append(args, "--compliance", compliance)
		}
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "--output-formats", outputFormats)
		}
		if outputDir := request.GetString("output_directory", ""); outputDir != "" {
			args = append(args, "--output-dir", outputDir)
		}
		if excludeChecks := request.GetString("exclude_checks", ""); excludeChecks != "" {
			args = append(args, "--exclude-checks", excludeChecks)
		}
		if includeChecks := request.GetString("include_checks", ""); includeChecks != "" {
			args = append(args, "--include-checks", includeChecks)
		}
		
		return executeShipCommand(args)
	})

	// Prowler list checks tool
	listChecksTool := mcp.NewTool("prowler_list_checks",
		mcp.WithDescription("List all available Prowler security checks"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider to list checks for"),
			mcp.Enum("aws", "azure", "gcp", "kubernetes"),
		),
		mcp.WithString("service",
			mcp.Description("Specific service to list checks for"),
		),
		mcp.WithString("compliance",
			mcp.Description("Filter by compliance framework"),
		),
		mcp.WithString("severity",
			mcp.Description("Filter by severity level"),
			mcp.Enum("critical", "high", "medium", "low"),
		),
	)
	s.AddTool(listChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "prowler", "--list-checks"}
		
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, "--provider", provider)
		}
		if service := request.GetString("service", ""); service != "" {
			args = append(args, "--service", service)
		}
		if compliance := request.GetString("compliance", ""); compliance != "" {
			args = append(args, "--compliance", compliance)
		}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		
		return executeShipCommand(args)
	})

	// Prowler get version tool
	getVersionTool := mcp.NewTool("prowler_get_version",
		mcp.WithDescription("Get Prowler version and capability information"),
		mcp.WithBoolean("detailed",
			mcp.Description("Include detailed capability information"),
		),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "prowler", "--version"}
		if request.GetBool("detailed", false) {
			args = append(args, "--verbose")
		}
		return executeShipCommand(args)
	})
}