package mcp

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddOpenVASTools adds OpenVAS/GVM vulnerability scanning tools to the MCP server
// Note: OpenVAS requires a running Greenbone Vulnerability Management (GVM) stack
func AddOpenVASTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	addOpenVASToolsDirect(s)
}

// addOpenVASToolsDirect adds OpenVAS tools using direct Dagger module calls
func addOpenVASToolsDirect(s *server.MCPServer) {
	// Get GVM system status
	getSystemStatusTool := mcp.NewTool("openvas_get_system_status",
		mcp.WithDescription("Get the status of a running GVM (Greenbone Vulnerability Management) instance. Requires a running GVM stack."),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
	)
	s.AddTool(getSystemStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.GetSystemStatus(ctx, socketPath, username, password)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get system status: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// List scans
	listScansTool := mcp.NewTool("openvas_list_scans",
		mcp.WithDescription("List all vulnerability scans from a GVM instance"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
	)
	s.AddTool(listScansTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.ListScans(ctx, socketPath, username, password)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list scans: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// List targets
	listTargetsTool := mcp.NewTool("openvas_list_targets",
		mcp.WithDescription("List all scan targets from a GVM instance"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
	)
	s.AddTool(listTargetsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.ListTargets(ctx, socketPath, username, password)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list targets: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Create target
	createTargetTool := mcp.NewTool("openvas_create_target",
		mcp.WithDescription("Create a new scan target in GVM"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
		mcp.WithString("target_name",
			mcp.Description("Name for the target"),
			mcp.Required(),
		),
		mcp.WithString("hosts",
			mcp.Description("Comma-separated list of hosts/IPs to scan"),
			mcp.Required(),
		),
		mcp.WithString("port_list_id",
			mcp.Description("ID of port list to use (optional)"),
		),
	)
	s.AddTool(createTargetTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")
		targetName := request.GetString("target_name", "")
		hosts := request.GetString("hosts", "")
		portListID := request.GetString("port_list_id", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}
		if targetName == "" || hosts == "" {
			return mcp.NewToolResultError("target_name and hosts are required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.CreateTarget(ctx, socketPath, username, password, targetName, hosts, portListID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create target: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Create scan
	createScanTool := mcp.NewTool("openvas_create_scan",
		mcp.WithDescription("Create a new vulnerability scan task in GVM"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
		mcp.WithString("scan_name",
			mcp.Description("Name for the scan task"),
			mcp.Required(),
		),
		mcp.WithString("target_id",
			mcp.Description("ID of the target to scan"),
			mcp.Required(),
		),
		mcp.WithString("config_id",
			mcp.Description("Scan configuration ID (default: Full and fast)"),
		),
	)
	s.AddTool(createScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")
		scanName := request.GetString("scan_name", "")
		targetID := request.GetString("target_id", "")
		configID := request.GetString("config_id", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}
		if scanName == "" || targetID == "" {
			return mcp.NewToolResultError("scan_name and target_id are required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.CreateScan(ctx, socketPath, username, password, scanName, targetID, configID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create scan: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Start scan
	startScanTool := mcp.NewTool("openvas_start_scan",
		mcp.WithDescription("Start an existing scan task in GVM"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
		mcp.WithString("task_id",
			mcp.Description("ID of the scan task to start"),
			mcp.Required(),
		),
	)
	s.AddTool(startScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")
		taskID := request.GetString("task_id", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}
		if taskID == "" {
			return mcp.NewToolResultError("task_id is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.StartScan(ctx, socketPath, username, password, taskID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to start scan: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Get scan status
	getScanStatusTool := mcp.NewTool("openvas_get_scan_status",
		mcp.WithDescription("Get the status of a scan task in GVM"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
		mcp.WithString("task_id",
			mcp.Description("ID of the scan task"),
			mcp.Required(),
		),
	)
	s.AddTool(getScanStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")
		taskID := request.GetString("task_id", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}
		if taskID == "" {
			return mcp.NewToolResultError("task_id is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.GetScanStatus(ctx, socketPath, username, password, taskID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get scan status: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Get scan results
	getScanResultsTool := mcp.NewTool("openvas_get_scan_results",
		mcp.WithDescription("Get the results of a completed scan from GVM"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
		mcp.WithString("report_id",
			mcp.Description("ID of the scan report"),
			mcp.Required(),
		),
		mcp.WithString("min_severity",
			mcp.Description("Minimum severity to include (0.0-10.0)"),
		),
	)
	s.AddTool(getScanResultsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")
		reportID := request.GetString("report_id", "")
		minSeverity := request.GetString("min_severity", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}
		if reportID == "" {
			return mcp.NewToolResultError("report_id is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.GetScanResults(ctx, socketPath, username, password, reportID, minSeverity)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get scan results: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// List scan configs
	listScanConfigsTool := mcp.NewTool("openvas_list_scan_configs",
		mcp.WithDescription("List available scan configurations in GVM"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
	)
	s.AddTool(listScanConfigsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.ListScanConfigs(ctx, socketPath, username, password)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list scan configs: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Stop scan
	stopScanTool := mcp.NewTool("openvas_stop_scan",
		mcp.WithDescription("Stop a running scan task in GVM"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
		mcp.WithString("task_id",
			mcp.Description("ID of the scan task to stop"),
			mcp.Required(),
		),
	)
	s.AddTool(stopScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")
		taskID := request.GetString("task_id", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}
		if taskID == "" {
			return mcp.NewToolResultError("task_id is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.StopScan(ctx, socketPath, username, password, taskID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to stop scan: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Delete scan
	deleteScanTool := mcp.NewTool("openvas_delete_scan",
		mcp.WithDescription("Delete a scan task from GVM"),
		mcp.WithString("socket_path",
			mcp.Description("Path to GVM socket (default: /run/gvmd/gvmd.sock)"),
		),
		mcp.WithString("username",
			mcp.Description("GVM username"),
			mcp.Required(),
		),
		mcp.WithString("password",
			mcp.Description("GVM password"),
			mcp.Required(),
		),
		mcp.WithString("task_id",
			mcp.Description("ID of the scan task to delete"),
			mcp.Required(),
		),
	)
	s.AddTool(deleteScanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		socketPath := request.GetString("socket_path", "/run/gvmd/gvmd.sock")
		username := request.GetString("username", "")
		password := request.GetString("password", "")
		taskID := request.GetString("task_id", "")

		if username == "" || password == "" {
			return mcp.NewToolResultError("username and password are required"), nil
		}
		if taskID == "" {
			return mcp.NewToolResultError("task_id is required"), nil
		}

		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to connect to dagger: %v", err)), nil
		}
		defer client.Close()

		module := modules.NewOpenVASModule(client)
		result, err := module.DeleteScan(ctx, socketPath, username, password, taskID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to delete scan: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}
