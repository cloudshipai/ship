package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTerrascanTools adds Terrascan (IaC security scanner) MCP tool implementations
func AddTerrascanTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Terrascan scan directory tool
	scanDirectoryTool := mcp.NewTool("terrascan_scan_directory",
		mcp.WithDescription("Scan directory for IaC security issues using Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "terrascan", directory}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Terrascan scan Terraform files tool
	scanTerraformTool := mcp.NewTool("terrascan_scan_terraform",
		mcp.WithDescription("Scan Terraform files specifically using Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanTerraformTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "terrascan", directory, "--type", "terraform"}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Terrascan scan Kubernetes manifests tool
	scanKubernetesTool := mcp.NewTool("terrascan_scan_kubernetes",
		mcp.WithDescription("Scan Kubernetes manifests using Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Kubernetes manifests"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanKubernetesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "terrascan", directory, "--type", "k8s"}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Terrascan scan with severity filter tool
	scanWithSeverityTool := mcp.NewTool("terrascan_scan_with_severity",
		mcp.WithDescription("Scan with minimum severity level using Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level"),
			mcp.Required(),
			mcp.Enum("LOW", "MEDIUM", "HIGH"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanWithSeverityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		severity := request.GetString("severity", "")
		args := []string{"security", "terrascan", directory, "--severity", severity}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Terrascan scan remote repository tool
	scanRemoteTool := mcp.NewTool("terrascan_scan_remote",
		mcp.WithDescription("Scan remote repository using Terrascan"),
		mcp.WithString("repo_url",
			mcp.Description("URL of the remote repository"),
			mcp.Required(),
		),
		mcp.WithString("repo_type",
			mcp.Description("Type of repository"),
			mcp.Enum("git", "s3", "gcs", "http"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanRemoteTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		repoURL := request.GetString("repo_url", "")
		args := []string{"security", "terrascan", "--remote-url", repoURL}
		if repoType := request.GetString("repo_type", ""); repoType != "" {
			args = append(args, "--remote-type", repoType)
		}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Terrascan scan with policy path tool
	scanWithPolicyTool := mcp.NewTool("terrascan_scan_with_policy",
		mcp.WithDescription("Scan using custom policy path with Terrascan"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan"),
			mcp.Required(),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to custom policies"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "xml", "junit-xml", "sarif"),
		),
	)
	s.AddTool(scanWithPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "terrascan", directory, "--policy-path", policyPath}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--output", format)
		}
		return executeShipCommand(args)
	})

	// Terrascan get version tool
	getVersionTool := mcp.NewTool("terrascan_get_version",
		mcp.WithDescription("Get Terrascan version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "terrascan", "--version"}
		return executeShipCommand(args)
	})
}