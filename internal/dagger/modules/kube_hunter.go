package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// KubeHunterModule runs kube-hunter for Kubernetes penetration testing
type KubeHunterModule struct {
	client *dagger.Client
	name   string
}

// NewKubeHunterModule creates a new kube-hunter module
func NewKubeHunterModule(client *dagger.Client) *KubeHunterModule {
	return &KubeHunterModule{
		client: client,
		name:   "kube-hunter",
	}
}

// ScanRemote scans remote Kubernetes cluster
func (m *KubeHunterModule) ScanRemote(ctx context.Context, remote string) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec([]string{
			"kube-hunter",
			"--remote", remote,
			"--report", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run kube-hunter remote scan: %w", err)
	}

	return output, nil
}

// ScanCIDR scans CIDR range for Kubernetes clusters
func (m *KubeHunterModule) ScanCIDR(ctx context.Context, cidr string) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec([]string{
			"kube-hunter",
			"--cidr", cidr,
			"--report", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run kube-hunter CIDR scan: %w", err)
	}

	return output, nil
}

// ScanInterface scans network interface
func (m *KubeHunterModule) ScanInterface(ctx context.Context, networkInterface string) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec([]string{
			"kube-hunter",
			"--interface", networkInterface,
			"--report", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run kube-hunter interface scan: %w", err)
	}

	return output, nil
}

// ScanPod runs kube-hunter as pod in cluster
func (m *KubeHunterModule) ScanPod(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-hunter:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"kube-hunter",
		"--pod",
		"--report", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run kube-hunter pod scan: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of kube-hunter
func (m *KubeHunterModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec([]string{"kube-hunter", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get kube-hunter version: %w", err)
	}

	return output, nil
}
