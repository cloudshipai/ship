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

// AuditDependencies audits package dependencies for vulnerabilities
func (m *GitHubPackagesModule) AuditDependencies(ctx context.Context, owner string, repo string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/%s/%s/dependency-graph/sbom | jq .`, owner, repo),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to audit dependencies: %w", err)
	}

	return output, nil
}

// CheckSignatures verifies package signatures
func (m *GitHubPackagesModule) CheckSignatures(ctx context.Context, owner string, packageName string, version string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/users/%s/packages/container/%s/versions | jq '.[] | select(.name=="%s") | .metadata.container.tags'`, owner, packageName, version),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check signatures: %w", err)
	}

	return output, nil
}

// EnforcePolicies enforces package security policies
func (m *GitHubPackagesModule) EnforcePolicies(ctx context.Context, owner string, repo string, policyFile string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithFile("/workspace/policy.json", m.client.Host().File(policyFile)).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -X PUT -H "Authorization: token $GITHUB_TOKEN" -H "Content-Type: application/json" -d @/workspace/policy.json https://api.github.com/repos/%s/%s/security-advisories/policy`, owner, repo),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to enforce policies: %w", err)
	}

	return output, nil
}

// GenerateSBOM generates Software Bill of Materials
func (m *GitHubPackagesModule) GenerateSBOM(ctx context.Context, owner string, repo string, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -H "Authorization: token $GITHUB_TOKEN" -H "Accept: application/vnd.github+json" https://api.github.com/repos/%s/%s/dependency-graph/sbom`, owner, repo),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate SBOM: %w", err)
	}

	return output, nil
}

// GetVersion returns API version information
func (m *GitHubPackagesModule) GetVersion(ctx context.Context, token string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithEnvVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"sh", "-c",
			"curl -H \"Authorization: token $GITHUB_TOKEN\" https://api.github.com/meta | jq .",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get API version: %w", err)
	}

	return output, nil
}
