package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// GitHubPackagesModule manages GitHub Packages security
type GitHubPackagesModule struct {
	client *dagger.Client
	name   string
}

// NewGitHubPackagesModule creates a new GitHub Packages module
func NewGitHubPackagesModule(client *dagger.Client) *GitHubPackagesModule {
	return &GitHubPackagesModule{
		client: client,
		name:   "github-packages",
	}
}

// ScanPackage scans a GitHub package for vulnerabilities
func (m *GitHubPackagesModule) ScanPackage(ctx context.Context, packageName string, version string, token string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"trivy",
			"image",
			"--format", "json",
			fmt.Sprintf("ghcr.io/%s:%s", packageName, version),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan package: %w", err)
	}

	return output, nil
}

// ListPackages lists packages in a repository
func (m *GitHubPackagesModule) ListPackages(ctx context.Context, owner string, repo string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/%s/%s/packages | jq .`, owner, repo),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list packages: %w", err)
	}

	return output, nil
}

// GetPackageVersions gets versions of a package
func (m *GitHubPackagesModule) GetPackageVersions(ctx context.Context, owner string, packageName string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/users/%s/packages/container/%s/versions | jq .`, owner, packageName),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get package versions: %w", err)
	}

	return output, nil
}
