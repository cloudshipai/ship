package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddLitmusTools adds Litmus chaos engineering MCP tool implementations using real CLI commands
func AddLitmusTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Litmus install using Helm
	installTool := mcp.NewTool("litmus_install",
		mcp.WithDescription("Install Litmus chaos engineering platform using Helm"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for Litmus installation (default: litmus)"),
		),
		mcp.WithString("release_name",
			mcp.Description("Helm release name (default: chaos)"),
		),
		mcp.WithBoolean("create_namespace",
			mcp.Description("Create namespace if it doesn't exist"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "litmus")
		releaseName := request.GetString("release_name", "chaos")
		
		// Add Litmus Helm repository and update
		repoArgs := []string{"sh", "-c", "helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/ && helm repo update"}
		_, err := executeShipCommand(repoArgs)
		if err != nil {
			return nil, err
		}
		
		// Create namespace if requested
		if request.GetBool("create_namespace", false) {
			nsArgs := []string{"kubectl", "create", "namespace", namespace}
			executeShipCommand(nsArgs) // Ignore error if namespace exists
		}
		
		// Install Litmus
		installArgs := []string{"helm", "install", releaseName, "litmuschaos/litmus", "--namespace", namespace}
		return executeShipCommand(installArgs)
	})

	// Connect chaos infrastructure using litmusctl
	connectInfraTool := mcp.NewTool("litmus_connect_chaos_infra",
		mcp.WithDescription("Connect chaos infrastructure using litmusctl"),
		mcp.WithString("project_id",
			mcp.Description("Project ID for the chaos infrastructure"),
		),
	)
	s.AddTool(connectInfraTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"litmusctl", "connect", "chaos-infra"}
		
		if projectID := request.GetString("project_id", ""); projectID != "" {
			args = append(args, "--project-id", projectID)
		}
		
		return executeShipCommand(args)
	})

	// Create project using litmusctl
	createProjectTool := mcp.NewTool("litmus_create_project",
		mcp.WithDescription("Create a new project using litmusctl"),
	)
	s.AddTool(createProjectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"litmusctl", "create", "project"}
		return executeShipCommand(args)
	})

	// Create chaos experiment using litmusctl
	createExperimentTool := mcp.NewTool("litmus_create_chaos_experiment",
		mcp.WithDescription("Create chaos experiment using litmusctl"),
		mcp.WithString("manifest_file",
			mcp.Description("Path to chaos experiment manifest file"),
			mcp.Required(),
		),
		mcp.WithString("project_id",
			mcp.Description("Project ID for the experiment"),
		),
		mcp.WithString("chaos_infra_id",
			mcp.Description("Chaos infrastructure ID"),
		),
	)
	s.AddTool(createExperimentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		manifestFile := request.GetString("manifest_file", "")
		args := []string{"litmusctl", "create", "chaos-experiment", "-f", manifestFile}
		
		if projectID := request.GetString("project_id", ""); projectID != "" {
			args = append(args, "--project-id", projectID)
		}
		if chaosInfraID := request.GetString("chaos_infra_id", ""); chaosInfraID != "" {
			args = append(args, "--chaos-infra-id", chaosInfraID)
		}
		
		return executeShipCommand(args)
	})

	// Run chaos experiment using litmusctl
	runExperimentTool := mcp.NewTool("litmus_run_chaos_experiment",
		mcp.WithDescription("Run chaos experiment using litmusctl"),
		mcp.WithString("experiment_id",
			mcp.Description("Chaos experiment ID to run"),
			mcp.Required(),
		),
		mcp.WithString("project_id",
			mcp.Description("Project ID"),
		),
	)
	s.AddTool(runExperimentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		experimentID := request.GetString("experiment_id", "")
		args := []string{"litmusctl", "run", "chaos-experiment", experimentID}
		
		if projectID := request.GetString("project_id", ""); projectID != "" {
			args = append(args, "--project-id", projectID)
		}
		
		return executeShipCommand(args)
	})

	// Get projects using litmusctl
	getProjectsTool := mcp.NewTool("litmus_get_projects",
		mcp.WithDescription("List projects using litmusctl"),
	)
	s.AddTool(getProjectsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"litmusctl", "get", "projects"}
		return executeShipCommand(args)
	})

	// Get chaos experiments using litmusctl
	getExperimentsTool := mcp.NewTool("litmus_get_chaos_experiments",
		mcp.WithDescription("List chaos experiments using litmusctl"),
		mcp.WithString("project_id",
			mcp.Description("Project ID to filter experiments"),
		),
	)
	s.AddTool(getExperimentsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"litmusctl", "get", "chaos-experiment"}
		
		if projectID := request.GetString("project_id", ""); projectID != "" {
			args = append(args, "--project-id", projectID)
		}
		
		return executeShipCommand(args)
	})

	// Get chaos infrastructure using litmusctl
	getChaosInfraTool := mcp.NewTool("litmus_get_chaos_infra",
		mcp.WithDescription("List chaos infrastructure using litmusctl"),
		mcp.WithString("project_id",
			mcp.Description("Project ID to filter infrastructure"),
		),
	)
	s.AddTool(getChaosInfraTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"litmusctl", "get", "chaos-infra"}
		
		if projectID := request.GetString("project_id", ""); projectID != "" {
			args = append(args, "--project-id", projectID)
		}
		
		return executeShipCommand(args)
	})

	// Configure litmusctl
	configSetAccountTool := mcp.NewTool("litmus_config_set_account",
		mcp.WithDescription("Setup ChaosCenter account configuration using litmusctl"),
	)
	s.AddTool(configSetAccountTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"litmusctl", "config", "set-account"}
		return executeShipCommand(args)
	})

	// Get litmusctl version
	versionTool := mcp.NewTool("litmus_version",
		mcp.WithDescription("Get litmusctl version information"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"litmusctl", "version"}
		return executeShipCommand(args)
	})

	// Apply chaos experiment manifest using kubectl
	applyChaosExperimentTool := mcp.NewTool("litmus_apply_chaos_experiment",
		mcp.WithDescription("Apply chaos experiment manifest using kubectl"),
		mcp.WithString("manifest_file",
			mcp.Description("Path to chaos experiment YAML manifest"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to apply the experiment"),
		),
	)
	s.AddTool(applyChaosExperimentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		manifestFile := request.GetString("manifest_file", "")
		args := []string{"kubectl", "apply", "-f", manifestFile}
		
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		}
		
		return executeShipCommand(args)
	})

	// Get chaos experiment results using kubectl
	getChaosResultsTool := mcp.NewTool("litmus_get_chaos_results",
		mcp.WithDescription("Get chaos experiment results using kubectl"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to check results"),
		),
		mcp.WithString("experiment_name",
			mcp.Description("Specific experiment name to get results for"),
		),
	)
	s.AddTool(getChaosResultsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubectl", "get", "chaosresult"}
		
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "-n", namespace)
		}
		if experimentName := request.GetString("experiment_name", ""); experimentName != "" {
			args = append(args, experimentName)
		}
		
		return executeShipCommand(args)
	})
}