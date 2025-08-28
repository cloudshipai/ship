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
		From("golang:1.21-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithExec([]string{"go", "install", "github.com/FairwindsOps/goldilocks/cmd/goldilocks@latest"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithEnvVariable("PATH", "/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"goldilocks",
		"recommendations",
		"--namespace", namespace,
		"--output", "json",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
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
		From("golang:1.21-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithExec([]string{"go", "install", "github.com/FairwindsOps/goldilocks/cmd/goldilocks@latest"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithEnvVariable("PATH", "/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"goldilocks",
		"create-vpas",
		"--namespace", namespace,
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create VPAs: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Goldilocks
func (m *GoldilocksModule) GetVersion(ctx context.Context) (string, error) {
	// Goldilocks doesn't have a simple version command, return the module info
	return "goldilocks@latest", nil
}

// InstallHelm installs Goldilocks using Helm
func (m *GoldilocksModule) InstallHelm(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	if namespace == "" {
		namespace = "goldilocks"
	}

	container := m.client.Container().
		From("alpine/helm:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"helm", "repo", "add", "fairwinds-stable", "https://charts.fairwinds.com/stable",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	}).WithExec([]string{
		"helm", "install", "goldilocks", "fairwinds-stable/goldilocks",
		"--namespace", namespace,
		"--create-namespace",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install goldilocks via helm: %w", err)
	}

	return output, nil
}

// EnableNamespace enables Goldilocks for a namespace
func (m *GoldilocksModule) EnableNamespace(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"kubectl", "label", "ns", namespace, "goldilocks.fairwinds.com/enabled=true",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to enable namespace for goldilocks: %w", err)
	}

	return output, nil
}

// Uninstall removes Goldilocks using Helm
func (m *GoldilocksModule) Uninstall(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	if namespace == "" {
		namespace = "goldilocks"
	}

	container := m.client.Container().
		From("alpine/helm:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"helm", "uninstall", "goldilocks", "--namespace", namespace,
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to uninstall goldilocks: %w", err)
	}

	return output, nil
}
