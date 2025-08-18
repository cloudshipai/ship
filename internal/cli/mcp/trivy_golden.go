package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTrivyGoldenTools adds enhanced Trivy for golden images MCP tool implementations using real trivy CLI commands
func AddTrivyGoldenTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Trivy image scan tool for golden images
	scanImageTool := mcp.NewTool("trivy_golden_scan_image",
		mcp.WithDescription("Scan container image for vulnerabilities using real trivy CLI"),
		mcp.WithString("image",
			mcp.Description("Container image to scan"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Comma-separated list of severities to include"),
			mcp.Enum("UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif", "template", "cyclonedx", "spdx", "github"),
		),
		mcp.WithString("scanners",
			mcp.Description("Comma-separated list of scanners to use"),
		),
		mcp.WithBoolean("exit_code",
			mcp.Description("Exit with non-zero code if vulnerabilities found"),
		),
		mcp.WithString("output",
			mcp.Description("Output file path"),
		),
	)
	s.AddTool(scanImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		image := request.GetString("image", "")
		args := []string{"trivy", "image"}
		
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if scanners := request.GetString("scanners", ""); scanners != "" {
			args = append(args, "--scanners", scanners)
		}
		if request.GetBool("exit_code", false) {
			args = append(args, "--exit-code", "1")
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		args = append(args, image)
		
		return executeShipCommand(args)
	})

	// Trivy filesystem scan tool for golden image builds
	scanFilesystemTool := mcp.NewTool("trivy_golden_scan_filesystem",
		mcp.WithDescription("Scan filesystem for vulnerabilities and misconfigurations using real trivy CLI"),
		mcp.WithString("path",
			mcp.Description("Path to scan"),
			mcp.Required(),
		),
		mcp.WithString("scanners",
			mcp.Description("Comma-separated list of scanners (vuln,secret,misconfig,license)"),
		),
		mcp.WithString("severity",
			mcp.Description("Comma-separated list of severities to include"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif", "template", "cyclonedx", "spdx"),
		),
		mcp.WithString("output",
			mcp.Description("Output file path"),
		),
	)
	s.AddTool(scanFilesystemTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := request.GetString("path", "")
		args := []string{"trivy", "fs"}
		
		if scanners := request.GetString("scanners", ""); scanners != "" {
			args = append(args, "--scanners", scanners)
		}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		args = append(args, path)
		
		return executeShipCommand(args)
	})

	// Trivy config scan tool for IaC compliance
	scanConfigTool := mcp.NewTool("trivy_golden_scan_config",
		mcp.WithDescription("Scan configuration files for misconfigurations using real trivy CLI"),
		mcp.WithString("path",
			mcp.Description("Path to configuration files"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Comma-separated list of severities to include"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif", "template"),
		),
		mcp.WithString("policy",
			mcp.Description("Path to custom policy file"),
		),
		mcp.WithString("output",
			mcp.Description("Output file path"),
		),
	)
	s.AddTool(scanConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := request.GetString("path", "")
		args := []string{"trivy", "config"}
		
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if policy := request.GetString("policy", ""); policy != "" {
			args = append(args, "--policy", policy)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		args = append(args, path)
		
		return executeShipCommand(args)
	})

	// Trivy SBOM generation tool
	generateSBOMTool := mcp.NewTool("trivy_golden_generate_sbom",
		mcp.WithDescription("Generate Software Bill of Materials (SBOM) for golden image using real trivy CLI"),
		mcp.WithString("image",
			mcp.Description("Container image to generate SBOM for"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("SBOM format"),
			mcp.Enum("cyclonedx", "spdx", "spdx-json", "github"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output file path"),
		),
	)
	s.AddTool(generateSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		image := request.GetString("image", "")
		format := request.GetString("format", "")
		args := []string{"trivy", "image", "--format", format}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		args = append(args, image)
		
		return executeShipCommand(args)
	})

	// Trivy secret detection tool
	scanSecretsTool := mcp.NewTool("trivy_golden_scan_secrets",
		mcp.WithDescription("Scan for secrets in golden image using real trivy CLI"),
		mcp.WithString("target",
			mcp.Description("Target to scan (image name or filesystem path)"),
			mcp.Required(),
		),
		mcp.WithString("target_type",
			mcp.Description("Type of target to scan"),
			mcp.Enum("image", "fs", "repo"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("table", "json", "sarif"),
		),
		mcp.WithString("output",
			mcp.Description("Output file path"),
		),
	)
	s.AddTool(scanSecretsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		targetType := request.GetString("target_type", "")
		args := []string{"trivy", targetType, "--scanners", "secret"}
		
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		args = append(args, target)
		
		return executeShipCommand(args)
	})

	// Trivy version tool
	versionTool := mcp.NewTool("trivy_version",
		mcp.WithDescription("Get Trivy version information using real trivy CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"trivy", "version"}
		return executeShipCommand(args)
	})
}