package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddOpenSCAPTools adds OpenSCAP (security compliance scanning) MCP tool implementations using direct Dagger calls
func AddOpenSCAPTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addOpenSCAPToolsDirect(s)
}

// addOpenSCAPToolsDirect adds OpenSCAP tools using direct Dagger module calls
func addOpenSCAPToolsDirect(s *server.MCPServer) {
	// OpenSCAP XCCDF evaluation tool
	xccdfEvalTool := mcp.NewTool("openscap_xccdf_eval",
		mcp.WithDescription("Evaluate XCCDF content for security compliance using oscap"),
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		xccdfFile := request.GetString("xccdf_file", "")
		if xccdfFile == "" {
			return mcp.NewToolResultError("xccdf_file is required"), nil
		}
		profile := request.GetString("profile", "")

		// Evaluate profile
		output, err := module.EvaluateProfile(ctx, xccdfFile, profile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap xccdf eval failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenSCAP OVAL evaluation tool
	ovalEvalTool := mcp.NewTool("openscap_oval_eval",
		mcp.WithDescription("Evaluate OVAL definitions using oscap"),
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		ovalFile := request.GetString("oval_file", "")
		if ovalFile == "" {
			return mcp.NewToolResultError("oval_file is required"), nil
		}
		resultsFile := request.GetString("results_file", "")
		variablesFile := request.GetString("variables_file", "")
		definitionId := request.GetString("definition_id", "")

		// Evaluate OVAL
		output, err := module.OvalEvaluate(ctx, ovalFile, resultsFile, variablesFile, definitionId)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap oval eval failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenSCAP generate XCCDF report tool
	xccdfGenerateReportTool := mcp.NewTool("openscap_xccdf_generate_report",
		mcp.WithDescription("Generate HTML report from XCCDF results using oscap"),
		mcp.WithString("results_file",
			mcp.Description("Path to XCCDF results XML file"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Path to save HTML report (default: stdout)"),
		),
	)
	s.AddTool(xccdfGenerateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		resultsFile := request.GetString("results_file", "")
		if resultsFile == "" {
			return mcp.NewToolResultError("results_file is required"), nil
		}

		// Generate report
		output, err := module.GenerateReport(ctx, resultsFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap xccdf generate report failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenSCAP generate XCCDF guide tool
	xccdfGenerateGuideTool := mcp.NewTool("openscap_xccdf_generate_guide",
		mcp.WithDescription("Generate HTML guide from XCCDF content using oscap"),
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		xccdfFile := request.GetString("xccdf_file", "")
		if xccdfFile == "" {
			return mcp.NewToolResultError("xccdf_file is required"), nil
		}
		profile := request.GetString("profile", "")
		outputFile := request.GetString("output_file", "")

		// Generate guide
		output, err := module.GenerateGuide(ctx, xccdfFile, profile, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap xccdf generate guide failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenSCAP DataStream validation tool
	dataStreamValidateTool := mcp.NewTool("openscap_ds_validate",
		mcp.WithDescription("Validate Source DataStream file using oscap"),
		mcp.WithString("datastream_file",
			mcp.Description("Path to Source DataStream file"),
			mcp.Required(),
		),
	)
	s.AddTool(dataStreamValidateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		datastreamFile := request.GetString("datastream_file", "")
		if datastreamFile == "" {
			return mcp.NewToolResultError("datastream_file is required"), nil
		}

		// Validate datastream
		output, err := module.ValidateDataStream(ctx, datastreamFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap ds sds-validate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenSCAP content validation tool
	validateContentTool := mcp.NewTool("openscap_validate",
		mcp.WithDescription("Validate SCAP content (XCCDF, OVAL, CPE, CVE) using oscap"),
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		contentFile := request.GetString("content_file", "")
		if contentFile == "" {
			return mcp.NewToolResultError("content_file is required"), nil
		}
		contentType := request.GetString("content_type", "")
		schematron := request.GetBool("schematron", false)

		// Validate content
		output, err := module.ValidateContent(ctx, contentFile, contentType, schematron)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap validate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenSCAP content information tool
	infoTool := mcp.NewTool("openscap_info",
		mcp.WithDescription("Display information about SCAP content using oscap"),
		mcp.WithString("content_file",
			mcp.Description("Path to SCAP content file"),
			mcp.Required(),
		),
	)
	s.AddTool(infoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		contentFile := request.GetString("content_file", "")
		if contentFile == "" {
			return mcp.NewToolResultError("content_file is required"), nil
		}

		// Get content info
		output, err := module.GetInfo(ctx, contentFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap info failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenSCAP XCCDF remediation tool
	xccdfRemediateTool := mcp.NewTool("openscap_xccdf_remediate",
		mcp.WithDescription("Apply remediation based on XCCDF results using oscap"),
		mcp.WithString("results_file",
			mcp.Description("Path to XCCDF results XML file"),
			mcp.Required(),
		),
	)
	s.AddTool(xccdfRemediateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		resultsFile := request.GetString("results_file", "")
		if resultsFile == "" {
			return mcp.NewToolResultError("results_file is required"), nil
		}

		// Apply remediation
		output, err := module.RemediateXCCDF(ctx, resultsFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap xccdf remediate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenSCAP OVAL report generation tool
	ovalGenerateReportTool := mcp.NewTool("openscap_oval_generate_report",
		mcp.WithDescription("Generate report from OVAL results using oscap"),
		mcp.WithString("oval_results_file",
			mcp.Description("Path to OVAL results XML file"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Path to save HTML report (default: stdout)"),
		),
	)
	s.AddTool(ovalGenerateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		ovalResultsFile := request.GetString("oval_results_file", "")
		if ovalResultsFile == "" {
			return mcp.NewToolResultError("oval_results_file is required"), nil
		}
		outputFile := request.GetString("output_file", "")

		// Generate OVAL report
		output, err := module.GenerateOvalReport(ctx, ovalResultsFile, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap oval generate report failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// OpenSCAP DataStream split tool
	dataStreamSplitTool := mcp.NewTool("openscap_ds_split",
		mcp.WithDescription("Split DataStream into component files using oscap"),
		mcp.WithString("datastream_file",
			mcp.Description("Path to Source DataStream file"),
			mcp.Required(),
		),
		mcp.WithString("output_dir",
			mcp.Description("Directory to save split files"),
		),
	)
	s.AddTool(dataStreamSplitTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewOpenSCAPModule(client)

		// Get parameters
		datastreamFile := request.GetString("datastream_file", "")
		if datastreamFile == "" {
			return mcp.NewToolResultError("datastream_file is required"), nil
		}
		outputDir := request.GetString("output_dir", "")

		// Split datastream
		output, err := module.SplitDataStream(ctx, datastreamFile, outputDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("oscap ds sds-split failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}