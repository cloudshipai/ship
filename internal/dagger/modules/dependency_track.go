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

// ScanSBOM scans a Software Bill of Materials (SBOM) file
func (m *DependencyTrackModule) ScanSBOM(ctx context.Context, sbomPath string) (string, error) {
	sbomFile := m.Client.Host().File(sbomPath)
	
	result := m.Client.Container().
		From("dependencytrack/apiserver:latest").
		WithFile("/app/sbom.json", sbomFile).
		WithExec([]string{
			"java", "-jar", "/app/dependency-track-apiserver.jar",
			"--analyze", "/app/sbom.json",
		})

	return result.Stdout(ctx)
}

// AnalyzeProject analyzes a project directory for dependencies
func (m *DependencyTrackModule) AnalyzeProject(ctx context.Context, projectPath string) (string, error) {
	projectDir := m.Client.Host().Directory(projectPath)
	
	result := m.Client.Container().
		From("dependencytrack/apiserver:latest").
		WithDirectory("/app/project", projectDir).
		WithWorkdir("/app/project").
		WithExec([]string{
			"java", "-jar", "/app/dependency-track-apiserver.jar",
			"--analyze", ".",
		})

	return result.Stdout(ctx)
}

// GenerateReport generates a vulnerability report
func (m *DependencyTrackModule) GenerateReport(ctx context.Context, projectPath string, format string) (string, error) {
	if format == "" {
		format = "json"
	}
	
	projectDir := m.Client.Host().Directory(projectPath)
	
	result := m.Client.Container().
		From("dependencytrack/apiserver:latest").
		WithDirectory("/app/project", projectDir).
		WithWorkdir("/app/project").
		WithExec([]string{
			"java", "-jar", "/app/dependency-track-apiserver.jar",
			"--report", "--format", format, ".",
		})

	return result.Stdout(ctx)
}

// ValidateComponents validates components against policies
func (m *DependencyTrackModule) ValidateComponents(ctx context.Context, projectPath string) (string, error) {
	projectDir := m.Client.Host().Directory(projectPath)
	
	result := m.Client.Container().
		From("dependencytrack/apiserver:latest").
		WithDirectory("/app/project", projectDir).
		WithWorkdir("/app/project").
		WithExec([]string{
			"java", "-jar", "/app/dependency-track-apiserver.jar",
			"--validate", ".",
		})

	return result.Stdout(ctx)
}

// TrackDependencies tracks dependencies and their lineage
func (m *DependencyTrackModule) TrackDependencies(ctx context.Context, projectPath string) (string, error) {
	projectDir := m.Client.Host().Directory(projectPath)
	
	result := m.Client.Container().
		From("dependencytrack/apiserver:latest").
		WithDirectory("/app/project", projectDir).
		WithWorkdir("/app/project").
		WithExec([]string{
			"java", "-jar", "/app/dependency-track-apiserver.jar",
			"--track", "--lineage", ".",
		})

	return result.Stdout(ctx)
}
