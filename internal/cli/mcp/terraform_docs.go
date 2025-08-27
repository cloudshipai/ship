package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddTerraformDocsTools adds terraform-docs (documentation generator) MCP tool implementations using direct Dagger calls
func AddTerraformDocsTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTerraformDocsToolsDirect(s)
}

// addTerraformDocsToolsDirect adds terraform-docs tools using direct Dagger module calls
func addTerraformDocsToolsDirect(s *server.MCPServer) {
	// Terraform-docs generate tool
	generateTool := mcp.NewTool("terraform_docs_generate",
		mcp.WithDescription("Generate Terraform module documentation"),
		mcp.WithString("module_path",
			mcp.Description("Path to Terraform module directory"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format for documentation"),
			mcp.Enum("markdown", "json", "yaml", "xml", "adoc", "pretty", "table"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path (optional, prints to stdout if not specified)"),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to terraform-docs configuration file"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Generate documentation for all modules recursively"),
		),
		mcp.WithBoolean("sort",
			mcp.Description("Sort items"),
		),
		mcp.WithString("header_from",
			mcp.Description("Path to file to use as header"),
		),
		mcp.WithString("footer_from",
			mcp.Description("Path to file to use as footer"),
		),
	)
	s.AddTool(generateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		modulePath := request.GetString("module_path", "")
		outputFormat := request.GetString("output_format", "")
		outputFile := request.GetString("output_file", "")
		configFile := request.GetString("config_file", "")
		recursive := request.GetBool("recursive", false)
		sort := request.GetBool("sort", false)
		headerFrom := request.GetString("header_from", "")
		footerFrom := request.GetString("footer_from", "")

		if modulePath == "" {
			return mcp.NewToolResultError("module_path is required"), nil
		}

		// Set default format
		if outputFormat == "" {
			outputFormat = "markdown"
		}

		// Create terraform-docs module
		terraformDocsModule := modules.NewTerraformDocsModule(client)

		// Set up options
		opts := modules.TerraformDocsOptions{
			OutputFormat: outputFormat,
			OutputFile:   outputFile,
			ConfigFile:   configFile,
			Recursive:    recursive,
			Sort:         sort,
			HeaderFrom:   headerFrom,
			FooterFrom:   footerFrom,
		}

		// Generate documentation
		result, err := terraformDocsModule.Generate(ctx, modulePath, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("terraform-docs generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Terraform-docs validate tool
	validateTool := mcp.NewTool("terraform_docs_validate",
		mcp.WithDescription("Validate that Terraform module documentation is up to date"),
		mcp.WithString("module_path",
			mcp.Description("Path to Terraform module directory"),
			mcp.Required(),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to terraform-docs configuration file"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Validate all modules recursively"),
		),
	)
	s.AddTool(validateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		modulePath := request.GetString("module_path", "")
		configFile := request.GetString("config_file", "")
		recursive := request.GetBool("recursive", false)

		if modulePath == "" {
			return mcp.NewToolResultError("module_path is required"), nil
		}

		// Create terraform-docs module
		terraformDocsModule := modules.NewTerraformDocsModule(client)

		// Set up options
		opts := modules.TerraformDocsValidateOptions{
			ConfigFile: configFile,
			Recursive:  recursive,
		}

		// Validate documentation
		result, err := terraformDocsModule.Validate(ctx, modulePath, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("terraform-docs validation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}