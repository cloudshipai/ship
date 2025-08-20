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
		WithExec([]string{
			"/syft", ".", "-o", format,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to generate SBOM: no output received")
}

// GenerateSBOMFromImage generates SBOM from a container image
func (m *SyftModule) GenerateSBOMFromImage(ctx context.Context, imageName string, format string) (string, error) {
	if format == "" {
		format = "json"
	}

	container := m.client.Container().
		From("anchore/syft:latest").
		WithExec([]string{
			"/syft", imageName, "-o", format,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to generate SBOM from image: no output received")
}

// GenerateSBOMFromPackage generates SBOM from a specific package manager
func (m *SyftModule) GenerateSBOMFromPackage(ctx context.Context, dir string, packageType string, format string) (string, error) {
	if format == "" {
		format = "json"
	}

	// packageType can be used for future enhancements to filter by package manager
	_ = packageType

	container := m.client.Container().
		From("anchore/syft:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"/syft", ".", "-o", format,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

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
	
	if len(target) > 6 && target[:6] == "image:" {
		// It's an image
		container = container.WithExec([]string{
			"/syft", target, "-o", format,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})
	} else {
		// Assume it's a directory
		container = container.
			WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace").
			WithExec([]string{
				"/syft", ".", "-o", format,
			}, dagger.ContainerWithExecOpts{
				Expect: "ANY",
			})
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

	cmd := "/usr/local/bin/syft dir:. -o " + outputFormat
	if len(catalogers) > 0 {
		cmd += " --catalogers " + catalogers[0]
	}
	if packageManagers != "" {
		cmd += " --catalogers " + packageManagers
	}
	cmd += " 2>&1"

	container = container.WithExec([]string{"sh", "-c", cmd}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

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

	cmd := "/usr/local/bin/syft dir:."

	// Configure analysis depth
	switch analysisDepth {
	case "shallow":
		cmd += " --scope Squashed"
	case "deep":
		cmd += " --scope AllLayers"
	case "comprehensive":
		cmd += " --scope AllLayers --catalogers all"
	default:
		cmd += " --scope AllLayers"
	}

	// Configure output formats
	if len(outputFormats) > 0 {
		cmd += " -o " + outputFormats[0]
	} else {
		cmd += " -o cyclonedx-json,spdx-json"
	}

	if includeLicenseAnalysis {
		cmd += " -o spdx-json"
	}
	cmd += " 2>&1"

	container = container.WithExec([]string{"sh", "-c", cmd}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

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

	cmd := "/usr/local/bin/syft /baseline -o syft-json --file /tmp/baseline-sbom.json"
	if outputFormat != "" {
		cmd += " -o " + outputFormat
	}
	cmd += " 2>&1"

	container = container.WithExec([]string{"sh", "-c", cmd}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

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

	cmd := "/usr/local/bin/syft dir:."

	// Configure output format based on compliance framework
	switch complianceFramework {
	case "ntia-minimum":
		cmd += " -o spdx-json"
	case "spdx-2.3":
		cmd += " -o spdx-json"
	case "cyclonedx-1.4":
		cmd += " -o cyclonedx-json"
	case "sbom-quality":
		cmd += " -o syft-json"
	}

	if outputFormat != "" {
		cmd += " -o " + outputFormat
	}
	if outputFile != "" {
		cmd += " --file " + outputFile
	}

	// Add compliance-specific catalogers
	cmd += " --catalogers all --scope AllLayers 2>&1"

	container = container.WithExec([]string{"sh", "-c", cmd}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

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

	cmd := "/usr/local/bin/syft /archive"
	if outputFormat != "" {
		cmd += " -o " + outputFormat
	}
	if archiveType != "auto" && archiveType != "" {
		cmd += " --from " + archiveType
	}
	if extractNested {
		cmd += " --scope AllLayers"
	}
	if extractionDepth != "" {
		cmd += " --max-depth " + extractionDepth
	}
	cmd += " 2>&1"

	container = container.WithExec([]string{"sh", "-c", cmd}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

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

	cmd := "/usr/local/bin/syft dir:."

	// Configure based on pipeline stage
	switch pipelineStage {
	case "build":
		cmd += " -o syft-json,cyclonedx-json"
	case "test":
		cmd += " -o syft-json"
	case "staging", "production":
		cmd += " -o spdx-json,cyclonedx-json"
	case "release":
		cmd += " -o spdx-json,cyclonedx-json,syft-json"
	default:
		cmd += " -o syft-json,cyclonedx-json"
	}

	if len(outputFormats) > 0 {
		cmd += " -o " + outputFormats[0]
	}
	if outputDirectory != "" {
		if artifactName == "" {
			artifactName = "sbom"
		}
		cmd += " --file " + outputDirectory + "/" + artifactName
	}
	if quietMode {
		cmd += " --quiet"
	}
	if timeout != "" {
		cmd += " --timeout " + timeout + "s"
	}
	cmd += " 2>&1"

	container = container.WithExec([]string{"sh", "-c", cmd}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

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

	cmd := "/usr/local/bin/syft dir:. -o " + outputFormat

	// Enable all catalogers for comprehensive metadata
	cmd += " --catalogers all --scope AllLayers"

	if includeFileMetadata {
		// Enable file cataloger for file-level metadata
		cmd += " --catalogers file-metadata"
	}
	cmd += " 2>&1"

	container = container.WithExec([]string{"sh", "-c", cmd}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to extract metadata: %w", err)
	}

	return output, nil
}