package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCosignAdvancedTools adds advanced Cosign workflows using real CLI capabilities
func AddCosignAdvancedTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Cosign keyless signing with OIDC
	signKeylessTool := mcp.NewTool("cosign_advanced_sign_keyless",
		mcp.WithDescription("Sign container image using keyless signing with OIDC"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to sign"),
			mcp.Required(),
		),
		mcp.WithString("identity_regex",
			mcp.Description("Identity regex for verification"),
		),
		mcp.WithString("issuer",
			mcp.Description("OIDC issuer URL"),
		),
	)
	s.AddTool(signKeylessTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"cosign", "sign", imageRef}
		return executeShipCommand(args)
	})

	// Cosign verify with certificate identity
	verifyIdentityTool := mcp.NewTool("cosign_advanced_verify_identity",
		mcp.WithDescription("Verify container image signature with certificate identity"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to verify"),
			mcp.Required(),
		),
		mcp.WithString("certificate_identity",
			mcp.Description("Certificate identity to verify"),
		),
		mcp.WithString("certificate_identity_regexp",
			mcp.Description("Certificate identity regex pattern"),
		),
		mcp.WithString("certificate_oidc_issuer",
			mcp.Description("Certificate OIDC issuer"),
		),
	)
	s.AddTool(verifyIdentityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"cosign", "verify"}
		
		if certIdentity := request.GetString("certificate_identity", ""); certIdentity != "" {
			args = append(args, "--certificate-identity", certIdentity)
		}
		if certIdentityRegexp := request.GetString("certificate_identity_regexp", ""); certIdentityRegexp != "" {
			args = append(args, "--certificate-identity-regexp", certIdentityRegexp)
		}
		if certOidcIssuer := request.GetString("certificate_oidc_issuer", ""); certOidcIssuer != "" {
			args = append(args, "--certificate-oidc-issuer", certOidcIssuer)
		}
		
		args = append(args, imageRef)
		return executeShipCommand(args)
	})

	// Cosign upload eBPF program
	uploadEbpfTool := mcp.NewTool("cosign_advanced_upload_ebpf",
		mcp.WithDescription("Upload eBPF program to OCI registry"),
		mcp.WithString("ebpf_path",
			mcp.Description("Path to eBPF program file"),
			mcp.Required(),
		),
		mcp.WithString("registry_url",
			mcp.Description("OCI registry URL"),
			mcp.Required(),
		),
	)
	s.AddTool(uploadEbpfTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ebpfPath := request.GetString("ebpf_path", "")
		registryUrl := request.GetString("registry_url", "")
		args := []string{"cosign", "upload", "blob", "-f", ebpfPath, registryUrl}
		return executeShipCommand(args)
	})

	// Cosign attest with predicate type
	attestWithTypeTool := mcp.NewTool("cosign_advanced_attest_type",
		mcp.WithDescription("Create attestation with specific predicate type"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference"),
			mcp.Required(),
		),
		mcp.WithString("predicate_type",
			mcp.Description("Predicate type URI (e.g., https://slsa.dev/provenance/v0.2)"),
			mcp.Required(),
		),
		mcp.WithString("predicate_file",
			mcp.Description("Path to predicate JSON file"),
			mcp.Required(),
		),
		mcp.WithString("key",
			mcp.Description("Path to signing key"),
		),
	)
	s.AddTool(attestWithTypeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		predicateType := request.GetString("predicate_type", "")
		predicateFile := request.GetString("predicate_file", "")
		
		args := []string{"cosign", "attest", "--predicate", predicateFile, "--type", predicateType}
		
		if key := request.GetString("key", ""); key != "" {
			args = append(args, "--key", key)
		}
		
		args = append(args, imageRef)
		return executeShipCommand(args)
	})

	// Cosign verify attestation with type and policy
	verifyAttestationAdvancedTool := mcp.NewTool("cosign_advanced_verify_attestation",
		mcp.WithDescription("Verify attestation with specific type and policy"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference"),
			mcp.Required(),
		),
		mcp.WithString("type",
			mcp.Description("Attestation type to verify"),
		),
		mcp.WithString("policy",
			mcp.Description("Policy file path for verification"),
		),
		mcp.WithString("key",
			mcp.Description("Public key for verification"),
		),
	)
	s.AddTool(verifyAttestationAdvancedTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"cosign", "verify-attestation"}
		
		if attestationType := request.GetString("type", ""); attestationType != "" {
			args = append(args, "--type", attestationType)
		}
		if policy := request.GetString("policy", ""); policy != "" {
			args = append(args, "--policy", policy)
		}
		if key := request.GetString("key", ""); key != "" {
			args = append(args, "--key", key)
		}
		
		args = append(args, imageRef)
		return executeShipCommand(args)
	})

	// Cosign verify with offline bundle
	verifyOfflineTool := mcp.NewTool("cosign_advanced_verify_offline",
		mcp.WithDescription("Verify signatures using offline bundle"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference"),
			mcp.Required(),
		),
		mcp.WithString("bundle",
			mcp.Description("Path to offline bundle file"),
			mcp.Required(),
		),
		mcp.WithString("certificate_identity",
			mcp.Description("Certificate identity to verify"),
		),
		mcp.WithString("certificate_oidc_issuer",
			mcp.Description("Certificate OIDC issuer"),
		),
	)
	s.AddTool(verifyOfflineTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		bundle := request.GetString("bundle", "")
		
		args := []string{"cosign", "verify", "--bundle", bundle}
		
		if certIdentity := request.GetString("certificate_identity", ""); certIdentity != "" {
			args = append(args, "--certificate-identity", certIdentity)
		}
		if certOidcIssuer := request.GetString("certificate_oidc_issuer", ""); certOidcIssuer != "" {
			args = append(args, "--certificate-oidc-issuer", certOidcIssuer)
		}
		
		args = append(args, imageRef)
		return executeShipCommand(args)
	})

	// Cosign triangulate signatures
	triangulateTool := mcp.NewTool("cosign_advanced_triangulate",
		mcp.WithDescription("Get signature image reference for a given image"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference"),
			mcp.Required(),
		),
	)
	s.AddTool(triangulateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"cosign", "triangulate", imageRef}
		return executeShipCommand(args)
	})

	// Cosign clean signatures
	cleanTool := mcp.NewTool("cosign_advanced_clean",
		mcp.WithDescription("Clean signatures from a given image"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference"),
			mcp.Required(),
		),
		mcp.WithString("type",
			mcp.Description("Type of signatures to clean (signature, attestation, sbom, all)"),
			mcp.Enum("signature", "attestation", "sbom", "all"),
		),
	)
	s.AddTool(cleanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"cosign", "clean"}
		
		if cleanType := request.GetString("type", ""); cleanType != "" {
			args = append(args, "--type", cleanType)
		}
		
		args = append(args, imageRef)
		return executeShipCommand(args)
	})
}