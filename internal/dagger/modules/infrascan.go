package modules

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
)

// InfraScanModule runs Trivy for security scanning of Terraform code
// Using Trivy instead of InfraScan as it provides better Terraform security scanning
type InfraScanModule struct {
	client *dagger.Client
	name   string
}

// NewInfraScanModule creates a new InfraScan module (using Trivy)
func NewInfraScanModule(client *dagger.Client) *InfraScanModule {
	return &InfraScanModule{
		client: client,
		name:   "trivy",
	}
}

// ScanDirectory scans a directory for security issues
func (m *InfraScanModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	// Mount the directory and run Trivy
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"trivy",
			"fs",
			".",
			"--scanners", "misconfig",
			"--format", "json",
			"--severity", "HIGH,CRITICAL,MEDIUM,LOW",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run trivy: %w", err)
	}

	return output, nil
}

// ScanFile scans a specific Terraform file
func (m *InfraScanModule) ScanFile(ctx context.Context, filePath string) (string, error) {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)

	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"trivy",
			"fs",
			filename,
			"--scanners", "misconfig",
			"--format", "json",
			"--severity", "HIGH,CRITICAL,MEDIUM,LOW",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run trivy on file: %w", err)
	}

	return output, nil
}

// ScanWithRules scans using custom rule set
func (m *InfraScanModule) ScanWithRules(ctx context.Context, dir string, rulesFile string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir))

	// If rules file is provided, mount it
	if rulesFile != "" {
		container = container.WithFile("/policy.rego", m.client.Host().File(rulesFile))
		container = container.WithExec([]string{
			"trivy",
			"fs",
			"/workspace",
			"--scanners", "misconfig",
			"--format", "json",
			"--severity", "HIGH,CRITICAL,MEDIUM,LOW",
			"--config-policy", "/policy.rego",
		})
	} else {
		container = container.WithExec([]string{
			"trivy",
			"fs",
			"/workspace",
			"--scanners", "misconfig",
			"--format", "json",
			"--severity", "HIGH,CRITICAL,MEDIUM,LOW",
		})
	}

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run trivy with rules: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Trivy
func (m *InfraScanModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{"trivy", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get trivy version: %w", err)
	}

	return output, nil
}
