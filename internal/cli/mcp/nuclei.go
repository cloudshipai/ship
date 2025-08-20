package mcp

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddNucleiTools adds Nuclei vulnerability scanning tools to the MCP server
func AddNucleiTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Scan URL
	scanURLTool := mcp.NewTool("nuclei_scan_url",
		mcp.WithDescription("Scan a URL for vulnerabilities using Nuclei templates"),
		mcp.WithString("url",
			mcp.Description("Target URL to scan"),
		),
		mcp.WithString("severity",
			mcp.Description("Filter by severity (info, low, medium, high, critical)"),
		),
	)
	s.AddTool(scanURLTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url := request.GetString("url", "")
		severity := request.GetString("severity", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNucleiModule(client)
		result, err := module.ScanURL(ctx, url, severity)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nuclei scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Scan with template
	scanWithTemplateTool := mcp.NewTool("nuclei_scan_with_template",
		mcp.WithDescription("Scan using specific Nuclei template(s)"),
		mcp.WithString("url",
			mcp.Description("Target URL to scan"),
		),
		mcp.WithString("template",
			mcp.Description("Template path or name"),
		),
	)
	s.AddTool(scanWithTemplateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url := request.GetString("url", "")
		template := request.GetString("template", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNucleiModule(client)
		result, err := module.ScanWithTemplate(ctx, url, template)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nuclei scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Scan with tags
	scanWithTagsTool := mcp.NewTool("nuclei_scan_with_tags",
		mcp.WithDescription("Scan using specific vulnerability tags"),
		mcp.WithString("url",
			mcp.Description("Target URL to scan"),
		),
		mcp.WithString("tags",
			mcp.Description("Comma-separated list of tags"),
		),
	)
	s.AddTool(scanWithTagsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url := request.GetString("url", "")
		tags := request.GetString("tags", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNucleiModule(client)
		result, err := module.ScanWithTags(ctx, url, tags)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nuclei scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Update templates
	updateTemplatesTool := mcp.NewTool("nuclei_update_templates",
		mcp.WithDescription("Update Nuclei vulnerability templates"),
	)
	s.AddTool(updateTemplatesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNucleiModule(client)
		result, err := module.UpdateTemplates(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nuclei update templates failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Validate template
	validateTemplateTool := mcp.NewTool("nuclei_validate_template",
		mcp.WithDescription("Validate a Nuclei template"),
		mcp.WithString("template_path",
			mcp.Description("Path to template file to validate"),
		),
	)
	s.AddTool(validateTemplateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templatePath := request.GetString("template_path", "")

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNucleiModule(client)
		result, err := module.ValidateTemplate(ctx, templatePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nuclei template validation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Generate report
	generateReportTool := mcp.NewTool("nuclei_generate_report",
		mcp.WithDescription("Generate a vulnerability scan report"),
		mcp.WithString("url",
			mcp.Description("Target URL to scan"),
		),
		mcp.WithString("report_type",
			mcp.Description("Report format (json, markdown, sarif)"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url := request.GetString("url", "")
		reportType := request.GetString("report_type", "")
		if reportType == "" {
			reportType = "json"
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewNucleiModule(client)
		result, err := module.GenerateReport(ctx, url, reportType)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nuclei report generation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}
