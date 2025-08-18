package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddKubeBenchTools adds Kube-bench (Kubernetes CIS benchmark) MCP tool implementations using direct Dagger calls
func AddKubeBenchTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addKubeBenchToolsDirect(s)
}

// addKubeBenchToolsDirect adds Kube-bench tools using direct Dagger module calls
func addKubeBenchToolsDirect(s *server.MCPServer) {
	// Kube-bench run tool
	runTool := mcp.NewTool("kube_bench_run",
		mcp.WithDescription("Run CIS Kubernetes benchmark using kube-bench"),
		mcp.WithString("targets",
			mcp.Description("Comma-delimited list of targets to run (master, node, etcd, controlplane, policies)"),
		),
		mcp.WithString("benchmark",
			mcp.Description("Manually specify CIS benchmark version"),
			mcp.Enum("gke-1.0", "gke-1.2.0", "gke-1.6.0", "ack-1.0", "tkgi-1.2.53", "rke-cis-1.7", "rke2-cis-1.7", "k3s-cis-1.7", "rh-0.7", "rh-1.0"),
		),
		mcp.WithString("version",
			mcp.Description("Manually specify Kubernetes version"),
		),
		mcp.WithString("config_dir",
			mcp.Description("Config directory"),
		),
		mcp.WithString("config",
			mcp.Description("Config file path"),
		),
		mcp.WithBoolean("json",
			mcp.Description("Output results in JSON format"),
		),
		mcp.WithBoolean("junit",
			mcp.Description("Output results in JUnit format"),
		),
		mcp.WithString("outputfile",
			mcp.Description("Write results to output file when using JSON or JUnit"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(runTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeBenchModule(client)

		// Get parameters
		kubeconfig := request.GetString("kubeconfig", "")
		targets := request.GetString("targets", "")
		
		// Determine which specific function to call based on targets
		var output string
		if targets != "" {
			// Check for specific target types
			if strings.Contains(targets, "master") {
				output, err = module.RunMasterBenchmark(ctx, kubeconfig)
			} else if strings.Contains(targets, "node") {
				output, err = module.RunNodeBenchmark(ctx, kubeconfig)
			} else {
				// General benchmark with targets
				output, err = module.RunBenchmark(ctx, kubeconfig)
			}
		} else {
			// General benchmark run
			output, err = module.RunBenchmark(ctx, kubeconfig)
		}

		// If specific output format requested, use custom output function
		if request.GetBool("json", false) || request.GetBool("junit", false) {
			outputFormat := ""
			if request.GetBool("json", false) {
				outputFormat = "json"
			} else if request.GetBool("junit", false) {
				outputFormat = "junit"
			}
			outputFile := request.GetString("outputfile", "")
			
			output, err = module.RunWithCustomOutput(ctx, kubeconfig, outputFormat, outputFile)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-bench run failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-bench with specific checks
	runChecksTool := mcp.NewTool("kube_bench_run_checks",
		mcp.WithDescription("Run specific CIS benchmark checks using kube-bench"),
		mcp.WithString("checks",
			mcp.Description("Comma-delimited list of specific check IDs to run"),
			mcp.Required(),
		),
		mcp.WithString("benchmark",
			mcp.Description("Manually specify CIS benchmark version"),
		),
		mcp.WithString("config_dir",
			mcp.Description("Config directory"),
		),
		mcp.WithBoolean("json",
			mcp.Description("Output results in JSON format"),
		),
		mcp.WithString("outputfile",
			mcp.Description("Write results to output file"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(runChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeBenchModule(client)

		// Get parameters
		checks := request.GetString("checks", "")
		if checks == "" {
			return mcp.NewToolResultError("checks is required"), nil
		}
		
		kubeconfig := request.GetString("kubeconfig", "")

		// Run with specific checks
		output, err := module.RunWithChecks(ctx, kubeconfig, checks)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-bench run with checks failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-bench with skip checks
	runSkipTool := mcp.NewTool("kube_bench_run_skip",
		mcp.WithDescription("Run CIS benchmark with skipped checks using kube-bench"),
		mcp.WithString("skip",
			mcp.Description("Comma-delimited list of check IDs to skip"),
			mcp.Required(),
		),
		mcp.WithString("targets",
			mcp.Description("Comma-delimited list of targets to run"),
		),
		mcp.WithString("benchmark",
			mcp.Description("Manually specify CIS benchmark version"),
		),
		mcp.WithBoolean("json",
			mcp.Description("Output results in JSON format"),
		),
		mcp.WithString("outputfile",
			mcp.Description("Write results to output file"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(runSkipTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeBenchModule(client)

		// Get parameters
		skip := request.GetString("skip", "")
		if skip == "" {
			return mcp.NewToolResultError("skip is required"), nil
		}
		
		kubeconfig := request.GetString("kubeconfig", "")

		// Run with skip
		output, err := module.RunWithSkip(ctx, kubeconfig, skip)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-bench run with skip failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-bench with custom output
	runCustomOutputTool := mcp.NewTool("kube_bench_run_custom_output",
		mcp.WithDescription("Run CIS benchmark with custom output format using kube-bench"),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Required(),
			mcp.Enum("json", "junit", "text", "asff"),
		),
		mcp.WithString("outputfile",
			mcp.Description("Write results to output file"),
		),
		mcp.WithString("targets",
			mcp.Description("Comma-delimited list of targets to run"),
		),
		mcp.WithString("benchmark",
			mcp.Description("Manually specify CIS benchmark version"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(runCustomOutputTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeBenchModule(client)

		// Get parameters
		outputFormat := request.GetString("output_format", "")
		if outputFormat == "" {
			return mcp.NewToolResultError("output_format is required"), nil
		}
		
		outputFile := request.GetString("outputfile", "")
		kubeconfig := request.GetString("kubeconfig", "")

		// Run with custom output
		output, err := module.RunWithCustomOutput(ctx, kubeconfig, outputFormat, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-bench custom output failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-bench version
	versionTool := mcp.NewTool("kube_bench_version",
		mcp.WithDescription("Get kube-bench version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeBenchModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get kube-bench version: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-bench ASFF output
	runAsffTool := mcp.NewTool("kube_bench_run_asff",
		mcp.WithDescription("Run CIS benchmark with AWS Security Hub ASFF output using kube-bench"),
		mcp.WithString("aws_account",
			mcp.Description("AWS account ID"),
		),
		mcp.WithString("aws_region",
			mcp.Description("AWS region"),
		),
		mcp.WithString("cluster_arn",
			mcp.Description("EKS cluster ARN"),
		),
		mcp.WithString("outputfile",
			mcp.Description("Write ASFF results to output file"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
	)
	s.AddTool(runAsffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeBenchModule(client)

		// Get kubeconfig
		kubeconfig := request.GetString("kubeconfig", "")

		// Run ASFF
		output, err := module.RunASFF(ctx, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-bench ASFF output failed: %v", err)), nil
		}

		// Add AWS metadata if provided
		awsAccount := request.GetString("aws_account", "")
		awsRegion := request.GetString("aws_region", "")
		clusterArn := request.GetString("cluster_arn", "")
		
		if awsAccount != "" || awsRegion != "" || clusterArn != "" {
			metadata := fmt.Sprintf("\n# AWS Metadata:\n# Account: %s\n# Region: %s\n# Cluster ARN: %s\n", 
				awsAccount, awsRegion, clusterArn)
			output = metadata + output
		}

		return mcp.NewToolResultText(output), nil
	})
}