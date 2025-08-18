package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddKubeHunterTools adds Kube-hunter (Kubernetes penetration testing) MCP tool implementations using real CLI commands
func AddKubeHunterTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Kube-hunter remote scan tool
	remoteScanTool := mcp.NewTool("kube_hunter_remote_scan",
		mcp.WithDescription("Scan specific IP addresses or DNS names using kube-hunter --remote"),
		mcp.WithString("target",
			mcp.Description("IP address or DNS name to scan"),
			mcp.Required(),
		),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting (exploit testing)"),
		),
		mcp.WithString("log",
			mcp.Description("Log level"),
			mcp.Enum("DEBUG", "INFO", "WARNING"),
		),
		mcp.WithString("report",
			mcp.Description("Report format"),
			mcp.Enum("json"),
		),
		mcp.WithString("dispatch",
			mcp.Description("Output method"),
			mcp.Enum("stdout", "http"),
		),
	)
	s.AddTool(remoteScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target := request.GetString("target", "")
		args := []string{"kube-hunter", "--remote", target}
		
		if request.GetBool("active", false) {
			args = append(args, "--active")
		}
		if log := request.GetString("log", ""); log != "" {
			args = append(args, "--log", log)
		}
		if report := request.GetString("report", ""); report != "" {
			args = append(args, "--report", report)
		}
		if dispatch := request.GetString("dispatch", ""); dispatch != "" {
			args = append(args, "--dispatch", dispatch)
		}
		
		return executeShipCommand(args)
	})

	// Kube-hunter CIDR scan tool
	cidrScanTool := mcp.NewTool("kube_hunter_cidr_scan",
		mcp.WithDescription("Scan IP range using kube-hunter --cidr"),
		mcp.WithString("cidr",
			mcp.Description("CIDR range to scan (e.g., 192.168.0.0/24)"),
			mcp.Required(),
		),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting (exploit testing)"),
		),
		mcp.WithBoolean("mapping",
			mcp.Description("Show only network mapping of Kubernetes nodes"),
		),
		mcp.WithBoolean("quick",
			mcp.Description("Limit subnet scanning to /24 CIDR"),
		),
		mcp.WithString("log",
			mcp.Description("Log level"),
			mcp.Enum("DEBUG", "INFO", "WARNING"),
		),
		mcp.WithString("report",
			mcp.Description("Report format"),
			mcp.Enum("json"),
		),
	)
	s.AddTool(cidrScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cidr := request.GetString("cidr", "")
		args := []string{"kube-hunter", "--cidr", cidr}
		
		if request.GetBool("active", false) {
			args = append(args, "--active")
		}
		if request.GetBool("mapping", false) {
			args = append(args, "--mapping")
		}
		if request.GetBool("quick", false) {
			args = append(args, "--quick")
		}
		if log := request.GetString("log", ""); log != "" {
			args = append(args, "--log", log)
		}
		if report := request.GetString("report", ""); report != "" {
			args = append(args, "--report", report)
		}
		
		return executeShipCommand(args)
	})

	// Kube-hunter interface scan tool
	interfaceScanTool := mcp.NewTool("kube_hunter_interface_scan",
		mcp.WithDescription("Scan all local network interfaces using kube-hunter --interface"),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting (exploit testing)"),
		),
		mcp.WithBoolean("quick",
			mcp.Description("Limit subnet scanning to /24 CIDR"),
		),
		mcp.WithString("log",
			mcp.Description("Log level"),
			mcp.Enum("DEBUG", "INFO", "WARNING"),
		),
		mcp.WithString("report",
			mcp.Description("Report format"),
			mcp.Enum("json"),
		),
	)
	s.AddTool(interfaceScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kube-hunter", "--interface"}
		
		if request.GetBool("active", false) {
			args = append(args, "--active")
		}
		if request.GetBool("quick", false) {
			args = append(args, "--quick")
		}
		if log := request.GetString("log", ""); log != "" {
			args = append(args, "--log", log)
		}
		if report := request.GetString("report", ""); report != "" {
			args = append(args, "--report", report)
		}
		
		return executeShipCommand(args)
	})

	// Kube-hunter pod scan tool
	podScanTool := mcp.NewTool("kube_hunter_pod_scan",
		mcp.WithDescription("Scan from within a Kubernetes pod using kube-hunter --pod"),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting (exploit testing)"),
		),
		mcp.WithBoolean("k8s_auto_discover_nodes",
			mcp.Description("Query Kubernetes for all nodes and scan them"),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
		),
		mcp.WithString("service_account_token",
			mcp.Description("JWT Bearer token of service account"),
		),
		mcp.WithString("log",
			mcp.Description("Log level"),
			mcp.Enum("DEBUG", "INFO", "WARNING"),
		),
		mcp.WithString("report",
			mcp.Description("Report format"),
			mcp.Enum("json"),
		),
	)
	s.AddTool(podScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kube-hunter", "--pod"}
		
		if request.GetBool("active", false) {
			args = append(args, "--active")
		}
		if request.GetBool("k8s_auto_discover_nodes", false) {
			args = append(args, "--k8s-auto-discover-nodes")
		}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		if token := request.GetString("service_account_token", ""); token != "" {
			args = append(args, "--service-account-token", token)
		}
		if log := request.GetString("log", ""); log != "" {
			args = append(args, "--log", log)
		}
		if report := request.GetString("report", ""); report != "" {
			args = append(args, "--report", report)
		}
		
		return executeShipCommand(args)
	})

	// Kube-hunter list tests tool
	listTestsTool := mcp.NewTool("kube_hunter_list_tests",
		mcp.WithDescription("List available tests using kube-hunter --list"),
		mcp.WithBoolean("active",
			mcp.Description("Include active hunting tests in the list"),
		),
		mcp.WithBoolean("raw_hunter_names",
			mcp.Description("Show raw hunter class names"),
		),
	)
	s.AddTool(listTestsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"kube-hunter", "--list"}
		
		if request.GetBool("active", false) {
			args = append(args, "--active")
		}
		if request.GetBool("raw_hunter_names", false) {
			args = append(args, "--raw-hunter-names")
		}
		
		return executeShipCommand(args)
	})

	// Kube-hunter custom hunters tool
	customHuntersTool := mcp.NewTool("kube_hunter_custom_hunters",
		mcp.WithDescription("Run specific hunters using kube-hunter --custom"),
		mcp.WithString("hunters",
			mcp.Description("Space-separated list of hunter class names"),
			mcp.Required(),
		),
		mcp.WithBoolean("active",
			mcp.Description("Enable active hunting (exploit testing)"),
		),
		mcp.WithString("target",
			mcp.Description("Target IP address or DNS name (for remote scanning)"),
		),
		mcp.WithString("log",
			mcp.Description("Log level"),
			mcp.Enum("DEBUG", "INFO", "WARNING"),
		),
	)
	s.AddTool(customHuntersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		hunters := request.GetString("hunters", "")
		args := []string{"kube-hunter", "--custom", hunters}
		
		if request.GetBool("active", false) {
			args = append(args, "--active")
		}
		if target := request.GetString("target", ""); target != "" {
			args = append(args, "--remote", target)
		}
		if log := request.GetString("log", ""); log != "" {
			args = append(args, "--log", log)
		}
		
		return executeShipCommand(args)
	})
}