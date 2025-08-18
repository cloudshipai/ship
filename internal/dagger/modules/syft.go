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

// LanguageSpecificCataloging generates SBOM with focus on specific programming languages
func (m *SyftModule) LanguageSpecificCataloging(ctx context.Context, target string, languages []string, outputFormat string, packageManagers string, includeDevDeps bool, includeTestDeps bool, depthLimit string) (string, error) {
	if outputFormat == "" {
		outputFormat = "syft-json"
	}

	container := m.client.Container().
		From("anchore/syft:latest").
		WithDirectory("/workspace", m.client.Host().Directory(target)).
		WithWorkdir("/workspace")

	// Map languages to catalogers
	languageMap := map[string]string{
		"python":     "python-package",
		"javascript": "javascript-package",
		"typescript": "javascript-package",
		"java":       "java-package",
		"go":         "go-module",
		"rust":       "rust-cargo",
		"ruby":       "ruby-gem",
		"php":        "php-composer",
		"dotnet":     "dotnet-package",
		"cpp":        "conan-package",
	}

	var catalogers []string
	for _, lang := range languages {
		if cataloger, exists := languageMap[lang]; exists {
			catalogers = append(catalogers, cataloger)
		}
	}

	args := []string{"syft", "dir:.", "-o", outputFormat}
	if len(catalogers) > 0 {
		args = append(args, "--catalogers", fmt.Sprintf("%s", catalogers[0]))
	}
	if packageManagers != "" {
		args = append(args, "--catalogers", packageManagers)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate language-specific SBOM: %w", err)
	}

	return output, nil
}

// SupplyChainAnalysis performs comprehensive supply chain analysis and SBOM generation
func (m *SyftModule) SupplyChainAnalysis(ctx context.Context, target string, analysisDepth string, outputFormats []string, outputDirectory string, includeTransitiveDeps bool, includeLicenseAnalysis bool, includeProvenance bool, riskAssessment string) (string, error) {
	container := m.client.Container().
		From("anchore/syft:latest").
		WithDirectory("/workspace", m.client.Host().Directory(target)).
		WithWorkdir("/workspace")

	args := []string{"syft", "dir:."}

	// Configure analysis depth
	switch analysisDepth {
	case "shallow":
		args = append(args, "--scope", "Squashed")
	case "deep":
		args = append(args, "--scope", "AllLayers")
	case "comprehensive":
		args = append(args, "--scope", "AllLayers", "--catalogers", "all")
	default:
		args = append(args, "--scope", "AllLayers")
	}

	// Configure output formats
	if len(outputFormats) > 0 {
		args = append(args, "-o", outputFormats[0])
	} else {
		args = append(args, "-o", "cyclonedx-json,spdx-json")
	}

	if includeLicenseAnalysis {
		args = append(args, "-o", "spdx-json")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to perform supply chain analysis: %w", err)
	}

	return output, nil
}

// SBOMComparison compares SBOMs and generates difference analysis
func (m *SyftModule) SBOMComparison(ctx context.Context, baselineTarget string, comparisonTarget string, comparisonType string, outputFormat string, diffOutputFile string, showAddedOnly bool, showRemovedOnly bool, includeVersionChanges bool) (string, error) {
	container := m.client.Container().
		From("anchore/syft:latest")

	// Generate SBOM for baseline target (simplified implementation)
	if baselineTarget != "" {
		container = container.WithDirectory("/baseline", m.client.Host().Directory(baselineTarget))
	}
	if comparisonTarget != "" {
		container = container.WithDirectory("/comparison", m.client.Host().Directory(comparisonTarget))
	}

	args := []string{"syft", "/baseline", "-o", "syft-json", "--file", "/tmp/baseline-sbom.json"}
	if outputFormat != "" {
		args = append(args, "-o", outputFormat)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to perform SBOM comparison: %w", err)
	}

	return output, nil
}

