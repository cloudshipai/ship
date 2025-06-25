package cli

import (
	"fmt"
	"os"

	"github.com/cloudship/ship/internal/dagger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var terraformToolsCmd = &cobra.Command{
	Use:   "terraform-tools",
	Short: "Run Terraform analysis tools",
	Long:  `Run various Terraform analysis tools including cost estimation, security scanning, and documentation generation`,
}

var costAnalysisCmd = &cobra.Command{
	Use:   "cost-analysis [plan-file]",
	Short: "Analyze Terraform costs using OpenInfraQuote",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCostAnalysis,
}

var securityScanCmd = &cobra.Command{
	Use:   "security-scan [directory]",
	Short: "Scan Terraform code for security issues using InfraScan",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSecurityScan,
}

var generateDocsCmd = &cobra.Command{
	Use:   "generate-docs [directory]",
	Short: "Generate documentation for Terraform modules",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGenerateDocs,
}

var lintCmd = &cobra.Command{
	Use:   "lint [directory]",
	Short: "Lint Terraform code using TFLint",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runLint,
}

var checkovScanCmd = &cobra.Command{
	Use:   "checkov-scan [directory]",
	Short: "Scan Terraform code for security issues using Checkov",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCheckovScan,
}

var infracostCmd = &cobra.Command{
	Use:   "cost-estimate [directory]",
	Short: "Estimate cloud costs using Infracost",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInfracost,
}

func init() {
	rootCmd.AddCommand(terraformToolsCmd)
	terraformToolsCmd.AddCommand(costAnalysisCmd)
	terraformToolsCmd.AddCommand(securityScanCmd)
	terraformToolsCmd.AddCommand(generateDocsCmd)
	terraformToolsCmd.AddCommand(lintCmd)
	terraformToolsCmd.AddCommand(checkovScanCmd)
	terraformToolsCmd.AddCommand(infracostCmd)
}

func runCostAnalysis(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Check if path exists
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create OpenInfraQuote module
	module := engine.NewOpenInfraQuoteModule()

	// Check if it's a file or directory
	fileInfo, _ := os.Stat(path)
	var output string

	if fileInfo.IsDir() {
		fmt.Printf("Analyzing Terraform directory: %s\n", path)
		output, err = module.AnalyzeDirectory(ctx, path)
	} else {
		fmt.Printf("Analyzing Terraform plan file: %s\n", path)
		output, err = module.AnalyzePlan(ctx, path)
	}

	if err != nil {
		return fmt.Errorf("cost analysis failed: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Println("\n✓ Cost analysis completed!")
	fmt.Printf("\nResults:\n%s\n", output)

	return nil
}

func runSecurityScan(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create InfraScan module
	module := engine.NewInfraScanModule()

	fmt.Printf("Scanning Terraform code in: %s\n", dir)
	output, err := module.ScanDirectory(ctx, dir)
	if err != nil {
		return fmt.Errorf("security scan failed: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Println("\n✓ Security scan completed!")
	fmt.Printf("\nResults:\n%s\n", output)

	return nil
}

func runGenerateDocs(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create terraform-docs module
	module := engine.NewTerraformDocsModule()

	fmt.Printf("Generating documentation for: %s\n", dir)
	output, err := module.GenerateMarkdown(ctx, dir)
	if err != nil {
		return fmt.Errorf("documentation generation failed: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Println("\n✓ Documentation generated!")
	fmt.Printf("\n%s\n", output)

	return nil
}

func runLint(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create TFLint module
	module := engine.NewTFLintModule()

	fmt.Printf("Linting Terraform code in: %s\n", dir)
	output, err := module.LintDirectory(ctx, dir)
	if err != nil {
		return fmt.Errorf("linting failed: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Println("\n✓ Linting completed!")
	fmt.Printf("\nResults:\n%s\n", output)

	return nil
}

func runCheckovScan(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create Checkov module
	module := engine.NewCheckovModule()

	fmt.Printf("Scanning Terraform code with Checkov in: %s\n", dir)
	output, err := module.ScanDirectory(ctx, dir)
	if err != nil {
		return fmt.Errorf("Checkov scan failed: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Println("\n✓ Checkov scan completed!")
	fmt.Printf("\nResults:\n%s\n", output)

	return nil
}

func runInfracost(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create Infracost module
	module := engine.NewInfracostModule()

	// Check if INFRACOST_API_KEY is set
	if os.Getenv("INFRACOST_API_KEY") == "" {
		fmt.Println("Warning: INFRACOST_API_KEY is not set. Using free tier limits.")
	}

	fmt.Printf("Estimating costs for Terraform code in: %s\n", dir)
	output, err := module.GenerateTableReport(ctx, dir)
	if err != nil {
		return fmt.Errorf("cost estimation failed: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Println("\n✓ Cost estimation completed!")
	fmt.Printf("\n%s\n", output)

	return nil
}