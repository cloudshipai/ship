package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddOSVScannerTools adds OSV Scanner (Open Source Vulnerability scanner) MCP tool implementations
func AddOSVScannerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// OSV Scanner scan directory tool
	scanDirectoryTool := mcp.NewTool("osv_scanner_scan_directory",
		mcp.WithDescription("Scan directory for open source vulnerabilities using OSV Scanner"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively scan subdirectories"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (json, table, sarif)"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "osv-scanner", "scan"}
		if directory := request.GetString("directory", ""); directory != "" {
			args = append(args, directory)
		}
		if request.GetBool("recursive", false) {
			args = append(args, "--recursive")
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		return executeShipCommand(args)
	})

	// OSV Scanner scan lockfile tool
	scanLockfileTool := mcp.NewTool("osv_scanner_scan_lockfile",
		mcp.WithDescription("Scan specific lockfile for vulnerabilities"),
		mcp.WithString("lockfile_path",
			mcp.Description("Path to lockfile (package-lock.json, requirements.txt, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("ecosystem",
			mcp.Description("Package ecosystem (npm, pip, maven, etc.)"),
		),
	)
	s.AddTool(scanLockfileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lockfilePath := request.GetString("lockfile_path", "")
		args := []string{"security", "osv-scanner", "scan-lockfile", lockfilePath}
		if ecosystem := request.GetString("ecosystem", ""); ecosystem != "" {
			args = append(args, "--ecosystem", ecosystem)
		}
		return executeShipCommand(args)
	})

	// OSV Scanner scan SBOM tool
	scanSBOMTool := mcp.NewTool("osv_scanner_scan_sbom",
		mcp.WithDescription("Scan SBOM file for vulnerabilities"),
		mcp.WithString("sbom_path",
			mcp.Description("Path to SBOM file"),
			mcp.Required(),
		),
		mcp.WithString("sbom_format",
			mcp.Description("SBOM format (cyclonedx, spdx)"),
		),
	)
	s.AddTool(scanSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sbomPath := request.GetString("sbom_path", "")
		args := []string{"security", "osv-scanner", "scan-sbom", sbomPath}
		if sbomFormat := request.GetString("sbom_format", ""); sbomFormat != "" {
			args = append(args, "--format", sbomFormat)
		}
		return executeShipCommand(args)
	})

	// OSV Scanner scan commit tool
	scanCommitTool := mcp.NewTool("osv_scanner_scan_commit",
		mcp.WithDescription("Scan git commit for vulnerabilities"),
		mcp.WithString("commit_hash",
			mcp.Description("Git commit hash to scan"),
			mcp.Required(),
		),
		mcp.WithString("repository_path",
			mcp.Description("Path to git repository (default: current directory)"),
		),
	)
	s.AddTool(scanCommitTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		commitHash := request.GetString("commit_hash", "")
		args := []string{"security", "osv-scanner", "scan-commit", commitHash}
		if repoPath := request.GetString("repository_path", ""); repoPath != "" {
			args = append(args, "--repo", repoPath)
		}
		return executeShipCommand(args)
	})

	// OSV Scanner generate report tool
	generateReportTool := mcp.NewTool("osv_scanner_generate_report",
		mcp.WithDescription("Generate comprehensive vulnerability report"),
		mcp.WithString("scan_target",
			mcp.Description("Target to scan (directory, file, or SBOM)"),
			mcp.Required(),
		),
		mcp.WithString("report_format",
			mcp.Description("Report format (html, json, sarif, csv)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for the report"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scanTarget := request.GetString("scan_target", "")
		args := []string{"security", "osv-scanner", "report", scanTarget}
		if reportFormat := request.GetString("report_format", ""); reportFormat != "" {
			args = append(args, "--format", reportFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// OSV Scanner check package tool
	checkPackageTool := mcp.NewTool("osv_scanner_check_package",
		mcp.WithDescription("Check specific package for vulnerabilities"),
		mcp.WithString("package_name",
			mcp.Description("Name of the package to check"),
			mcp.Required(),
		),
		mcp.WithString("package_version",
			mcp.Description("Version of the package"),
			mcp.Required(),
		),
		mcp.WithString("ecosystem",
			mcp.Description("Package ecosystem (npm, pip, maven, etc.)"),
			mcp.Required(),
		),
	)
	s.AddTool(checkPackageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName := request.GetString("package_name", "")
		packageVersion := request.GetString("package_version", "")
		ecosystem := request.GetString("ecosystem", "")
		args := []string{"security", "osv-scanner", "check", "--package", packageName, "--version", packageVersion, "--ecosystem", ecosystem}
		return executeShipCommand(args)
	})

	// OSV Scanner get version tool
	getVersionTool := mcp.NewTool("osv_scanner_get_version",
		mcp.WithDescription("Get OSV Scanner version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "osv-scanner", "--version"}
		return executeShipCommand(args)
	})
}