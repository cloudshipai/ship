package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// TerrascanModule runs Terrascan for IaC security scanning
type TerrascanModule struct {
	client *dagger.Client
	name   string
}

// NewTerrascanModule creates a new Terrascan module
func NewTerrascanModule(client *dagger.Client) *TerrascanModule {
	return &TerrascanModule{
		client: client,
		name:   "terrascan",
	}
}

// ScanDirectory scans a directory for IaC security issues using Terrascan
func (m *TerrascanModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Terrascan returns non-zero exit code when violations are found
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan: %w", err)
	}

	return output, nil
}

// ScanTerraform scans Terraform files specifically
func (m *TerrascanModule) ScanTerraform(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-i", "terraform", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan on terraform: %w", err)
	}

	return output, nil
}

// ScanKubernetes scans Kubernetes manifests
func (m *TerrascanModule) ScanKubernetes(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-i", "k8s", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan on kubernetes: %w", err)
	}

	return output, nil
}

// ScanCloudFormation scans CloudFormation templates
func (m *TerrascanModule) ScanCloudFormation(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-i", "cloudformation", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan on cloudformation: %w", err)
	}

	return output, nil
}

// ScanDockerfiles scans Dockerfile for security issues
func (m *TerrascanModule) ScanDockerfiles(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-i", "docker", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan on docker: %w", err)
	}

	return output, nil
}

// ScanWithSeverity scans with a specific severity threshold
func (m *TerrascanModule) ScanWithSeverity(ctx context.Context, dir string, severity string, iacType string) (string, error) {
	args := []string{"terrascan", "scan", "-d", ".", "-o", "json"}
	
	if iacType != "" {
		args = append(args, "-i", iacType)
	}
	
	if severity != "" {
		args = append(args, "--severity", severity)
	}

	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan with severity %s: %w", severity, err)
	}

	return output, nil
}

// ScanRemote scans remote repository
func (m *TerrascanModule) ScanRemote(ctx context.Context, repoURL string, repoType string, outputFormat string) (string, error) {
	args := []string{"terrascan", "scan", "-r", repoURL}
	if repoType != "" {
		args = append(args, "--remote-type", repoType)
	}
	if outputFormat != "" {
		args = append(args, "--output", outputFormat)
	}

	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan remote repository: %w", err)
	}

	return output, nil
}

// ScanWithPolicy scans using custom policy path
func (m *TerrascanModule) ScanWithPolicy(ctx context.Context, dir string, policyPath string, outputFormat string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithDirectory("/policies", m.client.Host().Directory(policyPath)).
		WithWorkdir("/workspace")

	args := []string{"terrascan", "scan", "-d", ".", "--policy-path", "/policies"}
	if outputFormat != "" {
		args = append(args, "--output", outputFormat)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan with custom policy: %w", err)
	}

	return output, nil
}

// ComprehensiveIaCScan performs comprehensive Infrastructure as Code security scanning
func (m *TerrascanModule) ComprehensiveIaCScan(ctx context.Context, target string, iacType string, outputFormat string, outputFile string, severityThreshold string, policyTypes string, excludeRules string, verbose bool, showPassed bool) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(target)).
		WithWorkdir("/workspace")

	args := []string{"terrascan", "scan", "-i", iacType, "-d", "."}
	if outputFormat != "" {
		args = append(args, "-o", outputFormat)
	}
	if outputFile != "" {
		args = append(args, "--output-file", outputFile)
	}
	if severityThreshold != "" {
		args = append(args, "--severity", severityThreshold)
	}
	if policyTypes != "" {
		args = append(args, "--policy-type", policyTypes)
	}
	if excludeRules != "" {
		args = append(args, "--skip-rules", excludeRules)
	}
	if verbose {
		args = append(args, "-v")
	}
	if showPassed {
		args = append(args, "--show-passed")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed comprehensive IaC scan: %w", err)
	}

	return output, nil
}

