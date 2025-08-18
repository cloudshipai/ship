package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddInTotoTools adds in-toto (supply chain attestation) MCP tool implementations using real CLI tools
func AddInTotoTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// in-toto-run tool
	runStepTool := mcp.NewTool("in_toto_run_step",
		mcp.WithDescription("Run in-toto supply chain step with attestation using in-toto-run"),
		mcp.WithString("step_name",
			mcp.Description("Name for link metadata file"),
			mcp.Required(),
		),
		mcp.WithString("command",
			mcp.Description("Command to execute for this step"),
			mcp.Required(),
		),
		mcp.WithString("signing_key",
			mcp.Description("Path to signing key in PKCS8/PEM format"),
		),
		mcp.WithString("gpg_keyid",
			mcp.Description("GPG keyid to sign metadata"),
		),
		mcp.WithString("materials",
			mcp.Description("Paths to files/directories to record before command"),
		),
		mcp.WithString("products",
			mcp.Description("Paths to files/directories to record after command"),
		),
		mcp.WithBoolean("record_streams",
			mcp.Description("Capture stdout/stderr"),
		),
		mcp.WithBoolean("no_command",
			mcp.Description("Generate metadata without executing command"),
		),
	)
	s.AddTool(runStepTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stepName := request.GetString("step_name", "")
		command := request.GetString("command", "")
		args := []string{"in-toto-run", "-n", stepName}
		
		if signingKey := request.GetString("signing_key", ""); signingKey != "" {
			args = append(args, "--signing-key", signingKey)
		}
		if gpgKeyid := request.GetString("gpg_keyid", ""); gpgKeyid != "" {
			args = append(args, "-g", gpgKeyid)
		}
		if materials := request.GetString("materials", ""); materials != "" {
			args = append(args, "-m", materials)
		}
		if products := request.GetString("products", ""); products != "" {
			args = append(args, "-p", products)
		}
		if request.GetBool("record_streams", false) {
			args = append(args, "-s")
		}
		if request.GetBool("no_command", false) {
			args = append(args, "-x")
		} else {
			args = append(args, "--", command)
		}
		
		return executeShipCommand(args)
	})

	// in-toto-verify tool
	verifySupplyChainTool := mcp.NewTool("in_toto_verify",
		mcp.WithDescription("Verify in-toto supply chain integrity using in-toto-verify"),
		mcp.WithString("layout_path",
			mcp.Description("Path to supply chain layout file"),
			mcp.Required(),
		),
		mcp.WithString("layout_key_paths",
			mcp.Description("Comma-separated paths to layout verification keys"),
		),
		mcp.WithString("link_dir",
			mcp.Description("Directory containing link metadata"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose output"),
		),
	)
	s.AddTool(verifySupplyChainTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		layoutPath := request.GetString("layout_path", "")
		args := []string{"in-toto-verify", "-l", layoutPath}
		
		if layoutKeyPaths := request.GetString("layout_key_paths", ""); layoutKeyPaths != "" {
			args = append(args, "-k", layoutKeyPaths)
		}
		if linkDir := request.GetString("link_dir", ""); linkDir != "" {
			args = append(args, "-d", linkDir)
		}
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		
		return executeShipCommand(args)
	})

	// in-toto-record tool
	recordTool := mcp.NewTool("in_toto_record",
		mcp.WithDescription("Record in-toto step metadata using in-toto-record"),
		mcp.WithString("step_name",
			mcp.Description("Name for link metadata file"),
			mcp.Required(),
		),
		mcp.WithString("operation",
			mcp.Description("Record operation"),
			mcp.Required(),
			mcp.Enum("start", "stop"),
		),
		mcp.WithString("signing_key",
			mcp.Description("Path to signing key in PKCS8/PEM format"),
		),
		mcp.WithString("gpg_keyid",
			mcp.Description("GPG keyid to sign metadata"),
		),
		mcp.WithString("materials",
			mcp.Description("Paths to files/directories to record as materials"),
		),
		mcp.WithString("products",
			mcp.Description("Paths to files/directories to record as products"),
		),
	)
	s.AddTool(recordTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stepName := request.GetString("step_name", "")
		operation := request.GetString("operation", "")
		args := []string{"in-toto-record", operation, "-n", stepName}
		
		if signingKey := request.GetString("signing_key", ""); signingKey != "" {
			args = append(args, "--signing-key", signingKey)
		}
		if gpgKeyid := request.GetString("gpg_keyid", ""); gpgKeyid != "" {
			args = append(args, "-g", gpgKeyid)
		}
		if materials := request.GetString("materials", ""); materials != "" {
			args = append(args, "-m", materials)
		}
		if products := request.GetString("products", ""); products != "" {
			args = append(args, "-p", products)
		}
		
		return executeShipCommand(args)
	})

	// in-toto-sign tool
	signTool := mcp.NewTool("in_toto_sign",
		mcp.WithDescription("Sign in-toto metadata using in-toto-sign"),
		mcp.WithString("metadata_file",
			mcp.Description("Path to metadata file to sign"),
			mcp.Required(),
		),
		mcp.WithString("signing_key",
			mcp.Description("Path to signing key in PKCS8/PEM format"),
		),
		mcp.WithString("gpg_keyid",
			mcp.Description("GPG keyid to sign metadata"),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file for signed metadata"),
		),
	)
	s.AddTool(signTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		metadataFile := request.GetString("metadata_file", "")
		args := []string{"in-toto-sign", "-f", metadataFile}
		
		if signingKey := request.GetString("signing_key", ""); signingKey != "" {
			args = append(args, "--signing-key", signingKey)
		}
		if gpgKeyid := request.GetString("gpg_keyid", ""); gpgKeyid != "" {
			args = append(args, "-g", gpgKeyid)
		}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "-o", outputFile)
		}
		
		return executeShipCommand(args)
	})
}