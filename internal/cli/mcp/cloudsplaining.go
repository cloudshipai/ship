package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCloudsplainingTools adds Cloudsplaining (AWS IAM policy scanner) MCP tool implementations
func AddCloudsplainingTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addCloudsplainingToolsDirect(s)
}

// addCloudsplainingToolsDirect implements direct Dagger calls for Cloudsplaining tools
func addCloudsplainingToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		profile := request.GetString("profile", "")
		includeNonDefaultVersions := request.GetBool("include_non_default_policy_versions", false)

		// Create Cloudsplaining module and download
		cloudsplainingModule := modules.NewCloudsplainingModule(client)
		result, err := cloudsplainingModule.Download(ctx, profile, includeNonDefaultVersions)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudsplaining download failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		inputFile := request.GetString("input_file", "")
		exclusionsFile := request.GetString("exclusions_file", "")
		output := request.GetString("output", "")

		// Create Cloudsplaining module and scan account data
		cloudsplainingModule := modules.NewCloudsplainingModule(client)
		result, err := cloudsplainingModule.ScanAccountData(ctx, inputFile, exclusionsFile, output)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudsplaining scan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		inputFile := request.GetString("input_file", "")

		// Create Cloudsplaining module and scan policy file
		cloudsplainingModule := modules.NewCloudsplainingModule(client)
		result, err := cloudsplainingModule.ScanPolicyFile(ctx, inputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudsplaining scan policy file failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Cloudsplaining create exclusions file
	createExclusionsTool := mcp.NewTool("cloudsplaining_create_exclusions_file",
		mcp.WithDescription("Create exclusions file template"),
	)
	s.AddTool(createExclusionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create Cloudsplaining module and create exclusions file
		cloudsplainingModule := modules.NewCloudsplainingModule(client)
		result, err := cloudsplainingModule.CreateExclusionsFile(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudsplaining create exclusions file failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		outputFile := request.GetString("output_file", "")

		// Create Cloudsplaining module and create multi-account config
		cloudsplainingModule := modules.NewCloudsplainingModule(client)
		result, err := cloudsplainingModule.CreateMultiAccountConfig(ctx, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudsplaining create multi-account config failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		configFile := request.GetString("config_file", "")
		profile := request.GetString("profile", "")
		roleName := request.GetString("role_name", "")
		outputBucket := request.GetString("output_bucket", "")
		outputDirectory := request.GetString("output_directory", "")

		// Create Cloudsplaining module and scan multi-account
		cloudsplainingModule := modules.NewCloudsplainingModule(client)
		result, err := cloudsplainingModule.ScanMultiAccount(ctx, configFile, profile, roleName, outputBucket, outputDirectory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cloudsplaining scan multi-account failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}