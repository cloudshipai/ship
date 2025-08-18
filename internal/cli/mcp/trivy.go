package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddTrivyTools adds Trivy (universal vulnerability scanner) MCP tool implementations using direct Dagger calls
func AddTrivyTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTrivyToolsDirect(s)
}

// addTrivyToolsDirect adds Trivy tools using direct Dagger module calls
func addTrivyToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		imageName := request.GetString("image_name", "")
		if imageName == "" {
			return mcp.NewToolResultError("image_name is required"), nil
		}

		// Note: Dagger ScanImage function uses fixed parameters, some MCP parameters not directly supported
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" && outputFormat != "json" {
			return mcp.NewToolResultError(fmt.Sprintf("Warning: output_format '%s' not supported, using json", outputFormat)), nil
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			return mcp.NewToolResultError("Warning: output_file parameter not supported with direct Dagger calls"), nil
		}
		if scanners := request.GetString("scanners", ""); scanners != "" {
			return mcp.NewToolResultError("Warning: scanners parameter not supported with direct Dagger calls"), nil
		}
		if request.GetBool("ignore_unfixed", false) {
			return mcp.NewToolResultError("Warning: ignore_unfixed parameter not supported with direct Dagger calls"), nil
		}

		// Scan image
		output, err := module.ScanImage(ctx, imageName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy image scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		dir := request.GetString("directory", ".")

		// Note: Dagger ScanFilesystem function uses fixed parameters, some MCP parameters not directly supported
		if severity := request.GetString("severity", ""); severity != "" && severity != "HIGH,CRITICAL" {
			return mcp.NewToolResultError(fmt.Sprintf("Warning: severity '%s' not supported, using HIGH,CRITICAL", severity)), nil
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" && outputFormat != "json" {
			return mcp.NewToolResultError(fmt.Sprintf("Warning: output_format '%s' not supported, using json", outputFormat)), nil
		}
		if scanners := request.GetString("scanners", ""); scanners != "" {
			return mcp.NewToolResultError("Warning: scanners parameter not supported with direct Dagger calls"), nil
		}
		if skipDirs := request.GetString("skip_dirs", ""); skipDirs != "" {
			return mcp.NewToolResultError("Warning: skip_dirs parameter not supported with direct Dagger calls"), nil
		}
		if request.GetBool("include_dev_deps", false) {
			return mcp.NewToolResultError("Warning: include_dev_deps parameter not supported with direct Dagger calls"), nil
		}

		// Scan filesystem
		output, err := module.ScanFilesystem(ctx, dir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy filesystem scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		repoURL := request.GetString("repo_url", "")
		if repoURL == "" {
			return mcp.NewToolResultError("repo_url is required"), nil
		}

		// Note: Dagger ScanRepository function uses fixed parameters, some MCP parameters not directly supported
		if branch := request.GetString("branch", ""); branch != "" {
			return mcp.NewToolResultError("Warning: branch parameter not supported with direct Dagger calls"), nil
		}
		if commit := request.GetString("commit", ""); commit != "" {
			return mcp.NewToolResultError("Warning: commit parameter not supported with direct Dagger calls"), nil
		}
		if severity := request.GetString("severity", ""); severity != "" && severity != "HIGH,CRITICAL" {
			return mcp.NewToolResultError(fmt.Sprintf("Warning: severity '%s' not supported, using HIGH,CRITICAL", severity)), nil
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" && outputFormat != "json" {
			return mcp.NewToolResultError(fmt.Sprintf("Warning: output_format '%s' not supported, using json", outputFormat)), nil
		}
		if scanners := request.GetString("scanners", ""); scanners != "" {
			return mcp.NewToolResultError("Warning: scanners parameter not supported with direct Dagger calls"), nil
		}

		// Scan repository
		output, err := module.ScanRepository(ctx, repoURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy repository scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		dir := request.GetString("directory", ".")

		// Note: Dagger ScanConfig function uses fixed parameters, some MCP parameters not directly supported
		if severity := request.GetString("severity", ""); severity != "" && severity != "HIGH,CRITICAL" {
			return mcp.NewToolResultError(fmt.Sprintf("Warning: severity '%s' not supported, using HIGH,CRITICAL", severity)), nil
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" && outputFormat != "json" {
			return mcp.NewToolResultError(fmt.Sprintf("Warning: output_format '%s' not supported, using json", outputFormat)), nil
		}
		if policyBundle := request.GetString("policy_bundle", ""); policyBundle != "" {
			return mcp.NewToolResultError("Warning: policy_bundle parameter not supported with direct Dagger calls"), nil
		}
		if configPolicy := request.GetString("config_policy", ""); configPolicy != "" {
			return mcp.NewToolResultError("Warning: config_policy parameter not supported with direct Dagger calls"), nil
		}

		// Scan config
		output, err := module.ScanConfig(ctx, dir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy config scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		sbomPath := request.GetString("sbom_path", "")
		if sbomPath == "" {
			return mcp.NewToolResultError("sbom_path is required"), nil
		}

		// Get additional parameters
		severity := request.GetString("severity", "")
		outputFormat := request.GetString("output_format", "")
		outputFile := request.GetString("output_file", "")
		ignoreUnfixed := request.GetBool("ignore_unfixed", false)

		// Scan SBOM
		output, err := module.ScanSBOM(ctx, sbomPath, severity, outputFormat, outputFile, ignoreUnfixed)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy SBOM scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		target := request.GetString("target", "cluster")
		clusterContext := request.GetString("cluster_context", "")
		namespace := request.GetString("namespace", "")
		severity := request.GetString("severity", "")
		outputFormat := request.GetString("output_format", "")
		scanners := request.GetString("scanners", "")
		includeImages := request.GetBool("include_images", false)

		// Scan Kubernetes
		output, err := module.ScanKubernetes(ctx, target, clusterContext, namespace, severity, outputFormat, scanners, includeImages)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy Kubernetes scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		targetType := request.GetString("target_type", "")
		target := request.GetString("target", "")
		sbomFormat := request.GetString("sbom_format", "")
		outputFile := request.GetString("output_file", "")
		includeDevDeps := request.GetBool("include_dev_deps", false)

		if targetType == "" {
			return mcp.NewToolResultError("target_type is required"), nil
		}
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		if sbomFormat == "" {
			return mcp.NewToolResultError("sbom_format is required"), nil
		}

		// Generate SBOM
		output, err := module.GenerateSBOM(ctx, target, targetType, sbomFormat, outputFile, includeDevDeps)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy SBOM generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		targetType := request.GetString("target_type", "")
		target := request.GetString("target", "")
		severity := request.GetString("severity", "")
		vulnType := request.GetString("vuln_type", "")
		ignoreFile := request.GetString("ignore_file", "")
		ignoreUnfixed := request.GetBool("ignore_unfixed", false)
		exitCode := request.GetString("exit_code", "")

		if targetType == "" {
			return mcp.NewToolResultError("target_type is required"), nil
		}
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		// Scan with filters
		output, err := module.ScanWithFilters(ctx, target, targetType, severity, vulnType, ignoreFile, ignoreUnfixed, exitCode)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy scan with filters failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		operation := request.GetString("operation", "")
		skipUpdate := request.GetBool("skip_update", false)
		cacheDir := request.GetString("cache_dir", "")

		if operation == "" {
			return mcp.NewToolResultError("operation is required"), nil
		}

		// Database operation
		output, err := module.DatabaseUpdate(ctx, operation, skipUpdate, cacheDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy database operation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		listenPort := request.GetString("listen_port", "")
		listenAddress := request.GetString("listen_address", "")
		debug := request.GetBool("debug", false)
		token := request.GetString("token", "")

		// Server mode
		output, err := module.ServerMode(ctx, listenPort, listenAddress, debug, token)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy server mode failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		targetType := request.GetString("target_type", "")
		target := request.GetString("target", "")
		serverURL := request.GetString("server_url", "")
		token := request.GetString("token", "")
		outputFormat := request.GetString("output_format", "")

		if targetType == "" {
			return mcp.NewToolResultError("target_type is required"), nil
		}
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		if serverURL == "" {
			return mcp.NewToolResultError("server_url is required"), nil
		}

		// Client scan
		output, err := module.ClientScan(ctx, target, targetType, serverURL, token, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy client scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		action := request.GetString("action", "")
		pluginName := request.GetString("plugin_name", "")

		if action == "" {
			return mcp.NewToolResultError("action is required"), nil
		}

		// Plugin management
		output, err := module.PluginManagement(ctx, action, pluginName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy plugin management failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get parameters
		inputSBOM := request.GetString("input_sbom", "")
		outputFormat := request.GetString("output_format", "")
		outputFile := request.GetString("output_file", "")

		if inputSBOM == "" {
			return mcp.NewToolResultError("input_sbom is required"), nil
		}
		if outputFormat == "" {
			return mcp.NewToolResultError("output_format is required"), nil
		}

		// Convert SBOM
		output, err := module.ConvertSBOM(ctx, inputSBOM, outputFormat, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy SBOM conversion failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Trivy get version tool
	getVersionTool := mcp.NewTool("trivy_get_version",
		mcp.WithDescription("Get Trivy version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}