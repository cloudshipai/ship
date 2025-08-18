package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddSLSAVerifierTools adds SLSA Verifier MCP tool implementations using direct Dagger calls
func AddSLSAVerifierTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addSLSAVerifierToolsDirect(s)
}

// addSLSAVerifierToolsDirect adds SLSA Verifier tools using direct Dagger module calls
func addSLSAVerifierToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSLSAVerifierModule(client)

		// Get parameters
		artifact := request.GetString("artifact", "")
		provenancePath := request.GetString("provenance_path", "")
		sourceURI := request.GetString("source_uri", "")
		sourceTag := request.GetString("source_tag", "")
		sourceBranch := request.GetString("source_branch", "")
		builderID := request.GetString("builder_id", "")
		printProvenance := request.GetBool("print_provenance", false)

		// Verify artifact
		output, err := module.VerifyArtifact(ctx, artifact, provenancePath, sourceURI, sourceTag, sourceBranch, builderID, printProvenance)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SLSA verifier artifact verification failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSLSAVerifierModule(client)

		// Get parameters
		image := request.GetString("image", "")
		sourceURI := request.GetString("source_uri", "")
		sourceTag := request.GetString("source_tag", "")
		sourceBranch := request.GetString("source_branch", "")
		builderID := request.GetString("builder_id", "")
		printProvenance := request.GetBool("print_provenance", false)

		// Verify image
		output, err := module.VerifyImage(ctx, image, sourceURI, sourceTag, sourceBranch, builderID, printProvenance)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SLSA verifier image verification failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSLSAVerifierModule(client)

		// Get parameters
		packageTarball := request.GetString("package_tarball", "")
		attestationsPath := request.GetString("attestations_path", "")
		packageName := request.GetString("package_name", "")
		packageVersion := request.GetString("package_version", "")
		sourceURI := request.GetString("source_uri", "")
		printProvenance := request.GetBool("print_provenance", false)

		// Verify npm package
		output, err := module.VerifyNpmPackage(ctx, packageTarball, attestationsPath, packageName, packageVersion, sourceURI, printProvenance)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SLSA verifier npm package verification failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// SLSA Verifier version tool
	versionTool := mcp.NewTool("slsa_verifier_version",
		mcp.WithDescription("Get SLSA Verifier version information using real slsa-verifier CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewSLSAVerifierModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("SLSA verifier get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}