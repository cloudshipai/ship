package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// KubeBenchModule runs kube-bench for Kubernetes security benchmarks
type KubeBenchModule struct {
	client *dagger.Client
	name   string
}

// NewKubeBenchModule creates a new kube-bench module
func NewKubeBenchModule(client *dagger.Client) *KubeBenchModule {
	return &KubeBenchModule{
		client: client,
		name:   "kube-bench",
	}
}

// RunBenchmark runs CIS Kubernetes benchmark
func (m *KubeBenchModule) RunBenchmark(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-bench:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"kube-bench",
		"--json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run kube-bench: %w", err)
	}

	return output, nil
}

// RunMasterBenchmark runs benchmark for master node
func (m *KubeBenchModule) RunMasterBenchmark(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-bench:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"kube-bench",
		"master",
		"--json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run kube-bench master: %w", err)
	}

	return output, nil
}

// RunNodeBenchmark runs benchmark for worker node
func (m *KubeBenchModule) RunNodeBenchmark(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-bench:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"kube-bench",
		"node",
		"--json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run kube-bench node: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of kube-bench
func (m *KubeBenchModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-bench:latest").
		WithExec([]string{"kube-bench", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get kube-bench version: %w", err)
	}

	return output, nil
}
