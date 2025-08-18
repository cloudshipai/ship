package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddOSVScannerTools adds OSV Scanner (Open Source Vulnerability scanner) MCP tool implementations using direct Dagger calls
func AddOSVScannerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addOSVScannerToolsDirect(s)
}

// addOSVScannerToolsDirect adds OSV Scanner tools using direct Dagger module calls
func addOSVScannerToolsDirect(s *server.MCPServer) {
	// OSV Scanner scan source directory tool
	scanSourceTool := mcp.NewTool("osv_scanner_scan_source",
		mcp.WithDescription("Scan source directory for open source vulnerabilities using osv-scanner"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively scan subdirectories"),
		),
		mcp.WithString("output",
			mcp.Description("Output file path to save results"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "sarif"),
		),
		mcp.WithString("config",
			mcp.Description("Override configuration file path"),
		),
		mcp.WithBoolean("licenses",
			mcp.Description("Include license scanning"),
		),
	)
	s.AddTool(scanSourceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		directory := request.GetString("directory", ".")

		// Scan directory - use the ScanDirectory function from Dagger
		output, err := module.ScanDirectory(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OSV Scanner scan container image tool
	scanImageTool := mcp.NewTool("osv_scanner_scan_image",
		mcp.WithDescription("Scan container image for vulnerabilities using osv-scanner"),
		mcp.WithString("image",
			mcp.Description("Container image name and tag (e.g., my-image:latest)"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output file path to save results"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "sarif"),
		),
		mcp.WithString("config",
			mcp.Description("Override configuration file path"),
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
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		image := request.GetString("image", "")
		if image == "" {
			return mcp.NewToolResultError("image is required"), nil
		}
		output := request.GetString("output", "")
		format := request.GetString("format", "")
		config := request.GetString("config", "")

		// Scan image
		result, err := module.ScanImage(ctx, image, output, format, config)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner image scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// OSV Scanner scan lockfile tool
	scanLockfileTool := mcp.NewTool("osv_scanner_scan_lockfile",
		mcp.WithDescription("Scan specific lockfile for vulnerabilities using osv-scanner"),
		mcp.WithString("lockfile_path",
			mcp.Description("Path to lockfile (package-lock.json, go.mod, requirements.txt, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output file path to save results"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "sarif"),
		),
		mcp.WithBoolean("all_packages",
			mcp.Description("Output all packages in JSON format"),
		),
	)
	s.AddTool(scanLockfileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		lockfilePath := request.GetString("lockfile_path", "")
		if lockfilePath == "" {
			return mcp.NewToolResultError("lockfile_path is required"), nil
		}

		// Scan lockfile
		output, err := module.ScanLockfile(ctx, lockfilePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner lockfile scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OSV Scanner scan manifest tool
	scanManifestTool := mcp.NewTool("osv_scanner_scan_manifest",
		mcp.WithDescription("Scan package manifest file for vulnerabilities using osv-scanner"),
		mcp.WithString("manifest_path",
			mcp.Description("Path to manifest file (package.json, go.mod, pom.xml, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output file path to save results"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "sarif"),
		),
	)
	s.AddTool(scanManifestTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		manifestPath := request.GetString("manifest_path", "")
		if manifestPath == "" {
			return mcp.NewToolResultError("manifest_path is required"), nil
		}
		output := request.GetString("output", "")
		format := request.GetString("format", "")

		// Scan manifest
		result, err := module.ScanManifest(ctx, manifestPath, output, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner manifest scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// OSV Scanner license scanning tool
	licenseScanTool := mcp.NewTool("osv_scanner_license_scan",
		mcp.WithDescription("Scan for license compliance using osv-scanner"),
		mcp.WithString("path",
			mcp.Description("Path to scan for licenses"),
			mcp.Required(),
		),
		mcp.WithString("allowed_licenses",
			mcp.Description("Comma-separated list of allowed licenses (e.g., MIT,Apache-2.0)"),
		),
		mcp.WithString("output",
			mcp.Description("Output file path to save results"),
		),
	)
	s.AddTool(licenseScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}
		allowedLicenses := request.GetString("allowed_licenses", "")
		output := request.GetString("output", "")

		// Scan licenses
		result, err := module.LicenseScan(ctx, path, allowedLicenses, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner license scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// OSV Scanner offline scanning tool
	offlineScanTool := mcp.NewTool("osv_scanner_offline_scan",
		mcp.WithDescription("Scan using offline vulnerability databases with osv-scanner"),
		mcp.WithString("path",
			mcp.Description("Path to scan"),
			mcp.Required(),
		),
		mcp.WithString("offline_db_path",
			mcp.Description("Path to offline vulnerability databases"),
		),
		mcp.WithBoolean("download_databases",
			mcp.Description("Download offline databases before scanning"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively scan subdirectories"),
		),
	)
	s.AddTool(offlineScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}
		offlineDbPath := request.GetString("offline_db_path", "")
		downloadDatabases := request.GetBool("download_databases", false)
		recursive := request.GetBool("recursive", false)

		// Offline scan
		output, err := module.OfflineScan(ctx, path, offlineDbPath, downloadDatabases, recursive)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner offline scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OSV Scanner guided remediation tool
	fixTool := mcp.NewTool("osv_scanner_fix",
		mcp.WithDescription("Apply guided remediation for vulnerabilities using osv-scanner"),
		mcp.WithString("manifest_path",
			mcp.Description("Path to package manifest file (package.json, etc.)"),
		),
		mcp.WithString("lockfile_path",
			mcp.Description("Path to lockfile (package-lock.json, etc.)"),
		),
		mcp.WithString("strategy",
			mcp.Description("Remediation strategy"),
			mcp.Enum("in-place", "relock"),
		),
		mcp.WithString("max_depth",
			mcp.Description("Maximum dependency depth to consider"),
		),
		mcp.WithString("min_severity",
			mcp.Description("Minimum vulnerability severity to fix"),
		),
		mcp.WithBoolean("ignore_dev",
			mcp.Description("Ignore development dependencies"),
		),
	)
	s.AddTool(fixTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		manifestPath := request.GetString("manifest_path", "")
		lockfilePath := request.GetString("lockfile_path", "")
		strategy := request.GetString("strategy", "")
		maxDepth := request.GetString("max_depth", "")
		minSeverity := request.GetString("min_severity", "")
		ignoreDev := request.GetBool("ignore_dev", false)

		// Apply fix
		output, err := module.Fix(ctx, manifestPath, lockfilePath, strategy, maxDepth, minSeverity, ignoreDev)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner fix failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OSV Scanner serve report tool
	serveReportTool := mcp.NewTool("osv_scanner_serve_report",
		mcp.WithDescription("Generate and serve HTML vulnerability report locally using osv-scanner"),
		mcp.WithString("path",
			mcp.Description("Path to scan and generate report for"),
			mcp.Required(),
		),
		mcp.WithString("port",
			mcp.Description("Port to serve the report on (default: 8080)"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively scan subdirectories"),
		),
	)
	s.AddTool(serveReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}
		port := request.GetString("port", "")
		recursive := request.GetBool("recursive", false)

		// Serve report
		output, err := module.ServeReport(ctx, path, port, recursive)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner serve report failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OSV Scanner SBOM scan tool
	scanSBOMTool := mcp.NewTool("osv_scanner_scan_sbom",
		mcp.WithDescription("Scan SBOM file for vulnerabilities using osv-scanner"),
		mcp.WithString("sbom_path",
			mcp.Description("Path to SBOM file (SPDX or CycloneDX format)"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Output file path to save results"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "sarif"),
		),
	)
	s.AddTool(scanSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		sbomPath := request.GetString("sbom_path", "")
		if sbomPath == "" {
			return mcp.NewToolResultError("sbom_path is required"), nil
		}

		// Scan SBOM
		output, err := module.ScanSBOM(ctx, sbomPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner SBOM scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OSV Scanner verbosity control tool
	verboseScanTool := mcp.NewTool("osv_scanner_verbose_scan",
		mcp.WithDescription("Run OSV Scanner with verbose logging using osv-scanner"),
		mcp.WithString("path",
			mcp.Description("Path to scan"),
			mcp.Required(),
		),
		mcp.WithString("verbosity",
			mcp.Description("Logging verbosity level"),
			mcp.Enum("error", "warn", "info", "debug"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively scan subdirectories"),
		),
	)
	s.AddTool(verboseScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}
		verbosity := request.GetString("verbosity", "")
		recursive := request.GetBool("recursive", false)

		// Verbose scan
		output, err := module.VerboseScan(ctx, path, verbosity, recursive)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("osv-scanner verbose scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}