// ComplianceFrameworkScan scans against compliance frameworks
func (m *TerrascanModule) ComplianceFrameworkScan(ctx context.Context, target string, complianceFramework string, iacType string, outputFormat string, outputFile string, includeSeverityDetails bool) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(target)).
		WithWorkdir("/workspace")

	args := []string{"terrascan", "scan", "-i", iacType, "-d", ".", "--policy-type", complianceFramework}
	if outputFormat != "" {
		args = append(args, "-o", outputFormat)
	}
	if outputFile != "" {
		args = append(args, "--output-file", outputFile)
	}
	if includeSeverityDetails {
		args = append(args, "--show-passed", "-v")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed compliance framework scan: %w", err)
	}

	return output, nil
}

// RemoteRepositoryScan performs advanced remote repository scanning
func (m *TerrascanModule) RemoteRepositoryScan(ctx context.Context, repoURL string, repoType string, iacType string, branch string, sshKeyPath string, accessToken string, configPath string, outputFormat string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest")

	if sshKeyPath != "" {
		container = container.WithFile("/ssh_key", m.client.Host().File(sshKeyPath))
	}
	if configPath != "" {
		container = container.WithFile("/config.yaml", m.client.Host().File(configPath))
	}

	args := []string{"terrascan", "scan", "-r", repoURL, "-t", repoType, "-i", iacType}
	if branch != "" {
		args = append(args, "--remote-branch", branch)
	}
	if sshKeyPath != "" {
		args = append(args, "--ssh-key", "/ssh_key")
	}
	if accessToken != "" {
		args = append(args, "--access-token", accessToken)
	}
	if configPath != "" {
		args = append(args, "-c", "/config.yaml")
	}
	if outputFormat != "" {
		args = append(args, "-o", outputFormat)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed remote repository scan: %w", err)
	}

	return output, nil
}

// CustomPolicyManagement manages and validates custom Terrascan policies
func (m *TerrascanModule) CustomPolicyManagement(ctx context.Context, action string, policyPath string, target string, iacType string, testDataPath string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest")

	if policyPath != "" {
		container = container.WithDirectory("/policies", m.client.Host().Directory(policyPath))
	}
	if target != "" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(target))
	}
	if testDataPath != "" {
		container = container.WithDirectory("/testdata", m.client.Host().Directory(testDataPath))
	}

	var args []string
	switch action {
	case "validate":
		args = []string{"terrascan", "init", "--policy-path", "/policies"}
	case "test":
		args = []string{"terrascan", "scan", "--policy-path", "/policies", "-d", "/testdata"}
	case "scan-with-custom":
		args = []string{"terrascan", "scan", "-i", iacType, "-d", "/workspace", "--policy-path", "/policies"}
	default:
		args = []string{"terrascan", "--help"}
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed custom policy management: %w", err)
	}

	return output, nil
}

// CICDPipelineIntegration performs optimized IaC security scanning for CI/CD pipelines
func (m *TerrascanModule) CICDPipelineIntegration(ctx context.Context, target string, iacType string, pipelineStage string, gatePolicy string, outputFormat string, outputFile string, failOnViolations bool, baselineFile string, quietMode bool) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(target)).
		WithWorkdir("/workspace")

	if baselineFile != "" {
		container = container.WithFile("/baseline.json", m.client.Host().File(baselineFile))
	}

	args := []string{"terrascan", "scan", "-i", iacType, "-d", "."}

	// Configure based on pipeline stage and gate policy
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

	if outputFormat != "" {
		args = append(args, "-o", outputFormat)
	}
	if outputFile != "" {
		args = append(args, "--output-file", outputFile)
	}
	if baselineFile != "" {
		args = append(args, "--baseline", "/baseline.json")
	}
	if quietMode {
		args = append(args, "--quiet")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed CI/CD pipeline integration: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Terrascan
func (m *TerrascanModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithExec([]string{"terrascan", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get terrascan version: %w", err)
	}

	return output, nil
}