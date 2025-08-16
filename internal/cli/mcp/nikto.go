package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddNiktoTools adds Nikto (web vulnerability scanner) MCP tool implementations
func AddNiktoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Nikto scan host tool
	scanHostTool := mcp.NewTool("nikto_scan_host",
		mcp.WithDescription("Scan web host for vulnerabilities using Nikto"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "html", "txt"),
		),
	)
	s.AddTool(scanHostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"security", "nikto", "--host", host}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Nikto SSL scan tool
	scanSSLTool := mcp.NewTool("nikto_scan_ssl",
		mcp.WithDescription("Scan host with SSL/TLS analysis"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithNumber("port",
			mcp.Description("Port number for SSL scan"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "html", "txt"),
		),
	)
	s.AddTool(scanSSLTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"security", "nikto", "--host", host, "--ssl"}
		if port := request.GetInt("port", 0); port > 0 {
			args = append(args, "--port", string(rune(port)))
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Nikto scan with custom options tool
	scanWithOptionsTool := mcp.NewTool("nikto_scan_with_options",
		mcp.WithDescription("Scan with custom Nikto options"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("plugins",
			mcp.Description("Comma-separated list of plugins to use"),
		),
		mcp.WithString("tuning",
			mcp.Description("Tuning options for scan types"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "html", "txt"),
		),
	)
	s.AddTool(scanWithOptionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"security", "nikto", "--host", host}
		if plugins := request.GetString("plugins", ""); plugins != "" {
			args = append(args, "--plugins", plugins)
		}
		if tuning := request.GetString("tuning", ""); tuning != "" {
			args = append(args, "--tuning", tuning)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Nikto scan file or directory tool
	scanFileTool := mcp.NewTool("nikto_scan_file",
		mcp.WithDescription("Scan using hosts from file"),
		mcp.WithString("hosts_file",
			mcp.Description("Path to file containing list of hosts"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "html", "txt"),
		),
	)
	s.AddTool(scanFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		hostsFile := request.GetString("hosts_file", "")
		args := []string{"security", "nikto", "--hosts-file", hostsFile}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Nikto get version tool
	getVersionTool := mcp.NewTool("nikto_get_version",
		mcp.WithDescription("Get Nikto version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "nikto", "--version"}
		return executeShipCommand(args)
	})
}