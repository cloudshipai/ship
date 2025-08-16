package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTfstateReaderTools adds Terraform state reader MCP tool implementations
func AddTfstateReaderTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Tfstate read resources tool
	readResourcesTool := mcp.NewTool("tfstate_read_resources",
		mcp.WithDescription("Read and list resources from Terraform state file"),
		mcp.WithString("state_file",
			mcp.Description("Path to Terraform state file"),
			mcp.Required(),
		),
		mcp.WithString("resource_type",
			mcp.Description("Filter by resource type (optional)"),
		),
	)
	s.AddTool(readResourcesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stateFile := request.GetString("state_file", "")
		args := []string{"terraform", "tfstate-reader", "resources", stateFile}
		if resourceType := request.GetString("resource_type", ""); resourceType != "" {
			args = append(args, "--type", resourceType)
		}
		return executeShipCommand(args)
	})

	// Tfstate analyze dependencies tool
	analyzeDependenciesTool := mcp.NewTool("tfstate_analyze_dependencies",
		mcp.WithDescription("Analyze resource dependencies in Terraform state"),
		mcp.WithString("state_file",
			mcp.Description("Path to Terraform state file"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (graph, json, table)"),
		),
	)
	s.AddTool(analyzeDependenciesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stateFile := request.GetString("state_file", "")
		args := []string{"terraform", "tfstate-reader", "dependencies", stateFile}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Tfstate extract outputs tool
	extractOutputsTool := mcp.NewTool("tfstate_extract_outputs",
		mcp.WithDescription("Extract output values from Terraform state"),
		mcp.WithString("state_file",
			mcp.Description("Path to Terraform state file"),
			mcp.Required(),
		),
		mcp.WithString("output_name",
			mcp.Description("Specific output name to extract (optional)"),
		),
	)
	s.AddTool(extractOutputsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stateFile := request.GetString("state_file", "")
		args := []string{"terraform", "tfstate-reader", "outputs", stateFile}
		if outputName := request.GetString("output_name", ""); outputName != "" {
			args = append(args, "--output", outputName)
		}
		return executeShipCommand(args)
	})

	// Tfstate compare states tool
	compareStatesTool := mcp.NewTool("tfstate_compare_states",
		mcp.WithDescription("Compare two Terraform state files for differences"),
		mcp.WithString("state_file_1",
			mcp.Description("First Terraform state file"),
			mcp.Required(),
		),
		mcp.WithString("state_file_2",
			mcp.Description("Second Terraform state file"),
			mcp.Required(),
		),
	)
	s.AddTool(compareStatesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stateFile1 := request.GetString("state_file_1", "")
		stateFile2 := request.GetString("state_file_2", "")
		args := []string{"terraform", "tfstate-reader", "compare", stateFile1, stateFile2}
		return executeShipCommand(args)
	})

	// Tfstate export tool
	exportTool := mcp.NewTool("tfstate_export",
		mcp.WithDescription("Export Terraform state to different formats"),
		mcp.WithString("state_file",
			mcp.Description("Path to Terraform state file"),
			mcp.Required(),
		),
		mcp.WithString("export_format",
			mcp.Description("Export format (json, yaml, csv)"),
			mcp.Required(),
		),
	)
	s.AddTool(exportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stateFile := request.GetString("state_file", "")
		exportFormat := request.GetString("export_format", "")
		args := []string{"terraform", "tfstate-reader", "export", stateFile, "--format", exportFormat}
		return executeShipCommand(args)
	})

	// Tfstate reader get version tool
	getVersionTool := mcp.NewTool("tfstate_reader_get_version",
		mcp.WithDescription("Get Terraform state reader version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform", "tfstate-reader", "--version"}
		return executeShipCommand(args)
	})
}