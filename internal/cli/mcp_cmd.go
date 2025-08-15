// MCP Server implementation for Ship CLI

package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"github.com/cloudshipai/ship/internal/telemetry"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp [tool]",
	Short: "Start MCP server for a specific tool or all tools",
	Long: `Start an MCP server that exposes specific Ship CLI tools for AI assistants.

Available tools:
  # Terraform Tools
  lint         - TFLint for syntax and best practices
  checkov      - Checkov security scanning
  trivy        - Trivy security scanning
  cost         - OpenInfraQuote cost analysis
  docs         - terraform-docs documentation
  diagram      - InfraMap diagram generation
  
  # Security Tools
  gitleaks     - Secret detection with Gitleaks
  grype        - Vulnerability scanning with Grype
  syft         - SBOM generation with Syft
  prowler      - Multi-cloud security assessment
  trufflehog   - Verified secret detection
  cosign       - Container signing and verification
  actionlint   - GitHub Actions workflow linting
  semgrep      - Static analysis security scanning
  hadolint     - Dockerfile security linting
  cfn-nag      - CloudFormation security scanning
  conftest     - OPA policy testing
  git-secrets  - Git secrets scanning
  kube-bench   - Kubernetes security benchmarks
  kube-hunter  - Kubernetes penetration testing
  zap          - Web application security testing
  falco        - Runtime security monitoring
  nikto        - Web server security scanning
  openscap     - Security compliance scanning
  ossf-scorecard - Open Source Security Foundation scorecard
  scout-suite  - Multi-cloud security auditing
  steampipe    - Cloud infrastructure queries
  powerpipe    - Infrastructure benchmarking
  velero       - Kubernetes backup and disaster recovery
  goldilocks   - Kubernetes resource recommendations
  allstar      - GitHub security policy enforcement
  rekor        - Software supply chain transparency
  osv-scanner  - Open Source Vulnerability scanning
  license-detector - Software license detection
  registry     - Container registry operations
  cosign-golden - Enhanced Cosign for golden images
  history-scrub - Git history cleaning and secret removal
  iac-plan     - Infrastructure as Code planning
  slsa-verifier - SLSA provenance verification
  in-toto      - Supply chain attestation
  gatekeeper   - OPA Gatekeeper policy validation
  kubescape    - Kubernetes security scanning
  dockle       - Container image linting
  sops         - Secrets management with Mozilla SOPS
  
  # Cloud & Infrastructure Tools
  cloudquery   - Cloud asset inventory
  custodian    - Cloud governance engine
  terraformer  - Infrastructure import and management
  infracost    - Infrastructure cost estimation
  inframap     - Infrastructure visualization
  infrascan    - Infrastructure security scanning
  aws-iam-rotation - AWS IAM credential rotation
  tfstate-reader - Terraform state analysis
  packer       - Machine image building
  fleet        - GitOps for Kubernetes
  kuttl        - Kubernetes testing framework
  litmus       - Chaos engineering for Kubernetes
  cert-manager - Certificate management
  step-ca      - Certificate authority operations
  check-ssl-cert - SSL certificate validation
  k8s-network-policy - Kubernetes network policy management
  kyverno      - Kubernetes policy management
  kyverno-multitenant - Multi-tenant Kyverno policies
  github-admin - GitHub administration tools
  github-packages - GitHub Packages security
  trivy-golden - Enhanced Trivy for golden images
  
  # AWS IAM Tools
  cloudsplaining - AWS IAM security assessment
  parliament     - AWS IAM policy linting
  pmapper        - AWS IAM privilege mapping
  policy-sentry  - AWS IAM policy generation
  
  # Collections
  terraform  - All Terraform tools
  security   - All security tools
  aws-iam    - All AWS IAM tools
  cloud      - All cloud infrastructure tools
  kubernetes - All Kubernetes tools
  all        - All tools (default if no tool specified)

Examples:
  ship mcp lint      # MCP server for just TFLint
  ship mcp checkov   # MCP server for just Checkov
  ship mcp all       # MCP server for all tools`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMCPServer,
}

