package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// CloudQueryModule runs CloudQuery for cloud asset inventory
type CloudQueryModule struct {
	client *dagger.Client
	name   string
}

// NewCloudQueryModule creates a new CloudQuery module
func NewCloudQueryModule(client *dagger.Client) *CloudQueryModule {
	return &CloudQueryModule{
		client: client,
		name:   "cloudquery",
	}
}

// SyncWithConfig syncs cloud resources using configuration
func (m *CloudQueryModule) SyncWithConfig(ctx context.Context, configPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithDirectory("/config", m.client.Host().Directory(configPath)).
		WithExec([]string{
			"cloudquery",
			"sync",
			"/config/config.yml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery sync: %w", err)
	}

	return output, nil
}

// ValidateConfig validates CloudQuery configuration
func (m *CloudQueryModule) ValidateConfig(ctx context.Context, configPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithDirectory("/config", m.client.Host().Directory(configPath)).
		WithExec([]string{
			"cloudquery",
			"validate-config",
			"/config/config.yml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate cloudquery config: %w", err)
	}

	return output, nil
}

// ListProviders lists available CloudQuery providers
func (m *CloudQueryModule) ListProviders(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithExec([]string{
			"cloudquery",
			"provider",
			"list",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list cloudquery providers: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of CloudQuery
func (m *CloudQueryModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithExec([]string{"cloudquery", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cloudquery version: %w", err)
	}

	return output, nil
}
