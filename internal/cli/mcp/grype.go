package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGrypeTools adds Grype (vulnerability scanner for container images and filesystems) MCP tool implementations
func AddGrypeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		target := request.GetString("target", "")
		args := []string{"grype", target}
		
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "-o", output)
		}
		if failOn := request.GetString("fail_on", ""); failOn != "" {
			args = append(args, "--fail-on", failOn)
		}
		if request.GetBool("only_fixed", false) {
			args = append(args, "--only-fixed")
		}
		if request.GetBool("only_notfixed", false) {
			args = append(args, "--only-notfixed")
		}
		if exclude := request.GetString("exclude", ""); exclude != "" {
			args = append(args, "--exclude", exclude)
		}
		if scope := request.GetString("scope", ""); scope != "" {
			args = append(args, "--scope", scope)
		}
		if request.GetBool("add_cpes_if_none", false) {
			args = append(args, "--add-cpes-if-none")
		}
		if distro := request.GetString("distro", ""); distro != "" {
			args = append(args, "--distro", distro)
		}
		
		return executeShipCommand(args)
	})

	// Grype database status
	dbStatusTool := mcp.NewTool("grype_db_status",
		mcp.WithDescription("Report current status of Grype's vulnerability database"),
	)
	s.AddTool(dbStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"grype", "db", "status"}
		return executeShipCommand(args)
	})

	// Grype database check
	dbCheckTool := mcp.NewTool("grype_db_check",
		mcp.WithDescription("Check if updates are available for the vulnerability database"),
	)
	s.AddTool(dbCheckTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"grype", "db", "check"}
		return executeShipCommand(args)
	})

	// Grype database update
	dbUpdateTool := mcp.NewTool("grype_db_update",
		mcp.WithDescription("Update the vulnerability database to the latest version"),
	)
	s.AddTool(dbUpdateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"grype", "db", "update"}
		return executeShipCommand(args)
	})

	// Grype database list
	dbListTool := mcp.NewTool("grype_db_list",
		mcp.WithDescription("Show databases available for download"),
	)
	s.AddTool(dbListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"grype", "db", "list"}
		return executeShipCommand(args)
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
		archivePath := request.GetString("archive_path", "")
		args := []string{"grype", "db", "import", archivePath}
		return executeShipCommand(args)
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
		cveId := request.GetString("cve_id", "")
		// Note: This would typically pipe JSON results to grype explain
		// For Ship CLI integration, we'll just show the command structure
		// Grype explain requires piping scan results and is not a standalone command
		// Providing guidance on how to use explain feature
		args := []string{"echo", "To explain a CVE, pipe scan results to grype explain: echo '<scan-json>' | grype explain --id " + cveId}
		return executeShipCommand(args)
	})

	// Grype version
	versionTool := mcp.NewTool("grype_version",
		mcp.WithDescription("Display Grype version information"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"grype", "--version"}
		return executeShipCommand(args)
	})
}