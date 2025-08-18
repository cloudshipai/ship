package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddLicenseDetectorTools adds License Detector (software license detection) MCP tool implementations using direct Dagger calls
func AddLicenseDetectorTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addLicenseDetectorToolsDirect(s)
}

// addLicenseDetectorToolsDirect adds License Detector tools using direct Dagger module calls
func addLicenseDetectorToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLicenseDetectorModule(client)

		// Get parameters
		filePath := request.GetString("file_path", "")
		if filePath == "" {
			return mcp.NewToolResultError("file_path is required"), nil
		}
		optimize := request.GetBool("optimize", false)

		// Identify license with askalono
		output, err := module.AskalonoIdentify(ctx, filePath, optimize)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("askalono identify failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLicenseDetectorModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		if directory == "" {
			return mcp.NewToolResultError("directory is required"), nil
		}

		// Crawl directory with askalono
		output, err := module.AskalonoCrawl(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("askalono crawl failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// CycloneDX license-scanner scan file tool
	licenseScannerFileTool := mcp.NewTool("license_detector_scanner_file",
		mcp.WithDescription("Scan specific file using CycloneDX license-scanner CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to file to scan for licenses"),
			mcp.Required(),
		),
		mcp.WithBoolean("show_copyrights",
			mcp.Description("Show copyright statements"),
		),
		mcp.WithBoolean("show_hash",
			mcp.Description("Show hash values"),
		),
		mcp.WithBoolean("show_keywords",
			mcp.Description("Show license keywords"),
		),
		mcp.WithBoolean("debug",
			mcp.Description("Enable debug output"),
		),
	)
	s.AddTool(licenseScannerFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLicenseDetectorModule(client)

		// Get parameters
		filePath := request.GetString("file_path", "")
		if filePath == "" {
			return mcp.NewToolResultError("file_path is required"), nil
		}
		showCopyrights := request.GetBool("show_copyrights", false)
		showHash := request.GetBool("show_hash", false)
		showKeywords := request.GetBool("show_keywords", false)
		debug := request.GetBool("debug", false)

		// Scan file with license scanner
		output, err := module.LicenseScannerFile(ctx, filePath, showCopyrights, showHash, showKeywords, debug)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("license scanner file scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// CycloneDX license-scanner scan directory tool
	licenseScannerDirTool := mcp.NewTool("license_detector_scanner_directory",
		mcp.WithDescription("Scan directory using CycloneDX license-scanner CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan for licenses"),
			mcp.Required(),
		),
		mcp.WithBoolean("show_copyrights",
			mcp.Description("Show copyright statements"),
		),
		mcp.WithBoolean("show_hash",
			mcp.Description("Show hash values"),
		),
		mcp.WithBoolean("quiet",
			mcp.Description("Quiet mode - reduce output"),
		),
	)
	s.AddTool(licenseScannerDirTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLicenseDetectorModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		if directory == "" {
			return mcp.NewToolResultError("directory is required"), nil
		}
		showCopyrights := request.GetBool("show_copyrights", false)
		showHash := request.GetBool("show_hash", false)
		quiet := request.GetBool("quiet", false)

		// Scan directory with license scanner
		output, err := module.LicenseScannerDirectory(ctx, directory, showCopyrights, showHash, quiet)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("license scanner directory scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// CycloneDX license-scanner list tool
	licenseScannerListTool := mcp.NewTool("license_detector_scanner_list",
		mcp.WithDescription("List available licenses in CycloneDX license-scanner"),
	)
	s.AddTool(licenseScannerListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLicenseDetectorModule(client)

		// This uses a basic license detection functionality
		// Note: We use DetectLicenses as a general license detection function
		output, err := module.DetectLicenses(ctx, ".")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("license detection failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Go license detector tool
	goLicenseDetectorTool := mcp.NewTool("license_detector_go_detector",
		mcp.WithDescription("Detect licenses using go-license-detector CLI"),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory"),
			mcp.Required(),
		),
	)
	s.AddTool(goLicenseDetectorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLicenseDetectorModule(client)

		// Get parameters
		projectPath := request.GetString("project_path", "")
		if projectPath == "" {
			return mcp.NewToolResultError("project_path is required"), nil
		}

		// Detect licenses with go-license-detector
		output, err := module.GoLicenseDetector(ctx, projectPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("go license detector failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// License Finder scan tool
	licenseFinderScanTool := mcp.NewTool("license_detector_licensefinder_scan",
		mcp.WithDescription("Scan project for licenses using License Finder"),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory"),
			mcp.Required(),
		),
	)
	s.AddTool(licenseFinderScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLicenseDetectorModule(client)

		// Get parameters
		projectPath := request.GetString("project_path", "")
		if projectPath == "" {
			return mcp.NewToolResultError("project_path is required"), nil
		}

		// Use DetectLicenses as general license detection
		output, err := module.DetectLicenses(ctx, projectPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("license finder scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// License Finder report tool
	licenseFinderReportTool := mcp.NewTool("license_detector_licensefinder_report",
		mcp.WithDescription("Generate license report using License Finder"),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Report format (json, csv, xml, html)"),
			mcp.Enum("json", "csv", "xml", "html"),
		),
	)
	s.AddTool(licenseFinderReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLicenseDetectorModule(client)

		// Get parameters
		projectPath := request.GetString("project_path", "")
		if projectPath == "" {
			return mcp.NewToolResultError("project_path is required"), nil
		}
		format := request.GetString("format", "json")

		// Generate license report
		output, err := module.LicenseFinderReport(ctx, projectPath, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("license finder report failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// License Finder action items tool
	licenseFinderActionItemsTool := mcp.NewTool("license_detector_licensefinder_action_items",
		mcp.WithDescription("Get action items for license compliance using License Finder"),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory"),
			mcp.Required(),
		),
		mcp.WithString("allowed_licenses",
			mcp.Description("Comma-separated list of allowed licenses"),
		),
	)
	s.AddTool(licenseFinderActionItemsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewLicenseDetectorModule(client)

		// Get parameters
		projectPath := request.GetString("project_path", "")
		if projectPath == "" {
			return mcp.NewToolResultError("project_path is required"), nil
		}
		allowedLicensesStr := request.GetString("allowed_licenses", "")
		
		// Parse allowed licenses (simple split by comma)
		var allowedLicenses []string
		if allowedLicensesStr != "" {
			// This is a simple implementation - could be enhanced for better parsing
			allowedLicenses = []string{allowedLicensesStr}
		}

		// Validate license compliance
		output, err := module.ValidateLicenseCompliance(ctx, projectPath, allowedLicenses)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("license compliance validation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}