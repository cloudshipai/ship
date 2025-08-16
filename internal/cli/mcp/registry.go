package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

// ToolInfo contains information about a tool
type ToolInfo struct {
	Name        string
	Description string
	Category    string
	AddFunc     func(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc)
}

// ToolRegistry contains all available tools organized by category
var ToolRegistry = map[string][]ToolInfo{
	"security": {
		{Name: "gitleaks", Description: "Secret detection in code and git history", AddFunc: AddGitleaksTools},
		{Name: "trivy", Description: "Comprehensive vulnerability scanner", AddFunc: AddTrivyTools},
		{Name: "grype", Description: "Container vulnerability scanner", AddFunc: AddGrypeTools},
		{Name: "syft", Description: "SBOM generation tool", AddFunc: AddSyftTools},
		{Name: "checkov", Description: "Infrastructure as code static analysis", AddFunc: AddCheckovTools},
		{Name: "terrascan", Description: "IaC security scanner", AddFunc: AddTerrascanTools},
		{Name: "semgrep", Description: "Static analysis for security", AddFunc: AddSemgrepTools},
		{Name: "actionlint", Description: "GitHub Actions workflow linter", AddFunc: AddActionlintTools},
		{Name: "hadolint", Description: "Dockerfile linter", AddFunc: AddHadolintTools},
		{Name: "conftest", Description: "OPA policy testing", AddFunc: AddConftestTools},
		{Name: "kube-bench", Description: "Kubernetes CIS benchmark", AddFunc: AddKubeBenchTools},
		{Name: "kube-hunter", Description: "Kubernetes penetration testing", AddFunc: AddKubeHunterTools},
		{Name: "falco", Description: "Runtime security monitoring", AddFunc: AddFalcoTools},
		{Name: "nikto", Description: "Web server security scanner", AddFunc: AddNiktoTools},
		{Name: "zap", Description: "OWASP ZAP web application scanner", AddFunc: AddZapTools},
		{Name: "git-secrets", Description: "Git repository secret scanner", AddFunc: AddGitSecretsTools},
		{Name: "trufflehog", Description: "Advanced secret scanning with verification", AddFunc: AddTrufflehogTools},
		{Name: "kubescape", Description: "Kubernetes security scanner", AddFunc: AddKubescapeTools},
		{Name: "dockle", Description: "Container image linter", AddFunc: AddDockleTools},
		{Name: "sops", Description: "Secrets management", AddFunc: AddSOPSTools},
		{Name: "ossf-scorecard", Description: "OSSF security scorecard", AddFunc: AddOSSFScorecardTools},
		{Name: "steampipe", Description: "Cloud asset querying with SQL", AddFunc: AddSteampipeTools},
		{Name: "allstar", Description: "Kubernetes security policy enforcement", AddFunc: AddAllstarTools},
		{Name: "cfn-nag", Description: "CloudFormation security linter", AddFunc: AddCfnNagTools},
		{Name: "gatekeeper", Description: "OPA Gatekeeper policy validation", AddFunc: AddGatekeeperTools},
		{Name: "history-scrub", Description: "Git history cleaning and secret removal", AddFunc: AddHistoryScrubTools},
		{Name: "license-detector", Description: "Software license detection", AddFunc: AddLicenseDetectorTools},
		{Name: "openscap", Description: "Security compliance scanning", AddFunc: AddOpenSCAPTools},
		{Name: "osv-scanner", Description: "Open Source Vulnerability scanning", AddFunc: AddOSVScannerTools},
		{Name: "scout-suite", Description: "Multi-cloud security auditing", AddFunc: AddScoutSuiteTools},
		{Name: "powerpipe", Description: "Infrastructure benchmarking", AddFunc: AddPowerpipeTools},
		{Name: "container-registry", Description: "Container registry operations", AddFunc: AddContainerRegistryTools},
		{Name: "infrascan", Description: "Infrastructure security scanning", AddFunc: AddInfrascanTools},
		{Name: "check-ssl-cert", Description: "SSL certificate validation", AddFunc: AddCheckSSLCertTools},
		{Name: "step-ca", Description: "Certificate authority operations", AddFunc: AddStepCATools},
		{Name: "github-admin", Description: "GitHub administration tools", AddFunc: AddGitHubAdminTools},
		{Name: "github-packages", Description: "GitHub Packages security", AddFunc: AddGitHubPackagesTools},
		{Name: "trivy-golden", Description: "Enhanced Trivy for golden images", AddFunc: AddTrivyGoldenTools},
	},
	"terraform": {
		{Name: "tflint", Description: "Terraform linter", AddFunc: AddTfLintTools},
		{Name: "terraform-docs", Description: "Terraform documentation generator", AddFunc: AddTerraformDocsTools},
		{Name: "infracost", Description: "Infrastructure cost estimation", AddFunc: AddInfracostTools},
		{Name: "inframap", Description: "Infrastructure visualization", AddFunc: AddInfraMapTools},
		{Name: "iac-plan", Description: "Infrastructure as code planning", AddFunc: AddIacPlanTools},
		{Name: "terraformer", Description: "Infrastructure import and management", AddFunc: AddTerraformerTools},
		{Name: "tfstate-reader", Description: "Terraform state analysis", AddFunc: AddTfstateReaderTools},
		{Name: "openinfraquote", Description: "Infrastructure cost estimation", AddFunc: AddOpenInfraQuoteTools},
	},
	"kubernetes": {
		{Name: "velero", Description: "Kubernetes backup and restore", AddFunc: AddVeleroTools},
		{Name: "goldilocks", Description: "Kubernetes resource recommendations", AddFunc: AddGoldilocksTools},
		{Name: "fleet", Description: "GitOps for Kubernetes", AddFunc: AddFleetTools},
		{Name: "kuttl", Description: "Kubernetes testing framework", AddFunc: AddKuttlTools},
		{Name: "litmus", Description: "Chaos engineering for Kubernetes", AddFunc: AddLitmusTools},
		{Name: "cert-manager", Description: "Certificate management", AddFunc: AddCertManagerTools},
		{Name: "k8s-network-policy", Description: "Kubernetes network policy management", AddFunc: AddK8sNetworkPolicyTools},
		{Name: "kyverno", Description: "Kubernetes policy management", AddFunc: AddKyvernoTools},
		{Name: "kyverno-multitenant", Description: "Multi-tenant Kyverno policies", AddFunc: AddKyvernoMultitenantTools},
	},
	"cloud": {
		{Name: "cloudquery", Description: "Cloud asset inventory", AddFunc: AddCloudQueryTools},
		{Name: "custodian", Description: "Cloud governance engine", AddFunc: AddCustodianTools},
		{Name: "packer", Description: "Machine image building", AddFunc: AddPackerTools},
	},
	"supply-chain": {
		{Name: "cosign", Description: "Container signing and verification", AddFunc: AddCosignTools},
		{Name: "cosign-golden", Description: "Advanced container signing workflows", AddFunc: AddCosignGoldenTools},
		{Name: "sigstore-policy-controller", Description: "Sigstore policy enforcement", AddFunc: AddSigstorePolicyControllerTools},
		{Name: "guac", Description: "Graph for Understanding Artifact Composition", AddFunc: AddGuacTools},
		{Name: "rekor", Description: "Transparency log", AddFunc: AddRekorTools},
		{Name: "in-toto", Description: "Supply chain attestation", AddFunc: AddInTotoTools},
		{Name: "slsa-verifier", Description: "SLSA provenance verification", AddFunc: AddSLSAVerifierTools},
		{Name: "dependency-track", Description: "OWASP Dependency-Track SBOM analysis", AddFunc: AddDependencyTrackTools},
	},
	"aws": {
		{Name: "cloudsplaining", Description: "AWS IAM policy scanner", AddFunc: AddCloudsplainingTools},
		{Name: "parliament", Description: "AWS IAM policy linter", AddFunc: AddParliamentTools},
		{Name: "pmapper", Description: "AWS IAM privilege escalation analysis", AddFunc: AddPMapperTools},
		{Name: "policy-sentry", Description: "AWS IAM policy generator", AddFunc: AddPolicySentryTools},
		{Name: "prowler", Description: "Multi-cloud security assessment", AddFunc: AddProwlerTools},
		{Name: "aws-iam-rotation", Description: "AWS IAM credential rotation", AddFunc: AddAWSIAMRotationTools},
		{Name: "aws-pricing", Description: "AWS pricing and cost calculator", AddFunc: AddAWSPricingTools},
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