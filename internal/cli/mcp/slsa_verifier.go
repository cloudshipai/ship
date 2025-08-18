package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSLSAVerifierTools adds SLSA Verifier MCP tool implementations using real slsa-verifier CLI commands
func AddSLSAVerifierTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// SLSA Verifier verify artifact tool
	verifyArtifactTool := mcp.NewTool("slsa_verifier_verify_artifact",
		mcp.WithDescription("Verify SLSA provenance for artifact using real slsa-verifier CLI"),
		mcp.WithString("artifact",
			mcp.Description("Path to artifact to verify"),
			mcp.Required(),
		),
		mcp.WithString("provenance_path",
			mcp.Description("Path to SLSA provenance file"),
			mcp.Required(),
		),
		mcp.WithString("source_uri",
			mcp.Description("Expected source URI (e.g., github.com/owner/repo)"),
			mcp.Required(),
		),
		mcp.WithString("source_tag",
			mcp.Description("Expected source tag for verification"),
		),
		mcp.WithString("source_branch",
			mcp.Description("Expected source branch for verification"),
		),
		mcp.WithString("builder_id",
			mcp.Description("Unique builder ID for verification"),
		),
		mcp.WithBoolean("print_provenance",
			mcp.Description("Output verified provenance"),
		),
	)
	s.AddTool(verifyArtifactTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		artifact := request.GetString("artifact", "")
		provenancePath := request.GetString("provenance_path", "")
		sourceURI := request.GetString("source_uri", "")
		args := []string{"slsa-verifier", "verify-artifact", artifact, "--provenance-path", provenancePath, "--source-uri", sourceURI}
		
		if sourceTag := request.GetString("source_tag", ""); sourceTag != "" {
			args = append(args, "--source-tag", sourceTag)
		}
		if sourceBranch := request.GetString("source_branch", ""); sourceBranch != "" {
			args = append(args, "--source-branch", sourceBranch)
		}
		if builderID := request.GetString("builder_id", ""); builderID != "" {
			args = append(args, "--builder-id", builderID)
		}
		if request.GetBool("print_provenance", false) {
			args = append(args, "--print-provenance")
		}
		
		return executeShipCommand(args)
	})

	// SLSA Verifier verify container image tool
	verifyImageTool := mcp.NewTool("slsa_verifier_verify_image",
		mcp.WithDescription("Verify SLSA provenance for container image using real slsa-verifier CLI"),
		mcp.WithString("image",
			mcp.Description("Container image digest to verify"),
			mcp.Required(),
		),
		mcp.WithString("source_uri",
			mcp.Description("Expected source URI (e.g., github.com/owner/repo)"),
			mcp.Required(),
		),
		mcp.WithString("source_tag",
			mcp.Description("Expected source tag for verification"),
		),
		mcp.WithString("source_branch",
			mcp.Description("Expected source branch for verification"),
		),
		mcp.WithString("builder_id",
			mcp.Description("Unique builder ID for verification"),
		),
		mcp.WithBoolean("print_provenance",
			mcp.Description("Output verified provenance"),
		),
	)
	s.AddTool(verifyImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		image := request.GetString("image", "")
		sourceURI := request.GetString("source_uri", "")
		args := []string{"slsa-verifier", "verify-image", image, "--source-uri", sourceURI}
		
		if sourceTag := request.GetString("source_tag", ""); sourceTag != "" {
			args = append(args, "--source-tag", sourceTag)
		}
		if sourceBranch := request.GetString("source_branch", ""); sourceBranch != "" {
			args = append(args, "--source-branch", sourceBranch)
		}
		if builderID := request.GetString("builder_id", ""); builderID != "" {
			args = append(args, "--builder-id", builderID)
		}
		if request.GetBool("print_provenance", false) {
			args = append(args, "--print-provenance")
		}
		
		return executeShipCommand(args)
	})

	// SLSA Verifier verify npm package tool (experimental)
	verifyNpmTool := mcp.NewTool("slsa_verifier_verify_npm_package",
		mcp.WithDescription("Verify SLSA provenance for npm package using real slsa-verifier CLI (experimental)"),
		mcp.WithString("package_tarball",
			mcp.Description("Path to npm package tarball to verify"),
			mcp.Required(),
		),
		mcp.WithString("attestations_path",
			mcp.Description("Path to attestations file"),
			mcp.Required(),
		),
		mcp.WithString("package_name",
			mcp.Description("NPM package name"),
			mcp.Required(),
		),
		mcp.WithString("package_version",
			mcp.Description("NPM package version"),
			mcp.Required(),
		),
		mcp.WithString("source_uri",
			mcp.Description("Expected source URI (e.g., github.com/owner/repo)"),
		),
		mcp.WithBoolean("print_provenance",
			mcp.Description("Output verified provenance"),
		),
	)
	s.AddTool(verifyNpmTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageTarball := request.GetString("package_tarball", "")
		attestationsPath := request.GetString("attestations_path", "")
		packageName := request.GetString("package_name", "")
		packageVersion := request.GetString("package_version", "")
		
		args := []string{"slsa-verifier", "verify-npm-package", packageTarball, "--attestations-path", attestationsPath, "--package-name", packageName, "--package-version", packageVersion}
		
		if sourceURI := request.GetString("source_uri", ""); sourceURI != "" {
			args = append(args, "--source-uri", sourceURI)
		}
		if request.GetBool("print_provenance", false) {
			args = append(args, "--print-provenance")
		}
		
		return executeShipCommand(args)
	})

	// SLSA Verifier version tool
	versionTool := mcp.NewTool("slsa_verifier_version",
		mcp.WithDescription("Get SLSA Verifier version information using real slsa-verifier CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"slsa-verifier", "version"}
		return executeShipCommand(args)
	})
}