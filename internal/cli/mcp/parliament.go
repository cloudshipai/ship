package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddParliamentTools adds Parliament (AWS IAM policy linter) MCP tool implementations
func AddParliamentTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Parliament lint policy file tool
	lintPolicyFileTool := mcp.NewTool("parliament_lint_file",
		mcp.WithDescription("Lint AWS IAM policy file using Parliament"),
		mcp.WithString("policy_path",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
	)
	s.AddTool(lintPolicyFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "parliament", "--file", policyPath}
		return executeShipCommand(args)
	})

	// Parliament lint policy directory tool
	lintPolicyDirectoryTool := mcp.NewTool("parliament_lint_directory",
		mcp.WithDescription("Lint AWS IAM policy files in directory using Parliament"),
		mcp.WithString("directory",
			mcp.Description("Directory containing IAM policy files"),
			mcp.Required(),
		),
	)
	s.AddTool(lintPolicyDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"security", "parliament", "--directory", directory}
		return executeShipCommand(args)
	})

	// Parliament lint policy string tool
	lintPolicyStringTool := mcp.NewTool("parliament_lint_string",
		mcp.WithDescription("Lint AWS IAM policy JSON string using Parliament"),
		mcp.WithString("policy_json",
			mcp.Description("IAM policy JSON string"),
			mcp.Required(),
		),
	)
	s.AddTool(lintPolicyStringTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyJSON := request.GetString("policy_json", "")
		args := []string{"security", "parliament", "--string", policyJSON}
		return executeShipCommand(args)
	})

	// Parliament lint with community auditors tool
	lintWithCommunityAuditorsTool := mcp.NewTool("parliament_lint_community",
		mcp.WithDescription("Lint AWS IAM policy with community auditors using Parliament"),
		mcp.WithString("policy_path",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
	)
	s.AddTool(lintWithCommunityAuditorsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "parliament", "--file", policyPath, "--include-community"}
		return executeShipCommand(args)
	})

	// Parliament lint with private auditors tool
	lintWithPrivateAuditorsTool := mcp.NewTool("parliament_lint_private",
		mcp.WithDescription("Lint AWS IAM policy with private auditors using Parliament"),
		mcp.WithString("policy_path",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
		mcp.WithString("auditors_path",
			mcp.Description("Path to private auditors directory"),
			mcp.Required(),
		),
	)
	s.AddTool(lintWithPrivateAuditorsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		auditorsPath := request.GetString("auditors_path", "")
		args := []string{"security", "parliament", "--file", policyPath, "--private-auditors", auditorsPath}
		return executeShipCommand(args)
	})

	// Parliament lint with severity filter tool
	lintWithSeverityFilterTool := mcp.NewTool("parliament_lint_severity",
		mcp.WithDescription("Lint AWS IAM policy with severity filter using Parliament"),
		mcp.WithString("policy_path",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
		mcp.WithString("min_severity",
			mcp.Description("Minimum severity level (LOW, MEDIUM, HIGH, CRITICAL)"),
			mcp.Required(),
			mcp.Enum("LOW", "MEDIUM", "HIGH", "CRITICAL"),
		),
	)
	s.AddTool(lintWithSeverityFilterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		minSeverity := request.GetString("min_severity", "")
		args := []string{"security", "parliament", "--file", policyPath, "--min-severity", minSeverity}
		return executeShipCommand(args)
	})
}