package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// AllstarModule runs Allstar for GitHub security policy enforcement
type AllstarModule struct {
	client *dagger.Client
	name   string
}

// NewAllstarModule creates a new Allstar module
func NewAllstarModule(client *dagger.Client) *AllstarModule {
	return &AllstarModule{
		client: client,
		name:   "allstar",
	}
}

// ScanRepository scans a GitHub repository for security policies
func (m *AllstarModule) ScanRepository(ctx context.Context, repoURL string, configPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/ossf/allstar:latest")

	if configPath != "" {
		container = container.WithDirectory("/config", m.client.Host().Directory(configPath))
	}

	container = container.WithExec([]string{
		"allstar",
		"--repo", repoURL,
		"--config", "/config",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run allstar scan: %w", err)
	}

	return output, nil
}

// ValidateConfig validates Allstar configuration
func (m *AllstarModule) ValidateConfig(ctx context.Context, configPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/ossf/allstar:latest").
		WithDirectory("/config", m.client.Host().Directory(configPath)).
		WithExec([]string{
			"allstar",
			"--validate-config", "/config",
			"--output", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate allstar config: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Allstar
func (m *AllstarModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/ossf/allstar:latest").
		WithExec([]string{"allstar", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get allstar version: %w", err)
	}

	return output, nil
}
