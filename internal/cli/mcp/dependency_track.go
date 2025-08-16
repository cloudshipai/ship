package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddDependencyTrackTools adds Dependency Track (software component analysis) MCP tool implementations
func AddDependencyTrackTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Dependency Track upload BOM tool
	uploadBOMTool := mcp.NewTool("dependency_track_upload_bom",
		mcp.WithDescription("Upload Software Bill of Materials to Dependency Track"),
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
	)
	s.AddTool(uploadBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bomPath := request.GetString("bom_path", "")
		projectName := request.GetString("project_name", "")
		args := []string{"security", "dependency-track", "upload", bomPath, "--project", projectName}
		if projectVersion := request.GetString("project_version", ""); projectVersion != "" {
			args = append(args, "--version", projectVersion)
		}
		return executeShipCommand(args)
	})

	// Dependency Track scan vulnerabilities tool
	scanVulnerabilitiesTool := mcp.NewTool("dependency_track_scan_vulnerabilities",
		mcp.WithDescription("Scan project for vulnerabilities in Dependency Track"),
		mcp.WithString("project_name",
			mcp.Description("Project name to scan"),
			mcp.Required(),
		),
		mcp.WithString("project_version",
			mcp.Description("Project version to scan"),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level (info, low, medium, high, critical)"),
		),
	)
	s.AddTool(scanVulnerabilitiesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName := request.GetString("project_name", "")
		args := []string{"security", "dependency-track", "scan", projectName}
		if projectVersion := request.GetString("project_version", ""); projectVersion != "" {
			args = append(args, "--version", projectVersion)
		}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		return executeShipCommand(args)
	})

	// Dependency Track generate report tool
	generateReportTool := mcp.NewTool("dependency_track_generate_report",
		mcp.WithDescription("Generate vulnerability report for project"),
		mcp.WithString("project_name",
			mcp.Description("Project name"),
			mcp.Required(),
		),
		mcp.WithString("project_version",
			mcp.Description("Project version"),
		),
		mcp.WithString("report_format",
			mcp.Description("Report format (json, csv, xml, pdf)"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for the report"),
		),
	)
	s.AddTool(generateReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName := request.GetString("project_name", "")
		args := []string{"security", "dependency-track", "report", projectName}
		if projectVersion := request.GetString("project_version", ""); projectVersion != "" {
			args = append(args, "--version", projectVersion)
		}
		if reportFormat := request.GetString("report_format", ""); reportFormat != "" {
			args = append(args, "--format", reportFormat)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// Dependency Track list projects tool
	listProjectsTool := mcp.NewTool("dependency_track_list_projects",
		mcp.WithDescription("List all projects in Dependency Track"),
		mcp.WithBoolean("show_metrics",
			mcp.Description("Include project metrics in output"),
		),
	)
	s.AddTool(listProjectsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "dependency-track", "list-projects"}
		if request.GetBool("show_metrics", false) {
			args = append(args, "--metrics")
		}
		return executeShipCommand(args)
	})

	// Dependency Track analyze components tool
	analyzeComponentsTool := mcp.NewTool("dependency_track_analyze_components",
		mcp.WithDescription("Analyze components for security issues and policy violations"),
		mcp.WithString("project_name",
			mcp.Description("Project name to analyze"),
			mcp.Required(),
		),
		mcp.WithString("project_version",
			mcp.Description("Project version"),
		),
		mcp.WithBoolean("include_policy_violations",
			mcp.Description("Include policy violations in analysis"),
		),
	)
	s.AddTool(analyzeComponentsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName := request.GetString("project_name", "")
		args := []string{"security", "dependency-track", "analyze", projectName}
		if projectVersion := request.GetString("project_version", ""); projectVersion != "" {
			args = append(args, "--version", projectVersion)
		}
		if request.GetBool("include_policy_violations", false) {
			args = append(args, "--include-policy")
		}
		return executeShipCommand(args)
	})

	// Dependency Track create project tool
	createProjectTool := mcp.NewTool("dependency_track_create_project",
		mcp.WithDescription("Create new project in Dependency Track"),
		mcp.WithString("project_name",
			mcp.Description("Name for the new project"),
			mcp.Required(),
		),
		mcp.WithString("project_version",
			mcp.Description("Initial version for the project"),
		),
		mcp.WithString("description",
			mcp.Description("Project description"),
		),
	)
	s.AddTool(createProjectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName := request.GetString("project_name", "")
		args := []string{"security", "dependency-track", "create-project", projectName}
		if projectVersion := request.GetString("project_version", ""); projectVersion != "" {
			args = append(args, "--version", projectVersion)
		}
		if description := request.GetString("description", ""); description != "" {
			args = append(args, "--description", description)
		}
		return executeShipCommand(args)
	})

	// Dependency Track get project metrics tool
	getMetricsTool := mcp.NewTool("dependency_track_get_metrics",
		mcp.WithDescription("Get security metrics for a project"),
		mcp.WithString("project_name",
			mcp.Description("Project name"),
			mcp.Required(),
		),
		mcp.WithString("project_version",
			mcp.Description("Project version"),
		),
	)
	s.AddTool(getMetricsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName := request.GetString("project_name", "")
		args := []string{"security", "dependency-track", "metrics", projectName}
		if projectVersion := request.GetString("project_version", ""); projectVersion != "" {
			args = append(args, "--version", projectVersion)
		}
		return executeShipCommand(args)
	})
}