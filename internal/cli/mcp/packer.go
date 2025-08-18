package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddPackerTools adds Packer (machine image building) MCP tool implementations using real CLI commands
func AddPackerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		templateFile := request.GetString("template_file", "")
		args := []string{"packer", "build"}
		
		if vars := request.GetString("var", ""); vars != "" {
			args = append(args, "-var", vars)
		}
		if varFile := request.GetString("var_file", ""); varFile != "" {
			args = append(args, "-var-file", varFile)
		}
		if only := request.GetString("only", ""); only != "" {
			args = append(args, "-only", only)
		}
		if except := request.GetString("except", ""); except != "" {
			args = append(args, "-except", except)
		}
		if request.GetBool("force", false) {
			args = append(args, "-force")
		}
		if request.GetBool("debug", false) {
			args = append(args, "-debug")
		}
		if parallelBuilds := request.GetString("parallel_builds", ""); parallelBuilds != "" {
			args = append(args, "-parallel-builds", parallelBuilds)
		}
		if onError := request.GetString("on_error", ""); onError != "" {
			args = append(args, "-on-error", onError)
		}
		
		args = append(args, templateFile)
		return executeShipCommand(args)
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
		templateFile := request.GetString("template_file", "")
		args := []string{"packer", "validate"}
		
		if vars := request.GetString("var", ""); vars != "" {
			args = append(args, "-var", vars)
		}
		if varFile := request.GetString("var_file", ""); varFile != "" {
			args = append(args, "-var-file", varFile)
		}
		if only := request.GetString("only", ""); only != "" {
			args = append(args, "-only", only)
		}
		if except := request.GetString("except", ""); except != "" {
			args = append(args, "-except", except)
		}
		if request.GetBool("syntax_only", false) {
			args = append(args, "-syntax-only")
		}
		if request.GetBool("evaluate_datasources", false) {
			args = append(args, "-evaluate-datasources")
		}
		
		args = append(args, templateFile)
		return executeShipCommand(args)
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
		templateFile := request.GetString("template_file", "")
		args := []string{"packer", "inspect"}
		
		if request.GetBool("machine_readable", false) {
			args = append(args, "-machine-readable")
		}
		
		args = append(args, templateFile)
		return executeShipCommand(args)
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
		templateFile := request.GetString("template_file", "")
		args := []string{"packer", "fix"}
		
		if request.GetBool("validate", false) {
			args = append(args, "-validate")
		}
		
		args = append(args, templateFile)
		return executeShipCommand(args)
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
		templateFile := request.GetString("template_file", "")
		args := []string{"packer", "console"}
		
		if vars := request.GetString("var", ""); vars != "" {
			args = append(args, "-var", vars)
		}
		if varFile := request.GetString("var_file", ""); varFile != "" {
			args = append(args, "-var-file", varFile)
		}
		
		args = append(args, templateFile)
		return executeShipCommand(args)
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
		args := []string{"packer", "fmt"}
		
		if request.GetBool("check", false) {
			args = append(args, "-check")
		}
		if request.GetBool("diff", false) {
			args = append(args, "-diff")
		}
		if request.GetBool("write", true) {
			args = append(args, "-write")
		}
		if request.GetBool("recursive", false) {
			args = append(args, "-recursive")
		}
		
		path := request.GetString("path", ".")
		args = append(args, path)
		
		return executeShipCommand(args)
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
		configFile := request.GetString("config_file", "")
		args := []string{"packer", "init"}
		
		if request.GetBool("upgrade", false) {
			args = append(args, "-upgrade")
		}
		
		args = append(args, configFile)
		return executeShipCommand(args)
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
		subcommand := request.GetString("subcommand", "")
		args := []string{"packer", "plugins", subcommand}
		
		switch subcommand {
		case "install", "remove":
			if pluginName := request.GetString("plugin_name", ""); pluginName != "" {
				args = append(args, pluginName)
			}
			if version := request.GetString("version", ""); version != "" {
				args = append(args, version)
			}
		case "required":
			if configFile := request.GetString("config_file", ""); configFile != "" {
				args = append(args, configFile)
			}
		}
		
		return executeShipCommand(args)
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
		templateFile := request.GetString("template_file", "")
		args := []string{"packer", "hcl2_upgrade"}
		
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-output-file", outputFile)
		}
		if request.GetBool("with_annotations", false) {
			args = append(args, "-with-annotations")
		}
		
		args = append(args, templateFile)
		return executeShipCommand(args)
	})

	// Packer version tool
	versionTool := mcp.NewTool("packer_version",
		mcp.WithDescription("Get Packer version information using real packer CLI"),
		mcp.WithBoolean("machine_readable",
			mcp.Description("Output in machine-readable format"),
		),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"packer", "version"}
		
		if request.GetBool("machine_readable", false) {
			args = []string{"packer", "-machine-readable", "version"}
		}
		
		return executeShipCommand(args)
	})
}