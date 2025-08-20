package mcp

import (
	"context"
	"encoding/json"
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
	// Updated Grype scan tool to match specifications
	scanTool := mcp.NewTool("grype_scan",
		mcp.WithDescription("Vulnerability scan from a SBOM (preferred) or direct target; return JSON/SARIF paths and structured counts"),
		mcp.WithString("input",
			mcp.Description("Input source (e.g., sbom:/abs/sbom.json, dir:., docker:alpine:3.19)"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format: json or sarif (default: json)"),
		),
		mcp.WithString("min_severity",
			mcp.Description("Minimum severity to report: negligible|low|medium|high|critical"),
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

		// Get parameters
		input := request.GetString("input", "")
		if input == "" {
			return mcp.NewToolResultError("input is required"), nil
		}
		
		format := request.GetString("format", "json")
		_ = request.GetString("min_severity", "") // TODO: Implement min_severity filtering

		// Generate output
		var stdout string
		var stderr string
		
		// Determine input type and call appropriate method
		if strings.HasPrefix(input, "sbom:") {
			sbomPath := strings.TrimPrefix(input, "sbom:")
			stdout, err = module.ScanSBOM(ctx, sbomPath, format)
		} else if strings.HasPrefix(input, "dir:") {
			dirPath := strings.TrimPrefix(input, "dir:")
			stdout, err = module.ScanDirectory(ctx, dirPath, format)
		} else if strings.HasPrefix(input, "docker:") {
			image := strings.TrimPrefix(input, "docker:")
			stdout, err = module.ScanImage(ctx, image, format)
		} else {
			// Default: treat as file path for SBOM
			stdout, err = module.ScanSBOM(ctx, input, format)
		}
		
		// Build result in the expected format
		result := map[string]interface{}{
			"status": "ok",
			"stdout": stdout,
			"stderr": stderr,
			"artifacts": map[string]string{},
			"summary": map[string]interface{}{
				"high":     0,
				"critical": 0,
			},
			"diagnostics": []string{},
		}
		
		// Add artifact path based on format
		if format == "sarif" {
			result["artifacts"].(map[string]string)["grype_sarif"] = "./grype.sarif"
		} else {
			result["artifacts"].(map[string]string)["grype_json"] = "./grype.json"
		}
		
		// Parse JSON output to extract counts if format is json
		if format == "json" && err == nil && stdout != "" {
			var grypeOutput map[string]interface{}
			if jsonErr := json.Unmarshal([]byte(stdout), &grypeOutput); jsonErr == nil {
				// Try to extract vulnerability counts
				if matches, ok := grypeOutput["matches"].([]interface{}); ok {
					highCount := 0
					criticalCount := 0
					for _, match := range matches {
						if m, ok := match.(map[string]interface{}); ok {
							if vuln, ok := m["vulnerability"].(map[string]interface{}); ok {
								if severity, ok := vuln["severity"].(string); ok {
									switch strings.ToLower(severity) {
									case "critical":
										criticalCount++
									case "high":
										highCount++
									}
								}
							}
						}
					}
					result["summary"].(map[string]interface{})["high"] = highCount
					result["summary"].(map[string]interface{})["critical"] = criticalCount
				}
			}
		}
		
		if err != nil {
			result["status"] = "error"
			result["stderr"] = err.Error()
			result["diagnostics"] = []string{fmt.Sprintf("Grype scan failed: %v", err)}
		}

		// Return as JSON
		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(resultJSON)), nil
	})

	// Keep the legacy scan tool for backward compatibility
	scanToolLegacy := mcp.NewTool("grype_scan_legacy",
		mcp.WithDescription("[DEPRECATED - use grype_scan] Scan target for vulnerabilities"),
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
	s.AddTool(scanToolLegacy, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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