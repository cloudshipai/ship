package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddVeleroTools adds Velero (Kubernetes backup and restore) MCP tool implementations
func AddVeleroTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Velero create backup tool
	createBackupTool := mcp.NewTool("velero_create_backup",
		mcp.WithDescription("Create Kubernetes cluster backup using Velero"),
		mcp.WithString("backup_name",
			mcp.Description("Name for the backup"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
			mcp.Required(),
		),
		mcp.WithString("include_namespaces",
			mcp.Description("Comma-separated list of namespaces to include"),
		),
		mcp.WithString("exclude_namespaces",
			mcp.Description("Comma-separated list of namespaces to exclude"),
		),
	)
	s.AddTool(createBackupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		backupName := request.GetString("backup_name", "")
		kubeconfig := request.GetString("kubeconfig", "")
		args := []string{"security", "velero", "--create-backup", backupName, "--kubeconfig", kubeconfig}
		if includeNs := request.GetString("include_namespaces", ""); includeNs != "" {
			args = append(args, "--include-namespaces", includeNs)
		}
		if excludeNs := request.GetString("exclude_namespaces", ""); excludeNs != "" {
			args = append(args, "--exclude-namespaces", excludeNs)
		}
		return executeShipCommand(args)
	})

	// Velero list backups tool
	listBackupsTool := mcp.NewTool("velero_list_backups",
		mcp.WithDescription("List Velero backups"),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
			mcp.Required(),
		),
	)
	s.AddTool(listBackupsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		kubeconfig := request.GetString("kubeconfig", "")
		args := []string{"security", "velero", "--list-backups", "--kubeconfig", kubeconfig}
		return executeShipCommand(args)
	})

	// Velero restore backup tool
	restoreBackupTool := mcp.NewTool("velero_restore_backup",
		mcp.WithDescription("Restore from Velero backup"),
		mcp.WithString("backup_name",
			mcp.Description("Name of the backup to restore from"),
			mcp.Required(),
		),
		mcp.WithString("restore_name",
			mcp.Description("Name for the restore operation"),
			mcp.Required(),
		),
		mcp.WithString("kubeconfig",
			mcp.Description("Path to kubeconfig file"),
			mcp.Required(),
		),
		mcp.WithString("include_namespaces",
			mcp.Description("Comma-separated list of namespaces to include in restore"),
		),
	)
	s.AddTool(restoreBackupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		backupName := request.GetString("backup_name", "")
		restoreName := request.GetString("restore_name", "")
		kubeconfig := request.GetString("kubeconfig", "")
		args := []string{"security", "velero", "--restore", backupName, "--restore-name", restoreName, "--kubeconfig", kubeconfig}
		if includeNs := request.GetString("include_namespaces", ""); includeNs != "" {
			args = append(args, "--include-namespaces", includeNs)
		}
		return executeShipCommand(args)
	})

	// Velero get version tool
	getVersionTool := mcp.NewTool("velero_get_version",
		mcp.WithDescription("Get Velero version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "velero", "--version"}
		return executeShipCommand(args)
	})
}