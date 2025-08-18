package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitHubPackagesTools adds GitHub Packages management MCP tool implementations using gh CLI
func AddGitHubPackagesTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addGitHubPackagesToolsDirect(s)
}

// addGitHubPackagesToolsDirect implements direct Dagger calls for GitHub Packages tools
func addGitHubPackagesToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		organization := request.GetString("organization", "")
		packageType := request.GetString("package_type", "")

		// Create GitHub Packages module and list packages
		githubPackagesModule := modules.NewGitHubPackagesModule(client)
		result, err := githubPackagesModule.ListPackagesSimple(ctx, organization, packageType)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github packages list packages failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		organization := request.GetString("organization", "")
		packageName := request.GetString("package_name", "")
		version := request.GetString("version", "")

		// Create GitHub Packages module and audit dependencies
		githubPackagesModule := modules.NewGitHubPackagesModule(client)
		result, err := githubPackagesModule.AuditDependenciesSimple(ctx, organization, packageName, version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github packages audit dependencies failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		organization := request.GetString("organization", "")
		packageName := request.GetString("package_name", "")
		version := request.GetString("version", "")

		// Create GitHub Packages module and check signatures
		githubPackagesModule := modules.NewGitHubPackagesModule(client)
		result, err := githubPackagesModule.CheckSignaturesSimple(ctx, organization, packageName, version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github packages check signatures failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		organization := request.GetString("organization", "")
		policyFile := request.GetString("policy_file", "")

		// Create GitHub Packages module and enforce policies
		githubPackagesModule := modules.NewGitHubPackagesModule(client)
		result, err := githubPackagesModule.EnforcePoliciesSimple(ctx, organization, policyFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github packages enforce policies failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		organization := request.GetString("organization", "")
		packageName := request.GetString("package_name", "")
		outputFormat := request.GetString("output_format", "")

		// Create GitHub Packages module and generate SBOM
		githubPackagesModule := modules.NewGitHubPackagesModule(client)
		result, err := githubPackagesModule.GenerateSBOMSimple(ctx, organization, packageName, outputFormat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github packages generate SBOM failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// GitHub packages get version tool
	getVersionTool := mcp.NewTool("github_packages_get_version",
		mcp.WithDescription("Get GitHub Packages security tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create GitHub Packages module and get version
		githubPackagesModule := modules.NewGitHubPackagesModule(client)
		result, err := githubPackagesModule.GetVersionSimple(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("github packages get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}