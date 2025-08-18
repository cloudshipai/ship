package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKubeBenchTools adds Kube-bench (Kubernetes CIS benchmark) MCP tool implementations using real CLI commands
func AddKubeBenchTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
	)
	s.AddTool(runTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kube-bench", "run"}
		
		if targets := request.GetString("targets", ""); targets != "" {
			args = append(args, "--targets", targets)
		}
		if benchmark := request.GetString("benchmark", ""); benchmark != "" {
			args = append(args, "--benchmark", benchmark)
		}
		if version := request.GetString("version", ""); version != "" {
			args = append(args, "--version", version)
		}
		if configDir := request.GetString("config_dir", ""); configDir != "" {
			args = append(args, "--config-dir", configDir)
		}
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json", false) {
			args = append(args, "--json")
		}
		if request.GetBool("junit", false) {
			args = append(args, "--junit")
		}
		if outputfile := request.GetString("outputfile", ""); outputfile != "" {
			args = append(args, "--outputfile", outputfile)
		}
		
		return executeShipCommand(args)
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
	)
	s.AddTool(runChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		checks := request.GetString("checks", "")
		args := []string{"kube-bench", "--check", checks}
		
		if benchmark := request.GetString("benchmark", ""); benchmark != "" {
			args = append(args, "--benchmark", benchmark)
		}
		if configDir := request.GetString("config_dir", ""); configDir != "" {
			args = append(args, "--config-dir", configDir)
		}
		if request.GetBool("json", false) {
			args = append(args, "--json")
		}
		if outputfile := request.GetString("outputfile", ""); outputfile != "" {
			args = append(args, "--outputfile", outputfile)
		}
		
		return executeShipCommand(args)
	})

	// Kube-bench with skipped checks
	runSkipTool := mcp.NewTool("kube_bench_run_skip",
		mcp.WithDescription("Run CIS benchmark skipping specific checks using kube-bench"),
		mcp.WithString("skip",
			mcp.Description("Comma-delimited list of check IDs or groups to skip"),
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
	)
	s.AddTool(runSkipTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		skip := request.GetString("skip", "")
		args := []string{"kube-bench", "--skip", skip}
		
		if targets := request.GetString("targets", ""); targets != "" {
			args = append(args, "run", "--targets", targets)
		}
		if benchmark := request.GetString("benchmark", ""); benchmark != "" {
			args = append(args, "--benchmark", benchmark)
		}
		if request.GetBool("json", false) {
			args = append(args, "--json")
		}
		if outputfile := request.GetString("outputfile", ""); outputfile != "" {
			args = append(args, "--outputfile", outputfile)
		}
		
		return executeShipCommand(args)
	})

	// Kube-bench with output options
	runCustomOutputTool := mcp.NewTool("kube_bench_run_custom_output",
		mcp.WithDescription("Run CIS benchmark with custom output options using kube-bench"),
		mcp.WithString("targets",
			mcp.Description("Comma-delimited list of targets to run"),
		),
		mcp.WithString("benchmark",
			mcp.Description("Manually specify CIS benchmark version"),
		),
		mcp.WithBoolean("noremediations",
			mcp.Description("Disable printing of remediations section"),
		),
		mcp.WithBoolean("noresults",
			mcp.Description("Disable printing of results section"),
		),
		mcp.WithBoolean("nosummary",
			mcp.Description("Disable printing of summary section"),
		),
		mcp.WithBoolean("nototals",
			mcp.Description("Disable printing of totals for failed, passed checks"),
		),
		mcp.WithBoolean("scored",
			mcp.Description("Run only scored CIS checks"),
		),
		mcp.WithBoolean("unscored",
			mcp.Description("Run only unscored CIS checks"),
		),
		mcp.WithString("exit_code",
			mcp.Description("Specify exit code for when checks fail"),
		),
	)
	s.AddTool(runCustomOutputTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kube-bench"}
		
		if targets := request.GetString("targets", ""); targets != "" {
			args = append(args, "run", "--targets", targets)
		}
		if benchmark := request.GetString("benchmark", ""); benchmark != "" {
			args = append(args, "--benchmark", benchmark)
		}
		if request.GetBool("noremediations", false) {
			args = append(args, "--noremediations")
		}
		if request.GetBool("noresults", false) {
			args = append(args, "--noresults")
		}
		if request.GetBool("nosummary", false) {
			args = append(args, "--nosummary")
		}
		if request.GetBool("nototals", false) {
			args = append(args, "--nototals")
		}
		if request.GetBool("scored", false) {
			args = append(args, "--scored")
		}
		if request.GetBool("unscored", false) {
			args = append(args, "--unscored")
		}
		if exitCode := request.GetString("exit_code", ""); exitCode != "" {
			args = append(args, "--exit-code", exitCode)
		}
		
		return executeShipCommand(args)
	})

	// Kube-bench version
	versionTool := mcp.NewTool("kube_bench_version",
		mcp.WithDescription("Get kube-bench version information"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kube-bench", "version"}
		return executeShipCommand(args)
	})

	// Kube-bench with AWS Security Hub integration
	runAsffTool := mcp.NewTool("kube_bench_run_asff",
		mcp.WithDescription("Run CIS benchmark and send results to AWS Security Hub using kube-bench"),
		mcp.WithString("targets",
			mcp.Description("Comma-delimited list of targets to run"),
		),
		mcp.WithString("benchmark",
			mcp.Description("Manually specify CIS benchmark version"),
		),
		mcp.WithBoolean("asff",
			mcp.Description("Send results to AWS Security Hub"),
		),
	)
	s.AddTool(runAsffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kube-bench"}
		
		if targets := request.GetString("targets", ""); targets != "" {
			args = append(args, "run", "--targets", targets)
		}
		if benchmark := request.GetString("benchmark", ""); benchmark != "" {
			args = append(args, "--benchmark", benchmark)
		}
		if request.GetBool("asff", false) {
			args = append(args, "--asff")
		}
		
		return executeShipCommand(args)
	})
}