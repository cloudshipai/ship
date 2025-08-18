package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCheckSSLCertTools adds SSL certificate validation MCP tool implementations
func AddCheckSSLCertTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Check SSL certificate for host tool
	checkHostTool := mcp.NewTool("check_ssl_cert_host",
		mcp.WithDescription("Check SSL certificate for a remote host"),
		mcp.WithString("host",
			mcp.Description("Server hostname or IP address to check"),
			mcp.Required(),
		),
		mcp.WithString("port",
			mcp.Description("TCP port (default: 443)"),
		),
		mcp.WithString("warning",
			mcp.Description("Days before expiry to warn (default: 20)"),
		),
		mcp.WithString("critical",
			mcp.Description("Days before expiry for critical (default: 15)"),
		),
		mcp.WithString("protocol",
			mcp.Description("Protocol to use (https, smtp, ftp, imap, pop3, etc.)"),
		),
		mcp.WithBoolean("selfsigned",
			mcp.Description("Allow self-signed certificates"),
		),
	)
	s.AddTool(checkHostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"check_ssl_cert", "-H", host}
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "-p", port)
		}
		if warning := request.GetString("warning", ""); warning != "" {
			args = append(args, "-w", warning)
		}
		if critical := request.GetString("critical", ""); critical != "" {
			args = append(args, "-c", critical)
		}
		if protocol := request.GetString("protocol", ""); protocol != "" {
			args = append(args, "-P", protocol)
		}
		if request.GetBool("selfsigned", false) {
			args = append(args, "-s")
		}
		return executeShipCommand(args)
	})

	// Check SSL certificate from file tool
	checkFileTool := mcp.NewTool("check_ssl_cert_file",
		mcp.WithDescription("Check SSL certificate from a local file"),
		mcp.WithString("file",
			mcp.Description("Path to certificate file (PEM, DER, PKCS12, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("warning",
			mcp.Description("Days before expiry to warn (default: 20)"),
		),
		mcp.WithString("critical",
			mcp.Description("Days before expiry for critical (default: 15)"),
		),
		mcp.WithBoolean("selfsigned",
			mcp.Description("Allow self-signed certificates"),
		),
	)
	s.AddTool(checkFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		file := request.GetString("file", "")
		args := []string{"check_ssl_cert", "-f", file}
		if warning := request.GetString("warning", ""); warning != "" {
			args = append(args, "-w", warning)
		}
		if critical := request.GetString("critical", ""); critical != "" {
			args = append(args, "-c", critical)
		}
		if request.GetBool("selfsigned", false) {
			args = append(args, "-s")
		}
		return executeShipCommand(args)
	})

	// Check SSL certificate with chain validation tool
	checkChainTool := mcp.NewTool("check_ssl_cert_chain",
		mcp.WithDescription("Check SSL certificate with chain validation"),
		mcp.WithString("host",
			mcp.Description("Server hostname or IP address to check"),
			mcp.Required(),
		),
		mcp.WithString("rootcert",
			mcp.Description("Path to root certificate for validation"),
		),
		mcp.WithBoolean("check_chain",
			mcp.Description("Verify certificate chain integrity"),
		),
		mcp.WithBoolean("noauth",
			mcp.Description("Ignore authority warnings"),
		),
	)
	s.AddTool(checkChainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"check_ssl_cert", "-H", host}
		if rootcert := request.GetString("rootcert", ""); rootcert != "" {
			args = append(args, "-r", rootcert)
		}
		if request.GetBool("check_chain", false) {
			args = append(args, "--check-chain")
		}
		if request.GetBool("noauth", false) {
			args = append(args, "-A")
		}
		return executeShipCommand(args)
	})

	// Check SSL certificate with fingerprint tool
	checkFingerprintTool := mcp.NewTool("check_ssl_cert_fingerprint",
		mcp.WithDescription("Check SSL certificate fingerprint"),
		mcp.WithString("host",
			mcp.Description("Server hostname or IP address to check"),
			mcp.Required(),
		),
		mcp.WithString("fingerprint",
			mcp.Description("Expected certificate fingerprint"),
			mcp.Required(),
		),
	)
	s.AddTool(checkFingerprintTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		fingerprint := request.GetString("fingerprint", "")
		args := []string{"check_ssl_cert", "-H", host, "--fingerprint", fingerprint}
		return executeShipCommand(args)
	})

	// Check SSL certificate with all checks tool
	checkAllTool := mcp.NewTool("check_ssl_cert_all",
		mcp.WithDescription("Check SSL certificate with all optional checks enabled"),
		mcp.WithString("host",
			mcp.Description("Server hostname or IP address to check"),
			mcp.Required(),
		),
		mcp.WithString("timeout",
			mcp.Description("Connection timeout in seconds"),
		),
		mcp.WithBoolean("debug",
			mcp.Description("Produce debugging output"),
		),
	)
	s.AddTool(checkAllTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"check_ssl_cert", "-H", host, "--all"}
		if timeout := request.GetString("timeout", ""); timeout != "" {
			args = append(args, "--timeout", timeout)
		}
		if request.GetBool("debug", false) {
			args = append(args, "-d")
		}
		return executeShipCommand(args)
	})

	// Check SSL cert get version tool
	getVersionTool := mcp.NewTool("check_ssl_cert_get_version",
		mcp.WithDescription("Get check_ssl_cert version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"check_ssl_cert", "--version"}
		return executeShipCommand(args)
	})
}