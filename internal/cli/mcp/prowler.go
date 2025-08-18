package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddProwlerTools adds Prowler (multi-cloud security assessment) MCP tool implementations using real CLI commands
func AddProwlerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Prowler scan AWS tool
	scanAWSTool := mcp.NewTool("prowler_aws",
		mcp.WithDescription("Scan AWS account for security issues using real prowler CLI"),
		mcp.WithString("profile",
			mcp.Description("AWS CLI profile to use"),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of AWS regions to scan"),
		),
		mcp.WithString("checks",
			mcp.Description("Comma-separated list of specific checks to run"),
		),
		mcp.WithString("services",
			mcp.Description("Comma-separated list of AWS services to scan"),
		),
		mcp.WithString("compliance",
			mcp.Description("Compliance framework to check against"),
		),
		mcp.WithString("output_formats",
			mcp.Description("Output formats (csv, json-asff, json-ocsf, html)"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Output directory for results"),
		),
	)
	s.AddTool(scanAWSTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"prowler", "aws"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if regions := request.GetString("regions", ""); regions != "" {
			args = append(args, "--regions", regions)
		}
		if checks := request.GetString("checks", ""); checks != "" {
			args = append(args, "--checks", checks)
		}
		if services := request.GetString("services", ""); services != "" {
			args = append(args, "--services", services)
		}
		if compliance := request.GetString("compliance", ""); compliance != "" {
			args = append(args, "--compliance", compliance)
		}
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "--output-formats", outputFormats)
		}
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			args = append(args, "--output-directory", outputDirectory)
		}
		
		return executeShipCommand(args)
	})

	// Prowler scan Azure tool
	scanAzureTool := mcp.NewTool("prowler_azure",
		mcp.WithDescription("Scan Azure subscription for security issues using real prowler CLI"),
		mcp.WithString("subscription_id",
			mcp.Description("Azure subscription ID to scan"),
		),
		mcp.WithString("checks",
			mcp.Description("Comma-separated list of specific checks to run"),
		),
		mcp.WithString("services",
			mcp.Description("Comma-separated list of Azure services to scan"),
		),
		mcp.WithString("compliance",
			mcp.Description("Compliance framework to check against"),
		),
		mcp.WithString("output_formats",
			mcp.Description("Output formats (csv, json-asff, json-ocsf, html)"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Output directory for results"),
		),
	)
	s.AddTool(scanAzureTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"prowler", "azure"}
		
		if subscriptionId := request.GetString("subscription_id", ""); subscriptionId != "" {
			args = append(args, "--subscription-id", subscriptionId)
		}
		if checks := request.GetString("checks", ""); checks != "" {
			args = append(args, "--checks", checks)
		}
		if services := request.GetString("services", ""); services != "" {
			args = append(args, "--services", services)
		}
		if compliance := request.GetString("compliance", ""); compliance != "" {
			args = append(args, "--compliance", compliance)
		}
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "--output-formats", outputFormats)
		}
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			args = append(args, "--output-directory", outputDirectory)
		}
		
		return executeShipCommand(args)
	})

	// Prowler scan GCP tool
	scanGCPTool := mcp.NewTool("prowler_gcp",
		mcp.WithDescription("Scan GCP project for security issues using real prowler CLI"),
		mcp.WithString("project_id",
			mcp.Description("GCP project ID to scan"),
		),
		mcp.WithString("checks",
			mcp.Description("Comma-separated list of specific checks to run"),
		),
		mcp.WithString("services",
			mcp.Description("Comma-separated list of GCP services to scan"),
		),
		mcp.WithString("compliance",
			mcp.Description("Compliance framework to check against"),
		),
		mcp.WithString("output_formats",
			mcp.Description("Output formats (csv, json-asff, json-ocsf, html)"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Output directory for results"),
		),
	)
	s.AddTool(scanGCPTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"prowler", "gcp"}
		
		if projectId := request.GetString("project_id", ""); projectId != "" {
			args = append(args, "--project-id", projectId)
		}
		if checks := request.GetString("checks", ""); checks != "" {
			args = append(args, "--checks", checks)
		}
		if services := request.GetString("services", ""); services != "" {
			args = append(args, "--services", services)
		}
		if compliance := request.GetString("compliance", ""); compliance != "" {
			args = append(args, "--compliance", compliance)
		}
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "--output-formats", outputFormats)
		}
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			args = append(args, "--output-directory", outputDirectory)
		}
		
		return executeShipCommand(args)
	})

	// Prowler scan Kubernetes tool
	scanKubernetesTool := mcp.NewTool("prowler_kubernetes",
		mcp.WithDescription("Scan Kubernetes cluster for security issues using real prowler CLI"),
		mcp.WithString("kubeconfig_path",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("context",
			mcp.Description("Kubernetes context to use"),
		),
		mcp.WithString("checks",
			mcp.Description("Comma-separated list of specific checks to run"),
		),
		mcp.WithString("compliance",
			mcp.Description("Compliance framework to check against"),
		),
		mcp.WithString("output_formats",
			mcp.Description("Output formats (csv, json-asff, json-ocsf, html)"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Output directory for results"),
		),
	)
	s.AddTool(scanKubernetesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"prowler", "kubernetes"}
		
		if kubeconfigPath := request.GetString("kubeconfig_path", ""); kubeconfigPath != "" {
			args = append(args, "--kubeconfig", kubeconfigPath)
		}
		if context := request.GetString("context", ""); context != "" {
			args = append(args, "--context", context)
		}
		if checks := request.GetString("checks", ""); checks != "" {
			args = append(args, "--checks", checks)
		}
		if compliance := request.GetString("compliance", ""); compliance != "" {
			args = append(args, "--compliance", compliance)
		}
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "--output-formats", outputFormats)
		}
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			args = append(args, "--output-directory", outputDirectory)
		}
		
		return executeShipCommand(args)
	})

	// Prowler list checks tool
	listChecksTool := mcp.NewTool("prowler_list_checks",
		mcp.WithDescription("List available security checks using real prowler CLI"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, azure, gcp, kubernetes)"),
			mcp.Enum("aws", "azure", "gcp", "kubernetes"),
		),
		mcp.WithString("service",
			mcp.Description("Specific service to list checks for"),
		),
		mcp.WithString("compliance",
			mcp.Description("Filter by compliance framework"),
		),
	)
	s.AddTool(listChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"prowler"}
		
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, provider)
		}
		
		args = append(args, "--list-checks")
		
		if service := request.GetString("service", ""); service != "" {
			args = append(args, "--service", service)
		}
		if compliance := request.GetString("compliance", ""); compliance != "" {
			args = append(args, "--compliance", compliance)
		}
		
		return executeShipCommand(args)
	})

	// Prowler list services tool
	listServicesTool := mcp.NewTool("prowler_list_services",
		mcp.WithDescription("List available services using real prowler CLI"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, azure, gcp, kubernetes)"),
			mcp.Enum("aws", "azure", "gcp", "kubernetes"),
		),
	)
	s.AddTool(listServicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"prowler"}
		
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, provider)
		}
		
		args = append(args, "--list-services")
		return executeShipCommand(args)
	})

	// Prowler list compliance tool
	listComplianceTool := mcp.NewTool("prowler_list_compliance",
		mcp.WithDescription("List available compliance frameworks using real prowler CLI"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, azure, gcp, kubernetes)"),
			mcp.Enum("aws", "azure", "gcp", "kubernetes"),
		),
	)
	s.AddTool(listComplianceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"prowler"}
		
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, provider)
		}
		
		args = append(args, "--list-compliance")
		return executeShipCommand(args)
	})

	// Prowler dashboard tool
	dashboardTool := mcp.NewTool("prowler_dashboard",
		mcp.WithDescription("Launch Prowler dashboard using real prowler CLI"),
		mcp.WithString("port",
			mcp.Description("Port to run dashboard on (default: 11666)"),
		),
		mcp.WithString("host",
			mcp.Description("Host to bind dashboard to (default: 127.0.0.1)"),
		),
	)
	s.AddTool(dashboardTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"prowler", "dashboard"}
		
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "--port", port)
		}
		if host := request.GetString("host", ""); host != "" {
			args = append(args, "--host", host)
		}
		
		return executeShipCommand(args)
	})

}