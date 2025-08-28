package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

// ToolInfo contains information about a tool
type ToolInfo struct {
	Name         string
	Description  string
	Category     string
	AddFunc      func(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc)
	HasVariables bool // Indicates if the tool requires variables (like AWS credentials)
}

// ToolRegistry contains all available tools organized by category
var ToolRegistry = map[string][]ToolInfo{
	"security": {
		{Name: "trivy", Description: "Comprehensive vulnerability scanner", AddFunc: AddTrivyTools, HasVariables: false},
		{Name: "grype", Description: "Container vulnerability scanner", AddFunc: AddGrypeTools, HasVariables: false},
		{Name: "syft", Description: "SBOM generation tool", AddFunc: AddSyftTools, HasVariables: false},
		{Name: "checkov", Description: "Infrastructure as code static analysis", AddFunc: AddCheckovTools, HasVariables: false},
		{Name: "terrascan", Description: "IaC security scanner", AddFunc: AddTerrascanTools, HasVariables: false},
		{Name: "tfsec", Description: "Terraform-specific security scanner", AddFunc: AddTfsecTools, HasVariables: false},
		{Name: "semgrep", Description: "Static analysis for security", AddFunc: AddSemgrepTools, HasVariables: false},
		{Name: "actionlint", Description: "GitHub Actions workflow linter", AddFunc: AddActionlintTools, HasVariables: false},
		{Name: "hadolint", Description: "Dockerfile linter", AddFunc: AddHadolintTools, HasVariables: false},
		{Name: "conftest", Description: "OPA policy testing", AddFunc: AddConftestTools, HasVariables: false},
		{Name: "kube-bench", Description: "Kubernetes CIS benchmark", AddFunc: AddKubeBenchTools, HasVariables: false},
		{Name: "kube-hunter", Description: "Kubernetes penetration testing", AddFunc: AddKubeHunterTools, HasVariables: false},
		{Name: "falco", Description: "Runtime security monitoring", AddFunc: AddFalcoTools, HasVariables: false},
		{Name: "nuclei", Description: "Fast vulnerability scanner with community templates", AddFunc: AddNucleiTools, HasVariables: false},
		{Name: "zap", Description: "OWASP ZAP web application scanner", AddFunc: AddZapTools, HasVariables: false},
		{Name: "nmap", Description: "Network exploration and security auditing", AddFunc: AddNmapTools, HasVariables: false},
		{Name: "git-secrets", Description: "Git repository secret scanner", AddFunc: AddGitSecretsTools, HasVariables: false},
		{Name: "trufflehog", Description: "Advanced secret scanning with verification", AddFunc: AddTrufflehogTools, HasVariables: false},
		{Name: "gitleaks", Description: "Fast secret scanning for git repositories", AddFunc: AddGitleaksTools, HasVariables: false},
		{Name: "kubescape", Description: "Kubernetes security scanner", AddFunc: AddKubescapeTools, HasVariables: false},
		{Name: "dockle", Description: "Container image linter", AddFunc: AddDockleTools, HasVariables: false},
		{Name: "sops", Description: "Secrets management", AddFunc: AddSOPSTools, HasVariables: true},
		{Name: "ossf-scorecard", Description: "OSSF security scorecard", AddFunc: AddOSSFScorecardTools, HasVariables: false},
		{Name: "steampipe", Description: "Cloud asset querying with SQL", AddFunc: AddSteampipeTools, HasVariables: true},
		{Name: "cfn-nag", Description: "CloudFormation security linter", AddFunc: AddCfnNagTools, HasVariables: false},
		{Name: "gatekeeper", Description: "OPA Gatekeeper policy validation", AddFunc: AddGatekeeperTools, HasVariables: false},
		{Name: "license-detector", Description: "Software license detection", AddFunc: AddLicenseDetectorTools, HasVariables: false},
		{Name: "openscap", Description: "Security compliance scanning", AddFunc: AddOpenSCAPTools, HasVariables: false},
		{Name: "scout-suite", Description: "Multi-cloud security auditing", AddFunc: AddScoutSuiteTools, HasVariables: true},
		{Name: "powerpipe", Description: "Infrastructure benchmarking", AddFunc: AddPowerpipeTools, HasVariables: true},
		{Name: "infrascan", Description: "Infrastructure security scanning", AddFunc: AddInfrascanTools, HasVariables: true},
		{Name: "github-admin", Description: "GitHub administration tools", AddFunc: AddGitHubAdminTools, HasVariables: true},
		{Name: "github-packages", Description: "GitHub Packages security", AddFunc: AddGitHubPackagesTools, HasVariables: true},
	},
	"development": {
		{Name: "opencode", Description: "OpenCode AI coding assistant", AddFunc: AddOpenCodeTools, HasVariables: true},
	},
	"terraform": {
		{Name: "tflint", Description: "Terraform linter", AddFunc: AddTfLintTools, HasVariables: false},
		{Name: "terraform-docs", Description: "Terraform documentation generator", AddFunc: AddTerraformDocsTools, HasVariables: false},
		{Name: "inframap", Description: "Infrastructure visualization", AddFunc: AddInfraMapTools, HasVariables: false},
		{Name: "iac-plan", Description: "Infrastructure as code planning", AddFunc: AddIacPlanTools, HasVariables: false},
		{Name: "terraformer", Description: "Infrastructure import and management", AddFunc: AddTerraformerTools, HasVariables: true},
		{Name: "tfstate-reader", Description: "Terraform state analysis", AddFunc: AddTfstateReaderTools, HasVariables: false},
		{Name: "openinfraquote", Description: "Infrastructure cost estimation", AddFunc: AddOpenInfraQuoteTools, HasVariables: true},
	},
	"kubernetes": {
		{Name: "velero", Description: "Kubernetes backup and restore", AddFunc: AddVeleroTools, HasVariables: true},
		{Name: "goldilocks", Description: "Kubernetes resource recommendations", AddFunc: AddGoldilocksTools, HasVariables: false},
		{Name: "fleet", Description: "GitOps for Kubernetes", AddFunc: AddFleetTools, HasVariables: false},
		{Name: "kuttl", Description: "Kubernetes testing framework", AddFunc: AddKuttlTools, HasVariables: false},
		{Name: "litmus", Description: "Chaos engineering for Kubernetes", AddFunc: AddLitmusTools, HasVariables: false},
		{Name: "cert-manager", Description: "Certificate management", AddFunc: AddCertManagerTools, HasVariables: false},
		{Name: "k8s-network-policy", Description: "Kubernetes network policy management", AddFunc: AddK8sNetworkPolicyTools, HasVariables: false},
		{Name: "kyverno", Description: "Kubernetes policy management", AddFunc: AddKyvernoTools, HasVariables: false},
		{Name: "kyverno-multitenant", Description: "Multi-tenant Kyverno policies", AddFunc: AddKyvernoMultitenantTools, HasVariables: false},
	},
	"cloud": {
		{Name: "cloudquery", Description: "Cloud asset inventory", AddFunc: AddCloudQueryTools, HasVariables: true},
		{Name: "custodian", Description: "Cloud governance engine", AddFunc: AddCustodianTools, HasVariables: true},
		{Name: "packer", Description: "Machine image building", AddFunc: AddPackerTools, HasVariables: true},
	},
	"supply-chain": {
		{Name: "cosign", Description: "Container signing and verification", AddFunc: AddCosignTools, HasVariables: true},
		{Name: "dependency-track", Description: "OWASP Dependency-Track SBOM analysis", AddFunc: AddDependencyTrackTools, HasVariables: true},
	},
	"aws": {
		{Name: "cloudsplaining", Description: "AWS IAM policy scanner", AddFunc: AddCloudsplainingTools, HasVariables: true},
		{Name: "parliament", Description: "AWS IAM policy linter", AddFunc: AddParliamentTools, HasVariables: true},
		{Name: "pmapper", Description: "AWS IAM privilege escalation analysis", AddFunc: AddPMapperTools, HasVariables: true},
		{Name: "policy-sentry", Description: "AWS IAM policy generator", AddFunc: AddPolicySentryTools, HasVariables: true},
		{Name: "prowler", Description: "Multi-cloud security assessment", AddFunc: AddProwlerTools, HasVariables: true},
		{Name: "aws-iam-rotation", Description: "AWS IAM credential rotation", AddFunc: AddAWSIAMRotationTools, HasVariables: true},
		{Name: "aws-pricing", Description: "AWS pricing and cost calculator", AddFunc: AddAWSPricingTools, HasVariables: true},
	},
}

// RegisterAllTools registers all tools with the MCP server
func RegisterAllTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	for _, tools := range ToolRegistry {
		for _, tool := range tools {
			if tool.AddFunc != nil {
				tool.AddFunc(s, executeShipCommand)
			}
		}
	}
}

// RegisterToolsByCategory registers tools from a specific category
func RegisterToolsByCategory(category string, s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	if tools, ok := ToolRegistry[category]; ok {
		for _, tool := range tools {
			if tool.AddFunc != nil {
				tool.AddFunc(s, executeShipCommand)
			}
		}
	}
}

// RegisterToolByName registers a specific tool by name
func RegisterToolByName(name string, s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	for _, tools := range ToolRegistry {
		for _, tool := range tools {
			if tool.Name == name && tool.AddFunc != nil {
				tool.AddFunc(s, executeShipCommand)
				return
			}
		}
	}
}
