package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitHubPackagesTools adds GitHub Packages security MCP tool implementations
func AddGitHubPackagesTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// GitHub packages scan vulnerabilities tool
	scanVulnerabilitiesTool := mcp.NewTool("github_packages_scan_vulnerabilities",
		mcp.WithDescription("Scan GitHub Packages for security vulnerabilities"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("package_type",
			mcp.Description("Package type (npm, docker, maven, nuget)"),
		),
		mcp.WithString("package_name",
			mcp.Description("Specific package name to scan (optional)"),
		),
	)
	s.AddTool(scanVulnerabilitiesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		args := []string{"security", "github-packages", "scan-vulnerabilities", organization}
		if packageType := request.GetString("package_type", ""); packageType != "" {
			args = append(args, "--type", packageType)
		}
		if packageName := request.GetString("package_name", ""); packageName != "" {
			args = append(args, "--package", packageName)
		}
		return executeShipCommand(args)
	})

	// GitHub packages audit dependencies tool
	auditDependenciesTool := mcp.NewTool("github_packages_audit_dependencies",
		mcp.WithDescription("Audit package dependencies for security issues"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("package_name",
			mcp.Description("Package name to audit"),
			mcp.Required(),
		),
		mcp.WithString("version",
			mcp.Description("Package version to audit (optional)"),
		),
	)
	s.AddTool(auditDependenciesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		packageName := request.GetString("package_name", "")
		args := []string{"security", "github-packages", "audit-dependencies", organization, packageName}
		if version := request.GetString("version", ""); version != "" {
			args = append(args, "--version", version)
		}
		return executeShipCommand(args)
	})

	// GitHub packages check signatures tool
	checkSignaturesTool := mcp.NewTool("github_packages_check_signatures",
		mcp.WithDescription("Verify package signatures and attestations"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("package_name",
			mcp.Description("Package name to verify"),
			mcp.Required(),
		),
		mcp.WithString("version",
			mcp.Description("Package version to verify"),
			mcp.Required(),
		),
	)
	s.AddTool(checkSignaturesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		packageName := request.GetString("package_name", "")
		version := request.GetString("version", "")
		args := []string{"security", "github-packages", "check-signatures", organization, packageName, version}
		return executeShipCommand(args)
	})

	// GitHub packages enforce policies tool
	enforcePolicesTool := mcp.NewTool("github_packages_enforce_policies",
		mcp.WithDescription("Enforce security policies for GitHub Packages"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("policy_file",
			mcp.Description("Path to policy configuration file"),
			mcp.Required(),
		),
	)
	s.AddTool(enforcePolicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		policyFile := request.GetString("policy_file", "")
		args := []string{"security", "github-packages", "enforce-policies", organization, "--policy-file", policyFile}
		return executeShipCommand(args)
	})

	// GitHub packages generate SBOM tool
	generateSBOMTool := mcp.NewTool("github_packages_generate_sbom",
		mcp.WithDescription("Generate Software Bill of Materials for packages"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("package_name",
			mcp.Description("Package name to generate SBOM for"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("SBOM output format (spdx, cyclone-dx)"),
		),
	)
	s.AddTool(generateSBOMTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		packageName := request.GetString("package_name", "")
		args := []string{"security", "github-packages", "generate-sbom", organization, packageName}
		if format := request.GetString("output_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// GitHub packages get version tool
	getVersionTool := mcp.NewTool("github_packages_get_version",
		mcp.WithDescription("Get GitHub Packages security tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "github-packages", "--version"}
		return executeShipCommand(args)
	})
}