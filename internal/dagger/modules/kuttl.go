package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// KuttlModule runs KUTTL for Kubernetes testing
type KuttlModule struct {
	client *dagger.Client
	name   string
}

// NewKuttlModule creates a new KUTTL module
func NewKuttlModule(client *dagger.Client) *KuttlModule {
	return &KuttlModule{
		client: client,
		name:   "kuttl",
	}
}

// RunTest runs KUTTL tests
func (m *KuttlModule) RunTest(ctx context.Context, testPath string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("kudobuilder/kuttl:latest").
		WithDirectory("/tests", m.client.Host().Directory(testPath))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"kubectl", "kuttl",
		"test",
		"/tests",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run kuttl tests: %w", err)
	}

	return output, nil
}

// ValidateTest validates test configuration
func (m *KuttlModule) ValidateTest(ctx context.Context, testPath string) (string, error) {
	container := m.client.Container().
		From("kudobuilder/kuttl:latest").
		WithDirectory("/tests", m.client.Host().Directory(testPath)).
		WithExec([]string{
			"kubectl", "kuttl",
			"test",
			"--dry-run",
			"/tests",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate kuttl tests: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of KUTTL
func (m *KuttlModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("kudobuilder/kuttl:latest").
		WithExec([]string{"kubectl-kuttl", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get kuttl version: %w", err)
	}

	return output, nil
}
