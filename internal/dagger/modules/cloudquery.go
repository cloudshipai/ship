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

const cloudQueryBinary = "/app/cloudquery"

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
			cloudQueryBinary,
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
			cloudQueryBinary,
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
			cloudQueryBinary,
			"provider",
			"list",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list cloudquery providers: %w", err)
	}

	return output, nil
}

// MigrateConfig updates destination schema
func (m *CloudQueryModule) MigrateConfig(ctx context.Context, configPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithDirectory("/config", m.client.Host().Directory(configPath)).
		WithExec([]string{
			cloudQueryBinary,
			"migrate",
			"/config/config.yml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery migrate: %w", err)
	}

	return output, nil
}

// InitConfig generates initial configuration
func (m *CloudQueryModule) InitConfig(ctx context.Context, source string, destination string) (string, error) {
	args := []string{"cloudquery", "init"}
	if source != "" {
		args = append(args, "--source", source)
	}
	if destination != "" {
		args = append(args, "--destination", destination)
	}

	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery init: %w", err)
	}

	return output, nil
}

// TestConnection tests plugin connections
func (m *CloudQueryModule) TestConnection(ctx context.Context, configPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithDirectory("/config", m.client.Host().Directory(configPath)).
		WithExec([]string{
			cloudQueryBinary,
			"test-connection",
			"/config/config.yml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery test-connection: %w", err)
	}

	return output, nil
}

// GetTables generates table documentation
func (m *CloudQueryModule) GetTables(ctx context.Context, source string, outputDir string, format string) (string, error) {
	args := []string{"cloudquery", "tables"}
	if source != "" {
		args = append(args, source)
	}
	if outputDir != "" {
		args = append(args, "--output-dir", outputDir)
	}
	if format != "" {
		args = append(args, "--format", format)
	}

	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery tables: %w", err)
	}

	return output, nil
}

// Login to CloudQuery Hub
func (m *CloudQueryModule) Login(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithExec([]string{"cloudquery", "login"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery login: %w", err)
	}

	return output, nil
}

// Logout from CloudQuery Hub
func (m *CloudQueryModule) Logout(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithExec([]string{"cloudquery", "logout"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery logout: %w", err)
	}

	return output, nil
}

// InstallPlugin installs a CloudQuery plugin
func (m *CloudQueryModule) InstallPlugin(ctx context.Context, pluginName string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithExec([]string{"cloudquery", "plugin", "install", pluginName})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install cloudquery plugin: %w", err)
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

// SyncWithOptions syncs cloud resources with configurable options
func (m *CloudQueryModule) SyncWithOptions(ctx context.Context, configPath string, logLevel string, noMigrate bool) (string, error) {
	args := []string{"cloudquery", "sync", "/config"}
	
	if logLevel != "" {
		args = append(args, "--log-level", logLevel)
	}
	if noMigrate {
		args = append(args, "--no-migrate")
	}

	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithFile("/config", m.client.Host().File(configPath)).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery sync with options: %w", err)
	}

	return output, nil
}

// MigrateWithOptions migrates schema with configurable options
func (m *CloudQueryModule) MigrateWithOptions(ctx context.Context, configPath string, logLevel string) (string, error) {
	args := []string{"cloudquery", "migrate", "/config"}
	
	if logLevel != "" {
		args = append(args, "--log-level", logLevel)
	}

	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithFile("/config", m.client.Host().File(configPath)).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery migrate with options: %w", err)
	}

	return output, nil
}

// Switch between CloudQuery contexts or configurations
func (m *CloudQueryModule) Switch(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/cloudquery/cloudquery:latest").
		WithExec([]string{"cloudquery", "switch"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run cloudquery switch: %w", err)
	}

	return output, nil
}
