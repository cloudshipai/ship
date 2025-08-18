package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddParliamentTools adds Parliament (AWS IAM policy linter) MCP tool implementations using real CLI commands
func AddParliamentTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Parliament lint policy file tool
	lintPolicyFileTool := mcp.NewTool("parliament_lint_file",
		mcp.WithDescription("Lint AWS IAM policy file using real parliament CLI"),
		mcp.WithString("policy_path",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Path to custom configuration file"),
		),
		mcp.WithBoolean("json_output",
			mcp.Description("Output results in JSON format"),
		),
	)
	s.AddTool(lintPolicyFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"parliament", "--file", policyPath}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json_output", false) {
			args = append(args, "--json")
		}
		
		return executeShipCommand(args)
	})

	// Parliament lint policy directory tool
	lintPolicyDirectoryTool := mcp.NewTool("parliament_lint_directory",
		mcp.WithDescription("Lint AWS IAM policy files in directory using real parliament CLI"),
		mcp.WithString("directory",
			mcp.Description("Directory containing IAM policy files"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Path to custom configuration file"),
		),
		mcp.WithBoolean("json_output",
			mcp.Description("Output results in JSON format"),
		),
		mcp.WithString("include_policy_extension",
			mcp.Description("File extension to include (e.g., json)"),
		),
		mcp.WithString("exclude_pattern",
			mcp.Description("Pattern to exclude (regex)"),
		),
	)
	s.AddTool(lintPolicyDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"parliament", "--directory", directory}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json_output", false) {
			args = append(args, "--json")
		}
		if extension := request.GetString("include_policy_extension", ""); extension != "" {
			args = append(args, "--include_policy_extension", extension)
		}
		if pattern := request.GetString("exclude_pattern", ""); pattern != "" {
			args = append(args, "--exclude_pattern", pattern)
		}
		
		return executeShipCommand(args)
	})

	// Parliament lint policy string tool
	lintPolicyStringTool := mcp.NewTool("parliament_lint_string",
		mcp.WithDescription("Lint AWS IAM policy JSON string using real parliament CLI"),
		mcp.WithString("policy_json",
			mcp.Description("IAM policy JSON string"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Path to custom configuration file"),
		),
		mcp.WithBoolean("json_output",
			mcp.Description("Output results in JSON format"),
		),
	)
	s.AddTool(lintPolicyStringTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyJSON := request.GetString("policy_json", "")
		args := []string{"parliament", "--string", policyJSON}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json_output", false) {
			args = append(args, "--json")
		}
		
		return executeShipCommand(args)
	})

	// Parliament lint with community auditors tool
	lintWithCommunityAuditorsTool := mcp.NewTool("parliament_lint_community",
		mcp.WithDescription("Lint AWS IAM policy with community auditors using real parliament CLI"),
		mcp.WithString("policy_path",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Path to custom configuration file"),
		),
		mcp.WithBoolean("json_output",
			mcp.Description("Output results in JSON format"),
		),
	)
	s.AddTool(lintWithCommunityAuditorsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"parliament", "--file", policyPath, "--include-community-auditors"}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json_output", false) {
			args = append(args, "--json")
		}
		
		return executeShipCommand(args)
	})

	// Parliament lint with private auditors tool
	lintWithPrivateAuditorsTool := mcp.NewTool("parliament_lint_private",
		mcp.WithDescription("Lint AWS IAM policy with private auditors using real parliament CLI"),
		mcp.WithString("policy_path",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
		mcp.WithString("private_auditors",
			mcp.Description("Path to private auditors directory"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Path to custom configuration file"),
		),
		mcp.WithBoolean("json_output",
			mcp.Description("Output results in JSON format"),
		),
	)
	s.AddTool(lintWithPrivateAuditorsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		privateAuditors := request.GetString("private_auditors", "")
		args := []string{"parliament", "--file", policyPath, "--private_auditors", privateAuditors}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json_output", false) {
			args = append(args, "--json")
		}
		
		return executeShipCommand(args)
	})

	// Parliament lint AWS managed policies tool
	lintAwsManagedPolicesTool := mcp.NewTool("parliament_lint_aws_managed",
		mcp.WithDescription("Lint AWS managed policies using real parliament CLI"),
		mcp.WithString("config",
			mcp.Description("Path to custom configuration file"),
		),
		mcp.WithBoolean("json_output",
			mcp.Description("Output results in JSON format"),
		),
	)
	s.AddTool(lintAwsManagedPolicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"parliament", "--aws-managed-policies"}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json_output", false) {
			args = append(args, "--json")
		}
		
		return executeShipCommand(args)
	})

	// Parliament lint auth details file tool
	lintAuthDetailsFileTool := mcp.NewTool("parliament_lint_auth_details",
		mcp.WithDescription("Lint AWS IAM authorization details file using real parliament CLI"),
		mcp.WithString("auth_details_file",
			mcp.Description("Path to AWS IAM authorization details file"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Path to custom configuration file"),
		),
		mcp.WithBoolean("json_output",
			mcp.Description("Output results in JSON format"),
		),
	)
	s.AddTool(lintAuthDetailsFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		authDetailsFile := request.GetString("auth_details_file", "")
		args := []string{"parliament", "--auth-details-file", authDetailsFile}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json_output", false) {
			args = append(args, "--json")
		}
		
		return executeShipCommand(args)
	})

	// Parliament comprehensive analysis tool
	comprehensiveAnalysisTool := mcp.NewTool("parliament_comprehensive_analysis",
		mcp.WithDescription("Comprehensive IAM policy analysis with all auditors using real parliament CLI"),
		mcp.WithString("policy_path",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
		mcp.WithString("private_auditors",
			mcp.Description("Path to private auditors directory"),
		),
		mcp.WithString("config",
			mcp.Description("Path to custom configuration file"),
		),
		mcp.WithBoolean("json_output",
			mcp.Description("Output results in JSON format"),
		),
	)
	s.AddTool(comprehensiveAnalysisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"parliament", "--file", policyPath, "--include-community-auditors"}
		
		if privateAuditors := request.GetString("private_auditors", ""); privateAuditors != "" {
			args = append(args, "--private_auditors", privateAuditors)
		}
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json_output", false) {
			args = append(args, "--json")
		}
		
		return executeShipCommand(args)
	})

	// Parliament batch directory analysis tool
	batchDirectoryAnalysisTool := mcp.NewTool("parliament_batch_directory_analysis",
		mcp.WithDescription("Batch analysis of multiple policy directories using real parliament CLI"),
		mcp.WithString("base_directory",
			mcp.Description("Base directory containing policy subdirectories"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Path to custom configuration file"),
		),
		mcp.WithString("private_auditors",
			mcp.Description("Path to private auditors directory"),
		),
		mcp.WithBoolean("json_output",
			mcp.Description("Output results in JSON format"),
		),
		mcp.WithString("include_policy_extension",
			mcp.Description("File extension to include (default: json)"),
		),
		mcp.WithString("exclude_pattern",
			mcp.Description("Pattern to exclude from analysis"),
		),
	)
	s.AddTool(batchDirectoryAnalysisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		baseDirectory := request.GetString("base_directory", "")
		args := []string{"parliament", "--directory", baseDirectory, "--include-community-auditors"}
		
		if privateAuditors := request.GetString("private_auditors", ""); privateAuditors != "" {
			args = append(args, "--private_auditors", privateAuditors)
		}
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if request.GetBool("json_output", false) {
			args = append(args, "--json")
		}
		if extension := request.GetString("include_policy_extension", ""); extension != "" {
			args = append(args, "--include_policy_extension", extension)
		}
		if pattern := request.GetString("exclude_pattern", ""); pattern != "" {
			args = append(args, "--exclude_pattern", pattern)
		}
		
		return executeShipCommand(args)
	})
}