package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGrypeTools adds Grype (container vulnerability scanner) MCP tool implementations
func AddGrypeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Grype scan directory tool
	scanDirTool := mcp.NewTool("grype_scan_directory",
		mcp.WithDescription("Scan a directory for vulnerabilities using Grype"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "cyclonedx", "sarif"),
		),
	)
	s.AddTool(scanDirTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "grype"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Grype scan image tool
	scanImageTool := mcp.NewTool("grype_scan_image",
		mcp.WithDescription("Scan a container image for vulnerabilities using Grype"),
		mcp.WithString("image_name",
			mcp.Description("Container image name to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "cyclonedx", "sarif"),
		),
	)
	s.AddTool(scanImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"security", "grype", imageName}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Grype scan SBOM tool
	scanSBOMTool := mcp.NewTool("grype_scan_sbom",
		mcp.WithDescription("Scan SBOM file for vulnerabilities using Grype"),
		mcp.WithString("sbom_path",
			mcp.Description("Path to SBOM file"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "cyclonedx", "sarif"),
		),
	)
	s.AddTool(scanSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sbomPath := request.GetString("sbom_path", "")
		args := []string{"security", "grype", "sbom:" + sbomPath}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Grype scan with severity filter tool
	scanWithSeverityTool := mcp.NewTool("grype_scan_severity",
		mcp.WithDescription("Scan with severity threshold using Grype"),
		mcp.WithString("target",
			mcp.Description("Target to scan (directory, image, or SBOM)"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level"),
			mcp.Required(),
			mcp.Enum("negligible", "low", "medium", "high", "critical"),
		),
	)
	s.AddTool(scanWithSeverityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		severity := request.GetString("severity", "")
		args := []string{"security", "grype", target, "--fail-on", severity}
		return executeShipCommand(args)
	})

	// Grype scan with exclusions tool
	scanWithExclusionsTool := mcp.NewTool("grype_scan_with_exclusions",
		mcp.WithDescription("Scan with package or vulnerability exclusions using Grype"),
		mcp.WithString("target",
			mcp.Description("Target to scan (directory, image, or SBOM)"),
			mcp.Required(),
		),
		mcp.WithString("exclude_patterns",
			mcp.Description("Comma-separated exclusion patterns (package names or CVE IDs)"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "cyclonedx", "sarif"),
		),
	)
	s.AddTool(scanWithExclusionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		excludePatterns := request.GetString("exclude_patterns", "")
		args := []string{"security", "grype", target}
		
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		
		// Add exclusions
		for _, pattern := range strings.Split(excludePatterns, ",") {
			if strings.TrimSpace(pattern) != "" {
				args = append(args, "--exclude", strings.TrimSpace(pattern))
			}
		}
		return executeShipCommand(args)
	})

	// Grype scan only fixed vulnerabilities tool
	scanOnlyFixedTool := mcp.NewTool("grype_scan_only_fixed",
		mcp.WithDescription("Scan and report only vulnerabilities with available fixes"),
		mcp.WithString("target",
			mcp.Description("Target to scan (directory, image, or SBOM)"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "cyclonedx", "sarif"),
		),
	)
	s.AddTool(scanOnlyFixedTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"security", "grype", target, "--only-fixed"}
		
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Grype scan with platform specification tool
	scanWithPlatformTool := mcp.NewTool("grype_scan_with_platform",
		mcp.WithDescription("Scan multi-platform container image with specific platform"),
		mcp.WithString("image_name",
			mcp.Description("Container image name to scan"),
			mcp.Required(),
		),
		mcp.WithString("platform",
			mcp.Description("Platform specification (e.g., linux/amd64, linux/arm64)"),
			mcp.Required(),
		),
		mcp.WithString("scope",
			mcp.Description("Search scope for vulnerability detection"),
			mcp.Enum("Squashed", "AllLayers"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "cyclonedx", "sarif"),
		),
	)
	s.AddTool(scanWithPlatformTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		platform := request.GetString("platform", "")
		args := []string{"security", "grype", imageName, "--platform", platform}
		
		if scope := request.GetString("scope", ""); scope != "" {
			args = append(args, "--scope", scope)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Grype generate comprehensive report tool
	generateReportTool := mcp.NewTool("grype_generate_report",
		mcp.WithDescription("Generate comprehensive vulnerability report with analysis"),
		mcp.WithString("target",
			mcp.Description("Target to scan and report on"),
			mcp.Required(),
		),
		mcp.WithString("report_type",
			mcp.Description("Type of report to generate"),
			mcp.Enum("executive", "technical", "compliance", "security-team"),
		),
		mcp.WithString("severity_threshold",
			mcp.Description("Minimum severity to include in report"),
			mcp.Enum("negligible", "low", "medium", "high", "critical"),
		),
		mcp.WithBoolean("include_fixed",
			mcp.Description("Include vulnerabilities with available fixes"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"security", "grype", target, "--output", "json"}
		
		if severityThreshold := request.GetString("severity_threshold", ""); severityThreshold != "" {
			args = append(args, "--fail-on", severityThreshold)
		}
		if request.GetBool("include_fixed", false) {
			args = append(args, "--only-fixed")
		}
		
		// Add report-specific formatting
		reportType := request.GetString("report_type", "technical")
		args = append(args, "--report-type", reportType)
		
		return executeShipCommand(args)
	})

	// Grype get version tool
	getVersionTool := mcp.NewTool("grype_get_version",
		mcp.WithDescription("Get Grype version and database information"),
		mcp.WithBoolean("show_db_info",
			mcp.Description("Include vulnerability database information"),
		),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "grype", "--version"}
		if request.GetBool("show_db_info", false) {
			args = append(args, "--db-status")
		}
		return executeShipCommand(args)
	})
}