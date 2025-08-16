package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddSLSAVerifierTools adds SLSA Verifier (SLSA attestation verification) MCP tool implementations
func AddSLSAVerifierTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// SLSA Verifier verify artifact tool
	verifyArtifactTool := mcp.NewTool("slsa_verifier_verify_artifact",
		mcp.WithDescription("Verify SLSA provenance for artifact"),
		mcp.WithString("artifact_path",
			mcp.Description("Path to artifact to verify"),
			mcp.Required(),
		),
		mcp.WithString("provenance_path",
			mcp.Description("Path to SLSA provenance file"),
			mcp.Required(),
		),
		mcp.WithString("source_uri",
			mcp.Description("Expected source URI for verification"),
		),
	)
	s.AddTool(verifyArtifactTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		artifactPath := request.GetString("artifact_path", "")
		provenancePath := request.GetString("provenance_path", "")
		args := []string{"security", "slsa-verifier", "verify", artifactPath, "--provenance", provenancePath}
		if sourceURI := request.GetString("source_uri", ""); sourceURI != "" {
			args = append(args, "--source-uri", sourceURI)
		}
		return executeShipCommand(args)
	})

	// SLSA Verifier verify container tool
	verifyContainerTool := mcp.NewTool("slsa_verifier_verify_container",
		mcp.WithDescription("Verify SLSA provenance for container image"),
		mcp.WithString("image_name",
			mcp.Description("Container image name to verify"),
			mcp.Required(),
		),
		mcp.WithString("source_uri",
			mcp.Description("Expected source URI for verification"),
			mcp.Required(),
		),
		mcp.WithString("source_tag",
			mcp.Description("Expected source tag for verification"),
		),
	)
	s.AddTool(verifyContainerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		sourceURI := request.GetString("source_uri", "")
		args := []string{"security", "slsa-verifier", "verify-image", imageName, "--source-uri", sourceURI}
		if sourceTag := request.GetString("source_tag", ""); sourceTag != "" {
			args = append(args, "--source-tag", sourceTag)
		}
		return executeShipCommand(args)
	})

	// SLSA Verifier verify npm package tool
	verifyNpmTool := mcp.NewTool("slsa_verifier_verify_npm",
		mcp.WithDescription("Verify SLSA provenance for npm package"),
		mcp.WithString("package_name",
			mcp.Description("NPM package name to verify"),
			mcp.Required(),
		),
		mcp.WithString("package_version",
			mcp.Description("NPM package version to verify"),
			mcp.Required(),
		),
		mcp.WithString("source_uri",
			mcp.Description("Expected source URI for verification"),
		),
	)
	s.AddTool(verifyNpmTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName := request.GetString("package_name", "")
		packageVersion := request.GetString("package_version", "")
		args := []string{"security", "slsa-verifier", "verify-npm", packageName, "--version", packageVersion}
		if sourceURI := request.GetString("source_uri", ""); sourceURI != "" {
			args = append(args, "--source-uri", sourceURI)
		}
		return executeShipCommand(args)
	})

	// SLSA Verifier validate provenance tool
	validateProvenanceTool := mcp.NewTool("slsa_verifier_validate_provenance",
		mcp.WithDescription("Validate SLSA provenance format and content"),
		mcp.WithString("provenance_path",
			mcp.Description("Path to SLSA provenance file"),
			mcp.Required(),
		),
		mcp.WithString("schema_version",
			mcp.Description("Expected SLSA schema version"),
		),
	)
	s.AddTool(validateProvenanceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provenancePath := request.GetString("provenance_path", "")
		args := []string{"security", "slsa-verifier", "validate", provenancePath}
		if schemaVersion := request.GetString("schema_version", ""); schemaVersion != "" {
			args = append(args, "--schema-version", schemaVersion)
		}
		return executeShipCommand(args)
	})

	// SLSA Verifier check build tool
	checkBuildTool := mcp.NewTool("slsa_verifier_check_build",
		mcp.WithDescription("Check build requirements against SLSA levels"),
		mcp.WithString("provenance_path",
			mcp.Description("Path to SLSA provenance file"),
			mcp.Required(),
		),
		mcp.WithString("slsa_level",
			mcp.Description("Minimum SLSA level to check (1, 2, 3)"),
		),
	)
	s.AddTool(checkBuildTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provenancePath := request.GetString("provenance_path", "")
		args := []string{"security", "slsa-verifier", "check-build", provenancePath}
		if slsaLevel := request.GetString("slsa_level", ""); slsaLevel != "" {
			args = append(args, "--level", slsaLevel)
		}
		return executeShipCommand(args)
	})

	// SLSA Verifier generate policy tool
	generatePolicyTool := mcp.NewTool("slsa_verifier_generate_policy",
		mcp.WithDescription("Generate SLSA verification policy template"),
		mcp.WithString("policy_type",
			mcp.Description("Type of policy to generate (basic, strict, custom)"),
			mcp.Required(),
		),
		mcp.WithString("output_file",
			mcp.Description("Output file path for the policy"),
		),
	)
	s.AddTool(generatePolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyType := request.GetString("policy_type", "")
		args := []string{"security", "slsa-verifier", "generate-policy", "--type", policyType}
		if outputFile := request.GetString("output_file", ""); outputFile != "" {
			args = append(args, "--output", outputFile)
		}
		return executeShipCommand(args)
	})

	// SLSA Verifier get version tool
	getVersionTool := mcp.NewTool("slsa_verifier_get_version",
		mcp.WithDescription("Get SLSA Verifier version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "slsa-verifier", "--version"}
		return executeShipCommand(args)
	})
}