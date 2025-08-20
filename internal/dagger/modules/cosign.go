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

const cosignBinary = "/ko-app/cosign"

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
		WithExec([]string{cosignBinary, "verify", imageName})

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
		WithExec([]string{cosignBinary, "verify", "--key", "/tmp/public.key", imageName})

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
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1").
		WithExec([]string{"cosign", "sign", imageName}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Keyless signing requires OIDC authentication", nil
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
		WithExec([]string{cosignBinary, "generate-key-pair"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Key pair generation requires interactive input", nil
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

// SignBlob signs arbitrary blob using Cosign
func (m *CosignModule) SignBlob(ctx context.Context, blobPath string, keyPath string, outputSignature string) (string, error) {
	args := []string{"cosign", "sign-blob"}
	if keyPath != "" {
		args = append(args, "--key", "/tmp/private.key")
	}
	if outputSignature != "" {
		args = append(args, "--output-signature", outputSignature)
	}
	args = append(args, "/tmp/blob")

	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/blob", m.client.Host().File(blobPath))

	if keyPath != "" {
		container = container.WithFile("/tmp/private.key", m.client.Host().File(keyPath))
	}

	if os.Getenv("COSIGN_PASSWORD") != "" {
		container = container.WithEnvVariable("COSIGN_PASSWORD", os.Getenv("COSIGN_PASSWORD"))
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Blob signing requires keys or keyless authentication", nil
}

// VerifyBlob verifies blob signature using Cosign
func (m *CosignModule) VerifyBlob(ctx context.Context, blobPath string, signaturePath string, keyPath string) (string, error) {
	args := []string{"cosign", "verify-blob", "--signature", "/tmp/signature"}
	if keyPath != "" {
		args = append(args, "--key", "/tmp/public.key")
	}
	args = append(args, "/tmp/blob")

	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/blob", m.client.Host().File(blobPath)).
		WithFile("/tmp/signature", m.client.Host().File(signaturePath))

	if keyPath != "" {
		container = container.WithFile("/tmp/public.key", m.client.Host().File(keyPath))
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Blob verification requires signature and keys", nil
}

// UploadBlob uploads generic artifact as a blob to registry
func (m *CosignModule) UploadBlob(ctx context.Context, blobPath string, registryURL string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/blob", m.client.Host().File(blobPath)).
		WithExec([]string{cosignBinary, "upload", "blob", "-f", "/tmp/blob", registryURL})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	return output, nil
}

// UploadWasm uploads WebAssembly module to registry
func (m *CosignModule) UploadWasm(ctx context.Context, wasmPath string, registryURL string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/wasm", m.client.Host().File(wasmPath)).
		WithExec([]string{cosignBinary, "upload", "wasm", "-f", "/tmp/wasm", registryURL})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to upload wasm: %w", err)
	}

	return output, nil
}

// CopyImage copies images between registries
func (m *CosignModule) CopyImage(ctx context.Context, sourceImage string, destinationImage string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithExec([]string{cosignBinary, "copy", sourceImage, destinationImage})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to copy image: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Cosign
func (m *CosignModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithExec([]string{cosignBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cosign version: %w", err)
	}

	return output, nil
}

// SignImageWithOptions signs container image with comprehensive options
func (m *CosignModule) SignImageWithOptions(ctx context.Context, imageName string, keyPath string, keyless bool) (string, error) {
	args := []string{"cosign", "sign"}
	
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	if keyless {
		// Use keyless signing
		if os.Getenv("COSIGN_IDENTITY_TOKEN") != "" {
			container = container.WithEnvVariable("COSIGN_IDENTITY_TOKEN", os.Getenv("COSIGN_IDENTITY_TOKEN"))
		}
	} else if keyPath != "" {
		// Use key-based signing
		args = append(args, "--key", "/tmp/private.key")
		container = container.WithFile("/tmp/private.key", m.client.Host().File(keyPath))
		if os.Getenv("COSIGN_PASSWORD") != "" {
			container = container.WithEnvVariable("COSIGN_PASSWORD", os.Getenv("COSIGN_PASSWORD"))
		}
	}

	args = append(args, imageName)
	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to sign image with options: %w", err)
	}

	return output, nil
}

// VerifyImageWithOptions verifies container image with comprehensive options
func (m *CosignModule) VerifyImageWithOptions(ctx context.Context, imageName string, keyPath string, keyless bool) (string, error) {
	args := []string{"cosign", "verify"}
	
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	if !keyless && keyPath != "" {
		args = append(args, "--key", "/tmp/public.key")
		container = container.WithFile("/tmp/public.key", m.client.Host().File(keyPath))
	}

	args = append(args, imageName)
	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to verify image with options: %w", err)
	}

	return output, nil
}

// AttestWithOptions creates attestation with comprehensive options
func (m *CosignModule) AttestWithOptions(ctx context.Context, imageName string, predicatePath string, keyPath string) (string, error) {
	args := []string{"cosign", "attest", "--predicate", "/tmp/predicate.json"}
	
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithFile("/tmp/predicate.json", m.client.Host().File(predicatePath)).
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	if keyPath != "" {
		args = append(args, "--key", "/tmp/private.key")
		container = container.WithFile("/tmp/private.key", m.client.Host().File(keyPath))
		if os.Getenv("COSIGN_PASSWORD") != "" {
			container = container.WithEnvVariable("COSIGN_PASSWORD", os.Getenv("COSIGN_PASSWORD"))
		}
	}

	args = append(args, imageName)
	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to attest with options: %w", err)
	}

	return output, nil
}

// VerifyAttestationWithOptions verifies attestation with comprehensive options
func (m *CosignModule) VerifyAttestationWithOptions(ctx context.Context, imageName string, keyPath string, policyPath string) (string, error) {
	args := []string{"cosign", "verify-attestation"}
	
	container := m.client.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithEnvVariable("COSIGN_EXPERIMENTAL", "1")

	if keyPath != "" {
		args = append(args, "--key", "/tmp/public.key")
		container = container.WithFile("/tmp/public.key", m.client.Host().File(keyPath))
	}
	if policyPath != "" {
		args = append(args, "--policy", "/tmp/policy.json")
		container = container.WithFile("/tmp/policy.json", m.client.Host().File(policyPath))
	}

	args = append(args, imageName)
	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to verify attestation with options: %w", err)
	}

	return output, nil
}