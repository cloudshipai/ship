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

const gitBinary = "/usr/bin/git"

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
		From("depop/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			gitBinary, "secrets",
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
		From("depop/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"bash", "-c",
			"git secrets --register-aws && git secrets --scan",
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
		From("depop/git-secrets:latest").
		WithExec([]string{"git", "secrets", "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	// git-secrets doesn't have a --version flag, return image info
	if output == "" && stderr != "" {
		// Version flag not supported, return default
		return "git-secrets-latest", nil
	}
	
	if output != "" {
		return output, nil
	}

	return "git-secrets-latest", nil
}

// ScanHistory scans git history for secrets
func (m *GitSecretsModule) ScanHistory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("depop/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"git", "secrets", "--scan-history"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan git history: %w", err)
	}

	return output, nil
}

// InstallHooks installs git-secrets hooks
func (m *GitSecretsModule) InstallHooks(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("depop/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"git", "secrets", "--install"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install git-secrets hooks: %w", err)
	}

	return output, nil
}

// AddPattern adds a secret pattern
func (m *GitSecretsModule) AddPattern(ctx context.Context, dir string, pattern string, isAllowed bool) (string, error) {
	args := []string{"git", "secrets", "--add"}
	if isAllowed {
		args = append(args, "--allowed", pattern)
	} else {
		args = append(args, pattern)
	}

	container := m.client.Container().
		From("depop/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to add pattern: %w", err)
	}

	return output, nil
}

// ListConfig lists git-secrets configuration
func (m *GitSecretsModule) ListConfig(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("depop/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"git", "secrets", "--list"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list git-secrets config: %w", err)
	}

	return output, nil
}

// RegisterAWS registers AWS patterns
func (m *GitSecretsModule) RegisterAWS(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("depop/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"git", "secrets", "--register-aws"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to register AWS patterns: %w", err)
	}

	return output, nil
}

// AddAllowedPattern adds an allowed pattern to prevent false positives (MCP compatible)
func (m *GitSecretsModule) AddAllowedPattern(ctx context.Context, dir string, pattern string) (string, error) {
	container := m.client.Container().
		From("depop/git-secrets:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"git", "secrets", "--add", "-a", pattern})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to add allowed pattern: %w", err)
	}

	return output, nil
}
