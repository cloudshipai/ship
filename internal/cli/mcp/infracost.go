package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddInfracostTools adds Infracost MCP tool implementations using real CLI commands
func AddInfracostTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		path := request.GetString("path", "")
		args := []string{"infracost", "breakdown", "--path", path}
		
		if configFile := request.GetString("config_file", ""); configFile != "" {
			args = append(args, "--config-file", configFile)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if outFile := request.GetString("out_file", ""); outFile != "" {
			args = append(args, "--out-file", outFile)
		}
		if request.GetBool("show_skipped", false) {
			args = append(args, "--show-skipped")
		}
		if terraformVarFile := request.GetString("terraform_var_file", ""); terraformVarFile != "" {
			args = append(args, "--terraform-var-file", terraformVarFile)
		}
		if terraformVar := request.GetString("terraform_var", ""); terraformVar != "" {
			args = append(args, "--terraform-var", terraformVar)
		}
		if terraformWorkspace := request.GetString("terraform_workspace", ""); terraformWorkspace != "" {
			args = append(args, "--terraform-workspace", terraformWorkspace)
		}
		if usageFile := request.GetString("usage_file", ""); usageFile != "" {
			args = append(args, "--usage-file", usageFile)
		}
		
		return executeShipCommand(args)
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
		path := request.GetString("path", "")
		args := []string{"infracost", "diff", "--path", path}
		
		if compareTo := request.GetString("compare_to", ""); compareTo != "" {
			args = append(args, "--compare-to", compareTo)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if outFile := request.GetString("out_file", ""); outFile != "" {
			args = append(args, "--out-file", outFile)
		}
		if request.GetBool("show_skipped", false) {
			args = append(args, "--show-skipped")
		}
		if terraformVarFile := request.GetString("terraform_var_file", ""); terraformVarFile != "" {
			args = append(args, "--terraform-var-file", terraformVarFile)
		}
		if terraformVar := request.GetString("terraform_var", ""); terraformVar != "" {
			args = append(args, "--terraform-var", terraformVar)
		}
		if usageFile := request.GetString("usage_file", ""); usageFile != "" {
			args = append(args, "--usage-file", usageFile)
		}
		
		return executeShipCommand(args)
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
		path := request.GetString("path", "")
		args := []string{"infracost", "output", "--path", path}
		
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if outFile := request.GetString("out_file", ""); outFile != "" {
			args = append(args, "--out-file", outFile)
		}
		if request.GetBool("show_skipped", false) {
			args = append(args, "--show-skipped")
		}
		if request.GetBool("show_all_projects", false) {
			args = append(args, "--show-all-projects")
		}
		
		return executeShipCommand(args)
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
		path := request.GetString("path", "")
		args := []string{"infracost", "upload", "--path", path}
		
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		
		return executeShipCommand(args)
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
		args := []string{"infracost", "configure", "set", setting, value}
		
		return executeShipCommand(args)
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
		repoPath := request.GetString("repo_path", "")
		templatePath := request.GetString("template_path", "")
		args := []string{"infracost", "generate", "config", "--repo-path", repoPath, "--template-path", templatePath}
		
		return executeShipCommand(args)
	})

	// Infracost auth login tool
	authTool := mcp.NewTool("infracost_auth_login",
		mcp.WithDescription("Get a free API key or log in to existing account using infracost auth login"),
	)
	s.AddTool(authTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"infracost", "auth", "login"}
		return executeShipCommand(args)
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
		path := request.GetString("path", "")
		repo := request.GetString("repo", "")
		pr := request.GetString("pull_request", "")
		
		args := []string{"infracost", "comment", "github", "--path", path, "--repo", repo, "--pull-request", pr}
		
		if token := request.GetString("github_token", ""); token != "" {
			args = append(args, "--github-token", token)
		}
		if behavior := request.GetString("behavior", ""); behavior != "" {
			args = append(args, "--behavior", behavior)
		}
		
		return executeShipCommand(args)
	})
}

