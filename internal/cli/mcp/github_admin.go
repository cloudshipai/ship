package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddGitHubAdminTools adds GitHub administration MCP tool implementations
func AddGitHubAdminTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// GitHub admin audit repositories tool
	auditRepositoriesTool := mcp.NewTool("github_admin_audit_repositories",
		mcp.WithDescription("Audit GitHub repositories for security and compliance"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("audit_type",
			mcp.Description("Type of audit (security, compliance, permissions)"),
		),
	)
	s.AddTool(auditRepositoriesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		args := []string{"security", "github-admin", "audit-repos", organization}
		if auditType := request.GetString("audit_type", ""); auditType != "" {
			args = append(args, "--type", auditType)
		}
		return executeShipCommand(args)
	})

	// GitHub admin manage branch protection tool
	manageBranchProtectionTool := mcp.NewTool("github_admin_manage_branch_protection",
		mcp.WithDescription("Manage branch protection rules across repositories"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("branch_pattern",
			mcp.Description("Branch pattern to protect (e.g., main, master)"),
			mcp.Required(),
		),
		mcp.WithString("protection_rules",
			mcp.Description("Protection rules to apply (reviews, checks, etc.)"),
		),
	)
	s.AddTool(manageBranchProtectionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		branchPattern := request.GetString("branch_pattern", "")
		args := []string{"security", "github-admin", "branch-protection", organization, branchPattern}
		if rules := request.GetString("protection_rules", ""); rules != "" {
			args = append(args, "--rules", rules)
		}
		return executeShipCommand(args)
	})

	// GitHub admin manage team permissions tool
	manageTeamPermissionsTool := mcp.NewTool("github_admin_manage_team_permissions",
		mcp.WithDescription("Manage team permissions across repositories"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("team_name",
			mcp.Description("Team name to manage"),
			mcp.Required(),
		),
		mcp.WithString("permission_level",
			mcp.Description("Permission level (read, write, admin)"),
			mcp.Required(),
		),
	)
	s.AddTool(manageTeamPermissionsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		teamName := request.GetString("team_name", "")
		permissionLevel := request.GetString("permission_level", "")
		args := []string{"security", "github-admin", "team-permissions", organization, teamName, permissionLevel}
		return executeShipCommand(args)
	})

	// GitHub admin enforce security policies tool
	enforceSecurityPoliciesTool := mcp.NewTool("github_admin_enforce_security_policies",
		mcp.WithDescription("Enforce organization-wide security policies"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("policy_set",
			mcp.Description("Security policy set to enforce (basic, strict, enterprise)"),
			mcp.Required(),
		),
	)
	s.AddTool(enforceSecurityPoliciesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		policySet := request.GetString("policy_set", "")
		args := []string{"security", "github-admin", "enforce-policies", organization, "--policy-set", policySet}
		return executeShipCommand(args)
	})

	// GitHub admin generate compliance report tool
	generateComplianceReportTool := mcp.NewTool("github_admin_generate_compliance_report",
		mcp.WithDescription("Generate compliance report for GitHub organization"),
		mcp.WithString("organization",
			mcp.Description("GitHub organization name"),
			mcp.Required(),
		),
		mcp.WithString("report_format",
			mcp.Description("Report format (json, csv, pdf)"),
		),
	)
	s.AddTool(generateComplianceReportTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		organization := request.GetString("organization", "")
		args := []string{"security", "github-admin", "compliance-report", organization}
		if format := request.GetString("report_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// GitHub admin get version tool
	getVersionTool := mcp.NewTool("github_admin_get_version",
		mcp.WithDescription("Get GitHub admin tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "github-admin", "--version"}
		return executeShipCommand(args)
	})
}