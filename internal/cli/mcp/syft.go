package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSyftTools adds Syft (SBOM generation) MCP tool implementations
func AddSyftTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Syft generate SBOM from directory tool
	generateDirTool := mcp.NewTool("syft_generate_sbom_directory",
		mcp.WithDescription("Generate SBOM from a directory using Syft"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "spdx-json", "cyclonedx-json", "table"),
		),
	)
	s.AddTool(generateDirTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "syft"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Syft generate SBOM from image tool
	generateImageTool := mcp.NewTool("syft_generate_sbom_image",
		mcp.WithDescription("Generate SBOM from a container image using Syft"),
		mcp.WithString("image_name",
			mcp.Description("Container image name to scan"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "spdx-json", "cyclonedx-json", "table"),
		),
	)
	s.AddTool(generateImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"security", "syft", imageName}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Syft generate SBOM for specific package type tool
	generatePackageTool := mcp.NewTool("syft_generate_sbom_package",
		mcp.WithDescription("Generate SBOM for specific package type using Syft"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("package_type",
			mcp.Description("Specific package type to scan"),
			mcp.Required(),
			mcp.Enum("npm", "pip", "gem", "go", "java", "rust", "dotnet"),
		),
	)
	s.AddTool(generatePackageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "syft"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		packageType := request.GetString("package_type", "")
		args = append(args, "--catalogers", packageType)
		return executeShipCommand(args)
	})

	// Syft generate comprehensive SBOM with multiple formats tool
	generateComprehensiveTool := mcp.NewTool("syft_generate_comprehensive_sbom",
		mcp.WithDescription("Generate comprehensive SBOM with advanced options"),
		mcp.WithString("target",
			mcp.Description("Target image or directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("SBOM output format"),
			mcp.Enum("cyclonedx-json", "spdx-json", "syft-json", "github-json", "spdx-tag-value", "cyclonedx-xml", "table"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("catalogers",
			mcp.Description("Comma-separated catalogers (python,javascript,go,java,ruby)"),
		),
		mcp.WithString("scope",
			mcp.Description("Search scope for cataloging packages"),
			mcp.Enum("Squashed", "AllLayers"),
		),
		mcp.WithString("platform",
			mcp.Description("Platform for multi-platform images (linux/amd64, linux/arm64)"),
		),
	)
	s.AddTool(generateComprehensiveTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"security", "syft", target}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--file", outputFile)
		}
		if catalogers := request.GetString("catalogers", ""); catalogers != "" {
			args = append(args, "--catalogers", catalogers)
		}
		if scope := request.GetString("scope", ""); scope != "" {
			args = append(args, "--scope", scope)
		}
		if platform := request.GetString("platform", ""); platform != "" {
			args = append(args, "--platform", platform)
		}
		return executeShipCommand(args)
	})

	// Syft scan with exclusions tool
	scanWithExclusionsTool := mcp.NewTool("syft_scan_with_exclusions",
		mcp.WithDescription("Scan with path exclusions using Syft"),
		mcp.WithString("target",
			mcp.Description("Target to scan"),
			mcp.Required(),
		),
		mcp.WithString("exclude_paths",
			mcp.Description("Comma-separated paths to exclude"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "cyclonedx-json", "spdx-json"),
		),
	)
	s.AddTool(scanWithExclusionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		excludePaths := request.GetString("exclude_paths", "")
		args := []string{"security", "syft", target}
		
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--output", outputFormat)
		}
		
		// Add exclusions
		for _, path := range strings.Split(excludePaths, ",") {
			if strings.TrimSpace(path) != "" {
				args = append(args, "--exclude", strings.TrimSpace(path))
			}
		}
		return executeShipCommand(args)
	})

	// Syft analyze existing SBOM tool
	analyzeSBOMTool := mcp.NewTool("syft_analyze_sbom",
		mcp.WithDescription("Analyze packages from existing SBOM file"),
		mcp.WithString("sbom_file",
			mcp.Description("Path to SBOM file to analyze"),
			mcp.Required(),
		),
		mcp.WithString("analysis_type",
			mcp.Description("Type of analysis to perform"),
			mcp.Enum("summary", "detailed", "security-focused", "licensing"),
		),
		mcp.WithString("filter_language",
			mcp.Description("Filter packages by programming language"),
		),
	)
	s.AddTool(analyzeSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sbomFile := request.GetString("sbom_file", "")
		args := []string{"security", "syft", "--analyze-sbom", sbomFile}
		
		if analysisType := request.GetString("analysis_type", ""); analysisType != "" {
			args = append(args, "--analysis-type", analysisType)
		}
		if filterLang := request.GetString("filter_language", ""); filterLang != "" {
			args = append(args, "--filter-language", filterLang)
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
		args := []string{"security", "syft", "--version"}
		if request.GetBool("detailed", false) {
			args = append(args, "--help")
		}
		return executeShipCommand(args)
	})
}