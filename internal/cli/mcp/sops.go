package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSOPSTools adds SOPS (Secrets OPerationS) MCP tool implementations
func AddSOPSTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// SOPS encrypt file tool
	encryptFileTool := mcp.NewTool("sops_encrypt_file",
		mcp.WithDescription("Encrypt file using SOPS"),
		mcp.WithString("file_path",
			mcp.Description("Path to file to encrypt"),
			mcp.Required(),
		),
		mcp.WithString("key_type",
			mcp.Description("Key type (pgp, kms, azure-kv, gcp-kms, age)"),
			mcp.Enum("pgp", "kms", "azure-kv", "gcp-kms", "age"),
		),
		mcp.WithString("key_id",
			mcp.Description("Key ID or identifier"),
		),
	)
	s.AddTool(encryptFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"security", "sops", "--encrypt", filePath}
		if keyType := request.GetString("key_type", ""); keyType != "" {
			args = append(args, "--key-type", keyType)
		}
		if keyId := request.GetString("key_id", ""); keyId != "" {
			args = append(args, "--key-id", keyId)
		}
		return executeShipCommand(args)
	})

	// SOPS decrypt file tool
	decryptFileTool := mcp.NewTool("sops_decrypt_file",
		mcp.WithDescription("Decrypt file using SOPS"),
		mcp.WithString("file_path",
			mcp.Description("Path to encrypted file to decrypt"),
			mcp.Required(),
		),
		mcp.WithString("output_path",
			mcp.Description("Output path for decrypted file"),
		),
	)
	s.AddTool(decryptFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"security", "sops", "--decrypt", filePath}
		if outputPath := request.GetString("output_path", ""); outputPath != "" {
			args = append(args, "--output", outputPath)
		}
		return executeShipCommand(args)
	})

	// SOPS rotate keys tool
	rotateKeysTool := mcp.NewTool("sops_rotate_keys",
		mcp.WithDescription("Rotate encryption keys for SOPS file"),
		mcp.WithString("file_path",
			mcp.Description("Path to SOPS encrypted file"),
			mcp.Required(),
		),
		mcp.WithString("new_key_id",
			mcp.Description("New key ID to rotate to"),
		),
	)
	s.AddTool(rotateKeysTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"security", "sops", "--rotate", filePath}
		if newKeyId := request.GetString("new_key_id", ""); newKeyId != "" {
			args = append(args, "--add-key", newKeyId)
		}
		return executeShipCommand(args)
	})

	// SOPS edit file tool
	editFileTool := mcp.NewTool("sops_edit_file",
		mcp.WithDescription("Edit encrypted file using SOPS"),
		mcp.WithString("file_path",
			mcp.Description("Path to SOPS encrypted file to edit"),
			mcp.Required(),
		),
		mcp.WithString("editor",
			mcp.Description("Editor to use (vim, nano, etc.)"),
		),
	)
	s.AddTool(editFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"security", "sops", "--edit", filePath}
		if editor := request.GetString("editor", ""); editor != "" {
			args = append(args, "--editor", editor)
		}
		return executeShipCommand(args)
	})

	// SOPS generate config tool
	generateConfigTool := mcp.NewTool("sops_generate_config",
		mcp.WithDescription("Generate SOPS configuration file"),
		mcp.WithString("output_path",
			mcp.Description("Output path for SOPS config file"),
		),
		mcp.WithString("key_type",
			mcp.Description("Default key type for configuration"),
			mcp.Enum("pgp", "kms", "azure-kv", "gcp-kms", "age"),
		),
	)
	s.AddTool(generateConfigTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "sops", "--generate-config"}
		if outputPath := request.GetString("output_path", ""); outputPath != "" {
			args = append(args, "--output", outputPath)
		}
		if keyType := request.GetString("key_type", ""); keyType != "" {
			args = append(args, "--key-type", keyType)
		}
		return executeShipCommand(args)
	})

	// SOPS validate file tool
	validateFileTool := mcp.NewTool("sops_validate_file",
		mcp.WithDescription("Validate SOPS encrypted file integrity"),
		mcp.WithString("file_path",
			mcp.Description("Path to SOPS encrypted file to validate"),
			mcp.Required(),
		),
	)
	s.AddTool(validateFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"security", "sops", "--validate", filePath}
		return executeShipCommand(args)
	})
}