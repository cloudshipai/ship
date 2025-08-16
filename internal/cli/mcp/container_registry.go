package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddContainerRegistryTools adds container registry MCP tool implementations
func AddContainerRegistryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Registry scan tool
	scanRegistryTool := mcp.NewTool("registry_scan",
		mcp.WithDescription("Scan container registry for vulnerabilities and compliance"),
		mcp.WithString("registry_url",
			mcp.Description("Container registry URL to scan"),
			mcp.Required(),
		),
		mcp.WithString("credentials",
			mcp.Description("Registry credentials (optional)"),
		),
	)
	s.AddTool(scanRegistryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		registryURL := request.GetString("registry_url", "")
		args := []string{"security", "registry", "scan", registryURL}
		if creds := request.GetString("credentials", ""); creds != "" {
			args = append(args, "--credentials", creds)
		}
		return executeShipCommand(args)
	})

	// Registry list images tool
	listImagesTool := mcp.NewTool("registry_list_images",
		mcp.WithDescription("List images in container registry"),
		mcp.WithString("registry_url",
			mcp.Description("Container registry URL"),
			mcp.Required(),
		),
		mcp.WithString("repository",
			mcp.Description("Repository name (optional)"),
		),
	)
	s.AddTool(listImagesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		registryURL := request.GetString("registry_url", "")
		args := []string{"cloud", "registry", "list", registryURL}
		if repo := request.GetString("repository", ""); repo != "" {
			args = append(args, "--repository", repo)
		}
		return executeShipCommand(args)
	})

	// Registry cleanup tool
	cleanupTool := mcp.NewTool("registry_cleanup",
		mcp.WithDescription("Clean up unused images and tags in container registry"),
		mcp.WithString("registry_url",
			mcp.Description("Container registry URL"),
			mcp.Required(),
		),
		mcp.WithString("age",
			mcp.Description("Age threshold for cleanup (e.g., 30d, 7d)"),
		),
	)
	s.AddTool(cleanupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		registryURL := request.GetString("registry_url", "")
		args := []string{"cloud", "registry", "cleanup", registryURL}
		if age := request.GetString("age", ""); age != "" {
			args = append(args, "--age", age)
		}
		return executeShipCommand(args)
	})

	// Registry mirror tool
	mirrorTool := mcp.NewTool("registry_mirror",
		mcp.WithDescription("Mirror images between container registries"),
		mcp.WithString("source_registry",
			mcp.Description("Source registry URL"),
			mcp.Required(),
		),
		mcp.WithString("target_registry",
			mcp.Description("Target registry URL"),
			mcp.Required(),
		),
		mcp.WithString("image_filter",
			mcp.Description("Image filter pattern (optional)"),
		),
	)
	s.AddTool(mirrorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sourceRegistry := request.GetString("source_registry", "")
		targetRegistry := request.GetString("target_registry", "")
		args := []string{"cloud", "registry", "mirror", sourceRegistry, targetRegistry}
		if filter := request.GetString("image_filter", ""); filter != "" {
			args = append(args, "--filter", filter)
		}
		return executeShipCommand(args)
	})

	// Registry get version tool
	getVersionTool := mcp.NewTool("registry_get_version",
		mcp.WithDescription("Get container registry tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloud", "registry", "--version"}
		return executeShipCommand(args)
	})
}