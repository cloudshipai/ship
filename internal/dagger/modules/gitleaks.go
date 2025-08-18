package modules

import (
	"context"
	"path/filepath"

	"dagger.io/dagger"
)

// GitleaksModule runs Gitleaks for secret detection
type GitleaksModule struct {
	client *dagger.Client
	name   string
}

const gitleaksBinary = "/usr/bin/gitleaks"

// NewGitleaksModule creates a new Gitleaks module
func NewGitleaksModule(client *dagger.Client) *GitleaksModule {
	return &GitleaksModule{
		client: client,
		name:   "gitleaks",
	}
}

// ScanDirectory scans a directory for secrets using Gitleaks
func (m *GitleaksModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("zricethezav/gitleaks:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{gitleaksBinary, "dir", ".", "--report-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Gitleaks returns non-zero exit code when secrets are found
		// Try to get the output anyway
		output, _ = container.Stderr(ctx)
		return output, nil
	}

	return output, nil
}

// ScanFile scans a specific file for secrets
func (m *GitleaksModule) ScanFile(ctx context.Context, filePath string) (string, error) {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)

	container := m.client.Container().
		From("zricethezav/gitleaks:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{gitleaksBinary, "dir", filename, "--report-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr output if stdout fails
		output, _ = container.Stderr(ctx)
		return output, nil
	}

	return output, nil
}

// ScanGitRepo scans a git repository for secrets
func (m *GitleaksModule) ScanGitRepo(ctx context.Context, repoDir string) (string, error) {
	container := m.client.Container().
		From("zricethezav/gitleaks:latest").
		WithDirectory("/workspace", m.client.Host().Directory(repoDir)).
		WithWorkdir("/workspace").
		WithExec([]string{gitleaksBinary, "git", ".", "--report-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr output if stdout fails
		output, _ = container.Stderr(ctx)
		return output, nil
	}

	return output, nil
}

// ScanWithConfig scans using a custom Gitleaks configuration
func (m *GitleaksModule) ScanWithConfig(ctx context.Context, dir string, configFile string) (string, error) {
	container := m.client.Container().
		From("zricethezav/gitleaks:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{gitleaksBinary, "dir", ".", "--config", configFile, "--report-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr output if stdout fails
		output, _ = container.Stderr(ctx)
		return output, nil
	}

	return output, nil
}

// ScanStdin scans secrets from stdin input
func (m *GitleaksModule) ScanStdin(ctx context.Context, input string) (string, error) {
	container := m.client.Container().
		From("zricethezav/gitleaks:latest").
		WithExec([]string{"sh", "-c", "echo '" + input + "' | gitleaks stdin --report-format json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr output if stdout fails
		output, _ = container.Stderr(ctx)
		return output, nil
	}

	return output, nil
}

// GetVersion returns the Gitleaks version
func (m *GitleaksModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("zricethezav/gitleaks:latest").
		WithExec([]string{gitleaksBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", err
	}

	return output, nil
}