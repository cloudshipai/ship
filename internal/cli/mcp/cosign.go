package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCosignTools adds Cosign (container signing and verification) MCP tool implementations
func AddCosignTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Cosign sign container image tool
	signImageTool := mcp.NewTool("cosign_sign_image",
		mcp.WithDescription("Sign container image using Cosign"),
		mcp.WithString("image_name",
			mcp.Description("Container image name to sign"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to signing key"),
		),
		mcp.WithBoolean("keyless",
			mcp.Description("Use keyless signing with OIDC"),
		),
	)
	s.AddTool(signImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"cosign", "sign"}
		if keyPath := request.GetString("key_path", ""); keyPath != "" {
			args = append(args, "--key", keyPath)
		}
		args = append(args, imageName)
		return executeShipCommand(args)
	})

	// Cosign verify container image tool
	verifyImageTool := mcp.NewTool("cosign_verify_image",
		mcp.WithDescription("Verify container image signature using Cosign"),
		mcp.WithString("image_name",
			mcp.Description("Container image name to verify"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to verification key"),
		),
		mcp.WithBoolean("keyless",
			mcp.Description("Use keyless verification"),
		),
	)
	s.AddTool(verifyImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"cosign", "verify"}
		if keyPath := request.GetString("key_path", ""); keyPath != "" {
			args = append(args, "--key", keyPath)
		}
		args = append(args, imageName)
		return executeShipCommand(args)
	})

	// Cosign generate key pair tool
	generateKeyTool := mcp.NewTool("cosign_generate_key",
		mcp.WithDescription("Generate Cosign key pair for signing"),
		mcp.WithString("output_path",
			mcp.Description("Output path for key files"),
		),
	)
	s.AddTool(generateKeyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cosign", "generate-key-pair"}
		return executeShipCommand(args)
	})

	// Cosign attest and sign tool
	attestTool := mcp.NewTool("cosign_attest",
		mcp.WithDescription("Create and sign attestation for container image"),
		mcp.WithString("image_name",
			mcp.Description("Container image name"),
			mcp.Required(),
		),
		mcp.WithString("predicate_path",
			mcp.Description("Path to predicate file"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to signing key"),
		),
	)
	s.AddTool(attestTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		predicatePath := request.GetString("predicate_path", "")
		args := []string{"cosign", "attest", "--predicate", predicatePath}
		if keyPath := request.GetString("key_path", ""); keyPath != "" {
			args = append(args, "--key", keyPath)
		}
		args = append(args, imageName)
		return executeShipCommand(args)
	})

	// Cosign verify attestation tool
	verifyAttestationTool := mcp.NewTool("cosign_verify_attestation",
		mcp.WithDescription("Verify attestation for container image"),
		mcp.WithString("image_name",
			mcp.Description("Container image name"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to verification key"),
		),
		mcp.WithString("policy_path",
			mcp.Description("Path to policy file for verification"),
		),
	)
	s.AddTool(verifyAttestationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"cosign", "verify-attestation"}
		if keyPath := request.GetString("key_path", ""); keyPath != "" {
			args = append(args, "--key", keyPath)
		}
		if policyPath := request.GetString("policy_path", ""); policyPath != "" {
			args = append(args, "--policy", policyPath)
		}
		args = append(args, imageName)
		return executeShipCommand(args)
	})

	// Cosign sign blob tool
	signBlobTool := mcp.NewTool("cosign_sign_blob",
		mcp.WithDescription("Sign arbitrary blob using Cosign"),
		mcp.WithString("blob_path",
			mcp.Description("Path to blob file to sign"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to signing key"),
		),
		mcp.WithString("output_signature",
			mcp.Description("Output path for signature"),
		),
	)
	s.AddTool(signBlobTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blobPath := request.GetString("blob_path", "")
		args := []string{"cosign", "sign-blob"}
		if keyPath := request.GetString("key_path", ""); keyPath != "" {
			args = append(args, "--key", keyPath)
		}
		if outputSig := request.GetString("output_signature", ""); outputSig != "" {
			args = append(args, "--output-signature", outputSig)
		}
		args = append(args, blobPath)
		return executeShipCommand(args)
	})

	// Cosign verify blob tool
	verifyBlobTool := mcp.NewTool("cosign_verify_blob",
		mcp.WithDescription("Verify blob signature using Cosign"),
		mcp.WithString("blob_path",
			mcp.Description("Path to blob file"),
			mcp.Required(),
		),
		mcp.WithString("signature_path",
			mcp.Description("Path to signature file"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to verification key"),
		),
	)
	s.AddTool(verifyBlobTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blobPath := request.GetString("blob_path", "")
		signaturePath := request.GetString("signature_path", "")
		args := []string{"cosign", "verify-blob", "--signature", signaturePath}
		if keyPath := request.GetString("key_path", ""); keyPath != "" {
			args = append(args, "--key", keyPath)
		}
		args = append(args, blobPath)
		return executeShipCommand(args)
	})

	// Cosign upload blob tool
	uploadBlobTool := mcp.NewTool("cosign_upload_blob",
		mcp.WithDescription("Upload generic artifact as a blob to registry"),
		mcp.WithString("blob_path",
			mcp.Description("Path to blob file to upload"),
			mcp.Required(),
		),
		mcp.WithString("registry_url",
			mcp.Description("Registry URL to upload to"),
			mcp.Required(),
		),
	)
	s.AddTool(uploadBlobTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		blobPath := request.GetString("blob_path", "")
		registryURL := request.GetString("registry_url", "")
		args := []string{"cosign", "upload", "blob", "-f", blobPath, registryURL}
		return executeShipCommand(args)
	})

	// Cosign upload wasm tool
	uploadWasmTool := mcp.NewTool("cosign_upload_wasm",
		mcp.WithDescription("Upload WebAssembly module to registry"),
		mcp.WithString("wasm_path",
			mcp.Description("Path to WebAssembly file to upload"),
			mcp.Required(),
		),
		mcp.WithString("registry_url",
			mcp.Description("Registry URL to upload to"),
			mcp.Required(),
		),
	)
	s.AddTool(uploadWasmTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wasmPath := request.GetString("wasm_path", "")
		registryURL := request.GetString("registry_url", "")
		args := []string{"cosign", "upload", "wasm", "-f", wasmPath, registryURL}
		return executeShipCommand(args)
	})

	// Cosign download/copy image tool (using copy command)
	copyImageTool := mcp.NewTool("cosign_copy_image",
		mcp.WithDescription("Copy images between registries"),
		mcp.WithString("source_image",
			mcp.Description("Source image reference"),
			mcp.Required(),
		),
		mcp.WithString("destination_image",
			mcp.Description("Destination image reference"),
			mcp.Required(),
		),
	)
	s.AddTool(copyImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sourceImage := request.GetString("source_image", "")
		destinationImage := request.GetString("destination_image", "")
		args := []string{"cosign", "copy", sourceImage, destinationImage}
		return executeShipCommand(args)
	})

	// Cosign sign wasm tool
	signWasmTool := mcp.NewTool("cosign_sign_wasm",
		mcp.WithDescription("Sign WebAssembly module using Cosign"),
		mcp.WithString("wasm_artifact",
			mcp.Description("WebAssembly artifact reference to sign"),
			mcp.Required(),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to signing key"),
		),
		mcp.WithBoolean("keyless",
			mcp.Description("Use keyless signing with OIDC"),
		),
	)
	s.AddTool(signWasmTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wasmArtifact := request.GetString("wasm_artifact", "")
		args := []string{"cosign", "sign"}
		if keyPath := request.GetString("key_path", ""); keyPath != "" {
			args = append(args, "--key", keyPath)
		}
		args = append(args, wasmArtifact)
		return executeShipCommand(args)
	})

	// Cosign get version tool
	getVersionTool := mcp.NewTool("cosign_get_version",
		mcp.WithDescription("Get Cosign version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"cosign", "version"}
		return executeShipCommand(args)
	})
}