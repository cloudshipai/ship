package modules

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// TruffleHogModule runs TruffleHog for verified secret detection
type TruffleHogModule struct {
	client *dagger.Client
	name   string
}

// NewTruffleHogModule creates a new TruffleHog module
func NewTruffleHogModule(client *dagger.Client) *TruffleHogModule {
	return &TruffleHogModule{
		client: client,
		name:   "trufflehog",
	}
}

// ScanDirectory scans a directory for secrets using TruffleHog
func (m *TruffleHogModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"trufflehog", "filesystem", ".", "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// TruffleHog returns non-zero exit code when secrets are found
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// ScanGitRepo scans a Git repository for secrets
func (m *TruffleHogModule) ScanGitRepo(ctx context.Context, repoURL string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithExec([]string{"trufflehog", "git", repoURL, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan git repo: %w", err)
	}

	return output, nil
}

// ScanGitHub scans a GitHub repository for secrets
func (m *TruffleHogModule) ScanGitHub(ctx context.Context, repo string, token string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	if token != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", token)
	} else if os.Getenv("GITHUB_TOKEN") != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	}

	container = container.WithExec([]string{"trufflehog", "github", "--repo", repo, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan github repo: %w", err)
	}

	return output, nil
}

// ScanGitHubOrg scans an entire GitHub organization for secrets
func (m *TruffleHogModule) ScanGitHubOrg(ctx context.Context, org string, token string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest")

	if token != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", token)
	} else if os.Getenv("GITHUB_TOKEN") != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	}

	container = container.WithExec([]string{"trufflehog", "github", "--org", org, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan github org: %w", err)
	}

	return output, nil
}

// ScanDockerImage scans a Docker image for secrets
func (m *TruffleHogModule) ScanDockerImage(ctx context.Context, imageName string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithExec([]string{"trufflehog", "docker", "--image", imageName, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan docker image: %w", err)
	}

	return output, nil
}

// ScanS3 scans an S3 bucket for secrets
func (m *TruffleHogModule) ScanS3(ctx context.Context, bucket string) (string, error) {
	container := m.client.Container().
		From("trufflesecurity/trufflehog:latest").
		WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
		WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION")).
		WithExec([]string{"trufflehog", "s3", "--bucket", bucket, "--json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan s3 bucket: %w", err)
	}

	return output, nil
}

// ScanWithVerification scans with verification enabled for found secrets
func (m *TruffleHogModule) ScanWithVerification(ctx context.Context, target string, targetType string) (string, error) {
	var args []string

	switch targetType {
	case "filesystem":
		args = []string{"trufflehog", "filesystem", target, "--json", "--verify"}
	case "git":
		args = []string{"trufflehog", "git", target, "--json", "--verify"}
	case "github":
		args = []string{"trufflehog", "github", "--repo", target, "--json", "--verify"}
	default:
		args = []string{"trufflehog", "filesystem", target, "--json", "--verify"}
	}

	container := m.client.Container().From("trufflesecurity/trufflehog:latest")

	if targetType == "filesystem" {
		container = container.
			WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		args[2] = "."
	}

	if os.Getenv("GITHUB_TOKEN") != "" {
		container = container.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}