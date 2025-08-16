package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTerraformDocsTools adds Terraform Docs MCP tool implementations
func AddTerraformDocsTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Terraform docs generate markdown tool
	generateMarkdownTool := mcp.NewTool("terraform_docs_generate_markdown",
		mcp.WithDescription("Generate Terraform documentation in Markdown format"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for documentation"),
		),
	)
	s.AddTool(generateMarkdownTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs", "--format", "markdown"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// Terraform docs generate JSON tool
	generateJSONTool := mcp.NewTool("terraform_docs_generate_json",
		mcp.WithDescription("Generate Terraform documentation in JSON format"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for documentation"),
		),
	)
	s.AddTool(generateJSONTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs", "--format", "json"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// Terraform docs generate with config tool
	generateWithConfigTool := mcp.NewTool("terraform_docs_generate_config",
		mcp.WithDescription("Generate Terraform documentation using configuration file"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to terraform-docs configuration file"),
			mcp.Required(),
		),
	)
	s.AddTool(generateWithConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs", "--config", request.GetString("config_file", "")}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// Terraform docs generate table tool
	generateTableTool := mcp.NewTool("terraform_docs_generate_table",
		mcp.WithDescription("Generate Terraform documentation in table format"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for documentation"),
		),
	)
	s.AddTool(generateTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs", "--format", "table"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// Terraform docs inject documentation tool
	injectDocumentationTool := mcp.NewTool("terraform_docs_inject_documentation",
		mcp.WithDescription("Inject Terraform documentation into existing files"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output_file",
			mcp.Description("File to inject documentation into (e.g., README.md)"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Documentation format"),
			mcp.Enum("markdown", "asciidoc", "json"),
		),
	)
	s.AddTool(injectDocumentationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs", "--output-mode", "inject"}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--output-format", format)
		} else {
			args = append(args, "--output-format", "markdown")
		}
		outputFile := request.GetString("output_file", "")
		args = append(args, "--output-file", outputFile)
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// Terraform docs generate AsciiDoc tool
	generateAsciiDocTool := mcp.NewTool("terraform_docs_generate_asciidoc",
		mcp.WithDescription("Generate Terraform documentation in AsciiDoc format"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for documentation"),
		),
	)
	s.AddTool(generateAsciiDocTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs", "--format", "asciidoc"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// Terraform docs generate markdown table tool
	generateMarkdownTableTool := mcp.NewTool("terraform_docs_generate_markdown_table",
		mcp.WithDescription("Generate Terraform documentation as markdown table"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for documentation"),
		),
	)
	s.AddTool(generateMarkdownTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs", "markdown", "table"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		return executeShipCommand(args)
	})

	// Terraform docs recursive generation tool
	generateRecursiveTool := mcp.NewTool("terraform_docs_generate_recursive",
		mcp.WithDescription("Generate Terraform documentation recursively for all modules"),
		mcp.WithString("directory",
			mcp.Description("Root directory containing Terraform modules (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Documentation format"),
			mcp.Enum("markdown", "asciidoc", "json", "table"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Enable recursive mode for all subdirectories"),
		),
	)
	s.AddTool(generateRecursiveTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		format := request.GetString("format", "markdown")
		args := []string{"tf", "docs", "--format", format}
		if request.GetBool("recursive", false) {
			args = append(args, "--recursive")
		}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		return executeShipCommand(args)
	})

	// Terraform docs get version tool
	getVersionTool := mcp.NewTool("terraform_docs_get_version",
		mcp.WithDescription("Get terraform-docs version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs", "--version"}
		return executeShipCommand(args)
	})
}