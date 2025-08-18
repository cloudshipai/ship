package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddNiktoTools adds Nikto (web vulnerability scanner) MCP tool implementations using real CLI commands
func AddNiktoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Nikto basic scan tool
	scanHostTool := mcp.NewTool("nikto_scan_host",
		mcp.WithDescription("Scan web host for vulnerabilities using real nikto CLI"),
		mcp.WithString("host",
			mcp.Description("Target host to scan (e.g., example.com or 192.168.1.100)"),
			mcp.Required(),
		),
		mcp.WithString("port",
			mcp.Description("Port number to scan (default: 80 for HTTP, 443 for HTTPS)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("csv", "htm", "xml", "txt"),
		),
	)
	s.AddTool(scanHostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"nikto.pl", "-h", host}
		
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "-p", port)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "-Format", format)
		}
		
		return executeShipCommand(args)
	})

	// Nikto SSL scan tool
	scanSSLTool := mcp.NewTool("nikto_scan_ssl",
		mcp.WithDescription("Scan HTTPS host with SSL/TLS using real nikto CLI"),
		mcp.WithString("host",
			mcp.Description("Target HTTPS host to scan"),
			mcp.Required(),
		),
		mcp.WithString("port",
			mcp.Description("SSL port number (default: 443)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("csv", "htm", "xml", "txt"),
		),
	)
	s.AddTool(scanSSLTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"nikto", "-h", host, "-ssl"}
		
		if port := request.GetString("port", ""); port != "" {
			args = append(args, "-p", port)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "-Format", format)
		}
		
		return executeShipCommand(args)
	})

	// Nikto scan with tuning options tool
	scanWithTuningTool := mcp.NewTool("nikto_scan_tuning",
		mcp.WithDescription("Scan with specific vulnerability tuning using real nikto CLI"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("tuning",
			mcp.Description("Tuning options: 1=interesting files, 2=misconfigs, 3=info disclosure, 4=injection, 8=command exec, 9=sql injection"),
		),
		mcp.WithString("display",
			mcp.Description("Display options: 1=redirects, 2=cookies, 3=200 responses, 4=auth required, D=debug, E=errors, P=progress, V=verbose"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("csv", "htm", "xml", "txt"),
		),
	)
	s.AddTool(scanWithTuningTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"nikto.pl", "-h", host}
		
		if tuning := request.GetString("tuning", ""); tuning != "" {
			args = append(args, "-Tuning", tuning)
		}
		if display := request.GetString("display", ""); display != "" {
			args = append(args, "-Display", display)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "-Format", format)
		}
		
		return executeShipCommand(args)
	})

	// Nikto scan multiple hosts from file tool
	scanHostsFileTool := mcp.NewTool("nikto_scan_hosts_file",
		mcp.WithDescription("Scan multiple hosts from file using real nikto CLI"),
		mcp.WithString("hosts_file",
			mcp.Description("Path to file containing list of hosts (one per line)"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("csv", "htm", "xml", "txt"),
		),
	)
	s.AddTool(scanHostsFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		hostsFile := request.GetString("hosts_file", "")
		args := []string{"nikto", "-h", hostsFile}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "-Format", format)
		}
		
		return executeShipCommand(args)
	})

	// Nikto scan with authentication tool
	scanWithAuthTool := mcp.NewTool("nikto_scan_auth",
		mcp.WithDescription("Scan with basic authentication using real nikto CLI"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("auth",
			mcp.Description("Authentication in format username:password or username:password:realm"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("csv", "htm", "xml", "txt"),
		),
	)
	s.AddTool(scanWithAuthTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		auth := request.GetString("auth", "")
		args := []string{"nikto", "-h", host, "-id", auth}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "-Format", format)
		}
		
		return executeShipCommand(args)
	})

	// Nikto scan with proxy tool
	scanWithProxyTool := mcp.NewTool("nikto_scan_proxy",
		mcp.WithDescription("Scan through proxy using real nikto CLI"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("proxy",
			mcp.Description("Proxy URL (e.g., http://127.0.0.1:8080)"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("csv", "htm", "xml", "txt"),
		),
	)
	s.AddTool(scanWithProxyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		proxy := request.GetString("proxy", "")
		args := []string{"nikto", "-h", host, "-useproxy", proxy}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "-Format", format)
		}
		
		return executeShipCommand(args)
	})

	// Nikto scan with evasion techniques tool
	scanWithEvasionTool := mcp.NewTool("nikto_scan_evasion",
		mcp.WithDescription("Scan with evasion techniques using real nikto CLI"),
		mcp.WithString("host",
			mcp.Description("Target host to scan"),
			mcp.Required(),
		),
		mcp.WithString("evasion",
			mcp.Description("Evasion techniques: 1=random URI encoding, 2=directory self-reference, 3=premature URL ending, 4=prepend long random string, 5=fake parameter, 6=TAB as request spacer, 7=change case, 8=Windows directory separator"),
		),
		mcp.WithString("timeout",
			mcp.Description("Request timeout in seconds (default: 10)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("csv", "htm", "xml", "txt"),
		),
	)
	s.AddTool(scanWithEvasionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"nikto.pl", "-h", host}
		
		if evasion := request.GetString("evasion", ""); evasion != "" {
			args = append(args, "-evasion", evasion)
		}
		if timeout := request.GetString("timeout", ""); timeout != "" {
			args = append(args, "-timeout", timeout)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "-Format", format)
		}
		
		return executeShipCommand(args)
	})

	// Nikto database check tool
	dbCheckTool := mcp.NewTool("nikto_database_check",
		mcp.WithDescription("Check Nikto scan database for errors using real nikto CLI"),
	)
	s.AddTool(dbCheckTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"nikto.pl", "-dbcheck"}
		return executeShipCommand(args)
	})

	// Nikto update database tool
	updateTool := mcp.NewTool("nikto_update",
		mcp.WithDescription("Update Nikto database using real nikto CLI"),
	)
	s.AddTool(updateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"nikto.pl", "-update"}
		return executeShipCommand(args)
	})

	// Nikto version tool
	versionTool := mcp.NewTool("nikto_version",
		mcp.WithDescription("Get Nikto version information using real nikto CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"nikto.pl", "-Version"}
		return executeShipCommand(args)
	})

	// Nikto find only mode tool
	findOnlyTool := mcp.NewTool("nikto_find_only",
		mcp.WithDescription("Find HTTP(S) ports without performing security scan using real nikto CLI"),
		mcp.WithString("host",
			mcp.Description("Target host to check for HTTP(S) ports"),
			mcp.Required(),
		),
	)
	s.AddTool(findOnlyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host := request.GetString("host", "")
		args := []string{"nikto", "-h", host, "-findonly"}
		return executeShipCommand(args)
	})
}