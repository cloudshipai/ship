package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddDockleTools adds Dockle (container image linter) MCP tool implementations
func AddDockleTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addDockleToolsDirect(s)
}

// addDockleToolsDirect implements direct Dagger calls for Dockle tools
func addDockleToolsDirect(s *server.MCPServer) {
	// Dockle scan image tool
	scanImageTool := mcp.NewTool("dockle_scan_image",
		mcp.WithDescription("Scan container image for security and best practices using Dockle"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageRef := request.GetString("image_ref", "")

		// Create Dockle module and scan image
		dockleModule := modules.NewDockleModule(client)
		result, err := dockleModule.ScanImageString(ctx, imageRef)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("dockle scan image failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		tarballPath := request.GetString("tarball_path", "")

		// Create Dockle module and scan tarball
		dockleModule := modules.NewDockleModule(client)
		result, err := dockleModule.ScanTarballString(ctx, tarballPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("dockle scan tarball failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageRef := request.GetString("image_ref", "")
		outputFile := request.GetString("output_file", "")

		// Create Dockle module and scan image with JSON output
		dockleModule := modules.NewDockleModule(client)
		result, err := dockleModule.ScanImageJSON(ctx, imageRef, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("dockle scan json failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		tarballPath := request.GetString("tarball_path", "")
		outputFile := request.GetString("output_file", "")

		// Create Dockle module and scan tarball with JSON output
		dockleModule := modules.NewDockleModule(client)
		result, err := dockleModule.ScanTarballJSON(ctx, tarballPath, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("dockle scan tarball json failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}