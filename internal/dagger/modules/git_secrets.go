package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// GitSecretsModule runs git-secrets for scanning git repositories for secrets
type GitSecretsModule struct {
	client *dagger.Client
	name   string
}

// NewGitSecretsModule creates a new git-secrets module
func NewGitSecretsModule(client *dagger.Client) *GitSecretsModule {
	return &GitSecretsModule{
		client: client,
		name:   "git-secrets",
	}
}

// ScanRepository scans a git repository for secrets
func (m *GitSecretsModule) ScanRepository(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"git-secrets",
			"--scan",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)

	if output != "" || stderr != "" {
		result := output
		if stderr != "" {
			result += "\n" + stderr
		}
		return result, nil
	}

	return "", fmt.Errorf("failed to run git-secrets: no output received")
}

// ScanWithAwsProviders scans with AWS secret patterns
func (m *GitSecretsModule) ScanWithAwsProviders(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"bash", "-c",
			"git-secrets --register-aws && git-secrets --scan",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)

	if output != "" || stderr != "" {
		result := output
		if stderr != "" {
			result += "\n" + stderr
		}
		return result, nil
	}

	return "", fmt.Errorf("failed to run git-secrets with AWS patterns: no output received")
}

// GetVersion returns the version of git-secrets
func (m *GitSecretsModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("cloudshipai/git-secrets:latest").
		WithExec([]string{"git-secrets", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get git-secrets version: %w", err)
	}

	return output, nil
}
