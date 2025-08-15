package modules

import (
	"context"
	"dagger.io/dagger"
)

// GuacModule provides GUAC (Graph for Understanding Artifact Composition) capabilities
type GuacModule struct {
	Client *dagger.Client
}

// NewGuacModule creates a new GUAC module
func NewGuacModule(client *dagger.Client) *GuacModule {
	return &GuacModule{
		Client: client,
	}
}

// IngestSBOM ingests an SBOM into the GUAC graph
func (m *GuacModule) IngestSBOM(ctx context.Context, sbomPath string) (string, error) {
	sbomFile := m.Client.Host().File(sbomPath)
	
	result := m.Client.Container().
		From("ghcr.io/guacsec/guac:latest").
		WithFile("/app/sbom.json", sbomFile).
		WithExec([]string{
			"guacone", "collect", "files", "/app/sbom.json",
		})

	return result.Stdout(ctx)
}

// AnalyzeArtifact analyzes an artifact and its dependencies
func (m *GuacModule) AnalyzeArtifact(ctx context.Context, artifactPath string) (string, error) {
	artifactFile := m.Client.Host().File(artifactPath)
	
	result := m.Client.Container().
		From("ghcr.io/guacsec/guac:latest").
		WithFile("/app/artifact", artifactFile).
		WithExec([]string{
			"guacone", "collect", "image", "/app/artifact",
		})

	return result.Stdout(ctx)
}

// QueryDependencies queries the GUAC graph for dependency information
func (m *GuacModule) QueryDependencies(ctx context.Context, packageName string) (string, error) {
	result := m.Client.Container().
		From("ghcr.io/guacsec/guac:latest").
		WithExec([]string{
			"guacone", "query", "dependencies", packageName,
		})

	return result.Stdout(ctx)
}

// QueryVulnerabilities queries vulnerabilities for a package
func (m *GuacModule) QueryVulnerabilities(ctx context.Context, packageName string) (string, error) {
	result := m.Client.Container().
		From("ghcr.io/guacsec/guac:latest").
		WithExec([]string{
			"guacone", "query", "vulnerabilities", packageName,
		})

	return result.Stdout(ctx)
}

// GenerateGraph generates a dependency graph visualization
func (m *GuacModule) GenerateGraph(ctx context.Context, projectPath string) (string, error) {
	projectDir := m.Client.Host().Directory(projectPath)
	
	result := m.Client.Container().
		From("ghcr.io/guacsec/guac:latest").
		WithDirectory("/app/project", projectDir).
		WithWorkdir("/app/project").
		WithExec([]string{
			"guacone", "collect", "files", ".",
		}).
		WithExec([]string{
			"guacone", "visualizer", "--output", "graph.dot",
		})

	return result.Stdout(ctx)
}

// AnalyzeImpact analyzes the impact of a vulnerability across the dependency graph
func (m *GuacModule) AnalyzeImpact(ctx context.Context, vulnID string) (string, error) {
	result := m.Client.Container().
		From("ghcr.io/guacsec/guac:latest").
		WithExec([]string{
			"guacone", "query", "impact", vulnID,
		})

	return result.Stdout(ctx)
}

// CollectFiles collects and processes multiple files into the GUAC graph
func (m *GuacModule) CollectFiles(ctx context.Context, projectPath string) (string, error) {
	projectDir := m.Client.Host().Directory(projectPath)
	
	result := m.Client.Container().
		From("ghcr.io/guacsec/guac:latest").
		WithDirectory("/app/project", projectDir).
		WithWorkdir("/app/project").
		WithExec([]string{
			"guacone", "collect", "files", ".",
		})

	return result.Stdout(ctx)
}

// ValidateAttestation validates software attestations
func (m *GuacModule) ValidateAttestation(ctx context.Context, attestationPath string) (string, error) {
	attestationFile := m.Client.Host().File(attestationPath)
	
	result := m.Client.Container().
		From("ghcr.io/guacsec/guac:latest").
		WithFile("/app/attestation.json", attestationFile).
		WithExec([]string{
			"guacone", "collect", "attestation", "/app/attestation.json",
		})

	return result.Stdout(ctx)
}
