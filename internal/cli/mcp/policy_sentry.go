package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddPolicySentryTools adds Policy Sentry (AWS IAM policy generator) MCP tool implementations using direct Dagger calls
func AddPolicySentryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addPolicySentryToolsDirect(s)
}

// addPolicySentryToolsDirect adds Policy Sentry tools using direct Dagger module calls
func addPolicySentryToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPolicySentryModule(client)

		// Get parameters
		templateType := request.GetString("template_type", "")
		if templateType == "" {
			return mcp.NewToolResultError("template_type is required"), nil
		}
		outputFile := request.GetString("output_file", "")

		// Create template
		output, err := module.CreateTemplate(ctx, templateType, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("policy sentry create template failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPolicySentryModule(client)

		// Get parameters
		inputFile := request.GetString("input_file", "")
		if inputFile == "" {
			return mcp.NewToolResultError("input_file is required"), nil
		}

		// Note: Dagger module doesn't support minimize and format options
		if request.GetBool("minimize", false) || request.GetString("format", "") != "" {
			return mcp.NewToolResultError("minimize and format options not supported in Dagger module"), nil
		}

		// Write policy
		output, err := module.WritePolicy(ctx, inputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("policy sentry write policy failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Policy Sentry initialize tool
	initializeTool := mcp.NewTool("policy_sentry_initialize",
		mcp.WithDescription("Initialize Policy Sentry IAM database using real policy_sentry CLI"),
		mcp.WithBoolean("fetch",
			mcp.Description("Fetch latest AWS documentation from AWS docs"),
		),
	)
	s.AddTool(initializeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Note: Initialize function not available in Dagger module
		// This would require persistent database initialization
		return mcp.NewToolResultError("initialize function not supported in Dagger module - database is pre-initialized in container"), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPolicySentryModule(client)

		// Get parameters
		service := request.GetString("service", "")
		if service == "" {
			return mcp.NewToolResultError("service is required"), nil
		}

		// Note: Dagger module doesn't support advanced filter options
		if request.GetString("name", "") != "" || request.GetString("access_level", "") != "" ||
			request.GetString("condition", "") != "" || request.GetString("resource_type", "") != "" ||
			request.GetString("format", "") != "" {
			return mcp.NewToolResultError("advanced query options not supported in Dagger module"), nil
		}

		// Query action table
		output, err := module.QueryActionTable(ctx, service)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("policy sentry query action table failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPolicySentryModule(client)

		// Get parameters
		service := request.GetString("service", "")
		if service == "" {
			return mcp.NewToolResultError("service is required"), nil
		}

		// Note: Dagger module doesn't support name and format options
		if request.GetString("name", "") != "" || request.GetString("format", "") != "" {
			return mcp.NewToolResultError("name and format options not supported in Dagger module"), nil
		}

		// Query condition table
		output, err := module.QueryConditionTable(ctx, service)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("policy sentry query condition table failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Note: ARN table query not available in Dagger module
		// Only basic action and condition table queries are supported
		return mcp.NewToolResultError("ARN table query not supported in Dagger module"), nil
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
		// Note: Service table query not available in Dagger module
		// Only basic action and condition table queries are supported
		return mcp.NewToolResultError("service table query not supported in Dagger module"), nil
	})



}