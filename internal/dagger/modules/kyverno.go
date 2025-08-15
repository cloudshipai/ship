package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// KyvernoModule runs Kyverno for Kubernetes policy management
type KyvernoModule struct {
	client *dagger.Client
	name   string
}

// NewKyvernoModule creates a new Kyverno module
func NewKyvernoModule(client *dagger.Client) *KyvernoModule {
	return &KyvernoModule{
		client: client,
		name:   "kyverno",
	}
}

// ApplyPolicies applies Kyverno policies to cluster
func (m *KyvernoModule) ApplyPolicies(ctx context.Context, policiesPath string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/kyverno/kyverno-cli:latest").
		WithDirectory("/policies", m.client.Host().Directory(policiesPath))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"kyverno",
		"apply",
		"/policies",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to apply kyverno policies: %w", err)
	}

	return output, nil
}

// ValidatePolicies validates Kyverno policy syntax
func (m *KyvernoModule) ValidatePolicies(ctx context.Context, policiesPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/kyverno/kyverno-cli:latest").
		WithDirectory("/policies", m.client.Host().Directory(policiesPath)).
		WithExec([]string{
			"kyverno",
			"validate",
			"/policies",
			"--output", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate kyverno policies: %w", err)
	}

	return output, nil
}

// TestPolicies tests policies against resources
func (m *KyvernoModule) TestPolicies(ctx context.Context, policiesPath string, resourcesPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/kyverno/kyverno-cli:latest").
		WithDirectory("/policies", m.client.Host().Directory(policiesPath)).
		WithDirectory("/resources", m.client.Host().Directory(resourcesPath)).
		WithExec([]string{
			"kyverno",
			"test",
			"/policies",
			"--resource", "/resources",
			"--output", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to test kyverno policies: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Kyverno CLI
func (m *KyvernoModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/kyverno/kyverno-cli:latest").
		WithExec([]string{"kyverno", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get kyverno version: %w", err)
	}

	return output, nil
}
