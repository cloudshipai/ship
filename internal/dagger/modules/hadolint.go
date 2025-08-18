package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// HadolintModule runs Hadolint for Dockerfile linting
type HadolintModule struct {
	client *dagger.Client
	name   string
}

// NewHadolintModule creates a new Hadolint module
func NewHadolintModule(client *dagger.Client) *HadolintModule {
	return &HadolintModule{
		client: client,
		name:   "hadolint",
	}
}

// ScanDockerfile scans a Dockerfile for best practices
func (m *HadolintModule) ScanDockerfile(ctx context.Context, dockerfilePath string) (string, error) {
	container := m.client.Container().
		From("hadolint/hadolint:latest").
		WithFile("/workspace/Dockerfile", m.client.Host().File(dockerfilePath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"hadolint",
			"--format", "json",
			"Dockerfile",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run hadolint: no output received")
}

// ScanDirectory scans all Dockerfiles in a directory
func (m *HadolintModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("hadolint/hadolint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"find", ".", "-name", "Dockerfile*", "-exec", "hadolint", "--format", "json", "{}", "+",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run hadolint on directory: no output received")
}

// GetVersion returns the version of Hadolint
func (m *HadolintModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("hadolint/hadolint:latest").
		WithExec([]string{"hadolint", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get hadolint version: %w", err)
	}

	return output, nil
}

// ScanWithConfig scans using custom configuration
func (m *HadolintModule) ScanWithConfig(ctx context.Context, dockerfilePath string, configPath string) (string, error) {
	container := m.client.Container().
		From("hadolint/hadolint:latest").
		WithFile("/workspace/Dockerfile", m.client.Host().File(dockerfilePath)).
		WithFile("/workspace/.hadolint.yaml", m.client.Host().File(configPath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"hadolint",
			"--config", ".hadolint.yaml",
			"--format", "json",
			"Dockerfile",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan dockerfile with config: %w", err)
	}

	return output, nil
}

// ScanIgnoreRules scans while ignoring specific rules
func (m *HadolintModule) ScanIgnoreRules(ctx context.Context, dockerfilePath string, ignoreRules []string) (string, error) {
	args := []string{"hadolint", "--format", "json"}
	
	// Add ignore rules
	for _, rule := range ignoreRules {
		args = append(args, "--ignore", rule)
	}
	args = append(args, "/workspace/Dockerfile")

	container := m.client.Container().
		From("hadolint/hadolint:latest").
		WithFile("/workspace/Dockerfile", m.client.Host().File(dockerfilePath)).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan dockerfile with ignore rules: %w", err)
	}

	return output, nil
}
