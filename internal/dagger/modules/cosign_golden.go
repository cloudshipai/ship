package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// CosignGoldenModule runs enhanced Cosign operations for golden image pipelines
type CosignGoldenModule struct {
	client *dagger.Client
	name   string
}

const cosignGoldenBinary = "/ko-app/cosign"

// NewCosignGoldenModule creates a new Cosign Golden module
func NewCosignGoldenModule(client *dagger.Client) *CosignGoldenModule {
	return &CosignGoldenModule{
		client: client,
		name:   "cosign-golden",
	}
}

// SignKeyless signs container image using keyless OIDC authentication
func (m *CosignGoldenModule) SignKeyless(ctx context.Context, imageRef string, identity string, issuer string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	cmd := []string{cosignGoldenBinary, "sign", "--yes", imageRef}

	if identity != "" {
		cmd = append(cmd, "--certificate-identity", identity)
	}

	if issuer != "" {
		cmd = append(cmd, "--certificate-oidc-issuer", issuer)
	}

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to sign image with keyless authentication: %w", err)
	}

	return output, nil
}

// VerifyKeyless verifies container image signature using keyless verification
func (m *CosignGoldenModule) VerifyKeyless(ctx context.Context, imageRef string, identity string, issuer string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	cmd := []string{cosignGoldenBinary, "verify", "--output=json", imageRef}

	if identity != "" {
		cmd = append(cmd, "--certificate-identity", identity)
	}

	if issuer != "" {
		cmd = append(cmd, "--certificate-oidc-issuer", issuer)
	}

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to verify image signature: %w", err)
	}

	return output, nil
}

// SignGoldenPipeline signs golden image with pipeline-specific metadata
func (m *CosignGoldenModule) SignGoldenPipeline(ctx context.Context, imageRef string, buildMetadata map[string]string, securityAttestations map[string]string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	cmd := []string{cosignGoldenBinary, "sign", "--yes"}

	// Add build metadata annotations
	for key, value := range buildMetadata {
		cmd = append(cmd, "-a", fmt.Sprintf("build.%s=%s", key, value))
	}

	// Add security attestation annotations
	for key, value := range securityAttestations {
		cmd = append(cmd, "-a", fmt.Sprintf("security.%s=%s", key, value))
	}

	// Add pipeline-specific annotations
	cmd = append(cmd, "-a", "org.opencontainers.image.title=Golden Container Image")
	cmd = append(cmd, "-a", "io.cosign.signed=true")
	cmd = append(cmd, "-a", "pipeline.type=golden-ami-container")

	cmd = append(cmd, imageRef)

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to sign golden pipeline image: %w", err)
	}

	return output, nil
}

// GenerateAttestation generates and signs SLSA provenance or custom attestation
func (m *CosignGoldenModule) GenerateAttestation(ctx context.Context, imageRef string, attestationType string, predicateData string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1").
		WithNewFile("/predicate.json", predicateData)

	cmd := []string{cosignGoldenBinary, "attest", "--yes", "--predicate", "/predicate.json", "--type", attestationType, imageRef}

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate attestation: %w", err)
	}

	return output, nil
}

// VerifyAttestation verifies attestations attached to golden image
func (m *CosignGoldenModule) VerifyAttestation(ctx context.Context, imageRef string, attestationType string, identity string, issuer string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	cmd := []string{cosignGoldenBinary, "verify-attestation", "--output=json", "--type", attestationType}

	if identity != "" {
		cmd = append(cmd, "--certificate-identity", identity)
	}

	if issuer != "" {
		cmd = append(cmd, "--certificate-oidc-issuer", issuer)
	}

	cmd = append(cmd, imageRef)

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to verify attestation: %w", err)
	}

	return output, nil
}

// CopySignatures copies signatures and attestations from one image to another
func (m *CosignGoldenModule) CopySignatures(ctx context.Context, sourceRef string, destinationRef string, force bool) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest")

	cmd := []string{cosignGoldenBinary, "copy", sourceRef, destinationRef}

	if force {
		cmd = append(cmd, "--force")
	}

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to copy signatures: %w", err)
	}

	return output, nil
}

