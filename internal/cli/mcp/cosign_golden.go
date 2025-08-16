package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCosignGoldenTools adds Cosign Golden (advanced signing workflows) MCP tool implementations
func AddCosignGoldenTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Cosign sign keyless tool
	signKeylessTool := mcp.NewTool("cosign_golden_sign_keyless",
		mcp.WithDescription("Sign container image using keyless signing with OIDC"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to sign"),
			mcp.Required(),
		),
		mcp.WithString("identity",
			mcp.Description("OIDC identity for keyless signing"),
			mcp.Required(),
		),
		mcp.WithString("issuer",
			mcp.Description("OIDC issuer for keyless signing"),
			mcp.Required(),
		),
	)
	s.AddTool(signKeylessTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		identity := request.GetString("identity", "")
		issuer := request.GetString("issuer", "")
		args := []string{"security", "cosign-golden", "--sign-keyless", imageRef, "--identity", identity, "--issuer", issuer}
		return executeShipCommand(args)
	})

	// Cosign verify keyless tool
	verifyKeylessTool := mcp.NewTool("cosign_golden_verify_keyless",
		mcp.WithDescription("Verify container image signature using keyless verification"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to verify"),
			mcp.Required(),
		),
		mcp.WithString("identity",
			mcp.Description("OIDC identity for keyless verification"),
			mcp.Required(),
		),
		mcp.WithString("issuer",
			mcp.Description("OIDC issuer for keyless verification"),
			mcp.Required(),
		),
	)
	s.AddTool(verifyKeylessTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		identity := request.GetString("identity", "")
		issuer := request.GetString("issuer", "")
		args := []string{"security", "cosign-golden", "--verify-keyless", imageRef, "--identity", identity, "--issuer", issuer}
		return executeShipCommand(args)
	})

	// Cosign sign golden pipeline tool
	signGoldenPipelineTool := mcp.NewTool("cosign_golden_sign_pipeline",
		mcp.WithDescription("Sign container image using golden pipeline workflow"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference to sign"),
			mcp.Required(),
		),
		mcp.WithString("build_metadata",
			mcp.Description("Build metadata JSON string"),
		),
		mcp.WithString("security_attestations",
			mcp.Description("Security attestations JSON string"),
		),
	)
	s.AddTool(signGoldenPipelineTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"security", "cosign-golden", "--sign-pipeline", imageRef}
		if buildMetadata := request.GetString("build_metadata", ""); buildMetadata != "" {
			args = append(args, "--build-metadata", buildMetadata)
		}
		if securityAttestations := request.GetString("security_attestations", ""); securityAttestations != "" {
			args = append(args, "--security-attestations", securityAttestations)
		}
		return executeShipCommand(args)
	})

	// Cosign generate attestation tool
	generateAttestationTool := mcp.NewTool("cosign_golden_generate_attestation",
		mcp.WithDescription("Generate attestation for container image"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference"),
			mcp.Required(),
		),
		mcp.WithString("attestation_type",
			mcp.Description("Type of attestation (slsa, spdx, vuln, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("predicate_data",
			mcp.Description("Predicate data for attestation"),
			mcp.Required(),
		),
	)
	s.AddTool(generateAttestationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		attestationType := request.GetString("attestation_type", "")
		predicateData := request.GetString("predicate_data", "")
		args := []string{"security", "cosign-golden", "--generate-attestation", imageRef, "--type", attestationType, "--predicate", predicateData}
		return executeShipCommand(args)
	})

	// Cosign verify attestation tool
	verifyAttestationTool := mcp.NewTool("cosign_golden_verify_attestation",
		mcp.WithDescription("Verify attestation for container image"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference"),
			mcp.Required(),
		),
		mcp.WithString("attestation_type",
			mcp.Description("Type of attestation to verify"),
			mcp.Required(),
		),
		mcp.WithString("identity",
			mcp.Description("OIDC identity for verification"),
			mcp.Required(),
		),
		mcp.WithString("issuer",
			mcp.Description("OIDC issuer for verification"),
			mcp.Required(),
		),
	)
	s.AddTool(verifyAttestationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		attestationType := request.GetString("attestation_type", "")
		identity := request.GetString("identity", "")
		issuer := request.GetString("issuer", "")
		args := []string{"security", "cosign-golden", "--verify-attestation", imageRef, "--type", attestationType, "--identity", identity, "--issuer", issuer}
		return executeShipCommand(args)
	})

	// Cosign copy signatures tool
	copySignaturesTool := mcp.NewTool("cosign_golden_copy_signatures",
		mcp.WithDescription("Copy signatures from source to destination image"),
		mcp.WithString("source_ref",
			mcp.Description("Source container image reference"),
			mcp.Required(),
		),
		mcp.WithString("destination_ref",
			mcp.Description("Destination container image reference"),
			mcp.Required(),
		),
		mcp.WithBoolean("force",
			mcp.Description("Force copy even if destination has signatures"),
		),
	)
	s.AddTool(copySignaturesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sourceRef := request.GetString("source_ref", "")
		destinationRef := request.GetString("destination_ref", "")
		args := []string{"security", "cosign-golden", "--copy", sourceRef, destinationRef}
		if force := request.GetBool("force", false); force {
			args = append(args, "--force")
		}
		return executeShipCommand(args)
	})

	// Cosign tree view tool
	treeViewTool := mcp.NewTool("cosign_golden_tree_view",
		mcp.WithDescription("Display tree view of image signatures and attestations"),
		mcp.WithString("image_ref",
			mcp.Description("Container image reference"),
			mcp.Required(),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format (text, json)"),
			mcp.Enum("text", "json"),
		),
	)
	s.AddTool(treeViewTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageRef := request.GetString("image_ref", "")
		args := []string{"security", "cosign-golden", "--tree", imageRef}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "--format", outputFormat)
		}
		return executeShipCommand(args)
	})

	// Cosign get version tool
	getVersionTool := mcp.NewTool("cosign_golden_get_version",
		mcp.WithDescription("Get Cosign Golden version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "cosign-golden", "--version"}
		return executeShipCommand(args)
	})
}