package mcp

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddNmapTools adds Nmap network scanning tools to the MCP server
func AddNmapTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addNmapToolsDirect(s)
}

// addNmapToolsDirect adds Nmap tools using direct Dagger module calls
func addNmapToolsDirect(s *server.MCPServer) {
	// Scan host
	scanHostTool := mcp.NewTool("nmap_scan_host",
		mcp.WithDescription("Perform basic host scanning"),
		mcp.WithString("target",
			mcp.Description("Target host or IP address"),
			mcp.Required(),
		),
		mcp.WithString("scan_type",
			mcp.Description("Type of scan to perform"),
		),
	)
	s.AddTool(scanHostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		scanType := request.GetString("scan_type", "")

		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNmapModule(client)
		result, err := module.ScanHost(ctx, target, scanType)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Port scan
	portScanTool := mcp.NewTool("nmap_port_scan",
		mcp.WithDescription("Scan specific ports on a host"),
		mcp.WithString("target",
			mcp.Description("Target host or IP address"),
			mcp.Required(),
		),
		mcp.WithString("ports",
			mcp.Description("Ports to scan (e.g., '80,443' or '1-1000')"),
		),
	)
	s.AddTool(portScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		ports := request.GetString("ports", "")

		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNmapModule(client)
		result, err := module.PortScan(ctx, target, ports)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Service detection
	serviceDetectionTool := mcp.NewTool("nmap_service_detection",
		mcp.WithDescription("Detect services and versions running on target"),
		mcp.WithString("target",
			mcp.Description("Target host or IP address"),
			mcp.Required(),
		),
	)
	s.AddTool(serviceDetectionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")

		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNmapModule(client)
		result, err := module.ServiceDetection(ctx, target)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("service detection failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Vulnerability scan
	vulnScanTool := mcp.NewTool("nmap_vulnerability_scan",
		mcp.WithDescription("Scan for vulnerabilities using NSE scripts"),
		mcp.WithString("target",
			mcp.Description("Target host or IP address"),
			mcp.Required(),
		),
		mcp.WithString("script_category",
			mcp.Description("NSE script category"),
		),
	)
	s.AddTool(vulnScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		scriptCategory := request.GetString("script_category", "safe")

		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNmapModule(client)
		result, err := module.VulnerabilityScan(ctx, target, scriptCategory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("vulnerability scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Network discovery
	networkDiscoveryTool := mcp.NewTool("nmap_network_discovery",
		mcp.WithDescription("Discover hosts on a network"),
		mcp.WithString("network",
			mcp.Description("Network range (e.g., 192.168.1.0/24)"),
			mcp.Required(),
		),
	)
	s.AddTool(networkDiscoveryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		network := request.GetString("network", "")

		if network == "" {
			return mcp.NewToolResultError("network is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNmapModule(client)
		result, err := module.NetworkDiscovery(ctx, network)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("discovery failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Script scan
	scriptScanTool := mcp.NewTool("nmap_script_scan",
		mcp.WithDescription("Run specific NSE scripts"),
		mcp.WithString("target",
			mcp.Description("Target host or IP address"),
			mcp.Required(),
		),
		mcp.WithString("script",
			mcp.Description("NSE script name (e.g., http-headers, ssl-cert)"),
			mcp.Required(),
		),
	)
	s.AddTool(scriptScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		script := request.GetString("script", "")

		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}
		if script == "" {
			return mcp.NewToolResultError("script is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNmapModule(client)
		result, err := module.ScriptScan(ctx, target, script)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("script scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Get version
	versionTool := mcp.NewTool("nmap_get_version",
		mcp.WithDescription("Get Nmap version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNmapModule(client)
		result, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get version: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}