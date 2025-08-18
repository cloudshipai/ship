package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddVeleroTools adds Velero (Kubernetes backup and restore) MCP tool implementations using direct Dagger calls
func AddVeleroTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addVeleroToolsDirect(s)
}

// addVeleroToolsDirect adds Velero tools using direct Dagger module calls
func addVeleroToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		provider := request.GetString("provider", "")
		bucket := request.GetString("bucket", "")
		region := request.GetString("region", "")
		secretFile := request.GetString("secret_file", "")
		useNodeAgent := request.GetBool("use_node_agent", false)
		namespace := request.GetString("namespace", "")
		noSecret := request.GetBool("no_secret", false)
		plugins := request.GetString("plugins", "")

		if bucket == "" {
			return mcp.NewToolResultError("bucket is required"), nil
		}

		// Install Velero
		output, err := module.Install(ctx, provider, bucket, region, secretFile, useNodeAgent, namespace, noSecret, plugins)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero install failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		schedule := request.GetString("schedule", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		if schedule == "" {
			return mcp.NewToolResultError("schedule is required"), nil
		}

		includeNamespaces := request.GetString("include_namespaces", "")
		excludeNamespaces := request.GetString("exclude_namespaces", "")
		ttl := request.GetString("ttl", "")
		labels := request.GetString("labels", "")
		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Create schedule
		output, err := module.CreateSchedule(ctx, name, schedule, includeNamespaces, excludeNamespaces, ttl, labels, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero create schedule failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		fromSchedule := request.GetString("from_schedule", "")
		includeNamespaces := request.GetString("include_namespaces", "")
		excludeNamespaces := request.GetString("exclude_namespaces", "")
		labels := request.GetString("labels", "")
		wait := request.GetBool("wait", false)
		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Create backup
		output, err := module.CreateBackupAdvanced(ctx, name, fromSchedule, includeNamespaces, excludeNamespaces, labels, wait, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero backup create failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		fromBackup := request.GetString("from_backup", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		if fromBackup == "" {
			return mcp.NewToolResultError("from_backup is required"), nil
		}

		namespaceMappings := request.GetString("namespace_mappings", "")
		includeNamespaces := request.GetString("include_namespaces", "")
		excludeNamespaces := request.GetString("exclude_namespaces", "")
		restorePVs := request.GetBool("restore_pvs", false)
		wait := request.GetBool("wait", false)
		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Create restore
		output, err := module.CreateRestoreAdvanced(ctx, name, fromBackup, namespaceMappings, includeNamespaces, excludeNamespaces, restorePVs, wait, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero restore create failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		output := request.GetString("output", "")
		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Get backups
		result, err := module.GetBackups(ctx, output, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero get backups failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Describe backup
		result, err := module.DescribeBackup(ctx, name, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero describe backup failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		output := request.GetString("output", "")
		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Get restores
		result, err := module.GetRestores(ctx, output, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero get restores failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Describe restore
		result, err := module.DescribeRestore(ctx, name, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero describe restore failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		output := request.GetString("output", "")
		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Get schedules
		result, err := module.GetSchedules(ctx, output, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero get schedules failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		confirm := request.GetBool("confirm", false)
		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Delete backup
		result, err := module.DeleteBackup(ctx, name, confirm, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero delete backup failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		confirm := request.GetBool("confirm", false)
		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Delete schedule
		result, err := module.DeleteSchedule(ctx, name, confirm, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero delete schedule failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		provider := request.GetString("provider", "")
		bucket := request.GetString("bucket", "")
		config := request.GetString("config", "")

		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		if provider == "" {
			return mcp.NewToolResultError("provider is required"), nil
		}
		if bucket == "" {
			return mcp.NewToolResultError("bucket is required"), nil
		}

		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Create backup location
		result, err := module.CreateBackupLocation(ctx, name, provider, bucket, config, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero create backup location failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		output := request.GetString("output", "")
		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Get backup locations
		result, err := module.GetBackupLocations(ctx, output, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero get backup locations failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Velero get version tool
	getVersionTool := mcp.NewTool("velero_version",
		mcp.WithDescription("Get Velero version information using velero CLI"),
		mcp.WithBoolean("client_only",
			mcp.Description("Show client version only"),
		),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Note: client_only parameter not directly supported in Dagger GetVersion function (it's hardcoded as client-only)
		if !request.GetBool("client_only", false) {
			return mcp.NewToolResultError("Warning: Dagger GetVersion always returns client-only version"), nil
		}

		// Get version
		result, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Get backup logs
		result, err := module.GetBackupLogs(ctx, name, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero backup logs failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewVeleroModule(client)

		// Get parameters
		name := request.GetString("name", "")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		kubeconfig := "" // Note: kubeconfig parameter not available in MCP interface

		// Get restore logs
		result, err := module.GetRestoreLogs(ctx, name, kubeconfig)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Velero restore logs failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}