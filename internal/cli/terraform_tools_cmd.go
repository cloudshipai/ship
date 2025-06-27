package cli

import (
	"fmt"
	"os"
	"path/filepath"

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

	// Add output file flags
	generateDocsCmd.Flags().StringP("output", "o", "", "Output file to save documentation (default: print to stdout)")
	generateDocsCmd.Flags().StringP("filename", "f", "README.md", "Filename to save documentation as")

	costAnalysisCmd.Flags().StringP("output", "o", "", "Output file to save cost analysis (default: print to stdout)")
	costAnalysisCmd.Flags().StringP("format", "", "json", "Output format: json, table")
	costAnalysisCmd.Flags().StringP("region", "r", "us-east-1", "AWS region for pricing (e.g., us-east-1, us-west-2)")

	securityScanCmd.Flags().StringP("output", "o", "", "Output file to save security scan results (default: print to stdout)")
	securityScanCmd.Flags().StringP("format", "", "json", "Output format: json, table, sarif")

	lintCmd.Flags().StringP("output", "o", "", "Output file to save lint results (default: print to stdout)")
	lintCmd.Flags().StringP("format", "", "default", "Output format: default, json, compact")

	checkovScanCmd.Flags().StringP("output", "o", "", "Output file to save scan results (default: print to stdout)")
	checkovScanCmd.Flags().StringP("format", "", "cli", "Output format: cli, json, junit, sarif")

	infracostCmd.Flags().StringP("output", "o", "", "Output file to save cost estimation (default: print to stdout)")
	infracostCmd.Flags().StringP("format", "", "table", "Output format: json, table, html")
}

// saveOrPrintOutput saves output to a file or prints to stdout
func saveOrPrintOutput(output, outputFile string, successMsg string) error {
	if outputFile != "" {
		// Ensure output directory exists
		if dir := filepath.Dir(outputFile); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
		}

		// Write to file
		err := os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		green := color.New(color.FgGreen)
		green.Printf("✓ %s\n", successMsg)
		fmt.Printf("Output saved to: %s\n", outputFile)
	} else {
		// Print to stdout
		green := color.New(color.FgGreen)
		green.Printf("\n✓ %s\n", successMsg)
		fmt.Printf("\n%s\n", output)
	}
	return nil
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

	// Get region flag
	region, _ := cmd.Flags().GetString("region")

	// Check if it's a file or directory
	fileInfo, _ := os.Stat(path)
	var output string

	if fileInfo.IsDir() {
		fmt.Printf("Analyzing Terraform directory: %s\n", path)
		fmt.Printf("Using AWS region: %s\n", region)
		output, err = module.AnalyzeDirectory(ctx, path, region)
	} else {
		fmt.Printf("Analyzing Terraform plan file: %s\n", path)
		fmt.Printf("Using AWS region: %s\n", region)
		output, err = module.AnalyzePlan(ctx, path, region)
	}

	if err != nil {
		return fmt.Errorf("cost analysis failed: %w", err)
	}

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Cost analysis completed!")
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

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Security scan completed!")
}

func runGenerateDocs(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Default to current directory
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Get flags
	outputFile, _ := cmd.Flags().GetString("output")
	filename, _ := cmd.Flags().GetString("filename")

	// If output flag is not set but we have a filename, use it
	if outputFile == "" && filename != "README.md" {
		outputFile = filename
	} else if outputFile == "" && filename == "README.md" {
		// Auto-generate output file path
		outputFile = filepath.Join(dir, "README.md")
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

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Documentation generated!")
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

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Linting completed!")
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

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Checkov scan completed!")
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

	// Get output flag
	outputFile, _ := cmd.Flags().GetString("output")

	// Save or print output
	return saveOrPrintOutput(output, outputFile, "Cost estimation completed!")
}
