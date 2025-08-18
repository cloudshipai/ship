package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSemgrepTools adds Semgrep (advanced static analysis for code security) MCP tool implementations using real semgrep CLI commands
func AddSemgrepTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		args := []string{"semgrep", "scan"}
		
		target := request.GetString("target", ".")
		args = append(args, target)
		
		ruleset := request.GetString("ruleset", "p/security-audit")
		args = append(args, "--config", ruleset)
		
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		if excludePaths := request.GetString("exclude_paths", ""); excludePaths != "" {
			for _, path := range strings.Split(excludePaths, ",") {
				if strings.TrimSpace(path) != "" {
					args = append(args, "--exclude", strings.TrimSpace(path))
				}
			}
		}
		if request.GetBool("verbose", false) {
			args = append(args, "--verbose")
		}
		if request.GetBool("fail_on_findings", false) {
			args = append(args, "--error")
		}
		
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		language := request.GetString("language", "")
		args := []string{"semgrep", "scan", target, "--config", "p/" + language}
		
		if securityCategory := request.GetString("security_category", ""); securityCategory != "" {
			args = append(args, "--config", "p/" + securityCategory)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if request.GetBool("include_experimental", false) {
			args = append(args, "--config", "p/experimental")
		}
		if confidence := request.GetString("confidence", ""); confidence != "" {
			args = append(args, "--confidence", confidence)
		}
		
		return executeShipCommand(args)
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
		args := []string{"semgrep", "scan"}
		
		target := request.GetString("target", ".")
		args = append(args, target)
		
		configPolicy := request.GetString("config_policy", "p/ci")
		args = append(args, "--config", configPolicy)
		
		if baselineRef := request.GetString("baseline_ref", ""); baselineRef != "" {
			args = append(args, "--baseline-ref", baselineRef)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		if request.GetBool("diff_aware", false) {
			args = append(args, "--diff-depth", "1")
		}
		if request.GetBool("fail_open", false) {
			args = append(args, "--disable-version-check")
		}
		if timeout := request.GetString("timeout", ""); timeout != "" {
			args = append(args, "--timeout", timeout)
		}
		if request.GetBool("quiet", false) {
			args = append(args, "--quiet")
		}
		
		return executeShipCommand(args)
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
		action := request.GetString("action", "")
		rulesPath := request.GetString("rules_path", "")
		
		switch action {
		case "validate":
			args := []string{"semgrep", "validate", "--config", rulesPath}
			if request.GetBool("strict", false) {
				args = append(args, "--strict")
			}
			return executeShipCommand(args)
		case "test":
			args := []string{"semgrep", "test", "--config", rulesPath}
			if testFiles := request.GetString("test_files", ""); testFiles != "" {
				for _, file := range strings.Split(testFiles, ",") {
					if strings.TrimSpace(file) != "" {
						args = append(args, strings.TrimSpace(file))
					}
				}
			}
			return executeShipCommand(args)
		case "scan":
			target := request.GetString("target", ".")
			args := []string{"semgrep", "scan", target, "--config", rulesPath}
			if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
				args = append(args, "--output", outputFormat)
			}
			return executeShipCommand(args)
		default:
			args := []string{"semgrep", "--help"}
			return executeShipCommand(args)
		}
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
		target := request.GetString("target", "")
		args := []string{"semgrep", "scan", target}
		
		configPolicy := request.GetString("config_policy", "auto")
		args = append(args, "--config", configPolicy)
		
		if maxMemory := request.GetString("max_memory", ""); maxMemory != "" {
			args = append(args, "--max-memory", maxMemory)
		}
		if maxTargetBytes := request.GetString("max_target_bytes", ""); maxTargetBytes != "" {
			args = append(args, "--max-target-bytes", maxTargetBytes)
		}
		if jobs := request.GetString("jobs", ""); jobs != "" {
			args = append(args, "--jobs", jobs)
		}
		if timeout := request.GetString("timeout", ""); timeout != "" {
			args = append(args, "--timeout", timeout)
		}
		if request.GetBool("enable_metrics", false) {
			args = append(args, "--metrics")
		}
		if request.GetBool("optimizations", false) {
			args = append(args, "--optimizations")
		}
		if excludePatterns := request.GetString("exclude_patterns", ""); excludePatterns != "" {
			for _, pattern := range strings.Split(excludePatterns, ",") {
				if strings.TrimSpace(pattern) != "" {
					args = append(args, "--exclude", strings.TrimSpace(pattern))
				}
			}
		}
		
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		args := []string{"semgrep", "scan", directory, "--config", "p/secrets"}
		
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		
		if excludePatterns := request.GetString("exclude_patterns", ""); excludePatterns != "" {
			for _, pattern := range strings.Split(excludePatterns, ",") {
				if strings.TrimSpace(pattern) != "" {
					args = append(args, "--exclude", strings.TrimSpace(pattern))
				}
			}
		}
		
		return executeShipCommand(args)
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
		directory := request.GetString("directory", "")
		args := []string{"semgrep", "scan", directory, "--config", "p/owasp-top-ten"}
		
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		
		if language := request.GetString("language_focus", ""); language != "" {
			args = append(args, "--config", "p/"+language)
		}
		
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		researchMode := request.GetString("research_mode", "")
		args := []string{"semgrep", "scan", target}
		
		// Configure research-specific rulesets
		switch researchMode {
		case "cve-analysis":
			args = append(args, "--config", "p/cwe-top-25", "--config", "p/owasp-top-ten")
		case "pattern-discovery":
			args = append(args, "--config", "p/security-audit", "--config", "p/experimental")
		case "exploit-detection":
			args = append(args, "--config", "p/security-audit", "--config", "p/insecure-transport")
		case "zero-day-hunting":
			args = append(args, "--config", "p/r2c-security-audit", "--config", "p/experimental")
		}
		
		if languageFocus := request.GetString("language_focus", ""); languageFocus != "" {
			args = append(args, "--config", "p/"+languageFocus)
		}
		if vulnTypes := request.GetString("vulnerability_types", ""); vulnTypes != "" {
			for _, vulnType := range strings.Split(vulnTypes, ",") {
				if strings.TrimSpace(vulnType) != "" {
					args = append(args, "--config", "p/"+strings.TrimSpace(vulnType))
				}
			}
		}
		if request.GetBool("include_experimental", false) {
			args = append(args, "--config", "p/experimental")
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		complianceFramework := request.GetString("compliance_framework", "")
		args := []string{"semgrep", "scan", target, "--config", "p/"+complianceFramework}
		
		if industryFocus := request.GetString("industry_focus", ""); industryFocus != "" {
			args = append(args, "--config", "p/"+industryFocus)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		if severityThreshold := request.GetString("severity_threshold", ""); severityThreshold != "" {
			args = append(args, "--severity", severityThreshold)
		}
		if request.GetBool("include_remediation", false) {
			args = append(args, "--sarif")
		}
		
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		reportType := request.GetString("report_type", "")
		args := []string{"semgrep", "scan", target}
		
		// Configure report-specific rulesets
		switch reportType {
		case "executive-summary":
			args = append(args, "--config", "p/security-audit", "--config", "p/owasp-top-ten")
		case "detailed-technical":
			args = append(args, "--config", "p/security-audit", "--config", "p/cwe-top-25", "--config", "p/experimental")
		case "developer-focused":
			args = append(args, "--config", "p/security-audit", "--config", "p/secrets")
		case "compliance-audit":
			args = append(args, "--config", "p/owasp-top-ten", "--config", "p/cwe-top-25")
		}
		
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "--output", outputFormats)
		} else {
			args = append(args, "--output", "json")
		}
		
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			args = append(args, "--output-file", outputDirectory+"/semgrep-report.json")
		}
		if request.GetBool("include_metrics", false) {
			args = append(args, "--metrics")
		}
		if baselineComparison := request.GetString("baseline_comparison", ""); baselineComparison != "" {
			args = append(args, "--baseline-ref", baselineComparison)
		}
		
		return executeShipCommand(args)
	})

	// Semgrep get version tool
	getVersionTool := mcp.NewTool("semgrep_get_version",
		mcp.WithDescription("Get Semgrep version and configuration information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"semgrep", "--version"}
		return executeShipCommand(args)
	})
}