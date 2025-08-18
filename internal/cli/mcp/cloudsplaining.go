package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCloudsplainingTools adds Cloudsplaining (AWS IAM policy scanner) MCP tool implementations
func AddCloudsplainingTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Cloudsplaining download account data
	downloadTool := mcp.NewTool("cloudsplaining_download",
		mcp.WithDescription("Download AWS account authorization data"),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use for download"),
		),
		mcp.WithBoolean("include_non_default_policy_versions",
			mcp.Description("Include non-default policy versions"),
		),
	)
	s.AddTool(downloadTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloudsplaining", "download"}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if request.GetBool("include_non_default_policy_versions", false) {
			args = append(args, "--include-non-default-policy-versions")
		}
		
		return executeShipCommand(args)
	})

	// Cloudsplaining scan account data
	scanTool := mcp.NewTool("cloudsplaining_scan",
		mcp.WithDescription("Scan downloaded account authorization data"),
		mcp.WithString("input_file",
			mcp.Description("Path to downloaded account authorization data file"),
			mcp.Required(),
		),
		mcp.WithString("exclusions_file",
			mcp.Description("Path to exclusions file"),
		),
		mcp.WithString("output",
			mcp.Description("Output directory for scan results"),
		),
	)
	s.AddTool(scanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputFile := request.GetString("input_file", "")
		args := []string{"cloudsplaining", "scan", "--input-file", inputFile}
		
		if exclusionsFile := request.GetString("exclusions_file", ""); exclusionsFile != "" {
			args = append(args, "--exclusions-file", exclusionsFile)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		
		return executeShipCommand(args)
	})

	// Cloudsplaining scan policy file
	scanPolicyFileTool := mcp.NewTool("cloudsplaining_scan_policy_file",
		mcp.WithDescription("Scan a specific IAM policy file"),
		mcp.WithString("input_file",
			mcp.Description("Path to IAM policy JSON file"),
			mcp.Required(),
		),
		mcp.WithString("exclusions_file",
			mcp.Description("Path to exclusions file"),
		),
	)
	s.AddTool(scanPolicyFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputFile := request.GetString("input_file", "")
		args := []string{"cloudsplaining", "scan-policy-file", "--input-file", inputFile}
		
		if exclusionsFile := request.GetString("exclusions_file", ""); exclusionsFile != "" {
			args = append(args, "--exclusions-file", exclusionsFile)
		}
		
		return executeShipCommand(args)
	})

	// Cloudsplaining create exclusions file
	createExclusionsTool := mcp.NewTool("cloudsplaining_create_exclusions_file",
		mcp.WithDescription("Create exclusions file template"),
	)
	s.AddTool(createExclusionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cloudsplaining", "create-exclusions-file"}
		return executeShipCommand(args)
	})

	// Cloudsplaining create multi-account config
	createMultiAccountConfigTool := mcp.NewTool("cloudsplaining_create_multi_account_config",
		mcp.WithDescription("Create multi-account configuration file"),
		mcp.WithString("output_file",
			mcp.Description("Output file path for configuration"),
			mcp.Required(),
		),
	)
	s.AddTool(createMultiAccountConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		outputFile := request.GetString("output_file", "")
		args := []string{"cloudsplaining", "create-multi-account-config-file", "-o", outputFile}
		return executeShipCommand(args)
	})

	// Cloudsplaining scan multi-account
	scanMultiAccountTool := mcp.NewTool("cloudsplaining_scan_multi_account",
		mcp.WithDescription("Scan multiple AWS accounts"),
		mcp.WithString("config_file",
			mcp.Description("Path to multi-account configuration file"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("role_name",
			mcp.Description("IAM role name to assume"),
		),
		mcp.WithString("output_bucket",
			mcp.Description("S3 bucket for output"),
		),
		mcp.WithString("output_directory",
			mcp.Description("Local output directory"),
		),
	)
	s.AddTool(scanMultiAccountTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configFile := request.GetString("config_file", "")
		args := []string{"cloudsplaining", "scan-multi-account", "-c", configFile}
		
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if roleName := request.GetString("role_name", ""); roleName != "" {
			args = append(args, "--role-name", roleName)
		}
		if outputBucket := request.GetString("output_bucket", ""); outputBucket != "" {
			args = append(args, "--output-bucket", outputBucket)
		}
		if outputDirectory := request.GetString("output_directory", ""); outputDirectory != "" {
			args = append(args, "--output-directory", outputDirectory)
		}
		
		return executeShipCommand(args)
	})
}