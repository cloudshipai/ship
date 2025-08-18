package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddDependencyTrackTools adds Dependency Track (software component analysis) MCP tool implementations
func AddDependencyTrackTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		bomPath := request.GetString("bom_path", "")
		projectName := request.GetString("project_name", "")
		args := []string{"dtrack-cli", "--bom-path", bomPath, "--project-name", projectName}
		
		if projectVersion := request.GetString("project_version", ""); projectVersion != "" {
			args = append(args, "--project-version", projectVersion)
		}
		if serverUrl := request.GetString("server_url", ""); serverUrl != "" {
			args = append(args, "--server", serverUrl)
		}
		if apiKey := request.GetString("api_key", ""); apiKey != "" {
			args = append(args, "--api-key", apiKey)
		}
		
		// Add auto-create by default for easier usage
		args = append(args, "--auto-create", "true")
		return executeShipCommand(args)
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
		bomPath := request.GetString("bom_path", "")
		serverUrl := request.GetString("server_url", "")
		apiKey := request.GetString("api_key", "")
		
		args := []string{"curl", "-X", "POST", serverUrl + "/api/v1/bom",
			"-H", "Content-Type: multipart/form-data",
			"-H", "X-Api-Key: " + apiKey,
			"-F", "bom=@" + bomPath}
		
		if projectName := request.GetString("project_name", ""); projectName != "" {
			args = append(args, "-F", "projectName=" + projectName)
		}
		if projectVersion := request.GetString("project_version", ""); projectVersion != "" {
			args = append(args, "-F", "projectVersion=" + projectVersion)
		}
		if request.GetBool("auto_create", false) {
			args = append(args, "-F", "autoCreate=true")
		}
		
		return executeShipCommand(args)
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
		projectType := request.GetString("project_type", "")
		
		var args []string
		switch projectType {
		case "npm":
			args = []string{"cyclonedx-npm"}
		case "maven":
			args = []string{"mvn", "org.cyclonedx:cyclonedx-maven-plugin:makeBom"}
		case "gradle":
			args = []string{"gradle", "cyclonedxBom"}
		case "pip":
			args = []string{"cyclonedx-py"}
		case "composer":
			args = []string{"cyclonedx-php", "composer"}
		case "dotnet":
			args = []string{"cyclonedx", "dotnet"}
		default:
			args = []string{"cyclonedx-npm"}
		}
		
		if projectPath := request.GetString("project_path", ""); projectPath != "" && projectType == "npm" {
			args = append(args, "-o", projectPath)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" && projectType == "npm" {
			args = append(args, "-o", outputFile)
		}
		
		return executeShipCommand(args)
	})
}