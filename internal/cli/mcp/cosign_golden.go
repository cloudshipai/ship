package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCosignAdvancedTools adds advanced Cosign workflows using real CLI capabilities
func AddCosignAdvancedTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addCosignAdvancedToolsDirect(s)
}

// addCosignAdvancedToolsDirect implements direct Dagger calls for advanced Cosign tools
func addCosignAdvancedToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageRef := request.GetString("image_ref", "")
		identityRegex := request.GetString("identity_regex", "")
		issuer := request.GetString("issuer", "")

		// Create Cosign Golden module and sign keyless
		cosignModule := modules.NewCosignGoldenModule(client)
		result, err := cosignModule.SignKeyless(ctx, imageRef, identityRegex, issuer)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign advanced sign keyless failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageRef := request.GetString("image_ref", "")
		certIdentity := request.GetString("certificate_identity", "")
		certIdentityRegexp := request.GetString("certificate_identity_regexp", "")
		certOidcIssuer := request.GetString("certificate_oidc_issuer", "")

		// Create Cosign Golden module and verify identity
		cosignModule := modules.NewCosignGoldenModule(client)
		result, err := cosignModule.VerifyIdentity(ctx, imageRef, certIdentity, certIdentityRegexp, certOidcIssuer)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign advanced verify identity failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		ebpfPath := request.GetString("ebpf_path", "")
		registryURL := request.GetString("registry_url", "")

		// Create Cosign Golden module and upload eBPF
		cosignModule := modules.NewCosignGoldenModule(client)
		result, err := cosignModule.UploadEBPF(ctx, ebpfPath, registryURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign advanced upload ebpf failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageRef := request.GetString("image_ref", "")
		predicateType := request.GetString("predicate_type", "")
		predicateFile := request.GetString("predicate_file", "")
		key := request.GetString("key", "")

		// Create Cosign Golden module and attest with type
		cosignModule := modules.NewCosignGoldenModule(client)
		result, err := cosignModule.AttestWithType(ctx, imageRef, predicateType, predicateFile, key)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign advanced attest type failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageRef := request.GetString("image_ref", "")
		attestationType := request.GetString("type", "")
		policy := request.GetString("policy", "")
		key := request.GetString("key", "")

		// Create Cosign Golden module and verify attestation advanced
		cosignModule := modules.NewCosignGoldenModule(client)
		result, err := cosignModule.VerifyAttestationAdvanced(ctx, imageRef, attestationType, policy, key)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign advanced verify attestation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageRef := request.GetString("image_ref", "")
		bundle := request.GetString("bundle", "")
		certIdentity := request.GetString("certificate_identity", "")
		certOidcIssuer := request.GetString("certificate_oidc_issuer", "")

		// Create Cosign Golden module and verify offline
		cosignModule := modules.NewCosignGoldenModule(client)
		result, err := cosignModule.VerifyOffline(ctx, imageRef, bundle, certIdentity, certOidcIssuer)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign advanced verify offline failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageRef := request.GetString("image_ref", "")

		// Create Cosign Golden module and triangulate
		cosignModule := modules.NewCosignGoldenModule(client)
		result, err := cosignModule.Triangulate(ctx, imageRef)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign advanced triangulate failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageRef := request.GetString("image_ref", "")
		cleanType := request.GetString("type", "")

		// Create Cosign Golden module and clean
		cosignModule := modules.NewCosignGoldenModule(client)
		result, err := cosignModule.Clean(ctx, imageRef, cleanType)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign advanced clean failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}