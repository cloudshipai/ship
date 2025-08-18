package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddInfrascanTools adds Infrascan (AWS infrastructure mapping) MCP tool implementations using direct Dagger calls
func AddInfrascanTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addInfrascanToolsDirect(s)
}

// addInfrascanToolsDirect adds Infrascan tools using direct Dagger module calls
func addInfrascanToolsDirect(s *server.MCPServer) {
	// Infrascan scan tool
	scanTool := mcp.NewTool("infrascan_scan",
		mcp.WithDescription("Scan AWS infrastructure and generate system map using infrascan scan"),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of AWS regions to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_dir",
			mcp.Description("Output directory for scan results"),
			mcp.Required(),
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
		module := modules.NewInfraScanModule(client)

		// Get parameters
		regionsStr := request.GetString("regions", "")
		if regionsStr == "" {
			return mcp.NewToolResultError("regions is required"), nil
		}
		
		outputDir := request.GetString("output_dir", "")
		if outputDir == "" {
			return mcp.NewToolResultError("output_dir is required"), nil
		}

		// Parse regions
		regions := []string{}
		for _, region := range strings.Split(regionsStr, ",") {
			region = strings.TrimSpace(region)
			if region != "" {
				regions = append(regions, region)
			}
		}

		// Scan AWS infrastructure
		output, err := module.ScanAWSInfrastructure(ctx, regions, outputDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to scan infrastructure: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Infrascan graph tool
	graphTool := mcp.NewTool("infrascan_graph",
		mcp.WithDescription("Generate graph from infrascan scan results using infrascan graph"),
		mcp.WithString("input_dir",
			mcp.Description("Input directory containing scan results"),
			mcp.Required(),
		),
	)
	s.AddTool(graphTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfraScanModule(client)

		// Get input directory
		inputDir := request.GetString("input_dir", "")
		if inputDir == "" {
			return mcp.NewToolResultError("input_dir is required"), nil
		}

		// Generate graph
		output, err := module.GenerateGraph(ctx, inputDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to generate graph: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Infrascan render tool
	renderTool := mcp.NewTool("infrascan_render",
		mcp.WithDescription("Render infrastructure graph using infrascan render"),
		mcp.WithString("input_file",
			mcp.Description("Path to graph JSON file"),
			mcp.Required(),
		),
		mcp.WithBoolean("browser",
			mcp.Description("Open graph in browser"),
		),
	)
	s.AddTool(renderTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfraScanModule(client)

		// Get parameters
		inputFile := request.GetString("input_file", "")
		if inputFile == "" {
			return mcp.NewToolResultError("input_file is required"), nil
		}
		
		openBrowser := request.GetBool("browser", false)

		// Render graph
		output, err := module.RenderGraph(ctx, inputFile, openBrowser)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to render graph: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}