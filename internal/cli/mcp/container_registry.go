package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddContainerRegistryTools adds container registry operations using Docker CLI
func AddContainerRegistryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addContainerRegistryToolsDirect(s)
}

// addContainerRegistryToolsDirect implements direct Dagger calls for container registry tools
func addContainerRegistryToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		registry := request.GetString("registry", "")
		username := request.GetString("username", "")
		password := request.GetString("password", "")

		// Create container registry module and login
		registryModule := modules.NewContainerRegistryModule(client)
		result, err := registryModule.Login(ctx, registry, username, password)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("docker login failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		image := request.GetString("image", "")

		// Create container registry module and push image
		registryModule := modules.NewContainerRegistryModule(client)
		result, err := registryModule.PushImage(ctx, image)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("docker push failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		image := request.GetString("image", "")

		// Create container registry module and pull image
		registryModule := modules.NewContainerRegistryModule(client)
		result, err := registryModule.PullImage(ctx, image)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("docker pull failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		repository := request.GetString("repository", "")
		all := request.GetBool("all", false)

		// Create container registry module and list images
		registryModule := modules.NewContainerRegistryModule(client)
		result, err := registryModule.ListImages(ctx, repository, all)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("docker images failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		sourceImage := request.GetString("source_image", "")
		targetImage := request.GetString("target_image", "")

		// Create container registry module and tag image
		registryModule := modules.NewContainerRegistryModule(client)
		result, err := registryModule.TagImage(ctx, sourceImage, targetImage)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("docker tag failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Docker logout
	logoutTool := mcp.NewTool("docker_logout",
		mcp.WithDescription("Logout from container registry"),
		mcp.WithString("registry",
			mcp.Description("Registry URL to logout from"),
		),
	)
	s.AddTool(logoutTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		registry := request.GetString("registry", "")

		// Create container registry module and logout
		registryModule := modules.NewContainerRegistryModule(client)
		result, err := registryModule.Logout(ctx, registry)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("docker logout failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}