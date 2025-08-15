package modules

import (
	"context"
	"dagger.io/dagger"
)

// SigstorePolicyControllerModule provides Sigstore Policy Controller capabilities
type SigstorePolicyControllerModule struct {
	Client *dagger.Client
}

// NewSigstorePolicyControllerModule creates a new Sigstore Policy Controller module
func NewSigstorePolicyControllerModule(client *dagger.Client) *SigstorePolicyControllerModule {
	return &SigstorePolicyControllerModule{
		Client: client,
	}
}

// ValidatePolicy validates a ClusterImagePolicy
func (m *SigstorePolicyControllerModule) ValidatePolicy(ctx context.Context, policyPath string) (string, error) {
	policyFile := m.Client.Host().File(policyPath)
	
	result := m.Client.Container().
		From("ghcr.io/sigstore/policy-controller:latest").
		WithFile("/app/policy.yaml", policyFile).
		WithExec([]string{
			"policy-controller", "validate", "/app/policy.yaml",
		})

	return result.Stdout(ctx)
}

// TestPolicy tests a policy against an image
func (m *SigstorePolicyControllerModule) TestPolicy(ctx context.Context, policyPath string, imageName string) (string, error) {
	policyFile := m.Client.Host().File(policyPath)
	
	result := m.Client.Container().
		From("ghcr.io/sigstore/policy-controller:latest").
		WithFile("/app/policy.yaml", policyFile).
		WithExec([]string{
			"policy-controller", "test",
			"--policy", "/app/policy.yaml",
			"--image", imageName,
		})

	return result.Stdout(ctx)
}

// VerifySignature verifies an image signature against policies
func (m *SigstorePolicyControllerModule) VerifySignature(ctx context.Context, imageName string, publicKeyPath string) (string, error) {
	var result *dagger.Container
	
	if publicKeyPath != "" {
		publicKey := m.Client.Host().File(publicKeyPath)
		result = m.Client.Container().
			From("ghcr.io/sigstore/policy-controller:latest").
			WithFile("/app/pubkey.pem", publicKey).
			WithExec([]string{
				"policy-controller", "verify",
				"--image", imageName,
				"--public-key", "/app/pubkey.pem",
			})
	} else {
		result = m.Client.Container().
			From("ghcr.io/sigstore/policy-controller:latest").
			WithExec([]string{
				"policy-controller", "verify",
				"--image", imageName,
			})
	}

	return result.Stdout(ctx)
}

// GeneratePolicyTemplate generates a policy template
func (m *SigstorePolicyControllerModule) GeneratePolicyTemplate(ctx context.Context, namespace string, keyRef string) (string, error) {
	result := m.Client.Container().
		From("ghcr.io/sigstore/policy-controller:latest").
		WithExec([]string{
			"policy-controller", "generate",
			"--namespace", namespace,
			"--key-ref", keyRef,
		})

	return result.Stdout(ctx)
}

// ValidateManifest validates a Kubernetes manifest against signing policies
func (m *SigstorePolicyControllerModule) ValidateManifest(ctx context.Context, manifestPath string, policyPath string) (string, error) {
	manifestFile := m.Client.Host().File(manifestPath)
	policyFile := m.Client.Host().File(policyPath)
	
	result := m.Client.Container().
		From("ghcr.io/sigstore/policy-controller:latest").
		WithFile("/app/manifest.yaml", manifestFile).
		WithFile("/app/policy.yaml", policyFile).
		WithExec([]string{
			"policy-controller", "validate-manifest",
			"--manifest", "/app/manifest.yaml",
			"--policy", "/app/policy.yaml",
		})

	return result.Stdout(ctx)
}

// CheckCompliance checks if images in a directory comply with policies
func (m *SigstorePolicyControllerModule) CheckCompliance(ctx context.Context, manifestsPath string, policyPath string) (string, error) {
	manifestsDir := m.Client.Host().Directory(manifestsPath)
	policyFile := m.Client.Host().File(policyPath)
	
	result := m.Client.Container().
		From("ghcr.io/sigstore/policy-controller:latest").
		WithDirectory("/app/manifests", manifestsDir).
		WithFile("/app/policy.yaml", policyFile).
		WithWorkdir("/app/manifests").
		WithExec([]string{
			"policy-controller", "check-compliance",
			"--manifests", ".",
			"--policy", "/app/policy.yaml",
		})

	return result.Stdout(ctx)
}

// ListPolicies lists all available policies in a directory
func (m *SigstorePolicyControllerModule) ListPolicies(ctx context.Context, policiesPath string) (string, error) {
	policiesDir := m.Client.Host().Directory(policiesPath)
	
	result := m.Client.Container().
		From("ghcr.io/sigstore/policy-controller:latest").
		WithDirectory("/app/policies", policiesDir).
		WithWorkdir("/app/policies").
		WithExec([]string{
			"policy-controller", "list-policies", ".",
		})

	return result.Stdout(ctx)
}

// AuditImages audits images for signing compliance
func (m *SigstorePolicyControllerModule) AuditImages(ctx context.Context, namespace string, policyPath string) (string, error) {
	policyFile := m.Client.Host().File(policyPath)
	
	result := m.Client.Container().
		From("ghcr.io/sigstore/policy-controller:latest").
		WithFile("/app/policy.yaml", policyFile).
		WithExec([]string{
			"policy-controller", "audit",
			"--namespace", namespace,
			"--policy", "/app/policy.yaml",
		})

	return result.Stdout(ctx)
}
