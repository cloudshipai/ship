package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddZapTools adds OWASP ZAP (web application security scanner) MCP tool implementations
func AddZapTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// ZAP passive scan tool
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
		target := request.GetString("target", "")
		args := []string{"security", "zap", target, "--scan-type", "passive"}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		args := []string{"security", "zap", target, "--scan-type", "active"}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		args := []string{"security", "zap", target, "--scan-type", "spider"}
		if maxDepth := request.GetInt("max_depth", 0); maxDepth > 0 {
			args = append(args, "--max-depth", string(rune(maxDepth)))
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		args := []string{"security", "zap", target, "--scan-type", "full"}
		if maxDuration := request.GetInt("max_duration", 0); maxDuration > 0 {
			args = append(args, "--max-duration", string(rune(maxDuration)))
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
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
		target := request.GetString("target", "")
		args := []string{"security", "zap", target, "--scan-type", "baseline"}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// ZAP get version tool
	getVersionTool := mcp.NewTool("zap_get_version",
		mcp.WithDescription("Get OWASP ZAP version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "zap", "--version"}
		return executeShipCommand(args)
	})
}