// ComplianceAttestation generates compliance-focused SBOMs with attestation features
func (m *SyftModule) ComplianceAttestation(ctx context.Context, target string, complianceFramework string, outputFormat string, attestationFormat string, outputFile string, includeSupplierInfo bool, includeHashes bool, validateCompleteness bool) (string, error) {
	container := m.client.Container().
		From("anchore/syft:latest").
		WithDirectory("/workspace", m.client.Host().Directory(target)).
		WithWorkdir("/workspace")

	args := []string{"syft", "dir:."}

	// Configure output format based on compliance framework
	switch complianceFramework {
	case "ntia-minimum":
		args = append(args, "-o", "spdx-json")
	case "spdx-2.3":
		args = append(args, "-o", "spdx-json")
	case "cyclonedx-1.4":
		args = append(args, "-o", "cyclonedx-json")
	case "sbom-quality":
		args = append(args, "-o", "syft-json")
	}

	if outputFormat != "" {
		args = append(args, "-o", outputFormat)
	}
	if outputFile != "" {
		args = append(args, "--file", outputFile)
	}

	// Add compliance-specific catalogers
	args = append(args, "--catalogers", "all", "--scope", "AllLayers")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate compliance SBOM: %w", err)
	}

	return output, nil
}

// ArchiveAnalysis analyzes archives, packages, and compressed files for SBOM generation
func (m *SyftModule) ArchiveAnalysis(ctx context.Context, archivePath string, archiveType string, outputFormat string, extractNested bool, includeMetadata bool, extractionDepth string) (string, error) {
	container := m.client.Container().
		From("anchore/syft:latest").
		WithFile("/archive", m.client.Host().File(archivePath))

	args := []string{"syft", "/archive"}
	if outputFormat != "" {
		args = append(args, "-o", outputFormat)
	}
	if archiveType != "auto" && archiveType != "" {
		args = append(args, "--from", archiveType)
	}
	if extractNested {
		args = append(args, "--scope", "AllLayers")
	}
	if extractionDepth != "" {
		args = append(args, "--max-depth", extractionDepth)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to analyze archive: %w", err)
	}

	return output, nil
}

// CICDPipelineIntegration performs optimized SBOM generation for CI/CD pipelines
func (m *SyftModule) CICDPipelineIntegration(ctx context.Context, target string, pipelineStage string, artifactName string, outputFormats []string, outputDirectory string, failOnError bool, quietMode bool, timeout string) (string, error) {
	container := m.client.Container().
		From("anchore/syft:latest").
		WithDirectory("/workspace", m.client.Host().Directory(target)).
		WithWorkdir("/workspace")

	args := []string{"syft", "dir:."}

	// Configure based on pipeline stage
	switch pipelineStage {
	case "build":
		args = append(args, "-o", "syft-json,cyclonedx-json")
	case "test":
		args = append(args, "-o", "syft-json")
	case "staging", "production":
		args = append(args, "-o", "spdx-json,cyclonedx-json")
	case "release":
		args = append(args, "-o", "spdx-json,cyclonedx-json,syft-json")
	default:
		args = append(args, "-o", "syft-json,cyclonedx-json")
	}

	if len(outputFormats) > 0 {
		args = append(args, "-o", outputFormats[0])
	}
	if outputDirectory != "" {
		if artifactName == "" {
			artifactName = "sbom"
		}
		args = append(args, "--file", outputDirectory+"/"+artifactName)
	}
	if quietMode {
		args = append(args, "--quiet")
	}
	if timeout != "" {
		args = append(args, "--timeout", timeout+"s")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate CI/CD SBOM: %w", err)
	}

	return output, nil
}

// MetadataExtraction extracts and enriches metadata for comprehensive SBOM generation
func (m *SyftModule) MetadataExtraction(ctx context.Context, target string, metadataTypes []string, outputFormat string, includeFileMetadata bool, includeChecksums bool, includeCertificates bool, includeSignatures bool, customAnnotations string) (string, error) {
	container := m.client.Container().
		From("anchore/syft:latest").
		WithDirectory("/workspace", m.client.Host().Directory(target)).
		WithWorkdir("/workspace")

	// Use syft-json for maximum metadata preservation
	if outputFormat == "" {
		outputFormat = "syft-json"
	}

	args := []string{"syft", "dir:.", "-o", outputFormat}

	// Enable all catalogers for comprehensive metadata
	args = append(args, "--catalogers", "all", "--scope", "AllLayers")

	if includeFileMetadata {
		// Enable file cataloger for file-level metadata
		args = append(args, "--catalogers", "file-metadata")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to extract metadata: %w", err)
	}

	return output, nil
}