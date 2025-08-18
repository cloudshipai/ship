package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddTfLintTools adds TFLint MCP tool implementations using direct Dagger calls
func AddTfLintTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTfLintToolsDirect(s)
}

// addTfLintToolsDirect adds TFLint tools using direct Dagger module calls
func addTfLintToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTFLintModule(client)

		// Get parameters
		chdir := request.GetString("chdir", ".")
		config := request.GetString("config", "")

		// Note: format and recursive parameters not supported in current Dagger implementation
		if request.GetString("format", "") != "" {
			return mcp.NewToolResultError("Warning: format parameter is not supported with direct Dagger calls"), nil
		}
		if request.GetBool("recursive", false) {
			return mcp.NewToolResultError("Warning: recursive parameter is not supported with direct Dagger calls"), nil
		}

		// Use LintWithConfig if config is provided, otherwise LintDirectory
		var output string
		if config != "" {
			output, err = module.LintWithConfig(ctx, chdir, config)
		} else {
			output, err = module.LintDirectory(ctx, chdir)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TFLint linting failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTFLintModule(client)

		// Get parameters
		chdir := request.GetString("chdir", ".")
		enableRule := request.GetString("enable_rule", "")
		disableRule := request.GetString("disable_rule", "")

		// Build rule arrays
		var enableRules, disableRules []string
		if enableRule != "" {
			enableRules = []string{enableRule}
		}
		if disableRule != "" {
			disableRules = []string{disableRule}
		}

		// Lint with rules
		output, err := module.LintWithRules(ctx, chdir, enableRules, disableRules)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TFLint rules linting failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// TFLint init plugins tool
	initTool := mcp.NewTool("tflint_init",
		mcp.WithDescription("Initialize TFLint plugins using real tflint CLI"),
		mcp.WithString("chdir",
			mcp.Description("Change working directory before initializing"),
		),
	)
	s.AddTool(initTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTFLintModule(client)

		// Get parameters
		chdir := request.GetString("chdir", ".")

		// Initialize plugins
		err = module.InitPlugins(ctx, chdir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TFLint init failed: %v", err)), nil
		}

		return mcp.NewToolResultText("TFLint plugins initialized successfully"), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTFLintModule(client)

		// Get parameters
		chdir := request.GetString("chdir", ".")
		varFile := request.GetString("var_file", "")
		format := request.GetString("format", "json")

		// Lint with variable file
		output, err := module.LintWithVarFile(ctx, chdir, varFile, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TFLint var file linting failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTFLintModule(client)

		// Get parameters
		chdir := request.GetString("chdir", ".")
		variable := request.GetString("var", "")
		format := request.GetString("format", "json")

		// Lint with variable
		output, err := module.LintWithVar(ctx, chdir, variable, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TFLint variable linting failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// TFLint version tool
	versionTool := mcp.NewTool("tflint_version",
		mcp.WithDescription("Get TFLint version information using real tflint CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTFLintModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TFLint get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}