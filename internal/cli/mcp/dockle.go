package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddDockleTools adds Dockle (container image linter) MCP tool implementations
func AddDockleTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Dockle scan image tool
	scanImageTool := mcp.NewTool("dockle_scan_image",
		mcp.WithDescription("Scan container image for security and best practices using Dockle"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"dockle", imageRef}
		return executeShipCommand(args)
	})

	// Dockle scan tarball tool
	scanTarballTool := mcp.NewTool("dockle_scan_tarball",
		mcp.WithDescription("Scan container image tarball using Dockle"),
		mcp.WithString("tarball_path",
			mcp.Description("Path to container image tarball"),
			mcp.Required(),
		),
	)
	s.AddTool(scanTarballTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tarballPath := request.GetString("tarball_path", "")
		args := []string{"dockle", "--input", tarballPath}
		return executeShipCommand(args)
	})

	// Dockle scan with JSON output
	scanJsonTool := mcp.NewTool("dockle_scan_json",
		mcp.WithDescription("Scan container image and output results in JSON format"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for JSON results"),
		),
	)
	s.AddTool(scanJsonTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"dockle", "-f", "json"}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		args = append(args, imageRef)
		return executeShipCommand(args)
	})

	// Dockle scan image from tarball with JSON output
	scanTarballJsonTool := mcp.NewTool("dockle_scan_tarball_json",
		mcp.WithDescription("Scan container image tarball and output results in JSON format"),
		mcp.WithString("tarball_path",
			mcp.Description("Path to container image tarball"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for JSON results"),
		),
	)
	s.AddTool(scanTarballJsonTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tarballPath := request.GetString("tarball_path", "")
		args := []string{"dockle", "-f", "json", "--input", tarballPath}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		return executeShipCommand(args)
	})
}