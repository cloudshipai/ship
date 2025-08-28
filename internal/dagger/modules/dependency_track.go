package modules

import (
	"context"
	"dagger.io/dagger"
)

// DependencyTrackModule provides OWASP Dependency-Track scanning capabilities
type DependencyTrackModule struct {
	Client *dagger.Client
}

// Common binary paths for dependency track tools
const (
	dtrackCliBinary = "/usr/local/bin/dtrack-cli"
	dependencyTrackSyftBinary = "syft"
	cyclonedxNpmBinary = "/usr/local/bin/cyclonedx-npm"
	cyclonedxPyBinary = "/usr/local/bin/cyclonedx-py"
)

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
			dtrackCliBinary, 
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
		WithExec([]string{dependencyTrackSyftBinary, ".", "-o", "cyclonedx-json=sbom.json"}).
		WithExec([]string{
			dtrackCliBinary,
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
		WithExec([]string{dependencyTrackSyftBinary, ".", "-o", "cyclonedx-json=sbom.json"}).
		WithExec([]string{
			dtrackCliBinary,
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
		WithExec([]string{dependencyTrackSyftBinary, ".", "-o", "cyclonedx-json=sbom.json"}).
		WithExec([]string{
			dtrackCliBinary,
			"--bom-path", "sbom.json",
			"--project-name", projectName,
			"--project-version", projectVersion,
			"--auto-create", "true",
		})

	return result.Stdout(ctx)
}

// UploadBOM uploads BOM file to Dependency Track using dtrack-cli
func (m *DependencyTrackModule) UploadBOM(ctx context.Context, bomPath string, projectName string, projectVersion string, serverURL string, apiKey string) (string, error) {
	bomFile := m.Client.Host().File(bomPath)
	
	args := []string{
		"dtrack-cli",
		"--bom-path", "/app/bom.json",
		"--project-name", projectName,
		"--auto-create", "true",
	}
	
	if projectVersion != "" {
		args = append(args, "--project-version", projectVersion)
	}
	if serverURL != "" {
		args = append(args, "--server", serverURL)
	}
	if apiKey != "" {
		args = append(args, "--api-key", apiKey)
	}

	result := m.Client.Container().
		From("node:alpine").
		WithExec([]string{"npm", "install", "-g", "@fjbarrena/dtrack-cli"}).
		WithFile("/app/bom.json", bomFile).
		WithExec(args)

	return result.Stdout(ctx)
}

// UploadBOMAPI uploads BOM to Dependency Track via REST API using curl
func (m *DependencyTrackModule) UploadBOMAPI(ctx context.Context, bomPath string, serverURL string, apiKey string, projectName string, projectVersion string, autoCreate bool) (string, error) {
	bomFile := m.Client.Host().File(bomPath)
	
	args := []string{
		"curl", "-X", "POST", serverURL + "/api/v1/bom",
		"-H", "Content-Type: multipart/form-data",
		"-H", "X-Api-Key: " + apiKey,
		"-F", "bom=@/app/bom.json",
	}
	
	if projectName != "" {
		args = append(args, "-F", "projectName=" + projectName)
	}
	if projectVersion != "" {
		args = append(args, "-F", "projectVersion=" + projectVersion)
	}
	if autoCreate {
		args = append(args, "-F", "autoCreate=true")
	}

	result := m.Client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithFile("/app/bom.json", bomFile).
		WithExec(args)

	return result.Stdout(ctx)
}

// GenerateBOM generates CycloneDX BOM using various build tools
func (m *DependencyTrackModule) GenerateBOM(ctx context.Context, projectType string, projectPath string, outputFile string) (string, error) {
	projectDir := m.Client.Host().Directory(projectPath)
	
	var result *dagger.Container
	
	switch projectType {
	case "npm":
		result = m.Client.Container().
			From("node:alpine").
			WithExec([]string{"npm", "install", "-g", "@cyclonedx/cyclonedx-npm"}).
			WithDirectory("/app/project", projectDir).
			WithWorkdir("/app/project").
			WithExec([]string{cyclonedxNpmBinary, "-o", "bom.json"})
		
	case "maven":
		result = m.Client.Container().
			From("maven:3-openjdk-11").
			WithDirectory("/app/project", projectDir).
			WithWorkdir("/app/project").
			WithExec([]string{"mvn", "org.cyclonedx:cyclonedx-maven-plugin:makeBom"})
		
	case "gradle":
		result = m.Client.Container().
			From("gradle:jdk11").
			WithDirectory("/app/project", projectDir).
			WithWorkdir("/app/project").
			WithExec([]string{"gradle", "cyclonedxBom"})
		
	case "pip":
		result = m.Client.Container().
			From("python:3.9-alpine").
			WithExec([]string{"pip", "install", "cyclonedx-bom"}).
			WithDirectory("/app/project", projectDir).
			WithWorkdir("/app/project").
			WithExec([]string{cyclonedxPyBinary, "-o", "bom.json"})
		
	case "composer":
		result = m.Client.Container().
			From("composer:latest").
			WithExec([]string{"composer", "global", "require", "cyclonedx/cyclonedx-php-composer"}).
			WithDirectory("/app/project", projectDir).
			WithWorkdir("/app/project").
			WithExec([]string{"cyclonedx-php", "composer"})
		
	case "dotnet":
		result = m.Client.Container().
			From("mcr.microsoft.com/dotnet/sdk:6.0").
			WithExec([]string{"dotnet", "tool", "install", "--global", "CycloneDX"}).
			WithDirectory("/app/project", projectDir).
			WithWorkdir("/app/project").
			WithExec([]string{"cyclonedx", "dotnet"})
		
	default:
		// Default to npm
		result = m.Client.Container().
			From("node:alpine").
			WithExec([]string{"npm", "install", "-g", "@cyclonedx/cyclonedx-npm"}).
			WithDirectory("/app/project", projectDir).
			WithWorkdir("/app/project").
			WithExec([]string{cyclonedxNpmBinary, "-o", "bom.json"})
	}

	return result.Stdout(ctx)
}
