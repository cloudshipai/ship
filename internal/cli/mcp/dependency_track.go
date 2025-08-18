package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddDependencyTrackTools adds Dependency Track (software component analysis) MCP tool implementations
func AddDependencyTrackTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addDependencyTrackToolsDirect(s)
}

// addDependencyTrackToolsDirect implements direct Dagger calls for Dependency Track tools
func addDependencyTrackToolsDirect(s *server.MCPServer) {
	// Dependency Track upload BOM tool using dtrack-cli
	uploadBOMTool := mcp.NewTool("dependency_track_upload_bom",
		mcp.WithDescription("Upload Software Bill of Materials to Dependency Track using dtrack-cli"),
		mcp.WithString("bom_path",
			mcp.Description("Path to BOM file (CycloneDX or SPDX format)"),
			mcp.Required(),
		),
		mcp.WithString("project_name",
			mcp.Description("Project name in Dependency Track"),
			mcp.Required(),
		),
		mcp.WithString("project_version",
			mcp.Description("Project version"),
		),
		mcp.WithString("server_url",
			mcp.Description("Dependency Track server URL"),
		),
		mcp.WithString("api_key",
			mcp.Description("API key for authentication"),
		),
	)
	s.AddTool(uploadBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		bomPath := request.GetString("bom_path", "")
		projectName := request.GetString("project_name", "")
		projectVersion := request.GetString("project_version", "")
		serverUrl := request.GetString("server_url", "")
		apiKey := request.GetString("api_key", "")

		// Create Dependency Track module and upload BOM
		dependencyTrackModule := modules.NewDependencyTrackModule(client)
		result, err := dependencyTrackModule.UploadBOM(ctx, bomPath, projectName, projectVersion, serverUrl, apiKey)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("dependency track upload bom failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Dependency Track upload BOM via API (alternative using curl)
	uploadBOMApiTool := mcp.NewTool("dependency_track_upload_bom_api",
		mcp.WithDescription("Upload BOM to Dependency Track via REST API using curl"),
		mcp.WithString("bom_path",
			mcp.Description("Path to BOM file (CycloneDX or SPDX format)"),
			mcp.Required(),
		),
		mcp.WithString("server_url",
			mcp.Description("Dependency Track server URL"),
			mcp.Required(),
		),
		mcp.WithString("api_key",
			mcp.Description("API key for authentication"),
			mcp.Required(),
		),
		mcp.WithString("project_name",
			mcp.Description("Project name in Dependency Track"),
		),
		mcp.WithString("project_version",
			mcp.Description("Project version"),
		),
		mcp.WithBoolean("auto_create",
			mcp.Description("Auto-create project if it doesn't exist"),
		),
	)
	s.AddTool(uploadBOMApiTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		bomPath := request.GetString("bom_path", "")
		serverUrl := request.GetString("server_url", "")
		apiKey := request.GetString("api_key", "")
		projectName := request.GetString("project_name", "")
		projectVersion := request.GetString("project_version", "")
		autoCreate := request.GetBool("auto_create", false)

		// Create Dependency Track module and upload BOM via API
		dependencyTrackModule := modules.NewDependencyTrackModule(client)
		result, err := dependencyTrackModule.UploadBOMAPI(ctx, bomPath, serverUrl, apiKey, projectName, projectVersion, autoCreate)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("dependency track upload bom api failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Generate CycloneDX BOM using CycloneDX CLI tools
	generateBOMTool := mcp.NewTool("dependency_track_generate_bom",
		mcp.WithDescription("Generate CycloneDX BOM using cyclonedx-cli tools"),
		mcp.WithString("project_type",
			mcp.Description("Project type"),
			mcp.Enum("npm", "maven", "gradle", "pip", "composer", "dotnet"),
			mcp.Required(),
		),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output BOM file path"),
		),
	)
	s.AddTool(generateBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		projectType := request.GetString("project_type", "")
		projectPath := request.GetString("project_path", ".")
		outputFile := request.GetString("output_file", "bom.json")

		// Create Dependency Track module and generate BOM
		dependencyTrackModule := modules.NewDependencyTrackModule(client)
		result, err := dependencyTrackModule.GenerateBOM(ctx, projectType, projectPath, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("dependency track generate bom failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}