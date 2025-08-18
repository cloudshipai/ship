package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddKubeHunterTools adds Kube-hunter (Kubernetes penetration testing) MCP tool implementations using direct Dagger calls
func AddKubeHunterTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addKubeHunterToolsDirect(s)
}

// addKubeHunterToolsDirect adds Kube-hunter tools using direct Dagger module calls
func addKubeHunterToolsDirect(s *server.MCPServer) {
	// Kube-hunter remote scan
	remoteScanTool := mcp.NewTool("kube_hunter_remote_scan",
		mcp.WithDescription("Scan remote Kubernetes cluster using kube-hunter"),
		mcp.WithString("remote",
			mcp.Description("Remote hostname or IP address to scan"),
			mcp.Required(),
		),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting (attempts to exploit vulnerabilities)"),
		),
		mcp.WithString("report",
			mcp.Description("Report format: json, yaml, table"),
			mcp.Enum("json", "yaml", "table"),
		),
		mcp.WithString("log_level",
			mcp.Description("Log level: debug, info, warning"),
			mcp.Enum("debug", "info", "warning"),
		),
	)
	s.AddTool(remoteScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeHunterModule(client)

		// Get parameters
		remote := request.GetString("remote", "")
		if remote == "" {
			return mcp.NewToolResultError("remote is required"), nil
		}

		active := request.GetBool("active", false)
		reportFormat := request.GetString("report", "json")

		// Run remote scan
		output, err := module.ScanRemote(ctx, remote, active, reportFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-hunter remote scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-hunter CIDR scan
	cidrScanTool := mcp.NewTool("kube_hunter_cidr_scan",
		mcp.WithDescription("Scan CIDR range for Kubernetes clusters using kube-hunter"),
		mcp.WithString("cidr",
			mcp.Description("CIDR range to scan (e.g., 192.168.1.0/24)"),
			mcp.Required(),
		),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting"),
		),
		mcp.WithString("report",
			mcp.Description("Report format: json, yaml, table"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(cidrScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeHunterModule(client)

		// Get parameters
		cidr := request.GetString("cidr", "")
		if cidr == "" {
			return mcp.NewToolResultError("cidr is required"), nil
		}

		active := request.GetBool("active", false)
		reportFormat := request.GetString("report", "json")

		// Run CIDR scan
		output, err := module.ScanCIDR(ctx, cidr, active, reportFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-hunter CIDR scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-hunter interface scan
	interfaceScanTool := mcp.NewTool("kube_hunter_interface_scan",
		mcp.WithDescription("Scan network interface using kube-hunter"),
		mcp.WithString("interface",
			mcp.Description("Network interface to scan (e.g., eth0)"),
			mcp.Required(),
		),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting"),
		),
		mcp.WithString("report",
			mcp.Description("Report format: json, yaml, table"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(interfaceScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeHunterModule(client)

		// Get parameters
		networkInterface := request.GetString("interface", "")
		if networkInterface == "" {
			return mcp.NewToolResultError("interface is required"), nil
		}

		active := request.GetBool("active", false)
		reportFormat := request.GetString("report", "json")

		// Run interface scan
		output, err := module.ScanInterface(ctx, networkInterface, active, reportFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-hunter interface scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-hunter pod scan
	podScanTool := mcp.NewTool("kube_hunter_pod_scan",
		mcp.WithDescription("Run kube-hunter as pod in Kubernetes cluster"),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting"),
		),
		mcp.WithString("report",
			mcp.Description("Report format: json, yaml, table"),
			mcp.Enum("json", "yaml", "table"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithBoolean("quick",
			mcp.Description("Perform quick scan"),
		),
	)
	s.AddTool(podScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeHunterModule(client)

		// Get parameters
		kubeconfig := request.GetString("kubeconfig", "")
		active := request.GetBool("active", false)
		reportFormat := request.GetString("report", "json")

		// Run pod scan
		output, err := module.ScanPod(ctx, kubeconfig, active, reportFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-hunter pod scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-hunter list tests
	listTestsTool := mcp.NewTool("kube_hunter_list_tests",
		mcp.WithDescription("List all available kube-hunter tests"),
		mcp.WithBoolean("active",
			mcp.Description("Include active hunting tests"),
		),
	)
	s.AddTool(listTestsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeHunterModule(client)

		// Get parameters
		showActive := request.GetBool("active", false)

		// List tests
		output, err := module.ListTests(ctx, showActive)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list kube-hunter tests: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Kube-hunter custom hunters
	customHuntersTool := mcp.NewTool("kube_hunter_custom_hunters",
		mcp.WithDescription("Run kube-hunter with specific hunters enabled/disabled"),
		mcp.WithString("target",
			mcp.Description("Target to scan (IP, hostname, CIDR, or interface)"),
			mcp.Required(),
		),
		mcp.WithString("include_hunters",
			mcp.Description("Comma-separated list of hunter types to include"),
		),
		mcp.WithString("exclude_hunters",
			mcp.Description("Comma-separated list of hunter types to exclude"),
		),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting"),
		),
	)
	s.AddTool(customHuntersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewKubeHunterModule(client)

		// Get parameters
		target := request.GetString("target", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		// Parse include hunters
		includeHunters := []string{}
		if include := request.GetString("include_hunters", ""); include != "" {
			for _, h := range strings.Split(include, ",") {
				h = strings.TrimSpace(h)
				if h != "" {
					includeHunters = append(includeHunters, h)
				}
			}
		}

		// Parse exclude hunters
		excludeHunters := []string{}
		if exclude := request.GetString("exclude_hunters", ""); exclude != "" {
			for _, h := range strings.Split(exclude, ",") {
				h = strings.TrimSpace(h)
				if h != "" {
					excludeHunters = append(excludeHunters, h)
				}
			}
		}

		active := request.GetBool("active", false)

		// Run custom hunters
		output, err := module.RunCustomHunters(ctx, target, includeHunters, excludeHunters, active)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("kube-hunter custom hunters failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}