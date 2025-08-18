package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddLitmusTools adds Litmus chaos engineering MCP tool implementations using direct Dagger calls
func AddLitmusTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addLitmusToolsDirect(s)
}

// addLitmusToolsDirect adds Litmus tools using direct Dagger module calls
func addLitmusToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get parameters
		namespace := request.GetString("namespace", "litmus")
		releaseName := request.GetString("release_name", "chaos")
		createNamespace := request.GetBool("create_namespace", false)

		// Install Litmus
		output, err := module.Install(ctx, namespace, releaseName, createNamespace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("litmus install failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Connect chaos infrastructure using litmusctl
	connectInfraTool := mcp.NewTool("litmus_connect_chaos_infra",
		mcp.WithDescription("Connect chaos infrastructure using litmusctl"),
		mcp.WithString("project_id",
			mcp.Description("Litmus project ID"),
			mcp.Required(),
		),
	)
	s.AddTool(connectInfraTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get parameters
		projectID := request.GetString("project_id", "")
		if projectID == "" {
			return mcp.NewToolResultError("project_id is required"), nil
		}

		// Connect chaos infrastructure
		output, err := module.ConnectChaosInfra(ctx, projectID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("connect chaos infra failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Create project using litmusctl
	createProjectTool := mcp.NewTool("litmus_create_project",
		mcp.WithDescription("Create a new Litmus project using litmusctl"),
	)
	s.AddTool(createProjectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Create project
		output, err := module.CreateProject(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("create project failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Create chaos experiment using litmusctl
	createExperimentTool := mcp.NewTool("litmus_create_chaos_experiment",
		mcp.WithDescription("Create chaos experiment using litmusctl"),
		mcp.WithString("manifest_file",
			mcp.Description("Path to chaos experiment manifest file"),
			mcp.Required(),
		),
		mcp.WithString("project_id",
			mcp.Description("Litmus project ID"),
			mcp.Required(),
		),
		mcp.WithString("chaos_infra_id",
			mcp.Description("Chaos infrastructure ID"),
			mcp.Required(),
		),
	)
	s.AddTool(createExperimentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get parameters
		manifestFile := request.GetString("manifest_file", "")
		if manifestFile == "" {
			return mcp.NewToolResultError("manifest_file is required"), nil
		}
		projectID := request.GetString("project_id", "")
		if projectID == "" {
			return mcp.NewToolResultError("project_id is required"), nil
		}
		chaosInfraID := request.GetString("chaos_infra_id", "")
		if chaosInfraID == "" {
			return mcp.NewToolResultError("chaos_infra_id is required"), nil
		}

		// Create chaos experiment
		output, err := module.CreateChaosExperiment(ctx, manifestFile, projectID, chaosInfraID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("create chaos experiment failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Run chaos experiment using litmusctl
	runExperimentTool := mcp.NewTool("litmus_run_chaos_experiment",
		mcp.WithDescription("Run chaos experiment using litmusctl"),
		mcp.WithString("experiment_id",
			mcp.Description("Chaos experiment ID"),
			mcp.Required(),
		),
		mcp.WithString("project_id",
			mcp.Description("Litmus project ID"),
			mcp.Required(),
		),
	)
	s.AddTool(runExperimentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get parameters
		experimentID := request.GetString("experiment_id", "")
		if experimentID == "" {
			return mcp.NewToolResultError("experiment_id is required"), nil
		}
		projectID := request.GetString("project_id", "")
		if projectID == "" {
			return mcp.NewToolResultError("project_id is required"), nil
		}

		// Run chaos experiment
		output, err := module.RunChaosExperiment(ctx, experimentID, projectID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("run chaos experiment failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Get projects using litmusctl
	getProjectsTool := mcp.NewTool("litmus_get_projects",
		mcp.WithDescription("Get all Litmus projects using litmusctl"),
	)
	s.AddTool(getProjectsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get projects
		output, err := module.GetProjects(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get projects failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Get chaos experiments using litmusctl
	getExperimentsTool := mcp.NewTool("litmus_get_chaos_experiments",
		mcp.WithDescription("Get chaos experiments using litmusctl"),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(getExperimentsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get parameters
		kubeconfig := request.GetString("kubeconfig", "")

		// Get experiments
		output, err := module.GetExperiments(ctx, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get experiments failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Get chaos infrastructure using litmusctl
	getChaosInfraTool := mcp.NewTool("litmus_get_chaos_infra",
		mcp.WithDescription("Get chaos infrastructure using litmusctl"),
		mcp.WithString("project_id",
			mcp.Description("Litmus project ID"),
			mcp.Required(),
		),
	)
	s.AddTool(getChaosInfraTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get parameters
		projectID := request.GetString("project_id", "")
		if projectID == "" {
			return mcp.NewToolResultError("project_id is required"), nil
		}

		// Get chaos infrastructure
		output, err := module.GetChaosInfra(ctx, projectID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get chaos infra failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Config set account using litmusctl
	configSetAccountTool := mcp.NewTool("litmus_config_set_account",
		mcp.WithDescription("Configure account using litmusctl"),
	)
	s.AddTool(configSetAccountTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Config set account
		output, err := module.ConfigSetAccount(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("config set account failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Get version using litmusctl
	versionTool := mcp.NewTool("litmus_version",
		mcp.WithDescription("Get Litmus CLI version using litmusctl"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Apply chaos experiment using kubectl
	applyChaosExperimentTool := mcp.NewTool("litmus_apply_chaos_experiment",
		mcp.WithDescription("Apply chaos experiment manifest using kubectl"),
		mcp.WithString("manifest_file",
			mcp.Description("Path to chaos experiment manifest file"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace (default: litmus)"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(applyChaosExperimentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get parameters
		manifestFile := request.GetString("manifest_file", "")
		if manifestFile == "" {
			return mcp.NewToolResultError("manifest_file is required"), nil
		}
		namespace := request.GetString("namespace", "litmus")
		kubeconfig := request.GetString("kubeconfig", "")

		// Apply chaos experiment
		output, err := module.ApplyChaosExperiment(ctx, manifestFile, namespace, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("apply chaos experiment failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Get chaos results using litmusctl
	getChaosResultsTool := mcp.NewTool("litmus_get_chaos_results",
		mcp.WithDescription("Get chaos experiment results using litmusctl"),
		mcp.WithString("experiment_name",
			mcp.Description("Name of the chaos experiment"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(getChaosResultsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLitmusModule(client)

		// Get parameters
		experimentName := request.GetString("experiment_name", "")
		if experimentName == "" {
			return mcp.NewToolResultError("experiment_name is required"), nil
		}
		kubeconfig := request.GetString("kubeconfig", "")

		// Get chaos results
		output, err := module.GetChaosResults(ctx, experimentName, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get chaos results failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}