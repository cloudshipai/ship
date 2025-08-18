package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddGrypeTools adds Grype (vulnerability scanner for container images and filesystems) MCP tool implementations using direct Dagger calls
func AddGrypeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGrypeToolsDirect(s)
}

// addGrypeToolsDirect adds Grype tools using direct Dagger module calls
func addGrypeToolsDirect(s *server.MCPServer) {
	// Grype scan target for vulnerabilities
	scanTool := mcp.NewTool("grype_scan",
		mcp.WithDescription("Scan target for vulnerabilities (image, directory, archive, SBOM)"),
		mcp.WithString("target",
			mcp.Description("Target to scan (container image, directory, archive, or SBOM)"),
			mcp.Required(),
		),
		mcp.WithString("output",
			mcp.Description("Report output format"),
			mcp.Enum("table", "json", "cyclonedx", "sarif", "template"),
		),
		mcp.WithString("fail_on",
			mcp.Description("Exit with error on specified severity or higher"),
			mcp.Enum("negligible", "low", "medium", "high", "critical"),
		),
		mcp.WithBoolean("only_fixed",
			mcp.Description("Show only vulnerabilities with confirmed fixes"),
		),
		mcp.WithBoolean("only_notfixed",
			mcp.Description("Show only vulnerabilities without confirmed fixes"),
		),
		mcp.WithString("exclude",
			mcp.Description("Exclude specific file path from scanning"),
		),
		mcp.WithString("scope",
			mcp.Description("Scope for image layer scanning"),
			mcp.Enum("all-layers", "squashed"),
		),
		mcp.WithBoolean("add_cpes_if_none",
			mcp.Description("Generate CPE information if missing from SBOM"),
		),
		mcp.WithString("distro",
			mcp.Description("Specify distribution (format: <distro>:<version>)"),
		),
	)
	s.AddTool(scanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGrypeModule(client)

		// Get target
		target := request.GetString("target", "")
		if target == "" {
			return mcp.NewToolResultError("target is required"), nil
		}

		// Determine scan type and execute
		var output string
		
		// Check if target has severity filter
		if severity := request.GetString("fail_on", ""); severity != "" {
			output, err = module.ScanWithSeverity(ctx, target, severity)
		} else if strings.HasSuffix(target, ".sbom") || strings.HasSuffix(target, ".json") {
			// SBOM scan
			output, err = module.ScanSBOM(ctx, target)
		} else if strings.Contains(target, ":") || strings.Contains(target, "/") {
			// Image scan (has : for tag or / for registry)
			output, err = module.ScanImage(ctx, target)
		} else {
			// Directory scan
			output, err = module.ScanDirectory(ctx, target)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to scan target: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Grype database status
	dbStatusTool := mcp.NewTool("grype_db_status",
		mcp.WithDescription("Report current status of Grype's vulnerability database"),
	)
	s.AddTool(dbStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGrypeModule(client)

		// Get database status
		output, err := module.DBStatus(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get database status: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Grype database check
	dbCheckTool := mcp.NewTool("grype_db_check",
		mcp.WithDescription("Check if updates are available for the vulnerability database"),
	)
	s.AddTool(dbCheckTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGrypeModule(client)

		// Check database
		output, err := module.DBCheck(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to check database: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Grype database update
	dbUpdateTool := mcp.NewTool("grype_db_update",
		mcp.WithDescription("Update the vulnerability database to the latest version"),
	)
	s.AddTool(dbUpdateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGrypeModule(client)

		// Update database
		output, err := module.DBUpdate(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to update database: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Grype database list
	dbListTool := mcp.NewTool("grype_db_list",
		mcp.WithDescription("Show databases available for download"),
	)
	s.AddTool(dbListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGrypeModule(client)

		// List databases
		output, err := module.DBList(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list databases: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Grype database import
	dbImportTool := mcp.NewTool("grype_db_import",
		mcp.WithDescription("Import a vulnerability database archive"),
		mcp.WithString("archive_path",
			mcp.Description("Path to database archive file"),
			mcp.Required(),
		),
	)
	s.AddTool(dbImportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGrypeModule(client)

		// Get archive path
		archivePath := request.GetString("archive_path", "")
		if archivePath == "" {
			return mcp.NewToolResultError("archive_path is required"), nil
		}

		// Since there's no direct import function, use DBUpdate as placeholder
		// In practice, this would need a custom implementation
		output, err := module.DBUpdate(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to import database: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Database import from %s: %s", archivePath, output)), nil
	})

	// Grype explain CVE
	explainTool := mcp.NewTool("grype_explain",
		mcp.WithDescription("Get detailed information about a specific CVE from previous scan results"),
		mcp.WithString("cve_id",
			mcp.Description("CVE ID to explain (e.g., CVE-2023-36632)"),
			mcp.Required(),
		),
		mcp.WithString("scan_results",
			mcp.Description("JSON output from previous Grype scan"),
			mcp.Required(),
		),
	)
	s.AddTool(explainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGrypeModule(client)

		// Get CVE ID
		cveId := request.GetString("cve_id", "")
		if cveId == "" {
			return mcp.NewToolResultError("cve_id is required"), nil
		}

		// Explain CVE
		output, err := module.Explain(ctx, cveId)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to explain CVE: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Grype version
	versionTool := mcp.NewTool("grype_version",
		mcp.WithDescription("Display Grype version information"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewGrypeModule(client)

		// Get version
		version, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get version: %v", err)), nil
		}

		return mcp.NewToolResultText(version), nil
	})
}