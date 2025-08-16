package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddPolicySentryTools adds Policy Sentry MCP tool implementations
func AddPolicySentryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Policy Sentry create template tool
	createTemplateTool := mcp.NewTool("policy_sentry_create_template",
		mcp.WithDescription("Create IAM policy template using Policy Sentry"),
		mcp.WithString("template_type",
			mcp.Description("Type of template (crud, actions, service)"),
			mcp.Required(),
			mcp.Enum("crud", "actions", "service"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for template"),
		),
	)
	s.AddTool(createTemplateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateType := request.GetString("template_type", "")
		args := []string{"security", "policy-sentry", "create-template", "--type", templateType}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// Policy Sentry write policy tool
	writePolicyTool := mcp.NewTool("policy_sentry_write_policy",
		mcp.WithDescription("Write IAM policy from input file using Policy Sentry"),
		mcp.WithString("input_file",
			mcp.Description("Input YAML file path"),
			mcp.Required(),
		),
	)
	s.AddTool(writePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputFile := request.GetString("input_file", "")
		args := []string{"security", "policy-sentry", "write-policy", "--input", inputFile}
		return executeShipCommand(args)
	})

	// Policy Sentry write policy from template tool
	writePolicyFromTemplateTool := mcp.NewTool("policy_sentry_write_from_template",
		mcp.WithDescription("Write IAM policy from template YAML using Policy Sentry"),
		mcp.WithString("template_yaml",
			mcp.Description("Template YAML content"),
			mcp.Required(),
		),
	)
	s.AddTool(writePolicyFromTemplateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateYAML := request.GetString("template_yaml", "")
		args := []string{"security", "policy-sentry", "write-policy", "--template", templateYAML}
		return executeShipCommand(args)
	})

	// Policy Sentry write policy with actions tool
	writePolicyWithActionsTool := mcp.NewTool("policy_sentry_write_with_actions",
		mcp.WithDescription("Write IAM policy with specific actions using Policy Sentry"),
		mcp.WithString("actions",
			mcp.Description("Comma-separated list of IAM actions"),
			mcp.Required(),
		),
		mcp.WithString("resource_arns",
			mcp.Description("Comma-separated list of resource ARNs"),
			mcp.Required(),
		),
	)
	s.AddTool(writePolicyWithActionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		actions := request.GetString("actions", "")
		resourceArns := request.GetString("resource_arns", "")
		args := []string{"security", "policy-sentry", "write-policy", "--actions", actions, "--resources", resourceArns}
		return executeShipCommand(args)
	})

	// Policy Sentry write policy with CRUD tool
	writePolicyWithCRUDTool := mcp.NewTool("policy_sentry_write_with_crud",
		mcp.WithDescription("Write IAM policy with CRUD access levels using Policy Sentry"),
		mcp.WithString("resource_arns",
			mcp.Description("Comma-separated list of resource ARNs"),
			mcp.Required(),
		),
		mcp.WithString("access_levels",
			mcp.Description("Comma-separated list of access levels (read, write, list, tagging, permissions-management)"),
			mcp.Required(),
		),
	)
	s.AddTool(writePolicyWithCRUDTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resourceArns := request.GetString("resource_arns", "")
		accessLevels := request.GetString("access_levels", "")
		args := []string{"security", "policy-sentry", "write-policy", "--crud", "--resources", resourceArns, "--access-levels", accessLevels}
		return executeShipCommand(args)
	})

	// Policy Sentry query action table tool
	queryActionTableTool := mcp.NewTool("policy_sentry_query_actions",
		mcp.WithDescription("Query AWS service action table using Policy Sentry"),
		mcp.WithString("service",
			mcp.Description("AWS service name (e.g., s3, ec2, iam)"),
			mcp.Required(),
		),
	)
	s.AddTool(queryActionTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service := request.GetString("service", "")
		args := []string{"security", "policy-sentry", "query", "actions", "--service", service}
		return executeShipCommand(args)
	})

	// Policy Sentry query condition table tool
	queryConditionTableTool := mcp.NewTool("policy_sentry_query_conditions",
		mcp.WithDescription("Query AWS service condition table using Policy Sentry"),
		mcp.WithString("service",
			mcp.Description("AWS service name (e.g., s3, ec2, iam)"),
			mcp.Required(),
		),
	)
	s.AddTool(queryConditionTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service := request.GetString("service", "")
		args := []string{"security", "policy-sentry", "query", "conditions", "--service", service}
		return executeShipCommand(args)
	})
}