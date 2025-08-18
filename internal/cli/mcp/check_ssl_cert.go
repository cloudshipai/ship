package mcp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCheckSSLCertTools adds SSL certificate validation MCP tool implementations
func AddCheckSSLCertTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addCheckSSLCertToolsDirect(s)
}

// addCheckSSLCertToolsDirect implements direct Dagger calls for SSL certificate checking tools
func addCheckSSLCertToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		host := request.GetString("host", "")
		portStr := request.GetString("port", "443")
		port, _ := strconv.Atoi(portStr)
		warningStr := request.GetString("warning", "20")
		warningDays, _ := strconv.Atoi(warningStr)
		criticalStr := request.GetString("critical", "15")
		criticalDays, _ := strconv.Atoi(criticalStr)
		protocol := request.GetString("protocol", "")
		allowSelfSigned := request.GetBool("selfsigned", false)

		// Create SSL cert module and check host certificate
		sslCertModule := modules.NewCheckSSLCertModule(client)
		result, err := sslCertModule.CheckCertificateWithAdvancedOptions(ctx, host, port, protocol, warningDays, criticalDays, allowSelfSigned, "", false, false, 0, false)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("check SSL certificate for host failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		file := request.GetString("file", "")
		warningStr := request.GetString("warning", "20")
		warningDays, _ := strconv.Atoi(warningStr)
		criticalStr := request.GetString("critical", "15")
		criticalDays, _ := strconv.Atoi(criticalStr)
		allowSelfSigned := request.GetBool("selfsigned", false)

		// Create SSL cert module and check certificate from file
		sslCertModule := modules.NewCheckSSLCertModule(client)
		result, err := sslCertModule.CheckCertificateFromFile(ctx, file, warningDays, criticalDays, allowSelfSigned)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("check SSL certificate from file failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		host := request.GetString("host", "")
		rootCert := request.GetString("rootcert", "")
		checkChain := request.GetBool("check_chain", false)
		ignoreAuth := request.GetBool("noauth", false)

		// Create SSL cert module and check certificate chain
		sslCertModule := modules.NewCheckSSLCertModule(client)
		result, err := sslCertModule.CheckCertificateWithAdvancedOptions(ctx, host, 443, "", 0, 0, false, rootCert, checkChain, ignoreAuth, 0, false)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("check SSL certificate chain failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		host := request.GetString("host", "")
		fingerprint := request.GetString("fingerprint", "")

		// Create SSL cert module and check certificate fingerprint
		sslCertModule := modules.NewCheckSSLCertModule(client)
		result, err := sslCertModule.CheckCertificateFingerprint(ctx, host, 443, fingerprint)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("check SSL certificate fingerprint failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		host := request.GetString("host", "")
		timeoutStr := request.GetString("timeout", "0")
		timeout, _ := strconv.Atoi(timeoutStr)
		debug := request.GetBool("debug", false)

		// Create SSL cert module and run comprehensive check
		sslCertModule := modules.NewCheckSSLCertModule(client)
		result, err := sslCertModule.CheckCertificateComprehensive(ctx, host, 443, timeout, debug)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("comprehensive SSL certificate check failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Check SSL cert get version tool
	getVersionTool := mcp.NewTool("check_ssl_cert_get_version",
		mcp.WithDescription("Get check_ssl_cert version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create SSL cert module and get version
		sslCertModule := modules.NewCheckSSLCertModule(client)
		result, err := sslCertModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get SSL cert tool version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}