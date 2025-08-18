package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddContainerRegistryTools adds container registry operations using Docker CLI
func AddContainerRegistryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Docker login to registry
	loginTool := mcp.NewTool("docker_login",
		mcp.WithDescription("Login to container registry using Docker"),
		mcp.WithString("registry",
			mcp.Description("Registry URL (e.g., docker.io, ghcr.io)"),
		),
		mcp.WithString("username",
			mcp.Description("Registry username"),
		),
		mcp.WithString("password",
			mcp.Description("Registry password or token"),
		),
	)
	s.AddTool(loginTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"docker", "login"}
		
		if registry := request.GetString("registry", ""); registry != "" {
			args = append(args, registry)
		}
		if username := request.GetString("username", ""); username != "" {
			args = append(args, "--username", username)
		}
		if password := request.GetString("password", ""); password != "" {
			args = append(args, "--password", password)
		}
		
		return executeShipCommand(args)
	})

	// Docker push image
	pushTool := mcp.NewTool("docker_push",
		mcp.WithDescription("Push image to container registry"),
		mcp.WithString("image",
			mcp.Description("Image name and tag to push"),
			mcp.Required(),
		),
	)
	s.AddTool(pushTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		image := request.GetString("image", "")
		args := []string{"docker", "push", image}
		return executeShipCommand(args)
	})

	// Docker pull image
	pullTool := mcp.NewTool("docker_pull",
		mcp.WithDescription("Pull image from container registry"),
		mcp.WithString("image",
			mcp.Description("Image name and tag to pull"),
			mcp.Required(),
		),
	)
	s.AddTool(pullTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		image := request.GetString("image", "")
		args := []string{"docker", "pull", image}
		return executeShipCommand(args)
	})

	// Docker images list
	imagesTool := mcp.NewTool("docker_images",
		mcp.WithDescription("List local Docker images"),
		mcp.WithString("repository",
			mcp.Description("Repository name filter"),
		),
		mcp.WithBoolean("all",
			mcp.Description("Show all images (including intermediate)"),
		),
	)
	s.AddTool(imagesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"docker", "images"}
		
		if repository := request.GetString("repository", ""); repository != "" {
			args = append(args, repository)
		}
		if request.GetBool("all", false) {
			args = append(args, "--all")
		}
		
		return executeShipCommand(args)
	})

	// Docker tag image
	tagTool := mcp.NewTool("docker_tag",
		mcp.WithDescription("Create a tag for an image"),
		mcp.WithString("source_image",
			mcp.Description("Source image name and tag"),
			mcp.Required(),
		),
		mcp.WithString("target_image",
			mcp.Description("Target image name and tag"),
			mcp.Required(),
		),
	)
	s.AddTool(tagTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sourceImage := request.GetString("source_image", "")
		targetImage := request.GetString("target_image", "")
		args := []string{"docker", "tag", sourceImage, targetImage}
		return executeShipCommand(args)
	})

	// Docker logout
	logoutTool := mcp.NewTool("docker_logout",
		mcp.WithDescription("Logout from container registry"),
		mcp.WithString("registry",
			mcp.Description("Registry URL to logout from"),
		),
	)
	s.AddTool(logoutTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"docker", "logout"}
		
		if registry := request.GetString("registry", ""); registry != "" {
			args = append(args, registry)
		}
		
		return executeShipCommand(args)
	})
}