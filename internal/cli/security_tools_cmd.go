package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	shipdagger "github.com/cloudshipai/ship/internal/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var securityToolsCmd = &cobra.Command{
	Use:   "security",
	Short: "Run security analysis tools",
	Long:  `Run various security analysis tools including secret detection, vulnerability scanning, and compliance checking`,
}

var gitleaksCmd = &cobra.Command{
	Use:   "gitleaks [directory]",
	Short: "Scan for secrets using Gitleaks",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGitleaks,
}

var grypeScanCmd = &cobra.Command{
	Use:   "grype [target]",
	Short: "Scan for vulnerabilities using Grype",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGrype,
}

var syftScanCmd = &cobra.Command{
	Use:   "syft [target]",
	Short: "Generate SBOM using Syft",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSyft,
}

var prowlerScanCmd = &cobra.Command{
	Use:   "prowler [provider]",
	Short: "Multi-cloud security assessment using Prowler",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runProwler,
}

var truffleHogCmd = &cobra.Command{
	Use:   "trufflehog [target]",
	Short: "Verified secret detection using TruffleHog",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runTruffleHog,
}

var cosignCmd = &cobra.Command{
	Use:   "cosign [command]",
	Short: "Container signing and verification using Cosign",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runCosign,
}

var cloudsplainingCmd = &cobra.Command{
	Use:   "cloudsplaining [command]",
	Short: "AWS IAM security assessment using Cloudsplaining",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runCloudsplaining,
}

var parliamentCmd = &cobra.Command{
	Use:   "parliament [policy-file]",
	Short: "AWS IAM policy linting using Parliament",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runParliament,
}

var pmapperCmd = &cobra.Command{
	Use:   "pmapper [command]",
	Short: "AWS IAM privilege mapping using PMapper",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runPMapper,
}

var policySentryCmd = &cobra.Command{
	Use:   "policy-sentry [command]",
	Short: "AWS IAM policy generation using Policy Sentry",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runPolicySentry,
}

var actionlintCmd = &cobra.Command{
	Use:   "actionlint [directory]",
	Short: "Lint GitHub Actions workflows using actionlint",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runActionlint,
}

var trivyCmd = &cobra.Command{
	Use:   "trivy [target]",
	Short: "Comprehensive vulnerability scanning using Trivy",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runTrivy,
}

var semgrepCmd = &cobra.Command{
	Use:   "semgrep [directory]",
	Short: "Static analysis using Semgrep",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSemgrep,
}

var hadolintCmd = &cobra.Command{
	Use:   "hadolint [dockerfile]",
	Short: "Dockerfile linting using Hadolint",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runHadolint,
}

var cfnNagCmd = &cobra.Command{
	Use:   "cfn-nag [template]",
	Short: "CloudFormation security scanning using cfn-nag",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCfnNag,
}

var conftestCmd = &cobra.Command{
	Use:   "conftest [directory]",
	Short: "OPA policy testing using Conftest",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runConftest,
}

var gitSecretsCmd = &cobra.Command{
	Use:   "git-secrets [directory]",
	Short: "Git secrets scanning using git-secrets",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGitSecrets,
}

var kubeBenchCmd = &cobra.Command{
	Use:   "kube-bench",
	Short: "Kubernetes security benchmarks using kube-bench",
	RunE:  runKubeBench,
}

var kubeHunterCmd = &cobra.Command{
	Use:   "kube-hunter [target]",
	Short: "Kubernetes penetration testing using kube-hunter",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runKubeHunter,
}

var zapCmd = &cobra.Command{
	Use:   "zap [target]",
	Short: "Web application security testing using OWASP ZAP",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runZap,
}

var falcoCmd = &cobra.Command{
	Use:   "falco",
	Short: "Runtime security monitoring using Falco",
	RunE:  runFalco,
}

var slsaVerifierCmd = &cobra.Command{
	Use:   "slsa-verifier [command]",
	Short: "SLSA provenance verification using slsa-verifier",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSLSAVerifier,
}

var inTotoCmd = &cobra.Command{
	Use:   "in-toto [command]",
	Short: "Supply chain attestation using in-toto",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runInToto,
}

var gatekeeperCmd = &cobra.Command{
	Use:   "gatekeeper [command]",
	Short: "OPA Gatekeeper policy validation",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runGatekeeper,
}

var kubescapeCmd = &cobra.Command{
	Use:   "kubescape [command]",
	Short: "Kubernetes security scanning using Kubescape",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runKubescape,
}

var dockleCmd = &cobra.Command{
	Use:   "dockle [image]",
	Short: "Container image linting using Dockle",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runDockle,
}

var sopsCmd = &cobra.Command{
	Use:   "sops [command]",
	Short: "Secrets management using Mozilla SOPS",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSOPS,
}

