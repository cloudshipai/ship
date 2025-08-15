package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// RegistryModule runs container registry operations
type RegistryModule struct {
	client *dagger.Client
	name   string
}

// NewRegistryModule creates a new registry module
func NewRegistryModule(client *dagger.Client) *RegistryModule {
	return &RegistryModule{
		client: client,
		name:   "registry",
	}
}

// ScanRegistry scans container registry for vulnerabilities
func (m *RegistryModule) ScanRegistry(ctx context.Context, registryURL string, repository string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			"trivy",
			"image",
			"--format", "json",
			"--server", registryURL,
			repository,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan registry: %w", err)
	}

	return output, nil
}

// ListRepositories lists repositories in registry
func (m *RegistryModule) ListRepositories(ctx context.Context, registryURL string, username string, password string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -u %s:%s %s/v2/_catalog | jq .`, username, password, registryURL),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list repositories: %w", err)
	}

	return output, nil
}

// GetImageTags gets tags for an image
func (m *RegistryModule) GetImageTags(ctx context.Context, registryURL string, repository string, username string, password string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"}).
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`curl -u %s:%s %s/v2/%s/tags/list | jq .`, username, password, registryURL, repository),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get image tags: %w", err)
	}

	return output, nil
}

// CheckImageSecurity checks image security properties
func (m *RegistryModule) CheckImageSecurity(ctx context.Context, imageName string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			"trivy",
			"image",
			"--format", "json",
			"--scanners", "vuln,secret,config",
			imageName,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check image security: %w", err)
	}

	return output, nil
}
