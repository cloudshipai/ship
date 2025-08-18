package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTerraformDocsTools adds Terraform Docs MCP tool implementations using real terraform-docs CLI
func AddTerraformDocsTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Terraform docs generate markdown tool
	generateMarkdownTool := mcp.NewTool("terraform_docs_markdown",
		mcp.WithDescription("Generate Terraform documentation in Markdown format using real terraform-docs CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for documentation"),
		),
		mcp.WithString("output_mode",
			mcp.Description("Output mode (inject, replace)"),
			mcp.Enum("inject", "replace"),
		),
	)
	s.AddTool(generateMarkdownTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-docs", "markdown"}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		if outputMode := request.GetString("output_mode", ""); outputMode != "" {
			args = append(args, "--output-mode", outputMode)
		}
		
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		} else {
			args = append(args, ".")
		}
		
		return executeShipCommand(args)
	})

	// Terraform docs markdown table tool
	markdownTableTool := mcp.NewTool("terraform_docs_markdown_table",
		mcp.WithDescription("Generate Terraform documentation as markdown table using real terraform-docs CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for documentation"),
		),
		mcp.WithString("output_mode",
			mcp.Description("Output mode (inject, replace)"),
			mcp.Enum("inject", "replace"),
		),
	)
	s.AddTool(markdownTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-docs", "markdown", "table"}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		if outputMode := request.GetString("output_mode", ""); outputMode != "" {
			args = append(args, "--output-mode", outputMode)
		}
		
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		} else {
			args = append(args, ".")
		}
		
		return executeShipCommand(args)
	})

	// Terraform docs with config tool
	generateWithConfigTool := mcp.NewTool("terraform_docs_with_config",
		mcp.WithDescription("Generate Terraform documentation using configuration file with real terraform-docs CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to .terraform-docs.yml configuration file"),
		),
	)
	s.AddTool(generateWithConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-docs"}
		
		if configFile := request.GetString("config_file", ""); configFile != "" {
			args = append(args, "--config", configFile)
		}
		
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		} else {
			args = append(args, ".")
		}
		
		return executeShipCommand(args)
	})






	// Terraform docs version tool
	getVersionTool := mcp.NewTool("terraform_docs_version",
		mcp.WithDescription("Get terraform-docs version information using real terraform-docs CLI"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-docs", "--version"}
		return executeShipCommand(args)
	})
}