func init() {
	rootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().Int("port", 0, "Port to listen on (0 for stdio)")
	mcpCmd.Flags().String("host", "localhost", "Host to bind to")
	mcpCmd.Flags().Bool("stdio", true, "Use stdio transport (default)")
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	host, _ := cmd.Flags().GetString("host")
	useStdio, _ := cmd.Flags().GetBool("stdio")

	// Determine which tool to serve
	toolName := "all"
	if len(args) > 0 {
		toolName = args[0]
	}

	// Track MCP command usage
	telemetry.TrackMCPCommand(toolName)

	// Create MCP server
	serverName := fmt.Sprintf("ship-%s", toolName)
	s := server.NewMCPServer(serverName, "1.0.0")

	// Add specific tools based on argument
	switch toolName {
	// Terraform Tools
	case "lint":
		addLintTool(s)
	case "checkov":
		addCheckovTool(s)
	case "trivy":
		addTrivyTool(s)
	case "cost":
		addCostTool(s)
	case "docs":
		addDocsTool(s)
	case "diagram":
		addDiagramTool(s)

	// Security Tools
	case "gitleaks":
		addGitleaksTool(s)
	case "grype":
		addGrypeTool(s)
	case "syft":
		addSyftTool(s)
	case "prowler":
		addProwlerTool(s)
	case "trufflehog":
		addTruffleHogTool(s)
	case "cosign":
		addCosignTool(s)
	case "actionlint":
		addActionlintTool(s)
	case "semgrep":
		addSemgrepTool(s)
	case "hadolint":
		addHadolintTool(s)
	case "cfn-nag":
		addCfnNagTool(s)
	case "conftest":
		addConftestTool(s)
	case "git-secrets":
		addGitSecretsTool(s)
	case "kube-bench":
		addKubeBenchTool(s)
	case "kube-hunter":
		addKubeHunterTool(s)
	case "zap":
		addZapTool(s)
	case "falco":
		addFalcoTool(s)
	case "nikto":
		addNiktoTool(s)
	case "openscap":
		addOpenSCAPTool(s)
	case "ossf-scorecard":
		addOSSFScorecardTool(s)
	case "scout-suite":
		addScoutSuiteTool(s)
	case "steampipe":
		addSteampipeTool(s)
	case "powerpipe":
		addPowerpipeTool(s)
	case "velero":
		addVeleroTool(s)
	case "goldilocks":
		addGoldilocksTool(s)
	case "allstar":
		addAllstarTool(s)
	case "rekor":
		addRekorTool(s)
	case "osv-scanner":
		addOSVScannerTool(s)
	case "license-detector":
		addLicenseDetectorTool(s)
	case "registry":
		addRegistryTool(s)
	case "cosign-golden":
		addCosignGoldenTool(s)
	case "history-scrub":
		addHistoryScrubTool(s)
	case "iac-plan":
		addIacPlanTool(s)
	case "slsa-verifier":
		addSLSAVerifierTool(s)
	case "in-toto":
		addInTotoTool(s)
	case "gatekeeper":
		addGatekeeperTool(s)
	case "kubescape":
		addKubescapeTool(s)
	case "dockle":
		addDockleTool(s)
	case "sops":
		addSOPSTool(s)

	// Cloud & Infrastructure Tools
	case "cloudquery":
		addCloudQueryTool(s)
	case "custodian":
		addCustodianTool(s)
	case "terraformer":
		addTerraformerTool(s)
	case "infracost":
		addInfracostTool(s)
	case "inframap":
		addInframapTool(s)
	case "infrascan":
		addInfrascanTool(s)
	case "aws-iam-rotation":
		addAWSIAMRotationTool(s)
	case "tfstate-reader":
		addTfstateReaderTool(s)
	case "packer":
		addPackerTool(s)
	case "fleet":
		addFleetTool(s)
	case "kuttl":
		addKuttlTool(s)
	case "litmus":
		addLitmusTool(s)
	case "cert-manager":
		addCertManagerTool(s)
	case "step-ca":
		addStepCATool(s)
	case "check-ssl-cert":
		addCheckSSLCertTool(s)
	case "k8s-network-policy":
		addK8sNetworkPolicyTool(s)
	case "kyverno":
		addKyvernoTool(s)
	case "kyverno-multitenant":
		addKyvernoMultitenantTool(s)
	case "github-admin":
		addGitHubAdminTool(s)
	case "github-packages":
		addGitHubPackagesTool(s)
	case "trivy-golden":
		addTrivyGoldenTool(s)

	// AWS IAM Tools
	case "cloudsplaining":
		addCloudsplainingTool(s)
	case "parliament":
		addParliamentTool(s)
	case "pmapper":
		addPMapperTool(s)
	case "policy-sentry":
		addPolicySentryTool(s)

	// Collections
	case "terraform":
		addTerraformTools(s)
	case "security":
		addSecurityTools(s)
	case "aws-iam":
		addAWSIAMTools(s)
	case "cloud":
		addCloudTools(s)
	case "kubernetes":
		addKubernetesTools(s)
	case "all":
		addTerraformTools(s)
		addSecurityTools(s)
		addAWSIAMTools(s)
		addCloudTools(s)
		addKubernetesTools(s)
	default:
		return fmt.Errorf("unknown tool: %s. Use 'ship mcp --help' to see available tools", toolName)
	}

	// Add resources for documentation and help
	addResources(s)

	// Add prompts only for 'all' mode
	if toolName == "all" {
		addPrompts(s)
	}

	// Start server
	if useStdio || port == 0 {
		fmt.Fprintf(os.Stderr, "Starting %s MCP server on stdio...\n", serverName)
		return server.ServeStdio(s)
	} else {
		fmt.Fprintf(os.Stderr, "Starting %s MCP server on %s:%d...\n", serverName, host, port)
		return fmt.Errorf("HTTP server not implemented in this version, use --stdio")
	}
}