func init() {
	rootCmd.AddCommand(securityToolsCmd)
	securityToolsCmd.AddCommand(gitleaksCmd)
	securityToolsCmd.AddCommand(grypeScanCmd)
	securityToolsCmd.AddCommand(syftScanCmd)
	securityToolsCmd.AddCommand(prowlerScanCmd)
	securityToolsCmd.AddCommand(truffleHogCmd)
	securityToolsCmd.AddCommand(cosignCmd)
	securityToolsCmd.AddCommand(cloudsplainingCmd)
	securityToolsCmd.AddCommand(parliamentCmd)
	securityToolsCmd.AddCommand(pmapperCmd)
	securityToolsCmd.AddCommand(policySentryCmd)
	securityToolsCmd.AddCommand(actionlintCmd)
	securityToolsCmd.AddCommand(trivyCmd)
	securityToolsCmd.AddCommand(semgrepCmd)
	securityToolsCmd.AddCommand(hadolintCmd)
	securityToolsCmd.AddCommand(cfnNagCmd)
	securityToolsCmd.AddCommand(conftestCmd)
	securityToolsCmd.AddCommand(gitSecretsCmd)
	securityToolsCmd.AddCommand(kubeBenchCmd)
	securityToolsCmd.AddCommand(kubeHunterCmd)
	securityToolsCmd.AddCommand(zapCmd)
	securityToolsCmd.AddCommand(falcoCmd)
	securityToolsCmd.AddCommand(slsaVerifierCmd)
	securityToolsCmd.AddCommand(inTotoCmd)
	securityToolsCmd.AddCommand(gatekeeperCmd)
	securityToolsCmd.AddCommand(kubescapeCmd)
	securityToolsCmd.AddCommand(dockleCmd)
	securityToolsCmd.AddCommand(sopsCmd)

	// Gitleaks flags
	gitleaksCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	gitleaksCmd.Flags().StringP("config", "c", "", "Path to Gitleaks configuration file")
	gitleaksCmd.Flags().BoolP("git", "g", false, "Scan git repository history")

	// Grype flags
	grypeScanCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	grypeScanCmd.Flags().StringP("severity", "s", "", "Minimum severity level (negligible, low, medium, high, critical)")
	grypeScanCmd.Flags().StringP("format", "f", "json", "Output format (json, table, cyclonedx)")

	// Syft flags
	syftScanCmd.Flags().StringP("output", "o", "", "Output file to save SBOM (default: print to stdout)")
	syftScanCmd.Flags().StringP("format", "f", "json", "Output format (json, spdx-json, cyclonedx-json, table)")
	syftScanCmd.Flags().StringP("package-type", "p", "", "Package type filter (npm, python, go, java)")

	// Prowler flags
	prowlerScanCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	prowlerScanCmd.Flags().StringP("region", "r", "us-east-1", "AWS region")
	prowlerScanCmd.Flags().StringP("compliance", "c", "", "Compliance framework (cis, pci, gdpr, hipaa)")
	prowlerScanCmd.Flags().StringP("services", "s", "", "Specific services to scan")

	// TruffleHog flags
	truffleHogCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	truffleHogCmd.Flags().StringP("type", "t", "filesystem", "Scan type (filesystem, git, github, docker, s3)")
	truffleHogCmd.Flags().BoolP("verify", "v", false, "Verify found secrets")
	truffleHogCmd.Flags().StringP("token", "", "", "GitHub token for repository access")

	// Cosign flags
	cosignCmd.Flags().StringP("key", "k", "", "Path to public/private key file")
	cosignCmd.Flags().BoolP("keyless", "", false, "Use keyless signing/verification")
	cosignCmd.Flags().StringP("output", "o", "", "Output file to save results")

	// Cloudsplaining flags
	cloudsplainingCmd.Flags().StringP("profile", "p", "default", "AWS profile to use")
	cloudsplainingCmd.Flags().StringP("output", "o", "", "Output file to save results")
	cloudsplainingCmd.Flags().StringP("minimize", "m", "", "Statement ID to minimize")

	// Parliament flags
	parliamentCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	parliamentCmd.Flags().BoolP("community", "c", false, "Include community auditors")
	parliamentCmd.Flags().StringP("auditors", "a", "", "Path to private auditors directory")
	parliamentCmd.Flags().StringP("severity", "s", "", "Minimum severity level")

	// PMapper flags
	pmapperCmd.Flags().StringP("profile", "p", "default", "AWS profile to use")
	pmapperCmd.Flags().StringP("output", "o", "", "Output file to save results")
	pmapperCmd.Flags().StringP("format", "f", "", "Output format for visualization")

	// Policy Sentry flags
	policySentryCmd.Flags().StringP("template-type", "t", "crud", "Template type (crud, actions)")
	policySentryCmd.Flags().StringP("input", "i", "", "Input YAML file")
	policySentryCmd.Flags().StringP("output", "o", "", "Output file to save policy")
	policySentryCmd.Flags().StringP("service", "s", "", "AWS service for query commands")

	// Actionlint flags
	actionlintCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	actionlintCmd.Flags().StringP("config", "c", "", "Path to actionlint configuration file")

	// Trivy flags
	trivyCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	trivyCmd.Flags().StringP("type", "t", "fs", "Scan type (fs, image, repo, config)")
	trivyCmd.Flags().StringP("severity", "s", "HIGH,CRITICAL", "Severity levels")

	// Semgrep flags
	semgrepCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	semgrepCmd.Flags().StringP("config", "c", "auto", "Semgrep configuration/ruleset")
	semgrepCmd.Flags().StringP("severity", "s", "ERROR", "Minimum severity level")

	// Hadolint flags
	hadolintCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	hadolintCmd.Flags().BoolP("directory", "d", false, "Scan all Dockerfiles in directory")

	// CFN-nag flags
	cfnNagCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	cfnNagCmd.Flags().StringP("rules", "r", "", "Path to custom rules directory")

	// Conftest flags
	conftestCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	conftestCmd.Flags().StringP("policy", "p", "", "Path to policy directory (required)")

	// Git-secrets flags
	gitSecretsCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	gitSecretsCmd.Flags().BoolP("aws", "a", false, "Include AWS secret patterns")

	// Kube-bench flags
	kubeBenchCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	kubeBenchCmd.Flags().StringP("kubeconfig", "k", "", "Path to kubeconfig file")
	kubeBenchCmd.Flags().StringP("node-type", "n", "", "Node type (master, node)")

	// Kube-hunter flags
	kubeHunterCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	kubeHunterCmd.Flags().StringP("scan-type", "t", "remote", "Scan type (remote, cidr, interface, pod)")
	kubeHunterCmd.Flags().StringP("kubeconfig", "k", "", "Path to kubeconfig file (for pod scan)")

	// ZAP flags
	zapCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	zapCmd.Flags().StringP("scan-type", "t", "baseline", "Scan type (baseline, full, api)")
	zapCmd.Flags().StringP("api-spec", "a", "", "Path to OpenAPI/Swagger spec file (for API scan)")
	zapCmd.Flags().IntP("max-duration", "m", 60, "Maximum scan duration in minutes (for full scan)")
	zapCmd.Flags().StringP("context", "c", "", "Path to ZAP context file")

	// Falco flags
	falcoCmd.Flags().StringP("output", "o", "", "Output file to save results (default: print to stdout)")
	falcoCmd.Flags().StringP("rules", "r", "", "Path to custom rules directory")
	falcoCmd.Flags().StringP("kubeconfig", "k", "", "Path to kubeconfig file")
	falcoCmd.Flags().BoolP("validate", "v", false, "Validate rules only")

	// SLSA Verifier flags
	slsaVerifierCmd.Flags().StringP("artifact", "a", "", "Path to artifact file")
	slsaVerifierCmd.Flags().StringP("provenance", "p", "", "Path to provenance file")
	slsaVerifierCmd.Flags().StringP("source-uri", "s", "", "Source URI for verification")
	slsaVerifierCmd.Flags().StringP("builder-id", "b", "", "Builder ID for verification")
	slsaVerifierCmd.Flags().BoolP("print-provenance", "", false, "Print provenance information")

	// in-toto flags
	inTotoCmd.Flags().StringP("step-name", "n", "", "Step name for attestation")
	inTotoCmd.Flags().StringP("key", "k", "", "Path to signing key")
	inTotoCmd.Flags().StringP("layout", "l", "", "Path to layout file")
	inTotoCmd.Flags().StringArrayP("materials", "m", []string{}, "Material patterns")
	inTotoCmd.Flags().StringArrayP("products", "", []string{}, "Product patterns")

	// Gatekeeper flags
	gatekeeperCmd.Flags().StringP("constraints", "c", "", "Path to constraints directory")
	gatekeeperCmd.Flags().StringP("templates", "t", "", "Path to constraint templates directory")
	gatekeeperCmd.Flags().StringP("resources", "r", "", "Path to resources directory")
	gatekeeperCmd.Flags().StringP("format", "f", "pretty", "Output format")
	gatekeeperCmd.Flags().BoolP("verbose", "v", false, "Verbose output")

	// Kubescape flags
	kubescapeCmd.Flags().StringP("framework", "f", "nsa", "Security framework to use")
	kubescapeCmd.Flags().StringP("format", "", "pretty-printer", "Output format")
	kubescapeCmd.Flags().StringP("severity", "s", "", "Severity threshold")
	kubescapeCmd.Flags().StringP("namespace", "n", "", "Namespace to scan")
	kubescapeCmd.Flags().StringP("kubeconfig", "k", "", "Path to kubeconfig file")

	// Dockle flags
	dockleCmd.Flags().StringP("format", "f", "json", "Output format")
	dockleCmd.Flags().StringP("exit-level", "e", "warn", "Exit level")
	dockleCmd.Flags().StringArrayP("accept-key", "", []string{}, "Accept specific check keys")
	dockleCmd.Flags().StringArrayP("ignore", "i", []string{}, "Ignore specific findings")

	// SOPS flags
	sopsCmd.Flags().StringP("kms", "", "", "KMS ARN for encryption")
	sopsCmd.Flags().StringP("pgp", "", "", "PGP fingerprint for encryption")
	sopsCmd.Flags().StringP("age", "", "", "Age public key for encryption")
	sopsCmd.Flags().StringP("key-file", "k", "", "Path to key file")
	sopsCmd.Flags().BoolP("in-place", "i", false, "Edit file in place")
}

