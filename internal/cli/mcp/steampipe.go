package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSteampipeTools adds Steampipe (cloud asset querying) MCP tool implementations
// NOTE: Steampipe is typically configured as an external MCP server via npx @turbot/steampipe-mcp
// These tools provide Dagger-based execution as an alternative
func AddSteampipeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Steampipe query tool
	queryTool := mcp.NewTool("steampipe_query",
		mcp.WithDescription("Execute SQL query against cloud resources using Steampipe"),
		mcp.WithString("query",
			mcp.Description("SQL query to execute"),
			mcp.Required(),
		),
		mcp.WithString("plugin",
			mcp.Description("Steampipe plugin to use"),
			mcp.Required(),
			mcp.Enum("aws", "azure", "gcp", "kubernetes", "github", "slack"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("json", "csv", "table"),
		),
	)
	s.AddTool(queryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.GetString("query", "")
		plugin := request.GetString("plugin", "")
		args := []string{"security", "steampipe", "--query", query, "--plugin", plugin}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})

	// Steampipe query from file tool
	queryFromFileTool := mcp.NewTool("steampipe_query_from_file",
		mcp.WithDescription("Execute SQL queries from file against cloud resources"),
		mcp.WithString("query_file",
			mcp.Description("Path to SQL query file"),
			mcp.Required(),
		),
		mcp.WithString("plugin",
			mcp.Description("Steampipe plugin to use"),
			mcp.Required(),
			mcp.Enum("aws", "azure", "gcp", "kubernetes", "github", "slack"),
		),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("json", "csv", "table"),
		),
	)
	s.AddTool(queryFromFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		queryFile := request.GetString("query_file", "")
		plugin := request.GetString("plugin", "")
		args := []string{"security", "steampipe", "--query-file", queryFile, "--plugin", plugin}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})

	// Steampipe list plugins tool
	listPluginsTool := mcp.NewTool("steampipe_list_plugins",
		mcp.WithDescription("List available Steampipe plugins"),
	)
	s.AddTool(listPluginsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "steampipe", "--list-plugins"}
		return executeShipCommand(args)
	})

	// Steampipe cost analysis tool
	costAnalysisTool := mcp.NewTool("steampipe_cost_analysis",
		mcp.WithDescription("Run predefined cost analysis queries for AWS resources"),
		mcp.WithString("analysis_type",
			mcp.Description("Type of cost analysis to perform"),
			mcp.Required(),
			mcp.Enum("ec2_idle_instances", "ebs_unattached_volumes", "s3_unused_buckets", "cloudwatch_unused_log_groups", "elb_unused_load_balancers", "rds_idle_instances", "all_cost_issues"),
		),
		mcp.WithString("region",
			mcp.Description("AWS region to analyze (default: all regions)"),
		),
		mcp.WithString("min_age_days",
			mcp.Description("Minimum age in days for resources (default: 7)"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "csv", "table"),
		),
	)
	s.AddTool(costAnalysisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		analysisType := request.GetString("analysis_type", "")
		args := []string{"security", "steampipe", "--cost-analysis", analysisType}
		
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		if minAge := request.GetString("min_age_days", ""); minAge != "" {
			args = append(args, "--min-age", minAge)
		}
		if output := request.GetString("output_format", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})

	// Steampipe compliance check tool
	complianceCheckTool := mcp.NewTool("steampipe_compliance_check",
		mcp.WithDescription("Run compliance checks against AWS resources using Steampipe benchmarks"),
		mcp.WithString("benchmark",
			mcp.Description("Compliance benchmark to run"),
			mcp.Required(),
			mcp.Enum("aws_thrifty", "aws_compliance", "aws_foundational_security", "cis_v140"),
		),
		mcp.WithString("tags",
			mcp.Description("Comma-separated tags to filter controls"),
		),
		mcp.WithString("where_clause",
			mcp.Description("SQL WHERE clause for additional filtering"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "csv", "table"),
		),
	)
	s.AddTool(complianceCheckTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		benchmark := request.GetString("benchmark", "")
		args := []string{"security", "steampipe", "--benchmark", benchmark}
		
		if tags := request.GetString("tags", ""); tags != "" {
			args = append(args, "--tags", tags)
		}
		if whereClause := request.GetString("where_clause", ""); whereClause != "" {
			args = append(args, "--where", whereClause)
		}
		if output := request.GetString("output_format", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})

	// Steampipe security assessment tool
	securityAssessmentTool := mcp.NewTool("steampipe_security_assessment",
		mcp.WithDescription("Comprehensive security assessment of cloud resources"),
		mcp.WithString("assessment_type",
			mcp.Description("Type of security assessment"),
			mcp.Required(),
			mcp.Enum("public_resources", "unencrypted_data", "overprivileged_access", "network_exposure", "compliance_gaps"),
		),
		mcp.WithString("cloud_provider",
			mcp.Description("Cloud provider to assess"),
			mcp.Required(),
			mcp.Enum("aws", "azure", "gcp"),
		),
		mcp.WithString("severity_threshold",
			mcp.Description("Minimum severity to report"),
			mcp.Enum("low", "medium", "high", "critical"),
		),
	)
	s.AddTool(securityAssessmentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		assessmentType := request.GetString("assessment_type", "")
		cloudProvider := request.GetString("cloud_provider", "")
		args := []string{"security", "steampipe", "--security-assessment", assessmentType, "--provider", cloudProvider}
		
		if severity := request.GetString("severity_threshold", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		return executeShipCommand(args)
	})

	// Steampipe resource inventory tool
	resourceInventoryTool := mcp.NewTool("steampipe_resource_inventory",
		mcp.WithDescription("Generate comprehensive cloud resource inventory"),
		mcp.WithString("cloud_provider",
			mcp.Description("Cloud provider for inventory"),
			mcp.Required(),
			mcp.Enum("aws", "azure", "gcp", "kubernetes"),
		),
		mcp.WithString("resource_types",
			mcp.Description("Comma-separated resource types to include"),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated regions to include"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("json", "csv", "table"),
		),
		mcp.WithBoolean("include_metadata",
			mcp.Description("Include detailed resource metadata"),
		),
	)
	s.AddTool(resourceInventoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cloudProvider := request.GetString("cloud_provider", "")
		args := []string{"security", "steampipe", "--inventory", "--provider", cloudProvider}
		
		if resourceTypes := request.GetString("resource_types", ""); resourceTypes != "" {
			args = append(args, "--resource-types", resourceTypes)
		}
		if regions := request.GetString("regions", ""); regions != "" {
			args = append(args, "--regions", regions)
		}
		if output := request.GetString("output_format", ""); output != "" {
			args = append(args, "--output", output)
		}
		if request.GetBool("include_metadata", false) {
			args = append(args, "--include-metadata")
		}
		return executeShipCommand(args)
	})

	// Steampipe get version tool
	getVersionTool := mcp.NewTool("steampipe_get_version",
		mcp.WithDescription("Get Steampipe version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "steampipe", "--version"}
		return executeShipCommand(args)
	})
}