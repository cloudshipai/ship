package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddDockleTools adds Dockle (container image linter) MCP tool implementations
func AddDockleTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Dockle scan image tool
	scanImageTool := mcp.NewTool("dockle_scan_image",
		mcp.WithDescription("Scan container image for security and best practices using Dockle"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to scan"),
			mcp.Required(),
		),
	)
	s.AddTool(scanImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"security", "dockle", "--image", imageRef}
		return executeShipCommand(args)
	})

	// Dockle scan tarball tool
	scanTarballTool := mcp.NewTool("dockle_scan_tarball",
		mcp.WithDescription("Scan container image tarball using Dockle"),
		mcp.WithString("tarball_path",
			mcp.Description("Path to container image tarball"),
			mcp.Required(),
		),
	)
	s.AddTool(scanTarballTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tarballPath := request.GetString("tarball_path", "")
		args := []string{"security", "dockle", "--tarball", tarballPath}
		return executeShipCommand(args)
	})

	// Dockle scan dockerfile tool
	scanDockerfileTool := mcp.NewTool("dockle_scan_dockerfile",
		mcp.WithDescription("Scan Dockerfile for best practices using Dockle"),
		mcp.WithString("dockerfile_path",
			mcp.Description("Path to Dockerfile"),
			mcp.Required(),
		),
	)
	s.AddTool(scanDockerfileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dockerfilePath := request.GetString("dockerfile_path", "")
		args := []string{"security", "dockle", "--dockerfile", dockerfilePath}
		return executeShipCommand(args)
	})

	// Dockle generate config tool
	generateConfigTool := mcp.NewTool("dockle_generate_config",
		mcp.WithDescription("Generate Dockle configuration template"),
		mcp.WithString("output_path",
			mcp.Description("Output path for configuration file"),
		),
	)
	s.AddTool(generateConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "dockle", "--generate-config"}
		if outputPath := request.GetString("output_path", ""); outputPath != "" {
			args = append(args, "--output", outputPath)
		}
		return executeShipCommand(args)
	})

	// Dockle list checks tool
	listChecksTool := mcp.NewTool("dockle_list_checks",
		mcp.WithDescription("List all available Dockle security checks"),
	)
	s.AddTool(listChecksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "dockle", "--list-checks"}
		return executeShipCommand(args)
	})

	// Dockle scan with policy tool
	scanWithPolicyTool := mcp.NewTool("dockle_scan_with_policy",
		mcp.WithDescription("Scan container image with custom policy using Dockle"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to scan"),
			mcp.Required(),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to custom policy file"),
			mcp.Required(),
		),
	)
	s.AddTool(scanWithPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "dockle", "--image", imageRef, "--policy", policyPath}
		return executeShipCommand(args)
	})
}