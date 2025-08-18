package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// SLSAVerifierModule runs SLSA Verifier for provenance verification
type SLSAVerifierModule struct {
	client *dagger.Client
	name   string
}

const slsaVerifierBinary = "/usr/local/bin/slsa-verifier"

// NewSLSAVerifierModule creates a new SLSA Verifier module
func NewSLSAVerifierModule(client *dagger.Client) *SLSAVerifierModule {
	return &SLSAVerifierModule{
		client: client,
		name:   slsaVerifierBinary,
	}
}

// VerifyArtifact verifies SLSA provenance for binary artifacts
func (m *SLSAVerifierModule) VerifyArtifact(ctx context.Context, artifactPath string, provenancePath string, sourceURI string, sourceTag string, sourceBranch string, builderID string, printProvenance bool) (string, error) {
	container := m.client.Container().
		From("ghcr.io/slsa-framework/slsa-verifier:latest").
		WithFile("/workspace/artifact", m.client.Host().File(artifactPath)).
		WithFile("/workspace/provenance", m.client.Host().File(provenancePath)).
		WithWorkdir("/workspace")

	args := []string{slsaVerifierBinary, "verify-artifact", "/workspace/artifact", "--provenance-path", "/workspace/provenance", "--source-uri", sourceURI}

	if sourceTag != "" {
		args = append(args, "--source-tag", sourceTag)
	}
	if sourceBranch != "" {
		args = append(args, "--source-branch", sourceBranch)
	}
	if builderID != "" {
		args = append(args, "--builder-id", builderID)
	}
	if printProvenance {
		args = append(args, "--print-provenance")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run SLSA verifier artifact verification: %w", err)
	}

	return output, nil
}

// VerifyImage verifies SLSA provenance for container images
func (m *SLSAVerifierModule) VerifyImage(ctx context.Context, image string, sourceURI string, sourceTag string, sourceBranch string, builderID string, printProvenance bool) (string, error) {
	container := m.client.Container().
		From("ghcr.io/slsa-framework/slsa-verifier:latest")

	args := []string{slsaVerifierBinary, "verify-image", image, "--source-uri", sourceURI}

	if sourceTag != "" {
		args = append(args, "--source-tag", sourceTag)
	}
	if sourceBranch != "" {
		args = append(args, "--source-branch", sourceBranch)
	}
	if builderID != "" {
		args = append(args, "--builder-id", builderID)
	}
	if printProvenance {
		args = append(args, "--print-provenance")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run SLSA verifier image verification: %w", err)
	}

	return output, nil
}

// VerifyNpmPackage verifies SLSA provenance for npm packages (experimental)
func (m *SLSAVerifierModule) VerifyNpmPackage(ctx context.Context, packageTarball string, attestationsPath string, packageName string, packageVersion string, sourceURI string, printProvenance bool) (string, error) {
	container := m.client.Container().
		From("ghcr.io/slsa-framework/slsa-verifier:latest").
		WithFile("/workspace/package.tgz", m.client.Host().File(packageTarball)).
		WithFile("/workspace/attestations", m.client.Host().File(attestationsPath)).
		WithWorkdir("/workspace")

	args := []string{slsaVerifierBinary, "verify-npm-package", "/workspace/package.tgz", "--attestations-path", "/workspace/attestations", "--package-name", packageName, "--package-version", packageVersion}

	if sourceURI != "" {
		args = append(args, "--source-uri", sourceURI)
	}
	if printProvenance {
		args = append(args, "--print-provenance")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run SLSA verifier npm package verification: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of SLSA Verifier
func (m *SLSAVerifierModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/slsa-framework/slsa-verifier:latest").
		WithExec([]string{slsaVerifierBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get SLSA verifier version: %w", err)
	}

	return output, nil
}