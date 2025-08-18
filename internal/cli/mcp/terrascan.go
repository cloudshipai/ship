package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTerrascanTools adds Terrascan (IaC security scanner) MCP tool implementations using real terrascan CLI commands
func AddTerrascanTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		directory := request.GetString("directory", "")
		args := []string{"terrascan", "scan", "-d", directory}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		args := []string{"terrascan", "scan", "-i", "terraform", "-d", directory}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		args := []string{"terrascan", "scan", "-i", "k8s", "-d", directory}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		severity := request.GetString("severity", "")
		args := []string{"terrascan", "scan", "-d", directory, "--severity", severity}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		repoURL := request.GetString("repo_url", "")
		args := []string{"terrascan", "scan", "-r", repoURL}
		if repoType := request.GetString("repo_type", ""); repoType != "" {
			args = append(args, "--remote-type", repoType)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		policyPath := request.GetString("policy_path", "")
		args := []string{"terrascan", "scan", "-d", directory, "--policy-path", policyPath}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		iacType := request.GetString("iac_type", "")
		args := []string{"terrascan", "scan", "-i", iacType, "-d", target}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "-o", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		if severityThreshold := request.GetString("severity_threshold", ""); severityThreshold != "" {
			args = append(args, "--severity", severityThreshold)
		}
		if policyTypes := request.GetString("policy_types", ""); policyTypes != "" {
			args = append(args, "--policy-type", policyTypes)
		}
		if excludeRules := request.GetString("exclude_rules", ""); excludeRules != "" {
			args = append(args, "--skip-rules", excludeRules)
		}
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		if request.GetBool("show_passed", false) {
			args = append(args, "--show-passed")
		}
		
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		complianceFramework := request.GetString("compliance_framework", "")
		iacType := request.GetString("iac_type", "terraform")
		args := []string{"terrascan", "scan", "-i", iacType, "-d", target, "--policy-type", complianceFramework}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "-o", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		if request.GetBool("include_severity_details", false) {
			args = append(args, "--show-passed", "-v")
		}
		
		return executeShipCommand(args)
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
		repoURL := request.GetString("repository_url", "")
		repoType := request.GetString("repository_type", "git")
		iacType := request.GetString("iac_type", "terraform")
		args := []string{"terrascan", "scan", "-r", repoURL, "-t", repoType, "-i", iacType}
		
		if branch := request.GetString("branch", ""); branch != "" {
			args = append(args, "--remote-branch", branch)
		}
		if sshKeyPath := request.GetString("ssh_key_path", ""); sshKeyPath != "" {
			args = append(args, "--ssh-key", sshKeyPath)
		}
		if accessToken := request.GetString("access_token", ""); accessToken != "" {
			args = append(args, "--access-token", accessToken)
		}
		if configPath := request.GetString("config_path", ""); configPath != "" {
			args = append(args, "-c", configPath)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "-o", outputFormat)
		}
		
		return executeShipCommand(args)
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
		action := request.GetString("action", "")
		
		switch action {
		case "validate":
			policyPath := request.GetString("policy_path", "")
			args := []string{"terrascan", "init", "--policy-path", policyPath}
			return executeShipCommand(args)
		case "test":
			policyPath := request.GetString("policy_path", "")
			testDataPath := request.GetString("test_data_path", "")
			args := []string{"terrascan", "scan", "--policy-path", policyPath, "-d", testDataPath}
			return executeShipCommand(args)
		case "list":
			args := []string{"terrascan", "--help"}
			return executeShipCommand(args)
		case "scan-with-custom":
			target := request.GetString("target", "")
			policyPath := request.GetString("policy_path", "")
			iacType := request.GetString("iac_type", "terraform")
			args := []string{"terrascan", "scan", "-i", iacType, "-d", target, "--policy-path", policyPath}
			return executeShipCommand(args)
		default:
			args := []string{"terrascan", "--help"}
			return executeShipCommand(args)
		}
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
		target := request.GetString("target", "")
		iacType := request.GetString("iac_type", "")
		args := []string{"terrascan", "scan", "-i", iacType, "-d", target}
		
		// Configure based on pipeline stage and gate policy
		pipelineStage := request.GetString("pipeline_stage", "build")
		gatePolicy := request.GetString("gate_policy", "standard")
		
		switch pipelineStage {
		case "pre-commit":
			if gatePolicy == "strict" {
				args = append(args, "--severity", "MEDIUM")
			} else {
				args = append(args, "--severity", "HIGH")
			}
		case "build":
			args = append(args, "--severity", "HIGH")
		case "test":
			args = append(args, "--severity", "MEDIUM")
		case "staging", "production":
			args = append(args, "--severity", "LOW")
		}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "-o", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		if baselineFile := request.GetString("baseline_file", ""); baselineFile != "" {
			args = append(args, "--baseline", baselineFile)
		}
		if request.GetBool("quiet_mode", false) {
			args = append(args, "--quiet")
		}
		
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		cloudProvider := request.GetString("cloud_provider", "")
		iacType := request.GetString("iac_type", "terraform")
		args := []string{"terrascan", "scan", "-i", iacType, "-d", target, "--policy-type", cloudProvider}
		
		if securityCategories := request.GetString("security_categories", ""); securityCategories != "" {
			args = append(args, "--categories", securityCategories)
		}
		if serviceFocus := request.GetString("service_focus", ""); serviceFocus != "" {
			args = append(args, "--services", serviceFocus)
		}
		if request.GetBool("include_best_practices", false) {
			args = append(args, "--include-best-practices")
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "-o", outputFormat)
		}
		
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		iacType := request.GetString("iac_type", "")
		args := []string{"terrascan", "scan", "-i", iacType, "-d", target}
		
		scanMode := request.GetString("scan_mode", "balanced")
		switch scanMode {
		case "fast":
			args = append(args, "--severity", "HIGH", "--skip-rules", "low-priority")
		case "thorough":
			args = append(args, "--severity", "LOW", "--show-passed")
		case "balanced":
			args = append(args, "--severity", "MEDIUM")
		}
		
		if parallelWorkers := request.GetString("parallel_workers", ""); parallelWorkers != "" {
			args = append(args, "--workers", parallelWorkers)
		}
		if maxFileSize := request.GetString("max_file_size", ""); maxFileSize != "" {
			args = append(args, "--max-file-size", maxFileSize)
		}
		if request.GetBool("skip_large_files", false) {
			args = append(args, "--skip-large-files")
		}
		if request.GetBool("enable_caching", false) {
			args = append(args, "--cache")
		}
		if excludeDirs := request.GetString("exclude_dirs", ""); excludeDirs != "" {
			args = append(args, "--exclude-dirs", excludeDirs)
		}
		if request.GetBool("enable_metrics", false) {
			args = append(args, "--metrics")
		}
		
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		iacType := request.GetString("iac_type", "")
		reportType := request.GetString("report_type", "")
		args := []string{"terrascan", "scan", "-i", iacType, "-d", target}
		
		// Configure report-specific settings
		switch reportType {
		case "executive-summary":
			args = append(args, "--severity", "HIGH", "-o", "html")
		case "technical-detail":
			args = append(args, "--severity", "LOW", "--show-passed", "-o", "json")
		case "compliance-audit":
			args = append(args, "--policy-type", "all", "-o", "sarif")
		case "risk-assessment":
			args = append(args, "--severity", "MEDIUM", "-o", "yaml")
		}
		
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "-o", outputFormats)
		}
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			args = append(args, "--output-file", outputDirectory+"/terrascan-report")
		}
		if request.GetBool("include_remediation", false) {
			args = append(args, "--include-remediation")
		}
		if request.GetBool("include_policy_details", false) {
			args = append(args, "-v")
		}
		if baselineComparison := request.GetString("baseline_comparison", ""); baselineComparison != "" {
			args = append(args, "--baseline", baselineComparison)
		}
		
		return executeShipCommand(args)
	})

	// Terrascan get version tool
	getVersionTool := mcp.NewTool("terrascan_get_version",
		mcp.WithDescription("Get Terrascan version and supported IaC types"),
		mcp.WithBoolean("show_supported_types",
			mcp.Description("Include supported IaC types information"),
		),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terrascan", "version"}
		if request.GetBool("show_supported_types", false) {
			args = append(args, "--list")
		}
		return executeShipCommand(args)
	})
}