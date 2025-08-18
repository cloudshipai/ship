package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddRekorTools adds Rekor (transparency log) MCP tool implementations using real rekor-cli commands
func AddRekorTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Rekor upload artifact tool
	uploadTool := mcp.NewTool("rekor_upload_artifact",
		mcp.WithDescription("Upload artifact to Rekor transparency log using real rekor-cli"),
		mcp.WithString("artifact",
			mcp.Description("Path or URL to artifact file"),
			mcp.Required(),
		),
		mcp.WithString("signature",
			mcp.Description("Path or URL to signature file"),
			mcp.Required(),
		),
		mcp.WithString("public_key",
			mcp.Description("Path or URL to public key file"),
			mcp.Required(),
		),
		mcp.WithString("type",
			mcp.Description("Type of entry"),
		),
		mcp.WithString("pki_format",
			mcp.Description("Format of signature/public key"),
			mcp.Enum("pgp", "minisign", "x509", "ssh", "tuf"),
		),
	)
	s.AddTool(uploadTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"rekor-cli", "upload"}
		if artifact := request.GetString("artifact", ""); artifact != "" {
			args = append(args, "--artifact", artifact)
		}
		if signature := request.GetString("signature", ""); signature != "" {
			args = append(args, "--signature", signature)
		}
		if publicKey := request.GetString("public_key", ""); publicKey != "" {
			args = append(args, "--public-key", publicKey)
		}
		if entryType := request.GetString("type", ""); entryType != "" {
			args = append(args, "--type", entryType)
		}
		if pkiFormat := request.GetString("pki_format", ""); pkiFormat != "" {
			args = append(args, "--pki-format", pkiFormat)
		}
		return executeShipCommand(args)
	})

	// Rekor search transparency log tool
	searchTool := mcp.NewTool("rekor_search_log",
		mcp.WithDescription("Search Rekor transparency log using real rekor-cli"),
		mcp.WithString("artifact",
			mcp.Description("Path or URL to artifact file"),
		),
		mcp.WithString("public_key",
			mcp.Description("Path or URL to public key file"),
		),
		mcp.WithString("sha",
			mcp.Description("SHA512, SHA256, or SHA1 sum of artifact"),
		),
		mcp.WithString("email",
			mcp.Description("Email associated with public key"),
		),
		mcp.WithString("pki_format",
			mcp.Description("Format of public key (required when using public-key)"),
			mcp.Enum("pgp", "minisign", "x509", "ssh", "tuf"),
		),
		mcp.WithString("operator",
			mcp.Description("Search operator"),
			mcp.Enum("and", "or"),
		),
	)
	s.AddTool(searchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"rekor-cli", "search"}
		if artifact := request.GetString("artifact", ""); artifact != "" {
			args = append(args, "--artifact", artifact)
		}
		if publicKey := request.GetString("public_key", ""); publicKey != "" {
			args = append(args, "--public-key", publicKey)
		}
		if sha := request.GetString("sha", ""); sha != "" {
			args = append(args, "--sha", sha)
		}
		if email := request.GetString("email", ""); email != "" {
			args = append(args, "--email", email)
		}
		if pkiFormat := request.GetString("pki_format", ""); pkiFormat != "" {
			args = append(args, "--pki-format", pkiFormat)
		}
		if operator := request.GetString("operator", ""); operator != "" {
			args = append(args, "--operator", operator)
		}
		return executeShipCommand(args)
	})

	// Rekor get entry by UUID tool
	getByUuidTool := mcp.NewTool("rekor_get_by_uuid",
		mcp.WithDescription("Get entry from Rekor transparency log by UUID using real rekor-cli"),
		mcp.WithString("uuid",
			mcp.Description("UUID of entry to retrieve"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("", "tle"),
		),
	)
	s.AddTool(getByUuidTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		uuid := request.GetString("uuid", "")
		args := []string{"rekor-cli", "get", "--uuid", uuid}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Rekor get entry by log index tool
	getByIndexTool := mcp.NewTool("rekor_get_by_index",
		mcp.WithDescription("Get entry from Rekor transparency log by log index using real rekor-cli"),
		mcp.WithString("log_index",
			mcp.Description("Log index of entry to retrieve"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("", "tle"),
		),
	)
	s.AddTool(getByIndexTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logIndex := request.GetString("log_index", "")
		args := []string{"rekor-cli", "get", "--log-index", logIndex}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Rekor verify entry by UUID tool
	verifyByUuidTool := mcp.NewTool("rekor_verify_by_uuid",
		mcp.WithDescription("Verify entry in Rekor transparency log by UUID using real rekor-cli"),
		mcp.WithString("uuid",
			mcp.Description("UUID of entry to verify"),
			mcp.Required(),
		),
	)
	s.AddTool(verifyByUuidTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		uuid := request.GetString("uuid", "")
		args := []string{"rekor-cli", "verify", "--uuid", uuid}
		return executeShipCommand(args)
	})

	// Rekor verify entry by log index tool
	verifyByIndexTool := mcp.NewTool("rekor_verify_by_index",
		mcp.WithDescription("Verify entry in Rekor transparency log by log index using real rekor-cli"),
		mcp.WithString("log_index",
			mcp.Description("Log index of entry to verify"),
			mcp.Required(),
		),
	)
	s.AddTool(verifyByIndexTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logIndex := request.GetString("log_index", "")
		args := []string{"rekor-cli", "verify", "--log-index", logIndex}
		return executeShipCommand(args)
	})

	// Rekor verify artifact tool
	verifyArtifactTool := mcp.NewTool("rekor_verify_artifact",
		mcp.WithDescription("Verify artifact in Rekor transparency log using real rekor-cli"),
		mcp.WithString("artifact",
			mcp.Description("Path or URL to artifact file"),
			mcp.Required(),
		),
		mcp.WithString("signature",
			mcp.Description("Path or URL to signature file"),
		),
		mcp.WithString("public_key",
			mcp.Description("Path or URL to public key file"),
		),
		mcp.WithString("type",
			mcp.Description("Type of entry to verify"),
		),
		mcp.WithString("pki_format",
			mcp.Description("Format of signature/public key"),
			mcp.Enum("pgp", "minisign", "x509", "ssh", "tuf"),
		),
	)
	s.AddTool(verifyArtifactTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"rekor-cli", "verify"}
		if artifact := request.GetString("artifact", ""); artifact != "" {
			args = append(args, "--artifact", artifact)
		}
		if signature := request.GetString("signature", ""); signature != "" {
			args = append(args, "--signature", signature)
		}
		if publicKey := request.GetString("public_key", ""); publicKey != "" {
			args = append(args, "--public-key", publicKey)
		}
		if entryType := request.GetString("type", ""); entryType != "" {
			args = append(args, "--type", entryType)
		}
		if pkiFormat := request.GetString("pki_format", ""); pkiFormat != "" {
			args = append(args, "--pki-format", pkiFormat)
		}
		return executeShipCommand(args)
	})
}