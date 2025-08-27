package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddTfLintTools adds TFLint (Terraform linter) MCP tool implementations using direct Dagger calls
func AddTfLintTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTfLintToolsDirect(s)
}

// addTfLintToolsDirect adds TFLint tools using direct Dagger module calls
func addTfLintToolsDirect(s *server.MCPServer) {
	// TFLint check tool
	checkTool := mcp.NewTool("tflint_check",
		mcp.WithDescription("Run TFLint to check Terraform configuration for issues"),
		mcp.WithString("source_path",
			mcp.Description("Path to Terraform configuration directory"),
			mcp.Required(),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to TFLint configuration file"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("default", "json", "checkstyle", "junit", "compact", "sarif"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Run recursively over subdirectories"),
		),
		mcp.WithString("enable_rule",
			mcp.Description("Enable a specific rule"),
		),
		mcp.WithString("disable_rule",
			mcp.Description("Disable a specific rule"),
		),
		mcp.WithString("only",
			mcp.Description("Run only specified rules"),
		),
		mcp.WithString("var_file",
			mcp.Description("Terraform variables file"),
		),
		mcp.WithString("var",
			mcp.Description("Set Terraform variables"),
		),
		mcp.WithBoolean("fix",
			mcp.Description("Automatically fix issues where possible"),
		),
	)
	s.AddTool(checkTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		sourcePath := request.GetString("source_path", "")
		configFile := request.GetString("config_file", "")
		format := request.GetString("format", "")
		recursive := request.GetBool("recursive", false)
		enableRule := request.GetString("enable_rule", "")
		disableRule := request.GetString("disable_rule", "")
		only := request.GetString("only", "")
		varFile := request.GetString("var_file", "")
		varValue := request.GetString("var", "")
		fix := request.GetBool("fix", false)

		if sourcePath == "" {
			return mcp.NewToolResultError("source_path is required"), nil
		}

		// Create TFLint module
		tflintModule := modules.NewTFLintModule(client)

		// Set up options
		opts := modules.TFLintOptions{
			ConfigFile:  configFile,
			Format:      format,
			Recursive:   recursive,
			EnableRule:  enableRule,
			DisableRule: disableRule,
			Only:        only,
			VarFile:     varFile,
			Var:         varValue,
			Fix:         fix,
		}

		// Run TFLint
		result, err := tflintModule.Check(ctx, sourcePath, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TFLint check failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// TFLint init tool
	initTool := mcp.NewTool("tflint_init",
		mcp.WithDescription("Initialize TFLint in a Terraform configuration directory"),
		mcp.WithString("source_path",
			mcp.Description("Path to Terraform configuration directory"),
			mcp.Required(),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to TFLint configuration file"),
		),
		mcp.WithBoolean("upgrade",
			mcp.Description("Upgrade plugins to the latest available version"),
		),
	)
	s.AddTool(initTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		sourcePath := request.GetString("source_path", "")
		configFile := request.GetString("config_file", "")
		upgrade := request.GetBool("upgrade", false)

		if sourcePath == "" {
			return mcp.NewToolResultError("source_path is required"), nil
		}

		// Create TFLint module
		tflintModule := modules.NewTFLintModule(client)

		// Set up options
		opts := modules.TFLintInitOptions{
			ConfigFile: configFile,
			Upgrade:    upgrade,
		}

		// Initialize TFLint
		result, err := tflintModule.Init(ctx, sourcePath, opts)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TFLint init failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}