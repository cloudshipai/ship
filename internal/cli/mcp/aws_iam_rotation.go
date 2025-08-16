package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddAWSIAMRotationTools adds AWS IAM rotation MCP tool implementations
func AddAWSIAMRotationTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// AWS IAM rotate credentials tool
	rotateCredentialsTool := mcp.NewTool("aws_iam_rotate_credentials",
		mcp.WithDescription("Rotate AWS IAM user credentials (access keys)"),
		mcp.WithString("username",
			mcp.Description("IAM username to rotate credentials for"),
			mcp.Required(),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
	)
	s.AddTool(rotateCredentialsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		username := request.GetString("username", "")
		args := []string{"aws", "iam-rotation", "rotate-credentials", username}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		return executeShipCommand(args)
	})

	// AWS IAM rotate service account tool
	rotateServiceAccountTool := mcp.NewTool("aws_iam_rotate_service_account",
		mcp.WithDescription("Rotate AWS service account credentials and update applications"),
		mcp.WithString("service_account",
			mcp.Description("Service account name"),
			mcp.Required(),
		),
		mcp.WithString("applications",
			mcp.Description("Comma-separated list of applications to update"),
		),
	)
	s.AddTool(rotateServiceAccountTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		serviceAccount := request.GetString("service_account", "")
		args := []string{"aws", "iam-rotation", "rotate-service-account", serviceAccount}
		if apps := request.GetString("applications", ""); apps != "" {
			args = append(args, "--applications", apps)
		}
		return executeShipCommand(args)
	})

	// AWS IAM schedule rotation tool
	scheduleRotationTool := mcp.NewTool("aws_iam_schedule_rotation",
		mcp.WithDescription("Schedule automated IAM credential rotation"),
		mcp.WithString("username",
			mcp.Description("IAM username for scheduled rotation"),
			mcp.Required(),
		),
		mcp.WithString("schedule",
			mcp.Description("Rotation schedule (e.g., 30d, 90d)"),
			mcp.Required(),
		),
	)
	s.AddTool(scheduleRotationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		username := request.GetString("username", "")
		schedule := request.GetString("schedule", "")
		args := []string{"aws", "iam-rotation", "schedule", username, "--schedule", schedule}
		return executeShipCommand(args)
	})

	// AWS IAM check rotation status tool
	checkRotationStatusTool := mcp.NewTool("aws_iam_check_rotation_status",
		mcp.WithDescription("Check IAM credential rotation status and history"),
		mcp.WithString("username",
			mcp.Description("IAM username to check (optional, checks all if not specified)"),
		),
	)
	s.AddTool(checkRotationStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"aws", "iam-rotation", "status"}
		if username := request.GetString("username", ""); username != "" {
			args = append(args, "--username", username)
		}
		return executeShipCommand(args)
	})

	// AWS IAM emergency rotation tool
	emergencyRotationTool := mcp.NewTool("aws_iam_emergency_rotation",
		mcp.WithDescription("Perform emergency rotation of compromised IAM credentials"),
		mcp.WithString("username",
			mcp.Description("IAM username for emergency rotation"),
			mcp.Required(),
		),
		mcp.WithString("notification_target",
			mcp.Description("SNS topic or email for notifications"),
		),
	)
	s.AddTool(emergencyRotationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		username := request.GetString("username", "")
		args := []string{"aws", "iam-rotation", "emergency", username}
		if target := request.GetString("notification_target", ""); target != "" {
			args = append(args, "--notify", target)
		}
		return executeShipCommand(args)
	})

	// AWS IAM rotation get version tool
	getVersionTool := mcp.NewTool("aws_iam_rotation_get_version",
		mcp.WithDescription("Get AWS IAM rotation tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"aws", "iam-rotation", "--version"}
		return executeShipCommand(args)
	})
}