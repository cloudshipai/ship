package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddLicenseDetectorTools adds License Detector (software license detection) MCP tool implementations
func AddLicenseDetectorTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// License Detector scan directory tool
	scanDirectoryTool := mcp.NewTool("license_detector_scan_directory",
		mcp.WithDescription("Scan directory for software licenses"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively scan subdirectories"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (json, csv, table)"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "license-detector", "scan"}
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

	// License Detector scan file tool
	scanFileTool := mcp.NewTool("license_detector_scan_file",
		mcp.WithDescription("Scan specific file for license information"),
		mcp.WithString("file_path",
			mcp.Description("Path to file to scan"),
			mcp.Required(),
		),
		mcp.WithString("confidence_threshold",
			mcp.Description("Minimum confidence threshold (0.0-1.0)"),
		),
	)
	s.AddTool(scanFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"security", "license-detector", "scan-file", filePath}
		if threshold := request.GetString("confidence_threshold", ""); threshold != "" {
			args = append(args, "--threshold", threshold)
		}
		return executeShipCommand(args)
	})

	// License Detector check compliance tool
	checkComplianceTool := mcp.NewTool("license_detector_check_compliance",
		mcp.WithDescription("Check license compliance against policy"),
		mcp.WithString("directory",
			mcp.Description("Directory to check (default: current directory)"),
		),
		mcp.WithString("policy_file",
			mcp.Description("Path to license policy file"),
			mcp.Required(),
		),
		mcp.WithBoolean("fail_on_violation",
			mcp.Description("Fail if policy violations are found"),
		),
	)
	s.AddTool(checkComplianceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyFile := request.GetString("policy_file", "")
		args := []string{"security", "license-detector", "check", "--policy", policyFile}
		if directory := request.GetString("directory", ""); directory != "" {
			args = append(args, directory)
		}
		if request.GetBool("fail_on_violation", false) {
			args = append(args, "--fail-on-violation")
		}
		return executeShipCommand(args)
	})

	// License Detector generate report tool
	generateReportTool := mcp.NewTool("license_detector_generate_report",
		mcp.WithDescription("Generate comprehensive license report"),
		mcp.WithString("directory",
			mcp.Description("Directory to analyze (default: current directory)"),
		),
		mcp.WithString("report_format",
			mcp.Description("Report format (html, pdf, json, csv)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for the report"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "license-detector", "report"}
		if directory := request.GetString("directory", ""); directory != "" {
			args = append(args, directory)
		}
		if reportFormat := request.GetString("report_format", ""); reportFormat != "" {
			args = append(args, "--format", reportFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// License Detector list supported licenses tool
	listLicensesTool := mcp.NewTool("license_detector_list_licenses",
		mcp.WithDescription("List all supported license types"),
		mcp.WithBoolean("show_details",
			mcp.Description("Show detailed license information"),
		),
	)
	s.AddTool(listLicensesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "license-detector", "list"}
		if request.GetBool("show_details", false) {
			args = append(args, "--details")
		}
		return executeShipCommand(args)
	})

	// License Detector validate license tool
	validateLicenseTool := mcp.NewTool("license_detector_validate_license",
		mcp.WithDescription("Validate specific license text"),
		mcp.WithString("license_text",
			mcp.Description("License text to validate"),
			mcp.Required(),
		),
		mcp.WithString("expected_license",
			mcp.Description("Expected license type for comparison"),
		),
	)
	s.AddTool(validateLicenseTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		licenseText := request.GetString("license_text", "")
		args := []string{"security", "license-detector", "validate", "--text", licenseText}
		if expectedLicense := request.GetString("expected_license", ""); expectedLicense != "" {
			args = append(args, "--expected", expectedLicense)
		}
		return executeShipCommand(args)
	})
}