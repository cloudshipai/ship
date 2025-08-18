package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddProwlerTools adds Prowler (multi-cloud security assessment) MCP tool implementations using direct Dagger calls
func AddProwlerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addProwlerToolsDirect(s)
}

// addProwlerToolsDirect adds Prowler tools using direct Dagger module calls
func addProwlerToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewProwlerModule(client)

		// Get parameters
		regions := request.GetString("regions", "us-east-1")
		services := request.GetString("services", "")
		compliance := request.GetString("compliance", "")

		// Note: Dagger module doesn't support profile, checks, output_formats, output_directory
		if request.GetString("profile", "") != "" || request.GetString("checks", "") != "" ||
			request.GetString("output_formats", "") != "" || request.GetString("output_directory", "") != "" {
			return mcp.NewToolResultError("profile, checks, output_formats, and output_directory options not supported in Dagger module"), nil
		}

		// Choose appropriate scan method
		var output string
		if compliance != "" {
			output, err = module.ScanWithCompliance(ctx, "aws", compliance, regions)
		} else if services != "" {
			output, err = module.ScanSpecificServices(ctx, "aws", services, regions)
		} else {
			output, err = module.ScanAWS(ctx, "aws", regions)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("prowler AWS scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewProwlerModule(client)

		// Note: Dagger module doesn't support subscription_id, checks, services, compliance, output options
		if request.GetString("subscription_id", "") != "" || request.GetString("checks", "") != "" ||
			request.GetString("services", "") != "" || request.GetString("compliance", "") != "" ||
			request.GetString("output_formats", "") != "" || request.GetString("output_directory", "") != "" {
			return mcp.NewToolResultError("advanced Azure options not supported in Dagger module - uses environment variables"), nil
		}

		// Scan Azure
		output, err := module.ScanAzure(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("prowler Azure scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewProwlerModule(client)

		// Get parameters
		projectId := request.GetString("project_id", "")

		// Note: Dagger module doesn't support checks, services, compliance, output options
		if request.GetString("checks", "") != "" || request.GetString("services", "") != "" ||
			request.GetString("compliance", "") != "" || request.GetString("output_formats", "") != "" ||
			request.GetString("output_directory", "") != "" {
			return mcp.NewToolResultError("advanced GCP options not supported in Dagger module"), nil
		}

		// Scan GCP
		output, err := module.ScanGCP(ctx, projectId)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("prowler GCP scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewProwlerModule(client)

		// Get parameters
		kubeconfigPath := request.GetString("kubeconfig_path", "")
		if kubeconfigPath == "" {
			return mcp.NewToolResultError("kubeconfig_path is required"), nil
		}

		// Note: Dagger module doesn't support context, checks, compliance, output options
		if request.GetString("context", "") != "" || request.GetString("checks", "") != "" ||
			request.GetString("compliance", "") != "" || request.GetString("output_formats", "") != "" ||
			request.GetString("output_directory", "") != "" {
			return mcp.NewToolResultError("advanced Kubernetes options not supported in Dagger module"), nil
		}

		// Scan Kubernetes
		output, err := module.ScanKubernetes(ctx, kubeconfigPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("prowler Kubernetes scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewProwlerModule(client)

		// Get parameters
		provider := request.GetString("provider", "")

		// Note: Dagger module doesn't support service and compliance filtering
		if request.GetString("service", "") != "" || request.GetString("compliance", "") != "" {
			return mcp.NewToolResultError("service and compliance filtering not supported in Dagger module"), nil
		}

		// List checks
		output, err := module.ListChecks(ctx, provider)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("prowler list checks failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewProwlerModule(client)

		// Get parameters
		provider := request.GetString("provider", "")

		// List services
		output, err := module.ListServices(ctx, provider)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("prowler list services failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewProwlerModule(client)

		// Get parameters
		provider := request.GetString("provider", "")

		// List compliance
		output, err := module.ListCompliance(ctx, provider)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("prowler list compliance failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Note: Dashboard server functionality not available in Dagger module
		// Dagger module only supports dashboard generation from scan results
		return mcp.NewToolResultError("dashboard server not supported in Dagger module - only dashboard generation from scan results"), nil
	})

}