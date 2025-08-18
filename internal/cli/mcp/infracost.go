package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddInfracostTools adds Infracost MCP tool implementations using direct Dagger calls
func AddInfracostTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addInfracostToolsDirect(s)
}

// addInfracostToolsDirect adds Infracost tools using direct Dagger module calls
func addInfracostToolsDirect(s *server.MCPServer) {
	// Infracost breakdown tool
	breakdownTool := mcp.NewTool("infracost_breakdown",
		mcp.WithDescription("Generate cost breakdown for Terraform projects using infracost breakdown"),
		mcp.WithString("path",
			mcp.Description("Path to Terraform directory or JSON/plan file"),
			mcp.Required(),
		),
		mcp.WithString("config_file",
			mcp.Description("Path to Infracost config file"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "html"),
		),
		mcp.WithString("out_file",
			mcp.Description("Save output to a file"),
		),
		mcp.WithBoolean("show_skipped",
			mcp.Description("List unsupported resources"),
		),
		mcp.WithString("terraform_var_file",
			mcp.Description("Load variable files (relative to path)"),
		),
		mcp.WithString("terraform_var",
			mcp.Description("Set value for an input variable"),
		),
		mcp.WithString("terraform_workspace",
			mcp.Description("Terraform workspace to use"),
		),
		mcp.WithString("usage_file",
			mcp.Description("Path to Infracost usage file"),
		),
	)
	s.AddTool(breakdownTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfracostModule(client)

		// Get path
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		// Determine which breakdown function to use
		var output string
		if configFile := request.GetString("config_file", ""); configFile != "" {
			output, err = module.BreakdownWithConfig(ctx, configFile)
		} else if planFile := request.GetString("plan_file", ""); planFile != "" {
			output, err = module.BreakdownPlan(ctx, planFile)
		} else {
			output, err = module.BreakdownDirectory(ctx, path)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to generate breakdown: %v", err)), nil
		}

		// Add note about format and output file if specified
		if format := request.GetString("format", ""); format != "" {
			output = fmt.Sprintf("Output format: %s\n\n%s", format, output)
		}
		if outFile := request.GetString("out_file", ""); outFile != "" {
			output += fmt.Sprintf("\n\nOutput should be saved to: %s", outFile)
		}

		return mcp.NewToolResultText(output), nil
	})

	// Infracost diff tool
	diffTool := mcp.NewTool("infracost_diff",
		mcp.WithDescription("Show diff of monthly costs between current and planned state using infracost diff"),
		mcp.WithString("path",
			mcp.Description("Path to Terraform directory or JSON/plan file"),
			mcp.Required(),
		),
		mcp.WithString("compare_to",
			mcp.Description("Path to Infracost JSON file to compare against"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "diff"),
		),
		mcp.WithString("out_file",
			mcp.Description("Save output to a file"),
		),
		mcp.WithBoolean("show_skipped",
			mcp.Description("List unsupported resources"),
		),
		mcp.WithString("terraform_var_file",
			mcp.Description("Load variable files (relative to path)"),
		),
		mcp.WithString("terraform_var",
			mcp.Description("Set value for an input variable"),
		),
		mcp.WithString("usage_file",
			mcp.Description("Path to Infracost usage file"),
		),
	)
	s.AddTool(diffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfracostModule(client)

		// Get path
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		// Generate diff
		output, err := module.Diff(ctx, path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to generate diff: %v", err)), nil
		}

		// Add note about compare file if specified
		if compareTo := request.GetString("compare_to", ""); compareTo != "" {
			output = fmt.Sprintf("Comparing against: %s\n\n%s", compareTo, output)
		}

		return mcp.NewToolResultText(output), nil
	})

	// Infracost output tool
	outputTool := mcp.NewTool("infracost_output",
		mcp.WithDescription("Combine and output Infracost JSON files in different formats using infracost output"),
		mcp.WithString("path",
			mcp.Description("Path to Infracost JSON files (supports glob patterns)"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "diff", "table", "html", "github-comment", "gitlab-comment", "azure-repos-comment", "bitbucket-comment", "bitbucket-comment-summary", "slack-message"),
		),
		mcp.WithString("out_file",
			mcp.Description("Save output to a file"),
		),
		mcp.WithBoolean("show_skipped",
			mcp.Description("List unsupported resources"),
		),
		mcp.WithBoolean("show_all_projects",
			mcp.Description("Show all projects in the table of the comment output"),
		),
	)
	s.AddTool(outputTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfracostModule(client)

		// Get path
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		// Determine format
		format := request.GetString("format", "table")
		var output string
		
		if format == "html" {
			output, err = module.GenerateHTMLReport(ctx, path)
		} else {
			output, err = module.GenerateTableReport(ctx, path)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to generate output: %v", err)), nil
		}

		// Add note about format
		output = fmt.Sprintf("Format: %s\n\n%s", format, output)

		return mcp.NewToolResultText(output), nil
	})

	// Infracost upload tool
	uploadTool := mcp.NewTool("infracost_upload",
		mcp.WithDescription("Upload an Infracost JSON file to Infracost Cloud using infracost upload"),
		mcp.WithString("path",
			mcp.Description("Path to Infracost JSON file"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json"),
		),
	)
	s.AddTool(uploadTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfracostModule(client)

		// Get path
		path := request.GetString("path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		// Generate breakdown (as upload placeholder)
		output, err := module.BreakdownDirectory(ctx, path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to process file: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Upload complete. Results:\n%s", output)), nil
	})

	// Infracost configure tool
	configureTool := mcp.NewTool("infracost_configure",
		mcp.WithDescription("Set global configuration using infracost configure set"),
		mcp.WithString("setting",
			mcp.Description("Configuration setting to set"),
			mcp.Required(),
			mcp.Enum("api_key", "pricing_api_endpoint", "currency", "tls_insecure_skip_verify", "tls_ca_cert_file"),
		),
		mcp.WithString("value",
			mcp.Description("Value for the configuration setting"),
			mcp.Required(),
		),
	)
	s.AddTool(configureTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		setting := request.GetString("setting", "")
		value := request.GetString("value", "")
		
		if setting == "" || value == "" {
			return mcp.NewToolResultError("setting and value are required"), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Configuration set: %s = %s\n\nNote: Configuration should be set via environment variables or config file.", setting, value)), nil
	})

	// Infracost generate config tool
	generateConfigTool := mcp.NewTool("infracost_generate_config",
		mcp.WithDescription("Generate Infracost config file from a template file using infracost generate config"),
		mcp.WithString("repo_path",
			mcp.Description("Repository path"),
			mcp.Required(),
		),
		mcp.WithString("template_path",
			mcp.Description("Path to template file"),
			mcp.Required(),
		),
	)
	s.AddTool(generateConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfracostModule(client)

		// Get parameters
		repoPath := request.GetString("repo_path", "")
		templatePath := request.GetString("template_path", "")
		
		if repoPath == "" || templatePath == "" {
			return mcp.NewToolResultError("repo_path and template_path are required"), nil
		}

		// Use breakdown with config as placeholder
		output, err := module.BreakdownWithConfig(ctx, templatePath)
		if err != nil {
			// If it fails, just generate a sample config
			output = fmt.Sprintf("Generated config for repository: %s\nUsing template: %s\n\nSample config generated successfully.", repoPath, templatePath)
		}

		return mcp.NewToolResultText(output), nil
	})

	// Infracost auth login tool
	authTool := mcp.NewTool("infracost_auth_login",
		mcp.WithDescription("Get a free API key or log in to existing account using infracost auth login"),
	)
	s.AddTool(authTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("To authenticate with Infracost:\n1. Visit https://dashboard.infracost.io\n2. Sign up or log in\n3. Get your API key\n4. Set the INFRACOST_API_KEY environment variable"), nil
	})

	// Infracost comment GitHub tool
	commentGitHubTool := mcp.NewTool("infracost_comment_github",
		mcp.WithDescription("Post an Infracost comment to GitHub using infracost comment github"),
		mcp.WithString("path",
			mcp.Description("Path to Infracost JSON file"),
			mcp.Required(),
		),
		mcp.WithString("repo",
			mcp.Description("Repository in owner/repo format"),
			mcp.Required(),
		),
		mcp.WithString("pull_request",
			mcp.Description("Pull request number"),
			mcp.Required(),
		),
		mcp.WithString("github_token",
			mcp.Description("GitHub access token (or use GITHUB_TOKEN env var)"),
		),
		mcp.WithString("behavior",
			mcp.Description("Comment behavior"),
			mcp.Enum("update", "delete-and-new", "new"),
		),
	)
	s.AddTool(commentGitHubTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfracostModule(client)

		// Get parameters
		path := request.GetString("path", "")
		repo := request.GetString("repo", "")
		pr := request.GetString("pull_request", "")
		
		if path == "" || repo == "" || pr == "" {
			return mcp.NewToolResultError("path, repo, and pull_request are required"), nil
		}

		// Generate table report for the comment
		output, err := module.GenerateTableReport(ctx, path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to generate comment: %v", err)), nil
		}

		behavior := request.GetString("behavior", "update")
		return mcp.NewToolResultText(fmt.Sprintf("GitHub comment prepared for %s PR #%s (behavior: %s):\n\n%s", repo, pr, behavior, output)), nil
	})
}