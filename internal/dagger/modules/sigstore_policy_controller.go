package modules

import (
	"context"
	"dagger.io/dagger"
)

// SigstorePolicyControllerModule provides Sigstore Policy Controller capabilities
type SigstorePolicyControllerModule struct {
	Client *dagger.Client
}

const policyTesterBinary = "/usr/local/bin/policy-tester"
const kubectlBinary = "/usr/bin/kubectl"

// NewSigstorePolicyControllerModule creates a new Sigstore Policy Controller module
func NewSigstorePolicyControllerModule(client *dagger.Client) *SigstorePolicyControllerModule {
	return &SigstorePolicyControllerModule{
		Client: client,
	}
}

// ValidatePolicy validates a ClusterImagePolicy using locally built policy-tester
func (m *SigstorePolicyControllerModule) ValidatePolicy(ctx context.Context, policyPath string) (string, error) {
	policyFile := m.Client.Host().File(policyPath)
	
	result := m.Client.Container().
		From("golang:alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git", "make"}).
		WithExec([]string{"git", "clone", "https://github.com/sigstore/policy-controller.git", "/src"}).
		WithWorkdir("/src").
		WithExec([]string{"make", "policy-tester"}).
		WithFile("/app/policy.yaml", policyFile).
		WithExec([]string{
			"./policy-tester", "--policy", "/app/policy.yaml",
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
			"policy-tester", "--policy", "/app/policy.yaml", "--image", imageName,
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

// GetVersion gets the policy-tester version
func (m *SigstorePolicyControllerModule) GetVersion(ctx context.Context) (string, error) {
	result := m.Client.Container().
		From("ghcr.io/sigstore/policy-controller:latest").
		WithExec([]string{policyTesterBinary, "--version"})

	return result.Stdout(ctx)
}

// CreatePolicy creates a ClusterImagePolicy using kubectl
func (m *SigstorePolicyControllerModule) CreatePolicy(ctx context.Context, policyPath string, namespace string) (string, error) {
	policyFile := m.Client.Host().File(policyPath)
	
	args := []string{kubectlBinary, "apply", "-f", "/app/policy.yaml"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	result := m.Client.Container().
		From("bitnami/kubectl:latest").
		WithFile("/app/policy.yaml", policyFile).
		WithExec(args)

	return result.Stdout(ctx)
}

// ListClusterImagePolicies lists ClusterImagePolicies using kubectl
func (m *SigstorePolicyControllerModule) ListClusterImagePolicies(ctx context.Context, outputFormat string) (string, error) {
	args := []string{kubectlBinary, "get", "clusterimagepolicy"}
	if outputFormat != "" {
		args = append(args, "-o", outputFormat)
	}

	result := m.Client.Container().
		From("bitnami/kubectl:latest").
		WithExec(args)

	return result.Stdout(ctx)
}

// DeletePolicy deletes a ClusterImagePolicy using kubectl
func (m *SigstorePolicyControllerModule) DeletePolicy(ctx context.Context, policyName string) (string, error) {
	result := m.Client.Container().
		From("bitnami/kubectl:latest").
		WithExec([]string{kubectlBinary, "delete", "clusterimagepolicy", policyName})

	return result.Stdout(ctx)
}

// EnableNamespace enables Sigstore policy enforcement for namespace using kubectl
func (m *SigstorePolicyControllerModule) EnableNamespace(ctx context.Context, namespace string) (string, error) {
	result := m.Client.Container().
		From("bitnami/kubectl:latest").
		WithExec([]string{kubectlBinary, "label", "namespace", namespace, "policy.sigstore.dev/include=true"})

	return result.Stdout(ctx)
}

// DisableNamespace disables Sigstore policy enforcement for namespace using kubectl
func (m *SigstorePolicyControllerModule) DisableNamespace(ctx context.Context, namespace string) (string, error) {
	result := m.Client.Container().
		From("bitnami/kubectl:latest").
		WithExec([]string{kubectlBinary, "label", "namespace", namespace, "policy.sigstore.dev/exclude=true"})

	return result.Stdout(ctx)
}

// GetNamespaceStatus gets Sigstore policy enforcement status for namespace using kubectl
func (m *SigstorePolicyControllerModule) GetNamespaceStatus(ctx context.Context, namespace string) (string, error) {
	result := m.Client.Container().
		From("bitnami/kubectl:latest").
		WithExec([]string{kubectlBinary, "get", "namespace", namespace, "--show-labels"})

	return result.Stdout(ctx)
}

// DescribePolicy describes a ClusterImagePolicy using kubectl
func (m *SigstorePolicyControllerModule) DescribePolicy(ctx context.Context, policyName string) (string, error) {
	result := m.Client.Container().
		From("bitnami/kubectl:latest").
		WithExec([]string{kubectlBinary, "describe", "clusterimagepolicy", policyName})

	return result.Stdout(ctx)
}
