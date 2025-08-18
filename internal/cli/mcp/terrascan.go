package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddTerrascanTools adds Terrascan (IaC security scanner) MCP tool implementations using direct Dagger calls
func AddTerrascanTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTerrascanToolsDirect(s)
}

// addTerrascanToolsDirect adds Terrascan tools using direct Dagger module calls
func addTerrascanToolsDirect(s *server.MCPServer) {
	// Terrascan scan directory tool
	scanDirectoryTool := mcp.NewTool("terrascan_scan_directory",
		mcp.WithDescription("Scan directory for IaC security issues using Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		directory := request.GetString("directory", "")

		// Scan directory
		output, err := module.ScanDirectory(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan directory scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan scan Terraform files tool
	scanTerraformTool := mcp.NewTool("terrascan_scan_terraform",
		mcp.WithDescription("Scan Terraform files specifically using Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanTerraformTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		directory := request.GetString("directory", "")

		// Scan Terraform
		output, err := module.ScanTerraform(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan Terraform scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan scan Kubernetes manifests tool
	scanKubernetesTool := mcp.NewTool("terrascan_scan_kubernetes",
		mcp.WithDescription("Scan Kubernetes manifests using Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Kubernetes manifests"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
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
		module := modules.NewTerrascanModule(client)

		// Get parameters
		directory := request.GetString("directory", "")

		// Scan Kubernetes
		output, err := module.ScanKubernetes(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan Kubernetes scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan scan with severity filter tool
	scanWithSeverityTool := mcp.NewTool("terrascan_scan_with_severity",
		mcp.WithDescription("Scan with minimum severity level using Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level"),
			mcp.Required(),
			mcp.Enum("LOW", "MEDIUM", "HIGH"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanWithSeverityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		severity := request.GetString("severity", "")

		// Scan with severity
		output, err := module.ScanWithSeverity(ctx, directory, severity, "")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan severity scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan scan remote repository tool
	scanRemoteTool := mcp.NewTool("terrascan_scan_remote",
		mcp.WithDescription("Scan remote repository using Terrascan"),
		mcp.WithString("repo_url",
			mcp.Description("URL of the remote repository"),
			mcp.Required(),
		),
		mcp.WithString("repo_type",
			mcp.Description("Type of repository"),
			mcp.Enum("git", "s3", "gcs", "http"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanRemoteTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		repoURL := request.GetString("repo_url", "")
		repoType := request.GetString("repo_type", "")
		outputFormat := request.GetString("output_format", "json")

		// Scan remote repository
		output, err := module.ScanRemote(ctx, repoURL, repoType, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan remote scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan scan with policy path tool
	scanWithPolicyTool := mcp.NewTool("terrascan_scan_with_policy",
		mcp.WithDescription("Scan using custom policy path with Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to custom policies"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanWithPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		policyPath := request.GetString("policy_path", "")
		outputFormat := request.GetString("output_format", "json")

		// Scan with custom policy
		output, err := module.ScanWithPolicy(ctx, directory, policyPath, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan policy scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan comprehensive IaC security scanning tool
	comprehensiveIaCScanTool := mcp.NewTool("terrascan_comprehensive_iac_scan",
		mcp.WithDescription("Comprehensive Infrastructure as Code security scanning"),
		mcp.WithString("target",
			mcp.Description("Target to scan (directory, file, or remote repository)"),
			mcp.Required(),
		),
		mcp.WithString("iac_type",
			mcp.Description("Infrastructure as Code type"),
			mcp.Enum("terraform", "k8s", "helm", "kustomize", "dockercompose", "arm", "cfn"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for results"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif", "csv", "html"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for results"),
		),
		mcp.WithString("severity_threshold",
			mcp.Description("Minimum severity threshold"),
			mcp.Enum("LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
		mcp.WithString("policy_types",
			mcp.Description("Comma-separated policy types to include"),
		),
		mcp.WithString("exclude_rules",
			mcp.Description("Comma-separated rule IDs to exclude"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose output"),
		),
		mcp.WithBoolean("show_passed",
			mcp.Description("Show passed checks in output"),
		),
	)
	s.AddTool(comprehensiveIaCScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		target := request.GetString("target", "")
		iacType := request.GetString("iac_type", "")
		outputFormat := request.GetString("output_format", "json")
		outputFile := request.GetString("output_file", "")
		severityThreshold := request.GetString("severity_threshold", "")
		policyTypes := request.GetString("policy_types", "")
		excludeRules := request.GetString("exclude_rules", "")
		verbose := request.GetBool("verbose", false)
		showPassed := request.GetBool("show_passed", false)

		// Comprehensive IaC scan
		output, err := module.ComprehensiveIaCScan(ctx, target, iacType, outputFormat, outputFile, severityThreshold, policyTypes, excludeRules, verbose, showPassed)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan comprehensive scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan compliance framework scanning tool
	complianceFrameworkScanTool := mcp.NewTool("terrascan_compliance_framework_scan",
		mcp.WithDescription("Scan against compliance frameworks using Terrascan"),
		mcp.WithString("target",
			mcp.Description("Target to scan for compliance"),
			mcp.Required(),
		),
		mcp.WithString("compliance_framework",
			mcp.Description("Compliance framework to validate against"),
			mcp.Enum("cis", "nist", "pci", "gdpr", "hipaa", "iso27001", "sox", "aws-foundational"),
			mcp.Required(),
		),
		mcp.WithString("iac_type",
			mcp.Description("Infrastructure as Code type"),
			mcp.Enum("terraform", "k8s", "helm", "cfn", "arm"),
		),
		mcp.WithString("output_format",
			mcp.Description("Compliance report format"),
			mcp.Enum("json", "yaml", "sarif", "junit-xml", "html"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for compliance report"),
		),
		mcp.WithBoolean("include_severity_details",
			mcp.Description("Include detailed severity analysis"),
		),
	)
	s.AddTool(complianceFrameworkScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		target := request.GetString("target", "")
		complianceFramework := request.GetString("compliance_framework", "")
		iacType := request.GetString("iac_type", "terraform")
		outputFormat := request.GetString("output_format", "json")
		outputFile := request.GetString("output_file", "")
		includeSeverityDetails := request.GetBool("include_severity_details", false)

		// Compliance framework scan
		output, err := module.ComplianceFrameworkScan(ctx, target, complianceFramework, iacType, outputFormat, outputFile, includeSeverityDetails)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan compliance scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan remote repository scanning tool
	remoteRepositoryScanTool := mcp.NewTool("terrascan_remote_repository_scan",
		mcp.WithDescription("Advanced remote repository scanning with authentication"),
		mcp.WithString("repository_url",
			mcp.Description("Remote repository URL to scan"),
			mcp.Required(),
		),
		mcp.WithString("repository_type",
			mcp.Description("Remote repository type"),
			mcp.Enum("git", "s3", "gcs", "azure-repos", "bitbucket", "tarball"),
		),
		mcp.WithString("iac_type",
			mcp.Description("Infrastructure as Code type"),
			mcp.Enum("terraform", "k8s", "helm", "kustomize", "dockercompose"),
		),
		mcp.WithString("branch",
			mcp.Description("Specific branch to scan"),
		),
		mcp.WithString("ssh_key_path",
			mcp.Description("Path to SSH private key for authentication"),
		),
		mcp.WithString("access_token",
			mcp.Description("Access token for repository authentication"),
		),
		mcp.WithString("config_path",
			mcp.Description("Path to Terrascan configuration file"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for scan results"),
			mcp.Enum("json", "yaml", "sarif", "junit-xml"),
		),
	)
	s.AddTool(remoteRepositoryScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		repoURL := request.GetString("repository_url", "")
		repoType := request.GetString("repository_type", "git")
		iacType := request.GetString("iac_type", "terraform")
		branch := request.GetString("branch", "")
		sshKeyPath := request.GetString("ssh_key_path", "")
		accessToken := request.GetString("access_token", "")
		configPath := request.GetString("config_path", "")
		outputFormat := request.GetString("output_format", "json")

		// Remote repository scan
		output, err := module.RemoteRepositoryScan(ctx, repoURL, repoType, iacType, branch, sshKeyPath, accessToken, configPath, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan remote repository scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan custom policy management tool
	customPolicyManagementTool := mcp.NewTool("terrascan_custom_policy_management",
		mcp.WithDescription("Manage and validate custom Terrascan policies"),
		mcp.WithString("action",
			mcp.Description("Policy management action"),
			mcp.Enum("validate", "test", "list", "scan-with-custom"),
			mcp.Required(),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to custom policy file or directory"),
		),
		mcp.WithString("target",
			mcp.Description("Target to scan (when action=scan-with-custom)"),
		),
		mcp.WithString("iac_type",
			mcp.Description("Infrastructure as Code type for policy validation"),
			mcp.Enum("terraform", "k8s", "helm", "cfn"),
		),
		mcp.WithString("test_data_path",
			mcp.Description("Path to test data for policy validation"),
		),
	)
	s.AddTool(customPolicyManagementTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		action := request.GetString("action", "")
		policyPath := request.GetString("policy_path", "")
		target := request.GetString("target", "")
		iacType := request.GetString("iac_type", "terraform")
		testDataPath := request.GetString("test_data_path", "")

		// Custom policy management
		output, err := module.CustomPolicyManagement(ctx, action, policyPath, target, iacType, testDataPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan custom policy management failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan CI/CD pipeline integration tool
	cicdPipelineIntegrationTool := mcp.NewTool("terrascan_cicd_pipeline_integration",
		mcp.WithDescription("Optimized IaC security scanning for CI/CD pipelines"),
		mcp.WithString("target",
			mcp.Description("Target for CI/CD scanning"),
			mcp.Required(),
		),
		mcp.WithString("iac_type",
			mcp.Description("Infrastructure as Code type"),
			mcp.Enum("terraform", "k8s", "helm", "kustomize"),
			mcp.Required(),
		),
		mcp.WithString("pipeline_stage",
			mcp.Description("CI/CD pipeline stage"),
			mcp.Enum("pre-commit", "build", "test", "staging", "production"),
		),
		mcp.WithString("gate_policy",
			mcp.Description("Quality gate policy for pipeline"),
			mcp.Enum("strict", "standard", "permissive"),
		),
		mcp.WithString("output_format",
			mcp.Description("CI-friendly output format"),
			mcp.Enum("json", "sarif", "junit-xml", "csv"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for CI artifacts"),
		),
		mcp.WithBoolean("fail_on_violations",
			mcp.Description("Fail CI pipeline on policy violations"),
		),
		mcp.WithString("baseline_file",
			mcp.Description("Baseline file for incremental scanning"),
		),
		mcp.WithBoolean("quiet_mode",
			mcp.Description("Suppress progress output for CI logs"),
		),
	)
	s.AddTool(cicdPipelineIntegrationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		target := request.GetString("target", "")
		iacType := request.GetString("iac_type", "")
		pipelineStage := request.GetString("pipeline_stage", "build")
		gatePolicy := request.GetString("gate_policy", "standard")
		outputFormat := request.GetString("output_format", "json")
		outputFile := request.GetString("output_file", "")
		failOnViolations := request.GetBool("fail_on_violations", false)
		baselineFile := request.GetString("baseline_file", "")
		quietMode := request.GetBool("quiet_mode", false)

		// CI/CD pipeline integration
		output, err := module.CICDPipelineIntegration(ctx, target, iacType, pipelineStage, gatePolicy, outputFormat, outputFile, failOnViolations, baselineFile, quietMode)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan CI/CD integration failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan cloud provider specific scanning tool
	cloudProviderScanTool := mcp.NewTool("terrascan_cloud_provider_scan",
		mcp.WithDescription("Cloud provider specific security scanning"),
		mcp.WithString("target",
			mcp.Description("Target to scan"),
			mcp.Required(),
		),
		mcp.WithString("cloud_provider",
			mcp.Description("Cloud provider to focus scanning on"),
			mcp.Enum("aws", "azure", "gcp", "kubernetes", "github"),
			mcp.Required(),
		),
		mcp.WithString("iac_type",
			mcp.Description("Infrastructure as Code type"),
			mcp.Enum("terraform", "cfn", "arm", "k8s", "helm"),
		),
		mcp.WithString("security_categories",
			mcp.Description("Comma-separated security categories to focus on"),
		),
		mcp.WithString("service_focus",
			mcp.Description("Specific cloud services to focus scanning on"),
		),
		mcp.WithBoolean("include_best_practices",
			mcp.Description("Include cloud provider best practices checks"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for results"),
			mcp.Enum("json", "yaml", "sarif", "html"),
		),
	)
	s.AddTool(cloudProviderScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		target := request.GetString("target", "")
		cloudProvider := request.GetString("cloud_provider", "")
		iacType := request.GetString("iac_type", "terraform")
		securityCategories := request.GetString("security_categories", "")
		serviceFocus := request.GetString("service_focus", "")
		includeBestPractices := request.GetBool("include_best_practices", false)
		outputFormat := request.GetString("output_format", "json")

		// Cloud provider scan
		output, err := module.CloudProviderScan(ctx, target, cloudProvider, iacType, securityCategories, serviceFocus, includeBestPractices, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan cloud provider scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan performance optimization tool
	performanceOptimizationTool := mcp.NewTool("terrascan_performance_optimization",
		mcp.WithDescription("High-performance IaC scanning with optimization features"),
		mcp.WithString("target",
			mcp.Description("Target for optimized scanning"),
			mcp.Required(),
		),
		mcp.WithString("iac_type",
			mcp.Description("Infrastructure as Code type"),
			mcp.Enum("terraform", "k8s", "helm", "cfn"),
			mcp.Required(),
		),
		mcp.WithString("scan_mode",
			mcp.Description("Scanning optimization mode"),
			mcp.Enum("fast", "thorough", "balanced"),
		),
		mcp.WithString("parallel_workers",
			mcp.Description("Number of parallel scanning workers"),
		),
		mcp.WithString("max_file_size",
			mcp.Description("Maximum file size to scan (e.g., 10MB)"),
		),
		mcp.WithBoolean("skip_large_files",
			mcp.Description("Skip scanning of large files"),
		),
		mcp.WithBoolean("enable_caching",
			mcp.Description("Enable scanning result caching"),
		),
		mcp.WithString("exclude_dirs",
			mcp.Description("Comma-separated directories to exclude"),
		),
		mcp.WithBoolean("enable_metrics",
			mcp.Description("Enable performance metrics collection"),
		),
	)
	s.AddTool(performanceOptimizationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		target := request.GetString("target", "")
		iacType := request.GetString("iac_type", "")
		scanMode := request.GetString("scan_mode", "balanced")
		parallelWorkers := request.GetString("parallel_workers", "")
		maxFileSize := request.GetString("max_file_size", "")
		skipLargeFiles := request.GetBool("skip_large_files", false)
		enableCaching := request.GetBool("enable_caching", false)
		excludeDirs := request.GetString("exclude_dirs", "")
		enableMetrics := request.GetBool("enable_metrics", false)

		// Performance optimization scan
		output, err := module.PerformanceOptimization(ctx, target, iacType, scanMode, parallelWorkers, maxFileSize, skipLargeFiles, enableCaching, excludeDirs, enableMetrics)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan performance optimization failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan comprehensive reporting tool
	comprehensiveReportingTool := mcp.NewTool("terrascan_comprehensive_reporting",
		mcp.WithDescription("Generate comprehensive IaC security reports with analytics"),
		mcp.WithString("target",
			mcp.Description("Target for comprehensive analysis"),
			mcp.Required(),
		),
		mcp.WithString("iac_type",
			mcp.Description("Infrastructure as Code type"),
			mcp.Enum("terraform", "k8s", "helm", "cfn", "arm"),
			mcp.Required(),
		),
		mcp.WithString("report_type",
			mcp.Description("Type of comprehensive report"),
			mcp.Enum("executive-summary", "technical-detail", "compliance-audit", "risk-assessment"),
			mcp.Required(),
		),
		mcp.WithString("output_formats",
			mcp.Description("Comma-separated output formats"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Directory for all report outputs"),
		),
		mcp.WithBoolean("include_remediation",
			mcp.Description("Include remediation guidance in reports"),
		),
		mcp.WithBoolean("include_trends",
			mcp.Description("Include security trend analysis"),
		),
		mcp.WithString("baseline_comparison",
			mcp.Description("Path to baseline scan for comparison"),
		),
		mcp.WithBoolean("include_policy_details",
			mcp.Description("Include detailed policy information"),
		),
	)
	s.AddTool(comprehensiveReportingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get parameters
		target := request.GetString("target", "")
		iacType := request.GetString("iac_type", "")
		reportType := request.GetString("report_type", "")
		outputFormats := request.GetString("output_formats", "")
		outputDirectory := request.GetString("output_directory", "")
		includeRemediation := request.GetBool("include_remediation", false)
		includeTrends := request.GetBool("include_trends", false)
		baselineComparison := request.GetString("baseline_comparison", "")
		includePolicyDetails := request.GetBool("include_policy_details", false)

		// Comprehensive reporting
		output, err := module.ComprehensiveReporting(ctx, target, iacType, reportType, outputFormats, outputDirectory, includeRemediation, includeTrends, baselineComparison, includePolicyDetails)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan comprehensive reporting failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terrascan get version tool
	getVersionTool := mcp.NewTool("terrascan_get_version",
		mcp.WithDescription("Get Terrascan version and supported IaC types"),
		mcp.WithBoolean("show_supported_types",
			mcp.Description("Include supported IaC types information"),
		),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerrascanModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terrascan get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}