// TreeView displays signature and attestation tree for golden image
func (m *CosignGoldenModule) TreeView(ctx context.Context, imageRef string, outputFormat string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest")

	cmd := []string{cosignGoldenBinary, "tree", imageRef}

	if outputFormat == "json" {
		cmd = append(cmd, "--json")
	}

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate tree view: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Cosign
func (m *CosignGoldenModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithExec([]string{"cosign", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cosign version: %w", err)
	}

	return output, nil
}

// UploadEBPF uploads eBPF program to OCI registry
func (m *CosignGoldenModule) UploadEBPF(ctx context.Context, ebpfPath string, registryURL string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/ebpf", m.client.Host().File(ebpfPath)).
		WithExec([]string{"cosign", "upload", "blob", "-f", "/tmp/ebpf", registryURL})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to upload eBPF program: %w", err)
	}

	return output, nil
}

// AttestWithType creates attestation with specific predicate type
func (m *CosignGoldenModule) AttestWithType(ctx context.Context, imageRef string, predicateType string, predicateFile string, keyPath string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/predicate.json", m.client.Host().File(predicateFile)).
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	cmd := []string{cosignGoldenBinary, "attest", "--predicate", "/tmp/predicate.json", "--type", predicateType}

	if keyPath != "" {
		container = container.WithFile("/tmp/private.key", m.client.Host().File(keyPath))
		cmd = append(cmd, "--key", "/tmp/private.key")
	}

	cmd = append(cmd, imageRef)
	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to attest with type: %w", err)
	}

	return output, nil
}

// VerifyAttestationAdvanced verifies attestation with specific type and policy
func (m *CosignGoldenModule) VerifyAttestationAdvanced(ctx context.Context, imageRef string, attestationType string, policyPath string, keyPath string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	cmd := []string{cosignGoldenBinary, "verify-attestation", "--output=json"}

	if attestationType != "" {
		cmd = append(cmd, "--type", attestationType)
	}
	if policyPath != "" {
		container = container.WithFile("/tmp/policy.json", m.client.Host().File(policyPath))
		cmd = append(cmd, "--policy", "/tmp/policy.json")
	}
	if keyPath != "" {
		container = container.WithFile("/tmp/public.key", m.client.Host().File(keyPath))
		cmd = append(cmd, "--key", "/tmp/public.key")
	}

	cmd = append(cmd, imageRef)
	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to verify attestation advanced: %w", err)
	}

	return output, nil
}

// VerifyOffline verifies signatures using offline bundle
func (m *CosignGoldenModule) VerifyOffline(ctx context.Context, imageRef string, bundlePath string, certIdentity string, certOidcIssuer string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/bundle", m.client.Host().File(bundlePath)).
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	cmd := []string{cosignGoldenBinary, "verify", "--bundle", "/tmp/bundle", "--output=json"}

	if certIdentity != "" {
		cmd = append(cmd, "--certificate-identity", certIdentity)
	}
	if certOidcIssuer != "" {
		cmd = append(cmd, "--certificate-oidc-issuer", certOidcIssuer)
	}

	cmd = append(cmd, imageRef)
	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to verify offline: %w", err)
	}

	return output, nil
}

// Triangulate gets signature image reference for a given image
func (m *CosignGoldenModule) Triangulate(ctx context.Context, imageRef string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithExec([]string{"cosign", "triangulate", imageRef})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to triangulate signatures: %w", err)
	}

	return output, nil
}

// Clean cleans signatures from a given image
func (m *CosignGoldenModule) Clean(ctx context.Context, imageRef string, cleanType string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest")

	cmd := []string{cosignGoldenBinary, "clean"}

	if cleanType != "" {
		cmd = append(cmd, "--type", cleanType)
	}

	cmd = append(cmd, imageRef)
	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to clean signatures: %w", err)
	}

	return output, nil
}

// VerifyIdentity verifies container image signature with certificate identity
func (m *CosignGoldenModule) VerifyIdentity(ctx context.Context, imageRef string, certIdentity string, certIdentityRegexp string, certOidcIssuer string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	cmd := []string{cosignGoldenBinary, "verify", "--output=json"}

	if certIdentity != "" {
		cmd = append(cmd, "--certificate-identity", certIdentity)
	}
	if certIdentityRegexp != "" {
		cmd = append(cmd, "--certificate-identity-regexp", certIdentityRegexp)
	}
	if certOidcIssuer != "" {
		cmd = append(cmd, "--certificate-oidc-issuer", certOidcIssuer)
	}

	cmd = append(cmd, imageRef)
	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to verify identity: %w", err)
	}

	return output, nil
}
