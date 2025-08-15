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

	cmd := []string{"cosign", "sign", "--yes", imageRef}

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

	cmd := []string{"cosign", "verify", "--output=json", imageRef}

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

	cmd := []string{"cosign", "sign", "--yes"}

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

	cmd := []string{"cosign", "attest", "--yes", "--predicate", "/predicate.json", "--type", attestationType, imageRef}

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

	cmd := []string{"cosign", "verify-attestation", "--output=json", "--type", attestationType}

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

	cmd := []string{"cosign", "copy", sourceRef, destinationRef}

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

	cmd := []string{"cosign", "tree", imageRef}

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
