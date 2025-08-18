package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// kyvernoMultitenantBinary is the path to the kyverno binary in the container
const kyvernoMultitenantBinary = "/usr/local/bin/kyverno"

// KyvernoMultitenantModule runs Kyverno for multi-tenant environments
type KyvernoMultitenantModule struct {
	client *dagger.Client
	name   string
}

// NewKyvernoMultitenantModule creates a new Kyverno multitenant module
func NewKyvernoMultitenantModule(client *dagger.Client) *KyvernoMultitenantModule {
	return &KyvernoMultitenantModule{
		client: client,
		name:   "kyverno-multitenant",
	}
}

// CreateTenantPolicies creates tenant isolation policies
func (m *KyvernoMultitenantModule) CreateTenantPolicies(ctx context.Context, tenantName string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/kyverno/kyverno-cli:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	// Generate tenant isolation policy
	policyYAML := fmt.Sprintf(`
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: tenant-%s-isolation
spec:
  validationFailureAction: enforce
  background: false
  rules:
  - name: restrict-cross-tenant-access
    match:
      any:
      - resources:
          kinds:
          - Pod
          namespaces:
          - "%s-*"
    validate:
      message: "Cross-tenant access not allowed"
      pattern:
        spec:
          securityContext:
            runAsNonRoot: true
`, tenantName, tenantName)

	container = container.
		WithNewFile("/policy.yaml", policyYAML).
		WithExec([]string{
			kyvernoMultitenantBinary,
			"apply", "/policy.yaml",
			"--output", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create tenant policies: %w", err)
	}

	return output, nil
}

// ValidateMultitenantSetup validates multi-tenant setup
func (m *KyvernoMultitenantModule) ValidateMultitenantSetup(ctx context.Context, tenantsConfig string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/kyverno/kyverno-cli:latest").
		WithFile("/tenants.yaml", m.client.Host().File(tenantsConfig))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		kyvernoMultitenantBinary,
		"test", "/tenants.yaml",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate multitenant setup: %w", err)
	}

	return output, nil
}

// CreateTenantNamespace creates a namespace for a tenant
func (m *KyvernoMultitenantModule) CreateTenantNamespace(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"kubectl", "create", "namespace", namespace})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create tenant namespace: %w", err)
	}

	return output, nil
}

// CreateResourceQuota creates resource quota for a tenant
func (m *KyvernoMultitenantModule) CreateResourceQuota(ctx context.Context, namespace string, cpuLimit string, memoryLimit string, kubeconfig string) (string, error) {
	quotaYaml := fmt.Sprintf(`apiVersion: v1
kind: ResourceQuota
metadata:
  name: tenant-quota
  namespace: %s
spec:
  hard:
    requests.cpu: %s
    requests.memory: %s
    limits.cpu: %s
    limits.memory: %s`, namespace, cpuLimit, memoryLimit, cpuLimit, memoryLimit)

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithNewFile("/tmp/quota.yaml", quotaYaml)

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"kubectl", "apply", "-f", "/tmp/quota.yaml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create resource quota: %w", err)
	}

	return output, nil
}

// ListTenantNamespaces lists namespaces with tenant labels
func (m *KyvernoMultitenantModule) ListTenantNamespaces(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"kubectl", "get", "namespaces", "-l", "tenant", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list tenant namespaces: %w", err)
	}

	return output, nil
}

// GetTenantPolicies gets policies for a specific tenant
func (m *KyvernoMultitenantModule) GetTenantPolicies(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"kubectl", "get", "policies", "-n", namespace, "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get tenant policies: %w", err)
	}

	return output, nil
}
