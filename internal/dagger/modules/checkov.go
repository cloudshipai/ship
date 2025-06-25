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
			"checkov",
			"--directory", ".",
			"--output", "json",
			"--framework", "terraform",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run checkov: %w", err)
	}

	return output, nil
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
			"checkov",
			"--file", filename,
			"--output", "json",
			"--framework", "terraform",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run checkov on file: %w", err)
	}

	return output, nil
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

// GetVersion returns the version of Checkov
func (m *CheckovModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("bridgecrew/checkov:latest").
		WithExec([]string{"checkov", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get checkov version: %w", err)
	}

	return output, nil
}