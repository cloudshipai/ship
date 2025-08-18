package modules

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
)

// CheckovModule runs Checkov for multi-cloud security scanning
type CheckovModule struct {
	client *dagger.Client
	name   string
}

// Checkov is installed via pip in the container
const checkovBinary = "/usr/local/bin/checkov"

// NewCheckovModule creates a new Checkov module
func NewCheckovModule(client *dagger.Client) *CheckovModule {
	return &CheckovModule{
		client: client,
		name:   "checkov",
	}
}

// ScanDirectory scans a directory for security issues
func (m *CheckovModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			checkovBinary,
			"--directory", ".",
			"--output", "json",
			"--framework", "terraform",
		}, dagger.ContainerWithExecOpts{
			// Checkov returns non-zero exit code when it finds issues, which is expected
			Expect: "ANY",
		})

	// Always get the output, even if checkov "failed" (found issues)
	output, _ := container.Stdout(ctx)

	// If we got output, it's a success (checkov ran and produced results)
	if output != "" {
		return output, nil
	}

	// Only return error if we got no output at all
	return "", fmt.Errorf("failed to run checkov: no output received")
}

// ScanFile scans a specific file for security issues
func (m *CheckovModule) ScanFile(ctx context.Context, filePath string) (string, error) {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			checkovBinary,
			"--file", filename,
			"--output", "json",
			"--framework", "terraform",
		}, dagger.ContainerWithExecOpts{
			// Checkov returns non-zero exit code when it finds issues, which is expected
			Expect: "ANY",
		})

	// Always get the output, even if checkov "failed" (found issues)
	output, _ := container.Stdout(ctx)

	// If we got output, it's a success (checkov ran and produced results)
	if output != "" {
		return output, nil
	}

	// Only return error if we got no output at all
	return "", fmt.Errorf("failed to run checkov on file: no output received")
}

// ScanWithPolicy scans using custom policies
func (m *CheckovModule) ScanWithPolicy(ctx context.Context, dir string, policyPath string) (string, error) {
	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir))

	args := []string{
		"checkov",
		"--directory", "/workspace",
		"--output", "json",
		"--framework", "terraform",
	}

	// If policy path is provided, mount it
	if policyPath != "" {
		container = container.WithDirectory("/policies", m.client.Host().Directory(policyPath))
		args = append(args, "--external-checks-dir", "/policies")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run checkov with policy: %w", err)
	}

	return output, nil
}

// ScanMultiFramework scans for multiple cloud frameworks
func (m *CheckovModule) ScanMultiFramework(ctx context.Context, dir string, frameworks []string) (string, error) {
	args := []string{
		"checkov",
		"--directory", ".",
		"--output", "json",
	}

	// Add frameworks
	for _, framework := range frameworks {
		args = append(args, "--framework", framework)
	}

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run checkov with multiple frameworks: %w", err)
	}

	return output, nil
}

// ScanWithSeverity scans filtering by severity levels
func (m *CheckovModule) ScanWithSeverity(ctx context.Context, dir string, severities []string) (string, error) {
	args := []string{
		"checkov",
		"--directory", ".",
		"--output", "json",
		"--framework", "terraform",
	}

	// Add severity filter
	if len(severities) > 0 {
		for _, severity := range severities {
			args = append(args, "--check", severity)
		}
	}

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run checkov with severity filter: %w", err)
	}

	return output, nil
}

// ScanWithSkips scans while skipping specific checks
func (m *CheckovModule) ScanWithSkips(ctx context.Context, dir string, skipChecks []string) (string, error) {
	args := []string{
		"checkov",
		"--directory", ".",
		"--output", "json",
		"--framework", "terraform",
	}

	// Add skip checks
	for _, check := range skipChecks {
		args = append(args, "--skip-check", check)
	}

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run checkov with skip checks: %w", err)
	}

	return output, nil
}

