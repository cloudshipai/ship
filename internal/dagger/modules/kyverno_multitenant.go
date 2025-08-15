package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

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
			"kyverno",
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
		"kyverno",
		"test", "/tenants.yaml",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate multitenant setup: %w", err)
	}

	return output, nil
}
