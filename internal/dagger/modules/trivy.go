package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// TrivyModule runs Trivy for comprehensive vulnerability scanning
type TrivyModule struct {
	client *dagger.Client
	name   string
}

// NewTrivyModule creates a new Trivy module
func NewTrivyModule(client *dagger.Client) *TrivyModule {
	return &TrivyModule{
		client: client,
		name:   "trivy",
	}
}

// ScanImage scans a container image for vulnerabilities
func (m *TrivyModule) ScanImage(ctx context.Context, imageName string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			"trivy",
			"image",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			imageName,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan image: %w", err)
	}

	return output, nil
}

// ScanFilesystem scans a filesystem for vulnerabilities
func (m *TrivyModule) ScanFilesystem(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"trivy",
			"fs",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			".",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan filesystem: %w", err)
	}

	return output, nil
}

// ScanRepository scans a git repository
func (m *TrivyModule) ScanRepository(ctx context.Context, repoURL string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			"trivy",
			"repo",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			repoURL,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan repository: %w", err)
	}

	return output, nil
}

// ScanConfig scans configuration files for misconfigurations
func (m *TrivyModule) ScanConfig(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"trivy",
			"config",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			".",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan config: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Trivy
func (m *TrivyModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{"trivy", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get trivy version: %w", err)
	}

	return output, nil
}
