package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddPolicySentryTools adds Policy Sentry (AWS IAM policy generator) MCP tool implementations using real CLI commands
func AddPolicySentryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Policy Sentry create template tool
	createTemplateTool := mcp.NewTool("policy_sentry_create_template",
		mcp.WithDescription("Create IAM policy template using real policy_sentry CLI"),
		mcp.WithString("template_type",
			mcp.Description("Type of template (crud, actions)"),
			mcp.Required(),
			mcp.Enum("crud", "actions"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for template"),
		),
	)
	s.AddTool(createTemplateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateType := request.GetString("template_type", "")
		args := []string{"policy_sentry", "create-template", "--template-type", templateType}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output-file", outputFile)
		}
		return executeShipCommand(args)
	})

	// Policy Sentry write policy tool
	writePolicyTool := mcp.NewTool("policy_sentry_write_policy",
		mcp.WithDescription("Write IAM policy from input YAML file using real policy_sentry CLI"),
		mcp.WithString("input_file",
			mcp.Description("Path of the YAML file for generating policies"),
			mcp.Required(),
		),
		mcp.WithBoolean("minimize",
			mcp.Description("Minimize policy statements"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("yaml", "json", "terraform"),
		),
	)
	s.AddTool(writePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputFile := request.GetString("input_file", "")
		args := []string{"policy_sentry", "write-policy", "--input-file", inputFile}
		
		if request.GetBool("minimize", false) {
			args = append(args, "--minimize")
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--fmt", format)
		}
		
		return executeShipCommand(args)
	})

	// Policy Sentry initialize tool
	initializeTool := mcp.NewTool("policy_sentry_initialize",
		mcp.WithDescription("Initialize Policy Sentry IAM database using real policy_sentry CLI"),
		mcp.WithBoolean("fetch",
			mcp.Description("Fetch latest AWS documentation from AWS docs"),
		),
	)
	s.AddTool(initializeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"policy_sentry", "initialize"}
		
		if request.GetBool("fetch", false) {
			args = append(args, "--fetch")
		}
		
		return executeShipCommand(args)
	})



	// Policy Sentry query action table tool
	queryActionTableTool := mcp.NewTool("policy_sentry_query_action_table",
		mcp.WithDescription("Query AWS service action table using real policy_sentry CLI"),
		mcp.WithString("service",
			mcp.Description("AWS service name (e.g., s3, ec2, iam)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("IAM Action name"),
		),
		mcp.WithString("access_level",
			mcp.Description("Access level filter"),
			mcp.Enum("read", "write", "list", "tagging", "permissions-management"),
		),
		mcp.WithString("condition",
			mcp.Description("Condition key filter"),
		),
		mcp.WithString("resource_type",
			mcp.Description("Resource type filter"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("yaml", "json"),
		),
	)
	s.AddTool(queryActionTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service := request.GetString("service", "")
		args := []string{"policy_sentry", "query", "action-table", "--service", service}
		
		if name := request.GetString("name", ""); name != "" {
			args = append(args, "--name", name)
		}
		if accessLevel := request.GetString("access_level", ""); accessLevel != "" {
			args = append(args, "--access-level", accessLevel)
		}
		if condition := request.GetString("condition", ""); condition != "" {
			args = append(args, "--condition", condition)
		}
		if resourceType := request.GetString("resource_type", ""); resourceType != "" {
			args = append(args, "--resource-type", resourceType)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--fmt", format)
		}
		
		return executeShipCommand(args)
	})

	// Policy Sentry query condition table tool
	queryConditionTableTool := mcp.NewTool("policy_sentry_query_condition_table",
		mcp.WithDescription("Query AWS service condition table using real policy_sentry CLI"),
		mcp.WithString("service",
			mcp.Description("AWS service name (e.g., s3, ec2, iam)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Condition key name"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("yaml", "json"),
		),
	)
	s.AddTool(queryConditionTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service := request.GetString("service", "")
		args := []string{"policy_sentry", "query", "condition-table", "--service", service}
		
		if name := request.GetString("name", ""); name != "" {
			args = append(args, "--name", name)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--fmt", format)
		}
		
		return executeShipCommand(args)
	})

	// Policy Sentry query ARN table tool
	queryArnTableTool := mcp.NewTool("policy_sentry_query_arn_table",
		mcp.WithDescription("Query AWS service ARN table using real policy_sentry CLI"),
		mcp.WithString("service",
			mcp.Description("AWS service name (e.g., s3, ec2, iam)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Resource ARN type name"),
		),
		mcp.WithBoolean("list_arn_types",
			mcp.Description("List ARN types"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("yaml", "json"),
		),
	)
	s.AddTool(queryArnTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service := request.GetString("service", "")
		args := []string{"policy_sentry", "query", "arn-table", "--service", service}
		
		if name := request.GetString("name", ""); name != "" {
			args = append(args, "--name", name)
		}
		if request.GetBool("list_arn_types", false) {
			args = append(args, "--list-arn-types")
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--fmt", format)
		}
		
		return executeShipCommand(args)
	})

	// Policy Sentry query service table tool
	queryServiceTableTool := mcp.NewTool("policy_sentry_query_service_table",
		mcp.WithDescription("Query AWS service table using real policy_sentry CLI"),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("yaml", "json", "csv"),
		),
	)
	s.AddTool(queryServiceTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"policy_sentry", "query", "service-table"}
		
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--fmt", format)
		}
		
		return executeShipCommand(args)
	})



}