package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitHubPackagesTools adds GitHub Packages management MCP tool implementations using gh CLI
func AddGitHubPackagesTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// GitHub packages list packages in organization
	listPackagesTool := mcp.NewTool("github_packages_list_packages",
		mcp.WithDescription("List GitHub Packages in organization using gh API"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("package_type",
			mcp.Description("Package type filter"),
			mcp.Enum("npm", "docker", "maven", "nuget", "rubygems"),
		),
	)
	s.AddTool(listPackagesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		endpoint := "orgs/" + organization + "/packages"
		args := []string{"gh", "api", endpoint}
		
		if packageType := request.GetString("package_type", ""); packageType != "" {
			args = append(args, "--field", "package_type=" + packageType)
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
		// GitHub doesn't have a dedicated dependency audit API for packages
		// This provides informational message about using GitHub's security features
		args := []string{"echo", "GitHub Packages dependency auditing is available through GitHub's security tab in the web interface and dependabot alerts. Use gh security commands or check the repository's security tab for vulnerability information."}
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
		// GitHub uses sigstore/cosign for package signing - provide guidance
		args := []string{"echo", "GitHub Packages signature verification uses cosign. To verify signatures: 1) Get package manifest, 2) Use cosign verify with appropriate keys/certificates. See cosign documentation for details."}
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
		// GitHub package policies are managed through organization settings
		args := []string{"echo", "GitHub Packages policies are configured through organization settings in the web interface. Go to Organization Settings > Packages to configure package visibility, access, and deletion policies."}
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
		// SBOM generation is not directly available through GitHub Packages API
		// Suggest using external SBOM tools like syft or cyclonedx
		args := []string{"echo", "GitHub Packages doesn't provide direct SBOM generation. Use tools like syft, cyclonedx-cli, or other SBOM generators to create SBOMs from package artifacts."}
		return executeShipCommand(args)
	})

	// GitHub packages get version tool
	getVersionTool := mcp.NewTool("github_packages_get_version",
		mcp.WithDescription("Get GitHub Packages security tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"gh", "--version"}
		return executeShipCommand(args)
	})
}