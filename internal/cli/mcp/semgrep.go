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

// AddSemgrepTools adds Semgrep (advanced static analysis for code security) MCP tool implementations using direct Dagger calls
func AddSemgrepTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addSemgrepToolsDirect(s)
}

// addSemgrepToolsDirect adds Semgrep tools using direct Dagger module calls
func addSemgrepToolsDirect(s *server.MCPServer) {
	// Semgrep security audit scan tool
	securityAuditScanTool := mcp.NewTool("semgrep_security_audit_scan",
		mcp.WithDescription("Comprehensive security audit scan using Semgrep"),
		mcp.WithString("target",
			mcp.Description("Directory or file to scan (default: current directory)"),
		),
		mcp.WithString("ruleset",
			mcp.Description("Security ruleset to use"),
			mcp.Enum("auto", "p/security-audit", "p/owasp-top-ten", "p/cwe-top-25", "p/r2c-security-audit"),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level"),
			mcp.Enum("INFO", "WARNING", "ERROR"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text", "sarif", "junit-xml", "gitlab-sast", "csv"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("exclude_paths",
			mcp.Description("Comma-separated paths to exclude"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose output"),
		),
		mcp.WithBoolean("fail_on_findings",
			mcp.Description("Exit with error code on findings"),
		),
	)
	s.AddTool(securityAuditScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get parameters
		target := request.GetString("target", ".")
		ruleset := request.GetString("ruleset", "p/security-audit")
		severity := request.GetString("severity", "")
		outputFormat := request.GetString("output_format", "")
		outputFile := request.GetString("output_file", "")
		excludePathsStr := request.GetString("exclude_paths", "")
		verbose := request.GetBool("verbose", false)
		failOnFindings := request.GetBool("fail_on_findings", false)

		// Parse exclude paths
		var excludePaths []string
		if excludePathsStr != "" {
			for _, path := range strings.Split(excludePathsStr, ",") {
				if trimmed := strings.TrimSpace(path); trimmed != "" {
					excludePaths = append(excludePaths, trimmed)
				}
			}
		}

		// Run security audit scan
		output, err := module.SecurityAuditScan(ctx, target, ruleset, severity, outputFormat, outputFile, excludePaths, verbose, failOnFindings)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep security audit scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep language-specific security scan tool
	languageSpecificScanTool := mcp.NewTool("semgrep_language_specific_scan",
		mcp.WithDescription("Language-specific security analysis with Semgrep"),
		mcp.WithString("target",
			mcp.Description("Directory or file to scan"),
			mcp.Required(),
		),
		mcp.WithString("language",
			mcp.Description("Programming language to focus on"),
			mcp.Enum("python", "javascript", "typescript", "java", "go", "ruby", "php", "c", "cpp", "csharp", "rust"),
			mcp.Required(),
		),
		mcp.WithString("security_category",
			mcp.Description("Security category to focus on"),
			mcp.Enum("security", "owasp-top-ten", "cwe-top-25", "injection", "crypto", "xss", "sqli"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text", "sarif", "junit-xml", "csv"),
		),
		mcp.WithBoolean("include_experimental",
			mcp.Description("Include experimental rules"),
		),
		mcp.WithString("confidence",
			mcp.Description("Minimum confidence level"),
			mcp.Enum("low", "medium", "high"),
		),
	)
	s.AddTool(languageSpecificScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get parameters
		target := request.GetString("target", "")
		language := request.GetString("language", "")
		securityCategory := request.GetString("security_category", "")
		outputFormat := request.GetString("output_format", "")
		includeExperimental := request.GetBool("include_experimental", false)
		confidence := request.GetString("confidence", "")

		// Run language-specific scan
		output, err := module.LanguageSpecificScan(ctx, target, language, securityCategory, outputFormat, includeExperimental, confidence)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep language-specific scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep CI/CD integration scan tool
	cicdIntegrationScanTool := mcp.NewTool("semgrep_cicd_integration_scan",
		mcp.WithDescription("Optimized Semgrep scan for CI/CD pipelines"),
		mcp.WithString("target",
			mcp.Description("Directory or file to scan (default: current directory)"),
		),
		mcp.WithString("baseline_ref",
			mcp.Description("Git baseline reference for differential scanning"),
		),
		mcp.WithString("output_format",
			mcp.Description("CI/CD friendly output format"),
			mcp.Enum("json", "sarif", "junit-xml", "gitlab-sast", "github-actions"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for CI artifacts"),
		),
		mcp.WithString("config_policy",
			mcp.Description("Configuration policy for CI scans"),
			mcp.Enum("auto", "p/ci", "p/security-audit", "p/owasp-top-ten"),
		),
		mcp.WithBoolean("diff_aware",
			mcp.Description("Enable diff-aware scanning for faster CI runs"),
		),
		mcp.WithBoolean("fail_open",
			mcp.Description("Succeed if Semgrep fails to run"),
		),
		mcp.WithString("timeout",
			mcp.Description("Timeout in seconds for CI environments"),
		),
		mcp.WithBoolean("quiet",
			mcp.Description("Suppress progress output for CI"),
		),
	)
	s.AddTool(cicdIntegrationScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get parameters
		target := request.GetString("target", ".")
		baselineRef := request.GetString("baseline_ref", "")
		outputFormat := request.GetString("output_format", "")
		outputFile := request.GetString("output_file", "")
		configPolicy := request.GetString("config_policy", "p/ci")
		diffAware := request.GetBool("diff_aware", false)
		failOpen := request.GetBool("fail_open", false)
		timeout := request.GetString("timeout", "")
		quiet := request.GetBool("quiet", false)

		// Run CI/CD integration scan
		output, err := module.CICDIntegrationScan(ctx, target, baselineRef, outputFormat, outputFile, configPolicy, diffAware, failOpen, timeout, quiet)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep CI/CD scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep custom rule management tool
	customRuleManagementTool := mcp.NewTool("semgrep_custom_rule_management",
		mcp.WithDescription("Manage and validate custom Semgrep rules"),
		mcp.WithString("action",
			mcp.Description("Rule management action"),
			mcp.Enum("validate", "test", "lint", "scan"),
			mcp.Required(),
		),
		mcp.WithString("rules_path",
			mcp.Description("Path to custom rules file or directory"),
			mcp.Required(),
		),
		mcp.WithString("target",
			mcp.Description("Target directory or file for scanning (when action=scan)"),
		),
		mcp.WithString("test_files",
			mcp.Description("Comma-separated test files for rule validation"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for results"),
			mcp.Enum("json", "text", "sarif"),
		),
		mcp.WithBoolean("strict",
			mcp.Description("Enable strict validation mode"),
		),
	)
	s.AddTool(customRuleManagementTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get parameters
		action := request.GetString("action", "")
		rulesPath := request.GetString("rules_path", "")
		target := request.GetString("target", ".")
		testFilesStr := request.GetString("test_files", "")
		outputFormat := request.GetString("output_format", "")
		strict := request.GetBool("strict", false)

		// Parse test files
		var testFiles []string
		if testFilesStr != "" {
			for _, file := range strings.Split(testFilesStr, ",") {
				if trimmed := strings.TrimSpace(file); trimmed != "" {
					testFiles = append(testFiles, trimmed)
				}
			}
		}

		// Run custom rule management
		output, err := module.CustomRuleManagement(ctx, action, rulesPath, target, testFiles, outputFormat, strict)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep custom rule management failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep performance optimized scan tool
	performanceOptimizedScanTool := mcp.NewTool("semgrep_performance_optimized_scan",
		mcp.WithDescription("High-performance Semgrep scan with optimization features"),
		mcp.WithString("target",
			mcp.Description("Directory or file to scan"),
			mcp.Required(),
		),
		mcp.WithString("config_policy",
			mcp.Description("Configuration policy for optimized scanning"),
			mcp.Enum("auto", "p/security-audit", "p/owasp-top-ten", "p/r2c-security-audit"),
		),
		mcp.WithString("max_memory",
			mcp.Description("Maximum memory usage in MB"),
		),
		mcp.WithString("max_target_bytes",
			mcp.Description("Maximum file size to scan in bytes"),
		),
		mcp.WithString("jobs",
			mcp.Description("Number of parallel jobs"),
		),
		mcp.WithString("timeout",
			mcp.Description("Timeout per file in seconds"),
		),
		mcp.WithBoolean("enable_metrics",
			mcp.Description("Enable performance metrics collection"),
		),
		mcp.WithBoolean("optimizations",
			mcp.Description("Enable all performance optimizations"),
		),
		mcp.WithString("exclude_patterns",
			mcp.Description("Comma-separated patterns to exclude for performance"),
		),
	)
	s.AddTool(performanceOptimizedScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get parameters
		target := request.GetString("target", "")
		configPolicy := request.GetString("config_policy", "auto")
		maxMemory := request.GetString("max_memory", "")
		maxTargetBytes := request.GetString("max_target_bytes", "")
		jobs := request.GetString("jobs", "")
		timeout := request.GetString("timeout", "")
		enableMetrics := request.GetBool("enable_metrics", false)
		optimizations := request.GetBool("optimizations", false)
		excludePatternsStr := request.GetString("exclude_patterns", "")

		// Parse exclude patterns
		var excludePatterns []string
		if excludePatternsStr != "" {
			for _, pattern := range strings.Split(excludePatternsStr, ",") {
				if trimmed := strings.TrimSpace(pattern); trimmed != "" {
					excludePatterns = append(excludePatterns, trimmed)
				}
			}
		}

		// Run performance optimized scan
		output, err := module.PerformanceOptimizedScan(ctx, target, configPolicy, maxMemory, maxTargetBytes, jobs, timeout, enableMetrics, optimizations, excludePatterns)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep performance scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep secrets scanning tool
	scanSecretsTool := mcp.NewTool("semgrep_scan_secrets",
		mcp.WithDescription("Specialized secrets scanning using Semgrep"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan for secrets"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "sarif", "gitlab-secrets"),
		),
		mcp.WithString("exclude_patterns",
			mcp.Description("Comma-separated exclude patterns"),
		),
	)
	s.AddTool(scanSecretsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		outputFormat := request.GetString("output_format", "")
		excludePatternsStr := request.GetString("exclude_patterns", "")

		// Parse exclude patterns
		var excludePatterns []string
		if excludePatternsStr != "" {
			for _, pattern := range strings.Split(excludePatternsStr, ",") {
				if trimmed := strings.TrimSpace(pattern); trimmed != "" {
					excludePatterns = append(excludePatterns, trimmed)
				}
			}
		}

		// Run secrets scan
		output, err := module.ScanSecrets(ctx, directory, outputFormat, excludePatterns)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep secrets scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep OWASP Top 10 scan tool
	scanOWASPTool := mcp.NewTool("semgrep_scan_owasp_top10",
		mcp.WithDescription("Scan for OWASP Top 10 vulnerabilities using Semgrep"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "sarif", "text"),
		),
		mcp.WithString("language_focus",
			mcp.Description("Focus on specific language"),
			mcp.Enum("python", "javascript", "java", "go", "php", "ruby"),
		),
	)
	s.AddTool(scanOWASPTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		outputFormat := request.GetString("output_format", "")
		languageFocus := request.GetString("language_focus", "")

		// Run OWASP Top 10 scan
		output, err := module.ScanOWASPTop10(ctx, directory, outputFormat, languageFocus)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep OWASP scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep vulnerability research tool
	vulnerabilityResearchTool := mcp.NewTool("semgrep_vulnerability_research",
		mcp.WithDescription("Advanced vulnerability research and pattern discovery with Semgrep"),
		mcp.WithString("target",
			mcp.Description("Directory or file to analyze"),
			mcp.Required(),
		),
		mcp.WithString("research_mode",
			mcp.Description("Vulnerability research mode"),
			mcp.Enum("cve-analysis", "pattern-discovery", "exploit-detection", "zero-day-hunting"),
			mcp.Required(),
		),
		mcp.WithString("language_focus",
			mcp.Description("Language to focus research on"),
			mcp.Enum("python", "javascript", "java", "go", "php", "ruby", "c", "cpp"),
		),
		mcp.WithString("vulnerability_types",
			mcp.Description("Comma-separated vulnerability types to research"),
		),
		mcp.WithBoolean("include_experimental",
			mcp.Description("Include experimental detection rules"),
		),
		mcp.WithString("output_format",
			mcp.Description("Research output format"),
			mcp.Enum("json", "sarif", "text"),
		),
	)
	s.AddTool(vulnerabilityResearchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get parameters
		target := request.GetString("target", "")
		researchMode := request.GetString("research_mode", "")
		languageFocus := request.GetString("language_focus", "")
		vulnerabilityTypesStr := request.GetString("vulnerability_types", "")
		includeExperimental := request.GetBool("include_experimental", false)
		outputFormat := request.GetString("output_format", "")

		// Parse vulnerability types
		var vulnerabilityTypes []string
		if vulnerabilityTypesStr != "" {
			for _, vulnType := range strings.Split(vulnerabilityTypesStr, ",") {
				if trimmed := strings.TrimSpace(vulnType); trimmed != "" {
					vulnerabilityTypes = append(vulnerabilityTypes, trimmed)
				}
			}
		}

		// Run vulnerability research
		output, err := module.VulnerabilityResearch(ctx, target, researchMode, languageFocus, vulnerabilityTypes, includeExperimental, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep vulnerability research failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep compliance scanning tool
	complianceScanningTool := mcp.NewTool("semgrep_compliance_scanning",
		mcp.WithDescription("Compliance-focused security scanning with Semgrep"),
		mcp.WithString("target",
			mcp.Description("Directory or file to scan for compliance"),
			mcp.Required(),
		),
		mcp.WithString("compliance_framework",
			mcp.Description("Compliance framework to validate against"),
			mcp.Enum("owasp-top-ten", "cwe-top-25", "pci-dss", "nist", "iso27001", "hipaa"),
			mcp.Required(),
		),
		mcp.WithString("industry_focus",
			mcp.Description("Industry-specific compliance requirements"),
			mcp.Enum("fintech", "healthcare", "government", "retail", "automotive"),
		),
		mcp.WithString("output_format",
			mcp.Description("Compliance report format"),
			mcp.Enum("json", "sarif", "junit-xml", "csv"),
		),
		mcp.WithString("output_file",
			mcp.Description("Compliance report output file"),
		),
		mcp.WithBoolean("include_remediation",
			mcp.Description("Include remediation guidance in report"),
		),
		mcp.WithString("severity_threshold",
			mcp.Description("Minimum severity for compliance violations"),
			mcp.Enum("INFO", "WARNING", "ERROR"),
		),
	)
	s.AddTool(complianceScanningTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get parameters
		target := request.GetString("target", "")
		complianceFramework := request.GetString("compliance_framework", "")
		industryFocus := request.GetString("industry_focus", "")
		outputFormat := request.GetString("output_format", "")
		outputFile := request.GetString("output_file", "")
		includeRemediation := request.GetBool("include_remediation", false)
		severityThreshold := request.GetString("severity_threshold", "")

		// Run compliance scanning
		output, err := module.ComplianceScanning(ctx, target, complianceFramework, industryFocus, outputFormat, outputFile, includeRemediation, severityThreshold)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep compliance scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep comprehensive reporting tool
	comprehensiveReportingTool := mcp.NewTool("semgrep_comprehensive_reporting",
		mcp.WithDescription("Generate comprehensive security analysis reports with Semgrep"),
		mcp.WithString("target",
			mcp.Description("Directory or file to analyze"),
			mcp.Required(),
		),
		mcp.WithString("report_type",
			mcp.Description("Type of comprehensive report to generate"),
			mcp.Enum("executive-summary", "detailed-technical", "developer-focused", "compliance-audit"),
			mcp.Required(),
		),
		mcp.WithString("output_formats",
			mcp.Description("Comma-separated output formats"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Directory to save all report files"),
		),
		mcp.WithBoolean("include_metrics",
			mcp.Description("Include performance and coverage metrics"),
		),
		mcp.WithBoolean("include_trends",
			mcp.Description("Include security trend analysis"),
		),
		mcp.WithString("baseline_comparison",
			mcp.Description("Path to baseline report for comparison"),
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
		module := modules.NewSemgrepModule(client)

		// Get parameters
		target := request.GetString("target", "")
		reportType := request.GetString("report_type", "")
		outputFormatsStr := request.GetString("output_formats", "")
		outputDirectory := request.GetString("output_directory", "")
		includeMetrics := request.GetBool("include_metrics", false)
		includeTrends := request.GetBool("include_trends", false)
		baselineComparison := request.GetString("baseline_comparison", "")

		// Parse output formats
		var outputFormats []string
		if outputFormatsStr != "" {
			for _, format := range strings.Split(outputFormatsStr, ",") {
				if trimmed := strings.TrimSpace(format); trimmed != "" {
					outputFormats = append(outputFormats, trimmed)
				}
			}
		}

		// Run comprehensive reporting
		output, err := module.ComprehensiveReporting(ctx, target, reportType, outputFormats, outputDirectory, includeMetrics, includeTrends, baselineComparison)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep comprehensive reporting failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Semgrep get version tool
	getVersionTool := mcp.NewTool("semgrep_get_version",
		mcp.WithDescription("Get Semgrep version and configuration information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSemgrepModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("semgrep get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}