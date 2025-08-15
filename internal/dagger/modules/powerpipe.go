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

// NewPowerpipeModule creates a new Powerpipe module
func NewPowerpipeModule(client *dagger.Client) *PowerpipeModule {
	return &PowerpipeModule{
		client: client,
		name:   "powerpipe",
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
		"powerpipe",
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
		"powerpipe",
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
		"powerpipe",
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
		WithExec([]string{"powerpipe", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get powerpipe version: %w", err)
	}

	return output, nil
}
