package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddOSVScannerTools adds OSV Scanner (Open Source Vulnerability scanner) MCP tool implementations using real CLI commands
func AddOSVScannerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// OSV Scanner scan source directory tool
	scanSourceTool := mcp.NewTool("osv_scanner_scan_source",
		mcp.WithDescription("Scan source directory for open source vulnerabilities using real osv-scanner CLI"),
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
		args := []string{"osv-scanner", "scan", "source"}
		
		if request.GetBool("recursive", false) {
			args = append(args, "-r")
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("licenses", false) {
			args = append(args, "--licenses")
		}
		
		directory := request.GetString("directory", ".")
		args = append(args, directory)
		
		return executeShipCommand(args)
	})

	// OSV Scanner scan container image tool
	scanImageTool := mcp.NewTool("osv_scanner_scan_image",
		mcp.WithDescription("Scan container image for vulnerabilities using real osv-scanner CLI"),
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
		image := request.GetString("image", "")
		args := []string{"osv-scanner", "scan", "image", image}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		
		return executeShipCommand(args)
	})

	// OSV Scanner scan lockfile tool
	scanLockfileTool := mcp.NewTool("osv_scanner_scan_lockfile",
		mcp.WithDescription("Scan specific lockfile for vulnerabilities using real osv-scanner CLI"),
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
		lockfilePath := request.GetString("lockfile_path", "")
		args := []string{"osv-scanner", "-L", lockfilePath}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if request.GetBool("all_packages", false) {
			args = append(args, "--all-packages")
		}
		
		return executeShipCommand(args)
	})

	// OSV Scanner scan manifest tool
	scanManifestTool := mcp.NewTool("osv_scanner_scan_manifest",
		mcp.WithDescription("Scan package manifest file for vulnerabilities using real osv-scanner CLI"),
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
		manifestPath := request.GetString("manifest_path", "")
		args := []string{"osv-scanner", "-M", manifestPath}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		
		return executeShipCommand(args)
	})

	// OSV Scanner license scanning tool
	licenseScanTool := mcp.NewTool("osv_scanner_license_scan",
		mcp.WithDescription("Scan for license compliance using real osv-scanner CLI"),
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
		path := request.GetString("path", "")
		args := []string{"osv-scanner", "--licenses"}
		
		if allowedLicenses := request.GetString("allowed_licenses", ""); allowedLicenses != "" {
			args = []string{"osv-scanner", "--licenses=" + allowedLicenses}
		}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		
		args = append(args, path)
		return executeShipCommand(args)
	})

	// OSV Scanner offline scanning tool
	offlineScanTool := mcp.NewTool("osv_scanner_offline_scan",
		mcp.WithDescription("Scan using offline vulnerability databases with real osv-scanner CLI"),
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
		path := request.GetString("path", "")
		args := []string{"osv-scanner", "--offline"}
		
		if request.GetBool("download_databases", false) {
			args = append(args, "--download-offline-databases")
		}
		if offlineDbPath := request.GetString("offline_db_path", ""); offlineDbPath != "" {
			args = append(args, "--offline-vulnerabilities", offlineDbPath)
		}
		if request.GetBool("recursive", false) {
			args = append(args, "-r")
		}
		
		args = append(args, path)
		return executeShipCommand(args)
	})

	// OSV Scanner guided remediation tool
	fixTool := mcp.NewTool("osv_scanner_fix",
		mcp.WithDescription("Apply guided remediation for vulnerabilities using real osv-scanner CLI"),
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
		args := []string{"osv-scanner", "fix"}
		
		if manifestPath := request.GetString("manifest_path", ""); manifestPath != "" {
			args = append(args, "-M", manifestPath)
		}
		if lockfilePath := request.GetString("lockfile_path", ""); lockfilePath != "" {
			args = append(args, "-L", lockfilePath)
		}
		if strategy := request.GetString("strategy", ""); strategy != "" {
			args = append(args, "--strategy", strategy)
		}
		if maxDepth := request.GetString("max_depth", ""); maxDepth != "" {
			args = append(args, "--max-depth", maxDepth)
		}
		if minSeverity := request.GetString("min_severity", ""); minSeverity != "" {
			args = append(args, "--min-severity", minSeverity)
		}
		if request.GetBool("ignore_dev", false) {
			args = append(args, "--ignore-dev")
		}
		
		return executeShipCommand(args)
	})

	// OSV Scanner serve report tool
	serveReportTool := mcp.NewTool("osv_scanner_serve_report",
		mcp.WithDescription("Generate and serve HTML vulnerability report locally using real osv-scanner CLI"),
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
		path := request.GetString("path", "")
		args := []string{"osv-scanner", "--serve"}
		
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "--port", port)
		}
		if request.GetBool("recursive", false) {
			args = append(args, "-r")
		}
		
		args = append(args, path)
		return executeShipCommand(args)
	})

	// OSV Scanner SBOM scan tool
	scanSBOMTool := mcp.NewTool("osv_scanner_scan_sbom",
		mcp.WithDescription("Scan SBOM file for vulnerabilities using real osv-scanner CLI"),
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
		sbomPath := request.GetString("sbom_path", "")
		args := []string{"osv-scanner", "--sbom", sbomPath}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		
		return executeShipCommand(args)
	})

	// OSV Scanner verbosity control tool
	verboseScanTool := mcp.NewTool("osv_scanner_verbose_scan",
		mcp.WithDescription("Run OSV Scanner with verbose logging using real osv-scanner CLI"),
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
		path := request.GetString("path", "")
		args := []string{"osv-scanner"}
		
		if verbosity := request.GetString("verbosity", ""); verbosity != "" {
			args = append(args, "--verbosity", verbosity)
		}
		if request.GetBool("recursive", false) {
			args = append(args, "-r")
		}
		
		args = append(args, path)
		return executeShipCommand(args)
	})
}