package mcp

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSnykTools adds Snyk vulnerability scanning tools to the MCP server
func AddSnykTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addSnykToolsDirect(s)
}

// addSnykToolsDirect adds Snyk tools using direct Dagger module calls
func addSnykToolsDirect(s *server.MCPServer) {
	// Test project
	testProjectTool := mcp.NewTool("snyk_test_project",
		mcp.WithDescription("Test a project for vulnerabilities"),
		mcp.WithString("directory",
			mcp.Description("Directory containing project files"),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity threshold"),
		),
	)
	s.AddTool(testProjectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		severity := request.GetString("severity", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewSnykModule(client)
		result, err := module.TestProject(ctx, dir, severity)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("test failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Test container
	testContainerTool := mcp.NewTool("snyk_test_container",
		mcp.WithDescription("Test a container image for vulnerabilities"),
		mcp.WithString("image",
			mcp.Description("Container image name"),
			mcp.Required(),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity threshold"),
		),
	)
	s.AddTool(testContainerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		image := request.GetString("image", "")
		severity := request.GetString("severity", "")

		if image == "" {
			return mcp.NewToolResultError("image is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewSnykModule(client)
		result, err := module.TestContainer(ctx, image, severity)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("test failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Test IaC
	testIaCTool := mcp.NewTool("snyk_test_iac",
		mcp.WithDescription("Test Infrastructure as Code files for security issues"),
		mcp.WithString("directory",
			mcp.Description("Directory containing IaC files"),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity threshold"),
		),
	)
	s.AddTool(testIaCTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		severity := request.GetString("severity", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewSnykModule(client)
		result, err := module.TestIaC(ctx, dir, severity)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("test failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Test code
	testCodeTool := mcp.NewTool("snyk_test_code",
		mcp.WithDescription("Test source code for security vulnerabilities (SAST)"),
		mcp.WithString("directory",
			mcp.Description("Directory containing source code"),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity threshold"),
		),
	)
	s.AddTool(testCodeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		severity := request.GetString("severity", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewSnykModule(client)
		result, err := module.TestCode(ctx, dir, severity)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("test failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Generate SBOM
	generateSBOMTool := mcp.NewTool("snyk_generate_sbom",
		mcp.WithDescription("Generate Software Bill of Materials"),
		mcp.WithString("directory",
			mcp.Description("Directory containing project"),
		),
		mcp.WithString("format",
			mcp.Description("SBOM format"),
		),
	)
	s.AddTool(generateSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dir := request.GetString("directory", ".")
		format := request.GetString("format", "cyclonedx1.4+json")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewSnykModule(client)
		result, err := module.GenerateSBOM(ctx, dir, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SBOM generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Get version
	versionTool := mcp.NewTool("snyk_get_version",
		mcp.WithDescription("Get Snyk CLI version"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewSnykModule(client)
		result, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get version: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}