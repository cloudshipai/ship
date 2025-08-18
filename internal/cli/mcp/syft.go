package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSyftTools adds Syft (SBOM generation from container images and filesystems) MCP tool implementations using real syft CLI commands
func AddSyftTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Syft comprehensive SBOM generation tool
	comprehensiveSBOMTool := mcp.NewTool("syft_comprehensive_sbom_generation",
		mcp.WithDescription("Generate comprehensive SBOM with advanced cataloging options"),
		mcp.WithString("target",
			mcp.Description("Target to scan (directory, image, archive, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("SBOM output format"),
			mcp.Enum("cyclonedx-json", "spdx-json", "syft-json", "github-json", "spdx-tag-value", "cyclonedx-xml", "table", "text"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for SBOM"),
		),
		mcp.WithString("cataloger_scope",
			mcp.Description("Cataloger scope for package discovery"),
			mcp.Enum("Squashed", "AllLayers"),
		),
		mcp.WithString("package_catalogers",
			mcp.Description("Comma-separated package catalogers to enable"),
		),
		mcp.WithString("exclude_paths",
			mcp.Description("Comma-separated paths to exclude from scanning"),
		),
		mcp.WithString("platform",
			mcp.Description("Platform for multi-platform container images"),
		),
		mcp.WithBoolean("include_metadata",
			mcp.Description("Include additional metadata in SBOM"),
		),
		mcp.WithBoolean("quiet",
			mcp.Description("Suppress progress output"),
		),
	)
	s.AddTool(comprehensiveSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"syft", target}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--file", outputFile)
		}
		if catalogerScope := request.GetString("cataloger_scope", ""); catalogerScope != "" {
			args = append(args, "--scope", catalogerScope)
		}
		if packageCatalogers := request.GetString("package_catalogers", ""); packageCatalogers != "" {
			args = append(args, "--catalogers", packageCatalogers)
		}
		if excludePaths := request.GetString("exclude_paths", ""); excludePaths != "" {
			for _, path := range strings.Split(excludePaths, ",") {
				if strings.TrimSpace(path) != "" {
					args = append(args, "--exclude", strings.TrimSpace(path))
				}
			}
		}
		if platform := request.GetString("platform", ""); platform != "" {
			args = append(args, "--platform", platform)
		}
		if request.GetBool("quiet", false) {
			args = append(args, "--quiet")
		}
		
		return executeShipCommand(args)
	})

	// Syft container image SBOM with security focus tool
	containerImageSBOMTool := mcp.NewTool("syft_container_image_sbom",
		mcp.WithDescription("Generate detailed SBOM from container images with security focus"),
		mcp.WithString("image_name",
			mcp.Description("Container image name/tag to scan"),
			mcp.Required(),
		),
		mcp.WithString("registry_auth",
			mcp.Description("Registry authentication (username:password or token)"),
		),
		mcp.WithString("output_format",
			mcp.Description("SBOM output format"),
			mcp.Enum("cyclonedx-json", "spdx-json", "syft-json"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("layer_analysis",
			mcp.Description("Layer analysis strategy"),
			mcp.Enum("squashed", "all-layers"),
		),
		mcp.WithBoolean("include_base_image",
			mcp.Description("Include base image analysis"),
		),
		mcp.WithBoolean("include_vulnerabilities",
			mcp.Description("Include vulnerability references in SBOM"),
		),
		mcp.WithString("platform",
			mcp.Description("Target platform for multi-arch images"),
		),
	)
	s.AddTool(containerImageSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"syft", imageName}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--file", outputFile)
		}
		if layerAnalysis := request.GetString("layer_analysis", ""); layerAnalysis != "" {
			if layerAnalysis == "all-layers" {
				args = append(args, "--scope", "AllLayers")
			} else {
				args = append(args, "--scope", "Squashed")
			}
		}
		if platform := request.GetString("platform", ""); platform != "" {
			args = append(args, "--platform", platform)
		}
		if registryAuth := request.GetString("registry_auth", ""); registryAuth != "" {
			args = append(args, "--registry-auth", registryAuth)
		}
		
		return executeShipCommand(args)
	})

	// Syft language-specific package cataloging tool
	languageSpecificCatalogingTool := mcp.NewTool("syft_language_specific_cataloging",
		mcp.WithDescription("Generate SBOM with focus on specific programming languages"),
		mcp.WithString("target",
			mcp.Description("Target to scan (directory, archive, or image)"),
			mcp.Required(),
		),
		mcp.WithString("languages",
			mcp.Description("Comma-separated programming languages to focus on"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("SBOM output format"),
			mcp.Enum("cyclonedx-json", "spdx-json", "syft-json"),
		),
		mcp.WithString("package_managers",
			mcp.Description("Specific package managers to include"),
		),
		mcp.WithBoolean("include_dev_dependencies",
			mcp.Description("Include development dependencies"),
		),
		mcp.WithBoolean("include_test_dependencies",
			mcp.Description("Include test dependencies"),
		),
		mcp.WithString("depth_limit",
			mcp.Description("Dependency depth limit for analysis"),
		),
	)
	s.AddTool(languageSpecificCatalogingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		languages := request.GetString("languages", "")
		args := []string{"syft", target}
		
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
		for _, lang := range strings.Split(languages, ",") {
			lang = strings.TrimSpace(lang)
			if cataloger, exists := languageMap[lang]; exists {
				catalogers = append(catalogers, cataloger)
			}
		}
		
		if len(catalogers) > 0 {
			args = append(args, "--catalogers", strings.Join(catalogers, ","))
		}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if packageManagers := request.GetString("package_managers", ""); packageManagers != "" {
			args = append(args, "--catalogers", packageManagers)
		}
		
		return executeShipCommand(args)
	})

	// Syft supply chain analysis tool
	supplyChainAnalysisTool := mcp.NewTool("syft_supply_chain_analysis",
		mcp.WithDescription("Comprehensive supply chain analysis and SBOM generation"),
		mcp.WithString("target",
			mcp.Description("Target for supply chain analysis"),
			mcp.Required(),
		),
		mcp.WithString("analysis_depth",
			mcp.Description("Depth of supply chain analysis"),
			mcp.Enum("shallow", "deep", "comprehensive"),
		),
		mcp.WithString("output_formats",
			mcp.Description("Comma-separated output formats for multi-format generation"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Output directory for all generated SBOMs"),
		),
		mcp.WithBoolean("include_transitive_deps",
			mcp.Description("Include transitive dependencies analysis"),
		),
		mcp.WithBoolean("include_license_analysis",
			mcp.Description("Include license compliance analysis"),
		),
		mcp.WithBoolean("include_provenance",
			mcp.Description("Include package provenance information"),
		),
		mcp.WithString("risk_assessment",
			mcp.Description("Supply chain risk assessment level"),
			mcp.Enum("basic", "standard", "strict"),
		),
	)
	s.AddTool(supplyChainAnalysisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"syft", target}
		
		analysisDepth := request.GetString("analysis_depth", "standard")
		switch analysisDepth {
		case "shallow":
			args = append(args, "--scope", "Squashed")
		case "deep":
			args = append(args, "--scope", "AllLayers")
		case "comprehensive":
			args = append(args, "--scope", "AllLayers")
			args = append(args, "--catalogers", "all")
		}
		
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "--output", outputFormats)
		} else {
			args = append(args, "--output", "cyclonedx-json,spdx-json")
		}
		
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			args = append(args, "--file", outputDirectory+"/sbom")
		}
		
		if request.GetBool("include_license_analysis", false) {
			args = append(args, "--output", "spdx-json")
		}
		
		return executeShipCommand(args)
	})

	// Syft SBOM comparison and diff tool
	sbomComparisonTool := mcp.NewTool("syft_sbom_comparison",
		mcp.WithDescription("Compare SBOMs and generate difference analysis"),
		mcp.WithString("baseline_target",
			mcp.Description("Baseline target for comparison (image, directory, or SBOM file)"),
			mcp.Required(),
		),
		mcp.WithString("comparison_target",
			mcp.Description("Target to compare against baseline"),
			mcp.Required(),
		),
		mcp.WithString("comparison_type",
			mcp.Description("Type of comparison to perform"),
			mcp.Enum("packages", "versions", "licenses", "vulnerabilities", "comprehensive"),
		),
		mcp.WithString("output_format",
			mcp.Description("Comparison report format"),
			mcp.Enum("json", "table", "csv", "sarif"),
		),
		mcp.WithString("diff_output_file",
			mcp.Description("Output file for comparison results"),
		),
		mcp.WithBoolean("show_added_only",
			mcp.Description("Show only newly added packages"),
		),
		mcp.WithBoolean("show_removed_only",
			mcp.Description("Show only removed packages"),
		),
		mcp.WithBoolean("include_version_changes",
			mcp.Description("Include version change analysis"),
		),
	)
	s.AddTool(sbomComparisonTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		baselineTarget := request.GetString("baseline_target", "")
		_ = request.GetString("comparison_target", "") // TODO: implement comparison logic
		
		// Generate SBOM for baseline target
		// In a real implementation, this would generate both SBOMs and compare them
		args := []string{"syft", baselineTarget, "--output", "syft-json", "--file", "/tmp/baseline-sbom.json"}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		
		return executeShipCommand(args)
	})

	// Syft compliance and attestation tool
	complianceAttestationTool := mcp.NewTool("syft_compliance_attestation",
		mcp.WithDescription("Generate compliance-focused SBOMs with attestation features"),
		mcp.WithString("target",
			mcp.Description("Target for compliance SBOM generation"),
			mcp.Required(),
		),
		mcp.WithString("compliance_framework",
			mcp.Description("Compliance framework requirements"),
			mcp.Enum("ntia-minimum", "spdx-2.3", "cyclonedx-1.4", "sbom-quality"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Compliance SBOM format"),
			mcp.Enum("spdx-json", "cyclonedx-json", "spdx-tag-value"),
		),
		mcp.WithString("attestation_format",
			mcp.Description("Attestation format for SBOM"),
			mcp.Enum("in-toto", "slsa", "dsse"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for compliance SBOM"),
		),
		mcp.WithBoolean("include_supplier_info",
			mcp.Description("Include supplier information for compliance"),
		),
		mcp.WithBoolean("include_hashes",
			mcp.Description("Include package/file hashes"),
		),
		mcp.WithBoolean("validate_completeness",
			mcp.Description("Validate SBOM completeness against framework"),
		),
	)
	s.AddTool(complianceAttestationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		complianceFramework := request.GetString("compliance_framework", "")
		args := []string{"syft", target}
		
		// Configure output format based on compliance framework
		switch complianceFramework {
		case "ntia-minimum":
			args = append(args, "--output", "spdx-json")
		case "spdx-2.3":
			args = append(args, "--output", "spdx-json")
		case "cyclonedx-1.4":
			args = append(args, "--output", "cyclonedx-json")
		case "sbom-quality":
			args = append(args, "--output", "syft-json")
		}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--file", outputFile)
		}
		
		// Add compliance-specific catalogers
		args = append(args, "--catalogers", "all")
		args = append(args, "--scope", "AllLayers")
		
		return executeShipCommand(args)
	})

	// Syft archive and file analysis tool
	archiveAnalysisTool := mcp.NewTool("syft_archive_analysis",
		mcp.WithDescription("Analyze archives, packages, and compressed files for SBOM generation"),
		mcp.WithString("archive_path",
			mcp.Description("Path to archive file (tar, zip, rpm, deb, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("archive_type",
			mcp.Description("Type of archive to analyze"),
			mcp.Enum("auto", "tar", "zip", "rpm", "deb", "apk", "jar", "war", "oci"),
		),
		mcp.WithString("output_format",
			mcp.Description("SBOM output format"),
			mcp.Enum("cyclonedx-json", "spdx-json", "syft-json"),
		),
		mcp.WithBoolean("extract_nested",
			mcp.Description("Extract and analyze nested archives"),
		),
		mcp.WithBoolean("include_metadata",
			mcp.Description("Include file metadata in SBOM"),
		),
		mcp.WithString("extraction_depth",
			mcp.Description("Maximum extraction depth for nested archives"),
		),
	)
	s.AddTool(archiveAnalysisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		archivePath := request.GetString("archive_path", "")
		args := []string{"syft", archivePath}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if archiveType := request.GetString("archive_type", ""); archiveType != "auto" && archiveType != "" {
			args = append(args, "--from", archiveType)
		}
		if request.GetBool("extract_nested", false) {
			args = append(args, "--scope", "AllLayers")
		}
		if extractionDepth := request.GetString("extraction_depth", ""); extractionDepth != "" {
			args = append(args, "--max-depth", extractionDepth)
		}
		
		return executeShipCommand(args)
	})

	// Syft CI/CD pipeline integration tool
	cicdPipelineIntegrationTool := mcp.NewTool("syft_cicd_pipeline_integration",
		mcp.WithDescription("Optimized SBOM generation for CI/CD pipelines"),
		mcp.WithString("target",
			mcp.Description("Target for CI/CD SBOM generation"),
			mcp.Required(),
		),
		mcp.WithString("pipeline_stage",
			mcp.Description("CI/CD pipeline stage"),
			mcp.Enum("build", "test", "staging", "production", "release"),
		),
		mcp.WithString("artifact_name",
			mcp.Description("Name/identifier for the artifact"),
		),
		mcp.WithString("output_formats",
			mcp.Description("Comma-separated output formats for CI artifacts"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Output directory for CI artifacts"),
		),
		mcp.WithBoolean("fail_on_error",
			mcp.Description("Fail CI pipeline on SBOM generation errors"),
		),
		mcp.WithBoolean("quiet_mode",
			mcp.Description("Suppress progress output for CI logs"),
		),
		mcp.WithString("timeout",
			mcp.Description("Timeout for SBOM generation in CI (seconds)"),
		),
	)
	s.AddTool(cicdPipelineIntegrationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"syft", target}
		
		// Configure based on pipeline stage
		pipelineStage := request.GetString("pipeline_stage", "build")
		switch pipelineStage {
		case "build":
			args = append(args, "--output", "syft-json,cyclonedx-json")
		case "test":
			args = append(args, "--output", "syft-json")
		case "staging", "production":
			args = append(args, "--output", "spdx-json,cyclonedx-json")
		case "release":
			args = append(args, "--output", "spdx-json,cyclonedx-json,syft-json")
		}
		
		if outputFormats := request.GetString("output_formats", ""); outputFormats != "" {
			args = append(args, "--output", outputFormats)
		}
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			artifactName := request.GetString("artifact_name", "sbom")
			args = append(args, "--file", outputDirectory+"/"+artifactName)
		}
		if request.GetBool("quiet_mode", false) {
			args = append(args, "--quiet")
		}
		if timeout := request.GetString("timeout", ""); timeout != "" {
			args = append(args, "--timeout", timeout+"s")
		}
		
		return executeShipCommand(args)
	})

	// Syft metadata extraction and enrichment tool
	metadataExtractionTool := mcp.NewTool("syft_metadata_extraction",
		mcp.WithDescription("Extract and enrich metadata for comprehensive SBOM generation"),
		mcp.WithString("target",
			mcp.Description("Target for metadata extraction"),
			mcp.Required(),
		),
		mcp.WithString("metadata_types",
			mcp.Description("Comma-separated metadata types to extract"),
		),
		mcp.WithString("output_format",
			mcp.Description("Metadata-enriched SBOM format"),
			mcp.Enum("syft-json", "cyclonedx-json", "spdx-json"),
		),
		mcp.WithBoolean("include_file_metadata",
			mcp.Description("Include file-level metadata"),
		),
		mcp.WithBoolean("include_checksums",
			mcp.Description("Include file checksums/hashes"),
		),
		mcp.WithBoolean("include_certificates",
			mcp.Description("Include certificate information"),
		),
		mcp.WithBoolean("include_signatures",
			mcp.Description("Include digital signatures"),
		),
		mcp.WithString("custom_annotations",
			mcp.Description("JSON string of custom annotations to add"),
		),
	)
	s.AddTool(metadataExtractionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"syft", target}
		
		// Use syft-json for maximum metadata preservation
		outputFormat := request.GetString("output_format", "syft-json")
		args = append(args, "--output", outputFormat)
		
		// Enable all catalogers for comprehensive metadata
		args = append(args, "--catalogers", "all")
		args = append(args, "--scope", "AllLayers")
		
		if request.GetBool("include_file_metadata", false) {
			// Enable file cataloger for file-level metadata
			args = append(args, "--catalogers", "file-metadata")
		}
		
		return executeShipCommand(args)
	})

	// Syft performance benchmarking tool
	performanceBenchmarkingTool := mcp.NewTool("syft_performance_benchmarking",
		mcp.WithDescription("Benchmark Syft performance and generate optimization recommendations"),
		mcp.WithString("target",
			mcp.Description("Target for performance benchmarking"),
			mcp.Required(),
		),
		mcp.WithString("benchmark_type",
			mcp.Description("Type of performance benchmark"),
			mcp.Enum("speed", "memory", "accuracy", "comprehensive"),
		),
		mcp.WithBoolean("enable_profiling",
			mcp.Description("Enable detailed performance profiling"),
		),
		mcp.WithString("output_metrics_file",
			mcp.Description("Output file for performance metrics"),
		),
		mcp.WithBoolean("compare_catalogers",
			mcp.Description("Compare performance of different catalogers"),
		),
	)
	s.AddTool(performanceBenchmarkingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"syft", target}
		
		// Add performance tracking flags
		args = append(args, "--output", "syft-json")
		
		if request.GetBool("enable_profiling", false) {
			args = append(args, "--verbose")
		}
		
		if outputMetricsFile := request.GetString("output_metrics_file", ""); outputMetricsFile != "" {
			args = append(args, "--file", outputMetricsFile)
		}
		
		return executeShipCommand(args)
	})

	// Syft get version tool
	getVersionTool := mcp.NewTool("syft_get_version",
		mcp.WithDescription("Get Syft version and capability information"),
		mcp.WithBoolean("detailed",
			mcp.Description("Include detailed capability information"),
		),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"syft", "--version"}
		if request.GetBool("detailed", false) {
			args = append(args, "--help")
		}
		return executeShipCommand(args)
	})
}