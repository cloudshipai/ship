package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddCosignTools adds Cosign (container signing and verification) MCP tool implementations
func AddCosignTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addCosignToolsDirect(s)
}

// addCosignToolsDirect implements direct Dagger calls for Cosign tools
func addCosignToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageName := request.GetString("image_name", "")
		keyPath := request.GetString("key_path", "")
		keyless := request.GetBool("keyless", false)

		// Create Cosign module and sign image with options
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.SignImageWithOptions(ctx, imageName, keyPath, keyless)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign sign image failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageName := request.GetString("image_name", "")
		keyPath := request.GetString("key_path", "")
		keyless := request.GetBool("keyless", false)

		// Create Cosign module and verify image with options
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.VerifyImageWithOptions(ctx, imageName, keyPath, keyless)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign verify image failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Cosign generate key pair tool
	generateKeyTool := mcp.NewTool("cosign_generate_key",
		mcp.WithDescription("Generate Cosign key pair for signing"),
		mcp.WithString("output_path",
			mcp.Description("Output path for key files"),
		),
	)
	s.AddTool(generateKeyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		outputPath := request.GetString("output_path", "/tmp")

		// Create Cosign module and generate key pair
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.GenerateKeyPair(ctx, outputPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign generate key failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageName := request.GetString("image_name", "")
		predicatePath := request.GetString("predicate_path", "")
		keyPath := request.GetString("key_path", "")

		// Create Cosign module and attest with options
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.AttestWithOptions(ctx, imageName, predicatePath, keyPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign attest failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		imageName := request.GetString("image_name", "")
		keyPath := request.GetString("key_path", "")
		policyPath := request.GetString("policy_path", "")

		// Create Cosign module and verify attestation with options
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.VerifyAttestationWithOptions(ctx, imageName, keyPath, policyPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign verify attestation failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		blobPath := request.GetString("blob_path", "")
		keyPath := request.GetString("key_path", "")
		outputSignature := request.GetString("output_signature", "")

		// Create Cosign module and sign blob
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.SignBlob(ctx, blobPath, keyPath, outputSignature)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign sign blob failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		blobPath := request.GetString("blob_path", "")
		signaturePath := request.GetString("signature_path", "")
		keyPath := request.GetString("key_path", "")

		// Create Cosign module and verify blob
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.VerifyBlob(ctx, blobPath, signaturePath, keyPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign verify blob failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		blobPath := request.GetString("blob_path", "")
		registryURL := request.GetString("registry_url", "")

		// Create Cosign module and upload blob
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.UploadBlob(ctx, blobPath, registryURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign upload blob failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		wasmPath := request.GetString("wasm_path", "")
		registryURL := request.GetString("registry_url", "")

		// Create Cosign module and upload wasm
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.UploadWasm(ctx, wasmPath, registryURL)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign upload wasm failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		sourceImage := request.GetString("source_image", "")
		destinationImage := request.GetString("destination_image", "")

		// Create Cosign module and copy image
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.CopyImage(ctx, sourceImage, destinationImage)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign copy image failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Get parameters
		wasmArtifact := request.GetString("wasm_artifact", "")
		keyPath := request.GetString("key_path", "")
		keyless := request.GetBool("keyless", false)

		// Create Cosign module and sign wasm artifact (reuse SignImageWithOptions)
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.SignImageWithOptions(ctx, wasmArtifact, keyPath, keyless)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign sign wasm failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Cosign get version tool
	getVersionTool := mcp.NewTool("cosign_get_version",
		mcp.WithDescription("Get Cosign version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create Cosign module and get version
		cosignModule := modules.NewCosignModule(client)
		result, err := cosignModule.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cosign get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}