package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddSyftTools adds Syft (SBOM generation from container images and filesystems) MCP tool implementations using direct Dagger calls
func AddSyftTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Add the new unified syft_sbom tool first
	addNewSyftSBOMTool(s)
	
	// Keep existing tools for backward compatibility
	addSyftToolsDirect(s)
}

// addNewSyftSBOMTool adds the new unified SBOM generation tool
func addNewSyftSBOMTool(s *server.MCPServer) {
	// Syft SBOM generation tool - unified interface
	sbomTool := mcp.NewTool("syft_sbom",
		mcp.WithDescription("Generate CycloneDX or SPDX SBOM from a directory, image, or archive"),
		mcp.WithString("target",
			mcp.Description("Target to scan (e.g., dir:., docker:alpine:3.19, oci-archive:/path/image.tar)"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format: cyclonedx-json or spdx-json (default: cyclonedx-json)"),
		),
		mcp.WithString("output_path",
			mcp.Description("Where to write SBOM (default: ./sbom.cdx.json or ./sbom.spdx.json based on format)"),
		),
	)
	s.AddTool(sbomTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		target := request.GetString("target", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		
		format := request.GetString("format", "cyclonedx-json")
		outputPath := request.GetString("output_path", "")
		
		// Set default output path based on format
		if outputPath == "" {
			if format == "spdx-json" {
				outputPath = "./sbom.spdx.json"
			} else {
				outputPath = "./sbom.cdx.json"
			}
		}

		// Generate SBOM
		var stdout string
		var stderr string
		
		// Determine target type and call appropriate method
		if strings.HasPrefix(target, "dir:") {
			dirPath := strings.TrimPrefix(target, "dir:")
			stdout, err = module.GenerateSBOMFromDirectory(ctx, dirPath, format)
		} else if strings.HasPrefix(target, "docker:") || strings.HasPrefix(target, "registry:") {
			stdout, err = module.GenerateSBOMFromImage(ctx, target, format)
		} else if strings.HasPrefix(target, "oci-archive:") {
			// For archives, we'll use the archive analysis method if it exists
			archivePath := strings.TrimPrefix(target, "oci-archive:")
			// Try to use archive analysis, fall back to treating as directory
			stdout, err = module.ArchiveAnalysis(ctx, archivePath, "oci", format, false, true, "")
			if err != nil {
				// Fallback: treat as directory scan
				stdout, err = module.GenerateSBOMFromDirectory(ctx, archivePath, format)
			}
		} else {
			// Default: treat as directory
			stdout, err = module.GenerateSBOMFromDirectory(ctx, target, format)
		}
		
		// Build result in the expected format
		result := map[string]interface{}{
			"status": "ok",
			"stdout": stdout,
			"stderr": stderr,
			"artifacts": map[string]string{},
			"summary": map[string]interface{}{},
			"diagnostics": []string{},
		}
		
		// Add artifact path
		if format == "spdx-json" {
			result["artifacts"].(map[string]string)["sbom_spdx"] = outputPath
		} else {
			result["artifacts"].(map[string]string)["sbom_cyclonedx"] = outputPath
		}
		
		if err != nil {
			result["status"] = "error"
			result["stderr"] = err.Error()
			result["diagnostics"] = []string{fmt.Sprintf("Syft SBOM generation failed: %v", err)}
		}

		// Return as JSON
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
}

// addSyftToolsDirect adds Syft tools using direct Dagger module calls
func addSyftToolsDirect(s *server.MCPServer) {
	// Syft generate SBOM from directory tool
	generateSBOMDirectoryTool := mcp.NewTool("syft_generate_sbom_directory",
		mcp.WithDescription("Generate SBOM from directory using real Syft CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory path to scan"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format (json, spdx-json, cyclonedx-json, table)"),
		),
	)
	s.AddTool(generateSBOMDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		format := request.GetString("format", "json")

		// Generate SBOM from directory
		stdout, err := module.GenerateSBOMFromDirectory(ctx, directory, format)
		
		// Build result in the expected format
		result := map[string]interface{}{
			"status": "ok",
			"stdout": stdout,
			"stderr": "",
			"artifacts": map[string]string{},
			"summary": map[string]interface{}{},
			"diagnostics": []string{},
		}
		
		// Add artifact path based on format
		if format == "spdx-json" {
			result["artifacts"].(map[string]string)["sbom_spdx"] = "./sbom.spdx.json"
		} else {
			result["artifacts"].(map[string]string)["sbom_cyclonedx"] = "./sbom.cdx.json"
		}
		
		if err != nil {
			result["status"] = "error"
			result["stderr"] = err.Error()
			result["diagnostics"] = []string{fmt.Sprintf("Syft generate SBOM from directory failed: %v", err)}
		}

		// Return as JSON
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})

	// Syft generate SBOM from image tool
	generateSBOMImageTool := mcp.NewTool("syft_generate_sbom_image",
		mcp.WithDescription("Generate SBOM from container image using real Syft CLI"),
		mcp.WithString("image",
			mcp.Description("Container image name to scan"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format (json, spdx-json, cyclonedx-json, table)"),
		),
	)
	s.AddTool(generateSBOMImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		image := request.GetString("image", "")
		format := request.GetString("format", "json")

		// Generate SBOM from image
		output, err := module.GenerateSBOMFromImage(ctx, image, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft generate SBOM from image failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Syft generate SBOM from package tool
	generateSBOMPackageTool := mcp.NewTool("syft_generate_sbom_package",
		mcp.WithDescription("Generate SBOM from package using real Syft CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory containing packages"),
			mcp.Required(),
		),
		mcp.WithString("package_type",
			mcp.Description("Package type (npm, pip, gem, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format (json, spdx-json, cyclonedx-json, table)"),
		),
	)
	s.AddTool(generateSBOMPackageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		packageType := request.GetString("package_type", "")
		format := request.GetString("format", "json")

		// Generate SBOM from package
		output, err := module.GenerateSBOMFromPackage(ctx, directory, packageType, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft generate SBOM from package failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Syft generate attestations tool
	generateAttestationsTool := mcp.NewTool("syft_generate_attestations",
		mcp.WithDescription("Generate attestations using real Syft CLI"),
		mcp.WithString("target",
			mcp.Description("Target to generate attestations for"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Attestation format"),
		),
	)
	s.AddTool(generateAttestationsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		target := request.GetString("target", "")
		format := request.GetString("format", "json")

		// Generate attestations
		output, err := module.GenerateAttestations(ctx, target, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft generate attestations failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Syft language specific cataloging tool
	languageSpecificCatalogingTool := mcp.NewTool("syft_language_specific_cataloging",
		mcp.WithDescription("Language-specific cataloging using real Syft CLI"),
		mcp.WithString("target",
			mcp.Description("Target to scan"),
			mcp.Required(),
		),
		mcp.WithString("languages",
			mcp.Description("Comma-separated list of languages"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
		),
		mcp.WithString("package_managers",
			mcp.Description("Package managers to use"),
		),
		mcp.WithBoolean("include_dev_deps",
			mcp.Description("Include development dependencies"),
		),
		mcp.WithBoolean("include_test_deps",
			mcp.Description("Include test dependencies"),
		),
		mcp.WithString("depth_limit",
			mcp.Description("Depth limit for scanning"),
		),
	)
	s.AddTool(languageSpecificCatalogingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		target := request.GetString("target", "")
		languagesStr := request.GetString("languages", "")
		var languages []string
		if languagesStr != "" {
			languages = strings.Split(languagesStr, ",")
		}
		outputFormat := request.GetString("output_format", "json")
		packageManagers := request.GetString("package_managers", "")
		includeDevDeps := request.GetBool("include_dev_deps", false)
		includeTestDeps := request.GetBool("include_test_deps", false)
		depthLimit := request.GetString("depth_limit", "")

		// Language specific cataloging
		output, err := module.LanguageSpecificCataloging(ctx, target, languages, outputFormat, packageManagers, includeDevDeps, includeTestDeps, depthLimit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft language specific cataloging failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Syft supply chain analysis tool
	supplyChainAnalysisTool := mcp.NewTool("syft_supply_chain_analysis",
		mcp.WithDescription("Supply chain analysis using real Syft CLI"),
		mcp.WithString("target",
			mcp.Description("Target to analyze"),
			mcp.Required(),
		),
		mcp.WithString("analysis_depth",
			mcp.Description("Analysis depth"),
		),
		mcp.WithString("output_formats",
			mcp.Description("Comma-separated output formats"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Output directory"),
		),
		mcp.WithBoolean("include_transitive_deps",
			mcp.Description("Include transitive dependencies"),
		),
		mcp.WithBoolean("include_license_analysis",
			mcp.Description("Include license analysis"),
		),
		mcp.WithBoolean("include_provenance",
			mcp.Description("Include provenance information"),
		),
		mcp.WithString("risk_assessment",
			mcp.Description("Risk assessment level"),
		),
	)
	s.AddTool(supplyChainAnalysisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		target := request.GetString("target", "")
		analysisDepth := request.GetString("analysis_depth", "")
		outputFormatsStr := request.GetString("output_formats", "")
		var outputFormats []string
		if outputFormatsStr != "" {
			outputFormats = strings.Split(outputFormatsStr, ",")
		}
		outputDirectory := request.GetString("output_directory", "")
		includeTransitiveDeps := request.GetBool("include_transitive_deps", false)
		includeLicenseAnalysis := request.GetBool("include_license_analysis", false)
		includeProvenance := request.GetBool("include_provenance", false)
		riskAssessment := request.GetString("risk_assessment", "")

		// Supply chain analysis
		output, err := module.SupplyChainAnalysis(ctx, target, analysisDepth, outputFormats, outputDirectory, includeTransitiveDeps, includeLicenseAnalysis, includeProvenance, riskAssessment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft supply chain analysis failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Syft SBOM comparison tool
	sbomComparisonTool := mcp.NewTool("syft_sbom_comparison",
		mcp.WithDescription("SBOM comparison using real Syft CLI"),
		mcp.WithString("baseline_target",
			mcp.Description("Baseline target for comparison"),
			mcp.Required(),
		),
		mcp.WithString("comparison_target",
			mcp.Description("Comparison target"),
			mcp.Required(),
		),
		mcp.WithString("comparison_type",
			mcp.Description("Type of comparison"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
		),
		mcp.WithString("diff_output_file",
			mcp.Description("Diff output file"),
		),
		mcp.WithBoolean("show_added_only",
			mcp.Description("Show only added items"),
		),
		mcp.WithBoolean("show_removed_only",
			mcp.Description("Show only removed items"),
		),
		mcp.WithBoolean("include_version_changes",
			mcp.Description("Include version changes"),
		),
	)
	s.AddTool(sbomComparisonTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		baselineTarget := request.GetString("baseline_target", "")
		comparisonTarget := request.GetString("comparison_target", "")
		comparisonType := request.GetString("comparison_type", "")
		outputFormat := request.GetString("output_format", "json")
		diffOutputFile := request.GetString("diff_output_file", "")
		showAddedOnly := request.GetBool("show_added_only", false)
		showRemovedOnly := request.GetBool("show_removed_only", false)
		includeVersionChanges := request.GetBool("include_version_changes", false)

		// SBOM comparison
		output, err := module.SBOMComparison(ctx, baselineTarget, comparisonTarget, comparisonType, outputFormat, diffOutputFile, showAddedOnly, showRemovedOnly, includeVersionChanges)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft SBOM comparison failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Syft compliance attestation tool
	complianceAttestationTool := mcp.NewTool("syft_compliance_attestation",
		mcp.WithDescription("Compliance attestation using real Syft CLI"),
		mcp.WithString("target",
			mcp.Description("Target for compliance attestation"),
			mcp.Required(),
		),
		mcp.WithString("compliance_framework",
			mcp.Description("Compliance framework"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
		),
		mcp.WithString("attestation_format",
			mcp.Description("Attestation format"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file"),
		),
		mcp.WithBoolean("include_supplier_info",
			mcp.Description("Include supplier information"),
		),
		mcp.WithBoolean("include_hashes",
			mcp.Description("Include hashes"),
		),
		mcp.WithBoolean("validate_completeness",
			mcp.Description("Validate completeness"),
		),
	)
	s.AddTool(complianceAttestationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		target := request.GetString("target", "")
		complianceFramework := request.GetString("compliance_framework", "")
		outputFormat := request.GetString("output_format", "json")
		attestationFormat := request.GetString("attestation_format", "")
		outputFile := request.GetString("output_file", "")
		includeSupplierInfo := request.GetBool("include_supplier_info", false)
		includeHashes := request.GetBool("include_hashes", false)
		validateCompleteness := request.GetBool("validate_completeness", false)

		// Compliance attestation
		output, err := module.ComplianceAttestation(ctx, target, complianceFramework, outputFormat, attestationFormat, outputFile, includeSupplierInfo, includeHashes, validateCompleteness)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft compliance attestation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Syft archive analysis tool
	archiveAnalysisTool := mcp.NewTool("syft_archive_analysis",
		mcp.WithDescription("Archive analysis using real Syft CLI"),
		mcp.WithString("archive_path",
			mcp.Description("Archive path to analyze"),
			mcp.Required(),
		),
		mcp.WithString("archive_type",
			mcp.Description("Archive type"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
		),
		mcp.WithBoolean("extract_nested",
			mcp.Description("Extract nested archives"),
		),
		mcp.WithBoolean("include_metadata",
			mcp.Description("Include metadata"),
		),
		mcp.WithString("extraction_depth",
			mcp.Description("Extraction depth"),
		),
	)
	s.AddTool(archiveAnalysisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		archivePath := request.GetString("archive_path", "")
		archiveType := request.GetString("archive_type", "")
		outputFormat := request.GetString("output_format", "json")
		extractNested := request.GetBool("extract_nested", false)
		includeMetadata := request.GetBool("include_metadata", false)
		extractionDepth := request.GetString("extraction_depth", "")

		// Archive analysis
		output, err := module.ArchiveAnalysis(ctx, archivePath, archiveType, outputFormat, extractNested, includeMetadata, extractionDepth)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft archive analysis failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Syft CI/CD pipeline integration tool
	cicdPipelineIntegrationTool := mcp.NewTool("syft_cicd_pipeline_integration",
		mcp.WithDescription("CI/CD pipeline integration using real Syft CLI"),
		mcp.WithString("target",
			mcp.Description("Target for CI/CD integration"),
			mcp.Required(),
		),
		mcp.WithString("pipeline_stage",
			mcp.Description("Pipeline stage"),
		),
		mcp.WithString("artifact_name",
			mcp.Description("Artifact name"),
		),
		mcp.WithString("output_formats",
			mcp.Description("Comma-separated output formats"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Output directory"),
		),
		mcp.WithBoolean("fail_on_error",
			mcp.Description("Fail on error"),
		),
		mcp.WithBoolean("quiet_mode",
			mcp.Description("Quiet mode"),
		),
		mcp.WithString("timeout",
			mcp.Description("Timeout"),
		),
	)
	s.AddTool(cicdPipelineIntegrationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		target := request.GetString("target", "")
		pipelineStage := request.GetString("pipeline_stage", "")
		artifactName := request.GetString("artifact_name", "")
		outputFormatsStr := request.GetString("output_formats", "")
		var outputFormats []string
		if outputFormatsStr != "" {
			outputFormats = strings.Split(outputFormatsStr, ",")
		}
		outputDirectory := request.GetString("output_directory", "")
		failOnError := request.GetBool("fail_on_error", false)
		quietMode := request.GetBool("quiet_mode", false)
		timeout := request.GetString("timeout", "")

		// CI/CD pipeline integration
		output, err := module.CICDPipelineIntegration(ctx, target, pipelineStage, artifactName, outputFormats, outputDirectory, failOnError, quietMode, timeout)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft CI/CD pipeline integration failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Syft metadata extraction tool
	metadataExtractionTool := mcp.NewTool("syft_metadata_extraction",
		mcp.WithDescription("Metadata extraction using real Syft CLI"),
		mcp.WithString("target",
			mcp.Description("Target for metadata extraction"),
			mcp.Required(),
		),
		mcp.WithString("metadata_types",
			mcp.Description("Comma-separated metadata types"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
		),
		mcp.WithBoolean("include_file_metadata",
			mcp.Description("Include file metadata"),
		),
		mcp.WithBoolean("include_checksums",
			mcp.Description("Include checksums"),
		),
		mcp.WithBoolean("include_certificates",
			mcp.Description("Include certificates"),
		),
		mcp.WithBoolean("include_signatures",
			mcp.Description("Include signatures"),
		),
		mcp.WithString("custom_annotations",
			mcp.Description("Custom annotations"),
		),
	)
	s.AddTool(metadataExtractionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSyftModule(client)

		// Get parameters
		target := request.GetString("target", "")
		metadataTypesStr := request.GetString("metadata_types", "")
		var metadataTypes []string
		if metadataTypesStr != "" {
			metadataTypes = strings.Split(metadataTypesStr, ",")
		}
		outputFormat := request.GetString("output_format", "json")
		includeFileMetadata := request.GetBool("include_file_metadata", false)
		includeChecksums := request.GetBool("include_checksums", false)
		includeCertificates := request.GetBool("include_certificates", false)
		includeSignatures := request.GetBool("include_signatures", false)
		customAnnotations := request.GetString("custom_annotations", "")

		// Metadata extraction
		output, err := module.MetadataExtraction(ctx, target, metadataTypes, outputFormat, includeFileMetadata, includeChecksums, includeCertificates, includeSignatures, customAnnotations)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Syft metadata extraction failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}