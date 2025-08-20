package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddOSVScannerTools adds OSV Scanner (Open Source Vulnerability scanner) MCP tool implementations using direct Dagger calls
func AddOSVScannerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Add the new osv_scan tool first
	addNewOSVScanTool(s)
	
	// Keep existing tools for backward compatibility
	addOSVScannerToolsDirect(s)
}

// addNewOSVScanTool adds the new unified OSV scanning tool
func addNewOSVScanTool(s *server.MCPServer) {
	// OSV scan tool - unified interface
	scanTool := mcp.NewTool("osv_scan",
		mcp.WithDescription("OSV-Scanner over source/SBOM/image; optional license allowlist"),
		mcp.WithString("mode",
			mcp.Description("Scan mode: source, sbom, or image"),
			mcp.Required(),
		),
		mcp.WithString("path_or_ref",
			mcp.Description("Repository path, SBOM path, or image reference"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format: json, sarif, table, markdown, html (default: json)"),
		),
		mcp.WithString("licenses_allowlist",
			mcp.Description("Comma-separated SPDX license IDs (e.g., MIT,Apache-2.0,BSD-3-Clause)"),
		),
	)
	s.AddTool(scanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOSVScannerModule(client)

		// Get parameters
		mode := request.GetString("mode", "")
		if mode == "" {
			return mcp.NewToolResultError("mode is required"), nil
		}
		
		pathOrRef := request.GetString("path_or_ref", "")
		if pathOrRef == "" {
			return mcp.NewToolResultError("path_or_ref is required"), nil
		}
		
		format := request.GetString("format", "json")
		licensesAllowlist := request.GetString("licenses_allowlist", "")

		// Generate output
		var stdout string
		var stderr string
		
		// Execute based on mode
		switch mode {
		case "source":
			stdout, err = module.ScanSource(ctx, pathOrRef, format, licensesAllowlist)
		case "sbom":
			stdout, err = module.ScanSBOM(ctx, pathOrRef, format, licensesAllowlist)
		case "image":
			stdout, err = module.ScanImage(ctx, pathOrRef, format, licensesAllowlist)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("invalid mode: %s (must be source, sbom, or image)", mode)), nil
		}
		
		// Build result in the expected format
		result := map[string]interface{}{
			"status": "ok",
			"stdout": stdout,
			"stderr": stderr,
			"artifacts": map[string]string{},
			"summary": map[string]interface{}{
				"high":               0,
				"critical":           0,
				"license_violations": 0,
			},
			"diagnostics": []string{},
		}
		
		// Add artifact path based on format
		switch format {
		case "sarif":
			result["artifacts"].(map[string]string)["osv_sarif"] = "./osv.sarif"
		case "json":
			result["artifacts"].(map[string]string)["osv_json"] = "./osv.json"
		default:
			// For other formats, we just return the stdout
		}
		
		// Parse JSON output to extract counts if format is json
		if format == "json" && err == nil && stdout != "" {
			var osvOutput map[string]interface{}
			if jsonErr := json.Unmarshal([]byte(stdout), &osvOutput); jsonErr == nil {
				// Try to extract vulnerability counts from OSV output
				if results, ok := osvOutput["results"].([]interface{}); ok {
					highCount := 0
					criticalCount := 0
					for range results {
						// OSV scanner vulnerability counting logic
						// This is simplified - OSV has different structure than Grype
						highCount++ // Placeholder - would need actual severity parsing
					}
					result["summary"].(map[string]interface{})["high"] = highCount
					result["summary"].(map[string]interface{})["critical"] = criticalCount
				}
			}
		}
		
		if err != nil {
			result["status"] = "error"
			result["stderr"] = err.Error()
			result["diagnostics"] = []string{fmt.Sprintf("OSV scan failed: %v", err)}
		}

		// Return as JSON
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})
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
		_ = request.GetString("output", "") // Not used in new signature
		format := request.GetString("format", "")
		_ = request.GetString("config", "") // Not used in new signature

		// Scan image - using old signature for compatibility
		// The new ScanImage with format and licenses is used in osv_scan tool
		result, err := module.ScanImage(ctx, image, format, "")
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

		// Scan SBOM - using the old method signature for backward compatibility
		// The new ScanSBOM method with format and licenses is used in osv_scan tool
		output, err := module.ScanSBOM(ctx, sbomPath, "", "")
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