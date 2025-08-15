package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// LitmusModule runs Litmus for chaos engineering
type LitmusModule struct {
	client *dagger.Client
	name   string
}

// NewLitmusModule creates a new Litmus module
func NewLitmusModule(client *dagger.Client) *LitmusModule {
	return &LitmusModule{
		client: client,
		name:   "litmus",
	}
}

// CreateExperiment creates a chaos experiment
func (m *LitmusModule) CreateExperiment(ctx context.Context, experimentPath string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest").
		WithFile("/experiment.yaml", m.client.Host().File(experimentPath))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"litmusctl",
		"create",
		"experiment",
		"-f", "/experiment.yaml",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create litmus experiment: %w", err)
	}

	return output, nil
}

// GetExperiments lists chaos experiments
func (m *LitmusModule) GetExperiments(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"litmusctl",
		"get",
		"experiments",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get litmus experiments: %w", err)
	}

	return output, nil
}

// GetChaosResults gets chaos experiment results
func (m *LitmusModule) GetChaosResults(ctx context.Context, experimentName string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"litmusctl",
		"get",
		"chaosresults",
		experimentName,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chaos results: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Litmus
func (m *LitmusModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest").
		WithExec([]string{"litmusctl", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get litmus version: %w", err)
	}

	return output, nil
}
