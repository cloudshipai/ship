package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddScanCodeTools adds ScanCode Toolkit MCP tool implementations
func AddScanCodeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// ScanCode licenses tool
	licensesTool := mcp.NewTool("scancode_licenses",
		mcp.WithDescription("High-accuracy license detection over a path; output pretty JSON"),
		mcp.WithString("path",
			mcp.Description("Path to scan for licenses"),
			mcp.Required(),
		),
		mcp.WithString("output_path",
			mcp.Description("Where to write the JSON output (default: ./scancode.json)"),
		),
		mcp.WithString("extra_flags",
			mcp.Description("Additional flags to pass to scancode (comma-separated)"),
		),
	)
	s.AddTool(licensesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewScanCodeModule(client)

		// Get parameters
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}
		outputPath := request.GetString("output_path", "./scancode.json")
		
		// Parse extra flags if provided
		var extraFlags []string
		if flagsStr := request.GetString("extra_flags", ""); flagsStr != "" {
			// Simple comma split for flags
			for _, flag := range splitFlags(flagsStr) {
				if flag != "" {
					extraFlags = append(extraFlags, flag)
				}
			}
		}

		// Run license scan
		stdout, err := module.LicenseScan(ctx, path, outputPath, extraFlags)
		
		// Build result
		result := map[string]interface{}{
			"status": "ok",
			"stdout": stdout,
			"stderr": "",
			"artifacts": map[string]string{
				"scancode_json": outputPath,
			},
		}
		
		if err != nil {
			result["status"] = "error"
			result["stderr"] = err.Error()
			result["diagnostics"] = []string{fmt.Sprintf("ScanCode license scan failed: %v", err)}
		}

		// Return as JSON
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})

	// ScanCode help tool
	helpTool := mcp.NewTool("scancode_help",
		mcp.WithDescription("Get ScanCode help information"),
	)
	s.AddTool(helpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewScanCodeModule(client)

		// Get help
		output, err := module.GetHelp(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scancode help failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// ScanCode version tool
	versionTool := mcp.NewTool("scancode_version",
		mcp.WithDescription("Get ScanCode version information"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewScanCodeModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scancode version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}

// splitFlags splits a comma-separated string of flags
func splitFlags(s string) []string {
	var result []string
	current := ""
	inQuote := false
	
	for _, r := range s {
		if r == '"' {
			inQuote = !inQuote
			current += string(r)
		} else if r == ',' && !inQuote {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	
	if current != "" {
		result = append(result, current)
	}
	
	return result
}