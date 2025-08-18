package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddPackerTools adds Packer (machine image building) MCP tool implementations using direct Dagger calls
func AddPackerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addPackerToolsDirect(s)
}

// addPackerToolsDirect adds Packer tools using direct Dagger module calls
func addPackerToolsDirect(s *server.MCPServer) {
	// Packer build tool
	buildTool := mcp.NewTool("packer_build",
		mcp.WithDescription("Build machine images using real packer CLI"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
		mcp.WithString("var",
			mcp.Description("Set template variable (key=value format)"),
		),
		mcp.WithString("var_file",
			mcp.Description("Path to variables file"),
		),
		mcp.WithString("only",
			mcp.Description("Comma-separated list of builds to run only"),
		),
		mcp.WithString("except",
			mcp.Description("Comma-separated list of builds to skip"),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force build even with existing artifacts"),
		),
		mcp.WithBoolean("debug",
			mcp.Description("Enable debug mode with step-by-step pausing"),
		),
		mcp.WithString("parallel_builds",
			mcp.Description("Number of parallel builds (0 for no limit)"),
		),
		mcp.WithString("on_error",
			mcp.Description("Action on build error"),
			mcp.Enum("cleanup", "abort", "ask"),
		),
	)
	s.AddTool(buildTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPackerModule(client)

		// Get parameters
		templateFile := request.GetString("template_file", "")
		if templateFile == "" {
			return mcp.NewToolResultError("template_file is required"), nil
		}
		varFile := request.GetString("var_file", "")

		// Note: Dagger module only supports basic build with template and var-file
		// Advanced options like only/except/force/debug are not supported
		if request.GetString("only", "") != "" || request.GetString("except", "") != "" ||
			request.GetBool("force", false) || request.GetBool("debug", false) ||
			request.GetString("parallel_builds", "") != "" || request.GetString("on_error", "") != "" {
			return mcp.NewToolResultError("advanced build options not supported in Dagger module"), nil
		}

		// Build image
		output, err := module.BuildImage(ctx, templateFile, varFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer build failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Packer validate tool
	validateTool := mcp.NewTool("packer_validate",
		mcp.WithDescription("Validate Packer configuration template using real packer CLI"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
		mcp.WithString("var",
			mcp.Description("Set template variable (key=value format)"),
		),
		mcp.WithString("var_file",
			mcp.Description("Path to variables file"),
		),
		mcp.WithString("only",
			mcp.Description("Comma-separated list of builds to validate only"),
		),
		mcp.WithString("except",
			mcp.Description("Comma-separated list of builds to skip validation"),
		),
		mcp.WithBoolean("syntax_only",
			mcp.Description("Check syntax only without validating configuration"),
		),
		mcp.WithBoolean("evaluate_datasources",
			mcp.Description("Evaluate all data sources (HCL2 templates only)"),
		),
	)
	s.AddTool(validateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPackerModule(client)

		// Get parameters
		templateFile := request.GetString("template_file", "")
		if templateFile == "" {
			return mcp.NewToolResultError("template_file is required"), nil
		}

		// Note: Dagger module only supports basic validate
		// Advanced options like vars/only/except/syntax_only are not supported
		if request.GetString("var", "") != "" || request.GetString("var_file", "") != "" ||
			request.GetString("only", "") != "" || request.GetString("except", "") != "" ||
			request.GetBool("syntax_only", false) || request.GetBool("evaluate_datasources", false) {
			return mcp.NewToolResultError("advanced validate options not supported in Dagger module"), nil
		}

		// Validate template
		output, err := module.ValidateTemplate(ctx, templateFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer validate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Packer inspect tool
	inspectTool := mcp.NewTool("packer_inspect",
		mcp.WithDescription("Inspect and analyze Packer template configuration using real packer CLI"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
		mcp.WithBoolean("machine_readable",
			mcp.Description("Output in machine-readable format"),
		),
	)
	s.AddTool(inspectTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPackerModule(client)

		// Get parameters
		templateFile := request.GetString("template_file", "")
		if templateFile == "" {
			return mcp.NewToolResultError("template_file is required"), nil
		}
		machineReadable := request.GetBool("machine_readable", false)

		// Inspect template
		output, err := module.InspectTemplate(ctx, templateFile, machineReadable)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer inspect failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Packer fix tool
	fixTool := mcp.NewTool("packer_fix",
		mcp.WithDescription("Fix and upgrade Packer template to current version using real packer CLI"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
		mcp.WithBoolean("validate",
			mcp.Description("Validate template after fixing"),
		),
	)
	s.AddTool(fixTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPackerModule(client)

		// Get parameters
		templateFile := request.GetString("template_file", "")
		if templateFile == "" {
			return mcp.NewToolResultError("template_file is required"), nil
		}
		validateFlag := request.GetBool("validate", false)

		// Fix template
		output, err := module.FixTemplate(ctx, templateFile, validateFlag)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer fix failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Packer console tool
	consoleTool := mcp.NewTool("packer_console",
		mcp.WithDescription("Open Packer console for template debugging using real packer CLI"),
		mcp.WithString("template_file",
			mcp.Description("Path to Packer template file"),
			mcp.Required(),
		),
		mcp.WithString("var",
			mcp.Description("Set template variable (key=value format)"),
		),
		mcp.WithString("var_file",
			mcp.Description("Path to variables file"),
		),
	)
	s.AddTool(consoleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPackerModule(client)

		// Get parameters
		templateFile := request.GetString("template_file", "")
		if templateFile == "" {
			return mcp.NewToolResultError("template_file is required"), nil
		}
		vars := request.GetString("var", "")
		varFile := request.GetString("var_file", "")

		// Open console
		output, err := module.Console(ctx, templateFile, vars, varFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer console failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Packer fmt tool
	fmtTool := mcp.NewTool("packer_fmt",
		mcp.WithDescription("Format Packer template files using real packer CLI"),
		mcp.WithString("path",
			mcp.Description("Path to template file or directory (default: current directory)"),
		),
		mcp.WithBoolean("check",
			mcp.Description("Check if input is formatted without writing output"),
		),
		mcp.WithBoolean("diff",
			mcp.Description("Display diffs instead of rewriting files"),
		),
		mcp.WithBoolean("write",
			mcp.Description("Write result to source file (default: true)"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Process directory recursively"),
		),
	)
	s.AddTool(fmtTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPackerModule(client)

		// Get parameters
		path := request.GetString("path", ".")

		// Note: Dagger module only supports basic format
		// Advanced options like check/diff/write/recursive are not supported
		if request.GetBool("check", false) || request.GetBool("diff", false) ||
			request.GetBool("write", true) || request.GetBool("recursive", false) {
			return mcp.NewToolResultError("advanced format options not supported in Dagger module"), nil
		}

		// Format template (assuming path is a file)
		output, err := module.FormatTemplate(ctx, path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer fmt failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Packer init tool
	initTool := mcp.NewTool("packer_init",
		mcp.WithDescription("Initialize Packer configuration and install required plugins using real packer CLI"),
		mcp.WithString("config_file",
			mcp.Description("Path to Packer configuration file"),
			mcp.Required(),
		),
		mcp.WithBoolean("upgrade",
			mcp.Description("Upgrade plugins to latest compatible versions"),
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
		module := modules.NewPackerModule(client)

		// Get parameters
		configFile := request.GetString("config_file", "")
		if configFile == "" {
			return mcp.NewToolResultError("config_file is required"), nil
		}
		upgrade := request.GetBool("upgrade", false)

		// Initialize configuration
		output, err := module.InitConfiguration(ctx, configFile, upgrade)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer init failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Packer plugins tool
	pluginsTool := mcp.NewTool("packer_plugins",
		mcp.WithDescription("Manage Packer plugins using real packer CLI"),
		mcp.WithString("subcommand",
			mcp.Description("Plugin subcommand"),
			mcp.Enum("install", "remove", "required"),
			mcp.Required(),
		),
		mcp.WithString("plugin_name",
			mcp.Description("Name of the plugin (for install/remove commands)"),
		),
		mcp.WithString("version",
			mcp.Description("Plugin version constraint"),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to Packer configuration file (for required command)"),
		),
	)
	s.AddTool(pluginsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPackerModule(client)

		// Get parameters
		subcommand := request.GetString("subcommand", "")
		if subcommand == "" {
			return mcp.NewToolResultError("subcommand is required"), nil
		}
		pluginName := request.GetString("plugin_name", "")
		version := request.GetString("version", "")
		configFile := request.GetString("config_file", "")

		// Manage plugins
		output, err := module.ManagePlugins(ctx, subcommand, pluginName, version, configFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer plugins failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Packer hcl2_upgrade tool
	hcl2UpgradeTool := mcp.NewTool("packer_hcl2_upgrade",
		mcp.WithDescription("Upgrade JSON Packer template to HCL2 using real packer CLI"),
		mcp.WithString("template_file",
			mcp.Description("Path to JSON Packer template file"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for HCL2 template"),
		),
		mcp.WithBoolean("with_annotations",
			mcp.Description("Include helpful comments in output"),
		),
	)
	s.AddTool(hcl2UpgradeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPackerModule(client)

		// Get parameters
		templateFile := request.GetString("template_file", "")
		if templateFile == "" {
			return mcp.NewToolResultError("template_file is required"), nil
		}
		outputFile := request.GetString("output_file", "")
		withAnnotations := request.GetBool("with_annotations", false)

		// Upgrade to HCL2
		output, err := module.HCL2Upgrade(ctx, templateFile, outputFile, withAnnotations)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer hcl2_upgrade failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Packer version tool
	versionTool := mcp.NewTool("packer_version",
		mcp.WithDescription("Get Packer version information using real packer CLI"),
		mcp.WithBoolean("machine_readable",
			mcp.Description("Output in machine-readable format"),
		),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewPackerModule(client)

		// Note: machine_readable option not supported in Dagger module
		if request.GetBool("machine_readable", false) {
			return mcp.NewToolResultError("machine_readable option not supported in Dagger module"), nil
		}

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("packer version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}