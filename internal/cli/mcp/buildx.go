package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddBuildXTools adds Docker BuildX tools to the MCP server
func AddBuildXTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// BuildX Build - Build an OCI image using Docker BuildX
	buildxBuildTool := mcp.NewTool("buildx_build",
		mcp.WithDescription("Build an OCI image using Docker BuildX with multi-platform support"),
		mcp.WithString("src_dir", mcp.Required(), mcp.Description("Source directory path containing Dockerfile")),
		mcp.WithString("tag", mcp.Required(), mcp.Description("Image tag (e.g., myapp:latest)")),
		mcp.WithString("platform", mcp.Description("Target platform(s) (default: linux/amd64, can be comma-separated like linux/amd64,linux/arm64)")),
		mcp.WithString("dockerfile_path", mcp.Description("Path to Dockerfile relative to src_dir (default: .)")),
	)

	s.AddTool(buildxBuildTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		srcDir := request.GetString("src_dir", "")
		tag := request.GetString("tag", "")
		platform := request.GetString("platform", "linux/amd64")
		dockerfilePath := request.GetString("dockerfile_path", ".")

		args := []string{
			"buildx", "build",
			"--src-dir", srcDir,
			"--tag", tag,
			"--platform", platform,
			"--dockerfile-path", dockerfilePath,
		}

		return executeShipCommand(args)
	})

	// BuildX Publish - Build and publish an OCI image to registry
	buildxPublishTool := mcp.NewTool("buildx_publish",
		mcp.WithDescription("Build and publish an OCI image to a container registry using Docker BuildX"),
		mcp.WithString("src_dir", mcp.Required(), mcp.Description("Source directory path containing Dockerfile")),
		mcp.WithString("tag", mcp.Required(), mcp.Description("Image tag to publish (e.g., myregistry.com/myapp:latest)")),
		mcp.WithString("platform", mcp.Description("Target platform(s) (default: linux/amd64, can be comma-separated)")),
		mcp.WithString("username", mcp.Required(), mcp.Description("Registry username")),
		mcp.WithString("password", mcp.Required(), mcp.Description("Registry password or token")),
		mcp.WithString("registry", mcp.Description("Container registry URL (default: docker.io)")),
		mcp.WithString("dockerfile_path", mcp.Description("Path to Dockerfile relative to src_dir (default: .)")),
	)

	s.AddTool(buildxPublishTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		srcDir := request.GetString("src_dir", "")
		tag := request.GetString("tag", "")
		platform := request.GetString("platform", "linux/amd64")
		username := request.GetString("username", "")
		password := request.GetString("password", "")
		registry := request.GetString("registry", "docker.io")
		dockerfilePath := request.GetString("dockerfile_path", ".")

		args := []string{
			"buildx", "publish",
			"--src-dir", srcDir,
			"--tag", tag,
			"--platform", platform,
			"--username", username,
			"--password", password,
			"--registry", registry,
			"--dockerfile-path", dockerfilePath,
		}

		return executeShipCommand(args)
	})

	// BuildX Dev - Get a development environment with BuildX installed
	buildxDevTool := mcp.NewTool("buildx_dev",
		mcp.WithDescription("Set up a development environment container with Docker BuildX installed"),
		mcp.WithString("src_dir", mcp.Description("Source directory to mount in the dev environment (optional)")),
	)

	s.AddTool(buildxDevTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		srcDir := request.GetString("src_dir", "")

		args := []string{"buildx", "dev"}
		if srcDir != "" {
			args = append(args, "--src-dir", srcDir)
		}

		return executeShipCommand(args)
	})

	// BuildX Version - Get BuildX version information
	buildxVersionTool := mcp.NewTool("buildx_version",
		mcp.WithDescription("Get Docker BuildX version information"),
	)

	s.AddTool(buildxVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"buildx", "version"}
		return executeShipCommand(args)
	})
}