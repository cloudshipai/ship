package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddVeleroTools adds Velero (Kubernetes backup and restore) MCP tool implementations using real CLI commands
func AddVeleroTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Velero install with provider configuration
	installTool := mcp.NewTool("velero_install",
		mcp.WithDescription("Install Velero with storage provider configuration using real velero CLI"),
		mcp.WithString("provider",
			mcp.Description("Storage provider (aws, gcp, azure, vsphere)"),
			mcp.Enum("aws", "gcp", "azure", "vsphere"),
		),
		mcp.WithString("bucket",
			mcp.Description("Backup storage bucket name"),
			mcp.Required(),
		),
		mcp.WithString("region",
			mcp.Description("Cloud region"),
		),
		mcp.WithString("secret_file",
			mcp.Description("Path to credentials file"),
		),
		mcp.WithBoolean("use_node_agent",
			mcp.Description("Enable filesystem backup via node agent"),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for Velero"),
		),
		mcp.WithBoolean("no_secret",
			mcp.Description("Use cloud identity instead of static credentials"),
		),
		mcp.WithString("plugins",
			mcp.Description("Velero plugins to install (comma-separated)"),
		),
	)
	s.AddTool(installTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"velero", "install"}
		
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, "--provider", provider)
		}
		if bucket := request.GetString("bucket", ""); bucket != "" {
			args = append(args, "--bucket", bucket)
		}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--backup-location-config", "region="+region)
		}
		if secretFile := request.GetString("secret_file", ""); secretFile != "" {
			args = append(args, "--secret-file", secretFile)
		}
		if request.GetBool("use_node_agent", false) {
			args = append(args, "--use-node-agent")
		}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if request.GetBool("no_secret", false) {
			args = append(args, "--no-secret")
		}
		if plugins := request.GetString("plugins", ""); plugins != "" {
			args = append(args, "--plugins", plugins)
		}
		args = append(args, "--wait")
		
		return executeShipCommand(args)
	})

	// Velero create backup schedule
	createScheduleTool := mcp.NewTool("velero_create_schedule",
		mcp.WithDescription("Create a backup schedule using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Schedule name"),
			mcp.Required(),
		),
		mcp.WithString("schedule",
			mcp.Description("Cron expression (e.g., '0 2 * * *')"),
			mcp.Required(),
		),
		mcp.WithString("include_namespaces",
			mcp.Description("Comma-separated namespaces to include"),
		),
		mcp.WithString("exclude_namespaces",
			mcp.Description("Comma-separated namespaces to exclude"),
		),
		mcp.WithString("ttl",
			mcp.Description("Time to live for backups (e.g., 720h)"),
		),
		mcp.WithString("labels",
			mcp.Description("Comma-separated labels (key=value)"),
		),
	)
	s.AddTool(createScheduleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		schedule := request.GetString("schedule", "")
		args := []string{"velero", "schedule", "create", name, "--schedule", schedule}
		
		if ttl := request.GetString("ttl", ""); ttl != "" {
			args = append(args, "--ttl", ttl)
		}
		if includeNs := request.GetString("include_namespaces", ""); includeNs != "" {
			args = append(args, "--include-namespaces", includeNs)
		}
		if excludeNs := request.GetString("exclude_namespaces", ""); excludeNs != "" {
			args = append(args, "--exclude-namespaces", excludeNs)
		}
		if labels := request.GetString("labels", ""); labels != "" {
			args = append(args, "--labels", labels)
		}
		
		return executeShipCommand(args)
	})

	// Velero create backup on demand
	createBackupTool := mcp.NewTool("velero_backup_create",
		mcp.WithDescription("Create an on-demand backup using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Backup name"),
			mcp.Required(),
		),
		mcp.WithString("from_schedule",
			mcp.Description("Create backup from existing schedule"),
		),
		mcp.WithString("include_namespaces",
			mcp.Description("Comma-separated namespaces to include"),
		),
		mcp.WithString("exclude_namespaces",
			mcp.Description("Comma-separated namespaces to exclude"),
		),
		mcp.WithString("labels",
			mcp.Description("Comma-separated labels (key=value)"),
		),
		mcp.WithBoolean("wait",
			mcp.Description("Wait for backup to complete"),
		),
	)
	s.AddTool(createBackupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		args := []string{"velero", "backup", "create", name}
		
		if fromSchedule := request.GetString("from_schedule", ""); fromSchedule != "" {
			args = append(args, "--from-schedule", fromSchedule)
		}
		if includeNs := request.GetString("include_namespaces", ""); includeNs != "" {
			args = append(args, "--include-namespaces", includeNs)
		}
		if excludeNs := request.GetString("exclude_namespaces", ""); excludeNs != "" {
			args = append(args, "--exclude-namespaces", excludeNs)
		}
		if labels := request.GetString("labels", ""); labels != "" {
			args = append(args, "--labels", labels)
		}
		if request.GetBool("wait", false) {
			args = append(args, "--wait")
		}
		
		return executeShipCommand(args)
	})

	// Velero restore from backup
	createRestoreTool := mcp.NewTool("velero_restore_create",
		mcp.WithDescription("Create a restore from backup using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Restore name"),
			mcp.Required(),
		),
		mcp.WithString("from_backup",
			mcp.Description("Backup to restore from"),
			mcp.Required(),
		),
		mcp.WithString("namespace_mappings",
			mcp.Description("Map source namespaces to target (src1:tgt1,src2:tgt2)"),
		),
		mcp.WithString("include_namespaces",
			mcp.Description("Comma-separated namespaces to include"),
		),
		mcp.WithString("exclude_namespaces",
			mcp.Description("Comma-separated namespaces to exclude"),
		),
		mcp.WithBoolean("restore_pvs",
			mcp.Description("Restore persistent volumes"),
		),
		mcp.WithBoolean("wait",
			mcp.Description("Wait for restore to complete"),
		),
	)
	s.AddTool(createRestoreTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		fromBackup := request.GetString("from_backup", "")
		args := []string{"velero", "restore", "create", name, "--from-backup", fromBackup}
		
		if namespaceMappings := request.GetString("namespace_mappings", ""); namespaceMappings != "" {
			args = append(args, "--namespace-mappings", namespaceMappings)
		}
		if includeNs := request.GetString("include_namespaces", ""); includeNs != "" {
			args = append(args, "--include-namespaces", includeNs)
		}
		if excludeNs := request.GetString("exclude_namespaces", ""); excludeNs != "" {
			args = append(args, "--exclude-namespaces", excludeNs)
		}
		if request.GetBool("restore_pvs", false) {
			args = append(args, "--restore-volumes=true")
		}
		if request.GetBool("wait", false) {
			args = append(args, "--wait")
		}
		
		return executeShipCommand(args)
	})

	// Velero get backups
	getBackupsTool := mcp.NewTool("velero_backup_get",
		mcp.WithDescription("Get list of backups using velero CLI"),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(getBackupsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"velero", "backup", "get"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "-o", output)
		}
		
		return executeShipCommand(args)
	})

	// Velero describe backup
	describeBackupTool := mcp.NewTool("velero_backup_describe",
		mcp.WithDescription("Describe a specific backup using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Backup name"),
			mcp.Required(),
		),
	)
	s.AddTool(describeBackupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		args := []string{"velero", "backup", "describe", name}
		return executeShipCommand(args)
	})

	// Velero get restores  
	getRestoresTool := mcp.NewTool("velero_restore_get",
		mcp.WithDescription("Get list of restores using velero CLI"),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(getRestoresTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"velero", "restore", "get"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "-o", output)
		}
		
		return executeShipCommand(args)
	})

	// Velero describe restore
	describeRestoreTool := mcp.NewTool("velero_restore_describe",
		mcp.WithDescription("Describe a specific restore using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Restore name"),
			mcp.Required(),
		),
	)
	s.AddTool(describeRestoreTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		args := []string{"velero", "restore", "describe", name}
		return executeShipCommand(args)
	})

	// Velero get schedules
	getSchedulesTool := mcp.NewTool("velero_schedule_get",
		mcp.WithDescription("Get list of backup schedules using velero CLI"),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(getSchedulesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"velero", "schedule", "get"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "-o", output)
		}
		
		return executeShipCommand(args)
	})

	// Velero delete backup
	deleteBackupTool := mcp.NewTool("velero_backup_delete",
		mcp.WithDescription("Delete a backup using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Backup name to delete"),
			mcp.Required(),
		),
		mcp.WithBoolean("confirm",
			mcp.Description("Confirm deletion without prompting"),
		),
	)
	s.AddTool(deleteBackupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		args := []string{"velero", "backup", "delete", name}
		
		if request.GetBool("confirm", false) {
			args = append(args, "--confirm")
		}
		
		return executeShipCommand(args)
	})

	// Velero delete schedule
	deleteScheduleTool := mcp.NewTool("velero_schedule_delete",
		mcp.WithDescription("Delete a backup schedule using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Schedule name to delete"),
			mcp.Required(),
		),
		mcp.WithBoolean("confirm",
			mcp.Description("Confirm deletion without prompting"),
		),
	)
	s.AddTool(deleteScheduleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		args := []string{"velero", "schedule", "delete", name}
		
		if request.GetBool("confirm", false) {
			args = append(args, "--confirm")
		}
		
		return executeShipCommand(args)
	})

	// Velero backup location create
	backupLocationCreateTool := mcp.NewTool("velero_backup_location_create",
		mcp.WithDescription("Create backup storage location using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Backup location name"),
			mcp.Required(),
		),
		mcp.WithString("provider",
			mcp.Description("Storage provider (aws, gcp, azure)"),
			mcp.Required(),
		),
		mcp.WithString("bucket",
			mcp.Description("Storage bucket name"),
			mcp.Required(),
		),
		mcp.WithString("config",
			mcp.Description("Provider-specific configuration (key=value,key2=value2)"),
		),
	)
	s.AddTool(backupLocationCreateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		provider := request.GetString("provider", "")
		bucket := request.GetString("bucket", "")
		
		args := []string{"velero", "backup-location", "create", name, "--provider", provider, "--bucket", bucket}
		
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		
		return executeShipCommand(args)
	})

	// Velero backup location get
	backupLocationGetTool := mcp.NewTool("velero_backup_location_get",
		mcp.WithDescription("Get backup storage locations using velero CLI"),
		mcp.WithString("output",
			mcp.Description("Output format"),
			mcp.Enum("json", "yaml", "table"),
		),
	)
	s.AddTool(backupLocationGetTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"velero", "backup-location", "get"}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "-o", output)
		}
		
		return executeShipCommand(args)
	})

	// Velero get version tool
	getVersionTool := mcp.NewTool("velero_version",
		mcp.WithDescription("Get Velero version information using velero CLI"),
		mcp.WithBoolean("client_only",
			mcp.Description("Show client version only"),
		),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"velero", "version"}
		if request.GetBool("client_only", false) {
			args = append(args, "--client-only")
		}
		return executeShipCommand(args)
	})

	// Velero backup logs
	backupLogsTool := mcp.NewTool("velero_backup_logs",
		mcp.WithDescription("Get backup logs using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Backup name"),
			mcp.Required(),
		),
	)
	s.AddTool(backupLogsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		args := []string{"velero", "backup", "logs", name}
		return executeShipCommand(args)
	})

	// Velero restore logs
	restoreLogsTool := mcp.NewTool("velero_restore_logs",
		mcp.WithDescription("Get restore logs using velero CLI"),
		mcp.WithString("name",
			mcp.Description("Restore name"),
			mcp.Required(),
		),
	)
	s.AddTool(restoreLogsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetString("name", "")
		args := []string{"velero", "restore", "logs", name}
		return executeShipCommand(args)
	})
}