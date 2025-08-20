package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// SemgrepModule runs Semgrep for static analysis
type SemgrepModule struct {
	client *dagger.Client
	name   string
}

// NewSemgrepModule creates a new Semgrep module
func NewSemgrepModule(client *dagger.Client) *SemgrepModule {
	return &SemgrepModule{
		client: client,
		name:   "semgrep",
	}
}

// ScanDirectory scans a directory with Semgrep rules
func (m *SemgrepModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(dir)).
		WithWorkdir("/src").
		WithExec([]string{
			"semgrep",
			"--config=auto",
			"--json",
			"--severity=ERROR",
			".",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run semgrep: no output received")
}

// ScanWithRuleset scans with specific ruleset
func (m *SemgrepModule) ScanWithRuleset(ctx context.Context, dir string, ruleset string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(dir)).
		WithWorkdir("/src").
		WithExec([]string{
			"semgrep",
			"--config", ruleset,
			"--json",
			".",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run semgrep with ruleset: no output received")
}

// ScanFile scans a specific file
func (m *SemgrepModule) ScanFile(ctx context.Context, filePath string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithFile("/workspace/target.file", m.client.Host().File(filePath)).
		WithWorkdir("/src").
		WithExec([]string{
			"semgrep",
			"--config=auto",
			"--json",
			"target.file",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run semgrep on file: no output received")
}

// GetVersion returns the version of Semgrep
func (m *SemgrepModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithExec([]string{"semgrep", "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get semgrep version: %w", err)
	}

	return output, nil
}

// LanguageSpecificScan performs language-specific security analysis
func (m *SemgrepModule) LanguageSpecificScan(ctx context.Context, target string, language string, securityCategory string, outputFormat string, includeExperimental bool, confidence string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(target))

	args := []string{"semgrep", "scan", "/src", "--config", "p/" + language}
	if securityCategory != "" {
		args = append(args, "--config", "p/" + securityCategory)
	}
	if outputFormat == "json" {
		args = append(args, "--json")
	}
	if includeExperimental {
		args = append(args, "--config", "p/experimental")
	}
	if confidence != "" {
		args = append(args, "--confidence", confidence)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr for more info
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return "", fmt.Errorf("language-specific scan failed: %s", stderr)
		}
		return "", fmt.Errorf("language-specific scan failed: %w", err)
	}

	// Empty output might mean no findings, which is valid
	if output == "" {
		return "No security issues found", nil
	}
	return output, nil
}

// CICDIntegrationScan performs optimized scan for CI/CD pipelines
func (m *SemgrepModule) CICDIntegrationScan(ctx context.Context, target string, baselineRef string, outputFormat string, outputFile string, configPolicy string, diffAware bool, failOpen bool, timeout string, quiet bool) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(target))

	if configPolicy == "" {
		configPolicy = "p/ci"
	}

	args := []string{"semgrep", "scan", "/src", "--config", configPolicy}
	if baselineRef != "" {
		args = append(args, "--baseline-ref", baselineRef)
	}
	if outputFormat == "json" {
		args = append(args, "--json")
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	if diffAware {
		args = append(args, "--diff-depth", "1")
	}
	if failOpen {
		args = append(args, "--disable-version-check")
	}
	if timeout != "" {
		args = append(args, "--timeout", timeout)
	}
	if quiet {
		args = append(args, "--quiet")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run CI/CD scan: no output received")
}

// CustomRuleManagement manages and validates custom Semgrep rules
func (m *SemgrepModule) CustomRuleManagement(ctx context.Context, action string, rulesPath string, target string, testFiles []string, outputFormat string, strict bool) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest")

	if rulesPath != "" {
		container = container.WithFile("/rules.yaml", m.client.Host().File(rulesPath))
	}
	if target != "" {
		container = container.WithDirectory("/src", m.client.Host().Directory(target))
	}

	var args []string
	switch action {
	case "validate":
		args = []string{"semgrep", "validate", "--config", "/rules.yaml"}
		if strict {
			args = append(args, "--strict")
		}
	case "test":
		args = []string{"semgrep", "test", "--config", "/rules.yaml"}
		for _, file := range testFiles {
			container = container.WithFile("/test_"+file, m.client.Host().File(file))
			args = append(args, "/test_"+file)
		}
	case "scan":
		args = []string{"semgrep", "scan", "/src", "--config", "/rules.yaml"}
		if outputFormat == "json" {
			args = append(args, "--json")
		}
	default:
		args = []string{"semgrep", "--help"}
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr for more info
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return "", fmt.Errorf("custom rule management failed: %s", stderr)
		}
		return "", fmt.Errorf("custom rule management failed: %w", err)
	}

	// For test/validate actions, output might be minimal
	if output == "" && (action == "test" || action == "validate") {
		return "Rules validated successfully", nil
	}
	if output == "" {
		return "No findings", nil
	}
	return output, nil
}

