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

// NewSteampipeModule creates a new Steampipe module
func NewSteampipeModule(client *dagger.Client) *SteampipeModule {
	return &SteampipeModule{
		client: client,
		name:   "steampipe",
	}
}

// Query executes a SQL query against cloud resources
func (m *SteampipeModule) Query(ctx context.Context, query string, plugin string) (string, error) {
	container := m.client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{
			"steampipe", "plugin", "install", plugin,
		}).
		WithExec([]string{
			"steampipe", "query",
			"--output", "json",
			query,
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
		From("turbot/steampipe:latest").
		WithFile("/query.sql", m.client.Host().File(queryFile)).
		WithExec([]string{
			"steampipe", "plugin", "install", plugin,
		}).
		WithExec([]string{
			"steampipe", "query",
			"--output", "json",
			"/query.sql",
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
		From("turbot/steampipe:latest").
		WithExec([]string{
			"steampipe", "plugin", "list",
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
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get steampipe version: %w", err)
	}

	return output, nil
}
