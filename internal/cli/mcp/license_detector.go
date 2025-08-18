package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddLicenseDetectorTools adds License Detector (software license detection) MCP tool implementations using real CLI commands
func AddLicenseDetectorTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Askalono license detection tool
	askalonoIdentifyTool := mcp.NewTool("license_detector_askalono_identify",
		mcp.WithDescription("Identify license in a file using askalono CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to license file to identify"),
			mcp.Required(),
		),
		mcp.WithBoolean("optimize",
			mcp.Description("Optimize detection for files with headers/footers"),
		),
	)
	s.AddTool(askalonoIdentifyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"askalono", "id", filePath}
		
		if request.GetBool("optimize", false) {
			args = []string{"askalono", "id", "--optimize", filePath}
		}
		
		return executeShipCommand(args)
	})

	// Askalono crawl directory tool
	askalonoCrawlTool := mcp.NewTool("license_detector_askalono_crawl",
		mcp.WithDescription("Crawl directory for license files using askalono CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory to crawl for license files"),
			mcp.Required(),
		),
	)
	s.AddTool(askalonoCrawlTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"askalono", "crawl", directory}
		return executeShipCommand(args)
	})

	// CycloneDX license-scanner scan file tool
	licenseScannerFileTool := mcp.NewTool("license_detector_scanner_file",
		mcp.WithDescription("Scan specific file using CycloneDX license-scanner CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to file to scan"),
			mcp.Required(),
		),
		mcp.WithBoolean("show_copyrights",
			mcp.Description("Show copyright information"),
		),
		mcp.WithBoolean("show_hash",
			mcp.Description("Output file hash"),
		),
		mcp.WithBoolean("show_keywords",
			mcp.Description("Flag keywords"),
		),
		mcp.WithBoolean("debug",
			mcp.Description("Enable debug logging"),
		),
	)
	s.AddTool(licenseScannerFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"license-scanner", "--file", filePath}
		
		if request.GetBool("show_copyrights", false) {
			args = append(args, "--copyrights")
		}
		if request.GetBool("show_hash", false) {
			args = append(args, "--hash")
		}
		if request.GetBool("show_keywords", false) {
			args = append(args, "--keywords")
		}
		if request.GetBool("debug", false) {
			args = append(args, "--debug")
		}
		
		return executeShipCommand(args)
	})

	// CycloneDX license-scanner scan directory tool
	licenseScannerDirTool := mcp.NewTool("license_detector_scanner_directory",
		mcp.WithDescription("Scan directory using CycloneDX license-scanner CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithBoolean("show_copyrights",
			mcp.Description("Show copyright information"),
		),
		mcp.WithBoolean("show_hash",
			mcp.Description("Output file hash"),
		),
		mcp.WithBoolean("quiet",
			mcp.Description("Suppress logging"),
		),
	)
	s.AddTool(licenseScannerDirTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"license-scanner", "--dir", directory}
		
		if request.GetBool("show_copyrights", false) {
			args = append(args, "--copyrights")
		}
		if request.GetBool("show_hash", false) {
			args = append(args, "--hash")
		}
		if request.GetBool("quiet", false) {
			args = append(args, "--quiet")
		}
		
		return executeShipCommand(args)
	})

	// CycloneDX license-scanner list templates tool
	licenseScannerListTool := mcp.NewTool("license_detector_scanner_list",
		mcp.WithDescription("List license templates using CycloneDX license-scanner CLI"),
	)
	s.AddTool(licenseScannerListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"license-scanner", "--list"}
		return executeShipCommand(args)
	})

	// Go license detector tool
	goLicenseDetectorTool := mcp.NewTool("license_detector_go_detector",
		mcp.WithDescription("Detect project license using go-license-detector CLI"),
		mcp.WithString("path",
			mcp.Description("Path to project directory or GitHub repository URL"),
			mcp.Required(),
		),
	)
	s.AddTool(goLicenseDetectorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := request.GetString("path", "")
		args := []string{"license-detector", path}
		return executeShipCommand(args)
	})

	// LicenseFinder scan dependencies tool
	licenseFinderScanTool := mcp.NewTool("license_detector_licensefinder_scan",
		mcp.WithDescription("Scan project dependencies using LicenseFinder CLI"),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory (default: current directory)"),
		),
		mcp.WithString("decisions_file",
			mcp.Description("Path to decisions file"),
		),
	)
	s.AddTool(licenseFinderScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"license_finder"}
		
		if projectPath := request.GetString("project_path", ""); projectPath != "" {
			// Change to project directory
			args = []string{"sh", "-c", "cd " + projectPath + " && license_finder"}
		}
		if decisionsFile := request.GetString("decisions_file", ""); decisionsFile != "" {
			args = append(args, "--decisions_file", decisionsFile)
		}
		
		return executeShipCommand(args)
	})

	// LicenseFinder generate report tool
	licenseFinderReportTool := mcp.NewTool("license_detector_licensefinder_report",
		mcp.WithDescription("Generate license report using LicenseFinder CLI"),
		mcp.WithString("format",
			mcp.Description("Report format (text, csv, html, markdown)"),
			mcp.Enum("text", "csv", "html", "markdown"),
		),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory (default: current directory)"),
		),
	)
	s.AddTool(licenseFinderReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"license_finder", "report"}
		
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		
		if projectPath := request.GetString("project_path", ""); projectPath != "" {
			// Change to project directory
			fullCommand := "cd " + projectPath + " && " + "license_finder report"
			if format := request.GetString("format", ""); format != "" {
				fullCommand += " --format " + format
			}
			args = []string{"sh", "-c", fullCommand}
		}
		
		return executeShipCommand(args)
	})

	// LicenseFinder action items tool
	licenseFinderActionItemsTool := mcp.NewTool("license_detector_licensefinder_action_items",
		mcp.WithDescription("Show dependencies needing approval using LicenseFinder CLI"),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory (default: current directory)"),
		),
	)
	s.AddTool(licenseFinderActionItemsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"license_finder", "action_items"}
		
		if projectPath := request.GetString("project_path", ""); projectPath != "" {
			args = []string{"sh", "-c", "cd " + projectPath + " && license_finder action_items"}
		}
		
		return executeShipCommand(args)
	})
}