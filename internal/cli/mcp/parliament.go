package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddParliamentTools adds Parliament (AWS IAM policy linter) MCP tool implementations using direct Dagger calls
func AddParliamentTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addParliamentToolsDirect(s)
}

// addParliamentToolsDirect adds Parliament tools using direct Dagger module calls
func addParliamentToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewParliamentModule(client)

		// Get parameters
		policyPath := request.GetString("policy_path", "")
		if policyPath == "" {
			return mcp.NewToolResultError("policy_path is required"), nil
		}

		// Note: Dagger module doesn't support config and json_output options for basic linting
		if request.GetString("config", "") != "" || request.GetBool("json_output", false) {
			return mcp.NewToolResultError("config and json_output options not supported in Dagger module for basic file linting"), nil
		}

		// Lint policy file
		output, err := module.LintPolicyFile(ctx, policyPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parliament lint failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewParliamentModule(client)

		// Get parameters
		directory := request.GetString("directory", "")
		if directory == "" {
			return mcp.NewToolResultError("directory is required"), nil
		}

		// Note: Dagger module doesn't support advanced directory options
		if request.GetString("config", "") != "" || request.GetBool("json_output", false) ||
			request.GetString("include_policy_extension", "") != "" || request.GetString("exclude_pattern", "") != "" {
			return mcp.NewToolResultError("advanced directory options not supported in Dagger module"), nil
		}

		// Lint policy directory
		output, err := module.LintPolicyDirectory(ctx, directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parliament directory lint failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewParliamentModule(client)

		// Get parameters
		policyJSON := request.GetString("policy_json", "")
		if policyJSON == "" {
			return mcp.NewToolResultError("policy_json is required"), nil
		}

		// Note: Dagger module doesn't support config and json_output options for string linting
		if request.GetString("config", "") != "" || request.GetBool("json_output", false) {
			return mcp.NewToolResultError("config and json_output options not supported in Dagger module for string linting"), nil
		}

		// Lint policy string
		output, err := module.LintPolicyString(ctx, policyJSON)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parliament string lint failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewParliamentModule(client)

		// Get parameters
		policyPath := request.GetString("policy_path", "")
		if policyPath == "" {
			return mcp.NewToolResultError("policy_path is required"), nil
		}

		// Note: Dagger module doesn't support config and json_output options for community auditors
		if request.GetString("config", "") != "" || request.GetBool("json_output", false) {
			return mcp.NewToolResultError("config and json_output options not supported in Dagger module for community auditors"), nil
		}

		// Lint with community auditors
		output, err := module.LintWithCommunityAuditors(ctx, policyPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parliament community auditors lint failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewParliamentModule(client)

		// Get parameters
		policyPath := request.GetString("policy_path", "")
		if policyPath == "" {
			return mcp.NewToolResultError("policy_path is required"), nil
		}
		privateAuditors := request.GetString("private_auditors", "")
		if privateAuditors == "" {
			return mcp.NewToolResultError("private_auditors is required"), nil
		}

		// Note: Dagger module doesn't support config and json_output options for private auditors
		if request.GetString("config", "") != "" || request.GetBool("json_output", false) {
			return mcp.NewToolResultError("config and json_output options not supported in Dagger module for private auditors"), nil
		}

		// Lint with private auditors
		output, err := module.LintWithPrivateAuditors(ctx, policyPath, privateAuditors)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parliament private auditors lint failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewParliamentModule(client)

		// Get parameters
		config := request.GetString("config", "")
		jsonOutput := request.GetBool("json_output", false)

		// Lint AWS managed policies
		output, err := module.LintAWSManagedPolicies(ctx, config, jsonOutput)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parliament AWS managed policies lint failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewParliamentModule(client)

		// Get parameters
		authDetailsFile := request.GetString("auth_details_file", "")
		if authDetailsFile == "" {
			return mcp.NewToolResultError("auth_details_file is required"), nil
		}
		config := request.GetString("config", "")
		jsonOutput := request.GetBool("json_output", false)

		// Lint auth details file
		output, err := module.LintAuthDetailsFile(ctx, authDetailsFile, config, jsonOutput)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parliament auth details lint failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewParliamentModule(client)

		// Get parameters
		policyPath := request.GetString("policy_path", "")
		if policyPath == "" {
			return mcp.NewToolResultError("policy_path is required"), nil
		}
		privateAuditors := request.GetString("private_auditors", "")
		config := request.GetString("config", "")
		jsonOutput := request.GetBool("json_output", false)

		// Comprehensive analysis
		output, err := module.ComprehensiveAnalysis(ctx, policyPath, privateAuditors, config, jsonOutput)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parliament comprehensive analysis failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewParliamentModule(client)

		// Get parameters
		baseDirectory := request.GetString("base_directory", "")
		if baseDirectory == "" {
			return mcp.NewToolResultError("base_directory is required"), nil
		}
		privateAuditors := request.GetString("private_auditors", "")
		config := request.GetString("config", "")
		jsonOutput := request.GetBool("json_output", false)
		includeExtension := request.GetString("include_policy_extension", "")
		excludePattern := request.GetString("exclude_pattern", "")

		// Batch directory analysis
		output, err := module.BatchDirectoryAnalysis(ctx, baseDirectory, config, privateAuditors, jsonOutput, includeExtension, excludePattern)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("parliament batch directory analysis failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}