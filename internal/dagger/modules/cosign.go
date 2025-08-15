package modules

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// CosignModule runs Cosign for container signing and verification
type CosignModule struct {
	client *dagger.Client
	name   string
}

// NewCosignModule creates a new Cosign module
func NewCosignModule(client *dagger.Client) *CosignModule {
	return &CosignModule{
		client: client,
		name:   "cosign",
	}
}

// VerifyImage verifies a signed container image
func (m *CosignModule) VerifyImage(ctx context.Context, imageName string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1").
		WithExec([]string{"cosign", "verify", imageName})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to verify image: %w", err)
	}

	return output, nil
}

// VerifyImageWithKey verifies an image with a specific public key
func (m *CosignModule) VerifyImageWithKey(ctx context.Context, imageName string, publicKeyPath string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/public.key", m.client.Host().File(publicKeyPath)).
		WithExec([]string{"cosign", "verify", "--key", "/tmp/public.key", imageName})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to verify image with key: %w", err)
	}

	return output, nil
}

// SignImage signs a container image (requires authentication)
func (m *CosignModule) SignImage(ctx context.Context, imageName string, privateKeyPath string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/private.key", m.client.Host().File(privateKeyPath)).
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	if os.Getenv("COSIGN_PASSWORD") != "" {
		container = container.WithEnvVariable("COSIGN_PASSWORD", os.Getenv("COSIGN_PASSWORD"))
	}

	container = container.WithExec([]string{"cosign", "sign", "--key", "/tmp/private.key", imageName})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to sign image: %w", err)
	}

	return output, nil
}

// SignImageKeyless signs an image using keyless signing (OIDC)
func (m *CosignModule) SignImageKeyless(ctx context.Context, imageName string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	if os.Getenv("COSIGN_IDENTITY_TOKEN") != "" {
		container = container.WithEnvVariable("COSIGN_IDENTITY_TOKEN", os.Getenv("COSIGN_IDENTITY_TOKEN"))
	}

	container = container.WithExec([]string{"cosign", "sign", imageName})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to sign image keyless: %w", err)
	}

	return output, nil
}

// VerifyAttestation verifies attestations for an image
func (m *CosignModule) VerifyAttestation(ctx context.Context, imageName string, attestationType string) (string, error) {
	args := []string{"cosign", "verify-attestation", imageName}
	
	if attestationType != "" {
		args = append(args, "--type", attestationType)
	}

	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to verify attestation: %w", err)
	}

	return output, nil
}

// GenerateKeyPair generates a new signing key pair
func (m *CosignModule) GenerateKeyPair(ctx context.Context, outputDir string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithDirectory("/workspace", m.client.Host().Directory(outputDir)).
		WithWorkdir("/workspace").
		WithExec([]string{"cosign", "generate-key-pair"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to generate key pair: %w", err)
	}

	return output, nil
}

// AttestSBOM creates an SBOM attestation for an image
func (m *CosignModule) AttestSBOM(ctx context.Context, imageName string, sbomPath string, privateKeyPath string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/sbom.json", m.client.Host().File(sbomPath)).
		WithFile("/tmp/private.key", m.client.Host().File(privateKeyPath)).
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	if os.Getenv("COSIGN_PASSWORD") != "" {
		container = container.WithEnvVariable("COSIGN_PASSWORD", os.Getenv("COSIGN_PASSWORD"))
	}

	container = container.WithExec([]string{
		"cosign", "attest", "--predicate", "/tmp/sbom.json", 
		"--key", "/tmp/private.key", "--type", "spdx", imageName,
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to attest SBOM: %w", err)
	}

	return output, nil
}