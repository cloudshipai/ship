package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// GoldilocksModule runs Goldilocks for Kubernetes resource recommendations
type GoldilocksModule struct {
	client *dagger.Client
	name   string
}

// NewGoldilocksModule creates a new Goldilocks module
func NewGoldilocksModule(client *dagger.Client) *GoldilocksModule {
	return &GoldilocksModule{
		client: client,
		name:   "goldilocks",
	}
}

// GetRecommendations gets resource recommendations
func (m *GoldilocksModule) GetRecommendations(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("us-docker.pkg.dev/fairwinds-ops/oss/goldilocks:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"goldilocks",
		"recommendations",
		"--namespace", namespace,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get goldilocks recommendations: %w", err)
	}

	return output, nil
}

// CreateVPA creates Vertical Pod Autoscaler resources
func (m *GoldilocksModule) CreateVPA(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("us-docker.pkg.dev/fairwinds-ops/oss/goldilocks:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"goldilocks",
		"create-vpas",
		"--namespace", namespace,
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create VPAs: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Goldilocks
func (m *GoldilocksModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("us-docker.pkg.dev/fairwinds-ops/oss/goldilocks:latest").
		WithExec([]string{"goldilocks", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get goldilocks version: %w", err)
	}

	return output, nil
}
