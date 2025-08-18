package modules

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// TruffleHogModule runs TruffleHog for verified secret detection
type TruffleHogModule struct {
	client *dagger.Client
	name   string
}

const trufflehogBinary = "/usr/local/bin/trufflehog"

// NewTruffleHogModule creates a new TruffleHog module
func NewTruffleHogModule(client *dagger.Client) *TruffleHogModule {
	return &TruffleHogModule{
		client: client,
		name:   trufflehogBinary,
	}
}

// ScanDirectory scans a directory for secrets using TruffleHog
func (m *TruffleHogModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{trufflehogBinary, "filesystem", ".", "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// TruffleHog returns non-zero exit code when secrets are found
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// ScanGitRepo scans a Git repository for secrets
func (m *TruffleHogModule) ScanGitRepo(ctx context.Context, repoURL string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithExec([]string{trufflehogBinary, "git", repoURL, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan git repo: %w", err)
	}

	return output, nil
}

// ScanGitHub scans a GitHub repository for secrets
func (m *TruffleHogModule) ScanGitHub(ctx context.Context, repo string, token string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	if token != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", token)
	} else if os.Getenv("GITHUB_TOKEN") != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	}

	container = container.WithExec([]string{trufflehogBinary, "github", "--repo", repo, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan github repo: %w", err)
	}

	return output, nil
}

// ScanGitHubOrg scans an entire GitHub organization for secrets
func (m *TruffleHogModule) ScanGitHubOrg(ctx context.Context, org string, token string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	if token != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", token)
	} else if os.Getenv("GITHUB_TOKEN") != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	}

	container = container.WithExec([]string{trufflehogBinary, "github", "--org", org, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan github org: %w", err)
	}

	return output, nil
}

// ScanDockerImage scans a Docker image for secrets
func (m *TruffleHogModule) ScanDockerImage(ctx context.Context, imageName string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithExec([]string{trufflehogBinary, "docker", "--image", imageName, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan docker image: %w", err)
	}

	return output, nil
}

// ScanS3 scans an S3 bucket for secrets
func (m *TruffleHogModule) ScanS3(ctx context.Context, bucket string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
		WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION")).
		WithExec([]string{trufflehogBinary, "s3", "--bucket", bucket, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan s3 bucket: %w", err)
	}

	return output, nil
}

// ScanWithVerification scans with verification enabled for found secrets
func (m *TruffleHogModule) ScanWithVerification(ctx context.Context, target string, targetType string) (string, error) {
	var args []string

	switch targetType {
	case "filesystem":
		args = []string{trufflehogBinary, "filesystem", target, "--json", "--verify"}
	case "git":
		args = []string{trufflehogBinary, "git", target, "--json", "--verify"}
	case "github":
		args = []string{trufflehogBinary, "github", "--repo", target, "--json", "--verify"}
	default:
		args = []string{trufflehogBinary, "filesystem", target, "--json", "--verify"}
	}

	container := m.client.Container().From("trufflesecurity/trufflehog:latest")

	if targetType == "filesystem" {
		container = container.
			WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		args[2] = "."
	}

	if os.Getenv("GITHUB_TOKEN") != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// ScanGitAdvanced scans git repository with advanced filtering options
func (m *TruffleHogModule) ScanGitAdvanced(ctx context.Context, repoURL string, branch string, sinceDate string, untilDate string, onlyVerified bool, outputFormat string, excludePaths []string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	args := []string{trufflehogBinary, "git", repoURL}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	if sinceDate != "" {
		args = append(args, "--since", sinceDate)
	}
	if untilDate != "" {
		args = append(args, "--until", untilDate)
	}
	if onlyVerified {
		args = append(args, "--only-verified")
	}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}
	for _, path := range excludePaths {
		args = append(args, "--exclude", path)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// ScanFilesystemAdvanced scans filesystem with advanced options and exclusions
func (m *TruffleHogModule) ScanFilesystemAdvanced(ctx context.Context, path string, onlyVerified bool, excludePaths []string, outputFormat string, maxDepth string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithDirectory("/workspace", m.client.Host().Directory(path)).
		WithWorkdir("/workspace")

	args := []string{trufflehogBinary, "filesystem", "."}
	if onlyVerified {
		args = append(args, "--only-verified")
	}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}
	if maxDepth != "" {
		args = append(args, "--max-depth", maxDepth)
	}
	for _, excludePath := range excludePaths {
		args = append(args, "--exclude", excludePath)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// ScanDockerAdvanced scans Docker image with advanced verification options
func (m *TruffleHogModule) ScanDockerAdvanced(ctx context.Context, image string, onlyVerified bool, outputFormat string, layers string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	args := []string{trufflehogBinary, "docker", image}
	if onlyVerified {
		args = append(args, "--only-verified")
	}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}
	if layers != "" {
		args = append(args, "--layers", layers)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// ComprehensiveSecretDetection performs comprehensive secret detection with advanced filtering
func (m *TruffleHogModule) ComprehensiveSecretDetection(ctx context.Context, target string, sourceType string, outputFormat string, outputFile string, onlyVerified bool, includeDetectors bool, confidenceLevel string, excludePaths []string, includePaths []string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	if sourceType == "filesystem" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		target = "."
	}

	args := []string{trufflehogBinary, sourceType, target}
	if onlyVerified {
		args = append(args, "--only-verified")
	}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	if includeDetectors {
		args = append(args, "--include-detectors")
	}
	if confidenceLevel != "" {
		args = append(args, "--confidence", confidenceLevel)
	}
	for _, path := range excludePaths {
		args = append(args, "--exclude", path)
	}
	for _, path := range includePaths {
		args = append(args, "--include", path)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// CloudStorageScanning scans cloud storage for secrets
func (m *TruffleHogModule) CloudStorageScanning(ctx context.Context, cloudProvider string, resourceIdentifier string, credentialsProfile string, region string, recursive bool, filePatterns string, onlyVerified bool, maxFileSize string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	var args []string
	switch cloudProvider {
	case "s3":
		args = []string{trufflehogBinary, "s3", "--bucket", resourceIdentifier}
	case "gcs":
		args = []string{trufflehogBinary, "gcs", "--bucket", resourceIdentifier}
	case "azure-storage":
		args = []string{trufflehogBinary, "azure", "--container", resourceIdentifier}
	default:
		args = []string{trufflehogBinary, cloudProvider, resourceIdentifier}
	}

	if credentialsProfile != "" {
		args = append(args, "--credentials", credentialsProfile)
	}
	if region != "" {
		args = append(args, "--region", region)
	}
	if recursive {
		args = append(args, "--recursive")
	}
	if onlyVerified {
		args = append(args, "--only-verified")
	}
	if filePatterns != "" {
		args = append(args, "--include-patterns", filePatterns)
	}
	if maxFileSize != "" {
		args = append(args, "--max-file-size", maxFileSize)
	}

	// Add AWS credentials if available
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	}
	if os.Getenv("AWS_REGION") != "" {
		container = container.WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// CustomDetectorManagement manages and uses custom secret detectors
func (m *TruffleHogModule) CustomDetectorManagement(ctx context.Context, action string, detectorConfig string, target string, detectorPattern string, includeBuiltin bool) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	var args []string
	switch action {
	case "list":
		args = []string{trufflehogBinary, "--list-detectors"}
	case "validate":
		container = container.WithFile("/config.yaml", m.client.Host().File(detectorConfig))
		args = []string{trufflehogBinary, "--validate-config", "/config.yaml"}
	case "scan-with-custom":
		container = container.WithFile("/config.yaml", m.client.Host().File(detectorConfig)).
			WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		args = []string{trufflehogBinary, "filesystem", ".", "--config", "/config.yaml"}
		if includeBuiltin {
			args = append(args, "--include-builtin")
		}
	case "test-detector":
		args = []string{trufflehogBinary, "--test-detector", detectorPattern}
	default:
		args = []string{trufflehogBinary, "--help"}
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// EnterpriseGitScanning performs enterprise-grade git repository scanning
func (m *TruffleHogModule) EnterpriseGitScanning(ctx context.Context, gitSource string, repository string, authentication string, scanMode string, commitRange string, branches string, includeForks bool, includeIssues bool, includePullRequests bool, outputFormat string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	var args []string
	switch gitSource {
	case "github":
		args = []string{trufflehogBinary, "github", "--repo", repository}
	case "gitlab":
		args = []string{trufflehogBinary, "gitlab", "--repo", repository}
	case "bitbucket":
		args = []string{trufflehogBinary, "bitbucket", "--repo", repository}
	default:
		args = []string{trufflehogBinary, "git", repository}
	}

	if authentication != "" {
		args = append(args, "--token", authentication)
		container = container.WithEnvVariable("GITHUB_TOKEN", authentication)
	}
	if commitRange != "" {
		args = append(args, "--commit-range", commitRange)
	}
	if branches != "" {
		args = append(args, "--branches", branches)
	}
	if includeForks {
		args = append(args, "--include-forks")
	}
	if includeIssues {
		args = append(args, "--include-issues")
	}
	if includePullRequests {
		args = append(args, "--include-prs")
	}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// CICDPipelineIntegration performs optimized secret scanning for CI/CD pipelines
func (m *TruffleHogModule) CICDPipelineIntegration(ctx context.Context, scanTarget string, scanType string, baselineFile string, outputFormat string, outputFile string, failOnVerified bool, failOnUnverified bool, quietMode bool, timeout string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithDirectory("/workspace", m.client.Host().Directory(scanTarget)).
		WithWorkdir("/workspace")

	args := []string{trufflehogBinary, "filesystem", "."}

	// Configure based on scan type
	switch scanType {
	case "pre-commit":
		args = append(args, "--only-verified", "--max-depth", "1")
	case "post-commit":
		args = append(args, "--include-detectors")
	case "pull-request":
		args = append(args, "--only-verified", "--format", "sarif")
	case "release":
		args = append(args, "--only-verified", "--include-detectors")
	default:
		args = append(args, "--include-detectors")
	}

	if baselineFile != "" {
		container = container.WithFile("/baseline.json", m.client.Host().File(baselineFile))
		args = append(args, "--baseline", "/baseline.json")
	}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	if quietMode {
		args = append(args, "--quiet")
	}
	if timeout != "" {
		args = append(args, "--timeout", timeout+"s")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// PerformanceOptimization performs high-performance secret scanning with optimization features
func (m *TruffleHogModule) PerformanceOptimization(ctx context.Context, target string, sourceType string, concurrency string, maxFileSize string, bufferSize string, skipBinaries bool, enableSampling bool, memoryLimit string, enableMetrics bool) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	if sourceType == "filesystem" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		target = "."
	}

	args := []string{trufflehogBinary, sourceType, target}
	if concurrency != "" {
		args = append(args, "--concurrency", concurrency)
	}
	if maxFileSize != "" {
		args = append(args, "--max-file-size", maxFileSize)
	}
	if bufferSize != "" {
		args = append(args, "--buffer-size", bufferSize)
	}
	if skipBinaries {
		args = append(args, "--skip-binaries")
	}
	if enableSampling {
		args = append(args, "--enable-sampling")
	}
	if memoryLimit != "" {
		args = append(args, "--memory-limit", memoryLimit)
	}
	if enableMetrics {
		args = append(args, "--metrics")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// ComprehensiveReporting generates comprehensive secret scanning reports with analytics
func (m *TruffleHogModule) ComprehensiveReporting(ctx context.Context, target string, sourceType string, reportType string, outputFormats string, outputDirectory string, includeVerificationStatus bool, includeRiskAssessment bool, baselineComparison string, includeTrends bool) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	if sourceType == "filesystem" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		target = "."
	}

	args := []string{trufflehogBinary, sourceType, target}

	// Configure report-specific settings
	switch reportType {
	case "executive-summary":
		args = append(args, "--only-verified", "--format", "json")
	case "technical-detail":
		args = append(args, "--include-detectors", "--format", "jsonl")
	case "compliance-audit":
		args = append(args, "--only-verified", "--format", "sarif")
	case "remediation-guide":
		args = append(args, "--include-detectors", "--format", "json")
	}

	if outputFormats != "" {
		args = append(args, "--format", outputFormats)
	}
	if outputDirectory != "" {
		args = append(args, "--output", outputDirectory+"/trufflehog-report")
	}
	if includeVerificationStatus {
		args = append(args, "--include-verification")
	}
	if baselineComparison != "" {
		container = container.WithFile("/baseline.json", m.client.Host().File(baselineComparison))
		args = append(args, "--baseline", "/baseline.json")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// GetVersion returns the version of TruffleHog
func (m *TruffleHogModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithExec([]string{trufflehogBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get trufflehog version: %w", err)
	}

	return output, nil
}