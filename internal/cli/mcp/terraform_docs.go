package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddTerraformDocsTools adds Terraform Docs MCP tool implementations using direct Dagger calls
func AddTerraformDocsTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTerraformDocsToolsDirect(s)
}

// addTerraformDocsToolsDirect adds Terraform Docs tools using direct Dagger module calls
func addTerraformDocsToolsDirect(s *server.MCPServer) {
	// Terraform docs generate markdown tool
	generateMarkdownTool := mcp.NewTool("terraform_docs_markdown",
		mcp.WithDescription("Generate Terraform documentation in Markdown format using real terraform-docs CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
	)
	s.AddTool(generateMarkdownTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerraformDocsModule(client)

		// Get parameters
		directory := request.GetString("directory", "")

		// Generate markdown documentation
		output, err := module.GenerateMarkdown(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraform docs markdown generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terraform docs generate JSON tool
	generateJSONTool := mcp.NewTool("terraform_docs_json",
		mcp.WithDescription("Generate Terraform documentation in JSON format using real terraform-docs CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
	)
	s.AddTool(generateJSONTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerraformDocsModule(client)

		// Get parameters
		directory := request.GetString("directory", "")

		// Generate JSON documentation
		output, err := module.GenerateJSON(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraform docs JSON generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terraform docs generate table tool
	generateTableTool := mcp.NewTool("terraform_docs_table",
		mcp.WithDescription("Generate Terraform documentation as table using real terraform-docs CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
	)
	s.AddTool(generateTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerraformDocsModule(client)

		// Get parameters
		directory := request.GetString("directory", "")

		// Generate table documentation
		output, err := module.GenerateTable(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraform docs table generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terraform docs generate with config tool
	generateWithConfigTool := mcp.NewTool("terraform_docs_with_config",
		mcp.WithDescription("Generate Terraform documentation with custom config using real terraform-docs CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("config_file",
			mcp.Description("Configuration file path"),
			mcp.Required(),
		),
	)
	s.AddTool(generateWithConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerraformDocsModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		configFile := request.GetString("config_file", "")

		// Generate documentation with config
		output, err := module.GenerateWithConfig(ctx, directory, configFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraform docs config generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terraform docs version tool (using extra function from Dagger module)
	versionTool := mcp.NewTool("terraform_docs_version",
		mcp.WithDescription("Get terraform-docs version information using real terraform-docs CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerraformDocsModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraform docs get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}