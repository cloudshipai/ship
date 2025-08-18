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

// NetfetchScan scans network policies using netfetch
func (m *K8sNetworkPolicyModule) NetfetchScan(ctx context.Context, kubeconfig string, outputFormat string) (string, error) {
	args := []string{"netfetch", "scan"}
	if outputFormat != "" {
		args = append(args, "--output", outputFormat)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/deggja/netfetch/releases/latest/download/netfetch-linux-amd64 -o /usr/local/bin/netfetch && chmod +x /usr/local/bin/netfetch"})

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run netfetch scan: %w", err)
	}

	return output, nil
}

// NetpolEval evaluates network policies using netpol-analyzer
func (m *K8sNetworkPolicyModule) NetpolEval(ctx context.Context, kubeconfig string, namespace string) (string, error) {
	args := []string{"netpol", "eval"}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/np-guard/netpol-analyzer/releases/latest/download/netpol-linux-amd64 -o /usr/local/bin/netpol && chmod +x /usr/local/bin/netpol"})

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run netpol eval: %w", err)
	}

	return output, nil
}

// NetpolList lists network policies using netpol-analyzer
func (m *K8sNetworkPolicyModule) NetpolList(ctx context.Context, kubeconfig string, namespace string) (string, error) {
	args := []string{"netpol", "list"}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/np-guard/netpol-analyzer/releases/latest/download/netpol-linux-amd64 -o /usr/local/bin/netpol && chmod +x /usr/local/bin/netpol"})

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run netpol list: %w", err)
	}

	return output, nil
}
