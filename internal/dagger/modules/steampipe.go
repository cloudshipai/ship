package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// SteampipeModule runs Steampipe for cloud asset querying
type SteampipeModule struct {
	client *dagger.Client
	name   string
}

const steampipeBinary = "/usr/local/bin/steampipe"

// NewSteampipeModule creates a new Steampipe module
func NewSteampipeModule(client *dagger.Client) *SteampipeModule {
	return &SteampipeModule{
		client: client,
		name:   steampipeBinary,
	}
}

// Query executes a SQL query against cloud resources
func (m *SteampipeModule) Query(ctx context.Context, query string, plugin string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithExec([]string{
			steampipeBinary, "plugin", "install", plugin,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithExec([]string{
			steampipeBinary, "query",
			"--output", "json",
			query,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to execute steampipe query: %w", err)
	}

	return output, nil
}

// QueryFromFile executes queries from a file
func (m *SteampipeModule) QueryFromFile(ctx context.Context, queryFile string, plugin string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithFile("/query.sql", m.client.Host().File(queryFile)).
		WithExec([]string{
			steampipeBinary, "plugin", "install", plugin,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithExec([]string{
			steampipeBinary, "query",
			"--output", "json",
			"/query.sql",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to execute steampipe query from file: %w", err)
	}

	return output, nil
}

// ListPlugins lists available plugins
func (m *SteampipeModule) ListPlugins(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithExec([]string{
			steampipeBinary, "plugin", "list",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list steampipe plugins: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Steampipe
func (m *SteampipeModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithExec([]string{steampipeBinary, "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get steampipe version: %w", err)
	}

	return output, nil
}

// QueryInteractive starts an interactive query session
func (m *SteampipeModule) QueryInteractive(ctx context.Context, plugin string) (string, error) {
	// Interactive mode doesn't work well in containers
	// Return a helpful message instead
	return "Interactive mode is not supported in containerized environment. Use Query() or QueryFromFile() instead.", nil
}

// InstallPlugin installs a Steampipe plugin
func (m *SteampipeModule) InstallPlugin(ctx context.Context, plugin string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithExec([]string{
			steampipeBinary, "plugin", "install", plugin,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install plugin: %w", err)
	}

	return output, nil
}

// UpdatePlugin updates a Steampipe plugin
func (m *SteampipeModule) UpdatePlugin(ctx context.Context, plugin string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithExec([]string{
			steampipeBinary, "plugin", "update", plugin,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to update plugin: %w", err)
	}

	return output, nil
}

// UninstallPlugin uninstalls a Steampipe plugin
func (m *SteampipeModule) UninstallPlugin(ctx context.Context, plugin string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithExec([]string{
			steampipeBinary, "plugin", "uninstall", plugin,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to uninstall plugin: %w", err)
	}

	return output, nil
}

// StartService starts the Steampipe service
func (m *SteampipeModule) StartService(ctx context.Context, port int) (string, error) {
	args := []string{steampipeBinary, "service", "start"}
	if port > 0 {
		args = append(args, "--port", fmt.Sprintf("%d", port))
	}

	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Service mode requires persistent container
		return "Service mode requires persistent container. Use Query() for one-off queries.", nil
	}

	return output, nil
}

// GetServiceStatus gets the status of the Steampipe service
func (m *SteampipeModule) GetServiceStatus(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithExec([]string{
			steampipeBinary, "service", "status",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Service status not applicable in ephemeral containers
		return "Service not running (containerized environment). Use Query() for one-off queries.", nil
	}

	return output, nil
}

// StopService stops the Steampipe service
func (m *SteampipeModule) StopService(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/turbot/steampipe:latest").
		WithExec([]string{
			steampipeBinary, "service", "stop",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Service stop not applicable in ephemeral containers
		return "Service stop not needed (containerized environment).", nil
	}

	return output, nil
}
