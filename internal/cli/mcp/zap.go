package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddZapTools adds OWASP ZAP (web application security scanner) MCP tool implementations using direct Dagger calls
func AddZapTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addZapToolsDirect(s)
}

// addZapToolsDirect adds ZAP tools using direct Dagger module calls
func addZapToolsDirect(s *server.MCPServer) {
	// ZAP passive scan tool (using baseline scan)
	passiveScanTool := mcp.NewTool("zap_passive_scan",
		mcp.WithDescription("Perform passive security scan using OWASP ZAP"),
		mcp.WithString("target",
			mcp.Description("Target URL to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "xml", "html"),
		),
	)
	s.AddTool(passiveScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewZapModule(client)

		// Get parameters
		target := request.GetString("target", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		outputFormat := request.GetString("output_format", "")
		if outputFormat != "" {
			return mcp.NewToolResultError("Warning: output_format parameter not supported with direct Dagger calls"), nil
		}

		// Perform baseline scan (passive)
		output, err := module.BaselineScan(ctx, target)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("ZAP passive scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// ZAP active scan tool
	activeScanTool := mcp.NewTool("zap_active_scan",
		mcp.WithDescription("Perform active security scan using OWASP ZAP"),
		mcp.WithString("target",
			mcp.Description("Target URL to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "xml", "html"),
		),
	)
	s.AddTool(activeScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewZapModule(client)

		// Get parameters
		target := request.GetString("target", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		outputFormat := request.GetString("output_format", "")
		if outputFormat != "" {
			return mcp.NewToolResultError("Warning: output_format parameter not supported with direct Dagger calls"), nil
		}

		// Perform full scan (active)
		output, err := module.FullScan(ctx, target, 60) // Default 60 minutes
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("ZAP active scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// ZAP spider scan tool
	spiderScanTool := mcp.NewTool("zap_spider_scan",
		mcp.WithDescription("Perform spider crawl and scan using OWASP ZAP"),
		mcp.WithString("target",
			mcp.Description("Target URL to spider"),
			mcp.Required(),
		),
		mcp.WithNumber("max_depth",
			mcp.Description("Maximum spider depth"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "xml", "html"),
		),
	)
	s.AddTool(spiderScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewZapModule(client)

		// Get parameters
		target := request.GetString("target", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		maxDepth := request.GetInt("max_depth", 0)
		outputFormat := request.GetString("output_format", "")

		// Perform spider scan
		output, err := module.SpiderScan(ctx, target, maxDepth, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("ZAP spider scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// ZAP full scan tool
	fullScanTool := mcp.NewTool("zap_full_scan",
		mcp.WithDescription("Perform comprehensive security scan using OWASP ZAP"),
		mcp.WithString("target",
			mcp.Description("Target URL to scan"),
			mcp.Required(),
		),
		mcp.WithNumber("max_duration",
			mcp.Description("Maximum scan duration in minutes"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "xml", "html"),
		),
	)
	s.AddTool(fullScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewZapModule(client)

		// Get parameters
		target := request.GetString("target", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		maxDuration := request.GetInt("max_duration", 60) // Default 60 minutes
		outputFormat := request.GetString("output_format", "")
		if outputFormat != "" {
			return mcp.NewToolResultError("Warning: output_format parameter not supported with direct Dagger calls"), nil
		}

		// Perform full scan
		output, err := module.FullScan(ctx, target, maxDuration)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("ZAP full scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// ZAP baseline scan tool
	baselineScanTool := mcp.NewTool("zap_baseline_scan",
		mcp.WithDescription("Perform baseline security scan using OWASP ZAP"),
		mcp.WithString("target",
			mcp.Description("Target URL to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "xml", "html"),
		),
	)
	s.AddTool(baselineScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewZapModule(client)

		// Get parameters
		target := request.GetString("target", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		outputFormat := request.GetString("output_format", "")
		if outputFormat != "" {
			return mcp.NewToolResultError("Warning: output_format parameter not supported with direct Dagger calls"), nil
		}

		// Perform baseline scan
		output, err := module.BaselineScan(ctx, target)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("ZAP baseline scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// ZAP get version tool
	getVersionTool := mcp.NewTool("zap_get_version",
		mcp.WithDescription("Get OWASP ZAP version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewZapModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("ZAP get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}