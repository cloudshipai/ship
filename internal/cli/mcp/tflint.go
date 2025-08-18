package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTfLintTools adds TFLint MCP tool implementations using real tflint CLI commands
func AddTfLintTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// TFLint basic lint tool
	lintTool := mcp.NewTool("tflint_lint",
		mcp.WithDescription("Lint Terraform files using real tflint CLI"),
		mcp.WithString("chdir",
			mcp.Description("Change working directory before linting"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("default", "json", "checkstyle", "junit", "compact", "sarif"),
		),
		mcp.WithString("config",
			mcp.Description("Path to TFLint configuration file"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Run command in each directory recursively"),
		),
	)
	s.AddTool(lintTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tflint"}
		
		if chdir := request.GetString("chdir", ""); chdir != "" {
			args = append(args, "--chdir="+chdir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format="+format)
		}
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config="+config)
		}
		if request.GetBool("recursive", false) {
			args = append(args, "--recursive")
		}
		
		return executeShipCommand(args)
	})

	// TFLint lint with rules tool
	lintWithRulesTool := mcp.NewTool("tflint_lint_with_rules",
		mcp.WithDescription("Lint Terraform files with specific rules enabled/disabled using real tflint CLI"),
		mcp.WithString("chdir",
			mcp.Description("Change working directory before linting"),
		),
		mcp.WithString("enable_rule",
			mcp.Description("Enable specific rule"),
		),
		mcp.WithString("disable_rule",
			mcp.Description("Disable specific rule"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("default", "json", "checkstyle", "junit", "compact", "sarif"),
		),
	)
	s.AddTool(lintWithRulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tflint"}
		
		if chdir := request.GetString("chdir", ""); chdir != "" {
			args = append(args, "--chdir="+chdir)
		}
		if enableRule := request.GetString("enable_rule", ""); enableRule != "" {
			args = append(args, "--enable-rule="+enableRule)
		}
		if disableRule := request.GetString("disable_rule", ""); disableRule != "" {
			args = append(args, "--disable-rule="+disableRule)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format="+format)
		}
		
		return executeShipCommand(args)
	})

	// TFLint init plugins tool
	initTool := mcp.NewTool("tflint_init",
		mcp.WithDescription("Initialize TFLint plugins using real tflint CLI"),
		mcp.WithString("chdir",
			mcp.Description("Change working directory before initializing"),
		),
	)
	s.AddTool(initTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tflint", "--init"}
		
		if chdir := request.GetString("chdir", ""); chdir != "" {
			args = append(args, "--chdir="+chdir)
		}
		
		return executeShipCommand(args)
	})

	// TFLint with variable file tool
	lintWithVarFileTool := mcp.NewTool("tflint_lint_with_var_file",
		mcp.WithDescription("Lint Terraform files with variable files using real tflint CLI"),
		mcp.WithString("chdir",
			mcp.Description("Change working directory before linting"),
		),
		mcp.WithString("var_file",
			mcp.Description("Path to Terraform variable file"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("default", "json", "checkstyle", "junit", "compact", "sarif"),
		),
	)
	s.AddTool(lintWithVarFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		varFile := request.GetString("var_file", "")
		args := []string{"tflint", "--var-file=" + varFile}
		
		if chdir := request.GetString("chdir", ""); chdir != "" {
			args = append(args, "--chdir="+chdir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format="+format)
		}
		
		return executeShipCommand(args)
	})

	// TFLint with variables tool
	lintWithVarTool := mcp.NewTool("tflint_lint_with_var",
		mcp.WithDescription("Lint Terraform files with individual variables using real tflint CLI"),
		mcp.WithString("chdir",
			mcp.Description("Change working directory before linting"),
		),
		mcp.WithString("var",
			mcp.Description("Set Terraform variable (format: 'key=value')"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("default", "json", "checkstyle", "junit", "compact", "sarif"),
		),
	)
	s.AddTool(lintWithVarTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		variable := request.GetString("var", "")
		args := []string{"tflint", "--var=" + variable}
		
		if chdir := request.GetString("chdir", ""); chdir != "" {
			args = append(args, "--chdir="+chdir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format="+format)
		}
		
		return executeShipCommand(args)
	})

	// TFLint version tool
	versionTool := mcp.NewTool("tflint_version",
		mcp.WithDescription("Get TFLint version information using real tflint CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tflint", "--version"}
		return executeShipCommand(args)
	})
}