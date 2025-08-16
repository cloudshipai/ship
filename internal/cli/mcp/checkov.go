package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCheckovTools adds Checkov MCP tool implementations
func AddCheckovTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Checkov scan directory tool
	scanDirTool := mcp.NewTool("checkov_scan_directory",
		mcp.WithDescription("Scan a directory for security issues using Checkov"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
	)
	s.AddTool(scanDirTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// Checkov scan file tool
	scanFileTool := mcp.NewTool("checkov_scan_file",
		mcp.WithDescription("Scan a specific file for security issues using Checkov"),
		mcp.WithString("file_path",
			mcp.Description("Path to the file to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"tf", "checkov", "--file", filePath}
		return executeShipCommand(args)
	})

	// Checkov scan with policy tool
	scanWithPolicyTool := mcp.NewTool("checkov_scan_with_policy",
		mcp.WithDescription("Scan with custom policy using Checkov"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to custom policy file"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		policyPath := request.GetString("policy_path", "")
		args = append(args, "--policy", policyPath)
		return executeShipCommand(args)
	})

	// Checkov scan multi-framework tool
	scanMultiFrameworkTool := mcp.NewTool("checkov_scan_multi_framework",
		mcp.WithDescription("Scan with multiple frameworks using Checkov"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("frameworks",
			mcp.Description("Comma-separated list of frameworks"),
			mcp.Required(),
			mcp.Enum("terraform", "cloudformation", "kubernetes", "dockerfile", "arm"),
		),
	)
	s.AddTool(scanMultiFrameworkTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		frameworks := request.GetString("frameworks", "")
		args = append(args, "--framework", frameworks)
		return executeShipCommand(args)
	})

	// Checkov scan with severity tool
	scanWithSeverityTool := mcp.NewTool("checkov_scan_with_severity",
		mcp.WithDescription("Scan with severity threshold using Checkov"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("severities",
			mcp.Description("Comma-separated list of severities"),
			mcp.Required(),
			mcp.Enum("LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
	)
	s.AddTool(scanWithSeverityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		severities := request.GetString("severities", "")
		args = append(args, "--severity", severities)
		return executeShipCommand(args)
	})

	// Checkov scan with skips tool
	scanWithSkipsTool := mcp.NewTool("checkov_scan_with_skips",
		mcp.WithDescription("Scan with skipped checks using Checkov"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("skip_checks",
			mcp.Description("Comma-separated list of check IDs to skip"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithSkipsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		skipChecks := request.GetString("skip_checks", "")
		args = append(args, "--skip-check", skipChecks)
		return executeShipCommand(args)
	})

	// Checkov scan container image tool
	scanContainerImageTool := mcp.NewTool("checkov_scan_container_image",
		mcp.WithDescription("Scan container image for security issues using Checkov"),
		mcp.WithString("image_name",
			mcp.Description("Container image name to scan"),
			mcp.Required(),
		),
		mcp.WithString("dockerfile_path",
			mcp.Description("Path to Dockerfile for additional context"),
		),
	)
	s.AddTool(scanContainerImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"tf", "checkov", "--docker-image", imageName}
		if dockerfilePath := request.GetString("dockerfile_path", ""); dockerfilePath != "" {
			args = append(args, "--dockerfile-path", dockerfilePath)
		}
		return executeShipCommand(args)
	})

	// Checkov scan SCA packages tool
	scanSCAPackagesTool := mcp.NewTool("checkov_scan_sca_packages",
		mcp.WithDescription("Scan for software composition analysis using Checkov"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan for package files (default: current directory)"),
		),
	)
	s.AddTool(scanSCAPackagesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov", "--framework", "sca_package"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, "--directory", dir)
		}
		return executeShipCommand(args)
	})

	// Checkov scan secrets tool
	scanSecretsTool := mcp.NewTool("checkov_scan_secrets",
		mcp.WithDescription("Scan for secrets in code using Checkov"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan for secrets (default: current directory)"),
		),
	)
	s.AddTool(scanSecretsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov", "--framework", "secrets"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, "--directory", dir)
		}
		return executeShipCommand(args)
	})

	// Checkov scan with specific checks tool
	scanWithSpecificChecksTool := mcp.NewTool("checkov_scan_with_specific_checks",
		mcp.WithDescription("Scan with only specific checks enabled using Checkov"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("checks",
			mcp.Description("Comma-separated list of check IDs to run"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithSpecificChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, "--directory", dir)
		}
		checks := request.GetString("checks", "")
		args = append(args, "--check", checks)
		return executeShipCommand(args)
	})

	// Checkov get version tool
	getVersionTool := mcp.NewTool("checkov_get_version",
		mcp.WithDescription("Get Checkov version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov", "--version"}
		return executeShipCommand(args)
	})
}