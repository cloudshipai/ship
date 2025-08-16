package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGuacTools adds GUAC (Graph for Understanding Artifact Composition) MCP tool implementations
func AddGuacTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// GUAC ingest SBOM tool
	ingestSBOMTool := mcp.NewTool("guac_ingest_sbom",
		mcp.WithDescription("Ingest SBOM into GUAC knowledge graph"),
		mcp.WithString("sbom_path",
			mcp.Description("Path to SBOM file"),
			mcp.Required(),
		),
	)
	s.AddTool(ingestSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sbomPath := request.GetString("sbom_path", "")
		args := []string{"security", "guac", "--ingest", sbomPath}
		return executeShipCommand(args)
	})

	// GUAC analyze artifact tool
	analyzeArtifactTool := mcp.NewTool("guac_analyze_artifact",
		mcp.WithDescription("Analyze artifact using GUAC"),
		mcp.WithString("artifact_path",
			mcp.Description("Path to artifact file"),
			mcp.Required(),
		),
	)
	s.AddTool(analyzeArtifactTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		artifactPath := request.GetString("artifact_path", "")
		args := []string{"security", "guac", "--analyze", artifactPath}
		return executeShipCommand(args)
	})

	// GUAC query dependencies tool
	queryDependenciesTool := mcp.NewTool("guac_query_dependencies",
		mcp.WithDescription("Query package dependencies using GUAC"),
		mcp.WithString("package_name",
			mcp.Description("Package name to query"),
			mcp.Required(),
		),
	)
	s.AddTool(queryDependenciesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName := request.GetString("package_name", "")
		args := []string{"security", "guac", "--query-deps", packageName}
		return executeShipCommand(args)
	})

	// GUAC query vulnerabilities tool
	queryVulnsTool := mcp.NewTool("guac_query_vulnerabilities",
		mcp.WithDescription("Query vulnerabilities for package using GUAC"),
		mcp.WithString("package_name",
			mcp.Description("Package name to query vulnerabilities for"),
			mcp.Required(),
		),
	)
	s.AddTool(queryVulnsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName := request.GetString("package_name", "")
		args := []string{"security", "guac", "--query-vulns", packageName}
		return executeShipCommand(args)
	})

	// GUAC generate graph tool
	generateGraphTool := mcp.NewTool("guac_generate_graph",
		mcp.WithDescription("Generate dependency graph using GUAC"),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory"),
			mcp.Required(),
		),
	)
	s.AddTool(generateGraphTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectPath := request.GetString("project_path", "")
		args := []string{"security", "guac", "--generate-graph", projectPath}
		return executeShipCommand(args)
	})

	// GUAC analyze impact tool
	analyzeImpactTool := mcp.NewTool("guac_analyze_impact",
		mcp.WithDescription("Analyze vulnerability impact using GUAC"),
		mcp.WithString("vuln_id",
			mcp.Description("Vulnerability ID to analyze"),
			mcp.Required(),
		),
	)
	s.AddTool(analyzeImpactTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		vulnID := request.GetString("vuln_id", "")
		args := []string{"security", "guac", "--analyze-impact", vulnID}
		return executeShipCommand(args)
	})

	// GUAC collect files tool
	collectFilesTool := mcp.NewTool("guac_collect_files",
		mcp.WithDescription("Collect and analyze project files using GUAC"),
		mcp.WithString("project_path",
			mcp.Description("Path to project directory"),
			mcp.Required(),
		),
	)
	s.AddTool(collectFilesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectPath := request.GetString("project_path", "")
		args := []string{"security", "guac", "--collect", projectPath}
		return executeShipCommand(args)
	})

	// GUAC validate attestation tool
	validateAttestationTool := mcp.NewTool("guac_validate_attestation",
		mcp.WithDescription("Validate attestation using GUAC"),
		mcp.WithString("attestation_path",
			mcp.Description("Path to attestation file"),
			mcp.Required(),
		),
	)
	s.AddTool(validateAttestationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		attestationPath := request.GetString("attestation_path", "")
		args := []string{"security", "guac", "--validate", attestationPath}
		return executeShipCommand(args)
	})
}