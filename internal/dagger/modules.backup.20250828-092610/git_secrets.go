package modules

import (
	"context"

	"dagger.io/dagger"
)

// GitSecretsModule runs git-secrets for AWS credential scanning
type GitSecretsModule struct {
	client *dagger.Client
	name   string
}

// GitSecretsScanOptions contains options for git-secrets scanning
type GitSecretsScanOptions struct {
	OutputFormat string
	ScanHistory  bool
	Recursive    bool
}

// NewGitSecretsModule creates a new git-secrets module
func NewGitSecretsModule(client *dagger.Client) *GitSecretsModule {
	return &GitSecretsModule{
		client: client,
		name:   "git-secrets",
	}
}

// Scan runs git-secrets scan on the provided directory
func (m *GitSecretsModule) Scan(ctx context.Context, sourcePath string, opts GitSecretsScanOptions) (string, error) {
	args := []string{"git-secrets", "--scan"}

	// Add scan options
	if opts.Recursive {
		args = append(args, "-r")
	}
	if opts.ScanHistory {
		args = append(args, "--scan-history")
	}

	// Add source path
	args = append(args, ".")

	container := m.client.Container().
		From(getImageTag("git-secrets", "trufflesecurity/secrets:latest")).
		WithDirectory("/workspace", m.client.Host().Directory(sourcePath)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	stdout, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)

	if stdout != "" {
		return stdout, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "No secrets found", nil
}

// InstallHooks installs git-secrets hooks in a repository
func (m *GitSecretsModule) InstallHooks(ctx context.Context, repoPath string, force bool) (string, error) {
	args := []string{"git-secrets", "--install"}
	
	if force {
		args = append(args, "--force")
	}

	container := m.client.Container().
		From(getImageTag("git-secrets", "trufflesecurity/secrets:latest")).
		WithDirectory("/workspace", m.client.Host().Directory(repoPath)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output == "" {
		output = "Git-secrets hooks installed successfully"
	}

	return output, nil
}