package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddRekorTools adds Rekor (transparency log) MCP tool implementations
func AddRekorTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Rekor upload artifact tool
	uploadTool := mcp.NewTool("rekor_upload_artifact",
		mcp.WithDescription("Upload artifact to Rekor transparency log"),
		mcp.WithString("artifact_path",
			mcp.Description("Path to artifact file to upload"),
			mcp.Required(),
		),
		mcp.WithString("signature_path",
			mcp.Description("Path to signature file"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text"),
		),
	)
	s.AddTool(uploadTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		artifactPath := request.GetString("artifact_path", "")
		signaturePath := request.GetString("signature_path", "")
		args := []string{"security", "rekor", "--upload", artifactPath, "--signature", signaturePath}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Rekor search transparency log tool
	searchTool := mcp.NewTool("rekor_search_log",
		mcp.WithDescription("Search Rekor transparency log"),
		mcp.WithString("query",
			mcp.Description("Search query for artifacts"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text"),
		),
	)
	s.AddTool(searchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.GetString("query", "")
		args := []string{"security", "rekor", "--search", query}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Rekor get entry tool
	getEntryTool := mcp.NewTool("rekor_get_entry",
		mcp.WithDescription("Get entry from Rekor transparency log"),
		mcp.WithString("log_index",
			mcp.Description("Log index of entry to retrieve"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text"),
		),
	)
	s.AddTool(getEntryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logIndex := request.GetString("log_index", "")
		args := []string{"security", "rekor", "--get", logIndex}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Rekor verify entry tool
	verifyTool := mcp.NewTool("rekor_verify_entry",
		mcp.WithDescription("Verify entry in Rekor transparency log"),
		mcp.WithString("artifact_path",
			mcp.Description("Path to artifact file to verify"),
			mcp.Required(),
		),
		mcp.WithString("signature_path",
			mcp.Description("Path to signature file"),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "text"),
		),
	)
	s.AddTool(verifyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		artifactPath := request.GetString("artifact_path", "")
		signaturePath := request.GetString("signature_path", "")
		args := []string{"security", "rekor", "--verify", artifactPath, "--signature", signaturePath}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})
}