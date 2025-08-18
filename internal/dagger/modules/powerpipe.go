package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// PowerpipeModule runs Powerpipe for security dashboards
type PowerpipeModule struct {
	client *dagger.Client
	name   string
}

const powerpipeBinary = "/usr/local/bin/powerpipe"

// NewPowerpipeModule creates a new Powerpipe module
func NewPowerpipeModule(client *dagger.Client) *PowerpipeModule {
	return &PowerpipeModule{
		client: client,
		name:   powerpipeBinary,
	}
}

// RunBenchmark runs a security benchmark
func (m *PowerpipeModule) RunBenchmark(ctx context.Context, benchmark string, modPath string) (string, error) {
	container := m.client.Container().
		From("turbot/powerpipe:latest")

	if modPath != "" {
		container = container.WithDirectory("/mod", m.client.Host().Directory(modPath))
	}

	container = container.WithExec([]string{
		powerpipeBinary,
		"benchmark", "run", benchmark,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run powerpipe benchmark: %w", err)
	}

	return output, nil
}

// RunControl runs a specific control
func (m *PowerpipeModule) RunControl(ctx context.Context, control string, modPath string) (string, error) {
	container := m.client.Container().
		From("turbot/powerpipe:latest")

	if modPath != "" {
		container = container.WithDirectory("/mod", m.client.Host().Directory(modPath))
	}

	container = container.WithExec([]string{
		powerpipeBinary,
		"control", "run", control,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run powerpipe control: %w", err)
	}

	return output, nil
}

// ListBenchmarks lists available benchmarks
func (m *PowerpipeModule) ListBenchmarks(ctx context.Context, modPath string) (string, error) {
	container := m.client.Container().
		From("turbot/powerpipe:latest")

	if modPath != "" {
		container = container.WithDirectory("/mod", m.client.Host().Directory(modPath))
	}

	container = container.WithExec([]string{
		powerpipeBinary,
		"benchmark", "list",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list powerpipe benchmarks: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Powerpipe
func (m *PowerpipeModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("turbot/powerpipe:latest").
		WithExec([]string{powerpipeBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get powerpipe version: %w", err)
	}

	return output, nil
}

// RunQuery executes a Powerpipe query
func (m *PowerpipeModule) RunQuery(ctx context.Context, query string, modPath string) (string, error) {
	container := m.client.Container().
		From("turbot/powerpipe:latest")

	if modPath != "" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(modPath)).
			WithWorkdir("/workspace")
	}

	container = container.WithExec([]string{
		powerpipeBinary,
		"query", "run", query,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run powerpipe query: %w", err)
	}

	return output, nil
}

// ListQueries lists available queries
func (m *PowerpipeModule) ListQueries(ctx context.Context, modPath string) (string, error) {
	container := m.client.Container().
		From("turbot/powerpipe:latest")

	if modPath != "" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(modPath)).
			WithWorkdir("/workspace")
	}

	container = container.WithExec([]string{
		powerpipeBinary,
		"query", "list",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list powerpipe queries: %w", err)
	}

	return output, nil
}

// StartServer starts the Powerpipe server
func (m *PowerpipeModule) StartServer(ctx context.Context, modPath string, port int) (string, error) {
	container := m.client.Container().
		From("turbot/powerpipe:latest")

	if modPath != "" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(modPath)).
			WithWorkdir("/workspace")
	}

	args := []string{powerpipeBinary, "server"}
	if port > 0 {
		args = append(args, "--port", fmt.Sprintf("%d", port))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start powerpipe server: %w", err)
	}

	return output, nil
}

// ListDashboards lists available dashboards
func (m *PowerpipeModule) ListDashboards(ctx context.Context, modPath string) (string, error) {
	container := m.client.Container().
		From("turbot/powerpipe:latest")

	if modPath != "" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(modPath)).
			WithWorkdir("/workspace")
	}

	container = container.WithExec([]string{
		powerpipeBinary,
		"dashboard", "list",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list powerpipe dashboards: %w", err)
	}

	return output, nil
}