// Individual tool functions
func addLintTool(s *server.MCPServer) {
	lintTool := mcp.NewTool("lint",
		mcp.WithDescription("Run TFLint on Terraform code to check for syntax errors and best practices"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: default, json, compact"),
			mcp.Enum("default", "json", "compact"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save lint results"),
		),
	)

	s.AddTool(lintTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "lint"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addCheckovTool(s *server.MCPServer) {
	checkovTool := mcp.NewTool("checkov",
		mcp.WithDescription("Run Checkov security scan on Terraform code for policy compliance"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: cli, json, junit, sarif"),
			mcp.Enum("cli", "json", "junit", "sarif"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save scan results"),
		),
	)

	s.AddTool(checkovTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addTrivyTool(s *server.MCPServer) {
	trivyTool := mcp.NewTool("trivy",
		mcp.WithDescription("Run Trivy security scan on Terraform code using Trivy"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
	)

	s.AddTool(trivyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "trivy"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}

		return executeShipCommand(args)
	})
}

func addCostTool(s *server.MCPServer) {
	costTool := mcp.NewTool("cost",
		mcp.WithDescription("Analyze infrastructure costs using OpenInfraQuote"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing (e.g., us-east-1, us-west-2)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: json, table"),
			mcp.Enum("json", "table"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save cost analysis"),
		),
	)

	s.AddTool(costTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "cost"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addDocsTool(s *server.MCPServer) {
	docsTool := mcp.NewTool("docs",
		mcp.WithDescription("Generate documentation for Terraform modules using terraform-docs"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("filename",
			mcp.Description("Filename to save documentation as (default README.md)"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save documentation"),
		),
	)

	s.AddTool(docsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if filename := request.GetString("filename", ""); filename != "" {
			args = append(args, "--filename", filename)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addDiagramTool(s *server.MCPServer) {
	diagramTool := mcp.NewTool("diagram",
		mcp.WithDescription("Generate infrastructure diagrams from Terraform state"),
		mcp.WithString("input",
			mcp.Description("Input directory or file containing Terraform files (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: png, svg, pdf, dot"),
			mcp.Enum("png", "svg", "pdf", "dot"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save diagram"),
		),
		mcp.WithBoolean("hcl",
			mcp.Description("Generate from HCL files instead of state file"),
		),
		mcp.WithString("provider",
			mcp.Description("Filter by specific provider (aws, google, azurerm)"),
			mcp.Enum("aws", "google", "azurerm"),
		),
	)

	s.AddTool(diagramTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "diagram"}

		if input := request.GetString("input", ""); input != "" {
			args = append(args, input)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if hcl := request.GetBool("hcl", false); hcl {
			args = append(args, "--hcl")
		}
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, "--provider", provider)
		}

		return executeShipCommand(args)
	})
}

func addTerraformTools(s *server.MCPServer) {
	// Terraform Lint Tool
	lintTool := mcp.NewTool("lint",
		mcp.WithDescription("Run TFLint on Terraform code to check for syntax errors and best practices"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: default, json, compact"),
			mcp.Enum("default", "json", "compact"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save lint results"),
		),
	)

	s.AddTool(lintTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "lint"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})

	// Terraform Security Scan Tool (Checkov)
	checkovTool := mcp.NewTool("checkov",
		mcp.WithDescription("Run Checkov security scan on Terraform code for policy compliance"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: cli, json, junit, sarif"),
			mcp.Enum("cli", "json", "junit", "sarif"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save scan results"),
		),
	)

	s.AddTool(checkovTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "checkov"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})

	// Terraform Security Scan Tool (Alternative)
	securityTool := mcp.NewTool("trivy",
		mcp.WithDescription("Run alternative security scan on Terraform code using Trivy"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
	)

	s.AddTool(securityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "trivy"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		// Push functionality disabled during staging
		// if push := request.GetBool("push", false); push {
		//	args = append(args, "--push")
		// }

		return executeShipCommand(args)
	})

	// Terraform Cost Analysis Tool
	costAnalysisTool := mcp.NewTool("cost",
		mcp.WithDescription("Analyze infrastructure costs using OpenInfraQuote"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for pricing (e.g., us-east-1, us-west-2)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: json, table"),
			mcp.Enum("json", "table"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save cost analysis"),
		),
	)

	s.AddTool(costAnalysisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "cost"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})

	// Terraform Documentation Tool
	docsTool := mcp.NewTool("docs",
		mcp.WithDescription("Generate documentation for Terraform modules using terraform-docs"),
		mcp.WithString("directory",
			mcp.Description("Directory containing Terraform files (default: current directory)"),
		),
		mcp.WithString("filename",
			mcp.Description("Filename to save documentation as (default README.md)"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save documentation"),
		),
	)

	s.AddTool(docsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "docs"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if filename := request.GetString("filename", ""); filename != "" {
			args = append(args, "--filename", filename)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})

	// Terraform Diagram Generation Tool
	diagramTool := mcp.NewTool("diagram",
		mcp.WithDescription("Generate infrastructure diagrams from Terraform state"),
		mcp.WithString("input",
			mcp.Description("Input directory or file containing Terraform files (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: png, svg, pdf, dot"),
			mcp.Enum("png", "svg", "pdf", "dot"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save diagram"),
		),
		mcp.WithBoolean("hcl",
			mcp.Description("Generate from HCL files instead of state file"),
		),
		mcp.WithString("provider",
			mcp.Description("Filter by specific provider (aws, google, azurerm)"),
			mcp.Enum("aws", "google", "azurerm"),
		),
	)

	s.AddTool(diagramTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tf", "diagram"}

		if input := request.GetString("input", ""); input != "" {
			args = append(args, input)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		if hcl := request.GetBool("hcl", false); hcl {
			args = append(args, "--hcl")
		}
		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, "--provider", provider)
		}

		return executeShipCommand(args)
	})
}

// Investigation tools removed to focus on Terraform analysis workflows

func addResources(s *server.MCPServer) {
	// Help resource
	helpResource := mcp.NewResource("ship://help",
		"Ship CLI Help",
		mcp.WithResourceDescription("Complete help and usage information for Ship CLI"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(helpResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		result, err := executeShipCommand([]string{"--help"})
		if err != nil {
			return nil, err
		}

		// Extract text from result - the result should be a simple text response
		var helpText string
		if result != nil && len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(mcp.TextContent); ok {
				helpText = textContent.Text
			}
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "ship://help",
				MIMEType: "text/markdown",
				Text:     helpText,
			},
		}, nil
	})

	// Available tools resource
	toolsResource := mcp.NewResource("ship://tools",
		"Available Ship CLI Tools",
		mcp.WithResourceDescription("List of all available Ship CLI tools and their capabilities"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(toolsResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content := `# Ship CLI Tools

## Terraform Tools
- **lint**: Run TFLint for syntax and best practices
- **checkov**: Security scanning with Checkov
- **trivy**: Alternative security scanning with Trivy
- **cost**: Cost analysis with OpenInfraQuote
- **docs**: Generate documentation with terraform-docs
- **diagram**: Generate infrastructure diagrams with InfraMap


## Examples
- ` + "`ship tf lint`" + ` - Lint current directory
- ` + "`ship tf diagram . --hcl --format png`" + ` - Generate infrastructure diagram
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "ship://tools",
				MIMEType: "text/markdown",
				Text:     content,
			},
		}, nil
	})
}

func addPrompts(s *server.MCPServer) {
	// Security audit prompt
	securityPrompt := mcp.NewPrompt("security_audit",
		mcp.WithPromptDescription("Comprehensive security audit of cloud infrastructure"),
		mcp.WithArgument("provider",
			mcp.ArgumentDescription("Cloud provider to audit (aws, azure, gcp)"),
		),
	)

	s.AddPrompt(securityPrompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "Comprehensive security audit workflow",
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: `Please perform a comprehensive security audit of my Terraform infrastructure. Follow these steps:

1. Run terraform_checkov_scan to identify security issues in infrastructure-as-code
2. Run terraform_security_scan for additional security analysis  
3. Use terraform_lint to check for configuration best practices

4. Summarize all findings with:
   - Critical security issues requiring immediate attention
   - Recommendations for improvement
   - Best practices to implement

Please be thorough and provide actionable recommendations.`,
					},
				},
			},
		}, nil
	})

	// Cost optimization prompt
	costPrompt := mcp.NewPrompt("cost_optimization",
		mcp.WithPromptDescription("Identify cost optimization opportunities"),
		mcp.WithArgument("provider",
			mcp.ArgumentDescription("Cloud provider to analyze (aws, azure, gcp)"),
		),
	)

	s.AddPrompt(costPrompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "Cost optimization analysis workflow",
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: `Help me optimize costs for my Terraform infrastructure:

1. Use terraform_cost_analysis to analyze current cost projections
2. Review Terraform configurations for cost optimization opportunities
3. Use terraform_lint to identify inefficient resource configurations

4. Provide a prioritized list of cost-saving recommendations:
   - Quick wins (resource rightsizing, unused resources)
   - Medium-term optimizations (reserved instances, storage classes)
   - Long-term architectural improvements

Include estimated cost savings where possible.`,
					},
				},
			},
		}, nil
	})
}

// Maximum tokens allowed in MCP response (conservative estimate)
const maxMCPTokens = 20000

// Rough estimation: 1 token ≈ 4 characters for typical text
const charsPerToken = 4

// Security Tools
func addSecurityTools(s *server.MCPServer) {
	addGitleaksTool(s)
	addGrypeTool(s)
	addSyftTool(s)
	addProwlerTool(s)
	addTruffleHogTool(s)
	addCosignTool(s)
	addActionlintTool(s)
	addSemgrepTool(s)
	addHadolintTool(s)
	addCfnNagTool(s)
	addConftestTool(s)
	addGitSecretsTool(s)
	addKubeBenchTool(s)
	addKubeHunterTool(s)
	addZapTool(s)
	addFalcoTool(s)
	addNiktoTool(s)
	addOpenSCAPTool(s)
	addOSSFScorecardTool(s)
	addScoutSuiteTool(s)
	addAllstarTool(s)
	addRekorTool(s)
	addOSVScannerTool(s)
	addLicenseDetectorTool(s)
	addRegistryTool(s)
	addCosignGoldenTool(s)
	addHistoryScrubTool(s)
	addTrivyGoldenTool(s)
	addSLSAVerifierTool(s)
	addInTotoTool(s)
	addGatekeeperTool(s)
	addKubescapeTool(s)
	addDockleTool(s)
	addSOPSTool(s)
}

func addGitleaksTool(s *server.MCPServer) {
	gitleaksTool := mcp.NewTool("gitleaks",
		mcp.WithDescription("Scan for secrets using Gitleaks"),
		mcp.WithString("directory",
			mcp.Description("Directory to scan (default: current directory)"),
		),
		mcp.WithString("config",
			mcp.Description("Path to Gitleaks configuration file"),
		),
		mcp.WithBoolean("git",
			mcp.Description("Scan git repository history"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save results"),
		),
	)

	s.AddTool(gitleaksTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gitleaks"}

		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if git := request.GetBool("git", false); git {
			args = append(args, "--git")
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addGrypeTool(s *server.MCPServer) {
	grypeTool := mcp.NewTool("grype",
		mcp.WithDescription("Scan for vulnerabilities using Grype"),
		mcp.WithString("target",
			mcp.Description("Target to scan (directory, image:name, or sbom file)"),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level"),
			mcp.Enum("negligible", "low", "medium", "high", "critical"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "table", "cyclonedx"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save results"),
		),
	)

	s.AddTool(grypeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "grype"}

		if target := request.GetString("target", ""); target != "" {
			args = append(args, target)
		}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addSyftTool(s *server.MCPServer) {
	syftTool := mcp.NewTool("syft",
		mcp.WithDescription("Generate SBOM using Syft"),
		mcp.WithString("target",
			mcp.Description("Target to scan (directory or image:name)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format"),
			mcp.Enum("json", "spdx-json", "cyclonedx-json", "table"),
		),
		mcp.WithString("package-type",
			mcp.Description("Package type filter"),
			mcp.Enum("npm", "python", "go", "java"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save SBOM"),
		),
	)

	s.AddTool(syftTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "syft"}

		if target := request.GetString("target", ""); target != "" {
			args = append(args, target)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if packageType := request.GetString("package-type", ""); packageType != "" {
			args = append(args, "--package-type", packageType)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addProwlerTool(s *server.MCPServer) {
	prowlerTool := mcp.NewTool("prowler",
		mcp.WithDescription("Multi-cloud security assessment using Prowler"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider to scan"),
			mcp.Enum("aws", "azure", "gcp"),
		),
		mcp.WithString("region",
			mcp.Description("AWS region for scanning"),
		),
		mcp.WithString("compliance",
			mcp.Description("Compliance framework"),
			mcp.Enum("cis", "pci", "gdpr", "hipaa"),
		),
		mcp.WithString("services",
			mcp.Description("Specific services to scan"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save results"),
		),
	)

	s.AddTool(prowlerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "prowler"}

		if provider := request.GetString("provider", ""); provider != "" {
			args = append(args, provider)
		}
		if region := request.GetString("region", ""); region != "" {
			args = append(args, "--region", region)
		}
		if compliance := request.GetString("compliance", ""); compliance != "" {
			args = append(args, "--compliance", compliance)
		}
		if services := request.GetString("services", ""); services != "" {
			args = append(args, "--services", services)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addTruffleHogTool(s *server.MCPServer) {
	truffleHogTool := mcp.NewTool("trufflehog",
		mcp.WithDescription("Verified secret detection using TruffleHog"),
		mcp.WithString("target",
			mcp.Description("Target to scan (directory, repository URL, etc.)"),
		),
		mcp.WithString("type",
			mcp.Description("Scan type"),
			mcp.Enum("filesystem", "git", "github", "docker", "s3"),
		),
		mcp.WithBoolean("verify",
			mcp.Description("Verify found secrets"),
		),
		mcp.WithString("token",
			mcp.Description("GitHub token for repository access"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save results"),
		),
	)

	s.AddTool(truffleHogTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "trufflehog"}

		if target := request.GetString("target", ""); target != "" {
			args = append(args, target)
		}
		if scanType := request.GetString("type", ""); scanType != "" {
			args = append(args, "--type", scanType)
		}
		if verify := request.GetBool("verify", false); verify {
			args = append(args, "--verify")
		}
		if token := request.GetString("token", ""); token != "" {
			args = append(args, "--token", token)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addCosignTool(s *server.MCPServer) {
	cosignTool := mcp.NewTool("cosign",
		mcp.WithDescription("Container signing and verification using Cosign"),
		mcp.WithString("command",
			mcp.Description("Cosign command to execute"),
			mcp.Enum("verify", "sign", "generate-key-pair"),
		),
		mcp.WithString("image",
			mcp.Description("Container image to sign/verify"),
		),
		mcp.WithString("key",
			mcp.Description("Path to public/private key file"),
		),
		mcp.WithBoolean("keyless",
			mcp.Description("Use keyless signing/verification"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save results"),
		),
	)

	s.AddTool(cosignTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		command := request.GetString("command", "verify")
		args := []string{"security", "cosign", command}

		if image := request.GetString("image", ""); image != "" {
			args = append(args, image)
		}
		if key := request.GetString("key", ""); key != "" {
			args = append(args, "--key", key)
		}
		if keyless := request.GetBool("keyless", false); keyless {
			args = append(args, "--keyless")
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

// AWS IAM Tools
func addAWSIAMTools(s *server.MCPServer) {
	addCloudsplainingTool(s)
	addParliamentTool(s)
	addPMapperTool(s)
	addPolicySentryTool(s)
}

func addCloudsplainingTool(s *server.MCPServer) {
	cloudsplainingTool := mcp.NewTool("cloudsplaining",
		mcp.WithDescription("AWS IAM security assessment using Cloudsplaining"),
		mcp.WithString("command",
			mcp.Description("Cloudsplaining command to execute"),
			mcp.Enum("scan-account", "scan-policy"),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("policy-file",
			mcp.Description("IAM policy file to scan (for scan-policy command)"),
		),
		mcp.WithString("minimize",
			mcp.Description("Statement ID to minimize"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save results"),
		),
	)

	s.AddTool(cloudsplainingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		command := request.GetString("command", "scan-account")
		args := []string{"security", "cloudsplaining", command}

		if policyFile := request.GetString("policy-file", ""); policyFile != "" && command == "scan-policy" {
			args = append(args, policyFile)
		}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if minimize := request.GetString("minimize", ""); minimize != "" {
			args = append(args, "--minimize", minimize)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addParliamentTool(s *server.MCPServer) {
	parliamentTool := mcp.NewTool("parliament",
		mcp.WithDescription("AWS IAM policy linting using Parliament"),
		mcp.WithString("policy-file",
			mcp.Description("IAM policy file to lint"),
		),
		mcp.WithBoolean("community",
			mcp.Description("Include community auditors"),
		),
		mcp.WithString("auditors",
			mcp.Description("Path to private auditors directory"),
		),
		mcp.WithString("severity",
			mcp.Description("Minimum severity level"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save results"),
		),
	)

	s.AddTool(parliamentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "parliament"}

		if policyFile := request.GetString("policy-file", ""); policyFile != "" {
			args = append(args, policyFile)
		}
		if community := request.GetBool("community", false); community {
			args = append(args, "--community")
		}
		if auditors := request.GetString("auditors", ""); auditors != "" {
			args = append(args, "--auditors", auditors)
		}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addPMapperTool(s *server.MCPServer) {
	pmapperTool := mcp.NewTool("pmapper",
		mcp.WithDescription("AWS IAM privilege mapping using PMapper"),
		mcp.WithString("command",
			mcp.Description("PMapper command to execute"),
			mcp.Enum("create-graph", "query", "privesc", "admin", "list", "visualize"),
		),
		mcp.WithString("profile",
			mcp.Description("AWS profile to use"),
		),
		mcp.WithString("principal",
			mcp.Description("Principal name for queries"),
		),
		mcp.WithString("action",
			mcp.Description("Action for access queries"),
		),
		mcp.WithString("resource",
			mcp.Description("Resource ARN for access queries"),
		),
		mcp.WithString("format",
			mcp.Description("Output format for visualization"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save results"),
		),
	)

	s.AddTool(pmapperTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		command := request.GetString("command", "create-graph")
		args := []string{"security", "pmapper", command}

		if principal := request.GetString("principal", ""); principal != "" {
			args = append(args, principal)
		}
		if action := request.GetString("action", ""); action != "" {
			args = append(args, action)
		}
		if resource := request.GetString("resource", ""); resource != "" {
			args = append(args, resource)
		}
		if profile := request.GetString("profile", ""); profile != "" {
			args = append(args, "--profile", profile)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func addPolicySentryTool(s *server.MCPServer) {
	policySentryTool := mcp.NewTool("policy-sentry",
		mcp.WithDescription("AWS IAM policy generation using Policy Sentry"),
		mcp.WithString("command",
			mcp.Description("Policy Sentry command to execute"),
			mcp.Enum("create-template", "write-policy", "query-actions", "query-conditions"),
		),
		mcp.WithString("template-type",
			mcp.Description("Template type for create-template"),
			mcp.Enum("crud", "actions"),
		),
		mcp.WithString("input",
			mcp.Description("Input YAML file for write-policy"),
		),
		mcp.WithString("service",
			mcp.Description("AWS service for query commands"),
		),
		mcp.WithString("output",
			mcp.Description("Output file to save results"),
		),
	)

	s.AddTool(policySentryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		command := request.GetString("command", "create-template")
		args := []string{"security", "policy-sentry", command}

		if templateType := request.GetString("template-type", ""); templateType != "" {
			args = append(args, "--template-type", templateType)
		}
		if input := request.GetString("input", ""); input != "" {
			args = append(args, "--input", input)
		}
		if service := request.GetString("service", ""); service != "" {
			args = append(args, "--service", service)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}

		return executeShipCommand(args)
	})
}

func executeShipCommand(args []string) (*mcp.CallToolResult, error) {
	// Get the current binary path
	executable, err := os.Executable()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get executable path: %v", err)), nil
	}

	// Execute the ship command
	cmd := exec.Command(executable, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Command failed: %s\n\nOutput:\n%s", err.Error(), string(output))), nil
	}

	outputStr := string(output)

	// Check if output needs to be chunked
	if needsChunking(outputStr) {
		return createChunkedResponse(outputStr), nil
	}

	return mcp.NewToolResultText(outputStr), nil
}

func needsChunking(text string) bool {
	return utf8.RuneCountInString(text) > (maxMCPTokens * charsPerToken)
}

func createChunkedResponse(text string) *mcp.CallToolResult {
	maxChunkSize := maxMCPTokens * charsPerToken

	// Split text into chunks, preferring to break at newlines
	chunks := smartChunk(text, maxChunkSize)

	if len(chunks) <= 1 {
		return mcp.NewToolResultText(text)
	}

	// Create a summary response with information about chunking
	summary := fmt.Sprintf(`Output is large (%d characters, ~%d tokens) and has been summarized.

TOTAL CHUNKS: %d

SUMMARY OF FIRST CHUNK:
%s

--- [Content continues in additional chunks] ---

To see the full output, you can:
1. Run the command with smaller scope (specific directory/subset)
2. Use filtering options if available
3. Process chunks individually if needed

FIRST CHUNK PREVIEW (showing first %d characters):
%s`,
		utf8.RuneCountInString(text),
		utf8.RuneCountInString(text)/charsPerToken,
		len(chunks),
		getChunkSummary(chunks[0]),
		maxChunkSize/4, // Show 1/4 of max chunk size as preview
		truncateText(chunks[0], maxChunkSize/4),
	)

	return mcp.NewToolResultText(summary)
}

func smartChunk(text string, maxSize int) []string {
	if utf8.RuneCountInString(text) <= maxSize {
		return []string{text}
	}

	var chunks []string
	lines := strings.Split(text, "\n")

	currentChunk := ""
	currentSize := 0

	for _, line := range lines {
		lineSize := utf8.RuneCountInString(line) + 1 // +1 for newline

		// If adding this line would exceed the chunk size, start a new chunk
		if currentSize+lineSize > maxSize && currentChunk != "" {
			chunks = append(chunks, strings.TrimSuffix(currentChunk, "\n"))
			currentChunk = line + "\n"
			currentSize = lineSize
		} else {
			currentChunk += line + "\n"
			currentSize += lineSize
		}
	}

	// Add the last chunk if it has content
	if currentChunk != "" {
		chunks = append(chunks, strings.TrimSuffix(currentChunk, "\n"))
	}

	return chunks
}

func getChunkSummary(chunk string) string {
	lines := strings.Split(chunk, "\n")
	if len(lines) == 0 {
		return "Empty content"
	}

	// Try to identify the type of content
	firstLine := strings.TrimSpace(lines[0])

	if strings.Contains(chunk, "CRITICAL") || strings.Contains(chunk, "HIGH") {
		return "Security scan results with findings"
	} else if strings.Contains(chunk, "resource \"") {
		return "Terraform configuration analysis"
	} else if strings.Contains(chunk, "$") && strings.Contains(chunk, "cost") {
		return "Cost analysis results"
	} else if strings.Contains(chunk, "Error:") || strings.Contains(chunk, "Warning:") {
		return "Tool output with errors/warnings"
	} else {
		return fmt.Sprintf("Tool output starting with: %s", truncateText(firstLine, 100))
	}
}

func truncateText(text string, maxLen int) string {
	if utf8.RuneCountInString(text) <= maxLen {
		return text
	}

	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}

	return string(runes[:maxLen]) + "..."
}

// Collection functions for new tool categories
func addCloudTools(s *server.MCPServer) {
	addCloudQueryTool(s)
	addCustodianTool(s)
	addTerraformerTool(s)
	addInfracostTool(s)
	addInframapTool(s)
	addInfrascanTool(s)
	addAWSIAMRotationTool(s)
	addTfstateReaderTool(s)
	addPackerTool(s)
	addSteampipeTool(s)
	addPowerpipeTool(s)
	addIacPlanTool(s)
}

func addKubernetesTools(s *server.MCPServer) {
	addKubeBenchTool(s)
	addKubeHunterTool(s)
	addFalcoTool(s)
	addFleetTool(s)
	addKuttlTool(s)
	addLitmusTool(s)
	addCertManagerTool(s)
	addK8sNetworkPolicyTool(s)
	addKyvernoTool(s)
	addKyvernoMultitenantTool(s)
	addVeleroTool(s)
	addGoldilocksTool(s)
}

// Individual MCP tool functions for all remaining tools
func addActionlintTool(s *server.MCPServer) {
	tool := mcp.NewTool("actionlint",
		mcp.WithDescription("Lint GitHub Actions workflows"),
		mcp.WithString("directory", mcp.Description("Directory to scan (default: current directory)")),
		mcp.WithString("output", mcp.Description("Output file to save results")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "actionlint"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

func addSemgrepTool(s *server.MCPServer) {
	tool := mcp.NewTool("semgrep",
		mcp.WithDescription("Static analysis security scanning"),
		mcp.WithString("directory", mcp.Description("Directory to scan")),
		mcp.WithString("config", mcp.Description("Semgrep configuration")),
		mcp.WithString("output", mcp.Description("Output file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "semgrep"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if config := request.GetString("config", ""); config != "" {
			args = append(args, "--config", config)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

func addHadolintTool(s *server.MCPServer) {
	tool := mcp.NewTool("hadolint",
		mcp.WithDescription("Dockerfile security linting"),
		mcp.WithString("dockerfile", mcp.Description("Dockerfile to lint")),
		mcp.WithString("output", mcp.Description("Output file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "hadolint"}
		if dockerfile := request.GetString("dockerfile", ""); dockerfile != "" {
			args = append(args, dockerfile)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

func addCfnNagTool(s *server.MCPServer) {
	tool := mcp.NewTool("cfn-nag",
		mcp.WithDescription("CloudFormation security scanning"),
		mcp.WithString("template", mcp.Description("CloudFormation template")),
		mcp.WithString("output", mcp.Description("Output file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "cfn-nag"}
		if template := request.GetString("template", ""); template != "" {
			args = append(args, template)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

func addConftestTool(s *server.MCPServer) {
	tool := mcp.NewTool("conftest",
		mcp.WithDescription("OPA policy testing"),
		mcp.WithString("directory", mcp.Description("Directory to test")),
		mcp.WithString("policy", mcp.Description("Policy directory")),
		mcp.WithString("output", mcp.Description("Output file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "conftest"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if policy := request.GetString("policy", ""); policy != "" {
			args = append(args, "--policy", policy)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

func addGitSecretsTool(s *server.MCPServer) {
	tool := mcp.NewTool("git-secrets",
		mcp.WithDescription("Git secrets scanning"),
		mcp.WithString("directory", mcp.Description("Directory to scan")),
		mcp.WithString("output", mcp.Description("Output file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "git-secrets"}
		if dir := request.GetString("directory", ""); dir != "" {
			args = append(args, dir)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

func addKubeBenchTool(s *server.MCPServer) {
	tool := mcp.NewTool("kube-bench",
		mcp.WithDescription("Kubernetes security benchmarks"),
		mcp.WithString("output", mcp.Description("Output file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kube-bench"}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

func addKubeHunterTool(s *server.MCPServer) {
	tool := mcp.NewTool("kube-hunter",
		mcp.WithDescription("Kubernetes penetration testing"),
		mcp.WithString("target", mcp.Description("Target to scan")),
		mcp.WithString("output", mcp.Description("Output file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kube-hunter"}
		if target := request.GetString("target", ""); target != "" {
			args = append(args, target)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

func addZapTool(s *server.MCPServer) {
	tool := mcp.NewTool("zap",
		mcp.WithDescription("Web application security testing"),
		mcp.WithString("target", mcp.Description("Target URL")),
		mcp.WithString("output", mcp.Description("Output file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "zap"}
		if target := request.GetString("target", ""); target != "" {
			args = append(args, target)
		}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

func addFalcoTool(s *server.MCPServer) {
	tool := mcp.NewTool("falco",
		mcp.WithDescription("Runtime security monitoring"),
		mcp.WithString("output", mcp.Description("Output file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "falco"}
		if output := request.GetString("output", ""); output != "" {
			args = append(args, "--output", output)
		}
		return executeShipCommand(args)
	})
}

// Remaining security tools with simplified definitions
func addNiktoTool(s *server.MCPServer) {
	tool := mcp.NewTool("nikto", mcp.WithDescription("Web server security scanning"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "nikto"})
	})
}

func addOpenSCAPTool(s *server.MCPServer) {
	tool := mcp.NewTool("openscap", mcp.WithDescription("Security compliance scanning"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "openscap"})
	})
}

func addOSSFScorecardTool(s *server.MCPServer) {
	tool := mcp.NewTool("ossf-scorecard", mcp.WithDescription("Open Source Security Foundation scorecard"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "ossf-scorecard"})
	})
}

func addScoutSuiteTool(s *server.MCPServer) {
	tool := mcp.NewTool("scout-suite", mcp.WithDescription("Multi-cloud security auditing"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "scout-suite"})
	})
}

func addSteampipeTool(s *server.MCPServer) {
	tool := mcp.NewTool("steampipe", mcp.WithDescription("Cloud infrastructure queries"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "steampipe"})
	})
}

func addPowerpipeTool(s *server.MCPServer) {
	tool := mcp.NewTool("powerpipe", mcp.WithDescription("Infrastructure benchmarking"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "powerpipe"})
	})
}

func addVeleroTool(s *server.MCPServer) {
	tool := mcp.NewTool("velero", mcp.WithDescription("Kubernetes backup and disaster recovery"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "velero"})
	})
}

func addGoldilocksTool(s *server.MCPServer) {
	tool := mcp.NewTool("goldilocks", mcp.WithDescription("Kubernetes resource recommendations"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "goldilocks"})
	})
}

func addAllstarTool(s *server.MCPServer) {
	tool := mcp.NewTool("allstar", mcp.WithDescription("GitHub security policy enforcement"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "allstar"})
	})
}

func addRekorTool(s *server.MCPServer) {
	tool := mcp.NewTool("rekor", mcp.WithDescription("Software supply chain transparency"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "rekor"})
	})
}

func addOSVScannerTool(s *server.MCPServer) {
	tool := mcp.NewTool("osv-scanner", mcp.WithDescription("Open Source Vulnerability scanning"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "osv-scanner"})
	})
}

func addLicenseDetectorTool(s *server.MCPServer) {
	tool := mcp.NewTool("license-detector", mcp.WithDescription("Software license detection"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "license-detector"})
	})
}

func addRegistryTool(s *server.MCPServer) {
	tool := mcp.NewTool("registry", mcp.WithDescription("Container registry operations"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "registry"})
	})
}

func addCosignGoldenTool(s *server.MCPServer) {
	tool := mcp.NewTool("cosign-golden", mcp.WithDescription("Enhanced Cosign for golden images"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "cosign-golden"})
	})
}

func addHistoryScrubTool(s *server.MCPServer) {
	tool := mcp.NewTool("history-scrub", mcp.WithDescription("Git history cleaning and secret removal"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "history-scrub"})
	})
}

func addTrivyGoldenTool(s *server.MCPServer) {
	tool := mcp.NewTool("trivy-golden", mcp.WithDescription("Enhanced Trivy for golden images"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "trivy-golden"})
	})
}

// Cloud & Infrastructure Tools
func addCloudQueryTool(s *server.MCPServer) {
	tool := mcp.NewTool("cloudquery", mcp.WithDescription("Cloud asset inventory"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "cloudquery"})
	})
}

func addCustodianTool(s *server.MCPServer) {
	tool := mcp.NewTool("custodian", mcp.WithDescription("Cloud governance engine"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "custodian"})
	})
}

func addTerraformerTool(s *server.MCPServer) {
	tool := mcp.NewTool("terraformer", mcp.WithDescription("Infrastructure import and management"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "terraformer"})
	})
}

func addInfracostTool(s *server.MCPServer) {
	tool := mcp.NewTool("infracost", mcp.WithDescription("Infrastructure cost estimation"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "infracost"})
	})
}

func addInframapTool(s *server.MCPServer) {
	tool := mcp.NewTool("inframap", mcp.WithDescription("Infrastructure visualization"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "inframap"})
	})
}

func addInfrascanTool(s *server.MCPServer) {
	tool := mcp.NewTool("infrascan", mcp.WithDescription("Infrastructure security scanning"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "infrascan"})
	})
}

func addAWSIAMRotationTool(s *server.MCPServer) {
	tool := mcp.NewTool("aws-iam-rotation", mcp.WithDescription("AWS IAM credential rotation"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "aws-iam-rotation"})
	})
}

func addTfstateReaderTool(s *server.MCPServer) {
	tool := mcp.NewTool("tfstate-reader", mcp.WithDescription("Terraform state analysis"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "tfstate-reader"})
	})
}

func addPackerTool(s *server.MCPServer) {
	tool := mcp.NewTool("packer", mcp.WithDescription("Machine image building"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "packer"})
	})
}

func addFleetTool(s *server.MCPServer) {
	tool := mcp.NewTool("fleet", mcp.WithDescription("GitOps for Kubernetes"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "fleet"})
	})
}

func addKuttlTool(s *server.MCPServer) {
	tool := mcp.NewTool("kuttl", mcp.WithDescription("Kubernetes testing framework"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "kuttl"})
	})
}

func addLitmusTool(s *server.MCPServer) {
	tool := mcp.NewTool("litmus", mcp.WithDescription("Chaos engineering for Kubernetes"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "litmus"})
	})
}

func addCertManagerTool(s *server.MCPServer) {
	tool := mcp.NewTool("cert-manager", mcp.WithDescription("Certificate management"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "cert-manager"})
	})
}

func addStepCATool(s *server.MCPServer) {
	tool := mcp.NewTool("step-ca", mcp.WithDescription("Certificate authority operations"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "step-ca"})
	})
}

func addCheckSSLCertTool(s *server.MCPServer) {
	tool := mcp.NewTool("check-ssl-cert", mcp.WithDescription("SSL certificate validation"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "check-ssl-cert"})
	})
}

func addK8sNetworkPolicyTool(s *server.MCPServer) {
	tool := mcp.NewTool("k8s-network-policy", mcp.WithDescription("Kubernetes network policy management"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "k8s-network-policy"})
	})
}

func addKyvernoTool(s *server.MCPServer) {
	tool := mcp.NewTool("kyverno", mcp.WithDescription("Kubernetes policy management"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "kyverno"})
	})
}

func addKyvernoMultitenantTool(s *server.MCPServer) {
	tool := mcp.NewTool("kyverno-multitenant", mcp.WithDescription("Multi-tenant Kyverno policies"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "kyverno-multitenant"})
	})
}

func addGitHubAdminTool(s *server.MCPServer) {
	tool := mcp.NewTool("github-admin", mcp.WithDescription("GitHub administration tools"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "github-admin"})
	})
}

func addGitHubPackagesTool(s *server.MCPServer) {
	tool := mcp.NewTool("github-packages", mcp.WithDescription("GitHub Packages security"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "github-packages"})
	})
}

func addIacPlanTool(s *server.MCPServer) {
	tool := mcp.NewTool("iac-plan", mcp.WithDescription("Infrastructure as Code planning"))
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return executeShipCommand([]string{"security", "iac-plan"})
	})
}

// New high-priority security tools
func addSLSAVerifierTool(s *server.MCPServer) {
	tool := mcp.NewTool("slsa-verifier",
		mcp.WithDescription("SLSA provenance verification for supply chain security"),
		mcp.WithString("command", mcp.Description("SLSA command"), mcp.Enum("verify-artifact", "verify-image", "generate-policy")),
		mcp.WithString("artifact", mcp.Description("Path to artifact file")),
		mcp.WithString("provenance", mcp.Description("Path to provenance file")),
		mcp.WithString("source-uri", mcp.Description("Source URI for verification")),
		mcp.WithString("builder-id", mcp.Description("Builder ID for verification")),
		mcp.WithBoolean("print-provenance", mcp.Description("Print provenance information")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "slsa-verifier"}
		if command := request.GetString("command", ""); command != "" {
			args = append(args, command)
		}
		if artifact := request.GetString("artifact", ""); artifact != "" {
			args = append(args, "--artifact", artifact)
		}
		if provenance := request.GetString("provenance", ""); provenance != "" {
			args = append(args, "--provenance", provenance)
		}
		if sourceURI := request.GetString("source-uri", ""); sourceURI != "" {
			args = append(args, "--source-uri", sourceURI)
		}
		if builderID := request.GetString("builder-id", ""); builderID != "" {
			args = append(args, "--builder-id", builderID)
		}
		if printProvenance := request.GetBool("print-provenance", false); printProvenance {
			args = append(args, "--print-provenance")
		}
		return executeShipCommand(args)
	})
}

func addInTotoTool(s *server.MCPServer) {
	tool := mcp.NewTool("in-toto",
		mcp.WithDescription("Supply chain attestation using in-toto"),
		mcp.WithString("command", mcp.Description("in-toto command"), mcp.Enum("run", "verify", "record", "generate-layout")),
		mcp.WithString("step-name", mcp.Description("Step name for attestation")),
		mcp.WithString("key", mcp.Description("Path to signing key")),
		mcp.WithString("layout", mcp.Description("Path to layout file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "in-toto"}
		if command := request.GetString("command", ""); command != "" {
			args = append(args, command)
		}
		if stepName := request.GetString("step-name", ""); stepName != "" {
			args = append(args, "--step-name", stepName)
		}
		if key := request.GetString("key", ""); key != "" {
			args = append(args, "--key", key)
		}
		if layout := request.GetString("layout", ""); layout != "" {
			args = append(args, "--layout", layout)
		}
		return executeShipCommand(args)
	})
}

func addGatekeeperTool(s *server.MCPServer) {
	tool := mcp.NewTool("gatekeeper",
		mcp.WithDescription("OPA Gatekeeper policy validation"),
		mcp.WithString("command", mcp.Description("Gatekeeper command"), mcp.Enum("validate", "test", "generate-template", "sync", "analyze")),
		mcp.WithString("constraints", mcp.Description("Path to constraints directory")),
		mcp.WithString("templates", mcp.Description("Path to constraint templates directory")),
		mcp.WithString("resources", mcp.Description("Path to resources directory")),
		mcp.WithString("format", mcp.Description("Output format")),
		mcp.WithBoolean("verbose", mcp.Description("Verbose output")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "gatekeeper"}
		if command := request.GetString("command", ""); command != "" {
			args = append(args, command)
		}
		if constraints := request.GetString("constraints", ""); constraints != "" {
			args = append(args, "--constraints", constraints)
		}
		if templates := request.GetString("templates", ""); templates != "" {
			args = append(args, "--templates", templates)
		}
		if resources := request.GetString("resources", ""); resources != "" {
			args = append(args, "--resources", resources)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if verbose := request.GetBool("verbose", false); verbose {
			args = append(args, "--verbose")
		}
		return executeShipCommand(args)
	})
}

func addKubescapeTool(s *server.MCPServer) {
	tool := mcp.NewTool("kubescape",
		mcp.WithDescription("Kubernetes security scanning using Kubescape"),
		mcp.WithString("command", mcp.Description("Kubescape command"), mcp.Enum("cluster", "manifests", "helm", "repo", "report")),
		mcp.WithString("framework", mcp.Description("Security framework")),
		mcp.WithString("format", mcp.Description("Output format")),
		mcp.WithString("severity", mcp.Description("Severity threshold")),
		mcp.WithString("namespace", mcp.Description("Namespace to scan")),
		mcp.WithString("kubeconfig", mcp.Description("Path to kubeconfig file")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "kubescape"}
		if command := request.GetString("command", ""); command != "" {
			args = append(args, command)
		}
		if framework := request.GetString("framework", ""); framework != "" {
			args = append(args, "--framework", framework)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if severity := request.GetString("severity", ""); severity != "" {
			args = append(args, "--severity", severity)
		}
		if namespace := request.GetString("namespace", ""); namespace != "" {
			args = append(args, "--namespace", namespace)
		}
		if kubeconfig := request.GetString("kubeconfig", ""); kubeconfig != "" {
			args = append(args, "--kubeconfig", kubeconfig)
		}
		return executeShipCommand(args)
	})
}

func addDockleTool(s *server.MCPServer) {
	tool := mcp.NewTool("dockle",
		mcp.WithDescription("Container image linting using Dockle"),
		mcp.WithString("image", mcp.Description("Container image to scan")),
		mcp.WithString("format", mcp.Description("Output format")),
		mcp.WithString("exit-level", mcp.Description("Exit level")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "dockle"}
		if image := request.GetString("image", ""); image != "" {
			args = append(args, image)
		}
		if format := request.GetString("format", ""); format != "" {
			args = append(args, "--format", format)
		}
		if exitLevel := request.GetString("exit-level", ""); exitLevel != "" {
			args = append(args, "--exit-level", exitLevel)
		}
		return executeShipCommand(args)
	})
}

func addSOPSTool(s *server.MCPServer) {
	tool := mcp.NewTool("sops",
		mcp.WithDescription("Secrets management using Mozilla SOPS"),
		mcp.WithString("command", mcp.Description("SOPS command"), mcp.Enum("encrypt", "decrypt", "rotate", "edit", "generate-config", "validate")),
		mcp.WithString("file", mcp.Description("File to operate on")),
		mcp.WithString("kms", mcp.Description("KMS ARN for encryption")),
		mcp.WithString("pgp", mcp.Description("PGP fingerprint for encryption")),
		mcp.WithString("age", mcp.Description("Age public key for encryption")),
		mcp.WithBoolean("in-place", mcp.Description("Edit file in place")),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"security", "sops"}
		if command := request.GetString("command", ""); command != "" {
			args = append(args, command)
		}
		if file := request.GetString("file", ""); file != "" {
			args = append(args, file)
		}
		if kms := request.GetString("kms", ""); kms != "" {
			args = append(args, "--kms", kms)
		}
		if pgp := request.GetString("pgp", ""); pgp != "" {
			args = append(args, "--pgp", pgp)
		}
		if age := request.GetString("age", ""); age != "" {
			args = append(args, "--age", age)
		}
		if inPlace := request.GetBool("in-place", false); inPlace {
			args = append(args, "--in-place")
		}
		return executeShipCommand(args)
	})
}
