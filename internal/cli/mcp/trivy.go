package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTrivyTools adds Trivy (universal vulnerability scanner) MCP tool implementations using real trivy CLI commands
func AddTrivyTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Trivy scan image tool
	scanImageTool := mcp.NewTool("trivy_scan_image",
		mcp.WithDescription("Scan container image for vulnerabilities using Trivy"),
		mcp.WithString("image_name",
			mcp.Description("Container image name to scan"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Severity levels to include"),
			mcp.Enum("UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif", "template", "cyclonedx", "spdx", "spdx-json", "github", "cosign-vuln"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("scanners",
			mcp.Description("Comma-separated scanners to use"),
		),
		mcp.WithBoolean("ignore_unfixed",
			mcp.Description("Ignore unfixed vulnerabilities"),
		),
	)
	s.AddTool(scanImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"trivy", "image", imageName}
		
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		if scanners := request.GetString("scanners", ""); scanners != "" {
			args = append(args, "--scanners", scanners)
		}
		if request.GetBool("ignore_unfixed", false) {
			args = append(args, "--ignore-unfixed")
		}
		
		return executeShipCommand(args)
	})

	// Trivy scan filesystem tool
	scanFilesystemTool := mcp.NewTool("trivy_scan_filesystem",
		mcp.WithDescription("Scan filesystem for vulnerabilities using Trivy"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("severity",
			mcp.Description("Severity levels to include"),
			mcp.Enum("UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif", "template", "cyclonedx", "spdx", "spdx-json"),
		),
		mcp.WithString("scanners",
			mcp.Description("Comma-separated scanners: vuln,secret,misconfig,license"),
		),
		mcp.WithString("skip_dirs",
			mcp.Description("Comma-separated directories to skip"),
		),
		mcp.WithBoolean("include_dev_deps",
			mcp.Description("Include development dependencies"),
		),
	)
	s.AddTool(scanFilesystemTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"trivy", "fs"}
		dir := request.GetString("directory", ".")
		args = append(args, dir)
		
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if scanners := request.GetString("scanners", ""); scanners != "" {
			args = append(args, "--scanners", scanners)
		}
		if skipDirs := request.GetString("skip_dirs", ""); skipDirs != "" {
			args = append(args, "--skip-dirs", skipDirs)
		}
		if request.GetBool("include_dev_deps", false) {
			args = append(args, "--include-dev-deps")
		}
		
		return executeShipCommand(args)
	})

	// Trivy scan repository tool
	scanRepositoryTool := mcp.NewTool("trivy_scan_repository",
		mcp.WithDescription("Scan git repository for vulnerabilities using Trivy"),
		mcp.WithString("repo_url",
			mcp.Description("Git repository URL to scan"),
			mcp.Required(),
		),
		mcp.WithString("branch",
			mcp.Description("Git branch to scan (default: default branch)"),
		),
		mcp.WithString("commit",
			mcp.Description("Specific commit to scan"),
		),
		mcp.WithString("severity",
			mcp.Description("Severity levels to include"),
			mcp.Enum("UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif", "template", "cyclonedx", "spdx", "spdx-json"),
		),
		mcp.WithString("scanners",
			mcp.Description("Comma-separated scanners: vuln,secret,misconfig,license"),
		),
	)
	s.AddTool(scanRepositoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoURL := request.GetString("repo_url", "")
		args := []string{"trivy", "repo", repoURL}
		
		if branch := request.GetString("branch", ""); branch != "" {
			args = append(args, "--branch", branch)
		}
		if commit := request.GetString("commit", ""); commit != "" {
			args = append(args, "--commit", commit)
		}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if scanners := request.GetString("scanners", ""); scanners != "" {
			args = append(args, "--scanners", scanners)
		}
		
		return executeShipCommand(args)
	})

	// Trivy scan config tool
	scanConfigTool := mcp.NewTool("trivy_scan_config",
		mcp.WithDescription("Scan configuration files for security issues using Trivy"),
		mcp.WithString("directory",
			mcp.Description("Directory containing configuration files (default: current directory)"),
		),
		mcp.WithString("severity",
			mcp.Description("Severity levels to include"),
			mcp.Enum("UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif", "template"),
		),
		mcp.WithString("policy_bundle",
			mcp.Description("Path to policy bundle or URL"),
		),
		mcp.WithString("config_policy",
			mcp.Description("Comma-separated list of config policies to check"),
		),
	)
	s.AddTool(scanConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"trivy", "config"}
		dir := request.GetString("directory", ".")
		args = append(args, dir)
		
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if policyBundle := request.GetString("policy_bundle", ""); policyBundle != "" {
			args = append(args, "--policy-bundle", policyBundle)
		}
		if configPolicy := request.GetString("config_policy", ""); configPolicy != "" {
			args = append(args, "--config-policy", configPolicy)
		}
		
		return executeShipCommand(args)
	})

	// Trivy scan SBOM tool
	scanSBOMTool := mcp.NewTool("trivy_scan_sbom",
		mcp.WithDescription("Scan SBOM file for vulnerabilities using Trivy"),
		mcp.WithString("sbom_path",
			mcp.Description("Path to SBOM file"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Severity levels to include"),
			mcp.Enum("UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif", "template", "cyclonedx", "spdx", "spdx-json"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithBoolean("ignore_unfixed",
			mcp.Description("Ignore unfixed vulnerabilities"),
		),
	)
	s.AddTool(scanSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sbomPath := request.GetString("sbom_path", "")
		args := []string{"trivy", "sbom", sbomPath}
		
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		if request.GetBool("ignore_unfixed", false) {
			args = append(args, "--ignore-unfixed")
		}
		
		return executeShipCommand(args)
	})

	// Trivy scan Kubernetes tool
	scanKubernetesTool := mcp.NewTool("trivy_scan_kubernetes",
		mcp.WithDescription("Scan Kubernetes cluster for vulnerabilities using Trivy"),
		mcp.WithString("target",
			mcp.Description("Kubernetes target to scan"),
			mcp.Enum("cluster", "all", "workload", "node"),
		),
		mcp.WithString("cluster_context",
			mcp.Description("Kubernetes cluster context (default: current context)"),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to scan (default: all namespaces)"),
		),
		mcp.WithString("severity",
			mcp.Description("Severity levels to include"),
			mcp.Enum("UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif"),
		),
		mcp.WithString("scanners",
			mcp.Description("Comma-separated scanners: vuln,misconfig"),
		),
		mcp.WithBoolean("include_images",
			mcp.Description("Include container image scanning"),
		),
	)
	s.AddTool(scanKubernetesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "cluster")
		args := []string{"trivy", "k8s", target}
		
		if clusterContext := request.GetString("cluster_context", ""); clusterContext != "" {
			args = append(args, "--context", clusterContext)
		}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		if scanners := request.GetString("scanners", ""); scanners != "" {
			args = append(args, "--scanners", scanners)
		}
		if request.GetBool("include_images", false) {
			args = append(args, "--include-images")
		}
		
		return executeShipCommand(args)
	})

	// Trivy generate SBOM tool
	generateSBOMTool := mcp.NewTool("trivy_generate_sbom",
		mcp.WithDescription("Generate Software Bill of Materials (SBOM) using Trivy"),
		mcp.WithString("target",
			mcp.Description("Target to generate SBOM for (image, filesystem, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("target_type",
			mcp.Description("Type of target"),
			mcp.Enum("image", "fs", "repo"),
			mcp.Required(),
		),
		mcp.WithString("sbom_format",
			mcp.Description("SBOM format to generate"),
			mcp.Enum("cyclonedx", "spdx", "spdx-json"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for SBOM"),
		),
		mcp.WithBoolean("include_dev_deps",
			mcp.Description("Include development dependencies"),
		),
	)
	s.AddTool(generateSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		targetType := request.GetString("target_type", "")
		target := request.GetString("target", "")
		sbomFormat := request.GetString("sbom_format", "")
		args := []string{"trivy", targetType, "--format", sbomFormat, target}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		if request.GetBool("include_dev_deps", false) {
			args = append(args, "--include-dev-deps")
		}
		
		return executeShipCommand(args)
	})

	// Trivy scan with filters tool
	scanWithFiltersTool := mcp.NewTool("trivy_scan_with_filters",
		mcp.WithDescription("Scan with advanced filtering options using Trivy"),
		mcp.WithString("target",
			mcp.Description("Target to scan"),
			mcp.Required(),
		),
		mcp.WithString("target_type",
			mcp.Description("Type of target"),
			mcp.Enum("image", "fs", "repo", "config"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Comma-separated severity levels"),
		),
		mcp.WithString("vuln_type",
			mcp.Description("Comma-separated vulnerability types: os,library"),
		),
		mcp.WithString("ignore_file",
			mcp.Description("Path to .trivyignore file"),
		),
		mcp.WithBoolean("ignore_unfixed",
			mcp.Description("Ignore unfixed vulnerabilities"),
		),
		mcp.WithString("exit_code",
			mcp.Description("Exit code when vulnerabilities found"),
		),
	)
	s.AddTool(scanWithFiltersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		targetType := request.GetString("target_type", "")
		target := request.GetString("target", "")
		args := []string{"trivy", targetType, target}
		
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if vulnType := request.GetString("vuln_type", ""); vulnType != "" {
			args = append(args, "--vuln-type", vulnType)
		}
		if ignoreFile := request.GetString("ignore_file", ""); ignoreFile != "" {
			args = append(args, "--ignorefile", ignoreFile)
		}
		if request.GetBool("ignore_unfixed", false) {
			args = append(args, "--ignore-unfixed")
		}
		if exitCode := request.GetString("exit_code", ""); exitCode != "" {
			args = append(args, "--exit-code", exitCode)
		}
		
		return executeShipCommand(args)
	})

	// Trivy database operations tool
	databaseOperationsTool := mcp.NewTool("trivy_database_operations",
		mcp.WithDescription("Perform Trivy vulnerability database operations"),
		mcp.WithString("operation",
			mcp.Description("Database operation to perform"),
			mcp.Enum("download", "update", "reset", "clean"),
			mcp.Required(),
		),
		mcp.WithBoolean("skip_update",
			mcp.Description("Skip database update before scanning"),
		),
		mcp.WithString("cache_dir",
			mcp.Description("Custom cache directory path"),
		),
	)
	s.AddTool(databaseOperationsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		operation := request.GetString("operation", "")
		
		switch operation {
		case "download", "update":
			args := []string{"trivy", "image", "--download-db-only", "alpine"}
			if cacheDir := request.GetString("cache_dir", ""); cacheDir != "" {
				args = append(args, "--cache-dir", cacheDir)
			}
			return executeShipCommand(args)
		case "reset", "clean":
			args := []string{"trivy", "clean", "--all"}
			if cacheDir := request.GetString("cache_dir", ""); cacheDir != "" {
				args = append(args, "--cache-dir", cacheDir)
			}
			return executeShipCommand(args)
		default:
			args := []string{"trivy", "--help"}
			return executeShipCommand(args)
		}
	})

	// Trivy server mode tool
	serverModeTool := mcp.NewTool("trivy_server_mode",
		mcp.WithDescription("Run Trivy in server mode for client-server scanning"),
		mcp.WithString("listen_port",
			mcp.Description("Port to listen on (default: 4954)"),
		),
		mcp.WithString("listen_address",
			mcp.Description("Address to listen on (default: 0.0.0.0)"),
		),
		mcp.WithBoolean("debug",
			mcp.Description("Enable debug mode"),
		),
		mcp.WithString("token",
			mcp.Description("Authentication token for server access"),
		),
	)
	s.AddTool(serverModeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"trivy", "server"}
		
		if listenPort := request.GetString("listen_port", ""); listenPort != "" {
			args = append(args, "--listen", "0.0.0.0:"+listenPort)
		}
		if listenAddress := request.GetString("listen_address", ""); listenAddress != "" {
			args = append(args, "--listen", listenAddress)
		}
		if request.GetBool("debug", false) {
			args = append(args, "--debug")
		}
		if token := request.GetString("token", ""); token != "" {
			args = append(args, "--token", token)
		}
		
		return executeShipCommand(args)
	})

	// Trivy client scan tool
	clientScanTool := mcp.NewTool("trivy_client_scan",
		mcp.WithDescription("Scan using Trivy client mode (connect to Trivy server)"),
		mcp.WithString("target",
			mcp.Description("Target to scan"),
			mcp.Required(),
		),
		mcp.WithString("target_type",
			mcp.Description("Type of target"),
			mcp.Enum("image", "fs", "repo"),
			mcp.Required(),
		),
		mcp.WithString("server_url",
			mcp.Description("Trivy server URL (e.g., http://localhost:4954)"),
			mcp.Required(),
		),
		mcp.WithString("token",
			mcp.Description("Authentication token for server access"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif"),
		),
	)
	s.AddTool(clientScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		targetType := request.GetString("target_type", "")
		target := request.GetString("target", "")
		serverURL := request.GetString("server_url", "")
		args := []string{"trivy", targetType, "--server", serverURL, target}
		
		if token := request.GetString("token", ""); token != "" {
			args = append(args, "--token", token)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		
		return executeShipCommand(args)
	})

	// Trivy plugin management tool
	pluginManagementTool := mcp.NewTool("trivy_plugin_management",
		mcp.WithDescription("Manage Trivy plugins for extended functionality"),
		mcp.WithString("action",
			mcp.Description("Plugin action to perform"),
			mcp.Enum("list", "install", "uninstall", "upgrade", "info"),
			mcp.Required(),
		),
		mcp.WithString("plugin_name",
			mcp.Description("Plugin name or URL (for install/uninstall/upgrade/info)"),
		),
	)
	s.AddTool(pluginManagementTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		action := request.GetString("action", "")
		args := []string{"trivy", "plugin", action}
		
		if pluginName := request.GetString("plugin_name", ""); pluginName != "" && action != "list" {
			args = append(args, pluginName)
		}
		
		return executeShipCommand(args)
	})

	// Trivy convert SBOM tool
	convertSBOMTool := mcp.NewTool("trivy_convert_sbom",
		mcp.WithDescription("Convert SBOM between different formats using Trivy"),
		mcp.WithString("input_sbom",
			mcp.Description("Path to input SBOM file"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output SBOM format"),
			mcp.Enum("cyclonedx", "spdx", "spdx-json"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
	)
	s.AddTool(convertSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputSBOM := request.GetString("input_sbom", "")
		outputFormat := request.GetString("output_format", "")
		args := []string{"trivy", "convert", "--format", outputFormat, inputSBOM}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		
		return executeShipCommand(args)
	})

	// Trivy get version tool
	getVersionTool := mcp.NewTool("trivy_get_version",
		mcp.WithDescription("Get Trivy version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"trivy", "--version"}
		return executeShipCommand(args)
	})
}