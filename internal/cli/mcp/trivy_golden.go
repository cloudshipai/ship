package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTrivyGoldenTools adds enhanced Trivy for golden images MCP tool implementations
func AddTrivyGoldenTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Trivy golden scan base image tool
	scanBaseImageTool := mcp.NewTool("trivy_golden_scan_base_image",
		mcp.WithDescription("Scan base image for golden image creation using enhanced Trivy"),
		mcp.WithString("base_image",
			mcp.Description("Base container image to scan"),
			mcp.Required(),
		),
		mcp.WithString("severity_threshold",
			mcp.Description("Severity threshold (LOW, MEDIUM, HIGH, CRITICAL)"),
		),
	)
	s.AddTool(scanBaseImageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		baseImage := request.GetString("base_image", "")
		args := []string{"security", "trivy-golden", "scan-base", baseImage}
		if severity := request.GetString("severity_threshold", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		return executeShipCommand(args)
	})

	// Trivy golden create policy tool
	createPolicyTool := mcp.NewTool("trivy_golden_create_policy",
		mcp.WithDescription("Create security policy for golden image compliance"),
		mcp.WithString("policy_name",
			mcp.Description("Name of the security policy"),
			mcp.Required(),
		),
		mcp.WithString("compliance_framework",
			mcp.Description("Compliance framework (cis, nist, pci-dss)"),
			mcp.Required(),
		),
	)
	s.AddTool(createPolicyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policyName := request.GetString("policy_name", "")
		framework := request.GetString("compliance_framework", "")
		args := []string{"security", "trivy-golden", "create-policy", policyName, "--framework", framework}
		return executeShipCommand(args)
	})

	// Trivy golden validate compliance tool
	validateComplianceTool := mcp.NewTool("trivy_golden_validate_compliance",
		mcp.WithDescription("Validate golden image compliance against security policies"),
		mcp.WithString("image_name",
			mcp.Description("Golden image to validate"),
			mcp.Required(),
		),
		mcp.WithString("policy_file",
			mcp.Description("Path to compliance policy file"),
			mcp.Required(),
		),
	)
	s.AddTool(validateComplianceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		policyFile := request.GetString("policy_file", "")
		args := []string{"security", "trivy-golden", "validate-compliance", imageName, "--policy", policyFile}
		return executeShipCommand(args)
	})

	// Trivy golden generate attestation tool
	generateAttestationTool := mcp.NewTool("trivy_golden_generate_attestation",
		mcp.WithDescription("Generate security attestation for golden image"),
		mcp.WithString("image_name",
			mcp.Description("Golden image to attest"),
			mcp.Required(),
		),
		mcp.WithString("attestation_format",
			mcp.Description("Attestation format (in-toto, slsa)"),
		),
	)
	s.AddTool(generateAttestationTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		args := []string{"security", "trivy-golden", "generate-attestation", imageName}
		if format := request.GetString("attestation_format", ""); format != "" {
			args = append(args, "--format", format)
		}
		return executeShipCommand(args)
	})

	// Trivy golden benchmark tool
	benchmarkTool := mcp.NewTool("trivy_golden_benchmark",
		mcp.WithDescription("Run security benchmark against golden image"),
		mcp.WithString("image_name",
			mcp.Description("Golden image to benchmark"),
			mcp.Required(),
		),
		mcp.WithString("benchmark_type",
			mcp.Description("Benchmark type (cis-docker, cis-kubernetes)"),
			mcp.Required(),
		),
	)
	s.AddTool(benchmarkTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		imageName := request.GetString("image_name", "")
		benchmarkType := request.GetString("benchmark_type", "")
		args := []string{"security", "trivy-golden", "benchmark", imageName, "--type", benchmarkType}
		return executeShipCommand(args)
	})

	// Trivy golden get version tool
	getVersionTool := mcp.NewTool("trivy_golden_get_version",
		mcp.WithDescription("Get enhanced Trivy golden image tool version information"),
	)
	s.AddTool(getVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "trivy-golden", "--version"}
		return executeShipCommand(args)
	})
}