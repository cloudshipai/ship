package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

type SOPSModule struct {
	client *dagger.Client
}

const sopsBinary = "sops"

func NewSOPSModule(client *dagger.Client) *SOPSModule {
	return &SOPSModule{
		client: client,
	}
}

// EncryptFile encrypts a file using SOPS
func (m *SOPSModule) EncryptFile(ctx context.Context, filePath string, kmsArn string, pgpFingerprint string, agePublicKey string, outputFile string, inPlace bool) (string, error) {
	container := m.client.Container().
		From("mozilla/sops:latest").
		WithFile("/workspace/input", m.client.Host().File(filePath)).
		WithWorkdir("/workspace")

	args := []string{sopsBinary, "--encrypt"}

	if kmsArn != "" {
		args = append(args, "--kms", kmsArn)
	}
	if pgpFingerprint != "" {
		args = append(args, "--pgp", pgpFingerprint)
	}
	if agePublicKey != "" {
		args = append(args, "--age", agePublicKey)
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	if inPlace {
		args = append(args, "--in-place")
	}

	args = append(args, "/workspace/input")

	result := container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})
	output, err := result.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt file with SOPS: %w", err)
	}

	return output, nil
}

// DecryptFile decrypts a SOPS-encrypted file
func (m *SOPSModule) DecryptFile(ctx context.Context, filePath string, outputFile string) (string, error) {
	container := m.client.Container().
		From("mozilla/sops:latest").
		WithFile("/workspace/encrypted", m.client.Host().File(filePath)).
		WithWorkdir("/workspace")

	args := []string{sopsBinary, "--decrypt"}

	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}

	args = append(args, "/workspace/encrypted")

	result := container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})
	output, err := result.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt file with SOPS: %w", err)
	}

	return output, nil
}

// UpdateKeys rotates/updates encryption keys for SOPS files
func (m *SOPSModule) UpdateKeys(ctx context.Context, filePath string, addKms string, addPgp string, addAge string, inPlace bool) (string, error) {
	container := m.client.Container().
		From("mozilla/sops:latest").
		WithFile("/workspace/encrypted", m.client.Host().File(filePath)).
		WithWorkdir("/workspace")

	args := []string{sopsBinary, "--rotate"}

	if addKms != "" {
		args = append(args, "--add-kms", addKms)
	}
	if addPgp != "" {
		args = append(args, "--add-pgp", addPgp)
	}
	if addAge != "" {
		args = append(args, "--add-age", addAge)
	}
	if inPlace {
		args = append(args, "--in-place")
	}

	args = append(args, "/workspace/encrypted")

	result := container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})
	output, err := result.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to update keys with SOPS: %w", err)
	}

	return output, nil
}

// EditFile shows decrypted content for editing (interactive editing not supported in containers)
func (m *SOPSModule) EditFile(ctx context.Context, filePath string) (string, error) {
	container := m.client.Container().
		From("mozilla/sops:latest").
		WithFile("/workspace/encrypted", m.client.Host().File(filePath)).
		WithWorkdir("/workspace")

	// Note: Interactive editing is limited in containerized environments
	// This command will show the decrypted content
	args := []string{sopsBinary, "--decrypt", "/workspace/encrypted"}

	result := container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})
	output, err := result.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to show file content with SOPS: %w", err)
	}

	return output, nil
}

// GetVersion gets the SOPS version
func (m *SOPSModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("mozilla/sops:latest")

	result := container.WithExec([]string{sopsBinary, "--version"}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})
	output, err := result.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get SOPS version: %w", err)
	}

	return output, nil
}

// PublishKeys publishes/exports keys for sharing (shows public key info)
func (m *SOPSModule) PublishKeys(ctx context.Context, keyType string, keyPath string) (string, error) {
	container := m.client.Container().
		From("mozilla/sops:latest").
		WithWorkdir("/workspace")

	if keyPath != "" {
		container = container.WithFile("/workspace/keyfile", m.client.Host().File(keyPath))
	}

	args := []string{sopsBinary, "--help"}

	if keyType == "pgp" && keyPath != "" {
		// For PGP, show key information using gpg
		container = container.From("alpine:latest").
			WithExec([]string{"apk", "add", "--no-cache", "gnupg"})
		args = []string{"gpg", "--import", "/workspace/keyfile", "&&", "gpg", "--list-keys"}
	}

	result := container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})
	output, err := result.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to publish keys: %w", err)
	}

	return output, nil
}

