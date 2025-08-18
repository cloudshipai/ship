package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddInfrascanTools adds Infrascan (AWS infrastructure mapping) MCP tool implementations using real CLI commands
func AddInfrascanTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		regions := request.GetString("regions", "")
		outputDir := request.GetString("output_dir", "")
		
		args := []string{"infrascan", "scan", "-o", outputDir}
		
		// Add regions
		for _, region := range strings.Split(regions, ",") {
			region = strings.TrimSpace(region)
			if region != "" {
				args = append(args, "--region", region)
			}
		}
		
		return executeShipCommand(args)
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
		inputDir := request.GetString("input_dir", "")
		args := []string{"infrascan", "graph", "-i", inputDir}
		return executeShipCommand(args)
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
		inputFile := request.GetString("input_file", "")
		args := []string{"infrascan", "render", "-i", inputFile}
		
		if request.GetBool("browser", false) {
			args = append(args, "--browser")
		}
		
		return executeShipCommand(args)
	})


}