func runGitleaks(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	config, _ := cmd.Flags().GetString("config")
	gitScan, _ := cmd.Flags().GetBool("git")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewGitleaksModule(engine.GetClient())

	var result string
	if config != "" {
		result, err = module.ScanWithConfig(ctx, dir, config)
	} else if gitScan {
		result, err = module.ScanGitRepo(ctx, dir)
	} else {
		result, err = module.ScanDirectory(ctx, dir)
	}

	if err != nil {
		return fmt.Errorf("gitleaks scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runGrype(cmd *cobra.Command, args []string) error {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	severity, _ := cmd.Flags().GetString("severity")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewGrypeModule(engine.GetClient())

	var result string
	if severity != "" {
		result, err = module.ScanWithSeverity(ctx, target, severity)
	} else if target[:6] == "image:" {
		result, err = module.ScanImage(ctx, target[6:])
	} else {
		result, err = module.ScanDirectory(ctx, target)
	}

	if err != nil {
		return fmt.Errorf("grype scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runSyft(cmd *cobra.Command, args []string) error {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	packageType, _ := cmd.Flags().GetString("package-type")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewSyftModule(engine.GetClient())

	var result string
	if packageType != "" {
		result, err = module.GenerateSBOMFromPackage(ctx, target, packageType, format)
	} else if target[:6] == "image:" {
		result, err = module.GenerateSBOMFromImage(ctx, target[6:], format)
	} else {
		result, err = module.GenerateSBOMFromDirectory(ctx, target, format)
	}

	if err != nil {
		return fmt.Errorf("syft scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runProwler(cmd *cobra.Command, args []string) error {
	provider := "aws"
	if len(args) > 0 {
		provider = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	region, _ := cmd.Flags().GetString("region")
	compliance, _ := cmd.Flags().GetString("compliance")
	services, _ := cmd.Flags().GetString("services")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewProwlerModule(engine.GetClient())

	var result string
	switch provider {
	case "aws":
		if compliance != "" {
			result, err = module.ScanWithCompliance(ctx, provider, compliance, region)
		} else if services != "" {
			result, err = module.ScanSpecificServices(ctx, provider, services, region)
		} else {
			result, err = module.ScanAWS(ctx, provider, region)
		}
	case "azure":
		result, err = module.ScanAzure(ctx)
	case "gcp":
		projectId := os.Getenv("GOOGLE_CLOUD_PROJECT")
		if projectId == "" {
			return fmt.Errorf("GOOGLE_CLOUD_PROJECT environment variable required for GCP scanning")
		}
		result, err = module.ScanGCP(ctx, projectId)
	default:
		return fmt.Errorf("unsupported provider: %s (supported: aws, azure, gcp)", provider)
	}

	if err != nil {
		return fmt.Errorf("prowler scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runTruffleHog(cmd *cobra.Command, args []string) error {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	scanType, _ := cmd.Flags().GetString("type")
	verify, _ := cmd.Flags().GetBool("verify")
	token, _ := cmd.Flags().GetString("token")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewTruffleHogModule(engine.GetClient())

	var result string
	if verify {
		result, err = module.ScanWithVerification(ctx, target, scanType)
	} else {
		switch scanType {
		case "filesystem":
			result, err = module.ScanDirectory(ctx, target)
		case "git":
			result, err = module.ScanGitRepo(ctx, target)
		case "github":
			result, err = module.ScanGitHub(ctx, target, token)
		case "docker":
			result, err = module.ScanDockerImage(ctx, target)
		case "s3":
			result, err = module.ScanS3(ctx, target)
		default:
			result, err = module.ScanDirectory(ctx, target)
		}
	}

	if err != nil {
		return fmt.Errorf("trufflehog scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runCosign(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("cosign command required (verify, sign, generate-key-pair)")
	}

	command := args[0]
	target := ""
	if len(args) > 1 {
		target = args[1]
	}

	keyPath, _ := cmd.Flags().GetString("key")
	keyless, _ := cmd.Flags().GetBool("keyless")
	output, _ := cmd.Flags().GetString("output")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewCosignModule(engine.GetClient())

	var result string
	switch command {
	case "verify":
		if target == "" {
			return fmt.Errorf("image name required for verify command")
		}
		if keyPath != "" {
			result, err = module.VerifyImageWithKey(ctx, target, keyPath)
		} else {
			result, err = module.VerifyImage(ctx, target)
		}
	case "sign":
		if target == "" {
			return fmt.Errorf("image name required for sign command")
		}
		if keyless {
			result, err = module.SignImageKeyless(ctx, target)
		} else if keyPath != "" {
			result, err = module.SignImage(ctx, target, keyPath)
		} else {
			return fmt.Errorf("either --keyless or --key required for signing")
		}
	case "generate-key-pair":
		outputDir := "."
		if output != "" {
			outputDir = filepath.Dir(output)
		}
		result, err = module.GenerateKeyPair(ctx, outputDir)
	default:
		return fmt.Errorf("unsupported cosign command: %s", command)
	}

	if err != nil {
		return fmt.Errorf("cosign %s failed: %w", command, err)
	}

	return handleOutput(result, output)
}

func runCloudsplaining(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("cloudsplaining command required (scan-account, scan-policy)")
	}

	command := args[0]
	profile, _ := cmd.Flags().GetString("profile")
	output, _ := cmd.Flags().GetString("output")
	minimize, _ := cmd.Flags().GetString("minimize")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewCloudsplainingModule(engine.GetClient())

	var result string
	switch command {
	case "scan-account":
		if minimize != "" {
			result, err = module.ScanWithMinimization(ctx, profile, minimize)
		} else {
			result, err = module.ScanAccountAuthorization(ctx, profile)
		}
	case "scan-policy":
		if len(args) < 2 {
			return fmt.Errorf("policy file path required for scan-policy command")
		}
		result, err = module.ScanPolicyFile(ctx, args[1])
	default:
		return fmt.Errorf("unsupported cloudsplaining command: %s", command)
	}

	if err != nil {
		return fmt.Errorf("cloudsplaining %s failed: %w", command, err)
	}

	return handleOutput(result, output)
}

func runParliament(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("policy file path required")
	}

	policyFile := args[0]
	output, _ := cmd.Flags().GetString("output")
	community, _ := cmd.Flags().GetBool("community")
	auditorsPath, _ := cmd.Flags().GetString("auditors")
	severity, _ := cmd.Flags().GetString("severity")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewParliamentModule(engine.GetClient())

	var result string
	if auditorsPath != "" {
		result, err = module.LintWithPrivateAuditors(ctx, policyFile, auditorsPath)
	} else if community {
		result, err = module.LintWithCommunityAuditors(ctx, policyFile)
	} else if severity != "" {
		result, err = module.LintWithSeverityFilter(ctx, policyFile, severity)
	} else {
		result, err = module.LintPolicyFile(ctx, policyFile)
	}

	if err != nil {
		return fmt.Errorf("parliament lint failed: %w", err)
	}

	return handleOutput(result, output)
}

func runPMapper(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("pmapper command required (create-graph, query, privesc, admin, list)")
	}

	command := args[0]
	profile, _ := cmd.Flags().GetString("profile")
	output, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewPMapperModule(engine.GetClient())

	var result string
	switch command {
	case "create-graph":
		result, err = module.CreateGraph(ctx, profile)
	case "query":
		if len(args) < 4 {
			return fmt.Errorf("query requires: principal action [resource]")
		}
		principal := args[1]
		action := args[2]
		resource := ""
		if len(args) > 3 {
			resource = args[3]
		}
		result, err = module.QueryAccess(ctx, profile, principal, action, resource)
	case "privesc":
		if len(args) < 2 {
			return fmt.Errorf("privesc requires principal name")
		}
		result, err = module.FindPrivilegeEscalation(ctx, profile, args[1])
	case "admin":
		if len(args) < 2 {
			return fmt.Errorf("admin requires principal name")
		}
		result, err = module.CheckAdminAccess(ctx, profile, args[1])
	case "list":
		result, err = module.ListPrincipals(ctx, profile)
	case "visualize":
		result, err = module.VisualizeGraph(ctx, profile, format)
	default:
		return fmt.Errorf("unsupported pmapper command: %s", command)
	}

	if err != nil {
		return fmt.Errorf("pmapper %s failed: %w", command, err)
	}

	return handleOutput(result, output)
}

func runPolicySentry(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("policy-sentry command required (create-template, write-policy, query)")
	}

	command := args[0]
	templateType, _ := cmd.Flags().GetString("template-type")
	input, _ := cmd.Flags().GetString("input")
	output, _ := cmd.Flags().GetString("output")
	service, _ := cmd.Flags().GetString("service")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewPolicySentryModule(engine.GetClient())

	var result string
	switch command {
	case "create-template":
		result, err = module.CreateTemplate(ctx, templateType, output)
	case "write-policy":
		if input == "" {
			return fmt.Errorf("input file required for write-policy command")
		}
		result, err = module.WritePolicy(ctx, input)
	case "query-actions":
		if service == "" {
			return fmt.Errorf("service required for query-actions command")
		}
		result, err = module.QueryActionTable(ctx, service)
	case "query-conditions":
		if service == "" {
			return fmt.Errorf("service required for query-conditions command")
		}
		result, err = module.QueryConditionTable(ctx, service)
	default:
		return fmt.Errorf("unsupported policy-sentry command: %s", command)
	}

	if err != nil {
		return fmt.Errorf("policy-sentry %s failed: %w", command, err)
	}

	return handleOutput(result, output)
}

func runActionlint(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	output, _ := cmd.Flags().GetString("output")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewActionlintModule(engine.GetClient())
	result, err := module.ScanDirectory(ctx, dir)
	if err != nil {
		return fmt.Errorf("actionlint scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runTrivy(cmd *cobra.Command, args []string) error {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	scanType, _ := cmd.Flags().GetString("type")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewTrivyModule(engine.GetClient())

	var result string
	switch scanType {
	case "image":
		result, err = module.ScanImage(ctx, target)
	case "repo":
		result, err = module.ScanRepository(ctx, target)
	case "config":
		result, err = module.ScanConfig(ctx, target)
	default: // fs
		result, err = module.ScanFilesystem(ctx, target)
	}

	if err != nil {
		return fmt.Errorf("trivy scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runSemgrep(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	config, _ := cmd.Flags().GetString("config")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewSemgrepModule(engine.GetClient())

	var result string
	if config != "auto" {
		result, err = module.ScanWithRuleset(ctx, dir, config)
	} else {
		result, err = module.ScanDirectory(ctx, dir)
	}

	if err != nil {
		return fmt.Errorf("semgrep scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runHadolint(cmd *cobra.Command, args []string) error {
	target := "Dockerfile"
	if len(args) > 0 {
		target = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	directory, _ := cmd.Flags().GetBool("directory")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewHadolintModule(engine.GetClient())

	var result string
	if directory {
		result, err = module.ScanDirectory(ctx, target)
	} else {
		result, err = module.ScanDockerfile(ctx, target)
	}

	if err != nil {
		return fmt.Errorf("hadolint scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runCfnNag(cmd *cobra.Command, args []string) error {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	rules, _ := cmd.Flags().GetString("rules")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewCfnNagModule(engine.GetClient())

	var result string
	if rules != "" {
		result, err = module.ScanWithRules(ctx, target, rules)
	} else {
		// Check if target is a file or directory
		if filepath.Ext(target) != "" {
			result, err = module.ScanTemplate(ctx, target)
		} else {
			result, err = module.ScanDirectory(ctx, target)
		}
	}

	if err != nil {
		return fmt.Errorf("cfn-nag scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runConftest(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	policy, _ := cmd.Flags().GetString("policy")

	if policy == "" {
		return fmt.Errorf("policy directory is required (use -p flag)")
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewConftestModule(engine.GetClient())
	result, err := module.TestWithPolicy(ctx, dir, policy)
	if err != nil {
		return fmt.Errorf("conftest scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runGitSecrets(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	aws, _ := cmd.Flags().GetBool("aws")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewGitSecretsModule(engine.GetClient())

	var result string
	if aws {
		result, err = module.ScanWithAwsProviders(ctx, dir)
	} else {
		result, err = module.ScanRepository(ctx, dir)
	}

	if err != nil {
		return fmt.Errorf("git-secrets scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runKubeBench(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	nodeType, _ := cmd.Flags().GetString("node-type")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewKubeBenchModule(engine.GetClient())

	var result string
	switch nodeType {
	case "master":
		result, err = module.RunMasterBenchmark(ctx, kubeconfig)
	case "node":
		result, err = module.RunNodeBenchmark(ctx, kubeconfig)
	default:
		result, err = module.RunBenchmark(ctx, kubeconfig)
	}

	if err != nil {
		return fmt.Errorf("kube-bench scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runKubeHunter(cmd *cobra.Command, args []string) error {
	target := ""
	if len(args) > 0 {
		target = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	scanType, _ := cmd.Flags().GetString("scan-type")
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewKubeHunterModule(engine.GetClient())

	var result string
	switch scanType {
	case "remote":
		if target == "" {
			return fmt.Errorf("target required for remote scan")
		}
		result, err = module.ScanRemote(ctx, target)
	case "cidr":
		if target == "" {
			return fmt.Errorf("CIDR required for CIDR scan")
		}
		result, err = module.ScanCIDR(ctx, target)
	case "interface":
		if target == "" {
			return fmt.Errorf("interface required for interface scan")
		}
		result, err = module.ScanInterface(ctx, target)
	case "pod":
		result, err = module.ScanPod(ctx, kubeconfig)
	default:
		return fmt.Errorf("unsupported scan type: %s (supported: remote, cidr, interface, pod)", scanType)
	}

	if err != nil {
		return fmt.Errorf("kube-hunter scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runZap(cmd *cobra.Command, args []string) error {
	target := ""
	if len(args) > 0 {
		target = args[0]
	}

	if target == "" {
		return fmt.Errorf("target URL required")
	}

	output, _ := cmd.Flags().GetString("output")
	scanType, _ := cmd.Flags().GetString("scan-type")
	apiSpec, _ := cmd.Flags().GetString("api-spec")
	maxDuration, _ := cmd.Flags().GetInt("max-duration")
	contextFile, _ := cmd.Flags().GetString("context")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewZapModule(engine.GetClient())

	var result string
	switch scanType {
	case "baseline":
		if contextFile != "" {
			result, err = module.ScanWithContext(ctx, target, contextFile)
		} else {
			result, err = module.BaselineScan(ctx, target)
		}
	case "full":
		result, err = module.FullScan(ctx, target, maxDuration)
	case "api":
		if apiSpec == "" {
			return fmt.Errorf("API spec file required for API scan (use --api-spec)")
		}
		result, err = module.ApiScan(ctx, target, apiSpec)
	default:
		return fmt.Errorf("unsupported scan type: %s (supported: baseline, full, api)", scanType)
	}

	if err != nil {
		return fmt.Errorf("ZAP scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runFalco(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	rules, _ := cmd.Flags().GetString("rules")
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	validate, _ := cmd.Flags().GetBool("validate")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewFalcoModule(engine.GetClient())

	var result string
	if validate {
		if rules == "" {
			return fmt.Errorf("rules directory required for validation (use --rules)")
		}
		result, err = module.ValidateRules(ctx, rules)
	} else if rules != "" {
		result, err = module.RunWithCustomRules(ctx, rules, kubeconfig)
	} else {
		result, err = module.RunWithDefaultRules(ctx, kubeconfig)
	}

	if err != nil {
		return fmt.Errorf("falco scan failed: %w", err)
	}

	return handleOutput(result, output)
}

func runSLSAVerifier(cmd *cobra.Command, args []string) error {
	command := args[0]
	output, _ := cmd.Flags().GetString("output")
	artifact, _ := cmd.Flags().GetString("artifact")
	provenance, _ := cmd.Flags().GetString("provenance")
	sourceURI, _ := cmd.Flags().GetString("source-uri")
	builderID, _ := cmd.Flags().GetString("builder-id")
	printProvenance, _ := cmd.Flags().GetBool("print-provenance")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewSLSAVerifierModule(engine.GetClient())

	opts := []modules.SLSAVerifierOption{}
	if sourceURI != "" {
		opts = append(opts, modules.WithSourceURI(sourceURI))
	}
	if builderID != "" {
		opts = append(opts, modules.WithBuilderID(builderID))
	}
	if printProvenance {
		opts = append(opts, modules.WithPrintProvenance(printProvenance))
	}

	var container *dagger.Container
	switch command {
	case "verify-artifact":
		if artifact == "" {
			return fmt.Errorf("artifact path required for verify-artifact (use --artifact)")
		}
		container, err = module.VerifyProvenance(ctx, artifact, provenance, opts...)
	case "verify-image":
		if len(args) < 2 {
			return fmt.Errorf("image reference required for verify-image")
		}
		container, err = module.VerifyImage(ctx, args[1], opts...)
	case "generate-policy":
		container, err = module.GeneratePolicy(ctx, opts...)
	default:
		return fmt.Errorf("unknown command: %s (supported: verify-artifact, verify-image, generate-policy)", command)
	}

	if err != nil {
		return fmt.Errorf("SLSA verifier failed: %w", err)
	}

	result, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get output: %w", err)
	}

	return handleOutput(result, output)
}

func runInToto(cmd *cobra.Command, args []string) error {
	command := args[0]
	output, _ := cmd.Flags().GetString("output")
	stepName, _ := cmd.Flags().GetString("step-name")
	keyPath, _ := cmd.Flags().GetString("key")
	layoutPath, _ := cmd.Flags().GetString("layout")
	materials, _ := cmd.Flags().GetStringArray("materials")
	products, _ := cmd.Flags().GetStringArray("products")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewInTotoModule(engine.GetClient())

	opts := []modules.InTotoOption{}
	if keyPath != "" {
		opts = append(opts, modules.WithKeyPath(keyPath))
	}
	if len(materials) > 0 {
		opts = append(opts, modules.WithMaterials(materials))
	}
	if len(products) > 0 {
		opts = append(opts, modules.WithProducts(products))
	}

	var container *dagger.Container
	switch command {
	case "run":
		if stepName == "" {
			return fmt.Errorf("step name required for run command (use --step-name)")
		}
		if len(args) < 2 {
			return fmt.Errorf("command to execute required for run")
		}
		container, err = module.RunStep(ctx, stepName, args[1:], opts...)
	case "verify":
		if layoutPath == "" {
			return fmt.Errorf("layout file required for verify (use --layout)")
		}
		container, err = module.VerifySupplyChain(ctx, layoutPath, opts...)
	case "record":
		if stepName == "" {
			return fmt.Errorf("step name required for record (use --step-name)")
		}
		container, err = module.RecordMetadata(ctx, stepName, opts...)
	case "generate-layout":
		container, err = module.GenerateLayout(ctx, opts...)
	default:
		return fmt.Errorf("unknown command: %s (supported: run, verify, record, generate-layout)", command)
	}

	if err != nil {
		return fmt.Errorf("in-toto operation failed: %w", err)
	}

	result, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get output: %w", err)
	}

	return handleOutput(result, output)
}

func runGatekeeper(cmd *cobra.Command, args []string) error {
	command := args[0]
	output, _ := cmd.Flags().GetString("output")
	constraintsDir, _ := cmd.Flags().GetString("constraints")
	templatesDir, _ := cmd.Flags().GetString("templates")
	resourcesDir, _ := cmd.Flags().GetString("resources")
	format, _ := cmd.Flags().GetString("format")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewGatekeeperModule(engine.GetClient())

	opts := []modules.GatekeeperOption{}
	if constraintsDir != "" {
		opts = append(opts, modules.WithConstraintsDir(constraintsDir))
	}
	if templatesDir != "" {
		opts = append(opts, modules.WithTemplatesDir(templatesDir))
	}
	if format != "" {
		opts = append(opts, modules.WithFormat(format))
	}
	if verbose {
		opts = append(opts, modules.WithVerbose(verbose))
	}

	var container *dagger.Container
	switch command {
	case "validate":
		if resourcesDir == "" {
			return fmt.Errorf("resources directory required for validate (use --resources)")
		}
		container, err = module.ValidateConstraints(ctx, resourcesDir, opts...)
	case "test":
		if resourcesDir == "" {
			return fmt.Errorf("tests directory required for test (use --resources)")
		}
		container, err = module.TestConstraints(ctx, resourcesDir, opts...)
	case "generate-template":
		if len(args) < 2 {
			return fmt.Errorf("template name required for generate-template")
		}
		container, err = module.GenerateConstraintTemplate(ctx, args[1], opts...)
	case "sync":
		container, err = module.SyncConstraints(ctx, opts...)
	case "analyze":
		container, err = module.AnalyzeViolations(ctx, opts...)
	default:
		return fmt.Errorf("unknown command: %s (supported: validate, test, generate-template, sync, analyze)", command)
	}

	if err != nil {
		return fmt.Errorf("Gatekeeper operation failed: %w", err)
	}

	result, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get output: %w", err)
	}

	return handleOutput(result, output)
}

func runKubescape(cmd *cobra.Command, args []string) error {
	command := args[0]
	output, _ := cmd.Flags().GetString("output")
	framework, _ := cmd.Flags().GetString("framework")
	format, _ := cmd.Flags().GetString("format")
	severity, _ := cmd.Flags().GetString("severity")
	namespace, _ := cmd.Flags().GetString("namespace")
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewKubescapeModule(engine.GetClient())

	opts := []modules.KubescapeOption{}
	if framework != "" {
		opts = append(opts, modules.WithFramework(framework))
	}
	if format != "" {
		opts = append(opts, modules.WithKubescapeFormat(format))
	}
	if severity != "" {
		opts = append(opts, modules.WithSeverityThreshold(severity))
	}
	if namespace != "" {
		opts = append(opts, modules.WithKubescapeNamespace(namespace))
	}
	if kubeconfig != "" {
		opts = append(opts, modules.WithKubescapeKubeconfig(kubeconfig))
	}
	if output != "" {
		opts = append(opts, modules.WithKubescapeOutput(output))
	}

	var container *dagger.Container
	switch command {
	case "cluster":
		container, err = module.ScanCluster(ctx, opts...)
	case "manifests":
		if len(args) < 2 {
			return fmt.Errorf("manifests directory required for manifests scan")
		}
		container, err = module.ScanManifests(ctx, args[1], opts...)
	case "helm":
		if len(args) < 2 {
			return fmt.Errorf("helm chart path required for helm scan")
		}
		container, err = module.ScanHelm(ctx, args[1], opts...)
	case "repo":
		if len(args) < 2 {
			return fmt.Errorf("repository path required for repo scan")
		}
		container, err = module.ScanRepository(ctx, args[1], opts...)
	case "report":
		container, err = module.GenerateReport(ctx, opts...)
	default:
		return fmt.Errorf("unknown command: %s (supported: cluster, manifests, helm, repo, report)", command)
	}

	if err != nil {
		return fmt.Errorf("Kubescape scan failed: %w", err)
	}

	result, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get output: %w", err)
	}

	return handleOutput(result, output)
}

func runDockle(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	exitLevel, _ := cmd.Flags().GetString("exit-level")
	acceptKeys, _ := cmd.Flags().GetStringArray("accept-key")
	ignores, _ := cmd.Flags().GetStringArray("ignore")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewDockleModule(engine.GetClient())

	opts := []modules.DockleOption{}
	if format != "" {
		opts = append(opts, modules.WithDockleFormat(format))
	}
	if exitLevel != "" {
		opts = append(opts, modules.WithExitLevel(exitLevel))
	}
	if len(acceptKeys) > 0 {
		opts = append(opts, modules.WithAcceptKey(acceptKeys))
	}
	if len(ignores) > 0 {
		opts = append(opts, modules.WithDockleIgnore(ignores))
	}
	if output != "" {
		opts = append(opts, modules.WithDockleOutput(output))
	}

	var container *dagger.Container
	if len(args) > 0 {
		// Scan specific image
		container, err = module.ScanImage(ctx, args[0], opts...)
	} else {
		return fmt.Errorf("image reference required for Dockle scan")
	}

	if err != nil {
		return fmt.Errorf("Dockle scan failed: %w", err)
	}

	result, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get output: %w", err)
	}

	return handleOutput(result, output)
}

func runSOPS(cmd *cobra.Command, args []string) error {
	command := args[0]
	output, _ := cmd.Flags().GetString("output")
	kmsARN, _ := cmd.Flags().GetString("kms")
	pgpFingerprint, _ := cmd.Flags().GetString("pgp")
	agePublicKey, _ := cmd.Flags().GetString("age")
	keyFile, _ := cmd.Flags().GetString("key-file")
	inPlace, _ := cmd.Flags().GetBool("in-place")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	ctx := context.Background()
	engine, err := shipdagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	module := modules.NewSOPSModule(engine.GetClient())

	opts := []modules.SOPSOption{}
	if kmsARN != "" {
		opts = append(opts, modules.WithKMSARN(kmsARN))
	}
	if pgpFingerprint != "" {
		opts = append(opts, modules.WithGPGFingerprint(pgpFingerprint))
	}
	if agePublicKey != "" {
		opts = append(opts, modules.WithAgePublicKey(agePublicKey))
	}
	if keyFile != "" {
		opts = append(opts, modules.WithAgeKeyFile(keyFile))
	}
	if inPlace {
		opts = append(opts, modules.WithInPlace(inPlace))
	}
	if output != "" {
		opts = append(opts, modules.WithSOPSOutput(output))
	}

	var container *dagger.Container
	switch command {
	case "encrypt":
		if len(args) < 2 {
			return fmt.Errorf("file path required for encrypt")
		}
		container, err = module.EncryptFile(ctx, args[1], opts...)
	case "decrypt":
		if len(args) < 2 {
			return fmt.Errorf("file path required for decrypt")
		}
		container, err = module.DecryptFile(ctx, args[1], opts...)
	case "rotate":
		if len(args) < 2 {
			return fmt.Errorf("file path required for rotate")
		}
		container, err = module.RotateKeys(ctx, args[1], opts...)
	case "edit":
		if len(args) < 2 {
			return fmt.Errorf("file path required for edit")
		}
		container, err = module.EditFile(ctx, args[1], opts...)
	case "generate-config":
		container, err = module.GenerateConfig(ctx, opts...)
	case "validate":
		if len(args) < 2 {
			return fmt.Errorf("file path required for validate")
		}
		container, err = module.ValidateFile(ctx, args[1], opts...)
	default:
		return fmt.Errorf("unknown command: %s (supported: encrypt, decrypt, rotate, edit, generate-config, validate)", command)
	}

	if err != nil {
		return fmt.Errorf("SOPS operation failed: %w", err)
	}

	result, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get output: %w", err)
	}

	return handleOutput(result, output)
}

func handleOutput(result string, outputFile string) error {
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(result), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		color.Green("Results saved to %s", outputFile)
	} else {
		fmt.Print(result)
	}
	return nil
}
