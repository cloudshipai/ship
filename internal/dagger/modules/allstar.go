package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// AllstarModule provides information about Allstar (GitHub App, not a CLI tool)
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

// GetInfo provides information about Allstar since it's a GitHub App, not a CLI tool
func (m *AllstarModule) GetInfo(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{
			"echo", 
			"IMPORTANT: Allstar is a GitHub App (not a CLI tool). Install from https://github.com/apps/allstar and configure via .allstar/ repository.",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get allstar info: %w", err)
	}

	return output, nil
}