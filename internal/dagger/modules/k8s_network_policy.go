package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// K8sNetworkPolicyModule runs Kubernetes network policy tools
type K8sNetworkPolicyModule struct {
	client *dagger.Client
	name   string
}

// NewK8sNetworkPolicyModule creates a new Kubernetes network policy module
func NewK8sNetworkPolicyModule(client *dagger.Client) *K8sNetworkPolicyModule {
	return &K8sNetworkPolicyModule{
		client: client,
		name:   "k8s-network-policy",
	}
}

// AnalyzePolicies analyzes network policies in the cluster
func (m *K8sNetworkPolicyModule) AnalyzePolicies(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"kubectl", "get", "networkpolicies",
		"--namespace", namespace,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to analyze network policies: %w", err)
	}

	return output, nil
}

// ValidatePolicy validates a network policy
func (m *K8sNetworkPolicyModule) ValidatePolicy(ctx context.Context, policyPath string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("kubectl:latest").
		WithFile("/policy.yaml", m.client.Host().File(policyPath))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"kubectl", "apply",
		"--dry-run=client",
		"--validate=true",
		"-f", "/policy.yaml",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate network policy: %w", err)
	}

	return output, nil
}

// TestConnectivity tests network connectivity between pods
func (m *K8sNetworkPolicyModule) TestConnectivity(ctx context.Context, sourceNamespace string, targetNamespace string, targetService string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("nicolaka/netshoot:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("nc -zv %s.%s.svc.cluster.local 80 || echo 'Connection failed'", targetService, targetNamespace),
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to test connectivity: %w", err)
	}

	return output, nil
}
