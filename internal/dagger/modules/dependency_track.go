package modules

import (
	"context"
	"dagger.io/dagger"
)

// DependencyTrackModule provides OWASP Dependency-Track scanning capabilities
type DependencyTrackModule struct {
	Client *dagger.Client
}

// NewDependencyTrackModule creates a new DependencyTrack module
func NewDependencyTrackModule(client *dagger.Client) *DependencyTrackModule {
	return &DependencyTrackModule{
		Client: client,
	}
}

// ScanSBOM scans a Software Bill of Materials (SBOM) file using dtrack-cli
func (m *DependencyTrackModule) ScanSBOM(ctx context.Context, sbomPath string) (string, error) {
	sbomFile := m.Client.Host().File(sbomPath)
	
	result := m.Client.Container().
		From("node:alpine").
		WithExec([]string{"npm", "install", "-g", "@fjbarrena/dtrack-cli"}).
		WithFile("/app/sbom.json", sbomFile).
		WithExec([]string{
			"dtrack-cli", 
			"--bom-path", "/app/sbom.json",
			"--project-name", "default-project",
			"--project-version", "latest",
		})

	return result.Stdout(ctx)
}

// AnalyzeProject generates and uploads SBOM for a project directory
func (m *DependencyTrackModule) AnalyzeProject(ctx context.Context, projectPath string, projectName string, projectVersion string) (string, error) {
	projectDir := m.Client.Host().Directory(projectPath)
	
	// First generate SBOM using syft, then upload to Dependency Track
	result := m.Client.Container().
		From("node:alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"npm", "install", "-g", "@fjbarrena/dtrack-cli"}).
		WithExec([]string{"sh", "-c", "curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin"}).
		WithDirectory("/app/project", projectDir).
		WithWorkdir("/app/project").
		WithExec([]string{"syft", ".", "-o", "cyclonedx-json=sbom.json"}).
		WithExec([]string{
			"dtrack-cli",
			"--bom-path", "sbom.json",
			"--project-name", projectName,
			"--project-version", projectVersion,
			"--auto-create", "true",
		})

	return result.Stdout(ctx)
}

// GenerateReport generates a vulnerability report using dtrack-cli (requires existing project)
func (m *DependencyTrackModule) GenerateReport(ctx context.Context, projectName string, projectVersion string) (string, error) {
	result := m.Client.Container().
		From("node:alpine").
		WithExec([]string{"npm", "install", "-g", "@fjbarrena/dtrack-cli"}).
		WithExec([]string{
			"sh", "-c", 
			"echo 'Note: dtrack-cli primarily uploads BOMs. For reports, use Dependency Track web UI or API directly.'",
		})

	return result.Stdout(ctx)
}

// ValidateComponents uploads SBOM for validation (policy evaluation happens server-side)
func (m *DependencyTrackModule) ValidateComponents(ctx context.Context, projectPath string, projectName string, projectVersion string) (string, error) {
	projectDir := m.Client.Host().Directory(projectPath)
	
	result := m.Client.Container().
		From("node:alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"npm", "install", "-g", "@fjbarrena/dtrack-cli"}).
		WithExec([]string{"sh", "-c", "curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin"}).
		WithDirectory("/app/project", projectDir).
		WithWorkdir("/app/project").
		WithExec([]string{"syft", ".", "-o", "cyclonedx-json=sbom.json"}).
		WithExec([]string{
			"dtrack-cli",
			"--bom-path", "sbom.json",
			"--project-name", projectName,
			"--project-version", projectVersion,
			"--auto-create", "true",
		})

	return result.Stdout(ctx)
}

// TrackDependencies uploads SBOM to track dependencies (tracking happens server-side)
func (m *DependencyTrackModule) TrackDependencies(ctx context.Context, projectPath string, projectName string, projectVersion string) (string, error) {
	projectDir := m.Client.Host().Directory(projectPath)
	
	result := m.Client.Container().
		From("node:alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"npm", "install", "-g", "@fjbarrena/dtrack-cli"}).
		WithExec([]string{"sh", "-c", "curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin"}).
		WithDirectory("/app/project", projectDir).
		WithWorkdir("/app/project").
		WithExec([]string{"syft", ".", "-o", "cyclonedx-json=sbom.json"}).
		WithExec([]string{
			"dtrack-cli",
			"--bom-path", "sbom.json",
			"--project-name", projectName,
			"--project-version", projectVersion,
			"--auto-create", "true",
		})

	return result.Stdout(ctx)
}
