package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddTrivyGoldenTools adds enhanced Trivy for golden images MCP tool implementations using direct Dagger calls
func AddTrivyGoldenTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTrivyGoldenToolsDirect(s)
}

// addTrivyGoldenToolsDirect adds Trivy golden tools using direct Dagger module calls
func addTrivyGoldenToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyGoldenModule(client)

		// Get parameters
		image := request.GetString("image", "")
		if image == "" {
			return mcp.NewToolResultError("image is required"), nil
		}

		severity := request.GetString("severity", "")
		format := request.GetString("format", "")
		scanners := request.GetString("scanners", "")
		exitCode := request.GetBool("exit_code", false)
		output := request.GetString("output", "")

		// Scan image
		result, err := module.ScanImageBasic(ctx, image, severity, format, scanners, exitCode, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy golden image scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyGoldenModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		scanners := request.GetString("scanners", "")
		severity := request.GetString("severity", "")
		format := request.GetString("format", "")
		output := request.GetString("output", "")

		// Scan filesystem
		result, err := module.ScanFilesystemBasic(ctx, path, scanners, severity, format, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy golden filesystem scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyGoldenModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		severity := request.GetString("severity", "")
		format := request.GetString("format", "")
		policy := request.GetString("policy", "")
		output := request.GetString("output", "")

		// Scan config
		result, err := module.ScanConfigBasic(ctx, path, severity, format, policy, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy golden config scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyGoldenModule(client)

		// Get parameters
		image := request.GetString("image", "")
		format := request.GetString("format", "")
		output := request.GetString("output", "")

		if image == "" {
			return mcp.NewToolResultError("image is required"), nil
		}
		if format == "" {
			return mcp.NewToolResultError("format is required"), nil
		}

		// Generate SBOM
		result, err := module.GenerateSBOMBasic(ctx, image, format, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy golden SBOM generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyGoldenModule(client)

		// Get parameters
		target := request.GetString("target", "")
		targetType := request.GetString("target_type", "")
		format := request.GetString("format", "")
		output := request.GetString("output", "")

		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		if targetType == "" {
			return mcp.NewToolResultError("target_type is required"), nil
		}

		// Scan secrets
		result, err := module.ScanSecretsBasic(ctx, target, targetType, format, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy golden secrets scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Trivy version tool
	versionTool := mcp.NewTool("trivy_version",
		mcp.WithDescription("Get Trivy version information using real trivy CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTrivyGoldenModule(client)

		// Get version
		result, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Trivy golden get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}