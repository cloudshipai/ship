package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddLitmusTools adds Litmus chaos engineering MCP tool implementations
func AddLitmusTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Litmus install tool
	installTool := mcp.NewTool("litmus_install",
		mcp.WithDescription("Install Litmus chaos engineering platform"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for Litmus installation"),
		),
		mcp.WithString("version",
			mcp.Description("Litmus version to install"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "litmus", "install"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if version := request.GetString("version", ""); version != "" {
			args = append(args, "--version", version)
		}
		return executeShipCommand(args)
	})

	// Litmus create experiment tool
	createExperimentTool := mcp.NewTool("litmus_create_experiment",
		mcp.WithDescription("Create Litmus chaos experiment"),
		mcp.WithString("experiment_name",
			mcp.Description("Name of the chaos experiment"),
			mcp.Required(),
		),
		mcp.WithString("target_app",
			mcp.Description("Target application for chaos experiment"),
			mcp.Required(),
		),
		mcp.WithString("chaos_type",
			mcp.Description("Type of chaos (pod-delete, cpu-hog, memory-hog, network-loss)"),
			mcp.Required(),
		),
	)
	s.AddTool(createExperimentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		experimentName := request.GetString("experiment_name", "")
		targetApp := request.GetString("target_app", "")
		chaosType := request.GetString("chaos_type", "")
		args := []string{"kubernetes", "litmus", "create-experiment", experimentName, "--target", targetApp, "--chaos-type", chaosType}
		return executeShipCommand(args)
	})

	// Litmus run experiment tool
	runExperimentTool := mcp.NewTool("litmus_run_experiment",
		mcp.WithDescription("Run Litmus chaos experiment"),
		mcp.WithString("experiment_name",
			mcp.Description("Name of the chaos experiment to run"),
			mcp.Required(),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
	)
	s.AddTool(runExperimentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		experimentName := request.GetString("experiment_name", "")
		args := []string{"kubernetes", "litmus", "run", experimentName}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		return executeShipCommand(args)
	})

	// Litmus get results tool
	getResultsTool := mcp.NewTool("litmus_get_results",
		mcp.WithDescription("Get Litmus chaos experiment results"),
		mcp.WithString("experiment_name",
			mcp.Description("Name of the chaos experiment"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (json, yaml, table)"),
		),
	)
	s.AddTool(getResultsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		experimentName := request.GetString("experiment_name", "")
		args := []string{"kubernetes", "litmus", "results", experimentName}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Litmus list experiments tool
	listExperimentsTool := mcp.NewTool("litmus_list_experiments",
		mcp.WithDescription("List available Litmus chaos experiments"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to list experiments from"),
		),
	)
	s.AddTool(listExperimentsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "litmus", "list"}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		return executeShipCommand(args)
	})

	// Litmus get version tool
	getVersionTool := mcp.NewTool("litmus_get_version",
		mcp.WithDescription("Get Litmus version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kubernetes", "litmus", "--version"}
		return executeShipCommand(args)
	})
}