// PerformanceOptimizedScan performs high-performance scan with optimization features
func (m *SemgrepModule) PerformanceOptimizedScan(ctx context.Context, target string, configPolicy string, maxMemory string, maxTargetBytes string, jobs string, timeout string, enableMetrics bool, optimizations bool, excludePatterns []string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(target))

	if configPolicy == "" {
		configPolicy = "auto"
	}

	args := []string{"semgrep", "scan", "/src", "--config", configPolicy}
	if maxMemory != "" {
		args = append(args, "--max-memory", maxMemory)
	}
	if maxTargetBytes != "" {
		args = append(args, "--max-target-bytes", maxTargetBytes)
	}
	if jobs != "" {
		args = append(args, "--jobs", jobs)
	}
	if timeout != "" {
		args = append(args, "--timeout", timeout)
	}
	if enableMetrics {
		args = append(args, "--metrics")
	}
	if optimizations {
		args = append(args, "--optimizations")
	}
	for _, pattern := range excludePatterns {
		args = append(args, "--exclude", pattern)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr for more info
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return "", fmt.Errorf("performance optimized scan failed: %s", stderr)
		}
		return "", fmt.Errorf("performance optimized scan failed: %w", err)
	}

	// Empty output might mean no findings, which is valid
	if output == "" {
		return "No security issues found (optimized scan)", nil
	}
	return output, nil
}

// ScanSecrets performs specialized secrets scanning
func (m *SemgrepModule) ScanSecrets(ctx context.Context, directory string, outputFormat string, excludePatterns []string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(directory))

	args := []string{"semgrep", "scan", "/src", "--config", "p/secrets"}
	if outputFormat == "json" {
		args = append(args, "--json")
	}
	for _, pattern := range excludePatterns {
		args = append(args, "--exclude", pattern)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run secrets scan: no output received")
}

// ScanOWASPTop10 scans for OWASP Top 10 vulnerabilities
func (m *SemgrepModule) ScanOWASPTop10(ctx context.Context, directory string, outputFormat string, languageFocus string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(directory))

	args := []string{"semgrep", "scan", "/src", "--config", "p/owasp-top-ten"}
	if outputFormat == "json" {
		args = append(args, "--json")
	}
	if languageFocus != "" {
		args = append(args, "--config", "p/"+languageFocus)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run OWASP Top 10 scan: no output received")
}

// VulnerabilityResearch performs advanced vulnerability research and pattern discovery
func (m *SemgrepModule) VulnerabilityResearch(ctx context.Context, target string, researchMode string, languageFocus string, vulnerabilityTypes []string, includeExperimental bool, outputFormat string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(target))

	args := []string{"semgrep", "scan", "/src"}

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

	if languageFocus != "" {
		args = append(args, "--config", "p/"+languageFocus)
	}
	for _, vulnType := range vulnerabilityTypes {
		args = append(args, "--config", "p/"+vulnType)
	}
	if includeExperimental {
		args = append(args, "--config", "p/experimental")
	}
	if outputFormat == "json" {
		args = append(args, "--json")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run vulnerability research: no output received")
}

// ComplianceScanning performs compliance-focused security scanning
func (m *SemgrepModule) ComplianceScanning(ctx context.Context, target string, complianceFramework string, industryFocus string, outputFormat string, outputFile string, includeRemediation bool, severityThreshold string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(target))

	args := []string{"semgrep", "scan", "/src", "--config", "p/"+complianceFramework}
	if industryFocus != "" {
		args = append(args, "--config", "p/"+industryFocus)
	}
	if outputFormat == "json" {
		args = append(args, "--json")
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	if severityThreshold != "" {
		args = append(args, "--severity", severityThreshold)
	}
	if includeRemediation {
		args = append(args, "--sarif")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr for more info
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return "", fmt.Errorf("compliance scanning failed: %s", stderr)
		}
		return "", fmt.Errorf("compliance scanning failed: %w", err)
	}

	// Empty output might mean no compliance issues, which is valid
	if output == "" {
		return "No compliance issues found", nil
	}
	return output, nil
}

// ComprehensiveReporting generates comprehensive security analysis reports
func (m *SemgrepModule) ComprehensiveReporting(ctx context.Context, target string, reportType string, outputFormats []string, outputDirectory string, includeMetrics bool, includeTrends bool, baselineComparison string) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(target))

	args := []string{"semgrep", "scan", "/src"}

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

	if len(outputFormats) > 0 && outputFormats[0] == "json" {
		args = append(args, "--json")
	} else if len(outputFormats) == 0 {
		args = append(args, "--json")
	}

	if outputDirectory != "" {
		args = append(args, "--output-file", outputDirectory+"/semgrep-report.json")
	}
	if includeMetrics {
		args = append(args, "--metrics")
	}
	if baselineComparison != "" {
		args = append(args, "--baseline-ref", baselineComparison)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run comprehensive reporting: no output received")
}

// SecurityAuditScan performs comprehensive security audit scan
func (m *SemgrepModule) SecurityAuditScan(ctx context.Context, target string, ruleset string, severity string, outputFormat string, outputFile string, excludePaths []string, verbose bool, failOnFindings bool) (string, error) {
	container := m.client.Container().
		From("semgrep/semgrep:latest").
		WithDirectory("/src", m.client.Host().Directory(target))

	if ruleset == "" {
		ruleset = "p/security-audit"
	}

	args := []string{"semgrep", "scan", "/src", "--config", ruleset}
	if severity != "" {
		args = append(args, "--severity", severity)
	}
	if outputFormat == "json" {
		args = append(args, "--json")
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	for _, path := range excludePaths {
		args = append(args, "--exclude", path)
	}
	if verbose {
		args = append(args, "--verbose")
	}
	if failOnFindings {
		args = append(args, "--error")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr for more info
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return "", fmt.Errorf("security audit scan failed: %s", stderr)
		}
		return "", fmt.Errorf("security audit scan failed: %w", err)
	}

	// Empty output might mean no security issues, which is valid
	if output == "" {
		return "No security issues found in audit", nil
	}
	return output, nil
}
