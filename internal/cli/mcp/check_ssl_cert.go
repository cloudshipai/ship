package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCheckSSLCertTools adds SSL certificate validation MCP tool implementations
func AddCheckSSLCertTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Check SSL certificate tool
	checkCertificateTool := mcp.NewTool("check_ssl_cert_validate",
		mcp.WithDescription("Validate SSL certificate for domain or IP"),
		mcp.WithString("target",
			mcp.Description("Domain name or IP address to check"),
			mcp.Required(),
		),
		mcp.WithString("port",
			mcp.Description("Port number (default: 443)"),
		),
		mcp.WithString("timeout",
			mcp.Description("Connection timeout (e.g., 10s, 30s)"),
		),
	)
	s.AddTool(checkCertificateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"security", "check-ssl-cert", "validate", target}
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "--port", port)
		}
		if timeout := request.GetString("timeout", ""); timeout != "" {
			args = append(args, "--timeout", timeout)
		}
		return executeShipCommand(args)
	})

	// Check SSL certificate expiry tool
	checkExpiryTool := mcp.NewTool("check_ssl_cert_expiry",
		mcp.WithDescription("Check SSL certificate expiration date"),
		mcp.WithString("target",
			mcp.Description("Domain name or IP address to check"),
			mcp.Required(),
		),
		mcp.WithString("warning_days",
			mcp.Description("Days before expiry to warn (default: 30)"),
		),
	)
	s.AddTool(checkExpiryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"security", "check-ssl-cert", "expiry", target}
		if warningDays := request.GetString("warning_days", ""); warningDays != "" {
			args = append(args, "--warning", warningDays)
		}
		return executeShipCommand(args)
	})

	// Check SSL certificate chain tool
	checkChainTool := mcp.NewTool("check_ssl_cert_chain",
		mcp.WithDescription("Validate SSL certificate chain"),
		mcp.WithString("target",
			mcp.Description("Domain name or IP address to check"),
			mcp.Required(),
		),
		mcp.WithString("ca_file",
			mcp.Description("Path to CA bundle file (optional)"),
		),
	)
	s.AddTool(checkChainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"security", "check-ssl-cert", "chain", target}
		if caFile := request.GetString("ca_file", ""); caFile != "" {
			args = append(args, "--ca-file", caFile)
		}
		return executeShipCommand(args)
	})

	// Check SSL certificate file tool
	checkFileTool := mcp.NewTool("check_ssl_cert_file",
		mcp.WithDescription("Validate SSL certificate from file"),
		mcp.WithString("cert_file",
			mcp.Description("Path to certificate file"),
			mcp.Required(),
		),
		mcp.WithString("key_file",
			mcp.Description("Path to private key file (optional)"),
		),
	)
	s.AddTool(checkFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		certFile := request.GetString("cert_file", "")
		args := []string{"security", "check-ssl-cert", "file", certFile}
		if keyFile := request.GetString("key_file", ""); keyFile != "" {
			args = append(args, "--key", keyFile)
		}
		return executeShipCommand(args)
	})

	// Check SSL certificate bulk tool
	checkBulkTool := mcp.NewTool("check_ssl_cert_bulk",
		mcp.WithDescription("Check multiple SSL certificates from a list"),
		mcp.WithString("targets_file",
			mcp.Description("File containing list of targets to check"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (json, csv, table)"),
		),
	)
	s.AddTool(checkBulkTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		targetsFile := request.GetString("targets_file", "")
		args := []string{"security", "check-ssl-cert", "bulk", targetsFile}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Check SSL cert get version tool
	getVersionTool := mcp.NewTool("check_ssl_cert_get_version",
		mcp.WithDescription("Get SSL certificate checker version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "check-ssl-cert", "--version"}
		return executeShipCommand(args)
	})
}