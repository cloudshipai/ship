package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTfLintTools adds TFLint (Terraform linting) MCP tool implementations
func AddTfLintTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// TFLint lint directory tool
	lintDirectoryTool := mcp.NewTool("tflint_lint_directory",
		mcp.WithDescription("Lint all Terraform files in a directory using TFLint"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files to lint"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "compact", "checkstyle"),
		),
	)
	s.AddTool(lintDirectoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"terraform-tools", "tflint", directory}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// TFLint lint specific file tool
	lintFileTool := mcp.NewTool("tflint_lint_file",
		mcp.WithDescription("Lint a specific Terraform file using TFLint"),
		mcp.WithString("file_path",
			mcp.Description("Path to Terraform file to lint"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "compact", "checkstyle"),
		),
	)
	s.AddTool(lintFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"terraform-tools", "tflint", "--file", filePath}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// TFLint lint with custom config tool
	lintWithConfigTool := mcp.NewTool("tflint_lint_with_config",
		mcp.WithDescription("Lint Terraform files using custom TFLint configuration"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to TFLint configuration file"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "compact", "checkstyle"),
		),
	)
	s.AddTool(lintWithConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		configFile := request.GetString("config_file", "")
		args := []string{"terraform-tools", "tflint", directory, "--config", configFile}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// TFLint lint with rules tool
	lintWithRulesTool := mcp.NewTool("tflint_lint_with_rules",
		mcp.WithDescription("Lint Terraform files with specific rules enabled/disabled"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("enable_rules",
			mcp.Description("Comma-separated list of rules to enable"),
		),
		mcp.WithString("disable_rules",
			mcp.Description("Comma-separated list of rules to disable"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "compact", "checkstyle"),
		),
	)
	s.AddTool(lintWithRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"terraform-tools", "tflint", directory}
		if enableRules := request.GetString("enable_rules", ""); enableRules != "" {
			args = append(args, "--enable-rules", enableRules)
		}
		if disableRules := request.GetString("disable_rules", ""); disableRules != "" {
			args = append(args, "--disable-rules", disableRules)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// TFLint init plugins tool
	initPluginsTool := mcp.NewTool("tflint_init_plugins",
		mcp.WithDescription("Initialize TFLint plugins for a Terraform project"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
	)
	s.AddTool(initPluginsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"terraform-tools", "tflint", directory, "--init"}
		return executeShipCommand(args)
	})

	// TFLint validate format tool
	validateFormatTool := mcp.NewTool("tflint_validate_format",
		mcp.WithDescription("Validate Terraform format and syntax using TFLint"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Recursively lint subdirectories"),
		),
	)
	s.AddTool(validateFormatTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"terraform-tools", "tflint", directory, "--validate"}
		if request.GetBool("recursive", false) {
			args = append(args, "--recursive")
		}
		return executeShipCommand(args)
	})

	// TFLint fix issues tool
	fixIssuesTool := mcp.NewTool("tflint_fix_issues",
		mcp.WithDescription("Automatically fix TFLint issues where possible"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "compact", "checkstyle"),
		),
	)
	s.AddTool(fixIssuesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		args := []string{"terraform-tools", "tflint", directory, "--fix"}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// TFLint with var file tool
	lintWithVarFileTool := mcp.NewTool("tflint_lint_with_var_file",
		mcp.WithDescription("Lint Terraform files with variable files"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("var_file",
			mcp.Description("Path to Terraform variable file"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "compact", "checkstyle"),
		),
	)
	s.AddTool(lintWithVarFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		varFile := request.GetString("var_file", "")
		args := []string{"terraform-tools", "tflint", directory, "--var-file", varFile}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// TFLint with chdir tool
	lintWithChdirTool := mcp.NewTool("tflint_lint_with_chdir",
		mcp.WithDescription("Lint Terraform files after changing working directory"),
		mcp.WithString("target_directory",
			mcp.Description("Target directory to change to before linting"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "compact", "checkstyle"),
		),
	)
	s.AddTool(lintWithChdirTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		targetDir := request.GetString("target_directory", "")
		args := []string{"terraform-tools", "tflint", "--chdir", targetDir}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// TFLint only specific rule tool
	lintOnlyRuleTool := mcp.NewTool("tflint_lint_only_rule",
		mcp.WithDescription("Lint with only a specific rule enabled"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files"),
			mcp.Required(),
		),
		mcp.WithString("rule",
			mcp.Description("Specific rule to enable"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "compact", "checkstyle"),
		),
	)
	s.AddTool(lintOnlyRuleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory := request.GetString("directory", "")
		rule := request.GetString("rule", "")
		args := []string{"terraform-tools", "tflint", directory, "--only", rule}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// TFLint print version tool
	printVersionTool := mcp.NewTool("tflint_print_version",
		mcp.WithDescription("Print TFLint version information"),
	)
	s.AddTool(printVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform-tools", "tflint", "--version"}
		return executeShipCommand(args)
	})
}