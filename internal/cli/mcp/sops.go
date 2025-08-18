package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSOPSTools adds SOPS (Secrets OPerationS) MCP tool implementations using real sops CLI commands
func AddSOPSTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// SOPS encrypt file tool
	encryptFileTool := mcp.NewTool("sops_encrypt_file",
		mcp.WithDescription("Encrypt file using real sops CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to file to encrypt"),
			mcp.Required(),
		),
		mcp.WithString("pgp",
			mcp.Description("PGP key fingerprint for encryption"),
		),
		mcp.WithString("kms",
			mcp.Description("AWS KMS key ARN for encryption"),
		),
		mcp.WithString("age",
			mcp.Description("Age recipient public key for encryption"),
		),
		mcp.WithString("gcp_kms",
			mcp.Description("GCP KMS key for encryption"),
		),
		mcp.WithString("azure_kv",
			mcp.Description("Azure Key Vault URL for encryption"),
		),
		mcp.WithString("input_type",
			mcp.Description("Input file format"),
			mcp.Enum("yaml", "json", "env", "ini"),
		),
		mcp.WithString("output_type",
			mcp.Description("Output file format"),
			mcp.Enum("yaml", "json", "env", "ini"),
		),
	)
	s.AddTool(encryptFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"sops", "--encrypt"}
		
		if pgp := request.GetString("pgp", ""); pgp != "" {
			args = append(args, "--pgp", pgp)
		}
		if kms := request.GetString("kms", ""); kms != "" {
			args = append(args, "--kms", kms)
		}
		if age := request.GetString("age", ""); age != "" {
			args = append(args, "--age", age)
		}
		if gcpKms := request.GetString("gcp_kms", ""); gcpKms != "" {
			args = append(args, "--gcp-kms", gcpKms)
		}
		if azureKv := request.GetString("azure_kv", ""); azureKv != "" {
			args = append(args, "--azure-kv", azureKv)
		}
		if inputType := request.GetString("input_type", ""); inputType != "" {
			args = append(args, "--input-type", inputType)
		}
		if outputType := request.GetString("output_type", ""); outputType != "" {
			args = append(args, "--output-type", outputType)
		}
		args = append(args, filePath)
		
		return executeShipCommand(args)
	})

	// SOPS decrypt file tool
	decryptFileTool := mcp.NewTool("sops_decrypt_file",
		mcp.WithDescription("Decrypt file using real sops CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to encrypted file to decrypt"),
			mcp.Required(),
		),
		mcp.WithString("output_type",
			mcp.Description("Output file format"),
			mcp.Enum("yaml", "json", "env", "ini"),
		),
		mcp.WithString("extract",
			mcp.Description("Extract specific key from decrypted file"),
		),
	)
	s.AddTool(decryptFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"sops", "--decrypt"}
		
		if outputType := request.GetString("output_type", ""); outputType != "" {
			args = append(args, "--output-type", outputType)
		}
		if extract := request.GetString("extract", ""); extract != "" {
			args = append(args, "--extract", extract)
		}
		args = append(args, filePath)
		
		return executeShipCommand(args)
	})

	// SOPS update keys tool
	updateKeysTool := mcp.NewTool("sops_update_keys",
		mcp.WithDescription("Update encryption keys for SOPS file using real sops CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to SOPS encrypted file"),
			mcp.Required(),
		),
		mcp.WithString("add_pgp",
			mcp.Description("Add PGP key fingerprint"),
		),
		mcp.WithString("rm_pgp",
			mcp.Description("Remove PGP key fingerprint"),
		),
		mcp.WithString("add_kms",
			mcp.Description("Add AWS KMS key ARN"),
		),
		mcp.WithString("rm_kms",
			mcp.Description("Remove AWS KMS key ARN"),
		),
		mcp.WithString("add_age",
			mcp.Description("Add age recipient public key"),
		),
		mcp.WithString("rm_age",
			mcp.Description("Remove age recipient public key"),
		),
		mcp.WithBoolean("yes",
			mcp.Description("Automatically confirm changes"),
		),
	)
	s.AddTool(updateKeysTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"sops", "updatekeys"}
		
		if addPgp := request.GetString("add_pgp", ""); addPgp != "" {
			args = append(args, "--add-pgp", addPgp)
		}
		if rmPgp := request.GetString("rm_pgp", ""); rmPgp != "" {
			args = append(args, "--rm-pgp", rmPgp)
		}
		if addKms := request.GetString("add_kms", ""); addKms != "" {
			args = append(args, "--add-kms", addKms)
		}
		if rmKms := request.GetString("rm_kms", ""); rmKms != "" {
			args = append(args, "--rm-kms", rmKms)
		}
		if addAge := request.GetString("add_age", ""); addAge != "" {
			args = append(args, "--add-age", addAge)
		}
		if rmAge := request.GetString("rm_age", ""); rmAge != "" {
			args = append(args, "--rm-age", rmAge)
		}
		if request.GetBool("yes", false) {
			args = append(args, "-y")
		}
		args = append(args, filePath)
		
		return executeShipCommand(args)
	})

	// SOPS edit file tool
	editFileTool := mcp.NewTool("sops_edit_file",
		mcp.WithDescription("Edit encrypted file using real sops CLI"),
		mcp.WithString("file_path",
			mcp.Description("Path to SOPS encrypted file to edit"),
			mcp.Required(),
		),
	)
	s.AddTool(editFileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath := request.GetString("file_path", "")
		args := []string{"sops", filePath}
		return executeShipCommand(args)
	})

	// SOPS version tool
	versionTool := mcp.NewTool("sops_version",
		mcp.WithDescription("Get SOPS version information using real sops CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"sops", "--version"}
		return executeShipCommand(args)
	})

	// SOPS publish keys tool
	publishKeysTool := mcp.NewTool("sops_publish_keys",
		mcp.WithDescription("Publish keys to keyserver using real sops CLI"),
		mcp.WithString("keyserver",
			mcp.Description("Keyserver URL to publish to"),
		),
	)
	s.AddTool(publishKeysTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"sops", "publish"}
		
		if keyserver := request.GetString("keyserver", ""); keyserver != "" {
			args = append(args, "--keyserver", keyserver)
		}
		
		return executeShipCommand(args)
	})
}