// ScanDockerImage scans Docker container image for vulnerabilities
func (m *CheckovModule) ScanDockerImage(ctx context.Context, dockerImage string, dockerfilePath string, output string) (string, error) {
	args := []string{checkovBinary, "--docker-image", dockerImage, "--framework", "sca_image"}
	if dockerfilePath != "" {
		args = append(args, "--dockerfile-path", "/workspace/Dockerfile")
	}
	if output != "" {
		args = append(args, "--output", output)
	}

	container := m.client.Container().
		From("bridgecrew/checkov:latest")

	if dockerfilePath != "" {
		container = container.WithFile("/workspace/Dockerfile", m.client.Host().File(dockerfilePath))
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output_result, _ := container.Stdout(ctx)
	if output_result != "" {
		return output_result, nil
	}

	return "", fmt.Errorf("failed to scan docker image: no output received")
}

// ScanPackages scans package dependencies for vulnerabilities
func (m *CheckovModule) ScanPackages(ctx context.Context, dir string, output string) (string, error) {
	args := []string{checkovBinary, "--directory", ".", "--framework", "sca_package"}
	if output != "" {
		args = append(args, "--output", output)
	}

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output_result, _ := container.Stdout(ctx)
	if output_result != "" {
		return output_result, nil
	}

	return "", fmt.Errorf("failed to scan packages: no output received")
}

// ScanSecrets scans for hardcoded secrets in code
func (m *CheckovModule) ScanSecrets(ctx context.Context, dir string, output string) (string, error) {
	args := []string{checkovBinary, "--directory", ".", "--framework", "secrets"}
	if output != "" {
		args = append(args, "--output", output)
	}

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output_result, _ := container.Stdout(ctx)
	if output_result != "" {
		return output_result, nil
	}

	return "", fmt.Errorf("failed to scan secrets: no output received")
}

// ScanWithConfig scans using configuration file
func (m *CheckovModule) ScanWithConfig(ctx context.Context, dir string, configFile string) (string, error) {
	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithFile("/workspace/config.yml", m.client.Host().File(configFile)).
		WithWorkdir("/workspace").
		WithExec([]string{checkovBinary, "--directory", ".", "--config-file", "config.yml"}, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to scan with config: no output received")
}

// CreateConfig generates configuration file from current settings
func (m *CheckovModule) CreateConfig(ctx context.Context, configPath string) (string, error) {
	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithExec([]string{checkovBinary, "--create-config", "/workspace/config.yml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create config: %w", err)
	}

	return output, nil
}

// ScanWithExternalModules scans with external module downloading enabled
func (m *CheckovModule) ScanWithExternalModules(ctx context.Context, dir string, downloadExternalModules bool) (string, error) {
	args := []string{checkovBinary, "--directory", "."}
	if downloadExternalModules {
		args = append(args, "--download-external-modules", "true")
	}

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to scan with external modules: no output received")
}

// GetVersion returns the version of Checkov
func (m *CheckovModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithExec([]string{checkovBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get checkov version: %w", err)
	}

	return output, nil
}

// ScanDirectoryWithOptions scans a directory with configurable options
func (m *CheckovModule) ScanDirectoryWithOptions(ctx context.Context, dir string, framework string, output string, compact bool, quiet bool) (string, error) {
	args := []string{checkovBinary, "--directory", "."}
	
	if framework != "" {
		args = append(args, "--framework", framework)
	}
	if output != "" {
		args = append(args, "--output", output)
	} else {
		args = append(args, "--output", "json")
	}
	if compact {
		args = append(args, "--compact")
	}
	if quiet {
		args = append(args, "--quiet")
	}

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	result, _ := container.Stdout(ctx)
	if result != "" {
		return result, nil
	}

	return "", fmt.Errorf("failed to scan directory with options: no output received")
}

// ScanFileWithOptions scans a file with configurable options
func (m *CheckovModule) ScanFileWithOptions(ctx context.Context, filePath string, framework string, output string) (string, error) {
	args := []string{checkovBinary, "--file", "/workspace/input"}
	
	if framework != "" {
		args = append(args, "--framework", framework)
	}
	if output != "" {
		args = append(args, "--output", output)
	} else {
		args = append(args, "--output", "json")
	}

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithFile("/workspace/input", m.client.Host().File(filePath)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	result, _ := container.Stdout(ctx)
	if result != "" {
		return result, nil
	}

	return "", fmt.Errorf("failed to scan file with options: no output received")
}

// ScanWithSpecificChecks scans with specific checks enabled or disabled
func (m *CheckovModule) ScanWithSpecificChecks(ctx context.Context, dir string, checks string, skipChecks string) (string, error) {
	args := []string{checkovBinary, "--directory", "."}
	
	if checks != "" {
		args = append(args, "--check", checks)
	}
	if skipChecks != "" {
		args = append(args, "--skip-check", skipChecks)
	}
	args = append(args, "--output", "json")

	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	result, _ := container.Stdout(ctx)
	if result != "" {
		return result, nil
	}

	return "", fmt.Errorf("failed to scan with specific checks: no output received")
}
