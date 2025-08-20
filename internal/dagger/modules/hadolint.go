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
			"/bin/hadolint", "--format", "json", "Dockerfile",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run hadolint: no output received")
}

// ScanDirectory scans all Dockerfiles in a directory (simplified to scan Dockerfile in root)
func (m *HadolintModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("hadolint/hadolint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"/bin/hadolint", "--format", "json", "Dockerfile",
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
		WithExec([]string{
			"/bin/hadolint", "--version",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to get hadolint version: no output received")
}

// ScanWithConfig scans Dockerfile with custom configuration
func (m *HadolintModule) ScanWithConfig(ctx context.Context, dockerfilePath string, configPath string) (string, error) {
	container := m.client.Container().
		From("hadolint/hadolint:latest").
		WithFile("/workspace/Dockerfile", m.client.Host().File(dockerfilePath)).
		WithFile("/workspace/.hadolint.yaml", m.client.Host().File(configPath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"/bin/hadolint", "--config", ".hadolint.yaml", "--format", "json", "Dockerfile",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run hadolint with config: no output received")
}

// ScanIgnoreRules scans with specific rules ignored (alias for compatibility)
func (m *HadolintModule) ScanIgnoreRules(ctx context.Context, dockerfilePath string, ignoredRules []string) (string, error) {
	return m.ScanWithRules(ctx, dockerfilePath, ignoredRules)
}

// ScanWithRules scans with specific rules ignored
func (m *HadolintModule) ScanWithRules(ctx context.Context, dockerfilePath string, ignoredRules []string) (string, error) {
	args := []string{"/bin/hadolint"}
	for _, rule := range ignoredRules {
		args = append(args, "--ignore", rule)
	}
	args = append(args, "--format", "json", "Dockerfile")

	container := m.client.Container().
		From("hadolint/hadolint:latest").
		WithFile("/workspace/Dockerfile", m.client.Host().File(dockerfilePath)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run hadolint with rules: no output received")
}