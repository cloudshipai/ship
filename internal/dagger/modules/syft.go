package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// SyftModule runs Syft for SBOM generation
type SyftModule struct {
	client *dagger.Client
	name   string
}

// NewSyftModule creates a new Syft module
func NewSyftModule(client *dagger.Client) *SyftModule {
	return &SyftModule{
		client: client,
		name:   "syft",
	}
}

// GenerateSBOMFromDirectory generates SBOM from a directory
func (m *SyftModule) GenerateSBOMFromDirectory(ctx context.Context, dir string, format string) (string, error) {
	if format == "" {
		format = "json"
	}

	container := m.client.Container().
		From("anchore/syft:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"syft", "dir:.", "-o", format})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate SBOM: %w", err)
	}

	return output, nil
}

// GenerateSBOMFromImage generates SBOM from a container image
func (m *SyftModule) GenerateSBOMFromImage(ctx context.Context, imageName string, format string) (string, error) {
	if format == "" {
		format = "json"
	}

	container := m.client.Container().
		From("anchore/syft:latest").
		WithExec([]string{"syft", imageName, "-o", format})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate SBOM from image: %w", err)
	}

	return output, nil
}

// GenerateSBOMFromPackage generates SBOM from a specific package manager
func (m *SyftModule) GenerateSBOMFromPackage(ctx context.Context, dir string, packageType string, format string) (string, error) {
	if format == "" {
		format = "json"
	}

	var source string
	switch packageType {
	case "npm", "yarn":
		source = "dir:."
	case "pip", "python":
		source = "dir:."
	case "go":
		source = "dir:."
	case "maven", "gradle":
		source = "dir:."
	default:
		source = "dir:."
	}

	container := m.client.Container().
		From("anchore/syft:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"syft", source, "-o", format})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate SBOM for %s: %w", packageType, err)
	}

	return output, nil
}

// GenerateAttestations generates SBOM with attestations
func (m *SyftModule) GenerateAttestations(ctx context.Context, target string, format string) (string, error) {
	if format == "" {
		format = "json"
	}

	container := m.client.Container().From("anchore/syft:latest")
	
	if target[:6] != "image:" {
		// Assume it's a directory
		container = container.
			WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace").
			WithExec([]string{"syft", "dir:.", "-o", format, "--source-name", target})
	} else {
		// It's an image
		container = container.WithExec([]string{"syft", target, "-o", format})
	}

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate SBOM with attestations: %w", err)
	}

	return output, nil
}