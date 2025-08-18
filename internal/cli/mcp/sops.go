package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddSOPSTools adds SOPS (Secrets OPerationS) MCP tool implementations using direct Dagger calls
func AddSOPSTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addSOPSToolsDirect(s)
}

// addSOPSToolsDirect adds SOPS tools using direct Dagger module calls
func addSOPSToolsDirect(s *server.MCPServer) {
	// SOPS encrypt file tool
	encryptFileTool := mcp.NewTool("sops_encrypt_file",
		mcp.WithDescription("Encrypt file using real SOPS CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to file to encrypt"),
			mcp.Required(),
		),
		mcp.WithString("kms_arn",
			mcp.Description("AWS KMS key ARN for encryption"),
		),
		mcp.WithString("pgp_fingerprint",
			mcp.Description("PGP key fingerprint for encryption"),
		),
		mcp.WithString("age_public_key",
			mcp.Description("Age recipient public key for encryption"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for encrypted content"),
		),
		mcp.WithBoolean("in_place",
			mcp.Description("Encrypt file in place"),
		),
	)
	s.AddTool(encryptFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSOPSModule(client)

		// Get parameters
		filePath := request.GetString("file_path", "")
		kmsArn := request.GetString("kms_arn", "")
		pgpFingerprint := request.GetString("pgp_fingerprint", "")
		agePublicKey := request.GetString("age_public_key", "")
		outputFile := request.GetString("output_file", "")
		inPlace := request.GetBool("in_place", false)

		// Encrypt file
		output, err := module.EncryptFile(ctx, filePath, kmsArn, pgpFingerprint, agePublicKey, outputFile, inPlace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SOPS encrypt file failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// SOPS decrypt file tool
	decryptFileTool := mcp.NewTool("sops_decrypt_file",
		mcp.WithDescription("Decrypt file using real SOPS CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to encrypted file to decrypt"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for decrypted content"),
		),
	)
	s.AddTool(decryptFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSOPSModule(client)

		// Get parameters
		filePath := request.GetString("file_path", "")
		outputFile := request.GetString("output_file", "")

		// Decrypt file
		output, err := module.DecryptFile(ctx, filePath, outputFile)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SOPS decrypt file failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// SOPS update keys tool
	updateKeysTool := mcp.NewTool("sops_update_keys",
		mcp.WithDescription("Update encryption keys for SOPS file using real SOPS CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to SOPS encrypted file"),
			mcp.Required(),
		),
		mcp.WithString("add_kms",
			mcp.Description("Add AWS KMS key ARN"),
		),
		mcp.WithString("add_pgp",
			mcp.Description("Add PGP key fingerprint"),
		),
		mcp.WithString("add_age",
			mcp.Description("Add age recipient public key"),
		),
		mcp.WithBoolean("in_place",
			mcp.Description("Update file in place"),
		),
	)
	s.AddTool(updateKeysTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSOPSModule(client)

		// Get parameters
		filePath := request.GetString("file_path", "")
		addKms := request.GetString("add_kms", "")
		addPgp := request.GetString("add_pgp", "")
		addAge := request.GetString("add_age", "")
		inPlace := request.GetBool("in_place", false)

		// Update keys
		output, err := module.UpdateKeys(ctx, filePath, addKms, addPgp, addAge, inPlace)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SOPS update keys failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// SOPS edit file tool
	editFileTool := mcp.NewTool("sops_edit_file",
		mcp.WithDescription("Show decrypted content (edit mode not supported in containers) using real SOPS CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to SOPS encrypted file to show/edit"),
			mcp.Required(),
		),
	)
	s.AddTool(editFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSOPSModule(client)

		// Get parameters
		filePath := request.GetString("file_path", "")

		// Show file content
		output, err := module.EditFile(ctx, filePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SOPS edit file failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// SOPS version tool
	versionTool := mcp.NewTool("sops_version",
		mcp.WithDescription("Get SOPS version information using real SOPS CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSOPSModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SOPS get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// SOPS publish keys tool
	publishKeysTool := mcp.NewTool("sops_publish_keys",
		mcp.WithDescription("Show public key information using real SOPS CLI"),
		mcp.WithString("key_type",
			mcp.Description("Type of key to publish"),
			mcp.Enum("pgp", "age", "kms"),
		),
		mcp.WithString("key_path",
			mcp.Description("Path to key file (for PGP/age keys)"),
		),
	)
	s.AddTool(publishKeysTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSOPSModule(client)

		// Get parameters
		keyType := request.GetString("key_type", "")
		keyPath := request.GetString("key_path", "")

		// Publish keys
		output, err := module.PublishKeys(ctx, keyType, keyPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SOPS publish keys failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}