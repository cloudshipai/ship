package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddRekorTools adds Rekor (transparency log) MCP tool implementations using direct Dagger calls
func AddRekorTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addRekorToolsDirect(s)
}

// addRekorToolsDirect adds Rekor tools using direct Dagger module calls
func addRekorToolsDirect(s *server.MCPServer) {
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
			mcp.Description("Path or URL to public key file - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("type",
			mcp.Description("Type of entry - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("pki_format",
			mcp.Description("Format of signature/public key - NOTE: not supported in current Dagger module"),
			mcp.Enum("pgp", "minisign", "x509", "ssh", "tuf"),
		),
	)
	s.AddTool(uploadTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewRekorModule(client)

		// Get parameters
		artifact := request.GetString("artifact", "")
		signature := request.GetString("signature", "")

		// Check for unsupported parameters
		if request.GetString("public_key", "") != "" || request.GetString("type", "") != "" || request.GetString("pki_format", "") != "" {
			return mcp.NewToolResultError("public_key, type, and pki_format options not supported in current Dagger module"), nil
		}

		// Upload artifact
		output, err := module.Upload(ctx, artifact, signature)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("rekor upload failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Rekor search transparency log tool
	searchTool := mcp.NewTool("rekor_search_log",
		mcp.WithDescription("Search Rekor transparency log using real rekor-cli"),
		mcp.WithString("artifact",
			mcp.Description("Path or URL to artifact file"),
		),
		mcp.WithString("public_key",
			mcp.Description("Path or URL to public key file - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("sha",
			mcp.Description("SHA512, SHA256, or SHA1 sum of artifact - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("email",
			mcp.Description("Email associated with public key - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("pki_format",
			mcp.Description("Format of public key (required when using public-key) - NOTE: not supported in current Dagger module"),
			mcp.Enum("pgp", "minisign", "x509", "ssh", "tuf"),
		),
		mcp.WithString("operator",
			mcp.Description("Search operator - NOTE: not supported in current Dagger module"),
			mcp.Enum("and", "or"),
		),
	)
	s.AddTool(searchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewRekorModule(client)

		// Get parameters
		query := request.GetString("artifact", "")
		if query == "" {
			return mcp.NewToolResultError("artifact parameter is required"), nil
		}

		// Check for unsupported parameters
		if request.GetString("public_key", "") != "" || request.GetString("sha", "") != "" ||
			request.GetString("email", "") != "" || request.GetString("pki_format", "") != "" ||
			request.GetString("operator", "") != "" {
			return mcp.NewToolResultError("only artifact search supported in current Dagger module"), nil
		}

		// Search transparency log
		output, err := module.Search(ctx, query)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("rekor search failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Rekor get entry by UUID tool
	getByUuidTool := mcp.NewTool("rekor_get_by_uuid",
		mcp.WithDescription("Get entry from Rekor transparency log by UUID using real rekor-cli"),
		mcp.WithString("uuid",
			mcp.Description("UUID of entry to retrieve"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format - NOTE: only json supported in current Dagger module"),
			mcp.Enum("", "tle"),
		),
	)
	s.AddTool(getByUuidTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewRekorModule(client)

		// Get parameters
		uuid := request.GetString("uuid", "")

		// Check for unsupported format
		if format := request.GetString("format", ""); format != "" && format != "json" {
			return mcp.NewToolResultError("only json format supported in current Dagger module"), nil
		}

		// Get entry by UUID
		output, err := module.GetByUUID(ctx, uuid)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("rekor get by UUID failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Rekor get entry by log index tool
	getByIndexTool := mcp.NewTool("rekor_get_by_index",
		mcp.WithDescription("Get entry from Rekor transparency log by log index using real rekor-cli"),
		mcp.WithString("log_index",
			mcp.Description("Log index of entry to retrieve"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format - NOTE: only json supported in current Dagger module"),
			mcp.Enum("", "tle"),
		),
	)
	s.AddTool(getByIndexTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewRekorModule(client)

		// Get parameters
		logIndex := request.GetString("log_index", "")

		// Check for unsupported format
		if format := request.GetString("format", ""); format != "" && format != "json" {
			return mcp.NewToolResultError("only json format supported in current Dagger module"), nil
		}

		// Get entry by log index
		output, err := module.Get(ctx, logIndex)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("rekor get by index failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewRekorModule(client)

		// Get parameters
		uuid := request.GetString("uuid", "")

		// Verify entry by UUID
		output, err := module.VerifyByUUID(ctx, uuid)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("rekor verify by UUID failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewRekorModule(client)

		// Get parameters
		logIndex := request.GetString("log_index", "")

		// Verify entry by log index
		output, err := module.VerifyByIndex(ctx, logIndex)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("rekor verify by index failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
			mcp.Required(),
		),
		mcp.WithString("public_key",
			mcp.Description("Path or URL to public key file - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("type",
			mcp.Description("Type of entry to verify - NOTE: not supported in current Dagger module"),
		),
		mcp.WithString("pki_format",
			mcp.Description("Format of signature/public key - NOTE: not supported in current Dagger module"),
			mcp.Enum("pgp", "minisign", "x509", "ssh", "tuf"),
		),
	)
	s.AddTool(verifyArtifactTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewRekorModule(client)

		// Get parameters
		artifact := request.GetString("artifact", "")
		signature := request.GetString("signature", "")

		// Check for unsupported parameters
		if request.GetString("public_key", "") != "" || request.GetString("type", "") != "" || request.GetString("pki_format", "") != "" {
			return mcp.NewToolResultError("public_key, type, and pki_format options not supported in current Dagger module"), nil
		}

		// Verify artifact
		output, err := module.Verify(ctx, artifact, signature)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("rekor verify artifact failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}