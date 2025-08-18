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

// Install installs Kyverno using Helm
func (m *KyvernoModule) Install(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	if namespace == "" {
		namespace = "kyverno"
	}

	container := m.client.Container().
		From("alpine/helm:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"helm", "repo", "add", "kyverno", "https://kyverno.github.io/kyverno/",
	}).WithExec([]string{
		"helm", "install", "kyverno", "kyverno/kyverno",
		"--namespace", namespace,
		"--create-namespace",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install kyverno: %w", err)
	}

	return output, nil
}

// ListPolicies lists Kyverno policies
func (m *KyvernoModule) ListPolicies(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	args := []string{"kubectl", "get", "policies"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "--all-namespaces")
	}
	args = append(args, "-o", "json")

	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list kyverno policies: %w", err)
	}

	return output, nil
}

// GetPolicyReports gets Kyverno policy reports
func (m *KyvernoModule) GetPolicyReports(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	args := []string{"kubectl", "get", "policyreports"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "--all-namespaces")
	}
	args = append(args, "-o", "json")

	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get policy reports: %w", err)
	}

	return output, nil
}

// GetStatus gets Kyverno installation status
func (m *KyvernoModule) GetStatus(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	if namespace == "" {
		namespace = "kyverno"
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"kubectl", "get", "pods", "-n", namespace, "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get kyverno status: %w", err)
	}

	return output, nil
}

// CreateClusterRole creates necessary RBAC cluster role for Kyverno
func (m *KyvernoModule) CreateClusterRole(ctx context.Context, kubeconfig string) (string, error) {
	clusterRoleYAML := `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kyverno-cluster-role
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kyverno-cluster-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kyverno-cluster-role
subjects:
- kind: ServiceAccount
  name: kyverno
  namespace: kyverno`

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithNewFile("/workspace/cluster-role.yaml", clusterRoleYAML)

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"kubectl", "apply", "-f", "/workspace/cluster-role.yaml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create cluster role: %w", err)
	}

	return output, nil
}
