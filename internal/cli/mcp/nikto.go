package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddNiktoTools adds Nikto (web vulnerability scanner) MCP tool implementations using direct Dagger calls
func AddNiktoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addNiktoToolsDirect(s)
}

// addNiktoToolsDirect adds Nikto tools using direct Dagger module calls
func addNiktoToolsDirect(s *server.MCPServer) {
	// Nikto basic scan tool
	scanHostTool := mcp.NewTool("nikto_scan_host",
		mcp.WithDescription("Scan web host for vulnerabilities using nikto"),
		mcp.WithString("host",
			mcp.Description("Target host to scan (e.g., example.com or 192.168.1.100)"),
			mcp.Required(),
		),
	)
	s.AddTool(scanHostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Get parameters
		host := request.GetString("host", "")
		if host == "" {
			return mcp.NewToolResultError("host is required"), nil
		}

		// Scan host
		output, err := module.ScanHost(ctx, host)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto SSL scan tool
	scanSSLTool := mcp.NewTool("nikto_scan_ssl",
		mcp.WithDescription("Scan HTTPS host with SSL/TLS using nikto"),
		mcp.WithString("host",
			mcp.Description("Target HTTPS host to scan"),
			mcp.Required(),
		),
		mcp.WithNumber("port",
			mcp.Description("SSL port number (default: 443)"),
		),
	)
	s.AddTool(scanSSLTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Get parameters
		host := request.GetString("host", "")
		if host == "" {
			return mcp.NewToolResultError("host is required"), nil
		}
		port := request.GetInt("port", 443)

		// Scan with SSL
		output, err := module.ScanWithSSL(ctx, host, port)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto SSL scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto scan with tuning options tool
	scanWithTuningTool := mcp.NewTool("nikto_scan_tuning",
		mcp.WithDescription("Scan with specific vulnerability tuning using nikto"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("tuning",
			mcp.Description("Tuning options: 1=interesting files, 2=misconfigs, 3=info disclosure, 4=injection, 8=command exec, 9=sql injection"),
		),
	)
	s.AddTool(scanWithTuningTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Get parameters
		host := request.GetString("host", "")
		if host == "" {
			return mcp.NewToolResultError("host is required"), nil
		}
		tuning := request.GetString("tuning", "")

		// Scan with tuning
		output, err := module.ScanWithTuning(ctx, host, tuning)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto scan with tuning failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto scan multiple hosts from file tool
	scanHostsFileTool := mcp.NewTool("nikto_scan_hosts_file",
		mcp.WithDescription("Scan multiple hosts from file using nikto"),
		mcp.WithString("hosts_file",
			mcp.Description("Path to file containing list of hosts (one per line)"),
			mcp.Required(),
		),
	)
	s.AddTool(scanHostsFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Get parameters
		hostsFile := request.GetString("hosts_file", "")
		if hostsFile == "" {
			return mcp.NewToolResultError("hosts_file is required"), nil
		}

		// Scan hosts file
		output, err := module.ScanHostsFile(ctx, hostsFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto scan hosts file failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto scan with authentication tool
	scanWithAuthTool := mcp.NewTool("nikto_scan_auth",
		mcp.WithDescription("Scan with basic authentication using nikto"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("auth_method",
			mcp.Description("Authentication method (basic, ntlm, etc.)"),
		),
		mcp.WithString("credentials",
			mcp.Description("Authentication credentials in format username:password"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithAuthTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Get parameters
		host := request.GetString("host", "")
		if host == "" {
			return mcp.NewToolResultError("host is required"), nil
		}
		authMethod := request.GetString("auth_method", "")
		credentials := request.GetString("credentials", "")
		if credentials == "" {
			return mcp.NewToolResultError("credentials is required"), nil
		}

		// Scan with auth
		output, err := module.ScanWithAuth(ctx, host, authMethod, credentials)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto scan with auth failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto scan with proxy tool
	scanWithProxyTool := mcp.NewTool("nikto_scan_proxy",
		mcp.WithDescription("Scan through proxy using nikto"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("proxy_host",
			mcp.Description("Proxy host/IP address"),
			mcp.Required(),
		),
		mcp.WithString("proxy_port",
			mcp.Description("Proxy port number"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithProxyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Get parameters
		host := request.GetString("host", "")
		if host == "" {
			return mcp.NewToolResultError("host is required"), nil
		}
		proxyHost := request.GetString("proxy_host", "")
		if proxyHost == "" {
			return mcp.NewToolResultError("proxy_host is required"), nil
		}
		proxyPort := request.GetString("proxy_port", "")
		if proxyPort == "" {
			return mcp.NewToolResultError("proxy_port is required"), nil
		}

		// Scan with proxy
		output, err := module.ScanWithProxy(ctx, host, proxyHost, proxyPort)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto scan with proxy failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto scan with evasion techniques tool
	scanWithEvasionTool := mcp.NewTool("nikto_scan_evasion",
		mcp.WithDescription("Scan with evasion techniques using nikto"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("evasion_level",
			mcp.Description("Evasion techniques: 1=random URI encoding, 2=directory self-reference, 3=premature URL ending, 4=prepend long random string, 5=fake parameter, 6=TAB as request spacer, 7=change case, 8=Windows directory separator"),
		),
	)
	s.AddTool(scanWithEvasionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Get parameters
		host := request.GetString("host", "")
		if host == "" {
			return mcp.NewToolResultError("host is required"), nil
		}
		evasionLevel := request.GetString("evasion_level", "")

		// Scan with evasion
		output, err := module.ScanWithEvasion(ctx, host, evasionLevel)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto scan with evasion failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto database check tool
	dbCheckTool := mcp.NewTool("nikto_database_check",
		mcp.WithDescription("Check Nikto scan database for errors using nikto"),
	)
	s.AddTool(dbCheckTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Check database
		output, err := module.DatabaseCheck(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto database check failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto update database tool
	updateTool := mcp.NewTool("nikto_update",
		mcp.WithDescription("Update Nikto database using nikto"),
	)
	s.AddTool(updateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Update database
		output, err := module.UpdateDatabase(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto update failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto version tool
	versionTool := mcp.NewTool("nikto_version",
		mcp.WithDescription("Get Nikto version information using nikto"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Nikto find only mode tool
	findOnlyTool := mcp.NewTool("nikto_find_only",
		mcp.WithDescription("Find HTTP(S) ports without performing security scan using nikto"),
		mcp.WithString("host",
			mcp.Description("Target host to check for HTTP(S) ports"),
			mcp.Required(),
		),
	)
	s.AddTool(findOnlyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewNiktoModule(client)

		// Get parameters
		host := request.GetString("host", "")
		if host == "" {
			return mcp.NewToolResultError("host is required"), nil
		}

		// Find only scan
		output, err := module.FindOnly(ctx, host)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nikto find only failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}