package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddOpenSCAPTools adds OpenSCAP (security compliance scanning) MCP tool implementations using real CLI commands
func AddOpenSCAPTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// OpenSCAP XCCDF evaluation tool
	xccdfEvalTool := mcp.NewTool("openscap_xccdf_eval",
		mcp.WithDescription("Evaluate XCCDF content for security compliance using real oscap CLI"),
		mcp.WithString("xccdf_file",
			mcp.Description("Path to XCCDF file or DataStream"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("Security profile ID to evaluate (e.g., ospp, stig, cis)"),
		),
		mcp.WithString("results_file",
			mcp.Description("Path to save results XML file"),
		),
		mcp.WithString("report_file",
			mcp.Description("Path to save HTML report file"),
		),
		mcp.WithString("cpe_file",
			mcp.Description("Path to CPE dictionary file"),
		),
		mcp.WithBoolean("fetch_remote_resources",
			mcp.Description("Download remote resources during evaluation"),
		),
	)
	s.AddTool(xccdfEvalTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		xccdfFile := request.GetString("xccdf_file", "")
		args := []string{"oscap", "xccdf", "eval"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if resultsFile := request.GetString("results_file", ""); resultsFile != "" {
			args = append(args, "--results", resultsFile)
		}
		if reportFile := request.GetString("report_file", ""); reportFile != "" {
			args = append(args, "--report", reportFile)
		}
		if cpeFile := request.GetString("cpe_file", ""); cpeFile != "" {
			args = append(args, "--cpe", cpeFile)
		}
		if request.GetBool("fetch_remote_resources", false) {
			args = append(args, "--fetch-remote-resources")
		}
		
		args = append(args, xccdfFile)
		return executeShipCommand(args)
	})

	// OpenSCAP OVAL evaluation tool
	ovalEvalTool := mcp.NewTool("openscap_oval_eval",
		mcp.WithDescription("Evaluate OVAL definitions using real oscap CLI"),
		mcp.WithString("oval_file",
			mcp.Description("Path to OVAL definitions file"),
			mcp.Required(),
		),
		mcp.WithString("results_file",
			mcp.Description("Path to save OVAL results XML file"),
		),
		mcp.WithString("variables_file",
			mcp.Description("Path to OVAL variables file"),
		),
		mcp.WithString("definition_id",
			mcp.Description("Specific OVAL definition ID to evaluate"),
		),
	)
	s.AddTool(ovalEvalTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ovalFile := request.GetString("oval_file", "")
		args := []string{"oscap", "oval", "eval"}
		
		if resultsFile := request.GetString("results_file", ""); resultsFile != "" {
			args = append(args, "--results", resultsFile)
		}
		if variablesFile := request.GetString("variables_file", ""); variablesFile != "" {
			args = append(args, "--variables", variablesFile)
		}
		if definitionId := request.GetString("definition_id", ""); definitionId != "" {
			args = append(args, "--id", definitionId)
		}
		
		args = append(args, ovalFile)
		return executeShipCommand(args)
	})

	// OpenSCAP generate XCCDF report tool
	xccdfGenerateReportTool := mcp.NewTool("openscap_xccdf_generate_report",
		mcp.WithDescription("Generate HTML report from XCCDF results using real oscap CLI"),
		mcp.WithString("results_file",
			mcp.Description("Path to XCCDF results XML file"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Path to save HTML report (default: stdout)"),
		),
	)
	s.AddTool(xccdfGenerateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resultsFile := request.GetString("results_file", "")
		args := []string{"oscap", "xccdf", "generate", "report", resultsFile}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			// Redirect output to file
			args = []string{"sh", "-c", "oscap xccdf generate report " + resultsFile + " > " + outputFile}
		}
		
		return executeShipCommand(args)
	})

	// OpenSCAP generate XCCDF guide tool
	xccdfGenerateGuideTool := mcp.NewTool("openscap_xccdf_generate_guide",
		mcp.WithDescription("Generate HTML guide from XCCDF content using real oscap CLI"),
		mcp.WithString("xccdf_file",
			mcp.Description("Path to XCCDF file or DataStream"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("Security profile ID for guide generation"),
		),
		mcp.WithString("output_file",
			mcp.Description("Path to save HTML guide (default: stdout)"),
		),
	)
	s.AddTool(xccdfGenerateGuideTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		xccdfFile := request.GetString("xccdf_file", "")
		args := []string{"oscap", "xccdf", "generate", "guide"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		
		args = append(args, xccdfFile)
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			// Redirect output to file
			command := "oscap xccdf generate guide"
			if profile := request.GetString("profile", ""); profile != "" {
				command += " --profile " + profile
			}
			command += " " + xccdfFile + " > " + outputFile
			args = []string{"sh", "-c", command}
		}
		
		return executeShipCommand(args)
	})

	// OpenSCAP DataStream validation tool
	dataStreamValidateTool := mcp.NewTool("openscap_ds_validate",
		mcp.WithDescription("Validate Source DataStream file using real oscap CLI"),
		mcp.WithString("datastream_file",
			mcp.Description("Path to Source DataStream file"),
			mcp.Required(),
		),
	)
	s.AddTool(dataStreamValidateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		datastreamFile := request.GetString("datastream_file", "")
		args := []string{"oscap", "ds", "sds-validate", datastreamFile}
		return executeShipCommand(args)
	})

	// OpenSCAP content validation tool
	validateContentTool := mcp.NewTool("openscap_validate",
		mcp.WithDescription("Validate SCAP content (XCCDF, OVAL, CPE, CVE) using real oscap CLI"),
		mcp.WithString("content_file",
			mcp.Description("Path to SCAP content file"),
			mcp.Required(),
		),
		mcp.WithString("content_type",
			mcp.Description("Content type to validate"),
			mcp.Enum("xccdf", "oval", "cpe", "cve"),
		),
		mcp.WithBoolean("schematron",
			mcp.Description("Use Schematron validation for OVAL (more thorough)"),
		),
	)
	s.AddTool(validateContentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		contentFile := request.GetString("content_file", "")
		contentType := request.GetString("content_type", "")
		
		var args []string
		if contentType != "" {
			args = []string{"oscap", contentType, "validate"}
			if contentType == "oval" && request.GetBool("schematron", false) {
				args = append(args, "--schematron")
			}
		} else {
			// Auto-detect validation type
			args = []string{"oscap", "info", contentFile}
		}
		
		args = append(args, contentFile)
		return executeShipCommand(args)
	})

	// OpenSCAP content information tool
	infoTool := mcp.NewTool("openscap_info",
		mcp.WithDescription("Display information about SCAP content using real oscap CLI"),
		mcp.WithString("content_file",
			mcp.Description("Path to SCAP content file"),
			mcp.Required(),
		),
	)
	s.AddTool(infoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		contentFile := request.GetString("content_file", "")
		args := []string{"oscap", "info", contentFile}
		return executeShipCommand(args)
	})

	// OpenSCAP XCCDF remediation tool
	xccdfRemediateTool := mcp.NewTool("openscap_xccdf_remediate",
		mcp.WithDescription("Apply remediation based on XCCDF results using real oscap CLI"),
		mcp.WithString("results_file",
			mcp.Description("Path to XCCDF results XML file"),
			mcp.Required(),
		),
	)
	s.AddTool(xccdfRemediateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resultsFile := request.GetString("results_file", "")
		args := []string{"oscap", "xccdf", "remediate", resultsFile}
		return executeShipCommand(args)
	})

	// OpenSCAP OVAL report generation tool
	ovalGenerateReportTool := mcp.NewTool("openscap_oval_generate_report",
		mcp.WithDescription("Generate report from OVAL results using real oscap CLI"),
		mcp.WithString("oval_results_file",
			mcp.Description("Path to OVAL results XML file"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Path to save HTML report (default: stdout)"),
		),
	)
	s.AddTool(ovalGenerateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ovalResultsFile := request.GetString("oval_results_file", "")
		args := []string{"oscap", "oval", "generate", "report", ovalResultsFile}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			// Redirect output to file
			args = []string{"sh", "-c", "oscap oval generate report " + ovalResultsFile + " > " + outputFile}
		}
		
		return executeShipCommand(args)
	})

	// OpenSCAP DataStream split tool
	dataStreamSplitTool := mcp.NewTool("openscap_ds_split",
		mcp.WithDescription("Split DataStream into component files using real oscap CLI"),
		mcp.WithString("datastream_file",
			mcp.Description("Path to Source DataStream file"),
			mcp.Required(),
		),
		mcp.WithString("output_dir",
			mcp.Description("Directory to save split files"),
		),
	)
	s.AddTool(dataStreamSplitTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		datastreamFile := request.GetString("datastream_file", "")
		args := []string{"oscap", "ds", "sds-split"}
		
		if outputDir := request.GetString("output_dir", ""); outputDir != "" {
			args = append(args, "--output-dir", outputDir)
		}
		
		args = append(args, datastreamFile)
		return executeShipCommand(args)
	})
}