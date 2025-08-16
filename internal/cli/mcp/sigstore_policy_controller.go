package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSigstorePolicyControllerTools adds Sigstore Policy Controller MCP tool implementations
func AddSigstorePolicyControllerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Validate policy tool
	validatePolicyTool := mcp.NewTool("sigstore_validate_policy",
		mcp.WithDescription("Validate Sigstore policy syntax and structure"),
		mcp.WithString("policy_path",
			mcp.Description("Path to policy file"),
			mcp.Required(),
		),
	)
	s.AddTool(validatePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "sigstore-policy-controller", "--validate", policyPath}
		return executeShipCommand(args)
	})

	// Test policy tool
	testPolicyTool := mcp.NewTool("sigstore_test_policy",
		mcp.WithDescription("Test Sigstore policy against an image"),
		mcp.WithString("policy_path",
			mcp.Description("Path to policy file"),
			mcp.Required(),
		),
		mcp.WithString("image_name",
			mcp.Description("Container image name to test"),
			mcp.Required(),
		),
	)
	s.AddTool(testPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyPath := request.GetString("policy_path", "")
		imageName := request.GetString("image_name", "")
		args := []string{"security", "sigstore-policy-controller", "--test", policyPath, "--image", imageName}
		return executeShipCommand(args)
	})

	// Verify signature tool
	verifySignatureTool := mcp.NewTool("sigstore_verify_signature",
		mcp.WithDescription("Verify container image signature"),
		mcp.WithString("image_name",
			mcp.Description("Container image name"),
			mcp.Required(),
		),
		mcp.WithString("public_key_path",
			mcp.Description("Path to public key file"),
			mcp.Required(),
		),
	)
	s.AddTool(verifySignatureTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		publicKeyPath := request.GetString("public_key_path", "")
		args := []string{"security", "sigstore-policy-controller", "--verify", imageName, "--public-key", publicKeyPath}
		return executeShipCommand(args)
	})

	// Generate policy template tool
	generateTemplateTool := mcp.NewTool("sigstore_generate_template",
		mcp.WithDescription("Generate Sigstore policy template"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
			mcp.Required(),
		),
		mcp.WithString("key_ref",
			mcp.Description("Key reference for policy generation"),
			mcp.Required(),
		),
	)
	s.AddTool(generateTemplateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		keyRef := request.GetString("key_ref", "")
		args := []string{"security", "sigstore-policy-controller", "--generate", "--namespace", namespace, "--key-ref", keyRef}
		return executeShipCommand(args)
	})

	// Validate manifest tool
	validateManifestTool := mcp.NewTool("sigstore_validate_manifest",
		mcp.WithDescription("Validate Kubernetes manifest against policy"),
		mcp.WithString("manifest_path",
			mcp.Description("Path to manifest file"),
			mcp.Required(),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to policy file"),
			mcp.Required(),
		),
	)
	s.AddTool(validateManifestTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		manifestPath := request.GetString("manifest_path", "")
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "sigstore-policy-controller", "--validate-manifest", manifestPath, "--policy", policyPath}
		return executeShipCommand(args)
	})

	// Check compliance tool
	checkComplianceTool := mcp.NewTool("sigstore_check_compliance",
		mcp.WithDescription("Check manifests compliance with policy"),
		mcp.WithString("manifests_path",
			mcp.Description("Path to manifests directory"),
			mcp.Required(),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to policy file"),
			mcp.Required(),
		),
	)
	s.AddTool(checkComplianceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		manifestsPath := request.GetString("manifests_path", "")
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "sigstore-policy-controller", "--check-compliance", manifestsPath, "--policy", policyPath}
		return executeShipCommand(args)
	})

	// List policies tool
	listPoliciesTool := mcp.NewTool("sigstore_list_policies",
		mcp.WithDescription("List available Sigstore policies"),
		mcp.WithString("policies_path",
			mcp.Description("Path to policies directory"),
			mcp.Required(),
		),
	)
	s.AddTool(listPoliciesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policiesPath := request.GetString("policies_path", "")
		args := []string{"security", "sigstore-policy-controller", "--list", policiesPath}
		return executeShipCommand(args)
	})

	// Audit images tool
	auditImagesTool := mcp.NewTool("sigstore_audit_images",
		mcp.WithDescription("Audit container images in namespace"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
			mcp.Required(),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to policy file"),
			mcp.Required(),
		),
	)
	s.AddTool(auditImagesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := request.GetString("namespace", "")
		policyPath := request.GetString("policy_path", "")
		args := []string{"security", "sigstore-policy-controller", "--audit", "--namespace", namespace, "--policy", policyPath}
		return executeShipCommand(args)
	})
}