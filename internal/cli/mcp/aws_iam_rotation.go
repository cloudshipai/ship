package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddAWSIAMRotationTools adds AWS IAM rotation MCP tool implementations
func AddAWSIAMRotationTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addAWSIAMRotationToolsDirect(s)
}

// addAWSIAMRotationToolsDirect implements direct Dagger calls for AWS IAM rotation tools
func addAWSIAMRotationToolsDirect(s *server.MCPServer) {
	// AWS IAM list access keys tool
	listAccessKeysTool := mcp.NewTool("aws_iam_list_access_keys",
		mcp.WithDescription("List access keys for an IAM user"),
		mcp.WithString("user_name",
			mcp.Description("IAM username to list access keys for"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
	)
	s.AddTool(listAccessKeysTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		username := request.GetString("user_name", "")
		profile := request.GetString("profile", "")

		// Create AWS IAM rotation module and list access keys
		awsModule := modules.NewAWSIAMRotationModule(client)
		result, err := awsModule.ListAccessKeys(ctx, username, profile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("list access keys failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// AWS IAM create access key tool
	createAccessKeyTool := mcp.NewTool("aws_iam_create_access_key",
		mcp.WithDescription("Create a new access key for an IAM user"),
		mcp.WithString("user_name",
			mcp.Description("IAM username to create access key for"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
	)
	s.AddTool(createAccessKeyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		username := request.GetString("user_name", "")
		profile := request.GetString("profile", "")

		// Create AWS IAM rotation module and create access key (using RotateAccessKeys function)
		awsModule := modules.NewAWSIAMRotationModule(client)
		result, err := awsModule.RotateAccessKeys(ctx, username, profile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("create access key failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// AWS IAM update access key tool
	updateAccessKeyTool := mcp.NewTool("aws_iam_update_access_key",
		mcp.WithDescription("Update the status of an access key (Active/Inactive)"),
		mcp.WithString("access_key_id",
			mcp.Description("Access key ID to update"),
			mcp.Required(),
		),
		mcp.WithString("status",
			mcp.Description("New status for the access key"),
			mcp.Enum("Active", "Inactive"),
			mcp.Required(),
		),
		mcp.WithString("user_name",
			mcp.Description("IAM username that owns the access key"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
	)
	s.AddTool(updateAccessKeyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		accessKeyId := request.GetString("access_key_id", "")
		status := request.GetString("status", "")
		username := request.GetString("user_name", "")
		profile := request.GetString("profile", "")

		// Create AWS IAM rotation module and update access key
		awsModule := modules.NewAWSIAMRotationModule(client)
		result, err := awsModule.UpdateAccessKey(ctx, username, accessKeyId, status, profile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("update access key failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// AWS IAM delete access key tool
	deleteAccessKeyTool := mcp.NewTool("aws_iam_delete_access_key",
		mcp.WithDescription("Delete an access key for an IAM user"),
		mcp.WithString("access_key_id",
			mcp.Description("Access key ID to delete"),
			mcp.Required(),
		),
		mcp.WithString("user_name",
			mcp.Description("IAM username that owns the access key"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
	)
	s.AddTool(deleteAccessKeyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		accessKeyId := request.GetString("access_key_id", "")
		username := request.GetString("user_name", "")
		profile := request.GetString("profile", "")

		// Create AWS IAM rotation module and delete access key
		awsModule := modules.NewAWSIAMRotationModule(client)
		result, err := awsModule.DeleteAccessKey(ctx, username, accessKeyId, profile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("delete access key failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// AWS IAM get access key last used tool
	getAccessKeyLastUsedTool := mcp.NewTool("aws_iam_get_access_key_last_used",
		mcp.WithDescription("Get information about when an access key was last used"),
		mcp.WithString("access_key_id",
			mcp.Description("Access key ID to check"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
	)
	s.AddTool(getAccessKeyLastUsedTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		accessKeyId := request.GetString("access_key_id", "")
		profile := request.GetString("profile", "")

		// Create AWS IAM rotation module and get access key last used info
		awsModule := modules.NewAWSIAMRotationModule(client)
		result, err := awsModule.GetAccessKeyLastUsed(ctx, accessKeyId, profile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get access key last used failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// AWS CLI version tool
	getVersionTool := mcp.NewTool("aws_iam_get_version",
		mcp.WithDescription("Get AWS CLI version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create AWS IAM rotation module and get version
		awsModule := modules.NewAWSIAMRotationModule(client)
		result, err := awsModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("get AWS CLI version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}