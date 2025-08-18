package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddAWSIAMRotationTools adds AWS IAM rotation MCP tool implementations
func AddAWSIAMRotationTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		username := request.GetString("user_name", "")
		args := []string{"aws", "iam", "list-access-keys", "--user-name", username}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		return executeShipCommand(args)
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
		username := request.GetString("user_name", "")
		args := []string{"aws", "iam", "create-access-key", "--user-name", username}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		return executeShipCommand(args)
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
		accessKeyId := request.GetString("access_key_id", "")
		status := request.GetString("status", "")
		username := request.GetString("user_name", "")
		args := []string{"aws", "iam", "update-access-key", "--access-key-id", accessKeyId, "--status", status, "--user-name", username}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		return executeShipCommand(args)
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
		accessKeyId := request.GetString("access_key_id", "")
		username := request.GetString("user_name", "")
		args := []string{"aws", "iam", "delete-access-key", "--access-key-id", accessKeyId, "--user-name", username}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		return executeShipCommand(args)
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
		accessKeyId := request.GetString("access_key_id", "")
		args := []string{"aws", "iam", "get-access-key-last-used", "--access-key-id", accessKeyId}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		return executeShipCommand(args)
	})

	// AWS CLI version tool
	getVersionTool := mcp.NewTool("aws_iam_get_version",
		mcp.WithDescription("Get AWS CLI version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"aws", "--version"}
		return executeShipCommand(args)
	})
}